#!/usr/bin/env python3
"""Generate PUBLIC synthetic vectors with the independent, hash-pinned Python core.
Usage: python3 scripts/generate-mobile-license-vectors.py /path/to/fairplay-decrypt/core.py
Requires cryptography. Never takes an IPA, capture or real credential as input.
"""
import hashlib
import importlib.util
import json
from pathlib import Path
import sys
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

source = Path(sys.argv[1])
if hashlib.sha256(source.read_bytes()).hexdigest() != '9150ee023ae46846fdb9b929fc2da0c2f15488f299f8072aacd0d074dd9a21c4':
    raise SystemExit('reference core hash mismatch')
spec = importlib.util.spec_from_file_location('mobile_reference', source)
core = importlib.util.module_from_spec(spec)
sys.modules[spec.name] = core
spec.loader.exec_module(core)


def atom(tag, payload):
    return (8 + len(payload)).to_bytes(4, 'big') + tag + payload


def aes_encrypt(key, iv, data):
    op = Cipher(algorithms.AES(key), modes.CBC(iv)).encryptor()
    return op.update(data) + op.finalize()


vectors = []
for index, page in enumerate((0, 1, 17, 0xffff, 0x10000, 0xffffffff)):
    for length in (440, 448):
        credential = hashlib.sha256(f'PUBLIC credential {index}'.encode()).digest()[:16]
        inner = hashlib.sha256(f'PUBLIC inner {index}'.encode()).digest()[:16]
        iv = hashlib.sha256(f'PUBLIC iv {index}'.encode()).digest()[:16]
        name = b'Public fixture\0ignored suffix'
        # Only the header fields checked by upstream are asserted by this vector.
        # No claim is made that these synthetic private records are real licenses.
        clear = bytearray(length if length == 448 else 432)
        clear[:4] = b'itun'
        clear[20:24] = b'key '
        clear[24:40] = inner
        clear[44:48] = b'iviv'
        encrypted = aes_encrypt(bytes(x ^ 0x30 for x in credential), hashlib.md5(name.split(b'\0')[0] + iv).digest(), clear)
        if length == 440:
            encrypted += bytes(8)
        sinf = atom(b'sinf', atom(b'schi', atom(b'name', name) + atom(b'iviv', iv) + atom(b'priv', encrypted)))
        work = core.derive_work16(credential, sinf)
        record = hashlib.sha256(f'PUBLIC record {index}'.encode()).digest()
        key = core.derive_page_key(record[:16], page, work)
        plain = bytes((i * 37 + index) % 256 for i in range(4096))
        ciphertext = aes_encrypt(key, record[16:], plain)
        vectors.append(dict(sinf=sinf.hex(), credential=credential.hex(), work=work.hex(), record=record.hex(), page=page, key=key.hex(), ciphertext=ciphertext.hex(), plaintext_sha256=hashlib.sha256(plain).hexdigest()))
output = Path(__file__).resolve().parents[1] / 'MapleSyrup/NativeSAP/mobilelicense/testdata/reference.json'
output.parent.mkdir(exist_ok=True)
output.write_text(json.dumps(vectors, indent=2) + '\n')
