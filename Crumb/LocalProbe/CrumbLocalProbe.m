#import "CrumbLocalProbe.h"
#import "../../SharedDecrypt/CrumbDecrypt.h"
#import "../../NativeDevice/include/CrumbDevice.h"
#import "Vendor/worse_query.h"
#import <CommonCrypto/CommonDigest.h>
#import <Security/Security.h>
#import <TargetConditionals.h>
#if TARGET_OS_IOS
#import <UIKit/UIKit.h>
#endif
#import <dlfcn.h>
#import <objc/message.h>
#import <stdatomic.h>
#import <sys/stat.h>
#import <fcntl.h>
#import <unistd.h>
#import <sys/sysctl.h>

static void set_idle(BOOL disabled) {
#if TARGET_OS_IOS
  dispatch_async(dispatch_get_main_queue(),^{UIApplication.sharedApplication.idleTimerDisabled=disabled;});
#else
  (void)disabled;
#endif
}
static NSString *os_version(void) {
#if TARGET_OS_IOS
  return UIDevice.currentDevice.systemVersion;
#else
  return NSProcessInfo.processInfo.operatingSystemVersionString;
#endif
}
static atomic_bool processHasGrant=false;
static NSString *stage_name(int stage) {
  switch(stage) {
    case CRUMB_OPEN:return @"open";case CRUMB_INSPECT:return @"inspect";
    case CRUMB_MAP:return @"mmap";case CRUMB_PARSE:return @"parse";
    case CRUMB_COPY:return @"copy";case CRUMB_SIGNATURE:return @"signature";
    case CRUMB_REMAP:return @"remap";case CRUMB_READ:return @"page_read";
    case CRUMB_VERIFY:return @"verify_pages";case CRUMB_PATCH:return @"patch_output";
    case CRUMB_COMPLETE:return @"image_complete";default:return @"unknown";
  }
}
static NSDictionary *metrics_json(const CrumbDecryptMetrics *m) {
  return @{@"stage":stage_name(m->stage),@"errno":@(m->error_code),
    @"signature_attempted":@(m->signature_attempted),@"signature_result":@(m->signature_result),@"signature_errno":@(m->signature_errno),
    @"first_remap_result":@(m->first_remap_result),@"first_remap_errno":@(m->first_remap_errno),@"remap_attempts":@(m->remap_attempts),
    @"remap_attempted":@(m->remap_attempted),@"remap_result":@(m->remap_result),@"remap_errno":@(m->remap_errno),
    @"pages_expected":@(m->pages_expected),@"pages_read":@(m->pages_read),
    @"changed_pages":@(m->changed_pages),@"first_fail_offset":@(m->first_fail_offset),
    @"mach_result":@(m->page_result),@"read_errno":@(m->read_errno),@"write_errno":@(m->write_errno),
    @"entropy_before":@(m->entropy_before),@"entropy_after":@(m->entropy_after),
    @"codewords_before":@(m->oldwords),@"codewords_after":@(m->newwords),@"cancelled":@(m->cancelled)};
}
@interface CrumbLocalProbe () {
@public
  atomic_bool _cancelled;
}
@property(nonatomic) NSURL *directory;
@property(nonatomic, readwrite) NSURL *outputURL;
@property(nonatomic, readwrite) NSURL *reportURL;
@property(nonatomic) NSMutableDictionary *report;
@property(nonatomic, copy) void (^progress)(NSString *);
@property(nonatomic) NSString *lastProgress;
@end
static bool cancelled(void *ctx) {return atomic_load(&((__bridge CrumbLocalProbe *)ctx)->_cancelled);}
static void native_checkpoint(void *ctx,const CrumbDecryptMetrics *m);
@implementation CrumbLocalProbe
- (instancetype)initWithDirectory:(NSURL *)directory {
  if((self=[super init])) {
    atomic_init(&_cancelled,false);
    _directory=directory;
    NSFileManager *fm=NSFileManager.defaultManager;
    [fm createDirectoryAtURL:directory withIntermediateDirectories:YES attributes:@{NSFilePosixPermissions:@0700} error:nil];
    NSURL *excluded=directory;[excluded setResourceValue:@YES forKey:NSURLIsExcludedFromBackupKey error:nil];
    [self recover];
  }return self;
}
- (void)cancel {atomic_store(&_cancelled,true);}
- (BOOL)saveJSON:(NSDictionary *)json to:(NSURL *)url {
  NSData *data=[NSJSONSerialization dataWithJSONObject:json options:NSJSONWritingPrettyPrinted|NSJSONWritingSortedKeys error:nil];
  return data && [data writeToURL:url options:NSDataWritingAtomic|NSDataWritingFileProtectionComplete error:nil];
}
- (void)checkpoint:(NSString *)stage metrics:(NSDictionary *)metrics {
  self.report[@"current_stage"]=stage;
  if(metrics)self.report[@"last_native"]=metrics;
  if(![self saveJSON:self.report to:self.reportURL])atomic_store(&_cancelled,true);
  if(![stage isEqual:self.lastProgress]){self.lastProgress=stage;if(self.progress)self.progress(stage);}
}
- (void)recover {
  NSFileManager *fm=NSFileManager.defaultManager;
  NSURL *active=[self.directory URLByAppendingPathComponent:@"active.json"];
  NSDictionary *state=[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:active]?:[NSData data] options:0 error:nil];
  NSString *job=state[@"run_id"];
  if([job isKindOfClass:NSString.class] && [[NSUUID alloc]initWithUUIDString:job]) {
    NSURL *report=[self.directory URLByAppendingPathComponent:[NSString stringWithFormat:@"report-%@.json",job]];
    NSMutableDictionary *r=[[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:report]?:[NSData data] options:NSJSONReadingMutableContainers error:nil] mutableCopy];
    if(r){r[@"outcome"]=@"interrupted";[self saveJSON:r to:report];}
  }
  [fm removeItemAtURL:active error:nil];
  NSURL *sharedLedger=[self.directory URLByAppendingPathComponent:@"shared-job.json"];
  NSDictionary *sharedState=[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:sharedLedger]?:[NSData data] options:0 error:nil];
  NSString *sharedName=sharedState[@"name"], *sharedGroup=sharedState[@"group"];
  NSString *ownGroup=[NSString stringWithFormat:@"group.%@.worsequery",NSBundle.mainBundle.bundleIdentifier];
  if([sharedGroup isEqual:ownGroup] && [sharedName isKindOfClass:NSString.class] && [sharedName hasPrefix:@"CrumbProbe-"] && [[NSUUID alloc]initWithUUIDString:[sharedName substringFromIndex:11]]) {
    NSURL *sharedRoot=[fm containerURLForSecurityApplicationGroupIdentifier:sharedGroup];
    if(sharedRoot)[fm removeItemAtURL:[sharedRoot URLByAppendingPathComponent:sharedName] error:nil];
  }
  [fm removeItemAtURL:sharedLedger error:nil];
  // Only UUID job directories in our own private root are reaped.
  NSURL *runs=[self.directory URLByAppendingPathComponent:@"runs"];
  for(NSURL *url in [fm contentsOfDirectoryAtURL:runs includingPropertiesForKeys:nil options:0 error:nil])
    if([[NSUUID alloc]initWithUUIDString:url.lastPathComponent])[fm removeItemAtURL:url error:nil];
  NSURL *ledger=[self.directory URLByAppendingPathComponent:@"group-link.json"];
  NSDictionary *link=[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:ledger]?:[NSData data] options:0 error:nil];
  NSString *group=link[@"group"],*expected=link[@"target"];
  NSString *own=[@"group." stringByAppendingFormat:@"%@.worsequery",NSBundle.mainBundle.bundleIdentifier];
  if([group isEqual:own] && [expected isKindOfClass:NSString.class]) {
    NSURL *root=[fm containerURLForSecurityApplicationGroupIdentifier:group];
    NSURL *caches=[root URLByAppendingPathComponent:@"Library/Caches"];
    NSString *actual=caches?[fm destinationOfSymbolicLinkAtPath:caches.path error:nil]:nil;
    if([actual isEqual:expected]) {
      if(unlink(caches.fileSystemRepresentation)==0)[fm createDirectoryAtURL:caches withIntermediateDirectories:YES attributes:nil error:nil];
    }
  }
  [fm removeItemAtURL:ledger error:nil];
}
- (NSDictionary *)archive:(NSDictionary *)request {
  NSData *data=[NSJSONSerialization dataWithJSONObject:request options:0 error:nil];
  NSString *json=[[NSString alloc]initWithData:data encoding:NSUTF8StringEncoding];
  if(!json)return @{@"ok":@NO,@"code":@"invalid_archive_request"};
  char *raw=crumb_probe_archive(json.UTF8String,(__bridge void *)self,cancelled);
  if(!raw)return @{@"ok":@NO,@"code":@"archive_no_response"};
  NSData *response=[NSData dataWithBytes:raw length:strlen(raw)];crumb_device_string_free(raw);
  return [NSJSONSerialization JSONObjectWithData:response options:0 error:nil]?:@{@"ok":@NO,@"code":@"archive_invalid_response"};
}
- (NSArray *)effectiveGroups {
  // Resolve SecTask dynamically: the public SDK does not declare these SPI types.
  typedef CFTypeRef (*CreateTask)(CFAllocatorRef);
  typedef CFTypeRef (*CopyEntitlement)(CFTypeRef,CFStringRef,CFErrorRef *);
  CreateTask create=(CreateTask)dlsym(RTLD_DEFAULT,"SecTaskCreateFromSelf");
  CopyEntitlement copy=(CopyEntitlement)dlsym(RTLD_DEFAULT,"SecTaskCopyValueForEntitlement");
  if(!create||!copy)return nil;
  CFTypeRef task=create(kCFAllocatorDefault);if(!task)return nil;
  CFTypeRef value=copy(task,CFSTR("com.apple.security.application-groups"),NULL);CFRelease(task);
  id groups=CFBridgingRelease(value);return [groups isKindOfClass:NSArray.class]?groups:nil;
}
- (NSString *)hash:(NSURL *)url {
  errno=0;int fd=open(url.fileSystemRepresentation,O_RDONLY|O_NOFOLLOW|O_CLOEXEC);if(fd<0)return nil;
  struct stat st;if(fstat(fd,&st)){int code=errno;close(fd);errno=code;return nil;}
  if(!S_ISREG(st.st_mode)||st.st_size<0||st.st_size>512LL*1024*1024){close(fd);errno=EFBIG;return nil;}
  CC_SHA256_CTX hash;CC_SHA256_Init(&hash);uint8_t buffer[65536];ssize_t n;
  while(!atomic_load(&_cancelled) && (n=read(fd,buffer,sizeof(buffer)))>0)CC_SHA256_Update(&hash,buffer,(CC_LONG)n);
  int saved=errno;close(fd);if(atomic_load(&_cancelled)||n<0){errno=saved;return nil;}
  uint8_t digest[32];CC_SHA256_Final(digest,&hash);NSMutableString *s=[NSMutableString string];
  for(int i=0;i<32;i++)[s appendFormat:@"%02x",digest[i]];return s;
}
- (NSDictionary *)readCanary:(NSURL *)url expected:(NSData *)expected {
  errno=0;int fd=open(url.fileSystemRepresentation,O_RDONLY|O_NOFOLLOW|O_CLOEXEC);int code=fd<0?errno:0;
  BOOL matches=NO;
  if(fd>=0){uint8_t b[64];ssize_t n=read(fd,b,sizeof(b));matches=n==(ssize_t)expected.length && memcmp(b,expected.bytes,expected.length)==0;close(fd);}
  return @{@"open":@(fd>=0),@"errno":@(code),@"matches":@(matches)};
}
- (int64_t)grant:(NSURL *)target group:(NSString *)group result:(NSMutableDictionary *)result {
  if(atomic_load(&_cancelled))return -256;
  NSFileManager *fm=NSFileManager.defaultManager;
  NSURL *root=[fm containerURLForSecurityApplicationGroupIdentifier:group];
  if(!root){result[@"code"]=@"group_container_unavailable";return -1;}
  NSURL *caches=[root URLByAppendingPathComponent:@"Library/Caches"];
  struct stat st;
  if(lstat(caches.fileSystemRepresentation,&st)==0) {
    if(!S_ISDIR(st.st_mode)){result[@"code"]=@"group_caches_not_directory";return -256;}
    NSArray *contents=[fm contentsOfDirectoryAtURL:caches includingPropertiesForKeys:nil options:0 error:nil];
    if(!contents || contents.count){result[@"code"]=@"group_caches_not_empty";return -256;}
  }else if(errno!=ENOENT){result[@"code"]=@"group_caches_uninspectable";return -256;}
  [fm createDirectoryAtURL:caches.URLByDeletingLastPathComponent withIntermediateDirectories:YES attributes:nil error:nil];
  NSURL *ledger=[self.directory URLByAppendingPathComponent:@"group-link.json"];
  if(![self saveJSON:@{@"group":group,@"target":target.path} to:ledger]){result[@"code"]=@"grant_ledger_failed";return -256;}
  [self checkpoint:@"sandbox_extension" metrics:nil];
  BOOL alt=NSProcessInfo.processInfo.operatingSystemVersion.majorVersion==26 && [target.path containsString:@"/Containers/Shared/AppGroup/"];
  struct stat targetStatus;errno=0;
  int targetStat=lstat(target.fileSystemRepresentation,&targetStatus);
  result[@"target_stat_result"]=@(targetStat);result[@"target_stat_errno"]=@(targetStat?errno:0);
  if(atomic_load(&_cancelled))return -256;
  int64_t handle=worse_query(target.path,group,alt,false);result[@"result"]=@(handle);
  if(handle>=0){atomic_store(&processHasGrant,true);result[@"code"]=@"extension_consumed";}
  else result[@"code"]=@"extension_rejected";
  // Verify cleanup, including the upstream's early-return paths.
  NSString *link=[fm destinationOfSymbolicLinkAtPath:caches.path error:nil];
  if([link isEqual:target.path])unlink(caches.fileSystemRepresentation);
  if(![fm fileExistsAtPath:caches.path])[fm createDirectoryAtURL:caches withIntermediateDirectories:YES attributes:nil error:nil];
  struct stat after;BOOL clean=lstat(caches.fileSystemRepresentation,&after)==0 && S_ISDIR(after.st_mode);
  result[@"caches_restored"]=@(clean);
  if(clean)[fm removeItemAtURL:ledger error:nil];
  else{atomic_store(&_cancelled,true);result[@"code"]=@"group_cleanup_failed";}
  return handle;
}
- (NSURL *)installedBundle:(NSString *)bundleID result:(NSMutableDictionary *)result {
  [self checkpoint:@"installed_discovery" metrics:nil];
  @try {
    Class cls=NSClassFromString(@"LSApplicationProxy");SEL selector=NSSelectorFromString(@"applicationProxyForIdentifier:");
    if(![cls respondsToSelector:selector]){result[@"code"]=@"discovery_api_unavailable";return nil;}
    id proxy=((id(*)(id,SEL,id))objc_msgSend)(cls,selector,bundleID);
    SEL urlSelector=NSSelectorFromString(@"bundleURL");
    NSURL *url=[proxy respondsToSelector:urlSelector]?((id(*)(id,SEL))objc_msgSend)(proxy,urlSelector):nil;
    if(![url isKindOfClass:NSURL.class] || !url.fileURL){result[@"code"]=@"discovery_blocked";return nil;}
    url=url.URLByResolvingSymlinksInPath;
    if(![url.path hasPrefix:@"/private/var/containers/Bundle/Application/"] && ![url.path hasPrefix:@"/var/containers/Bundle/Application/"]){result[@"code"]=@"unexpected_bundle_location";return nil;}
    result[@"code"]=@"bundle_location_resolved";return url;
  }@catch(NSException *e){(void)e;result[@"code"]=@"discovery_blocked";return nil;}
}
- (BOOL)attempt:(NSString *)kind access:(NSString *)access source:(NSURL *)source
       manifest:(NSDictionary *)manifest work:(NSURL *)work group:(NSString *)group {
  NSMutableDictionary *attempt=[@{@"source":kind,@"access":access,@"images":[NSMutableArray array],@"outcome":@"started"} mutableCopy];
  [self.report[@"attempts"] addObject:attempt];[self checkpoint:@"attempt_start" metrics:nil];
  NSFileManager *fm=NSFileManager.defaultManager;int64_t handle=-1;
  NSURL *candidate=[work URLByAppendingPathComponent:NSUUID.UUID.UUIDString isDirectory:YES];
  NSMutableArray *receipts=[NSMutableArray array],*outputs=[NSMutableArray array];BOOL success=YES;
  @try {
    if([access isEqual:@"worseQuery"]){
      NSMutableDictionary *grant=[NSMutableDictionary dictionary];attempt[@"grant"]=grant;
      handle=[self grant:source group:group result:grant];
      if(handle<0){attempt[@"outcome"]=@"grant_unavailable";return NO;}
    }
    if([kind isEqual:@"installedApp"]){
      NSDictionary *d=[NSDictionary dictionaryWithContentsOfURL:[source URLByAppendingPathComponent:@"Info.plist"]];
      if(![d[@"CFBundleIdentifier"] isEqual:manifest[@"bundle_id"]] || ![d[@"CFBundleShortVersionString"] isEqual:manifest[@"version"]] || ![d[@"CFBundleVersion"] isEqual:manifest[@"build"]]) {
        attempt[@"outcome"]=d?@"installed_version_mismatch":@"installed_metadata_unreadable";return NO;
      }
    }
    if(![fm createDirectoryAtURL:candidate withIntermediateDirectories:YES attributes:@{NSFilePosixPermissions:@0700} error:nil]){attempt[@"outcome"]=@"workspace_failed";return NO;}
    NSUInteger index=0;
    for(NSDictionary *image in manifest[@"images"]){
      if(atomic_load(&_cancelled)){success=NO;break;}
      NSMutableDictionary *record=[@{@"index":@(index),@"image_sha256":image[@"sha256"],@"uuid":image[@"uuid"]?:NSNull.null,@"expected_size":image[@"size"],@"matching_sinf_staged":image[@"matching_sinf"],@"sinf_size":image[@"sinf_size"]} mutableCopy];
      [attempt[@"images"] addObject:record];self.report[@"active_image"]=@(index);
      NSString *relative=[kind isEqual:@"installedApp"]?image[@"app_relative"]:image[@"relative"];
      NSURL *file=[source URLByAppendingPathComponent:relative];
      NSString *root=[source.URLByResolvingSymlinksInPath.path stringByAppendingString:@"/"];
      if(![file.URLByResolvingSymlinksInPath.path hasPrefix:root]){record[@"outcome"]=@"image_path_escaped";success=NO;break;}
      if([kind isEqual:@"installedApp"]) {
        NSURL *sinf=[[file.URLByDeletingLastPathComponent URLByAppendingPathComponent:@"SC_Info"] URLByAppendingPathComponent:[file.lastPathComponent stringByAppendingString:@".sinf"]];
        struct stat status;errno=0;int present=lstat(sinf.fileSystemRepresentation,&status);int code=present?errno:0;
        record[@"installed_sinf_metadata"]=@{@"stat_result":@(present),@"errno":@(code),@"size":@(present?0:status.st_size)};
        record[@"license_context"]=@"installed_location_daemon_access_unobserved";
      }else record[@"license_context"]=@"original_staged_license";
      [self checkpoint:@"source_hash" metrics:nil];
      NSString *hash=[self hash:file];
      if(![hash isEqual:image[@"sha256"]]) {
        int code=errno;record[@"outcome"]=hash?@"installed_image_mismatch":@"source_unreadable";record[@"errno"]=@(code);success=NO;break;
      }
      if(![kind isEqual:@"installedApp"] && ![image[@"matching_sinf"] boolValue]){record[@"outcome"]=@"matching_license_missing";success=NO;break;}
      NSURL *destination=[candidate URLByAppendingPathComponent:[NSString stringWithFormat:@"%lu",(unsigned long)index]];
      CrumbDecryptHooks hooks={.context=(__bridge void *)self,.cancelled=cancelled,.checkpoint=native_checkpoint};CrumbDecryptMetrics metrics;
      int result=crumb_decrypt_image(file.fileSystemRepresentation,destination.fileSystemRepresentation,&metrics,&hooks);
      record[@"native"]=metrics_json(&metrics);
      record[@"outcome"]=result==1?@"pages_verified":@"native_failed";
      [self checkpoint:@"image_result" metrics:record[@"native"]];
      // Detect replacement of the source while native work was in progress.
      if(result!=1 || ![[self hash:file] isEqual:hash]){success=NO;break;}
      [receipts addObject:@{@"relative":image[@"relative"],@"image_sha256":hash,@"verified":@YES,@"page_size":@(getpagesize()),@"pages_expected":@(metrics.pages_expected),@"pages_read":@(metrics.pages_read)}];
      [outputs addObject:destination];index++;
    }
    if(!success || receipts.count!=[manifest[@"images"] count]){attempt[@"outcome"]=atomic_load(&_cancelled)?@"cancelled":@"incomplete";return NO;}
    NSArray *images=manifest[@"images"];
    for(NSUInteger i=0;i<outputs.count;i++) {
      NSURL *dest=[[work URLByAppendingPathComponent:@"archive/output"] URLByAppendingPathComponent:images[i][@"relative"]];
      NSDictionary *attrs=[fm attributesOfItemAtPath:dest.path error:nil];
      if(!attrs || ![fm removeItemAtURL:dest error:nil] || ![fm moveItemAtURL:outputs[i] toURL:dest error:nil] || ![fm setAttributes:@{NSFilePosixPermissions:attrs[NSFilePosixPermissions]} ofItemAtPath:dest.path error:nil]){
        attempt[@"outcome"]=@"output_staging_failed";atomic_store(&_cancelled,true);return NO;
      }
    }
    [self checkpoint:@"archive_verification" metrics:nil];
    NSDictionary *verified=[self archive:@{@"operation":@"finalize",@"work":[[work URLByAppendingPathComponent:@"archive"] path],@"receipts":receipts}];
    if(![verified[@"ok"] boolValue]){attempt[@"outcome"]=@"archive_verification_failed";attempt[@"archive_error"]=verified;atomic_store(&_cancelled,true);return NO;}
    NSURL *outputsRoot=[self.directory URLByAppendingPathComponent:@"outputs"];
    [fm createDirectoryAtURL:outputsRoot withIntermediateDirectories:YES attributes:nil error:nil];
    NSURL *final=[outputsRoot URLByAppendingPathComponent:[self.report[@"run_id"] stringByAppendingString:@".ipa"]];
    if(![fm moveItemAtURL:[NSURL fileURLWithPath:verified[@"output_path"]] toURL:final error:nil]){attempt[@"outcome"]=@"publication_failed";atomic_store(&_cancelled,true);return NO;}
    self.outputURL=final;attempt[@"outcome"]=@"complete";
    self.report[@"verified_images"]=verified[@"image_count"];self.report[@"verified_files"]=verified[@"file_count"];
    self.report[@"outcome"]=[kind isEqual:@"installedApp"]?@"installed_app_success":@"downloaded_ipa_success";
    return YES;
  }@finally {
    [fm removeItemAtURL:candidate error:nil];
    if(handle>=0)worse_query_release(handle);
  }
}
- (NSDictionary *)runIPA:(NSURL *)ipa bundleID:(NSString *)bundleID progress:(void (^)(NSString *))progress {
  atomic_store(&_cancelled,false);self.outputURL=nil;self.progress=progress;self.lastProgress=nil;
  NSString *job=NSUUID.UUID.UUIDString;
  self.reportURL=[self.directory URLByAppendingPathComponent:[NSString stringWithFormat:@"report-%@.json",job]];
  char model[256]={0};size_t modelSize=sizeof(model);sysctlbyname("hw.machine",model,&modelSize,NULL,0);
  char build[256]={0};size_t buildSize=sizeof(build);sysctlbyname("kern.osversion",build,&buildSize,NULL,0);
  self.report=[@{@"schema":@1,@"run_id":job,@"outcome":@"running",@"attempts":[NSMutableArray array],
    @"environment":@{@"model":@(model),@"os":os_version(),@"build":@(build),@"page_size":@(getpagesize()),
      @"crumb_build":[NSBundle.mainBundle objectForInfoDictionaryKey:@"CFBundleVersion"]?:@"unknown",
      @"worse_query_revision":@"d1540e20cd82b974d891e0551ced6d6aeb570ac2",@"baseline_contaminated":@(atomic_load(&processHasGrant))},
    @"daemon_observation":@"unavailable_without_external_trace",@"plaintext_authentication":@"not_established"} mutableCopy];
  NSFileManager *fm=NSFileManager.defaultManager;
  NSURL *work=[[self.directory URLByAppendingPathComponent:@"runs"] URLByAppendingPathComponent:job];
  NSURL *active=[self.directory URLByAppendingPathComponent:@"active.json"];
  [self saveJSON:@{@"run_id":job} to:active];
  [fm createDirectoryAtURL:work withIntermediateDirectories:YES attributes:@{NSFilePosixPermissions:@0700} error:nil];
  set_idle(YES);
  @try {
    [self checkpoint:@"archive_preparation" metrics:nil];
    NSDictionary *prep=[self archive:@{@"operation":@"prepare",@"work":[[work URLByAppendingPathComponent:@"archive"] path],@"ipa":ipa.path,@"bundle_id":bundleID}];
    if(![prep[@"ok"] boolValue]){self.report[@"outcome"]=@"preparation_failed";self.report[@"archive_error"]=prep;return self.report;}
    NSDictionary *manifest=prep[@"manifest"];
    self.report[@"image_count"]=@([manifest[@"images"] count]);
    NSURL *buildURL=[NSBundle.mainBundle URLForResource:@"CrumbProbeBuild" withExtension:@"json"];
    NSDictionary *buildInfo=buildURL?[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:buildURL]?:[NSData data] options:0 error:nil]:nil;
    self.report[@"source_revision"]=[buildInfo[@"source_revision"] isKindOfClass:NSString.class]?buildInfo[@"source_revision"]:@"unrecorded";
    NSString *group=[NSString stringWithFormat:@"group.%@.worsequery",NSBundle.mainBundle.bundleIdentifier];
    NSArray *groups=[self effectiveGroups];BOOL authorized=[groups containsObject:group];
    NSURL *groupRoot=authorized?[fm containerURLForSecurityApplicationGroupIdentifier:group]:nil;
    self.report[@"app_group"]=@{@"entitlements_readable":@(groups!=nil),@"dedicated_group_authorized":@(authorized),@"container_available":@(groupRoot!=nil)};
    NSMutableDictionary *discovery=[NSMutableDictionary dictionary];
    NSURL *installed=[self installedBundle:bundleID result:discovery];self.report[@"installed_discovery"]=discovery;
    NSURL *input=[work URLByAppendingPathComponent:@"archive/input"];
    BOOL contaminated=atomic_load(&processHasGrant);
    if(!contaminated) {
      if([self attempt:@"downloadedIPA" access:@"baseline" source:input manifest:manifest work:work group:group])return self.report;
      if(installed && [self attempt:@"installedApp" access:@"baseline" source:installed manifest:manifest work:work group:group])return self.report;
    }else self.report[@"baseline_skipped"]=@"restart_crumb_for_clean_baseline";
    if(groupRoot && !atomic_load(&_cancelled)) {
      NSURL *canary=[work URLByAppendingPathComponent:@"canary"];NSData *bytes=[job dataUsingEncoding:NSUTF8StringEncoding];
      [bytes writeToURL:canary options:NSDataWritingAtomic error:nil];
      NSMutableDictionary *control=[NSMutableDictionary dictionary];control[@"before"]=[self readCanary:canary expected:bytes];
      int64_t handle=[self grant:[NSURL fileURLWithPath:NSHomeDirectory() isDirectory:YES] group:group result:control];
      control[@"after"]=[self readCanary:canary expected:bytes];if(handle>=0)worse_query_release(handle);
      control[@"note"]=@"own_container_control_does_not_demonstrate_new_cross_app_access";self.report[@"canary"]=control;
      if(!atomic_load(&_cancelled)) {
        if([self attempt:@"downloadedIPA" access:@"worseQuery" source:input manifest:manifest work:work group:group])return self.report;
        if(installed && [self attempt:@"installedApp" access:@"worseQuery" source:installed manifest:manifest work:work group:group])return self.report;
      }
      if(!atomic_load(&_cancelled)) {
        // One distinct storage-class experiment in a dedicated authorized group.
        NSString *sharedJob=[@"CrumbProbe-" stringByAppendingString:job];
        NSURL *shared=[groupRoot URLByAppendingPathComponent:sharedJob isDirectory:YES];
        NSURL *ledger=[self.directory URLByAppendingPathComponent:@"shared-job.json"];
        if([self saveJSON:@{@"group":group,@"name":sharedJob} to:ledger]) {
          @try {
            [self checkpoint:@"shared_container_copy" metrics:nil];
            BOOL created=[fm createDirectoryAtURL:shared withIntermediateDirectories:NO attributes:@{NSFilePosixPermissions:@0700,NSFileProtectionKey:NSFileProtectionComplete} error:nil];
            BOOL excluded=created && [shared setResourceValue:@YES forKey:NSURLIsExcludedFromBackupKey error:nil];
            NSDictionary *copied=excluded?[self archive:@{@"operation":@"relocate",@"from":input.path,@"to":shared.path}]:@{@"ok":@NO};
            if([copied[@"ok"] boolValue]) {
              [self attempt:@"ownAppGroup" access:@"authorizedGroup" source:shared manifest:manifest work:work group:group];
            }else self.report[@"shared_container"]=@"copy_failed";
          }@finally{[fm removeItemAtURL:shared error:nil];[fm removeItemAtURL:ledger error:nil];}
        }
      }
    }
    if(!self.outputURL)self.report[@"outcome"]=atomic_load(&_cancelled)?@"cancelled_or_local_failure":@"no_verified_route";
  }@catch(NSException *e){(void)e;self.report[@"outcome"]=@"local_exception";}
  @finally {
    if(atomic_load(&_cancelled)) {
      if(self.outputURL)[fm removeItemAtURL:self.outputURL error:nil];
      self.outputURL=nil;self.report[@"outcome"]=@"cancelled_or_local_failure";
    }
    [self checkpoint:@"finished" metrics:nil];
    [fm removeItemAtURL:work error:nil];[fm removeItemAtURL:active error:nil];
    set_idle(NO);
    self.progress=nil;
  }
  return self.report;
}
@end
static void native_checkpoint(void *ctx,const CrumbDecryptMetrics *m) {
  if(m->stage==CRUMB_READ && m->pages_read%32)return;
  [(__bridge CrumbLocalProbe *)ctx checkpoint:stage_name(m->stage) metrics:metrics_json(m)];
}
