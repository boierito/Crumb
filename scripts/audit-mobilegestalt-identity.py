#!/usr/bin/env python3
"""Inspect the pinned native hardware-query registry; no hardware query is sent.

The normalized reference comes from public 24A446 firmware. This establishes
registry metadata and original call edges, not the phone's effective permissions
or a successful license/SINF operation.
"""
import argparse
import hashlib
import importlib.util
import json
import struct
from pathlib import Path

PIN = "73e327d8d64410c4542f56717fb3a9d9a2b094b7851621ec44251a937dd9d8b6"


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--image", required=True, type=Path)
    args = p.parse_args()
    data = args.image.read_bytes()
    if hashlib.sha256(data).hexdigest() != PIN:
        raise ValueError("Expected the pinned normalized 24A446 MobileGestalt image")
    spec = importlib.util.spec_from_file_location(
        "image_reader", Path(__file__).with_name("audit-coreml-identity.py")
    )
    reader = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(reader)
    segments, _ = reader.parse_image(data)

    def read(address, size):
        return reader.read_vm(data, segments, address, size)

    def call(pc, target):
        (word,) = struct.unpack("<I", read(pc, 4))
        if word & 0xFC000000 != 0x94000000:
            raise ValueError("Expected original direct BL instruction")
        imm = word & 0x3FFFFFF
        if imm & (1 << 25):
            imm -= 1 << 26
        if pc + (imm << 2) != target:
            raise ValueError("Original call edge differs")
        return {"site": hex(pc), "target": hex(target)}

    row_address = 0x1EC27B560
    row = read(row_address, 64)
    pointers = struct.unpack_from("<7Q", row)
    key = b"re6Zb+zwFKJNlkQTUeT+/w"
    if read(pointers[0], len(key) + 1) != key + b"\0":
        raise ValueError("Registry key differs")
    if pointers[1:4] != (0x1AF9F0D30, 0x1AF9EEA44, 0x1AF9F0A4C):
        raise ValueError("Registry callbacks differ")
    entitlement = b"com.apple.private.MobileGestalt.AllowedProtectedKeys"
    if read(0x1AFA6E942, len(entitlement) + 1) != entitlement + b"\0":
        raise ValueError("Entitlement constant differs")
    edges = [call(0x1AF9EBFBC, 0x1AF9EB588), call(0x1AF9EEA5C, 0x1AF9EEA6C)]
    (word,) = struct.unpack("<I", read(0x1AF9EBFE4, 4))
    imm = word & 0x3FFFFFF
    if imm & (1 << 25):
        imm -= 1 << 26
    if word & 0xFC000000 != 0x14000000 or 0x1AF9EBFE4 + (imm << 2) != 0x1AF9EB848:
        raise ValueError("MGCopyAnswer tail call differs")
    result = {
        "schema": "crumb-mobilegestalt-identity-reference-v1",
        "build": "24A446",
        "normalized_image_sha256": PIN,
        "public_query_name": "UniqueDeviceID",
        "public_registry_key": key.decode(),
        "native_query_name_to_registry_hash_replayed": False,
        "registry_record": hex(row_address),
        "registry_record_sha256": hashlib.sha256(row).hexdigest(),
        "getter": hex(pointers[1]),
        "permission_callback": hex(pointers[2]),
        "permission_report_callback": hex(pointers[3]),
        "record_flags_byte_0x38": row[0x38],
        "record_cache_index_0x3a": struct.unpack_from("<h", row, 0x3A)[0],
        "record_type_0x3c": struct.unpack_from("<I", row, 0x3C)[0],
        "entitlement_constant": entitlement.decode(),
        "entitlement_string_reference": "0x1af9f0acc",
        "original_direct_call_edges": edges,
        "external_sandbox_call_target_independently_resolved": False,
        "full_permission_callback_replayed": False,
        "installed_crumb_entitlements_measured": False,
        "device_query_executed": False,
        "sinf_validated": False,
        "generic_decryption_demonstrated": False,
    }
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
