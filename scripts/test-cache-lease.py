#!/usr/bin/env python3
"""Real filesystem transactions; no sandbox extension or FairPlay simulated."""
import os, pathlib, shutil, subprocess, tempfile
root = pathlib.Path(__file__).resolve().parents[1]
cc = os.environ.get('CC') or shutil.which('clang') or shutil.which('cc')
if not cc:
    raise SystemExit('A C compiler is required')
with tempfile.TemporaryDirectory(prefix='crumb-lease-') as directory:
    binary = str(pathlib.Path(directory) / 'lease-tests')
    subprocess.run([cc, '-std=gnu11', '-Wall', '-Wextra', '-Werror', '-g',
                    '-fsanitize=address,undefined', '-fno-omit-frame-pointer',
                    str(root/'tests/CrumbCacheLeaseTests.c'), str(root/'SharedDecrypt/CrumbCacheLease.c'),
                    '-o', binary], check=True)
    subprocess.run([binary], check=True)
