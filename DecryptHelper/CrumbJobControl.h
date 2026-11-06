#import <Foundation/Foundation.h>
BOOL crumb_matches_cancel_url(NSURL *url, NSString *jobID);
BOOL crumb_store_cancellation(NSURL *documents, NSString *jobID, NSError **error);
BOOL crumb_job_was_cancelled(NSURL *documents, NSString *jobID);
