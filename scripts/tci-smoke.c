#define _GNU_SOURCE
#include "deny-exec.h"
#include "TCIProbe.h"
int main(void) {
    if (deny_executable_memory() != 0) { perror("seccomp"); return 2; }
    waffle_tci_result result = waffle_probe_tci();
    printf("TCI with PROT_EXEC denied: status=%d RAX=%llu instruction-hooks=%u\n",
           result.error, (unsigned long long)result.guest_rax, result.instruction_hooks);
    return result.error != 0 || result.guest_rax != 42 || result.instruction_hooks != 4;
}
