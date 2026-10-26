package main

// Optional alternate FairPlay decoder, distinct from the Mac enrollment-model
// runtime. A 16-byte credential must have been legitimately recovered for the
// Apple Account that owns this IPA; legacy SIDB is NOT such a credential.
import (
    "os"

    "github.com/majd/ipatool/v2/internal/fairplay"
    credentialengine "github.com/majd/ipatool/v2/wafflebridge/credentialengine"
)

// No encryption material or filesystem paths may appear in JSON diagnostics.
func decryptIPACredential(source, destination, credentialPath string) (any, error) {
    if source == "" || destination == "" || credentialPath == "" || destination == source || destination == credentialPath {
        return nil, recoveryFailure("credential-input", "credential_paths_invalid", nil)
    }
    info, err := os.Lstat(credentialPath)
    if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 || info.Size() != credentialengine.CredentialSize {
        return nil, recoveryFailure("credential-input", "credential_private_16byte_file_required", nil)
    }
    credential, err := os.ReadFile(credentialPath)
    if err != nil || len(credential) != credentialengine.CredentialSize {
        return nil, recoveryFailure("credential-input", "credential_unreadable", nil)
    }
    defer clear(credential)
    geometry, err := fairplay.InspectPortableIPA(source)
    if err != nil { return nil, recoveryFailure("ipa-inspection", "ipa_inspection_failed", nil) }
    expected, _ := geometry["encrypted_images"].(int)
    if expected < 1 || expected > 256 {
        return nil, recoveryFailure("ipa-inspection", "no_supported_encrypted_images", nil)
    }
    out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
    if err != nil {
        return nil, recoveryFailure("output", "private_output_unavailable", nil)
    }
    completed := false
    defer func() {
        if !completed { _ = os.Remove(destination) }
    }()
    result, err := credentialengine.Decrypt(source, out, credential, false)
    if err != nil {
        out.Close()
        return nil, recoveryFailure("credential-decrypt", "credential_or_sinf_rejected", nil)
    }
    if err := out.Sync(); err != nil {
        out.Close()
        return nil, recoveryFailure("output", "output_write_failure", nil)
    }
    if err := out.Close(); err != nil { return nil, recoveryFailure("output", "output_write_failure", nil) }
    if result.Binaries != expected || result.Pages < 1 {
        return nil, recoveryFailure("validation", "encrypted_image_coverage_incomplete", nil)
    }
    verified, err := credentialengine.VerifyPatchedOutput(destination, result.PatchedPaths)
    if err != nil || verified != expected {
        return nil, recoveryFailure("validation", "processed_image_verification_failed", nil)
    }
    remaining := 0
    if stat, err := os.Stat(destination); err != nil || stat.Size() <= 0 {
        return nil, recoveryFailure("validation", "output_empty", nil)
    }
    completed = true
    return map[string]any{
        "engine": "credential16",
        "completed": true,
        "credential_sinf_matched": true,
        "all_encrypted_images_processed": true,
        "images": result.Binaries,
        "pages": result.Pages,
        "remaining_encrypted_images": remaining,
        "plaintext_authenticated": false,
    }, nil
}
