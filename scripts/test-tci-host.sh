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
  -I WaffleStore/MapleSyrup/Native scripts/tci-smoke.c WaffleStore/MapleSyrup/Native/TCIProbe.c \
  "$WORK/native/libunicorn.a" -lm -lpthread -o "$WORK/probe"
"$WORK/probe"
