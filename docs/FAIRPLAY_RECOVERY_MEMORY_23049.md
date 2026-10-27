# Authorized state and recovery memory contracts

The owner's 23048 replay returned HTTP 200 without Apple failure, imported both
real authorization bags, persisted account state, and returned
needs_account_diversity=false. Modern recovery found SIDB and SIDD, then stopped
in selection with raw_memory_unavailable. SINF validation was not attempted.
This proves original-producer authorization, not enrollment compatibility.

The shipped semantic profile identifies text+0x5d1a78 as an indirect call through
x8. A later failure in its external operation can retain that last generic PC.
Thus translating a replacement instruction at this PC is not justified.

Review of the entire external allocation/file path found an independent memory
contract violation in pinned raw_source_step.go: allocate maps max(1 MiB,
page-rounded request), but advances nextAllocation by only 1 MiB. Requests larger
than 1 MiB overlap subsequent allocations. mapZero preserves existing pages, so
a subsequent successful write can corrupt a still-live earlier allocation.
Executable mappings and file buffers can exceed this size, across applications.
No private device allocation sizes or failing pointer were present in the report;
this defect is not yet established as the cause of the owner's specific fault.

The overlay replaces that allocator as a unit: advance by the full mapped extent,
check alignment/overflow and refuse existing-page collisions before mutation.
The minimum extent for small requests is retained to avoid unrelated reference
layout changes. Only explicit supported allocation calls create memory. Unknown
pointers still fail; no license bytes, missing state files or pages are invented.
Free/munmap retain the reference lifetime behavior; allocator reuse is not claimed.
The file ABI's bounded read/mmap ranges and longest-suffix lookup remain unchanged.

An ABI transaction test allocates a live 3 MiB buffer, then opens, stats, maps,
reads and closes three distinct files larger than 2 MiB, and verifies every byte
of all live buffers after later operations. This sequence exposes the old overlap
without account keys. Invalid extent/collision tests require atomic rejection.
The normal complete portable and native/IPA CI pipelines are retained; these
checks are not evidence of real decrypted plaintext.

A failing memory operation may also report its last canonical external import
slot. The existing allowlisted slot formatter is reused; no private pointers,
registers, state bytes, file names or underlying exception text are exported.
This distinguishes an external ABI fault from a missing translated instruction.

The owner's authorized state remains on the phone and must not be discarded.
No fresh purchase, download, account authorization or arbitrary SIDB mutation is
needed to replay recovery. Real-input SINF validation and legacy/modern enrollment
compatibility remain required before claiming the first decrypted IPA.
