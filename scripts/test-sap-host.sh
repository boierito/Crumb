#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
if [[ $(uname -s) != Linux || $(uname -m) != x86_64 ]]; then
  echo 'This manual host SAP smoke test currently supports Linux x86_64.' >&2
  exit 2
fi
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
TCI_SHARED=ON bash "$ROOT/scripts/build-unicorn-tci.sh" "$WORK/native" host
curl -fL --retry 3 https://github.com/majd/ipatool/archive/3411d57f451f5111ae115641c22f7ed17bbd5fbe.tar.gz -o "$WORK/ipatool.tar.gz"
tar -xzf "$WORK/ipatool.tar.gz" -C "$WORK"
SOURCE="$WORK/ipatool-3411d57f451f5111ae115641c22f7ed17bbd5fbe"
cp "$ROOT/scripts/sap-host-probe_test.go" "$SOURCE/internal/sap/tci_probe_test.go"
python3 - "$SOURCE" "$WORK/native/libunicorn.so.2" <<'PY'
import pathlib,sys,json
path=pathlib.Path(sys.argv[1])/'internal/sap/unicorn/library_unix.go'
source=path.read_text()
start=source.index('\tpaths, err := cachedRuntimePaths(ctx)')
end=source.index('\tif err := ctx.Err()',start)
source=source[:start]+'\tpaths := struct { library string }{library: '+json.dumps(sys.argv[2])+'}\n'+source[end:]
path.write_text(source)
PY
cd "$SOURCE"
XDG_CACHE_HOME="$WORK/cache" go test ./internal/sap/machine \
  -run 'TestHardwareBlock|TestGuestServiceDispatchAndStackArguments|TestUnknownGuestServiceFailsClosed|TestGuestAllocatorReusesAndClearsFreedMemory' -v
XDG_CACHE_HOME="$WORK/cache" go test ./internal/sap -run TestTCIDynamicSAPSmoke -v -timeout 6m
echo 'Host SAP test passed. Physical jailed iOS execution and Apple login remain separate milestones.'
