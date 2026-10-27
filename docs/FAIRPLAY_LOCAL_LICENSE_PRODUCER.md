# Local license producer work — 2026-10-27

Base: `9824cefee55ff5aded88c660653c0f477daed118`, branch
`feature/fairplay-local-license-producer`. No historical branch was merged.

## Implemented

A separate mobile SINF/page consumer in `MapleSyrup/NativeSAP/mobilelicense`:
account-credential key confirmation against original SINF, exact SINF binding,
private working-material lifetime, and mobile page derivation/decryption.
The native source packaging includes it and checks its iOS/arm64 compilation;
CI tests it independently of the large translated macOS runtime. The consumer
is not wired into the app's download/decrypt path because no valid producer
exists yet. No new manual credential path or IPA export is provided.

The implementation retains the mobile algorithms rather than assuming that the
macOS `LogicalKey`, white-box model, and page primitives are interchangeable.
Twelve deterministic Python-reference vectors cover both accepted private
record sizes and six page indices. Those vectors verify algorithms only.

## Concrete producer dependencies found

Pinned EthanArbuckle source: `7f80db2ef37634073fab06b3aad1aa0cf24561a1`.

1. `credential-retrieval/credential/credential.c`, `run_retrieve`: reads
   `/private/var/mobile/Library/FairPlay/iTunes_Control/iTunes/IC-Info.sidb`,
   requiring header `80 06 00 03`, not Crumb's observed `00 07 00 0b`.
2. `run_observe_uid_layer` / `observe_uid_callback`: obtains a separate UID-layer
   buffer from the daemon. `observer.c:observe_fairplay_requests` calls
   `task_for_pid` and sets breakpoints around `IOConnectCallStructMethod`.
   The credential helper's entitlements include `platform-application`,
   `task_for_pid-allow`, and `com.apple.private.thread-set-state`. These are
   dependencies of that extraction procedure, not evidence that every software
   producer is impossible.
3. `keybag_clear_from_uid_layer` uses that buffer and a UDID-derived descriptor;
   `select_license` extracts account/key-selected `sess`, `scra`, `valu` records.
4. `credential-retrieval/driver/driver.c` opens service
   `com_apple_driver_FairPlayIOKit` and calls selector 9 or 21. The returned
   side result is used by `core3` before `core4_context` and `finish_credential`.
   Porting the latter transforms cannot replace the missing upstream inputs.

Londek's current README and credential mint implementation retain the privileged
bootstrap; it is not an independent solution for obtaining a local credential.

Crumb `kbsyncretained.go:ImportKeybag` invokes the StoreAgent wrapper at
`storeAgentBase+0xd5ce0` and returns a filesystem snapshot. It does not invoke a
mobile SINF license consumer or return a mobile credential. Its `clearScratch`
also means that persisting only the filesystem cannot be assumed to preserve
all intermediate session inputs. This does not establish that the required
account key is in scratch memory, nor that it can be extracted there.

## Still unresolved — no first IPA claim

- Binary-level dataflow from the legacy import to a compatible license-consumer
  operation has not been reconstructed or demonstrated.
- No conversion from legacy SIDB to modern enrollment or mobile credential has
  been established. No format bridge was fabricated.
- No Apple service authorization/entitlement restriction was tested on a stock
  iPhone in this environment. The impossibility of all allowed routes is not
  proven.
- No private capture or original IPA was present in this workspace. No new
  download was requested. No account-specific SINF or IPA decryption was tested.
- Existing `completed`/`plaintext_authenticated=false` semantics remain unchanged;
  this package makes no additional authenticity or successful-export claim.

Next required evidence is a source-identified StoreAgent/CoreFP license operation
with complete ABI and input provenance, executed on the original retained
transaction. Adding more translated macOS blocks, random SIDB modifications or
more app downloads does not supply that operation.

## Producer continuation: native legacy state initialization

Following the first consumer-only commit, investigation downloaded the exact
public Apple 10.9 assets through the existing hash-verifying ipatool loader.
No private capture or account material was downloaded.

Static evidence in StoreAgent SHA-256
`70ce036f9dbcbc04db9511ebd08de0dd3cbc35ccc9d44b089c90170cb5453c59`:

- `-[FairPlayHelper needDiversityForDSID:]` at text+0x88e65 calls text+0x11d490.
  Its arguments are the original uint32 global context, the same 24-byte
  FairPlayHWInfo structure used by initialization, and the 64-bit DSID. Every
  nonzero result becomes YES. This is a provisioning query, not key confirmation.
- The keybag wrapper text+0xd5ce0 resolves discriminator 0x736b3779 through
  text+0x1232a0, reaching text+0xc4810 in the unrelocated pinned image. These are
  local native routines, not the modern fairplayd enrollment consumer.
- `scripts/verify-storeagent-license-profile.py` rechecks the entire asset hash
  and exact caller/callee instructions defining the new query ABI.

Implemented `OpenRetainedLicenseState`: install an isolated, deep-copied legacy
filesystem before calling the original native global initializer. The default
constructor still initializes with its original behavior. The new path does
not generate KBSync, replay an old keybag against a new ephemeral session,
convert SIDB to another generation, or contact Apple for a new transaction.

A successful context initialization is insufficient. Restoration additionally
requires the original SIDB bytes to remain unchanged. A discarded, recreated or
rewritten SIDB fails with `producer_state_not_preserved`. An input header alone
must never be interpreted as a recovered account.

Actual native execution with a public synthetic MAC/DSID confirmed that the
unenrolled account reports diversity required. A deterministic four-byte
truncated SIDB takes the native rejection/cleanup path: the initializer calls
`_unlink`, which was previously unsupported. The new bounded VFS implementation
removes only an existing closed virtual file, clears its private bytes, denies
host paths and refuses unlink-while-open (the VFS has path-based descriptors,
not POSIX inode objects). It never reports unsupported semantics as success.
The restoration constructor detects the native discard and rejects it; the
original input remains unchanged.

The session retains its own hardware identity for the query and serializes
native calls against Close. New filesystem ownership/bounds, unlink isolation,
source-overlay drift and native rejection tests accompany this implementation.
The host SAP test now runs the producer probes using the non-JIT TCI runtime;
the CI `legacy-producer` job runs this same command.

Not yet demonstrated: restoring the owner's real enrolled SIDB, deriving its
mobile credential, or validating its SINF. The new constructor and query do not
claim any of these outcomes. They are an implementation of same-generation
state restoration and native account querying, not a crypto format bridge.

Local verification of this continuation: the full `scripts/test-sap-host.sh`
command passed with the compiled TCI shared library (producer probe 111.67 s,
KBSync smoke 45.58 s, SAP smoke 2.75 s). VFS ownership/unlink tests also passed
with `-race`. The ABI verifier and source-overlay contract tests passed. These
are real native-code executions with public test inputs, not enrollment or
SINF-decryption evidence for the owner's account.

## iPhone test build 23046

The laboratory now exposes **Test original StoreAgent state (no download)**.
Import the saved capture ZIP and its exact original encrypted IPA, then run this
once. It verifies the manifest and SINF association before initializing any
producer, restores a private filesystem copy, and invokes the same native
account-diversity query tested on TCI. The DSID is read privately from the
capture context; it is never included in the report. Old contexts without a
DSID report `capture_dsid_missing` rather than guessing it from the SINF.

The result reports `state_restored`, `sidb_preserved_after_initialization`,
`account_query_completed`, and `needs_account_diversity`. The fields
`sinf_model_verified`, `ready_to_decrypt` and `plaintext_authenticated` stay
false. A false diversity result is not evidence that a mobile credential was
recovered. Native initialization rejection reports a fixed error code, with no
native error text, state bytes, DSID, identity or paths.

This is the device acceptance test for the producer restoration path, not a
request to rerun the previous modern-fairplayd bridge or SIDB mutation probes.
