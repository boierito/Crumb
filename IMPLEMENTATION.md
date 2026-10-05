# WaffleStore login, versions, download and installation repair

This contribution starts from nxtcoreee3/WaffleStore d508e53. The original
project, UI, search, favorites, history and sideloading model are retained.
Code generation was performed with OpenAI Codex (**AI slop**); boierito
coordinated the work and reported the physical-device results in TESTING.md.

## Backend

The old MapleSyrup authentication no longer served the reported Apple flow.
The replacement follows majd/ipatool commit
3411d57f451f5111ae115641c22f7ed17bbd5fbe, fetched on October 5, 2026:

UI → MapleSAP → dynamic Store Bag → SAP initialization → exact signed request
→ Apple authentication/validated pod redirects → 2FA if requested → Keychain
session → app lookup/externalVersionId → Store descriptor → CDN → verified IPA.

The Bag supplies SAP endpoints, authentication URL and Store download/purchase
endpoints. Redirect destinations are validated before any credential replay.
Native XML/Document/Protocol replies are normalized; account/2FA errors remain
specific. Temporary HTML/empty HTTP failures use bounded automatic recovery
(12 HTTP attempts per endpoint, with one 120-second window per UI login), separate ephemeral
URLSessions and preserved cookies. Retry-After is respected, cancellation is
available, and the single -5000 logical retry follows ipatool. Signing and
redirects share that window; a slow signature cannot reset the network timeout.
The app retains a prepared SAP guest, ephemeral cookies and validated pod for
up to five minutes for the same account's retry/2FA. Each request is signed
freshly. This preparation is memory-only, contains no password/code, and closes
on success, cancellation, logout, account change, terminal error or expiration. No fake 2FA code
is supplied. Intermittent Apple login failures remain a known limitation.

DSID, passwordToken, storefront, pod, relevant cookies and stable six-byte GUID
are persisted in Keychain. Password/2FA input is never serialized. Logout clears
account/session/kbsync without regenerating machine identity. Cookies preserve
scope/path/Secure/expiry; legitimate Apple parent cookies are retained. Tokens
and account cookies are isolated from the CDN transport.

## SAP inside a jailed app

The current ipatool guest depends on Unicorn. Stock Unicorn generates executable
memory and cannot be the production jailed runtime. The build scripts restore
the QEMU TCI interpreter in a pinned Unicorn source tree, compiling it into a
static iOS library. Guest x86 instructions are interpreted rather than executed
as unsigned host code; no external JIT, attached debugger, private entitlements,
jailbreak or filesystem access outside the app container is required.

A Go/C ABI library provides SAP exchange/signing, StoreAgent kbsync and bounded
ZIP preparation/range inspection. It is not the ipatool CLI executable. Assets
use the reference loader/hash validation and app cache. Swift owns networking,
UI, Keychain and session state. Tests independently exercise TCI with executable
mappings denied. Memory/TCI smoke probes live only in tests/scripts, not the IPA;
the user-facing probe screens and exported tracing buffers have been removed.

## Versions, purchase and download

Search and lookup use the account's storefront/country. Current iOS version IDs
are resolved from MDM offers; historical IDs come from Apple's Store response.
The displayed number is read from the selected IPA's Info.plist through bounded
CDN ranges, with one worker for visible rows and selection priority. Missing
labels are identified as unavailable rather than invented. Confirmation uses
a separate sheet, avoiding the scrolled-row popup anchor. The range inspector
has an 8 MiB/20-second budget; the C ABI operation cannot be interrupted mid-call.

The preferred ent/download request uses account-bound kbsync and the exact
ipatool serial (five-byte prefix + hardwareID[2:]), X-Token/DSID/storefront and
Configurator 2.18 User-Agent. Other Store requests use the reference default UA.
Rejected cached kbsync is regenerated once. ent rejection tries the validated
pod before classifying the account; HTTP 401 alone does not prove token expiry.
The remaining Bag redownload/update fallbacks preserve the selected version.
An explicit pod SignInRequired response retains the account and shows a useful
error, without prescribing repeated logout as the only recovery.

Only verified free apps are acquired automatically. Paid apps/subscriptions
require the App Store. CDN retries, timeouts and HTTPS redirect validation are
separate from account authentication. Download has real byte progress. Package
preparation validates ZIP/CRC, provided MD5, app/bundle/platform and metadata's
externalVersionId, preserves Apple ZIP extras and places purchase metadata/SINF.
The actual version/build is read from Info.plist and compared with the inspected
version when available. Completed IPAs and nonsecret history live in Documents/
Downloads and can be exported through Files/Share Sheet. Auto-clean preserves them.

## Installation

Version review offers Download and install → verified IPA → Safari, without
visiting Downloaded apps. Download-only completion offers Install now/Export,
and the main screen offers Install latest download. Downloaded apps groups
technical details behind a disclosure and confirms deletion of the IPA plus
matching JSON sidecar. Other files, installed apps and app data are untouched.

The original public OTA mechanism is restored: Install →
Safari → itms-services → HTTPS manifest → one verified IPA on loopback.
The manifest generator remains api.palera.in. It receives only app metadata
and the loopback address, not account credentials or uploaded IPA bytes.

Telegraph binds to 127.0.0.1:9090, using randomized exact paths, read-only mapped
package data, validated single byte ranges and a private temporary file link.
Close/10-minute timeout/OS background expiry stops serving and retains the
original download. Only a permitted background task is requested. Opening the
URL or serving bytes never creates an installation-success history entry.
No installation/decryption/re-signing private API is introduced. iOS can refuse
FairPlay packages, licenses or older versions; export remains the fallback.

boierito reported successful installation with no visible errors on iOS 27.0.1.
Independent installed-version/data-preservation checks and iOS 26 remain untested.

## Build and IPA

Requires macOS, Xcode 26+, Go 1.25.1, CMake/Ninja/Python and network access for
pinned dependencies. scripts/build-unicorn-tci.sh and scripts/build-sap-native.sh
are invoked by the project's existing native build phase; CI provisions tools.

```sh
xcodebuild -project WaffleStore.xcodeproj -scheme WaffleStore \
  -configuration Release -sdk iphoneos -destination 'generic/platform=iOS' \
  -derivedDataPath build/DerivedData CODE_SIGNING_ALLOWED=NO \
  CODE_SIGNING_REQUIRED=NO CODE_SIGN_IDENTITY='' build
mkdir -p build/package/Payload
cp -R build/DerivedData/Build/Products/Release-iphoneos/WaffleStore.app build/package/Payload/
(cd build/package && zip -qry ../../WaffleStore-Release.ipa Payload)
python3 scripts/verify-ipa.py WaffleStore-Release.ipa
```

GitHub Actions builds Debug/Release unsigned IPAs on macOS and tests Swift,
Go package preparation and the no-executable-memory interpreter. Artifacts
include corresponding native sources and license notices. Version remains
numeric 2.3.0; the integration build is 23011. Sign with conventional sideloading,
keeping the same bundle ID/certificate when updating to preserve local state.

ipatool is MIT; Unicorn/TCI includes GPL-derived code. Preserve notices and
corresponding sources. Original WaffleStore has no explicit redistribution
license; upstream PR review does not establish permission to publish unrelated
redistributed releases. See THIRD_PARTY_NOTICES.md.

## Free-license preparation and UI recovery (23010, superseded below)

A verified zero catalog price now triggers `buyProduct` from the validated Bag
and authenticated pod before asking for version/download metadata. Purchase
always uses `appExtVrsId=0` / `price=0` / `STDQ`, matching ipatool's normal free
license request. Only `purchaseSuccess/status=0` or the explicit already-owned
code 5002 is accepted. The original requested externalVersionId is retained for
the following download. License success is memoized in the account-bound Store
session, not on disk; logout/recreating IPATool removes it. A cached license
rejected by Apple can be refreshed once, with no acquisition loop in one operation.
Paid/unknown-price apps are not automatically acquired; an owned paid app can
still download. Purchase failures retain Apple's sanitized information.

Previously only a final 9610 download failure initiated acquisition; intermediate
fallbacks could mask that signal. The preferred endpoint's structured 9610 now
propagates directly, and free apps do not depend on that signal to obtain a license.
There is no paid purchase, subscription acquisition, or replay after ambiguous
purchase failure. Device acceptance for previously unowned apps is pending.

The main terminal is collapsed by default; status/progress use compact rows and
user-facing wording. stdout is drained continuously into a bounded memory buffer
on MainActor, independent of the terminal view, so hiding it cannot fill the pipe.
Version/download modals use inline navigation titles and grouped rows. Technical
IDs are under Advanced/Details. Install/Export use the original PartyUI button
style, and downloaded-file deletion remains in the item menu/swipe with confirmation.
The clean integration line is based on build 23008; the optional investigation
branch/build 23009 remains separate. See docs/UI_REVIEW_23010.md for screenshot
findings and explicit visual/device validation gaps.

## Purchase recovery and readable logs (23011)

23010 made buyProduct a prerequisite for every known-free app. The user reported
that both unowned and owned apps failed at that step; no version metadata was
requested after a purchase rejection. 23011 restores download-first access checks.
Only a missing-license response (9610), or unavailable metadata for a verified
free app, leads to one acquisition per operation, followed by the original
selected externalVersionId. Paid/unknown-price apps are never acquired. The
accepted account-bound kbsync blob is retained after per-app license denial.

Verified reference differences: ipatool's purchaseRequest uses MZFinance buyProduct
with no guid URL query; the live anonymous Bag inspected on October 5 advertised
MZBuy buyProduct. 23010 reused that Bag route and appended a guid query. 23011
derives the MZFinance sibling from the strictly validated authentication URL and
authenticated pod. It retains GUID in the plist payload and sends the same
Configurator headers/body as ipatool. No fixed purchase host or speculative
redirect/endpoint fallback is added. The route/format mismatch is verified; live
Apple acceptance of its correction still requires device validation.

Purchase is not automatically replayed after ambiguous failure. Unlike ipatool's
HTTP-500-as-already-owned assumption, an empty purchase 500/incomplete plist
triggers one read-only descriptor verification. Access is accepted only if a
validated app/version response succeeds. Structured Apple denial, rate limits
and session/cancellation errors remain actionable. Successful descriptors for
owned apps never invoke purchase, so a broken purchase service cannot block them.

Fundamental purchase diagnostics now expose fixed route labels, HTTP/body type
and numeric failure codes, plus bounded error categories. They exclude body text,
URLs, credentials, tokens, signed CDN links and cookie values. No probe/debugging
screen is reintroduced. The Activity log remains optional, bounded and volatile.

PartyUI TerminalPlatter forces a 250-point frame; the extra inner 140-point frame
introduced in 23010 caused clipped/centered short content. Remove the inner frame,
use the original platter symmetrically, increase reading font and add Expand/Copy.
The expanded reader fills the sheet; automatic following is off by default so
new output cannot drag away the portion being read. Remove generic blue app.fill
placeholders. Use a spinner for active Store work and clear stale acquiring text
on failure. Screenshot/device visual validation is pending.
