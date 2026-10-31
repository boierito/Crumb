#!/usr/bin/env python3
"""Collect bounded public-cache call stubs without whole-OS symbol indexing.

Development only. Reads the already pinned cache selected by the acquisition
script. No service invocation, request synthesis or user license input.
"""
import argparse
import hashlib
import json
import re
import struct
import subprocess
from pathlib import Path

# Direct external calls observed in the pinned StoreServices wrapper bodies.
STUBS = (
    0x1D0154A30,
    0x1D0161E30,
    0x1D015CB80,
    0x1D015B5A0,
    0x1D015B5C0,
    0x1D015C650,
    0x1D015C660,
)
SYMBOLS = {
    "libMobileGestalt.dylib": ("_MGCopyAnswer",),
    "CoreFoundation": (
        "_CFStringGetSystemEncoding",
        "_CFStringCreateExternalRepresentation",
        "_CFDataGetBytePtr",
        "_CFDataGetLength",
    ),
    "libsystem_pthread.dylib": ("_pthread_mutex_lock", "_pthread_mutex_unlock"),
}


def pointer_slot(address, code):
    """Decode only ADRP/ADD/LDR unsigned stubs; refuse other shapes."""
    registers = {}
    for offset in range(0, min(len(code), 16), 4):
        (word,) = struct.unpack_from("<I", code, offset)
        rd, rn = word & 31, (word >> 5) & 31
        if word & 0x9F000000 == 0x90000000:  # ADRP
            imm = ((word >> 5) & 0x7FFFF) << 2 | ((word >> 29) & 3)
            if imm & (1 << 20):
                imm -= 1 << 21
            registers[rd] = ((address + offset) & ~4095) + (imm << 12)
        elif word & 0xFF800000 == 0x91000000 and rn in registers:
            registers[rd] = registers[rn] + (
                ((word >> 10) & 4095) << (12 if word & (1 << 22) else 0)
            )
        elif word & 0xFFC00000 == 0xF9400000 and rn in registers:
            return registers[rn] + (((word >> 10) & 4095) << 3)
        else:
            return None
    return None


def direct_target(address, code):
    """Decode exactly ADRP Xn; ADD Xn,Xn,#imm; BR Xn stubs."""
    if len(code) < 12:
        return None
    adrp, add, branch = struct.unpack_from("<3I", code)
    reg = adrp & 31
    if (
        reg == 31
        or adrp & 0x9F000000 != 0x90000000
        or add & 0xFF800000 != 0x91000000
        or add & 31 != reg
        or (add >> 5) & 31 != reg
        or branch != 0xD61F0000 | (reg << 5)
    ):
        return None
    imm = ((adrp >> 5) & 0x7FFFF) << 2 | ((adrp >> 29) & 3)
    if imm & (1 << 20):
        imm -= 1 << 21
    return (
        (address & ~4095)
        + (imm << 12)
        + (((add >> 10) & 4095) << (12 if add & (1 << 22) else 0))
    )


def resolve_direct_entries(report):
    """Resolve collected direct branches using independently looked-up symbols."""
    for entry in report["stubs"]:
        code = bytes.fromhex(entry["code_hex"])
        if hashlib.sha256(code).hexdigest() != entry["sha256"]:
            raise ValueError("Collected stub digest mismatch")
        direct = direct_target(int(entry["address"], 16), code)
        if direct is not None:
            entry.pop("unresolved", None)
            entry.update(target=hex(direct), target_kind="direct_branch")
        if "target" in entry:
            target = int(entry["target"], 16)
            matches = []
            for lookup in report["symbol_lookups"]:
                m = re.match(r"(0x[0-9a-fA-F]+):", lookup["tool_output"])
                if m and int(m[1], 16) == target:
                    matches.append(
                        {"image": lookup["image"], "symbol": lookup["symbol"]}
                    )
            entry["exact_symbol_matches"] = matches
    return report


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--ipsw", required=True)
    p.add_argument("--cache", required=True)
    p.add_argument("--output", type=Path, required=True)
    args = p.parse_args()

    def run(*argv, raw=False):
        result = subprocess.run(
            [args.ipsw, "--no-color", *map(str, argv)],
            capture_output=True,
            check=True,
            timeout=60,
        )
        return result.stdout if raw else result.stdout.decode()

    report = {
        "schema": "crumb-public-client-stubs-v1",
        "build": "24A446",
        "service_executed": False,
        "sinf_validated": False,
        "stubs": [],
        "symbol_lookups": [],
    }
    for address in STUBS:
        code = run(
            "dyld", "dump", args.cache, hex(address), "--size", 16, "--bytes", raw=True
        )
        if len(code) != 16:
            raise ValueError("Incomplete stub read")
        entry = {
            "address": hex(address),
            "code_hex": code.hex(),
            "sha256": hashlib.sha256(code).hexdigest(),
        }
        slot = pointer_slot(address, code)
        if slot is not None:
            raw = run(
                "dyld", "dump", args.cache, hex(slot), "--size", 8, "--bytes", raw=True
            )
            if len(raw) != 8:
                raise ValueError("Incomplete pointer read")
            decoded = run(
                "dyld", "dump", args.cache, hex(slot), "--count", 1, "--addr"
            ).strip()
            if not re.fullmatch(r"0x[0-9a-fA-F]+", decoded):
                raise ValueError("Unexpected decoded pointer output")
            entry.update(
                pointer_slot=hex(slot), pointer_raw_hex=raw.hex(), target=decoded
            )
        else:
            entry["unresolved"] = "stub shape outside bounded decoder"
        report["stubs"].append(entry)
    for image, symbols in SYMBOLS.items():
        for symbol in symbols:
            output = run("dyld", "symaddr", args.cache, symbol, "--image", image)
            report["symbol_lookups"].append(
                {"image": image, "symbol": symbol, "tool_output": output.strip()}
            )
    args.output.write_text(json.dumps(resolve_direct_entries(report), indent=2) + "\n")


if __name__ == "__main__":
    main()
