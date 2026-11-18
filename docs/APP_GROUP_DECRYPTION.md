# App Group decryption (experimental)

This branch replaces the normal decrypt flow with App Group staging, following
ipadn v2.2.0's approach. No ipadn source is copied. Crumb reuses its existing
read-only cryptid=2 FairPlay implementation, bounded thin-arm64 parser, full page
reads, entropy/code checks and cancellation. FAT/32-bit targets remain unsupported.

Sign Crumb in KSign with an App Group enabled in BOTH the app's signed entitlements
and its provisioning profile. The runtime discovers concrete group identifiers in
the embedded profile and asks iOS for the authorized container. If KSign does not
embed a profile listing the concrete group, set `FairPlayAppGroup` in Info.plist to
the actual authorized identifier when signing. Inventing a group identifier does
not authorize access. No fixed group from another developer is used.

From Downloads choose Decrypt, or Decrypt and install. The latter uses the existing
local signing identity and OTA server. Installing an encrypted original uses the
same original OTA method. None of these paths uses pairing, a helper or VPN. The
HTTPS manifest generator sees metadata and a tokenized loopback URL; IPA bytes
stay local. OTA cannot replace an existing App Store copy of the same bundle ID.
Export the verified output to KSign if preferred.

Keep Crumb open. A background expiration requests cancellation; the worker keeps
its workspace until it stops. The original IPA is retained on every failure.
The full archive directory is validated before extraction; links, traversal,
duplicate paths, excessive sizes and entry counts are rejected. Source CRCs and
repacked CRCs are checked, and repacked files are compared against the verified
plaintext tree before a download record is published. Original non-Payload files
and license resources are retained; missing nested SINF files are replicated.

FairPlay needs the device account licensed for the package and can reject remapping
on any OS. An accessible App Group is a prerequisite, not a guarantee. There is no
artificial iOS 27.0.1 gate on this experimental route. Physical-device results are
not implied by a successful build. The old native helper implementation is retained
as source for regression tests but is not entered by the normal UI.

Reference: https://github.com/mengangs/ipadn/tree/v2.2.0

## KSign group discovery

The app reads effective App Group entitlements from the running signed process
through dynamically resolved Security entry points, then tries explicit configuration
and concrete profile groups. If the runtime reader is unavailable, the latter two
sources remain supported. Wildcard profile groups cannot name an actual container.
Every candidate must also pass an actual write probe before decryption starts.

The decrypter settings screen contains normal usage guidance and signing settings.
Diagnostic readouts, diagnostic exports and the temporary manual group editor have
been removed after successful device testing. Existing saved overrides and the
Info.plist override remain recognized for signing compatibility. Neither override
can grant App Group access.
