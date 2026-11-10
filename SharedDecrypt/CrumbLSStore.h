#pragma once
#include <stddef.h>
#include <stdint.h>
#include <limits.h>
#define CRUMB_LS_MAX_PATHS 8
// CSStore2 layout evidence: kolbicz/WQFiles f149ba83050ef53ebdb9b73f2e83447f88057965.
// A path is a discovery candidate only; the caller must verify live bundle
// metadata and original encrypted-image hashes before attempting/exporting it.
typedef struct {uint64_t records, target_records, path_count; int error;} CrumbLSMetrics;
int crumb_ls_bundle_paths(const void *,size_t,const char *bundle_id,
                          char paths[CRUMB_LS_MAX_PATHS][PATH_MAX],CrumbLSMetrics *);
