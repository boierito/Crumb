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

func decodeLegacyEnrollment(context, consumer []byte) (legacyEnrollmentInput, error) {
	var identity struct {
		Schema   string `json:"schema"`
		Hardware []byte `json:"hardware_id"`
		Producer string `json:"producer"`
	}
	var state struct {
		Consumer  string            `json:"consumer"`
		Completed bool              `json:"emulator_completed"`
		Imported  bool              `json:"import_succeeded"`
		Status    int32             `json:"native_status"`
		Files     map[string][]byte `json:"state_files"`
	}
	result := legacyEnrollmentInput{}
	if json.Unmarshal(context, &identity) != nil || json.Unmarshal(consumer, &state) != nil || identity.Schema != "crumb-fairplay-context-v2" || identity.Producer != "retained-storeagent" || len(identity.Hardware) != 6 || state.Consumer != "storeagent-10.9-import-keybag" || !state.Completed || !state.Imported || state.Status != 0 || len(state.Files["SC Info.sidb"]) == 0 || len(state.Files) > 32 {
		clear(identity.Hardware)
		for _, data := range state.Files {
			clear(data)
		}
		return result, fmt.Errorf("accepted retained consumer and original hardware identity required")
	}
	total := 0
	for name, data := range state.Files {
		total += len(data)
		if path.Base(name) != name || strings.Contains(name, "\\") || !strings.HasPrefix(name, "SC Info.") || len(data) > 4<<20 || total > 4<<20 {
			clear(identity.Hardware)
			for _, data := range state.Files {
				clear(data)
			}
			return result, fmt.Errorf("invalid retained state files")
		}
	}
	result.hardware, result.files = identity.Hardware, state.Files
	return result, nil
}

func recoverCaptureEnrollment(source, ipa, destination string) (any, error) {
	input, err := readLegacyEnrollmentInput(source)
	if err != nil {
		return nil, err
	}
	defer input.clear()
	domains, err := fairplay.PortableEnrollmentInputs(ipa)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("target IPA SINF is not bound to captured transaction")
		}
		mappings := map[string][]byte{}
		for name, data := range domain.Files {
			mappings[name] = data
		}
		// Mount unchanged legacy state at the runtime's state locations. This
		// is a compatibility attempt; only the validation gates can accept it.
		for name, data := range input.files {
			mappings["/var/db/fpsd/SC Info/"+name] = data
			mappings["/Users/Shared/SC Info/"+name] = data
		}
		value, err := rawsource.Recover(domain.ExecutablePath, input.hardware, mappings)
		if err != nil {
			return nil, fmt.Errorf("legacy state not compatible with portable enrollment profile")
		}
		recovered = append(recovered, value)
	}
	return finishRecoveredEnrollments(recovered, ipa, destination)
}
