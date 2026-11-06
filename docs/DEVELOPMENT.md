# Development guide

The application project is `Crumb.xcodeproj`, with the shared `Crumb` scheme.
`README.md` describes the product; `IMPLEMENTATION.md` and `TESTING.md` retain
backend details and historical validation. `docs/ON_DEVICE_DECRYPTION.md` records
the on-device workflow and signing changes.

## Code layout

- `Crumb/Device/DeviceWorkflow.swift`: workflow state, task lifecycle, helper
  handoff, cancellation, result retrieval and installation orchestration.
- `Crumb/Device/SigningModels.swift`: value models shared by the workflow and UI.
- `Crumb/Device/SigningViews.swift`: configuration, signing review and report UI.
- `NativeDevice/src/resign.rs`: profile loading, signing plans, bundle signing,
  signed-artifact packaging and cleanup.
- `NativeDevice/src/resign/identifiers.rs`: identifier validation, profile scope
  matching and boundary-aware bundle-ID remapping.
- `NativeDevice/src/resign/entitlements.rs`: profile-authorized capabilities,
  identity changes and CloudKit report values.
- `NativeDevice/src/resign/tests.rs`: pure identifier/capability regressions.
- `NativeDevice/src/resign_tests.rs`: generated certificate/profile fixtures,
  actual CMS signing, resource preservation and cancellation integration tests.
- `NativeDevice/src/signing.rs`: certificate/profile inspection, existing
  entitlement decoding and profile permission matching.
- `DecryptHelper/`: helper UI, job isolation, image processing and IPA resources.

## Behavior to preserve

Use the same saved `.p12`, `.mobileprovision` and password for helper and IPA
signing. Preserve existing Keychain accounts and download sidecar schemas.
Keep original decrypted IPAs; install/export a separate signed artifact.
Keep resources and nested extensions, validate encryption state, and preserve
executable permissions. Cancellation must reset the UI without letting a retired
job alter its replacement. Connection checks must remain event-driven.

Profile permission shapes matter: scalar wildcards can authorize lists of
strings, and iCloud environment grants can list alternatives to a scalar app
selection. Empty container lists do not authorize a container. Preserve existing
authorized cloud identifiers before considering reviewed namespace changes.

The build 23039 refactor relocates existing Rust functions/tests and Swift value
models and formats signing views. It preserves their behavior, UI text, FFI
operations, persisted settings and signing rules from build 23038.

## Validation

```sh
cargo test --locked --manifest-path NativeDevice/Cargo.toml
cargo test --locked --release --manifest-path NativeDevice/Cargo.toml
cargo test --locked --manifest-path NativeDevice/Cargo.toml -p apple-codesign-quick --lib
python3 scripts/verify-crumb-branding.py
python3 scripts/test-device-macho.py
python3 scripts/test-device-file-access.py
python3 scripts/test-file-picker-injection.py
```

`Build IPA` additionally compiles against the iOS SDK, tests helper cancellation
and IPA resources on macOS, runs the Swift persistence/connection checks, builds
Debug and Release, and verifies each IPA. See `TESTING.md` for backend checks and
physical-device validation; host fixtures do not establish Apple acceptance.

## Main integration and rollback

`backup/main-before-decrypter-2026-11-06` preserves the pre-integration main commit
`ae40ff1d114a19c5b5c7543425fa0e6a5281527a`. Keep this branch and the development
branches when integrating. Use a merge commit to retain the feature history.
To undo the integration on shared main, revert the merge commit in a new PR
rather than rewriting main or removing the backup branch.

## Decrypted export integrity

`NativeDevice/src/export_audit.rs` compares every original archive entry with the
retrieved result, before saving it as a completed download. Unencrypted resources
must match exactly. Encrypted Mach-O files may differ only in the declared
ciphertext ranges and their cryptid fields, which must become zero. Signature
and entitlement bytes outside those ranges remain identical. Missing files,
empty directories, changed file modes and unexpected additions fail validation.
Only exact SINF replicas explicitly requested by the original manifest are allowed.

The helper's ZIP writer retains Payload directories (including empty ones) and
uses preserved file permissions instead of assigning executable mode to all files.
Keep the original encrypted download until result verification finishes.

Downloads → Export offers the original decrypted IPA and a separate signed IPA.
The signed choice uses the same profile-aware pipeline as Install, without opening
installation automatically. Its original source IPA is never rewritten. Signing
again in another application can change entitlements and is outside this audit;
the original Apple CMS is naturally invalid after decrypted code bytes change.
