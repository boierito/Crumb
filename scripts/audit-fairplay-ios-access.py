#!/usr/bin/env python3
"""Replay the *absent-entitlement* rejection branch of Apple's public iOS driver.

This verifies native control flow, not an iOS service contract or decryption.
No branch is patched. Only three entitlement lookups and a logger are mocked.
The driver binary is a local research input, never included in Crumb or CI.
"""
import argparse
import hashlib
import json
import struct
from pathlib import Path

PIN = '6bae7cc2a7e6808dcceefe1dabe16d97458d6431c0fdbb00e508357a92c482f1'
ENTRY = 0xfffffff009bdf890
LOOKUP = 0xfffffff009e1786c
LOGGER = 0xfffffff009e177ac
EXPECTED = ['com.apple.private.FairPlayIOKitUserClient.access',
            'com.apple.private.ProvInfoIOKitUserClient.access',
            'com.apple.private.FairPlayIOKitUserClient.Virtual.access']


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--driver', type=Path, required=True)
    args = parser.parse_args()
    data = args.driver.read_bytes()
    if hashlib.sha256(data).hexdigest() != PIN:
        raise SystemExit('Wrong binary; expected extracted iPhone17,1 / 24A446 FairPlayIOKit')
    from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM, UC_HOOK_CODE, UC_PROT_ALL
    from unicorn.arm64_const import UC_ARM64_REG_X0, UC_ARM64_REG_X1, UC_ARM64_REG_LR, UC_ARM64_REG_PC, UC_ARM64_REG_SP
    uc = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
    assert struct.unpack_from('<I', data)[0] == 0xfeedfacf
    pos = 32
    for _ in range(struct.unpack_from('<I', data, 16)[0]):
        kind, size = struct.unpack_from('<II', data, pos)
        assert size >= 8
        if kind == 0x19:
            address, length, offset, filelength = struct.unpack_from('<4Q', data, pos+24)
            start, end = address & ~4095, (address+length+4095) & ~4095
            if length:
                uc.mem_map(start, end-start, UC_PROT_ALL)
                uc.mem_write(address, data[offset:offset+filelength])
        pos += size
    stack, stop = 0x70000000, 0x10000
    uc.mem_map(stack, 0x10000, UC_PROT_ALL)
    uc.mem_map(stop, 4096, UC_PROT_ALL)
    uc.reg_write(UC_ARM64_REG_SP, stack+0x8000)
    uc.reg_write(UC_ARM64_REG_LR, stop)
    calls, statistics = [], {'instructions': 0, 'log_calls': 0}

    def cstring(address):
        result = bytearray()
        for i in range(256):
            byte = bytes(uc.mem_read(address+i, 1))
            if byte == b'\0':
                return result.decode('ascii')
            result += byte
        raise AssertionError('unterminated entitlement')

    def hook(engine, pc, size, _):
        statistics['instructions'] += 1
        if pc == LOOKUP:
            # Explicit model: absent entitlements return NULL. This is not a
            # measurement of an Apple service and does not grant any privilege.
            calls.append(cstring(engine.reg_read(UC_ARM64_REG_X1)))
            engine.reg_write(UC_ARM64_REG_X0, 0)
            engine.reg_write(UC_ARM64_REG_PC, engine.reg_read(UC_ARM64_REG_LR))
        elif pc == LOGGER:
            statistics['log_calls'] += 1
            assert 'entitlement not found' in cstring(engine.reg_read(UC_ARM64_REG_X0))
            engine.reg_write(UC_ARM64_REG_X0, 0)
            engine.reg_write(UC_ARM64_REG_PC, engine.reg_read(UC_ARM64_REG_LR))
        elif pc == stop:
            engine.emu_stop()
        else:
            word = struct.unpack('<I', bytes(engine.mem_read(pc, 4)))[0]
            if word in (0xd503237f, 0xd50323ff):  # PACIBSP / AUTIBSP
                engine.reg_write(UC_ARM64_REG_PC, pc+4)
            elif word == 0xd65f0fff:  # RETAB; authentication not modeled
                engine.reg_write(UC_ARM64_REG_PC, engine.reg_read(UC_ARM64_REG_LR))

    uc.hook_add(UC_HOOK_CODE, hook)
    uc.emu_start(ENTRY, stop+4, count=1000)
    assert calls == EXPECTED
    assert uc.reg_read(UC_ARM64_REG_PC) == stop
    assert uc.reg_read(UC_ARM64_REG_X0) == 0
    assert statistics['log_calls'] == 1
    print(json.dumps({'build': '24A446', 'device': 'iPhone17,1', 'sha256': PIN,
                      'entry': hex(ENTRY), 'modeled_entitlements': 'absent',
                      'queries': calls, 'initialization_accepted': False,
                      'branch_patched': False, 'service_results_measured': False,
                      'decryption_verified': False, **statistics}, indent=2))


if __name__ == '__main__':
    main()
