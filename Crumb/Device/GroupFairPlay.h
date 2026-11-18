#import <Foundation/Foundation.h>
void crumb_group_reset(void);
void crumb_group_cancel(void);
BOOL crumb_group_decrypt(NSURL *source, NSURL *destination, BOOL *encrypted, NSError **error);
