#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
export PATH="$HOME/.cargo/bin:$PATH"
rustup target add aarch64-apple-ios
cargo build --locked --release --target aarch64-apple-ios --manifest-path NativeDevice/Cargo.toml
mkdir -p build/Native
cp NativeDevice/target/aarch64-apple-ios/release/libcrumb_device.a build/Native/
HELPER=build/Native/CrumbDecryptTemplate
rm -rf "$HELPER"
rm -f build/Native/CrumbDecryptTemplate.zip
mkdir -p "$HELPER"
xcrun --sdk iphoneos clang -arch arm64 -miphoneos-version-min=16.4 -Os -fobjc-arc \
  -framework UIKit -framework Foundation -framework CoreGraphics -lz DecryptHelper/main.m DecryptHelper/CrumbMachO.c -o "$HELPER/CrumbDecrypt"
python3 - "$HELPER/Info.plist" <<'PY'
import plistlib,sys
with open(sys.argv[1],'wb') as f:
 plistlib.dump(dict(CFBundleIdentifier='com.certlium.crumb.decrypt',CFBundleExecutable='CrumbDecrypt',CFBundleName='Crumb Decrypt',CFBundleDisplayName='Crumb Decrypt',CFBundlePackageType='APPL',CFBundleVersion='1',CFBundleShortVersionString='1.0',MinimumOSVersion='27.0',UIDeviceFamily=[1,2],UIFileSharingEnabled=True,LSRequiresIPhoneOS=True,UILaunchScreen={},CFBundleURLTypes=[dict(CFBundleURLSchemes=['crumb-decrypt-helper'])]),f)
PY
# Store the unsigned template as opaque data, so a sideloader does not sign it
# as a nested app. Crumb unpacks, customizes and signs it on-device.
(cd "$HELPER" && zip -q -r ../CrumbDecryptTemplate.zip .)
mkdir -p "${TARGET_BUILD_DIR:?}/${UNLOCALIZED_RESOURCES_FOLDER_PATH:?}"
cp build/Native/CrumbDecryptTemplate.zip "$TARGET_BUILD_DIR/$UNLOCALIZED_RESOURCES_FOLDER_PATH/"

cp NativeDevice/vendor/apple-codesign-quick/COPYING.LESSER "$TARGET_BUILD_DIR/$UNLOCALIZED_RESOURCES_FOLDER_PATH/CrumbDevice-LGPL.txt"
