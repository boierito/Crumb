#!/usr/bin/env python3
"""Portable native-reader fault tests. Synthetic callbacks never prove FairPlay."""
import os, pathlib, shutil, subprocess, tempfile
root = pathlib.Path(__file__).resolve().parents[1]
cc = os.environ.get('CC') or shutil.which('clang') or shutil.which('cc')
if not cc:
    raise SystemExit('A C compiler is required')
with tempfile.TemporaryDirectory(prefix='crumb-reader-') as directory:
    binary = str(pathlib.Path(directory) / 'reader-tests')
    subprocess.run([cc, '-std=gnu11', '-Wall', '-Wextra', '-Werror', '-Wno-misleading-indentation', '-g',
                    '-fsanitize=address,undefined', '-fno-omit-frame-pointer',
                    str(root/'tests/CrumbDecryptTests.c'), str(root/'SharedDecrypt/CrumbDecrypt.c'),
                    str(root/'DecryptHelper/CrumbMachO.c'), '-lm', '-o', binary], check=True)
    subprocess.run([binary], check=True)
