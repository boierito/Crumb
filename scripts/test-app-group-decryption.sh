#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/.."
task_root=$(mktemp -d)
trap 'rm -rf "$task_root"' EXIT
mkdir -p "$task_root/Sources/ArchiveTests"
cp Crumb/Device/AppGroupDecryption.swift Crumb/Device/IPAExportPrivacy.swift tests/AppGroupDecryption/main.swift "$task_root/Sources/ArchiveTests/"
cat > "$task_root/Package.swift" <<'SWIFT'
// swift-tools-version: 5.9
import PackageDescription
let package = Package(name: "ArchiveTests", platforms: [.macOS(.v12)], dependencies: [
    .package(url: "https://github.com/weichsel/ZIPFoundation.git", exact: "0.9.20")
], targets: [.executableTarget(name: "ArchiveTests", dependencies: ["ZIPFoundation"])])
SWIFT
swift run --package-path "$task_root" ArchiveTests
