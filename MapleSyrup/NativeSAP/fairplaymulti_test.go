package fairplay

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImageSINFBoundaries(t *testing.T) {
	for _, test := range []struct {
		name, member string
		licenses     []string
		want         string
	}{
		{"main", "Payload/A.app/A", []string{"Payload/A.app/SC_Info/A.sinf"}, "Payload/A.app/SC_Info/A.sinf"},
		{"extension", "Payload/A.app/PlugIns/B.appex/B", []string{"Payload/A.app/SC_Info/A.sinf", "Payload/A.app/PlugIns/B.appex/SC_Info/B.sinf"}, "Payload/A.app/PlugIns/B.appex/SC_Info/B.sinf"},
		{"framework inherits owning app", "Payload/A.app/Frameworks/F.framework/F", []string{"Payload/A.app/SC_Info/A.sinf"}, "Payload/A.app/SC_Info/A.sinf"},
		{"extension shares app license", "Payload/A.app/PlugIns/B.appex/B", []string{"Payload/A.app/SC_Info/A.sinf"}, "Payload/A.app/SC_Info/A.sinf"},
		{"nested app cannot inherit", "Payload/A.app/Watch/W.app/W", []string{"Payload/A.app/SC_Info/A.sinf"}, ""},
		{"extension local ambiguity cannot inherit", "Payload/A.app/PlugIns/B.appex/B", []string{"Payload/A.app/SC_Info/A.sinf", "Payload/A.app/PlugIns/B.appex/SC_Info/C.sinf", "Payload/A.app/PlugIns/B.appex/SC_Info/D.sinf"}, ""},
		{"unrelated app", "Payload/A.app/A", []string{"Payload/B.app/SC_Info/B.sinf"}, ""},
		{"ambiguous exact", "Payload/A.app/A", []string{"Payload/A.app/SC_Info/A.sinf", "Payload/A.app/SC_Info/A.v5.sinf"}, ""},
		{"exact beats unrelated", "Payload/A.app/A", []string{"Payload/A.app/SC_Info/A.sinf", "Payload/A.app/SC_Info/B.sinf"}, "Payload/A.app/SC_Info/A.sinf"},
		{"ambiguous fallback", "Payload/A.app/Frameworks/F.framework/F", []string{"Payload/A.app/SC_Info/A.sinf", "Payload/A.app/SC_Info/B.sinf"}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			names := map[string]bool{}
			for _, name := range test.licenses {
				names[name] = true
			}
			got, err := imageSINFMember(test.member, names)
			if test.want == "" {
				if err == nil {
					t.Fatal("accepted missing or ambiguous association")
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("got %q, %v; want %q", got, err, test.want)
			}
		})
	}
}

func TestMultiSINFSelectsIndependentEnrollmentModels(t *testing.T) {
	root := t.TempDir()
	registry := filepath.Join(root, "models")
	if err := os.Mkdir(registry, 0700); err != nil {
		t.Fatal(err)
	}
	var keys [2][16]byte
	keys[0][0], keys[1][0] = 17, 29
	members := []string{"Payload/A.app/A", "Payload/A.app/PlugIns/B.appex/B"}
	licenses := []string{"Payload/A.app/SC_Info/A.sinf", "Payload/A.app/PlugIns/B.appex/SC_Info/B.sinf"}
	input := filepath.Join(root, "input.ipa")
	file, err := os.Create(input)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	for i := range keys {
		sinf, model := syntheticCurrentSINFV3Fixture(t, SINFSelector{User: uint32(i + 1), Key: uint32(i + 10)}, keys[i])
		if err := os.WriteFile(filepath.Join(registry, string(rune('a'+i))+".bin"), model, 0600); err != nil {
			t.Fatal(err)
		}
		entry, err := writer.Create(licenses[i])
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(sinf); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	archive, err := zip.OpenReader(input)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	files := map[string]*zip.File{}
	for _, entry := range archive.File {
		files[entry.Name] = entry
	}
	inventory := ipaInventory{images: []encryptedImage{{member: members[0]}, {member: members[0]}, {member: members[1]}}}
	resolved, err := resolveImageMaterials(inventory, files, registry)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 2 {
		t.Fatalf("got %d executable domains", len(resolved))
	}
	for i, member := range members {
		got := resolved[member]
		if got.sinfMember != licenses[i] || got.material.LogicalKey != logicalFromInnerKey(keys[i]) {
			t.Fatalf("executable %d selected another domain's material", i)
		}
	}
	if err := os.Remove(filepath.Join(registry, "b.bin")); err != nil {
		t.Fatal(err)
	}
	if _, err := resolveImageMaterials(inventory, files, registry); err == nil {
		t.Fatal("accepted incomplete model registry")
	}
}

// Match the reported topology: two encrypted ARM64 images, one application
// license and one SUPF per image. Put the extension first to catch order bias.
func sharedLicenseIPAFixture(t *testing.T, licenses map[string][]byte) string {
	t.Helper()
	target := filepath.Join(t.TempDir(), "shared.ipa")
	file, err := os.Create(target)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	put := func(name string, data []byte) {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	supf := make([]byte, 60)
	supf[0] = 5
	for offset, value := range map[int]uint32{4: 52, 8: 1, 12: 44, 16: cpuTypeARM64, 20: 0, 24: 1} {
		binary.BigEndian.PutUint32(supf[offset:], value)
	}
	for _, member := range []string{"Payload/A.app/PlugIns/B.appex/B", "Payload/A.app/A"} {
		blob := make([]byte, 8204)
		copy(blob, thinHeader(true, true))
		put(member, blob)
		put(filepath.ToSlash(filepath.Dir(member))+"/SC_Info/"+filepath.Base(member)+".supf", supf)
	}
	for name, data := range licenses {
		put(name, data)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return target
}

func TestPortableInputsSharedAppLicense(t *testing.T) {
	license := []byte("synthetic-license")
	source := sharedLicenseIPAFixture(t, map[string][]byte{"Payload/A.app/SC_Info/A.sinf": license})
	inspection, err := InspectPortableIPA(source)
	if err != nil || inspection["encrypted_images"] != 2 || inspection["sinf_count"] != 1 || inspection["supf_count"] != 2 {
		t.Fatalf("fixture geometry: %v / %v", inspection, err)
	}
	inputs, err := PortableEnrollmentInputs(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(inputs) != 1 || inputs[0].ExecutablePath != "/Payload/A.app/A" ||
		!bytes.Equal(inputs[0].Files["A.sinf"], license) || inputs[0].SINFSHA256 != digest(license) {
		t.Fatalf("wrong domain or representative: %v", inputs)
	}
	if _, found := inputs[0].Files["A.supf"]; !found {
		t.Fatal("missing owner's SUPF")
	}
	archive, err := zip.OpenReader(source)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	files := map[string]*zip.File{}
	for _, file := range archive.File {
		files[file.Name] = file
	}
	inventory := ipaInventory{images: []encryptedImage{{member: "Payload/A.app/A"}, {member: "Payload/A.app/PlugIns/B.appex/B"}}}
	calls := 0
	bound, err := bindImageMaterials(inventory, files, func(blob []byte) (KeyMaterial, error) {
		calls++
		if !bytes.Equal(blob, license) {
			t.Fatal("wrong license")
		}
		return KeyMaterial{}, nil
	})
	if err != nil || calls != 1 || len(bound) != 2 {
		t.Fatalf("shared domain binding: %v, calls=%d, images=%d", err, calls, len(bound))
	}
}

func TestPortableInputsSafeFailureCodes(t *testing.T) {
	for _, test := range []struct {
		licenses map[string][]byte
		code     string
	}{
		{map[string][]byte{}, "ipa_sinf_missing"},
		{map[string][]byte{"Payload/A.app/SC_Info/A.sinf": {}, "Payload/A.app/SC_Info/C.sinf": {}}, "ipa_sinf_ambiguous"},
	} {
		_, err := PortableEnrollmentInputs(sharedLicenseIPAFixture(t, test.licenses))
		var classified *EnrollmentInputError
		if !errors.As(err, &classified) || classified.Code != test.code || err.Error() != test.code {
			t.Fatalf("unsafe or wrong diagnostic: %v", err)
		}
	}
}

func TestDecryptionDiagnosticPreservesInnermostFixedStage(t *testing.T) {
	private := errors.New("PRIVATE-ERROR-AND-PATH")
	inner := decryptionFailure("arm64-audit", private)
	outer := decryptionFailure("image-decryption", inner)
	stage, ok := DecryptionFailureStage(outer)
	if !ok || stage != "arm64-audit" || !errors.Is(outer, private) {
		t.Fatal("lost typed failure or cause")
	}
	if stage, ok := DecryptionFailureStage(decryptionFailure("PRIVATE-STAGE", private)); ok || stage != "" {
		t.Fatal("exposed unrecognized stage")
	}
	if strings.Contains(outer.Error(), "PRIVATE") {
		t.Fatal("cause exposed by public error")
	}
}

func TestDecryptionDiagnosticUsesRealInventoryAndImageGates(t *testing.T) {
	root := t.TempDir()
	_, err := DecryptIPA(DecryptOptions{InputPath: filepath.Join(root, "PRIVATE-MISSING.ipa"), OutputPath: filepath.Join(root, "out.ipa"), ReportPath: filepath.Join(root, "report.json")})
	if stage, ok := DecryptionFailureStage(err); !ok || stage != "paths" {
		t.Fatalf("path gate: %s / %v", stage, err)
	}
	input := filepath.Join(root, "invalid.ipa")
	if err := os.WriteFile(input, []byte("PRIVATE-INVALID-ZIP"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err = DecryptIPA(DecryptOptions{InputPath: input, OutputPath: filepath.Join(root, "out.ipa"), ReportPath: filepath.Join(root, "report.json")})
	if stage, ok := DecryptionFailureStage(err); !ok || stage != "ipa-inventory" {
		t.Fatalf("inventory gate: %s / %v", stage, err)
	}
	_, err = decryptImage(nil, encryptedImage{}, nil, [16]byte{}, nil)
	if stage, ok := DecryptionFailureStage(err); !ok || stage != "image-geometry" {
		t.Fatalf("geometry gate: %s / %v", stage, err)
	}
	if strings.Contains(err.Error(), "PRIVATE") {
		t.Fatal("private error exposed")
	}
	if _, err := os.Stat(filepath.Join(root, "out.ipa")); !os.IsNotExist(err) {
		t.Fatal("invalid output published")
	}
}
