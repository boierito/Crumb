#!/usr/bin/env python3
"""Isolate the pinned pure-Go modern enrollment runtime (no guest execution)."""
import hashlib
import pathlib
import re
import sys

source, target = map(pathlib.Path, sys.argv[1:])
target.mkdir(parents=True, exist_ok=True)
names = ['oracle.go', 'main.go', 'dfa_aes.go', 'dfa_attack.go', 'dfa_finish.go',
         'dfa_linear.go', 'generated_26b118.go', 'input_size_default.go',
         'enrollment_frame_default.go', 'embedded_model_raw_source.go',
         'embedded_model_raw_source_test.go', 'raw_source_step.go',
         'raw_source_runtime_test.go', 'raw_source_vector_test.go']
names += [p.name for p in sorted(source.glob('generated_raw_source_*.go'))]
blocks = {}
instructions = {}
# All definitions come from the same hash-pinned generated binary profile.
# Keep a deterministic representative when phases share a block entry.
for name in sorted(names):
    if name.startswith('generated_raw_source_'):
        original = (source / name).read_text()
        definitions = re.findall(r'func (rawSource\w+Block([0-9a-f]+))\(', original)
        if definitions:
            for address, instruction in re.findall(r'// 0x([0-9a-f]+): ([^\n]+)', original):
                instruction = ' '.join(instruction.split())
                if address in instructions and instructions[address] != instruction:
                    raise SystemExit('conflicting translated instruction at a shared code offset')
                instructions[address] = instruction
        for function, address in definitions:
            blocks.setdefault(int(address, 16), function)
if len(blocks) != 4262:
    raise SystemExit('pinned CFG block inventory changed')
# Verified 26A434 public system daemon has the same pinned 26A428 hash.
# This exact omitted landing block is compiled ahead-of-time: no JIT and
# no raw executable pages loaded into the phone at runtime.
supplement_entry = 0x1008fff2c
if supplement_entry in blocks:
    raise SystemExit('pinned original CFG now covers the supplement; review upstream before merging')
blocks[supplement_entry] = 'rawSourceBootstrapBlock1008fff2c'
registry = ['// Generated from the pinned translated blocks; no native execution.',
            'package rawsource', 'var knownSourceBlocks = map[uint64]sourceBlock{']
for address, function in sorted(blocks.items()):
    registry.append(f'0x{address:x}: {{run: {function}, label: "text+0x{address - 0x100000000:x}"}},')
registry.append('}')
(target / 'known_blocks.go').write_text('\n'.join(registry) + '\n')
cfg_ends = {
    'generated_raw_source_account_entry_cfg.go': 0x10057ae74,
    'generated_raw_source_bootstrap_cfg.go': 0x1005b7438,
    'generated_raw_source_selection_cfg.go': 0x10008a644,
    'generated_raw_source_sinf_entry_cfg.go': 0x100234cdc,
}
for name in names:
    text = (source / name).read_text()
    text = re.sub(r'^//go:build[^\n]*\n', '', text)
    text = text.replace('package main', 'package rawsource', 1)
    if name == 'raw_source_step.go':
        # Instrument the existing lookup only; preserve longest-suffix matching
        # and all guest-visible return values. Never retain requested paths.
        field = 'type rawSourceRuntime struct {'
        lookup = 'return bestSuffix, bestData, bestSuffix != ""'
        if text.count(field) != 1 or text.count(lookup) != 1:
            raise SystemExit('raw runtime lookup contract changed')
        text = text.replace(field, field + '\n lookupCounts map[string]LookupCount\n lastKnownBlock uint64\n additionalBlocks int\n untranslatedTarget string\n targetClass string')
        text = text.replace(lookup, 'r.recordLookup(path, bestSuffix != "")\n\t' + lookup)
    if name in cfg_ends:
        # Preserve the original cases. Only their unknown-target branch may
        # resume an existing translated block from the shared inventory.
        pattern = r'next, stopped, err := (rawSource\w+Block([0-9a-f]+))\(s, runtime\)'
        text, count = re.subn(pattern, lambda m: f'next, stopped, err := callSourceBlock(s, runtime, 0x{m[2]}, {m[1]})', text)
        if count == 0:
            raise SystemExit('CFG dispatch call contract changed')
        old = '\t\tdefault:\n\t\t\treturn pc, blockCount, fmt.Errorf("bootstrap CFG reached unknown block %#x", pc)'
        if text.count(old) != 1:
            raise SystemExit('CFG unknown branch contract changed')
        fallback = f'''\t\tdefault:
            next, stopped, covered, err := resumeSourceBlock(s, runtime, pc, 0x{cfg_ends[name]:x})
            if err != nil {{ return pc, blockCount, err }}
            if covered {{
                if stopped {{ return next, blockCount, nil }}
                pc = next
                continue
            }}
            return pc, blockCount, fmt.Errorf("bootstrap CFG reached unknown block %#x", pc)'''
        text = text.replace(old, fallback)
    (target / name).write_text(text)
profile = 'macos26a428-account-entry-public.model.gz'
(target / 'profiles').mkdir(exist_ok=True)
resource = (source / 'profiles' / profile).read_bytes()
(target / 'profiles' / profile).write_bytes(resource)
expected = hashlib.sha256(resource).hexdigest()
api = pathlib.Path(__file__).resolve().parent.parent / 'MapleSyrup/NativeSAP/fairplayrawapi.go'
text = api.read_text().replace('RAW_RESOURCE_SHA256', expected)
(target / 'api.go').write_text(text)

(target / "api_test.go").write_text((api.parent / "fairplayrawapi_test.go").read_text())
(target / 'block_dispatch.go').write_text((api.parent / 'fairplayrawblocks.go').read_text())
(target / 'bootstrap_supplement.go').write_text((api.parent / 'fairplayrawsupplement.go').read_text())
(target / 'block_dispatch_test.go').write_text((api.parent / 'fairplayrawblocks_test.go').read_text())

for name, origin in [('LICENSE-IPATOOL', source.parent.parent / 'LICENSE'),
                     ('LICENSE-DARKPHOENIX.txt', source / 'LICENSE-DARKPHOENIX.txt'),
                     ('PROFILE-README.md', source / 'profiles/README.md')]:
    (target / name).write_bytes(origin.read_bytes())
