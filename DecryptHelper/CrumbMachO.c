#include "CrumbMachO.h"
#include <string.h>

static uint32_t word(const unsigned char *p) {
    return (uint32_t)p[0] | (uint32_t)p[1] << 8 |
           (uint32_t)p[2] << 16 | (uint32_t)p[3] << 24;
}

int crumb_macho(const unsigned char *b, size_t n, CrumbEncryption *out) {
    if (n < 4) return 0;
    uint32_t magic = word(b);
    if (magic == 0xbebafeca || magic == 0xbfbafeca || magic == 0xfeedface) return -1;
    if (magic != 0xfeedfacf) return 0;
    if (n < 32 || word(b + 4) != 0x0100000c) return -1;
    uint32_t count = word(b + 16), size = word(b + 20);
    if (size > 16 * 1024 * 1024 || size > n - 32 || count > size / 8) return -1;
    size_t end = 32 + size, p = 32;
    int found = 0;
    for (uint32_t i = 0; i < count; i++) {
        if (p > end || end - p < 8) return -1;
        uint32_t cmd = word(b + p), len = word(b + p + 4);
        if (len < 8 || len % 4 || len > end - p) return -1;
        if (cmd == 0x2c || cmd == 0x21) {
            if (len < (cmd == 0x2c ? 24 : 20)) return -1;
            uint32_t off = word(b + p + 8), length = word(b + p + 12), cid = word(b + p + 16);
            if (off > n || length > n - off) return -1;
            if (cid) {
                if (found || !length) return -1;
                found = 1;
                *out = (CrumbEncryption){word(b + 4), word(b + 8), off, length, p + 16};
            }
        }
        p += len;
    }
    return p == end ? found : -1;
}
