#import "CrumbMachO.h"
#import <CommonCrypto/CommonDigest.h>
#import <UIKit/UIKit.h>
#import <dlfcn.h>
#import <errno.h>
#import <fcntl.h>
#import <mach/mach.h>
#import <math.h>
#import <stdio.h>
#import <sys/mman.h>
#import <sys/stat.h>
#import <unistd.h>
#import <zlib.h>

static int remap_image(void *address, size_t length, uint32_t cpu,
                       uint32_t subtype) {
  typedef int (*Remap)(void *, size_t, uint32_t, uint32_t, uint32_t);
  Remap symbol = (Remap)dlsym(RTLD_DEFAULT, "mremap_encrypted");
  if (symbol)
    return symbol(address, length, 2, cpu, subtype);
#if defined(__arm64__)
  register uint64_t arg0 __asm__("x0") = (uint64_t)address;
  register uint64_t arg1 __asm__("x1") = length;
  register uint64_t arg2 __asm__("x2") = 2;
  register uint64_t arg3 __asm__("x3") = cpu;
  register uint64_t arg4 __asm__("x4") = subtype;
  register uint64_t number __asm__("x16") = 489;
  uint64_t failed;
  __asm__ volatile("svc #0x80\n\tcset %[failed], cs"
                   : "+r"(arg0), "+r"(arg1), [failed] "=r"(failed)
                   : "r"(arg2), "r"(arg3), "r"(arg4), "r"(number)
                   : "memory", "cc");
  if (failed) {
    errno = (int)arg0;
    return -1;
  }
  return (int)arg0;
#else
  return -1;
#endif
}
static NSError *failure(NSString *s) {
  return [NSError errorWithDomain:@"CrumbDecrypt"
                             code:1
                         userInfo:@{NSLocalizedDescriptionKey : s}];
}
static uint32_t codewords(const uint8_t *b, size_t n) {
  uint32_t count = 0, w;
  for (size_t p = 0; p + 4 <= n; p += 4) {
    memcpy(&w, b + p, 4);
    if (w == 0xd65f03c0 || w == 0xd503237f || w == 0xd65f0fff)
      count++;
  }
  return count;
}
static double histogram_entropy(uint64_t h[256], size_t n) {
  double e = 0;
  for (int i = 0; i < 256; i++)
    if (h[i]) {
      double p = (double)h[i] / n;
      e -= p * log2(p);
    }
  return e;
}
// Register the original Apple signature when present. The source stays
// read-only; this does not sign or alter the target. FairPlay still decides
// license access.
static BOOL register_signature(int fd, const uint8_t *header, size_t size) {
  uint32_t count = 0;
  memcpy(&count, header + 16, 4);
  size_t pos = 32;
  for (uint32_t i = 0; i < count; i++) {
    uint32_t command = 0, length = 0;
    memcpy(&command, header + pos, 4);
    memcpy(&length, header + pos + 4, 4);
    if (command == 0x1d && length >= 16) {
      uint32_t offset = 0, bytes = 0;
      memcpy(&offset, header + pos + 8, 4);
      memcpy(&bytes, header + pos + 12, 4);
      if (offset > size || bytes > size - offset || !bytes)
        return NO;
      fsignatures_t signatures = {.fs_file_start = 0,
                                  .fs_blob_start = (void *)(uintptr_t)offset,
                                  .fs_blob_size = bytes};
      return fcntl(fd, F_ADDFILESIGS_RETURN, &signatures) == 0;
    }
    pos += length;
  }
  return NO;
}
static BOOL decrypt(NSURL *source, NSURL *dest, BOOL *encrypted,
                    NSError **error) {
  int fd = open(source.fileSystemRepresentation, O_RDONLY | O_NOFOLLOW);
  if (fd < 0) {
    *error = failure(@"Cannot open installed image.");
    return NO;
  }
  struct stat st;
  if (fstat(fd, &st) || st.st_size < 0) {
    close(fd);
    *error = failure(@"Cannot inspect installed image.");
    return NO;
  }
  size_t size = (size_t)st.st_size;
  if (!size) {
    close(fd);
    return [NSFileManager.defaultManager copyItemAtURL:source
                                                 toURL:dest
                                                 error:error];
  }
  uint8_t *map = mmap(NULL, size, PROT_READ, MAP_PRIVATE, fd, 0);
  if (map == MAP_FAILED) {
    close(fd);
    *error = failure(@"Cannot map installed image.");
    return NO;
  }
  CrumbEncryption enc = {0};
  int kind = crumb_macho(map, size, &enc);
  *encrypted = kind == 1;
  if (kind < 0) {
    munmap(map, size);
    close(fd);
    *error = failure(@"Malformed or unsupported Mach-O (thin arm64 required).");
    return NO;
  }
  if (!kind) {
    munmap(map, size);
    close(fd);
    return [NSFileManager.defaultManager copyItemAtURL:source
                                                 toURL:dest
                                                 error:error];
  }
  size_t pg = (size_t)getpagesize();
  if (enc.offset % pg) {
    munmap(map, size);
    close(fd);
    *error = failure(@"Encrypted range is not page aligned.");
    return NO;
  }
  if (![NSFileManager.defaultManager copyItemAtURL:source
                                             toURL:dest
                                             error:error]) {
    munmap(map, size);
    close(fd);
    return NO;
  }
  // Installed resources may be read-only. Only the data-container copy is made
  // writable.
  if (![NSFileManager.defaultManager setAttributes:@{
        NSFilePosixPermissions : @(0700)
      }
                                      ofItemAtPath:dest.path
                                             error:error]) {
    munmap(map, size);
    close(fd);
    return NO;
  }
  int out = open(dest.fileSystemRepresentation, O_WRONLY | O_NOFOLLOW);
  if (out < 0) {
    munmap(map, size);
    close(fd);
    *error = failure(@"Cannot create decrypted image.");
    return NO;
  }
  if (register_signature(fd, map, size)) {
    munmap(map, size);
    map = mmap(NULL, size, PROT_READ, MAP_PRIVATE, fd, 0);
  }
  errno = 0;
  int remap_status =
      map == MAP_FAILED
          ? -1
          : remap_image(map + enc.offset, enc.length, enc.cpu, enc.subtype);
  int remap_errno = errno;
  int page_status = 0;
  BOOL ok = map != MAP_FAILED && remap_status == 0;
  uint8_t *plain = malloc(pg), *cipher = malloc(pg);
  uint64_t before[256] = {0}, after[256] = {0};
  uint64_t oldwords = 0, newwords = 0;
  BOOL changed = NO;
  if (!plain || !cipher)
    ok = NO;
  for (size_t p = 0; ok && p < enc.length; p += pg) {
    size_t len = MIN(pg, enc.length - p);
    vm_size_t got = 0;
    kern_return_t kr = vm_read_overwrite(mach_task_self(),
                                         (vm_address_t)(map + enc.offset + p),
                                         len, (vm_address_t)plain, &got);
    if (kr != KERN_SUCCESS || got != len ||
        pread(fd, cipher, len, enc.offset + p) != (ssize_t)len) {
      page_status = (int)kr;
      ok = NO;
      break;
    }
    changed |= memcmp(cipher, plain, len) != 0;
    oldwords += codewords(cipher, len);
    newwords += codewords(plain, len);
    for (size_t i = 0; i < len; i++) {
      before[cipher[i]]++;
      after[plain[i]]++;
    }
    if (pwrite(out, plain, len, enc.offset + p) != (ssize_t)len) {
      ok = NO;
      break;
    }
  }
  free(plain);
  free(cipher);
  if (map != MAP_FAILED)
    munmap(map, size);
  close(fd);
  // Conservative heuristics supplement complete fault-free FairPlay page reads.
  ok = ok && changed && histogram_entropy(before, enc.length) > 7.9 &&
       histogram_entropy(after, enc.length) < 7.5 && newwords > oldwords;
  if (ok) {
    uint32_t zero = 0;
    ok = pwrite(out, &zero, 4, enc.command) == 4;
  }
  if (ok)
    ok = fsync(out) == 0;
  close(out);
  if (!ok) {
    [NSFileManager.defaultManager removeItemAtURL:dest error:nil];
    *error = failure([NSString
        stringWithFormat:@"FairPlay verification failed (remap %d, errno %d, "
                         @"page %d). Check the licensed device account.",
                         remap_status, remap_errno, page_status]);
  }
  return ok;
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
static BOOL pack(NSURL *root, NSURL *out, NSMutableArray *hashes,
                 NSError **error) {
  NSFileManager *fm = NSFileManager.defaultManager;
  FILE *file = fopen(out.fileSystemRepresentation, "wb");
  if (!file) {
    *error = failure(@"Cannot create IPA.");
    return NO;
  }
  NSMutableData *central = [NSMutableData data];
  uint64_t offset = 0;
  uint16_t count = 0;
  BOOL ok = YES;
  NSArray *keys =
      @[ NSURLIsRegularFileKey, NSURLIsSymbolicLinkKey, NSURLFileSizeKey ];
  __block NSError *enumerationError = nil;
  uint8_t *buffer = malloc(1024 * 1024);
  if (!buffer) {
    fclose(file);
    *error = failure(@"Cannot allocate archive buffer.");
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
        *error = failure(@"Output contains a link.");
        ok = NO;
        break;
      }
      if (![attrs[NSURLIsRegularFileKey] boolValue])
        continue;
      uint64_t size = [attrs[NSURLFileSizeKey] unsignedLongLongValue];
      NSString *name = [url.path substringFromIndex:root.path.length + 1];
      NSData *n = [name dataUsingEncoding:NSUTF8StringEncoding];
      if (size > 512 * 1024 * 1024 || n.length > 65535 || count == 65535 ||
          offset + size + n.length + 30 > 4000000000ULL) {
        *error = failure(@"Result exceeds archive budget.");
        ok = NO;
        break;
      }
      FILE *input = fopen(url.fileSystemRepresentation, "rb");
      if (!input) {
        ok = NO;
        break;
      }
      CC_SHA256_CTX digest;
      CC_SHA256_Init(&digest);
      uint32_t crc = 0;
      uint64_t read = 0;
      size_t len;
      while ((len = fread(buffer, 1, 1024 * 1024, input))) {
        read += len;
        CC_SHA256_Update(&digest, buffer, (CC_LONG)len);
        crc = (uint32_t)crc32(crc, buffer, (uInt)len);
      }
      unsigned char hash[CC_SHA256_DIGEST_LENGTH];
      CC_SHA256_Final(hash, &digest);
      if (ferror(input) || read != size) {
        fclose(input);
        ok = NO;
        break;
      }
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
        fclose(input);
        ok = NO;
        break;
      }
      read = 0;
      while ((len = fread(buffer, 1, 1024 * 1024, input))) {
        read += len;
        if (fwrite(buffer, 1, len, file) != len) {
          ok = NO;
          break;
        }
      }
      if (ferror(input) || read != size)
        ok = NO;
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
      put32(central, 0100755U << 16);
      put32(central, (uint32_t)offset);
      [central appendData:n];
      offset += h.length + size;
      count++;
      [hashes addObject:@{@"path" : name, @"sha256" : hexhash(hash)}];
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
      *error = failure(@"IPA archive write failed.");
  }
  return ok;
}
@interface CrumbHelper : UIViewController
@property(nonatomic, strong) UILabel *label;
@end
@implementation CrumbHelper
- (void)viewDidLoad {
  [super viewDidLoad];
  self.view.backgroundColor = UIColor.systemBackgroundColor;
  self.label = [[UILabel alloc] initWithFrame:CGRectZero];
  self.label.translatesAutoresizingMaskIntoConstraints = NO;
  self.label.numberOfLines = 0;
  self.label.text = @"Decrypting licensed app… Keep this screen open.";
  [self.view addSubview:self.label];
  [NSLayoutConstraint activateConstraints:@[
    [self.label.leadingAnchor constraintEqualToAnchor:self.view.leadingAnchor
                                             constant:24],
    [self.label.trailingAnchor constraintEqualToAnchor:self.view.trailingAnchor
                                              constant:-24],
    [self.label.centerYAnchor constraintEqualToAnchor:self.view.centerYAnchor]
  ]];
  dispatch_async(dispatch_get_global_queue(QOS_CLASS_USER_INITIATED, 0), ^{
    @autoreleasepool {
      [self run];
    }
  });
}
- (void)run {
  NSFileManager *fm = NSFileManager.defaultManager;
  NSURL *docs = [fm URLsForDirectory:NSDocumentDirectory
                           inDomains:NSUserDomainMask]
                    .firstObject;
  [fm removeItemAtURL:[docs URLByAppendingPathComponent:@"result.json"]
                error:nil];
  [fm removeItemAtURL:[docs URLByAppendingPathComponent:@"result.ipa"]
                error:nil];
  NSURL *root = [docs URLByAppendingPathComponent:@"working" isDirectory:YES];
  [fm removeItemAtURL:root error:nil];
  NSError *error = nil;
  NSDictionary *job = [NSJSONSerialization
      JSONObjectWithData:[NSData
                             dataWithContentsOfURL:
                                 [NSBundle.mainBundle.bundleURL
                                     URLByAppendingPathComponent:@"job.json"]]
                 options:0
                   error:&error];
  NSString *app = job[@"app_name"];
  BOOL ok = job && [app isKindOfClass:NSString.class] &&
            [app.lastPathComponent isEqualToString:app] &&
            [app hasSuffix:@".app"];
  NSURL *src = [NSBundle.mainBundle.bundleURL
      URLByAppendingPathComponent:@"EmbeddedPayload/Target"
                      isDirectory:YES];
  NSURL *dest = [[root URLByAppendingPathComponent:@"Payload"]
      URLByAppendingPathComponent:app ?: @"Target.app"];
  __block NSError *enumerationError = nil;
  NSUInteger decrypted = 0;
  NSMutableArray *hashes = [NSMutableArray array];
  uint64_t total = 0;
  if (ok)
    ok = [fm createDirectoryAtURL:dest
        withIntermediateDirectories:YES
                         attributes:nil
                              error:&error];
  NSArray *keys = @[
    NSURLIsRegularFileKey, NSURLIsDirectoryKey, NSURLIsSymbolicLinkKey,
    NSURLFileSizeKey
  ];
  if (ok)
    for (NSURL *url in [fm enumeratorAtURL:src
                includingPropertiesForKeys:keys
                                   options:0
                              errorHandler:^BOOL(NSURL *u, NSError *e) {
                                enumerationError = e;
                                return NO;
                              }]) {
      @autoreleasepool {
        NSDictionary *a = [url resourceValuesForKeys:keys error:&error];
        NSString *relative = [url.path substringFromIndex:src.path.length + 1];
        NSURL *target = [dest URLByAppendingPathComponent:relative];
        if (!a || [a[NSURLIsSymbolicLinkKey] boolValue]) {
          ok = NO;
          break;
        }
        if ([a[NSURLIsDirectoryKey] boolValue]) {
          if (![fm createDirectoryAtURL:target
                  withIntermediateDirectories:YES
                                   attributes:nil
                                        error:&error]) {
            ok = NO;
            break;
          }
          continue;
        }
        if (![a[NSURLIsRegularFileKey] boolValue]) {
          ok = NO;
          break;
        }
        uint64_t size = [a[NSURLFileSizeKey] unsignedLongLongValue];
        total += size;
        if (size > 512 * 1024 * 1024 || total > 4000000000ULL) {
          error = failure(@"Payload exceeds the device memory/storage budget.");
          ok = NO;
          break;
        }
        [fm createDirectoryAtURL:target.URLByDeletingLastPathComponent
            withIntermediateDirectories:YES
                             attributes:nil
                                  error:&error];
        BOOL encrypted = NO;
        if (!decrypt(url, target, &encrypted, &error)) {
          ok = NO;
          break;
        }
        if (encrypted)
          decrypted++;
      }
    }
  if (enumerationError) {
    ok = NO;
    error = enumerationError;
  }
  if (ok && decrypted == 0) {
    error = failure(
        @"No encrypted images were found; this is not evidence of decryption.");
    ok = NO;
  }
  if (ok)
    ok = pack(root, [docs URLByAppendingPathComponent:@"result.ipa"], hashes,
              &error);
  NSDictionary *result = @{
    @"job_id" : job[@"job_id"] ?: @"",
    @"complete" : @(ok),
    @"decrypted" : @(decrypted),
    @"files" : hashes,
    @"message" : error.localizedDescription ?: @""
  };
  NSData *json = [NSJSONSerialization dataWithJSONObject:result
                                                 options:0
                                                   error:nil];
  [json writeToURL:[docs URLByAppendingPathComponent:@"result.json"]
           options:NSDataWritingAtomic
             error:nil];
  [fm removeItemAtURL:root error:nil];
  dispatch_async(dispatch_get_main_queue(), ^{
    self.label.text=ok?@"Decryption finished. Return to Crumb and tap Retrieve result.":error.localizedDescription?:@"Decryption failed. No completed IPA was produced.";
    if (ok) {
      NSURLComponents *callback = [NSURLComponents new];
      callback.scheme = @"crumb";
      callback.host = @"device-result";
      callback.queryItems =
          @[ [NSURLQueryItem queryItemWithName:@"job" value:job[@"job_id"]] ];
      [UIApplication.sharedApplication openURL:callback.URL
                                       options:@{}
                             completionHandler:nil];
    }
  });
}
@end
@interface Delegate : UIResponder <UIApplicationDelegate>
@property(nonatomic, strong) UIWindow *window;
@end
@implementation Delegate
- (BOOL)application:(UIApplication *)app
    didFinishLaunchingWithOptions:(NSDictionary *)options {
  self.window = [[UIWindow alloc] initWithFrame:UIScreen.mainScreen.bounds];
  self.window.rootViewController = [CrumbHelper new];
  [self.window makeKeyAndVisible];
  return YES;
}
- (BOOL)application:(UIApplication *)app
            openURL:(NSURL *)url
            options:(NSDictionary *)options {
  return YES;
}
@end
int main(int argc, char **argv) {
  @autoreleasepool {
    return UIApplicationMain(argc, argv, nil,
                             NSStringFromClass(Delegate.class));
  }
}
