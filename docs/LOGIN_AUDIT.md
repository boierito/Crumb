# Login audit — Crumb 23021

Compared on 2026-10-06 against `majd/ipatool` main at
`cde7d00355e152714377b953ec57438626d3cb5a`. Automated tests cannot establish
Apple's live acceptance or a measured device login speedup.

## Upstream changes missed by the previous comparison

Between the native source pin `3411d57` and current main, authentication changed:

- [`f9aa653`](https://github.com/majd/ipatool/commit/f9aa6536cd2db16ab93f663d8731e4300921c132),
  October 5: validate the initial Bag URL, then use
  `/WebObjects/MZFinance.woa/wa/authenticate/`, preserving host, port and query.
  Its comment explicitly associates the bare path with unusable responses,
  including 301 without Location. Crumb 23020 missed this correction.
- `93b527a`, October 6: country-derived `X-Apple-Store-Front` for phone-number
  accounts. This does not change the reported email login. Crumb does not claim
  the new phone-number support.

The native SAP and HTTP client did not change between these revisions. Keep
the native pin/TCI ABI/guest operations and user agent unchanged. Updating the
Go dependency alone would not fix the Bag path in Crumb's Swift implementation.

## Evidence from the supplied reports

In the 23009 reports, prepared SAP takes roughly 1.5–2.4 seconds, signing
50–90 ms. Most failed responses are empty/HTML 204/301/404/503, without a parsed
Apple credential failure. Eventually a valid 302 routes to a pod and a 200
contains a session. These responses do not prove a bad password or signer.

Crumb's automatic fallback pauses are 2, 4, 8, 15, 15… seconds:

| Trial | Pauses before pod | Pauses at pod before success |
|---|---:|---:|
| 7: five unusable initial replies | 44 seconds | 0 seconds |
| 8: six unusable initial replies | 59 seconds | 14 seconds |

These are sums of configured sleeps, not measured end-to-end times. Trial 8
therefore includes 73 seconds of waiting, before adding SAP/signing/transfers.
The old recovery policy amplifies routing failures. Upstream's path correction
matches the missing-Location symptom, but cannot explain every HTTP failure
without a real-device comparison. Signature length is not an acceptance test.

The old URL validator also compared `URLComponents.percentEncodedPath` to
`URL.path`. Foundation's URL.path can remove a trailing directory slash: the
Linux regression run reproduced rejection of a legitimate canonical endpoint.
This check could reject a valid pod/2FA URL even though both textual forms were
nominally allowed. Validation now uses the exact encoded HTTP path, accepting
only the two explicit allowed literals; encoded aliases remain rejected.

## Current request/state comparison

| Stage | Current ipatool | Crumb 23021 |
|---|---|---|
| Identity | Stable MAC → hardware ID/GUID | Stable six-byte Keychain identity for Bag, guest and payload |
| Bag | Discover, validate, canonicalize initial auth path | Same; supports canonical and pod URLs advertised by Bag |
| SAP | Cached assets → guest init → certificate → setup exchange | Same state flow through static TCI, sandbox cache, no unsigned fallback |
| Sign | Sign serialized plist for every send | Exact transmitted bytes signed afresh, including normalized 2FA code |
| Payload | appleId/attempt/guid/password/rmp/why | Same fields and email account headers; Swift XML encoding is signed byte-for-byte |
| Connections | Authentication transport, keepalives disabled | Isolated ephemeral login URLSessions, shared challenge jar; separate Store transport |
| Pod | Validate and replay same POST/body/attempt | Same; redirected path/query preserved exactly |
| Credential retry | Logical attempt 2 for initial -5000 | Same; parsed Apple errors never reclassified as temporary HTML |
| 2FA | Password+code, new signature | Same; RAM-only challenge, reuse live signer/pod within five minutes |
| HTTP recovery | Three attempts per request, 10/20s, Retry-After override | Existing automatic recovery retained, eleven retries shared across one login's routes/logical attempts, existing 120s scheduling deadline |
| Missing Location | Error | Existing bounded iOS retry at the current validated URL; no invented destination |
| Persistence | Keychain account, including password upstream | Account/token/cookies in Keychain, no password or code persistence |

## Changes and boundaries

`AuthenticationEndpoint.initial` validates before changing exactly the initial
allowed path. Foreign hosts, HTTP, user info, fragments, unsupported ports,
encoded aliases, double slashes and extra paths are rejected before repair.
Encoded routing query bytes (`%2F`, `+`) and advertised host/port are retained.
`SAPConfiguration.parse` applies this rule; `PreparedAppleLogin` consumes the
result. SAP-only native discovery still cannot bypass credential validation.
Apple's redirect destinations are validated and replayed as provided, not
rewritten to a guessed pod or changed into GET.

`AuthenticationRecoveryBudget` is created once per login/verification. Useful
routing replies do not consume retries, but changing pods or incrementing the
logical attempt cannot refill the eleven automatic repeat allowances. HTTP
and transient network failures share that budget. A final response may still
succeed when all eleven repeats have been used. Four-hop redirect and per-route
attempt limits remain. A new submitted 2FA code is a separate user action with
its own budget. Non-automatic clients retain the three-attempt reference policy.

Backoff and Retry-After handling are preserved for the device comparison. Sending
more frequently could worsen rate limits and obscure the path fix's effect.
There is no wait on a valid redirect, normal 2FA challenge or success. Signing
time counts against the existing scheduling deadline; fractional request
timeouts are no longer rounded up to one second. URLSession still has a separate
30-second resource cap; synchronous native asset/guest work has independent
limits, so 120 seconds is not a hard limit for the entire UI flow.

No new Release logs/debug UI, credential persistence, pooled login connections
or guest changes. Saved sessions need not be removed to update. First-time asset
downloads can still take longer than warm SAP preparation; measure separately.

## Verification

Nine new regressions cover canonical initial URLs and exact routing queries,
invalid URLs before normalization, the initial signed POST/exact pod redirect,
one prepared SAP handshake across challenge/verification, shared HTTP/network/
logical-attempt budgets, last-permitted-retry success and fractional timeout.
Existing tests cover unsafe redirects, cookies, Retry-After, codes, persistence,
cancel, Keychain and Store/download behavior. Guest signatures and Apple replies
in these tests are fixtures, not live authentication.

See the 23021 procedure in [TESTING.md](../TESTING.md) to compare normally signed
Release builds without debugger/JIT, separating cold preparation, time to 2FA
and time to the verified session. UI-level new-code delivery and live login
speed remain pending physical-device verification.
