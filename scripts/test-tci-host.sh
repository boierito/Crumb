#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if [[ $(uname -s) != Linux || $(uname -m) != x86_64 ]]; then
  echo 'This seccomp test requires Linux x86_64.' >&2
  exit 2
fi
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
bash scripts/build-unicorn-tci.sh "$WORK/native" host
cc -std=c11 -Wall -Wextra -Werror -I "$WORK/native/include" \
  -I scripts/fixtures/tci scripts/tci-smoke.c scripts/fixtures/tci/TCIProbe.c \
  "$WORK/native/libunicorn.a" -lm -lpthread -o "$WORK/probe"
"$WORK/probe"
