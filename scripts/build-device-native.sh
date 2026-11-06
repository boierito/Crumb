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
  -framework UIKit -framework Foundation -framework CoreGraphics -lz DecryptHelper/main.m DecryptHelper/CrumbJobControl.m DecryptHelper/CrumbMachO.c DecryptHelper/CrumbFileAccess.c DecryptHelper/CrumbLicenseResources.m DecryptHelper/CrumbIPAArchive.m -o "$HELPER/CrumbDecrypt"
sips -z 120 120 branding/Crumb-AppIcon-1024.png --out "$HELPER/Icon60@2x.png" >/dev/null
sips -z 180 180 branding/Crumb-AppIcon-1024.png --out "$HELPER/Icon60@3x.png" >/dev/null
sips -z 152 152 branding/Crumb-AppIcon-1024.png --out "$HELPER/Icon76@2x.png" >/dev/null
sips -z 167 167 branding/Crumb-AppIcon-1024.png --out "$HELPER/Icon83.5@2x.png" >/dev/null
python3 - "$HELPER/Info.plist" <<'PY'
import plistlib,sys,os
with open(sys.argv[1],'wb') as f:
 plistlib.dump(dict(CFBundleIcons=dict(CFBundlePrimaryIcon=dict(CFBundleIconFiles=['Icon60'],UIPrerenderedIcon=False)), **{'CFBundleIcons~ipad':dict(CFBundlePrimaryIcon=dict(CFBundleIconFiles=['Icon60','Icon76','Icon83.5'],UIPrerenderedIcon=False))}, CFBundleIdentifier='com.certlium.crumb.decrypt',CFBundleExecutable='CrumbDecrypt',CFBundleName='Crumb Decrypt',CFBundleDisplayName='Crumb Decrypt',CFBundlePackageType='APPL',CFBundleVersion=os.environ.get('CURRENT_PROJECT_VERSION','23037'),CFBundleShortVersionString='1.0',MinimumOSVersion='27.0',UIDeviceFamily=[1,2],UIFileSharingEnabled=True,LSRequiresIPhoneOS=True,UILaunchScreen={},CFBundleURLTypes=[dict(CFBundleURLSchemes=['crumb-decrypt-helper'])]),f)
PY
# Store the unsigned template as opaque data, so a sideloader does not sign it
# as a nested app. Crumb unpacks, customizes and signs it on-device.
(cd "$HELPER" && zip -q -r ../CrumbDecryptTemplate.zip .)
mkdir -p "${TARGET_BUILD_DIR:?}/${UNLOCALIZED_RESOURCES_FOLDER_PATH:?}"
cp build/Native/CrumbDecryptTemplate.zip "$TARGET_BUILD_DIR/$UNLOCALIZED_RESOURCES_FOLDER_PATH/"

cp NativeDevice/vendor/apple-codesign-quick/COPYING.LESSER "$TARGET_BUILD_DIR/$UNLOCALIZED_RESOURCES_FOLDER_PATH/CrumbDevice-LGPL.txt"
