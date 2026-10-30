# Original iOS daemon: Mach dispatch and caller-identity follow-up

Date: 2026-10-30. This is a bounded host replay of original 24A446 ARM64 code,
not an iPhone service invocation or a decryption correction. No private capture,
account, SINF, enrollment credential or purchased IPA is used by the replay.

## What changed in the investigation

The two entitlement helpers left unresolved by the
[static service audit](FAIRPLAY_SYSTEM_STORAGE_CONTRACT.md) now have a concrete
observed caller path: the native Mach 1000 handler reaches the numeric helper
with **`fairplay-client`**, then, when that entitlement is modeled as absent,
reaches the string helper with **`application-identifier`**. The name strings
are reconstructed by the original obfuscated code; they are not supplied as
arguments to the request or selected by patching a branch.

This rules out treating a missing `fairplay-client` result, by itself, as a
demonstrated rejection of every operation on that path. It does not establish
that any useful operation is authorized for Crumb. Security.framework responses
in this experiment are explicit fixtures, not measured signing identities.

The Mach 500 handler follows a different path: the original code extracts the
PID from the audit token, tests it against zero and rejects the tested positive
PIDs with **−42049**. This is distinct from the portable recovery's **−42057**
and the phone's `mremap_encrypted` **errno 12**. No relationship between those
errors is established merely by locating this check.

## Original image, dispatch and wire boundaries

The original daemon is the same independently acquired firmware input as the
static audit: SHA-256
`94d32a82fdc2a80c46b7734b20895a935836a676a34b1c3e863cdb479a68cb01`.
Its authenticated fixup records expose three MIG routine-table families:
1200–1203, 1000–1002 and a 500-range table. This investigation executes the
**500 and 1000 handlers**; it does not assign every family to a named launchd
service or prove the XPC service has the same authorization contract.

| Property | Request 500 | Request 1000 |
| --- | --- | --- |
| Native handler | `0x1000029e0` | `0x10000255c` |
| Request size | `0x50` | `0x3c` |
| Required complex bit | yes | yes |
| OOL descriptors | one | one |
| Native worker | `0x100761424` | `0x100005940` |
| Worker return observation | `0x100002b14` | `0x100002694` |

The replay supplies a valid **transport layout** and a deliberately invalid,
non-license OOL buffer and client word. Its contents are not a valid FairPlay
client request. An initial shared builder incorrectly reused request 500's
repeated-count field at +0x34 for request 1000. This is corrected: request 1000
has a separate four-byte in/out client word at +0x34 and its OOL count at +0x38.
The investigator now refuses a non-mutated request 1000 unless its **invalid
client-word fixture** is explicitly selected. A valid producer is still missing.
In particular, using a pathname-shaped dummy buffer for request 500 does not
establish that the native worker interprets that buffer as a filename.

The original handlers reject changed transport fields before reaching their
workers: a wrong message size returns MIG −304, an incorrect OOL descriptor
type or mismatched count returns −300, and an invalid trailer format returns
−309. These controls distinguish transport rejection from the subsequent
native worker result.

## Request 500: bounded PID-check evidence

At `0x100b2bd2c`, the original worker invokes `audit_token_to_au32` with only
the PID output pointer non-null. The helper's published implementation copies
`audit_token.val[5]` to that output. At `0x100b2bd30`–`0x100b2bd58`, the worker
tests the copied PID and selects an obfuscated branch target.

- PIDs 1, 123, 501 and 1000 each selected `0x1008f5410` and returned −42049.
- PID 123 still selected that branch and returned −42049 after the optional
  original initializer had run; the result was not caused by its missing
  allocation callback.
- An explicitly labeled **counterfactual** PID-zero fixture selected
  `0x100965eb0` and stopped at the first unmodeled `pthread_self` call. It did
  not produce a successful worker return, service authorization or decryption.

The kernel, not the sender's ordinary payload, supplies a real Mach audit
trailer. Apple's
[ipc_kmsg_init_trailer_and_sign](https://github.com/apple-oss-distributions/xnu/blob/f6217f891ac0bb64f3d375211650a4c1ff8ca1ea/osfmk/ipc/ipc_kmsg.c#L528)
uses `task_get_audit_token(sender)` for user senders and a kernel trailer
template for kernel-originated messages. Changing a field in a host fixture
does not provide that identity to Crumb. This source check is not a claim of
instruction-level equivalence with every path in the extracted iOS kernel.

This establishes a caller-PID restriction for the **tested request-500 path**.
It does not establish all operations on all FairPlay interfaces as privileged,
nor identify the failing internal operation behind the phone's remap error.

## Request 1000: alternative identity path and initialization dependency

With the missing-entitlement fixture, the original worker reaches:

1. `0x100001b40`, querying numeric `fairplay-client`;
2. `0x100001cec`, querying string `application-identifier`.

The second query is a real alternative in the original control flow. The test
supplies the deliberately fictional ordinary identifier
`TESTTEAM01.com.certlium.crumb`; it does not impersonate an Apple application,
grant an entitlement or claim to reproduce the phone's installed signature.

The first direct-handler exploration then encountered a null allocation-context
slot at `0x1011320f8`. Original function `0x100777114` read that slot and returned
null; continuing without proper initialization led to an unmapped write. That
host failure is **not** a device rejection or evidence of an inaccessible
credential. No dummy pointer was inserted into the slot.

The exploratory main-prefix trace located an original `pthread_once` callback,
`0x10073a790`, at the call through `0x1007854b4`. Replaying that original
initializer with the explicit host runtime models below returned after
6,302,431 instructions and populated the allocation context. Its final x0 is
not interpreted as a success/failure status: it is a once-initializer callback.
This is initialization of that component, not execution of the complete daemon
startup or establishment of an account/hardware context.

After this initialization, the same worker followed both identity queries,
allocated its working buffer through the now-initialized original allocator,
and returned **8** for the deliberately invalid 32-byte input. Execution then
stopped at the unmodeled `mig_deallocate` call in the handler. A completed
worker and a completed MIG reply are reported separately. Status 8 is not a
validated license, proof of client authorization or a known solution to −42057.

### Where the observed 8 actually originates

A complete 16,549-boundary register trace of the initialized request establishes
that the original code made **no native reads of the OOL input buffer** on this
path. Tracing the returned status backward found:

1. `0x100005b98` reads the request's +0x34 word; `0x100005b9c` copies its value
   `0x20` to the worker frame at `0x70017ea0` in this host layout.
2. `0x10083e76c` reads that same word, and `0x10083e770` masks it with `2`.
   The fixture's bit is clear. The executed branch selects the cleanup path.
3. `0x10083e798`–`0x10083e79c` set encoded status `0xc802ee72`, subsequently
   stored at context+0x24 by `0x1007faa98`.
4. The epilogue copies that word to the wrapper's status slot. Finally
   `0x100005d28` XORs it with `0xc802ee7a`, producing `8`.

Thus this observation diagnoses the **invalid research request**, not a rejected
SINF or a demonstrated entitlement denial. The in/out word is not interchangeable
with a length. Its source and lifecycle must be recovered from the matching
client. The investigator does not set the observed bit to bypass this check.
The neighboring Core ML MIG client's published decompilation also has a separate
in/out integer and length, but that different family is not adopted as the
missing mobile-license encoder.

## Replay models and verification limits

[`scripts/audit-fairplayd-mig.py`](../scripts/audit-fairplayd-mig.py) loads the
original image and a separately pinned ipsw fixup report. All 56,540 listed raw
fixup words are checked against the original memory before relocation. The
report's SHA-256 is
`6c5549d2cf4ec6cc7ba0beb841e4c05bf852caa8cc8d5582ac62e88c89a718a8`.
Regenerating it with ipsw v3.1.739 and `--no-color` was byte-identical.

The following are **models**, not measured Apple services:

- A 16 KiB page size, bounded anonymous read/write mappings, bounded allocations,
  and single-thread lock initialization. Allocated arenas are non-executable.
- A fixed clock and public deterministic startup entropy. Neither is an account
  credential, an enrollment secret or a proposed key for the user's license.
- Bounded stack probing, following Apple's arm64 `chkstk` x9-size convention.
- Audit-token field copies according to OpenBSM, opaque CFString/SecTask objects,
  an absent `fairplay-client`, and the explicit ordinary-identifier fixture.
- PAC operations as control transfers, without cryptographic authentication.
  Original image page protections are not enforced by this host model.

Every modeled external call is listed in the output. Unknown imports and unknown
entitlement queries stop execution. No native branch or instruction is patched.
Request execution is bounded independently of the initializer. Neither a budget
stop nor an emulation error counts as a native or cryptographic rejection.
Optional `--trace-request <local.json.gz>` retains at most 32,768 request
boundaries with registers, excluding the initializer; truncation is explicit.
The inspected 16,549-boundary trace was complete. It contains only public
firmware-derived state and these synthetic fixtures, and is not committed as an
account/enrollment artifact. Input-read and status-write provenance are reported
separately from any native result.

Primary runtime references:
[OpenBSM token copies](https://github.com/apple-oss-distributions/OpenBSM/blob/bb0681420e99d447f466ba70d38589213bb69fea/openbsm/libbsm/bsm_wrappers.c),
[Apple stack-probe assembly](https://github.com/apple-oss-distributions/libpthread/blob/42d026df5b07825070f60134b980a1ec2552dfee/src/pthread_asm.s).

Validation completed: four positive-PID controls, eight original-handler
transport-rejection controls, labeled PID-zero branch comparison, and the
initialized 500/1000 replays. Altering either pinned input was rejected even
under Python `-O`; PID zero without its counterfactual label was rejected.
The anonymous bounded results are in
[research/fairplayd-mig-24A446.json](research/fairplayd-mig-24A446.json).

## Reproduce and next boundary

Acquire the separately pinned public daemon, then generate the derived report:

```sh
ipsw --no-color macho info --fixups /research/fairplayd.H2 > /research/fixups.txt
python3 scripts/audit-fairplayd-mig.py \
  --daemon /research/fairplayd.H2 --fixups /research/fixups.txt \
  --request-1000 --ordinary-identifier TESTTEAM01.com.certlium.crumb \
  --prepare-runtime --fixture-invalid-client-word
```

Use `--request-500` for the PID control and `--mutation size`, `ool-type`,
`ool-count` or `trailer` for the transport controls. The local initializer has a
20-million-instruction ceiling; the request defaults to 300,000. The observed
initialized requests did not exhaust either limit. Dependencies are Unicorn and
Capstone on the research host; neither is added to the iPhone runtime here.

The next required producer is a **matching client-side request encoder**,
including its real in/out word, and an
operation that consumes the original mobile license. A correctly laid-out MIG
message with arbitrary bytes does not supply that protocol or valid crypto
context. The four-message Core ML family also cannot be substituted for it
merely because it lives in the same daemon. Identifying a real license-consuming
operation, its input sources, authorization and native output validation remains
necessary before an integrated device build is justified. Crumb's runtime,
private-data handling and SINF/export acceptance gates are unchanged.
