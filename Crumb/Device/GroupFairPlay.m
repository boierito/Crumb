// Reuse Crumb's checked FairPlay page reader without installing a helper.
#define CRUMB_FAIRPLAY_LIBRARY 1
#import "../../DecryptHelper/main.m"
#include "../../DecryptHelper/CrumbFileAccess.c"
#include "../../DecryptHelper/CrumbMachO.c"
#import "GroupFairPlay.h"
#import "FairPlayCore.h"
void crumb_group_reset(void) { atomic_store(&jobCancelled, false); }
void crumb_group_cancel(void) { atomic_store(&jobCancelled, true); }
// iPadN FairPlayCore is the only production decrypt engine. The old Crumb helper
// is kept for independent regression fixtures, not selected by the app.
BOOL crumb_group_decrypt(NSURL *source, NSURL *destination, BOOL *encrypted, NSError **error) {
    if (encrypted) *encrypted = NO;
    if (atomic_load(&jobCancelled)) {
        if (error) *error = [NSError errorWithDomain:@"CrumbFairPlay" code:-2
            userInfo:@{NSLocalizedDescriptionKey: @"Decryption cancelled."}];
        return NO;
    }
    fp_file_stats_t stats = {0};
    const int result = fp_decrypt_file_ex(source.fileSystemRepresentation,
        destination.fileSystemRepresentation, NULL, NULL, &stats);
    if (atomic_load(&jobCancelled)) {
        [NSFileManager.defaultManager removeItemAtURL:destination error:nil];
        if (error) *error = [NSError errorWithDomain:@"CrumbFairPlay" code:-2
            userInfo:@{NSLocalizedDescriptionKey: @"Decryption cancelled."}];
        return NO;
    }
    if (result < 0) {
        [NSFileManager.defaultManager removeItemAtURL:destination error:nil];
        if (error) *error = [NSError errorWithDomain:@"CrumbFairPlay" code:-1
            userInfo:@{NSLocalizedDescriptionKey:
                @"FairPlayCore rejected the image. Check the Apple license, App Group entitlement and iOS permissions."}];
        return NO;
    }
    if (encrypted) *encrypted = (result == 0 && stats.slices_ok > 0);
    return YES;
}

NSArray<NSString *> *crumb_group_signed_groups(void) {
    // These Security entry points are not declared in the public iOS SDK.
    // Resolve them at runtime and retain profile/configuration fallbacks.
    typedef CFTypeRef (*CreateTask)(CFAllocatorRef);
    typedef CFTypeRef (*CopyEntitlement)(CFTypeRef, CFStringRef, CFErrorRef *);
    CreateTask createTask = (CreateTask)dlsym(RTLD_DEFAULT, "SecTaskCreateFromSelf");
    CopyEntitlement copyEntitlement = (CopyEntitlement)dlsym(RTLD_DEFAULT, "SecTaskCopyValueForEntitlement");
    if (!createTask || !copyEntitlement) return nil;
    CFTypeRef task = createTask(kCFAllocatorDefault);
    if (!task) return nil;
    CFErrorRef error = NULL;
    CFTypeRef raw = copyEntitlement(task, CFSTR("com.apple.security.application-groups"), &error);
    CFRelease(task);
    if (error) { CFRelease(error); if (raw) CFRelease(raw); return nil; }
    if (!raw) return @[];
    id value = CFBridgingRelease(raw);
    if (![value isKindOfClass:NSArray.class]) return @[];
    NSMutableArray<NSString *> *groups = [NSMutableArray array];
    for (id group in value) {
        if ([group isKindOfClass:NSString.class] && [group length] > 0) [groups addObject:group];
    }
    return groups;
}
