# Clean-login test branch — Crumb Test 23023

Branch: `feature/2fa-clean-login-test`. Base main:
`01ad52dd53b725ee71dd7c1bfbacfd36cd14a196`. Main is not modified or merged into.
This is a normally resignable, no-JIT test IPA, not a production release.

## Isolation

- Display name: **Crumb Test**; bundle ID: `com.certlium.crumb.authtest`.
- Keychain service: `com.certlium.crumb.authtest.sap`, for identity/session/kbsync.
  It differs from production even if the signer shares Keychain access groups.
- Separate sandbox and `crumb-authtest` URL scheme. Keep this bundle ID when
  resigning; installing over production would defeat filesystem isolation.
- No production Keychain fallback, broad deletion, access-group manipulation,
  Apple trust revocation or changes to main. The legacy-key cleanup tag is also
  test-specific. Downloads/favourites are not removed by the reset.

## Install and test with your existing Apple ID

1. Download **CrumbAuthTest-Release** from this branch's successful Build IPA run,
   extract its IPA and resign/install with ksign/certificate or another regular
   sideloader. It appears beside production as **Crumb Test**. No JIT/debugger.
2. Open **Account → Reset test sign-in** and confirm. Wait for **New identity
   verified in Keychain**. If reset verification fails, do not assume a clean
   client: the login action remains disabled until a successful reset.
3. Enter the same Apple ID/password and sign in once. Do not repeatedly reset
   or submit during network work. Stop/wait if Apple rate limits the account.
4. If Apple requests 2FA, enter the actual six-digit trusted-device code and
   verify. A temporary transport failure keeps it for Retry verification;
   explicit Apple rejection requires another valid code. Request new code is
   an explicit password-only action with the existing cooldown; Apple controls
   delivery. Cancellation/account change cleanup is unchanged.
5. If login succeeds without 2FA, open **Test report**. The report distinguishes
   that outcome from verification and restoration. A verified local reset does
   **not** remove Apple's server-side trust or guarantee a fresh challenge.
6. Share/copy the report text, not screenshots of credentials. Report is held
   in bounded RAM only; save it before terminating the app. Session restoration
   on relaunch is explicitly labelled, not reported as a fresh-login test.

The reset makes no Apple network request. It requires no login/Store/download
work in progress, closes the guest and ephemeral transports before rotation,
clears the account/kbsync, deletes only the test identity and generates another.
It then verifies changed identity, read-back stability and absence of the two
account records. Reset errors never report success. No automatic reset per login.

Downloaded guest assets remain cached: they are loader input, not a saved Apple
authentication session. A new login still fetches Bag and initializes a new SAP
guest with the replacement identity. Keeping assets avoids confusing first-asset
download time with authentication/2FA latency.

## Report and limits

Report contains elapsed times/stages, numeric HTTP status, body category,
numeric Apple failure category, cookie **counts**, reset evidence and outcome.
It excludes password/code/email/DSID/token/GUID, raw bodies, signature contents,
cookie/header values, arbitrary error messages and URL/query strings. Only
complete allowlisted diagnostic formats are accepted; unknown formats are dropped.
Native runtime and authentication payload/UA/redirect/retry policy are unchanged.
The SDK still signs every request afresh and validates destinations before replay.

Tests cover report filtering/bounds/timing/outcomes, namespace rejection and
actual macOS Keychain rotation/deletion while retaining unrelated items. Those
fixtures are not live Apple acceptance and cannot demonstrate iOS Keychain
behavior, a forced 2FA challenge or measured physical-device login speed.

## Physical verification matrix

| Check | Automated | Physical iPhone iOS 27.0.1 |
|---|---|---|
| Test namespace differs from production | Source/fixture/IPA validation | Pending; keep production signed in and confirm it stays signed in |
| Identity changes and survives read-back | Real macOS Keychain fixture; runtime evidence | Pending; expect identity-changed=true and identity-persisted=true |
| Account/kbsync removed, guest/cookies closed | Keychain tests; controller guards | Pending; expect session-removed=true and kbsync-removed=true |
| Actual Apple 2FA challenge | Fixture flow | Apple decides; pending |
| Verification/resend/retry | Existing protocol/form fixtures | Pending |
| Normally signed IPA without JIT/debugger | macOS arm64 build and host TCI checks | Pending |
| Main ref and checkout unchanged | Git/GitHub SHA comparison | No production reset performed by this branch |

Do not merge the reset/namespace/report changes into main as a login fix. They
are temporary test instrumentation and cannot force server-side authentication.
