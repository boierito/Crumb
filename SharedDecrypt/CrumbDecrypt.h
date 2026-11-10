#pragma once
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
typedef enum {
  CRUMB_OPEN=1, CRUMB_INSPECT, CRUMB_MAP, CRUMB_PARSE, CRUMB_COPY,
  CRUMB_SIGNATURE, CRUMB_REMAP, CRUMB_READ, CRUMB_VERIFY, CRUMB_PATCH, CRUMB_COMPLETE
} CrumbDecryptStage;
typedef struct {
  int stage, error_code, signature_result, signature_errno, remap_result, remap_errno;
  int first_remap_result, first_remap_errno, remap_attempts;
  int page_result, read_errno, write_errno, cancelled;
  uint64_t pages_expected, pages_read, first_fail_offset, changed_pages;
  uint64_t oldwords, newwords;
  double entropy_before, entropy_after;
  bool encrypted, remap_attempted, signature_attempted;
} CrumbDecryptMetrics;
typedef struct {
  void *context;
  bool (*cancelled)(void *);
  void (*checkpoint)(void *, const CrumbDecryptMetrics *);
  // Test seam: production passes NULL for these three callbacks.
  int (*remap)(void *, void *, size_t, uint32_t, uint32_t);
  int (*read_page)(void *, const void *, void *, size_t, size_t *);
  int (*signature)(void *, int, const uint8_t *, size_t);
} CrumbDecryptHooks;
// Returns 1 for verified decryption, 0 for an unchanged resource, -1 on error.
// Source is always read-only. A failed destination is removed.
int crumb_decrypt_image(const char *, const char *, CrumbDecryptMetrics *, const CrumbDecryptHooks *);
