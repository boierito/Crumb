#import <Foundation/Foundation.h>
#include <stdio.h>
#import "../DecryptHelper/CrumbJobControl.h"
int main(void) { @autoreleasepool {
  NSString *a = @"00000000-0000-4000-8000-000000000001";
  NSString *b = @"00000000-0000-4000-8000-000000000002";
  NSURL *valid = [NSURL URLWithString:[@"crumb-decrypt-helper://cancel?job=" stringByAppendingString:a]];
  NSCAssert(crumb_matches_cancel_url(valid,a), @"Correct helper job must cancel");
  NSCAssert(!crumb_matches_cancel_url(valid,b), @"Old job must not cancel a new one");
  NSCAssert(!crumb_matches_cancel_url([NSURL URLWithString:[@"crumb-decrypt-helper://run?job=" stringByAppendingString:a]],a), @"Wrong command");
  NSCAssert(!crumb_matches_cancel_url([NSURL URLWithString:[@"https://cancel?job=" stringByAppendingString:a]],a), @"Wrong scheme");
  NSCAssert(!crumb_matches_cancel_url([NSURL URLWithString:[NSString stringWithFormat:@"crumb-decrypt-helper://cancel?job&job=%@",a]],a), @"Duplicate job fields");
  NSCAssert(!crumb_matches_cancel_url(valid,nil), @"Missing job");
  NSURL *root = [NSURL fileURLWithPath:[NSTemporaryDirectory() stringByAppendingPathComponent:NSUUID.UUID.UUIDString] isDirectory:YES];
  NSFileManager *fm = NSFileManager.defaultManager;
  NSCAssert([fm createDirectoryAtURL:root withIntermediateDirectories:YES attributes:nil error:nil], @"Fixture directory");
  NSCAssert(!crumb_job_was_cancelled(root,a), @"Fresh job");
  NSCAssert(crumb_store_cancellation(root,a,nil), @"Persistent cancellation");
  NSCAssert(crumb_job_was_cancelled(root,a), @"Restart must not resume cancelled work");
  NSCAssert(!crumb_job_was_cancelled(root,b), @"New jobs remain usable");
  NSCAssert(!crumb_store_cancellation(root,@"not-a-uuid",nil), @"Malformed identity");
  NSURL *marker = [root URLByAppendingPathComponent:@"cancelled-job"];
  [[NSMutableData dataWithLength:1024] writeToURL:marker atomically:YES];
  NSCAssert(!crumb_job_was_cancelled(root,a), @"Bounded token file");
  [fm removeItemAtURL:root error:nil];
  puts("Helper cancellation identity, persistence and restart regressions passed");
} return 0; }
