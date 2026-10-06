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
(at most eleven transient retries shared across routes/logical attempts, with one 120-second scheduling window per UI login/verification), separate ephemeral
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

## Build 23012: five tabs and catalog region

The root uses a native SwiftUI TabView: Search, Downloads, Favourites (center),
Apple Account and Settings. Search works before sign-in. Choosing a result, link,
ID or bundle ID opens the existing version sheet; if signed out it routes to
Account and resumes the selection after successful sign-in. Favourites stay
available offline; adding/removing uses the star on the resolved app summary.
Downloads keeps installation/export/deletion and now contains progress, cancel,
history and the optional full activity reader. No hamburger or overlay action
stack competes with the tab bar. Settings and credits retain their original content.
The native bar adapts to the OS, device and accessibility size; this is not a
pixel-for-pixel custom copy of DLipa. Login/SAP, license acquisition, version
selection, IPA verification and installation protocols are unchanged from 23011.

A globe in Search selects a public catalog country. The supported codes reuse
the MIT-attributed ipatool storefront mapping; only this non-secret preference is
stored in UserDefaults. Default is the authenticated account country (US when
signed out). Region names are localized and searchable. Stale searches are
cancelled when changing query/catalog or leaving Search; cancellation is not an
error and HTTP search failure is shown with retry instead of an empty result.

This is catalog browsing, not account-region migration. The current ipatool
search derives country from the authenticated storefront and offers no search
country override. Purchase/download eligibility still uses the saved Apple
account storefront and existing StoreSession. Selecting a foreign catalog never
rewrites that storefront, token, cookies or identity. A foreign-only app may be
unavailable to this account. Downloading it requires access accepted by Apple
with an eligible account; changing headers cannot establish that entitlement.
No automatic Apple country change or cross-region license bypass is implemented.

The user reports 23011 works well on their device. This supports preserving its
backend, but is not evidence of 23012 navigation or foreign-region downloads.

## Build 23013: correct the cross-region conclusion and metadata path

The 23012 limitation was overstated: an account storefront does not prove that
all apps found in another country are impossible to download. 23012 changed
search country but StoreSession.lookup and latestVersion forced account country,
so a foreign-only listing could fail before any authenticated download attempt.
23013 carries the selected country through StoreApp to both catalog lookups and
externalVersionId resolution, including version-label inspection and actual
selected-version download. Authentication credentials/storefront are unchanged;
the existing Apple download/license services decide eligibility. No speculative
storefront-header substitution is needed for this correction.

DLiPA evidence inspected October 5, 2026: repository AhmedBafkir/DLiPA commit
8bc61ee9d4642f2706b99e9012425c58077015b7 declares "Select an available storefronts"
but contains README/images, not application source. The published v1.4 IPA SHA256
is 72012beabf55a12d70c5ad3bbe95384e7fd450c8de1bf6fbc62f8844512eea0f.
Read-only inspection of its Objective-C metadata/call sites shows its selected
storefront stored separately under Storefronts_selected (selection callback at
0x100023958); search uses searchForApp:countryCode:limit:completion:. Download
uses downloadWithAppId:specificVer:completion: and the MZFinance volume endpoint,
with account DSID headers. Its purchase request obtains the storefront from
APStoreAccount.sharedInstance (0x10004a854 through 0x10004a888), not the selected
search preference. This supports trying a foreign catalog app with an existing
account; it is not proof that Apple accepts every cross-region acquisition.
No DLiPA binary/source code is copied into WaffleStore or redistributed.

Favourites optionally persist the catalog country; old records remain decodable
and use account country. Login handoff preserves the pending selection country.
Invalid catalog overrides fail before network access. Default requests retain
account-country behavior. Free acquisition remains once per operation after
license denial; paid/unknown-price apps are never automatically purchased.
These fixes allow genuine attempts, not invented licenses or claimed live Apple
acceptance. Real cross-region download/install remains a device test.

## Build 23014: explicit purchase-2059 recovery and centered confirmations

A physical 23013 report for Venmo returned Apple failureType 2059 and "Purchase
could not be completed." That is an explicit acquisition rejection, not proof
that all regional downloads are forbidden. It exposed a missing reference flow:
ipatool's purchaseWithParams/Purchase maps 2059 to ErrTemporarilyUnavailable and
tries pricingParameters GAME after STDQ for iOS. DLiPA v1.4 does the same: its
purchase branch compares 2059 at 0x10004abd8, then replaces pricingParameters with
GAME at 0x10004abfc–0x10004ac08, and reports failure if rejected again.

23014 implements this bounded reference behavior: after a missing license, a
verified free app receives STDQ; only an explicit 2059 enables one GAME request.
The price remains zero, with the same app ID, GUID, account token/storefront and
pod. This is not an account-country change, subscription purchase or entitlement
bypass. Another 2059, subscription/terms/session denial, timeout or ambiguous HTTP
failure is not retried through GAME. Ambiguous purchase errors still require the
existing read-only access verification; success still requires a validated Apple
reply and the requested externalVersionId. Paid/unknown-price apps remain blocked
from automatic acquisition. The user's live Venmo acceptance is still pending.

The 23013 reference investigation missed this 2059 branch; its conclusion about
DLiPA/ipatool acquisition parity was incomplete. The new tests cover successful
alternate acquisition, repeated denial with no loop, and subscription rejection.
The existing no-replay fixture now uses a non-2059 Apple denial, since 2059 has a
specific alternate-request meaning in both references.

The supplied sign-out/delete screenshots show iOS 27 adapting confirmationDialog
into a popover anchored far above the action. Both use native alert instead, with
an explicit Cancel and destructive button and shorter copy. Alert presentation
is centered and does not depend on row/menu geometry. No custom glass UI or API
is added. Physical-device appearance remains a manual verification.

## Build 23015: deeper DLiPA parity and download/installation presentation

23014 still returned 2059 in the user's device test. The completed caller/body
inspection is in docs/DLIPA_BINARY_AUDIT.md; the earlier pricing fallback alone
was insufficient. Regional acquisition now uses the nine-field/string-ID/2.20
request profile observed in DLiPA's binary. The default ipatool purchase profile
and login/SAP are preserved. This is a compatibility attempt whose live result
must be checked, not a confirmed cross-region fix. No extra replay or unvalidated
Location forwarding is introduced. A same-app/same-account DLiPA comparison is
needed to distinguish client differences from current Apple eligibility.

Search result labels now fill their row and use a Rectangle contentShape, so
blank space/icon/text all select the app. The active download section becomes one
VStack row: app name, short phase, actual progress/bytes/percentage and accessible
cancel control. This avoids List inserting dividers between every progress item.
Saved IPAs use separate sections/cards. No invented throughput/completion is shown.

The installation page now uses a responsive system-font card, actual app/version,
large Install button and adaptive dark/light colors. Metadata is HTML-escaped and
the itms-services target remains JSON-escaped and unchanged. The automatic handoff,
manual retry button, original HTTPS manifest service, loopback-only server and IPA
serving remain intact. Safari retains its native controls with adaptive tint; the
outer title is inline to avoid double large headers. Displaying the page still
never claims installation success. New template tests cover metadata injection
and target preservation. Visual/tap/real transfer checks require device validation.

## Build 23016 — remove unused history and confirm before deleting

The user verified 23015 search/download/export/install/UI on device, but regional acquisition still returns Apple 2059 (the reported number varies in subsequent messages). No further auth or purchase changes are made. Remove the unpopulated Downgrade History screen, state/model and JSON import/export controls; existing downloaded IPA records remain the download list. Historical UserDefaults bytes are left untouched. A destructive-role swipe button caused SwiftUI to optimistically remove a row before the alert, then restore it. Use a red-tinted ordinary swipe action to request confirmation; only a successful confirmed file deletion changes the list inside a short animation, respecting Reduce Motion. Native animation still needs physical-device validation.

## Build 23017

Settings no longer lists downloaded IPAs or global Export IPA/Clean Documents actions. Individual export remains in Downloads, which now owns a confirmed Delete all downloads action. Bulk cleanup only removes downloaded IPAs/matching sidecars inside Downloads, handles orphan IPAs/stale records, preserves unrelated files and rejects a symlink root. Disable bulk deletion during download; reload records even after partial failure. No broad Documents wipe.

## Build 23018 — distribution UI without debugging

Release has no Activity log UI, stdout pipe/capture, sign-in attempt counter or installed diagnostic callbacks. App console print statements compile only in Debug; SDK diagnostic hooks remain default no-ops for tests and reuse. Useful progress/cancel/errors, Keychain, SAP and Store recovery behavior are unchanged. Prebranding branch preserves current WaffleStore bundle ID/icon until a name is chosen. Public renaming/redistribution requires resolving missing upstream license; notices and corresponding sources remain intact.

## Crumb 1.0.0 (23019)

Crumb 1.0.0 (23019): display/product Crumb, bundle com.certlium.crumb, crumb URL scheme and new Keychain namespace. Original project, scheme, source tree and history retained; no new repo created. Exact approved icon exported 1024px RGB, default/dark; no alpha/pre-rounded mask. Existing WaffleStore can coexist; credentials/downloads/favourites do not migrate automatically. Native auth/SAP/Store behavior unchanged. CI names/product and artifact verifier updated. Credits/notices/AI disclosure retained; original installation provider is still identified honestly. Public redistribution permission remains unresolved.

## Crumb 1.0.0 (23020) — 2FA renewal and incomplete redirects

Crumb branding did not change the authentication request/SAP state machine (only independent Keychain namespace). The UI screenshot shows invalidRedirect, not proof of a bad 2FA code; older device reports included 301 HTML without Location. Distinguish incomplete allowed redirects from unsafe destinations: retry only the existing validated endpoint within existing attempts/deadline, fresh signature each request, preserve strict host/path/TLS checks. Exhaustion returns redirectUnavailable with useful message; retain same-account prepared signer for retry. An account-bound in-memory challenge snapshot keeps verification/resend independent of SecureField rendering; it is discarded on cancel/account change/logout/success, never serialized. New code action keeps the password in RAM only, clears previous code/challenge cookies/preparation and sends a fresh signed password-only login; Apple controls code delivery. 30-second user-action cooldown, cancellation/account change/success cleanup. No new SMS/private endpoint or guaranteed resend claim. Full-width primary/secondary 2FA controls and Verify code label. Added regression fixtures for 2FA missing/blank Location, bounds/deadline, unsafe destination and challenge renewal after rejection. Release debug remains absent. Physical first-login/new-code validation is pending.

## Crumb repository migration

Repository: https://github.com/boierito/Crumb. Merged its initial README commit with the full production 23020 history; no force replacement of the initial commit. README/current app source link adapted to Crumb. Removed the original author’s release-dispatch notification workflow; builds/tests remain repository-relative and require no custom secrets for ordinary artifacts. Upstream notices and redistribution status retained; no public binary release is created by this migration. Protocol/authentication/2FA unchanged from 23020.

## Crumb 1.0.0 (23021) — canonical initial authentication and shared recovery

Re-audited current ipatool `cde7d00355e152714377b953ec57438626d3cb5a`; the previous comparison missed `f9aa653` (October 5), which canonicalizes the initial Bag authentication URL to authenticate/ after strict validation, preserving the host and encoded query. Ported into SAPConfiguration/AuthenticationEndpoint; legitimate canonical/pod Bag URLs are accepted. HTTP path validation uses URLComponents.percentEncodedPath rather than comparing it to URL.path (which can strip a trailing slash). Unsafe/encoded paths remain rejected before normalization; Apple redirects retain their exact URL and POST body.

Factored the existing deadline/retry allowance into one AuthenticationRecoveryBudget per user login/verification. Changing pods or logical attempt does not refill automatic retries; HTTP/network errors share the allowance, normal routing/2FA incur no sleep. Reference non-automatic retry behavior and backoff/Retry-After remain; fractional remaining time no longer rounds up. Native SAP/TCI, user agent, cookie isolation, prepared signer reuse, RAM-only credentials and Release logging policy unchanged. Nine regressions added; see docs/LOGIN_AUDIT.md for trace arithmetic and precise limits. Live Apple login latency and first-login/new-code UI still require physical-device testing.

## Crumb 1.0.0 (23022) — typed sign-in and 2FA recovery

Crumb is an adaptation, not a one-to-one runtime/wire-byte copy of ipatool. The
current email authentication semantics follow upstream cde7d0; the jailed TCI,
URLSession, stable Keychain identity, password-free persistence and prepared
challenge reuse differ. docs/LOGIN_AUDIT.md records these differences explicitly.

The controller previously cleared the code on every error, even when Apple had
not rejected verification. AppleSignInForm replaces scattered pending credentials,
cookies and challenge flag with one account-bound RAM state. UI validation and
immutable submission share the same rules, preventing blank/malformed challenged
input from falling back to password-only login. Transient failures preserve the
code/challenge and valid SAP preparation; Retry verification resubmits the same
code. Only explicit verification rejection clears the code; account rejection
unlocks credentials. New-code delivery remains an explicit action with cooldown,
not an automatic reaction to a network error. Cancellation/account change/success
clear the challenge and form credentials; nothing new is persisted or logged.

Twelve state tests and a full form-to-protocol fixture cover these rules, including
one challenge followed by a bounded HTTP failure and a successful retry using the
same password+code body/pod. The existing wire protocol, recovery scheduling and
SAP runtime are unchanged. CI compilation/fixtures cannot establish Apple latency,
code expiry or first-login acceptance; TESTING.md includes the device checks.

## Branch-only Crumb Test 23023

feature/2fa-clean-login-test derives from main 01ad52d without modifying its
checkout/ref. Test bundle/service/URL scheme are separate. The explicit local
reset closes prepared guest/transports, blocks overlapping login/Store work,
clears only test account/kbsync/identity, rotates identity and verifies read-back
and absence of session records. It does not request an Apple challenge itself.
Guest assets remain cached; authentication protocol/runtime/retries are unchanged.
Branch-only diagnostic UI accepts fixed fields, never raw secret-bearing data,
and reports actual challenge/verification/restoration/no-challenge outcomes.
See docs/AUTH_TEST_BUILD.md for installation and physical limitations. Do not
merge test reset/instrumentation into main as a production authentication fix.
