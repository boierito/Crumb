package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/majd/ipatool/v2/internal/fairplay"
)

func modernCaptureFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Chmod(root, 0700); err != nil {
		t.Fatal(err)
	}
	manifest := fairplay.EnrollmentCaptureManifest{Schema: fairplay.EnrollmentCaptureSchema}
	manifest.App.Executable = "App"
	manifest.App.ExecutablePath = "/Applications/App.app/App"
	manifest.HardwareMAC.Path = "hardware/iomac.bin"
	for name, data := range map[string][]byte{
		"app/executable": {1, 2, 3}, "hardware/iomac.bin": {0, 1, 2, 3, 4, 5},
		"app/SC_Info/App.sinf": {6}, "app/SC_Info/App.supf": {7},
		"app/SC_Info/Other.sinf": {8}, "state/fpsd-sc-info/SC Info.sidb": {9},
	} {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, data, 0600); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(data)
		manifest.Files = append(manifest.Files, fairplay.EnrollmentCaptureFile{Path: name, Bytes: int64(len(data)), SHA256: hex.EncodeToString(sum[:])})
	}
	document, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), document, 0600); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestModernEnrollmentIntegrityAndRouting(t *testing.T) {
	root := modernCaptureFixture(t)
	manifest, files, err := readModernEnrollment(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, expanded := range []bool{false, true} {
		mappings, err := modernRuntimeFiles(manifest, files, expanded)
		if err != nil {
			t.Fatal(err)
		}
		if len(mappings["/App"]) != 3 || len(mappings["/var/db/fpsd/SC Info/SC Info.sidb"]) != 1 || len(mappings["App.sinf"]) != 1 {
			t.Fatal("missing required runtime mapping")
		}
		if _, exists := mappings["Other.sinf"]; exists {
			t.Fatal("expanded retry introduced another installed license")
		}
	}
	if err := os.WriteFile(filepath.Join(root, "app/executable"), []byte{4, 5, 6}, 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readModernEnrollment(root); err == nil {
		t.Fatal("accepted corrupted capture")
	}
}

func TestModernEnrollmentRejectsSymlinksAndLegacySchema(t *testing.T) {
	root := modernCaptureFixture(t)
	full := filepath.Join(root, "app/executable")
	if err := os.Remove(full); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "hardware/iomac.bin"), full); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readModernEnrollment(root); err == nil {
		t.Fatal("accepted symlink")
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte(`{"schema":"crumb-fairplay-capture-v1"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readModernEnrollment(root); err == nil {
		t.Fatal("legacy capture treated as modern state")
	}
}
