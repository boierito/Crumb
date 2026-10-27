# Original-SINF validation and complete scalar CAS semantics — 23052

## Evidence from the owner's 23051 diagnostic

`diagnostic-55.json` reports both authorization imports completed, one encrypted
executable domain, and no original-SINF validation attempt. Primary recovery
ended with `native_candidate_limit`, phase `sinf`; its generic-only replay ended
at `text+0x9885a0`, phase `account-parser`, with
`generic_unsupported_instruction`. Neither execution validated a mobile license.
The historical last generic PC and saturated AOT count are not fault causes.

The pinned immutable binary word at that replay PC is `0x88ebfe6c`,
`casal w11,w12,[x19]`. The base generic IR contains unsupported entries at all
61 scalar CAS sites. The primary AOT handler already implements that particular
CASAL; this explains why a generic replay could fail earlier without proving
that its earlier arithmetic disagreed with the primary execution.

## General changes

The generic interpreter now implements the scalar ARM64 CAS, CASA, CASL and
CASAL families for byte, halfword, 32-bit and 64-bit operands, including SP,
zero registers and aliased registers. Acquire/release have sequential semantics
in the one serialized private guest CPU. CASP, other unsupported atomics,
misaligned accesses and unmapped memory remain rejected. No instruction is
selected by a special-case PC. Unsupported base IR slots may use these
semantics only with instruction bytes from the already verified immutable text
of the same pinned source binary; guest-writable bytes are never executed.
The unchanged base-profile census and separate semantic-extension list keep
resource coverage reporting truthful. This neither introduces JIT nor obtains
any hardware or account secret.

Candidate validation now happens before the auxiliary source/frame gate.
Repeated visits to an interior native landmark are not distinct models: private
execution-local fingerprints suppress duplicate checks. The bound is 256
actual, distinct original-SINF checks, plus the existing native and instruction
budgets. Only naturally produced material is examined. No SIDB modification,
key enumeration, guessed key or externally supplied account material is added.

A model accepted by the original-SINF resolver is copied privately and retained
while actual native execution continues to seek a coherent source/frame proof.
Completion still requires both. The returned pack is exactly the accepted copy,
not a later scratch-buffer reconstruction. Failure clears the pack, checkpoint,
fingerprints and guest state. Safe diagnostics distinguish a SINF match from
runtime completion. One fresh generic-only replay also covers native-block and
instruction-budget rejections; a rejected replay's check count is reported
independently. No network authorization is repeated by this mechanism.

## Verification and remaining limitation

The differential suite now compares 3,610 complete programs with Unicorn ARM64,
including 256 new CAS cases: every width, ordering and comparison outcome,
SP/zero/aliased registers, all registers/flags/vectors and complete surrounding
memory. The production native-region runner passes these with new executable
memory mappings prohibited by seccomp. Tests also execute the exact reported
CASAL through the shipped real profile and verify coverage of all 61 CAS sites,
wrong-source rejection, alignment/memory failures, repeated candidate visits,
distinct-check bounds, accepted-copy ownership, cleanup and replay isolation.
These public/synthetic tests establish implementation contracts, not a real
FairPlay enrollment or IPA decryption.

The owner's encrypted IPA and private capture remain on the iPhone. There is no
local test input that proves the native producer derives the correct account
credential, or that the fixed source/frame output locations form a valid mobile
checkpoint. Those are still open until the original SINF and complete Mach-O
page/integrity checks pass on that device. Build 23052 must not be described as
a working IPA decryptor on the basis of compilation, coherent synthetic state,
or a successful CAS regression alone.
