# FairPlay offline evidence reuse — 2026-10-27

## Goal and delivered code

PR #12 implements a **session-scoped import** of previously exported private
FairPlay evidence and its **matching original encrypted IPA**. It removes the
need for repeated downloads and StoreAgent keybag exchanges during research
iterations, while leaving all cryptographic admission checks unchanged.

- `WaffleStore/App/FairPlayLabView.swift`, `FairPlayLab.analyzeCapture`:
  validates an imported `crumb-fairplay-capture-v3` ZIP with the native
  `analyze-capture` operation (manifest hashes, one transaction, verified
  integrity), copies it into a private temporary file, and remembers it as the
  active session's source. The previous UI only *displayed* ZIP analysis.
- `FairPlayLab.useExistingOriginalIPA`: accepts a user-chosen regular
  `.ipa` file from Files, enforces a 512 MiB upper bound, copies it into
  the protected lab directory, invokes the native `inspect` operation and,
  when a ZIP is already loaded, checks the ZIP/IPA association with
  `research-bridge`.
- `FairPlayLab.verifyImportedBinding`: invokes the **unchanged** authentic
  `research-bridge` operation, whose Go implementation checks the SHA-256
  of SINF from the retained transaction against the exact SINF in the
  encrypted IPA. Only a correctly bound audit result is accepted. A
  `sinf_rejected` audit status permits **research**, never decrypt-ready
  status; no model is created or exported.
- `analyzeCurrentCaptureFiles`: routes the existing audit/recovery/causality
  operations through the imported ZIP if present, otherwise uses the
  original in-session native capture/export path.
- The original `onDisappear` cleanup is unchanged. Imported evidence is
  ephemeral, excluded from backups with the temporary laboratory root,
  protected with `FileProtectionType.complete` and permissions 0600; the
  user continues to control the original files in Files.

## Usage

1. Before closing an *older* laboratory session, export **Private
   FairPlay capture ZIP** and keep its matching *original encrypted IPA*.
   The ZIP includes account-derived secrets: handle it as sensitive.
2. Install a newer signed Crumb IPA and open `FairPlay Laboratory`.
3. Select `Import an older capture ZIP for offline research`.
4. Select `Import existing encrypted IPA from Files (no download)`.
   Their SINF association is checked locally. Either import order is
   supported.
5. Run `Audit StoreAgent → fairplayd bridge (no download)`, or an
   appropriate bounded local recovery attempt. Do not assume model validity
   from a successful imported capture.

The original IPA must be **encrypted** and contain the same SINF license as
the capture. A decrypted, repersonalized, mismatched, or invalid IPA is
rejected by one of the native checks. If only the ZIP is available, there is
no supported automatic re-download of its original application.

## What this does *not* solve

The on-device legacy StoreAgent 10.9 guest writes an SIDB with header
`0x0007000b` but a geometry unlike UnfairPlay's independently verified
modern macOS enrollments; it did not produce a nonempty SIDD in the owner's
observed samples. Successful keybag import does **not** yield a usable
modern fairplayd model; a reproducible model from SIDB still fails strict
SINF private-atom decryption. Upstream UnfairPlay expects a matching
macOS-native enrollment with account and hardware state, whereas the Crumb
transaction is for the iOS StoreAgent route.

No jailbreak, JIT, Mac, external server, imported third-party account
credentials, fabricated SIDD, or optimistic decryption flag has been
introduced. Cryptographically valid SINF decoding and IPA verification remain
mandatory. Nothing in this release establishes or claims a fully decrypted
FairPlay IPA.

## Validation

The preexisting Go `researchCaptureBridge` binding is exercised by the
host test suite and rejects transaction mismatches; the iOS-specific
import/file-picker changes are compiled in the Debug and Release IPA jobs.
This document must not label iOS compilation or device-side workflows
successful until those jobs conclude.

Implementing an *authentic* iOS-only enrollment bridge would require
independently verified software semantics for every required identity,
keybag, state and SINF transform, plus a private real-device success
verified against complete Mach-O plaintext and CodeDirectory gates.
