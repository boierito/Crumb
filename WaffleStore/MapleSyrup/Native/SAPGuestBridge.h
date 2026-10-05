#ifndef WAFFLE_SAP_GUEST_BRIDGE_H
#define WAFFLE_SAP_GUEST_BRIDGE_H
#include <stdint.h>
#include <stddef.h>
int WaffleSAPOpen(char *cache, unsigned char *hardware, size_t length, uint64_t *handle);
int WaffleSAPExchange(uint64_t handle, uint32_t version, unsigned char *input, size_t length,
                     unsigned char **output, size_t *outputLength, int32_t *state);
int WaffleSAPSign(uint64_t handle, unsigned char *input, size_t length,
                 unsigned char **output, size_t *outputLength);
void WaffleSAPClose(uint64_t handle);
void WaffleSAPFree(unsigned char *pointer, size_t length);
#endif
