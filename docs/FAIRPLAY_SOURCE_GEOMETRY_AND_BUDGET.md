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

The pinned state transform accepts 192- and 224-byte inputs, whereas the
existing checkpoint gate reads 224 bytes. An attempted change to infer the
width from the source/frame equation failed its regression: both widths can
match the same frame. The equation alone does not identify the actual width,
so the speculative selection was withdrawn. Establishing the native source
length requires evidence from its producer or native calling convention.

The diagnostic does not establish the reason for the instruction-budget exit.
The private capture and original IPA remain on the owner's phone. No mobile
license or complete IPA decryption has been validated.
