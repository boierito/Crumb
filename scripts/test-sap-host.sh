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
cp "$ROOT/scripts/kbsync-host-probe_test.go" "$SOURCE/internal/sap/kbsync_tci_test.go"
python3 "$ROOT/scripts/prepare-retained-storeagent.py" "$SOURCE/internal/sap/machine"
for name in gueststate.go gueststate_test.go guestrestore.go guestrestore_test.go kbsyncretained.go; do
 cp "$ROOT/MapleSyrup/NativeSAP/$name" "$SOURCE/internal/sap/machine/"
done
cp "$ROOT/scripts/producer-native-probe_test.go" "$SOURCE/internal/sap/machine/"
cp "$ROOT/scripts/producer-capture-host_test.go" "$SOURCE/internal/sap/machine/"
# Reuse the production state validator for private account-state replay; compile
# it only as a test helper in the producer's package.
python3 - "$ROOT/MapleSyrup/NativeSAP/fairplayauthorizationstate.go" "$SOURCE/internal/sap/machine/authorization_state_host_test.go" <<'PYCODE'
import pathlib, sys
source = pathlib.Path(sys.argv[1]).read_text()
assert source.count('package main') == 1
pathlib.Path(sys.argv[2]).write_text(source.replace('package main', 'package machine', 1))
PYCODE


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
XDG_CACHE_HOME="$WORK/cache" go test ./internal/sap/machine -run 'TestGuestState|TestProducerNative' -v -timeout 6m
XDG_CACHE_HOME="$WORK/cache" go test ./internal/sap -run 'TestTCIDynamicSAPSmoke|TestTCIKBSyncSmoke' -v -timeout 6m
python3 "$ROOT/scripts/verify-storeagent-license-profile.py" "$WORK/cache/ipatool/sap/apple-storeagent-v1/storeagent"
python3 "$ROOT/scripts/verify-store-authorization-profile.py" "$WORK/cache/ipatool/sap/apple-storeagent-v1/storeagent" "$WORK/cache/ipatool/sap/apple-assets-v2/CommerceKit"
if [[ -n "${CRUMB_PRIVATE_CAPTURE:-}" ]]; then
  XDG_CACHE_HOME="$WORK/cache" go test ./internal/sap/machine -run '^TestRetainedPrivateCaptureProvisioning$' -count=1 -timeout 6m
fi
echo 'Host SAP test passed. Physical jailed iOS execution and Apple login remain separate milestones.'
