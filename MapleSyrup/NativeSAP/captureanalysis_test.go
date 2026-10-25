package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func captureFixture(t *testing.T, tamper, duplicate, traversal bool) string {
	t.Helper()
	atom := func(tag string, b []byte) []byte {
		out := make([]byte, 8+len(b))
		binary.BigEndian.PutUint32(out, uint32(len(out)))
		copy(out[4:], tag)
		copy(out[8:], b)
		return out
	}
	fields := append(atom("user", make([]byte, 4)), atom("key ", make([]byte, 4))...)
	fields = append(fields, atom("priv", make([]byte, 440))...)
	sinf := atom("sinf", atom("schi", fields))
	blobs := map[string][]byte{"kbsync.bin": make([]byte, 196), "keybag.bin": make([]byte, 888), "sinf-0.bin": sinf}
	entries := []map[string]string{}
	for name, b := range blobs {
		sum := sha256.Sum256(b)
		entries = append(entries, map[string]string{"name": name, "bytes": strconv.Itoa(len(b)), "sha256": hex.EncodeToString(sum[:])})
	}
	manifest, _ := json.Marshal(map[string]any{"schema": "crumb-fairplay-capture-v1", "files": entries})
	p := filepath.Join(t.TempDir(), "capture.zip")
	f, e := os.Create(p)
	if e != nil {
		t.Fatal(e)
	}
	z := zip.NewWriter(f)
	for name, b := range blobs {
		if tamper && name == "keybag.bin" {
			b = append([]byte{}, b...)
			b[0] = 1
		}
		w, _ := z.Create("run/" + name)
		w.Write(b)
	}
	w, _ := z.Create("run/manifest.json")
	w.Write(manifest)
	if duplicate {
		w, _ = z.Create("run/keybag.bin")
		w.Write([]byte{1})
	}
	if traversal {
		w, _ = z.Create("../bad")
		w.Write([]byte{1})
	}
	if e = z.Close(); e != nil {
		t.Fatal(e)
	}
	f.Close()
	return p
}
func TestAnalyzeCaptureIntegrityAndPrivacy(t *testing.T) {
	result, e := analyzeCapture(captureFixture(t, false, false, false))
	if e != nil {
		t.Fatal(e)
	}
	if result["ready_to_decrypt"] != false {
		t.Fatal("capture treated as model")
	}
	encoded, _ := json.Marshal(result)
	var values map[string]any
	json.Unmarshal(encoded, &values)
	if _, exists := values["selector"]; exists {
		t.Fatal("selector leaked")
	}
	for _, tc := range [][3]bool{{true, false, false}, {false, true, false}, {false, false, true}} {
		if _, e := analyzeCapture(captureFixture(t, tc[0], tc[1], tc[2])); e == nil {
			t.Fatal("unsafe capture accepted")
		}
	}
}
func TestPrivateCaptureExternal(t *testing.T) {
	p := os.Getenv("CRUMB_PRIVATE_CAPTURE")
	if p == "" {
		t.Skip("private fixture not supplied")
	}
	result, e := analyzeCapture(p)
	if e != nil {
		t.Fatal(e)
	}
	encoded, _ := json.Marshal(result)
	t.Log(string(encoded)) // Only the explicit non-secret summary.
}

 
// All v3 fixtures use public synthetic blobs; no account, device, or keybag
// captured from an iPhone is stored in the repository or test output.
func captureV3Fixture(t *testing.T, overrides map[string]any) string {
	t.Helper()
	atom := func(tag string, b []byte) []byte {
		out := make([]byte, 8+len(b))
		binary.BigEndian.PutUint32(out, uint32(len(out)))
		copy(out[4:], tag)
		copy(out[8:], b)
		return out
	}
	fields := append(atom("user", make([]byte, 4)), atom("key ", make([]byte, 4))...)
	fields = append(fields, atom("priv", make([]byte, 440))...)
	report := map[string]any{
		"consumer": "storeagent-10.9-import-keybag",
		"native_status": -42052,
		"emulator_completed": true,
		"import_succeeded": false,
		"portable_model_ready": false,
		"state_files": map[string][]byte{"fixture.sidb": []byte("SYNTHETIC-SECRET-CANARY")},
	}
	for k, v := range overrides {
		report[k] = v
	}
	consumer, err := json.Marshal(report)
	if err != nil { t.Fatal(err) }
	blobs := map[string][]byte{
		"kbsync.bin": make([]byte, 196),
		"keybag.bin": make([]byte, 904),
		"sinf-0.bin": atom("sinf", atom("schi", fields)),
		"context.json": []byte(`{"schema":"crumb-fairplay-context-v2"}`),
		"consumer-state.json": consumer,
	}
	entries := []map[string]string{}
	for name, blob := range blobs {
		checksum := sha256.Sum256(blob)
		entries = append(entries, map[string]string{"name": name, "bytes": strconv.Itoa(len(blob)),
			"sha256": hex.EncodeToString(checksum[:])})
	}
	manifest, err := json.Marshal(map[string]any{"schema": "crumb-fairplay-capture-v3", "files": entries})
	if err != nil { t.Fatal(err) }
	p := filepath.Join(t.TempDir(), "v3.zip")
	f, err := os.Create(p)
	if err != nil { t.Fatal(err) }
	w := zip.NewWriter(f)
	for name, blob := range blobs {
		entry, err := w.Create("run/" + name)
		if err != nil { t.Fatal(err) }
		if _, err = entry.Write(blob); err != nil { t.Fatal(err) }
	}
	entry, err := w.Create("run/manifest.json")
	if err != nil { t.Fatal(err) }
	if _, err = entry.Write(manifest); err != nil { t.Fatal(err) }
	if err = w.Close(); err != nil { t.Fatal(err) }
	if err = f.Close(); err != nil { t.Fatal(err) }
	return p
}

func TestAnalyzeV3ConsumerStagesAndPrivacy(t *testing.T) {
	cases := []struct {
		name string
		overrides map[string]any
		stage string
		missing string
	}{
		{"rejected", nil, "native_rejected", "accepted_keybag_import"},
		{"zero-return", map[string]any{"native_status": 0, "import_succeeded": true},
			"native_returned_zero", ""},
		{"emulator-failed", map[string]any{"native_status": 0, "emulator_completed": false,
			"state_files": map[string][]byte{}}, "emulator_failed", "completed_keybag_import"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := analyzeCapture(captureV3Fixture(t, tc.overrides))
			if err != nil { t.Fatal(err) }
			if result["ready_to_decrypt"] != false { t.Fatal("unvalidated model was accepted") }
			summary := result["transactions"].([]map[string]any)[0]
			if summary["keybag_consumer"] != tc.stage || summary["ready_to_recover"] != false {
				t.Fatalf("incorrect readiness for %s", tc.name)
			}
			missing := summary["missing_stages"].([]string)
			if tc.missing != "" && (len(missing) < 1 || missing[0] != tc.missing) {
				t.Fatalf("wrong missing stage: %v", missing)
			}
			if missing[len(missing)-1] != "validated_enrollment_model" {
				t.Fatal("unverified model treated as recovered")
			}
			if tc.name == "rejected" {
				if summary["consumer_native_status"] != int32(-42052) ||
					summary["consumer_state_bytes"] != len("SYNTHETIC-SECRET-CANARY") ||
					summary["consumer_state_nonempty_files"] != 1 {
					t.Fatal("consumer result counts do not match")
				}
			}
			data, err := json.Marshal(result)
			if err != nil { t.Fatal(err) }
			if strings.Contains(string(data), "SYNTHETIC-SECRET-CANARY") ||
				strings.Contains(string(data), "hardware_id") ||
				strings.Contains(string(data), "dsid") {
				t.Fatal("private state or identity leaked into the diagnostic")
			}
		})
	}
}

func TestAnalyzeV3RejectsContradictoryConsumerReports(t *testing.T) {
	cases := []map[string]any{
		{"native_status": -42052, "import_succeeded": true},
		{"emulator_completed": false, "import_succeeded": true},
		{"portable_model_ready": true},
		{"native_status": nil},
		{"state_files": map[string][]byte{"../escape.sidb": []byte("bad")}},
	}
	for i, tc := range cases {
		if _, err := analyzeCapture(captureV3Fixture(t, tc)); err == nil {
			t.Fatalf("accepted contradictory consumer report case %d", i)
		}
	}
}


// V6 shape is modeled with synthetic bytes: a successful native call can
// create a nonempty SIDB and a null/empty lock-file entry, but must never be
// mistaken for verified enrollment or permission to decrypt any IPA.
func TestAnalyzeV3AcceptedConsumerWithMixedVirtualFiles(t *testing.T) {
	const sentinel = "NEVER-EXPORT-PRIVATE-STATE-CANARY"
	sidb := make([]byte, 862)
	copy(sidb, []byte(sentinel))
	p := captureV3Fixture(t, map[string]any{
		"native_status":        0,
		"emulator_completed":    true,
		"import_succeeded":     true,
		"portable_model_ready": false,
		"state_files": map[string][]byte{
			".sid.lockFile": nil,
			"SC Info.sidb": sidb,
		},
		"virtual_io": map[string]any{
			"total_calls": 21,
			"truncated":   false,
			"recent": []map[string]any{
				{"operation": "_FSSetCatalogInfo", "outcome": "zero", "mask": 3072},
			},
		},
	})
	result, err := analyzeCapture(p)
	if err != nil {
		t.Fatal(err)
	}
	if result["ready_to_decrypt"] != false {
		t.Fatal("native import was incorrectly promoted to decryption readiness")
	}
	transactions := result["transactions"].([]map[string]any)
	if len(transactions) != 1 {
		t.Fatalf("expected exactly one synthetic transaction, got %d", len(transactions))
	}
	got := transactions[0]
	if got["keybag_consumer"] != "native_returned_zero" ||
		got["consumer_import_succeeded"] != true ||
		got["consumer_native_status"] != int32(0) ||
		got["consumer_state_files"] != 2 ||
		got["consumer_state_nonempty_files"] != 1 ||
		got["consumer_state_bytes"] != 862 {
		t.Fatalf("incorrect sanitized accepted-import counts: %+v", got)
	}
	bridge, ok := got["enrollment_bridge"].(map[string]any)
	if !ok || bridge["legacy_import_accepted"] != true ||
		bridge["legacy_sidb_observed"] != true ||
		bridge["legacy_sidb_bytes"] != 862 ||
		bridge["portable_raw_state_verified"] != false ||
		bridge["portable_enrollment_frame_verified"] != false ||
		bridge["validated_portable_model"] != false ||
		bridge["status"] != "missing_portable_enrollment_source" {
		t.Fatalf("successful legacy import must not be considered portable enrollment: %+v", bridge)
	}
	stages := got["missing_stages"].([]string)
	if len(stages) != 2 || stages[0] != "raw_state_orchestration" ||
		stages[1] != "validated_enrollment_model" {
		t.Fatalf("accepted keybag must still require distinct enrollment stages: %v", stages)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	for _, canary := range []string{sentinel, "SC Info.sidb", ".sid.lockFile", "virtual_io"} {
		if strings.Contains(string(encoded), canary) {
			t.Fatal("private virtual filesystem data or diagnostic leaked into summary")
		}
	}
}


// Real-world accepted native states have variable sizes; these are *public
// synthetic fixtures* and must not imply a fixed SIDB size or recovered key.
func TestAnalyzeAcceptedNativeSIDBVariableLengths(t *testing.T) {
	for _, size := range []int{862, 883} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			fixture := make([]byte, size)
			for i := range fixture { fixture[i] = byte((i * 13 + 7) & 0xff) }
			p := captureV3Fixture(t, map[string]any{
				"native_status": 0,
				"import_succeeded": true,
				"state_files": map[string][]byte{
					".sid.lockFile": nil,
					"SC Info.sidb": fixture,
				},
			})
			result, err := analyzeCapture(p)
			if err != nil { t.Fatal(err) }
			if result["ready_to_decrypt"] != false {
				t.Fatal("variable legacy state incorrectly promoted to portable model")
			}
			transactions := result["transactions"].([]map[string]any)
			if len(transactions) != 1 { t.Fatalf("unexpected transactions count %d", len(transactions)) }
			got := transactions[0]
			bridge, ok := got["enrollment_bridge"].(map[string]any)
			if !ok || bridge["legacy_sidb_bytes"] != size ||
				bridge["legacy_sidb_observed"] != true ||
				bridge["portable_raw_state_verified"] != false ||
				bridge["portable_enrollment_frame_verified"] != false ||
				bridge["validated_portable_model"] != false {
				t.Fatalf("invalid portability classification for synthetic size %d", size)
			}
		})
	}
}
