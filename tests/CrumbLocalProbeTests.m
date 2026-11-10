#import "../Crumb/LocalProbe/CrumbLocalProbe.h"
#import <assert.h>
#import <stdatomic.h>
#import <string.h>
#import "../SharedDecrypt/CrumbCacheLease.h"
#import <sys/stat.h>
#import <unistd.h>
@interface CrumbLocalProbe (GroupTesting)
- (NSString *)selectGroup:(NSArray *)groups diagnostics:(NSMutableDictionary *)report root:(NSURL *__autoreleasing *)root;
- (void)recover;
@end
@interface GroupTestProbe : CrumbLocalProbe
@property(nonatomic,strong) NSArray *testGroups;
@property(nonatomic,strong) NSDictionary *testRoots;
@property(nonatomic) NSUInteger resolutionCalls;
@end
@implementation GroupTestProbe
- (NSArray *)effectiveGroups { return self.testGroups; }
- (NSURL *)groupRoot:(NSString *)group { self.resolutionCalls++; return self.testRoots[group]; }
@end
static NSURL *make_group(NSURL *root,NSString *name) {
  NSURL *group=[root URLByAppendingPathComponent:name];
  assert([NSFileManager.defaultManager createDirectoryAtURL:[group URLByAppendingPathComponent:@"Library/Caches"] withIntermediateDirectories:YES attributes:nil error:nil]);
  return group;
}
static void test_group_selection_and_recovery(NSURL *root) {
  NSFileManager *fm=NSFileManager.defaultManager;
  NSURL *first=make_group(root,@"group-first"),*second=make_group(root,@"group-second");
  NSURL *unrelated=[first URLByAppendingPathComponent:@"Library/Caches/other-app-data"];
  assert([[@"preserve" dataUsingEncoding:NSUTF8StringEncoding] writeToURL:unrelated atomically:YES]);
  NSURL *privateRoot=[root URLByAppendingPathComponent:@"policy-runner"];
  GroupTestProbe *probe=[[GroupTestProbe alloc]initWithDirectory:privateRoot];
  probe.testGroups=@[@"group.test.2",@"group.test.1"];
  probe.testRoots=@{@"group.test.1":first,@"group.test.2":second};
  NSMutableDictionary *report=[NSMutableDictionary dictionary];NSURL *selectedRoot=nil;
  NSString *selected=[probe selectGroup:probe.testGroups diagnostics:report root:&selectedRoot];
  assert([selected isEqual:@"group.test.2"] && [selectedRoot isEqual:second]);
  assert([report[@"selected_group_authorized"] boolValue] && ![report[@"dedicated_group_authorized"] boolValue]);
  assert([report[@"effective_group_count"] intValue]==2 && [report[@"eligible_group_count"] intValue]==1);
  assert([[NSString stringWithContentsOfURL:unrelated encoding:NSUTF8StringEncoding error:nil] isEqual:@"preserve"]);
  NSString *json=[[NSString alloc]initWithData:[NSJSONSerialization dataWithJSONObject:report options:0 error:nil] encoding:NSUTF8StringEncoding];
  assert(![json containsString:@"group.test"] && ![json containsString:root.path]);
  // Only the effective signature is consulted, never an accessible ungranted guess.
  probe.testGroups=@[@"group.test.absent"];probe.resolutionCalls=0;report=[NSMutableDictionary dictionary];selectedRoot=nil;
  assert([probe selectGroup:probe.testGroups diagnostics:report root:&selectedRoot]==nil && probe.resolutionCalls==1);
  probe.testGroups=@[@"group.*",@1,@"../group.test.2"];probe.resolutionCalls=0;
  assert([probe selectGroup:probe.testGroups diagnostics:[NSMutableDictionary dictionary] root:&selectedRoot]==nil && probe.resolutionCalls==0);
  // A signer-renamed dedicated group takes priority over a generic allowed one.
  probe.testGroups=@[@"group.test.generic",@"group.test.original.worsequery"];
  probe.testRoots=@{@"group.test.generic":second,@"group.test.original.worsequery":second};
  report=[NSMutableDictionary dictionary];
  assert([[probe selectGroup:probe.testGroups diagnostics:report root:&selectedRoot] isEqual:@"group.test.original.worsequery"]);
  assert([report[@"selection"] isEqual:@"dedicated_alias"]);
  // Rehydrate a shared-profile lease from the private journal after process death.
  probe.testGroups=@[@"group.test.2"];probe.testRoots=@{@"group.test.2":second};
  NSString *identifier=NSUUID.UUID.UUIDString;
  NSURL *links=[privateRoot URLByAppendingPathComponent:@"cache-links"];
  assert([fm createDirectoryAtURL:links withIntermediateDirectories:YES attributes:nil error:nil]);
  NSURL *stage=[links URLByAppendingPathComponent:identifier];
  CrumbCacheLease lease;
  assert(crumb_cache_begin(&lease,second.fileSystemRepresentation,privateRoot.fileSystemRepresentation,stage.fileSystemRepresentation,identifier.UTF8String)==0);
  NSDictionary *ledger=@{@"schema":@1,@"group":@"group.test.2",@"identifier":identifier,@"target":privateRoot.path,
    @"had_caches":@(lease.had_caches),@"original_device":@(lease.original_device),@"original_inode":@(lease.original_inode),@"link_device":@(lease.link_device),@"link_inode":@(lease.link_inode)};
  NSURL *ledgerURL=[privateRoot URLByAppendingPathComponent:@"group-lease.json"];
  assert([[NSJSONSerialization dataWithJSONObject:ledger options:0 error:nil] writeToURL:ledgerURL atomically:YES]);
  assert(crumb_cache_activate(&lease)==0);close(lease.lock_fd);
  probe.testGroups=@[];[probe recover];assert([fm fileExistsAtPath:ledgerURL.path]);
  probe.testGroups=@[@"group.test.2"];[probe recover];
  struct stat status;assert(lstat([second URLByAppendingPathComponent:@"Library/Caches"].fileSystemRepresentation,&status)==0 && S_ISDIR(status.st_mode));
  assert((uint64_t)status.st_ino==lease.original_inode && ![fm fileExistsAtPath:ledgerURL.path]);
}
// Archive operations are independently tested with real ZIPs in Rust. Here a
// boundary stub exercises interruption recovery, cancellation and redaction.
char *crumb_probe_archive(const char *json,void *ctx,bool (*cancel)(void *)) {
  (void)json;
  return strdup(cancel(ctx)?"{\"ok\":false,\"code\":\"cancelled\"}":"{\"ok\":false,\"code\":\"archive_operation_failed\"}");
}
void crumb_device_string_free(char *s){free(s);}
static void write_json(NSURL *url,NSDictionary *json){
  assert([[NSJSONSerialization dataWithJSONObject:json options:0 error:nil] writeToURL:url atomically:YES]);
}
int main(void){@autoreleasepool {
  NSFileManager *fm=NSFileManager.defaultManager;
  NSURL *root=[[NSURL fileURLWithPath:NSTemporaryDirectory()] URLByAppendingPathComponent:NSUUID.UUID.UUIDString];
  assert([fm createDirectoryAtURL:root withIntermediateDirectories:YES attributes:nil error:nil]);
  NSString *job=NSUUID.UUID.UUIDString;
  NSURL *work=[[root URLByAppendingPathComponent:@"runs"] URLByAppendingPathComponent:job];
  assert([fm createDirectoryAtURL:work withIntermediateDirectories:YES attributes:nil error:nil]);
  NSURL *unrelated=[[root URLByAppendingPathComponent:@"runs"] URLByAppendingPathComponent:@"unrelated"];
  assert([fm createDirectoryAtURL:unrelated withIntermediateDirectories:YES attributes:nil error:nil]);
  NSURL *report=[root URLByAppendingPathComponent:[NSString stringWithFormat:@"report-%@.json",job]];
  write_json(report,@{@"schema":@1,@"run_id":job,@"current_stage":@"remap",@"outcome":@"running"});
  write_json([root URLByAppendingPathComponent:@"active.json"],@{@"run_id":job});
  CrumbLocalProbe *probe=[[CrumbLocalProbe alloc]initWithDirectory:root];
  NSDictionary *recovered=[NSJSONSerialization JSONObjectWithData:[NSData dataWithContentsOfURL:report] options:0 error:nil];
  assert([recovered[@"outcome"] isEqual:@"interrupted"] && [recovered[@"current_stage"] isEqual:@"remap"]);
  assert(![fm fileExistsAtPath:work.path] && [fm fileExistsAtPath:unrelated.path]);
  NSString *privateName=@"do-not-export-license-or-path";
  NSDictionary *result=[probe runIPA:[root URLByAppendingPathComponent:privateName] bundleID:@"test.private" progress:^(NSString *stage){(void)stage;}];
  assert([result[@"outcome"] isEqual:@"preparation_failed"] && probe.outputURL==nil);
  NSString *json=[[NSString alloc]initWithData:[NSJSONSerialization dataWithJSONObject:result options:0 error:nil] encoding:NSUTF8StringEncoding];
  assert(![json containsString:root.path] && ![json containsString:privateName] && ![json containsString:@"test.private"]);
  assert([fm fileExistsAtPath:probe.reportURL.path]);
  result=[probe runIPA:[root URLByAppendingPathComponent:privateName] bundleID:@"test.private" progress:^(NSString *stage){if([stage isEqual:@"archive_preparation"])[probe cancel];}];
  assert([result[@"outcome"] isEqual:@"cancelled_or_local_failure"] && probe.outputURL==nil);
  assert(![fm fileExistsAtPath:[root URLByAppendingPathComponent:@"active.json"].path]);
  test_group_selection_and_recovery(root);
  [fm removeItemAtURL:root error:nil];
  puts("local coordinator: interruption recovery, isolated cleanup, failure redaction and cooperative cancellation passed (archive boundary stub)");
}}
