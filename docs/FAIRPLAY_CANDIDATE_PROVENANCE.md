# Candidate origin audit

This optional host observer extends the pinned native reference replay. It is
not linked into Crumb, does not change recovery budgets or guest memory, and
cannot supply or validate a credential. The public recovery APIs are unchanged.

Run with an existing private reference trace:

```
python3 scripts/replay-fairplay-native-reference.py \
  --trace /private/native-reference.jsonl --audit-candidate-origins \
  --audit-sinf-header --compare-mobile-sinf /private/original.sinf
```

The trace and raw SINF must be regular private files accepted by the existing
replayer. Their association must already have been verified against the same
capture and original IPA. No private fixture is uploaded to CI.

The observer follows executed register reads/writes and memory reads/writes.
Origins include actual file read/mmap contracts, service byte changes, initial
context, instruction constants and image segments. Mach-O **initial protection
bits**, rather than an image address interval, determine the initial image
classification. In particular `__DATA_CONST` with initial protection 3 is
conservatively classified as mutable; its name alone is not proof of runtime
read-only protection. Subsequent native writes replace the initial attribution.

Lookup indices retain their origin: a byte loaded from an immutable white-box
table may still depend on SIDB-derived index registers. Attribution is a
conservative union, including address bases and condition registers reported
by Capstone. It is not a symbolic derivation or an exact backward slice.

At every `text+0x234cdc` visit, the observer checks the actual 16-byte native
frame input. Equality is tested only in private process memory. Reports use
ordinal candidate numbers, origin categories, public native instruction
offsets and consecutive-run counts. Distinct visits are counted even when the
material is unchanged. No guest addresses, paths, bytes or secret hashes are
serialized by this observer.

Important limits remain explicit in every report:

- Branch/control dependencies are not reconstructed. Unchanged candidate
  bytes do not prove that an input or earlier rejection had no effect.
- Register aliases and operand origins are conservative, not exact bit-level
  expressions. Initial pointer provenance can remain in an output's union.
- Unknown services' same-value memory writes cannot be recovered from changed
  pages. Known successful file reads and mappings are tagged independently of
  whether returned bytes differ from previous contents.
- PAC and service semantics retain the independent replay's existing limits.
- A complete causal chain and SINF validation are never inferred from origins.

Per-byte origin storage uses bounded sparse pages (at most 64 MiB of guest
address coverage). Overflow aborts the audit; it does not truncate origins and
continue with an apparent success. The optional Python observer gets a longer
host watchdog because it executes per instruction. The production iOS
execution budget is unchanged.

Regression tests execute real ARM64 loads, transformations, stores and indexed
table lookups using Unicorn. They also cover Mach-O permissions, missing writer
information, repeated candidates, register aliases and same-value file reads.
These tests establish observer behavior, not real license recovery.
