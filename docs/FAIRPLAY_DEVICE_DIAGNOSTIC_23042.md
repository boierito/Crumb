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
