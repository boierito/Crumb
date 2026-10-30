# First returned direct-runtime report, build 23057

The user's safe report contains one `data_container` attempt: return -1, errno
12, zero pages, and zero usable shared locations. It requests neither signing
nor installation. The input IPA/transaction is not identified in that report;
do not assume it was DD or correlate it to DD's prior successful embedded run.

## What it proves

- No page was accepted as decrypted and no verified output was produced.
- No shared-container native attempt ran. A zero discovery count cannot be
  interpreted as a failed decryption in a shared container, or as proof of
  exactly which entitlement/container-discovery operation was unavailable.
- The aggregate `preparation_or_validation` label was misleading: the report's
  actual attempt came from the native pipeline. It is not another observed
  portable-recovery -42057 result.

## What it does not prove

The 23057 C reader initializes `remap_result` to -1, and also writes
`remap_errno` if `mmap` fails before the syscall. Therefore this v1 report
cannot distinguish memory-mapping failure from native remapping failure.
The first interpretation that mremap had definitely run was too strong.

Even **if** mremap ran, errno 12 does not identify the FairPlay service cause.
The actual public iPhone17,1 iOS 27.0.1 / 24A446 kernel's syscall table resolves
entry 489 to `0xfffffff00b0d8a34`. At `0xfffffff00b0d8ca8` through
`0xfffffff00b0d8cb8`, its error branch returns EPERM for `0xe00002c1`, and
ENOMEM for other nonzero crypter-creation returns. Generic error, true no-memory
and not-permitted cases therefore collapse to the same errno.

`scripts/audit-fairplay-ios-errno.py --kernel <extracted kernel>` verifies the
pinned original kernel bytes and executes those five instructions independently
in Unicorn for four modeled callback returns. It does not measure the callback
return on the user's phone or prove which operation failed there.

Kernel SHA-256:
`fa07057bf34f33af1ec26681ad0c87c13e5aaa548a300740eddb0e92876e5f95`.
Acquisition: the official Apple IPSW and extraction described in
`FAIRPLAY_NO_REINSTALL_RESEARCH.md`. No private input is used by the audit.

## Implementation correction

The reader now records mapping-attempted, mapping-errno and remap-attempted
separately. The Swift stage classifier preserves legacy reports as ambiguous,
distinguishes mapping, remap and page-read failures, and marks aggregate native
attempts as such. Reports use schema v2; secrets/paths remain excluded. No
execution budget, cryptid policy, license data or decryption-success gate changes.

This corrects the observation contract; **it does not fix the native rejection**.
No new IPA is requested solely to obtain more diagnostic fields.

## Signature registration review

Crumb's direct reader currently does not invoke F_ADDFILESIGS_RETURN. Upstream
FairPlayed's `fp_probe_file` tries both registration modes for each mapping mode.
Consequently adding that call alone is not established as a causal fix for the
reported data-container restriction. No manufactured signature, re-signed input
or speculative registration fallback is presented as resolving this report.

The compatible explanation remains a native loader/service failure, potentially
including the upstream-reported file-access restriction. The exact cause is
unresolved; this report alone does not prove insufficient RAM, bad Apple Account
authentication, an invalid SINF, or a universal impossibility of other routes.

## Returned 23058 / v2 result

The next user report explicitly records `mapping_attempted=true`,
`mapping_errno=0`, `remap_attempted=true`, `remap_return=-1`,
`remap_errno=12`, and stage `native_remap`. This resolves the earlier
measurement ambiguity: mmap succeeded and the native remap failed before any
page was read. It still does not identify the crypter callback's original error.
There were zero authorized shared locations and no shared-location attempt.
The report does not identify the purchased app; it must not be treated as a
matched DD location experiment.

## Core ML broker follow-up, 2026-10-30

Header declarations alone were insufficient. The investigation additionally
inspected published **decompiler output for iOS 18.2 / 22C152**, pinned to
EthanArbuckle/iPhone17-1_18.2_22C152_Restore commit
`e26ed4563f78871c59d2d96856756a65d62517e5`. This is secondary reverse-engineering
evidence for that build, not verified Apple source or an executed iOS 27 test.
The output contains missing arguments, placeholder properties and incomplete
blocks; do not compile it or assume it is byte-accurate executable source.

The relevant call chain in that published representation is:

1. [MLKeyManager](https://github.com/EthanArbuckle/iPhone17-1_18.2_22C152_Restore/blob/e26ed4563f78871c59d2d96856756a65d62517e5/System/Library/Frameworks/CoreML.framework/MLKeyManager.m)
   extracts a key identifier from the model's `encryptionInfo` file. It looks
   for a persistent key blob in MLPersistentKeyStorage, otherwise asking
   CoreMLModelSecurityService to create one. This is not an App Store SINF parser.
2. [CoreMLModelSecurity](https://github.com/EthanArbuckle/iPhone17-1_18.2_22C152_Restore/blob/e26ed4563f78871c59d2d96856756a65d62517e5/System/Library/Frameworks/CoreML.framework/XPCServices/CoreMLModelSecurityService.xpc/CoreMLModelSecurity.m)
   shows key-request generation, a model-key-server response, and processing of
   that response into a persistent blob. Its identity helpers obtain signing
   identity from the XPC client's audit token. The intermediate identity-choice
   helper is not defined in this file; its full policy has not been verified.
3. [MLCloudSession](https://github.com/EthanArbuckle/iPhone17-1_18.2_22C152_Restore/blob/e26ed4563f78871c59d2d96856756a65d62517e5/System/Library/Frameworks/CoreML.framework/XPCServices/CoreMLModelSecurityService.xpc/MLCloudSession.m)
   names the model key store, `fetchKey2`, and request fields keyId,
   signedKeyRequest and teamId. No call was made to this server and no private
   transaction or license was sent to it. Nothing inspected establishes that
   it issues a model key for a purchased application's license.
4. [MLFairPlayKeyLoadingSession](https://github.com/EthanArbuckle/iPhone17-1_18.2_22C152_Restore/blob/e26ed4563f78871c59d2d96856756a65d62517e5/System/Library/Frameworks/CoreML.framework/MLFairPlayKeyLoadingSession.m)
   interprets the key identifier as a UUID and uses a stateful FairPlay session
   to process the key response. These 16 UUID bytes are an identifier, not
   evidence of the 16-byte App Store account credential.
5. [MLFairPlayDecryptSessionManager](https://github.com/EthanArbuckle/iPhone17-1_18.2_22C152_Restore/blob/e26ed4563f78871c59d2d96856756a65d62517e5/System/Library/Frameworks/CoreML.framework/MLFairPlayDecryptSessionManager.m)
   rejects a nil key blob or a blob of length >= 0x89, then copies at most
   136 bytes into its decrypt context, with a model-directory path and team
   identifier. A 440-byte private SINF cannot simply be supplied as that blob.
   This does not prove that every smaller byte string is valid or that a
   correct conversion is impossible; no such conversion or local producer has
   been demonstrated.

The explicit sandbox-extension consumer is
`secureModelWithContentsOfURL:sandboxExtensionToken:configuration:decryptCredential:withReply:`.
It consumes the token in the Core ML XPC service, loads an MLModel, returns a
secure-model object, and releases the token. The shown
`startDecryptionOfModelAtPath:usingKeyBlob:...` method takes no token. Neither
method shows forwarding an extension to fairplayd or exporting arbitrary
Mach-O plaintext. A permission consumed by this XPC service is not proof of a
permission consumed by fairplayd.

**Disposition:** this broker is not integrated as a speculative fallback. The
missing prerequisite is a demonstrated producer of a valid, account-license-
compatible model blob (or a different verified license handoff), plus verified
file-access behavior on the target OS. No truncated SINF, invented UUID/key,
zero blob, copied developer credential or successful RPC return is accepted as
decryption. These findings narrow this particular alternative; they do not
prove that all iOS-only routes are impossible.

No runtime correction or new decryption result follows from this follow-up.
The one-tap UI remains the desired interface; the native access/credential
prerequisite remains unresolved. Rebuilding the same reader would not test a
new causal hypothesis, so this investigation does not request another IPA
download or device run.
