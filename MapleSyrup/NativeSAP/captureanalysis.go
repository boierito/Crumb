package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/majd/ipatool/v2/internal/fairplay"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

// Structural validation only. No ciphertext is guessed to be a state slot or key.
func analyzeCapture(source string) (map[string]any, error) {
	z, err := zip.OpenReader(source)
	if err != nil {
		return nil, fmt.Errorf("capture_zip")
	}
	defer z.Close()
	if len(z.File) > 256 {
		return nil, fmt.Errorf("capture_limit")
	}
	files := map[string][]byte{}
	total := 0
	defer func() {
		for _, b := range files {
			clear(b)
		}
	}()
	for _, f := range z.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if f.Mode()&os.ModeSymlink != 0 || strings.Contains(f.Name, "\\") || strings.HasPrefix(f.Name, "/") || path.Clean(f.Name) != f.Name || strings.HasPrefix(f.Name, "../") {
			return nil, fmt.Errorf("capture_path")
		}
		if _, exists := files[f.Name]; exists {
			return nil, fmt.Errorf("capture_duplicate")
		}
		rd, e := f.Open()
		if e != nil {
			return nil, fmt.Errorf("capture_read")
		}
		b, e := io.ReadAll(io.LimitReader(rd, (8<<20)+1))
		rd.Close()
		if e != nil || len(b) > 8<<20 {
			return nil, fmt.Errorf("capture_limit")
		}
		total += len(b)
		if total > 32<<20 {
			return nil, fmt.Errorf("capture_limit")
		}
		files[f.Name] = b
	}
	reports := []map[string]any{}
	for name, data := range files {
		if path.Base(name) != "manifest.json" {
			continue
		}
		var manifest struct {
			Schema string `json:"schema"`
			Route  string `json:"route"`
			Files  []struct{ Name, Bytes, SHA256 string }
		}
		if json.Unmarshal(data, &manifest) != nil || (manifest.Schema != "crumb-fairplay-capture-v1" && manifest.Schema != "crumb-fairplay-capture-v2") {
			return nil, fmt.Errorf("capture_manifest")
		}
		prefix := path.Dir(name)
		seen := map[string]bool{}
		summary := map[string]any{"schema": manifest.Schema, "integrity": "verified", "keybag_consumer": "not implemented", "ready_to_recover": false, "generation_context_metadata": false}
		sinfs := 0
		for _, entry := range manifest.Files {
			if path.Base(entry.Name) != entry.Name || seen[entry.Name] {
				return nil, fmt.Errorf("capture_manifest_path")
			}
			seen[entry.Name] = true
			b, ok := files[path.Join(prefix, entry.Name)]
			if !ok {
				return nil, fmt.Errorf("capture_missing")
			}
			count, e := strconv.Atoi(entry.Bytes)
			sum := sha256.Sum256(b)
			if e != nil || count != len(b) || hex.EncodeToString(sum[:]) != entry.SHA256 {
				return nil, fmt.Errorf("capture_integrity")
			}
			switch {
			case entry.Name == "kbsync.bin":
				if len(b) == 0 {
					return nil, fmt.Errorf("capture_empty")
				}
				summary["kbsync_bytes"] = len(b)
			case entry.Name == "keybag.bin":
				if len(b) == 0 {
					return nil, fmt.Errorf("capture_empty")
				}
				summary["keybag_bytes"] = len(b)
			case entry.Name == "context.json":
				var context map[string]string
				if json.Unmarshal(b, &context) != nil || context["schema"] != "crumb-fairplay-context-v2" {
					return nil, fmt.Errorf("capture_context")
				}
				summary["generation_context_metadata"] = true
			case strings.HasPrefix(entry.Name, "sinf-") && strings.HasSuffix(entry.Name, ".bin"):
				info, e := fairplay.InspectSINF(b)
				if e != nil {
					return nil, fmt.Errorf("capture_sinf")
				}
				sinfs++
				summary["sinf_current_layout"] = info.CurrentLayout
				summary["sinf_private_bytes"] = info.PrivateBytes
			default:
				return nil, fmt.Errorf("capture_unexpected_file")
			}
		}
		if !seen["kbsync.bin"] || sinfs == 0 {
			return nil, fmt.Errorf("capture_incomplete")
		}
		// Refuse unmanifested files, rather than suggesting their integrity was checked.
		for file := range files {
			if path.Dir(file) == prefix && path.Base(file) != "manifest.json" && !seen[path.Base(file)] {
				return nil, fmt.Errorf("capture_unmanifested")
			}
		}
		summary["sinf_count"] = sinfs
		summary["missing_stages"] = []string{"keybag_consumer", "raw_state_orchestration", "validated_enrollment_model"}
		reports = append(reports, summary)
	}
	if len(reports) == 0 {
		return nil, fmt.Errorf("capture_manifest_missing")
	}
	return map[string]any{"transactions": reports, "transaction_count": len(reports), "ready_to_decrypt": false}, nil
}
