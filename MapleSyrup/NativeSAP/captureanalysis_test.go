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
