#!/usr/bin/env python3
import os, pathlib, shutil, subprocess, tempfile
root=pathlib.Path(__file__).resolve().parents[1]
cc=os.environ.get('CC') or shutil.which('clang') or shutil.which('cc')
with tempfile.TemporaryDirectory(prefix='crumb-ls-') as directory:
    binary=str(pathlib.Path(directory)/'ls-tests')
    subprocess.run([cc,'-std=gnu11','-Wall','-Wextra','-Werror','-Wno-misleading-indentation',
                    '-fsanitize=address,undefined','-g','-fno-omit-frame-pointer',
                    str(root/'tests/CrumbLSStoreTests.c'),str(root/'SharedDecrypt/CrumbLSStore.c'),'-o',binary],check=True)
    subprocess.run([binary],check=True)
