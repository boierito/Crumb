#!/usr/bin/env python3
"""Public binary evidence for authorization URL policy and response handling.
No account, capture or SINF input is read. Both inputs must be the exact Apple
assets already hash-verified by ipatool's loader.
"""
import hashlib
import json
import struct
import sys
from pathlib import Path


def verify(path, digest, checks):
    data = Path(path).read_bytes()
    if hashlib.sha256(data).hexdigest() != digest:
        raise SystemExit('Authorization profile hash mismatch')
    for offset, expected in checks.items():
        value = bytes.fromhex(expected)
        if data[offset:offset + len(value)] != value:
            raise SystemExit(f'Authorization profile mismatch at {offset:#x}')
    return data


agent = verify(sys.argv[1], '70ce036f9dbcbc04db9511ebd08de0dd3cbc35ccc9d44b089c90170cb5453c59', {
    0xa79f2: '488d15effc1a00',   # Bag.sendGUIDWithURL: looks up guid-urls
    0xa7d92: 'b901000000',       # regex matching is case insensitive
    0xa7da9: '488b3590c91900',   # regex input is NSURL.absoluteString
    0x59e73: '4d85f6',           # test failureType object
    0x59e76: '0f849f000000',     # no failureType skips customerMessage handling
})
commerce = verify(sys.argv[2], 'b84ff12c21987856c0a17b78f1ad82b73195a6dec5f3b208a17d245555a2c8a2', {
    0x1cecc: '488b353d832800',   # _runURLOperation adds standard query parameters
    0x1bdea: '488b353f932800',   # query builder calls sendGUIDWithURL:
    0x1bf11: '488d0d18142900',   # query key is guid
})


def cfstring_at_lea(data, offset, base, expected):
    target = offset + 7 + struct.unpack_from('<i', data, offset + 3)[0]
    pointer = struct.unpack_from('<Q', data, target + 16)[0] - base
    if data[pointer:pointer + len(expected) + 1] != expected + b'\0':
        raise SystemExit('Authorization protocol literal mismatch')


cfstring_at_lea(agent, 0xa79f2, 0x100000000, b'guid-urls')
cfstring_at_lea(commerce, 0x1bf11, 0, b'guid')
print(json.dumps(dict(bag_guid_policy_verified=True, guid_in_query_distinct_from_plist=True,
    customer_message_alone_is_not_native_failure=True, real_server_response_verified=False,
    sinf_validated=False)))
