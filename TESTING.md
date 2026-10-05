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
