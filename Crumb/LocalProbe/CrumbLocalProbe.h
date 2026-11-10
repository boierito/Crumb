#import <Foundation/Foundation.h>
NS_ASSUME_NONNULL_BEGIN
@interface CrumbLocalProbe : NSObject
- (instancetype)initWithDirectory:(NSURL *)directory;
// Must be called on a background serial queue. Returns only a sanitized report.
- (NSDictionary *)runIPA:(NSURL *)ipa bundleID:(NSString *)bundleID
               progress:(void (^)(NSString *stage))progress;
@property(nonatomic, readonly, nullable) NSURL *outputURL;
@property(nonatomic, readonly, nullable) NSURL *reportURL;
- (void)cancel;
@end
NS_ASSUME_NONNULL_END
