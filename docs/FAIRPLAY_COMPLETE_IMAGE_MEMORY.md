# Complete immutable image memory and engine contracts

The owner's 23049 report repeats 23048's selection failure with exactly
13,692,374 generic instructions and 292,521 shared block executions. The previous
allocation/mmap fixes did not resolve this particular replay. Authorization and
both genuine license-bag imports completed, but SINF validation never started.
That remains the last real device evidence; no decrypted IPA is established.

## Image coverage, independently measured

Actions run [40493389169](https://github.com/boierito/Crumb/actions/runs/40493389169)
compared every public model page against the exact supported system executable.
It verified full executable SHA-256, immutable segment permissions, and no
chained relocations in `__TEXT`. This is separate from the semantic profile hash.

| Region | Image pages | Captured pages | Restoration |
| --- | ---: | ---: | --- |
| `__TEXT` | 4,080 | 227 | Copy 3,853 missing pages from the verified image |
| `__DATA_CONST` | 128 | 128 | Preserve captured bytes |
| `__DATA` | 24 including zero-fill | 24 | Preserve captured bytes |
| Other captured regions | — | 93 | Preserve captured bytes |

The interpreter previously covered almost every instruction while having only
a partial immutable image to read as data. Restoring the entire exact immutable
segment addresses that coverage gap for all execution paths, rather than
mapping one faulting address. It does not establish that the owner's failed
access was in this segment: the previous diagnostic lacked the precise fault PC
and address classification. Missing mutable state is still a possible blocker.

`fairplayrawimage.go` verifies the compressed and decoded resources, enforces
bounded page geometry, and installs private copies only where the model has no
page. In particular, all six captured `__TEXT` pages that differ from file bytes
remain unchanged. Cleanup cannot corrupt a resource needed by a later recovery.
No private blob, key or arbitrary zero-filled range is introduced.

## Execution verification and fault provenance

The integrated comparison runs the actual 4,265 AOT dispatch entries and the
shipped generic profile from identical synthetic register/memory states, using
four deterministic seeds. Comparisons include registers, SP, NZCV, vectors,
returned PC and memory. Each function's actual instruction span bounds the
generic execution; a returned branch target may be an interior loop head and
must not be used as a stop boundary. Synthetic missing pages belong only to the
test fixture and never to the production recovery. Rejected fixtures are
reported separately and are not evidence of equivalent executions.

The pre-integration audit passed 17,028 executions, rejected 32 fixture cases,
and found zero divergences. It cannot prove equivalence for all possible states,
nor validate the producer's enrollment or any private license.

Generated AOT instruction stop predicates now also track the current public
instruction address without changing the stop decision or architectural state.
The generic interpreter tracks the same field. Execution failures can report
`fault_instruction`, `fault_context`, and (only for an actual external context)
`fault_external_slot`. Older last-generic/last-import fields remain historical.
No guest pointer, registers, file paths or account material are exported.
The relevant public block regression forces a failure at `text+0xf0a3c` after
a historical malloc and requires attribution to AOT rather than that malloc.

## Acceptance boundary

The full immutable image/snapshot test checks every inserted page against the
exact resource and every original page against its pre-install digest. It then
destroys one recovery's pages, creates another recovery, and confirms an unknown
heap pointer still fails. Complete portable, native, independent ARM64 and iOS
build pipelines remain required.

These changes supply build-global memory, not a new credential producer. The
StoreAgent 10.9 to fairplayd 26A428 enrollment relationship is still unproven.
The SINF, source/frame, Mach-O/page and output gates are unchanged. Only a real
local replay that validates the original SINF and every encrypted executable can
establish that the owner has obtained the first decrypted IPA.
