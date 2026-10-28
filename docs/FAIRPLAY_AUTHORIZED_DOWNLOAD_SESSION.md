# Retaining the authorized producer through a download

The previous laboratory download constructor always created an empty virtual
filesystem, even when account authorization had saved a matching SIDB/SIDD.
The separate recovery action could mount those saved bytes, but that did not
make them the producer of the download's KBSync or keybag consumption.

The retained laboratory generator now passes the existing account-private
state directory to `WaffleKBSyncOpenWithState`. Before generating mode-1 KBSync,
the bridge validates the saved state against the requested DSID and hardware,
restores it before native global initialization and checks the native diversity
query. It keeps that same producer until the bound download keybag is consumed.
A missing saved file follows the existing cold producer flow. Invalid, foreign
or still-unprovisioned saved state is rejected rather than silently ignored.

Only a completed, status-zero native import may persist its resulting state.
Persistence compares the original file digests under the shared bridge mutex
before replacing state atomically. A different account/identity, a replaced or
removed saved state, or state created during a cold download cannot be
silently overwritten. Persistence failure is reported separately; the private
transaction capture still receives its actual consumer output. State storage
and native import success never assert SINF or plaintext validity.

## Local evidence

The original private capture and its exported post-authorization account state
were reused without Apple authentication, new account secrets or a new IPA.

- Re-importing the original captured keybag into a restored producer completed
  with native status -42052. The account query before import did not require
  diversity. Although the private native filesystem changed, this rejected
  import was not persisted or sent to recovery. This does not identify the
  semantic cause of -42052 or prove a general session requirement by itself.
- The production `OpenRetainedAuthorizedKBSync` factory restored that same saved
  state, generated 196 bytes of mode-1 download KBSync, retained its producer
  and subsequently reported diversity not required. This verifies local
  preparation, not acceptance of a future Apple response or SINF validation.
- The prior download-constructor disconnection is established directly by its
  empty-filesystem initialization, independently of the rejected replay.
- Deterministic storage tests reject stale updates and wrong account/identity
  without changing the current state; they also cover a state file appearing
  or disappearing while a download is pending.

The optional producer test can repeat local mode-1 generation with
`CRUMB_PRIVATE_CAPTURE`, `CRUMB_PRIVATE_AUTHORIZED_STATE`,
`CRUMB_PRIVATE_DOWNLOAD_SYNC=1` and a private `CRUMB_PRODUCER_REPORT` destination
through `scripts/test-sap-host.sh`. No private inputs belong in public CI.

The remaining licensing proof is a genuine response associated with this
retained producer yielding material that validates the original mobile SINF.
The saved-state replay rejection remains unresolved; no new decrypted IPA is
claimed by this integration.
