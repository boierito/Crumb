#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if [[ $(uname -s) != Linux || $(uname -m) != x86_64 ]]; then
  echo 'This host reproduction requires Linux x86_64; use the iOS diagnostic on devices.' >&2
  exit 2
fi
UNICORN_PATH=$(python3 -c 'import pathlib,unicorn; print(pathlib.Path(unicorn.__file__).parent)')
TEST_DIR=$(mktemp -d)
trap 'rm -rf "$TEST_DIR"' EXIT
cc -std=c11 -Wall -Wextra -Werror -I "$UNICORN_PATH/include" scripts/unicorn-no-exec.c -ldl -o "$TEST_DIR/probe"
"$TEST_DIR/probe" "$UNICORN_PATH/lib/libunicorn.so.2" allow-exec
set +e
"$TEST_DIR/probe" "$UNICORN_PATH/lib/libunicorn.so.2" deny-exec > "$TEST_DIR/denied.log" 2>&1
RESULT=$?
set -e
cat "$TEST_DIR/denied.log"
if [[ $RESULT == 1 ]] && grep -q 'Could not allocate dynamic translator buffer' "$TEST_DIR/denied.log"; then
  echo 'Unicorn terminated the process while allocating its dynamic translator buffer (exit 1).'
elif [[ $RESULT != 10 ]]; then
  echo "Expected Unicorn allocation failure; received $RESULT." >&2
  exit 1
fi
echo 'PASS: Unicorn works with executable memory and fails when new executable mappings are denied.'
echo 'This reproduces a runtime dependency; it does not substitute for iOS 26/27 device validation.'
