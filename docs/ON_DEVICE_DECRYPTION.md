# Experimental on-device pipeline

This branch adds code, not physical-device acceptance. No licensed IPA has been
decrypted in this cloud session. Linux tests cannot exercise Apple's FairPlay,
iOS Settings pairing, certificate acceptance or the compiled UIKit helper.

## User flow

1. Settings → Device connection → Start pairing, then select Crumb in Settings →
   Privacy & Security → Developer Mode and enter the PIN. Approve Local Network.
   Pairing is exposed only on iOS 27.0 and 27.0.1 in this first build.
2. Enable LocalDevVPN and Check authenticated connection. The backend verifies
   the remote pairing and reads a random marker from this installation's own
   Documents through House Arrest before enabling operations. Interrupted remote helper uploads are removed from
   Crumb’s own PublicStaging subdirectory on the next installation attempt.
3. On iOS 27.0.1 only, import an RSA .p12 signing identity and its matching
   provisioning profile. The profile must permit the chosen helper bundle ID,
   list its signing certificate, be current and permit this device. No signing
   certificate is revoked, registered or obtained from Apple's portal here.
4. Download a licensed IPA with Crumb. In Downloads, open Decrypt on this device,
   enter the certificate password and prepare/install the helper. Crumb's Apple
   session is not a developer-portal signing session. The certificate/profile
   imports are necessary for this implementation, even with a free Apple ID.
5. Open Crumb Decrypt, trust its development identity if iOS requires it, and
   leave it in the foreground. Return to Crumb and Retrieve result when finished.
   The first version uses explicit foreground handoff rather than assuming iOS
   permits an unattended background decryption.
6. Export the verified output. Output must be re-signed before normal sideloading;
   this pipeline does not install the final decrypted third-party app or provision
   its extensions. Normal encrypted IPA installation/export remains available.

The helper occupies an app slot under free provisioning. It is a separate app,
not an update or replacement of Crumb. A new job replaces the same helper ID.
Deleting the helper removes its embedded target and result. Removing local
pairing keys does not remove the system's trusted-host entry. Signing imports
can be removed in Device connection.

## Why an installed helper

[FairPlayed](https://github.com/morzelowski/FairPlayed) documents successful
`mremap_encrypted` with CRYPTID_MODEL_ENCRYPTION on iOS 27.0.1. It also documents
that fairplayd cannot read the downloaded target in an app data container.
Crumb therefore embeds its target, including SINF material, as immutable helper
resources. The helper asks the operating system to decrypt it without mapping
executable memory. No decrypted bytes are invented and cryptid is not cleared on
a failed remap, page read or verification. The implementation is original; the
FairPlayed source is not vendored.

The device must have an Apple account licensed for the target. Pairing neither
supplies a license nor enables arbitrary reads of other apps. Apple's rejection
or a system-policy change remains a real failure. Version restrictions are
experimental entry conditions, not a claim that Crumb passed a physical test.

## Verification boundaries

- Native jobs have cancellation and a 600-second outer deadline. Pairing has a
  120-second per-connection timeout. Swift stops publication on cancellation and
  holds the worker until completion before deleting the signing workspace.
- Imported signing files and pairing credentials use a separate, non-synchronizing
  Keychain namespace. Password input is memory-only. Prepared files use protected
  Application Support; certificate staging copies are removed after preparation.
- Pinned idevice 89702099f4d7552f704f93b4526756abda319ceb, no dependency tracing
  enabled (the upstream parser has a debug record dump).
- Archive extraction rejects traversal, links, duplicate paths, size overflow and
  excessive entry counts. First version rejects FAT/32-bit images. Limits are 4 GB
  expanded and 512 MiB per output file. No partial IPA is reported as completed.
- Decryption uses a read-only map and cryptid=2, reads every encrypted page with
  vm_read_overwrite, requires changed bytes, lower plaintext entropy and
  additional arm64 return/authentication instructions. These are conservative
  heuristics, not a proof of application behavior. Sparse encrypted data or small
  images can be rejected even with a valid license.
- The helper patches cryptid only after those checks and hashes every output ZIP
  file. Retrieval checks the job UUID, full file hashes, CRC and absence of
  encrypted Mach-O commands before atomically publishing the output.
- There is no live-device test, execution test, final-app signing flow, automatic
  Apple developer login or blanket support for future iOS releases in this branch.

## Build and checks

Linux: `cargo test --locked --manifest-path NativeDevice/Cargo.toml` and
`python3 scripts/test-device-macho.py`. The vendored signer's tests include
preservation of embedded encrypted dylibs. macOS/Xcode uses the existing Build
IPA workflow, which compiles the Rust aarch64 static library, Swift and a UIKit
helper template. The template is stored in an opaque ZIP to avoid signing it as
an extension when Crumb itself is sideloaded. Corresponding sources are uploaded
with CI artifacts for the LGPL signer relinking requirement.

Physical acceptance must establish PIN pairing, same-installation verification,
valid-profile installation, actual licensed FairPlay page reads, failed-license
behavior, main binary/framework/extension completion, retrieval and independent
cryptid inspection. A device result is the remaining evidence needed to call
this a working Crumb decrypter.
