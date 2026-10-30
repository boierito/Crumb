"""Independent local-signing tests; ephemeral SELF-SIGNED test identities only.

No Apple identity, private app, SINF, FairPlay key or installed-success claim.
The fixture verifies CMS math and payload preservation, not iOS trust policy.
"""
import argparse
import ctypes
import datetime
import hashlib
from pathlib import Path
import plistlib
import struct
import subprocess
import tempfile


def command(*args):
    subprocess.run(args, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--library', type=Path, required=True)
    parser.add_argument('--executable', type=Path, required=True)
    args = parser.parse_args()
    library = ctypes.CDLL(str(args.library.resolve()))
    sign = library.CrumbSignPreparedBundle
    sign.argtypes = [ctypes.c_char_p]*6
    sign.restype = ctypes.c_int
    validate = library.CrumbValidateSigningIdentity
    validate.argtypes = [ctypes.c_char_p]*5
    validate.restype = ctypes.c_int
    with tempfile.TemporaryDirectory() as tmp:
        root = Path(tmp)
        command('openssl', 'req', '-newkey', 'rsa:2048', '-x509', '-days', '1', '-nodes',
                '-subj', '/CN=Crumb Fixture CA', '-keyout', str(root/'ca-key.pem'), '-out', str(root/'ca.pem'))
        command('openssl', 'req', '-newkey', 'rsa:2048', '-nodes', '-subj', '/CN=Crumb Fixture/OU=FIXTURETEAM',
                '-keyout', str(root/'key.pem'), '-out', str(root/'request.csr'))
        command('openssl', 'x509', '-req', '-in', str(root/'request.csr'), '-CA', str(root/'ca.pem'),
                '-CAkey', str(root/'ca-key.pem'), '-CAcreateserial', '-days', '1', '-out', str(root/'cert.pem'))
        command('openssl', 'pkcs12', '-export', '-inkey', str(root/'key.pem'), '-in', str(root/'cert.pem'),
                '-certfile', str(root/'ca.pem'), '-passout', 'pass:fixture-password', '-out', str(root/'identity.p12'))
        command('openssl', 'x509', '-in', str(root/'cert.pem'), '-outform', 'DER', '-out', str(root/'certificate.der'))
        profile = {'TeamIdentifier': ['FIXTURETEAM'],
                   'DeveloperCertificates': [(root/'certificate.der').read_bytes()],
                   'ExpirationDate': datetime.datetime.now(datetime.timezone.utc).replace(tzinfo=None)+datetime.timedelta(hours=1),
                   'Entitlements': {'application-identifier': 'FIXTURETEAM.com.certlium.crumb',
                                    'com.apple.developer.team-identifier': 'FIXTURETEAM', 'get-task-allow': True}}
        (root/'profile.plist').write_bytes(plistlib.dumps(profile))
        command('openssl', 'cms', '-sign', '-binary', '-nodetach', '-in', str(root/'profile.plist'),
                '-signer', str(root/'cert.pem'), '-inkey', str(root/'key.pem'), '-outform', 'DER', '-out', str(root/'profile.mobileprovision'))
        app = root/'Crumb.app'
        app.mkdir()
        (app/'Info.plist').write_bytes(plistlib.dumps({'CFBundleIdentifier': 'com.certlium.crumb',
            'CFBundleExecutable': 'Crumb', 'CFBundleName': 'Fixture', 'CFBundleVersion': '1', 'CFBundleShortVersionString': '1'}))
        executable = app/'Crumb';executable.write_bytes(args.executable.read_bytes())
        target = app/'CrumbSystemPayload/Target/PlugIns/Test.appex';target.mkdir(parents=True)
        (target/'Info.plist').write_bytes(plistlib.dumps({'CFBundleIdentifier':'original.app.extension','CFBundleExecutable':'Test'}))
        # Recognizable but deliberately incomplete Mach-O: signing it would fail.
        (target/'Test').write_bytes(bytes.fromhex('cffaedfe')+bytes(4092))
        (target/'SC_Info').mkdir();(target/'SC_Info/Test.sinf').write_bytes(b'synthetic license marker; no FairPlay material')
        snapshot = {p.relative_to(app).as_posix():p.read_bytes() for p in (app/'CrumbSystemPayload').rglob('*') if p.is_file()}
        arguments = [str(app).encode(), str(root/'identity.p12').encode(), str(root/'profile.mobileprovision').encode(),
                     b'fixture-password', b'com.certlium.crumb', b'FIXTURETEAM']
        # Preflight errors stay distinct and never mutate the installed bundle.
        checks = arguments[1:]
        assert validate(*checks) == 0
        assert validate(*[checks[0],checks[1],b'bad-password',checks[3],checks[4]]) == 11
        assert validate(*[checks[0],checks[1],checks[2],b'other.bundle',checks[4]]) == 18
        assert validate(*[checks[0],checks[1],checks[2],checks[3],b'OTHERTEAM']) == 17
        # Common sideloader export: RC2-40 + SHA1 PKCS#12, needs legacy provider.
        command('openssl','pkcs12','-export','-legacy','-inkey',str(root/'key.pem'),'-in',str(root/'cert.pem'),
                '-passout','pass:fixture-password','-out',str(root/'legacy.p12'))
        assert validate(str(root/'legacy.p12').encode(),*checks[1:]) == 0
        # Prove profile CMS integrity rather than merely extracting its plist.
        damaged = bytearray((root/'profile.mobileprovision').read_bytes()); damaged[-1] ^= 1
        (root/'damaged.mobileprovision').write_bytes(damaged)
        assert validate(checks[0],str(root/'damaged.mobileprovision').encode(),*checks[2:]) == 14
        for label,change,expected in [
            ('expired',{'ExpirationDate':datetime.datetime(2000,1,1)},15),
            ('membership',{'DeveloperCertificates':[]},16),
            ('scoped-wildcard',{'Entitlements':{'application-identifier':'FIXTURETEAM.com.certlium.*'}},0),
            ('wrong-wildcard',{'Entitlements':{'application-identifier':'FIXTURETEAM.com.other.*'}},18)]:
            altered = dict(profile);altered.update(change)
            (root/'altered.plist').write_bytes(plistlib.dumps(altered))
            command('openssl','cms','-sign','-binary','-nodetach','-in',str(root/'altered.plist'),
                    '-signer',str(root/'cert.pem'),'-inkey',str(root/'key.pem'),'-outform','DER','-out',str(root/'altered.mobileprovision'))
            assert validate(checks[0],str(root/'altered.mobileprovision').encode(),*checks[2:]) == expected,label
        before = hashlib.sha256(executable.read_bytes()).digest()
        assert sign(*[None]*6) == 1
        assert sign(*[arguments[0],arguments[1],arguments[2],b'wrong-password',arguments[4],arguments[5]]) != 0
        assert sign(*[arguments[0],arguments[1],arguments[2],arguments[3],b'wrong.bundle',arguments[5]]) != 0
        assert sign(*[arguments[0],arguments[1],arguments[2],arguments[3],arguments[4],b'OTHERTEAM']) != 0
        assert hashlib.sha256(executable.read_bytes()).digest() == before
        status = sign(*arguments)
        assert status == 0, 'local signing failed (code '+str(status)+')'
        assert snapshot == {p.relative_to(app).as_posix():p.read_bytes() for p in (app/'CrumbSystemPayload').rglob('*') if p.is_file()}
        data = executable.read_bytes();position = 32;signature = None
        assert struct.unpack_from('<I',data)[0] == 0xFEEDFACF
        for _ in range(struct.unpack_from('<I',data,16)[0]):
            kind,size = struct.unpack_from('<II',data,position)
            if kind == 0x1D:
                offset,length = struct.unpack_from('<II',data,position+8);signature=data[offset:offset+length]
            position += size
        assert signature is not None
        magic,length,count = struct.unpack_from('>III',signature)
        assert magic == 0xFADE0CC0
        directories = [];cms = None
        for i in range(count):
            slot,offset = struct.unpack_from('>II',signature,12+i*8)
            magic,length = struct.unpack_from('>II',signature,offset)
            blob = signature[offset:offset+length]
            if magic == 0xFADE0C02:
                directories.append(blob)
                hashes,_,_,slots,limit = struct.unpack_from('>IIIII',blob,16)
                hashsize,kind,_,power = struct.unpack_from('4B',blob,36)
                fn = {1:hashlib.sha1,2:hashlib.sha256}[kind]
                for n in range(slots):
                    assert fn(data[n*(1<<power):min((n+1)*(1<<power),limit)]).digest()[:hashsize] == blob[hashes+n*hashsize:hashes+(n+1)*hashsize]
            if slot == 0x10000: cms=blob[8:]
        assert cms is not None and directories
        (root/'cms.der').write_bytes(cms);(root/'code-directory.bin').write_bytes(directories[0])
        # -noverify skips certificate TRUST only; signature/digest verification
        # remains enabled. This test identity is intentionally not Apple-trusted.
        command('openssl', 'cms', '-verify', '-binary', '-inform', 'DER', '-in', str(root/'cms.der'),
                '-content', str(root/'code-directory.bin'), '-noverify', '-out', str(root/'verified.bin'))
        print('Local CMS signature, code-page hashes, invalid credentials and exact purchased-payload preservation passed; iOS trust/installation untested.')


if __name__ == '__main__':
    main()
