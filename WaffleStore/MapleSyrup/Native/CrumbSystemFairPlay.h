#ifndef CRUMB_SYSTEM_FAIRPLAY_H
#define CRUMB_SYSTEM_FAIRPLAY_H
#include <stdint.h>
typedef struct {
    uint64_t slice_offset, slice_size, crypt_offset, crypt_size, cryptid_field;
    uint32_t cpu_type, cpu_subtype;
} CrumbSystemRange;
typedef struct {
    uint64_t pages_total, pages_read, pages_changed;
    int32_t remap_result, remap_errno, read_error;
} CrumbSystemRangeResult;
// source is in the installed bundle. destination is a private copy of that
// complete file. Only the verified encrypted range and cryptid field are written.
int CrumbSystemDecryptRange(const char *source, const char *destination,
                           CrumbSystemRange range, CrumbSystemRangeResult *result);
#endif
