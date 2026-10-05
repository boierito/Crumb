# Current ipatool audit and port map

Reference HEAD fetched 2026-10-05:
`3411d57f451f5111ae115641c22f7ed17bbd5fbe`, dated 2026-10-01.
Links below are immutable:
https://github.com/majd/ipatool/tree/3411d57f451f5111ae115641c22f7ed17bbd5fbe

## Relevant changes after the first SAP fix

| Commit | Change | Port implication |
|---|---|---|
| abd86cb, Aug 28 | SAP signer replaces old auth | Bag → guest initialization → certificate/setup exchange → sign body |
| a53550f, Aug 29 | Bounded SAP guest timeout | Interpreter must enforce guest execution bounds, not just network timeout |
| ad88910, Sep 13 | Jailbroken iOS CLI build | Not a jailed-app solution; uses no-container entitlement and JIT memory patch |
| e5211d6, Sep 18 | Consumer catalog iOS version fallback | Latest externalVersionId resolution must handle catalogs |
| 735b689, Sep 19 | 2FA normalization and errors | Trim whitespace; exactly six ASCII digits; distinguish fresh-code-required replies |
| 2944880, Sep 29 | Reuse credentials across download retries | Avoid repeated login/2FA challenges per CDN retry |
| 1ec8b3e, Oct 1 | Preserve signed auth through redirects | Retain POST and exact plist body/attempt value; validate destination; bound redirects |
| e258fed, Oct 1 | Auth timeouts/transport retries | Bound per-request time; transient response and transport retry taxonomy |
| 0cc226a, Oct 1 | Stable machine identity | GUID and SAP hardware ID must share the same stable bytes |
| 3411d57, Oct 1 | kbsync download reliability | Adds another emulated StoreAgent flow; signer-only success is not sufficient for download |

## Signer and assets

`internal/sap/signer_local.go` loads assets, opens a guest Machine, initializes
it with hardware bytes, fetches the Bag certificate, calls Exchange(version
200), expects state 1, POSTs its buffer, then expects state 0. Sign returns
nonempty binary output; `pkg/http/client.go` Base64-encodes it as
X-Apple-ActionSignature over the exact serialized body sent on the wire.
The setup plist keys are sign-sap-setup-cert and sign-sap-setup-buffer.
HTTP 200 and nonempty plist Data fields are required; responses capped at 1 MiB.

`internal/sap/assets` downloads hash/size-pinned macOS 10.9 update components
from Apple: CommerceKit, CommerceCore, CoreFP, CoreFP.icxs, and for kbsync an
additional StoreAgent image. These are guest data, not host dlopen images.
The hardcoded update package is an **asset acquisition reference**, distinct
from dynamic auth/store endpoints. An iOS implementation must use a sandbox
cache, validate integrity, bound extraction/storage, and resolve distribution
rights. This branch does not download or bundle these proprietary files.

`machimage` parses x86-64 Mach-O, maps sections and relocates imports.
`machine/machine.go` uses guest ABI entry points:
`_cp2g1b9ro` initialization; `_Mib5yocT` exchange; `_Fc3vhtJDvr` signing;
`_IPaI1oem5iL` teardown; `_jEHf8Xzsv8K` disposal. CoreFP exports and CommerceCore
`_get_mac_address` are resolved separately. `shims*` implement allocation,
memory, platform callbacks, file-like access to asset data, crypto and guest
imports. These do not justify calling private host iOS frameworks.

Guest layout: return page 0x100000000; image bases 0x100000000000,
0x100040000000, 0x100080000000; scratch 32 MiB, heap 64 MiB, stack 8 MiB.
These guest virtual addresses must be map keys, never attempted sandbox host
addresses. SysV x86-64 call registers and stack args are implemented by invoke;
output buffers are bounds checked and disposed. A new interpreter must preserve
hooks, stop behavior, relocation, flags, indirect calls and floating/SIMD behavior
used by the actual images. A toy instruction demo cannot replace this runtime.

## Why current iOS ipatool cannot be directly embedded

`tools/build-ios.sh` explicitly says **jailbroken arm64 devices**. It statically
links Unicorn 2.1.4, exports API symbols for purego, and ldid-signs a standalone
CLI with `com.apple.private.security.no-container`.
`unicorn/library_ios.go` only changes how the library is found; it does not
remove generated executable code. Unicorn is a dynamic translator. The iOS
patch writes guest translation into host RW memory, then mprotects it to RX;
failure calls abort. It does not use an interpreter. Static linking and avoiding
simultaneous RWX therefore do not establish jailed compatibility.

Unicorn pinned source:
`8028ec436f2d9376525352dd38ed9ed6b9f6be10` (2.1.4).
`qemu/accel/tcg/translate-all.c:alloc_code_gen_buffer` requests executable
memory; `qemu/include/tcg/tcg.h:tcg_qemu_tb_exec` calls the generated host code.
The shipped tree has no `qemu/tcg/tci` or `qemu/tcg/tci.c` implementation.
Configure text mentions TCI but it is not a usable CMake interpreter option.
This is why adding a purported enable-interpreter flag is not a fix.

The reproducible Linux test in `scripts/test-no-exec.sh` runs the same Unicorn
version in a subprocess. With normal permissions MOV EAX,42 executes. After
seccomp denies PROT_EXEC mappings/protection, ordinary RW allocation still
works, RW→RX fails EPERM, and Unicorn exits allocating its dynamic translator
buffer. The test is evidence of an executable-memory dependency, **not a
measurement of Apple's actual sandbox**. The in-app native diagnostic probes
those permission requests without executing unsigned memory or aborting.

## Old → modern → adaptation

| Old WaffleStore | Modern ipatool | Action | Jailed adaptation |
|---|---|---|---|
| SHA1 Apple-ID GUID | machineIdentity from stable adapter MAC | Replace | Keychain-backed six-byte synthetic identity; Apple acceptance unverified |
| Bag auth-only/fallback | full SAPConfig from current Bag | Replace | Swift parser implemented; no unsigned fallback |
| No guest signer | SAP Machine + Unicorn | Replace | Need interpreter/AOT backend; current JIT rejected |
| JSON auth | signed serialized plist | Replace | SAPSession adapter implemented; actual signer missing |
| Recursive pod resolution | validated, bounded redirects | Replace | Dedicated URLSession delegate retaining body/method; future auth phase |
| Blocking Bool login | Login returning account/errors | Replace | Swift async state/result; after signer validation |
| Raw appended 2FA | normalized AuthCode, Apple failure taxonomy | Adapt | Preserve UI fields; don't persist code/appended password |
| Encrypted file credential payload | keychain account + kbsync cache | Replace | Native generic-password Keychain; separate identity/logout/session records |
| Search/favorites/history | lookup/search account metadata | Retain/adapt | Original public search/local data reusable |
| Manual/server version mapping | metadata list + IPA Info.plist authoritative version | Adapt | Verify selected external ID and actual IPA version |
| No purchase | free purchase + license handling | Add later | Reject paid acquisition and surface user interaction |
| Hardcoded volume download | Bag ent/download→volume→redownload/update fallback | Replace | kbsync through SAME interpreter; no token redirect forwarding |
| No kbsync | StoreAgent + hardware/DSID-bound blob | Add later | Interpreter image/ABI support and Keychain cache |
| Poll + sleep CDN | retries/resume/ranges/ZIP validation | Adapt | URLSessionDownloadDelegate progress; background lifecycle separately tested |
| Automatic localhost install | export/installation distinct | Separate later | Keep exported IPA; resign/install externally when supported |

## Authentication, HTTP and downloads

Login always fetches current Bag, initializes signer, then sends plist with
appleId, password+normalized code, guid, rmp, why and attempt. Up to four
validated auth redirects keep the same attempt/body. Login parsing preserves
Apple customer errors and requires passwordToken, DSID, storefront. A pod header
may be absent; unsigned fallback is forbidden. Upstream persists password for
re-login; this iOS adaptation should instead prompt again when needed rather
than persisting password+code.

The request retry loop allows three attempts with 10/20-second fallback delay,
maximum 30 seconds. Retry-After seconds/date takes precedence; excessive server
wait ends the operation. Empty/unusable 204, HTML/non-plist 403/404, 429 and 5xx
can retry; a parsed credential failure must not be relabeled transient. Timeout,
EOF, reset and broken pipe have explicit transport handling; cancellation is not
retried. This branch only ports SAP setup's own no-retry policy. Auth retry
tests/implementation remain future work, not claimed complete.

The October download code prefers Bag ent/download when kbsync is available.
It needs a six-byte GUID, numeric nonzero DSID, correct latest/selected external
version, Configurator 2.18 headers, serialNumber, and account-bound kbsync. Cache
only after a valid download response; retry a rejected cached blob once with a
fresh blob. Do not forward X-Token to redirect targets. There are bounded
contexts and response validation of app ID, requested external version and URL.
Then volume/redownload/update fallbacks have platform/version safeguards.
CDN/ZIP processing additionally validates framing, patches metadata and applies
SINF purchase data. It does **not** decrypt an IPA into something any sideloader
can necessarily install. Search and downloading encrypted store packages must
not be advertised as universal installation/downgrade capability.
