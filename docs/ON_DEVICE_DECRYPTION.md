# Build 23037: correctly match iCloud environment grants

Profiles may authorize iCloud environments as an array while the signed app
selects one string. The signer now recognizes this key-specific representation
and preserves the original selected environment when authorized. It does not
switch production/development environments or rewrite CloudKit container IDs.
Tests preserve all four reported iCloud capabilities when granted and reject
unrelated containers. This fixes a matcher defect, not confirmed Pinlist recovery.

The shareable signing report now includes per-bundle original/profile/signed
iCloud and ubiquity values. It contains no private key or certificate password.
This distinguishes missing profile authorization from erroneous entitlement loss.
Actual Pinlist profile/signature comparison and device runtime validation remain
necessary; preserving unprovisioned entitlements cannot grant original CloudKit access.

# Build 23036: one signing identity

A single .p12, .mobileprovision and saved password now sign both the helper and
local decrypted IPAs. The existing helper identity takes precedence on upgrade;
the old installation identity is a migration fallback only. There is no second
app-installation identity editor. Both settings entry points reuse one view.

CFBundleIdentifier is preserved by default and can be customized for installation.
Like zsign's default profile-entitlement template, an explicit profile identity
can be used for signing without forcing the app's CFBundleIdentifier to match it.
This applies to all profile types, not only enterprise profiles. Nested bundles
prefer an available matching profile and otherwise reuse the selected profile.
Crumb does not reserve a profile for itself/helper; actual replacement of the
running app or its managed helper is still rejected. The profile CMS is never
altered. Final profile/device/capability acceptance is determined by iOS.

Supported entitlements/resources/extensions are retained. Adjusted capabilities
are reviewed before signing. The original decrypted IPA is never overwritten;
the separate signed artifact can be installed/exported. Original and custom IDs,
main/extension signatures, explicit helper-profile reuse for non-enterprise and
enterprise profiles, resources, executable modes and cancellation are covered by
synthetic OpenSSL/CMS fixtures. These do not establish device installation.

Pinlist's provided crash report has an uncaught Objective-C exception inside
CloudKit; missing/incompatible iCloud authorization is a hypothesis, not a
confirmed decryption defect or resolved crash. A third-party signing identity
cannot create access to the original developer's CloudKit containers.

References consulted: zsign profile-derived entitlements in
https://github.com/zhlynn/zsign/blob/master/src/openssl.cpp and nested identifier
updates in https://github.com/zhlynn/zsign/blob/master/src/bundle.cpp; Feather's
Zsign submodule at https://github.com/claration/feather/tree/main/Zsign.
Implementation remains in the existing Rust signer; these projects are not vendored.

An integration fixture generates a temporary self-signed test certificate/profile,
signs a main app and extension with actual CMS, verifies output identifiers and
entitlements, shared-profile reuse, preview retry, cancellation, unchanged
source bytes and preservation of Settings resources/empty folders. It is not an
Apple-approved profile or a physical-device acceptance test.

A reported Settings crash remains unclassified without the target's .ips and
installation details. Package/cryptid/hash verification does not prove runtime
behavior; neither signing success nor dropping capabilities is a crash fix.

---

# Build 23034: stop and recover

Downloads always exposes Stop for an active decrypter session, including the
external-helper waiting state. Stop invalidates the callback generation, clears
pending metadata/queues and frees the UI immediately. It retains signing,
pairing, completed downloads and the managed-helper identity for the next job.

Cancelled native workers are retired independently. Their workspace is deleted
only after their terminal state; they cannot publish into a newer UI session.
The native signing/fetch code already checks cancellation before publication.

New helpers support a job-bound cancel URL and a Stop button. Cancellation is
persisted per job, checked between files, encrypted pages and verification stages,
and reported to Crumb. A later job does not inherit an old cancellation. Crumb
also unlocks with an older helper that lacks this protocol, without reopening
that older app. New jobs advertise the cancellation capability explicitly. iOS system calls
already in progress cannot be forcibly interrupted; cancellation is cooperative
at the next checkpoint, while Crumb itself is immediately available again.

Settings uses equal-height Telegram/GitHub controls with an internal border,
with a vertical layout for accessibility text sizes. Links point to
https://t.me/certlium and https://github.com/boierito/Crumb.

---

# Build 23033: result retrieval and connection scheduling

The 23032 pending job stored download metadata in `record`, the same field used
for the authenticated pairing plist. Fetch merged all pending-job fields over
credentials. This produced the reported `10.7.0.1:49152: io on plist` while
parsing local data, before TCP—not evidence of a VPN disconnect.

Download metadata now uses `download_record`; retrieval whitelists only job ID,
helper bundle and output path. Reserved authentication/routing fields cannot be
overridden by operation parameters. Earlier saved jobs remain readable.

Automatic probes are event-driven: one per foreground/network transition,
coalescing duplicate notifications, with no periodic setup loop. A pending
helper/result prevents background verification. Helper callbacks queue retrieval
if an operation is already running; completion drains it without losing the
callback. Native operations still establish their own authenticated tunnel;
this does not keep or control the external VPN app's connection.

Project/source/plist/scheme paths now consistently use Crumb. App bundle ID,
Keychain namespaces, saved pairing, signing files and download history remain
compatible. Prior source attribution is preserved in CREDITS.md and notices.

Regression coverage includes the original metadata/credential collision,
reserved-field protection, duplicate events, cellular/helper suppression and a
native invalid-record test proving that this failure occurs before TCP.

---

# Build 23032: user decrypter workflow

- One Decrypter setup screen with automatic connection checks, compact signing-file
  summaries, validated password persistence in device-only Keychain, automatic
  helper ID derived from provisioning-profile scope, and cellular/PIN guidance.
- Four version actions: encrypted/decrypted, each with download or download/install.
- Persistent encrypted and decrypted history with export and installation actions.
  Decrypted installation signs the main app and nested app extensions, validates
  profile scopes and original capabilities, then uses the existing iOS OTA handoff.
  A profile that cannot authorize an extension/capability fails without dropping it.
- The helper uses the approved branding/Crumb-AppIcon-1024.png artwork.
- User-provided FilePickerFix.dylib is embedded in Crumb (not target apps) through
  an arm64 LC_LOAD_DYLIB command before ordinary signing. The build checks its
  exact bytes and load path. Runtime behavior still needs physical testing.
- Older-iOS connection setup can import a computer-generated Lockdown pairing
  plist (e.g. mobiledevicepairing), stored in Keychain, and uses authenticated
  Lockdown over LocalDevVPN. Remote-pairing plist import is also accepted.
  Same-installation nonce validation remains mandatory for either transport.
  This adds an older-iOS connection route, not evidence of decryption support:
  the FairPlay method remains enabled only on iOS 27.0.1.

Physical testing must confirm saved-password reuse, automatic reconnection,
helper icon/progress, main/framework/extension decryption, persistent decrypted
history and signed installation. Credentials are never committed or logged.
Historical implementation notes follow; the 23032 workflow supersedes their UI.

---

# Current user-flow revision

Base: `feature/on-device-pairing-decryption` (4a3ca40).
Use LocalDevVPN with authenticated pairing and a separately signed helper.
Loopback is removed from the app flow; the older notes below describe history.
The native loopback fixtures remain for protocol regression coverage.

Pairing establishes trusted-host credentials. LocalDevVPN supplies the local
transport used for helper installation and result retrieval. Neither grants
code-signing privileges, an App Store license or decryption compatibility.
The helper still requires a matching certificate and provisioning profile.
Installation opens the helper automatically; successful processing returns to
Crumb and starts retrieval. iOS trust prompts still require user interaction.

On-device pairing is limited to 27.0–27.0.1 and decryption to 27.0.1.
Older iOS support would require a separately implemented computer-assisted
pairing/import transport and independent decryption validation. A VPN alone
cannot enable it. Do not remove version gates on the basis of pairing alone.

The helper now includes Crumb home-screen icons, Dynamic Type text, file-based
processing progress and separate resource-checking/packaging stages. Progress
is based on file count, not an estimate of time remaining. Console diagnostics
in the Swift app are restricted to DEBUG builds.

Preparation snapshots a complete target inventory before signing/installation.
Before packaging, the helper requires every expected file/directory, matching
sizes, and matching SHA256 for files that were originally unencrypted. Original
encrypted images must pass the existing actual decryption checks. This detects
missing extensions/resources even if the output ZIP has internally valid hashes.
Reinstall a newly prepared helper: older helpers lack the required inventory.

Final decrypted-app signing and installation are **not implemented**. Export is
available; the IPA needs a signer capable of provisioning the app and all nested
extensions/entitlements. The helper's profile does not generally authorize the
target app. Do not label export as an installable signed IPA.

Acceptance still requires macOS compilation and physical-device testing with
licensed main apps, frameworks and extensions. Linux parser/file-access tests
do not establish FairPlay success, UI behavior or installability.

---

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


## Build 23030: preserve installed license metadata for export

The device's new error identifies an extension SC_Info/*.supp resource:
open(O_RDONLY) returns EPERM despite mode 0644 and uid/gid 33. Ordinary mode
bits alone do not explain the denial. This is a license resource, not evidence
of a failed encrypted Mach-O remap. The helper was incorrectly trying to read
every installed SC_Info file while reconstructing the output.

Preparation now preserves SC_Info resources from the downloaded target before
installation, after SINF replication. Originals remain unmodified under Target
for FairPlay. Separate, encoded copies and a path/size/SHA256 manifest are sealed
into the helper as ordinary resources. The helper skips installed SC_Info
descendants and restores original metadata from these copies after processing
the images. It verifies decoded size, hash and absence of encrypted image
commands before writing. These restores do not increment the decrypted-image
count or alter cryptid; actual encrypted images still require the existing
FairPlay checks. Output verification still checks all ZIP file hashes and
encryption commands. No metadata is silently dropped.

Preservation limits are 4 MiB per license resource, 16 MiB total and 4096 entries.
Native tests cover main/extension byte preservation, unchanged originals,
encrypted metadata rejection, size limits and links. macOS runs the actual
Foundation restoration implementation against seven fixtures: normal restoration
without an installed SC_Info tree, hash mismatch, path traversal, an encrypted
image disguised as metadata, a symlink, excessive size and an empty resource. This does not
simulate FairPlay or assert that a full physical-device decryption passed.

Re-prepare and install the helper from Downloads after updating Crumb. The
already installed helper cannot gain the new snapshots from a Crumb update.
Direct Loopback retains the 23029 discovered-port behavior and remains pending
physical validation. This patch addresses the measured installed-resource read
failure; it does not establish a change in iOS's direct connection policy.


## Build 23031: preserve the original IPA root

Preparation previously moved only Payload/*.app into the installed helper.
Although the input ZIP was extracted in full, top-level META-INF,
iTunesMetadata.plist and other non-Payload files were left in the disposable
input workspace and never reconstructed.

The helper now receives sealed, encoded snapshots of every supported original
file outside Payload and an inventory of its directories. Original names and
bytes are preserved; absent metadata is not invented. Archive metadata is not
parsed or regenerated. Restoration streams bounded base64 chunks to temporary
files, validates size/SHA256 and encryption commands, then moves verified bytes
into the reconstructed IPA root. Root entries cannot overwrite Payload. The
existing 512 MiB per-file, 4 GB overall and thin-arm64 image restrictions remain.
License preservation retains its separate 4 MiB/file and 16 MiB total limits.

The ZIP producer now records non-Payload directories, including empty ones,
and hashes directory entries as empty data. Native verification checks entry
types and rejects missing or changed original root files/directories against
the preserved inventory. Job/result schema 3 requires a newly prepared helper.
The combined snapshot manifest is limited to 16 MiB.

Four new native cases cover original metadata/artwork bytes, empty root folders,
large streamed resources, missing/changed metadata and directory type checks.
The shared Foundation restore/archive code runs thirteen cases, including errors
that must survive internal autorelease pools. An independent
Python ZIP reader inspects its actual output for exact root paths, metadata
bytes, directory flags, hashes and CRC. Test inputs are synthetic resources,
not a simulated FairPlay success.

META-INF bytes are copied from the original package. Decryption and repacking
change the archive, so this does not generate a new valid Apple archive
signature. The exported app still needs normal re-signing before installation.
Re-prepare and install the helper after updating Crumb; previous helpers lack
these original root snapshots. Physical full decryption remains unconfirmed.
