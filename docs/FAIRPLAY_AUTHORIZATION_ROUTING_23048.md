# Authorization routing correction — build 23048

## What the 23047 device diagnostic establishes

The owner imported the saved original IPA and its matching capture on iOS
27.0.1. The capture analyzer accepted one transaction, the previous download
keybag import reported native status 0, and the saved state contained 884 bytes.
Authorization then stopped with `serverRejected`. No authorization reply bags
were imported and no SINF validation was reached in that attempt.

23047 mapped non-200 HTTP responses, invalid reply formats, nonzero failureType,
and the mere presence of customerMessage to the same error. The exported JSON
contains no HTTP status, failure code, response body or headers. Therefore it
cannot establish whether Apple refused the token, required ADI, redirected the
request, or whether the client misclassified a response. No account-specific
cause or cryptographic impossibility is inferred from this diagnostic.

## Verified protocol discrepancy

Original public Apple assets, already hash-pinned by ipatool:

- StoreAgent SHA256 `70ce036f9dbcbc04db9511ebd08de0dd3cbc35ccc9d44b089c90170cb5453c59`.
- CommerceKit SHA256 `b84ff12c21987856c0a17b78f1ad82b73195a6dec5f3b208a17d245555a2c8a2`.

CommerceKit `ISStoreURLOperation._runURLOperation` invokes
`_addStandardQueryParametersForURL:` at text+0x1cecc. The query builder at
text+0x1bcf7 asks `sendGUIDWithURL:` (text+0x1bdea), and sets the `guid` query key
(text+0x1bf11). StoreAgent `Bag.sendGUIDWithURL:` at text+0xa79c0 consults the
`guid-urls` dictionary, optional schemes and case-insensitive regex matches
against NSURL.absoluteString (text+0xa7da9); missing/non-dictionary policy falls
back to HTTPS. A read-only current Configurator Bag fetch advertises
`/commerce/machine/authorize` and a `/commerce/` GUID regex. 23047 provided GUID
in the plist but omitted it from this required query. These are distinct inputs.
This is a concrete construction bug, not proof that it caused the owner's
particular rejected response.

StoreAgent text+0x59e73..0x59e76 tests failureType and skips customerMessage
handling when absent. Reading customerMessage occurs at text+0x59ec7 inside
that error branch. 23047's unconditional rejection of any customerMessage,
including an empty field beside genuine reply bags, was stricter than native.

CommerceKit `ISURLOperation._handleRedirectRequest:response:` restores the active
request's body (text+0x1fb32..0x1fb4f) and method (text+0x1fb7d..0x1fbad).
`ISStoreURLOperation.shouldFollowRedirectWithRequest:returningError:` checks the
trusted URL before continuing. 23047's transport intentionally blocks automatic
URLSession redirects, but the authorization client lacked the explicit trusted
routing continuation that the native client supplies.

## Implementation

- Apply the Bag GUID policy to initial and redirected authorization URLs. Keep
  the captured virtual identity, reject conflicting/duplicate GUID parameters,
  and preserve existing encoded routing query bytes. Do not alter the signed
  typed plist or KBSync generation mode.
- Follow at most three explicit 301/302/307/308 routing redirects, only to
  validated HTTPS Apple buy hosts and authorization service paths. Preserve the
  exact POST body and its SAP signature; detect loops. Reject 303, other service
  paths, foreign hosts and identity changes. No retry of HTTP failures or ADI
  challenges is introduced.
- Classify HTTP rejection, Apple error, invalid reply and redirect errors
  separately. Surface existing fixed-field ResponseDiagnostic metadata through
  the dedicated bridge: status, format and bounded/redacted numeric failure
  code. Never log body text, tokens, identities, reply bags, messages or routing
  URLs. HTTP/embedded ADI challenges take precedence over redirect handling.
- Do not infer failure from customerMessage alone. Reply bags must still be
  nonempty, correctly typed, genuinely present and accepted by their distinct
  native import routines. The subsequent account query, every original SINF,
  source/frame and whole-IPA gates remain mandatory.

`verify-store-authorization-profile.py` checks both entire asset hashes, exact
native query/failure branch instructions and GUID-policy literal references.
The normal native producer CI job runs it against the downloaded pinned assets.
Swift fixtures exercise Bag policy, byte-exact routing/signatures, identity
mismatch, malicious redirects, loop/hop limits, response classification and
secret-free diagnostics. These fixtures are protocol checks, not real licenses.

No private capture, original IPA or account token was available in the workspace.
No authenticated authorization POST was sent here. The next owner-device attempt
reuses the same saved IPA/capture. A real SINF/IPA result, ADI challenge completion
and compatibility of legacy enrollment with the modern model remain unproven.
