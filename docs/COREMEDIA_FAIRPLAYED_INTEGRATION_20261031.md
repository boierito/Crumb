# CoreMedia + FairPlayed integration research (2026-10-31)

## Sources and exact scope

- FairPlayed: https://github.com/morzelowski/FairPlayed (read-only `mremap_encrypted`, cryptid 2).
- CoreMedia disclosure: https://github.com/ceoca69/cve-112706 (fork of https://github.com/MchinaMd57764/CVE-2026-112706-CoreMedia-fd-leak).
- Local base: `fix/fairplay-native-failure-contract` (PR #30). Existing `CrumbSystemDecryptRange` and the strict all-page export checks are **preserved without modifications**.

This change integrates CoreMedia's bytestream session protocol as a **separate, explicitly requested file-access diagnostic**. In the download's Decrypt screen, `Test CoreMedia file access` checks only a regular downloaded IPA belonging to Crumb's own app data container. It does not test arbitrary system paths, dump user files, use the service to access installed apps, or change FairPlay license material. Success requires matching the first 32 bytes from the daemon-vended file descriptor with the same IPA opened directly from Crumb. No file data is retained or exported, and the diagnostic returns only stage, coarse error and byte count.

The comparison is only a functional test of the CoreMedia request, **not** evidence of FairPlay license authorization, `fairplayd` access to Crumb's data container, or successful decryption. It does not feed a file descriptor into `mremap_encrypted` because the existing kernel route resolves the backing vnode and filesystem path separately. Treat these as separate permissions.

### Why not copy the full CVE PoC?

The disclosed CoreMedia primitive is a permissions vulnerability and has a limited demonstrated readable domain. The original 27-target matrix did not show arbitrary access to app bundles or other app containers. This integration does not expand its reach or include a general-purpose arbitrary-path file reader. On a device with a fixed OS release, the service may also reject the request.

For the existing FairPlayed pipeline, the installed-resource test previously recovered 178 encrypted pages from a licensed reference app on iOS 27.0.1. The one-tap `data_container` route still fails: the updated 23058 diagnostic shows `mmap` succeeded but `mremap_encrypted` returned -1/errno 12. A CoreMedia descriptor does not itself fix that kernel/service failure.

## Code

- `WaffleStore/App/CoreMediaLocalProbe.swift`: CoreMedia/FIG XPC handshake, bytestream descriptor request and a strict local-IPA-only comparison. No private bytes or absolute paths in returned summary.
- `WaffleStore/App/SystemFairPlayView.swift`: explicit per-downloaded-IPA diagnostic UI. Never treats CoreMedia results as decrypted output.
- `scripts/test-coremedia-local-probe.swift`: path boundary tests: own IPA accepted; outside container, remote URL, symlink escape, missing input and wrong extension all rejected; host call fails closed.
- `.github/workflows/sap-tests.yml`: macOS policy tests and iOS ARM64 typecheck. Standard IPA build workflow also validates native linkage/UI compilation via the PR trigger.

## How to test

1. Install a CI-built **development** IPA on an owned test device using an ordinary sideloader.
2. Download an IPA into Crumb, open its Decrypt view and choose `Test CoreMedia file access`. Copy the stage/detail summary, *not* private IPA/license content.
3. Run the original Decrypt action separately. Keep the existing `native_remap`, `native_mapping_errno`, `native_remap_errno` and page counters as the only evidence of that route.
4. Success criteria for the **diagnostic**: `matched`, 32 compared bytes. Success criteria for **decryption** remain all encrypted ARM64 slices and pages recovered and independently validated, producing a fully verified export. Do not infer the latter from the former.

## Limitations and release gate

- No device runtime test is possible in GitHub Actions. CI completion establishes compile/test status only.
- CoreMedia diagnostics do not rewrite the IPA, SINF, payload, or provisioning identity.
- No automatic fallback from a failed FairPlay read to CoreMedia is implemented, because no validated license-access bridge exists.
- Not a release-ready one-tap decryptor. Keep this as a draft research change until an actual iPhone report confirms the probe and a separate device report demonstrates complete decryption.
