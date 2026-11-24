# App Group decryption

Sign Crumb with an authorized App Group present in both the app entitlements and its provisioning profile. Crumb discovers authorized containers from the effective entitlements and concrete profile groups. An optional `FairPlayAppGroup` Info.plist configuration must name an actual authorized group; it cannot grant access.

From Downloads, choose Decrypt or Decrypt and install. Installation of decrypted IPAs uses the local signing identity and OTA workflow. Encrypted originals use the original OTA workflow. No external helper, pairing or VPN is required.

Keep Crumb open. Cancellation and failures retain the original IPA. App Groups and a device account licensed for the app are prerequisites, not guarantees of successful processing. Physical-device acceptance must be verified separately.

Crumb uses the iPadN FairPlayCore engine; see [third-party notices](../THIRD_PARTY_NOTICES.md) for provenance and licensing.
