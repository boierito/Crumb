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
5. Crumb opens its installed helper. Trust its development identity if iOS
   requires it, and leave it in the foreground. On success the helper returns to
   Crumb using a job-bound callback and retrieval starts automatically. Open
   Crumb Decrypt and Retrieve result remain available as recovery actions if the
   system refuses the handoff. Decryption runs in the foreground.
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


## Build 23024: connection diagnostics

Device feedback on build 23023 showed stored pairing keys and a failed Loopback
connection displaying only `device socket io failed`. That screenshot does not
identify the failing protocol stage or establish whether the LocalDevVPN route
works. No on-device IPA decryption has been established by this report.

Connection failures now include a fixed stage label (TCP, pairing handshake,
pairing verification, tunnel, RSD, app listing or same-installation proof).
Socket errors expose the I/O kind and OS error number. Other protocol errors
expose numeric library codes; response payloads and credential contents are
excluded. Tests cover stage/cause retention and private-payload exclusion.
The routing and pairing protocol are unchanged. This diagnostic change is not
proof that direct Loopback access is supported on the device.


## Build 23025: app Documents paths

Device feedback with LocalDevVPN selected reached AFC and reported
`afc error: Object not found`. The same-installation check incorrectly opened
`.crumb-device-proof` at the app-container root: VendDocuments allows access to
the Documents child without making that child the AFC root. The proof now opens
`/Documents/.crumb-device-proof`; the helper manifest and IPA retrieval use
`/Documents/result.json` and `/Documents/result.ipa` for the same reason.
The fresh nonce comparison remains required. Normal reads explicitly close
their remote descriptors, including after a failed small-document read.

A regression fixture speaks AFC to the real dependency client, reproduces the
old ObjectNotFound response, and checks successful reads and descriptor closure
for all three Documents paths. It is a protocol test, not a physical-device
acceptance test. Re-test Check authenticated connection with LocalDevVPN enabled;
helper installation and actual licensed decryption still require device testing.


## Build 23026: bounded dual-stack Loopback attempts

Direct mode previously tried only `127.0.0.1`. It now attempts the complete
authenticated tunnel over IPv4 and then `::1`, with a 15-second bound per
family. TCP acceptance does not count as success: pairing verification, TLS
and RSD must complete before selecting that route. Failures identify each
endpoint and its protocol stage. Direct mode never silently falls back to the
VPN address. Changing the selected route clears the previous Connected status
and requires a new authenticated check.

Three regression tests cover the destination allowlist, a TCP peer that accepts
then rejects the real RemotePairing handshake on both families, and a stalled
IPv4 handshake that must time out before IPv6 is tried. They do not prove
that iOS accepts a direct authenticated tunnel.

[Published iOS 27 probes](https://github.com/leeguooooo/iphone-use/blob/7e9c7bdb2bd708c852c8034953d27c50bcc1c3ca/experiments/phone-self-host/README.md)
report TCP acceptance followed by a RemotePairing reset on direct IPv4. Their
successful on-device path uses a Network Extension. This is evidence against
treating TCP reachability as VPN-free support, not proof of the cause on this
user's device. An integrated VPN extension is not part of this build.

For the physical comparison, establish Connected with LocalDevVPN selected
first. Then disconnect the VPN, select Loopback (experimental) and repeat the
check. Record either Connected or the complete IPv4/IPv6 errors. If iOS refuses
both direct routes, use LocalDevVPN for the decryption acceptance steps.


## Build 23027: helper profile validation

Device feedback reached local helper preparation and reported `Profile does
not permit the helper bundle identifier`. This can indicate an actual explicit
profile for a different app, including Crumb itself. The original validator also
incorrectly assumed the App ID prefix always equals the Team ID.
[Apple QA1879](https://developer.apple.com/library/archive/qa/qa1879/_index.html)
explicitly documents that these values can differ.

The helper application-identifier now uses the prefix in the profile's
entitlement, checked against ApplicationIdentifierPrefix. The Mach-O signing
team and team-identifier entitlement continue to use the Team ID. Exact and
terminal-wildcard bundle restrictions remain enforced; Crumb cannot be replaced
by the helper. Mismatch errors name both the permitted pattern and requested
helper ID, and explain when a profile is restricted to Crumb itself.

Profile scope and expiry are checked before unpacking the target IPA. A local
preparation error preserves the previously authenticated connection state.
Five regression tests cover different App ID/Team ID prefixes, wildcard scope,
actionable mismatch errors, the Crumb-only profile case and invalid prefixes.

The imported profile from the device has not been inspected. This change does
not establish which of the two mismatch causes applies to that profile. An
explicit profile for Crumb still cannot authorize com.certlium.crumb.decrypt;
the user needs a matching separate-helper profile or a permitted wildcard one.
No profile is modified or generated and no Apple developer provisioning is
automated by this change.


## Build 23028: installed-file diagnostics and failed-result recovery

Device feedback now shows an authenticated LocalDevVPN connection and a
launched helper. Direct IPv4 and IPv6 both reset the RemotePairing handshake
with OS 54; direct access is not established on that device. The helper reports
`Cannot open installed image`, which is an open() failure, not evidence that
FairPlay accepted or rejected that image. The failing path and errno were
previously discarded, so the underlying cause is not established.

The helper now preserves open errno, relative target path and available file
mode/ownership metadata. It retries EINTR and retains O_NOFOLLOW, O_RDONLY and
close-on-exec. It does not make installed files writable, follow rejected links
or mark an unreadable file decrypted. Four portable OS-level tests cover
read-only access, a missing path, a denied read and a rejected symlink. These
do not simulate iOS sandbox or FairPlay policy. The UIKit helper is compiled
and linked against the iOS SDK in both build configurations.

Retrieval checks the job identity before surfacing a bounded failed-helper
message and decrypted-image count. It rejects incomplete output and keeps the
original image error instead of replacing it with a generic message. Three
regressions cover partial failure, unrelated-job diagnostics and message bounds.

User signing files are inspected outside the repository, not committed or
packaged. No private key or password belongs in CI artifacts. Certificate
compatibility does not by itself prove decryption; the installed-file failure
needs a physical retest with the new diagnostics.


## Build 23029: own-device pairing port discovery

Direct mode previously assumed that the device's RemotePairing service always
used 49152. Apple can advertise a different listener port. Loopback now browses
_remotepairing._tcp for up to three seconds and accepts only a service resolving
to an address held by this iPhone (getifaddrs). Advertisements from other devices
do not supply a route. No personal service names or interface addresses are
printed or persisted. The UI shows the discovered port or lookup failure.

The advertised port is attempted before the default, over both 127.0.0.1 and
::1. Native inputs permit at most two valid ports, deduplicate them and retain
the fixed local host allowlist. All candidates still require pair verification,
TLS, RSD and the fresh same-installation proof. Complete tunnel attempts are
bounded to at most roughly 32 seconds across four candidates, after discovery.
Discovery can be cancelled and does not allow an older callback to start work
after cancellation. Bonjour and UIKit behavior require physical testing.

Two native regressions cover advertised-port ordering/local-only destinations
and malformed or excessive port lists. This change corrects port selection; it
is not proof that iOS accepts direct loopback on the user's phone. If the
verified advertised port also resets, the reset remains unresolved. The incoming
Settings pairing exchange is a pair-setup flow; the reviewed upstream responder
does not expose a documented device tunnel on that existing connection.

The installed helper must be prepared and installed again from Downloads to
receive the file diagnostics added in 23028. Merely opening the previously
installed helper continues running its old code. Signing files are unchanged.
