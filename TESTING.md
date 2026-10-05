# Testing — 2.3.0-dev.6 / build 23006

Build/test success is not equivalent to Apple acceptance or installation.
No physical device is connected to this environment. Device results below
were reported by the user; untested cells remain pending.

## Device matrix

| iOS | Device | Login | 2FA | Search | Versions | Purchase | Download | Export |
|---|---|---|---|---|---|---|---|---|
| 27.0.1 | iPhone, model unspecified; ksign/certificate, no JIT | Confirmed by user, 23002 | Confirmed challenge → success, 23002 | Cancellation logs observed; dev.3 fix pending | Dev.3 failed at ent/download; dev.4 pending | Pending | Latest/old pending | Pending |
| 26.x | iPhone | Pending | Pending | Pending | Pending | Pending | Pending | Pending |
| 27.x | iPad | Pending | Pending | Pending | Pending | Pending | Pending | Pending |

The user confirmed reopening restores the 23002 session. This checks local
Keychain restoration, not indefinite token validity. Exact device model,
Release configuration and debugger status were not explicitly reported.
Logout, wrong/expired 2FA, account-disabled, session-expired and UI error cases
are not claimed physically tested. See docs/evidence/ios27-authentication-user-report.md.

## Automated checks

```sh
swift test --package-path MapleSyrup/SAPKit           # macOS + real Keychain fixtures
go -C MapleSyrup/NativeSAP/packageipa test -race ./... # Go 1.25+
bash scripts/test-no-exec.sh                         # Linux stock Unicorn rejection
bash scripts/test-tci-host.sh                        # Linux TCI with PROT_EXEC denied
bash scripts/test-sap-host.sh                        # MANUAL network: SAP and synthetic kbsync
```

The manual network smoke uses no credentials, purchase or download request.
The kbsync host test produced 196 bytes for synthetic DSID 1 with the TCI runtime.
It was not run under seccomp; the no-executable-memory TCI test is separate.
Generation does not prove Apple acceptance of that blob. Test fixtures never
claim cryptographic validity of fake signatures/blobs.

Swift fixtures cover Bag, exact auth plist/signature bytes, 2FA/cookies, pod
redirects, status/retries/timeouts/cancellation, safe messages and Keychain;
Store adds account country, string/numeric IDs, strict app/version response
matching, CDN URL restrictions, cached kbsync rejection/regeneration, accepted-only
cache, pinned redownload → empty500 → Bag update, free-license Bag/pod purchase,
paid purchase refusal, Retry-After dates and DSID/GUID-bound real Keychain cache.
Go fixtures cover ZIP payload preservation, SINF replacement without duplicate
entries, actual Info.plist version, wrong bundle/version/MD5/path rejection,
missing sinfs and HTTP Range status/content-range/length enforcement.
CI builds both unsigned Debug/Release arm64 IPAs with Xcode 26 on macOS.
The complete app is not configured for Simulator. Build evidence is in GitHub Actions.

## Install dev.5 and test downloads on iOS 27

1. Use **Release 23005 / v2.3.0-dev.5**. Sign using the same certificate and
   bundle ID as 23002 to retain Keychain identity/session where the signing
   method supports it. Install with ksign, SideStore, AltStore, Sideloadly or a
   developer certificate, without JIT/debugger/TrollStore. Record device model,
   exact iOS, method and Release build number. Do not share credentials.
2. Reopen. If the saved session is inaccessible under the new signature, sign
   in normally. Enter a six-digit current trusted-device code only when asked.
   Do not manually append it to the password. Retry a transient 204/404 after
   waiting; do not repeatedly sign in to provoke Apple rate limits.
3. Start with a small **free app already owned** by your account. Search or
   paste its App Store URL/numeric ID/bundle ID. Tap **Choose version / download
   IPA**. Wait for latest externalVersionId, Bag and kbsync. Cold assets and
   interpreted guest work may take minutes; keep the app foregrounded.
4. Choose **Latest iOS build**. Expect a confirmed Info.plist version when CDN
   Range is supported. If it is unavailable, the dialog says so and download
   verifies the version afterwards. Tap Download IPA. Progress must report
   bytes received, then a ZIP/identity/purchase-data validation stage.
5. Success must say **IPA verified and saved**. App Info shows the actual
   bundle/version. Export IPA via ShareLink/menu or Settings → Downloaded IPAs.
   Use Save to Files; check the ZIP has Payload/<app>.app/Info.plist, matching
   bundle/version, iTunesMetadata and SINF when supplied. No installer is launched.
6. Reopen WaffleStore; the completed download must remain available. Also check
   Files → On My iPhone → WaffleStore → Downloads. Auto-Clean does not delete
   completed files. Verify the JSON sidecar maps ID↔actual version without
   DSID/token/kbsync/URL. IPA metadata contains the license Apple ID: keep it personal.
7. Repeat for an **older externalVersionId** returned by Apple. Confirm that
   the final version is older and that the correct ID is in the sidecar. If Apple
   does not list old IDs, enter a known ID manually; availability is not guaranteed.
8. Test a not-yet-owned free app: a required license may trigger Bag buyProduct
   with price=0. Test paid/unknown-price apps only to verify refusal to auto-buy;
   acquire licenses manually in App Store. No Arcade/subscription/payment automation.
9. Cancel during CDN transfer: no new completed IPA/history should appear.
   Keep an earlier download to confirm it is not overwritten. Native kbsync,
   Range inspection and ZIP preparation are bounded blocking calls; cancellation
   takes effect after they return and before a final file is committed.
10. If a step fails, Settings → **Export Store/download diagnostic** plus the
    visible sanitized error. Include app ID and selected externalVersionId if
    relevant. Never send HTTP bodies, signed URLs, password, cookies, token,
    signature, codes or kbsync. Native stages: 2 assets, 8 kbsync, 22 ranges,
    31 ZIP/limits, 32 MD5, 33 bundle/platform, 34 metadata, 35 SINF, 36 CRC, 37 storage.
11. Logout with no active operation; expect account/cache removed and identity
    retained. Sign in again. Session-invalid errors ask for fresh login because
    no password is persisted for unattended refresh.

## Manual edge cases and limits

| Case | Expected behavior | Evidence |
|---|---|---|
| HTTP204/404/429/5xx or timeout | Bounded retries; honor Retry-After; useful status | Auth/Store fixtures; real initial auth404/204 reported |
| Cached kbsync rejected | Clear cached blob; generate fresh once; preserve requested ID | Store fixture; physical pending |
| ent unavailable | Validated pod/Bag fallbacks; never silently switch selected ID | Store fixtures; physical pending |
| CDN redirects | HTTPS Apple CDN only, max8, no account credentials | Policy/range fixtures; physical pending |
| Wrong app/version/platform/checksum | Refuse completed IPA | Swift/Go fixtures; real rejection pending |
| Missing sinfs | Preserve archive without injecting imaginary license | Go fixture; physical pending |
| Expired CDN URL | Explicit failure; repeat selection to get fresh descriptor | Implemented; physical pending |
| App background/termination | Foreground session may pause/expire; repeat unfinished operation | Physical pending; no background/resume claim |
| Low disk/RAM/jetsam | No false success; completed prior files remain | Physical pending |
| Re-sign with different access group | Keychain may be inaccessible; fresh sign-in needed | Physical pending |
| Share Sheet / iPad popover | ShareLink works; Files export persists | Physical pending |
| Install/downgrade IPA | Depends on signing/DRM/iOS/receiving sideloader | Not implemented or claimed |

The App Store IPA may remain encrypted. Successfully sideloading WaffleStore
is not proof that SideStore/AltStore can install a protected downloaded IPA.
No installation-success entry is fabricated in the original downgrade history.

## Dev.4 regression checks on device

1. Upgrade using the same signing identity and bundle ID; reopen and keep the
   restored session. Do not logout as the first versions troubleshooting step.
2. Request latest versions for an app already owned. If ent fails, diagnostic
   must show recovery=ent-to-pod and the pod HTTP/failure result. Try a historical
   externalVersionId, download and export if versions succeeds.
3. If versions fails, copy Store probe v4. It should contain HTTP and numeric
   apple-failure or absent, and a stable category; no token or cookie values.
4. Logout only when deliberately testing login or a confirmed expired session.
   Login can succeed without 2FA when Apple returns valid session fields. If
   Apple challenges, enter its actual code; never append a guessed/stale code.
5. For continued login 204/404, copy the authentication v4 diagnostic. Separate
   URLSessions reduce connection reuse; acceptance still requires device proof.

New local HTTP fixtures check bare ent 401 → pod success, structured ent session
error → pod success, pod expiry preservation, bare pod 401 classification,
diagnostic redaction, ephemeral cookies across isolated URLSessions, and bare
Document login requiring a complete session. These are protocol regressions,
not tests against Apple's servers.

## Dev.5 device validation (pending)

Build 23004 was tested by the user: login succeeds after transient HTML/empty
HTTP failures, without another 2FA challenge; ent returns empty HTTP 401, pod
returns plist failureType 2042. Versions still does not work. Dev.5 corrects
the serial bytes and makes restored cookie application explicit; no Apple
acceptance is claimed from fixtures.

Upgrade to Release 23005 with the same certificate/bundle ID and retained account.
Request versions first. Copy Store probe v5 if it fails: request-cookie-count
shows whether each destination actually has applicable cookies. A zero count
can be legitimate for host/path-scoped cookies; do not broaden their domain.
If a fresh login is necessary, copy authentication probe v5 as well, including
cookie counts and redirect Location presence (values are withheld). Do not repeat
logout automatically for 2042. Continue latest/old download and Share Sheet
validation only if a valid descriptor is returned.

## Dev.6 physical-device checks (pending)

User report for 23005: fresh logout/login eventually allows versions and IPA
export; 20–30 manual login attempts, version labels/popup poor, installation
untested. Device model remains unspecified. Dev.6 must validate:

1. Fresh sign-in: one tap, automatic attempt counter in diagnostic, cookies
   retained, Cancel works, correct password/2FA stops appropriately; record the
   number of HTTP attempts. Do not claim backend reliability from fixtures.
2. Scroll to old entries: visible version numbers load gradually from IPA
   Info.plist; select a low row, scroll, review sheet remains stable, close/retry.
3. Download an old ID; verify Downloaded apps shows visible version/build and
   selected ID. If range inspection succeeded, final version must match it.
4. Open menu → Downloaded apps → Install this downloaded version → Request iOS
   installation. Keep Safari screen open and confirm iOS. Record the exact iOS
   result (success/refusal), not merely Safari launch or downloaded bytes.
5. Check installed Settings/App version when accepted. Reopen WaffleStore and
   export still works. Server Close/timeouts stop transfer without deleting IPA.
6. Repeat on iOS 26/27, with an owned free app. No paid purchase automation,
   decryption, sideloader capability or App Store downgrade acceptance is assumed.

Generator smoke: HTTP 200 and valid OTA plist for com.example.fixture, build 1,
loopback URL. No Apple login/IPA was uploaded. This checks generator format only.
