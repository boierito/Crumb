#!/usr/bin/env python3
"""Replay Core ML's identity-choice helper from pinned original Apple services.

Only control flow is checked. Objective-C calls and ARC are explicitly modeled;
no system service, signing identity, model key or App Store license is exercised.
The service binaries must be acquired separately from the public Apple IPSWs.
"""
import argparse
import hashlib
import json
import plistlib
import struct
from pathlib import Path


BUILDS = {
    '8820f7af22ea0742f108a358bbb709d5ec856d037008121d1af1807e2f757193': {
        'build': '22C152', 'os': '18.2', 'entry': 0x10000a24c,
        'signing': 0x10000bb20, 'team': 0x10000c660,
        'claim': 0x10000b750, 'autorelease': 0x10000b710,
        'return': 0x10000b740, 'error': 0x10000c040,
        'log': 0x10000bbe0, 'enabled': 0x10000b960,
        'release': 0x10000b7c0,
        'producer': (0x100009f30, 0x10000a24c),
        'cloud_request': (0x100007f24, 0x100008014),
    },
    '07b277507b80a27b3aeb5c954721f3223b657cd3c5b551e45fbb2f8d4afb4e5a': {
        'build': '24A446', 'os': '27.0.1', 'entry': 0x100006c14,
        'signing': 0x1000082c0, 'team': 0x100008e80,
        'claim': 0x100007eb4, 'autorelease': 0x100007e74,
        'return': 0x100007ea4, 'error': 0x100008820,
        'log': 0x100008380, 'enabled': 0x1000080c4,
        'release': 0x100007f24,
        'producer': (0x100006908, 0x100006c14),
        'cloud_request': (0x100004940, 0x100004a30),
    },
}


def parse_image(data):
    if len(data) < 32 or struct.unpack_from('<I', data)[0] != 0xfeedfacf:
        raise ValueError('Expected a thin ARM64 Mach-O')
    if struct.unpack_from('<I', data, 4)[0] != 0x100000c:
        raise ValueError('Expected ARM64')
    commands, size = struct.unpack_from('<II', data, 16)
    if size > len(data)-32:
        raise ValueError('Load-command bounds')
    pos, segments, signature = 32, [], None
    for _ in range(commands):
        if pos+8 > 32+size:
            raise ValueError('Truncated load command')
        kind, length = struct.unpack_from('<II', data, pos)
        if length < 8 or pos+length > 32+size:
            raise ValueError('Invalid load-command size')
        if kind == 0x19:
            if length < 72:
                raise ValueError('Invalid segment command')
            name = data[pos+8:pos+24].rstrip(b'\0').decode('ascii')
            va, memory, off, file_size = struct.unpack_from('<4Q', data, pos+24)
            if off > len(data) or file_size > len(data)-off:
                raise ValueError('Segment outside file')
            if name != '__PAGEZERO':
                segments.append((va, memory, off, file_size))
        elif kind == 0x1d:
            signature = struct.unpack_from('<II', data, pos+8)
        pos += length
    if pos != 32+size:
        raise ValueError('Inconsistent command size')
    return segments, signature


def read_vm(data, segments, address, size):
    for va, _, off, file_size in segments:
        if va <= address and size <= file_size-(address-va):
            return data[off+address-va:off+address-va+size]
    raise ValueError(f'Unbacked VM range {address:#x}+{size:#x}')


def selector(data, segments, stub):
    # The pinned Objective-C selector stub is ADRP x1; LDR x1,[x1,#imm].
    adrp, ldr = struct.unpack('<II', read_vm(data, segments, stub, 8))
    if adrp & 0x9f00001f != 0x90000001 or ldr & 0xffc003ff != 0xf9400021:
        raise ValueError('Unexpected Objective-C stub instructions')
    imm = ((adrp >> 5) & 0x7ffff) << 2 | ((adrp >> 29) & 3)
    if imm & (1 << 20):
        imm -= 1 << 21
    address = (stub & ~4095)+(imm << 12)+(((ldr >> 10) & 4095) << 3)
    raw, = struct.unpack('<Q', read_vm(data, segments, address, 8))
    # These two pinned images use an unauthenticated, unbound arm64e rebase.
    # The low 43-bit offset excludes the next-chain bits. Do not generalize to
    # arbitrary pointers or authenticated/bound chained-pointer formats.
    if raw >> 62:
        raise ValueError('Unsupported selector rebase format')
    target = 0x100000000+(raw & ((1 << 43)-1))
    result = bytearray()
    for i in range(512):
        b = read_vm(data, segments, target+i, 1)
        if b == b'\0':
            return result.decode('ascii')
        result += b
    raise ValueError('Unterminated selector')


def entitlements(data, signature):
    if signature is None:
        raise ValueError('Missing code signature')
    off, size = signature
    blob = data[off:off+size]
    if len(blob) != size or size < 12:
        raise ValueError('Invalid code signature bounds')
    magic, length, count = struct.unpack_from('>III', blob)
    if magic != 0xfade0cc0 or length > size or count > (length-12)//8:
        raise ValueError('Invalid signature superblob')
    for i in range(count):
        slot, start = struct.unpack_from('>II', blob, 12+i*8)
        if slot != 5:
            continue
        if start > length-8:
            raise ValueError('Invalid entitlement offset')
        tag, n = struct.unpack_from('>II', blob, start)
        if tag != 0xfade7171 or n < 8 or n > length-start:
            raise ValueError('Invalid entitlement blob')
        return plistlib.loads(blob[start+8:start+n])
    raise ValueError('XML entitlement slot missing')


def selector_calls(data, segments, start, end):
    """Resolve only original direct BL calls to the verified selector-stub form.

    This is a static call-site inventory, not proof that all branches execute.
    """
    result = []
    for pc in range(start, end, 4):
        word, = struct.unpack('<I', read_vm(data, segments, pc, 4))
        if word & 0xfc000000 != 0x94000000:
            continue
        imm = word & 0x3ffffff
        if imm & (1 << 25):
            imm -= 1 << 26
        target = pc+imm*4
        try:
            name = selector(data, segments, target)
        except ValueError:
            continue
        result.append({'call_site': hex(pc), 'selector': name})
    return result


def producer_contract(data, segments, build):
    producer = selector_calls(data, segments, *build['producer'])
    expected = [
        'generateKeyRequestForKeyIdentifier:teamIdentifier:error:',
        'fetchKeyResponseFromServerForKeyID:signedKeyRequest:error:',
        'generatePersistentKeyBlobFromKeyResponse:error:',
    ]
    selected = [c for c in producer if c['selector'] in expected]
    if [c['selector'] for c in selected] != expected:
        raise ValueError('Native model-blob producer call chain differs')
    cloud = selector_calls(data, segments, *build['cloud_request'])
    fields = ['setKeyId:', 'setSignedKeyRequest:', 'setTeamId:', 'setRawRequest:']
    selected_fields = [c for c in cloud if c['selector'] in fields]
    if [c['selector'] for c in selected_fields] != fields:
        raise ValueError('Native model-key request fields differ')
    pc = int(selected_fields[-1]['call_site'], 16)
    # Both originals explicitly set w2=0 just before setRawRequest:.
    if read_vm(data, segments, pc-4, 4) != bytes.fromhex('02008052'):
        raise ValueError('Native raw-request flag is not the audited false constant')
    return {'static_producer_call_sites': selected,
            'static_request_field_call_sites': selected_fields,
            'native_raw_request_flag': False,
            'network_request_executed': False}


def replay(data, segments, build, use_signing, identity_present):
    from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM, UC_HOOK_CODE
    from unicorn.arm64_const import (
        UC_ARM64_REG_X0, UC_ARM64_REG_X1, UC_ARM64_REG_X30,
        UC_ARM64_REG_SP, UC_ARM64_REG_PC,
    )
    uc = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
    for va, memory, off, size in segments:
        if memory > 16 << 20 or va % 4096:
            raise ValueError('Unexpected segment layout')
        uc.mem_map(va, (memory+4095) & ~4095)
        if size:
            uc.mem_write(va, data[off:off+size])
    stack, error_address = 0x70000000, 0x70001000
    uc.mem_map(stack, 0x20000)
    uc.reg_write(UC_ARM64_REG_SP, stack+0x18000)
    uc.reg_write(UC_ARM64_REG_X30, 0x71000000)
    uc.reg_write(UC_ARM64_REG_X0, int(use_signing))
    uc.reg_write(UC_ARM64_REG_X1, error_address)
    modeled_identity, modeled_error = 0x60001000, 0x60002000
    selected, calls, skipped_pac, returned = [], [], [], []

    def hook(uc, pc, size, _):
        if pc in (build['signing'], build['team']):
            selected.append('signing_identifier' if pc == build['signing'] else 'team_identifier')
            uc.reg_write(UC_ARM64_REG_X0, modeled_identity if identity_present else 0)
        elif pc in (build['claim'], build['autorelease']):
            pass  # Modeled ARC leaves x0 unchanged.
        elif pc == build['error']:
            uc.reg_write(UC_ARM64_REG_X0, modeled_error)
        elif pc == build['log']:
            uc.reg_write(UC_ARM64_REG_X0, 0x60003000)
        elif pc == build['enabled']:
            uc.reg_write(UC_ARM64_REG_X0, 0)  # Logging is modeled as disabled.
        elif pc == build['release']:
            pass
        elif pc == build['return']:
            returned.append(uc.reg_read(UC_ARM64_REG_X0))
            uc.emu_stop()
            return
        else:
            if not build['entry'] <= pc < build['entry']+0x100:
                raise ValueError(f'Unmodeled external instruction at {pc:#x}')
            instruction = bytes(uc.mem_read(pc, 4))
            if instruction in (bytes.fromhex('7f2303d5'), bytes.fromhex('ff2303d5')):
                skipped_pac.append(hex(pc))
                uc.reg_write(UC_ARM64_REG_PC, pc+4)
            return
        calls.append(hex(pc))
        uc.reg_write(UC_ARM64_REG_PC, uc.reg_read(UC_ARM64_REG_X30))

    uc.hook_add(UC_HOOK_CODE, hook)
    uc.emu_start(build['entry'], 0, count=200)
    expected = 'signing_identifier' if use_signing else 'team_identifier'
    if selected != [expected] or returned != [modeled_identity if identity_present else 0]:
        raise ValueError('Unexpected native identity selection or return')
    error, = struct.unpack('<Q', bytes(uc.mem_read(error_address, 8)))
    if error != (0 if identity_present else modeled_error):
        raise ValueError('Unexpected native error write')
    return {
        'uses_code_signing_identity': use_signing,
        'modeled_identity_present': identity_present,
        'selected': selected[0], 'native_returned_identity': identity_present,
        'native_error_written': bool(error), 'modeled_calls': calls,
        'pac_stack_instructions_modeled': skipped_pac,
    }


def audit(path):
    data = path.read_bytes()
    pin = hashlib.sha256(data).hexdigest()
    if pin not in BUILDS:
        raise ValueError('Service SHA-256 is not one of the two audited Apple images')
    build = BUILDS[pin]
    segments, signature = parse_image(data)
    for name, expected in [('signing', 'codeSigningIdentifierFromXPCConnection'),
                           ('team', 'teamIdentifierFromXPCConnection')]:
        if selector(data, segments, build[name]) != expected:
            raise ValueError('Objective-C selector does not match the modeled contract')
    ent = entitlements(data, signature)
    keys = ['com.apple.private.fpsd.client', 'com.apple.private.coreml.decrypt_allowed',
            'platform-application', 'fairplay-client']
    if any(k not in ent for k in keys):
        raise ValueError('Expected service entitlements missing')
    return {
        'schema': 'crumb-coreml-identity-audit-v1', 'build': build['build'],
        'os': build['os'], 'service_sha256': pin,
        'identity_helper': hex(build['entry']),
        'service_entitlements': {k: ent[k] for k in keys},
        'model_blob_contract': producer_contract(data, segments, build),
        'cases': [replay(data, segments, build, signing, present)
                  for signing in (False, True) for present in (False, True)],
        'evidence_scope': 'original ARM64 control flow with modeled Objective-C/ARC calls',
        'system_service_executed': False, 'license_conversion_verified': False,
        'decryption_verified': False,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--service', type=Path, required=True)
    args = parser.parse_args()
    try:
        result = audit(args.service)
    except (ValueError, OSError) as error:
        parser.exit(1, f'{error}\n')
    print(json.dumps(result, indent=2))


if __name__ == '__main__':
    main()
