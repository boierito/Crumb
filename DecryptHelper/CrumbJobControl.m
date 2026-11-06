#import "CrumbJobControl.h"
BOOL crumb_matches_cancel_url(NSURL *url, NSString *jobID) {
  if (![jobID isKindOfClass:NSString.class] || ![url.scheme isEqualToString:@"crumb-decrypt-helper"] ||
      ![url.host isEqualToString:@"cancel"] || ![[NSUUID alloc] initWithUUIDString:jobID]) return NO;
  NSArray *items = [NSURLComponents componentsWithURL:url resolvingAgainstBaseURL:NO].queryItems;
  NSString *candidate = nil;
  NSUInteger count = 0;
  for (NSURLQueryItem *item in items) if ([item.name isEqualToString:@"job"]) {
    if (++count != 1) return NO;
    candidate = item.value;
  }
  return candidate && [candidate isEqualToString:jobID];
}
BOOL crumb_store_cancellation(NSURL *documents, NSString *jobID, NSError **error) {
  if (![jobID isKindOfClass:NSString.class] || ![[NSUUID alloc] initWithUUIDString:jobID]) return NO;
  return [jobID writeToURL:[documents URLByAppendingPathComponent:@"cancelled-job"] atomically:YES encoding:NSUTF8StringEncoding error:error];
}
BOOL crumb_job_was_cancelled(NSURL *documents, NSString *jobID) {
  if (![jobID isKindOfClass:NSString.class]) return NO;
  NSURL *file = [documents URLByAppendingPathComponent:@"cancelled-job"];
  NSDictionary *values = [file resourceValuesForKeys:@[NSURLFileSizeKey, NSURLIsSymbolicLinkKey] error:nil];
  if (!values || [values[NSURLFileSizeKey] unsignedLongLongValue] > 64 || [values[NSURLIsSymbolicLinkKey] boolValue]) return NO;
  NSString *stored = [NSString stringWithContentsOfURL:file encoding:NSUTF8StringEncoding error:nil];
  return stored && [stored isEqualToString:jobID];
}
