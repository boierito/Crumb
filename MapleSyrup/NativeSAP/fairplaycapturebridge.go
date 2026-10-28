package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
)

type legacyEnrollmentInput struct {
	dsid     uint64
	hardware []byte
	files    map[string][]byte
	sinfs    map[string]bool
}

func (input legacyEnrollmentInput) clear() {
	clear(input.hardware)
	for _, data := range input.files {
		clear(data)
	}
}

// Validate the capture contract, then recheck the exact bytes used below against
// its manifest. The state remains an opaque container, never a sliced frame/key.
func readLegacyEnrollmentInput(source string) (legacyEnrollmentInput, error) {
	input, err := readLegacyEnrollmentInputUnchecked(source)
	if err != nil {
		return input, recoveryFailure("capture-input", "capture_input_invalid", err)
	}
	return input, nil
}

func readLegacyEnrollmentInputUnchecked(source string) (legacyEnrollmentInput, error) {
	var result legacyEnrollmentInput
	if _, err := analyzeCapture(source); err != nil {
		return result, err
	}
	archive, err := zip.OpenReader(source)
	if err != nil {
		return result, err
	}
	defer archive.Close()
	members := map[string]*zip.File{}
	var manifestFile *zip.File
	for _, file := range archive.File {
		members[file.Name] = file
		if path.Base(file.Name) == "manifest.json" {
			if manifestFile != nil {
				return result, fmt.Errorf("one retained transaction required")
			}
			manifestFile = file
		}
	}
	if manifestFile == nil {
		return result, fmt.Errorf("capture manifest required")
	}
	read := func(file *zip.File) ([]byte, error) {
		if file == nil {
			return nil, fmt.Errorf("capture member missing")
		}
		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		data, err := io.ReadAll(io.LimitReader(reader, (8<<20)+1))
		if err != nil || len(data) > 8<<20 {
			return nil, fmt.Errorf("capture limit")
		}
		return data, nil
	}
	document, err := read(manifestFile)
	if err != nil {
		return result, err
	}
	var manifest struct {
		Files []struct{ Name, Bytes, SHA256 string }
	}
	if json.Unmarshal(document, &manifest) != nil {
		return result, fmt.Errorf("capture manifest")
	}
	blobs := map[string][]byte{}
	defer func() {
		for _, data := range blobs {
			clear(data)
		}
	}()
	for _, entry := range manifest.Files {
		if entry.Name != "context.json" && entry.Name != "consumer-state.json" && !(strings.HasPrefix(entry.Name, "sinf-") && strings.HasSuffix(entry.Name, ".bin")) {
			continue
		}
		data, err := read(members[path.Join(path.Dir(manifestFile.Name), entry.Name)])
		if err != nil {
			return result, err
		}
		sum := sha256.Sum256(data)
		size, e := strconv.Atoi(entry.Bytes)
		if e != nil || size != len(data) || hex.EncodeToString(sum[:]) != entry.SHA256 {
			clear(data)
			return result, fmt.Errorf("capture integrity")
		}
		blobs[entry.Name] = data
	}
	result, err = decodeLegacyEnrollment(blobs["context.json"], blobs["consumer-state.json"])
	if err != nil {
		return result, err
	}
	result.sinfs = map[string]bool{}
	for name, data := range blobs {
		if strings.HasPrefix(name, "sinf-") {
			sum := sha256.Sum256(data)
			result.sinfs[hex.EncodeToString(sum[:])] = true
		}
	}
	return result, nil
}

func decodeLegacyEnrollment(context, consumer []byte) (result legacyEnrollmentInput, err error) {
	var identity struct {
		Schema   string `json:"schema"`
		Hardware []byte `json:"hardware_id"`
		Producer string `json:"producer"`
		DSID     string `json:"dsid"`
	}
	var state struct {
		Consumer  string            `json:"consumer"`
		Completed bool              `json:"emulator_completed"`
		Imported  bool              `json:"import_succeeded"`
		Status    int32             `json:"native_status"`
		Files     map[string][]byte `json:"state_files"`
	}
	defer func() {
		if err != nil {
			clear(identity.Hardware)
			for _, data := range state.Files {
				clear(data)
			}
		}
	}()
	if e := json.Unmarshal(context, &identity); e != nil {
		return result, recoveryFailure("identity", "capture_context_json_invalid", e)
	}
	if identity.Schema != "crumb-fairplay-context-v2" {
		return result, recoveryFailure("identity", "capture_context_schema_invalid", nil)
	}
	if identity.Producer != "retained-storeagent" {
		return result, recoveryFailure("identity", "capture_producer_mismatch", nil)
	}
	if len(identity.Hardware) != 6 {
		return result, recoveryFailure("identity", "capture_hardware_length_invalid", nil)
	}
	if e := json.Unmarshal(consumer, &state); e != nil {
		return result, recoveryFailure("consumer-state", "consumer_state_json_invalid", e)
	}
	if state.Consumer != "storeagent-10.9-import-keybag" {
		return result, recoveryFailure("consumer-state", "consumer_kind_mismatch", nil)
	}
	if !state.Completed || !state.Imported || state.Status != 0 {
		return result, recoveryFailure("consumer-state", "consumer_import_not_accepted", nil)
	}
	if len(state.Files["SC Info.sidb"]) == 0 {
		return result, recoveryFailure("consumer-state", "legacy_sidb_missing", nil)
	}
	if len(state.Files) > 32 {
		return result, recoveryFailure("consumer-state", "state_file_count_limit", nil)
	}
	total := 0
	for name, data := range state.Files {
		// Match the flat isolated guest filesystem contract. Lock names do
		// not have to share the SIDB's prefix. No host paths are exposed.
		if name == "" || name == "." || name == ".." || path.Base(name) != name || strings.Contains(name, "\\") {
			return result, recoveryFailure("state-files", "state_filename_invalid", nil)
		}
		total += len(data)
		if len(data) > 4<<20 || total > 4<<20 {
			return result, recoveryFailure("state-files", "state_byte_limit", nil)
		}
	}
	if identity.DSID != "" {
        result.dsid, err = strconv.ParseUint(identity.DSID, 10, 64)
        if err != nil || result.dsid == 0 {
            return result, recoveryFailure("identity", "capture_dsid_invalid", nil)
        }
    }
    result.hardware, result.files = identity.Hardware, state.Files
	return result, nil
}

// Verify only the actual license-to-capture association. This deliberately
// does NOT run a modern fairplayd interpreter, derive a model, or claim that
// the current StoreAgent enrollment can decrypt any binary. It can be used
// to reopen a genuine original IPA with its retained private transaction
// even when the cross-generation emulator is unsupported.
func capturedSINFMatchesDomains(sinfs map[string]bool, domains []fairplay.PortableEnrollmentInput) bool {
    return len(domains) == 1 && domains[0].SINFSHA256 != "" && sinfs[domains[0].SINFSHA256]
}

func verifyCaptureIPABinding(source, ipa string) (any, error) {
    input, err := readLegacyEnrollmentInput(source)
    if err != nil { return nil, recoveryFailure("offline-binding", "capture_input_invalid", err) }
    defer input.clear()
    domains, err := fairplay.PortableEnrollmentInputs(ipa)
    if err != nil { return nil, enrollmentInventoryFailure(err) }
    defer func() {
        for _, domain := range domains {
            for _, data := range domain.Files { clear(data) }
        }
    }()
    if !capturedSINFMatchesDomains(input.sinfs, domains) {
        return nil, recoveryFailure("offline-binding", "capture_ipa_mismatch", nil)
    }
    return map[string]any{
        "schema": "crumb-capture-ipa-binding-v1",
        "matched": true,
        "encrypted_domain_count": len(domains),
        "model_validated": false,
        "plaintext_authenticated": false,
        "private_material_exported": false,
    }, nil
}

func recoverCaptureEnrollment(source, ipa, destination string) (any, error) {
	return recoverCaptureWithAccountState(source, ipa, destination, "")
}

func recoverCaptureWithAccountState(source, ipa, destination, accountState string) (any, error) {
	input, err := readLegacyEnrollmentInput(source)
	if err != nil {
		return nil, err
	}
	defer input.clear()
	stateFiles := input.files
	if accountState != "" {
		stored, err := readAuthorizedStateFile(accountState, input.dsid, input.hardware)
		if err != nil { return nil, recoveryFailure("authorized-state", "authorized_state_invalid", err) }
		defer stored.clear()
		stateFiles = stored.Files
	}
	domains, err := fairplay.PortableEnrollmentInputs(ipa)
	if err != nil {
		return nil, enrollmentInventoryFailure(err)
	}
	defer func() {
		for _, domain := range domains {
			for _, data := range domain.Files {
				clear(data)
			}
		}
	}()
	var recovered []rawsource.Result
	defer func() {
		for _, value := range recovered {
			clear(value.Source)
			clear(value.Frame)
			clear(value.Pack)
		}
	}()
	for _, domain := range domains {
		if !input.sinfs[domain.SINFSHA256] {
			return nil, recoveryFailure("transaction-binding", "capture_ipa_mismatch", nil)
		}
		mappings := map[string][]byte{}
		for name, data := range domain.Files {
			mappings[name] = data
		}
		// Mount unchanged legacy state at the runtime's state locations. This
		// is a compatibility attempt; only the validation gates can accept it.
		for name, data := range stateFiles {
			mappings["/var/db/fpsd/SC Info/"+name] = data
			mappings["/Users/Shared/SC Info/"+name] = data
		}
		value, err := rawsource.RecoverValidated(domain.ExecutablePath, input.hardware, mappings, func(pack []byte) bool {
			valid, _ := probeModelSINF(domain.Files, pack)
			return valid
		})
		if err != nil {
			return nil, recoveryFailure("portable-runtime", "legacy_portable_runtime_rejected", err)
		}
		recovered = append(recovered, value)
	}
	result, err := finishRecoveredEnrollments(recovered, ipa, destination)
	if err != nil {
		return nil, recoveryFailure("recovered-validation", "recovered_output_validation_failed", err)
	}
	return result, nil
}

// Initial enrollment may cover several license domains. Compare complete SINF
// sets; never pick a license by ZIP order or assume one executable per IPA.
// Capture bootstrap requires the complete license set. Once a matching account
// state is saved, each later IPA is validated against its own original SINFs.
func capturedSINFMatchesAllDomains(sinfs map[string]bool, domains []fairplay.PortableEnrollmentInput) bool {
	if len(domains) == 0 || len(sinfs) == 0 { return false }
	seen := map[string]bool{}
	for _, domain := range domains {
		if domain.SINFSHA256 == "" || !sinfs[domain.SINFSHA256] { return false }
		seen[domain.SINFSHA256] = true
	}
	return len(seen) == len(sinfs)
}
