#!/usr/bin/env python3
"""Audit one public 26A428 CFG jump against the hash-pinned STARTUP model.

No customer capture, SIDB, keybag, account ID or iPhone is needed.
This is a static build-profile coverage check, not executable translation.
"""
import gzip
import hashlib
import io
import json
from pathlib import Path
import re
import struct
import sys
import tarfile

ROOT = Path(__file__).resolve().parent.parent
REV = "f7242cce58077a4af8db738ed295ee23ad8f067d"
ARCHIVE = ROOT / "third_party" / "archives" / f"ipatool-fairplay-{REV}.tar.gz"
EXPECTED_SHA = "e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39"
TARGET = 0x1008fff2c
SOURCE_BLOCK = 0x100b2f9d8
SIZE = 0x1000
expected_suffix = "/tools/fairplay-recovery/profiles/macos26a428-account-entry-public.model.gz"
cfg_suffix = "/tools/fairplay-recovery/generated_raw_source_bootstrap_cfg.go"

def checked_member(archive, suffix):
    names = [m for m in archive.getmembers() if m.isfile() and m.name.endswith(suffix)]
    if len(names) != 1:
        raise RuntimeError(f"Expected one pinned {suffix}, got {len(names)}")
    return archive.extractfile(names[0]).read()

def read_profile(data):
    stream = io.BytesIO(gzip.decompress(data))
    def take(n):
        result = stream.read(n)
        if len(result) != n:
            raise ValueError("truncated public model")
        return result
    if take(8) != b"FPGOOR1\\x00" or struct.unpack("<I", take(4))[0] != 2:
        raise ValueError("unexpected pinned model magic/version")
    binary_hash = take(32).hex()
    take(8 + 31 * 8 + 8 + 4)
    count = struct.unpack("<I", take(4))[0]
    if not 0 < count <= 8192:
        raise ValueError("invalid public page count")
    pages = {}
    for _ in range(count):
        address = struct.unpack("<Q", take(8))[0]
        if address % SIZE or address in pages:
            raise ValueError("invalid/duplicate public page address")
        pages[address] = take(SIZE)
    return binary_hash, pages

def main():
    if hashlib.sha256(ARCHIVE.read_bytes()).hexdigest() != EXPECTED_SHA:
        raise RuntimeError("modern FairPlay archive integrity check failed")
    with tarfile.open(ARCHIVE, "r:gz") as tar:
        model_data = checked_member(tar, expected_suffix)
        cfg_source = checked_member(tar, cfg_suffix).decode("utf-8")
    binary_hash, pages = read_profile(model_data)
    static_blocks = set(int(x, 16) for x in
        re.findall(r"func rawSourceBootstrapBlock([0-9a-f]+)\\(", cfg_source))
    translated_instructions = set(int(x, 16) for x in
        re.findall(r"// 0x([0-9a-f]+):", cfg_source))
    if SOURCE_BLOCK not in static_blocks or TARGET in static_blocks:
        raise RuntimeError("expected original missing-target regression no longer applies")
    start = cfg_source.index(f"func rawSourceBootstrapBlock{SOURCE_BLOCK:x}(")
    end = cfg_source.find("\\nfunc ", start + 5)
    function = cfg_source[start:end if end >= 0 else None]
    if "braa x10, x17" not in function or "ldr x10, [x30, w10, sxtw #3]" not in function:
        raise RuntimeError("source block is no longer the pinned indirect branch")
    page_base = TARGET & ~(SIZE - 1)
    page = pages.get(page_base)
    offset = TARGET - page_base
    raw_word = page[offset:offset+4] if page is not None else None
    neighbors = sorted((p for p in static_blocks if p <= TARGET), reverse=True)
    report = {
        "schema": "crumb-pinned-public-cfg-audit-v1",
        "profile_binary_sha256": binary_hash,
        "source": f"text+0x{SOURCE_BLOCK-0x100000000:x}",
        "source_is_authenticated_indirect_branch": True,
        "target": f"text+0x{TARGET-0x100000000:x}",
        "target_is_translated_block_entry": TARGET in static_blocks,
        "target_is_in_translated_instruction_trace": TARGET in translated_instructions,
        "target_page_present_in_public_model": page is not None,
        "target_instruction_word_nonzero": bool(raw_word and any(raw_word)),
        "nearest_previous_translated_block": (
            f"text+0x{neighbors[0]-0x100000000:x}" if neighbors else None
        ),
        "private_material_read": False,
    }
    print(json.dumps(report, sort_keys=True, indent=2))
    if len(sys.argv) == 2:
        Path(sys.argv[1]).write_text(json.dumps(report, sort_keys=True, indent=2) + "\\n")

if __name__ == "__main__":
    main()
