# FairPlay iPhone V9: verified checkpoint and precise SINF binding

Observed on 2026-10-26, unsigned Release IPA `test-23045-phase-checkpoint` signed locally by the owner. All fields below are sanitized diagnostics, not copies of any private capture.

## What the real run establishes

- Retained StoreAgent consumer: completed, status 0, legacy SIDB observed (873 bytes). This is **not** proof of modern enrollment compatibility.
- Raw execution: `generic-profile`, 13,683,433 instructions, completed, no JIT; nine state lookups with five misses.
- `checkpoint_compared_after_sinf=true`, `source_changed_after_sinf=true`, `frame_changed_after_sinf=false`.
- Source/frame checkpoint immediately after selection passes the **independent** `state.Evaluate` and all 14 `frame.Evaluate` comparisons: `state_evaluation_passed=true`, `source_frame_matched=true`.
- The pack length is 4184 bytes, identifying an `FPSENR2` candidate. A temporary candidate model is written with private mode 0600, then disposed of on downstream failure.
- **Failure**: `recover-capture/recovered-validation/model-sinf-binding`, `recovered_model_sinf_binding_failed`. `target_ipa_validated=false`. The original encrypted IPA was retained.

The phase checkpoint resolved a real coherence error. It has **not** established the correctness of the SINF AES key, suitability of the legacy SIDB, or any decrypted page.

## Next isolated change

The existing model/SINF provider error combined file permissions, model envelope, selector, SINF geometry, and cryptographic plaintext reconstruction into one safe but ambiguous stage. `fairplaymulti.go` now classifies an **original failure**, without retrying or choosing a different model:

| New nested stage | Meaning |
| --- | --- |
| `model-pack-input` | Cannot load a private model pack or directory |
| `model-sinf-structure` | The installed SINF selector cannot be parsed |
| `model-sinf-envelope` | Candidate model format or tables fail validation |
| `model-sinf-selector` | No unique pack matches the installed SINF selector |
| `model-sinf-decode` | Matching model fails native SINF plaintext validation |
| `model-sinf-resolution` | Some other resolver failure, still rejected |
| `model-sinf-binding` | Fallback stage for other binding errors |

All classifications are fixed labels. The existing resolver remains authoritative and unchanged. No selectors, private hashes, SINF payloads, derived keys, addresses, paths, or nested errors are emitted by `recover-capture`. Candidate packs and any processed IPA remain private and unpublished until target IPA validation succeeds.

The regression uses actual (synthetic) V3 SINF ciphertext and the *existing* provider in the local CI test harness, checking valid case, corrupted envelope, mismatched selector and corrupted key. This is **not** a simulated success for an actual App Store IPA.

## Evidence required next

The single classification from the next on-device attempt will determine whether the recovered `FPSENR2` pack is malformed, mismatched to the IPA selector, or cryptographically incoherent. On `model-sinf-decode`, work should focus on the SINF phase and provenance of the key material at `currentRawSourceSINFMaterialBytes`, and possible incompatibility of StoreAgent state with native `fairplayd`. Do not bypass `validateCurrentSINFV3Plaintext` or the legacy plaintext atom checks, and do not treat `target_ipa_validated=false` as success.

Do not publish private capture, SIDB, SINF, keybag, or model in GitHub or the chat.


## iPhone V10: reproduced specific model/SINF decode failure (2026-10-26)

Owner's new iPhone Release test used source commit `959f1238`:

- Actual retained state: 878-byte StoreAgent SIDB, 888-byte keybag, 196-byte kbsync, one 440-byte private SINF, consumer native status 0.
- `source_frame_matched=true`, 14 frame candidates checked, `state_evaluation_passed=true`, 224-byte source, 16-byte frame, 4184-byte `FPSENR2` candidate.
- `recovered_model_sinf_decode_failed`, stage `model-sinf-decode`. This means the installed SINF parsed correctly, the model envelope and selector were accepted, **but the model failed the decoder's plaintext checks**.
- `target_ipa_validated=false`; no decryption or authenticated pages. The same missing 5 out of 9 runtime state lookups occurred as in the preceding checkpoint run.

### Root-cause separation, without inventing private inputs

The native 26A428 enrollment engine retrieves 16 bytes of SINF key material from its current stack frame after running the SINF CFG. A compatible `FPSENR2` model decrypts 440-byte legacy-format `priv` via AES-CBC, decoding the expected atom sequence. Failure of that plaintext check does not independently identify which of these is responsible:

1. Legacy StoreAgent 10.9 SIDB data format/enrollment is incompatible with the modern macOS 26A428 enrollment runtime;
2. Modern state files required by this execution are absent or not mounted correctly;
3. Runtime interpretation or stack material extraction differs from the pinned native path;
4. The SINF is genuinely personalized to a different device/enrollment than the state used to derive this candidate.

The fact that the source/frame pass is *internal self-consistency* of the model's processing. It does not authenticate a modern enrollment or prove the SINF decryption key.

### Bounded investigation committed after the iPhone run

The rawsource runtime now splits its existing **aggregate** state-file lookups into a finite, sanitized `state_file_lookups` histogram: `sidb`, `sidd`, `sidf`, and `other`. Each holds only attempts and missing counts. No filename, path, file data, digest, key material, or user identity is retained or displayed. The original guest-visible file resolver, SIDB input, crypto routines, and all failure gates are **unchanged**.

This isolates a precise question: were the five failed lookups references to optional sidecar names, or to the main SIDB required for this profile? If all five misses are optional, the root hypothesis shifts toward state compatibility or CPU/key material semantics, **not** a missing SIDB path. This histogram is evidence of file *resolution*, not proof that guest code cryptographically consumed the material.

Do not repeat identical captures for unchanged binaries. Do not infer an iPhone-only decryption implementation from this instrumented runtime.


## iPhone V11: state-file routes isolated (2026-10-26, `test-23045-state-lookups`)

The device's genuine StoreAgent consumer returned an 869-byte legacy SIDB, native status zero and a 904-byte keybag. The experiment again reached the selection and SINF phases, produced a 4184-byte `FPSENR2` candidate and failed the **real** resolver: `recovered_model_sinf_decode_failed`, `target_ipa_validated=false`.

The safe runtime counters now narrow nine state-related lookups to:

| Type | Attempts | Missing |
| --- | ---: | ---: |
| `sidb` | 4 | 0 |
| `sidd` | 1 | 1 |
| `other` | 4 | 4 |

This **disproves the missing-SIDB-file-path hypothesis for these virtual routes**. It does NOT prove that the legacy SIDB is semantically understood by modern 26A428 fairplayd, and the five other missing lookups should not be labeled optional without a reference. Nor does a zero native StoreAgent status prove any model decryption.

### Probe correction

The previously shipped four-mode probe (dual paths, fpsd only, shared only, SIDB alone) checked only internal source/frame coherence; this is insufficient because all recent captures passed that check yet **failed** SINF decoding. `fairplayprobe.go` now calls the existing `fairplay.ResolveLogicalKeyWithModel` against the exact per-domain installed SINF and freshly recovered candidate for *each* routing mode. The underlying model and SINF bytes remain unchanged, and the strict resolver remains authoritative.

Each probe reports separately:
- `source_frame_verified`: original internal consistency check;
- `sinf_model_verified`: **only** true if the original SINF resolver really decodes and validates plaintext;
- a fixed label `sinf_model_validated`, `sinf_model_rejected`, or `sinf_license_not_unique`, never the private resolver cause, keys, selectors, paths or hashes.

The probe never runs IPA page decryption and never publishes output. A positive SINF probe would require a full IPA validation in the standard pipeline before any decryption success claim. Failure of all probe routes indicates virtual routing is not sufficient to repair the cryptographic model, and the next priority becomes the trust/protocol compatibility of legacy StoreAgent state and hardware identity.

Security: four bounded modes; use original owner capture only; no guessed state fields, synthetic keys, network export or cryptographic bypass.


## iPhone V12: complete four-route SINF probe (2026-10-26)

The signed owner iPhone tested the verified release test-23045-fourway-sinf-probe, using the same real StoreAgent transaction throughout four bounded diagnostic modes.

| Virtual state route | Source/frame | Strict model/SINF decryption | Code |
| --- | --- | --- | --- |
| dual | verified | rejected | sinf_model_rejected |
| fpsd only | verified | rejected | sinf_model_rejected |
| shared only | not reached (bootstrap memory unavailable) | not reached | raw_memory_unavailable |
| SIDB alone in both aliases | verified | rejected | sinf_model_rejected |

The original default recovery independently returned recovered_model_sinf_decode_failed and target_ipa_validated=false. The matching keybag was 904 bytes; the imported StoreAgent SIDB 863 bytes. The runtime resolved four of four SIDB path requests, missed one SIDD and four other state requests. No SINF plaintext or target IPA was validated.

**Conclusion:** the four tested virtual mount mappings and removal of optional guest sidecars did not repair SINF decryption. Source/frame coherence is an *internal* test, not proof of the correct SINF key. Repeating this four-way probe on the unchanged pipeline adds no useful information.

**Engineering pivot:** test a *genuine* StoreAgent 10.9-compatible license/key derivation path or implement the missing modern device enrollment producer; do not infer equivalence from legacy SIDB presence or native status zero. The existing 26A428 path has independently validated native Mac controls but has **not** been proven compatible with a synthetic/legacy StoreAgent virtual enrollment. No guessed key, alternate selector, skipped atom validation, or exported unverified decrypted IPA is acceptable.

Further work on the emulator should be driven by a reproducible difference between native and translated key derivation on a controlled authenticated enrollment. Keep raw SIDB, SINF, identity, key material and any private hashes entirely off public logs/repositories.
