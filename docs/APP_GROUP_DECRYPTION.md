# App Group decryption (experimental)

This branch replaces the normal decrypt flow with App Group staging, following
ipadn v2.2.0's approach. No ipadn source is copied. Crumb reuses its existing
bounded thin-arm64 parser, full page reads and cancellation. The FairPlay path
now follows the v2.3.0 entropy verdict and tries cryptid=2/read-only and
cryptid=1/read+execute, each with and without registering the original signature.
Ciphertext is snapshotted before any attempt into unlinked temporary storage.
A successful syscall, complete page reads, changed bytes, source entropy above
7.9, plaintext entropy below 7.5 and a drop of at least 0.5 remain mandatory.
ARM64 instruction landmarks are no longer a hard requirement. FAT/32-bit targets remain unsupported.

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


Reference reviewed on 2026-11-20: latest GitHub release iPadn v2.3.0,
commit 92540912281b1d11c57d0e00e114506ec459a494. Its updated verdict removes the
RET/PACIBSP requirement for data-only/obfuscated encrypted regions. This is an
original adaptation using Crumb's existing reader; no upstream source is vendored.
The unsupported FAT/32-bit restriction remains in this scoped correction.
Failure messages identify the relative file and attempted operation; no full
sandbox paths, account credentials or permanent diagnostic UI are added.
