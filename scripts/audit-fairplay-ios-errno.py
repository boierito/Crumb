#!/usr/bin/env python3
"""Verify iOS 27.0.1 syscall-489 error collapsing using the actual kernel bytes.

The callback return is modeled. This does not measure which callback error the
phone produced, or prove that a legacy report reached mremap_encrypted at all.
"""
import argparse
import hashlib
import json
import struct
from pathlib import Path

PIN = 'fa07057bf34f33af1ec26681ad0c87c13e5aaa548a300740eddb0e92876e5f95'
ADDRESS = 0xfffffff00b0d8ca8
EXPECTED = bytes.fromhex('285880520800bc72df02086b8801805200159f1a')


def mapping_bytes(data):
    assert struct.unpack_from('<I', data)[0] == 0xfeedfacf
    pos = 32
    header = None
    for _ in range(struct.unpack_from('<I', data, 16)[0]):
        kind, size = struct.unpack_from('<II', data, pos)
        if kind == 0x80000035:
            name = struct.unpack_from('<I', data, pos+24)[0]
            if data[pos+name:pos+size].split(b'\0')[0] == b'com.apple.kernel':
                header = struct.unpack_from('<Q', data, pos+16)[0]
        pos += size
    assert header is not None
    pos = header+32
    for _ in range(struct.unpack_from('<I', data, header+16)[0]):
        kind, size = struct.unpack_from('<II', data, pos)
        if kind == 0x19:
            address, _, offset, length = struct.unpack_from('<4Q', data, pos+24)
            if address <= ADDRESS and ADDRESS+len(EXPECTED) <= address+length:
                start = offset+ADDRESS-address
                result = data[start:start+len(EXPECTED)]
                assert result == EXPECTED
                return result
        pos += size
    raise ValueError('mapping sequence absent')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--kernel', type=Path, required=True)
    args = parser.parse_args()
    data = args.kernel.read_bytes()
    if hashlib.sha256(data).hexdigest() != PIN:
        raise SystemExit('Expected the extracted iPhone17,1 / 24A446 kernelcache')
    instructions = mapping_bytes(data)
    from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM
    from unicorn.arm64_const import UC_ARM64_REG_W0, UC_ARM64_REG_W22
    results = []
    for name, returned, expected in [('not_privileged', 0xe00002c1, 1),
                                      ('generic_error', 0xe00002bc, 12),
                                      ('no_memory', 0xe00002bd, 12),
                                      ('not_permitted', 0xe00002e2, 12)]:
        uc = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        uc.mem_map(ADDRESS & ~4095, 4096)
        uc.mem_write(ADDRESS, instructions)
        uc.reg_write(UC_ARM64_REG_W22, returned)
        uc.emu_start(ADDRESS, ADDRESS+len(instructions), count=5)
        observed = uc.reg_read(UC_ARM64_REG_W0)
        assert observed == expected
        results.append({'modeled_callback_result': name, 'value': hex(returned), 'darwin_errno': observed})
    print(json.dumps({'build': '24A446', 'kernel_sha256': PIN,
                      'mapping_address': hex(ADDRESS), 'executed_instructions_per_case': 5,
                      'results': results, 'phone_callback_error_identified': False,
                      'decryption_verified': False}, indent=2))


if __name__ == '__main__':
    main()
