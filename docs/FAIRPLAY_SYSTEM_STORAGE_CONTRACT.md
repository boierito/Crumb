# Imported IPA: system-service and storage investigation

Research date: 2026-10-30. This investigation has **not** fixed the native remap
failure and has not demonstrated a generic imported-IPA decryption route.
No new experimental phone build is warranted by these results alone.

## Original firmware evidence

Inputs were selectively extracted from the official iPhone17,1 iOS 27.0.1
24A446 restore image already identified in
[the direct-route investigation](FAIRPLAY_NO_REINSTALL_RESEARCH.md). The supplied
phone log identifies the model and OS version, not the exact installed build.
No Apple Account, user license, private capture or original IPA was used to
acquire these firmware inputs.

The original launchd plist is
`/System/Library/LaunchDaemons/com.apple.fairplayd.H2.plist`. It starts
`/usr/sbin/fairplayd.H2` with **UserName=mobile and GroupName=mobile**. Its declared
Mach services include `com.apple.unfreed`, HostSpecialPort 17, and
`com.apple.fairplayd.xpc`. A declared service name does not establish that an
ordinary client can invoke an arbitrary operation on it.

The daemon's embedded XML entitlements include
`com.apple.private.FairPlayIOKitUserClient.access`,
`com.apple.security.iokit-user-client-class`,
`com.apple.private.security.storage.fpsd`, and
`com.apple.security.system-container`. These are **daemon entitlements**, not
demonstrated requirements for every client of that daemon. This audit does not
verify the Apple code signature cryptographically or infer all effective
privileges from the XML blob alone.

The matching kernel's `fairplayd.H2` protobox profile, decoded with ipsw
v3.1.739, contains `(allow default)` and no explicit `file-read*` or
`file-write*` rules. This is one profile, not the entirety of iOS security policy.
It neither proves access to app data nor identifies the cause of errno 12.
The prior generic statement “the daemon sandbox denies the path” therefore
remains unproven for this device failure. System Policy, container protections,
license interpretation and crypter initialization have not been distinguished.

Two tempting changes are unsupported by this evidence:

- Making private license files world-readable: a different daemon UID is not
  established, and the launchd identity is `mobile`. The private permissions
  remain unchanged. This does not independently prove every access check passes.
- Adding `F_ADDFILESIGS` as the presumed missing step: the supplied **successful**
  embedded-target FairPlayed log recorded `sig=-` for both processed executables.
  Registration was not necessary on that successful path. Whether some distinct
  registration operation could affect another path remains untested.

## Caller identity helpers: located, not a validated alternate interface

The original daemon imports `open`, `open_dprotected_np`, `fcntl`, `fstat`,
`lstat`, `read`, and `container_system_path_for_identifier`. No named
`sandbox_extension_consume` or `sandbox_check` import was found. This does not
exclude an indirect equivalent or a security check in another component.

Two helper functions create a SecTask from a caller audit token and read a
named entitlement. Original direct calls are:

| Helper entry | SecTaskCreateWithAuditToken | SecTaskCopyValueForEntitlement |
| --- | --- | --- |
| `0x100001b40` | `0x100001bb4` | `0x100001bc8` |
| `0x100001cec` | `0x100001d58` | `0x100001d6c` |

The first helper converts a number or decimal string to a numeric output; the
second copies a string output. The entitlement name is a caller argument.
There are no direct immediate B/BL calls to these two helper entries in
`__text`. Authenticated pointer records referring to their offsets were found
at `0x101107ae0` and `0x101107be0`. The indirect callers, selected entitlement
names, operations and authorization consequences are **not resolved**. These
helpers must not be described as a demonstrated SINF-service entitlement gate,
nor exposed in Crumb as though an accepted IPC contract had been found.

## Background Assets: dynamic manifest, existing app-group requirement

The SDK does provide dynamic manifests. It would be incorrect to dismiss this
framework on the premise that all manifests must be embedded at build time.
Both `BAAssetPackManifest.initWithContentsOfURL:applicationGroupIdentifier:error:`
and `initFromData:applicationGroupIdentifier:error:` accept a local manifest.
Their documentation explicitly identifies the supplied **application group** as
the place to store downloaded unmanaged packs. `BAURLDownload` similarly
requires an application group and documents an **HTTPS** request URL.

`BAAssetPackManager` offers a file descriptor API as well as URL lookup. Its
shared manager also requires the corresponding managed downloader extension.
None of these contracts establishes a special FairPlay-readable location or
delegates the caller's file permissions to fairplayd. In the latest device
result, zero authorized dedicated shared locations were available. Creating a
manifest does not grant an App Group entitlement.

Consequently this is not an independent, currently justified fallback for the
tested installation. If a dedicated group is legitimately available, Crumb's
existing shared-container attempt can test that storage class directly, without
first adding managed downloads. Success there remains unverified. The SDK
contract alone does not establish the backing vnode or fairplayd's access to it.

Pinned SDK source (Apple-authored headers mirrored in the repository):

- [BAAssetPackManifest.h](https://github.com/xybp888/iOS-SDKs/blob/4a48e5409f8ea6ac422c44d91b86fa0f8ae010ff/iPhoneOS27.0.sdk/System/Library/Frameworks/BackgroundAssets.framework/Headers/BAAssetPackManifest.h)
- [BAAssetPackManager.h](https://github.com/xybp888/iOS-SDKs/blob/4a48e5409f8ea6ac422c44d91b86fa0f8ae010ff/iPhoneOS27.0.sdk/System/Library/Frameworks/BackgroundAssets.framework/Headers/BAAssetPackManager.h)
- [BAURLDownload.h](https://github.com/xybp888/iOS-SDKs/blob/4a48e5409f8ea6ac422c44d91b86fa0f8ae010ff/iPhoneOS27.0.sdk/System/Library/Frameworks/BackgroundAssets.framework/Headers/BAURLDownload.h)

## On-Demand Resources: no local registration path established

`NSBundleResourceRequest` requests existing resource tags from a bundle. The
SDK accepts an explicit bundle argument; it does not document a general
file-import/asset-pack-registration method on that class. The iOS 27 header
deprecates ODR in favor of Background Assets; deprecation is not removal.
Apple's swift-build tests show tag maps, pack manifests and hosting URLs produced
during the build. This is evidence about those build outputs, not a proof that
every runtime registration technique is impossible.

A GitHub lead,
[FlekDeck's ODR guest hooks](https://github.com/flekstore/FlekDeck/blob/6ba96e1181ffe6c21e70ecdf71ae7027256a5e8d/TweakLoader/ODR+GuestHooks.m),
does **not** register imported data with the system. It serves already-unpacked
files locally, replaces NSBundle resource lookup and skips ondemandd. Its own
comments attribute real registration to installd. Those comments are a
third-party observation, not independently verified iOS 27 service behavior.
Copying its hooks would leave Crumb's underlying file-location problem intact.

References:
[SDK NSBundle.h](https://github.com/xybp888/iOS-SDKs/blob/4a48e5409f8ea6ac422c44d91b86fa0f8ae010ff/iPhoneOS27.0.sdk/System/Library/Frameworks/Foundation.framework/Headers/NSBundle.h),
[Apple swift-build ODR tests](https://github.com/swiftlang/swift-build/blob/0d155c394d89ac27375a857b0e725211ae876812/Tests/SWBBuildSystemTests/OnDemandResourcesTests.swift).

## Reproduction and remaining boundary

Run `scripts/audit-fairplayd-contract.py --daemon <original fairplayd.H2>
--launchd <original launchd plist> --profile <matching decoded profile>`.
The script pins all three input hashes, parses the original Mach-O entitlement
and symbol tables, and enumerates the four direct entitlement-import calls.
Its anonymous output is committed in
[research/fairplayd-contract-24A446.json](research/fairplayd-contract-24A446.json).
The XML entitlements were compared to a separate extraction; changing one byte
in each of the three inputs was rejected. No phone service was executed by this
audit. The decoded profile is a derived artifact, and its hash is explicitly
separate from the original daemon and launchd hashes.

The unresolved system-path boundary is still the actual crypter/service failure
for the imported backing vnode. Neither a successful mmap, a generic ENOMEM,
nor the profile's default allow identifies that failure. The next native reverse
engineering step is to resolve the service dispatch and indirect callers above,
and determine whether a legitimate caller-supplied file/license handoff exists.
Its contract must include caller authorization and actual license consumption;
a service name or an untested path is insufficient. There is no new functional
decryption correction to put in a phone IPA from this investigation yet.
