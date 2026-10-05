# Testing

Build success is not Apple acceptance. No physical iOS device is connected to
this environment. Device outcomes below were reported by boierito.

| iOS | Device/signing | Login/2FA | Session reopen | Versions | IPA export | Install |
|---|---|---|---|---|---|---|
| 27.0.1 | iPhone, model unspecified; ksign/certificate, no JIT | Confirmed; intermittent HTML/empty failures still observed | Confirmed | Confirmed after renewed login | Confirmed | Reported working with no visible errors after build 23006 |
| 26.x | iPhone | Pending | Pending | Pending | Pending | Pending |
| 27.x | iPad | Pending | Pending | Pending | Pending | Pending |

Automatic recovery efficiency, all visible version labels/scroll cases,
independent installed-version identity and downgrade data retention still need
measurement. No universal app/device compatibility or end-to-end iOS 26 result
is claimed. Build 23008 adds direct installation/deletion and warm login preparation.
These changes require the device checks below; no measured reduction in Apple
HTTP failures or login time is claimed from host fixtures.

## Automated checks

```sh
swift test --package-path MapleSyrup/SAPKit
go -C MapleSyrup/NativeSAP/packageipa test -race ./...
bash scripts/test-no-exec.sh
bash scripts/test-tci-host.sh
```

Fixtures cover Bag/XML parsing, exact signed payload, credential/2FA handling,
cookies/pod redirects, secure persistence, bounded retries/Retry-After,
lookup/version IDs, ent recovery, free-license gating, ZIP/CRC/MD5/metadata/SINF,
actual IPA version/build, CDN ranges and OTA single ranges. Fake signatures and
kbsync fixtures do not prove cryptographic validity. Memory probes are test-only.
Host SAP smoke can be run separately with scripts/test-sap-host.sh; a synthetic
kbsync is not proof of an authenticated download. The HTTPS manifest generator
returned HTTP200/valid software-package plist for a dummy bundle/loopback URL;
this checks generator format only.

## Physical-device regression

1. Update a normally signed IPA with the same certificate/bundle ID. Restore
   session; check search/favorites/history and that no probe/export-trace UI remains.
2. Deliberately test a fresh login: automatic attempt counter, Cancel, correct
   password, actual Apple 2FA challenge when requested, wrong/expired code and
   rate limiting. Do not invent a code or repeat logout to mask Store failures.
3. Select latest and an old externalVersionId; scroll far down before reviewing.
   Numbers should load from IPA Info.plist; unavailable labels stay explicit.
4. Download and compare the displayed Info.plist version/build and selected ID
   in Downloaded apps; export through Files/Share Sheet and reopen the app.
5. Choose Download and install; after verification the Safari installer should
   open directly. Choose Download IPA only; completion should offer Install now
   and Export. Check Install latest download and Downloaded apps too. Keep Safari
   screen open and
   confirm iOS. Record the actual system result and installed version, rather
   than treating URL opening or served bytes as installation confirmation.
6. Delete one downloaded IPA with confirmation (also try swipe/delete/cancel).
   Confirm its IPA and JSON sidecar disappear, other downloads/exports remain,
   the latest-download shortcut refreshes, and the installed app/data remain.
7. Retry a temporary login failure and submit 2FA within five minutes: status
   should show Using prepared SAP session, skipping Bag/SAP setup. Check success,
   cancel, account change and expiration discard the preparation. Repeat after
   five minutes and confirm a fresh Bag/SAP setup. Measure elapsed time and
   manual attempts; verify wrong/expired code still gets an Apple error.
8. Verify server Close/timeouts stop serving without deleting the original IPA.
   Repeat with a previously licensed free app on iOS 26/27 as available.

## Build 23010: free-license acquisition and UI

Manual validation pending on a physical normally signed iPhone (iOS 26 and 27).
Earlier user-reported success does not validate these new changes.

- Use a free app never acquired by the test account. Choose version: confirm
  Obtaining free app precedes version loading and Apple's account library gains
  access. Select an older known version, download and check the saved IPA's
  version/build and externalVersionId; install and verify on-device version.
- Repeat with the owned app: version-label requests must not repeatedly acquire
  the same license within one Store session. Test another previously unowned app.
- Paid app not owned: no buyProduct request. Paid app owned: download continues.
  Unknown price: refuse automatic acquisition. Terms/age/subscription/account
  restrictions: show Apple's sanitized error, no download without confirmed license.
- Inspect all four reported screens: login (retry, cancellation, 2FA and error),
  versions (loading, acquisition failure/retry, deep scroll and confirmation),
  main actions (selected app versus labeled completed download), downloaded apps
  (Install, Export, Details, menu/swipe deletion and empty state).
- Verify Dynamic Type, VoiceOver labels, dark/light appearance, landscape, small
  iPhone and iPad. Screenshot review cannot establish contrast/accessibility compliance.
- Keep Activity log collapsed during several operations, then expand it; verify
  activity remains available and UI does not hang. Buffer is bounded to 32,768
  characters; restart clears it. It must never include credentials or signed URLs.

Fixture coverage: acquisition before download, exact older ID retained, already
owned only once, failed acquisition not cached/no download, missing license after
purchase has no loop, structured 9610 propagation and opt-out, no paid acquisition.
These fixtures do not prove live Apple purchase acceptance or iOS installation.

## Build 23011: recover purchase/download and log reading

These cases supersede 23010's acquisition-first expectation. Pending physical
iOS 27 verification, use normally signed Release IPA, preserving bundle/identity.

1. Owned free app (e.g. the previously working Pinlist): request versions and an
   older known version. Expected no Acquiring free app license before a successful
   descriptor, actual IPA version/build preserved, install/export unchanged.
2. Never-obtained free app: initial missing-license/unavailable response, then
   exactly one MZFinance/no-query purchase; confirm purchase response and selected
   version descriptor. Copy the complete log if it fails. Do not infer license
   success merely from the stage or from HTTP 500.
3. Owned paid app: download allowed. Unowned paid/unknown-price app: no purchase.
   Structured denial, invalid session, subscription/terms restrictions and rate
   limits: no automatic acquisition replay and no invented success.
4. Open Activity log: no 140-point clipping inside the original 250-point platter.
   Expand to full reading view; Copy; read older entries while output arrives;
   enable Follow latest and verify scrolling only when requested. Test Dynamic
   Type, VoiceOver, dark/light, small iPhone and iPad; not visually verified here.
5. Version app summary/download cards: no blue generic placeholder. Failure ends
   spinner and stale Obtaining free app text. Retrying starts with an access check.

Fixtures cover owned free/paid download without purchase, one acquisition after
9610/unavailable, reference MZFinance route/headers/no query, exact older ID,
5002 followed by verified download, failed/ambiguous purchase without replay or
fake ownership, opt-out and retention of accepted account kbsync after license denial.
Earlier device success is not validation of the new purchase path.

## Build 23012: navigation and catalog (device validation pending)

Use a Release IPA signed with the same bundle ID/certificate as 23011.

- Confirm five tabs, with Favourites in the center. Test small iPhone, iPad,
  light/dark, large text and VoiceOver. Native tab placement adapts to iPad/iOS.
- Search signed out, choose an app, sign in on Account; the chosen version sheet
  should open. Cancel login, switch tabs, then retry. Existing 2FA and cancellation
  behavior must remain intact. Sign out keeps saved downloads/favourites.
- Enter an App Store link, numeric ID and bundle ID; choose a version. Star the
  resolved app; open Favourites, return to its versions, swipe to remove.
- Download an older version, switch tabs during transfer, inspect progress and
  cancel from Downloads. Install directly after downloading and from Downloads;
  export/delete, open history, and expand/copy the optional activity log.
- Select Argentina, US and another supported catalog; results must use the chosen
  country, survive relaunch, and return to account country when reset. Changing
  region must never mutate account credentials/storefront. Retry network failure;
  rapid typing/tab changes must not produce stale results or stuck spinners.
- A foreign-only app may fail account lookup or license acquisition. Do not record
  catalog search success as proof of regional download eligibility.

| iOS | Device | Login/2FA | Search | Versions/purchase | Download/install/export | Tabs/region |
| --- | --- | --- | --- | --- | --- | --- |
| 27.0.1 | User iPhone, normally signed 23011 | User reports remaining minor login glitches | Reported working | User reports 23011 works well, exact cases unspecified | Earlier install/export confirmed | Not in 23011 |
| 27 | Physical iPhone, Release 23012 | Pending regression | Pending | Pending regression | Pending regression | Pending |
| 26 | Physical iPhone/iPad, Release 23012 | Pending | Pending | Pending | Pending | Pending |
