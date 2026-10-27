#!/usr/bin/env python3
"""Verify the public native caller that defines the producer diversity ABI.
Input must be the unmodified StoreAgent asset already verified by ipatool.
No private state, account identifiers or credential inputs are accepted.
"""
import hashlib
import json
import struct
from pathlib import Path
import sys

blob = Path(sys.argv[1]).read_bytes()
if hashlib.sha256(blob).hexdigest() != '70ce036f9dbcbc04db9511ebd08de0dd3cbc35ccc9d44b089c90170cb5453c59':
    raise SystemExit('StoreAgent profile mismatch')
# This pinned, thin x86_64 image has __TEXT vmaddr 0x100000000 and fileoff 0.
# -[FairPlayHelper needDiversityForDSID:] passes globalContext (edi), the
# FairPlayHWInfo address (rsi), and NSNumber's 64-bit DSID (rdx), then treats
# every nonzero native result as YES. This is not SINF key validation.
checks = {
    0x88e91: '4489f7',                 # mov edi,r14d
    0x88e94: '4889de',                 # mov rsi,rbx
    0x88e97: '4889c2',                 # mov rdx,rax
    0x88e9a: 'e8f1450900',             # call text+0x11d490
    0x88e9f: '85c00f95c0',             # test eax,eax; setne al
    0x11d4a4: '48899570fdffff',         # store 64-bit DSID
    0x11d4ab: '4889b578fdffff',         # store hardware pointer
    0x11d4b2: '89bd6cfdffff',           # store 32-bit context
    0x59888: '31d2',                   # authorization sync reserved argument 0
    0x5988a: 'b902000000',             # authorization KBSync mode 2
    0x5988f: 'e82cfb0600',             # call text+0xc93c0
    0x59bb2: '4885c0',                 # test optional keybag object
    0x59bb5: '0f843d010000',           # nil jumps to diversity field at 0x59cf8
    0x59d1b: '4885c0',                 # test optional diversitybag object
    0x59d1e: '0f842d010000',           # nil jumps to failureType at 0x59e51
    0x59d87: '4489f7',                 # diversity import context u32
    0x59d8a: '4889de',                 # diversity bytes pointer
    0x59d8d: '89c2',                   # diversity length u32
    0x59d8f: 'e80c410700',             # call text+0xcdea0

}
for offset, expected in checks.items():
    value = bytes.fromhex(expected)
    if blob[offset:offset+len(value)] != value:
        raise SystemExit(f'ABI contract mismatch at public text offset {offset:#x}')
# These literal CFStrings belong to the original authorization request and
# response handling. They are static binary metadata, never account values.
for offset, expected in [(0x251fe8,b'authorizeMachine'),(0x252048,b'needDiv'),(0x252068,b'1'),(0x252108,b'diversitybag')]:
    pointer=struct.unpack_from('<Q',blob,offset+16)[0]-0x100000000
    if blob[pointer:pointer+len(expected)+1] != expected+b'\0':
        raise SystemExit('authorization protocol literal mismatch')
print(json.dumps(dict(authorization_sync_mode=2, diversity_import='text+0xcdea0', profile='storeagent-10.9', query='needs-account-diversity',
    entry='text+0x11d490', arguments=['context-u32','hardware-info-pointer','dsid-u64'],
    result='native-nonzero-means-diversity', sinf_validation=False)))
