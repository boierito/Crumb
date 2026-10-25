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
