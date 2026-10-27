package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
		recovered, err = rawsource.RecoverValidated(manifest.App.ExecutablePath, files[manifest.HardwareMAC.Path], mappings, func(pack []byte) bool {
			valid, _ := probeModelSINF(mappings, pack)
			return valid
		})
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return finishRecoveredEnrollments([]rawsource.Result{recovered}, ipa, destination)
}

// Validation reports contain lengths, counts and fixed status labels only.
// They describe progress, never recovered state, keys, hashes or file paths.
type recoveredValidationDiagnostic struct {
	DomainCount            int  `json:"domain_count"`
	DomainsCrossChecked    int  `json:"domains_cross_checked"`
	SourceBytes            int  `json:"source_bytes"`
	FrameBytes             int  `json:"frame_bytes"`
	PackBytes              int  `json:"pack_bytes"`
	StateEvaluationPassed  bool `json:"state_evaluation_passed"`
	FrameCandidatesChecked int  `json:"frame_candidates_checked"`
	SourceFrameMatched     bool `json:"source_frame_matched"`
	ModelFilesWritten      int  `json:"model_files_written"`
	TargetIPAValidated     bool `json:"target_ipa_validated"`
}

func verifyRecoveredEnrollment(recovered rawsource.Result) error {
	return verifyRecoveredEnrollmentWithReport(recovered, &recoveredValidationDiagnostic{})
}

func verifyRecoveredEnrollmentWithReport(recovered rawsource.Result, report *recoveredValidationDiagnostic) error {
	report.SourceBytes, report.FrameBytes, report.PackBytes = len(recovered.Source), len(recovered.Frame), len(recovered.Pack)
	report.StateEvaluationPassed, report.SourceFrameMatched, report.FrameCandidatesChecked = false, false, 0
	transformed, err := state.Evaluate(recovered.Source, nil)
	if err != nil {
		return recoveryFailure("recovered-validation/state-evaluation", "recovered_state_evaluation_failed", err)
	}
	defer clear(transformed)
	report.StateEvaluationPassed = true
	for offset := 0; offset+16 <= len(transformed); offset += 16 {
		candidate, err := frame.Evaluate(transformed[offset:offset+16], nil)
		if err != nil {
			return recoveryFailure("recovered-validation/frame-evaluation", "recovered_frame_evaluation_failed", err)
		}
		report.FrameCandidatesChecked++
		report.SourceFrameMatched = bytes.Equal(candidate, recovered.Frame) || report.SourceFrameMatched
		clear(candidate)
	}
	if !report.SourceFrameMatched {
		return recoveryFailure("recovered-validation/source-frame-match", "recovered_source_frame_mismatch", nil)
	}
	return nil
}

func finishRecoveredEnrollments(recovered []rawsource.Result, ipa, destination string) (output any, failure error) {
	diagnostic := recoveredValidationDiagnostic{DomainCount: len(recovered)}
	var runtime *rawsource.RecoveryDiagnostic
	defer func() {
		if failure == nil {
			return
		}
		var classified *labStageError
		if !errors.As(failure, &classified) {
			failure = recoveryFailure("recovered-validation", "recovered_output_validation_failed", failure)
			errors.As(failure, &classified)
		}
		classified.validation, classified.runtime = &diagnostic, runtime
	}()
	if len(recovered) == 0 {
		return nil, recoveryFailure("recovered-validation/outputs", "recovered_outputs_missing", nil)
	}
	defer func() {
		for _, value := range recovered {
			clear(value.Source)
			clear(value.Frame)
			clear(value.Pack)
		}
	}()
	for _, value := range recovered {
		snapshot := value.Runtime
		runtime = &snapshot
		if err := verifyRecoveredEnrollmentWithReport(value, &diagnostic); err != nil {
			return nil, err
		}
		diagnostic.DomainsCrossChecked++
	}
	if err := os.Mkdir(destination, 0700); err != nil {
		return nil, recoveryFailure("recovered-validation/output-directory", "recovery_output_directory_failed", err)
	}
	complete := false
	defer func() {
		if !complete {
			os.RemoveAll(destination)
		}
	}()
	model := filepath.Join(destination, "models")
	if err := os.Mkdir(model, 0700); err != nil {
		return nil, recoveryFailure("recovered-validation/model-directory", "recovery_model_directory_failed", err)
	}
	// The authorization flow retains the processed IPA on success. The model
	// is only a temporary input to DecryptIPA and must not survive that flow.
	defer os.RemoveAll(model)
	seen := map[[32]byte]bool{}
	for index, value := range recovered {
		sum := sha256.Sum256(value.Pack)
		if seen[sum] {
			continue
		}
		seen[sum] = true
		if err := os.WriteFile(filepath.Join(model, fmt.Sprintf("model-%d.bin", index)), value.Pack, 0600); err != nil {
			return nil, recoveryFailure("recovered-validation/model-write", "recovery_model_write_failed", err)
		}
		diagnostic.ModelFilesWritten++
	}
	result, err := fairplay.DecryptIPA(fairplay.DecryptOptions{
		InputPath: ipa, OutputPath: filepath.Join(destination, "processed.ipa"),
		ReportPath: filepath.Join(destination, "private-report.json"), ModelPath: model,
	})
	if err != nil {
		stage, code := "target-ipa", "recovered_model_target_ipa_failed"
		if detail, ok := fairplay.DecryptionFailureStage(err); ok {
			stage, code = detail, "recovered_"+strings.ReplaceAll(detail, "-", "_")+"_failed"
		}
		return nil, recoveryFailure("recovered-validation/"+stage, code, err)
	}
	if err := os.RemoveAll(model); err != nil {
		return nil, recoveryFailure("recovered-validation/model-cleanup", "recovery_model_cleanup_failed", err)
	}
	complete = true
	diagnostic.TargetIPAValidated = true
	return map[string]any{
		"validation": diagnostic, "runtime": runtime,
		"portable_source_derived": true, "enrollment_frame_derived": true,
		"model_sinf_validated": true, "source_frame_cross_checked": true, "images": result.EncryptedImages,
		"pages": result.DecryptedPages, "completed": true,
		"plaintext_authenticated": false,
	}, nil
}
