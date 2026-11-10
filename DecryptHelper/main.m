#import "CrumbFileAccess.h"
#import "CrumbIPAArchive.h"
#import "CrumbLicenseResources.h"
#import "CrumbMachO.h"
#import "CrumbJobControl.h"
#import <stdatomic.h>
#import <stdbool.h>
static atomic_bool jobCancelled = false;
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

#import "../SharedDecrypt/CrumbDecrypt.h"
static NSError *failure(NSString *s) {
  return [NSError errorWithDomain:@"CrumbDecrypt" code:1
                         userInfo:@{NSLocalizedDescriptionKey:s}];
}
static bool helper_cancelled(void *unused) {
  (void)unused; return atomic_load(&jobCancelled);
}
static BOOL decrypt(NSURL *source, NSURL *dest, BOOL *encrypted, NSError **error) {
  CrumbDecryptMetrics metrics;
  CrumbDecryptHooks hooks={.cancelled=helper_cancelled};
  int status=crumb_decrypt_image(source.fileSystemRepresentation,
                               dest.fileSystemRepresentation,&metrics,&hooks);
  *encrypted=metrics.encrypted;
  if(status<0 && error) *error=failure([NSString stringWithFormat:
      @"Decryption failed (stage %d, errno %d, remap %d/%d, pages %llu/%llu, Mach %d).",
      metrics.stage,metrics.error_code,metrics.remap_result,metrics.remap_errno,
      (unsigned long long)metrics.pages_read,(unsigned long long)metrics.pages_expected,
      metrics.page_result]);
  return status>=0;
}
@interface CrumbHelper : UIViewController
@property(nonatomic, strong) UILabel *label;
@property(nonatomic, strong) UIProgressView *progress;
@property(nonatomic, strong) UILabel *titleLabel;
- (void)status:(NSString *)text progress:(float)progress;
- (void)cancelDecryption;

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
  UIButton *stop = [UIButton buttonWithType:UIButtonTypeSystem];
  [stop setTitle:NSLocalizedString(@"Stop",nil) forState:UIControlStateNormal];
  stop.titleLabel.font = [UIFont preferredFontForTextStyle:UIFontTextStyleHeadline];
  stop.tintColor = UIColor.systemRedColor;
  [stop.heightAnchor constraintGreaterThanOrEqualToConstant:44].active = YES;
  [stop addTarget:self action:@selector(cancelDecryption) forControlEvents:UIControlEventTouchUpInside];
  UIStackView *stack = [[UIStackView alloc] initWithArrangedSubviews:@[brand,self.titleLabel,self.label,self.progress,stop]];
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
    if (atomic_load(&jobCancelled)) return;
    self.label.text = text;
    [self.progress setProgress:progress animated:!UIAccessibilityIsReduceMotionEnabled()];
  });
}
- (void)cancelDecryption {
  atomic_store(&jobCancelled, true);
  UIApplication.sharedApplication.idleTimerDisabled = NO;
  NSDictionary *job = [NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:[NSBundle.mainBundle.bundleURL URLByAppendingPathComponent:@"job.json"]] options:0 error:nil];
  NSURL *docs = [NSFileManager.defaultManager URLsForDirectory:NSDocumentDirectory inDomains:NSUserDomainMask].firstObject;
  crumb_store_cancellation(docs,job[@"job_id"],nil);
  self.titleLabel.text = @"Stopped";
  self.label.text = @"Returning to Crumb…";
  NSURLComponents *callback = [NSURLComponents new];
  callback.scheme = @"crumb"; callback.host = @"device-cancelled";
  callback.queryItems = @[[NSURLQueryItem queryItemWithName:@"job" value:job[@"job_id"]]];
  [UIApplication.sharedApplication openURL:callback.URL options:@{} completionHandler:nil];
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
  if (crumb_job_was_cancelled(docs,job[@"job_id"])) atomic_store(&jobCancelled,true);
  NSString *app = job[@"app_name"];
  NSDictionary *licenses = job[@"license_resources"];
  NSDictionary *archiveResources = job[@"archive_resources"];
  NSArray *archiveDirectories = job[@"archive_directories"];
  NSDictionary *inventory = job[@"payload_inventory"];
  __block NSUInteger processed = 0;
  BOOL ok =
      [inventory isKindOfClass:NSDictionary.class] && inventory.count > 0 && inventory.count <= 100000 &&
      !atomic_load(&jobCancelled) && job && [app isKindOfClass:NSString.class] &&
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
        if (atomic_load(&jobCancelled)) { error = failure(@"Cancelled"); ok = NO; break; }
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
      if (atomic_load(&jobCancelled)) { ok = NO; break; }
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
      if (atomic_load(&jobCancelled)) { ok = NO; break; }
      if (!crumb_create_archive_directory(root, relative, &error)) {
        ok = NO;
        break;
      }
    }
  if (ok)
    for (NSString *relative in archiveResources) {
      if (atomic_load(&jobCancelled)) { ok = NO; break; }
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
    if (atomic_load(&jobCancelled)) { ok = NO; break; }
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
      uint8_t buffer[65536]; NSInteger count = 0;
      while (!atomic_load(&jobCancelled) && (count = [stream read:buffer maxLength:sizeof(buffer)]) > 0)
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
  if (atomic_load(&jobCancelled)) ok = NO;
  if (ok)
    ok = crumb_pack_ipa(root, [docs URLByAppendingPathComponent:@"result.ipa"],
                        hashes, &error);
  if (atomic_load(&jobCancelled)) {
    ok = NO; error = failure(@"Cancelled");
    [fm removeItemAtURL:[docs URLByAppendingPathComponent:@"result.ipa"] error:nil];
  }
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
    if (atomic_load(&jobCancelled)) { self.titleLabel.text = @"Stopped"; return; }
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
  NSDictionary *job = [NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:[NSBundle.mainBundle.bundleURL URLByAppendingPathComponent:@"job.json"]] options:0 error:nil];
  if (crumb_matches_cancel_url(url,job[@"job_id"])) {
    [(CrumbHelper *)self.window.rootViewController cancelDecryption];
  }
  return YES;
}
@end
int main(int argc, char **argv) {
  @autoreleasepool {
    return UIApplicationMain(argc, argv, nil,
                             NSStringFromClass(Delegate.class));
  }
}
