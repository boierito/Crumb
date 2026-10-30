"""Embed an owned IPA in the unsigned reference helper, locally only.

Never run with private inputs in CI. Output contains the original app's SINF.
No credentials, captured state or account password are read by this tool.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import plistlib
import stat
import zipfile


def safe_path(name):
    parts = name.rstrip('/').split('/')
    if not parts or '\\' in name or any(p in ('', '.', '..') for p in parts):
        raise ValueError('unsafe archive path')
    return PurePosixPath(*parts)


def index(archive):
    entries, total = {}, 0
    for item in archive.infolist():
        path = safe_path(item.filename)
        if item.is_dir():
            continue
        mode = item.external_attr >> 16
        if stat.S_ISLNK(mode) or (stat.S_IFMT(mode) not in (0, stat.S_IFREG)):
            raise ValueError('non-regular archive member')
        if path in entries:
            raise ValueError('duplicate archive member')
        total += item.file_size
        if item.file_size > 1 << 30 or total > 4 << 30 or len(entries) >= 100000:
            raise ValueError('archive resource limit')
        entries[path] = item
    return entries


def prepare(helper_path, original_path, output, stamp):
    os.umask(0o077)
    if output.exists():
        raise ValueError('refusing to overwrite an existing output')
    try:
        with zipfile.ZipFile(helper_path) as helper, zipfile.ZipFile(original_path) as original:
            helpers, sources = index(helper), index(original)
            hp = PurePosixPath('Payload/FairPlayed.app')
            info = plistlib.loads(helper.read(str(hp/'Info.plist')))
            if info.get('CFBundleIdentifier') != 'com.boierito.crumb.systemfairplay':
                raise ValueError('unexpected helper identity')
            if any('EmbeddedPayload' in p.parts for p in helpers):
                raise ValueError('helper already contains a payload')
            apps = {p.parts[1] for p in sources if len(p.parts) > 2 and p.parts[0] == 'Payload' and p.parts[1].endswith('.app')}
            if len(apps) != 1:
                raise ValueError('expected exactly one app')
            app = PurePosixPath('Payload', apps.pop())
            app_info = plistlib.loads(original.read(str(app/'Info.plist')))
            executable = app_info.get('CFBundleExecutable')
            if not isinstance(executable, str) or len(safe_path(executable).parts) != 1 or app/executable not in sources:
                raise ValueError('invalid app executable')
            manifest_path = app/'SC_Info/Manifest.plist'
            if manifest_path not in sources:
                raise ValueError('original license manifest missing')
            manifest = plistlib.loads(original.read(str(manifest_path)))
            sinfs = manifest.get('SinfPaths')
            if not isinstance(sinfs, list) or not sinfs:
                raise ValueError('original SINF path missing')
            main = app/safe_path(sinfs[0])
            if main not in sources:
                raise ValueError('original SINF missing')
            replicas = manifest.get('SinfReplicationPaths', [])
            if not isinstance(replicas, list) or len(replicas) > 4096:
                raise ValueError('invalid SINF replication manifest')
            replication = [safe_path(p) for p in replicas]
            target = hp/'EmbeddedPayload/Target'
            written, copied = set(), 0
            with zipfile.ZipFile(output, 'x', compression=zipfile.ZIP_DEFLATED, compresslevel=6) as dest:
                def put(path, data, mode=0o644):
                    if path in written:
                        raise ValueError('output member collision')
                    entry = zipfile.ZipInfo(str(path))
                    entry.create_system = 3
                    entry.external_attr = (stat.S_IFREG | mode) << 16
                    entry.compress_type = zipfile.ZIP_DEFLATED
                    dest.writestr(entry, data)
                    written.add(path)
                for path, member in helpers.items():
                    if '_CodeSignature' in path.parts or path.name == 'embedded.mobileprovision':
                        continue
                    data = helper.read(member)
                    if path == hp/'Info.plist':
                        info['CFBundleDisplayName'] = 'Crumb FP System Test'
                        info['CFBundleVersion'] = stamp
                        data = plistlib.dumps(info, fmt=plistlib.FMT_BINARY)
                    put(path, data, (member.external_attr >> 16) & 0o777 or 0o644)
                for path, member in sources.items():
                    if not path.is_relative_to(app):
                        continue
                    relative = path.relative_to(app)
                    if relative.parts[0] == 'Watch':
                        continue  # upstream route is iPhone-only; explicitly not a whole watchOS export
                    put(target/relative, original.read(member), (member.external_attr >> 16) & 0o777 or 0o644)
                    copied += 1
                for relative in replication:
                    if relative.parts[0] != 'Watch' and target/relative not in written:
                        put(target/relative, original.read(str(main)))
                put(hp/'CrumbSystemTest.json', json.dumps({
                    'upstream_commit': 'bad7df138839b5e0ac8b9f20cc4689140b8436fb',
                    'utc_build': stamp, 'route': 'installed-bundle-system-fairplay',
                    'license_validated': False, 'plaintext_authenticated': False,
                    'watch_content_included': False}).encode())
        os.chmod(output, 0o600)
        with zipfile.ZipFile(output) as archive:
            if archive.testzip() is not None:
                raise ValueError('output CRC validation failed')
        return {'sha256': hashlib.sha256(output.read_bytes()).hexdigest(), 'bytes': output.stat().st_size,
                'target_files': copied, 'license_validated': False, 'plaintext_authenticated': False}
    except BaseException:
        output.unlink(missing_ok=True)
        raise


if __name__ == '__main__':
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--helper', type=Path, required=True)
    p.add_argument('--original', type=Path, required=True)
    p.add_argument('--output', type=Path, required=True)
    p.add_argument('--utc-build', required=True)
    a = p.parse_args()
    if len(a.utc_build) != 12 or not a.utc_build.isascii() or not a.utc_build.isdigit():
        p.error('UTC build must be YYYYMMDDHHmm')
    print(json.dumps(prepare(a.helper, a.original, a.output, a.utc_build)))
