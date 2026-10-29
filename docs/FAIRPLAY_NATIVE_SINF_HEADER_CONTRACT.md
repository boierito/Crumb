# Native SINF header decision and mobile work-buffer agreement

The TT and DD captures, their exact original encrypted IPAs, and their exported
account states were tested locally. No private fixture, identifier, credential,
key, private digest, or application bytes are included in this repository.
No Apple request was made during the recovery/reference tests.

## What the failing branch checks

For this execution of the pinned 26A428 binary:

1. The block at `text+0x4ccdf4` decodes four bytes from a native buffer. Its
   four loads are at `4cce00`, `4cce10`, `4cce1c` and `4cce30`.
2. The instructions at `4cce3c..4cce4c` transform the decoded big-endian word
   `W` into `W XOR 0xce757bc5`. The transformation is bijective.
3. The expected word at `4cb590..4cb5bc` depends on the actual stack scalar.
   In the real TT and DD executions its decoded value is `0x6974756e` (`itun`).
4. The comparison at `4cb5c0` is false. Execution follows the native target
   `8c7ed8`, stores the encoded error, and ultimately returns -42057.

The comparison is observed, never changed. This identifies this occurrence's
header rejection; it does not assign a universal symbolic name to -42057 or
exclude other paths returning the same status. It is not an Apple HTTP error.

## Independent verification with real inputs

The entire recorded invocation was re-executed using the original pinned ARM64
bytes in Unicorn. Registers, flags, SP, PC, native memory-write coverage and
written-page contents matched at all **293** recorded service boundaries and
at the final return. The replay uses recorded service replies and the portable
PAC convention; neither service semantics nor PAC authentication are thereby
independently established.

A separate optional observer reads the same native-frame material used by the
production candidate callback at `text+0x234cdc`. It applies the public pinned
key mask, then computes the mobile work buffer with Python cryptography:

```
effective_key = native_material XOR public_model_key_mask
iv = MD5(name_before_NUL || visible_iviv)
ciphertext = priv[:-8] if len(priv) == 440 else priv
work_buffer = AES-CBC-Decrypt(effective_key, iv, ciphertext)[:432]
```

This is the same effective AES key as the documented mobile relation
`credential = model_key XOR 0x30`, without exporting a credential. It does not
search or guess keys. The independent mobile implementation is from
EthanArbuckle/fairplay-decrypt `core.py`, already pinned in the mobile consumer
README and used to generate the public vectors.

In both executions there were two distinct candidates. In DD the **second**
produced exactly the same **432 bytes** as the native buffer at the header
decode. That buffer still failed the native `itun` comparison. This establishes
agreement of this candidate extraction and mobile work-buffer computation with
native execution for this input. It does not validate the SINF or establish a
reusable account credential.

Consequently, replacing the header validator, changing cryptid, accepting a
zero status, or increasing the execution budget is unsupported. The execution
returned normally with a negative result; it was not cut off by the budget.

## Reproduction

Use private, matching inputs already accepted by `verify-capture-ipa`. The
optional raw SINF must be extracted byte-for-byte from that same IPA, stored as
a regular private file (0600), and remain local. No key input is accepted.

```
CRUMB_NATIVE_REFERENCE=1 bash scripts/build-fairplay-host.sh /private/new-build
CRUMB_NATIVE_REFERENCE_TRACE=/private/new-trace.jsonl \
  /private/new-build/crumb-fairplay-host \
  -capture /private/capture.zip -ipa /private/original.ipa \
  -authorized-state /private/account-state.json -output /private/new-output
python3 scripts/replay-fairplay-native-reference.py \
  --trace /private/new-trace.jsonl --audit-sinf-header \
  --compare-mobile-sinf /private/original.sinf
```

The optional mobile comparison requires `cryptography==46.0.3`. Its report
contains booleans and counts only. It cannot authenticate the external SINF's
association on its own; that is a separate input-binding requirement above.
A match of one block must never authorize license or plaintext export.
Python cannot guarantee erasure of intermediate key copies. These observers
run only on the research host and are not part of the iOS application.

## Regression coverage and controlled DD result

The public tests execute the original ARM64 instruction sequences against
4,101 deterministic scalar/word cases, test missing observation handling,
and compare mobile work buffers against independent public vectors. Audit
reports explicitly deny license validation.

Unlike TT's earlier cold download, DD reports authorization restored before
KBSync generation. The retained producer imported its keybag with status zero
and saved its state. The exported state matches the consumer state. Restoring
DD's state in the original StoreAgent producer and querying provisioning
completed with `NeedsDiversity == false`. DD contained no diversity bag; its
absence is not evidence of a dropped response field in this transaction.
Nevertheless, DD's current-state recovery and recovery using the previous TT
authorized state both return -42057 with two rejected candidates. Therefore,
restoring authorization before download did **not** correct this case. No new
download or authorization is needed to reproduce it.

An additional native replay confirmed the file service returned all **869**
requested SIDB bytes and all **200356** requested SIDD bytes. The successful
reads do not establish semantic compatibility of those files with this engine.
The absent `SC Info.sidn`, `.sido` and `.sidv` lookups do not by themselves prove
any of those files is required. No replacement files were fabricated.

The host-only candidate origin observer is documented in
[FAIRPLAY_CANDIDATE_PROVENANCE.md](FAIRPLAY_CANDIDATE_PROVENANCE.md).
Its conservative categories are not an exact backward slice or a demonstration
of which input causes the wrong material. A successful reference replay still
uses recorded services and the portable PAC convention. No SINF correction,
valid account credential or decrypted IPA has been demonstrated by these tests.
