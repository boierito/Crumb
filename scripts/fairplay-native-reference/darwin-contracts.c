#include <pthread.h>
#include <stdatomic.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/time.h>
#include <unistd.h>

static pthread_once_t control = PTHREAD_ONCE_INIT;
static _Atomic int calls;
static unsigned char inside[sizeof(control)];
static int payload[4];
static void initialize(void) {
    memcpy(inside, &control, sizeof(control));
    atomic_fetch_add(&calls, 1);
    for (int i = 0; i < 4; ++i) payload[i] = 41 + i;
}
static void *worker(void *unused) {
    (void)unused;
    pthread_once(&control, initialize);
    for (int i = 0; i < 4; ++i) if (payload[i] != 41 + i) return (void *)1;
    return NULL;
}
static void hexbytes(const unsigned char *p, size_t size) {
    for (size_t i = 0; i < size; ++i) printf("%02x", p[i]);
}
int main(void) {
#ifndef __APPLE__
    fprintf(stderr, "Darwin execution required\n");
    return 2;
#else
    unsigned char before[sizeof(control)];
    memcpy(before, &control, sizeof(control));
    pthread_t threads[8];
    for (int i = 0; i < 8; ++i) if (pthread_create(&threads[i], NULL, worker, NULL)) return 1;
    for (int i = 0; i < 8; ++i) {
        void *result = NULL;
        if (pthread_join(threads[i], &result) || result) return 1;
    }
    if (atomic_load(&calls) != 1) return 1;
    struct timeval tv;
    memset(&tv, 0xa5, sizeof(tv));
    if (gettimeofday(&tv, NULL) || tv.tv_usec < 0 || tv.tv_usec >= 1000000) return 1;
    char name[] = "/tmp/crumb-darwin-contracts-XXXXXX";
    int fd = mkstemp(name);
    if (fd < 0) return 1;
    unlink(name);
    unsigned char data[128] = {0};
    struct stat info;
    if (write(fd, data, sizeof(data)) != sizeof(data) || fstat(fd, &info)) { close(fd); return 1; }
    close(fd);
    if (!S_ISREG(info.st_mode) || info.st_size != sizeof(data)) return 1;
    printf("{\"schema\":\"crumb-native-darwin-contracts-v1\",\"once_size\":%zu,\"once_initial\":\"", sizeof(control));
    hexbytes(before, sizeof(before));
    printf("\",\"once_inside\":\""); hexbytes(inside, sizeof(inside));
    printf("\",\"once_final\":\""); hexbytes((unsigned char *)&control, sizeof(control));
    printf("\",\"initializer_calls\":%d,\"waiters_observed_initialized_data\":true,", atomic_load(&calls));
    printf("\"timeval_size\":%zu,\"timeval_usec_offset\":%zu,\"timeval_usec_size\":%zu,", sizeof(tv), offsetof(struct timeval, tv_usec), sizeof(tv.tv_usec));
    printf("\"timeval_bytes\":\""); hexbytes((unsigned char *)&tv, sizeof(tv));
    printf("\",\"stat_size\":%zu,\"stat_mode_offset\":%zu,\"stat_nlink_offset\":%zu,\"stat_size_offset\":%zu,\"stat_mtime_offset\":%zu,\"regular_file_checked\":true}\n", sizeof(info), offsetof(struct stat, st_mode), offsetof(struct stat, st_nlink), offsetof(struct stat, st_size), offsetof(struct stat, st_mtimespec));
    return 0;
#endif
}
