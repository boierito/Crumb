# Generic local system FairPlay pipeline

This implementation removes the per-target Mac/build-server preparation step
from the FairPlayed experiment. An ordinary iOS app cannot replace its own
installed, signed resources in place. Crumb therefore prepares and signs an
update containing the purchased target, asks iOS to install that update, and
uses the system FairPlay reader from the installed bundle. One target is
processed per installed update. No new Crumb source build is needed per target.

## Runtime sequence

1. Download with Crumb's existing account flow, or select an original IPA already
   in Files. The original download remains untouched.
2. From Downloaded apps, choose **Decrypt with iPhone FairPlay**. Alternatively,
   open Settings → **FairPlay on iPhone**. Preparation snapshots only Crumb's
   installed bundle, never its data container, keychain, account state or caches.
   A previous `CrumbSystemPayload` is excluded. The selected target and its
   original SINF files are added as resources.
3. Configure a locally held P12 certificate and provisioning profile matching
   the installed Crumb's bundle ID and signing team. These are app-signing
   credentials, not FairPlay decryption keys. The password is not saved. The
   files stay in protected, backup-excluded application support storage.
4. **Sign prepared Crumb locally** uses the pinned MIT zsign library and static
   OpenSSL 3.5.7 in process, without shell commands, JIT or a signing server.
   It checks the key/certificate pair, certificate expiry, profile certificate
   membership, profile expiry, bundle ID and current signing team. It does not
   bypass Apple's certificate trust, device eligibility or revocation decisions.
   Every purchased-payload byte is checked against the prepared package after
   signing; added, missing or changed payload files reject the package.
5. **Request installation** serves the signed update only on loopback. It retains
   Crumb's existing HTTPS manifest mechanism: api.palera.in receives the helper's
   app metadata and loopback URL, not IPA bytes, licenses or signing credentials.
   iOS must accept the signature/profile and the user must confirm installation.
   Serving an IPA or creating a local signature is never reported as installation.
6. Open the updated Crumb → Settings → **FairPlay on iPhone**. The installed
   target's bundle ID/version/build are shown. A different selected download
   cannot be processed using the previously installed target.
7. **Decrypt all supported executables and prepare export** inventories every
   Mach-O, including frameworks, extensions and FAT slices. Encrypted supported
   ARM64 slices use a read-only mapping and `mremap_encrypted(cryptid=2)` followed
   by fault-reporting self-process page reads. No other process is attached.
   Every protected page must read successfully and change from its ciphertext.
   Unknown encrypted architectures, malformed files and missing outputs reject
   the complete operation; no partial IPA is accepted.
8. The shared Go verifier requires unchanged file sizes and byte equality outside
   the original encrypted ranges and cryptid fields. It then performs the cleanup
   below and creates a fresh unsigned IPA. The result requires its own signature
   before installation; preparation/signing packages remain private because they
   intentionally retain the original license.

The iPhone's system Apple Account must hold the target license. Logging into
Crumb alone does not provision the system FairPlay account. Native failure does
not trigger guessed keys, modified SINF data or fabricated provider responses.

## License association

All declared original SINF paths must exist. Existing per-component licenses are
preserved. Missing manifest-directed replicas can be populated from the original
primary SINF when there is one license, or when all original license blobs are
identical. Distinct licenses with an unresolved replica mapping are rejected
before installation. Filename similarity is not treated as cryptographic proof
of a license association. Preparation never reports a validated license.

## What the unsigned export removes and checks

* All `SC_Info` directories, `_CodeSignature` directories, embedded provisioning
  profiles, iTunesMetadata plists, `_MASReceipt` and StoreKit receipt files.
* Every supported Mach-O's `LC_CODE_SIGNATURE` command and the entire referenced
  blob. Commands are compacted only within the existing header area, and vacated
  header/signature bytes are zeroed. File sizes, segment offsets and recovered
  protected bytes are preserved. Certificate chains and signing entitlements
  embedded in those blobs are therefore not carried into the unsigned export.
* Original ZIP comments, extra fields, ownership metadata and timestamps: fresh
  headers are emitted instead of copying personalized archive headers.
* Known account identifiers, when available from the current Crumb session, are
  checked in resource paths and streamed file contents, including UTF-8,
  UTF-16LE/BE and Base64 representations, across buffer boundaries. A match
  rejects the export; it is not silently rewritten in executable code.

`known_identifiers_checked` reports the number of inputs actually available. A
signed-out session cannot be reported as having checked identifiers it does not
possess. Reports contain counts and flags, not those inputs or private bytes.

This is a bounded metadata/identifier check, not a proof that an arbitrary app
contains no obfuscated watermark, server-defined personalization or unknown
account data. `no_account_data_guaranteed` remains **false**. It must not be
relabelled "100% anonymous" or "public-safe". Removing licensing data also does
not establish a right to redistribute an app. Capability entitlements may need
to be supplied appropriately during re-signing; launch and app/server behavior
are separate checks.

`plaintext_authenticated` remains **false**: the real DD comparison established
that its original CodeDirectory hashes cover ciphertext, not recovered text.
Neither clearing cryptid nor changing a page establishes plaintext authentication.

Public base-build validation rejects embedded target directories, SINF/SIDB/SIDD
files, P12 files and provisioning profiles. User captures and generated packages
are never repository or CI inputs. Existing project redistribution/licensing
limitations in THIRD_PARTY_NOTICES.md remain applicable.

## Verification and remaining device check

The package module and host CLI use exactly the code linked into Crumb:

```
go -C MapleSyrup/NativeSAP/packageipa test -race ./...
go -C MapleSyrup/NativeSAP/packageipa run ./cmd/systemipa prepare --source <Crumb.app> --ipa <original.ipa> --output <private-helper.ipa>
go -C MapleSyrup/NativeSAP/packageipa run ./cmd/systemipa plan --source <installed-target>
go -C MapleSyrup/NativeSAP/packageipa run ./cmd/systemipa export --source <installed-target> --overrides <native-files> --output <unsigned.ipa>
```

Deterministic tests cover multiple executables/domains, two FAT slices, unchanged
and incomplete pages, changes outside protected ranges, malformed commands,
overlapping slices, unsupported encrypted CPUs, metadata removal, signature-byte
zeroing, identifier matches across stream boundaries, exact license replication,
ambiguous license rejection, private target exclusion and post-signing payload
preservation. Unsupported hosts cannot return a successful native read.

The signing tests generate an ephemeral test CA/certificate unrelated to Apple
or FairPlay. They independently verify CMS signatures with OpenSSL and recompute
all code-page hashes, reject wrong passwords/bundle IDs/signing teams, and check
the target remains byte-identical. This proves the local signature construction
contract, not Apple trust or an accepted iPhone installation.

For DD, the shared export code was run against the real recovered executable
bytes: 2 executables, 178 changed 16 KiB pages, 43 output files, 2 removed embedded
signatures and one available known identifier checked. Input signing areas were
reconstructed to account for the previous helper's extension re-signing; this
was a replay/cleanup test, not a new native decryption run. Private bytes stayed
local. The actual native decryption proof remains the returned FairPlayed run
documented in FAIRPLAYED_SYSTEM_ROUTE.md.

The new integrated Crumb workflow still needs one complete device run with a
valid same-team signing identity: prepare → sign → accepted iOS update → native
reads → cleaned export. Builds and synthetic signing fixtures do not substitute
for that check. The current native route has real device evidence on iOS 27.0.1;
other iOS versions and encrypted unsupported architectures are not certified.
