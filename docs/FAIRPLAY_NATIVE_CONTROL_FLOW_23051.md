# Continuous native enrollment recovery

The owner's 23050 run completed authorization and progressed past the old
selection failure into the SINF-labelled path. It stopped at `text+0x8a680`,
`ldr w14,[x21,#0xb40]`. It had not validated a SINF or decrypted an IPA. The
report does not identify the private pointer or the predecessor of the forced
SINF entry, so it cannot establish why this particular context was unavailable.

## Two general control-flow defects

The 23050 production source path still called `runRawSourceKernelLoop("outer")`.
Its generated outer kernel contains 35,243 recorded instructions and 36
noncontiguous recorded edges. For example, `text+0x5b7a8c` is `braa x24,x17`,
but that kernel continues at `text+0x5f3bdc` without validating x24. Other traced
kernels similarly replace calculated branches with recorded paths or host loop
conditions. The previous 4,265-block differential test did not cover this outer
trace shortcut. Matching its arithmetic does not establish its control flow for
another account/state.

The shared CFG dispatcher also returns a block's local `stopped` flag to the
active phase. The real SINF fragment at `text+0x234cb8` returns
`text+0x234cdc,true`. When reused by selection, the previous fallback accepted
that flag as selection completion even though selection's boundary was
`text+0x8a644`. The orchestrator discarded the returned PC and started the next
phase at its fixed entry. That can restart a frame body with registers from its
continuation. This is a source-verifiable defect and has an actual-block
regression; the owner's exact predecessor was not recorded.

## Implemented execution contract

`fairplayrawflow.go` follows the actual PC continuously from the pinned account
entry. It dispatches verified AOT blocks, generic instructions, and the existing
supported external ABI. A local trace stop is a continuation, not stage success.
Phase names are observations at native landmarks; they never reset the PC, SP
or registers. Production recovery no longer executes the linear source kernels.
Unknown instructions, imports and memory remain errors. The whole run is
bounded by 8 million blocks and the existing 50 million generic instructions.

The original native material landmark can occur during different frame calls.
`RecoverValidated` therefore checks only naturally produced candidates, using
the production original-SINF resolver supplied by the account/IPA caller. A
rejected candidate continues at its actual PC. It never restarts a function,
patches a pointer, changes SIDB bytes or substitutes material. At most 64 native
candidate events are allowed. The non-validating research API retains the
reference's two naturally observed frame visits, without forcing their PCs.

Source/frame snapshots are retained only after the same independent state/frame
equations used by final validation succeed. If no coherent snapshot exists, no
candidate can complete recovery. All later selector, original SINF, whole-IPA,
Mach-O and page validation remains mandatory. `sinf_validation_attempted` is
updated when the original SINF resolver is actually invoked, including attempts
inside a recovery that subsequently fails.

On an eligible memory/control-flow rejection, one generic-only execution starts
from fresh public model, immutable image, CPU and private file copies. It does
not repeat network authorization or continue a failed block. A successful raw
candidate still passes the original SINF and final IPA gates. Rejected replays
return no private material and are reported independently; a budget rejection
does not imply agreement with the first execution. Each attempt clears its
private model/files/CPU state. Resource and unsupported-operation failures are
not repeatedly retried.

Typed memory failures retain addresses only inside the package and expose a
fixed region taxonomy. No address, register, key, private hash or requested path
is added to diagnostics.

## Verification beyond the reported failure

The actual native branch and real cross-phase SINF fragment have regressions.
Control-flow tests include loops, natural boundary visits, calls through x30,
returns, partial AOT fragments, rejected candidates, unknown destinations and
finite execution bounds. The independent ARM64 differential suite now runs the
production region runner, including additional complete loops and calls, rather
than assuming how many generic chunks should execute. Its 3,354 programs also
run with new executable mappings prohibited. The previous full AOT inventory,
immutable-image, offline capture and native producer checks are retained.

The production portable page engine is also cross-checked against all 12
independent EthanArbuckle mobile vectors and Crumb's separate mobile consumer:
inner SINF key to logical key, absolute-offset page-key derivation, AES-CBC page
plaintext and expected SHA-256. Both 440/448-byte private records and page indices
through uint32 maximum are covered. This proves agreement of these downstream
primitives on public fixtures, not acquisition of a real account credential or
decryption of the owner's application.

No real post-change device SINF/page result is available yet. These changes
correct execution assumptions; they do not prove StoreAgent 10.9 / fairplayd
26A428 enrollment compatibility or fulfill the first real decrypted-IPA goal.
