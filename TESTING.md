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
is claimed. Integration build 23007 removes probes/tracing from the app; its
cleaned UI/credits should receive an update/regression check on device.

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
5. Request installation from Downloaded apps, keep Safari screen open and
   confirm iOS. Record the actual system result and installed version, rather
   than treating URL opening or served bytes as installation confirmation.
6. Verify server Close/timeouts stop serving without deleting the original IPA.
   Repeat with a previously licensed free app on iOS 26/27 as available.
