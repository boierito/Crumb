#import <Foundation/Foundation.h>
NS_ASSUME_NONNULL_BEGIN
@interface CrumbInstalledDiscovery : NSObject
@property(nonatomic,readonly) NSMutableDictionary *report;
@property(nonatomic,readonly) int64_t parentHandle;
- (instancetype)initWithBundleID:(NSString *)bundleID appName:(NSString *)appName
                            work:(NSURL *)work
                           grant:(int64_t (^)(NSURL *,NSMutableDictionary *))grant
                         release:(void (^)(int64_t))release
                       cancelled:(BOOL (^)(void))cancelled
                        progress:(void (^)(NSString *))progress;
- (NSArray<NSURL *> *)run;
- (void)close;
@end
NS_ASSUME_NONNULL_END
