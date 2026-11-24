#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
export PATH="$HOME/.cargo/bin:$PATH"
rustup target add aarch64-apple-ios
cargo build --locked --release --target aarch64-apple-ios --manifest-path NativeDevice/Cargo.toml
mkdir -p build/Native
cp NativeDevice/target/aarch64-apple-ios/release/libcrumb_device.a build/Native/
mkdir -p "${TARGET_BUILD_DIR:?}/${UNLOCALIZED_RESOURCES_FOLDER_PATH:?}"
rm -f "$TARGET_BUILD_DIR/$UNLOCALIZED_RESOURCES_FOLDER_PATH/CrumbDecryptTemplate.zip"

cp NativeDevice/vendor/apple-codesign-quick/COPYING.LESSER "$TARGET_BUILD_DIR/$UNLOCALIZED_RESOURCES_FOLDER_PATH/CrumbDevice-LGPL.txt"
