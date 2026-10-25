// Copied into the pinned internal/fairplay package by the native build.
package fairplay

import (
	"archive/zip"
	"fmt"
	"os"
	"path"
	"strings"
)

// Bind by executable name first, then by the nearest enclosing bundle's sole
// SC_Info license. Never inherit a containing app's license across an appex/app
// boundary, and never guess between licenses based on ZIP order.
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
		return "", fmt.Errorf("%s: ambiguous executable SINF", member)
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
			return "", fmt.Errorf("%s: ambiguous bundle SINF", member)
		}
		if strings.HasSuffix(dir, ".app") || strings.HasSuffix(dir, ".appex") {
			break
		}
	}
	return "", fmt.Errorf("%s: no associated SINF", member)
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
		return nil, err
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
	seen := map[string]bool{}
	total := 0
	for _, image := range inventory.images {
		sinf, err := imageSINFMember(image.member, names)
		if err != nil {
			return nil, err
		}
		if seen[sinf] {
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
			if file == nil || file.UncompressedSize64 > 64<<20 {
				return nil, fmt.Errorf("portable enrollment input limit")
			}
			data, err := zipBytes(file)
			if err != nil {
				return nil, err
			}
			input.Files[item.key] = data
			total += len(data)
			if total > 128<<20 {
				return nil, fmt.Errorf("portable enrollment input limit")
			}
			if item.member == sinf {
				result[len(result)-1].SINFSHA256 = digest(data)
			}
		}
	}
	good = true
	return result, nil
}
