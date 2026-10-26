# Build 23045: recovered-output validation diagnostics

The 23044 iPhone report passed the generic-profile self-test and returned from
rawsource.Recover after approximately 63 seconds. Its next stage failed with
`recovered_output_validation_failed`. This establishes that this run passed the
runtime output extraction path, but does not establish correct execution,
legacy/modern enrollment compatibility, a valid enrollment model or decryption.

The old wrapper discarded the distinction between state evaluation, frame
matching, filesystem failures and target-IPA validation. It also lost the CPU
execution counters on successful raw recovery. No private capture was received;
the exact underlying rejection cannot be determined from that report alone.

Build 23045 reports all of those gates in one normal recovery attempt:

- Fixed `recovered-validation/*` stages and codes distinguish missing outputs,
  state evaluation, frame evaluation, frame mismatch, private output directory,
  model directory/write, and downstream decryption gates.
- Typed downstream stages identify platform, paths, inventory, model/SINF
  binding, archive reads/writes, SUPF, image geometry, CodeDirectory,
  page decryption, ARM64 plausibility, mutation bounds, sync and publication.
  An inner image failure keeps its specific stage through outer wrappers.
- `validation` includes lengths, frame candidates checked, domains cross-checked,
  source/frame agreement, model files written and final target validation.
- `runtime` preserves lookup counts, execution backend, instruction/block counts
  and public profile offsets after the raw interpreter has completed.
- The self-test includes `recovery_diagnostics_version: 2`.

Diagnostics contain no private state, source/frame/model bytes or hashes,
account identifiers, archive member names, filesystem paths or underlying error
strings. Successful prior gates do not change the final success criteria.
Models and processed output remain unpublished on any validation failure.
No CPU instruction semantics, file routing, identity, model selection or
cryptographic validation is changed by this build.

Validation uses public synthetic state/frame inputs through the real
independent evaluators, forced frame mismatch, invalid source, preexisting
output directory and an invalid ZIP after a successful frame cross-check.
Tests assert the exact failure class, retained CPU counters, no secret canary
in the public envelope, wiping of recovered buffers and output cleanup.
Existing decryption tests continue to exercise the unchanged validation gates.

Device test: install build 23045, confirm self-test diagnostic version 2,
then run normal on-device recovery once. Share only the Components, Enrollment
and Decryption logs. No four-way routing test or private capture export is
required to identify the failing gate.
