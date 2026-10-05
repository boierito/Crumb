# User report: build 23004

Saved account restored. Login attempts encountered HTML 404/403/503, empty 204,
and an HTML 301 rejected by endpoint/redirect policy. A later attempt received
302 plist, followed the pod, retried HTML 404, then received 200 plist and saved
DSID/token/storefront/pod without a fresh 2FA challenge.

Three versions attempts each generated kbsync, received ent/download HTTP 401
with empty body, then volumeStore pod HTTP 200 plist failureType 2042. Probe v4
called it session-expired-confirmed-by-Apple-response; reference ipatool calls
2042 SignInRequired. This does not prove token expiry or Apple kbsync acceptance.
Versions, purchase, download and export remain unvalidated.

Dev.5 corrects a serial-number byte omission and validates cookie application
in local transport fixtures. These fixes need another physical-device test.
