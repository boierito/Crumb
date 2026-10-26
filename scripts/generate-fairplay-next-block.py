#!/usr/bin/env python3
"""Generate a review candidate, never a release, from the exact public daemon.

Uses the translator in the integrity-checked vendored archive. No checkpoints,
account inputs, executable mappings, or runtime code generation are involved.
The generated function must be reviewed and differentially tested before use.
"""
import argparse
import hashlib
import importlib
import json
from pathlib import Path
import sys
import subprocess
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parent.parent
REV = "f7242cce58077a4af8db738ed295ee23ad8f067d"
ARCHIVE_HASH = "e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39"
BINARY_HASH = "1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784"
TARGET = 0x100A4BE38


def exact_system_binary(directory):
    for name in (
        "/System/Library/PrivateFrameworks/CoreFP.framework/Versions/A/fairplayd",
        "/System/Library/PrivateFrameworks/CoreFP.framework/fairplayd",
    ):
        source = Path(name)
        if not source.is_file() or source.is_symlink():
            continue
        candidates = [source]
        for arch in ("arm64e", "arm64"):
            thin = Path(directory) / ("fairplayd-" + arch)
            result = subprocess.run(["lipo", "-thin", arch, str(source), "-output", str(thin)], capture_output=True, timeout=30)
            if result.returncode == 0:
                candidates.append(thin)
        for path in candidates:
            if 0 < path.stat().st_size <= 128 * 1024 * 1024 and hashlib.sha256(path.read_bytes()).hexdigest() == BINARY_HASH:
                return path
    raise ValueError("exact system daemon not available; no candidate generated")


def block_instructions(database, decoder, translator, start, limit=256):
    """Only decode within a file-backed executable section; fail at any gap."""
    section = next((s for s in database.sections if s.contains_instructions
                    and s.vmaddr <= start < s.vmaddr + s.size), None)
    if section is None or start % 4:
        raise ValueError("target is not an aligned executable-section address")
    result = []
    for address in range(start, min(start + limit * 4, section.vmaddr + section.size), 4):
        offset = section.fileoff + address - section.vmaddr
        raw = database.image[offset:offset + 4]
        instruction = next(decoder.disasm(raw, address, count=1), None)
        if instruction is None or instruction.size != 4:
            raise ValueError("instruction cannot be decoded")
        result.append(instruction)
        if translator.branch_return_lines(decoder, instruction) is not None:
            return result
    raise ValueError("no supported block terminator within the bounded section")


def generate(database, decoder, translator, start=TARGET):
    instructions = block_instructions(database, decoder, translator, start)
    # Reuse the pinned generator's instruction semantics and per-PC stop checks.
    source = translator.generate_source(database, [start],
                                       {start: len(instructions) * 4}, 0)
    marker = f"func rawSourceBootstrapBlock{start:x}("
    function = source[source.index(marker):]
    imports = 'import "fmt"\n\n' if "fmt." in function else ""
    return ("// REVIEW CANDIDATE: not registered for execution or release.\n"
            f"// Exact daemon SHA256: {BINARY_HASH}\n"
            f"// Translator revision: {REV}\n"
            "package rawsource\n\n" + imports + function)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=Path, help="Exact thin daemon; defaults to build-host system daemon")
    parser.add_argument("--target", type=lambda value: int(value, 0), default=TARGET,
                        help="Instruction address in the exact executable section")
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--reference-output", type=Path,
                        help="Optional public instruction fixture for independent execution tests")
    args = parser.parse_args()
    archive = ROOT / "third_party/archives" / f"ipatool-fairplay-{REV}.tar.gz"
    if hashlib.sha256(archive.read_bytes()).hexdigest() != ARCHIVE_HASH:
        raise ValueError("translator archive integrity mismatch")
    with tempfile.TemporaryDirectory(prefix="crumb-aot-translator-") as directory:
        binary = args.binary if args.binary is not None else exact_system_binary(directory)
        if (binary.is_symlink() or not binary.is_file()
                or not 0 < binary.stat().st_size <= 128 * 1024 * 1024
                or hashlib.sha256(binary.read_bytes()).hexdigest() != BINARY_HASH):
            raise ValueError("exact thin daemon hash required")
        with tarfile.open(archive, "r:gz") as tar:
            tar.extractall(directory, filter="data")
        python = Path(directory) / f"ipatool-{REV}/tools/fairplay-recovery/python"
        sys.path.insert(0, str(python))
        translator = importlib.import_module("generate_macos26a428_bootstrap_cfg_go")
        from capstone import CS_ARCH_ARM64, CS_MODE_ARM, Cs
        from raw_macho_db import RawMachODatabase
        decoder = Cs(CS_ARCH_ARM64, CS_MODE_ARM)
        decoder.detail = True
        database = RawMachODatabase(binary)
        candidate = generate(database, decoder, translator, args.target)
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(candidate)
        if args.reference_output:
            instructions = block_instructions(database, decoder, translator, args.target)
            reference = {"binary_sha256": BINARY_HASH, "start": args.target,
                         "instructions": [{"address": i.address, "hex": bytes(i.bytes).hex(),
                                           "assembly": i.mnemonic + " " + i.op_str}
                                          for i in instructions]}
            args.reference_output.parent.mkdir(parents=True, exist_ok=True)
            args.reference_output.write_text(json.dumps(reference, indent=2) + "\n")
    print("Generated review candidate; dispatcher and release coverage unchanged.")


if __name__ == "__main__":
    main()
