# Crumb

<p align="center"><img src="branding/Crumb-AppIcon-1024.png" width="120" alt="Crumb icon"></p>

[![Build IPA](https://github.com/boierito/Crumb/actions/workflows/build-ipa.yml/badge.svg)](https://github.com/boierito/Crumb/actions/workflows/build-ipa.yml)
[![Protocol tests](https://github.com/boierito/Crumb/actions/workflows/sap-tests.yml/badge.svg)](https://github.com/boierito/Crumb/actions/workflows/sap-tests.yml)

Crumb lets you find App Store apps, choose a version and download its IPA on your
iPhone or iPad. It is based on WaffleStore and runs as a normally sideloaded app,
without jailbreak, TrollStore or external JIT.

## Features

- Sign in with your Apple ID and two-factor authentication. Sessions are stored
  in Keychain; passwords and verification codes are not saved.
- Search by app name, App Store link, app ID or bundle ID, and save favourites.
- Browse available versions and download the latest release or an older version.
- Obtain free-app licenses through your account when Apple permits it.
- Track downloads, check the downloaded version and build, and install or export
  individual IPAs through the share sheet.
- Manage downloaded IPAs and delete them individually or all at once.

App availability, licenses and installation remain subject to Apple's account
and device restrictions. Browsing another region does not guarantee a download.

## Installation

Download **Crumb-Release.ipa** from the [latest release](https://github.com/boierito/Crumb/releases/latest)
and sign it with SideStore, AltStore, Sideloadly, ksign with a certificate or another
compatible sideloader. Requires **iOS 16.4 or later**.

Sign in, find an app and choose the version you want. Downloaded IPAs are available
in **Downloads**, where you can install, export or delete them.

## Origins and credits

Crumb builds on [WaffleStore](https://github.com/nxtcoreee3/WaffleStore), whose roots
include [MuffinStore Jailed](https://github.com/mineek/MuffinStoreJailed-Public) and
[PancakeStore](https://github.com/jailbreakdotparty/PancakeStore). Its modern Apple
authentication and download backend adapts logic from
[ipatool](https://github.com/majd/ipatool) to run inside the iOS sandbox.

- **mineek, lunginspector, skadz, jailbreak.party and nxtcoreee3** — original projects
  and contributions that Crumb is based on.
- **[boierito](https://github.com/boierito)** — Crumb coordination and device testing.
- **[Majd Alfhaily / ipatool](https://github.com/majd/ipatool)** — modern backend reference, MIT licensed.
- **OpenAI Codex and ImageGen** — AI-assisted implementation and app icon. **AI slop.**

For build instructions and technical details, see [IMPLEMENTATION.md](IMPLEMENTATION.md)
and [TESTING.md](TESTING.md). Licensing and third-party attributions are documented
in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md); the project as a whole is not MIT licensed.

## Experimental IPA decryption

In Downloads, open an app’s options and choose **Decrypt IPA…**. Import a private
`fairplay-enrollment.bin` model generated with [DrWhax/ipatool](https://github.com/DrWhax/ipatool#fairplay-decryption)
on a supported Mac. The model must match the account and SINF selector of the
personalized IPA. Only supported modern enrollments from new accounts work;
older or migrated profiles are unsupported. Crumb does not recover enrollment keys.

Processing runs offline on the device. The original IPA is preserved. Export the
result before closing the sheet; temporary outputs and the imported model copy
are removed when finished. The exported IPA must be signed before installation.
The decoder is experimental; on-device compatibility still needs validation.
