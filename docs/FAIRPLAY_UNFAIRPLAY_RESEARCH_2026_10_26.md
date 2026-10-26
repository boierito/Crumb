# Crumb FairPlay: architecture and causal research, 2026-10-26

> Engineering evidence record; authorized, account-owned IPAs only. No private enrollment state, keys, SIDB/SINF data, original/decrypted IPA, hardware identity, selectors, account identifiers or hashes are stored in this document or in public CI output.

## Requirement, explicitly preserved

Entire workflow **on an ordinary iPhone**: no jailbreak (even once), no JIT, no privileged process memory or private iOS system files, no Mac or PC enrollment, no external decryption service, and no third-party account credential. This is an unsolved feasibility condition, not an already implemented capability.

## Primary upstream studied

**DrWhax/unfairplay**, pinned for this study to source commit [12bf939d38f28279de14cb5d318257f2092f2daf](https://github.com/DrWhax/unfairplay/tree/12bf939d38f28279de14cb5d318257f2092f2daf). The source incorporates a complete macOS native-enrollment collector, canonical input normalization, portable recovery and an offline IPA decrypter.

Relevant source, not just README:

- [internal/fairplay/enrollment_normalize.go](https://github.com/DrWhax/unfairplay/blob/12bf939d38f28279de14cb5d318257f2092f2daf/internal/fairplay/enrollment_normalize.go): contract `macos-current-sinf-v1`. `canonicalSIDB` treats the big-endian `0x0007000b` header as `macos-sidb-v7000b`; known byte counts 630 and 642 correspond to 192-byte sources, 993 and 1003 to 224-byte sources. **A recognized header with another length is itself unclassified**, and an unrecognized header is not automatically invalid data.
- [internal/recovery/raw_source_step.go](https://github.com/DrWhax/unfairplay/blob/12bf939d38f28279de14cb5d318257f2092f2daf/internal/recovery/raw_source_step.go): native 26A428 semantic reconstruction: account parser, bootstrap, outer writer, selection and SINF key material extraction from a build-specific frame.
- [internal/recovery/execution_profile.go](https://github.com/DrWhax/unfairplay/blob/12bf939d38f28279de14cb5d318257f2092f2daf/internal/recovery/execution_profile.go): build-dependent pointer offsets, code semantics and SINF masks. Support for 25G83, 26A428/26A434 and 26B118 is based on reference contracts, **not universal equivalence between different producers' SIDBs**.
- [internal/fairplay/current_sinf.go](https://github.com/DrWhax/unfairplay/blob/12bf939d38f28279de14cb5d318257f2092f2daf/internal/fairplay/current_sinf.go): `FPSENR1/2/3` model envelopes; V2 (4,184 bytes) selector and per-position input byte substitution; V3 additionally includes a different 432-byte ciphertext representation and codec. Both forms require actual SINF plaintext validation. Structurally matching a model is not sufficient.
- [internal/fairplay/provider.go](https://github.com/DrWhax/unfairplay/blob/12bf939d38f28279de14cb5d318257f2092f2daf/internal/fairplay/provider.go): model registries may contain more than one matching selector and must select using **actual SINF decryption**, not selector equality alone.
- [README](https://github.com/DrWhax/unfairplay/blob/12bf939d38f28279de14cb5d318257f2092f2daf/README.md): requires an *enrolled Mac*, its IOMACAddress, installed license, system FairPlay state and IPA from the same enrollment. Extracting a model on another host is allowed after genuine Mac collection, but still violates Crumb's stand-alone, iPhone-only acquisition constraint.

UnfairPlay does **not** demonstrate that merely importing an App Store keybag into a virtual, earlier StoreAgent implementation is equivalent to the state produced by an enrolled modern macOS FairPlay daemon.

## Confirmed Crumb execution: what passed and what failed

Real owner-device, 2026-10-26, last tested release `test-23045-fourway-sinf-probe`:

| Step | Result | Interpretation |
| --- | --- | --- |
| App Store entitlement/KBSync and keybag | Native StoreAgent consumer returns zero | Imported by StoreAgent 10.9 emulator; no guarantee of modern enrollment |
| Virtual StoreAgent SIDB | One nonempty state file, variable observed lengths: 863, 869 etc. | **Not any of UnfairPlay's four validated Mac geometries**; contents/header have not been inspected or disclosed |
| 26A428 execution | 13,683,433 generic instructions, no JIT | Semantically executed this path without a fatal guest error |
| SIDB file access | All four SIDB path requests resolved | A missing SIDB *path* is not the blocker |
| Source/Frame | 224/16 bytes, transformed frame matches one of 14 slots | Internal consistency **only**; not cryptographic provenance |
| FPSENR2 candidate | 4,184 bytes, accepted envelope and selector | Synthesizing an envelope is not proof of AES correctness |
| Real SINF decryption | `recovered_model_sinf_decode_failed` | **Decisive blocking check**, original IPA retained |
| Four state-file mount modes | dual, fpsd only and SIDB only all fail real SINF; shared only aborts bootstrap | Changing aliases or deleting sidecars is not a solution |

Do **not** identify 863 bytes alone as proof that the StoreAgent SIDB is invalid: variable-format containers could exist. The error establishes only lack of demonstrated **cross-generation equivalence**, not which specific instruction or key is wrong.

## Why the previous direction was stuck

The original path feeds an opaque StoreAgent 10.9 SIDB and a six-byte **software-generated Keychain identity** into UnfairPlay-derived `fairplayd` 26A428 semantics. A successful `open`/mmap and `native_status=0` on the producer do not prove that the consumer parses the intended enrollment or derives the correct cryptographic material. The SINF parser and selected model are already functioning. Repeated IPA downloads or alternate SIDB aliases cannot repair a missing producer-consumer contract.

Independent public alternatives (`EthanArbuckle/fairplay-decrypt`, `londek/ipadecrypt`) implement an offline decrypter given a reusable 16-byte FairPlay account credential. Their documented credential-acquisition stage uses a jailbroken device. Those extraction methods **are not** a no-jailbreak solution and will not become a required Crumb dependency.

## Implemented new research, not another route comparison

Files: `MapleSyrup/NativeSAP/fairplaycausal.go`, `fairplaycausal_test.go`, `fairplaylabcore.go`, `WaffleStore/App/FairPlayLabView.swift`.

Native operation `research-causality` is an *offline-only, bounded causal experiment* on **one already retained and integrity-checked transaction** and its original IPA. It never downloads another app, requests entitlement/keybag, or creates a decrypted IPA:

1. Verify the original capture manifest and the SINF-to-IPA binding; load opaque SIDB and original virtual identity in private memory.
2. Classify SIDB conservatively using UnfairPlay's **exact** known Mac header and length matrix: `macos_v7000b_known_192`, `macos_v7000b_known_224`, `macos_v7000b_unclassified_size`, `other_header`, `unclassified`. Only this *family label* and byte count may leave the process.
3. Run one genuine baseline `rawsource.Recover` with the same input, and apply independent frame and strict SINF validation. This does **not** decrypt IPA pages.
4. Run four experiments, each changing a single bit in a separate **temporary SIDB copy**: first byte, boundary after first word, middle byte and last byte. A fifth experiment changes the first bit of a copied **virtual** hardware identity. The original bytes never change.
5. For each experiment compare the newly generated source (224), frame (16) and model envelope against the genuine baseline **within the same process**, reporting only three `changed` booleans and fixed execution/validation statuses. No addresses, original or mutated bytes, hashes, model material, selectors, paths, private exceptions or plaintext leave the process. Temporary states are zeroed after use.
6. The app includes **Research SIDB causality (no download)**, reusing its existing capture. Unlike the previous button, there is no download-plus-fourway-probe shortcut.

Crumb's synthetic host tests check the upstream format matrix, frame dependence using the real public transforms, all causal labels, private copy integrity, fixed error categories and prohibition on private data in JSON. These tests prove the research instrument works; **they cannot prove a key was recovered**.

### Correct interpretation of later experimental outcomes

| Observation | Next inference |
| --- | --- |
| SIDB bit change changes Source, Frame or model | Some bytes influence the modern algorithm; still **no proof** of correct state format, real key or enrollment |
| All four sampled SIDB mutations leave outputs unchanged | These **sampled** bits are not observable in these outputs; run byte-use/taint instrumentation before concluding that *all* SIDB state is ignored |
| Changing synthetic identity changes the model | Model depends on supplied virtual identity; does **not** establish a hardware-backed Apple identity |
| Trial execution aborts after a mutation | Likely structural/state validation dependency; never replace errors with fabricated values |
| SINF validation passes on any real, unmodified profile | Advance to the **existing full IPA decryption and page/CodeDirectory verification gate**, without overriding fail-closed controls |

A causal observation that an invalid SIDB mutation changes material is **not itself a cryptographic breakthrough**. The experiment is valuable because it can reject or support a concrete mechanistic hypothesis on the original capture rather than accumulating generic status counters.

## Follow-on engineering directions, without changing the iPhone-only constraints

- **Producer analysis:** locate the actual StoreAgent 10.9 keybag import and persisted SIDB dataflow in the existing native emulation. Map native stores/loads, transform boundaries and normalization to the account-specific license. Do not assume modern `fairplayd` runs the same protocol.
- **Bridge proof:** build a reproducible, synthetic reference that demonstrates the exact relationship, if any, between StoreAgent's persisted structure and one canonical Mac enrollment input. If no bridge exists, retire the 26A428 consumer for StoreAgent inputs instead of repeatedly retrying it.
- **Independent account-credential feasibility:** isolate the unprivileged operations possible from an ordinary App Store client and determine whether they can produce a license-bound credential under the software-generated identity. The documented public alternatives cannot do this without external enrollment/privilege.
- **Portable backend:** retain UnfairPlay's strict `FPSENR2/FPSENR3`, selector, decrypted atom and IPA validation gates. Only port additional build/version semantics after there is a genuine enrollment producer; an additional AOT instruction profile by itself cannot invent private keys.
- **Go/no-go:** a full signed original enrollment + target IPA + private model must produce valid SINF *and* complete page and Mach-O validation. Without that evidence, status stays research-only, no decrypted file advertised or exported.

**Security/release rule:** No new testing IPA should be published solely because a counter, label, routing variant or synthetically tested transform changed. Ship the causal study only after native iOS compile and host tests; do not describe it as an operational decrypter.


## V13 owner-device finding, 2026-10-26 19:46 Argentina

Real `research-causality`, baseline with StoreAgent-native accepted keybag (888 bytes), KBSync (196 bytes), SIDB (878 bytes), single 440-byte SINF private region:

- **The SIDB reports big-endian header `0x0007000b`** (our strict `macos_v7000b_unclassified_size` label). The header agrees with UnfairPlay's known Mac SIDB family, but length 878 matches none of its independently checked Mac source geometries (630/642/993/1003). This is positive *structural-family* evidence, **not** evidence of a transferable model.
- Original execution `source_frame_verified=true`, `sinf_model_verified=false`.
- First-byte, offset-4, middle, last SIDB one-bit trials, as well as a one-bit change in the six-byte virtual hardware identity, all reported `status=runtime_rejected`.
- The original V13 instrument reported `source_changed=false`, `frame_changed=false` and `model_changed=false` for errors because it had **no result to compare**. These booleans are invalid to interpret as proof that source, frame or model ignored the changed bit. **No causal conclusion is available from these five failures alone.**

### Implemented correction on PR #12 (after device V13)

- Adds an **exact, unmodified baseline replay** first; `replay_verified=true` only if original Source, Frame, model and strict SINF validation results reproduce exactly.
- If the unmodified replay fails or changes, the operation stops before any mutations: the experiment is **inconclusive because the engine is not reproducible in that session**.
- Only if replay verified, the original five single-bit interventions run. Each failure includes only `failure_phase` and `failure_code` from strict **fixed-string, allowlisted** rawsource diagnostics. Unknown failures collapse to `other`; no exception messages, paths, addresses or private values are emitted.
- Synthetic tests cover unmodified replay success, baseline replay failure, fixed-code sanitization and rejection privacy. This does not supply any new key or claim real IPA decryption.

Research decision on next device report: **if replay fails**, investigate interpreter/model runtime reentry before blaming private SIDB checksums. **If replay passes but SIDB bit changes fail during account parsing**, investigate SIDB integrity, version and parser geometry rather than altering model/SINF validation. **If hardware-only mutation fails**, that is consistent with binding but does not establish which hardware field is cryptographically bound. An unchanged control followed by successful altered runs is the precondition for downstream Source/Frame/Pack causality claims.

These revisions require a new iOS build before the additional `replay_verified` and `failure_phase`/`failure_code` fields can appear in logs; V13 used the earlier app and its result cannot be retroactively reinterpreted.


## V14 on-device causal replay: verified deterministic engine, divergent unsupported imports (2026-10-26)

Owner-device capture from the released `test-23045-sidb-causality-replay` version:

- Keybag 904 bytes; KBSync 196 bytes; SINF private 440 bytes; StoreAgent 10.9 import succeeded with native status 0.
- SIDB 882 bytes, header `0x0007000b`. Header matches macOS format-family classification, but its source geometry is **not** one of the known 630/642/993/1003-byte native macOS contracts.
- Baseline source/frame coherent (`source_frame_verified=true`), strict SINF model validation fails (`sinf_model_verified=false`).
- **`replay_verified=true`: the same inputs reproduced Source, Frame, Pack and the same strict SINF result in an independent invocation. This rules out simple failure of state reset between invocations in this session.**
- Single-bit modifications:
  - `sidb_first`: `bootstrap/raw_unsupported_external`.
  - `sidb_header_end`, `sidb_middle`, `sidb_last`, `hardware_first`: `selection/raw_unsupported_external`.
  - All trial outputs `source_changed/frame_changed/model_changed=false` are **not real comparisons** because the trial never completed.

### Source-level root of `raw_unsupported_external`

The pinned UnfairPlay `internal/recovery/raw_source_step.go` `rawSourceRuntime.handle(s, pc)` has an explicit set of modeled external entrypoints (file `open/read/mmap/fstat`, selected IOKit/CF operations, allocator, pthread once and locks). Its `default:` **rejects** other calls with `unsupported bootstrap external %#x` rather than impersonating unknown system APIs. Crumb copies this implementation via `scripts/prepare-fairplay-raw.py`. `MapleSyrup/NativeSAP/fairplayrawapi.go` converts that error to the fixed category `raw_unsupported_external`, attached to the active `runEnrollmentSource` phase. Thus the two observed phases are genuine **semantic execution divergence into an unmodeled call**, not failed source/frame equality.

The currently redacted report cannot tell whether those calls are *legitimate, unmodeled imports* used only for rejection/cleanup or *invalid branch targets* triggered by modifying a protected serialized state. **Do not add blind `return 0` shims for all unknown calls.** Doing so could synthesize a fake model and weaken the strict SINF/IPA rejection gate. The 882-byte SIDB remains a **structurally related but cryptographically unvalidated** source, not a converted macOS enrollment.

### Engineering decision

Arbitrary one-bit mutations of opaque SIDB or device identity are now **exhausted as useful end-to-end trials**: once execution aborts, their downstream output comparison cannot answer key derivation. Work on a *valid-state* producer/consumer bridge instead: map the StoreAgent keybag import side effects and the SIDB format and authenticate the target SINF via the real validator. Avoid additional downloads and mutation-only IPA builds until a specific, independently tested import semantics or enrollment transform has been implemented.

The result does **not** support deriving a FairPlay decryption key on an ordinary iPhone today. UnfairPlay requires a separately enrolled native Mac source; our iPhone sandbox does not automatically provide equivalent private account/hardware material.


## Producer/consumer implementation audit — commit f78dbc91

### Confirmed implementation boundaries, not assumptions

**Producer (StoreAgent 10.9 guest on iPhone):**
- `MapleSyrup/NativeSAP/kbsyncretained.go`: `OpenRetainedKBSync` loads the existing pinned native StoreAgent image, establishes one virtual six-byte MAC/DSID identity, obtains KBSync, then `ImportKeybag` invokes the observed native entrypoint `storeAgentBase+0xd5ce0` on the *same retained guest*. It snapshots the virtual filesystem only after the return.
- `MapleSyrup/NativeSAP/gueststate.go`: the producer's virtual `SC Info` filesystem implements open/read/write/stat/lseek/fcntl/catalog metadata and a bounded 64-operation path-free trace. The real device snapshots show SIDB and an empty companion/lock file but no populated `SC Info.sidd`. An accepted native return code is **not** a completed FairPlay enrollment or decryption key.
- `MapleSyrup/NativeSAP/kbsyncbridge.go`: keybag is consumed only once; state files and category-only I/O metadata are persisted inside the protected private capture. The consumer does not synthesise modern daemon state.

**Consumer (UnfairPlay fairplayd 26A428 translated runtime):**
- `scripts/fetch-fairplay-raw.sh` pins the upstream source; `scripts/prepare-fairplay-raw.py` uses its original external-call switch and generated CFG and adds the verified generic instruction profile, not native macOS execution.
- Pinned `raw_source_step.go` supports specific file, IOKit/CF, pthread, and memory external-call stubs; its unknown external entrypoint explicitly returns an error. Five one-bit device experiments changed the control flow to these unsupported endpoints, making any missing Source/Frame/Pack comparisons inapplicable.
- `MapleSyrup/NativeSAP/fairplayrawapi.go` now preserves the **static public import slot** for this exact class of failure, only if its PC is within the pinned model's synthetic `F100000000000000` import space, 0x1000-aligned and within 4096 possible entries. All private pointers, guest registers, raw unsupported PCs, call arguments and return values remain unreported. The execution *still fails closed*; no generic externals are replaced by no-op stubs.
- For each of `account-parser`, `bootstrap`, `source`, `selection`, and `sinf`, aggregate virtual state-file lookup/missing counts are saved. These establish whether the daemon actually attempted to consult SIDB/SIDD, not whether cryptographic decryption succeeded.

### Native bridge auditor implementation

`MapleSyrup/NativeSAP/fairplaybridgeaudit.go` implements `research-bridge` via the **original**, checked capture and exactly bound original encrypted IPA. It reads the existing StoreAgent state *without modifying it*, reports file presence and conservative upstream SIDB family/geometry, executes actual UnfairPlay recovery once, reports only phase-level state-file lookup coverage and sanitized static external import slot when failed, and requires `verifyRecoveredEnrollment` **plus strict `probeModelSINF`** to claim a usable model. All transient buffers are cleared. It cannot export a model or key or authenticate plaintext on its own.

Added a specific on-device control **Audit StoreAgent → fairplayd bridge (no download)**, native Go routing, host + iOS build integration and tests for simulated 882-byte legacy SIDB, missing SIDD, no mutation of captured state, strict model rejection, unknown runtime error privacy, bounded import labels, and immutable phase trace snapshots.

### What is not implemented / cannot be asserted

- No evidence establishes a legitimate macOS 26A428 FairPlay enrollment model from the legacy StoreAgent 10.9 virtual SIDB alone. Their differing generations, state paths and 882-byte geometry prevent claiming an exact portable conversion merely from matching `0x0007000b`.
- No source-verified modern SIDD construction, modern account-key provisioning, license-bound model derivation, or authenticated executable plaintext is established. In particular **do not fabricate SIDD**, claim keybag import is a license key, or ignore SINF validation.
- The exact function behind a rejected external import slot can only be implemented once mapped to an independently identified public/binary import symbol and ABI plus a relevant reproducible test; treating an unknown call as success is prohibited.

The bridge audit is a real implementation of a safer *producer/consumer contract verifier*, **not** the fully functioning DRM decryption solution requested. All public fixture and iOS tests must pass before a new release is described as validated.
