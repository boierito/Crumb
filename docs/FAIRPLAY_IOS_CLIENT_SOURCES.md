# Matching iOS producer references

Research continued 2026-10-31. No device service or user license is executed by
this source inventory. It does not establish a working imported-IPA decryptor.

## Original standalone clients acquired

Both binaries were selectively read from the official iPhone17,1 / 24A446
BuildManifest OS filesystem component using the existing AEA/APFS reader.
Downloaded segments passed that reader's HMAC/SHA-256 checks; the complete
archive was not downloaded or verified. No user capture, account or IPA was
used. Hashes identify the original standalone executables, not normalized
shared-cache exports.

| Original filesystem path | SHA-256 |
| --- | --- |
| `/System/Library/PrivateFrameworks/AppStoreDaemon.framework/Support/appstored` | `4a9a5addf872e5ec6d9034e8a33d533890e325ce072153da5b636a97ed5a8636` |
| `/System/Library/PrivateFrameworks/iTunesStore.framework/Support/itunesstored` | `548ae647718f1c00839e4c9337b889bb0b159969db2255ba5694d07adf1e67b4` |

The model/OS association comes from the supplied phone log and firmware catalog;
the phone's exact installed firmware build was not measured by this inventory.

`itunesstored` imports `SSVFairPlaySAPSession`,
`SSVFairPlaySubscriptionController` and `SSVFairPlayGetHardwareInfo` from
**StoreServices**. These imports identify a real native producer reference.
They do not establish that these operations are available to Crumb or return
the material needed for the original mobile SINF.

`FairPlayDecryptSession` and the `dpInfo` processing classes are also present.
Their names do not establish Mach-O page decryption: package/stream decryption
must remain separate from consumption of a mobile executable's SINF.

## Concrete AppStoreDaemon broker boundary

The original `appstored` contains `ASDFairPlayServiceProtocol`, with:

- `generateKeybagRequestForDSID:completionHandler:`;
- `generateSubscriptionRequestForDSID:completionHandler:`;
- `importKeybag:completionHandler:`;
- `importSubscriptionKeybag:completionHandler:`.

The implementing Swift class is `_TtC9appstored15FairPlayService`. Its
`+entitlement` method at `0x1003c2578` constructs the small Swift string
**`FairPlay`** and bridges it to Objective-C. This is a service label; it must
not be mistaken for an entitlement key that grants access by adding that spelling
to Crumb.

Native `-[XPCServiceClient getFairPlayServiceWithReplyHandler:]`, at
`0x1002842c4`, calls the common service dispatcher at `0x100005c78`.
The dispatcher calls predicate `0x100005164` at `0x100005cf0`. A false result
branches at `0x100005cf4` to the error construction path. That path uses the
original constant string **`Not entitled for this service.`** and error code
`0x1f9` (505). No callback with a service is issued by that branch.

The predicate consults an entitlement-array membership result, an optional
service-specific Boolean entitlement, and a legacy Boolean entitlement.
The array's producer at `0x100005214` obtains the caller SecTask from an audit
token, then reads the original CFString at `0x100693738`:
**`com.apple.private.appstored`**. The Boolean fallback CFString at
`0x100691c18` is **`com.apple.itunesstored.private`**. The Boolean reader
`0x10000567c` calls `SecTaskCopyValueForEntitlement`, checks the returned type
and obtains its Boolean value. A missing value follows its false path.

These observations establish a **process-entitlement boundary on this broker**.
An Apple Account password or a successful HTTP download does not supply the
caller SecTask entitlements consumed by this code. This investigation does not
measure the installed Crumb process's entitlements, reproduce the complete
service catalog or grant either entitlement. It also does not prove every
lower-level StoreServices/FairPlay client operation has the same boundary.
No broker invocation or guessed entitlement is added to Crumb.

## Acquiring the actual request encoder

The remaining request-1000 in/out-word producer is described in
[the native dispatch investigation](FAIRPLAY_MIG_DISPATCH_RESEARCH.md).
Its valid value cannot be inferred from the rejected placeholder's bit test.

The separately pinned public Cryptex component is
`043-69419-713.dmg.aea`, 2,348,810,240 bytes, SHA-256
`23bd07e51c154ea4c1afd5f59c24e05ddeab2b6d0cbd380d4623c13fb5a93091`.
`scripts/acquire-fairplay-public-client.sh` runs only in a development macOS
runner. It verifies the pinned extractor and component, opens the public
filesystem read-only, and extracts selected CoreML, MediaToolbox and StoreServices
images with cache-slide normalization and their local symbols. The separate
DriverKit cache is excluded. No user inputs are sent to the runner.

Global stub-island indexing is deliberately omitted from this acquisition.
The inspected extractor already exports each requested image's local symbols;
its optional global index attempted to serialize over 21 million OS symbols,
which is unnecessary for obtaining these three references. This change affects
development acquisition, not Crumb's runtime or its decryption checks.

A completed acquisition must contain `acquisition.json` with image hashes.
Partial CI artifacts without that manifest are not treated as acquired clients.
Cache-slide-normalized exports are explicitly distinguished from the untouched
standalone binaries above. Obtaining the reference is not validation of a SINF.


The acquisition completed in [run 41148283988](https://github.com/boierito/Crumb/actions/runs/41148283988)
from commit `3ac7c76826e49e2a250632a068f4ca6d566e44eb`. The downloaded
artifact ZIP's SHA-256 was
`9e2b85243ea9db6470550e67fcd054f928a84a09dcc81e7778b2138fa68abcb4`.
Each extracted image was checked against the manifest:

| Normalized image | Bytes | SHA-256 |
| --- | ---: | --- |
| CoreML | 11164632 | `b820c809dd4cb7925a31846fa600bd59aad94c8a58c7eb03dc6f2490c15c2034` |
| MediaToolbox | 18919608 | `cd9dec1e47018cd539dc30eb96b82b2e823256ddcaacd816fc77a04f5cc5b611` |
| StoreServices | 4784528 | `e388e6f1bab28b88f36760f1aed899179cb3bc6181628699e1a38f6955d3dca5` |

## StoreServices producer data flow

Original local symbols identify the following wrapper bodies. These addresses
refer to the normalized 24A446 StoreServices image, not Crumb code.

| Wrapper | Address | Observed internal operation |
| --- | --- | --- |
| `SSVFairPlayContextIdentifier` | `0x1cfa1dbec` | Hardware-info wrapper and internal `0x1cf90b2e8`; cached 32-bit context |
| `SSVFairPlayGetHardwareInfo` | `0x1cfa1dd1c` | Requests the original string `UniqueDeviceID`; external call target remains to be resolved |
| `SSVFairPlayCopyKeyBagSyncData` | `0x1cfa1e140` | Context helper `0x1cfa1e1d0`, then internal `0x1cf9547ec` |
| `SSVFairPlayImportKeybagData` | `0x1cfa1e230` | Same context helper, NSData bytes/length, then internal `0x1cf95cc0c` |

The import wrapper returns true for nil input without invoking the actual
import; this is not evidence of importing a license. For nonnil input it tests
the internal result and processes returned 32-byte records. The record layout
has not been decoded and does not establish that a 16-byte account credential
is present. The 32-bit context identifier has not been connected to the MIG
request-1000 in/out word; equal widths do not establish that connection.

The exported image does not include the shared external stub islands. Guessing
external function names from argument shapes is insufficient. The acquisition
now collects seven bounded original call stubs, their referenced raw pointers,
slide-decoded targets and independent image-scoped symbol lookups. Unknown
stub shapes remain unresolved. This avoids the whole-OS symbol index and does
not fabricate external memory or service replies.

## Installed-file signature observation

In the FairPlayed reference `fp_attempt_describe`, `sig=-` means `addsigs` was
false. `sig=ok` is emitted only for an addsigs attempt with no recorded error.
The successful installed-target log therefore does not require that attempt's
explicit `F_ADDFILESIGS_RETURN` call. It does not establish the absence of
preexisting vnode/signature context created during installation, or that
registering an original signature would repair a data-container remap. No
runtime change or cause attribution follows from this observation alone.
