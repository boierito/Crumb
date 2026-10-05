# WaffleStore original backend audit

Original: `nxtcoreee3/WaffleStore` / fork `boierito/WaffleStore`,
`d508e532e5b40d6470bab003a6cc0026d1ca7c75`.
Tag `upstream-wafflestore-2.2.2` marks that **original HEAD**, not a claim that
HEAD is identical to upstream's historical `WaffleStore2.2.2` release tag.

## Project inventory

One application target and scheme, `WaffleStore`; bundle ID
`com.nxtcoreee3.WaffleStore`; deployment target iOS 16.4; Swift language mode 5;
project objectVersion 77 with a synchronized `WaffleStore/` source group.
Original marketing version is 2.2, build 2. Project metadata says Xcode 27.0;
the original actually builds with Xcode 26.0 after resolving Swift 6.2 PartyUI.
No custom signing entitlements file is configured in the original project.
The Info.plist enables Files document sharing and the `wafflestore` URL scheme.
The `MapleSyrup/` directory originally contains a share-extension controller
and storyboard; it is not a separate backend Swift package or active target.

SwiftPM direct products: PartyUI 1.1.3, Telegraph 0.40.0, Zip 2.1.2.
Transitive pinned packages: CocoaAsyncSocket 7.6.5 and HTTPParserC 9.2.0.
The original `ipabuild.sh` builds unsigned using xcodebuild, packages Payload,
strips signing, and removes its build directory. The CI workflow packages the
same app without destructive cleanup and keeps build logs.

## Existing flow and source map

```text
ContentView / NavigationButtons / AppData
  → IPATool (WaffleStore/Functions/IPATool.swift)
  → StoreClient.authenticate / getBagEndpoint
  → Apple Bag → native auth → DSID + passwordToken + storefront + pod
  → volumeStoreDownloadProduct
  → softwareVersionExternalIdentifiers
  → manual external ID OR third-party visible-version lookup
  → download metadata + URL + sinfs
  → CDN file in tmp/app.ipa → Zip / SC_Info / iTunesMetadata
  → Downgrader → loopback server → Safari → itms-services manifest
  → Settings / AppMenu Share Sheet export
```

| Concern | Original implementation | Finding |
|---|---|---|
| Apple ID/password | StoreClient fields; NavigationButtons appends code to password | UI stores credential strings in memory; appended 2FA is later persisted as password |
| GUID | generateGuid hashes Apple ID with fixed CAFEBABE seed | Not the same bytes as a SAP hardware identity |
| Bag | getBagEndpoint fetches bag.xml, extracts only authenticateAccount | Ignores SAP version/setup/certificate; silently falls back to hardcoded auth URL |
| Authentication | JSON body with form content type; URLSession redirects | No SAP initialization, no signature, no signed plist body |
| Async completion | authenticate launches Task and returns false immediately | A successful server response still cannot make the synchronous UI return true |
| 2FA | Matches customerMessage containing Configurator_message | Force-casts customerMessage and many success fields; wrong/empty replies can crash |
| Credentials | dsPersonId/passwordToken; download-queue-info.dsid; headers | Force-casts a queue object that need not exist; storefront force-unwrapped |
| Pod/redirects | Recursive attemptGetPod after automatic URLSession redirects | Unbounded application recursion; POST/body preservation is not explicitly controlled |
| Session/cookies | URLSession.shared plus copied HTTPCookies | Shared cookie jar; no invalidation/refresh policy |
| Persistence | EncryptedKeychainWrapper; Documents/authinfo encrypted using ECIES | Payload includes password and tokens; key normally in Keychain/SEP, but has Library/.authkey plaintext private-key fallback |
| Logout | nuke removes file/key, sleeps and exits | Missing file removal force-try can crash; cookies and in-memory sessions are not explicitly cleared |
| Search/ID/bundle | AppSearchView / fetchAppNameAndBundleId / iTunes lookup | Public lookup is independent of the broken auth; UI gates many entry points on isAuthenticated |
| Favorites/history | Favourites / DowngradeHistory JSON | Local sandbox functionality; retain |
| Versions | getVersionIDList uses authenticated download metadata | External IDs, not display versions; dependent on login/license/download endpoint |
| Visible version mapping | apis.bilin.eu.org/history in Downgrader | Third-party mapping is not Apple's authoritative selected-version verification |
| Purchase | None | Account must already have a license; no modern free license acquisition |
| Download | Hardcoded p{pod}-buy / volumeStoreDownloadProduct | No Bag download endpoints, kbsync or modern ent/download preference |
| Progress/errors | downloadTask.progress polled using sleep; unconditional final success | Can report completion after transport/move failure; blocks callers |
| IPA preparation | Zip + iTunesMetadata + SC_Info sinfs | Force-tries manifest before fallback, unsafe array/path assumptions; encrypted App Store IPA remains encrypted |
| Installation | api.palera.in manifest + localhost HTTP server + Safari | Download is coupled to installation; serving a page is recorded as downgrade success without OS installation confirmation |
| Export | presentShareSheet from PartyUI, tmp/app.ipa | Export gated on hasAppBeenServed rather than a verified downloaded file |
| Logging | LogView copies redirected stdout | Original logged all auth header values, including X-Token, and signed CDN URLs; those two leaks are redacted in this branch |

## Where compatibility broke

The exact code mismatch is at creation of the native auth HTTP request:
StoreClient sends unsigned JSON; current ipatool signs the **serialized plist
body bytes** using a SAP session and sends Base64 in X-Apple-ActionSignature.
Bag discovery alone, a header constant, a made-up signature or a new URL cannot
replace this. Redirects must preserve the signed request payload. After login,
old download assumptions also omit the October kbsync runtime.

This conclusion comes from comparing implementations and the upstream status
notice. No real credential exchange was performed here, so we do not claim to
have captured a particular account's Apple rejection, identified the rollout
day, or proved that every unsigned request always fails identically.

## iOS 26/27

Project metadata and history refer to Xcode/iOS 27. The checked-in code does not
contain an active iOS 27-specific availability branch. The existing Safari
installation route is retained, including its commented iOS 18 rationale.
UIApplication.windows/first root controller selection is deprecated and may
misbehave with scenes; progress/AppData changes also cross thread boundaries.
None of these has been removed or presented as recovered installation support.
Real jailed installation, app lifecycle, ATS/loopback, background downloads,
Share Sheet and Keychain behavior remain device tests for both iOS 26 and 27.

## Build evidence

Linux failed before compilation: xcodebuild is unavailable.
First macOS run failed resolving PartyUI because Xcode 16.4 has Swift 6.1.
With Xcode 26.0 the unmodified original builds Debug and Release and both IPA
artifacts are uploaded:
https://github.com/boierito/wafflestore/actions/runs/37300555931
See `docs/evidence/original-build-linux.log`, `original-ci-xcode16.log` and
`original-ci-xcode26.log`. These are build/toolchain outcomes, not Apple backend
tests. No UI or functional source change was needed for the original build.
