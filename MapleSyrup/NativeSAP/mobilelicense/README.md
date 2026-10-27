# Mobile license consumer — not a local credential producer

This package implements only account credential → SINF-checked work material →
4096-byte candidate plaintext. It has no file loader, network access, device
access, credential extraction, Mach-O mutation, ZIP writer, or export operation.
It is deliberately not connected to an automatic decryption fallback.

The code was selectively adapted from Crumb's experimental
`feature/fairplay-dual-engine-23039` commit
`a1b919695ec9defe8c196db589232cf65abcf5b6` (londek MIT implementation).
The experimental IPA writer, credential file loader and UI were not imported.
`Open` checks exactly the supported private-record markers of the mobile
reference. This is not a signature or full plaintext authentication check.
`Matches` binds the handle to the exact SINF bytes. `Close` clears retained work
material, but makes no guarantee about compiler copies or expanded AES keys.

## Independent reference

- EthanArbuckle/fairplay-decrypt commit `7f80db2ef37634073fab06b3aad1aa0cf24561a1`
- `core.py` SHA-256 `9150ee023ae46846fdb9b929fc2da0c2f15488f299f8072aacd0d074dd9a21c4`
- Cross-checked with londek/ipadecrypt commit `849a1290afda0d74e608f80859119391dbab9e86`

`testdata/reference.json` contains **public synthetic** credentials and licenses,
not customer data. Regenerate from that exact Python core with:

```
python3 scripts/generate-mobile-license-vectors.py /path/to/fairplay-decrypt/core.py
go -C MapleSyrup/NativeSAP/mobilelicense test -race ./...
```

The generator requires `cryptography`. Normal Go tests need no Python or network
and compare work material, page keys and plaintext hashes with the reference.
They cover 440/448-byte private records and page indices through UINT32_MAX.

## Missing prerequisite

A legitimate 16-byte account credential derived by Crumb from its own transaction
is **not implemented**. Do not populate it by slicing the SIDB, trying random
keys, loading a credential from another device, or treating import status 0 as
key recovery. Matching the SINF header does not authenticate a whole IPA.
