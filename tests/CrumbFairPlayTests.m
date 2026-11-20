#import <Foundation/Foundation.h>
#import <mach/mach.h>
#import <sys/mman.h>
#import <unistd.h>
#import <errno.h>
static int targetAttempt, calls;
static uint32_t observed[4];
static BOOL noChange, shortPage, failFlush, mutateSource;
static NSString *inputPath;
static int test_fsync(int fd) { if (failFlush) { errno = ENOSPC; return -1; } return fsync(fd); }
static kern_return_t test_vm_read(vm_map_t task, vm_address_t address, vm_size_t size,
                                  vm_address_t output, vm_size_t *got) {
    if (shortPage) { *got = 0; return KERN_SUCCESS; }
    return vm_read_overwrite(task, address, size, output, got);
}
#define CRUMB_FAIRPLAY_LIBRARY 1
#define CRUMB_FAIRPLAY_TESTING 1
#define fsync test_fsync
#define vm_read_overwrite test_vm_read
#import "../DecryptHelper/main.m"
#include "../DecryptHelper/CrumbMachO.c"
#include "../DecryptHelper/CrumbFileAccess.c"
#undef fsync
#undef vm_read_overwrite
static int remap_image(void *address, size_t length, uint32_t cryptid, uint32_t cpu, uint32_t subtype) {
    (void)cpu; (void)subtype;
    assert(calls < 4); observed[calls++] = cryptid;
    if (targetAttempt && calls != targetAttempt) { errno = EACCES; return -1; }
    if (!noChange) {
        assert(mprotect(address, length, PROT_READ | PROT_WRITE) == 0);
        memset(address, 0, length); // Valid low-entropy DATA, no ARM64 landmarks.
        if (mutateSource) {
            int fd = open(inputPath.fileSystemRepresentation, O_WRONLY);
            assert(fd >= 0 && pwrite(fd, address, length, getpagesize()) == (ssize_t)length);
            close(fd);
        }
    }
    return 0;
}
static void runCase(NSURL *root, int attempt, BOOL unchanged, BOOL shortRead, BOOL flushFailure, BOOL mutate) {
    calls = 0; targetAttempt = attempt; noChange = unchanged; shortPage = shortRead;
    failFlush = flushFailure; mutateSource = mutate;
    atomic_store(&jobCancelled, false);
    size_t pg = getpagesize();
    NSMutableData *image = [NSMutableData dataWithLength:pg * 2];
    uint32_t header[] = {0xfeedfacf,0x0100000c,0,2,1,24,0,0,0x2c,24,(uint32_t)pg,(uint32_t)pg,1,0};
    memcpy(image.mutableBytes, header, sizeof header);
    for (size_t i = 0; i < pg; i++) ((uint8_t *)image.mutableBytes)[pg + i] = i % 256;
    NSURL *source = [root URLByAppendingPathComponent:@"source"];
    NSURL *dest = [root URLByAppendingPathComponent:@"result"];
    [[NSFileManager defaultManager] removeItemAtURL:dest error:nil];
    assert([image writeToURL:source atomically:YES]); inputPath = source.path;
    NSError *error = nil; BOOL encrypted = NO;
    BOOL success = decrypt(source, dest, &encrypted, &error);
    BOOL expected = attempt > 0 && !unchanged && !shortRead && !flushFailure;
    assert(success == expected && encrypted);
    assert(calls == (expected || flushFailure ? attempt : 4));
    for (int i = 0; i < calls; i++) assert(observed[i] == (i < 2 ? 2 : 1));
    if (success) {
        NSData *result = [NSData dataWithContentsOfURL:dest];
        assert(result.length == image.length);
        uint32_t cryptid; memcpy(&cryptid, (const uint8_t *)result.bytes + 48, 4); assert(cryptid == 0);
        const uint8_t *bytes = result.bytes;
        for (size_t i = pg; i < pg * 2; i++) assert(bytes[i] == 0);
    } else {
        assert(error != nil && ![[NSFileManager defaultManager] fileExistsAtPath:dest.path]);
        assert(![error.localizedDescription containsString:@"licensed device account"]);
        if (shortRead) assert([error.localizedDescription containsString:@"page read failed"]);
        if (flushFailure) assert([error.localizedDescription containsString:@"finalize"]);
        if (unchanged) assert([error.localizedDescription containsString:@"changed no"]);
    }
    if (!mutate) assert([[NSData dataWithContentsOfURL:source] isEqualToData:image]);
}
int main(void) {
    @autoreleasepool {
        NSFileManager *fm = NSFileManager.defaultManager;
        NSURL *root = [fm.temporaryDirectory URLByAppendingPathComponent:NSUUID.UUID.UUIDString];
        assert([fm createDirectoryAtURL:root withIntermediateDirectories:YES attributes:nil error:nil]);
        for (int attempt = 1; attempt <= 4; attempt++) runCase(root, attempt, NO, NO, NO, NO);
        runCase(root, 0, YES, NO, NO, NO); // A successful syscall returning ciphertext is rejected.
        runCase(root, 1, NO, YES, NO, NO); // KERN_SUCCESS with short data is still rejected.
        runCase(root, 1, NO, NO, YES, NO); // No result survives an output fsync failure.
        runCase(root, 1, NO, NO, NO, YES); // Snapshot survives changed source/cache reference.
        assert(!crumb_fairplay_verified(true, 1, 2, true, 8, 1));
        assert(!crumb_fairplay_verified(true, 2, 2, true, 8, 7.6));
        assert(!crumb_fairplay_verified(true, 2, 2, true, 7.8, 1));
        assert(!crumb_fairplay_verified(false, 2, 2, true, 8, 1));
        [fm removeItemAtURL:root error:nil];
        puts("FairPlay tests passed: all modes, data-only plaintext, unchanged bytes, short pages, output failures and immutable ciphertext snapshot.");
    }
}
