#ifndef CRUMB_FILE_ACCESS_H
#define CRUMB_FILE_ACCESS_H
#include <stdint.h>

typedef struct {
  int error_code;
  int metadata_available;
  uint32_t mode;
  uint32_t owner;
  uint32_t group;
} CrumbFileAccess;

// Returns a read-only descriptor, or -1 with the original open errno retained.
// Never follows a final symlink. Failure metadata uses lstat, not a retry that
// would weaken that restriction.
int crumb_open_installed_image(const char *path, CrumbFileAccess *status);
#endif
