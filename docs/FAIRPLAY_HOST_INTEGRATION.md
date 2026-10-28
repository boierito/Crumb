# Production bridge integration outside the app

Build the portable production bridge without SwiftUI or the iOS C bridge:

```sh
bash scripts/build-fairplay-host.sh /private/new-build-directory
/private/new-build-directory/crumb-fairplay-host \
  --ipa /private/original.ipa \
  --capture /private/capture.zip \
  --output /private/new-processed-directory > /private/result.json
```

The build directory must not exist. Keep it and the inputs private. Building
may download Go dependencies; the capture replay itself does not contact Apple
or send inputs to a server. Do not submit private files to public CI.

This executable copies the same production `labRun` implementation and recovery,
validation and IPA processing packages used by Crumb. It calls selftest,
inspection, capture analysis, exact capture/IPA binding, and capture recovery
in order. The first failure stops the pipeline with exit status 1 and fixed
diagnostic fields. Source/frame and SINF validation remain mandatory. Only
successful validated recovery proceeds to the real IPA writer for all images.
The selftest is explicitly synthetic; it is not evidence of real decryption.
Models are temporary, and plaintext authentication is not asserted.

The host harness does not test Swift Keychain, UI lifecycle, signing,
installation, live Apple authorization, or the retained StoreAgent producer.
Those boundaries are explicitly reported instead of being mocked as successful.
An original capture predating authorization cannot reproduce the later account
state saved on an iPhone. The harness therefore does not claim to reproduce
that state or the complete iOS operating environment.

## Exporting a saved state from Crumb

After the native authorization operation reports that account state was saved,
FairPlay Lab prepares a separate private JSON copy, even if SINF recovery fails.
The persistent internal state is retained. The laboratory also offers
`Prepare PRIVATE authorized state (no new authorization)` to copy an existing
saved state without another request to Apple, followed by an explicit share
sheet. No upload is automatic. The file contains account-bound licensing data
and must remain private; the export schema does not contain login passwords,
Apple session cookies or second-factor codes. Export alone does not prove that
provisioning or SINF validation completed successfully.

Replay the exported file with the same capture identity and original IPA:

```sh
chmod 600 /private/PRIVATE-FairPlay-authorized-state.json
/private/new-build-directory/crumb-fairplay-host \
  --ipa /private/original.ipa --capture /private/capture.zip \
  --authorized-state /private/PRIVATE-FairPlay-authorized-state.json \
  --output /private/new-output-directory > /private/result.json
```

The loader verifies the saved account, hardware identity, schema, file hashes
and permissions. An explicitly supplied missing or invalid state fails rather
than silently falling back to the original capture. Every candidate must still
validate the target SINF before IPA processing can proceed.

`test-fairplay-host.sh` also compiles the host entry point after exercising the
portable contracts, so future changes must continue to build outside the app.

## Real original-capture findings

The matching uploaded IPA/capture passed inventory, capture integrity and exact
SINF binding. Recovery rejected its one candidate. After the exact native
epilogue wrote the four caller output scalars, it returned status -42057 to the
snapshot's null saved link register. Recognize this narrow return contract as
`native_entry_failed`, not as an instruction missing at address zero. A zero
status through this return still does not establish SINF validation and remains
an error. The meaning of -42057 has not been independently identified.

An additional replay restored the original capture in the real pinned
StoreAgent using the host TCI backend, with no JIT and no Apple authorization
request. Its native account query completed with `needs_account_diversity=true`.
This proves a missing provisioning condition in that original state, not the
cause of every recovery failure or impossibility of the project. It does not
test the later state saved by authorization on the phone.

The optional native test is reproducible through `scripts/test-sap-host.sh` by
setting `CRUMB_PRIVATE_CAPTURE` to a local capture ZIP and
`CRUMB_PRODUCER_REPORT` to a private output JSON file. Without private inputs it
is skipped. It does not contact Apple with account data or extract credentials.
Do not run private replays in public CI.

## Replaying the saved authorized state in the original producer

Set `CRUMB_PRIVATE_AUTHORIZED_STATE` alongside `CRUMB_PRIVATE_CAPTURE` and
`CRUMB_PRODUCER_REPORT` when running `scripts/test-sap-host.sh`. The optional
producer test uses the same production JSON validator as the portable replay:
it checks the account/hardware binding, schema, file hashes and private file
permissions before restoring the files in StoreAgent. It records whether an
authorized-state file was supplied. It does not perform another authorization,
fetch a license, or validate SINF. A timeout before the native account query
completes is an incomplete execution, not a negative provisioning result.
Run this private replay separately from large compiler workloads.

## Authorized-state replay, build 23053

A user-exported state passed the loader and exact original-IPA/capture binding.
It contains a 3,036-byte SIDB and 200,356-byte SIDD. Replaying the unchanged
production pipeline outside the app reproduced the device's failure exactly:
`generic_instruction_limit`, 50,000,000 interpreted instructions, last generic
PC `text+0x4886e0`, one rejected SINF candidate, and no output IPA. The
original producer's pre-authorization state and this later state are different
inputs; conclusions from the former cannot be substituted for the latter.

A serial replay in the original StoreAgent restored the authorized state and
completed its native account query with `needs_account_diversity=false`, with
no Apple authorization request. An earlier restore attempted during a large
compiler workload timed out before producing a comparable query result.

A code-address-only profile of the unchanged recovery attempt found 174,476
entries at `text+0x48866c` (132 instructions) and 12,521 entries at
`text+0x234cdc` (2,099 instructions). Together they account for approximately
99% of the 50-million instruction budget. Three further observed blocks at
`text+0x128818`, `text+0x5ea024`, and `text+0xa06948` complete the batch of
five AOT translations. The differential test compares ten public blocks across
120 synthetic initial states with independent ARM64 execution, including
nonzero scalar state and lookup tables for the 132-instruction block, registers,
flags, touched pages and returned PC. Execution bounds and original-SINF validation are
unchanged. This comparison establishes block semantics, not a working license
or successful whole-program recovery. Private inputs are not CI fixtures.

The AOT-only replay completed native execution instead of exhausting the
instruction budget. It checked two distinct candidates; neither validated the
original SINF. It returned `native_entry_failed` with status -42057 after
10,859,709 generic instructions, and created no output IPA. This does not
identify the meaning or cryptographic cause of that native status.

A private 30-second CPU profile located another cost: 93% of sampled CPU time
was inside the auxiliary source/frame verifier; loading and decompressing the
public frame model accounted for 84%. The public state/frame evaluator APIs now
load their immutable model template once, while each evaluation keeps a separate
copy-on-write CPU/memory state. Templates that need private enrollment patching
are explicitly rejected by this cache; the SIDB evaluator continues to load its
own per-call model. Concurrent evaluations are compared with a fresh-model
reference and must leave both caller inputs and the complete public model
unchanged. The original SINF and source/frame checks are still performed.

The signed-dispatch regression runs with both the configured dispatcher and a
forced generic fallback, including negative table indices, after promoting
its formerly unknown PC to AOT. CI's host CLI smoke build uses the same race
instrumentation as the preceding suite to reuse its large compiled packages.
