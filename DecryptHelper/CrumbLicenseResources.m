#import "CrumbLicenseResources.h"
#import "CrumbFileAccess.h"
#import "CrumbMachO.h"
#import <CommonCrypto/CommonDigest.h>
#import <errno.h>
#import <fcntl.h>
#import <string.h>
#import <sys/mman.h>
#import <sys/stat.h>
#import <unistd.h>

static NSError *resource_failure(NSString *message) {
  return [NSError errorWithDomain:@"CrumbPreservedResource"
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
static BOOL valid_path(NSString *relative, BOOL archive) {
  if (![relative isKindOfClass:NSString.class] || !relative.length ||
      relative.isAbsolutePath || [relative containsString:@"\\"])
    return NO;
  NSArray *parts = [relative componentsSeparatedByString:@"/"];
  for (NSString *part in parts)
    if (!part.length || [part isEqualToString:@"."] ||
        [part isEqualToString:@".."])
      return NO;
  if (archive)
    return
        [parts.firstObject caseInsensitiveCompare:@"Payload"] != NSOrderedSame;
  return [parts containsObject:@"SC_Info"];
}
BOOL crumb_create_archive_directory(NSURL *dest, NSString *relative,
                                    NSError **error) {
  if (!valid_path(relative, YES)) {
    *error = resource_failure(@"Invalid preserved archive directory.");
    return NO;
  }
  return [NSFileManager.defaultManager
             createDirectoryAtURL:[dest URLByAppendingPathComponent:relative
                                                        isDirectory:YES]
      withIntermediateDirectories:YES
                       attributes:nil
                            error:error];
}
static BOOL copy_resource(NSURL *bundle, NSString *relative,
                          NSDictionary *resource, NSURL *dest, BOOL archive,
                          NSError **error) {
  NSString *file = resource[@"file"], *expectedHash = resource[@"sha256"];
  NSNumber *expectedSize = resource[@"size"];
  uint64_t limit = archive ? 512ULL * 1024 * 1024 : 4ULL * 1024 * 1024;
  if (!valid_path(relative, archive) || ![file isKindOfClass:NSString.class] ||
      ![file.lastPathComponent isEqualToString:file] ||
      ![file.pathExtension isEqualToString:@"bin"] ||
      ![expectedHash isKindOfClass:NSString.class] ||
      expectedHash.length != 64 ||
      ![expectedSize isKindOfClass:NSNumber.class] ||
      expectedSize.unsignedLongLongValue > limit) {
    *error = resource_failure(@"Invalid preserved resource manifest.");
    return NO;
  }
  NSString *folder =
      archive ? @"CrumbArchiveResources" : @"CrumbLicenseResources";
  NSURL *source = [[bundle URLByAppendingPathComponent:folder isDirectory:YES]
      URLByAppendingPathComponent:file];
  CrumbFileAccess access;
  int input =
      crumb_open_installed_image(source.fileSystemRepresentation, &access);
  if (input < 0) {
    *error = resource_file_failure(@"Cannot open preserved resource", access);
    return NO;
  }
  size_t size = expectedSize.unsignedIntegerValue,
         encodedSize = ((size + 2) / 3) * 4;
  struct stat st;
  if (fstat(input, &st) || !S_ISREG(st.st_mode) || st.st_size < 0 ||
      (uint64_t)st.st_size != encodedSize) {
    close(input);
    *error = resource_failure(@"Preserved resource size/type mismatch.");
    return NO;
  }
  NSURL *temporary = [[NSURL fileURLWithPath:NSTemporaryDirectory()]
      URLByAppendingPathComponent:NSUUID.UUID.UUIDString];
  int output = open(temporary.fileSystemRepresentation,
                    O_RDWR | O_CREAT | O_EXCL | O_NOFOLLOW | O_CLOEXEC, 0600);
  if (output < 0) {
    access.error_code = errno;
    close(input);
    *error = resource_file_failure(@"Cannot create temporary resource", access);
    return NO;
  }
  BOOL ok = YES;
  size_t remaining = encodedSize, decoded = 0;
  CC_SHA256_CTX digest;
  CC_SHA256_Init(&digest);
  while (ok && remaining) {
    @autoreleasepool {
      size_t wanted = MIN((size_t)65536, remaining), got = 0;
      NSMutableData *encoded = [NSMutableData dataWithLength:wanted];
      while (got < wanted) {
        ssize_t n =
            read(input, (uint8_t *)encoded.mutableBytes + got, wanted - got);
        if (n < 0 && errno == EINTR)
          continue;
        if (n <= 0) {
          ok = NO;
          *error = resource_failure(@"Preserved resource read failed.");
          break;
        }
        got += (size_t)n;
      }
      if (ok) {
        NSData *data = [[NSData alloc] initWithBase64EncodedData:encoded
                                                         options:0];
        if (!data || data.length > size - decoded) {
          ok = NO;
          *error = resource_failure(@"Preserved resource decoding failed.");
        } else {
          CC_SHA256_Update(&digest, data.bytes, (CC_LONG)data.length);
          size_t written = 0;
          while (written < data.length) {
            ssize_t n = write(output, (const uint8_t *)data.bytes + written,
                              data.length - written);
            if (n < 0 && errno == EINTR)
              continue;
            if (n <= 0) {
              ok = NO;
              *error = resource_failure(@"Preserved resource write failed.");
              break;
            }
            written += (size_t)n;
          }
          decoded += data.length;
          remaining -= wanted;
        }
      }
    }
  }
  close(input);
  unsigned char hash[CC_SHA256_DIGEST_LENGTH];
  CC_SHA256_Final(hash, &digest);
  if (ok && (decoded != size ||
             ![resource_hexhash(hash) isEqualToString:expectedHash])) {
    ok = NO;
    *error = resource_failure(@"Preserved resource hash/length mismatch.");
  }
  if (ok && size) {
    void *map = mmap(NULL, size, PROT_READ, MAP_PRIVATE, output, 0);
    CrumbEncryption enc = {0};
    if (map == MAP_FAILED || crumb_macho(map, size, &enc) != 0) {
      ok = NO;
      *error = resource_failure(
          @"Preserved resource failed encryption verification.");
    }
    if (map != MAP_FAILED)
      munmap(map, size);
  }
  if (ok && fsync(output)) {
    ok = NO;
    *error = resource_failure(@"Cannot flush preserved resource.");
  }
  close(output);
  if (ok) {
    NSURL *target = [dest URLByAppendingPathComponent:relative];
    ok = [NSFileManager.defaultManager
                    createDirectoryAtURL:target.URLByDeletingLastPathComponent
             withIntermediateDirectories:YES
                              attributes:nil
                                   error:error] &&
         [NSFileManager.defaultManager moveItemAtURL:temporary
                                               toURL:target
                                               error:error];
  }
  if (!ok)
    [NSFileManager.defaultManager removeItemAtURL:temporary error:nil];
  return ok;
}
BOOL crumb_copy_license_resource(NSURL *bundle, NSString *relative,
                                 NSDictionary *resource, NSURL *dest,
                                 NSError **error) {
  return copy_resource(bundle, relative, resource, dest, NO, error);
}
BOOL crumb_copy_archive_resource(NSURL *bundle, NSString *relative,
                                 NSDictionary *resource, NSURL *dest,
                                 NSError **error) {
  return copy_resource(bundle, relative, resource, dest, YES, error);
}
