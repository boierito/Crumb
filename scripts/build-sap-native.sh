#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
OUTPUT="$ROOT/build/Native"
mkdir -p "$OUTPUT"
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
# Reuse the current SAP guest, loader and shims only, with an isolated C ABI.
# All downloaded source and generated modifications are retained for review.
curl -fL --retry 3 https://github.com/majd/ipatool/archive/3411d57f451f5111ae115641c22f7ed17bbd5fbe.tar.gz -o "$WORK/ipatool.tar.gz"
python3 - "$WORK/ipatool.tar.gz" <<'PY'
import sys,hashlib
with open(sys.argv[1], 'rb') as f: digest=hashlib.sha256(f.read()).hexdigest()
if digest != '5c940788df3b619b07ead5e592b5a06beecd96ccdb7570f4fbd09a02b41785b7':
    raise SystemExit('ipatool source integrity mismatch')
PY
tar -xzf "$WORK/ipatool.tar.gz" -C "$WORK"
SOURCE="$WORK/ipatool-3411d57f451f5111ae115641c22f7ed17bbd5fbe"
mkdir -p "$SOURCE/wafflebridge"
cp "$ROOT/MapleSyrup/NativeSAP/main.go" "$SOURCE/wafflebridge/main.go"
mkdir -p "$SOURCE/wafflebridge/packageipa"
cp "$ROOT/MapleSyrup/NativeSAP/packageipa/"*.go "$SOURCE/wafflebridge/packageipa/"
python3 - "$SOURCE" <<'PY'
import pathlib,sys
root=pathlib.Path(sys.argv[1])
for name in ['assets.go','storeagent.go']:
    path=root/'internal/sap/assets'/name
    source=path.read_text()
    old='root, err := os.UserCacheDir()\n\tif err != nil {'
    if source.count(old)!=1: raise SystemExit('Unexpected SAP cache layout')
    source=source.replace(old,'root := CacheRoot\n\tvar err error\n\tif root == "" { err = errors.New("explicit sandbox cache path is required") }\n\tif err != nil {')
    if '"errors"' not in source: source=source.replace('import (', 'import (\n\t"errors"', 1)
    # os is still used for actual cache I/O, so the existing import stays.
    if name=='assets.go': source += '\n// Configured once by the sandboxed C bridge before any asset load.\nvar CacheRoot string\n'
    path.write_text(source)
PY
# Portable offline decoder only; preserve the existing SAP guest revision.
FAIRPLAY_REV=e94fcd7362036a1b2269e5b555232f20d8cd273e
curl -fL --retry 3 "https://github.com/DrWhax/ipatool/archive/$FAIRPLAY_REV.tar.gz" -o "$WORK/fairplay.tar.gz"
python3 - "$WORK/fairplay.tar.gz" <<'CHECK'
import hashlib, sys
if hashlib.sha256(open(sys.argv[1], 'rb').read()).hexdigest() != '98f0b7e48bfed8dc1c940208ce84cd5287bd3c689e607de401e38937fc7df085':
    raise SystemExit('FairPlay source integrity mismatch')
CHECK
tar -xzf "$WORK/fairplay.tar.gz" -C "$WORK"
cp -R "$WORK/ipatool-$FAIRPLAY_REV/internal/fairplay" "$SOURCE/internal/fairplay"
# iOS has the same Unix permission checks; treat its filesystem as case insensitive.
python3 - "$SOURCE/internal/fairplay/files.go" <<'PATCH'
import pathlib, sys
p = pathlib.Path(sys.argv[1]); s = p.read_text()
s = s.replace('runtime.GOOS != "darwin" && runtime.GOOS != "linux"', 'runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "ios"')
s = s.replace('runtime.GOOS == "darwin" || runtime.GOOS == "windows"', 'runtime.GOOS == "darwin" || runtime.GOOS == "windows" || runtime.GOOS == "ios"')
p.write_text(s)
PATCH
(cd "$SOURCE" && go get golang.org/x/arch@v0.32.0 && go test ./internal/fairplay)
SDK=$(xcrun --sdk iphoneos --show-sdk-path)
CLANG=$(xcrun --sdk iphoneos --find clang)
cd "$SOURCE"
GOOS=ios GOARCH=arm64 CGO_ENABLED=1 CC="$CLANG" SDKROOT="$SDK" \
  CGO_CFLAGS="-isysroot $SDK -target arm64-apple-ios16.4" \
  CGO_LDFLAGS="-isysroot $SDK -target arm64-apple-ios16.4" \
  go build -buildmode=c-archive -trimpath -o "$OUTPUT/libWaffleSAP.a" ./wafflebridge
go list -deps -json ./wafflebridge > "$WORK/dependencies.json"
python3 - "$WORK/dependencies.json" "$OUTPUT/GoThirdPartyNotices.txt" "$(go env GOROOT)" <<'PY'
import json,pathlib,sys
text=pathlib.Path(sys.argv[1]).read_text(); decoder=json.JSONDecoder(); modules={}
while text.strip():
    value,end=decoder.raw_decode(text.lstrip()); text=text.lstrip()[end:]
    module=value.get('Module')
    if module and module.get('Dir'): modules[module['Path']]=pathlib.Path(module['Dir'])
pieces=['Go standard library\n'+(pathlib.Path(sys.argv[3])/'LICENSE').read_text()]
for name,root in sorted(modules.items()):
    found=[]
    for path in sorted(root.iterdir()):
        if path.is_file() and path.name.upper().startswith(('LICENSE','COPYING','NOTICE')):
            found.append(path.name+'\n'+path.read_text(errors='replace'))
    if not found: raise SystemExit('Missing license notice for linked module: '+name)
    pieces.append(name+'\n'+'\n'.join(found))
pathlib.Path(sys.argv[2]).write_text('\n\n'.join(pieces))
PY
cp "$OUTPUT/libWaffleSAP.h" "$OUTPUT/SAP-generated-ABI.h"
cp "$SOURCE/LICENSE" "$OUTPUT/IPATOOL-LICENSE"
tar -czf "$OUTPUT/ipatool-sap-corresponding-source.tar.gz" -C "$WORK" "$(basename "$SOURCE")"
echo 'Built SAP guest library (not the ipatool CLI).'
