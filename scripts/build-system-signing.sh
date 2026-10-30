#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
OUTPUT="$ROOT/build/Native"
mkdir -p "$OUTPUT"
FINGERPRINT=$( { shasum -a 256 "$ROOT/scripts/build-system-signing.sh" "$ROOT/scripts/prepare-system-signing.py" "$ROOT/scripts/system-signing/bridge.cpp"; xcrun --sdk iphoneos --show-sdk-version; } | shasum -a 256 | cut -d ' ' -f 1)
if [[ -f "$OUTPUT/signing-fingerprint" && -f "$OUTPUT/libCrumbSigning.a" && -f "$OUTPUT/libCrumbCrypto.a" ]] && [[ $(cat "$OUTPUT/signing-fingerprint") == "$FINGERPRINT" ]]; then exit 0; fi
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
curl -fLsS --retry 3 https://github.com/zhlynn/zsign/archive/9299685bab7ef99827e06e5dfda93af747386590.tar.gz -o "$WORK/zsign.tar.gz"
curl -fLsS --retry 3 https://github.com/openssl/openssl/archive/8cf17aaeb4599f8af87fefd810b5b5fee90fe69e.tar.gz -o "$WORK/openssl.tar.gz"
python3 - "$WORK" <<'PY'
import hashlib, pathlib,sys
root=pathlib.Path(sys.argv[1])
for name,pin in [('zsign','97bcc2dbbc0dc5f092dc4a33cbe1f9fdbd5c83bdc291da23db45044d4396270e'),('openssl','03154360ff8808ac38da7530cd69d211e105683e077e9835e48161493c3d0050')]:
    if hashlib.sha256((root/(name+'.tar.gz')).read_bytes()).hexdigest()!=pin:raise SystemExit('Signing dependency integrity mismatch')
PY
tar -xzf "$WORK/zsign.tar.gz" -C "$WORK"
tar -xzf "$WORK/openssl.tar.gz" -C "$WORK"
ZSIGN="$WORK/zsign-9299685bab7ef99827e06e5dfda93af747386590"
OPENSSL="$WORK/openssl-8cf17aaeb4599f8af87fefd810b5b5fee90fe69e"
python3 "$ROOT/scripts/prepare-system-signing.py" "$ZSIGN/src"
(cd "$OPENSSL" && perl Configure ios64-xcrun no-shared no-module no-dso no-tests no-apps no-asm -miphoneos-version-min=16.4 && make -j3 build_libs)
mkdir -p "$WORK/objects"
for source in "$ZSIGN/src/archo.cpp" "$ZSIGN/src/macho.cpp" "$ZSIGN/src/bundle.cpp" "$ZSIGN/src/openssl.cpp" "$ZSIGN/src/signing.cpp" "$ZSIGN"/src/common/{fs,json,log,sha,timer,util}.cpp "$ROOT/scripts/system-signing/bridge.cpp"; do
  object="$WORK/objects/$(basename "$source" .cpp).o"
  xcrun --sdk iphoneos clang++ -arch arm64 -miphoneos-version-min=16.4 -std=c++11 -O2 -DCRUMB_IOS_SIGNER \
    -I"$ZSIGN/src" -I"$ZSIGN/src/common" -I"$OPENSSL/include" -c "$source" -o "$object"
done
xcrun libtool -static -o "$OUTPUT/libCrumbSigning.a" "$WORK"/objects/*.o
cp "$OPENSSL/libcrypto.a" "$OUTPUT/libCrumbCrypto.a"
# Exact public sources and licenses accompany downloadable builds. No signing
# identity, certificate, account state or target IPA is ever a build input.
cp "$ZSIGN/LICENSE" "$OUTPUT/ZSIGN-LICENSE"
cp "$OPENSSL/LICENSE.txt" "$OUTPUT/OPENSSL-LICENSE"
tar -czf "$OUTPUT/system-signing-corresponding-source.tar.gz" -C "$WORK" zsign.tar.gz openssl.tar.gz -C "$ROOT" scripts/system-signing scripts/prepare-system-signing.py scripts/build-system-signing.sh
printf '%s\n' "$FINGERPRINT" > "$OUTPUT/signing-fingerprint"
