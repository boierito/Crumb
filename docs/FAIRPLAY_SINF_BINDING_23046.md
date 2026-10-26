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
