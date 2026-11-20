#pragma once
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
// Aligned with iPadn v2.3.0's entropy verdict. ARM64 instruction landmarks
// are not a hard requirement: an encrypted range can contain only data.
typedef struct { uint32_t cryptid; bool signature; } CrumbFairPlayMode;
static const CrumbFairPlayMode crumb_fairplay_modes[] = {
    {2, false}, {2, true}, {1, false}, {1, true}
};
static inline bool crumb_fairplay_verified(bool remapped, size_t read,
        size_t expected, bool changed, double before, double after) {
    return remapped && expected > 0 && read == expected && changed &&
           before > 7.9 && after < 7.5 && before - after >= 0.5;
}
