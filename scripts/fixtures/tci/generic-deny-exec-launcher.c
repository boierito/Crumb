#define _GNU_SOURCE
#include "deny-exec.h"
int main(int argc, char **argv) {
    if (argc < 2) return 2;
    if (deny_executable_memory() != 0) { perror("seccomp"); return 2; }
    execv(argv[1], argv + 1);
    perror("execv");
    return 2;
}
