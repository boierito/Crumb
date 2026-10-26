# FairPlay 23038 — multiprobe on one iPhone IPA

Base: 23037 (`feature/fairplay-iphone-local-flow-23037`, commit `29fc14d3c7ed`).
Branch: `feature/fairplay-multi-probe-23038`.

## Why this build exists

The real 23037 device log repeatedly rejected recovery at
`bootstrap/raw_unknown_control_flow`, last translated entry
`text+0xb2f9d8`, after successful SINF and SUPF lookups.
A second identical manual retry did not add actionable evidence.
The import status zero and 880-byte legacy SIDB are **not** a portable
enrollment or a valid decrypt model.

The new build compares four bounded executions from a **single**
verified StoreAgent capture, all performed in the iOS application sandbox.
No new Apple login or IPA download is needed to run the diagnostic again
while the lab session remains open.

## Two new UI operations

- **Download + four-way diagnostic (iPhone only)**: download/capture once,
  audit it, then run the full matrix automatically.
- **Run four-way diagnostic (reuse current capture)**: repeat the same matrix
  against the original saved transaction without network traffic.

For a normal download/decryption attempt, the 23037 path remains available.

## Exact matrix

The private state contents, hardware identity, matched IPA/SINF and
recovery executable are identical in every case. Only the virtual
filesystem maps change:

1. `dual`: complete state under both fpsd and Shared SC Info aliases.
2. `fpsd-only`: state exposed only at `/var/db/fpsd/SC Info/`.
3. `shared-only`: state exposed only at `/Users/Shared/SC Info/`.
4. `sidb-only`: only the unchanged `SC Info.sidb` exposed at both locations
   (the lock and other auxiliary files are not exposed).

This tests path/auxiliary state assumptions. It is **not** a method of
conjuring new cryptographic inputs or weakening model/IPA validation.
On framework-heavy IPAs the experiment samples at most two matched
license domains, bounds execution by checking a 45-second time budget
before each attempt, and returns an explicit truncated flag if needed.

## Returned report

Native operation: `probe-capture`.
Schema: `crumb-fairplay-local-probes-v1`.

The result reports only per-attempt mode, named phase, fixed error code,
last publicly translated code offset, attempted untranslated target
**only when aligned and within the pinned public code window**, generic
out-of-window target classification otherwise, bounded lookup counts,
shared-block counts, source/frame verification flag and duration.

When recovery succeeds, the probe checks source/frame consistency in
memory but does NOT export a model, call IPA decryption, or claim correct
plaintext. Secret buffers are cleared. The existing legacy recovery and
full model/SINF/IPA gate remain separate.

The strict iPhone app sandbox and file protection rules are unchanged.
Never share the separate **PRIVATE capture** ZIP.

## Device test once GitHub IPA builds successfully

1. Install the unsigned 23038 IPA after signing.
2. Log in with your own authorized Apple Account and choose one small
   entitled, freely available test app.
3. Tap **Download + four-way diagnostic (iPhone only)**.
4. Once completed, tap **Prepare exports**, and send only the sanitized
   `diagnostic.json` after reviewing it.
5. Check whether `untranslated_target` is a public `text+0x...` offset,
   a generic out-of-profile classification, or absent. Compare the four
   variant phase/code/target tuples, not only durations.
6. If all variants end on the same untranslated target, focus next on the
   CFG translator; if they diverge, inspect the relevant runtime file
   lookup contracts before touching crypto. If one verifies source/frame,
   run the separate **Attempt recovery from current capture** and require
   full SINF/IPA validation; do not accept a model based on probe success
   alone.

No private data needs to be sent for this investigation.

## Build acceleration

The GitHub `Build IPA` workflow now has a `native` job that builds
Unicorn TCI and the Go/FairPlay library once per workflow run, restores
and saves an input-hashed cache, and shares that exact build's artifacts
with Debug and Release jobs. The original `ensure-native.sh` fingerprint
still guards consumers; if toolchain or source inputs differ, it rebuilds.
The workflow also cancels obsolete queued/running builds of the same
branch or PR.

Note: the first cold native build remains expensive; cache hits on
subsequent compatible runs provide the larger expected speed improvement.
Actual time savings must be measured from CI and are not yet demonstrated
on this branch.

## Remaining limitation

The original pinned modern model is not necessarily compatible with a
legacy iOS StoreAgent SIDB. `unfairplay` documents enrolled-mac
recovery for specific macOS builds, **not** verified iPhone-only portable
enrollment. A fully decrypted authenticated IPA has not been produced.
