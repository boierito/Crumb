# Production bridge integration outside the app

Build the portable production bridge without SwiftUI or the iOS C bridge:

```sh
bash scripts/build-fairplay-host.sh /private/new-build-directory
/private/new-build-directory/crumb-fairplay-host \
  --ipa /private/original.ipa \
  --capture /private/capture.zip \
  --output /private/new-processed-directory > /private/result.json
```

The build directory must not exist. Keep it and the inputs private. Building
may download Go dependencies; the capture replay itself does not contact Apple
or send inputs to a server. Do not submit private files to public CI.

This executable copies the same production `labRun` implementation and recovery,
validation and IPA processing packages used by Crumb. It calls selftest,
inspection, capture analysis, exact capture/IPA binding, and capture recovery
in order. The first failure stops the pipeline with exit status 1 and fixed
diagnostic fields. Source/frame and SINF validation remain mandatory. Only
successful validated recovery proceeds to the real IPA writer for all images.
The selftest is explicitly synthetic; it is not evidence of real decryption.
Models are temporary, and plaintext authentication is not asserted.

The host harness does not test Swift Keychain, UI lifecycle, signing,
installation, live Apple authorization, or the retained StoreAgent producer.
Those boundaries are explicitly reported instead of being mocked as successful.
An original capture predating authorization cannot reproduce the later account
state saved on an iPhone. The harness therefore does not claim to reproduce
that state or the complete iOS operating environment.

`test-fairplay-host.sh` also compiles the host entry point after exercising the
portable contracts, so future changes must continue to build outside the app.
