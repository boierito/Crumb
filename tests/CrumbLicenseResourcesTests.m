#import "../DecryptHelper/CrumbLicenseResources.h"
#import <CommonCrypto/CommonDigest.h>
#import <errno.h>
#import <stdio.h>
#import <stdlib.h>

static void require(BOOL condition, NSString *message) {
  if (!condition) {
    fprintf(stderr, "Resource test failed: %s\n", message.UTF8String);
    exit(1);
  }
}
static NSString *sha(NSData *data) {
  unsigned char hash[CC_SHA256_DIGEST_LENGTH];
  CC_SHA256(data.bytes, (CC_LONG)data.length, hash);
  NSMutableString *text = [NSMutableString string];
  for (int i = 0; i < CC_SHA256_DIGEST_LENGTH; i++)
    [text appendFormat:@"%02x", hash[i]];
  return text;
}
static void runCase(int number) {
  NSFileManager *fm = NSFileManager.defaultManager;
  NSURL *root = [[NSURL fileURLWithPath:NSTemporaryDirectory()]
      URLByAppendingPathComponent:NSUUID.UUID.UUIDString
                      isDirectory:YES];
  NSURL *bundle = [root URLByAppendingPathComponent:@"helper" isDirectory:YES];
  NSURL *copies = [bundle URLByAppendingPathComponent:@"CrumbLicenseResources"
                                          isDirectory:YES];
  NSURL *dest = [root URLByAppendingPathComponent:@"output" isDirectory:YES];
  require([fm createDirectoryAtURL:copies
              withIntermediateDirectories:YES
                               attributes:nil
                                    error:nil],
          @"fixture directory");
  NSString *relative = @"PlugIns/Share.appex/SC_Info/Share.v4.supp";
  NSData *original = [@"original supplemental license bytes"
      dataUsingEncoding:NSUTF8StringEncoding];
  if (number == 7)
    original = [NSData data];
  if (number == 4) {
    uint32_t words[16] = {0xfeedfacf, 0x0100000c, 0,  0, 1, 24, 0, 0,
                          0x2c,       24,         56, 8, 1, 0,  0, 0};
    original = [NSData dataWithBytes:words length:sizeof(words)];
  }
  NSMutableDictionary *item = [@{
    @"file" : @"00000000.bin",
    @"size" : @(original.length),
    @"sha256" : sha(original)
  } mutableCopy];
  NSURL *saved = [copies URLByAppendingPathComponent:item[@"file"]];
  require([[original base64EncodedDataWithOptions:0] writeToURL:saved
                                                     atomically:YES],
          @"encoded fixture");
  if (number == 2)
    item[@"sha256"] = [@"0" stringByPaddingToLength:64
                                         withString:@"0"
                                    startingAtIndex:0];
  if (number == 3)
    relative = @"../escape/SC_Info/Share.sinf";
  if (number == 5) {
    NSURL *outside = [root URLByAppendingPathComponent:@"outside"];
    [[original base64EncodedDataWithOptions:0] writeToURL:outside
                                               atomically:YES];
    [fm removeItemAtURL:saved error:nil];
    require([fm createSymbolicLinkAtURL:saved
                     withDestinationURL:outside
                                  error:nil],
            @"link fixture");
  }
  if (number == 6)
    item[@"size"] = @(4 * 1024 * 1024 + 1);
  NSError *error = nil;
  BOOL ok = crumb_copy_license_resource(bundle, relative, item, dest, &error);
  if (number == 1 || number == 7) {
    // No installed Target/SC_Info exists at all: restoration must use only the
    // snapshot.
    require(ok, error.localizedDescription ?: @"restoration failed");
    require(
        [[NSData
            dataWithContentsOfURL:[dest URLByAppendingPathComponent:relative]]
            isEqual:original],
        @"original bytes preserved");
  } else {
    require(!ok && error != nil, @"invalid input was accepted");
    require(![fm fileExistsAtPath:dest.path],
            @"failed validation wrote output");
    require(
        ![fm
            fileExistsAtPath:[root URLByAppendingPathComponent:@"escape"].path],
        @"path escaped output");
    if (number == 5)
      require(error.code == ELOOP, @"symlink was followed");
  }
  [fm removeItemAtURL:root error:nil];
}
int main(void) {
  @autoreleasepool {
    for (int i = 1; i <= 7; i++)
      runCase(i);
    puts("7 license-resource restoration tests passed");
  }
  return 0;
}
