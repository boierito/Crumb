// Host reproduction of the SAP runtime prerequisite; not an iOS emulator.
// Load Unicorn FIRST, then deny new executable mappings with Linux seccomp.
#define _GNU_SOURCE
#include <unicorn/unicorn.h>
#include <dlfcn.h>
#include <errno.h>
#include <linux/filter.h>
#include <linux/seccomp.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/prctl.h>
#include <sys/syscall.h>
#include <unistd.h>

static int deny_executable_memory(void) {
    struct sock_filter rules[] = {
        BPF_STMT(BPF_LD | BPF_W | BPF_ABS, offsetof(struct seccomp_data, nr)),
        BPF_JUMP(BPF_JMP | BPF_JEQ | BPF_K, __NR_mmap, 0, 3),
        BPF_STMT(BPF_LD | BPF_W | BPF_ABS, offsetof(struct seccomp_data, args[2])),
        BPF_JUMP(BPF_JMP | BPF_JSET | BPF_K, PROT_EXEC, 0, 1),
        BPF_STMT(BPF_RET | BPF_K, SECCOMP_RET_ERRNO | EPERM),
        BPF_STMT(BPF_LD | BPF_W | BPF_ABS, offsetof(struct seccomp_data, nr)),
        BPF_JUMP(BPF_JMP | BPF_JEQ | BPF_K, __NR_mprotect, 0, 3),
        BPF_STMT(BPF_LD | BPF_W | BPF_ABS, offsetof(struct seccomp_data, args[2])),
        BPF_JUMP(BPF_JMP | BPF_JSET | BPF_K, PROT_EXEC, 0, 1),
        BPF_STMT(BPF_RET | BPF_K, SECCOMP_RET_ERRNO | EPERM),
        BPF_STMT(BPF_RET | BPF_K, SECCOMP_RET_ALLOW)
    };
    struct sock_fprog program = {sizeof(rules) / sizeof(rules[0]), rules};
    if (prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0) != 0) return -1;
    return prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER, &program);
}

int main(int argc, char **argv) {
    setvbuf(stdout, NULL, _IONBF, 0);
    if (argc != 3) return 2;
    void *library = dlopen(argv[1], RTLD_NOW | RTLD_LOCAL);
    if (!library) { fprintf(stderr, "Unicorn load failed: %s\n", dlerror()); return 2; }
    uc_err (*open_engine)(uc_arch, uc_mode, uc_engine **) = dlsym(library, "uc_open");
    uc_err (*map)(uc_engine *, uint64_t, size_t, uint32_t) = dlsym(library, "uc_mem_map");
    uc_err (*write_mem)(uc_engine *, uint64_t, const void *, size_t) = dlsym(library, "uc_mem_write");
    uc_err (*start)(uc_engine *, uint64_t, uint64_t, uint64_t, size_t) = dlsym(library, "uc_emu_start");
    uc_err (*read_register)(uc_engine *, int, void *) = dlsym(library, "uc_reg_read");
    uc_err (*close_engine)(uc_engine *) = dlsym(library, "uc_close");
    const char *(*error_text)(uc_err) = dlsym(library, "uc_strerror");
    if (!open_engine || !map || !write_mem || !start || !read_register || !close_engine || !error_text) return 2;
    int restricted = strcmp(argv[2], "deny-exec") == 0;
    if (restricted && deny_executable_memory() != 0) { perror("seccomp"); return 2; }
    if (restricted) {
        size_t page = (size_t)sysconf(_SC_PAGESIZE);
        void *data = mmap(NULL, page, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
        if (data == MAP_FAILED) return 2;
        errno = 0;
        int status = mprotect(data, page, PROT_READ | PROT_EXEC);
        printf("ipatool-iOS-RW-to-RX: status=%d errno=%d\n", status, errno);
        if (status != -1 || errno != EPERM) return 2;
        munmap(data, page);
    }
    uc_engine *engine = NULL;
    uc_err error = open_engine(UC_ARCH_X86, UC_MODE_64, &engine);
    if (error == UC_ERR_OK) error = map(engine, 0x1000, 0x1000, UC_PROT_ALL);
    const unsigned char code[] = {0xB8, 0x2A, 0, 0, 0}; // mov eax,42
    if (error == UC_ERR_OK) error = write_mem(engine, 0x1000, code, sizeof(code));
    if (error == UC_ERR_OK) error = start(engine, 0x1000, 0x1000 + sizeof(code), 1000000, 10);
    uint64_t value = 0;
    if (error == UC_ERR_OK) error = read_register(engine, UC_X86_REG_RAX, &value);
    printf("Unicorn-2.1.4: restricted=%d status=%d (%s) RAX=%llu\n",
           restricted, error, error_text(error), (unsigned long long)value);
    if (engine) close_engine(engine);
    dlclose(library);
    if (restricted) return error != UC_ERR_OK ? 10 : 1;
    return error == UC_ERR_OK && value == 42 ? 0 : 1;
}
