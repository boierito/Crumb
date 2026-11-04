#import "CrumbIPAArchive.h"
#import <CommonCrypto/CommonDigest.h>
#import <stdio.h>
#import <stdlib.h>
#import <zlib.h>
static NSError *archive_failure(NSString *message) {
  return [NSError errorWithDomain:@"CrumbIPAArchive"
                             code:1
                         userInfo:@{NSLocalizedDescriptionKey : message}];
}
static NSString *hexhash(unsigned char hash[CC_SHA256_DIGEST_LENGTH]) {
  NSMutableString *s = [NSMutableString string];
  for (int i = 0; i < CC_SHA256_DIGEST_LENGTH; i++)
    [s appendFormat:@"%02x", hash[i]];
  return s;
}
static void put16(NSMutableData *d, uint16_t n) {
  uint16_t x = CFSwapInt16HostToLittle(n);
  [d appendBytes:&x length:2];
}
static void put32(NSMutableData *d, uint32_t n) {
  uint32_t x = CFSwapInt32HostToLittle(n);
  [d appendBytes:&x length:4];
}
// Stored ZIP entries, CRC32 checked by the receiving Rust verifier. No ZIP64.
BOOL crumb_pack_ipa(NSURL *root, NSURL *out, NSMutableArray *hashes,
                    NSError **error) {
  NSFileManager *fm = NSFileManager.defaultManager;
  FILE *file = fopen(out.fileSystemRepresentation, "wb");
  if (!file) {
    *error = archive_failure(@"Cannot create IPA.");
    return NO;
  }
  NSMutableData *central = [NSMutableData data];
  uint64_t offset = 0;
  uint16_t count = 0;
  BOOL ok = YES;
  NSArray *keys = @[
    NSURLIsRegularFileKey, NSURLIsDirectoryKey, NSURLIsSymbolicLinkKey,
    NSURLFileSizeKey
  ];
  __block NSError *enumerationError = nil;
  uint8_t *buffer = malloc(1024 * 1024);
  if (!buffer) {
    fclose(file);
    *error = archive_failure(@"Cannot allocate archive buffer.");
    return NO;
  }
  for (NSURL *url in [fm enumeratorAtURL:root
              includingPropertiesForKeys:keys
                                 options:0
                            errorHandler:^BOOL(NSURL *u, NSError *e) {
                              enumerationError = e;
                              return NO;
                            }]) {
    @autoreleasepool {
      NSDictionary *attrs = [url resourceValuesForKeys:keys error:error];
      if (!attrs) {
        ok = NO;
        break;
      }
      if ([attrs[NSURLIsSymbolicLinkKey] boolValue]) {
        *error = archive_failure(@"Output contains a link.");
        ok = NO;
        break;
      }
      BOOL directory = [attrs[NSURLIsDirectoryKey] boolValue];
      if (![attrs[NSURLIsRegularFileKey] boolValue] && !directory)
        continue;
      NSString *name = [url.path substringFromIndex:root.path.length + 1];
      if (directory) {
        if ([name.pathComponents.firstObject isEqualToString:@"Payload"])
          continue;
        name = [name stringByAppendingString:@"/"];
      }
      uint64_t size =
          directory ? 0 : [attrs[NSURLFileSizeKey] unsignedLongLongValue];
      NSData *n = [name dataUsingEncoding:NSUTF8StringEncoding];
      if (size > 512 * 1024 * 1024 || n.length > 65535 || count == 65535 ||
          offset + size + n.length + 30 > 4000000000ULL) {
        *error = archive_failure(@"Result exceeds archive budget.");
        ok = NO;
        break;
      }
      FILE *input =
          directory ? NULL : fopen(url.fileSystemRepresentation, "rb");
      if (!directory && !input) {
        ok = NO;
        break;
      }
      CC_SHA256_CTX digest;
      CC_SHA256_Init(&digest);
      uint32_t crc = 0;
      uint64_t read = 0;
      size_t len;
      while (input && (len = fread(buffer, 1, 1024 * 1024, input))) {
        read += len;
        CC_SHA256_Update(&digest, buffer, (CC_LONG)len);
        crc = (uint32_t)crc32(crc, buffer, (uInt)len);
      }
      unsigned char hash[CC_SHA256_DIGEST_LENGTH];
      CC_SHA256_Final(hash, &digest);
      if (input && (ferror(input) || read != size)) {
        if (input)
          fclose(input);
        ok = NO;
        break;
      }
      if (input)
        rewind(input);
      NSMutableData *h = [NSMutableData data];
      put32(h, 0x04034b50);
      put16(h, 20);
      put16(h, 0x800);
      put16(h, 0);
      put16(h, 0);
      put16(h, 0);
      put32(h, crc);
      put32(h, (uint32_t)size);
      put32(h, (uint32_t)size);
      put16(h, (uint16_t)n.length);
      put16(h, 0);
      [h appendData:n];
      if (fwrite(h.bytes, 1, h.length, file) != h.length) {
        if (input)
          fclose(input);
        ok = NO;
        break;
      }
      read = 0;
      while (input && (len = fread(buffer, 1, 1024 * 1024, input))) {
        read += len;
        if (fwrite(buffer, 1, len, file) != len) {
          ok = NO;
          break;
        }
      }
      if (input && (ferror(input) || read != size))
        ok = NO;
      if (input)
        fclose(input);
      if (!ok)
        break;
      put32(central, 0x02014b50);
      put16(central, 0x0314);
      put16(central, 20);
      put16(central, 0x800);
      put16(central, 0);
      put16(central, 0);
      put16(central, 0);
      put32(central, crc);
      put32(central, (uint32_t)size);
      put32(central, (uint32_t)size);
      put16(central, (uint16_t)n.length);
      put16(central, 0);
      put16(central, 0);
      put16(central, 0);
      put16(central, 0);
      put32(central, ((directory ? 0040755U : 0100755U) << 16) |
                         (directory ? 0x10U : 0));
      put32(central, (uint32_t)offset);
      [central appendData:n];
      offset += h.length + size;
      count++;
      [hashes addObject:@{
        @"path" : name,
        @"sha256" : hexhash(hash),
        @"directory" : @(directory)
      }];
    }
  }
  if (enumerationError) {
    *error = enumerationError;
    ok = NO;
  }
  if (ok) {
    ok = fwrite(central.bytes, 1, central.length, file) == central.length;
    NSMutableData *end = [NSMutableData data];
    put32(end, 0x06054b50);
    put16(end, 0);
    put16(end, 0);
    put16(end, count);
    put16(end, count);
    put32(end, (uint32_t)central.length);
    put32(end, (uint32_t)offset);
    put16(end, 0);
    ok = ok && fwrite(end.bytes, 1, end.length, file) == end.length;
  }
  free(buffer);
  if (fclose(file) != 0)
    ok = NO;
  if (!ok) {
    [fm removeItemAtURL:out error:nil];
    if (!*error)
      *error = archive_failure(@"IPA archive write failed.");
  }
  return ok;
}
