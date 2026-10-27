#!/usr/bin/env bash
set -euo pipefail
# Bound compiler memory for the generated native enrollment CFG package.
export GOMEMLIMIT=${GOMEMLIMIT:-4GiB}
export GOGC=${GOGC:-30}
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
mkdir -p "$SOURCE/wafflebridge/mobilelicense"
cp "$ROOT/MapleSyrup/NativeSAP/mobilelicense/"*.go "$SOURCE/wafflebridge/mobilelicense/"
cp "$ROOT/MapleSyrup/NativeSAP/mobilelicense/LICENSE-UPSTREAM" "$SOURCE/wafflebridge/mobilelicense/"
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
FAIRPLAY_REV=ef98227e061838320424744d21a513baba7ade28
cp "$ROOT/third_party/archives/ipatool-fairplay-$FAIRPLAY_REV.tar.gz" "$WORK/fairplay.tar.gz"
python3 - "$WORK/fairplay.tar.gz" <<'CHECK'
import hashlib, sys
if hashlib.sha256(open(sys.argv[1], 'rb').read()).hexdigest() != '49f169ef2e14c89e3b55a9528057e7db66f236ca5425ff00e47d32f64dddf44d':
    raise SystemExit('FairPlay source integrity mismatch')
CHECK
tar -xzf "$WORK/fairplay.tar.gz" -C "$WORK"
cp -R "$WORK/ipatool-$FAIRPLAY_REV/internal/fairplay" "$SOURCE/internal/fairplay"
# iOS has the same Unix permission checks; treat its filesystem as case insensitive.
python3 "$ROOT/scripts/prepare-fairplay-package.py" "$SOURCE/internal/fairplay/files.go"
cp "$ROOT/MapleSyrup/NativeSAP/fairplayinspect.go" "$SOURCE/internal/fairplay/"
cp "$ROOT/MapleSyrup/NativeSAP/fairplaymulti.go" "$SOURCE/internal/fairplay/"
cp "$ROOT/MapleSyrup/NativeSAP/fairplaymulti_test.go" "$SOURCE/internal/fairplay/"
(cd "$SOURCE" && go get golang.org/x/arch@v0.32.0 && go test ./internal/fairplay)
python3 "$ROOT/scripts/prepare-fairplay-portable.py" "$WORK/ipatool-$FAIRPLAY_REV/tools/fairplay-recovery" "$SOURCE/wafflebridge/recovery"
bash "$ROOT/scripts/fetch-fairplay-raw.sh" "$SOURCE/wafflebridge/recovery/rawsource"
cp "$ROOT/MapleSyrup/NativeSAP/captureanalysis.go" "$SOURCE/wafflebridge/"
cp "$ROOT/MapleSyrup/NativeSAP/fairplayenrollment.go" "$ROOT/MapleSyrup/NativeSAP/fairplaycapturebridge.go" "$ROOT/MapleSyrup/NativeSAP/fairplayerrors.go" "$ROOT/MapleSyrup/NativeSAP/fairplayprobe.go" "$ROOT/MapleSyrup/NativeSAP/fairplaycausal.go" "$ROOT/MapleSyrup/NativeSAP/fairplaybridgeaudit.go" "$SOURCE/wafflebridge/"
mkdir -p "$SOURCE/wafflebridge/labtest"
cp "$ROOT/MapleSyrup/NativeSAP/fairplayenrollment.go" "$ROOT/MapleSyrup/NativeSAP/fairplayenrollment_test.go" "$ROOT/MapleSyrup/NativeSAP/fairplaycapturebridge.go" "$ROOT/MapleSyrup/NativeSAP/fairplaycapturebridge_test.go" "$ROOT/MapleSyrup/NativeSAP/fairplayerrors.go" "$ROOT/MapleSyrup/NativeSAP/fairplayprobe.go" "$ROOT/MapleSyrup/NativeSAP/fairplayprobe_test.go" "$ROOT/MapleSyrup/NativeSAP/fairplaycausal.go" "$ROOT/MapleSyrup/NativeSAP/fairplaycausal_test.go" "$ROOT/MapleSyrup/NativeSAP/fairplaybridgeaudit.go" "$ROOT/MapleSyrup/NativeSAP/fairplaybridgeaudit_test.go" "$SOURCE/wafflebridge/labtest/"
cp "$ROOT/MapleSyrup/NativeSAP/fairplaylab.go" "$ROOT/MapleSyrup/NativeSAP/fairplaylabcore.go" "$SOURCE/wafflebridge/"
cp "$ROOT/MapleSyrup/NativeSAP/captureanalysis.go" "$ROOT/MapleSyrup/NativeSAP/captureanalysis_test.go" "$ROOT/MapleSyrup/NativeSAP/fairplaylabcore.go" "$ROOT/MapleSyrup/NativeSAP/fairplaylab_test.go" "$SOURCE/wafflebridge/labtest/"
# Full rawsource/lab race tests run in the Linux fairplay-portable CI job.
# Compile the large raw CFG once for iOS here, rather than twice for host tests.
(cd "$SOURCE" && go test -p=1 ./wafflebridge/recovery/state ./wafflebridge/recovery/frame ./wafflebridge/recovery/sidb)
cp "$ROOT/MapleSyrup/NativeSAP/gueststate.go" "$SOURCE/internal/sap/machine/"
python3 "$ROOT/scripts/prepare-retained-storeagent.py" "$SOURCE/internal/sap/machine"
cp "$ROOT/MapleSyrup/NativeSAP/guestrestore.go" "$SOURCE/internal/sap/machine/"
cp "$ROOT/MapleSyrup/NativeSAP/kbsyncretained.go" "$SOURCE/internal/sap/machine/"
cp "$ROOT/MapleSyrup/NativeSAP/kbsyncbridge.go" "$SOURCE/wafflebridge/"
cp "$ROOT/MapleSyrup/NativeSAP/gueststate_test.go" "$ROOT/MapleSyrup/NativeSAP/guestrestore_test.go" "$SOURCE/internal/sap/machine/"
(cd "$SOURCE" && CGO_ENABLED=0 go test -run '^TestGuestState' ./internal/sap/machine)
(cd "$SOURCE" && CGO_ENABLED=0 GOOS=ios GOARCH=arm64 go build ./wafflebridge/mobilelicense)
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
cat "$WORK/ipatool-$FAIRPLAY_REV/tools/fairplay-recovery/LICENSE-DARKPHOENIX.txt" >> "$OUTPUT/GoThirdPartyNotices.txt"
cp "$WORK/ipatool-$FAIRPLAY_REV/tools/fairplay-recovery/LICENSE-DARKPHOENIX.txt" "$SOURCE/wafflebridge/recovery/"
cat "$ROOT/third_party/unfairplay-generic/LICENSE" "$ROOT/third_party/unfairplay-generic/NOTICE" >> "$OUTPUT/GoThirdPartyNotices.txt"
cp "$ROOT/third_party/unfairplay-generic/LICENSE" "$SOURCE/wafflebridge/recovery/rawsource/LICENSE-UNFAIRPLAY"
cp "$ROOT/third_party/unfairplay-generic/NOTICE" "$SOURCE/wafflebridge/recovery/rawsource/NOTICE-UNFAIRPLAY"
cp "$ROOT/third_party/unfairplay-generic/PROFILE.json" "$SOURCE/wafflebridge/recovery/rawsource/GENERIC-PROFILE-PROVENANCE.json"
cp "$OUTPUT/libWaffleSAP.h" "$OUTPUT/SAP-generated-ABI.h"
cp "$SOURCE/LICENSE" "$OUTPUT/IPATOOL-LICENSE"
tar -czf "$OUTPUT/ipatool-sap-corresponding-source.tar.gz" -C "$WORK" "$(basename "$SOURCE")"
echo 'Built SAP guest library (not the ipatool CLI).'
