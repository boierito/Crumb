#ifndef CRUMB_IPA_BRIDGE_H
#define CRUMB_IPA_BRIDGE_H
#include <stddef.h>
#include <stdint.h>

int CrumbPrepareIPA(char *source, char *destination, unsigned char *input, size_t length, unsigned char **output, size_t *outputLength);
int CrumbInspectIPA(char *url, char *bundle, unsigned char **output, size_t *outputLength);
void CrumbIPAFree(unsigned char *pointer, size_t length);
#endif
