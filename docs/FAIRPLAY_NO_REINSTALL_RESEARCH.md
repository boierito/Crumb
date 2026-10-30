# Direct runtime route: no per-download Crumb update

Status: research implementation, NOT a demonstrated solution to data-container
access. No direct-route device decryption has been observed. No installation or
signing fallback is invoked by the new Download → Decrypt coordinator.

## What actually changed

`PrepareSystemLocal` extracts the original purchased app into a new private data
directory. Both it and the old embedded-helper preparer now share the same
`systemLicensePlan`, including exact original SINF mapping and refusal to guess
between distinct licenses. Original binaries, metadata and licenses are retained
unchanged for the native call; replication follows the original manifest only.

Swift invokes the same read-only cryptid-2 native reader on these local files,
then the existing all-executable/all-page export validator. It never asks for a
certificate or invokes signing, helper packaging or an installer on this route.
The data tree and partial native outputs are removed after the attempt. A failed
request records only fixed numeric errors/counts, the route and the stage; no
path, license bytes or account identifiers are included. Failure cannot generate
a successful or partial decrypted IPA.

The old engine remains available in source for research and already-installed
resources; it is not an automatic fallback that reinstalls Crumb.

## Evidence and limits

1. **Observed on the user's phone:** the returned FairPlayed log shows successful
   reads of 178 protected pages from the installed bundle. It contains no matched
   experiment on a data-container copy. Do not claim this log independently
   demonstrates a data-container denial.
2. **Upstream's claim:** FairPlayed reports System Policy denying fairplayd access
   to data containers, including symlink variants. Its extraction and path-variant
   code is inspectable, but our supplied device log does not reproduce those
   variants. [Pinned upstream](https://github.com/morzelowski/fairplayed/tree/bad7df138839b5e0ac8b9f20cc4689140b8436fb).
3. **Kernel contract:** Apple's `mremap_encrypted` resolves the backing vnode with
   `fill_procregioninfo_onlymappedvnodes`, obtains `vn_getpath`, and passes the real
   filename, CPU and origin to `text_crypter_create`. Both cryptid 1 and 2 choose
   `com.apple.unfree`. Its syscall arguments do not include an SINF buffer,
   alternate license path or sandbox-extension token. A string alias, anonymous
   memory or merely changing cryptid does not establish an alternate license
   handoff. [Pinned XNU source](https://github.com/apple-oss-distributions/xnu/blob/f6217f891ac0bb64f3d375211650a4c1ff8ca1ea/bsd/kern/kern_mman.c).
   This source audit is not a claim of exact iOS 27 kernel binary equivalence.
4. **Direct driver alternative:** the extracted macOS 26A428 FairPlayIOKit
   `initWithTask` at `0xfffffe000a458e74` calls `copyClientEntitlement` for
   `com.apple.private.FairPlayIOKitUserClient.access`,
   `com.apple.private.ProvInfoIOKitUserClient.access` and
   `com.apple.private.FairPlayIOKitUserClient.Virtual.access`. At
   `0xfffffe000a458ee4`–`0xfffffe000a458efc`, all-null entitlement results select
   the rejection path; subsequent code checks the true singleton. This rules
   out treating this macOS user-client interface as an unrestricted portable
   crypto service. The exact imported/fixed binary SHA-256 is
   `61348129cfee02707e036e16bd257cf08426710e57cccd377edc1658a319a706`.
   It is NOT proof that every iOS service/interface has the same restriction.

No verified service accepting an arbitrary caller-supplied mobile SINF from an
ordinary iOS app was found in this investigation. No sandbox entitlement, UID key
or zero-filled substitute is fabricated. Merely issuing a file-provider access
permission to Crumb would not demonstrate that fairplayd consumes that permission;
that consumer behavior must be established before claiming it solves the issue.

Follow-up: [matching daemon and system-storage contracts](FAIRPLAY_SYSTEM_STORAGE_CONTRACT.md)
records the original daemon's mobile/mobile launch identity, matching protobox
profile, unresolved caller-entitlement helpers, and the limits of Background
Assets and ODR as alternate storage mechanisms. These observations do not
identify the cause of the phone's native remap error.

## Local verification with real DD data

The new direct extractor retained all 57 original DD files byte-for-byte and
planned its two encrypted ARM64 executables. Go race tests cover extraction,
license replication, absent license, ZIP corruption, non-overwrite and cleanup;
the pre-existing multi-license/layout/export suites still pass after sharing the
mapping implementation.

The 178 device-recovered plaintext pages were replayed through the new local tree
and existing export gates. For this replay ONLY, original non-encrypted bytes,
including original signature data, were retained; the prior helper's signature
changes were not accepted as native decryption. Output: 43 files, two executable
slices, 178 changed pages, two embedded signatures removed, four known metadata
entries excluded and one known identifier checked. This is **replay**, not fresh
system decryption, and does not verify that fairplayd can read the local directory.

A host cannot supply the iPhone's FairPlay daemon, account state or hardware. The
remaining unknown is the actual native service result for the local source on
that iPhone, or a verified alternative interface that legitimately accepts the
license. Until one succeeds, this route is not presented as a working generic
one-tap decryptor and no new IPA is offered as a solution to that limitation.

## Authorized shared-container hypothesis

The direct route can also try an **already granted dedicated Crumb** app-group container
after a native failure from the ordinary data container. The actual own-process
`com.apple.security.application-groups` entitlement is read with SecTask, and
Foundation must return the corresponding container URL. The only eligible namespace is `group.<current Crumb bundle ID>.fairplay`;
general-purpose sideloader groups are never used for private licenses. The
identifier must be present in the actual entitlement and Foundation must grant
access; a matching spelling alone grants nothing. No guessed accessible path,
modified profile, other app's container, global temporary path or fabricated
entitlement is accepted. With no actual authorized groups this route is skipped.

Each shared attempt uses a fresh UUID-named, protected, backup-excluded directory
and the same untouched target, native reader and export validator. It never
mixes pages from different attempts. Its own directory is removed afterward;
other group files are untouched. Results identify only the storage class, never
the group name/path. Preparation/validation failures do not trigger relocation
retries. This probes a distinct storage class, **not a demonstrated policy
exception**. It may be refused by the same system policy, and would not provide
universal support for ordinary installs without an authorized app group.

## Matching iOS kernel follow-up

The official iPhone17,1 iOS 27.0.1 IPSW is build **24A446**:
`https://updates.cdn-apple.com/2026FallFCS/59ffd938-75b6-41e5-af04-f4295ee8aeb3/iPhone17,1_27.0.1_24A446_Restore.ipsw`.
Its kernel reports `xnu-13432.2.10~2/RELEASE_ARM64_T8140`. Extracted FairPlayIOKit
version: 72.19.0; SHA-256:
`6bae7cc2a7e6808dcceefe1dabe16d97458d6431c0fdbb00e508357a92c482f1`.
Unlike the earlier macOS input, this kernel is for the model/OS version identified
by the supplied device log. The phone's exact installed build was not recorded
in that log; the firmware catalog supplied the 27.0.1 build association.

The three entitlement-name queries occur at `0xfffffff009bdf8c0` through
`0xfffffff009bdf8f4`. Their absent-results branch at `0xfffffff009bdf8fc` through
`0xfffffff009bdf914` selects the `entitlement not found` rejection. Entry inferred
from the actual prologue and identical gate structure: `0xfffffff009bdf890`.
The iOS binary is stripped; the entitlement-query call target has no exported
symbol name. Do not present the macOS symbol name as a resolved iOS symbol.

`scripts/audit-fairplay-ios-access.py --driver <extracted driver>` pins those exact
bytes and replays the absent-entitlement branch in Unicorn. With only the three
lookup results modeled as NULL and the logger mocked, the original 47 executed
instructions return false, without patching a branch. This verifies the rejection
control flow; it does not measure an iOS service, authenticate a license or prove
all possible interfaces inaccessible. PAC stack instructions are explicitly
modeled and are not a cryptographic authentication result.

A private cleanup ledger is written before any license is copied into the
dedicated shared directory. Completion removes it; restart also reaps only
UUID-named directories recorded in that private ledger. Unrelated shared-group
files are not enumerated or removed. A failed-attempt screen offers one ordinary
Share error details action containing the safe numeric report; there is no
debug-configuration menu or automatic transmission.
