#import "CrumbLocalProbe.h"
#import "../../SharedDecrypt/CrumbDecrypt.h"
#import "../../SharedDecrypt/CrumbCacheLease.h"
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
static BOOL remove_owned(NSURL *url) {
  if([NSFileManager.defaultManager removeItemAtURL:url error:nil])return YES;
  struct stat status;errno=0;
  return lstat(url.fileSystemRepresentation,&status)<0 && errno==ENOENT;
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
- (NSArray *)effectiveGroups;
- (NSURL *)groupRoot:(NSString *)group;
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
    if(sharedRoot && remove_owned([sharedRoot URLByAppendingPathComponent:sharedName]))[fm removeItemAtURL:sharedLedger error:nil];
  }else if(!sharedState)[fm removeItemAtURL:sharedLedger error:nil];
  NSURL *leaseLedger=[self.directory URLByAppendingPathComponent:@"group-lease.json"];
  NSDictionary *leaseState=[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:leaseLedger]?:[NSData data] options:0 error:nil];
  if(leaseState) {
    NSString *group=leaseState[@"group"], *identifier=leaseState[@"identifier"], *target=leaseState[@"target"];
    NSArray *effective=[self effectiveGroups];
    if([group isKindOfClass:NSString.class] && [effective containsObject:group] &&
       [identifier isKindOfClass:NSString.class] && [[NSUUID alloc]initWithUUIDString:identifier] &&
       [target isKindOfClass:NSString.class]) {
      NSURL *root=[self groupRoot:group];
      NSURL *stage=[[self.directory URLByAppendingPathComponent:@"cache-links"] URLByAppendingPathComponent:identifier];
      CrumbCacheLease lease={.lock_fd=-1};
      if(root && root.path.length<PATH_MAX && target.length<PATH_MAX && stage.path.length<PATH_MAX) {
        snprintf(lease.root,sizeof(lease.root),"%s",root.fileSystemRepresentation);
        snprintf(lease.target,sizeof(lease.target),"%s",target.fileSystemRepresentation);
        snprintf(lease.staging,sizeof(lease.staging),"%s",stage.fileSystemRepresentation);
        snprintf(lease.identifier,sizeof(lease.identifier),"%s",identifier.UTF8String);
        lease.had_caches=[leaseState[@"had_caches"] boolValue];
        lease.original_device=[leaseState[@"original_device"] unsignedLongLongValue];
        lease.original_inode=[leaseState[@"original_inode"] unsignedLongLongValue];
        lease.link_device=[leaseState[@"link_device"] unsignedLongLongValue];
        lease.link_inode=[leaseState[@"link_inode"] unsignedLongLongValue];
        if(crumb_cache_restore(&lease)==0)[fm removeItemAtURL:leaseLedger error:nil];
      }
    }
  }
  if(![fm fileExistsAtPath:leaseLedger.path]) {
    NSURL *links=[self.directory URLByAppendingPathComponent:@"cache-links"];
    for(NSURL *url in [fm contentsOfDirectoryAtURL:links includingPropertiesForKeys:nil options:0 error:nil]) {
      if([[NSUUID alloc]initWithUUIDString:url.lastPathComponent]) {
        struct stat status;
        if(lstat(url.fileSystemRepresentation,&status)==0 && S_ISLNK(status.st_mode))unlink(url.fileSystemRepresentation);
      }
    }
  }
  // Only UUID job directories in our own private root are reaped.
  NSURL *runs=[self.directory URLByAppendingPathComponent:@"runs"];
  for(NSURL *url in [fm contentsOfDirectoryAtURL:runs includingPropertiesForKeys:nil options:0 error:nil])
    if([[NSUUID alloc]initWithUUIDString:url.lastPathComponent])[fm removeItemAtURL:url error:nil];
  NSURL *ledger=[self.directory URLByAppendingPathComponent:@"group-link.json"];
  NSDictionary *link=[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:ledger]?:[NSData data] options:0 error:nil];
  NSString *group=link[@"group"],*expected=link[@"target"];
  NSString *own=[@"group." stringByAppendingFormat:@"%@.worsequery",NSBundle.mainBundle.bundleIdentifier];
  BOOL linkClean=NO;
  if([group isEqual:own] && [expected isKindOfClass:NSString.class]) {
    NSURL *root=[fm containerURLForSecurityApplicationGroupIdentifier:group];
    NSURL *caches=[root URLByAppendingPathComponent:@"Library/Caches"];
    NSString *actual=caches?[fm destinationOfSymbolicLinkAtPath:caches.path error:nil]:nil;
    if([actual isEqual:expected]) {
      if(unlink(caches.fileSystemRepresentation)==0)[fm createDirectoryAtURL:caches withIntermediateDirectories:YES attributes:nil error:nil];
    }
    struct stat status;
    linkClean=caches && lstat(caches.fileSystemRepresentation,&status)==0 && S_ISDIR(status.st_mode);
  }
  if(linkClean || !link)[fm removeItemAtURL:ledger error:nil];
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
- (NSURL *)groupRoot:(NSString *)group {
  return [NSFileManager.defaultManager containerURLForSecurityApplicationGroupIdentifier:group];
}
- (NSString *)selectGroup:(NSArray *)effective diagnostics:(NSMutableDictionary *)result root:(NSURL *__autoreleasing *)selectedRoot {
  NSString *expected=[NSString stringWithFormat:@"group.%@.worsequery",NSBundle.mainBundle.bundleIdentifier];
  NSMutableSet *valid=[NSMutableSet set];
  NSCharacterSet *invalid=[[NSCharacterSet characterSetWithCharactersInString:@"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-"] invertedSet];
  for(id value in effective) {
    if([value isKindOfClass:NSString.class] && [value hasPrefix:@"group."] && [value length]>6 &&
       [value length]<256 && [value rangeOfCharacterFromSet:invalid].location==NSNotFound &&
       ![value containsString:@".."] && ![value hasSuffix:@"."])[valid addObject:value];
  }
  NSArray *groups=[[valid allObjects] sortedArrayUsingComparator:^NSComparisonResult(NSString *a,NSString *b) {
    int rankA=[a isEqual:expected]?0:([a hasSuffix:@".worsequery"]?1:2);
    int rankB=[b isEqual:expected]?0:([b hasSuffix:@".worsequery"]?1:2);
    return rankA==rankB?[a compare:b]:(rankA<rankB?NSOrderedAscending:NSOrderedDescending);
  }];
  result[@"entitlements_readable"]=@(effective!=nil);
  result[@"effective_group_count"]=@(effective.count);
  result[@"dedicated_group_authorized"]=@([valid containsObject:expected]);
  result[@"selected_group_authorized"]=@NO;result[@"container_available"]=@NO;
  result[@"payload_storage"]=@"private_container_only";
  NSMutableArray *checks=[NSMutableArray array];result[@"candidates"]=checks;
  NSString *selected=nil;NSUInteger accessible=0,eligible=0,index=0;
  for(NSString *group in groups) {
    NSURL *root=[self groupRoot:group];
    NSMutableDictionary *check=[@{@"index":@(index),@"container_available":@(root!=nil)} mutableCopy];
    [checks addObject:check];
    if(root) {
      accessible++;errno=0;int status=crumb_cache_inspect(root.fileSystemRepresentation);
      check[@"cache_state"]=status==0?@"empty_or_missing":@"ineligible";check[@"errno"]=@(status?errno:0);
      if(!status) {
        eligible++;
        if(!selected) {
          selected=group;*selectedRoot=root;result[@"selected_index"]=@(index);
          result[@"selection"]=[group isEqual:expected]?@"dedicated_expected":([group hasSuffix:@".worsequery"]?@"dedicated_alias":@"profile_granted");
        }
      }
    }
    index++;
  }
  result[@"accessible_group_count"]=@(accessible);result[@"eligible_group_count"]=@(eligible);
  result[@"selected_group_authorized"]=@(selected!=nil);result[@"container_available"]=@(selected!=nil);
  return selected;
}
- (int64_t)grant:(NSURL *)target group:(NSString *)group result:(NSMutableDictionary *)result {
  if(atomic_load(&_cancelled))return -256;
  NSFileManager *fm=NSFileManager.defaultManager;
  if(![[self effectiveGroups] containsObject:group]){result[@"code"]=@"group_not_in_effective_signature";return -1;}
  NSURL *root=[self groupRoot:group];
  NSURL *ledger=[self.directory URLByAppendingPathComponent:@"group-lease.json"];
  if([fm fileExistsAtPath:ledger.path] || [fm fileExistsAtPath:[self.directory URLByAppendingPathComponent:@"group-link.json"].path]){result[@"code"]=@"previous_grant_cleanup_pending";return -256;}
  if(!root){result[@"code"]=@"group_container_unavailable";return -1;}
  NSURL *links=[self.directory URLByAppendingPathComponent:@"cache-links"];
  if(![fm createDirectoryAtURL:links withIntermediateDirectories:YES attributes:@{NSFilePosixPermissions:@0700} error:nil]){result[@"code"]=@"private_lease_staging_failed";return -256;}
  NSString *identifier=NSUUID.UUID.UUIDString;
  NSURL *stage=[links URLByAppendingPathComponent:identifier];
  CrumbCacheLease lease;
  if(crumb_cache_begin(&lease,root.fileSystemRepresentation,target.fileSystemRepresentation,stage.fileSystemRepresentation,identifier.UTF8String)) {
    result[@"code"]=@"cache_lease_unavailable";result[@"lease_errno"]=@(errno);return -256;
  }
  int64_t handle=-256;
  @try {
    NSDictionary *state=@{@"schema":@1,@"group":group,@"identifier":identifier,@"target":target.path,
      @"had_caches":@(lease.had_caches),@"original_device":@(lease.original_device),@"original_inode":@(lease.original_inode),
      @"link_device":@(lease.link_device),@"link_inode":@(lease.link_inode)};
    if(![self saveJSON:state to:ledger]){result[@"code"]=@"grant_ledger_failed";return -256;}
    [self checkpoint:@"sandbox_extension" metrics:nil];
    if(atomic_load(&_cancelled))return -256;
    if(crumb_cache_activate(&lease)){result[@"code"]=@"cache_lease_activation_failed";result[@"lease_errno"]=@(errno);return -256;}
    BOOL alt=NSProcessInfo.processInfo.operatingSystemVersion.majorVersion==26 && [target.path containsString:@"/Containers/Shared/AppGroup/"];
    struct stat targetStatus;errno=0;int targetStat=lstat(target.fileSystemRepresentation,&targetStatus);
    result[@"target_stat_result"]=@(targetStat);result[@"target_stat_errno"]=@(targetStat?errno:0);
    handle=worse_query_prepared(target.path,group,[[root URLByAppendingPathComponent:@"Library/Caches"] path],alt);result[@"result"]=@(handle);
    if(handle>=0){atomic_store(&processHasGrant,true);result[@"code"]=@"extension_consumed";}
    else result[@"code"]=@"extension_rejected";
  }@finally {
    int restored=crumb_cache_restore(&lease);result[@"caches_restored"]=@(restored==0);result[@"restore_errno"]=@(restored?errno:0);
    if(!restored)[fm removeItemAtURL:ledger error:nil];
    else{atomic_store(&_cancelled,true);result[@"code"]=@"group_cleanup_pending";}
  }
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
  self.report=[@{@"schema":@2,@"run_id":job,@"outcome":@"running",@"attempts":[NSMutableArray array],
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
    NSArray *groups=[self effectiveGroups];NSURL *groupRoot=nil;
    NSMutableDictionary *groupDiagnostics=[NSMutableDictionary dictionary];
    NSString *group=[self selectGroup:groups diagnostics:groupDiagnostics root:&groupRoot];
    self.report[@"app_group"]=groupDiagnostics;
    NSMutableDictionary *discovery=[NSMutableDictionary dictionary];
    NSURL *installed=[self installedBundle:bundleID result:discovery];self.report[@"installed_discovery"]=discovery;
    NSURL *input=[work URLByAppendingPathComponent:@"archive/input"];
    BOOL contaminated=atomic_load(&processHasGrant);
    if(!contaminated) {
      if([self attempt:@"downloadedIPA" access:@"baseline" source:input manifest:manifest work:work group:group])return self.report;
      if(installed && [self attempt:@"installedApp" access:@"baseline" source:installed manifest:manifest work:work group:group])return self.report;
    }else self.report[@"baseline_skipped"]=@"restart_crumb_for_clean_baseline";
    if(groupRoot && !atomic_load(&_cancelled)) {
      NSURL *control=[work URLByAppendingPathComponent:@"control"];
      [fm createDirectoryAtURL:control withIntermediateDirectories:NO attributes:@{NSFilePosixPermissions:@0700} error:nil];
      NSURL *canary=[control URLByAppendingPathComponent:@"canary"];NSData *bytes=[job dataUsingEncoding:NSUTF8StringEncoding];
      [bytes writeToURL:canary options:NSDataWritingAtomic error:nil];
      NSMutableDictionary *canaryReport=[NSMutableDictionary dictionary];canaryReport[@"before"]=[self readCanary:canary expected:bytes];
      int64_t handle=[self grant:control group:group result:canaryReport];
      canaryReport[@"after"]=[self readCanary:canary expected:bytes];if(handle>=0)worse_query_release(handle);
      canaryReport[@"note"]=@"own_container_control_does_not_demonstrate_new_cross_app_access";self.report[@"canary"]=canaryReport;
      if(!atomic_load(&_cancelled)) {
        if([self attempt:@"downloadedIPA" access:@"worseQuery" source:input manifest:manifest work:work group:group])return self.report;
        if(installed && [self attempt:@"installedApp" access:@"worseQuery" source:installed manifest:manifest work:work group:group])return self.report;
      }
      self.report[@"shared_container"]=@"skipped_private_payload_storage_only";
    }
    if(!self.outputURL)self.report[@"outcome"]=atomic_load(&_cancelled)?@"cancelled_or_local_failure":@"no_verified_route";
  }@catch(NSException *e){(void)e;self.report[@"outcome"]=@"local_exception";}
  @finally {
    if(atomic_load(&_cancelled)) {
      if(self.outputURL)[fm removeItemAtURL:self.outputURL error:nil];
      self.outputURL=nil;self.report[@"outcome"]=@"cancelled_or_local_failure";
    }
    self.report[@"workspace_cleanup_complete"]=@(remove_owned(work));
    [self checkpoint:@"finished" metrics:nil];
    [fm removeItemAtURL:active error:nil];
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
