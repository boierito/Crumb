#include "TCIProbe.h"
#include <unicorn/unicorn.h>

static void count_instruction(uc_engine *engine, uint64_t address, uint32_t size, void *opaque) {
    (void)engine; (void)address; (void)size;
    waffle_tci_result *result = opaque;
    result->instruction_hooks++;
}

// Exercises x86-64 guest call/return, stack memory, arithmetic and C hooks.
// This is NOT SAP and never makes a cryptographic signature claim.
waffle_tci_result waffle_probe_tci(void) {
    waffle_tci_result result = {0};
    uc_engine *engine = NULL;
    uc_hook hook;
    uc_err error = uc_open(UC_ARCH_X86, UC_MODE_64, &engine);
    if (error == UC_ERR_OK) error = uc_mem_map(engine, 0x1000, 0x1000, UC_PROT_ALL);
    if (error == UC_ERR_OK) error = uc_mem_map(engine, 0x8000, 0x1000, UC_PROT_READ | UC_PROT_WRITE);
    const unsigned char code[] = {
        0x48, 0xC7, 0xC0, 40, 0, 0, 0, // mov rax,40
        0xE8, 1, 0, 0, 0,             // call 0x100d
        0xF4,                          // stop before hlt at 0x100c
        0x48, 0x83, 0xC0, 2, 0xC3     // add rax,2; ret
    };
    if (error == UC_ERR_OK) error = uc_mem_write(engine, 0x1000, code, sizeof(code));
    uint64_t stack = 0x8FF0;
    if (error == UC_ERR_OK) error = uc_reg_write(engine, UC_X86_REG_RSP, &stack);
    if (error == UC_ERR_OK) error = uc_hook_add(engine, &hook, UC_HOOK_CODE, (void *)count_instruction, &result, 1, 0);
    if (error == UC_ERR_OK) error = uc_emu_start(engine, 0x1000, 0x100C, 1000000, 100);
    if (error == UC_ERR_OK) error = uc_reg_read(engine, UC_X86_REG_RAX, &result.guest_rax);
    result.error = error;
    if (engine) uc_close(engine);
    return result;
}
