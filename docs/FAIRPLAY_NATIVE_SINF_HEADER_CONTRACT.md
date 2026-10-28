# Native SINF header decision and mobile first-block agreement

The TT capture, its exact original encrypted IPA, and the subsequently exported
account state were tested locally. No private fixture, identifier, credential,
key, private digest, or application bytes are included in this repository.
No Apple request was made during the recovery/reference tests.

## What the failing branch checks

For this execution of the pinned 26A428 binary:

1. The block at `text+0x4ccdf4` decodes four bytes from a native buffer. Its
   four loads are at `4cce00`, `4cce10`, `4cce1c` and `4cce30`.
2. The instructions at `4cce3c..4cce4c` transform the decoded big-endian word
   `W` into `W XOR 0xce757bc5`. The transformation is bijective.
3. The expected word at `4cb590..4cb5bc` depends on the actual stack scalar.
   In the real TT execution its decoded value is `0x6974756e` (`itun`).
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
key mask, then computes the mobile first block with Python cryptography:

```
effective_key = native_material XOR public_model_key_mask
iv = MD5(name_before_NUL || visible_iviv)
first_block = AES-CBC-Decrypt(effective_key, iv, priv[0:16])
```

This is the same effective AES key as the documented mobile relation
`credential = model_key XOR 0x30`, without exporting a credential. It does not
search or guess keys. The independent mobile implementation is from
EthanArbuckle/fairplay-decrypt `core.py`, already pinned in the mobile consumer
README and used to generate the public vectors.

In the TT execution there were two distinct candidates. The **second** produced
exactly the same first 16 bytes as the native buffer at the header decode.
That buffer still failed the native `itun` comparison. This is evidence that
this candidate extraction/mobile first-block relation agrees with native
execution for this input. It does not validate the SINF, prove full-buffer
agreement, or establish a reusable account credential.

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

## Regression coverage and remaining experiment

The public tests execute the original ARM64 instruction sequences against
4,101 deterministic scalar/word cases, test missing observation handling,
and compare mobile first blocks against independently generated public
vectors. Both audit reports explicitly deny license validation.

TT's download log says `account state restored=false`. Authorization was
performed after that transaction. Retrying recovery against the saved account
state is not equivalent to producing the original transaction's KBSync and
keybag with a provisioned producer. The already implemented retained download
factory restores saved authorization before generating KBSync and consumes the
response with that same producer. A genuine response from that context has
not yet been tested in the TT evidence. This is the next controlled licensing
experiment, not a promised correction to -42057. It needs no new app build.

If that transaction imports successfully and still fails the same native
header check, the result must be recorded as a rejection of that route. It
must not trigger repeated downloads or inferred credential guesses.
