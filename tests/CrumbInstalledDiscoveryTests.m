#import "../Crumb/LocalProbe/CrumbInstalledDiscovery.h"
#import <assert.h>
#import <errno.h>
#import <string.h>
@interface CrumbInstalledDiscovery (Testing)
- (NSDictionary *)listBundleParent;
- (void)scanInodes;
- (void)processInodePath:(NSString *)path;
- (NSData *)readFile:(NSURL *)url limit:(size_t)limit error:(int *)error;
@end
@interface DiscoveryFixture : CrumbInstalledDiscovery
@property(nonatomic,strong) NSArray<NSString *> *paths;
@property(nonatomic,strong) NSData *database;
@property(nonatomic,strong) NSData *metadata;
@property(nonatomic) BOOL scanned;
@end
@implementation DiscoveryFixture
- (NSDictionary *)listBundleParent {return @{@"opened":@NO,@"errno":@(EACCES)};}
- (void)scanInodes {self.scanned=YES;for(NSString *path in self.paths)[self processInodePath:path];}
- (NSData *)readFile:(NSURL *)url limit:(size_t)limit error:(int *)error {
  (void)limit;NSData *data=[url.path hasSuffix:@".csstore"]?self.database:
    [url.path hasSuffix:@"Target.app/Info.plist"]?self.metadata:nil;
  *error=data?0:EACCES;return data;
}
@end
static void word(NSMutableData *b,NSUInteger p,uint32_t v){memcpy((uint8_t *)b.mutableBytes+p,&v,4);}
static NSUInteger string(NSMutableData *b,NSUInteger p,uint32_t sid,NSString *s){
  NSData *data=[s dataUsingEncoding:NSUTF8StringEncoding];word(b,p,sid);word(b,p+4,(uint32_t)data.length);
  memcpy((uint8_t *)b.mutableBytes+p+8,data.bytes,data.length);return (p+8+data.length+4)&~(NSUInteger)3;
}
int main(void){@autoreleasepool {
  NSUserDefaults *defaults=NSUserDefaults.standardUserDefaults;id savedHint=[defaults objectForKey:@"CrumbInstalledLSStoreHint"];
  [defaults removeObjectForKey:@"CrumbInstalledLSStoreHint"];
  NSString *target=@"/private/var/containers/Bundle/Application/01234567-89AB-CDEF-0123-456789ABCDEF/Target.app";
  NSString *store=@"/private/var/mobile/Containers/Data/InternalDaemon/FEDCBA98-7654-3210-FEDC-BA9876543210/Library/Caches/com.apple.LaunchServices-test.csstore";
  NSMutableData *data=[NSMutableData dataWithLength:16384];memcpy(data.mutableBytes,"bdsl",4);
  word(data,260,0x250);word(data,264,20<<2);word(data,276,10<<2);
  NSUInteger p=string(data,1024,10,@"test.selected");string(data,p,20,target);
  NSData *metadata=[NSPropertyListSerialization dataWithPropertyList:@{@"CFBundleIdentifier":@"test.selected"} format:NSPropertyListBinaryFormat_v1_0 options:0 error:nil];
  __block int grants=0;NSMutableArray *released=[NSMutableArray array];
  NSURL *work=[NSURL fileURLWithPath:NSTemporaryDirectory()];
  int64_t (^grant)(NSURL *,NSMutableDictionary *)=^int64_t(NSURL *url,NSMutableDictionary *report){
    grants++;report[@"code"]=@"boundary_stub";(void)url;return 40+grants;
  };
  DiscoveryFixture *probe=[[DiscoveryFixture alloc]initWithBundleID:@"test.selected" appName:@"Target.app" work:work grant:grant release:^(int64_t h){[released addObject:@(h)];} cancelled:^BOOL{return NO;} progress:^(NSString *stage){(void)stage;}];
  probe.paths=@[store];probe.database=data;probe.metadata=metadata;
  NSArray *found=[probe run];assert(found.count==1 && [[found[0] path] isEqual:target]);
  assert(grants==2 && [released containsObject:@42] && ![released containsObject:@41]);
  assert([probe.report[@"metadata_verified_count"] intValue]==1);
  NSString *report=[[NSString alloc]initWithData:[NSJSONSerialization dataWithJSONObject:probe.report options:0 error:nil] encoding:NSUTF8StringEncoding];
  assert(![report containsString:target] && ![report containsString:store] && ![report containsString:@"test.selected"]);
  [probe close];assert([released containsObject:@41]);
  [defaults removeObjectForKey:@"CrumbInstalledLSStoreHint"];
  // Visible paths and unreadable metadata remain unverified candidates.
  DiscoveryFixture *denied=[[DiscoveryFixture alloc]initWithBundleID:@"test.selected" appName:@"Target.app" work:work grant:grant release:^(int64_t h){(void)h;} cancelled:^BOOL{return NO;} progress:^(NSString *stage){(void)stage;}];
  denied.paths=@[target];assert([denied run].count==1);
  assert([denied.report[@"metadata_verified_count"] intValue]==0 && [denied.report[@"absence_not_proven"] boolValue]);[denied close];
  // An unrelated basename is not opened as the selected app.
  DiscoveryFixture *other=[[DiscoveryFixture alloc]initWithBundleID:@"test.selected" appName:@"Target.app" work:work grant:grant release:^(int64_t h){(void)h;} cancelled:^BOOL{return NO;} progress:^(NSString *stage){(void)stage;}];
  other.paths=@[[target stringByReplacingOccurrencesOfString:@"Target.app" withString:@"Other.app"]];other.metadata=metadata;
  assert([other run].count==0 && [other.report[@"outcome"] isEqual:@"not_located_within_probe"]);[other close];
  // Cancellation does not perform a registry grant or an inode scan.
  int before=grants;
  DiscoveryFixture *cancelled=[[DiscoveryFixture alloc]initWithBundleID:@"test.selected" appName:@"Target.app" work:work grant:grant release:^(int64_t h){(void)h;} cancelled:^BOOL{return YES;} progress:^(NSString *stage){(void)stage;}];
  assert([cancelled run].count==0 && grants==before && !cancelled.scanned);[cancelled close];
  if(savedHint)[defaults setObject:savedHint forKey:@"CrumbInstalledLSStoreHint"];else [defaults removeObjectForKey:@"CrumbInstalledLSStoreHint"];
  puts("installed discovery: registry candidates, identity gates, unreadable metadata, redaction, held-grant lifetime and cancellation passed (filesystem/token doubles)");
}}
