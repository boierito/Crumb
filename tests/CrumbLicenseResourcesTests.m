#import "../DecryptHelper/CrumbIPAArchive.h"
#import "../DecryptHelper/CrumbLicenseResources.h"
#import <CommonCrypto/CommonDigest.h>
#import <errno.h>
#import <stdio.h>
#import <stdlib.h>
#import <string.h>

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
static NSMutableDictionary *archiveItem(NSURL *copies, NSString *name,
                                        NSData *data) {
  require([[data base64EncodedDataWithOptions:0]
              writeToURL:[copies URLByAppendingPathComponent:name]
              atomically:YES],
          @"archive snapshot");
  return [@{
    @"file" : name,
    @"size" : @(data.length),
    @"sha256" : sha(data)
  } mutableCopy];
}
static void runArchiveCase(int number, NSString *fixturePath) {
  NSFileManager *fm = NSFileManager.defaultManager;
  NSURL *root = [[NSURL fileURLWithPath:NSTemporaryDirectory()]
      URLByAppendingPathComponent:NSUUID.UUID.UUIDString];
  NSURL *bundle = [root URLByAppendingPathComponent:@"helper"];
  NSURL *copies = [bundle URLByAppendingPathComponent:@"CrumbArchiveResources"];
  NSURL *dest = [root URLByAppendingPathComponent:@"working"];
  require([fm createDirectoryAtURL:copies
              withIntermediateDirectories:YES
                               attributes:nil
                                    error:nil],
          @"archive directory");
  NSData *data =
      [@"original signature bytes" dataUsingEncoding:NSUTF8StringEncoding];
  if (number == 3) {
    NSMutableData *large = [NSMutableData dataWithLength:5 * 1024 * 1024 + 14];
    memset(large.mutableBytes, 'x', large.length);
    data = large;
  }
  NSMutableDictionary *item = archiveItem(copies, @"00000000.bin", data);
  NSError *error = nil;
  if (number == 2) {
    for (NSString *path in @[
           @"../escape/META-INF/file", @"/absolute/META-INF/file",
           @"Payload/Other.app/Test", @"payload/Test.app/Test"
         ]) {
      require(!crumb_copy_archive_resource(bundle, path, item, dest, &error),
              @"unsafe root path accepted");
    }
    require(!crumb_create_archive_directory(dest, @"Payload", &error),
            @"root directory overwrote Payload");
    require(![fm fileExistsAtPath:dest.path], @"unsafe root path wrote output");
  } else {
    NSString *relative = number == 3 ? @"iTunesArtwork"
                                     : @"META-INF/com.apple.ZipMetadata.plist";
    require(crumb_copy_archive_resource(bundle, relative, item, dest, &error),
            error.localizedDescription ?: @"archive restore");
    require(
        [[NSData
            dataWithContentsOfURL:[dest URLByAppendingPathComponent:relative]]
            isEqual:data],
        @"root bytes changed");
    require(crumb_create_archive_directory(dest, @"META-INF/empty", &error),
            @"empty root directory");
    if (number == 3) {
      item[@"size"] = @(512ULL * 1024 * 1024 + 1);
      require(!crumb_copy_archive_resource(bundle, @"too-large", item, dest,
                                           &error),
              @"root size limit ignored");
    }
    if (number == 4) {
      NSMutableDictionary *inventory = [NSMutableDictionary dictionary];
      inventory[relative] = item;
      NSData *metadata = [@"bplist00original-store-fields"
          dataUsingEncoding:NSUTF8StringEncoding];
      NSMutableDictionary *metaItem =
          archiveItem(copies, @"00000001.bin", metadata);
      require(crumb_copy_archive_resource(bundle, @"iTunesMetadata.plist",
                                          metaItem, dest, &error),
              @"iTunes metadata restore");
      inventory[@"iTunesMetadata.plist"] = metaItem;
      NSURL *app = [dest URLByAppendingPathComponent:@"Payload/Test.app"];
      [fm createDirectoryAtURL:app
          withIntermediateDirectories:YES
                           attributes:nil
                                error:nil];
      [[@"plain test image" dataUsingEncoding:NSUTF8StringEncoding]
          writeToURL:[app URLByAppendingPathComponent:@"Test"]
          atomically:YES];
      NSURL *output = fixturePath
                          ? [NSURL fileURLWithPath:fixturePath]
                          : [root URLByAppendingPathComponent:@"fixture.ipa"];
      NSMutableArray *hashes = [NSMutableArray array];
      require(crumb_pack_ipa(dest, output, hashes, &error),
              error.localizedDescription ?: @"IPA packing");
      require(hashes.count == 5,
              [NSString stringWithFormat:@"unexpected root archive entries: %@",
                                         hashes]);
      NSDictionary *manifest = @{
        @"schema" : @3,
        @"files" : hashes,
        @"archive_resources" : inventory,
        @"archive_directories" : @[ @"META-INF", @"META-INF/empty" ]
      };
      NSData *json = [NSJSONSerialization dataWithJSONObject:manifest
                                                     options:0
                                                       error:nil];
      require([json writeToFile:[output.path stringByAppendingString:@".json"]
                     atomically:YES],
              @"fixture manifest");
    }
  }
  [fm removeItemAtURL:root error:nil];
}
int main(int argc, const char *argv[]) {
  @autoreleasepool {
    for (int i = 1; i <= 7; i++)
      runCase(i);
    for (int i = 1; i <= 4; i++)
      runArchiveCase(i,
                     argc > 1 ? [NSString stringWithUTF8String:argv[1]] : nil);
    puts("11 preserved-resource and IPA archive tests passed");
  }
  return 0;
}
