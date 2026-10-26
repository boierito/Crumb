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
		if json.Unmarshal(data, &manifest) != nil || (manifest.Schema != "crumb-fairplay-capture-v1" && manifest.Schema != "crumb-fairplay-capture-v2" && manifest.Schema != "crumb-fairplay-capture-v3") {
			return nil, fmt.Errorf("capture_manifest")
		}
		prefix := path.Dir(name)
		seen := map[string]bool{}
		consumerStage := "not_observed"
		summary := map[string]any{"schema": manifest.Schema, "integrity": "verified", "keybag_consumer": consumerStage, "ready_to_recover": false, "generation_context_metadata": false}
		sinfs := 0
        licenseFamily := ""
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
			case entry.Name == "consumer-state.json":
				var report struct {
					Consumer  string            `json:"consumer"`
					Status    *int32            `json:"native_status"`
					Completed *bool             `json:"emulator_completed"`
					Imported  *bool             `json:"import_succeeded"`
					Ready     *bool             `json:"portable_model_ready"`
					Files     map[string][]byte `json:"state_files"`
				}
				if json.Unmarshal(b, &report) != nil || report.Consumer == "" ||
					report.Status == nil || report.Completed == nil || report.Imported == nil ||
					report.Ready == nil || *report.Ready ||
					*report.Imported != (*report.Completed && *report.Status == 0) ||
					len(report.Files) > 32 {
					return nil, fmt.Errorf("capture_consumer_report")
				}
				// JSON encodes state-file bytes as base64. Clear decoded data, even on rejection.
				defer func() {
					for _, data := range report.Files {
						clear(data)
					}
				}()
				bytes, nonempty := 0, 0
				for name, data := range report.Files {
					if name == "" || name == "." || path.Base(name) != name ||
						strings.Contains(name, "\\") || len(data) > 4<<20 {
						return nil, fmt.Errorf("capture_consumer_report")
					}
					bytes += len(data)
					if len(data) > 0 {
						nonempty++
					}
					if bytes > 4<<20 {
						return nil, fmt.Errorf("capture_consumer_report")
					}
				}
				summary["consumer_native_status"] = *report.Status
				summary["consumer_emulator_completed"] = *report.Completed
				summary["consumer_import_succeeded"] = *report.Imported
				summary["consumer_state_files"] = len(report.Files)
				summary["consumer_state_nonempty_files"] = nonempty
				summary["consumer_state_bytes"] = bytes
				// This explicitly separates the legacy StoreAgent virtual state
				// from the modern, build-specific portable enrollment inputs.
				// Never infer a 16-byte frame or 192/224-byte raw-state source
				// from a successful call or an opaque SIDB container.
				legacySIDB := report.Files["SC Info.sidb"]
				summary["enrollment_bridge"] = map[string]any{
					"legacy_import_accepted":          *report.Imported,
					"legacy_sidb_observed":            len(legacySIDB) > 0,
					"legacy_sidb_bytes":               len(legacySIDB),
					"portable_raw_state_verified":      false,
					"portable_enrollment_frame_verified": false,
					"validated_portable_model":         false,
					"status":                           "missing_portable_enrollment_source",
				}
				switch {
				case !*report.Completed:
					consumerStage = "emulator_failed"
				case *report.Imported:
					consumerStage = "native_returned_zero"
				default:
					consumerStage = "native_rejected"
				}
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
                // SINF key ID is the license generation, not the account's
                // age. Report only a fixed family label; never emit user ID,
                // selector, private payload, or per-app identifier.
                selector, selectorErr := fairplay.ParseSINFSelector(b)
                if selectorErr != nil { return nil, fmt.Errorf("capture_sinf_selector") }
                family := "non-key5"
                if selector.Key == 5 { family = "key5" }
                if licenseFamily == "" { licenseFamily = family
                } else if licenseFamily != family { licenseFamily = "mixed" }
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
        summary["sinf_key_family"] = licenseFamily
		summary["keybag_consumer"] = consumerStage
		missing := []string{}
		switch consumerStage {
		case "not_observed":
			missing = append(missing, "observed_keybag_import")
		case "emulator_failed":
			missing = append(missing, "completed_keybag_import")
		case "native_rejected":
			missing = append(missing, "accepted_keybag_import")
		}
		summary["missing_stages"] = append(missing, "raw_state_orchestration", "validated_enrollment_model")
		reports = append(reports, summary)
	}
	if len(reports) == 0 {
		return nil, fmt.Errorf("capture_manifest_missing")
	}
	return map[string]any{"transactions": reports, "transaction_count": len(reports), "ready_to_decrypt": false}, nil
}
