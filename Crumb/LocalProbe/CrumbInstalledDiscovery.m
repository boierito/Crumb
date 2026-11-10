#import "CrumbInstalledDiscovery.h"
#import "../../SharedDecrypt/CrumbLSStore.h"
#import <dlfcn.h>
#import <dirent.h>
#import <fcntl.h>
#import <sys/stat.h>
#import <sys/mount.h>
#import <unistd.h>
static NSString *const bundleRoot=@"/private/var/containers/Bundle/Application";
static NSString *const daemonRoot=@"/private/var/mobile/Containers/Data/InternalDaemon";
static NSString *const storeHintKey=@"CrumbInstalledLSStoreHint";
@interface CrumbInstalledDiscovery ()
@property(nonatomic,readwrite) NSMutableDictionary *report;
@property(nonatomic,readwrite) int64_t parentHandle;
@property(nonatomic) NSString *bundleID;
@property(nonatomic) NSString *appName;
@property(nonatomic) NSURL *work;
@property(nonatomic,copy) int64_t (^grant)(NSURL *,NSMutableDictionary *);
@property(nonatomic,copy) void (^releaseGrant)(int64_t);
@property(nonatomic,copy) BOOL (^cancelled)(void);
@property(nonatomic,copy) void (^progress)(NSString *);
@property(nonatomic) NSMutableArray<NSURL *> *candidates;
@property(nonatomic) NSMutableSet<NSString *> *seenStores;
@property(nonatomic) NSTimeInterval deadline;
@end
@implementation CrumbInstalledDiscovery
- (instancetype)initWithBundleID:(NSString *)bundleID appName:(NSString *)appName work:(NSURL *)work
                           grant:(int64_t (^)(NSURL *,NSMutableDictionary *))grant release:(void (^)(int64_t))release
                       cancelled:(BOOL (^)(void))cancelled progress:(void (^)(NSString *))progress {
  if((self=[super init])) {
    _bundleID=bundleID;_appName=appName;_work=work;_grant=grant;_releaseGrant=release;_cancelled=cancelled;_progress=progress;
    _parentHandle=-1;_candidates=[NSMutableArray array];_seenStores=[NSMutableSet set];
    _report=[@{@"candidate_count":@0,@"metadata_verified_count":@0,@"stores":[NSMutableArray array],
      @"discovery_does_not_grant_bundle_access":@YES} mutableCopy];
  }return self;
}
- (BOOL)stopped {return self.cancelled() || NSProcessInfo.processInfo.systemUptime>=self.deadline;}
- (void)close {if(self.parentHandle>=0){self.releaseGrant(self.parentHandle);self.parentHandle=-1;}}
- (void)dealloc {[self close];}
- (NSData *)readFile:(NSURL *)url limit:(size_t)limit error:(int *)error {
  *error=0;int fd=open(url.fileSystemRepresentation,O_RDONLY|O_NOFOLLOW|O_CLOEXEC);
  if(fd<0){*error=errno;return nil;}
  struct stat before,after;
  if(fstat(fd,&before)){*error=errno;close(fd);return nil;}
  if(!S_ISREG(before.st_mode)||before.st_size<=0||(uint64_t)before.st_size>limit){*error=EFBIG;close(fd);return nil;}
  NSMutableData *data=[NSMutableData dataWithLength:(NSUInteger)before.st_size];size_t got=0;
  while(got<data.length && ![self stopped]) {
    ssize_t n=read(fd,(uint8_t *)data.mutableBytes+got,MIN((size_t)65536,data.length-got));
    if(n<0&&errno==EINTR)continue;
    if(n<=0){*error=n<0?errno:EIO;break;}got+=(size_t)n;
  }
  if(got!=data.length && !*error)*error=ECANCELED;
  if(fstat(fd,&after)){if(!*error)*error=errno;}
  else if(before.st_ino!=after.st_ino||before.st_size!=after.st_size||before.st_mtimespec.tv_sec!=after.st_mtimespec.tv_sec||before.st_mtimespec.tv_nsec!=after.st_mtimespec.tv_nsec){*error=ESTALE;}
  close(fd);return *error?nil:data;
}
- (BOOL)validBundlePath:(NSString *)path {
  NSString *prefix=[bundleRoot stringByAppendingString:@"/"];
  if(![path hasPrefix:prefix])return NO;
  NSArray *parts=[[path substringFromIndex:prefix.length] componentsSeparatedByString:@"/"];
  return parts.count==2 && [[NSUUID alloc]initWithUUIDString:parts[0]] &&
         [parts[1] hasSuffix:@".app"] && [parts[1] length]>4 && ![parts[1] containsString:@"\\"];
}
- (void)addCandidate:(NSString *)path source:(NSString *)source {
  if([self stopped]||self.candidates.count>=CRUMB_LS_MAX_PATHS)return;
  if([path hasPrefix:@"/var/"])path=[@"/private" stringByAppendingString:path];
  while([path hasSuffix:@"/"])path=[path substringToIndex:path.length-1];
  if(![self validBundlePath:path])return;
  NSURL *url=[NSURL fileURLWithPath:path isDirectory:YES];if([self.candidates containsObject:url])return;
  int error=0;NSData *metadata=[self readFile:[url URLByAppendingPathComponent:@"Info.plist"] limit:1024*1024 error:&error];
  NSDictionary *info=metadata?[NSPropertyListSerialization propertyListWithData:metadata options:0 format:nil error:nil]:nil;
  if(info && ![info isKindOfClass:NSDictionary.class])info=nil;
  if(info && ![info[@"CFBundleIdentifier"] isEqual:self.bundleID]) {
    self.report[@"identity_mismatch_count"]=@([self.report[@"identity_mismatch_count"] unsignedIntegerValue]+1);return;
  }
  // Unreadable metadata is retained as an unverified candidate, so the caller
  // can test a target-scoped extension before applying its strict identity gate.
  [self.candidates addObject:url];
  self.report[@"candidate_count"]=@(self.candidates.count);
  NSMutableArray *checks=self.report[@"candidate_checks"];
  if(!checks){checks=[NSMutableArray array];self.report[@"candidate_checks"]=checks;}
  [checks addObject:@{@"index":@(self.candidates.count-1),@"source":source,@"metadata_readable":@(info!=nil),@"metadata_errno":@(error),@"identity_verified":@(info!=nil)}];
  if(info)self.report[@"metadata_verified_count"]=@([self.report[@"metadata_verified_count"] unsignedIntegerValue]+1);
}
- (NSDictionary *)listBundleParent {
  errno=0;DIR *dir=opendir(bundleRoot.fileSystemRepresentation);int error=dir?0:errno;
  NSMutableDictionary *result=[@{@"opened":@(dir!=NULL),@"errno":@(error),@"entries_checked":@0} mutableCopy];
  if(!dir)return result;
  NSUInteger count=0;struct dirent *entry;
  while(![self stopped] && count<4096 && (entry=readdir(dir))) {
    NSString *name=[NSFileManager.defaultManager stringWithFileSystemRepresentation:entry->d_name length:strlen(entry->d_name)];
    if(![[NSUUID alloc]initWithUUIDString:name])continue;count++;
    NSString *path=[[bundleRoot stringByAppendingPathComponent:name] stringByAppendingPathComponent:self.appName];
    struct stat st;errno=0;
    if(lstat(path.fileSystemRepresentation,&st)==0 && S_ISDIR(st.st_mode))[self addCandidate:path source:@"bundle_directory"];
  }
  closedir(dir);result[@"entries_checked"]=@(count);result[@"budget_exhausted"]=@(count==4096||[self stopped]);return result;
}
- (BOOL)isStorePath:(NSString *)path container:(NSString *__autoreleasing *)container {
  NSString *prefix=[daemonRoot stringByAppendingString:@"/"];
  if(![path hasPrefix:prefix])return NO;
  NSArray *parts=[[path substringFromIndex:prefix.length] componentsSeparatedByString:@"/"];
  if(parts.count!=4 || ![[NSUUID alloc]initWithUUIDString:parts[0]] || ![parts[1] isEqual:@"Library"] || ![parts[2] isEqual:@"Caches"] ||
     ![parts[3] hasPrefix:@"com.apple.LaunchServices-"] || ![parts[3] hasSuffix:@".csstore"])return NO;
  if(container)*container=[daemonRoot stringByAppendingPathComponent:parts[0]];return YES;
}
- (void)readStore:(NSString *)path {
  NSString *container=nil;
  if([self stopped]||self.seenStores.count>=8||[self.seenStores containsObject:path]||![self isStorePath:path container:&container])return;
  [self.seenStores addObject:path];self.progress(@"installed_registry_read");
  NSMutableDictionary *result=[@{@"index":@(self.seenStores.count-1),@"layout":@"CSStore2_27.0_candidate"} mutableCopy];
  [self.report[@"stores"] addObject:result];
  NSMutableDictionary *grant=[NSMutableDictionary dictionary];result[@"grant"]=grant;
  int64_t handle=self.grant([NSURL fileURLWithPath:container isDirectory:YES],grant);
  @try {
    int error=0;NSData *data=[self readFile:[NSURL fileURLWithPath:path] limit:64*1024*1024 error:&error];
    result[@"read_errno"]=@(error);result[@"readable"]=@(data!=nil);
    if(!data)return;
    char paths[CRUMB_LS_MAX_PATHS][PATH_MAX];CrumbLSMetrics metrics;
    int count=crumb_ls_bundle_paths(data.bytes,data.length,self.bundleID.UTF8String,paths,&metrics);
    result[@"parser_result"]=@(count);result[@"parser_errno"]=@(metrics.error);result[@"matching_records"]=@(metrics.target_records);
    if(count<0)return;
    [NSUserDefaults.standardUserDefaults setObject:path forKey:storeHintKey];
    for(int i=0;i<count;i++) {
      NSString *candidate=[[NSString alloc]initWithUTF8String:paths[i]];
      if(candidate)[self addCandidate:candidate source:@"launchservices_store"];
    }
  }@finally {if(handle>=0)self.releaseGrant(handle);}
}
- (void)processInodePath:(NSString *)raw {
  NSString *path=[raw hasPrefix:@"/var/"]?[@"/private" stringByAppendingString:raw]:raw;
  if([path hasPrefix:[daemonRoot stringByAppendingString:@"/"]]){[self readStore:path];return;}
  NSString *prefix=[bundleRoot stringByAppendingString:@"/"];
  if(![path hasPrefix:prefix])return;
  NSArray *parts=[[path substringFromIndex:prefix.length] componentsSeparatedByString:@"/"];
  if(parts.count>=2 && [parts[1] isEqual:self.appName] && [[NSUUID alloc]initWithUUIDString:parts[0]]) {
    [self addCandidate:[[bundleRoot stringByAppendingPathComponent:parts[0]] stringByAppendingPathComponent:parts[1]] source:@"inode_path"];
  }
}
- (void)scanInodes {
  self.progress(@"installed_inode_scan");
  NSMutableDictionary *scan=[@{@"max_probes":@3000000,@"max_seconds":@20,@"probes":@0,@"paths_seen":@0} mutableCopy];self.report[@"scan"]=scan;
  typedef int (*GetPath)(void *,size_t,fsid_t *,uint64_t);
  GetPath getpath=(GetPath)dlsym(RTLD_DEFAULT,"fsgetpath");
  if(!getpath){scan[@"code"]=@"fsgetpath_unavailable";return;}
  NSURL *marker=[self.work URLByAppendingPathComponent:@"inode-anchor"];
  if(![[NSData data] writeToURL:marker options:NSDataWritingAtomic error:nil]){scan[@"code"]=@"anchor_failed";return;}
  struct stat st;struct statfs volume;
  if(stat(marker.fileSystemRepresentation,&st)||statfs(marker.fileSystemRepresentation,&volume)){scan[@"code"]=@"volume_lookup_failed";scan[@"errno"]=@(errno);unlink(marker.fileSystemRepresentation);return;}
  fsid_t fsid=volume.f_fsid;uint64_t probes=0,seen=0;NSTimeInterval started=NSProcessInfo.processInfo.systemUptime;
  uint64_t inode=st.st_ino;char buffer[PATH_MAX];
  while(inode>0 && probes<3000000 && ![self stopped] && ![self.report[@"metadata_verified_count"] unsignedIntegerValue]) {
    int n=getpath(buffer,sizeof(buffer),&fsid,inode--);probes++;
    if(n<=0||n>=(int)sizeof(buffer))continue;
    buffer[n]=0;size_t len=strnlen(buffer,sizeof(buffer));if(len>=sizeof(buffer))continue;seen++;
    NSString *path=[NSFileManager.defaultManager stringWithFileSystemRepresentation:buffer length:len];
    if(path)[self processInodePath:path];
  }
  unlink(marker.fileSystemRepresentation);scan[@"probes"]=@(probes);scan[@"paths_seen"]=@(seen);
  scan[@"elapsed_ms"]=@((NSProcessInfo.processInfo.systemUptime-started)*1000);
  scan[@"code"]=[self.report[@"metadata_verified_count"] unsignedIntegerValue]?@"verified_candidate_found":self.cancelled()?@"cancelled":@"bounded_scan_incomplete";
  scan[@"absence_not_proven"]=@YES;
}
- (NSArray<NSURL *> *)run {
  self.deadline=NSProcessInfo.processInfo.systemUptime+20;
  self.progress(@"installed_bundle_access");self.report[@"parent_before"]=[self listBundleParent];
  if(![self.report[@"metadata_verified_count"] unsignedIntegerValue] && ![self stopped]) {
    NSMutableDictionary *grant=[NSMutableDictionary dictionary];self.report[@"parent_grant"]=grant;
    self.parentHandle=self.grant([NSURL fileURLWithPath:bundleRoot isDirectory:YES],grant);
    self.report[@"parent_after"]=[self listBundleParent];
  }
  if(![self.report[@"metadata_verified_count"] unsignedIntegerValue] && ![self stopped]) {
    id hint=[NSUserDefaults.standardUserDefaults objectForKey:storeHintKey];
    if([hint isKindOfClass:NSString.class])[self readStore:hint];
  }
  if(![self.report[@"metadata_verified_count"] unsignedIntegerValue] && ![self stopped])[self scanInodes];
  self.report[@"outcome"]=self.candidates.count?@"candidates_found":self.cancelled()?@"cancelled":@"not_located_within_probe";
  self.report[@"absence_not_proven"]=@YES;return [self.candidates copy];
}
@end
