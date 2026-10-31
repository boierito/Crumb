# Direct FairPlay processing of an IPA imported at runtime

## Requirement

Install Crumb once. After that, select any user-owned encrypted App Store IPA
from iOS Files and attempt to decrypt it **without embedding its payload in the
installed Crumb.app, rebuilding Crumb, re-signing it, or reinstalling it**.

"Embedded" in upstream FairPlayed means **installed bundle resources**, not
native code, a Swift constant, or executable linking. iOS seals the installed
bundle. Adding files to the app data container does not make them part of its
installed bundle. Confusing these storage classes obscured the real obstacle.

## Changes in this branch

1. `Downloads → Import encrypted IPA from Files` opens a standard iOS
   document picker in copy mode.
2. `RuntimeIPAImport` validates a regular `.ipa` file (4-byte ZIP marker,
   extension, size and symlink rejection), then copies it to encrypted-at-rest,
   backup-excluded private app storage. The source in Files is never modified.
   All preparation is performed on the iPhone, after installing Crumb.
3. The Decrypt view now shows the imported filename and a separate Decrypt
   button. This calls the **existing native FairPlayed-based** direct route,
   specifically `SystemFairPlaySession.decrypt(identifiers:localIPA:description:)`.
   The original encrypted IPA stays available if decryption fails.
4. Input restoration persists only the private relative path and presentation
   name; stale per-download exports are not attributed to an imported IPA.
5. The CoreMedia probe from the underlying research branch can also be run
   on the imported local IPA; it remains an independent diagnostic, and is
   **not** considered a FairPlay authorization result.
6. Host unit tests validate copy/non-overwrite, invalid format, symlink,
   directory and preservation boundaries. An iOS arm64 typecheck and full
   Xcode IPA build compile the new user interface and adapter.

## Verified versus unverified

**Implementation:** runtime import and connection to the real native
FairPlay attempt now exist. **Not verified:** actual decryption of an imported
IPA on a stock iPhone.

Upstream FairPlayed states that its kernel/system FairPlay path looks up
`SC_Info/<binary>.sinf` by the executable's vnode path and that
`fairplayd` is blocked from app data containers. Previous Crumb testing
on iOS 27.0.1 confirmed `mmap` succeeded but
`mremap_encrypted(cryptid=2)` returned -1, errno 12, with zero decrypted
pages for the direct data-container path.

This change deliberately **does not** sneak a target into the signed
application bundle or require a reinstall on the user's behalf.
It also does not pretend a CoreMedia file descriptor, a sandbox-scoped
document picker, a rewritten SINF or a changed `cryptid` grants the FairPlay
service separate filesystem permissions.

The original FairPlayed flow — one embedded target per signed helper
installation — was the working reference. Achieving the requested **one-time
Crumb installation and unrestricted later IPA imports** requires a new,
verified authorization or allowed storage path for `fairplayd`, or another
proven licensed decryption interface. None has yet been demonstrated.
Until then, failure reports and all-page checks are strictly preserved, and
no decrypted IPA is published unless the native/export gates pass.

## Device validation

1. Install a CI-built development Crumb IPA (no target IPA in the bundle).
2. From Downloads tap **Import encrypted IPA from Files**, select an encrypted
   IPA you have an Apple ID license for, and confirm the display name.
3. Tap Decrypt. It must either produce a fully verified decrypted export, or
   fail with an explicit numeric native report. No silent helper installation.
4. If it fails, share only the **Share error details** report, excluding private
   credentials/IPA/SINF payload. Note mapping and native remap stages, errno,
   and decrypted-page counts.
5. Compare with downloading the *same exact IPA* inside Crumb, keeping the
   account/device/OS fixed. If both fail at native remap, the import method
   is not the cause.

### Source references
- https://github.com/morzelowski/FairPlayed
- https://github.com/ceoca69/cve-112706
- Crumb `docs/FAIRPLAY_NO_REINSTALL_RESEARCH.md`
- Crumb `docs/FAIRPLAY_DIRECT_ATTEMPT_23057.md`
