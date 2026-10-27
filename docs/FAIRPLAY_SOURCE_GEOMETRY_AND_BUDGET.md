# Device diagnostic 777: source geometry and execution budget

Build 23052 checked one native candidate against the original mobile SINF; the
resolver rejected it. The primary execution then reached the 50,000,000 generic
instruction budget in the SINF-labelled phase, while its fresh generic-only
replay independently reached the same budget in the source-labelled phase.
Neither is evidence that all possible native candidates were exhausted, or that
the account state is incompatible. The two engines reached different phases,
so their budget exits are not comparable completed executions. The original
SINF, source/frame and plaintext remain unverified.

The last primary instruction was text+0x100628. The surrounding public native
code is a byte-processing loop with a decrementing 64-bit counter and an
indirect branch. A differential test executes the *actual pinned native code*
through both the shipped generic interpreter and Unicorn, with synthetic
memory for lengths 1, 2, 15, 16, 17, 32, 440, 448, 888 and 4096. Both engines
terminate and agree on all registers, flags, vectors and complete mapped pages.
The private run's counter is unavailable, so this test does not establish why
that run reached its budget. Raising the budget alone would not establish the
validity of a recovered license.

The pinned state transform accepts source inputs of 192 and 224 bytes, but the
previous checkpoint acceptance only tried 224. Both widths are now evaluated
against the independent state/frame equation, preferring 224 if it matches and
otherwise retaining 192 only if it matches. This uses the existing native
source address and does not infer source size from the legacy SIDB header or
modify account data. The same chosen width is used for post-SINF comparison.
The original SINF resolver and all later IPA checks remain mandatory.

The 192-byte path is supported by the portable evaluator and tested with
synthetic input; the diagnostic does not establish that this account produces a
192-byte source, or that this change solves the instruction-budget exit. The
original private capture and IPA remain only on the owner's device. No mobile
license or complete IPA decryption has yet been validated.
