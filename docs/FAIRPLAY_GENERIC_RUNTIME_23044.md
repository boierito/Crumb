# FairPlay 23044: generic execution instead of per-address repairs

## Problem and implemented change

The previous runtime executed only a catalog of 4,265 translated block entries.
A valid indirect jump to another instruction therefore stopped with
`raw_unknown_control_flow`. Device builds 23042 and 23043 demonstrated this at
successive bootstrap continuations. Adding an individual handler did not fix
that design limitation.

The dispatcher now retains the validated AOT blocks as a fast path and falls
back to one generic semantic interpreter for any aligned instruction address
in the exact executable sections. An entry does not need to be a discovered
basic-block start. There is no new-address allowlist, per-address function generation, JIT,
or native execution of guest code on the phone.

The build-time compiler enumerates **every four-byte position** in every executable section (including call stubs),
converts supported instructions into fixed-size 16-byte semantic records, and
explicitly marks unsupported positions. Non-executable gaps are distinct and
are never accepted as code. Indirect and conditional targets are
computed from the actual guest registers, flags and memory. The same
interpreter handles subsequent new destinations automatically.

## Current coverage and limits

The exact source daemon is pinned by SHA256
`1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784`.
The checked profile covers `__text`, `__auth_stubs` and `__objc_stubs`.
Its flat span has six non-executable gap slots, which are rejected.
It has 3,444,036 executable instruction positions; 3,442,039 are supported
(99.942016%). **1,997 positions remain unsupported**, classified and exemplified
in `third_party/unfairplay-generic/CENSUS.json`. Some are architectural traps
or platform-specific operations. They are not silently skipped.

A new valid PC inside the profile no longer fails merely because no function
exists for that PC. An unsupported instruction produces
`generic_unsupported_instruction`; an out-of-profile target produces
`generic_pc_outside_profile`. Existing memory and host-service failures remain
separate. Nothing invents missing memory, state, a frame or a branch target.

This is **not** a complete implementation of every ARM64 instruction or future
Apple build. Future instruction families can require semantic support; future
binary versions require a matching, revalidated profile. This does remove the
need to repair every new code address individually.

The profile occupies 13,332,763 compressed bytes and 55,104,800 decompressed
bytes. It contains public semantic operations and code locations, no account
material and no Mach-O executable. The phone does not need a Mac capture,
jailbreak or JIT. The existing portable convention treats authenticated
branches as branches to their supplied targets; native PAC is not validated.

## Validation and reproducibility

- 3,290 independently executed synthetic ARM64 programs match the generic
  interpreter: conditional aliases, every condition/NZCV combination, W/X
  widths, zero registers, BLR through LR, signed memory indices, stack effects
  and vector logic. The tests compile profiles from synthetic Mach-O sections;
  they do not call per-PC AOT functions.
- The same programs pass with new mmap/mprotect PROT_EXEC mappings denied by
  Linux seccomp. The interpreter itself uses only Go data structures.
- Go race tests exercise unknown/interior landings, cross-section branches, non-executable gaps, conditionals, stopPC,
  malformed source/ABI/checksum, unsupported and out-of-profile slots, a
  bounded sequential chunk and a cumulative instruction budget.
- The **shipped real profile** passes integrity checks and directly executes
  the current `0x100a06948` device regression with synthetic signed table
  indices. That PC has no AOT handler. Registers, NZCV, SP and the complete
  stack page are checked.
- Limits: 4,096 sequential instructions per dispatch chunk, the existing phase
  loop limits, and 50 million cumulative generic instructions per recovery.
- Code profile generation:
  https://github.com/boierito/Crumb/actions/runs/40363976865
  Compiler revision: `8888748f37bb49236a1398aa65da549ddc5446d9`.
  Original MIT interpreter/compiler sources are pinned at UnfairPlay
  `c54e751775cbc2a54510c0fcafb0aef6b1619bc9`; source hashes and notices are retained.

The new profile compiler also covers conditional select/compare and vector
logic families generically, rather than maintaining address exceptions. The
conditional tests include AL/NV alias behavior and the BLR/LR target alias.

## Build and future changes

`fairplay-generic-profile.yml` independently regenerates the whole-section
profile and census on an exact development host. Native builds verify every
resource part, the combined resource hash, source hash, schema, section bounds,
interpreter ABI and payload checksum before embedding it. Native cache keys
include interpreter sources and profile parts. No account capture enters CI.

The generic release gate checks resource integrity and actual dispatcher
integration, not the presence of a function named after the current PC.
`Run self-tests` reports execution-profile readiness separately from enrollment
or decryption readiness. Recovery diagnostics report the selected backend,
generic block count, instruction count and the last validated public instruction
offset, without registers or memory bytes.

For a new source binary: regenerate the matching source/model profile, compile
all instruction slots, review the full unsupported-operation census, run the
ISA differential and negative tests, then build and test on device. Addresses
alone are not reasons to add application code.

## What this does not establish

Synthetic CPU correctness and successful compilation do not prove that the
legacy StoreAgent SIDB satisfies the modern FairPlay enrollment contract. The
source/frame/model checks and IPA output gates are unchanged. Real iPhone
recovery, a validated model and decrypted-app behavior still need testing.
