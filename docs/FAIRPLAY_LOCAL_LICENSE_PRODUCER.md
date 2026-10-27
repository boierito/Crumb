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

## Owner-device result and authorization producer investigation

The owner's build-23046 diagnostic on iOS 27.0.1 completed the original-producer
probe in 18,679 ms: exact capture/SINF binding passed, the SIDB remained unchanged
after initialization, and the native account query completed with diversity
required. The same transaction had 196-byte KBSync, 888-byte keybag, one 440-byte
SINF private record, and 884 bytes across two virtual files (one nonempty).
These sanitized observations do not identify the account, license selector or
any private bytes. The private capture itself is not present in this workspace.

This confirms the probe executed on the owner's iPhone without discarding its
saved SIDB. It does not prove that the account is empty, that all account state
was consumed, or that a license was decrypted. `sinf_model_verified=false` in
this particular probe is a deliberately false pending-stage flag: the probe
**does not attempt SINF cryptography**. It must not be described as a rejected
candidate-key test.

Following the only `needDiversityForDSID:` selector reference in the pinned
StoreAgent identifies a separate machine-authorization flow:

- `PurchaseController authorizeMachineWithInfo:clientIdentity:` resolves the
  Bag key `authorizeMachine` at text+0x5925f. The request-body block generates
  KBSync at text+0x5988f with arguments `(global, dsid, 0, 2, out, length)`.
  Original purchase/download callers instead use mode 1; Crumb's existing
  download generation matches those callers and is not changed to mode 2.
- At text+0x59a27 the request builder queries `needDiversityForDSID:`. Only when
  nonzero does it add the string value `"1"` under the `needDiv` key. It also
  sends typed KBSync data, GUID and machineName in a POST plist.
- The response handler imports `keybag` via text+0xd5ce0, then separately
  extracts `diversitybag` at text+0x59cff and imports it via text+0xcdea0. Its
  call at text+0x59d8f passes `(global-u32, bytes-pointer, length-u32)`.
- The original authorization path also invokes ADI request mutation. Static
  code in `ADIController _mutateURLRequest:onCurrentThreadWithDSID:` adds
  X-Apple-MD-M, X-Apple-MD and, when present, X-Apple-MD-S. This is not evidence
  that every server authorization requires those headers, but their protocol
  semantics must be investigated before claiming a complete emulation.

The public Bag for the Configurator profile currently resolves `authorizeMachine`
to Apple's `/commerce/machine/authorize` endpoint. That discovery is read-only;
no machine-authorization request was sent and no device slot was consumed.
The endpoint is not copied into the ordinary ent/download request, and `needDiv`
is not blindly added to that different protocol.

Implemented internal producer primitives `AuthorizationSync` and
`ImportDiversityBag`, using those source-identified native calls on the retained
session. They return private protocol material/native status; neither returns
a logical key or marks a SINF valid. The hash/ABI verifier now checks the
mode-2 call, distinct import call and original protocol literals. A native
host execution generated a 196-byte authorization sync and rejected a public
invalid diversity response with status -42018. These results validate the
entrypoints and rejection behavior only, not real enrollment.

Remaining critical work: reproduce the complete authorized request/session and
its ADI challenge semantics, obtain a genuine matching diversity reply locally,
import it in that same producer, and test the original SINF with the resulting
material. Diversity provisioning is a source-backed hypothesis, not a claim
that it alone solves mobile FairPlay. Preserve the owner's current capture and
IPA for those tests; repeating downloads or the completed restoration probe
cannot establish the missing cryptographic step.

The full non-JIT host command passed after adding these primitives: the original
producer restoration/rejection probe took 116.22 s, the authorization primitive
probe took 49.63 s, and the SAP/KBSync smoke suite took 51.86 s. The authorization
probe generated a 196-byte native sync and required an actual nonzero native
status (-42018) for invalid diversity data; an interpreter error cannot satisfy
that test. Closed-session checks and existing VFS tests also passed with -race.
No authorization request or real diversity import was performed by these tests.
