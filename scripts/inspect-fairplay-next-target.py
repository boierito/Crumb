#!/usr/bin/env python3
"""Coverage census for the NEXT on-device FairPlay bootstrap miss.

On a macOS 27 runner, inspect only the sha256-pinned system executable
already verified by Crumb. Never publish executable instructions, private
enrollment state, account identifiers, binary segments, or machine bytes.
This file performs a read-only discovery, NOT decryption or code translation.
"""
from __future__ import annotations

import hashlib
import json
import re
import shutil
import subprocess
import sys
from pathlib import Path

PINNED_SHA256 = "1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784"
TARGET = 0x100A4BE38
CANDIDATES = (
    "/System/Library/PrivateFrameworks/CoreFP.framework/Versions/A/fairplayd",
    "/System/Library/PrivateFrameworks/CoreFP.framework/fairplayd",
)

def run(argv: list[str]) -> tuple[int, str]:
    try:
        completed = subprocess.run(argv, capture_output=True, text=True, timeout=20, check=False)
        return completed.returncode, completed.stdout
    except (OSError, subprocess.TimeoutExpired):
        return -1, ""

def branch_kind(mnemonic: str) -> str:
    op = mnemonic.lower()
    if op in ("braa", "brab", "br", "ret", "retaa", "retab"):
        return "indirect"
    if op in ("bl", "blr", "blraa", "blrab"):
        return "call"
    if op in ("cbz", "cbnz", "tbz", "tbnz") or op.startswith("b."):
        return "conditional"
    if op == "b":
        return "direct"
    return ""

def summarize_disassembly(listing: str, start: int, end: int) -> dict:
    """Only finite structural information; deliberately discard all bytes."""
    addresses: list[int] = []
    branch_count = 0
    indirect = 0
    branch_first = ""
    for line in listing.splitlines():
        match = re.match(r"^\s*([0-9a-f]+):\s+[0-9a-f]{8}\s+([a-zA-Z][a-zA-Z0-9_.]*)", line)
        if not match:
            continue
        address = int(match.group(1), 16)
        if address < start or address >= end:
            continue
        addresses.append(address)
        kind = branch_kind(match.group(2))
        if kind:
            branch_count += 1
            indirect += kind == "indirect"
            if not branch_first:
                branch_first = kind
    addresses.sort()
    return {
        "target_is_disassemblable": bool(addresses and addresses[0] == start),
        "instructions_observed": len(addresses),
        "branch_instructions_observed": branch_count,
        "indirect_branches_observed": indirect,
        "first_branch_kind": branch_first or "not_observed",
        "sequential_word_alignment": (all(b-a == 4 for a, b in zip(addresses, addresses[1:]))
                                      if addresses else False),
    }

def main() -> None:
    status = {
        "schema": "crumb-public-target-coverage-v1",
        "target": "text+0xa4be38",
        "pinned_binary_sha256": PINNED_SHA256,
        "exact_binary_verified": False,
        "code_bytes_exported": False,
        "private_state_accessed": False,
        "jit_required": False,
        "jailbreak_required": False,
        "status": "exact_profile_not_found",
    }
    for file_name in CANDIDATES:
        source = Path(file_name)
        if not source.is_file() or source.is_symlink():
            continue
        if source.stat().st_size > 128 * 1024 * 1024:
            continue
        if hashlib.sha256(source.read_bytes()).hexdigest() != PINNED_SHA256:
            continue
        status["exact_binary_verified"] = True
        status["status"] = "disassembler_unavailable"
        rc, tool = run(["xcrun", "-f", "llvm-objdump"])
        executable = tool.strip() if rc == 0 else shutil.which("llvm-objdump")
        if not executable:
            break
        rc, stdout = run([executable, "--disassemble",
                          "--start-address="+hex(TARGET),
                          "--stop-address="+hex(TARGET+0x200), file_name])
        if rc == 0:
            status.update(summarize_disassembly(stdout, TARGET, TARGET+0x200))
            status["status"] = ("target_in_exact_binary" if status["target_is_disassemblable"]
                                else "target_not_disassemblable")
        break
    encoded = json.dumps(status, sort_keys=True, indent=2) + "\n"
    print(encoded)
    if len(sys.argv) > 1:
        destination = Path(sys.argv[1])
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_text(encoded)

if __name__ == "__main__":
    main()
