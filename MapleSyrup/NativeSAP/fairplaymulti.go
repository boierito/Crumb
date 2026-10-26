// Copied into the pinned internal/fairplay package by the native build.
package fairplay

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

// Bind by executable name first, then by the nearest enclosing bundle's sole
// SC_Info license. Extensions and frameworks can share their containing app's
// license. Never cross an app boundary or guess between licenses by ZIP order.
func imageSINFMember(member string, names map[string]bool) (string, error) {
	directory, base := path.Dir(member), path.Base(member)
	var exact []string
	for _, suffix := range []string{".sinf", ".v5.sinf"} {
		candidate := directory + "/SC_Info/" + base + suffix
		if names[candidate] {
			exact = append(exact, candidate)
		}
	}
	if len(exact) == 1 {
		return exact[0], nil
	}
	if len(exact) > 1 {
		return "", enrollmentInputFailure("ipa_sinf_ambiguous", nil)
	}
	for dir := directory; dir != "." && dir != "/"; dir = path.Dir(dir) {
		var candidates []string
		prefix := dir + "/SC_Info/"
		for name := range names {
			if strings.HasPrefix(name, prefix) && path.Dir(name) == dir+"/SC_Info" && strings.HasSuffix(name, ".sinf") {
				candidates = append(candidates, name)
			}
		}
		if len(candidates) == 1 {
			return candidates[0], nil
		}
		if len(candidates) > 1 {
			return "", enrollmentInputFailure("ipa_sinf_ambiguous", nil)
		}
		if strings.HasSuffix(dir, ".app") {
			break
		}
	}
	return "", enrollmentInputFailure("ipa_sinf_missing", nil)
}

// EnrollmentInputError exposes a fixed diagnostic code, never archive paths or
// license bytes. Keep the private cause for internal error chaining only.
type EnrollmentInputError struct {
	Code  string
	cause error
}

func (e *EnrollmentInputError) Error() string { return e.Code }
func (e *EnrollmentInputError) Unwrap() error { return e.cause }
func enrollmentInputFailure(code string, cause error) error {
	var classified *EnrollmentInputError
	if errors.As(cause, &classified) {
		return cause
	}
	return &EnrollmentInputError{Code: code, cause: cause}
}

type imageMaterial struct {
	sinfMember string
	material   KeyMaterial
}

func resolveImageMaterials(inventory ipaInventory, files map[string]*zip.File, modelPath string) (map[string]imageMaterial, error) {
	var packs [][]byte
	registry := false
	if modelPath != "" {
		info, err := os.Stat(modelPath)
		if err != nil {
			return nil, fmt.Errorf("inspect FairPlay enrollment model: %w", err)
		}
		registry = info.IsDir()
		if registry {
			packs, err = readEnrollmentModelRegistry(modelPath)
		} else {
			var pack []byte
			pack, err = readEnrollmentModel(modelPath)
			if err == nil {
				packs = [][]byte{pack}
			}
		}
		if err != nil {
			return nil, err
		}
		defer clearModelPacks(packs)
	}
	return bindImageMaterials(inventory, files, func(sinf []byte) (KeyMaterial, error) {
		if len(packs) == 0 {
			return ResolveLogicalKey(sinf)
		}
		if registry {
			return ResolveLogicalKeyWithModels(sinf, packs)
		}
		return ResolveLogicalKeyWithModel(sinf, packs[0])
	})
}

// Resolve every domain before creating output. All slices of one member share
// a license, while separate members can select different enrollment models.
func bindImageMaterials(inventory ipaInventory, files map[string]*zip.File, resolve func([]byte) (KeyMaterial, error)) (map[string]imageMaterial, error) {
	names := make(map[string]bool, len(files))
	for name := range files {
		names[name] = true
	}
	result := map[string]imageMaterial{}
	cache := map[string]KeyMaterial{}
	for _, image := range inventory.images {
		if _, found := result[image.member]; found {
			continue
		}
		name, err := imageSINFMember(image.member, names)
		if err != nil {
			return nil, err
		}
		material, found := cache[name]
		if !found {
			sinf, err := zipBytes(files[name])
			if err != nil {
				return nil, err
			}
			material, err = resolve(sinf)
			if err != nil {
				return nil, fmt.Errorf("%s: resolve associated SINF: %w", image.member, err)
			}
			if material.Evidence == nil {
				material.Evidence = map[string]any{}
			}
			material.Evidence["request_ipa_sha256"] = inventory.sha256
			material.Evidence["request_sinf_sha256"] = digest(sinf)
			cache[name] = material
		}
		result[image.member] = imageMaterial{name, material}
	}
	return result, nil
}

// PortableEnrollmentInput is private input material, not a diagnostic object.
type PortableEnrollmentInput struct {
	ExecutablePath, SINFSHA256 string
	Files                      map[string][]byte
}

func PortableEnrollmentInputs(source string) ([]PortableEnrollmentInput, error) {
	inventory, archive, err := openIPAInventory(source)
	if err != nil {
		return nil, enrollmentInputFailure("ipa_structure_invalid", err)
	}
	defer archive.Close()
	files := map[string]*zip.File{}
	names := map[string]bool{}
	for _, file := range archive.File {
		files[file.Name] = file
		names[file.Name] = true
	}
	var result []PortableEnrollmentInput
	good := false
	defer func() {
		if !good {
			for _, input := range result {
				for _, data := range input.Files {
					clear(data)
				}
			}
		}
	}()
	// Prefer the executable beside the owning SC_Info directory. An extension
	// preceding its main executable in ZIP order must not become the enrollment
	// target just because they share a license.
	owners := map[string]string{}
	for _, image := range inventory.images {
		sinf, err := imageSINFMember(image.member, names)
		if err != nil {
			return nil, err
		}
		previous := owners[sinf]
		score := func(member string) int {
			if path.Dir(member) == path.Dir(path.Dir(sinf)) {
				return 0
			}
			return strings.Count(member, "/") + 1
		}
		if previous == "" || score(image.member) < score(previous) ||
			(score(image.member) == score(previous) && image.member < previous) {
			owners[sinf] = image.member
		}
	}
	sort.Slice(inventory.images, func(i, j int) bool { return inventory.images[i].member < inventory.images[j].member })
	seen := map[string]bool{}
	total := 0
	for _, image := range inventory.images {
		sinf, err := imageSINFMember(image.member, names)
		if err != nil {
			return nil, err
		}
		if seen[sinf] || owners[sinf] != image.member {
			continue
		}
		seen[sinf] = true
		input := PortableEnrollmentInput{ExecutablePath: "/" + image.member, Files: map[string][]byte{}}
		result = append(result, input)
		for _, item := range []struct{ member, key string }{
			{image.member, "/" + path.Base(image.member)},
			{sinf, path.Base(image.member) + ".sinf"},
			{image.supfMember, path.Base(image.supfMember)},
		} {
			file := files[item.member]
			if file == nil {
				return nil, enrollmentInputFailure("ipa_input_member_missing", nil)
			}
			if file.UncompressedSize64 > 64<<20 {
				return nil, enrollmentInputFailure("ipa_input_size_limit", nil)
			}
			data, err := zipBytes(file)
			if err != nil {
				return nil, enrollmentInputFailure("ipa_input_read_failed", err)
			}
			input.Files[item.key] = data
			total += len(data)
			if total > 128<<20 {
				return nil, enrollmentInputFailure("ipa_input_size_limit", nil)
			}
			if item.member == sinf {
				result[len(result)-1].SINFSHA256 = digest(data)
			}
		}
	}
	good = true
	return result, nil
}

// decryptionStageError retains a private cause but exposes only a fixed stage.
// The wrapper preserves the innermost stage when image validation fails.
type decryptionStageError struct {
	stage string
	cause error
}

func (e *decryptionStageError) Error() string { return "decryption_" + e.stage + "_failed" }
func (e *decryptionStageError) Unwrap() error { return e.cause }
func decryptionFailure(stage string, cause error) error {
	var classified *decryptionStageError
	if errors.As(cause, &classified) {
		return cause
	}
	return &decryptionStageError{stage, cause}
}
func DecryptionFailureStage(err error) (string, bool) {
	var classified *decryptionStageError
	if !errors.As(err, &classified) {
		return "", false
	}
	switch classified.stage {
	case "platform", "paths", "output-directory", "ipa-inventory", "model-sinf-binding",
		"temporary-output", "archive-copy", "image-read", "image-inventory", "supf",
		"image-decryption", "archive-write", "output-sync", "report-write", "output-publish",
		"image-geometry", "codedirectory", "page-decryption", "arm64-audit", "mutation-bounds":
		return classified.stage, true
	}
	return "", false
}
