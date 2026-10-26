package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/state"
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

// Exercise real evaluators and filesystem gates, rather than mocking the stage
// labels: a single recovery must identify where accepted CPU output is rejected.
func syntheticRecoveredResult(t *testing.T) rawsource.Result {
	t.Helper()
	source := make([]byte, 224)
	for i := range source {
		source[i] = byte(i)
	}
	transformed, err := state.Evaluate(source, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(transformed)
	selected, err := frame.Evaluate(transformed[8*16:9*16], nil)
	if err != nil {
		t.Fatal(err)
	}
	return rawsource.Result{Source: source, Frame: selected, Pack: []byte("PRIVATE-PACK-CANARY"),
		Runtime: rawsource.RecoveryDiagnostic{Phase: "completed", ExecutionBackend: "generic-profile", GenericInstructions: 123456, Lookups: map[string]rawsource.LookupCount{"state": {Attempts: 1}}}}
}

func TestRecoveredValidationIdentifiesGateAndPreservesRuntime(t *testing.T) {
	for _, mode := range []string{"state", "frame", "directory", "ipa"} {
		t.Run(mode, func(t *testing.T) {
			value := syntheticRecoveredResult(t)
			destination := filepath.Join(t.TempDir(), "output")
			wantCode := ""
			switch mode {
			case "state":
				value.Source = []byte("PRIVATE-BAD-STATE")
				wantCode = "recovered_state_evaluation_failed"
			case "frame":
				value.Frame = []byte("PRIVATE-BAD-FRAME")
				wantCode = "recovered_source_frame_mismatch"
			case "directory":
				if err := os.Mkdir(destination, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(destination, "sentinel"), []byte("keep"), 0600); err != nil {
					t.Fatal(err)
				}
				wantCode = "recovery_output_directory_failed"
			case "ipa":
				wantCode = "recovered_ipa_inventory_failed"
			}
			ipa := filepath.Join(t.TempDir(), "PRIVATE-INVALID.ipa")
			if err := os.WriteFile(ipa, []byte("PRIVATE-BAD-ZIP"), 0600); err != nil {
				t.Fatal(err)
			}
			result, err := finishRecoveredEnrollments([]rawsource.Result{value}, ipa, destination)
			if err == nil || result != nil {
				t.Fatal("unvalidated result accepted")
			}
			fields := labFailureFields("recover-capture", fmt.Errorf("PRIVATE-WRAPPER: %w", recoveryFailure("recovered-validation", "recovered_output_validation_failed", err)))
			if fields["code"] != wantCode {
				t.Fatalf("got %v, want %s", fields, wantCode)
			}
			report := fields["validation"].(*recoveredValidationDiagnostic)
			runtime := fields["runtime"].(*rawsource.RecoveryDiagnostic)
			if runtime.Phase != "completed" || runtime.GenericInstructions != 123456 || runtime.Lookups["state"].Attempts != 1 {
				t.Fatal("lost runtime after successful execution")
			}
			if mode != "state" && (!report.StateEvaluationPassed || report.FrameCandidatesChecked != 14) {
				t.Fatalf("incomplete cross-check: %+v", report)
			}
			if (mode == "directory" || mode == "ipa") && (!report.SourceFrameMatched || report.DomainsCrossChecked != 1) {
				t.Fatal("lost successful prior gates")
			}
			if report.TargetIPAValidated {
				t.Fatal("failed validation marked as successful")
			}
			document, err := json.Marshal(fields)
			if err != nil || bytes.Contains(document, []byte("PRIVATE")) {
				t.Fatal("private input or error escaped")
			}
			for _, data := range [][]byte{value.Source, value.Frame, value.Pack} {
				if !bytes.Equal(data, make([]byte, len(data))) {
					t.Fatal("failed recovery retained private output")
				}
			}
			if mode == "directory" {
				if _, err := os.Stat(filepath.Join(destination, "sentinel")); err != nil {
					t.Fatal("deleted preexisting destination")
				}
			} else if _, err := os.Stat(destination); !os.IsNotExist(err) {
				t.Fatal("failed recovery left an output directory")
			}
		})
	}
}
