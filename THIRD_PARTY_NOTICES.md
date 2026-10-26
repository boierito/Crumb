# Attribution and redistribution status

The SAP v200 protocol/state flow in `MapleSyrup/SAPKit` is adapted from
[majd/ipatool](https://github.com/majd/ipatool), pinned at
`3411d57f451f5111ae115641c22f7ed17bbd5fbe` (2026-10-01).
The adaptation uses Swift protocol code and a statically linked guest-only
Go library, not an embedded ipatool executable. Experimental Unicorn/TCI is
statically linked in the IPA. Apple proprietary assets are **not** bundled;
the guest downloads hash-pinned assets into the app's cache when requested.
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

### DrWhax/ipatool portable FairPlay decoder

The native build imports `internal/fairplay` from DrWhax/ipatool commit
`e94fcd7362036a1b2269e5b555232f20d8cd273e` under its MIT license. Its license
and corresponding modified sources are included in the existing native source
archive. Crumb adds iOS to the Unix permission/platform checks and to case
insensitive path comparisons. The recovery helper and private enrollment keys
are not bundled. The existing SAP backend remains pinned independently.

FairPlay Laboratory updates the portable FairPlay source pin to
`ef98227e061838320424744d21a513baba7ade28`. The generated state/frame/SIDB
operations and sanitized embedded models are adapted to separate Go packages;
corresponding sources remain in the native source archive. The partial recovery
algorithm's Apache-2.0 attribution is in `LICENSE-DARKPHOENIX.txt` alongside those
sources. No private enrollment model is embedded in Crumb.

The modern native-enrollment runtime is adapted separately from DrWhax/ipatool
`f7242cce58077a4af8db738ed295ee23ad8f067d` (MIT), preserving the Apache-2.0
DarkPhoenix attribution. Its generated arithmetic/CFG code and sanitized
account-entry model are packaged in `wafflebridge/recovery/rawsource`. The
executable path, hardware identity, license and state files are private runtime
inputs. Their sources, profile README and licenses are included in the native
corresponding-source archive. The SAP guest and portable decoder pins remain
independent; no recovered account model is embedded.

### DrWhax/UnfairPlay generic semantic interpreter

The generic instruction-profile compiler and compact interpreter are adapted
from `DrWhax/unfairplay` commit `c54e751775cbc2a54510c0fcafb0aef6b1619bc9`
under its MIT license. Original source hashes, LICENSE and NOTICE are retained
in `third_party/unfairplay-generic`. Build-time adapters preserve those files
and include the generated interpreter/profile and notices in the corresponding
source archive. The bundled profile contains semantic operations and public
code locations for the exact hash-pinned build; it contains no account state
and no Mach-O executable. Profile coverage is explicitly censused, including
unsupported instructions. This does not establish enrollment compatibility or
successful decryption.
