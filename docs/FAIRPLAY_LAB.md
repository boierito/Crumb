# FairPlay laboratory — device test guide

This build tests portable components and collects inputs. It does **not** yet
create the missing enrollment state from an Apple login or a keybag.

1. Sign the Release IPA and install on iOS 16.4 or newer. Sign into Crumb.
2. Open Settings → FairPlay Laboratory → Run self-tests. The state, frame and
   SIDB vectors should all pass. These use public synthetic fixtures only.
3. Enter an App Store ID/link or bundle ID for a small free app you can download.
   Tap Download and capture transaction. Wait for completion. Log entries show
   matched kbsync/keybag lengths and SINF count; absent keybags and fallback
   routes are recorded. Each successful transaction has its own capture folder.
4. Tap Inspect IPA. Record SINF/SUPF counts and Mach-O candidates. Alternatively
   open Downloads → app options → FairPlay Laboratory → Use selected downloaded
   IPA. Existing downloads contain no newly captured Store transaction.
5. If you already have a matching private enrollment model, import it and test
   decryption. Without one, stop here: the missing raw-state → 192/224-byte state
   input stage is not implemented. Capturing a keybag is not recovery.
6. Prepare exports. Share the diagnostic for ordinary troubleshooting. Share the
   PRIVATE capture only with a trusted researcher: it contains account-derived
   material. No password, 2FA code, auth token, cookies or signed download URL is
   included. The IPA and imported model are not included in the capture ZIP.
7. Export a successful decrypted IPA, sign it, and record installation and launch
   results. Source CodeDirectory hashes check encrypted input; structural ARM64
   checks and cryptid=0 alone do not prove plaintext correctness.
8. Export before closing. Closing or Erase this test removes this session's
   temporary files; original downloads, account and other sessions are preserved.
   Repeat with a second app from the same account. Force-quitting may leave
   temporary files until iOS clears its temporary storage.

Portable source pin: DrWhax/ipatool ef98227e061838320424744d21a513baba7ade28.
The native build isolates state/frame/SIDB packages, retains their public vectors,
and patches the IPA page offset and first executable-page audit. Model resources
are pinned by the upstream fixture tests. SAP remains pinned independently.

Host tests do not establish real enrollment compatibility or iOS success. Device
logs intentionally do not expose raw native exceptions or secret evaluator output.

## Build 23023: capture analysis and retained producer

The step-3 **Import and analyze private capture ZIP** action accepts v1 and v2
captures. It checks manifest sizes/hashes, duplicate and unsafe ZIP paths, size
limits and SINF structure. Its report contains counts and readiness flags only.
It does not identify opaque keybag bytes as AES keys or raw SIDB state.

New laboratory downloads bypass the normal KBSync cache and keep the exact
StoreAgent producer alive until the validated ent reply has been bound to its
original kbsync. Binding checks byte identity and retains the keybag in memory;
it does not call an unverified FairPlay consumer. The emulator is closed before
the CDN transfer, and both blobs are cleared from that retained session.

Capture v2 adds `context.json` containing private hardware identity and account
DSID for replay research. This is metadata, not a guest-memory snapshot or
recovered enrollment key. It is included only in the PRIVATE capture ZIP and
excluded from public logs. Normal downloads keep their existing cache behavior.

The uploaded v1 control passed integrity checks and SINF inspection (440-byte
private payload). It lacks generation context metadata. The upstream raw-source
kernels at 0a851d5 remain research components: their dispatcher and raw-state
composition are unfinished, and they do not implement a keybag consumer.
Crumb deliberately keeps its validated portable pin until that connection can
be implemented and tested. No supplied private capture is committed or bundled.

## Build 23024: execute the native keybag consumer on iOS

Static analysis of the hash-verified StoreAgent binary found the same wrapper at
image offset `0xd5ce0` in three independent keybag handlers (authorization,
provisioning and purchase). Its native callers pass the helper's `globalContext`
as uint32, the NSData bytes pointer and its uint32 length. This is separate from
the macOS package `dpInfo` session initializer. No target binary is run natively.

The retained laboratory generator now calls that wrapper on the exact emulated
context after matching the transaction kbsync. The interpreter models the legacy
OS-version queries and a private, in-memory filesystem for flat files under the
emulated `SC Info` directory. Files are capped at 32; aggregate data at 4 MiB;
individual writes at 1 MiB. These handlers are enabled only for the laboratory
producer. They cannot read or write the host filesystem. Files and keybag bytes
are cleared on close. Single-emulator file locks validate descriptors and require
no host lock because each retained runtime is isolated and serialized.

The call runs away from MainActor. Capture v3 includes private
`consumer-state.json` with the native status, completion flag and any virtual
state-file bytes. Public logs show only status/counts. `import_succeeded` means
only that this wrapper returned zero; it is never reported as an enrollment
model or proof that any IPA was decrypted. No recovered key is embedded.

An offline replay of the supplied capture with a newly generated context reached
the consumer and returned -42023, with an empty virtual state file. That result
does not establish whether the live same-context transaction will be accepted.
The fresh transaction path now executes this actual operation, rather than only
retaining the response. Download and export remain available after a consumer
rejection. Portable enrollment recovery remains unimplemented.

## Build 23025: correct the virtual filesystem operations used by import

The supplied v3 control verified all manifest hashes and confirmed native
-42023 with no nonempty state. A local syscall trace found that F_GETPATH and
lstat still used the generic failing stubs, despite referring to virtual nodes.
The laboratory now models descriptor flags and F_GETPATH, stat/lstat of existing
virtual nodes, opaque FSRef round trips, and the requested Carbon catalog
permissions operations. Unknown security-object requests remain rejected.

Catalog layouts follow Files.h in the macOS 10.9 SDK: FSRef is opaque (80 bytes);
FSCatalogInfo is packed to 2-byte alignment, and the LP64 permissions mode is at
byte 66. These guest structures never become native host filesystem references.
Permission changes affect the isolated VFS only; private app captures retain
their real iOS file protection and 0600 permissions.

Regression tests cover the actual F_GETPATH/lstat failure, invalid descriptors,
FSRef lookup and round trips, the packed permissions layout, consistent inode
IDs, and permission updates. The offline v3 replay then wrote 287–306 bytes of
virtual state and advanced past the previously missing Carbon operations. Its
final native status became -42052 with no emulator exception. This replay uses
a freshly generated context, so it is not a substitute for a live same-context
transaction. Neither writing a state file nor native status zero is treated as
an enrollment model or authenticated plaintext.

Header reference: phracker/MacOSX-SDKs, MacOSX10.9.sdk, CarbonCore/Headers/Files.h
(Apple SDK declarations). No SDK header, Apple binary or private control data is
bundled by this change.


## Build 23026: distinguish failed import from successful enrollment

Capture v3 analysis now reports one of four observed consumer stages:
`not_observed`, `emulator_failed`, `native_rejected`, or
`native_returned_zero`. A successfully completed emulator call, a nonempty
state file, and a zero native return are separate observations. None can be
promoted to `ready_to_decrypt` without independently validated enrollment
state and plaintext verification. For a rejected native call, the next missing
stage is `accepted_keybag_import`, not the previously misleading
`keybag_consumer` implementation placeholder.

The capture importer checks native result consistency, state-file bounds,
and names; it publishes counts and status only. Synthetic v3 fixtures test
native status -42052, emulator failure, native zero return, inconsistent
status flags, and omission of private data from diagnostics. The retained
KBSync comparison copy is cleared, and the keybag copy is cleared immediately
after its single consumer call.

No private capture is checked in. Next experiment: test a **fresh, matched
transaction on-device** with build 23026 and compare its sanitized native
status and state sizes against the offline -42052 replay. Replaying v3 with
a newly constructed StoreAgent context cannot prove original-context import.


The retained guest also records the last 64 virtual filesystem operations as
fixed operation labels and coarse results (zero, nonzero, errno, negative status,
or emulator failure). Error entries can carry numeric errno only. No guest
arguments, paths, file data, register values, or account identifiers are
recorded. The `virtual_io` field is added to **private** consumer-state.json
only; the ordinary public laboratory summary remains counts/status only.
This can narrow down future native -42052 failures without exporting secrets.


## Build 23027: V4 device evidence and catalog-write diagnostics

An independently captured **private** V4 transaction reproduced a successfully
completed emulator call followed by native status `-50` (import rejected).
Its integrity-verified manifest showed 196-byte kbsync, 904-byte keybag, and
one SINF. The exported virtual state contained **two entries**: an empty
lock file and an **865-byte SIDB**. The public UI reported 0 files/0 bytes
because it tried to cast a JSON dictionary containing one `null` value to
`[String: String]`; the all-or-nothing conversion failed. The Swift summary
now counts mixed `null` and base64 byte strings correctly, without logging
any file contents.

The V4 guest trace contained 21 events. It reached `FSPathMakeRef` and
`FSGetCatalogInfo` with zero return values, then a negative return from
`FSSetCatalogInfo`. This is the **leading suspected source** of the observed
`-50` paramErr; it is not yet proven which parameter was rejected. The next
build annotates only that catalog-write event with one of
`unsupported_catalog_mask`, `missing_catalog_info`, or
`unknown_catalog_reference`, plus a bounded catalog flag mask (<=0xffff).
It does **not** log the FSRef, pointers, paths, SIDB bytes, DSID or keybag.
Unsupported operations still fail closed, rather than being silently accepted.

**Next test**: install build 23027 and capture a new transaction in the same
producer context. Compare the private catalog-write `detail` and `mask`
against the status and generated state sizes. Do not check private capture
archives into this public repository. Status zero and SIDB creation are
insufficient to mark any decrypted IPA as validated.


## Build 23028: replay V5's observed Carbon catalog flags

Private on-device V5 after build 23027 reported synthetic self-tests passing;
matched download artifacts: 196-byte kbsync, 888-byte keybag, and one SINF.
The manifest hashes validated and the virtual state had an empty lock file
plus a nonempty 878-byte SIDB. Native import still returned `-50`.

The private, bounded VFS trace identified the exact rejected call:
`_FSSetCatalogInfo` with `detail=unsupported_catalog_mask`,
`mask=3072` (`0xC00`). The old shim accepted only `0x400`.
Per Apple's historical Carbon File Manager `Files.h`, `0x400`
is `kFSCatInfoPermissions` and `0x800` is `kFSCatInfoFinderInfo`.
The FSCatalogInfo structure contains 16 bytes of Finder info at
offset 72 for this guest ABI, with permission fields beginning at 56.

The virtual filesystem now explicitly supports permission and Finder-info
fields in combination, persists the Finder metadata per virtual file, and
provides corresponding `FSGetCatalogInfo` readback. It applies only fields
selected by `whichInfo`, so updating Finder metadata cannot accidentally
zero permissions. Unsupported flag sets still fail closed.
Synthetic regression tests exercise `0xC00` roundtrips, Finder-only
preservation, and rejection of unsupported masks. This is **ABI
compatibility**, not decryption or a verified portable enrollment state.

Next on-device experiment: build 23028, use the same App Store app and
record whether `_FSSetCatalogInfo` becomes `zero` and which operation
comes next. Check sanitized native status, file counts/sizes and bounded
virtual-I/O trace without sharing raw blobs or identifiers in the repository.


## V6 / build 23028: on-device native import accepted

A fresh same-context transaction on iOS produced native status **0** with
`emulator_completed=true` and `import_succeeded=true`. All five manifest
entries in the private V6 ZIP passed SHA-256 and byte-length checks. The
bounded trace contained **21 operations**: `FSPathMakeRef`,
`FSGetCatalogInfo`, and `FSSetCatalogInfo(0xC00)` now all returned zero.
The private virtual state contained **two entries**, one empty lock file
and one **862-byte** SIDB. Input sizes were kbsync **196 bytes**, keybag
**904 bytes**, and one SINF.

**This closes the previously failing virtual catalog-write compatibility
step; it does not close FairPlay recovery.** In particular:
- The accepted StoreAgent keybag consumer is not a verified enrollment model.
- The captured SIDB is a product of a guest virtual filesystem, not a
  documented replacement for the modern fairplayd enrollment frame or
  its missing raw state/orchestration.
- On-device inspection still reports two encrypted images / 712
  encrypted pages. The portable model remains unavailable.
- The native status means *only* that the observed wrapper returned zero.
  Successful IPA plaintext authentication, signing and launch have not
  been demonstrated.

Next engineering stage: specify and validate the bridge from this accepted
guest state to the exact raw-state inputs and matching enrollment model
required by the pinned portable evaluator, or demonstrate that the legacy
StoreAgent output is incompatible. Do not treat arbitrary SIDB bytes as a
192/224-byte state input, and never publish account/device material.
A new **synthetic** v3 capture test guards accepted-import classification,
the presence of empty plus nonempty virtual files, and separation between
native success and decrypt readiness. The private V6 archive is not in Git.


## Build 23029: local enrollment-stage audit after successful V6 import

V6 demonstrated a correct native StoreAgent return value of zero and a
nonempty virtual legacy `SC Info.sidb`. It did **not** provide the required
modern fairplayd raw-state source, 16-byte enrollment frame, or an authenticated
portable model. The embedded macOS 26A428 evaluators are build-specific; an
unstructured legacy SIDB blob cannot be treated as an evaluator input.

The native capture analyzer now returns a **sanitized** `enrollment_bridge`
object with the observed legacy SIDB length and independent boolean gates for
portable raw-state source, enrollment frame, and model validation. All three
remain false until separately implemented and demonstrated. Synthetic tests
check that successful StoreAgent import does not make `ready_to_decrypt` true,
including mixed null lock files and nonempty virtual SIDB state. No real
account files, identifiers, or byte sequences are added to the repository.

The FairPlay Laboratory now has **Verify current capture — no import required**.
It exports the on-device current-session capture into an owner-only ZIP,
applies complete iOS file protection to that ZIP, validates the manifest and
SINF through the existing native analyzer, and displays **Enrollment logs**
containing counts and stage flags only. `Prepare exports` also runs the same
audit *before* generating `diagnostic.json`, so diagnostics include these
checks without forcing users to reimport their own private ZIP. Importing an
**older** capture remains optional. Private captured content is not exposed by
the enrollment log. The private ZIP still must not be checked into GitHub.

Next required engineering work is a rigorously validated bridge from a
matching enrollment to the pinned modern portable evaluator, including a
reproducible frame and raw-state source and an independent IPA decryption
authentication check. This build does not claim to recover that model.


## Build 23030: V7 diagnostics are stable, one local audit per transaction

The private V7 capture of build 23029 passed all **five** manifest SHA-256
checks. The retained StoreAgent consumer again returned native status **0**
with emulator completion and accepted import. It generated two virtual files:
an empty lock file and an **883-byte** legacy SIDB. The bounded trace shows
21 operations and a zero result for `FSSetCatalogInfo(0xC00)`.
Input sizes were 196-byte KBSync, 904-byte keybag, and one SINF. The
portable enrollment model remained unavailable, and two Mach-O images
still contained 712 encrypted pages.

Different real transactions produced differently sized legacy SIDB output;
synthetic 862- and 883-byte fixtures ensure no fixed-size assumption enters
the analyzer. As before, any accepted native status **does not** imply
portable 192/224-byte raw-state source availability, 16-byte enrollment frame
recovery, matching model validation or decrypted IPA authenticity.

The laboratory now caches a successful local integrity check per capture
transaction. Tapping `Verify current capture` repeatedly, then
`Prepare exports`, does not reexport the same private ZIP or append
the full JSON three times. A new transaction or `Erase this test` resets
the cache. A small read-only `Enrollment readiness` view shows native
import, legacy SIDB, raw-state, frame and model status independently; a
fallback download without a matched keybag is shown as unavailable
rather than failing the entire export.

This is a diagnostic and correctness iteration, **not** a decryption
implementation. The remaining research problem is the validated modern
enrollment source and cross-profile compatibility, with an independent
successful page-hash/plaintext check before promising portable recovery.
