# Crumb

<p align="center"><img src="branding/Crumb-AppIcon-1024.png" width="120" alt="Crumb icon"></p>

[![Build IPA](https://github.com/boierito/Crumb/actions/workflows/build-ipa.yml/badge.svg)](https://github.com/boierito/Crumb/actions/workflows/build-ipa.yml)
[![Protocol tests](https://github.com/boierito/Crumb/actions/workflows/sap-tests.yml/badge.svg)](https://github.com/boierito/Crumb/actions/workflows/sap-tests.yml)

Crumb is an iOS App Store downloader and version selector based on
[nxtcoreee3/WaffleStore](https://github.com/nxtcoreee3/WaffleStore).
It runs inside the normal iOS sandbox, without jailbreak, TrollStore or external JIT.
Bundle ID: **`com.certlium.crumb`**. Minimum iOS: **16.4**; development prioritizes iOS 27.

## Features

- Apple ID login and 2FA, with sessions and stable machine identity in Keychain.
- App search by name, App Store link/ID or bundle ID; favourites and five navigation tabs.
- Version selection by `externalVersionId`, with IPA version/build verification.
- Existing account license checks and acquisition attempts for verified free apps only.
- IPA downloads with progress, individual export, installation handoff and confirmed deletion/bulk cleanup.
- Release builds without activity logs, stdout capture or debug controls.

The modern SAP/authentication/Store flow is adapted from
[majd/ipatool](https://github.com/majd/ipatool). A statically linked interpreter
executes the guest SAP code; Crumb does not embed or launch the ipatool CLI.
Passwords and verification codes are never persisted. Proprietary Apple guest
assets are fetched at runtime and are not bundled in the repository or IPA.

## Install and use

1. Open a successful [Build IPA run](https://github.com/boierito/Crumb/actions/workflows/build-ipa.yml)
   and download its **Crumb-Release** artifact. Extract `Crumb-Release.ipa`.
2. Sign/install with ksign and a certificate, SideStore, AltStore, Sideloadly or a
   compatible conventional sideloader. Retain `com.certlium.crumb` where possible.
3. Sign in with your Apple ID. When Apple requests 2FA, enter its six-digit code
   and choose **Verify code**. **Request new code** starts a fresh signed challenge
   without retyping the password; Apple controls delivery. A 30-second cooldown
   prevents rapid repeats. A code can also be obtained in trusted-device account settings.
4. Search/select an app and version. Choose **Download and install** or download
   the IPA for later use. Confirm the iOS installation prompt and keep the installer open.
5. **Downloads** contains each IPA’s Install, Export and Details actions. Swipe or
   use its menu to delete one IPA; **Delete all downloads** removes downloaded
   files/records, without uninstalling apps or deleting their data.

Crumb installs separately from WaffleStore. Existing WaffleStore sessions,
favourites and downloaded IPAs are not automatically migrated. Updating Crumb
with the same signing identity/bundle ID preserves its own data and Keychain access.

## Current status and limits

Current source is Crumb **1.0.0, build 23020**. The originating WaffleStore
integration was tested by boierito on a physical iPhone running iOS 27.0.1 using
ksign/certificate without JIT, including login/2FA, versions, download/export and installation.
Crumb’s first-login/2FA recovery and new-code changes still require physical-device confirmation.
Automated tests use fixtures and do not demonstrate Apple’s acceptance of a live login.

- Apple can return intermittent empty/HTML authentication responses. Recovery is bounded and cancellable.
- Cross-region free-app acquisition remains unresolved: Apple can return **2059**.
  Catalog region selection does not promise a license or download from that region.
- No automatic paid purchases, FairPlay decryption or IPA re-signing.
- The original installer uses Safari/`itms-services` and **api.palera.in** for an HTTPS
  manifest containing app name/bundle/build and a loopback URL. Credentials and
  IPA bytes are not uploaded to that service. Serving bytes/opening Safari is not
  proof that iOS installed an app; export remains available if installation is rejected.
- iOS 26, iPad and complete app/data-retention compatibility remain to be validated.

See [TESTING.md](TESTING.md) for device scenarios and the verification matrix.

## Build

The original project structure is retained: `WaffleStore.xcodeproj`, scheme
`WaffleStore`; the product is **Crumb.app**. CI uses macOS 15, Xcode 26 and Go 1.25.1.

```sh
xcodebuild -project WaffleStore.xcodeproj -scheme WaffleStore \
  -configuration Release -sdk iphoneos -destination 'generic/platform=iOS' \
  -derivedDataPath build/DerivedData CODE_SIGNING_ALLOWED=NO build
```

GitHub Actions builds Debug/Release IPAs on push to main/feature/release/fix branches,
saves artifacts with corresponding native sources/notices, and supports manual runs.
`v*` tags create **draft prereleases**, not automatic public releases.
The protocol workflow runs Swift tests, safe file-deletion tests, Go race tests and
TCI checks with executable mappings denied. See [IMPLEMENTATION.md](IMPLEMENTATION.md).

## Credits, provenance and licenses

Crumb integration/coordination and iOS 27 device testing: [boierito](https://github.com/boierito).
Original contributors: mineek, lunginspector, skadz, jailbreak.party and nxtcoreee3,
including [MuffinStoreJailed](https://github.com/mineek/MuffinStoreJailed-Public),
[PancakeStore](https://github.com/jailbreakdotparty/PancakeStore) and WaffleStore.
Modern backend reference: [Majd Alfhaily / ipatool](https://github.com/majd/ipatool), MIT.

**AI slop:** developed with OpenAI Codex; the selected icon was generated with ImageGen.
Original history and attributions are preserved. See [branding](branding/README.md),
[CHANGED_FILES.md](CHANGED_FILES.md) and [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

WaffleStore’s original revision contains **no license granting independent redistribution**.
This repository does not invent that permission or relicense the entire project as MIT.
Upstream permission or replacement of unlicensed material remains necessary before
independent redistributed releases. Unicorn/TCI includes GPL-derived code; notices
and modified corresponding sources accompany build artifacts. Public source alone
is not a blanket redistribution grant. No proprietary Apple guest binaries are included.
