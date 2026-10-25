#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
DESTINATION=$1
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
REV=f7242cce58077a4af8db738ed295ee23ad8f067d
cp "$ROOT/third_party/archives/ipatool-fairplay-$REV.tar.gz" "$WORK/source.tar.gz"
python3 - "$WORK/source.tar.gz" <<'PY'
import hashlib, pathlib, sys
if hashlib.sha256(pathlib.Path(sys.argv[1]).read_bytes()).hexdigest() != 'e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39':
    raise SystemExit('Modern FairPlay source integrity mismatch')
PY
tar -xzf "$WORK/source.tar.gz" -C "$WORK"
python3 "$ROOT/scripts/prepare-fairplay-raw.py" "$WORK/ipatool-$REV/tools/fairplay-recovery" "$DESTINATION"
