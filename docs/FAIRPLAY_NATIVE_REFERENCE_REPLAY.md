# Replaying a real saved-state attempt with independent ARM64 execution

This host-only tool records one invocation of the production portable runtime
and replays its original ARM64 instructions in Unicorn. It compares integer
registers, flags, stack pointer, PC, write coverage and written memory at every
service boundary and the final return. The app's recovery gates are unchanged.
It is intended to distinguish instruction/interpreter divergence from problems
in the modeled services or the producer/consumer contract.

## Scope and limits

The reference uses recorded service replies after verifying the native inputs
and memory preceding each call. It does **not** independently validate those
service implementations. PAC instructions follow the portable runtime's existing
convention, so pointer-authentication correctness is not claimed. Original
instruction bytes are restored from the hash-pinned public image before this
normalization. No SINF validity or plaintext authenticity follows from a match.
The four missing caller-owned output pages are accepted only at their exact
pinned return-store instructions and sizes. Other unmapped accesses fail.

Use a serial host invocation while recording. The first native invocation is captured. A stopped or unsupported invocation is
marked incomplete, not accepted as a completed replay. Later fallback attempts
or license domains do not append to that first trace. The normal host pipeline
still checks every domain independently before producing an IPA.

## Build, record and replay

Install `unicorn==2.1.4` and `capstone==5.0.6` for the reference process.
Build into a new private directory:

```sh
CRUMB_NATIVE_REFERENCE=1 bash scripts/build-fairplay-host.sh /private/new-host
CRUMB_NATIVE_REFERENCE_TRACE=/private/real-attempt.jsonl \
  /private/new-host/crumb-fairplay-host \
  --ipa /private/original.ipa --capture /private/capture.zip \
  --authorized-state /private/authorized-state.json \
  --output /private/new-output-directory > /private/production-result.json
python3 scripts/replay-fairplay-native-reference.py \
  --trace /private/real-attempt.jsonl > /private/reference-result.jsonl
```

The trace contains private registers and memory, including licensing material.
It is created with mode 0600 and exclusive creation; keep it local and never
submit it to GitHub, public CI, or a release. The reference emits comparison
results and public code labels only. Build dependencies may be downloaded;
recording/replay do not authorize an account or contact Apple. Native reference
execution is a research tool on the host, not an iPhone runtime dependency.
The iOS build never installs these recording hooks.

## Real-input findings

The original IPA, matching capture and exported post-authorization state were
reused locally. No new download, account password or second-factor code was
used. Results:

- The original producer restored the state and reported diversity not required.
- The production recovery checked two candidates. Both reached private-payload
  decryption, then failed on plaintext atom 0 with invalid bounds.
- The branch selecting the encoded error was reached at text+0x4cb5c0; its
  comparison was false and it transferred to text+0x8c7ed8. The compared word
  was last changed by the operation at text+0x4cce4c. This locates a dataflow
  dependency, not the semantic meaning of the error.
- Independent execution of that 25-instruction block on the actual private
  snapshot matched registers, flags, SP, target PC and captured pages.
- The complete native replay matched all 293 service boundaries and final
  written memory/return under the same service contract. Native status was
  -42057; the original SINF was not validated.
- Dual and fpsd-only state routing reached the same rejection. Shared-only
  failed earlier and is not a comparable completed result. SIDB-only reached
  the rejection after one candidate rather than two.
- The documented version-2 model-key/account-credential relation was tested
  against the real mobile SINF. Both native candidates were rejected by the
  Go mobile consumer and independently by EthanArbuckle's Python implementation.
  This rejects this particular derivation attempt, not all mobile producers.
- Replacing the fixed modeled wall clock with the actual host wall clock did
  not resolve the rejection. A memory-based once-initialization experiment
  executed the native initializer and still rejected both candidates, including
  when combined with the actual wall clock. These experiments do not establish
  a safe production replacement for the modeled services.
- The native Darwin timeval padding correction was also exercised with the
  actual saved inputs, alone and combined with the actual wall clock and
  memory-based once experiment. Both completed with native status -42057 and
  two rejected SINF candidates. The correction fixes a demonstrated ABI
  discrepancy; it did not resolve the licensing failure in these trials.
- A modified-register trace was rejected at its first service boundary; a
  truncated real trace was rejected for its missing final event after 293
  calls. Neither negative case was accepted as a completed match.
- Native stat-buffer read auditing observed size, mode and modification-time
  reads. Both producer and consumer currently model modification times as zero;
  this observation does not establish a missing real timestamp or justify
  inventing one.

The unresolved boundary remains a compatible, source-identified producer of
material that validates the original mobile SINF. A matching native execution
of the existing snapshot does not prove that the snapshot, service contract,
legacy account state and target licensing operation form a compatible system.

## Native Darwin contract measurements

`Native Darwin reference contracts` runs the public C probe in
`scripts/fairplay-native-reference/darwin-contracts.c` on macOS without any
account inputs. It verifies once initialization/publication to eight waiters
and records the once-control, timeval and stat layout. Observations include
16-byte once/timeval structures, a 4-byte microseconds field at offset 8,
untouched timeval padding, and a 144-byte stat structure with size at offset 96.
Opaque once bytes are platform/thread state, not licensing keys; they must not
be treated as a reusable secret or copied from CI into a recovery attempt.
