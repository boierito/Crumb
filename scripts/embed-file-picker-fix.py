#!/usr/bin/env python3
"""Embed the user-provided picker fix into Crumb before normal signing.
Never relocate executable sections; fail if the linker reserved too little space.
"""
import hashlib
import pathlib
import plistlib
import shutil
import struct
import sys

LOAD_PATH = b"@executable_path/Frameworks/FilePickerFix.dylib\0"
SOURCE = pathlib.Path(__file__).resolve().parents[1] / "EmbeddedLibraries/FilePickerFix.dylib"
SOURCE_SHA256 = hashlib.sha256(SOURCE.read_bytes()).hexdigest()

def inject(data):
    data = bytearray(data)
    if len(data) < 32 or struct.unpack_from('<II',data) != (0xFEEDFACF,0x0100000C):
        raise ValueError('FilePickerFix requires a thin arm64 executable')
    ncmds, size = struct.unpack_from('<II',data,16)
    end = 32 + size
    if end > len(data): raise ValueError('Truncated load commands')
    offset, first_data = 32, len(data)
    for _ in range(ncmds):
        cmd, length = struct.unpack_from('<II',data,offset)
        if length < 8 or offset + length > end: raise ValueError('Invalid load command')
        if cmd == 0xC:
            start = struct.unpack_from('<I',data,offset+8)[0]
            if data[offset+start:offset+length].split(b'\0')[0] == LOAD_PATH[:-1]: return bytes(data)
        if cmd == 0x19:  # LC_SEGMENT_64
            if length < 72: raise ValueError('Invalid segment')
            sections = struct.unpack_from('<I',data,offset+64)[0]
            if 72 + sections*80 > length: raise ValueError('Invalid section table')
            for i in range(sections):
                section = offset+72+i*80
                at = struct.unpack_from('<I',data,section+48)[0]
                if at: first_data = min(first_data,at)
        offset += length
    if offset != end: raise ValueError('Invalid command count')
    length = (24 + len(LOAD_PATH) + 7) & ~7
    if end+length > first_data or end+length > len(data) or any(data[end:end+length]):
        raise ValueError('Insufficient Mach-O header padding for FilePickerFix')
    command = struct.pack('<6I',0xC,length,24,0,0,0) + LOAD_PATH
    data[end:end+length] = command.ljust(length,b'\0')
    struct.pack_into('<II',data,16,ncmds+1,size+length)
    return bytes(data)

if __name__ == '__main__':
    app = pathlib.Path(sys.argv[1])
    info = plistlib.loads((app/'Info.plist').read_bytes())
    binary = app/info['CFBundleExecutable']
    patched = inject(binary.read_bytes())
    directory = app/'Frameworks'
    directory.mkdir(exist_ok=True)
    shutil.copyfile(SOURCE,directory/'FilePickerFix.dylib')
    binary.write_bytes(patched)
    print('Embedded FilePickerFix:',SOURCE_SHA256)
