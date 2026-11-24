# Development

Open `Crumb.xcodeproj` with Xcode 26 and use the shared Crumb scheme. The build targets physical arm64 iOS devices (iOS 16.4+).

## Structure

- `Crumb/`: SwiftUI application, downloads, signing, App Groups and device workflow.
- `Crumb/MapleSyrup/`: iPadN SAP runtime integration and login.
- `MapleSyrup/SAPKit/`: Apple Store protocol library.
- `MapleSyrup/NativeIPA/`: Go IPA packaging library.
- `NativeDevice/`: Rust signing and local IPA validation.
- `scripts/`: native builds, pinned SAP runtime preparation and IPA verification.

## Build

Use the Build IPA GitHub Actions workflow, or install Xcode 26, Go 1.25.1 and Rust 1.99.1 locally:

```sh
python3 scripts/prepare-ipadn-sap-runtime.py
xcodebuild -project Crumb.xcodeproj -scheme Crumb -configuration Release \
  -sdk iphoneos -destination 'generic/platform=iOS' \
  CODE_SIGNING_ALLOWED=NO CODE_SIGNING_REQUIRED=NO build
```

CI builds Debug and Release, verifies the produced IPA and provides corresponding source archives. Successful compilation does not establish Apple login or physical-device acceptance.

See [App Groups](APP_GROUP_DECRYPTION.md) and [third-party notices](../THIRD_PARTY_NOTICES.md) for usage requirements and licenses. Historical development logs remain in Git history.
