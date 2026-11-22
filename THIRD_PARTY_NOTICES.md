# Attribution and redistribution status

## iPadN WebKit / WebAssembly SAP engine (2026-11-22)

The WebKit signing service in `Crumb/MapleSyrup/iPadNSAPSigner.swift`,
loopback server in `Crumb/MapleSyrup/iPadNLoopbackServer.swift`, and
related adapters are ported from [iPadN](https://github.com/mengangs/ipadn),
upstream commit `70f2ac36e7f87cacf82a45b58e6fa27ed23ecb80`,
under GPL-2.0-or-later. Its nine-file SAP runtime is fetched from that
pinned commit during CI, checked against each vendor SHA-256 in
`scripts/ipadn-sap-manifest.json`, and bundled as `SAPRuntime.zip`.
The runtime includes `sap.wasm`, JavaScript, `libCoreADI.so`, and
`libstoreservicescore.so`. Upstream's
`iPadn-iOS/vendor/SAPRuntime/PROVENANCE.json` identifies the latter
libraries as coming from the official Apple Music for Android package and
explicitly states that a complete license grant for every component has
not been established. **Technical provenance is not redistribution
permission**. Publishing binaries that include those files requires
a separate review of Apple and other upstream rights.

## Vendored iPadN FairPlayCore (2026-11-21)

The on-device FairPlay engine in `Crumb/Device/FairPlayCore.c` and
`Crumb/Device/FairPlayCore.h` is vendored from
[mengangs/ipadn](https://github.com/mengangs/ipadn) (GPL-2.0-or-later).
Upstream file hashes at import: C blob `e68bbbc3d22f21ce53df936a5df36e31a245669c`,
header blob `aaa012b8e22216999c82b5d2c82e4dcc131c8d12`.
License text: `third_party/iPadN-GPL-2.0.txt`.
The `GroupFairPlay.m` adapter connects this engine to the existing Crumb pipeline.
Publication of binary releases remains subject to review of the combined
GPL and unrelated third-party redistribution permissions described below.



The SAP v200 protocol/state flow in `MapleSyrup/SAPKit` is adapted from
[majd/ipatool](https://github.com/majd/ipatool), pinned at
`3411d57f451f5111ae115641c22f7ed17bbd5fbe` (2026-10-01).
The adaptation uses Swift protocol code and a statically linked guest-only
Go library, not an embedded ipatool executable. Experimental Unicorn/TCI is
statically linked in the IPA. The legacy guest downloads its own hash-pinned
assets into the app cache when requested. Separately, the iPadN SAPRuntime
bundle includes vendor binaries whose redistribution rights remain unverified
(see iPadN section above).
The Linux negative test downloads upstream Unicorn separately via pip.

## ipatool MIT notice

MIT License

Copyright (c) 2021 Majd Alfhaily

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## WaffleStore and future runtime work

At the original revision `d508e532e5b40d6470bab003a6cc0026d1ca7c75`, the
WaffleStore tree has no LICENSE/COPYING file granting redistribution rights.
Public source visibility alone is not an MIT license. Existing Mineek,
PancakeStore, jailbreak.party, Skadz and nxtcoreee3 attributions are retained.
Public redistribution of a modified binary needs clarification from the
relevant upstream rightsholders; this work does not grant those rights.
CI creates **draft** prereleases for development tags, not public releases.

The existing SwiftPM dependencies retain their upstream notices. See the pinned
`Package.resolved` for PartyUI, Zip, Telegraph, CocoaAsyncSocket and HTTPParserC.
Unicorn includes GPL-2.0 code; ipatool's MIT license does not relicense Unicorn.
The restored QEMU TCI sources preserve their GPL headers. Every CI artifact
includes the modified corresponding source for Unicorn/TCI, its COPYING file,
and the source/notice for the guest bridge. The embedded Go dependencies retain
their licenses (including purego Apache-2.0); the build generates their complete
notices from the dependency modules. The combined licensing and upstream
redistribution rights require review before public binary publication; providing
source alone does not resolve an incompatible combination or missing upstream
grant. Apple CommerceKit/CoreFP/StoreAgent binaries have their
own proprietary terms; downloading them from Apple does not grant a right to
bundle them in a redistributed IPA. No such binaries have been copied here.

## Experimental device pipeline

- jkcoxson/idevice: MIT, pinned to 89702099f4d7552f704f93b4526756abda319ceb.
- apple-codesign-quick 0.1.0 (Dadoum/apple-crates): LGPL-2.1-or-later.
  Corresponding modified source is in NativeDevice/vendor/apple-codesign-quick,
  including COPYING.LESSER and CRUMB_CHANGES.md. The change prevents re-signing
  of EmbeddedPayload dylibs. NativeDevice/Cargo.lock and build scripts permit
  rebuilding/relinking; CI provides the corresponding source archive.
- p12-keystore and other native dependencies retain their manifest licenses and
  exact locked versions. No FairPlayed source code is copied. Its public
  documentation informed the independent installed-resource experiment.


ZIPFoundation 0.9.20 (https://github.com/weichsel/ZIPFoundation)

MIT License

Copyright (c) 2017-2025 Thomas Zoechling (https://www.peakstep.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
