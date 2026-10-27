#!/usr/bin/env python3
"""Verify the complete public image/snapshot memory contract in isolation."""
import hashlib
import pathlib
import re
import subprocess
import tarfile
import tempfile

from fairplay_image_resource import read_image_text

ROOT = pathlib.Path(__file__).resolve().parent.parent
REV = 'f7242cce58077a4af8db738ed295ee23ad8f067d'
ARCHIVE_SHA = 'e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39'


def digest_literal(value):
    return '[32]byte{' + ','.join('0x'+value[i:i+2] for i in range(0, 64, 2)) + '}'


encoded, pin = read_image_text(ROOT)
archive = ROOT/'third_party/archives'/f'ipatool-fairplay-{REV}.tar.gz'
if hashlib.sha256(archive.read_bytes()).hexdigest() != ARCHIVE_SHA:
    raise ValueError('source archive integrity mismatch')
with tempfile.TemporaryDirectory(prefix='crumb-image-contract-') as directory:
    work = pathlib.Path(directory)
    with tarfile.open(archive) as source:
        prefix = f'ipatool-{REV}/tools/fairplay-recovery/'
        for name in ('oracle.go', 'generated_26b118.go', 'input_size_default.go',
                     'embedded_model_raw_source.go'):
            text = source.extractfile(prefix+name).read().decode()
            text = re.sub(r'^//go:build[^\n]*\n', '', text)
            (work/name).write_text(text.replace('package main', 'package rawsource', 1))
        (work/'profiles').mkdir()
        model = source.extractfile(prefix+'profiles/macos26a428-account-entry-public.model.gz').read()
        if hashlib.sha256(model).hexdigest() != pin['model_sha256']:
            raise ValueError('public snapshot integrity mismatch')
        (work/'profiles/macos26a428-account-entry-public.model.gz').write_bytes(model)
        runtime = source.extractfile(prefix+'raw_source_step.go').read().decode()
        definition = re.search(r'var rawSourceBinaryHash = \[32\]byte\{[^}]+\}', runtime).group(0)
    subprocess.run(['python3', str(ROOT/'scripts/prepare-fairplay-generic.py'),
                    'runtime', str(work)], check=True)
    for source_name, target_name in (('fairplayrawimage.go', 'image.go'),
                                     ('fairplayrawimage_test.go', 'image_test.go')):
        (work/target_name).write_bytes((ROOT/'MapleSyrup/NativeSAP'/source_name).read_bytes())
    (work/'profiles/text.bin.gz').write_bytes(encoded)
    image = pin['immutable_text']
    (work/'fixture.go').write_text('package rawsource\nimport _ "embed"\n'+definition+'''
type rawSourceRuntime struct { stopPC uint64 }
//go:embed profiles/text.bin.gz
var imageTextCompressed []byte
'''+f'const imageTextAddress uint64={image["address"]}\n'
       f'const imageTextSize uint64={image["bytes"]}\n'
       f'var imageTextCompressedSHA={digest_literal(image["compressed_sha256"])}\n'
       f'var imageTextRawSHA={digest_literal(image["sha256"])}\n')
    (work/'go.mod').write_text('module crumbimagecontract\n\ngo 1.26.0\n')
    subprocess.run(['go', 'test', '-race', '-count=1', '-v', './...'],
                   cwd=work, check=True, timeout=180)
