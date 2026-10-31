# iOS 27.2 security advisory and Crumb access research — 31 October 2026

## Sources
- Apple iOS 27.2 security advisory (released 26 October 2026): https://support.apple.com/en-us/149230
- CoreMedia 112706 public PoC: https://github.com/ceoca69/cve-112706
- Independent evidence and Wave 5 path matrix: https://github.com/MchinaMd57764/CVE-2026-112706-CoreMedia-fd-leak/blob/main/evidence/wave5_domain_map.txt
- FairPlayed kernel route: https://github.com/morzelowski/FairPlayed
- Crumb native failure diagnosis: `docs/FAIRPLAY_DIRECT_ATTEMPT_23057.md`.

## Shortlist from the Apple security advisory

| CVE | Component | Apple impact statement | Relevance to Crumb / certainty |
|---|---|---|---|
| CVE-2026-112706 | CoreMedia | App may access sensitive user data; permissions fixed | A real bytestream FD primitive with public PoC; demonstrated outside-sandbox coverage was limited to iTunes_Control/iTunes, not arbitrary bundles |
| CVE-2026-113896 | FileProvider | App may access sensitive user data; race condition | Potential file-access investigation, no verified public FairPlay or app-bundle capability |
| CVE-2026-112618 | Files | App may access sensitive user data; access issue | Generic access weakness, no verified public app-bundle/FairPlay capability |
| CVE-2026-112628 | Shortcuts | App may access sensitive user data; authorization issue | No demonstrated ability to write to an installed signed app bundle or grant FairPlay access |
| CVE-2026-113914 | App Store | App may access sensitive user data; authorization state | Authorization wording does not demonstrate ability to supply FairPlay license credentials |
| CVE-2026-112612 | Accounts | Physical attacker with previously trust-paired device may read/write arbitrary files | Excluded from on-device-only architecture; depends on trust-paired external device |

GitHub repository searches for the four candidate identifiers 113896, 112618, 112628 and 113914 returned no named repositories in the queried results on 2026-10-31. That is *not* proof that no implementation exists elsewhere. Apple advisory impact statements alone are not executable proof of capability.

## What is actually being tested

The `CoreMediaLocalProbe` component now tests two authorized, deliberately different locations, each against its own direct-open reference:

1. A previously downloaded `.ipa` inside **Crumb's own data container**.
2. **Crumb's own installed executable** from `Bundle.main.executableURL`, with exact identity and own-bundle checks.

Each attempt performs the same XPC bytestream handshake and compares only the first 32 bytes. It returns a stage, status, and match count; no filesystem path, content, or protected license data is exported. It cannot target another installed application or user content. The two probes are independent and are not used to alter native FairPlay processing.

This comparison tests the hypothesis that CoreMedia's bytestream broker applies different read authorization for files in an installed app bundle versus app data. It **does not** test `fairplayd` authorization because `fairplayd` is a separate process and its path/policy decisions are not implied by a descriptor delivered to Crumb.

The critical native failure remains: on an iPhone17,1/iOS 27.0.1, Crumb mapped its local IPA successfully and attempted `mremap_encrypted` but received -1/errno 12 and recovered zero pages. A passing CoreMedia probe, even for both local targets, does not resolve this.

## Release and test gates

- Mac-host tests enforce path allowlisting, own-executable identity and symlink rejection.
- iOS ARM64 typechecks the CoreMedia adapter.
- Full IPA build exercises UI integration and linkage.
- Device runtime results for both targets are required to conclude anything about CoreMedia behavior on 27.0.1.
- A separate full-page decrypted-output assertion remains necessary before declaring FairPlay success.
- No new CVE exploit chain, arbitrary-file extraction, license transformation, undocumented installation permissions or claim of a working decryptor is part of this PR.

### Recording a real iPhone trial

Capture only:
- iPhone model and iOS version;
- `Downloaded IPA` stage/detail/compared-byte count;
- `Crumb installed executable` stage/detail/compared-byte count;
- independently: existing FairPlay native stage, mapping/remap errno and pages changed.

Do not attach purchased IPA bytes, SINF payloads, Apple account credentials, private keys or raw filesystem paths to GitHub.
