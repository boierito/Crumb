#!/usr/bin/env python3
"""Verify the public native caller that defines the producer diversity ABI.
Input must be the unmodified StoreAgent asset already verified by ipatool.
No private state, account identifiers or credential inputs are accepted.
"""
import hashlib
import json
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
}
for offset, expected in checks.items():
    value = bytes.fromhex(expected)
    if blob[offset:offset+len(value)] != value:
        raise SystemExit(f'ABI contract mismatch at public text offset {offset:#x}')
print(json.dumps(dict(profile='storeagent-10.9', query='needs-account-diversity',
    entry='text+0x11d490', arguments=['context-u32','hardware-info-pointer','dsid-u64'],
    result='native-nonzero-means-diversity', sinf_validation=False)))
