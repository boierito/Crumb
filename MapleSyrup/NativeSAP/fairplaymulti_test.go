package fairplay

import (
	"archive/zip"
	"os"
	"path/filepath"
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
		{"extension cannot inherit", "Payload/A.app/PlugIns/B.appex/B", []string{"Payload/A.app/SC_Info/A.sinf"}, ""},
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
