# FairPlayed system route: controlled reference experiment

Reviewed upstream: https://github.com/morzelowski/fairplayed/tree/bad7df138839b5e0ac8b9f20cc4689140b8436fb

This route delegates page decryption to the iPhone's own kernel/fairplayd. It
does not derive a license key in Crumb or repair the emulated SINF -42057 path.
The upstream's only stated device confirmation is iOS 27.0.1; that claim has
not been independently reproduced here.

## Code findings

`fpcore.c` maps the original file read-only, requests `mremap_encrypted` with
cryptid 2, and reads each protected page using `vm_read_overwrite` on the same
process. Cryptid 2 selects the non-executable model-encryption mapping path;
there is no JIT or privileged process attachment in this attempt. A direct
ARM64 syscall fallback uses syscall 489 when dlsym cannot resolve the symbol.

`embed_ipa.sh` and `main.m:embeddedExport` use a target under the helper's
installed `EmbeddedPayload` directory. The target's top-level `.app` suffix
is removed so it is treated as resources. SINF replication follows the
original SC_Info/Manifest.plist. A file imported into Documents after helper
installation is a different experiment and is not equivalent to this route.

The important limitation is bundle location: upstream reports that fairplayd
cannot read the target's license from a normal data container. Consequently,
this reference helper handles a target embedded before signing/installing.
It is not yet a generic Crumb download-and-decrypt flow. Do not claim that
embedding solves arbitrary on-device downloads or writable bundle access.

## Preparation and privacy

Workflow `fairplayed-system-helper.yml` builds the exact public upstream
commit with Xcode 26.0, unsigned, ARM64. It contains no private target,
license, authorized state, credentials or capture. Successful run:
https://github.com/boierito/Crumb/actions/runs/41041285890

`scripts/prepare-fairplayed-private-ipa.py` adds the existing original target
and its SINF locally after CI. It validates archive paths, duplicates,
regular-file types, size limits, one app, original license presence, helper
identity and output ZIP CRCs. It refuses overwrites and deletes incomplete
output. Original binary bytes are copied unchanged. Manifest-directed SINF
replicas are added where absent. Watch content is explicitly excluded, as
this is an iPhone-only route. Private output must not be uploaded to a public
release or CI. No account capture or password is used in packaging.

The DD preparation verified that all 57 original included target files were
identical after embedding, including the original binary and SINF. This
checks preparation, not native decryption.

## Device procedure and interpretation

1. Sign the private helper IPA with a development identity and install it.
   This is a separate bundle ID (`com.boierito.crumb.systemfairplay`), not an
   update to Crumb. Do not remove or re-sign embedded target resources.
2. Enable Developer Mode. The iPhone's system Apple Account must own the
   target app. Signing into Crumb alone does not establish system licensing.
3. Open **Crumb FP System Test**, then tap **Export** once. Do not select
   Pick IPA or Docs IPAs for the embedded-payload test.
4. After completion, tap **Share** and save `fairplayed.log` plus any output
   IPA. A `-PARTIAL.ipa` is explicitly incomplete.

The original program calls a slice decrypted after successful remapping,
complete page reads, changed bytes, entropy reduction and increased counts
of selected ARM64 instructions. These last checks are heuristics, not
plaintext authentication. A `cryptid=0` output or a `-decrypted` filename is
not sufficient proof. Returned output must still be compared against the
original: complete Mach-O inventory, encrypted-page coverage, exact changed
ranges, unchanged other files and applicable structural/cryptographic
checks. The upstream copy operation tolerates unreadable resources, so
also inspect its skipped-file list before claiming an intact application.

## First returned DD device result, 2026-10-30

The returned device log identifies iPhone17,1 and iOS 27.0.1. Both encrypted
ARM64 executables completed the read-only cryptid-2 route without adding
file signatures: remapping returned zero, and all 177 main-executable pages
plus the extension's single page were read successfully. The output archive
was obtained privately and compared against the retained original.

Independent comparison found:

* All 178 encrypted 16 KiB pages changed. The main range is 2,899,968 bytes;
  the extension range is 16,384 bytes. Both cryptid fields changed from 1 to 0.
* Both Mach-O sizes stayed unchanged and their segments remain within file
  bounds. The main executable differs only in its original encrypted range
  and four-byte cryptid field. The extension also differs within its original
  code-signature blob; no other extension executable bytes changed.
* Capstone decoded all 2,220,760 main-executable __text bytes overlapping the
  encrypted range, and all 9,304 corresponding extension bytes, with zero
  undecodable bytes. The original ciphertext had 1,430,472 and 5,952
  undecodable bytes respectively. The recovered sections have 6,894 and 47
  return instructions. This is strong structural evidence, not plaintext
  authentication.
* ZIP CRC verification passed. Of 57 original files, 45 were exported; the
  12 missing files are exclusively under SC_Info. There are no extra files.
  Besides the two executable files, the extension's CodeResources differs.
  The other 42 exported files are byte-identical to the original. The
  signature-related differences are consistent with re-signing the embedded
  extension before installation, but this comparison alone does not prove
  which signing tool made them.
* Both original SHA-1 and SHA-256 CodeDirectories match all 712 encrypted
  4 KiB code slots across the two original executables. Neither matches any
  of those recovered plaintext slots. For this real sample, these original
  hashes authenticate the ciphertext, not the recovered plaintext.

This result provides concrete evidence of native system page decryption and
IPA export on the user's ordinary iPhone. It does not validate Crumb's
emulated SINF candidates, expose the system's keys, demonstrate independent
plaintext authentication or establish a generic download flow without
pre-install embedding. The exported binaries need appropriate re-signing
before installation; the original encrypted-content signatures do not
validate the modified plaintext.

The remaining integration issue is preparing an accessible target location
for subsequent downloads. Do not ask for more variations of the already
rejected emulator candidates to re-establish this native-route result.
