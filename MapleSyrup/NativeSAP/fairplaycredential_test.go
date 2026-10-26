package main

import (
    "os"
    "path/filepath"
    "testing"
)

func TestCredentialGateRejectsUnsafeFiles(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "credential.bin")
    result := filepath.Join(dir, "out.ipa")

    if _, err := decryptIPACredential("/missing.ipa", result, path); err == nil {
        t.Fatal("missing credential accepted")
    }
    if err := os.WriteFile(path, make([]byte, 15), 0600); err != nil { t.Fatal(err) }
    if _, err := decryptIPACredential("/missing.ipa", result, path); err == nil {
        t.Fatal("incorrect credential size accepted")
    }
    if err := os.WriteFile(path, make([]byte, 16), 0644); err != nil { t.Fatal(err) }
    if _, err := decryptIPACredential("/missing.ipa", result, path); err == nil {
        t.Fatal("credential with loose permissions accepted")
    }
    if _, err := os.Stat(result); !os.IsNotExist(err) {
        t.Fatal("unsafe input created an output IPA")
    }
}
