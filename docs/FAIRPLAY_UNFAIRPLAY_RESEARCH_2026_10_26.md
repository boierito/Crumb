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
