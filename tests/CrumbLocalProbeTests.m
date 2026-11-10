#import "../Crumb/LocalProbe/CrumbLocalProbe.h"
#import <assert.h>
#import <stdatomic.h>
#import <string.h>
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
  [fm removeItemAtURL:root error:nil];
  puts("local coordinator: interruption recovery, isolated cleanup, failure redaction and cooperative cancellation passed (archive boundary stub)");
}}
