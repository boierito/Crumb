#import "CrumbLicenseResources.h"
#import "CrumbFileAccess.h"
#import "CrumbMachO.h"
#import <CommonCrypto/CommonDigest.h>
#import <errno.h>
#import <stdio.h>
#import <string.h>
#import <sys/stat.h>
#import <unistd.h>

static NSError *resource_failure(NSString *message) {
  return [NSError errorWithDomain:@"CrumbLicenseResource"
                             code:1
                         userInfo:@{NSLocalizedDescriptionKey : message}];
}
static NSError *resource_file_failure(NSString *operation,
                                      CrumbFileAccess access) {
  NSString *message =
      [NSString stringWithFormat:@"%@: %s (errno %d, mode %04o).", operation,
                                 strerror(access.error_code), access.error_code,
                                 access.mode & 07777];
  return [NSError errorWithDomain:NSPOSIXErrorDomain
                             code:access.error_code
                         userInfo:@{NSLocalizedDescriptionKey : message}];
}
static NSString *
resource_hexhash(unsigned char digest[CC_SHA256_DIGEST_LENGTH]) {
  NSMutableString *text = [NSMutableString string];
  for (int i = 0; i < CC_SHA256_DIGEST_LENGTH; i++)
    [text appendFormat:@"%02x", digest[i]];
  return text;
}
BOOL crumb_copy_license_resource(NSURL *bundle, NSString *relative,
                                 NSDictionary *resource, NSURL *dest,
                                 NSError **error) {
  NSString *file = resource[@"file"], *expectedHash = resource[@"sha256"];
  NSNumber *expectedSize = resource[@"size"];
  BOOL validPath =
      [relative isKindOfClass:NSString.class] && relative.length &&
      !relative.isAbsolutePath && ![relative containsString:@"\\"] &&
      [[relative componentsSeparatedByString:@"/"] containsObject:@"SC_Info"];
  for (NSString *part in [relative componentsSeparatedByString:@"/"])
    if (!part.length || [part isEqualToString:@"."] ||
        [part isEqualToString:@".."])
      validPath = NO;
  if (!validPath || ![file isKindOfClass:NSString.class] ||
      ![file.lastPathComponent isEqualToString:file] ||
      ![file.pathExtension isEqualToString:@"bin"] ||
      ![expectedHash isKindOfClass:NSString.class] ||
      expectedHash.length != 64 ||
      ![expectedSize isKindOfClass:NSNumber.class] ||
      expectedSize.unsignedLongLongValue > 4 * 1024 * 1024) {
    *error = resource_failure(@"Invalid preserved license resource manifest.");
    return NO;
  }
  NSURL *source = [[bundle URLByAppendingPathComponent:@"CrumbLicenseResources"
                                           isDirectory:YES]
      URLByAppendingPathComponent:file];
  CrumbFileAccess access;
  int fd = crumb_open_installed_image(source.fileSystemRepresentation, &access);
  if (fd < 0) {
    *error = resource_file_failure(@"Cannot open preserved license resource",
                                   access);
    return NO;
  }
  struct stat st;
  size_t size = expectedSize.unsignedIntegerValue;
  size_t encodedSize = ((size + 2) / 3) * 4;
  if (fstat(fd, &st) || !S_ISREG(st.st_mode) || st.st_size < 0 ||
      (uint64_t)st.st_size != encodedSize) {
    close(fd);
    *error =
        resource_failure(@"Preserved license resource size/type mismatch.");
    return NO;
  }
  NSMutableData *encoded = [NSMutableData dataWithLength:encodedSize];
  size_t offset = 0;
  while (offset < encodedSize) {
    ssize_t n = read(fd, (uint8_t *)encoded.mutableBytes + offset,
                     MIN((size_t)65536, encodedSize - offset));
    if (n < 0 && errno == EINTR)
      continue;
    if (n <= 0) {
      access.error_code = n < 0 ? errno : EIO;
      close(fd);
      *error = resource_file_failure(@"Cannot read preserved license resource",
                                     access);
      return NO;
    }
    offset += (size_t)n;
  }
  close(fd);
  NSData *data = [[NSData alloc] initWithBase64EncodedData:encoded options:0];
  if (!data || data.length != size) {
    *error = resource_failure(@"Preserved license resource decoding failed.");
    return NO;
  }
  unsigned char digest[CC_SHA256_DIGEST_LENGTH];
  CC_SHA256(data.bytes, (CC_LONG)data.length, digest);
  CrumbEncryption enc = {0};
  if (![resource_hexhash(digest) isEqualToString:expectedHash] ||
      crumb_macho(data.bytes, data.length, &enc) != 0) {
    *error = resource_failure(
        @"Preserved license resource failed integrity/encryption "
        @"verification.");
    return NO;
  }
  NSURL *target = [dest URLByAppendingPathComponent:relative];
  return [NSFileManager.defaultManager
                    createDirectoryAtURL:target.URLByDeletingLastPathComponent
             withIntermediateDirectories:YES
                              attributes:nil
                                   error:error] &&
         [data writeToURL:target options:NSDataWritingAtomic error:error];
}
