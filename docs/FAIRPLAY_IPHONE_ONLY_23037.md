# Crumb FairPlay — iPhone-only implementation, iteration 23037

**Base:** `fix/fairplay-consumer-io` at `fdd8ac278e8344ef56250234021087bf348bf275` (23036).  
**Branch:** `feature/fairplay-iphone-local-flow-23037`.  
**Scope:** one-device workflow and fail-closed publication gate; no Mac, PC, jailbreak or decryption server is required to attempt this path.

## 1. What changed

The FairPlay laboratory now offers **Download + attempt local recovery on iPhone**. This action:

1. Downloads an eligible IPA using the existing Crumb account/Store transaction code.
2. Keeps the matched KBSync producer alive while the associated keybag is consumed by the StoreAgent emulator.
3. Saves the existing private, transaction-bound capture (not raw bytes in the public log).
4. Releases the emulated StoreAgent before the IPA CDN transfer, then packages the original IPA.
5. If **exactly one** new transaction was retained, exports and audits its capture, then invokes `recover-capture` on the device. Alternative download routes with no capture are skipped.
6. Logs the sanitized native response. A failure or missing material leaves the original IPA available and does not claim success.
7. Exposes `processed.ipa` only when the native bridge returns a complete recovery result, all required model/frame/SINF gates are explicitly true, and the output file has nonzero size.
8. Removes temporary recovery model files and reports. Unsuccessful output folders are deleted.

The original separate **Download and capture transaction** and **Attempt recovery from current capture** actions remain available for reproducible debugging. The new path simply reuses their shared implementation, preventing divergent manual and automatic semantics.

## 2. New fail-closed result contract

`WaffleStore/App/FairPlayLocalRecoveryGate.swift` requires all the following for an iPhone-capture-derived output:

- `operation = recover-capture`, `ok = true`, and a valid `result` dictionary.
- `completed = true`, `portable_source_derived = true`, and `enrollment_frame_derived = true`.
- `source_frame_cross_checked = true` and `model_sinf_validated = true`.
- `images > 0`, `pages > 0`, and a nonempty `processed.ipa`.
- `plaintext_authenticated = false` as an *explicit limitation*: this build has no independent plaintext authentication mechanism.

`ok=true` alone is **never sufficient**. Rejected or malformed native responses cannot produce a shareable processed IPA through this flow. This is a publication gate, **not** a proof that transformed executable pages are cryptographically correct.

## 3. Tests and CI

`scripts/test-fairplay-local-recovery-gate.swift` is run in the macOS **SAP protocol tests** workflow:

```sh
swiftc WaffleStore/App/FairPlayLocalRecoveryGate.swift \
  scripts/test-fairplay-local-recovery-gate.swift \
  -o /tmp/crumb-fairplay-local-gate
/tmp/crumb-fairplay-local-gate
```

It covers a complete synthetic envelope, missing/false recovery flags, zero processed pages/images, absent output, native error and wrong operation. It does not use real credentials, SIDB, SINF, models or plaintext. The normal iOS Debug/Release IPA workflow is unchanged and must also pass.

## 4. Device acceptance procedure

1. Install a newly built and signed Crumb IPA on an iPhone. GitHub Actions produces **unsigned** artifacts.
2. Authenticate with an account entitled to a small free test app.
3. In Settings > FairPlay Laboratory, run synthetic self-tests.
4. Enter the app ID and select **Download + attempt local recovery on iPhone**.
5. Inspect the sanitized `Transaction`, `Enrollment` and `Decryption` logs. Check exact recovery phase and code (especially `bootstrap/raw_unknown_control_flow`).
6. If recovery fails, verify that the original IPA remains available and no processed IPA can be exported.
7. If recovery succeeds, verify every readiness flag and distinguish **model consistency** from **independent plaintext authenticity**.
8. Repeat with another authorized app, including multi-image / shared-license layouts when available, without sharing any private captures.

**Unverified:** The 23036 shared-CFG change has not been demonstrated to pass the actual `bootstrap` barrier. This iteration does not import the newer `DrWhax/unfairplay` compact CFG implementation and does **not** claim to have decrypted a real IPA.

## 5. Next implementation: advance past bootstrap

The next data-driven change must be based on a real *sanitized* 23037 device report. If `raw_unknown_control_flow` persists, map its last translated public block label to the pinned CFG generator and the current UnfairPlay compact translator. Port **verified translated operations** and add a public synthetic regression test. Do not invent missing operations, slice arbitrary SIDB bytes into frames, or downgrade readiness checks.

If bootstrap passes but another phase fails, use the named phase and bounded file lookup counts before modifying state mappings. If a selector-bound model and every page check succeed, test the output against independently known plaintext before claiming complete decryption.

## 6. Security

- Keep `context.json`, `consumer-state.json`, keybag, KBSync, private state and enrollment models in the iPhone sandbox.
- Keep public diagnostics limited to fixed stage codes, bounded counts, sizes and readiness flags.
- Preserve `FileProtectionType.complete` and restricted file permissions.
- Never upload private captures to GitHub or distribute decrypted apps without the appropriate rights.
- Keep original IPA and recovered output separate.
