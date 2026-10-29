# Preserve the complete download license reply

The retained download previously extracted only the first keybag field and
ignored `diversitybag`. Authorization already imported that second bag through
a different native routine. The pinned StoreAgent purchase handler also has
separate keybag and diversitybag imports; the latter is not a keybag alias.
Dropping a supplied diversitybag is a concrete dataflow omission, independent
of whether it explains a particular SINF rejection.

The laboratory download now retains the two bag kinds separately from the
same validated ent response. Repeated identical aliases are accepted; different
values for the same kind, invalid encodings and values exceeding the native
1 MiB input bound are rejected. The download consumer still requires its
keybag; this does not invent a license for a response lacking that input.

The C bridge binds both bags to the exact retained KBSync producer. Native
consumption holds the producer lock across keybag then optional diversitybag
import. It stops on the first rejection, releases the intermediate snapshot,
and returns the final snapshot plus distinct native statuses. The bridge saves
state only after every supplied import completed with status zero and the
existing state-version check passes. Both retained bag copies are wiped after
consumption or close. The old keybag-only ABI remains available.

The private capture includes the diversitybag when supplied. The reader checks
its manifest digest and rejects missing/contradictory import evidence before
attempting recovery. No raw licensing bytes enter the diagnostic. SINF and
whole-IPA export gates remain mandatory.

## Real-input limits

TT's exact IPA/capture association and local reference tests were reused without
contacting Apple. The original StoreAgent session-initialization operation,
which accepts dpInfo for Mac packages, rejected the unmodified mobile SINF
with native -42001. That API has not been repurposed as an iOS decryptor.

Extending the existing native/mobile comparison from one block to the complete
432-byte TT private-record work buffer produced an exact match for candidate
2. The buffer still failed the native `itun` predicate and the invocation
returned -42057. This is agreement on a rejected decryption, not key recovery.
Services still use recorded replies and PAC follows the portable convention.

A read-only inspection of 33 live allocations after the original producer
restored authorized TT state and queried provisioning found no bounded
20-byte `sean/user/key /scra/sess/valu` records in that allocation subset.
This is not evidence that licensing state is absent elsewhere or that a
credential cannot be obtained. No protected bytes were changed by the scan.

TT's old capture has no original response envelope or diversitybag blob.
Consequently it cannot establish whether the server supplied one and the old
client discarded it. The new phone log (920-byte keybag, accepted import, two
rejected SINF candidates) likewise does not establish bag presence or source
identity without its matching private capture. This change is not claimed to
resolve that transaction's -42057 or to produce a decrypted IPA.

## Verification

Public deterministic tests cover separate bag kinds, duplicate/ambiguous
aliases, limits, exact request/reply capture, native import order, both failure
stages, status independence, intermediate-snapshot cleanup, keybag-only
compatibility and private-capture consistency. The native helper can be tested
without credentials or Apple access:

```
go test -race MapleSyrup/NativeSAP/kbsynclicense.go \
  MapleSyrup/NativeSAP/kbsynclicense_test.go
```

Native producer, Swift protocol, portable capture and iOS build checks run
through the normal CI workflows. Private captures and Apple secrets must not
be supplied to CI. Successful compilation does not validate the original SINF.
