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
if 0x100a4be38 in blocks:
    raise SystemExit('upstream now covers next supplement; review before merging')
blocks[0x100a4be38] = 'rawSourceBootstrapBlock100a4be38'
if 0x1008c292c in blocks:
    raise SystemExit('upstream covers exit supplement; review before merging')
blocks[0x1008c292c] = 'rawSourceBootstrapBlock1008c292c'
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
    if name == 'oracle.go':
        fault = 'return nil, fmt.Errorf("oracle memory address %#x is outside the model", address)'
        if text.count(fault) != 1:
            raise SystemExit('oracle unmapped memory contract changed')
        text = text.replace(fault, 'return nil, &sourceMemoryError{address: address, write: write}')
    if name == 'raw_source_step.go':
        # Replace the pinned bump allocator as a unit. Its old fixed 1 MiB
        # advance aliases every larger allocation with subsequent allocations.
        allocation = re.compile(r'func \(r \*rawSourceRuntime\) allocate\(s \*state, size uint64\) \(uint64, error\) \{.*?\n\}', re.S)
        text, count = allocation.subn('func (r *rawSourceRuntime) allocate(s *state, size uint64) (uint64, error) {\n return allocateSourceMemory(r, s, size)\n}', text)
        if count != 1:
            raise SystemExit('raw allocator contract changed')
        mmap_extent = 'accessibleLength := (length + 0x3FFF) &^ 0x3FFF\n\t\taddress, err := r.allocate(s, length)'
        if text.count(mmap_extent) != 1:
            raise SystemExit('raw mmap extent contract changed')
        text = text.replace(mmap_extent, 'if length == 0 || length > 0x10000000 { return pc, fmt.Errorf("invalid bootstrap mapping size") }\n\t\t' + mmap_extent.replace('r.allocate(s, length)', 'r.allocate(s, accessibleLength)'))
        # Instrument the existing lookup only; preserve longest-suffix matching
        # and all guest-visible return values. Never retain requested paths.
        field = 'type rawSourceRuntime struct {'
        lookup = 'return bestSuffix, bestData, bestSuffix != ""'
        if text.count(field) != 1 or text.count(lookup) != 1:
            raise SystemExit('raw runtime lookup contract changed')
        text = text.replace(field, field + '\n lookupCounts map[string]LookupCount\n stateFileLookupCounts map[string]LookupCount\n stateLookupsByPhase map[string]LookupCount\n phase string\n genericOnly bool\n observe func(uint64)\n boundary func() (bool,error)\n candidateChecks int\n candidateSINFVerified bool\n rwlocks map[uint64]int\n qosClass uint32\n qosRelative int32\n qosSet bool\n currentInstruction uint64\n executionContext string\n lastExternalSlot string\n unmodeledExternalSlot string\n lastKnownBlock uint64\n additionalBlocks int\n untranslatedTarget string\n targetClass string\n genericBlocks int\n genericInstructions uint64\n genericLastPC uint64\n checkpointCompared bool\n sourceChangedAfterSINF bool\n frameChangedAfterSINF bool')
        text = text.replace(' genericOnly bool\n', ' genericOnly bool\n byteLoop genericLoopObservation\n', 1)
        # Replace legacy no-op pthread operations as a group. The new model
        # accepts only lifecycle states supported by native reference tests.
        legacy = "\tcase 0xF10000000000E000, // pthread_rwlock_unlock\n\t\t0xF100000000040000, // pthread_rwlock_rdlock\n\t\t0xF10000000005C000, // pthread_rwlock_init\n\t\t0xF100000000065000: // pthread_rwlock_wrlock\n\t\ts.x[0] = 0"
        replacement = "\tcase 0xF10000000004C000, 0xF10000000000E000, 0xF100000000040000, 0xF10000000005C000, 0xF100000000065000:\n\t\tstatus, err := r.handleSourceRWLock(s, pc)\n\t\tif err != nil { r.unmodeledExternalSlot = sourceExternalImportSlot(pc); return pc, err }\n\t\ts.x[0] = status"
        if text.count(legacy) != 1: raise SystemExit('pthread ABI contract changed')
        text = text.replace(legacy, replacement)
        qos_anchor = "\tcase 0xF10000000004C000, 0xF10000000000E000, 0xF100000000040000, 0xF10000000005C000, 0xF100000000065000:"
        if text.count(qos_anchor) != 1: raise SystemExit('pthread QoS ABI anchor changed')
        text = text.replace(qos_anchor, "\tcase 0xF100000000057000: // pthread_set_qos_class_self_np\n\t\ts.x[0] = r.handleSourceQoS(s)\n" + qos_anchor)
        entry = 'returnAddress := s.x[30]'
        if text.count(entry) != 1: raise SystemExit('external entry contract changed')
        text = text.replace(entry, entry + '\n r.lastExternalSlot = sourceExternalImportSlot(pc)\n r.executionContext = "external"')
        text = text.replace(lookup, 'r.recordLookup(path, bestSuffix != "")\n\t' + lookup)
        # Record the rejected import's *static model slot*, never a pointer
        # or private guest register. No unsupported function is stubbed.
        rejected = 'return pc, fmt.Errorf("unsupported bootstrap external %#x", pc)'
        if text.count(rejected) != 1:
            raise SystemExit('pinned external dispatch changed; review required')
        text = text.replace(rejected, 'r.unmodeledExternalSlot = sourceExternalImportSlot(pc)\n\t\t' + rejected)

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

(target/'source_qos.go').write_bytes((api.parent/'fairplayrawqos.go').read_bytes())
(target/'source_qos_test.go').write_bytes((api.parent/'fairplayrawqos_test.go').read_bytes())

(target/'source_rwlock.go').write_bytes((api.parent/'fairplayrawrwlock.go').read_bytes())
(target/'source_rwlock_test.go').write_bytes((api.parent/'fairplayrawrwlock_test.go').read_bytes())

(target/'candidate_validation.go').write_bytes((api.parent/'fairplayrawcandidates.go').read_bytes())
(target/'candidate_validation_test.go').write_bytes((api.parent/'fairplayrawcandidates_test.go').read_bytes())

(target / "api_test.go").write_text((api.parent / "fairplayrawapi_test.go").read_text())
(target / "enrollment_checkpoint_test.go").write_text((api.parent / "fairplayrawcheckpoint_test.go").read_text())
(target / 'block_dispatch.go').write_text((api.parent / 'fairplayrawblocks.go').read_text())
(target / 'bootstrap_supplement.go').write_text((api.parent / 'fairplayrawsupplement.go').read_text())
(target / 'block_dispatch_test.go').write_text((api.parent / 'fairplayrawblocks_test.go').read_text())
(target / 'bootstrap_supplement_test.go').write_text((api.parent / 'fairplayrawsupplement_test.go').read_text())
(target / 'external_test.go').write_text((api.parent / 'fairplayrawexternal_test.go').read_text())

for name, origin in [('LICENSE-IPATOOL', source.parent.parent / 'LICENSE'),
                     ('LICENSE-DARKPHOENIX.txt', source / 'LICENSE-DARKPHOENIX.txt'),
                     ('PROFILE-README.md', source / 'profiles/README.md')]:
    (target / name).write_bytes(origin.read_bytes())

(target / 'bootstrap_next.go').write_text((api.parent / 'fairplayrawsupplement_next.go').read_text())
(target / 'bootstrap_next_test.go').write_text((api.parent / 'fairplayrawsupplement_next_test.go').read_text())
(target / 'testdata').mkdir(exist_ok=True)
(target / 'testdata/bootstrap-next.json').write_bytes((api.parent.parent.parent / 'scripts/fixtures/fairplay-bootstrap-next.json').read_bytes())

(target / 'bootstrap_exit.go').write_text((api.parent / 'fairplayrawsupplement_exit.go').read_text())
(target / 'bootstrap_exit_test.go').write_text((api.parent / 'fairplayrawsupplement_exit_test.go').read_text())

# Generic instruction profile: data covering every executable-section slot.
import subprocess, json
root = api.parent.parent.parent
subprocess.run([sys.executable, str(root/'scripts/prepare-fairplay-generic.py'), 'runtime', str(target)], check=True)
for source_name, target_name in [('fairplaygeneric_atomic.go','generic_atomic.go'), ('fairplaygeneric_atomic_test.go','generic_atomic_test.go'), ('fairplaygeneric.go','generic.go'), ('fairplaygeneric_test.go','generic_test.go'), ('fairplaygeneric_real_test.go','generic_real_test.go')]:
    (target/target_name).write_bytes((api.parent/source_name).read_bytes())
from fairplay_generic_resource import read_profile
compressed,pin=read_profile(root)
(target/'profiles/fairplay-generic.ir.gz').write_bytes(compressed)
(target/'generic_embed.go').write_text('package rawsource\nimport (_ "embed"; "crypto/sha256"; "encoding/hex")\n//go:embed profiles/fairplay-generic.ir.gz\nvar genericProfileCompressed []byte\nfunc init() {\n sum:=sha256.Sum256(genericProfileCompressed)\n if hex.EncodeToString(sum[:]) != "GENERIC_RESOURCE_SHA256" { return }\n image,err:=loadGenericInstructionProfile(genericProfileCompressed,rawSourceBinaryHash,genericKernelABI)\n if err==nil { text,e:=immutableImageText(); if e==nil && image.attachNativeText(rawSourceBinaryHash,imageTextAddress,text)==nil { genericImage=image } }\n}\n'.replace('GENERIC_RESOURCE_SHA256',pin['compressed_sha256']))

(target/'memory_contract.go').write_bytes((api.parent/'fairplayrawmemory.go').read_bytes())
(target/'memory_contract_test.go').write_bytes((api.parent/'fairplayrawmemory_test.go').read_bytes())

# Public, immutable image data fills only missing __TEXT pages. No writable
# state is reconstructed; account/heap/stack gaps retain fail-closed behavior.
from fairplay_image_resource import read_image_text
image_text, image_pin = read_image_text(root)
(target/'IMMUTABLE-IMAGE-PROVENANCE.json').write_bytes((root/'third_party/fairplay-immutable-image/IMAGE.json').read_bytes())
(target/'IMMUTABLE-IMAGE-README.md').write_bytes((root/'third_party/fairplay-immutable-image/README.md').read_bytes())
(target/'profiles/fairplay-immutable-text.bin.gz').write_bytes(image_text)
(target/'image_text.go').write_bytes((api.parent/'fairplayrawimage.go').read_bytes())
(target/'image_text_test.go').write_bytes((api.parent/'fairplayrawimage_test.go').read_bytes())
def go_digest(value):
    return '[32]byte{' + ','.join('0x'+value[i:i+2] for i in range(0,64,2)) + '}'
text_pin = image_pin['immutable_text']
(target/'image_text_embed.go').write_text('package rawsource\nimport _ "embed"\n'
    '//go:embed profiles/fairplay-immutable-text.bin.gz\nvar imageTextCompressed []byte\n'
    f'const imageTextAddress uint64 = {text_pin["address"]}\n'
    f'const imageTextSize uint64 = {text_pin["bytes"]}\n'
    f'var imageTextCompressedSHA = {go_digest(text_pin["compressed_sha256"])}\n'
    f'var imageTextRawSHA = {go_digest(text_pin["sha256"])}\n')

(target/'instruction_trace.go').write_bytes((api.parent/'fairplayrawtrace.go').read_bytes())
(target/'instruction_trace_test.go').write_bytes((api.parent/'fairplayrawtrace_test.go').read_bytes())
for path in sorted(target.glob('*.go')):
    if path.name.startswith('generated_raw_source_') or path.name.startswith('bootstrap_'):
        text = path.read_text()
        text = re.sub(r'if runtime != nil && runtime.stopPC == (0x[0-9a-f]+)',
                      r'if sourceInstructionStop(runtime, \1)', text)
        path.write_text(text)

# Compare both engines over the actual instruction span of each selected AOT
# function, including trace fragments that end before a native branch. Taking
# the *returned* PC as a bound would incorrectly stop at an interior loop head.
function_ends = {}
for path in sorted(target.glob('*.go')):
    text = path.read_text()
    functions = list(re.finditer(r'^func (\w+)\(', text, re.M))
    for index, function in enumerate(functions):
        if not function[1].startswith('rawSource') or 'Block' not in function[1]:
            continue
        body = text[function.end():functions[index+1].start() if index+1 < len(functions) else len(text)]
        addresses = re.findall(r'// (0x[0-9a-f]+):|sourceInstructionStop\(runtime, (0x[0-9a-f]+)\)', body)
        if addresses:
            function_ends[function[1]] = max(int(a or b, 16) for a, b in addresses)+4
if any(function not in function_ends for function in blocks.values()):
    raise SystemExit('AOT comparison instruction geometry missing')
(target/'engine_geometry_test.go').write_text('package rawsource\nvar auditBlockEnd=map[uint64]uint64{\n' +
    ''.join(f'0x{pc:x}:0x{function_ends[function]:x},\n' for pc, function in sorted(blocks.items())) + '}\n')
(target/'engine_comparison_test.go').write_bytes((api.parent/'fairplayrawengines_test.go').read_bytes())

for name in ('flow', 'replay', 'fault'):
    (target/f'native_{name}.go').write_bytes((api.parent/f'fairplayraw{name}.go').read_bytes())
    (target/f'native_{name}_test.go').write_bytes((api.parent/f'fairplayraw{name}_test.go').read_bytes())
