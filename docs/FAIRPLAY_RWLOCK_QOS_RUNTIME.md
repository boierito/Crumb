# Serialized guest pthread runtime: lock lifecycle and QoS

A real recovery on the owner's matching encrypted IPA and original private
capture reached a rejected SINF candidate and then failed at import slot 04c.
Public import inspection of the exact SHA-256-pinned daemon maps that slot to
`pthread_rwlock_destroy`. The corresponding object's signature was zero even
though this run had called `pthread_rwlock_init`: the earlier Crumb shim
returned success for init, read lock, write lock and unlock without mutating
memory. This is an emulator defect, not evidence of a valid mobile credential.

Replace those four no-ops together with a bounded, single-guest-CPU lock
lifecycle and the identified destroy operation. Fresh or static default locks,
default attributes, 8/16-byte alignment, repeated readers, an exclusive writer,
unlock, busy destroy and zero-signature EINVAL are modeled. The writer uses a
stable virtual thread identity for this single serialized guest; it is not a
hardware or account value. Unknown attributes, lock layouts, contended or
unobserved states and wraparound fail closed. Private guest pointers stay in
per-attempt runtime state and are cleared on cleanup.

On the owner's original capture, an experimental execution of this lifecycle
passed the old destroy failure and next reached slot 057,
`pthread_set_qos_class_self_np`. Model the validated QoS class and relative
priority in the virtual runtime, returning EINVAL for invalid combinations.
This records a guest scheduling property: it does not change host scheduling,
create executable memory, invoke privileged system services or supply a key.
No QoS getter is fabricated for an unknown thread.

Evidence: the exact public binary hash and import symbols were independently
checked in Actions run 40560226616. Native synthetic pthread lifecycles and
return codes were measured in runs 40562499242, 40563030141, 40563566337,
40563904429, 40564245570, 40564860228, 40565651582 and 40566543771. Deterministic Go tests
compare the state transitions, 40 QoS argument pairs, error paths and busy
returns. Public native fixtures contain no account data.

The capture ZIP is the original StoreAgent state, not a copy of the subsequent
account authorization state saved in the iPhone sandbox. The SINF candidate in
that experiment was rejected. Neither reaching a new import nor a successful
pthread operation proves that this state can produce a mobile license. All
original-SINF, source/frame and IPA validation gates remain mandatory.
