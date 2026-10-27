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
