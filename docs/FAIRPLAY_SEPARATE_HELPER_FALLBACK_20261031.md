# FairPlay: separately installed helper fallback (2026-10-31)

## Actual purpose

Allow **Crumb to remain installed and untouched** while it prepares and signs a different app named **Crumb FP Helper**, containing one purchased App Store IPA as *installed resources*. Once iOS installs that helper, open it and decrypt the installed target with the existing cryptid-2 native FairPlay reader, then share the fully verified output using Files/Share.

This is a **fallback**, not the desired permanent iPhone-only direct data-container decryption. It requires installing a different helper when the target changes, not reinstalling Crumb. Nothing in this change asserts that one installed helper supports an arbitrary later IPA without being reinstalled. It never uses JIT, jailbreak, an attached PC, or a cloud signing service.

## How it works

1. Crumb's native download/import returns an encrypted IPA owned by the account.
2. In the Decrypt view, choose **Prepare separate helper**.
3. `packageipa.PrepareSystemSeparateHelper` copies the installed Crumb app into a **new IPA** and changes the helper's root `CFBundleIdentifier` from `<crumb ID>` to `<crumb ID>.fphelper`. It also changes the display name and removes duplicate URL schemes, without modifying the installed Crumb application. It copies the original encrypted application and its unchanged SINF into `CrumbSystemPayload/Target/`, applying the exact manifest-directed SINF replication logic.
4. Configure an existing private P12 and mobileprovision. The provisioning profile must legitimately authorize the **separate helper identifier**, normally through a suitable wildcard profile. An exact-ID profile for Crumb will fail closed with signing error 18. The bundle and all nested app extension identifiers must be permitted. A user must not modify a provisioning profile to impersonate another app or claim an unsupported entitlement.
5. **Sign separate helper** runs local `CrumbSignPreparedBundle`, using the helper ID and the actual Crumb team. The archive validator compares every private purchased resource to its original bytes and refuses accidental changes.
6. **Install Crumb FP Helper** runs the existing loopback OTA installer and asks iOS for installation. This may be refused by Apple if Developer Mode, profile, device registration, trust, expiration or signing identity are invalid. Installation success is not inferred from sending the request.
7. Launch the separately installed `Crumb FP Helper` icon. Its startup detects `.fphelper` and opens `SystemFairPlayView` for the embedded `CrumbSystemPayload/Target`. Tap **Decrypt**.
8. The existing `SystemNativeProcessor` reads ciphertext from the **installed bundle**, not the Crumb data container. Every encrypted ARM64 range must map, remap, read all pages, change all pages and pass `ExportSystemResult` gates. A failure produces a report, never a falsely decrypted IPA. On success, use **Save IPA** to share a sanitized, unsigned export. It requires signing separately if installation is desired.

## Verified versus not verified

The original FairPlayed standalone reference has one real on-device positive control: 178 pages from an installed resource were recovered on iPhone 16 Pro / iOS 27.0.1. That was not this exact new signed-helper flow.

Host tests now verify standalone helper bundle separation, byte-for-byte target/SINF preservation, absence of original URL scheme registration and unchanged `Crumb.app` input. Native protocol and iOS app build CI are necessary gates but cannot emulate the iPhone's FairPlay service or authorize installation.

**No claim of completed decryption using this new auxiliary helper is made until a device report confirms installed-resource native remap success and whole-IPA export.**

## Limitations and privacy

- One encrypted target per separately prepared/installed helper revision. A new target requires updating that auxiliary installation, **not Crumb**.
- Requires an Apple account entitled to the app, a compatible signing certificate, an actual provisioning profile that covers the second app ID, and an install route accepted by iOS.
- A new helper is a separate installed app and cannot directly access Crumb's Documents. Save the resulting IPA through the system Share sheet to a user-chosen Files location. The helper contains purchased SINF in its installed resources; do not publish its signed IPA, license files or installed helper contents.
- Other previously installed Crumb data and the encrypted original are retained.
- An OTA loopback installer delegates installation to iOS; it does not bypass iOS signing checks or claim install success.
- No CVE, arbitrary package write, sandbox escape, or private entitlement is needed for the packaging or signing steps.

## Code

- `MapleSyrup/NativeSAP/packageipa/system_job.go`: original in-place preparation retained; new separate identity package variant.
- `MapleSyrup/NativeSAP/fairplaylabcore.go`: separate helper route.
- `WaffleStore/App/FairPlayLabView.swift`: explicit helper ID argument in native request.
- `WaffleStore/App/SystemFairPlayView.swift`: separate fallback and signing/installation actions.
- `WaffleStore/ContentView.swift`: helper startup goes to its installed decrypted-target screen.
- `MapleSyrup/NativeSAP/packageipa/system_test.go`: bytes/inventory/identity preservation tests.

## Critical validation on iPhone

1. Prepare helper from the same licensed IPA used in the successful FairPlayed installed-resource test, if available.
2. Confirm signing with a profile permitting `<Crumb bundle ID>.fphelper`, and perform the iOS helper installation.
3. Open `Crumb FP Helper`. Confirm the embedded target name appears and tap Decrypt.
4. Record `native_remap_return`, `native_remap_errno`, `pages_total`, `pages_changed`, and whether a successful Shareable unsigned IPA was produced.
5. If it fails, record only the safe numeric report; do not upload actual licensed IPA, private keys, profile, SINF or personal data to GitHub.
