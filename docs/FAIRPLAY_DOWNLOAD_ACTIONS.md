# Download actions and local signing, build 23055

The original catalog/account/download flow remains intact. After download and
in Downloads, the same four actions are available:

- Save encrypted IPA: share the original saved archive.
- Save & install encrypted: retain the original and request the original OTA installation.
- Decrypt IPA: prepare and sign a private Crumb update using the locally imported
  identity, request its installation, reopen Crumb and process its installed target.
- Decrypt & install: the same path, then sign a separate installation copy of
  the verified plaintext IPA and request its installation. The unsigned cleaned
  export stays separate. A profile allowing the output bundle ID is required;
  an optional separate app profile is supported. P12 password is in memory only
  and must be re-entered after restarting for the final signature.

The Crumb-update step is a requirement of the demonstrated system route, not
hidden by the UI. Installation is an iOS decision; serving bytes or creating a
CMS signature is not installation proof. Signing an app does not guarantee its
extensions/entitlements can work with an arbitrary profile. Profile trust,
revocation and device eligibility remain system checks.

A retained random job ID binds automatic resume to the installed update. A failed
or interrupted native run requires an explicit retry rather than silently
retrying on each launch. The existing real DD page evidence remains valid;
the newly integrated end-to-end build still needs device validation.

## Signing and import

Settings → Signing identity is shared with download actions. Picker imports use
UIDocumentPicker copy mode. The supplied FilePickerFix ARM64 binary is pinned,
embedded and loaded before the picker; see its third_party README for bounded
inspection and unresolved redistribution licensing.

Files → On My iPhone → Crumb → Signing Inbox is an alternative: copy one .p12
or .mobileprovision there and select its role in Settings. On successful private
storage, only the inbox copy is deleted. Other selected original files are not
intentionally removed by Crumb. The compatibility dylib has its own file-moving
behavior, so native provider-copy mode is used.

Preflight checks bounded P12 input, matching private key, certificate dates,
profile CMS signature integrity, profile expiration, certificate membership,
Crumb's current team and exact/trailing-wildcard application identifier. It does
not assert Apple trust. OpenSSL's built-in legacy provider supports older RC2
P12 exports. Errors contain safe classifications, not raw certificate details,
paths, private keys or password values. Import does not equal validation.

Private license-bearing helpers, signing identities, final signed copies and
unsigned cleaned exports remain distinct. Unknown app-defined identifiers and
independent plaintext authenticity are not certified. The unsigned export's
existing all-page checks remain mandatory.
