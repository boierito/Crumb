#!/usr/bin/env python3
"""Replay bounded original fairplayd MIG dispatch with explicit host fixtures.

Research only: no device service, private user state, license or SINF is used.
The original daemon and ipsw fixup report are separately supplied public-firmware
inputs pinned by SHA-256. Every external effect modeled here is listed in JSON.
Security results are fixtures, not measurements. Unknown imports stop execution.
PAC is modeled as control transfer; cryptographic authentication is not modeled.
The optional initializer runs original code with explicit clock/entropy/memory/
single-thread-lock fixtures. It is not a complete daemon or hardware context.
"""
import argparse, json, re, struct, hashlib
import gzip
from collections import deque
from pathlib import Path
from unicorn import (
    Uc,
    UC_ARCH_ARM64,
    UC_MODE_ARM,
    UC_HOOK_CODE,
    UC_HOOK_MEM_INVALID,
    UC_HOOK_MEM_READ,
    UC_HOOK_MEM_WRITE,
    UcError,
)
import unicorn.arm64_const as ar
from capstone import Cs, CS_ARCH_ARM64, CS_MODE_ARM

p = argparse.ArgumentParser(description=__doc__)
p.add_argument("--daemon", type=Path, required=True)
p.add_argument("--fixups", type=Path, required=True)
p.add_argument("--limit", type=int, default=300000)
group = p.add_mutually_exclusive_group(required=True)
group.add_argument("--request-500", action="store_true")
group.add_argument("--request-1000", action="store_true")
p.add_argument("--caller-pid", type=int, default=123)
p.add_argument("--counterfactual-kernel-token", action="store_true")
p.add_argument("--ordinary-identifier", choices=["TESTTEAM01.com.certlium.crumb"])
p.add_argument("--prepare-runtime", action="store_true")
p.add_argument(
    "--fixture-invalid-client-word",
    action="store_true",
    help="explicitly use the invalid 0x20 client word; never a valid license request",
)
p.add_argument(
    "--trace-request",
    type=Path,
    help="bounded local register trace of public fixtures only",
)
p.add_argument("--mutation", choices=["size", "ool-type", "ool-count", "trailer"])
args = p.parse_args()
if args.request_1000 and not args.mutation and not args.fixture_invalid_client_word:
    raise ValueError(
        "request 1000 needs a matching client-produced in/out word; only the explicitly invalid fixture is supported"
    )
if not 1 <= args.limit <= 3000000 or not 0 <= args.caller_pid <= 2147483647:
    raise ValueError("execution or PID bounds")
data = args.daemon.read_bytes()
if (
    hashlib.sha256(data).hexdigest()
    != "94d32a82fdc2a80c46b7734b20895a935836a676a34b1c3e863cdb479a68cb01"
):
    raise ValueError("wrong original daemon; expected pinned 24A446 input")
fixup_text = args.fixups.read_bytes()
if (
    hashlib.sha256(fixup_text).hexdigest()
    != "6c5549d2cf4ec6cc7ba0beb841e4c05bf852caa8cc8d5582ac62e88c89a718a8"
):
    raise ValueError("wrong derived fixup report; expected pinned ipsw v3.1.739 output")
uc = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
uc.ctl_set_cpu_model(ar.UC_CPU_ARM64_MAX)
segments = []
pos = 32
for _ in range(struct.unpack_from("<I", data, 16)[0]):
    kind, size = struct.unpack_from("<II", data, pos)
    if kind == 0x19:
        va, mem, off, length = struct.unpack_from("<4Q", data, pos + 0x18)
        if mem and va:
            uc.mem_map(va & ~4095, (va + mem + 4095 & ~4095) - (va & ~4095))
            uc.mem_write(va, data[off : off + length])
            segments.append((va, mem, off, length))
    pos += size
external = {}
names = {}
fixups = 0
for line in fixup_text.decode("utf-8").splitlines():
    m = re.match("(0x[0-9a-f]+):\\s+raw: (0x[0-9a-f]+)\\s+(.+)", line)
    if not m:
        continue
    va, raw = (int(m[1], 16), int(m[2], 16))
    if not struct.unpack("<Q", uc.mem_read(va, 8))[0] == raw:
        raise ValueError("research-model contract violated")
    if "rebase:" in m[3]:
        target = int(re.search("target: (0x[0-9a-f]+)", m[3])[1], 16)
    elif "bind24:" in m[3]:
        name = m[3].split("\t")[-1].split("/", 1)[-1]
        if name not in names:
            address = 0x60000000 + len(names) * 4096
            names[name] = address
            external[address] = name
            uc.mem_map(address, 4096)
        target = names[name]
    else:
        raise ValueError(line)
    uc.mem_write(va, struct.pack("<Q", target))
    fixups += 1
uc.mem_map(0x70000000, 0x20000)
uc.mem_map(0x71000000, 0x20000)
uc.mem_map(0x72000000, 4096)
uc.reg_write(ar.UC_ARM64_REG_SP, 0x70018000)
uc.reg_write(ar.UC_ARM64_REG_X30, 0x72000000)
for n in range(8):
    uc.reg_write(getattr(ar, "UC_ARM64_REG_X" + str(n)), 0x71000000 + n * 4096)
uc.reg_write(ar.UC_ARM64_REG_X0, 123)
if args.request_500 or args.request_1000:
    if not (
        args.caller_pid > 0
        or (args.caller_pid == 0 and args.counterfactual_kernel_token)
    ):
        raise ValueError("kernel token requires explicit counterfactual flag")
    args.entry = 0x1000029E0 if args.request_500 else 0x10000255C
    request = bytearray(160)
    size, request_id = (0x50, 500) if args.request_500 else (0x3C, 1000)
    struct.pack_into("<6I", request, 0, 0x80000000, size, 0, 123, 0, request_id)
    struct.pack_into("<I", request, 0x18, 1)
    filename = (
        b"/nonexistent/CrumbResearch.app/CrumbResearch\x00"
        if args.request_500
        else bytes(32)
    )
    struct.pack_into("<Q", request, 0x1C, 0x71002000)
    request[0x27] = 1
    struct.pack_into("<I", request, 0x28, len(filename))
    if args.request_500:
        struct.pack_into("<I", request, 0x34, len(filename))
        struct.pack_into("<3Q", request, 0x38, 0x100000C, 2, 0)
    else:
        # Unlike request 500, +0x34 is NOT another OOL count. Its valid
        # producer is unresolved. Preserve 0x20 only as a labeled rejection
        # control; do not set the observed bit merely to advance execution.
        struct.pack_into(
            "<I", request, 0x34, 0x20 if args.fixture_invalid_client_word else 0
        )
        struct.pack_into("<I", request, 0x38, len(filename))
    struct.pack_into("<2I", request, size, 0, 0x34)
    struct.pack_into(
        "<8I", request, size + 0x14, 501, 501, 501, 501, 501, args.caller_pid, 1, 1
    )
    if args.mutation == "size":
        struct.pack_into("<I", request, 4, size - 4)
    if args.mutation == "ool-type":
        request[0x27] = 0
    if args.mutation == "ool-count":
        struct.pack_into("<I", request, 0x28, len(filename) + 1)
    if args.mutation == "trailer":
        struct.pack_into("<I", request, size, 1)
    uc.mem_write(0x71000000, bytes(request))
    uc.mem_write(0x71002000, filename)
    uc.reg_write(ar.UC_ARM64_REG_X0, 0x71000000)
    uc.reg_write(ar.UC_ARM64_REG_X1, 0x71001000)
uc.mem_write(names["___stack_chk_guard"], struct.pack("<Q", 0x12345678))
uc.mem_write(names["_kCFAllocatorDefault"], bytes(8))
uc.mem_write(names["_kCFAllocatorNull"], bytes(8))
if args.ordinary_identifier:
    if not args.ordinary_identifier == "TESTTEAM01.com.certlium.crumb":
        raise ValueError("only the explicit public ordinary-app fixture is supported")
    uc.mem_write(0x71011000, args.ordinary_identifier.encode() + b"\x00")
md = Cs(CS_ARCH_ARM64, CS_MODE_ARM)
stats = {"instructions": 0, "pac_control_transfers_modeled": 0}
events = []
last = []
runtime = []
x20_history = []
x0_eight_origins = []
last_x0 = [None]
input_reads = []
status_writes = []
request_trace = deque(maxlen=32768)
trace_count = [0]
request_phase = [not args.prepare_runtime]
arena = [0x74000000]
allocations = {}


def cstr(a):
    try:
        b = bytes(uc.mem_read(a, 256)).split(b"\x00")[0]
        return b.decode("ascii") if b and all((32 <= c < 127 for c in b)) else None
    except (UcError, UnicodeError):
        return None


def regs():
    return [uc.reg_read(getattr(ar, "UC_ARM64_REG_X" + str(n))) for n in range(8)]


def hook(e, pc, size, _):
    stats["instructions"] += 1
    if args.trace_request and request_phase[0]:
        request_trace.append(
            [pc]
            + [e.reg_read(getattr(ar, "UC_ARM64_REG_X" + str(n))) for n in range(31)]
            + [e.reg_read(ar.UC_ARM64_REG_SP)]
        )
        trace_count[0] += 1
    value0 = e.reg_read(ar.UC_ARM64_REG_X0)
    if value0 == 8 and last_x0[0] != 8:
        x0_eight_origins.append(last[-1] if last else None)
        x0_eight_origins[:] = x0_eight_origins[-24:]
    last_x0[0] = value0
    value = e.reg_read(ar.UC_ARM64_REG_X20)
    if not x20_history or value != x20_history[-1]["value"]:
        x20_history.append({"after_pc": last[-1] if last else None, "value": value})
        x20_history[:] = x20_history[-0x18:]
    if pc == 0x72000000:
        events.append({"event": "returned", "x0": e.reg_read(ar.UC_ARM64_REG_X0)})
        e.emu_stop()
        return
    if pc in external:
        name = external[pc]
        if name in ("_malloc", "_malloc_type_malloc", "_calloc"):
            values = regs()
            amount = values[0] * values[1] if name == "_calloc" else values[0]
            if not 0 < amount <= 1024 * 1024:
                raise ValueError("research-model contract violated")
            length = amount + 16383 & ~16383
            address = arena[0]
            arena[0] += length
            if not arena[0] < 0x78000000:
                raise ValueError("research-model contract violated")
            uc.mem_map(address, length, 3)
            allocations[address] = amount
            runtime.append(
                {
                    "name": name,
                    "effect": "bounded host allocation, no executable permission or license material",
                    "bytes": amount,
                }
            )
            e.reg_write(ar.UC_ARM64_REG_X0, address)
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if name in ("_free", "_malloc_type_free"):
            address = e.reg_read(ar.UC_ARM64_REG_X0)
            if not (address == 0 or address in allocations):
                raise ValueError("research-model contract violated")
            if address:
                del allocations[address]
            runtime.append({"name": name, "effect": "release tracked host allocation"})
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if name == "_pthread_rwlock_init":
            if not e.reg_read(ar.UC_ARM64_REG_X1) == 0:
                raise ValueError("research-model contract violated")
            uc.mem_read(e.reg_read(ar.UC_ARM64_REG_X0), 256)
            runtime.append(
                {
                    "name": name,
                    "effect": "single-thread host lock fixture initialized; opaque storage not consumed natively",
                }
            )
            e.reg_write(ar.UC_ARM64_REG_X0, 0)
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if name == "_gettimeofday":
            if not e.reg_read(ar.UC_ARM64_REG_X1) == 0:
                raise ValueError("research-model contract violated")
            e.mem_write(
                e.reg_read(ar.UC_ARM64_REG_X0), struct.pack("<qii", 1793318400, 0, 0)
            )
            runtime.append(
                {
                    "name": name,
                    "effect": "explicit fixed host-clock fixture; not device time",
                }
            )
            e.reg_write(ar.UC_ARM64_REG_X0, 0)
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if name == "_arc4random":
            runtime.append(
                {
                    "name": name,
                    "effect": "public deterministic startup-entropy fixture, not an account credential or license key",
                }
            )
            e.reg_write(ar.UC_ARM64_REG_X0, 0x6A09E667)
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if name == "_mmap":
            values = regs()
            if not (
                values[0] == 0
                and values[2] == 3
                and (values[3] == 0x1002)
                and (values[4] & 0xFFFFFFFF == 0xFFFFFFFF)
                and (values[5] == 0)
            ):
                raise ValueError("research-model contract violated")
            length = values[1] + 16383 & ~16383
            if not (0 < length <= 64 * 1024 * 1024 and arena[0] + length < 0x78000000):
                raise ValueError("research-model contract violated")
            address = arena[0]
            arena[0] += length
            uc.mem_map(address, length, 3)
            runtime.append(
                {
                    "name": name,
                    "effect": "bounded anonymous read/write arena in host model",
                    "bytes": length,
                }
            )
            e.reg_write(ar.UC_ARM64_REG_X0, address)
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if name == "_sysconf" and e.reg_read(ar.UC_ARM64_REG_X0) == 29:
            runtime.append(
                {
                    "name": name,
                    "effect": "explicit 16 KiB page-size model for iPhone17,1",
                }
            )
            e.reg_write(ar.UC_ARM64_REG_X0, 16384)
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if args.request_1000 and name in (
            "_CFStringCreateWithCStringNoCopy",
            "_SecTaskCreateWithAuditToken",
            "_SecTaskCopyValueForEntitlement",
            "_CFRelease",
            "_CFStringGetTypeID",
            "_CFGetTypeID",
            "_CFStringGetCString",
        ):
            values = regs()
            if name == "_CFStringCreateWithCStringNoCopy":
                text = cstr(values[1])
                if not (text and len(text) < 200):
                    raise ValueError("research-model contract violated")
                if not values[2] == 0x600:
                    raise ValueError("research-model contract violated")
                runtime.append(
                    {
                        "name": name,
                        "effect": "opaque CFString handle for original entitlement name",
                        "string": text,
                    }
                )
                e.reg_write(ar.UC_ARM64_REG_X0, values[1])
            elif name == "_SecTaskCreateWithAuditToken":
                token = struct.unpack("<8I", e.mem_read(values[1], 32))
                if not token[5] == args.caller_pid:
                    raise ValueError("research-model contract violated")
                runtime.append(
                    {
                        "name": name,
                        "effect": "opaque task handle for synthetic caller, no privilege",
                        "pid": token[5],
                    }
                )
                e.reg_write(ar.UC_ARM64_REG_X0, 0x71010000)
            elif name == "_SecTaskCopyValueForEntitlement":
                if not values[0] == 0x71010000:
                    raise ValueError("research-model contract violated")
                if (
                    cstr(values[1]) == "application-identifier"
                    and args.ordinary_identifier
                ):
                    runtime.append(
                        {
                            "name": name,
                            "effect": "explicit ordinary-identifier fixture; not actual signing identity",
                            "value": args.ordinary_identifier,
                        }
                    )
                    e.reg_write(ar.UC_ARM64_REG_X0, 0x71011000)
                    e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
                    return
                if cstr(values[1]) != "fairplay-client":
                    events.append(
                        {"event": "unmodeled_entitlement", "name": cstr(values[1])}
                    )
                    e.emu_stop()
                    return
                runtime.append(
                    {
                        "name": name,
                        "effect": "explicit ABSENT entitlement model; not a measured service result",
                        "entitlement": "fairplay-client",
                    }
                )
                e.reg_write(ar.UC_ARM64_REG_X0, 0)
            elif name in ("_CFStringGetTypeID", "_CFGetTypeID"):
                if name == "_CFGetTypeID":
                    if not values[0] == 0x71011000:
                        raise ValueError("research-model contract violated")
                runtime.append(
                    {"name": name, "effect": "type of modeled string object"}
                )
                e.reg_write(ar.UC_ARM64_REG_X0, 0x33)
            elif name == "_CFStringGetCString":
                if not (values[0] == 0x71011000 and values[3] == 0x600):
                    raise ValueError("research-model contract violated")
                raw = args.ordinary_identifier.encode() + b"\x00"
                if not len(raw) <= values[2]:
                    raise ValueError("research-model contract violated")
                e.mem_write(values[1], raw)
                e.reg_write(ar.UC_ARM64_REG_X0, 1)
                runtime.append(
                    {
                        "name": name,
                        "effect": "bounded copy of ordinary identifier fixture",
                    }
                )
            else:
                runtime.append(
                    {"name": name, "effect": "release modeled opaque object"}
                )
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if external[pc] == "_audit_token_to_au32":
            values = struct.unpack(
                "<8I", e.mem_read(e.reg_read(ar.UC_ARM64_REG_X0), 32)
            )
            pointers = [
                e.reg_read(getattr(ar, "UC_ARM64_REG_X" + str(n))) for n in range(1, 8)
            ]
            terminal = struct.unpack(
                "<Q", e.mem_read(e.reg_read(ar.UC_ARM64_REG_SP), 8)
            )[0]
            if not terminal == 0:
                raise ValueError("terminal host lookup not modeled")
            for value, ptr in zip(values, pointers):
                if ptr:
                    e.mem_write(ptr, struct.pack("<I", value))
            runtime.append(
                {
                    "name": "_audit_token_to_au32",
                    "pid": values[5],
                    "effect": "token field copies; no privilege granted",
                }
            )
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        if external[pc] == "___chkstk_darwin":
            amount = e.reg_read(ar.UC_ARM64_REG_X9)
            sp = e.reg_read(ar.UC_ARM64_REG_SP)
            if not 0x70000000 <= sp - amount <= sp <= 0x70020000:
                raise ValueError("stack probe outside mapped stack")
            runtime.append(
                {
                    "name": "___chkstk_darwin",
                    "bytes": amount,
                    "effect": "verified mapped stack, registers preserved",
                }
            )
            e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
            return
        values = regs()
        events.append(
            {
                "event": "unmodeled_import",
                "name": external[pc],
                "arguments": [hex(x) for x in values],
                "ascii_arguments": [cstr(x) for x in values],
            }
        )
        e.emu_stop()
        return
    if pc in (0x100001B40, 0x100001CEC):
        events.append(
            {
                "event": "entitlement_helper",
                "pc": hex(pc),
                "name": cstr(e.reg_read(ar.UC_ARM64_REG_X1)),
            }
        )
    if pc == 0x100B2BD58:
        events.append(
            {
                "event": "pid_zero_branch",
                "pid": args.caller_pid,
                "target": hex(e.reg_read(ar.UC_ARM64_REG_X8)),
            }
        )
    if pc == 0x100739F14:
        events.append(
            {
                "event": "allocation_or_context_call",
                "target": hex(e.reg_read(ar.UC_ARM64_REG_X8)),
                "x0": hex(e.reg_read(ar.UC_ARM64_REG_X0)),
            }
        )
    if pc == 0x100739F18:
        events.append(
            {
                "event": "context_call_return",
                "x0": hex(e.reg_read(ar.UC_ARM64_REG_X0)),
                "tail": list(last),
            }
        )
    if pc == 0x100777278:
        events.append(
            {
                "event": "callback_slot_read",
                "address": hex(
                    e.reg_read(ar.UC_ARM64_REG_X1) + e.reg_read(ar.UC_ARM64_REG_X2)
                    & (1 << 64) - 1
                ),
            }
        )
    if (
        args.request_500
        and pc == 0x100002B14
        or (args.request_1000 and pc == 0x100002694)
    ):
        events.append(
            {
                "event": "native_request_return",
                "status_u32": e.reg_read(ar.UC_ARM64_REG_X0) & 0xFFFFFFFF,
            }
        )
    word = struct.unpack("<I", bytes(e.mem_read(pc, 4)))[0]
    last.append(hex(pc))
    last[:] = last[-12:]
    if word in (0xD503237F, 0xD50323FF, 0xD503233F, 0xD50323BF):
        e.reg_write(ar.UC_ARM64_REG_PC, pc + 4)
        stats["pac_control_transfers_modeled"] += 1
    elif word in (0xD65F0FFF, 0xD65F0BFF):
        e.reg_write(ar.UC_ARM64_REG_PC, e.reg_read(ar.UC_ARM64_REG_X30))
        stats["pac_control_transfers_modeled"] += 1
    elif word & 0xFFFFFC00 in (
        0xDAC10000,
        0xDAC10400,
        0xDAC10800,
        0xDAC10C00,
    ) or word & 0xFFFFFFE0 in (0xDAC123E0, 0xDAC127E0):
        e.reg_write(ar.UC_ARM64_REG_PC, pc + 4)
        stats["pac_control_transfers_modeled"] += 1
    elif word & 0xFEFFF800 in (0xD61F0800, 0xD63F0800) or word & 0xFEFFF800 in (
        0xD71F0800,
        0xD73F0800,
    ):
        target = e.reg_read(getattr(ar, "UC_ARM64_REG_X" + str(word >> 5 & 31)))
        if word & 0x00200000:
            e.reg_write(ar.UC_ARM64_REG_X30, pc + 4)
        e.reg_write(ar.UC_ARM64_REG_PC, target)
        stats["pac_control_transfers_modeled"] += 1


uc.hook_add(UC_HOOK_CODE, hook)


def invalid(e, access, address, size, value, _):
    events.append(
        {
            "event": "unmapped_memory",
            "address": hex(address),
            "size": size,
            "x8": hex(e.reg_read(ar.UC_ARM64_REG_X8)),
            "x20": hex(e.reg_read(ar.UC_ARM64_REG_X20)),
            "sp": hex(e.reg_read(ar.UC_ARM64_REG_SP)),
        }
    )
    return False


uc.hook_add(UC_HOOK_MEM_INVALID, invalid)


def input_read(e, access, address, size, value, _):
    if len(input_reads) < 64:
        input_reads.append(
            {
                "pc": hex(e.reg_read(ar.UC_ARM64_REG_PC)),
                "offset": address - 0x71002000,
                "size": size,
            }
        )


uc.hook_add(UC_HOOK_MEM_READ, input_read, begin=0x71002000, end=0x7100203F)


def status_write(e, access, address, size, value, _):
    status_writes.append(
        {
            "pc": hex(e.reg_read(ar.UC_ARM64_REG_PC)),
            "address": hex(address),
            "size": size,
            "value": hex(value & ((1 << (8 * size)) - 1)),
            "preceding_pcs": list(last),
        }
    )
    status_writes[:] = status_writes[-32:]


uc.hook_add(UC_HOOK_MEM_WRITE, status_write, begin=0x70017F00, end=0x70017F03)
try:
    if args.prepare_runtime:
        context = uc.context_save()
        uc.emu_start(0x10073A790, 0, count=20000000)
        if not events or events[-1]["event"] != "returned":
            raise ValueError("native initializer did not return; request not executed")
        events[-1]["event"] = "native_initializer_returned"
        events[-1]["return_value_is_status"] = False
        events.append(
            {
                "event": "callback_after_initializer",
                "slot": "0x1011320f8",
                "value": hex(struct.unpack("<Q", uc.mem_read(0x1011320F8, 8))[0]),
            }
        )
        uc.context_restore(context)
        request_phase[0] = True
    uc.emu_start(args.entry, 0, count=args.limit)
except UcError as ex:
    pc = uc.reg_read(ar.UC_ARM64_REG_PC)
    raw = (
        bytes(uc.mem_read(pc, 4))
        if any((v <= pc < v + m for v, m, _, _ in segments))
        else b""
    )
    events.append(
        {
            "event": "emulation_error",
            "error": str(ex),
            "pc": hex(pc),
            "instruction": [i.mnemonic + " " + i.op_str for i in md.disasm(raw, pc)],
            "registers": [hex(x) for x in regs()],
        }
    )
if args.trace_request:
    with gzip.open(args.trace_request, "wt", encoding="utf-8") as output:
        json.dump(
            {
                "schema": "crumb-fairplayd-public-register-trace-v1",
                "private_user_inputs_used": False,
                "truncated": trace_count[0] > len(request_trace),
                "fields": ["pc"] + ["x" + str(n) for n in range(31)] + ["sp"],
                "external_tokens": {hex(k): v for k, v in external.items()},
                "rows": list(request_trace),
            },
            output,
        )

print(
    json.dumps(
        {
            "entry": hex(args.entry),
            "fixups_verified": fixups,
            "initial_context": "synthetic, not captured process state",
            "counterfactual_kernel_token": args.counterfactual_kernel_token,
            "crypto_service_results_modeled": False,
            "security_framework_results_measured": False,
            "sinf_validated": False,
            "private_user_inputs_used": False,
            "valid_client_request_produced": False,
            "invalid_client_word_fixture": bool(
                args.request_1000 and args.fixture_invalid_client_word
            ),
            "image_page_protection_enforced": False,
            "pac_authenticated": False,
            "termination": (
                events[-1]["event"]
                if events
                and events[-1]["event"]
                in (
                    "returned",
                    "unmodeled_import",
                    "unmodeled_entitlement",
                    "emulation_error",
                )
                else "instruction_limit"
            ),
            "handler_returned": bool(events and events[-1]["event"] == "returned"),
            "reply_status_u32": (
                struct.unpack("<I", uc.mem_read(0x71001020, 4))[0]
                if events and events[-1]["event"] == "returned"
                else None
            ),
            "mutation": args.mutation,
            "modeled_runtime": runtime,
            "device_service_executed": False,
            **stats,
            "events": events,
            "last": last,
            "x20_history": x20_history,
            "x0_eight_origins": x0_eight_origins,
            "opaque_input_reads_first_64": input_reads,
            "worker_1000_status_slot_writes_last_32": status_writes,
        },
        indent=2,
    )
)
