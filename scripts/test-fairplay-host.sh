#!/usr/bin/env bash
# Synthetic public fixtures only: this does not prove real IPA decryption.
set -euo pipefail
# Upstream permission tests intentionally create a 0644 fixture inside t.TempDir.
# An inherited 0077 umask would turn it into 0600 and invalidate that test.
umask 022
# Generated raw-state CFGs are large; keep compiler GC within CI memory.
export GOMEMLIMIT=${GOMEMLIMIT:-4GiB}
export GOGC=${GOGC:-30}
ROOT=$(cd "$(dirname "$0")/.." && pwd)
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
REV=ef98227e061838320424744d21a513baba7ade28
curl -fL --retry 3 "https://github.com/DrWhax/ipatool/archive/$REV.tar.gz" -o "$WORK/fairplay.tar.gz"
python3 - "$WORK/fairplay.tar.gz" <<'PY'
import hashlib, pathlib, sys
if hashlib.sha256(pathlib.Path(sys.argv[1]).read_bytes()).hexdigest() != '49f169ef2e14c89e3b55a9528057e7db66f236ca5425ff00e47d32f64dddf44d':
    raise SystemExit('FairPlay source integrity mismatch')
PY
tar -xzf "$WORK/fairplay.tar.gz" -C "$WORK"
SOURCE="$WORK/ipatool-$REV"
# Isolate the portable package from upstream CLI/macOS dependencies.
mkdir -p "$WORK/test/internal" "$WORK/test/wafflebridge/labtest"
cp -R "$SOURCE/internal/fairplay" "$WORK/test/internal/"
python3 "$ROOT/scripts/prepare-fairplay-package.py" "$WORK/test/internal/fairplay/files.go"
cp "$ROOT/MapleSyrup/NativeSAP/fairplayinspect.go" "$WORK/test/internal/fairplay/"
cp "$ROOT/MapleSyrup/NativeSAP/fairplaymulti.go" "$WORK/test/internal/fairplay/"
cp "$ROOT/MapleSyrup/NativeSAP/fairplaymulti_test.go" "$WORK/test/internal/fairplay/"
python3 "$ROOT/scripts/prepare-fairplay-portable.py" "$SOURCE/tools/fairplay-recovery" "$WORK/test/wafflebridge/recovery"
bash "$ROOT/scripts/fetch-fairplay-raw.sh" "$WORK/test/wafflebridge/recovery/rawsource"
for name in fairplaycapturebridge.go fairplaycapturebridge_test.go fairplayenrollment.go fairplayenrollment_test.go captureanalysis.go captureanalysis_test.go fairplaylabcore.go fairplaylab_test.go; do
  cp "$ROOT/MapleSyrup/NativeSAP/$name" "$WORK/test/wafflebridge/labtest/"
done
cat > "$WORK/test/go.mod" <<'MOD'
module github.com/majd/ipatool/v2

go 1.26.0

require (
    golang.org/x/arch v0.32.0
    howett.net/plist v1.0.1
)
MOD
cd "$WORK/test"
go mod tidy
go test -p=1 -race -count=1 ./internal/fairplay ./wafflebridge/recovery/... ./wafflebridge/labtest
