#import "CrumbFileAccess.h"
#import "CrumbIPAArchive.h"
#import "CrumbLicenseResources.h"
#import "CrumbMachO.h"
#import <CommonCrypto/CommonDigest.h>
#import <UIKit/UIKit.h>
#import <dlfcn.h>
#import <errno.h>
#import <fcntl.h>
#import <mach/mach.h>
#import <math.h>
#import <stdio.h>
#import <string.h>
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
static NSError *file_failure(NSString *operation, CrumbFileAccess status) {
  NSString *details =
      status.metadata_available
          ? [NSString stringWithFormat:@"; mode %04o, uid %u, gid %u",
                                       status.mode & 07777, status.owner,
                                       status.group]
          : @"; metadata unavailable";
  NSString *message = [NSString
      stringWithFormat:@"%@: %s (errno %d%@).", operation,
                       strerror(status.error_code), status.error_code, details];
  return [NSError errorWithDomain:NSPOSIXErrorDomain
                             code:status.error_code
                         userInfo:@{NSLocalizedDescriptionKey : message}];
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
  CrumbFileAccess access;
  int fd = crumb_open_installed_image(source.fileSystemRepresentation, &access);
  if (fd < 0) {
    *error = file_failure(@"Cannot open installed image", access);
    return NO;
  }
  struct stat st;
  int inspected = fstat(fd, &st);
  if (inspected || st.st_size < 0) {
    access.error_code = inspected ? errno : EFBIG;
    close(fd);
    *error = file_failure(@"Cannot inspect installed image", access);
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
    access.error_code = errno;
    close(fd);
    *error = file_failure(@"Cannot map installed image", access);
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
@interface CrumbHelper : UIViewController
@property(nonatomic, strong) UILabel *label;
@property(nonatomic, strong) UIProgressView *progress;
@property(nonatomic, strong) UILabel *titleLabel;
- (void)status:(NSString *)text progress:(float)progress;

@end
@implementation CrumbHelper
- (void)viewDidLoad {
  [super viewDidLoad];
  UIApplication.sharedApplication.idleTimerDisabled = YES;
  self.view.backgroundColor = UIColor.systemBackgroundColor;
  self.view.tintColor = [UIColor colorWithRed:0 green:0.66 blue:0.88 alpha:1];
  UIImageView *icon = [[UIImageView alloc] initWithImage:[UIImage imageNamed:@"Icon60@3x.png"]];
  icon.contentMode = UIViewContentModeScaleAspectFit;
  icon.layer.cornerRadius = 22;
  icon.clipsToBounds = YES;
  icon.translatesAutoresizingMaskIntoConstraints = NO;
  UIView *brand = [UIView new];
  [brand addSubview:icon];
  [NSLayoutConstraint activateConstraints:@[
    [brand.heightAnchor constraintEqualToConstant:88],
    [icon.heightAnchor constraintEqualToConstant:88],
    [icon.widthAnchor constraintEqualToConstant:88],
    [icon.centerXAnchor constraintEqualToAnchor:brand.centerXAnchor],
    [icon.topAnchor constraintEqualToAnchor:brand.topAnchor]
  ]];
  self.titleLabel = [UILabel new];
  self.titleLabel.font = [UIFont preferredFontForTextStyle:UIFontTextStyleLargeTitle];
  self.titleLabel.adjustsFontForContentSizeCategory = YES;
  self.titleLabel.text = @"Crumb Decrypt";
  self.titleLabel.numberOfLines = 0;
  self.titleLabel.textAlignment = NSTextAlignmentCenter;
  self.label = [UILabel new];
  self.label.font = [UIFont preferredFontForTextStyle:UIFontTextStyleBody];
  self.label.adjustsFontForContentSizeCategory = YES;
  self.label.numberOfLines = 0;
  self.label.textAlignment = NSTextAlignmentCenter;
  self.label.text = @"Preparing app…";
  self.progress = [[UIProgressView alloc] initWithProgressViewStyle:UIProgressViewStyleDefault];
  self.progress.accessibilityLabel = @"Decryption progress";
  UIStackView *stack = [[UIStackView alloc] initWithArrangedSubviews:@[brand,self.titleLabel,self.label,self.progress]];
  stack.axis = UILayoutConstraintAxisVertical;
  stack.spacing = 24;
  stack.translatesAutoresizingMaskIntoConstraints = NO;
  UIScrollView *scroll = [UIScrollView new];
  scroll.translatesAutoresizingMaskIntoConstraints = NO;
  UIView *content = [UIView new];
  content.translatesAutoresizingMaskIntoConstraints = NO;
  [self.view addSubview:scroll];
  [scroll addSubview:content];
  [content addSubview:stack];
  [NSLayoutConstraint activateConstraints:@[
    [scroll.topAnchor constraintEqualToAnchor:self.view.safeAreaLayoutGuide.topAnchor],
    [scroll.bottomAnchor constraintEqualToAnchor:self.view.safeAreaLayoutGuide.bottomAnchor],
    [scroll.leadingAnchor constraintEqualToAnchor:self.view.safeAreaLayoutGuide.leadingAnchor],
    [scroll.trailingAnchor constraintEqualToAnchor:self.view.safeAreaLayoutGuide.trailingAnchor],
    [content.topAnchor constraintEqualToAnchor:scroll.contentLayoutGuide.topAnchor],
    [content.bottomAnchor constraintEqualToAnchor:scroll.contentLayoutGuide.bottomAnchor],
    [content.leadingAnchor constraintEqualToAnchor:scroll.contentLayoutGuide.leadingAnchor],
    [content.trailingAnchor constraintEqualToAnchor:scroll.contentLayoutGuide.trailingAnchor],
    [content.widthAnchor constraintEqualToAnchor:scroll.frameLayoutGuide.widthAnchor],
    [content.heightAnchor constraintGreaterThanOrEqualToAnchor:scroll.frameLayoutGuide.heightAnchor],
    [stack.leadingAnchor constraintEqualToAnchor:content.leadingAnchor constant:28],
    [stack.trailingAnchor constraintEqualToAnchor:content.trailingAnchor constant:-28],
    [stack.centerYAnchor constraintEqualToAnchor:content.centerYAnchor],
    [stack.topAnchor constraintGreaterThanOrEqualToAnchor:content.topAnchor constant:28],
    [stack.bottomAnchor constraintLessThanOrEqualToAnchor:content.bottomAnchor constant:-28]
  ]];
  dispatch_async(dispatch_get_global_queue(QOS_CLASS_USER_INITIATED, 0), ^{
    @autoreleasepool {
      [self run];
    }
  });
}
- (void)status:(NSString *)text progress:(float)progress {
  dispatch_async(dispatch_get_main_queue(), ^{
    self.label.text = text;
    [self.progress setProgress:progress animated:!UIAccessibilityIsReduceMotionEnabled()];
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
  NSDictionary *licenses = job[@"license_resources"];
  NSDictionary *archiveResources = job[@"archive_resources"];
  NSArray *archiveDirectories = job[@"archive_directories"];
  NSDictionary *inventory = job[@"payload_inventory"];
  __block NSUInteger processed = 0;
  BOOL ok =
      [inventory isKindOfClass:NSDictionary.class] && inventory.count > 0 && inventory.count <= 100000 &&
      job && [app isKindOfClass:NSString.class] &&
      [app.lastPathComponent isEqualToString:app] && [app hasSuffix:@".app"] &&
      [job[@"schema"] integerValue] == 3 &&
      [licenses isKindOfClass:NSDictionary.class] && licenses.count <= 4096 &&
      [archiveResources isKindOfClass:NSDictionary.class] &&
      [archiveDirectories isKindOfClass:NSArray.class] &&
      archiveResources.count <= 100000 && archiveDirectories.count <= 100000;
  if (!ok && !error)
    error = failure(@"Helper inputs are outdated or invalid. Prepare the "
                    @"helper again in Crumb.");
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
  NSDirectoryEnumerator *entries = nil;
  if (ok)
    entries = [fm enumeratorAtURL:src
        includingPropertiesForKeys:keys
                           options:0
                      errorHandler:^BOOL(NSURL *u, NSError *e) {
                        enumerationError = e;
                        return NO;
                      }];
  if (ok)
    for (NSURL *url in entries) {
      @autoreleasepool {
        NSString *relative = [url.path substringFromIndex:src.path.length + 1];
        // Leave installed SC_Info to FairPlay. Export the exact pre-install
        // metadata from sealed resource copies after processing the images.
        if ([relative.pathComponents containsObject:@"SC_Info"]) {
          [entries skipDescendants];
          continue;
        }
        NSDictionary *a = [url resourceValuesForKeys:keys error:&error];
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
          // Keep paths relative to the embedded target. Installed container
          // UUIDs and absolute sandbox paths are not needed for diagnosis.
          NSMutableDictionary *details =
              [error.userInfo mutableCopy] ?: [NSMutableDictionary dictionary];
          details[NSLocalizedDescriptionKey] =
              [NSString stringWithFormat:@"%@: %@", relative,
                                         error.localizedDescription
                                             ?: @"Image processing failed."];
          NSString *domain = error.domain ?: @"CrumbDecrypt";
          NSInteger code = error ? error.code : 1;
          error = [NSError errorWithDomain:domain code:code userInfo:details];
          ok = NO;
          break;
        }
        if (encrypted)
          decrypted++;
        processed++;
        if (processed == 1 || processed % 25 == 0)
        [self status:[NSString stringWithFormat:@"Processing your app · %lu files", (unsigned long)processed]
             progress:0.8f * MIN(1.0f, (float)processed / MAX(1, inventory.count))];
      }
    }
  if (enumerationError) {
    ok = NO;
    error = enumerationError;
  }
  [self status:@"Restoring app resources…" progress:0.82];
  uint64_t licenseTotal = 0;
  if (ok)
    for (NSString *relative in licenses) {
      NSDictionary *resource = licenses[relative];
      if (![resource isKindOfClass:NSDictionary.class] ||
          ![resource[@"size"] isKindOfClass:NSNumber.class]) {
        ok = NO;
        error = failure(@"Invalid license resource entry.");
        break;
      }
      uint64_t size = [resource[@"size"] unsignedLongLongValue];
      if (size > 4 * 1024 * 1024 || size > 4000000000ULL - total ||
          size > 16 * 1024 * 1024 - licenseTotal) {
        ok = NO;
        error = failure(@"License resources exceed the export budget.");
        break;
      }
      total += size;
      licenseTotal += size;
      if (!crumb_copy_license_resource(NSBundle.mainBundle.bundleURL, relative,
                                       resource, dest, &error)) {
        ok = NO;
        error = failure(
            [NSString stringWithFormat:@"%@: %@", relative,
                                       error.localizedDescription
                                           ?: @"License preservation failed."]);
        break;
      }
    }
  if (ok && decrypted == 0) {
    error = failure(
        @"No encrypted images were found; this is not evidence of decryption.");
    ok = NO;
  }
  if (ok)
    for (NSString *relative in archiveDirectories) {
      if (!crumb_create_archive_directory(root, relative, &error)) {
        ok = NO;
        break;
      }
    }
  if (ok)
    for (NSString *relative in archiveResources) {
      NSDictionary *resource = archiveResources[relative];
      if (![resource isKindOfClass:NSDictionary.class] ||
          ![resource[@"size"] isKindOfClass:NSNumber.class]) {
        ok = NO;
        error = failure(@"Invalid original archive resource.");
        break;
      }
      uint64_t size = [resource[@"size"] unsignedLongLongValue];
      if (size > 512 * 1024 * 1024 || size > 4000000000ULL - total) {
        ok = NO;
        error =
            failure(@"Original archive resources exceed the export budget.");
        break;
      }
      total += size;
      if (!crumb_copy_archive_resource(NSBundle.mainBundle.bundleURL, relative,
                                       resource, root, &error)) {
        ok = NO;
        error = failure([NSString
            stringWithFormat:@"%@: %@", relative,
                             error.localizedDescription
                                 ?: @"Original archive preservation failed."]);
        break;
      }
    }
  // Hashes of the generated ZIP alone cannot detect files lost during installation.
  // Compare every pre-install payload entry, including extensions and empty folders.
  [self status:@"Checking all original app files…" progress:0.9];
  if (ok) for (NSString *relative in inventory) {
    NSDictionary *expected = inventory[relative];
    NSURL *url = [dest URLByAppendingPathComponent:relative];
    NSDictionary *actual = [url resourceValuesForKeys:keys error:&error];
    BOOL directory = [expected[@"directory"] boolValue];
    if (!actual || [actual[NSURLIsSymbolicLinkKey] boolValue] ||
        [actual[NSURLIsDirectoryKey] boolValue] != directory ||
        (!directory && (![actual[NSURLIsRegularFileKey] boolValue] ||
        [actual[NSURLFileSizeKey] unsignedLongLongValue] != [expected[@"size"] unsignedLongLongValue]))) {
      ok = NO;
    } else if (!directory && ![expected[@"encrypted"] boolValue]) {
      // Stream hashes to keep large resources within the memory budget.
      NSInputStream *stream = [NSInputStream inputStreamWithURL:url];
      [stream open];
      CC_SHA256_CTX context; CC_SHA256_Init(&context);
      uint8_t buffer[65536]; NSInteger count;
      while ((count = [stream read:buffer maxLength:sizeof(buffer)]) > 0)
        CC_SHA256_Update(&context, buffer, (CC_LONG)count);
      [stream close];
      unsigned char bytes[CC_SHA256_DIGEST_LENGTH]; CC_SHA256_Final(bytes, &context);
      NSMutableString *hash = [NSMutableString string];
      for (NSUInteger i = 0; i < sizeof(bytes); i++) [hash appendFormat:@"%02x", bytes[i]];
      if (count < 0 || ![hash isEqualToString:expected[@"sha256"]]) ok = NO;
    }
    if (!ok) {
      error = failure([NSString stringWithFormat:@"App file missing or changed: %@. Prepare the helper again; no completed IPA was produced.", relative]);
      break;
    }
  }
  [self status:@"Packaging your IPA…" progress:0.95];
  if (ok)
    ok = crumb_pack_ipa(root, [docs URLByAppendingPathComponent:@"result.ipa"],
                        hashes, &error);
  NSDictionary *result = @{
    @"schema" : @3,
    @"job_id" : job[@"job_id"] ?: @"",
    @"complete" : @(ok),
    @"decrypted" : @(decrypted),
    @"files" : hashes,
    @"archive_resources" : archiveResources ?: @{},
    @"archive_directories" : archiveDirectories ?: @[],
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
    UIApplication.sharedApplication.idleTimerDisabled = NO;
    if (ok) self.progress.progress = 1;
    self.titleLabel.text = ok ? @"Your app is ready" : @"Could not finish";
    self.label.text=ok?@"Saving to Crumb…":error.localizedDescription?:@"Decryption failed. No completed IPA was produced.";
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
