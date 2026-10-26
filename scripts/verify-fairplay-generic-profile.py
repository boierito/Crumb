#!/usr/bin/env python3
"""Verify a whole-section generic profile and its integration, without PC handlers."""
from collections import Counter
import gzip,hashlib,json,pathlib,struct,sys
from fairplay_generic_resource import read_profile
ROOT=pathlib.Path(__file__).resolve().parent.parent

def verify(root):
    folder=root/'third_party/unfairplay-generic'
    compressed,pin=read_profile(root)
    data=gzip.decompress(compressed)
    if len(data)>128<<20 or len(data)<128:raise ValueError('generic profile size invalid')
    payload,checksum=data[:-32],data[-32:]
    if hashlib.sha256(payload).digest()!=checksum or checksum.hex()!=pin['profile_payload_sha256']:raise ValueError('generic profile payload integrity mismatch')
    if payload[:8]!=b'CRIR1\0\0\0' or struct.unpack_from('<II',payload,8)!=(2,16):raise ValueError('generic profile schema invalid')
    kernel=(folder/'compact_kernel.go').read_text();cfg=(folder/'compact_cfg.go').read_text()
    abi=hashlib.sha256((kernel+cfg+'crumb-generic-semantics-v2').encode()).hexdigest()
    if payload[16:48].hex()!=pin['binary_sha256'] or payload[48:80].hex()!=abi or pin['kernel_abi']!=abi:raise ValueError('generic source/kernel ABI mismatch')
    base,count,reserved=struct.unpack_from('<QII',payload,80)
    if reserved or base%4 or count<1 or len(payload)!=96+16*count or base!=pin['text_base'] or count!=pin['span_slots']:raise ValueError('generic executable-section geometry invalid')
    kinds=Counter(payload[96::16])
    if set(kinds)-set(range(7)):raise ValueError('generic record type invalid')
    if count<=kinds[6]:raise ValueError('profile has no executable slots')
    if count-kinds[6]!=pin['instruction_slots']:raise ValueError('executable slot census mismatch')
    if kinds[0]!=pin['unsupported_slots'] or count-kinds[0]-kinds[6]!=pin['supported_slots']:raise ValueError('generic coverage census mismatch')
    census=json.loads((folder/'CENSUS.json').read_text())
    if any(census[key]!=pin[key] for key in ['binary_sha256','kernel_abi','instruction_slots','supported_slots','unsupported_slots','profile_payload_sha256']):raise ValueError('census is not for the embedded profile')
    dispatcher=(root/'MapleSyrup/NativeSAP/fairplayrawblocks.go').read_text()
    preparation=(root/'scripts/prepare-fairplay-raw.py').read_text()
    if 'genericImage.contains(pc)' not in dispatcher or 'runGenericInstructionBlock(s, r, pc)' not in dispatcher:raise ValueError('generic fallback not integrated')
    if "generic_embed.go" not in preparation or 'read_profile(root)' not in preparation:raise ValueError('generic resource not included in native package')
    return {'schema':'crumb-generic-release-v1','profile_present':True,'instruction_slots':count-kinds[6],'span_slots':count,'supported_slots':count-kinds[0]-kinds[6],'unsupported_slots':kinds[0],'supported_percent':round(100*(count-kinds[0]-kinds[6])/(count-kinds[6]),6),'unknown_address_allowlist_required':False,'jit_required':False,'ready_for_device_test':True,'decryption_validated':False}

if __name__=='__main__':
    try:print(json.dumps(verify(ROOT),indent=2))
    except (ValueError,KeyError,OSError) as error:
        print(json.dumps({'ready_for_device_test':False,'decryption_validated':False,'reason':str(error)},indent=2));sys.exit(1)
