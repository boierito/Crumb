package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/state"
)

// The modern capture is a separate, explicit input. A legacy consumer ZIP is
// never silently reinterpreted as a modern enrollment capture.
func readModernEnrollment(root string) (fairplay.EnrollmentCaptureManifest, map[string][]byte, error) {
	var manifest fairplay.EnrollmentCaptureManifest
	files := map[string][]byte{}
	good := false
	defer func() {
		if !good {
			for _, data := range files {
				clear(data)
			}
		}
	}()
	read := func(relative string, limit int64) ([]byte, error) {
		if !filepath.IsLocal(relative) || strings.Contains(relative, "\\") || filepath.ToSlash(filepath.Clean(relative)) != relative {
			return nil, fmt.Errorf("invalid enrollment capture path")
		}
		current := root
		for _, part := range strings.Split(relative, "/") {
			current = filepath.Join(current, part)
			info, err := os.Lstat(current)
			if err != nil || info.Mode()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("invalid enrollment capture file")
			}
		}
		file, err := os.Open(current)
		if err != nil {
			return nil, fmt.Errorf("unreadable enrollment capture file")
		}
		defer file.Close()
		info, err := file.Stat()
		if err != nil || !info.Mode().IsRegular() || info.Size() > limit || info.Mode().Perm()&0077 != 0 {
			return nil, fmt.Errorf("invalid enrollment capture file permissions or size")
		}
		data, err := io.ReadAll(io.LimitReader(file, limit+1))
		if err != nil || int64(len(data)) > limit {
			return nil, fmt.Errorf("enrollment capture limit")
		}
		return data, nil
	}
	info, err := os.Lstat(root)
	if err != nil || !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return manifest, nil, fmt.Errorf("private modern enrollment directory required")
	}
	document, err := read("manifest.json", 1<<20)
	if err != nil {
		return manifest, nil, err
	}
	if json.Unmarshal(document, &manifest) != nil || manifest.Schema != fairplay.EnrollmentCaptureSchema || len(manifest.Files) == 0 || len(manifest.Files) > 256 {
		return manifest, nil, fmt.Errorf("modern enrollment capture schema required")
	}
	total := 0
	for _, entry := range manifest.Files {
		if entry.Path != "app/executable" && entry.Path != manifest.HardwareMAC.Path && !strings.HasPrefix(entry.Path, "app/SC_Info/") && !strings.HasPrefix(entry.Path, "state/") {
			continue
		}
		if _, exists := files[entry.Path]; exists {
			return manifest, nil, fmt.Errorf("duplicate enrollment capture entry")
		}
		if entry.Bytes < 0 || entry.Bytes > 64<<20 {
			return manifest, nil, fmt.Errorf("enrollment capture limit")
		}
		data, err := read(entry.Path, entry.Bytes)
		if err != nil {
			return manifest, nil, err
		}
		sum := sha256.Sum256(data)
		if int64(len(data)) != entry.Bytes || hex.EncodeToString(sum[:]) != entry.SHA256 {
			clear(data)
			return manifest, nil, fmt.Errorf("enrollment capture integrity")
		}
		files[entry.Path] = data
		total += len(data)
		if total > 128<<20 {
			return manifest, nil, fmt.Errorf("enrollment capture limit")
		}
	}
	if len(files["app/executable"]) == 0 || len(files[manifest.HardwareMAC.Path]) != 6 || manifest.App.ExecutablePath == "" || manifest.App.Executable == "" {
		return manifest, nil, fmt.Errorf("modern enrollment executable and hardware identity required")
	}
	good = true
	return manifest, files, nil
}

func modernRuntimeFiles(manifest fairplay.EnrollmentCaptureManifest, files map[string][]byte, expanded bool) (map[string][]byte, error) {
	mappings := map[string][]byte{"/" + manifest.App.Executable: files["app/executable"]}
	stateFiles, sinfs := 0, 0
	for relative, data := range files {
		for _, route := range []struct{ capture, virtual string }{
			{"app/SC_Info/", ""},
			{"state/fpsd-sc-info/", "/var/db/fpsd/SC Info/"},
			{"state/shared-sc-info/", "/Users/Shared/SC Info/"},
			{"state/adi/", "/var/db/fpsd/adi/"},
		} {
			if !strings.HasPrefix(relative, route.capture) {
				continue
			}
			suffix := strings.TrimPrefix(relative, route.capture)
			name := path.Base(suffix)
			if route.capture == "app/SC_Info/" {
				// The model builder requires one installed license, even in
				// the expanded retry. The IPA may contain several licenses.
				if strings.HasSuffix(name, ".sinf") && name != manifest.App.Executable+".sinf" {
					continue
				}
				if !expanded && name != manifest.App.Executable+".sinf" && name != manifest.App.Executable+".supf" && name != manifest.App.Executable+".v5.supf" {
					continue
				}
			} else {
				if !expanded && (route.capture != "state/fpsd-sc-info/" || !strings.HasPrefix(name, "SC Info.")) {
					continue
				}
				stateFiles++
			}
			key := route.virtual + suffix
			if _, exists := mappings[key]; exists {
				return nil, fmt.Errorf("ambiguous enrollment runtime mapping")
			}
			mappings[key] = data
			if strings.HasSuffix(key, ".sinf") {
				sinfs++
			}
			break
		}
	}
	if stateFiles == 0 || sinfs != 1 {
		return nil, fmt.Errorf("modern enrollment state and one installed SINF required")
	}
	return mappings, nil
}

func recoverModernEnrollment(source, ipa, destination string) (any, error) {
	manifest, files, err := readModernEnrollment(source)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, data := range files {
			clear(data)
		}
	}()
	if ipa == "" {
		return nil, fmt.Errorf("target IPA required")
	}
	var recovered rawsource.Result
	for _, expanded := range []bool{false, true} {
		mappings, mapErr := modernRuntimeFiles(manifest, files, expanded)
		if mapErr != nil {
			err = mapErr
			continue
		}
		recovered, err = rawsource.Recover(manifest.App.ExecutablePath, files[manifest.HardwareMAC.Path], mappings)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	defer clear(recovered.Source)
	defer clear(recovered.Frame)
	defer clear(recovered.Pack)
	// Independently exercise the existing straight-line state and frame
	// evaluators against the raw runtime's selected frame before using a pack.
	transformed, err := state.Evaluate(recovered.Source, nil)
	if err != nil {
		return nil, fmt.Errorf("derived portable state cross-check failed")
	}
	defer clear(transformed)
	matched := false
	for offset := 0; offset+16 <= len(transformed); offset += 16 {
		candidate, err := frame.Evaluate(transformed[offset:offset+16], nil)
		if err != nil {
			return nil, fmt.Errorf("derived enrollment frame cross-check failed")
		}
		matched = matched || bytes.Equal(candidate, recovered.Frame)
		clear(candidate)
	}
	if !matched {
		return nil, fmt.Errorf("derived state and enrollment frame disagree")
	}
	if err := os.Mkdir(destination, 0700); err != nil {
		return nil, fmt.Errorf("new private recovery directory required")
	}
	complete := false
	defer func() {
		if !complete {
			os.RemoveAll(destination)
		}
	}()
	model := filepath.Join(destination, "fairplay-enrollment.bin")
	if err := os.WriteFile(model, recovered.Pack, 0600); err != nil {
		return nil, err
	}
	result, err := fairplay.DecryptIPA(fairplay.DecryptOptions{
		InputPath: ipa, OutputPath: filepath.Join(destination, "processed.ipa"),
		ReportPath: filepath.Join(destination, "private-report.json"), ModelPath: model,
	})
	if err != nil {
		return nil, fmt.Errorf("recovered model target IPA validation failed")
	}
	complete = true
	return map[string]any{
		"portable_source_derived": true, "enrollment_frame_derived": true,
		"model_sinf_validated": true, "source_frame_cross_checked": true, "images": result.EncryptedImages,
		"pages": result.DecryptedPages, "completed": true,
		"plaintext_authenticated": false,
	}, nil
}
