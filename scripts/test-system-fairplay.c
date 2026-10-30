#include "CrumbSystemFairPlay.h"
#include <assert.h>
#include <errno.h>
// Unsupported hosts must never supply bytes or report a fake native success.
int main(void) {
    CrumbSystemRange range={0};CrumbSystemRangeResult result;
    assert(CrumbSystemDecryptRange(0,0,range,0)==-1);
    assert(CrumbSystemDecryptRange(0,0,range,&result)==-1);
    assert(result.remap_result==-1);
    assert(result.remap_errno==ENOTSUP);
    assert(result.pages_read==0&&result.pages_changed==0);
    return 0;
}
