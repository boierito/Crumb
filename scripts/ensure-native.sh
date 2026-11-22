#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if [[ ${PLATFORM_NAME:-iphoneos} != iphoneos ]]; then
  echo 'The native IPA package module targets physical arm64 iOS devices. Use Go and Swift host tests.' >&2
  exit 2
fi
FINGERPRINT=$( { shasum -a 256 scripts/build-ipa-native.sh $(find MapleSyrup/NativeIPA -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) | sort); xcrun --sdk iphoneos --show-sdk-version; } | shasum -a 256 | cut -d ' ' -f 1)
if [[ -f build/Native/fingerprint && -f build/Native/libCrumbIPA.a ]] && [[ $(cat build/Native/fingerprint) == "$FINGERPRINT" ]]; then
  bash scripts/build-device-native.sh
  exit 0
fi
bash scripts/build-ipa-native.sh
printf '%s\n' "$FINGERPRINT" > build/Native/fingerprint
bash scripts/build-device-native.sh
