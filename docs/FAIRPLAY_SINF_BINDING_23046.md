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
