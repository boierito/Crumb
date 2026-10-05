#ifndef WAFFLE_TCI_PROBE_H
#define WAFFLE_TCI_PROBE_H
#include <stdint.h>
typedef struct {
    int32_t error;
    uint64_t guest_rax;
    uint32_t instruction_hooks;
} waffle_tci_result;
waffle_tci_result waffle_probe_tci(void);
#endif
