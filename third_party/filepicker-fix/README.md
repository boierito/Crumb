# User-supplied FilePickerFix

`FilePickerFix.dylib` is the ARM64 binary explicitly supplied for this integration.
SHA-256: `4af2e076f2377ed64ddabe9f2dce32edca3d43db624b4707214fc0f1dd71d6eb`.

Crumb embeds it in Frameworks and dlopens it before creating the document picker.
It is signed with Crumb by the user's signer. Its original install-name is not
used for loading. No MobileSubstrate dependency is imported.

Static inspection finds Objective-C document-picker/NSURL swizzles, a file-moving
helper and path logging. Imports reference Foundation, CoreFoundation,
CoreServices, libobjc and libSystem; no network framework or socket/HTTP imports
were found. This is bounded inspection, not a source audit or guarantee of
behavior. Crumb requests provider copies and stores imported identity files in
protected, backup-excluded storage. Only the P12 private key is used in memory;
it is never intentionally included in a prepared helper or exported app.

The supplier did not include source code or a redistribution license. Permission
for a general public redistribution of this binary has NOT been established.
This dependency must be licensed or replaced with an appropriately licensed
implementation before describing the combined build as cleared for public
redistribution. The compatibility behavior still needs device verification.
