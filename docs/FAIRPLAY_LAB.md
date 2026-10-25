# FairPlay laboratory — device test guide

This build tests portable components and collects inputs. It does **not** yet
create the missing enrollment state from an Apple login or a keybag.

1. Sign the Release IPA and install on iOS 16.4 or newer. Sign into Crumb.
2. Open Settings → FairPlay Laboratory → Run self-tests. The state, frame and
   SIDB vectors should all pass. These use public synthetic fixtures only.
3. Enter an App Store ID/link or bundle ID for a small free app you can download.
   Tap Download and capture transaction. Wait for completion. Log entries show
   matched kbsync/keybag lengths and SINF count; absent keybags and fallback
   routes are recorded. Each successful transaction has its own capture folder.
4. Tap Inspect IPA. Record SINF/SUPF counts and Mach-O candidates. Alternatively
   open Downloads → app options → FairPlay Laboratory → Use selected downloaded
   IPA. Existing downloads contain no newly captured Store transaction.
5. If you already have a matching private enrollment model, import it and test
   decryption. Without one, stop here: the missing raw-state → 192/224-byte state
   input stage is not implemented. Capturing a keybag is not recovery.
6. Prepare exports. Share the diagnostic for ordinary troubleshooting. Share the
   PRIVATE capture only with a trusted researcher: it contains account-derived
   material. No password, 2FA code, auth token, cookies or signed download URL is
   included. The IPA and imported model are not included in the capture ZIP.
7. Export a successful decrypted IPA, sign it, and record installation and launch
   results. Source CodeDirectory hashes check encrypted input; structural ARM64
   checks and cryptid=0 alone do not prove plaintext correctness.
8. Export before closing. Closing or Erase this test removes this session's
   temporary files; original downloads, account and other sessions are preserved.
   Repeat with a second app from the same account. Force-quitting may leave
   temporary files until iOS clears its temporary storage.

Portable source pin: DrWhax/ipatool ef98227e061838320424744d21a513baba7ade28.
The native build isolates state/frame/SIDB packages, retains their public vectors,
and patches the IPA page offset and first executable-page audit. Model resources
are pinned by the upstream fixture tests. SAP remains pinned independently.

Host tests do not establish real enrollment compatibility or iOS success. Device
logs intentionally do not expose raw native exceptions or secret evaluator output.
