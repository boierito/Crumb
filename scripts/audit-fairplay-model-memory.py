#!/usr/bin/env python3
"""Compare the public snapshot with its exact Mach-O image, without enrollment.

This is a development-host census, not a missing-memory recovery fallback.
In particular it never initializes absent writable globals from file bytes.
"""
import argparse
import gzip
import hashlib
import importlib.util
import json
from pathlib import Path
import struct
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parent.parent
REV = "f7242cce58077a4af8db738ed295ee23ad8f067d"
ARCHIVE_SHA = "e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39"
MODEL_SHA = "e8121726eab5a07a2b389f110522df08f2bea1aec94935b89b21578d77ccc786"
PAGE = 4096


def census(image, model):
    if len(image) < 32 or struct.unpack_from("<I", image)[0] != 0xFEEDFACF:
        raise ValueError("thin little-endian Mach-O required")
    count, size = struct.unpack_from("<II", image, 16)
    end = 32 + size
    if count > 4096 or end > len(image):
        raise ValueError("invalid Mach-O load commands")
    segments, fixups = [], []
    at = 32
    for _ in range(count):
        if at + 8 > end:
            raise ValueError("truncated load command")
        command, length = struct.unpack_from("<II", image, at)
        if length < 8 or at + length > end:
            raise ValueError("invalid load command length")
        if command == 0x19:
            if length < 72:
                raise ValueError("truncated segment")
            name = image[at+8:at+24].rstrip(b"\0").decode("ascii")
            vm, virtual_size, offset, file_size, maximum, initial = struct.unpack_from("<QQQQII", image, at+24)
            if file_size > virtual_size or offset + file_size > len(image):
                raise ValueError("invalid segment extent")
            segments.append(dict(name=name, address=vm, virtual_bytes=virtual_size,
                                 offset=offset, file_bytes=file_size,
                                 maximum_protection=maximum, initial_protection=initial))
        if command in (0x80000034, 0x80000022, 0x22):
            fixups.append(hex(command))
        at += length
    if at != end:
        raise ValueError("load command geometry mismatch")
    if len(model) < 316 or model[:8] != b"FPGOOR1\0":
        raise ValueError("invalid model")
    page_count = struct.unpack_from("<I", model, 312)[0]
    if page_count > 32768 or 316 + page_count * (PAGE+8) > len(model):
        raise ValueError("invalid model page geometry")
    pages = {}
    for i in range(page_count):
        pos = 316 + i * (PAGE+8)
        address = struct.unpack_from("<Q", model, pos)[0]
        if address % PAGE or address in pages:
            raise ValueError("invalid model page address")
        pages[address] = model[pos+8:pos+8+PAGE]
    accounted = set()
    for segment in segments:
        start = segment["address"]
        file_end = start + segment["file_bytes"]
        full_start = (start + PAGE-1) & ~(PAGE-1)
        full_end = file_end & ~(PAGE-1)
        captured, equal, different = 0, 0, 0
        for address, page in pages.items():
            if start <= address < start + segment["virtual_bytes"]:
                accounted.add(address)
                captured += 1
                if full_start <= address < full_end:
                    offset = segment["offset"] + address-start
                    if image[offset:offset+PAGE] == page:
                        equal += 1
                    else:
                        different += 1
        segment.update(captured_pages=captured, identical_file_pages=equal,
                       changed_file_pages=different,
                       full_file_pages=max(0, (full_end-full_start)//PAGE))
    return dict(schema="crumb-public-model-memory-census-v1", model_pages=page_count,
                non_image_pages=len(pages.keys()-accounted), segments=segments,
                relocation_commands=fixups, production_pages_added=0,
                enrollment_compatibility_proven=False)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=Path)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    spec = importlib.util.spec_from_file_location("exact", ROOT/"scripts/generate-fairplay-next-block.py")
    exact = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(exact)
    archive = ROOT/"third_party/archives"/f"ipatool-fairplay-{REV}.tar.gz"
    if hashlib.sha256(archive.read_bytes()).hexdigest() != ARCHIVE_SHA:
        raise ValueError("source archive integrity mismatch")
    member = f"ipatool-{REV}/tools/fairplay-recovery/profiles/macos26a428-account-entry-public.model.gz"
    with tarfile.open(archive) as source:
        stream = source.extractfile(member)
        if stream is None:
            raise ValueError("public model absent")
        encoded = stream.read()
    if hashlib.sha256(encoded).hexdigest() != MODEL_SHA:
        raise ValueError("public model integrity mismatch")
    with tempfile.TemporaryDirectory(prefix="crumb-image-census-") as directory:
        path = args.binary or exact.exact_system_binary(directory)
        if path.is_symlink() or not path.is_file() or path.stat().st_size > 128 << 20:
            raise ValueError("invalid reference image")
        image = path.read_bytes()
        if hashlib.sha256(image).hexdigest() != exact.BINARY_HASH:
            raise ValueError("exact reference binary required")
        result = census(image, gzip.decompress(encoded))
    result["binary_sha256"] = exact.BINARY_HASH
    result["model_sha256"] = MODEL_SHA
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(result, indent=2)+"\n")


if __name__ == "__main__":
    main()
