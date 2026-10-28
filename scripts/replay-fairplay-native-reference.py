"""Replay one captured FairPlay invocation using independent native ARM64.

The trace stays local and must be private. Instruction execution and guest
memory are checked against Go at every recorded service boundary. Services
are replayed from the capture, not independently reimplemented; PAC operations
follow the portable runtime convention. No SINF or plaintext authenticity is
claimed by a matching replay.
"""

import argparse, base64, collections, gzip, hashlib, json, pathlib, stat, struct, sys, time
from capstone import Cs, CS_ARCH_ARM64, CS_MODE_ARM
from unicorn import (
    Uc,
    UC_ARCH_ARM64,
    UC_MODE_ARM,
    UC_HOOK_CODE,
    UC_HOOK_BLOCK,
    UC_HOOK_MEM_WRITE,
    UC_HOOK_MEM_WRITE_UNMAPPED,
    UC_ERR_FETCH_UNMAPPED,
    UcError,
    arm64_const as arm,
)

ROOT = pathlib.Path(__file__).resolve().parent.parent
from fairplay_image_resource import read_image_text
from fairplay_generic_resource import read_profile

p = argparse.ArgumentParser()
p.add_argument("--trace", required=True)
args = p.parse_args()
path = pathlib.Path(args.trace)
st = path.lstat()
assert stat.S_ISREG(st.st_mode) and st.st_mode & 0o077 == 0 and st.st_size <= 512 << 20
root = ROOT
encoded, pin = read_image_text(root)
native = gzip.decompress(encoded)
assert hashlib.sha256(native).hexdigest() == pin["immutable_text"]["sha256"]
compressed, irpin = read_profile(root)
ir = gzip.decompress(compressed)
assert hashlib.sha256(ir[:-32]).digest() == ir[-32:]
base, count = struct.unpack_from("<QI", ir, 80)
records = memoryview(ir)[96:-32]
assert len(records) == count * 16
md = Cs(CS_ARCH_ARM64, CS_MODE_ARM)
md.detail = True
patches = {}
image_base = 0x100000000
for i in range(count):
    if records[i * 16] == 6:
        continue
    pc = base + i * 4
    offset = pc - image_base
    if offset < 0 or offset + 4 > len(native):
        continue
    word = struct.unpack_from("<I", native, offset)[0]
    if not (
        word & 0xFE000000 == 0xD6000000
        or word >> 24 == 0xDA
        or word & 0xFFFFF000 == 0xD5032000
    ):
        continue
    ins = next(md.disasm(native[offset : offset + 4], pc, count=1), None)
    if ins is None:
        continue
    name = ins.mnemonic
    replacement = None
    if name.startswith(("braa", "brab", "blraa", "blrab")):
        reg = md.reg_name(ins.operands[0].reg)
        rn = {"fp": 29, "lr": 30}.get(reg)
        if rn is None:
            rn = int(reg[1:])
        replacement = (
            0xD63F0000 if name.startswith(("blraa", "blrab")) else 0xD61F0000
        ) | (rn << 5)
    elif name in ("retaa", "retab"):
        replacement = 0xD65F03C0
    elif name.startswith(("pac", "aut", "xpac")):
        replacement = 0xD503201F
    if replacement is not None:
        patches[pc] = struct.pack("<I", replacement)
print(
    json.dumps(
        {
            "stage": "reference_loaded",
            "pac_convention_patches": len(patches),
            "binary_pin_verified": True,
        }
    ),
    flush=True,
)
stream = path.open()
initial = json.loads(next(stream))
assert initial["Type"] == "initial" and initial["BinarySHA256"] == pin["binary_sha256"]
assert initial["CPU"]["PC"] == 0x100CD1C18


def validate_cpu(cpu):
    assert len(cpu["Registers"]) == 31 and all(
        type(v) is int and 0 <= v < 1 << 64 for v in cpu["Registers"]
    )
    assert len(cpu["Flags"]) == 4 and all(type(v) is bool for v in cpu["Flags"])
    assert all(type(cpu[n]) is int and 0 <= cpu[n] < 1 << 64 for n in ("PC", "SP"))


def external_slot(pc):
    return pc in (0x720000000000, 0xF10000000002500C) or (
        0xF100000000000000 <= pc < 0xF100000001000000
        and (pc - 0xF100000000000000) % 0x1000 == 0
    )


validate_cpu(initial["CPU"])

# Discover external dispatch slots without retaining private events in memory.
slots = set()
events = 0
for line in stream:
    e = json.loads(line)
    events += 1
    if e["Type"] == "external":
        validate_cpu(e["CPU"])
        validate_cpu(e["After"])
        assert external_slot(e["CPU"]["PC"])
        slots.add(e["CPU"]["PC"])
assert events > 0
stream.close()
stream = path.open()
next(stream)
u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
u.ctl_set_cpu_model(arm.UC_CPU_ARM64_MAX)
mapped = set()
dirty = set()
recent = collections.deque(maxlen=32)
patch_pages = {}
for pc, b in patches.items():
    patch_pages.setdefault(pc & ~4095, []).append((pc, b))


def map_pages(values):
    assert all(
        type(a) is int
        and 0 <= a <= (1 << 64) - 4096
        and a % 4096 == 0
        and len(b) == 4096
        for a, b in values.items()
    )
    new = sorted(set(values) - mapped)
    if new:
        start = last = new[0]
        for a in new[1:] + [None]:
            if a is not None and a == last + 4096:
                last = a
                continue
            u.mem_map(start, last - start + 4096)
            if a is not None:
                start = last = a
        mapped.update(new)
    for a, data in values.items():
        u.mem_write(a, bytes(data))
        for pc, b in patch_pages.get(a, []):
            u.mem_write(pc, b)


pages = {
    int(a): bytearray(base64.b64decode(b, validate=True))
    for a, b in initial["Pages"].items()
}
restored = 0
for i in range(count):
    if records[i * 16] == 6:
        continue
    pc = base + i * 4
    a = pc & ~4095
    if a not in pages:
        continue
    offset = pc - image_base
    if offset < 0 or offset + 4 > len(native):
        continue
    current = pages[a][pc & 4095 : (pc & 4095) + 4]
    original = native[offset : offset + 4]
    if current != original:
        pages[a][pc & 4095 : (pc & 4095) + 4] = original
        restored += 1
map_pages(pages)
del pages
del initial["Pages"]
stub_pages = {a: bytes(4096) for a in {pc & ~4095 for pc in slots} if a not in mapped}
map_pages(stub_pages)
for pc in slots:
    u.mem_write(pc, struct.pack("<I", 0xD65F03C0))
print(
    json.dumps(
        {
            "stage": "initial_state_mapped",
            "pages": len(mapped),
            "original_instruction_restorations": restored,
        }
    ),
    flush=True,
)
regs = [getattr(arm, "UC_ARM64_REG_X" + str(i)) for i in range(29)] + [
    arm.UC_ARM64_REG_X29,
    arm.UC_ARM64_REG_X30,
]


def set_cpu(cpu):
    for r, v in zip(regs, cpu["Registers"]):
        u.reg_write(r, v)
    u.reg_write(arm.UC_ARM64_REG_SP, cpu["SP"])
    u.reg_write(
        arm.UC_ARM64_REG_NZCV,
        sum(int(v) << (31 - i) for i, v in enumerate(cpu["Flags"])),
    )
    u.reg_write(arm.UC_ARM64_REG_PC, cpu["PC"])


def write_hook(u, access, address, size, value, data):
    for a in range(address & ~4095, ((address + size + 4095) // 4096) * 4096, 4096):
        dirty.add(a)


def block_hook(u, address, size, data):
    if image_base <= address < image_base + (16 << 20):
        recent.append("text+0x%x" % (address - image_base))


reached = None


def stop_hook(u, address, size, data):
    global reached
    reached = address
    u.emu_stop()


u.hook_add(UC_HOOK_MEM_WRITE, write_hook)
u.hook_add(UC_HOOK_BLOCK, block_hook)
u.hook_add(UC_HOOK_CODE, stop_hook, begin=0xF000000000000000, end=0xFFFFFFFFFFFFFFFF)
u.hook_add(UC_HOOK_CODE, stop_hook, begin=0x720000000000, end=0x720000000000)
output_slots = {
    0x700000240000: (0x1000E0580, 8),
    0x700000250000: (0x1000E0588, 4),
    0x700000260000: (0x1000E0594, 8),
    0x700000270000: (0x1000E059C, 4),
}
assert list(struct.unpack("<4Q", u.mem_read(0x700000100008, 32))) == list(output_slots)
assert all(a not in mapped for a in output_slots)
caller_pages_mapped = 0


def caller_write_fault(u, access, address, size, value, data):
    global caller_pages_mapped
    expected = output_slots.get(address)
    if (
        expected is None
        or expected != (u.reg_read(arm.UC_ARM64_REG_PC), size)
        or address in mapped
    ):
        return False
    u.mem_map(address, 4096)
    mapped.add(address)
    caller_pages_mapped += 1
    return True


u.hook_add(UC_HOOK_MEM_WRITE_UNMAPPED, caller_write_fault)
set_cpu(initial["CPU"])
calls = 0
started = time.monotonic()


def fail(kind, index, **details):
    report = {
        "schema": "crumb-whole-native-reference-v1",
        "matched": False,
        "failure": kind,
        "event": index,
        "validated_service_calls": calls,
        "recent_native_blocks": list(recent),
        "pac_authentication_verified": False,
        "service_abi_independently_verified": False,
        **details,
    }
    print(json.dumps(report), flush=True)
    sys.exit(1)


for index, line in enumerate(stream, 1):
    e = json.loads(line)
    reached = None
    validate_cpu(e["CPU"])
    if e["Type"] == "final" and e.get("Stop") == "incomplete":
        fail("recorded_invocation_incomplete", index)
    try:
        u.emu_start(
            u.reg_read(arm.UC_ARM64_REG_PC),
            e["CPU"]["PC"] if e["Type"] == "final" else 0,
            timeout=30_000_000,
        )
    except UcError as error:
        if not (
            e["Type"] == "final"
            and error.errno == UC_ERR_FETCH_UNMAPPED
            and u.reg_read(arm.UC_ARM64_REG_PC) == 0
        ):
            fail("native_execution_rejected", index, unicorn_error_code=error.errno)
    cpu = e["CPU"]
    pc = u.reg_read(arm.UC_ARM64_REG_PC)
    if pc != cpu["PC"]:
        fail(
            "control_flow_divergence_or_timeout",
            index,
            expected_external=e["Type"] == "external",
            returned_early=pc == 0,
        )
    differences = [
        i for i, (r, v) in enumerate(zip(regs, cpu["Registers"])) if u.reg_read(r) != v
    ]
    if differences:
        fail("register_divergence", index, register_indices=differences)
    flags = u.reg_read(arm.UC_ARM64_REG_NZCV)
    if [bool(flags & (1 << (31 - i))) for i in range(4)] != cpu["Flags"]:
        fail("flags_divergence", index)
    if u.reg_read(arm.UC_ARM64_REG_SP) != cpu["SP"]:
        fail("stack_divergence", index)
    expected = {int(a): v for a, v in (e.get("NativeWrites") or {}).items()}
    if dirty != set(expected):
        fail(
            "native_write_coverage_divergence",
            index,
            reference_pages=len(dirty),
            production_pages=len(expected),
        )
    for a, h in expected.items():
        if hashlib.sha256(u.mem_read(a, 4096)).hexdigest() != h:
            fail("native_written_memory_divergence", index)
    dirty.clear()
    if e["Type"] == "final":
        status = u.reg_read(arm.UC_ARM64_REG_X0) & 0xFFFFFFFF
        if status & 0x80000000:
            status -= 1 << 32
        print(
            json.dumps(
                {
                    "schema": "crumb-whole-native-reference-v1",
                    "matched": True,
                    "native_return_status": status if pc == 0 else None,
                    "caller_scalar_output_pages_mapped": caller_pages_mapped,
                    "validated_service_calls": calls,
                    "native_written_pages_checked": True,
                    "integer_registers_flags_sp_pc_checked": True,
                    "pac_authentication_verified": False,
                    "service_abi_independently_verified": False,
                    "sinf_validated_by_reference": False,
                    "seconds": round(time.monotonic() - started, 3),
                }
            ),
            flush=True,
        )
        break
    assert e["Type"] == "external"
    if e.get("Failed"):
        fail("captured_service_failed", index)
    for a in e.get("Removed") or []:
        if a in mapped:
            u.mem_unmap(a, 4096)
            mapped.remove(a)
    map_pages(
        {
            int(a): base64.b64decode(b, validate=True)
            for a, b in (e.get("Pages") or {}).items()
        }
    )
    set_cpu(e["After"])
    calls += 1
    if calls % 100 == 0:
        print(
            json.dumps({"stage": "replay_progress", "validated_service_calls": calls}),
            flush=True,
        )
else:
    fail("missing_final_event", index)
