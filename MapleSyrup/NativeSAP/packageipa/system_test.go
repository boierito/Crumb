package packageipa

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howett.net/plist"
)

func systemFixtureMachO(encrypted bool, subtype uint32) ([]byte, []byte) {
	data := bytes.Repeat([]byte{0x5a}, 65536)
	clear(data[:16384])
	le := binary.LittleEndian
	for i, v := range []uint32{0xfeedfacf, 0x100000c, subtype, 2, 3, 112, 0, 0} {
		le.PutUint32(data[4*i:], v)
	}
	le.PutUint32(data[32:], 0x19)
	le.PutUint32(data[36:], 72)
	copy(data[40:], "__TEXT")
	le.PutUint64(data[56:], 0x100000000)
	le.PutUint64(data[64:], 65536)
	le.PutUint64(data[80:], 65536)
	le.PutUint32(data[104:], 0x2c)
	le.PutUint32(data[108:], 24)
	le.PutUint32(data[112:], 16384)
	le.PutUint32(data[116:], 16384)
	if encrypted {
		le.PutUint32(data[120:], 1)
	}
	le.PutUint32(data[128:], 0x1d)
	le.PutUint32(data[132:], 16)
	le.PutUint32(data[136:], 49152)
	le.PutUint32(data[140:], 4096)
	copy(data[49152:], "personal signing identity: owner@fixture.invalid")
	plain := bytes.Clone(data)
	if encrypted {
		for i := 16384; i < 32768; i++ {
			plain[i] ^= 0xa7
		}
		le.PutUint32(plain[120:], 0)
	}
	return data, plain
}

func systemFixtureFat() ([]byte, []byte) {
	a, ap := systemFixtureMachO(true, 0)
	b, bp := systemFixtureMachO(true, 2)
	data, plain := make([]byte, 147456), make([]byte, 147456)
	be := binary.BigEndian
	be.PutUint32(data, 0xcafebabe)
	be.PutUint32(data[4:], 2)
	for i, off := range []uint32{16384, 81920} {
		q := data[8+i*20:]
		be.PutUint32(q, 0x100000c)
		be.PutUint32(q[4:], uint32(i*2))
		be.PutUint32(q[8:], off)
		be.PutUint32(q[12:], 65536)
		be.PutUint32(q[16:], 14)
	}
	copy(plain, data)
	copy(data[16384:], a)
	copy(data[81920:], b)
	copy(plain[16384:], ap)
	copy(plain[81920:], bp)
	return data, plain
}

func systemFixtureDirectories(t *testing.T) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	source, overrides := filepath.Join(root, "Target"), filepath.Join(root, "native")
	if e := os.MkdirAll(source, 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.MkdirAll(overrides, 0700); e != nil {
		t.Fatal(e)
	}
	return root, source, overrides
}
func systemFixtureWrite(t *testing.T, root, name string, data []byte) {
	t.Helper()
	p := filepath.Join(root, name)
	if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(p, data, 0600); e != nil {
		t.Fatal(e)
	}
}

func TestSystemExportMultipleDomainsAndFATRemovesIdentityWithoutChangingPlaintext(t *testing.T) {
	root, source, overrides := systemFixtureDirectories(t)
	main, clearMain := systemFixtureMachO(true, 0)
	fat, clearFat := systemFixtureFat()
	plain, _ := systemFixtureMachO(false, 0)
	for name, data := range map[string][]byte{"Main": main, "PlugIns/Share.appex/Share": fat, "Frameworks/Plain.framework/Plain": plain, "asset.bin": []byte("original asset"), "Info.plist": []byte("unchanged app info")} {
		systemFixtureWrite(t, source, name, data)
	}
	for name, data := range map[string][]byte{"Main": clearMain, "PlugIns/Share.appex/Share": clearFat} {
		systemFixtureWrite(t, overrides, name, data)
	}
	for _, name := range []string{"SC_Info/account.sinf", "PlugIns/Share.appex/SC_Info/account.sinf", "_CodeSignature/CodeResources", "PlugIns/Share.appex/embedded.mobileprovision", "iTunesMetadata.plist", "StoreKit/receipt", "_MASReceipt/receipt"} {
		systemFixtureWrite(t, source, name, []byte("owner@fixture.invalid"))
	}
	plan, e := PlanSystemTarget(source)
	if e != nil || len(plan.Files) != 2 || plan.EncryptedSlices != 3 {
		t.Fatalf("plan %+v: %v", plan, e)
	}
	output := filepath.Join(root, "clean.ipa")
	report, e := ExportSystemResult(source, overrides, output, []string{"owner@fixture.invalid"})
	if e != nil {
		t.Fatal(e)
	}
	if report.EncryptedSlices != 3 || report.ChangedPages != 3 || report.RemovedSignatures != 4 || report.Files != 5 || report.NoAccountDataGuaranteed || report.PlaintextAuthenticated {
		t.Fatalf("incorrect gate %+v", report)
	}
	z, e := zip.OpenReader(output)
	if e != nil {
		t.Fatal(e)
	}
	defer z.Close()
	for _, f := range z.File {
		if systemPrivatePath(f.Name) {
			t.Fatal("private metadata included")
		}
		r, e := f.Open()
		if e != nil {
			t.Fatal(e)
		}
		data, e := io.ReadAll(r)
		r.Close()
		if e != nil {
			t.Fatal(e)
		}
		if bytes.Contains(data, []byte("owner@fixture.invalid")) {
			t.Fatal("identity survived signature cleanup")
		}
		image, e := systemMachO(bytes.NewReader(data), int64(len(data)))
		if e != nil {
			t.Fatal(e)
		}
		if image != nil {
			for _, s := range image.slices {
				if len(s.signatures) != 0 {
					t.Fatal("signature command remains")
				}
				if !bytes.Equal(data[s.off+49152:s.off+53248], make([]byte, 4096)) {
					t.Fatal("signature bytes remain")
				}
			}
		}
		if strings.HasSuffix(f.Name, "/Main") && !bytes.Equal(data[16384:32768], clearMain[16384:32768]) {
			t.Fatal("plaintext modified during cleanup")
		}
	}
}

func TestSystemRejectsMissingUnchangedAndCorruptedNativeOutputs(t *testing.T) {
	for _, kind := range []string{"missing", "unchanged", "corrupt", "encrypted", "unsupported"} {
		t.Run(kind, func(t *testing.T) {
			root, source, overrides := systemFixtureDirectories(t)
			a, b := systemFixtureMachO(true, 0)
			switch kind {
			case "unchanged":
				copy(b[16384:32768], a[16384:32768])
			case "corrupt":
				b[40000] ^= 1
			case "encrypted":
				binary.LittleEndian.PutUint32(b[120:], 1)
			case "unsupported":
				binary.LittleEndian.PutUint32(a[4:], 0x200000c)
			}
			systemFixtureWrite(t, source, "Main", a)
			if kind != "missing" {
				systemFixtureWrite(t, overrides, "Main", b)
			}
			if kind == "unsupported" {
				if _, e := PlanSystemTarget(source); e == nil {
					t.Fatal("unsupported CPU accepted")
				}
				return
			}
			output := filepath.Join(root, "out.ipa")
			if _, e := ExportSystemResult(source, overrides, output, nil); e == nil {
				t.Fatal("invalid native output accepted")
			}
			if _, e := os.Stat(output); !os.IsNotExist(e) {
				t.Fatal("partial output remains")
			}
		})
	}
}

func TestSystemRejectsKnownIdentifiersAcrossChunkBoundaryAndInNames(t *testing.T) {
	for _, kind := range []string{"ascii", "utf16", "base64", "name"} {
		t.Run(kind, func(t *testing.T) {
			root, source, overrides := systemFixtureDirectories(t)
			a, b := systemFixtureMachO(true, 0)
			systemFixtureWrite(t, source, "Main", a)
			systemFixtureWrite(t, overrides, "Main", b)
			patterns, e := systemPatterns([]string{"owner@fixture.invalid"})
			if e != nil {
				t.Fatal(e)
			}
			index := 0
			if kind == "utf16" {
				index = 2
			}
			if kind == "base64" {
				index = 1
			}
			name := "asset"
			content := append(bytes.Repeat([]byte{0x55}, 65530), patterns[index]...)
			if kind == "name" {
				name = "owner@fixture.invalid"
				content = []byte("data")
			}
			systemFixtureWrite(t, source, name, content)
			output := filepath.Join(root, "out.ipa")
			if _, e := ExportSystemResult(source, overrides, output, []string{"owner@fixture.invalid"}); e == nil {
				t.Fatal("known identifier exported")
			}
			if _, e := os.Stat(output); !os.IsNotExist(e) {
				t.Fatal("rejected output remains")
			}
		})
	}
}

func TestSystemMalformedLayoutsAndNoOverwrite(t *testing.T) {
	for _, kind := range []string{"command", "signature", "crypt", "fat"} {
		t.Run(kind, func(t *testing.T) {
			a, _ := systemFixtureMachO(true, 0)
			switch kind {
			case "command":
				binary.LittleEndian.PutUint32(a[108:], 8)
			case "signature":
				binary.LittleEndian.PutUint32(a[136:], 0xffffffff)
			case "crypt":
				binary.LittleEndian.PutUint32(a[116:], 0xffffffff)
			case "fat":
				a, _ = systemFixtureFat()
				binary.BigEndian.PutUint32(a[8+20+8:], 16384)
			}
			if _, e := systemMachO(bytes.NewReader(a), int64(len(a))); e == nil {
				t.Fatal("malformed layout accepted")
			}
		})
	}
	root, source, overrides := systemFixtureDirectories(t)
	a, b := systemFixtureMachO(true, 0)
	systemFixtureWrite(t, source, "Main", a)
	systemFixtureWrite(t, overrides, "Main", b)
	p := filepath.Join(root, "existing.ipa")
	if e := os.WriteFile(p, []byte("keep"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := ExportSystemResult(source, overrides, p, nil); e == nil {
		t.Fatal("overwrote existing output")
	}
	data, _ := os.ReadFile(p)
	if string(data) != "keep" {
		t.Fatal("existing output removed")
	}
}

func TestSystemPrepareSnapshotsOwnBundleAndRetainsOriginalLicense(t *testing.T) {
	root := t.TempDir()
	bundle := filepath.Join(root, "Crumb.app")
	systemFixtureWrite(t, bundle, "Crumb", []byte("helper executable"))
	helperInfo, _ := plist.Marshal(Info{BundleID: "com.certlium.crumb", Version: "1", Executable: "Crumb"}, plist.BinaryFormat)
	systemFixtureWrite(t, bundle, "Info.plist", helperInfo)
	systemFixtureWrite(t, bundle, SystemPayload+"/old-secret", []byte("not forwarded"))
	info, _ := plist.Marshal(Info{BundleID: "test.app", Version: "1", Executable: "Main"}, plist.BinaryFormat)
	manifest, _ := plist.Marshal(map[string]any{"SinfPaths": []string{"SC_Info/Main.sinf"}, "SinfReplicationPaths": []string{"PlugIns/Share.appex/SC_Info/Share.sinf"}}, plist.BinaryFormat)
	a, _ := systemFixtureMachO(true, 0)
	input := filepath.Join(root, "original.ipa")
	if e := systemZIP(input, func(z *zip.Writer) error {
		for name, data := range map[string][]byte{"Info.plist": info, "Main": a, "SC_Info/Manifest.plist": manifest, "SC_Info/Main.sinf": []byte("original private license")} {
			if e := systemPut(z, "Payload/Test.app/"+name, 0600, bytes.NewReader(data)); e != nil {
				return e
			}
		}
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	output := filepath.Join(root, "prepared.ipa")
	job, e := PrepareSystemHelper(bundle, input, output)
	if e != nil {
		t.Fatal(e)
	}
	if !job.RequiresSigning || job.LicenseValidated || len(job.JobID) != 32 || job.Files != 4 {
		t.Fatalf("incorrect preparation %+v", job)
	}
	z, e := zip.OpenReader(output)
	if e != nil {
		t.Fatal(e)
	}
	defer z.Close()
	count := 0
	for _, f := range z.File {
		r, e := f.Open()
		if e != nil {
			t.Fatal(e)
		}
		data, e := io.ReadAll(r)
		r.Close()
		if e != nil {
			t.Fatal(e)
		}
		if strings.Contains(f.Name, "old-secret") {
			t.Fatal("previous target leaked")
		}
		if strings.HasSuffix(f.Name, ".sinf") {
			if string(data) != "original private license" {
				t.Fatal("license modified")
			}
			count++
		}
		if strings.HasSuffix(f.Name, "/Target/Main") && !bytes.Equal(data, a) {
			t.Fatal("ciphertext modified")
		}
	}
	if count != 2 {
		t.Fatal("manifest-directed license replication missing")
	}
	unpacked := filepath.Join(root, "unpacked")
	folder, e := UnpackSystemSigning(output, unpacked)
	if e != nil || folder.BundlePath != "Payload/Crumb.app" {
		t.Fatalf("unpack %+v: %v", folder, e)
	}
	if _, e := ArchiveSystemSigning(unpacked, output, filepath.Join(root, "preserved.ipa")); e != nil {
		t.Fatal(e)
	}
	systemFixtureWrite(t, filepath.Join(unpacked, folder.BundlePath), SystemPayload+"/Target/Main", []byte("modified target"))
	rejected := filepath.Join(root, "rejected.ipa")
	if _, e := ArchiveSystemSigning(unpacked, output, rejected); e == nil {
		t.Fatal("signer changed purchased executable")
	}
	if _, e := os.Stat(rejected); !os.IsNotExist(e) {
		t.Fatal("rejected signed archive remains")
	}
}
