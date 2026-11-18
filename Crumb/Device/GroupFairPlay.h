#import <Foundation/Foundation.h>
void crumb_group_reset(void);
void crumb_group_cancel(void);
BOOL crumb_group_decrypt(NSURL *source, NSURL *destination, BOOL *encrypted, NSError **error);
// Effective entitlements of this running, re-signed process, not just its profile.
NSArray<NSString *> * _Nullable crumb_group_signed_groups(void);
