"""Read the immutable public semantic profile; never account-derived state."""
import hashlib,json,pathlib

def read_profile(root):
    folder=pathlib.Path(root)/'third_party/unfairplay-generic'
    pin=json.loads((folder/'PROFILE.json').read_text())
    if pin['schema']!='crumb-generic-profile-pin-v1':raise ValueError('generic profile pin schema mismatch')
    parts=[]
    for index,part in enumerate(pin['parts']):
        if part['path']!=f'profile/part-{index:03d}':raise ValueError('invalid generic profile part path')
        path=folder/part['path']
        if path.is_symlink() or not path.is_file():raise ValueError('invalid generic profile part')
        data=path.read_bytes()
        if len(data)!=part['bytes'] or hashlib.sha256(data).hexdigest()!=part['sha256']:raise ValueError('generic profile part integrity mismatch')
        parts.append(data)
    compressed=b''.join(parts)
    if hashlib.sha256(compressed).hexdigest()!=pin['compressed_sha256']:raise ValueError('generic profile integrity mismatch')
    return compressed,pin
