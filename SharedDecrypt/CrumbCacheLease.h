#pragma once
#include <stdbool.h>
#include <stdint.h>
#include <limits.h>
typedef struct {
  char root[PATH_MAX], target[PATH_MAX], staging[PATH_MAX], identifier[37];
  uint64_t original_device, original_inode, link_device, link_inode;
  bool had_caches;
  int lock_fd;
} CrumbCacheLease;
// Read-only check. Empty/missing Caches is eligible; files, links, occupied or
// inaccessible directories fail with errno. No shared contents are removed.
int crumb_cache_inspect(const char *root);
// Creates only a lock and a symlink in the caller's private staging directory.
// Persist this structure in a private ledger before calling activate.
int crumb_cache_begin(CrumbCacheLease *,const char *root,const char *target,
                      const char *private_staging,const char *identifier);
int crumb_cache_activate(CrumbCacheLease *);
// Also accepts a ledger-reconstructed lease with lock_fd=-1. Restores the same
// original directory inode and contents; never recursively deletes group data.
int crumb_cache_restore(CrumbCacheLease *);
