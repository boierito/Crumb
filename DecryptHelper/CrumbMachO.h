#pragma once
#include <stddef.h>
#include <stdint.h>
typedef struct {uint32_t cpu, subtype, offset, length;size_t command;} CrumbEncryption;
// 1 encrypted thin arm64, 0 plain/resource, -1 malformed/unsupported.
int crumb_macho(const unsigned char *b,size_t n,CrumbEncryption *out);
