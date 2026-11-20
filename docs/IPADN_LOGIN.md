# iPadn login adaptation

Reference: mengangs/ipadn at 1a74050ac45455707ec3cadd93d310e8df3e427d.
Branch: feature/ipadn-login.

Crumb's existing login controller, Keychain session storage, SAP signer, validated
redirects and short-lived 2FA preparation now use these iPadn request semantics:

- XML plist content type, signing exactly the transmitted bytes.
- attempt=1 for password login, attempt=2 when a verification code is supplied.
- A decoded credential failure ends the operation rather than replaying the
  password with a synthetic second attempt.
- Storefront headers for Chinese and Indian phone-number Apple Accounts.
- DSID resolution from the plist, then X-Dsid, then the current response's
  mz_at0-/mz_at_ssl- cookies. Token and storefront remain mandatory.
- Non-plist HTTP 200 responses enter bounded transport recovery, preserving
  the signed payload and verification attempt.

This adapts protocol behavior; it does not copy iPadn source files. Crumb retains
its existing retry budget, isolated cookie transport and native SAP runtime.

## Remaining compatibility boundary

iPadn also supports auth.itunes.apple.com native endpoints and AMD device-proof
provisioning/reset. Those require additional native guest operations and are
not implemented by this adaptation. Crumb continues to reject those endpoints;
do not broaden the credential allowlist without implementing their proof flow.
This branch therefore does not provide full parity with every iPadn login path.

## Validation

Authentication fixtures cover request bytes, 2FA attempts, terminal credential
errors, phone storefront routing, response DSID fallback and HTTP 200 HTML
recovery. Run:

    swift test --package-path MapleSyrup/SAPKit

The current Linux workspace has no Swift toolchain, so these tests have not
been executed here. Xcode compilation and real-device Apple acceptance remain
required, including password login, correct/incorrect 2FA, cancellation and
session restoration. No claim of live Apple acceptance is made.
