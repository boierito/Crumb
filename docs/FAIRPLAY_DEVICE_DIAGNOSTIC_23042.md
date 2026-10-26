# Crumb 23042: iPhone bootstrap coverage progression

Device report supplied by the user. Only fixed labels, public code offsets and
aggregate counts are recorded. No capture, SIDB contents, identifiers or keys.

| Field | Previous device result | 23042 device result |
| --- | --- | --- |
| phase | bootstrap | bootstrap |
| failure | raw_unknown_control_flow | raw_unknown_control_flow |
| last translated block | text+0x8fff2c | text+0xa4be38 |
| next unsupported target | text+0xa4be38 | text+0x8c292c |
| shared block executions | 1 | 65 |

The local recovery attempt took 177 ms. SINF and SUPF were found without missing
lookups. Native keybag import was accepted (status 0), and one nonempty legacy
SIDB was observed (882 bytes). Source/frame/model validation remains false.
The original IPA was retained; no processed IPA was published.

The new translated block executed on iPhone. This proves runtime coverage
progression, not enrollment compatibility or decryption. The next engineering
task is to verify and translate 0x1008c292c from the exact daemon, then test its
semantics before another device build. A four-way repro is unnecessary to
identify this specific unmapped target.

## Implemented continuation (23043)

The exact hash-checked CI daemon confirms this is an executable block. Its
translation stores W24 at SP+0x11c, computes a signed table index from W12,
loads the next target, restores X21/X25 from X25/X27 and branches to that
loaded target. The source is included as `fairplayrawsupplement_exit.go`.

Independent ARM64 differential testing passed for 64 synthetic initial states,
comparing all registers, NZCV, SP, next PC and the entire stack page. The cases
include negative table indices, W-register overflow and varying SP. Focused Go
race tests passed for dispatch registration, full-state/stack preservation,
interior stopPC and unmapped stack/table rejection. As in the established
runtime, BRAA uses the portable BR target convention; native PAC is not tested.

Provenance: https://github.com/boierito/Crumb/actions/runs/40342839589,
`FairPlay-macos27-code-profile`, `next-block-8c292c.go.txt`. Binary SHA256:
`1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784`.

Build 23043 has new executable coverage. A successful build is still not
successful decryption; a device run is required to determine its next outcome.
