#!/usr/bin/env python3
"""Inspect the actual Foundation ZIP producer using an independent reader."""
import hashlib
import json
from pathlib import Path
import sys
import zipfile

path = Path(sys.argv[1])
manifest = json.loads(path.with_name(path.name + '.json').read_text())
with zipfile.ZipFile(path) as archive:
    assert archive.testzip() is None
    expected = {
        'META-INF/', 'META-INF/empty/', 'META-INF/com.apple.ZipMetadata.plist',
        'iTunesMetadata.plist', 'Payload/Test.app/Test',
    }
    assert set(archive.namelist()) == expected
    assert archive.read('iTunesMetadata.plist') == b'bplist00original-store-fields'
    assert archive.read('META-INF/com.apple.ZipMetadata.plist') == b'original signature bytes'
    assert archive.getinfo('META-INF/empty/').is_dir()
    assert len(manifest['files']) == len(expected)
    for item in manifest['files']:
        entry = archive.getinfo(item['path'])
        assert entry.is_dir() == item.get('directory', False)
        assert hashlib.sha256(archive.read(item['path'])).hexdigest() == item['sha256']
    for name, original in manifest['archive_resources'].items():
        data = archive.read(name)
        assert len(data) == original['size']
        assert hashlib.sha256(data).hexdigest() == original['sha256']
print('Original root layout, metadata bytes, empty directories, hashes and ZIP checksums verified')
