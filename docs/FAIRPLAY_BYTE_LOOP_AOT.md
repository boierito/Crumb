# Bounded ahead-of-time execution of the native byte loop

The account-authorized device diagnostic saw 145,579 consecutive generic
visits to the public native loop at text+0x1005ac with no observed counter
discontinuity. Its largest observed counter was between 4,097 and 1,048,576;
the attempt stopped at the global 50-million generic-instruction budget before
the loop completed. This does not identify the input, prove eventual SINF
success, or establish that other parts of the recovery are finite.

Translate the complete 65-instruction block ahead of time from the exact pinned
public daemon. The function performs the same memory operations and branch
decision in Go, including per-instruction stop checks; it is not runtime code
generation. Each execution is still accounted for by the existing bounded
native-region block limit (8 million blocks). Other generic instructions keep
their 50-million limit, and a generic-only replay keeps its independent limit.
Consequently this changes the distribution of work between the two limits;
it does not prove that either remaining budget is sufficient. No license or
cryptographic validation is relaxed.

The generator artifact was created in Actions run 40594362858 from the pinned
binary SHA-256. The checked fixture ties every instruction byte to the verified
immutable text. Deterministic tests execute the exact native bytes in Unicorn
and compare both Go engines across eleven synthetic sizes including 145,600
bytes, all registers, flags, vectors, entire mapped pages, branch target and
observed iteration counts. The AOT branch uses the existing portable model's
pointer-authentication convention; this is not proof of native PAC validation.
The real account-authorized state remains on the owner's iPhone, so real SINF
validation remains unverified until the integrated path is tested there.
