#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if [[ ${PLATFORM_NAME:-iphoneos} != iphoneos ]]; then
  echo 'The native SAP probe currently targets physical arm64 iOS devices. Use swift test for host protocol tests.' >&2
  exit 2
fi
FINGERPRINT=$( { shasum -a 256 scripts/prepare-unicorn-tci.py scripts/build-unicorn-tci.sh scripts/build-sap-native.sh MapleSyrup/NativeSAP/main.go; xcrun --sdk iphoneos --show-sdk-version; } | shasum -a 256 | cut -d ' ' -f 1)
if [[ -f build/Native/fingerprint && -f build/Native/libunicorn.a && -f build/Native/libWaffleSAP.a ]] && [[ $(cat build/Native/fingerprint) == "$FINGERPRINT" ]]; then
  exit 0
fi
bash scripts/build-unicorn-tci.sh
bash scripts/build-sap-native.sh
printf '%s\n' "$FINGERPRINT" > build/Native/fingerprint
