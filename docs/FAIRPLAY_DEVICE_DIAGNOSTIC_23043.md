# Crumb 23043: bootstrap continuation exercised on iPhone

User-supplied fixed labels and aggregate counts only; no private capture data.

| Field | 23042 | 23043 |
| --- | --- | --- |
| phase / failure | bootstrap / raw_unknown_control_flow | bootstrap / raw_unknown_control_flow |
| last translated block | text+0xa4be38 | text+0x8c292c |
| next unsupported target | text+0x8c292c | text+0xa06948 |
| shared block executions | 65 | 66 |

Recovery took 176 ms. SINF and SUPF lookups each succeeded once. Native keybag
import was accepted (status 0), with one nonempty legacy SIDB of 878 bytes.
Portable source/frame/model validation remains false. No processed IPA exists.

The 23043 continuation is exercised in the device runtime. This is code coverage
progression, not decryption. A four-way repeat is not needed to identify this
next missing destination. The next task is exact-binary translation of
0x100a06948 and examination of its continuations before another device build.
