#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
OUTPUT="$ROOT/build/Native"
mkdir -p "$OUTPUT"
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
SDK=$(xcrun --sdk iphoneos --show-sdk-path)
CLANG=$(xcrun --sdk iphoneos --find clang)
cd "$ROOT/MapleSyrup/NativeIPA"
GOOS=ios GOARCH=arm64 CGO_ENABLED=1 CC="$CLANG" SDKROOT="$SDK" \
  CGO_CFLAGS="-isysroot $SDK -target arm64-apple-ios16.4" \
  CGO_LDFLAGS="-isysroot $SDK -target arm64-apple-ios16.4" \
  go build -buildmode=c-archive -trimpath -o "$OUTPUT/libCrumbIPA.a" .
go list -deps -json ./... > "$WORK/dependencies.json"
python3 - "$WORK/dependencies.json" "$OUTPUT/GoThirdPartyNotices.txt" "$(go env GOROOT)" <<'PY'
import json,pathlib,sys
raw=pathlib.Path(sys.argv[1]).read_text()
decoder=json.JSONDecoder(); modules={}
while raw.strip():
    value,end=decoder.raw_decode(raw.lstrip()); raw=raw.lstrip()[end:]
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
python3 - "$ROOT/Crumb/Resources/ThirdPartyNotices.txt" "$OUTPUT/IPATOOL-LICENSE" <<'PY'
import pathlib,sys
s=pathlib.Path(sys.argv[1]).read_text()
s=s.split('## ipatool MIT notice\n',1)[1].split('\n## ',1)[0].strip()
assert s.startswith('MIT License')
pathlib.Path(sys.argv[2]).write_text(s+'\n')
PY
tar -czf "$OUTPUT/ipatool-ipa-corresponding-source.tar.gz" -C "$ROOT" MapleSyrup/NativeIPA
echo 'Built local streaming IPA library without legacy Go SAP or Unicorn.'
