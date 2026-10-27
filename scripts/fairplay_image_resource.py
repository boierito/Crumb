"""Read only the immutable image artifact proven by its generation census."""
import hashlib
import json


def read_image_text(root):
    base = root / 'third_party/fairplay-immutable-image'
    pin = json.loads((base / 'IMAGE.json').read_text())
    if (pin['binary_sha256'] != '1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784'
            or pin['immutable_text']['address'] != 0x100000000
            or pin['immutable_text']['bytes'] != 16711680
            or pin['immutable_text']['chained_fixups'] is not False):
        raise ValueError('immutable image source or segment contract mismatch')
    pieces = []
    for index, part in enumerate(pin['parts']):
        if part['path'] != f'part-{index:03d}':
            raise ValueError('invalid immutable image part sequence')
        data = (base / part['path']).read_bytes()
        if len(data) != part['bytes'] or hashlib.sha256(data).hexdigest() != part['sha256']:
            raise ValueError('immutable image part integrity mismatch')
        pieces.append(data)
    encoded = b''.join(pieces)
    if len(encoded) != pin['immutable_text']['compressed_bytes'] or hashlib.sha256(encoded).hexdigest() != pin['immutable_text']['compressed_sha256']:
        raise ValueError('immutable image integrity mismatch')
    return encoded, pin
