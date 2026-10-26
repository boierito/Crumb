# Crumb 23039 — iOS 27.0.1 device diagnostic

Observed on-device diagnostic (2026-10-26). Values below are non-secret readiness counts and public code offsets. **No private capture, SIDB bytes, keybag, Apple account ID, or app-specific entitlement information was collected or committed.**

## Previous vs current behavior

| Field | 23038 | 23039 |
|---|---|---|
| `phase` | `bootstrap` | `bootstrap` |
| `code` | `raw_unknown_control_flow` | `raw_unknown_control_flow` |
| `last_translated_block` | `text+0xb2f9d8` | `text+0x8fff2c` |
| `untranslated_target` | `text+0x8fff2c` | `text+0xa4be38` |
| `shared_block_executions` | 0 | 1 |
| `source_frame_verified` | false | false |
| `plaintext_authenticated` | false | false |

All four modes (dual, fpsd-only, shared-only, sidb-only) reached the same stop at `text+0xa4be38`. The tests finished in a bounded time: the four-way probe in 282 ms, with `truncated_at_time_budget: false`.

## Enrollment state, without secret material

- Native keybag consumer finished, `native_status: 0`, import succeeded.
- SIDB observed: nonempty, 863 bytes (variable size; **not** a frame, not proof of a model).
- SINF count 1; private field length 440; SINF family `non-key5` (not a statement about the Apple Account's age).
- `missing_stages: ["raw_state_orchestration", "validated_enrollment_model"]`.
- Original IPA retained; **no processed IPA or decrypted plaintext published**.
- No SINF or SUPF lookups were missing.
- Native bridge selftests for state/frame/SIDB passed, but these are synthetic vectors.

## Interpretation

The implementation exercised the additional verified bootstrap block added to 23039. Progress from one unmapped target to a later one is evidence the coverage extension is exercised, **not evidence that downstream cryptographic state is correct**. The next unresolved target is an additional missing branch, not a JIT entitlement failure or evidence that an old Apple Account caused the error.

Do not alter `ready_to_decrypt`, synthesize missing enrollment frames, bypass source/frame validation, or publish a processed IPA on this basis. Current runtime reports `ready_to_decrypt=false`.

## CI reference

Source SHA `a6fc705ef30e2c17bd2f0c2d1a3870628e4545f9`.
- Build IPA: https://github.com/boierito/Crumb/actions/runs/40264574258 — completed successfully.
- SAP protocol and portable regression tests: https://github.com/boierito/Crumb/actions/runs/40264574222 — completed successfully.

## Next checks

1. Independently validate the next code destination's provenance and semantics against the exact pinned binary. Do not guess branch targets or accept unmapped addresses as success.
2. Maintain bounded execution and sanitized, non-secret diagnostics.
3. Repeat device testing only when there is a reviewed functional change. Successful compilation alone is not successful decryption.
