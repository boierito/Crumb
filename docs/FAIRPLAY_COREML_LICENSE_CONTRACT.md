# Core ML license handoff: native follow-up

Date: 2026-10-30. Status: research verification, **not an IPA decryption fix**.

The direct Crumb reader reached `mremap_encrypted` and failed before reading any
page in the returned v2 device report. Core ML was investigated as a possible
alternative broker for post-installation downloads. It must not be called with
a guessed UUID, truncated SINF or invented key blob.

## Original firmware inputs

The original `CoreMLModelSecurityService` executable and bundle metadata were
acquired from the following public Apple restore images, rather than compiled
from the published decompiler output:

| OS / build | Original service SHA-256 | Identity helper |
| --- | --- | --- |
| 18.2 / 22C152 | `8820f7af22ea0742f108a358bbb709d5ec856d037008121d1af1807e2f757193` | `0x10000a24c` |
| 27.0.1 / 24A446 | `07b277507b80a27b3aeb5c954721f3223b657cd3c5b551e45fbb2f8d4afb4e5a` | `0x100006c14` |

- [Apple iPhone17,1 / 22C152 IPSW](https://updates.cdn-apple.com/2024FallFCS/fullrestores/072-43747/F8E440DF-FA15-49A6-B511-C155E67CD405/iPhone17,1_18.2_22C152_Restore.ipsw)
- [Apple iPhone17,1 / 24A446 IPSW](https://updates.cdn-apple.com/2026FallFCS/59ffd938-75b6-41e5-af04-f4295ee8aeb3/iPhone17,1_27.0.1_24A446_Restore.ipsw)

Path inside each filesystem image:
`/System/Library/Frameworks/CoreML.framework/XPCServices/CoreMLModelSecurityService.xpc/CoreMLModelSecurityService`.
The filesystem image is the BuildManifest `OS` component, not
`Cryptex1,SystemOS`. The acquired members are `092-08959-201.dmg.aea` and
`043-70165-666.dmg.aea`, respectively.

A local HTTP-range/AEA/APFS reader verified the root and cluster-header HMACs,
and the segment HMAC and SHA-256 for every downloaded filesystem segment. It
retrieved approximately 17.2 MB and 22.4 MB, respectively, rather than the whole
archives. This establishes integrity of the bytes read under Apple's public
firmware archive keys; it is not a claim that every archive byte was fetched or
that the service's CMS chain was independently verified. No account, purchased
IPA or private capture was used or transmitted during acquisition.

Reacquisition can use the pinned `blacktop/ipsw` 3.1.739 extractor with `--remote
--files --pattern` for the path above; host support for the filesystem mount or
APFS reader is required. Acquisition is a development/research operation and is
not introduced as a runtime dependency of Crumb.

## What the original code establishes

`scripts/audit-coreml-identity.py` accepts only the two exact service hashes. It
parses the Mach-O, resolves the original selector stubs and reads the signed
entitlement blob. It then runs the original identity-choice helper in Unicorn
for the four combinations of identity mode and modeled identity presence.

The Objective-C/ARC calls are modeled explicitly, logging is modeled as
disabled, and the PAC stack instructions are modeled. Unknown external
instructions fail the run. This is a control-flow comparison, **not execution
of Security.framework, the XPC service, the key server or FairPlay cryptography**.

The eight cases across the two builds confirm:

- With `usesCodeSigningIdentityForEncryption=false`, the helper selects
  `teamIdentifierFromXPCConnection`.
- With the flag true, it selects `codeSigningIdentifierFromXPCConnection`.
- A modeled missing identity causes the original helper to write an error and
  return nil. A modeled available identity is returned without an error write.

On 24A446, static disassembly additionally shows
`teamIdentifierFromXPCConnection` at `0x1000073b0` obtaining `currentConnection`,
`processIdentifier` and `auditToken`, then invoking
`teamIdentifierFor:processIdentifier:` at `0x100006fec`. That method calls
`SecTaskCreateWithAuditToken` at `0x100007028` and
`SecTaskCopyTeamIdentifier` at `0x100007040`. The Apple Account used by Crumb's
download is not the signing identity selected by this code.

The audit also verifies the original **static call sites**, in order, for:

1. `generateKeyRequestForKeyIdentifier:teamIdentifier:error:`;
2. `fetchKeyResponseFromServerForKeyID:signedKeyRequest:error:`;
3. `generatePersistentKeyBlobFromKeyResponse:error:`.

It independently resolves the request's `setKeyId:`, `setSignedKeyRequest:`,
`setTeamId:` and `setRawRequest:` call sites. Immediately before the last call,
both originals set w2 to zero. No actual request is sent and no response is
fabricated; locating a call site does not assert execution of all its branches.

The original service's entitlement blob contains
`com.apple.private.fpsd.client=true`,
`com.apple.private.coreml.decrypt_allowed=true`,
`platform-application=true` and `fairplay-client=703697`. These are the service's
rights. Their presence **does not prove that an ordinary Core ML client needs
those rights**, or that it can acquire them by re-signing a copy of the service.

## Credential prerequisite still missing

The corresponding published 22C152 framework decompiler output shows the
model-key identifier parsed as a UUID, a persistent-blob cache, a maximum blob
length of 136 bytes in MLFairPlayDecryptSessionManager, and separate model
archive/encryption metadata. These framework details have not yet been
independently verified in the original CoreML framework code. The service audit
above does not silently promote them to native instruction-level evidence.

Attempted selective extraction of the CoreML image from the separate cryptex
stopped because the local APFS reader rejected its root btree. No missing bytes,
object map entries or pointer tables were filled with substitutes. This limits
the native verification to the acquired XPC services, not the entire framework.

The missing producer remains concrete: a valid response accepted by
`generatePersistentKeyBlobFromKeyResponse:error:` for a model key compatible
with the **original purchased app's license**. Crumb's retained transaction
provides an App Store keybag and SINF, but no inspected operation demonstrates
that they produce this response or persistent blob. Sharing a name, 16-byte
identifier size, cryptid value or FairPlay subsystem does not establish the
cryptographic conversion.

The separate `secureModelWithContentsOfURL:...` API consumes a file-access
extension in the Core ML service and loads an MLModel; it is not evidence of
forwarding permission to fairplayd or exporting arbitrary Mach-O plaintext.
No such forwarding or export has been verified by this follow-up.

There is therefore no justified Core ML runtime fallback to integrate yet.
This is a blocker for this proposed handoff, not proof that every iOS-only
route is impossible. No new phone build, download, login or private key is
requested to repeat the already observed remap failure.

## Reproduce the local audit

Install Unicorn in the research environment, acquire either original service,
then run:

```sh
python3 scripts/audit-coreml-identity.py --service /private/research/CoreMLModelSecurityService
```

Results for the original two inputs are stored in
[`research/coreml-contract-22C152.json`](research/coreml-contract-22C152.json)
and [`research/coreml-contract-24A446.json`](research/coreml-contract-24A446.json).
They contain public service metadata and modeled-control-flow results only.
They explicitly report `license_conversion_verified=false` and
`decryption_verified=false`.

Validation: all eight original-code replay cases passed; both static producer
and request-field inventories passed; an altered original image was rejected
by the hash gate. No application runtime code or success gate was changed.
