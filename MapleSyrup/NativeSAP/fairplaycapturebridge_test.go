package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/state"
)

func TestLegacyCandidatePreservesWholeStateAndIdentity(t *testing.T) {
	hardware := []byte{1, 2, 3, 4, 5, 6}
	sidb := bytes.Repeat([]byte{17}, 883)
	context, _ := json.Marshal(map[string]any{"schema": "crumb-fairplay-context-v2", "producer": "retained-storeagent", "hardware_id": hardware})
	consumer, _ := json.Marshal(map[string]any{"consumer": "storeagent-10.9-import-keybag", "emulator_completed": true, "import_succeeded": true, "native_status": 0, "state_files": map[string][]byte{"SC Info.sidb": sidb, "SC Info.sidb.lock": {}}})
	input, err := decodeLegacyEnrollment(context, consumer)
	if err != nil {
		t.Fatal(err)
	}
	defer input.clear()
	if !bytes.Equal(input.hardware, hardware) || !bytes.Equal(input.files["SC Info.sidb"], sidb) {
		t.Fatal("identity or opaque SIDB changed")
	}
	if _, err := decodeLegacyEnrollment([]byte(`{"schema":"crumb-fairplay-context-v2","producer":"retained-storeagent","hardware_id":"AQID"}`), consumer); err == nil {
		t.Fatal("accepted wrong hardware length")
	}
	if _, err := readLegacyEnrollmentInput(captureV3Fixture(t, map[string]any{"native_status": 0, "import_succeeded": true})); err == nil {
		t.Fatal("accepted capture with no retained hardware identity")
	}
}

func TestRecoveredFrameGateChecksLaterSlotsAndRejectsMismatch(t *testing.T) {
	source := make([]byte, 224)
	for i := range source {
		source[i] = byte(i)
	}
	transformed, err := state.Evaluate(source, nil)
	if err != nil {
		t.Fatal(err)
	}
	selected, err := frame.Evaluate(transformed[8*16:9*16], nil)
	if err != nil {
		t.Fatal(err)
	}
	value := rawsource.Result{Source: source, Frame: selected}
	if err := verifyRecoveredEnrollment(value); err != nil {
		t.Fatal(err)
	}
	value.Frame = []byte{1, 2, 3}
	if err := verifyRecoveredEnrollment(value); err == nil {
		t.Fatal("accepted nonmatching frame")
	}
}

func retainedLockFixture(t *testing.T, lock string) string {
	t.Helper()
	source := captureV3Fixture(t, map[string]any{"native_status": 0, "import_succeeded": true, "state_files": map[string][]byte{"SC Info.sidb": bytes.Repeat([]byte{17}, 859), lock: nil}})
	archive, err := zip.OpenReader(source)
	if err != nil {
		t.Fatal(err)
	}
	blobs := map[string][]byte{}
	for _, file := range archive.File {
		r, e := file.Open()
		if e != nil {
			t.Fatal(e)
		}
		blobs[file.Name], e = io.ReadAll(r)
		r.Close()
		if e != nil {
			t.Fatal(e)
		}
	}
	archive.Close()
	blobs["run/context.json"], err = json.Marshal(map[string]string{"schema": "crumb-fairplay-context-v2", "producer": "retained-storeagent", "hardware_id": "AQIDBAUG", "dsid": "1", "keybag_consumer": "storeagent-10.9-import-keybag"})
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Schema string              `json:"schema"`
		Files  []map[string]string `json:"files"`
	}
	if err := json.Unmarshal(blobs["run/manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	for _, entry := range manifest.Files {
		data := blobs["run/"+entry["name"]]
		sum := sha256.Sum256(data)
		entry["bytes"] = strconv.Itoa(len(data))
		entry["sha256"] = hex.EncodeToString(sum[:])
	}
	blobs["run/manifest.json"], err = json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "retained.zip")
	file, err := os.Create(destination)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	for name, data := range blobs {
		entry, e := writer.Create(name)
		if e != nil {
			t.Fatal(e)
		}
		if _, e := entry.Write(data); e != nil {
			t.Fatal(e)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return destination
}

func TestRetainedCaptureAllowsFlatLockNames(t *testing.T) {
	for _, lock := range []string{"lock", ".lock", "SC Info.sidb.lock"} {
		t.Run(lock, func(t *testing.T) {
			source := retainedLockFixture(t, lock)
			if _, err := analyzeCapture(source); err != nil {
				t.Fatal(err)
			}
			input, err := readLegacyEnrollmentInput(source)
			if err != nil {
				t.Fatal(err)
			}
			defer input.clear()
			if len(input.files["SC Info.sidb"]) != 859 || !bytes.Equal(input.hardware, []byte{1, 2, 3, 4, 5, 6}) {
				t.Fatal("lost original SIDB or identity")
			}
			if data, exists := input.files[lock]; !exists || len(data) != 0 {
				t.Fatal("empty lock not preserved")
			}
			_, err = recoverCaptureEnrollment(source, filepath.Join(t.TempDir(), "missing.ipa"), filepath.Join(t.TempDir(), "output"))
			stage, code := labErrorDetails("recover-capture", err)
			if stage != "recover-capture/ipa-inventory" || code != "ipa_structure_invalid" {
				t.Fatalf("did not reach IPA inventory: %s / %s", stage, code)
			}
		})
	}
}

func TestRecoveryErrorDetailsNeverExposeCause(t *testing.T) {
	private := fmt.Errorf("PRIVATE-HARDWARE-AND-PATH-CANARY")
	err := recoveryFailure("identity", "capture_hardware_length_invalid", private)
	stage, code := labErrorDetails("recover-capture", fmt.Errorf("wrapped: %w", err))
	if stage != "recover-capture/identity" || code != "capture_hardware_length_invalid" {
		t.Fatal("lost safe error details")
	}
	document, _ := json.Marshal(map[string]string{"stage": stage, "code": code})
	if bytes.Contains(document, []byte("PRIVATE")) {
		t.Fatal("private cause leaked")
	}
}

func TestInventoryFailureOnlyExposesAllowedCodes(t *testing.T) {
	for _, test := range []struct{ input, want string }{
		{"ipa_sinf_missing", "ipa_sinf_missing"},
		{"ipa_sinf_ambiguous", "ipa_sinf_ambiguous"},
		{"PRIVATE-UNTRUSTED-CODE", "ipa_inventory_failed"},
	} {
		err := enrollmentInventoryFailure(fmt.Errorf("PRIVATE-PATH: %w", &fairplay.EnrollmentInputError{Code: test.input}))
		stage, code := labErrorDetails("recover-capture", err)
		if stage != "recover-capture/ipa-inventory" || code != test.want {
			t.Fatalf("unexpected safe diagnostic: %s / %s", stage, code)
		}
	}
}

func TestPortableRuntimeFailurePreservesPhaseThroughBridge(t *testing.T) {
	_, rawErr := rawsource.Recover("", nil, nil)
	err := recoveryFailure("portable-runtime", "legacy_portable_runtime_rejected", rawErr)
	stage, code := labErrorDetails("recover-capture", err)
	if stage != "recover-capture/portable-runtime/inputs" || code != "raw_invalid_inputs" {
		t.Fatalf("lost portable phase: %s / %s", stage, code)
	}
}

func TestFlatLockFixStillRejectsTraversal(t *testing.T) {
	context := []byte(`{"schema":"crumb-fairplay-context-v2","producer":"retained-storeagent","hardware_id":"AQIDBAUG"}`)
	for _, name := range []string{"../lock", "/lock", "lock\\escape", ".", ".."} {
		consumer, _ := json.Marshal(map[string]any{"consumer": "storeagent-10.9-import-keybag", "emulator_completed": true, "import_succeeded": true, "native_status": 0, "state_files": map[string][]byte{"SC Info.sidb": {17}, name: nil}})
		_, err := decodeLegacyEnrollment(context, consumer)
		stage, code := labErrorDetails("recover-capture", err)
		if stage != "recover-capture/state-files" || code != "state_filename_invalid" {
			t.Fatalf("unsafe path accepted or misclassified: %s / %s", stage, code)
		}
	}
}

func TestCapturedSINFPairsOnlyOneExactEncryptedDomain(t *testing.T) {
    captures := map[string]bool{"own-license-sha256": true}
    tests := []struct {
        name string
        domains []fairplay.PortableEnrollmentInput
        want bool
    }{
        {"matching-original", []fairplay.PortableEnrollmentInput{{SINFSHA256:"own-license-sha256"}}, true},
        {"other-personalization", []fairplay.PortableEnrollmentInput{{SINFSHA256:"other-license-sha256"}}, false},
        {"empty-digest", []fairplay.PortableEnrollmentInput{{SINFSHA256:""}}, false},
        {"no-encrypted-images", nil, false},
        {"ambiguous-multiple-images", []fairplay.PortableEnrollmentInput{
            {SINFSHA256:"own-license-sha256"}, {SINFSHA256:"own-license-sha256"},
        }, false},
    }
    for _, tt := range tests {
        t.Run(tt.name,func(t *testing.T){
            if got:=capturedSINFMatchesDomains(captures,tt.domains); got!=tt.want {
                t.Fatalf("incorrect exact SINF binding result for %s: got %v want %v",tt.name,got,tt.want)
            }
        })
    }
}

func TestBindingMismatchUsesFixedErrorWithoutPrivateMaterial(t *testing.T) {
    _, err := verifyCaptureIPABinding("", "PRIVATE-ORIGINAL-IPA-PATH")
    if err == nil { t.Fatal("missing private capture erroneously accepted") }
    stage, code := labErrorDetails("verify-capture-ipa",err)
    if stage != "verify-capture-ipa/offline-binding/capture-input" || code != "capture_input_invalid" {
        t.Fatalf("unexpected binding rejection: %s / %s",stage,code)
    }
    encoded,_:=json.Marshal(map[string]string{"stage":stage,"code":code})
    if bytes.Contains(encoded,[]byte("PRIVATE")) {t.Fatal("private file path leaked")}
}
