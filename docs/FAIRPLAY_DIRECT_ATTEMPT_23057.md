# First returned direct-runtime report, build 23057

The user's safe report contains one `data_container` attempt: return -1, errno
12, zero pages, and zero usable shared locations. It requests neither signing
nor installation. The input IPA/transaction is not identified in that report;
do not assume it was DD or correlate it to DD's prior successful embedded run.

## What it proves

- No page was accepted as decrypted and no verified output was produced.
- No shared-container native attempt ran. A zero discovery count cannot be
  interpreted as a failed decryption in a shared container, or as proof of
  exactly which entitlement/container-discovery operation was unavailable.
- The aggregate `preparation_or_validation` label was misleading: the report's
  actual attempt came from the native pipeline. It is not another observed
  portable-recovery -42057 result.

## What it does not prove

The 23057 C reader initializes `remap_result` to -1, and also writes
`remap_errno` if `mmap` fails before the syscall. Therefore this v1 report
cannot distinguish memory-mapping failure from native remapping failure.
The first interpretation that mremap had definitely run was too strong.

Even **if** mremap ran, errno 12 does not identify the FairPlay service cause.
The actual public iPhone17,1 iOS 27.0.1 / 24A446 kernel's syscall table resolves
entry 489 to `0xfffffff00b0d8a34`. At `0xfffffff00b0d8ca8` through
`0xfffffff00b0d8cb8`, its error branch returns EPERM for `0xe00002c1`, and
ENOMEM for other nonzero crypter-creation returns. Generic error, true no-memory
and not-permitted cases therefore collapse to the same errno.

`scripts/audit-fairplay-ios-errno.py --kernel <extracted kernel>` verifies the
pinned original kernel bytes and executes those five instructions independently
in Unicorn for four modeled callback returns. It does not measure the callback
return on the user's phone or prove which operation failed there.

Kernel SHA-256:
`fa07057bf34f33af1ec26681ad0c87c13e5aaa548a300740eddb0e92876e5f95`.
Acquisition: the official Apple IPSW and extraction described in
`FAIRPLAY_NO_REINSTALL_RESEARCH.md`. No private input is used by the audit.

## Implementation correction

The reader now records mapping-attempted, mapping-errno and remap-attempted
separately. The Swift stage classifier preserves legacy reports as ambiguous,
distinguishes mapping, remap and page-read failures, and marks aggregate native
attempts as such. Reports use schema v2; secrets/paths remain excluded. No
execution budget, cryptid policy, license data or decryption-success gate changes.

This corrects the observation contract; **it does not fix the native rejection**.
No new IPA is requested solely to obtain more diagnostic fields.

## Signature registration review

Crumb's direct reader currently does not invoke F_ADDFILESIGS_RETURN. Upstream
FairPlayed's `fp_probe_file` tries both registration modes for each mapping mode.
Consequently adding that call alone is not established as a causal fix for the
reported data-container restriction. No manufactured signature, re-signed input
or speculative registration fallback is presented as resolving this report.

The compatible explanation remains a native loader/service failure, potentially
including the upstream-reported file-access restriction. The exact cause is
unresolved; this report alone does not prove insufficient RAM, bad Apple Account
authentication, an invalid SINF, or a universal impossibility of other routes.
