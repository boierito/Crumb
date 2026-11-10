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
   the **actual installed signature must authorize an accessible App Group**.
   Schema 2 prefers `group.<installed Crumb bundle ID>.worsequery`, then a
   dedicated renamed `.worsequery` group, then another actually granted group
   with an accessible empty/missing Caches directory. No identifier is guessed
   or hardcoded from a user certificate. The Xcode entitlement request is
   included; an unsigned IPA or a matching group spelling cannot grant it.
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
4. If an eligible effective group exists, perform a private-canary-directory check,
   then independently try worse_query against the downloaded tree and selected
   installed bundle. The canary already belongs to Crumb; success does not prove
   new cross-app access. A failed canary does not by itself prove every target
   path is unsupported.
5. Keep all IPA and SINF resources private. Schema 2 skips the previous
   shared-payload storage experiment because signer-provided groups may be
   shared with other apps. Only cache-directory metadata is inspected during
   selection; no other apps' data or system group is enumerated.
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

Schema 2 reports contain a random run ID, source revision when built by the
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
| `app_group.dedicated_group_authorized=false` | The expected dedicated name is absent; other effective groups may still be usable |
| `effective_group_count` / `accessible_group_count` / `eligible_group_count` | Distinguishes signed grants, resolvable containers and empty-cache eligibility |
| `selected_group_authorized=true` | A candidate from the effective signature was selected; `selection` describes dedicated vs profile-provided |
| `candidates[].cache_state=ineligible` | Preserve numeric errno; occupied or inaccessible group data is not cleared |
| `cache_lease_unavailable` / `cache_lease_activation_failed` | Exclusive transaction could not start; not a FairPlay error |
| `caches_restored=false` | A private journal retains the original cache inode for recovery; foreign replacements are never overwritten |
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

## Schema 2: signing-service group compatibility (build 23040)

A returned build-23039 report showed `dedicated_group_authorized=false` and
`entitlements_readable=true`. A subsequent KSign profile view showed five
concrete App Groups with service-selected names. That view proves profile
permissions, not the exact installed process's grants. The old dedicated-name
check was too restrictive to test those alternatives. Schema 2 reports counts
and candidate outcomes from the actual process, without exporting group IDs.

The upstream file-mutating wrapper is replaced with `worse_query_prepared`:
it verifies the caller's symlink and only performs the original container query
and token consumption. `CrumbCacheLease` owns the file transaction. It accepts
only an empty/missing Caches directory, takes a cooperative process lock, writes
a private journal before changing Caches, preserves the original directory by
an atomic non-overwriting rename, then restores that exact inode. Contents that
arrive after selection or through an already-open directory descriptor are
preserved. If another app replaces Caches, the replacement and original backup
are retained without overwrite and the journal remains pending. A temporary
namespace change may still make simultaneous cache operations fail; occupied
candidates are skipped rather than emptied.

The link is prepared in Crumb's private storage; no IPA or SINF is copied into
the selected shared group. A cross-volume rename fails explicitly and rolls
back. Recovery rechecks the group's current effective entitlement before
restoring a recorded transaction. No certificate, profile or extra entitlement
is fabricated or installed by this change.

`python3 scripts/test-cache-lease.py` tests real filesystem transitions with
sanitizers: inode/mode preservation, occupied directories, late writes, crash
replay, exclusive locking, foreign replacements, existing staging files and
originally absent Caches. macOS tests additionally cover generic profile-group
selection, dedicated aliases, redaction, exclusion of unsigned/wildcard guesses
and JSON-journal recovery. These tests do not prove iOS token issuance.

## Schema 3: installed-bundle research (build 23041)

Use **Find installed app** on the existing without-helper screen. This is a
separate mode; it does not repeat the already-failed downloaded-tree remap.
It prepares the original IPA only for exact identity/resource comparisons.

The mode tries the existing target-specific LaunchServices lookup first, then
measures read/list access to the installed-bundle parent, including one scoped
worse_query request. It retains a successful parent grant until installed-file
attempts finish. If ordinary discovery remains blocked, it searches recent file
IDs on the same data volume, bounded to 3,000,000 probes / 20 seconds of the
cooperative discovery budget. A currently executing OS call can outlast that
budget. APFS file-ID ordering is not assumed to cover the whole filesystem:
`bounded_scan_incomplete` never proves the app is absent.

The inode paths are immediately filtered to the selected app's basename or
lsd's `Library/Caches/com.apple.LaunchServices-*.csstore` under an InternalDaemon
container. Other paths are discarded without logging. Registry reads request a
token only for the identified daemon container, cap files at 64 MiB, and reject
snapshots whose inode, size or mtime changed during the read. Registry contents
stay in memory; only a validated cache-path hint is saved privately.

The CSStore2 27.0 record/string layout comes from inspection of
https://github.com/kolbicz/WQFiles/blob/f149ba83050ef53ebdb9b73f2e83447f88057965/WQLaunchServices.m
It is a research format, not an Apple-supported API. A small independent parser
returns at most eight paths associated with the selected bundle ID. Path
strings must have a UUID and exactly one `.app` component under the installed
bundle root. Live bundle ID, version/build and encrypted-image SHA256 validation
remain mandatory before decryption/publication. A located but unreadable path
is recorded as an unverified candidate, not a working route. Metadata reads and
image reads happen at the original installed path; executable contents are not
copied elsewhere to ask FairPlay to decrypt them.

`installed_alternative` separates parent access, token issuance, store parsing,
search completeness, candidate identity and the subsequent native attempt. Raw
paths, other apps' names/IDs and registry bytes are not exported in reports.
The existing journalled cache lease is reused, and no installed-bundle files
are written. The `no_verified_installed_route` outcome preserves these distinct
barriers. Installed-path research does not create a daemon-readable location
for a downloaded IPA when the target app is not installed.

Host tests use a synthetic CSStore2 fixture and filesystem/token doubles;
`test-ls-store.py` runs the parser under address/undefined-behavior sanitizers.
Actual iOS lookup, `fsgetpath`, bundle-parent extensions, registry layout and
FairPlay reads still require physical validation.
