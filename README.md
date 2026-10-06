# Crumb

Crumb is the branded continuation of this WaffleStore integration, with bundle ID `com.certlium.crumb`. The original Xcode project/target structure and contributor history are retained. New identity means a separate app and Keychain sandbox: sign in again; existing WaffleStore downloads/favourites are not imported automatically.

A **jailed** App Store app downloader and version selector, based on
[MuffinStoreJailed](https://github.com/mineek/MuffinStoreJailed-Public) and
[PancakeStore](https://github.com/jailbreakdotparty/PancakeStore).
Supports iOS 16.4+ and conventional sideloading without jailbreak, TrollStore or JIT.

Apple authentication and Store downloads are adapted from the current
[majd/ipatool](https://github.com/majd/ipatool) implementation. SAP guest code runs
in a static interpreter inside the app; the ipatool command-line executable is
not bundled. Passwords and 2FA codes are not persisted. Sessions and stable
machine identity use Keychain.

## Use

1. Sign Crumb with your preferred sideloading method and install it.
2. Sign in with your Apple ID. Temporary Apple failures are retried automatically;
   you can cancel. Enter an actual 2FA code only when Apple requests one.
3. Search or enter an App Store link, ID or bundle ID and choose a version.
   Existing account access is checked first. Verified free apps are obtained
   only when a license is missing or unavailable, then the selected version is requested.
4. Review its version number and choose **Download and install** to open the
   Safari installer after verification, or **Download IPA only** to see Install/
   Export when it finishes. Keep the installer open and confirm iOS's prompt.
5. **Downloaded apps** also offers installation, export, version details and
   **Delete IPA** in the options menu or by swiping (with confirmation), plus **Delete all downloads**. Deleting a download does not uninstall
   the installed app or remove its data.

The user confirmed authentication, versions, IPA export and installation on
an iPhone running iOS 27.0.1, signed with ksign/certificate without JIT. This is
not a guarantee for every app/device/iOS version; see [TESTING.md](TESTING.md).
Apple can still return intermittent HTML/empty login failures. iOS decides
whether App Store protection, signing, licenses and downgrades permit installation.
Crumb does not decrypt/re-sign an App Store package or report installation
success merely because Safari opens or the local server serves bytes.

The original OTA method uses **api.palera.in** for an HTTPS manifest containing
app name/bundle/build and a loopback URL. It does not upload the IPA or Apple
credentials. Export remains available if iOS refuses installation.
Automatic license acquisition is restricted to verified free apps.

## Credits and development

Original creators and contributors: mineek, lunginspector, skadz and nxtcoreee3.
[boierito](https://github.com/boierito) coordinated the revival and performed
on-device iOS 27 testing. [majd/ipatool](https://github.com/majd/ipatool), MIT,
provided the modern SAP/authentication/Store reference.

This contribution is **AI slop**: generated with OpenAI Codex, with device
validation reported by boierito. See [IMPLEMENTATION.md](IMPLEMENTATION.md) for
architecture/build instructions, [CHANGED_FILES.md](CHANGED_FILES.md) for the
file inventory and [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) for licenses.
Original attributions and history are preserved. Upstream redistribution rights
must be clarified before independently publishing a redistributed release.

## Build

Build `WaffleStore.xcodeproj`, scheme `WaffleStore`; the product is `Crumb.app`. GitHub Actions builds unsigned `Crumb-Debug.ipa` and `Crumb-Release.ipa` on macOS. Sign Release normally with ksign/SideStore/AltStore/Sideloadly. Release has no activity log or stdout capture. This branch is prepared for a future repository; independent public redistribution requires resolving upstream license permissions first.
