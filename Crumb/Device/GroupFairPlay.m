// Reuse Crumb's checked FairPlay page reader without installing a helper.
#define CRUMB_FAIRPLAY_LIBRARY 1
#import "../../DecryptHelper/main.m"
#include "../../DecryptHelper/CrumbFileAccess.c"
#include "../../DecryptHelper/CrumbMachO.c"
#import "GroupFairPlay.h"
void crumb_group_reset(void) { atomic_store(&jobCancelled, false); }
void crumb_group_cancel(void) { atomic_store(&jobCancelled, true); }
BOOL crumb_group_decrypt(NSURL *source, NSURL *destination, BOOL *encrypted, NSError **error) {
    return decrypt(source, destination, encrypted, error);
}
