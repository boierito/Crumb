# Experimental LocalDevVPN access matrix (NO helper)

Base: `feature/on-device-pairing-decryption` from commit `4a3ca40d9a6766ad448f346f56cfa007bf43dde5`. Branch: `experiment/fairplay-vpn-no-helper-20261104`.

## Objective
Determine whether the authenticated *same-iPhone* RemotePairing tunnel exposes useful storage or diagnostic services **without** preparing, signing, installing or relaunching a per-IPA helper. This first step is an authorization/access experiment, **not** a FairPlay decryptor.

## How to run on iOS 27.0.1 (no Wi-Fi required for the workaround)
1. Install the experimental Crumb IPA with an ordinary sideloader. No helper signing material is needed.
2. With Cellular Data **ON**, open **LocalDevVPN**, connect it to Device IP `10.7.0.1/32`, and return to Crumb. iOS must allow Local Network access. Keep the VPN connected through all steps.
3. In **Settings → Device connection** use the existing on-device pairing once if not paired. Tap **Discover and save RemotePairing port** while cellular connectivity remains available, if discovery can find it. An unavailable result is not fatal; `49152` is the fallback. An alternate known port may be entered in **Manual pairing port (optional)**. The saved port survives switching cellular off and restarting Crumb.
4. Tap **Inspect VPN routing (no pairing required)**. It performs a packet-free UDP route lookup, classifying the source as `pdp_ip0`, LocalDevVPN-style `utun`, or unknown. This is **not** an authenticated tunnel test.
5. Open **Control Center** and turn **Cellular Data OFF** while leaving LocalDevVPN connected. Do not disconnect or reconnect LocalDevVPN after the phone has lost network access. The tunnel can continue serving local traffic without external connectivity, according to Madeira's documented iOS 27 workflow. If it disconnects, turn cellular data back on, reconnect the VPN and repeat.
6. Return to Crumb. Inspect the route again. Choose **LocalDevVPN** (not Loopback) and tap **Check authenticated connection**. The native worker verifies authenticated pairing, TLS, RSD, installed Crumb identity and a fresh Documents nonce.
7. Only if **Connected** appears, tap **Run VPN access diagnostics**. Copy/share the JSON summary; never share pairing records, P12/private keys, license files or IPA contents.
8. Turn Cellular Data **ON** after testing. This experiment cannot itself toggle the cellular radio from within the app.

### Interpretation and limits
- No iOS API allows Crumb to flip Cellular Data automatically. An iOS Shortcut could be added later, with the user importing and approving it.
- The VPN's `Connected` label alone does not establish routing or availability of the RemotePairing daemon. The route probe identifies where a packet *would* go; the authenticated connection check is the authority for actual device access.
- On iOS 27, direct TCP over `127.0.0.1` / `::1` may accept connections and then reject TLS/RemotePairing. Loopback mode is not a proven replacement.
- RemotePairing may listen on a dynamic Bonjour port instead of `49152`. During cellular-off runs, DNS-SD may be unavailable, so Crumb now reuses its saved port or a user-selected numeric override (1–65535), always constrained to its fixed local device address and normal authentication. No LAN host address is accepted.
- If `Pairing TCP: ConnectionRefused (OS 61)` continues with a valid port and LocalDevVPN route, the device daemon may not be listening or iOS may refuse the tunnel. Never report this as an authentication, certificate or FairPlay failure.

Sources:
- https://github.com/willfaust/Madeira/blob/main/docs/JIT.md (cellular-off workaround and route-probe idea)
- https://github.com/jkcoxson/LocalDevVPN/issues/31 (maintainer confirms cellular route limitation)
- https://github.com/playportdev/playport/issues/55 (dynamic RemotePairing port on iOS 27.0.1)

## Actual probes and what they mean
- `self_documents_mmap`: opens and read-only `mmap`s the existing same-installation pairing proof file. It must match the fresh nonce. Positive control for Crumb's own Documents. No FairPlay call is made.
- `rsd_advertised`: reports whether AFC, CoreDevice FileService and OsTrace are advertised by RSD. Advertisement alone does not guarantee access.
- `afc_public_staging`: creates a **small generated canary**, using only `PublicStaging/Crumb-VPN-FairPlay-Probe/probe.bin`, reads it back via AFC, and then attempts a separate read-only `mmap` at the expected local media path (`/var/mobile/Media/PublicStaging/Crumb-VPN-FairPlay-Probe/probe.bin`). The stage is removed even when the diagnostic fails. The reported `open`, `map`, `matches` and errno distinguish staging availability from app sandbox access. A successful AFC read **does not** imply Crumb or fairplayd can map or use that file.
- `coredevice_own_documents`: attempts a FileService session only for Crumb's **own** app data container, and a directory listing of `Documents`. Returns a count, not filenames. A session scoped to app data is not installed-bundle access.
- `ostrace_passive`: attempts a *3-second passive* OsTrace stream and returns event counts/mentions of FairPlay or system policy, **never raw messages**. It deliberately does not perform a decrypt concurrently, so observed mentions cannot be treated as the causal error for an `mremap_encrypted` attempt.

All service operations are bounded. No third-party app is installed, updated, queried for private paths or accessed. No arbitrary path is accepted from an IPA or the UI. No licensed IPA, SINF, SUPP, SIDB, credentials or certificate is copied by this experiment.

## Interpreting outcomes
| Result | What is supported | What is NOT supported |
|---|---|---|
| `afc_roundtrip_ok` | Host tunnel can transfer the canary | App access, FairPlay authorization |
| `app_local_mapping.open=false` | Direct media path not openable from Crumb | Whether fairplayd can read it |
| `app_local_mapping.map=true` | Crumb can map the staged canary | FairPlay can read SINF or decrypt it |
| `own_container_list_ok` | CoreDevice lists own Documents | FileService writes to installed app bundle |
| `passive_sample_only` | OsTrace accepted a bounded stream | Rejection reason for a specific IPA |

## Next experiments (not yet implemented)
1. Use the **same purchased test IPA SHA256** in an installed-bundle positive control vs runtime-imported Crumb data, without attributing the `errno 12` mapping to a specific failure. Physical results must include mapped/remap attempted, remap return, page counts, original SINF availability and masked diagnostic stage. Do not upload license material.
2. If a *new, actually readable* staging class is discovered, try only a read-only candidate there after verifying app-side open + mmap. CoreDevice currently documents AppDataContainer, AppGroupDataContainer, Temporary and SystemCrashLogs; it does not expose an arbitrary installed-bundle write domain. Do not infer one.
3. For a causal trace, start a bounded OsTrace observation **concurrently** with an existing direct-runtime FairPlay attempt and correlate timestamps without returning private messages. Passive sampling alone cannot prove cause.
4. Remove the helper dependency only after an iPhone 27.0.1 test recovers and verifies **all** encrypted Mach-O pages from an IPA imported after Crumb's one-time installation.

## CI validation
`cargo test --locked --manifest-path NativeDevice/Cargo.toml`; release Rust tests; iOS ARM64 build via the existing Build IPA workflow. Native host tests cover path restriction, read-only mapping success, bad contents, and removed file. GitHub Actions cannot test the self-pairing tunnel, iOS System Policy or native FairPlay decryption.

## Known boundaries
The PublicStaging local path is a hypothesis, not a guaranteed alias for the AFC domain on every iOS release. The experiment never assumes the tunnel grants `fairplayd` access. Errors are reported at service or POSIX-code granularity. The existing helper pipeline remains untouched in this feature's base, but this button neither invokes nor installs it.


## 2026-11-04: direct-runtime probe with correlated OsTrace (schema 2)

The physical 2026-11-04 report showed AFC/PublicStaging round-trip succeeded but the same file could not be opened by the Crumb process (\`EPERM=1\`). CoreDevice FileService was not advertised; idevice error code 21 is \`ServiceNotFound\`. The 128-event passive trace collected no FairPlay messages, but it did **not** overlap any FairPlay syscall.

Instead of trying another staging alias, the new option in **Downloaded apps → Decrypt on this device → Direct FairPlay without helper** now:

1. Uses the existing completed IPA from Crumb's downloads; it is not modified.
2. Scans archive entries with bounded limits and only accepts a well-formed encrypted, thin ARM64 Mach-O in \`Payload/*.app/\` with one valid encryption-info command; never trusts traversal paths or symlinks.
3. Extracts only that executable into Crumb's own random temporary file, hashes the encrypted executable, and attempts a strictly read-only \`mmap\` followed by a single native \`mremap_encrypted\` call (\`cryptid=2\`), without generating or publishing plaintext.
4. Immediately removes the temporary input file, including on failures.
5. If **Connected** is true, performs the same probe using the existing authenticated LocalDevVPN route while capturing an aggregate, time-bounded OsTrace sample. If not connected, runs the native probe **without VPN** and explicitly reports the absence of correlated logs.
6. Returns a JSON report identifying whether \`mmap\` was attempted and succeeded, whether native remap was attempted and its exact error, and an encrypted-image SHA-256 for comparison. It reports **zero verified decrypted pages** regardless of remap return because this diagnostic deliberately does not validate or export plaintext.

This is a deliberate measurement of the reported direct-native \`errno=12\` problem, not a workaround for FairPlay or a proven decryptor. No helper signing, installation, P12, jailbreak, JIT or computer is used by this test. The IPA must be lawfully obtained for the testing account. The OS trace output contains only aggregate event counts, not user data or system log strings.

**Physical evaluation:**
- On the iPhone's own downloaded encrypted IPA, run once with the tunnel disconnected and record \`local_fairplay.native\`.
- Establish LocalDevVPN with the cellular-off method only if desired and available, verify the authenticated connection, then run it again to obtain \`trace\` alongside the native result.
- Compare the same \`image_sha256\`, syscall result, page geometry and system mentions.
- If \`remap_result=0\`, the test still **does not** declare decryption successful. A separate validated all-page reader is required.
- A \`remap_errno=12\` (ENOMEM) is known to collapse multiple crypter-creation errors in this iOS branch, and still does not disclose the daemon's private cause by itself.

**Limitations:** Mach-O and IPA formats beyond encrypted thin ARM64 are rejected; only the first acceptable encrypted executable is tested. The trace's short time window is correlation, not a complete activity archive. Dynamic SINF resources are not placed into an installed system-trusted bundle by this probe. This intentionally leaves the original no-helper FairPlay policy question open.

