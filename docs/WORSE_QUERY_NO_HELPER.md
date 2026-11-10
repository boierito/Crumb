# In-process worse_query experiment

Branch: `experiment/worse-query-no-helper`, based on main
`fae4d4a0534fccfcced5dc4aae3da404c948c2a0`.

This branch implements a measurable experiment, not a demonstrated universal
decryptor. The physical target selected for this investigation is iOS 27.0.1.
The optional installed-app route is an accepted partial outcome.

## Evidence used

- worse_query is pinned to `d1540e20cd82b974d891e0551ced6d6aeb570ac2`:
  https://github.com/forcequitOS/worse_query/tree/d1540e20cd82b974d891e0551ced6d6aeb570ac2
  Its README lists data containers, requires an App Group, and gives the version
  string `iOS 26.0 - 26.7.2 / 27.2`. Neither installed-bundle access nor operation
  on this particular iOS 27.0.1 installation is assumed. The native test is
  explicit and records the actual OS build.
- Historical Crumb no-helper work:
  https://github.com/boierito/Crumb/blob/f45fce1c9a68f113aa17d202451e98f72eaf9e32/docs/FAIRPLAY_VPN_NO_HELPER_EXPERIMENT_20261104.md
  This documents a matched-SINF direct attempt returning errno 12 and correlated
  path/SINF policy-denial classifications. Errno 12 alone does not establish
  memory exhaustion or identify the rejecting process.
- Historical installed-resource control:
  https://github.com/boierito/Crumb/blob/f4bf77608b2b9ad3164447e4a37a34ff826d393e/docs/FAIRPLAYED_SYSTEM_ROUTE.md
  It documents 178 recovered pages on a stock device. This is prior recorded
  evidence, not a new control run by this branch.
- Native mechanism and license-path research:
  https://github.com/morzelowski/FairPlayed/tree/bad7df138839b5e0ac8b9f20cc4689140b8436fb
  A sandbox extension consumed by Crumb does not demonstrate a grant to
  fairplayd. The experiment measures Crumb's access and actual page recovery
  separately.

## Run once on the phone

1. Sign/install the experimental Crumb IPA normally. To exercise worse_query,
   the **actual installed signature and profile must authorize**
   `group.<installed Crumb bundle ID>.worsequery`. The Xcode entitlement request
   is included; an unsigned IPA or a matching group spelling cannot grant it.
   A sideloader that drops/remaps this capability can leave it unavailable.
   There is no requirement for helper signing files, pairing or LocalDevVPN.
2. Use an encrypted download belonging to the phone's licensed App Store
   account. Open Downloads → its menu → **Try without helper**. This account
   requirement is distinct from signing into Crumb.
3. Keep Crumb foregrounded. The same installed version is detected and checked
   automatically when iOS permits discovery. Do not install another target
   merely to diagnose a missing App Group or failed local archive preparation.
4. Share the test report. A completely verified IPA also appears in Downloads.
   A decrypted IPA still needs appropriate signing before ordinary installation.

All baseline candidates run before worse_query issues a token. A later test in
that process skips the baseline and asks for a Crumb restart for a clean
comparison: closing file descriptors and releasing a token are not treated as
proof that the process has returned to its original access state.

## Pipeline and stopping rules

1. Extract the IPA with bounded streaming, CRC checks, raw central-directory
   duplicate detection, safe paths and rejection of symlinks. Retain a private
   original inventory, including empty directories, modes and non-Payload files.
2. Keep the original executable basename and adjacent license files. Replicate
   only manifest-declared original SINF mappings; ambiguous distinct licenses
   are rejected rather than guessed. Replicas are used for the source experiment
   only; the export retains the original archive inventory.
3. Try the private downloaded tree without a grant. Discover only the selected
   installed bundle; check its bundle ID/version/build and exact encrypted image
   hashes. A mismatched installed version is not mixed with the download.
4. If the dedicated effective group exists, perform an own-container canary,
   then independently try worse_query against the downloaded tree and selected
   installed bundle. The canary already belongs to Crumb; success does not prove
   new cross-app access. A failed canary does not by itself prove every target
   path is unsupported.
5. Try one additional storage class in the already-authorized dedicated own
   App Group. It is a separate hypothesis, not proof that the daemon accepts
   shared containers. No other apps' data or system group is enumerated.
6. For every admitted encrypted image, use the shared read-only cryptid-2
   reader. First attempt without signature registration; record a registration
   retry separately if the remap fails. Read every encrypted page using
   `vm_read_overwrite`, recording Mach status, counts, changed bytes, entropy and
   codeword heuristics. Patch cryptid only after these checks pass.
7. A successful route must complete every image. Validate unchanged image size,
   unchanged bytes outside the encrypted range and cryptid, all original
   resources/modes, full inventory and a separately reopened ZIP's CRC/hash and
   encryption state before publication. Candidate routes never contribute pages
   to one another.

Each route is attempted at most once per run (with the single measured native
signature-registration retry). A first image failure ends that route. No helper,
self-reinstallation, transport fallback, executable mapping, JIT or arbitrary
system-file writes are triggered by this coordinator. The existing helper UI
remains separately available and shares the same reader implementation.

## Debug report and interpretation

Schema 1 reports contain a random run ID, source revision when built by the
repository script, device model/OS/build/page size, capability checks and a
bounded attempt/image list. They exclude account identifiers, full local paths,
SINF bytes/hashes, tokens, P12 material and raw OS logs. Image ciphertext SHA256
and Mach-O UUID identify comparable executable content.

The last checkpoint is atomically persisted before native operations. Page-read
checkpoints are persisted for the first page and every 32 pages; a crash can
therefore identify the operation without establishing the exact final page.
Restart marks a still-active report `interrupted` and cleans only this runner's
UUID workspaces and ledger-recorded group artifacts. Stop is cooperative and
waits for an already-entered OS call to return.

| Report | Interpretation |
| --- | --- |
| `app_group.dedicated_group_authorized=false` | Installed Crumb does not have the required grant; no conclusion about worse_query's kernel behavior |
| `entitlements_readable=false` | Effective capability inspection unavailable; no group authorization is inferred |
| `grant.result < 0` | Token acquisition failed; preserve target-stat errno and upstream result separately |
| `source_unreadable` | Crumb could not hash the original; native decryption was not established |
| `installed_version_mismatch` / `installed_image_mismatch` | Installed target is not a valid matched comparator |
| `matching_license_missing` | Staged licensed-file context is incomplete; not evidence of a daemon policy decision |
| Remap success, partial pages | No complete decrypted image; first read failure is reported |
| `downloaded_ipa_success` | All images/resources verified for a downloaded IPA; inspect attempt source/access to identify which route worked |
| `installed_app_success` | All images/resources verified using the matched installed app; accepted partial solution |
| `no_verified_route` | None of the attempted routes produced a complete result; inspect stages and prerequisites |

`daemon_observation=unavailable_without_external_trace` is deliberate. This route
has no remote logging transport. Crumb-side EPERM/ENOMEM cannot identify the
private daemon's exact denial reason. `plaintext_authentication=not_established`
also remains explicit: instruction/entropy checks and archive integrity are
useful structural evidence, not independent cryptographic plaintext proof.

## Automated verification and remaining physical boundary

Run `cargo test --locked --manifest-path NativeDevice/Cargo.toml` and
`python3 scripts/test-local-probe-reader.py`. Real ZIP fixtures cover preparation,
manifest replication, ambiguity, duplicate/traversal/corrupt inputs, multi-image
replay, mode/resource preservation, changed originals, missing extensions,
receipt geometry, cancellation and failure redaction. C fault callbacks cover
successful/failed remapping, signature retry, partial reads, cancellation,
unchanged ciphertext, symlinks, non-overwrite and malformed Mach-O. They are
synthetic controls, not a FairPlay implementation.

The Build IPA workflow compiles against the iOS SDK, exercises coordinator
interruption recovery/cancellation/redaction with a stubbed archive boundary,
runs existing native/parser/resource/persistence tests, then builds Debug and
Release IPAs. It also packages corresponding experimental source and notices.

Only the phone can establish effective post-sideload permissions, actual
worse_query behavior, App Store license acceptance and fresh protected-page
recovery. Simulator, Linux and replay results do not establish those facts.
FAT/32-bit images and existing archive budgets remain explicit unsupported
inputs; no version gates in the existing helper flow are removed.
