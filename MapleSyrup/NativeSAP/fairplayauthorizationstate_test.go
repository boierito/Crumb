package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/majd/ipatool/v2/internal/fairplay"
)

func TestAuthorizedStateExportPreservesOriginalAndAccountBinding(t *testing.T) {
	dir := t.TempDir()
	os.Chmod(dir, 0700)
	hardware := []byte{2, 0, 0, 0, 0, 1}
	if err := saveAuthorizedState(dir, 1, hardware, map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}, "SC Info.sidd": {1, 2, 3}}); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(dir, "account-state.json"))
	output := filepath.Join(dir, "PRIVATE-export.json")
	if err := exportAuthorizedState(dir, output, 1, hardware); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(output)
	if err != nil || info.Mode().Perm()&0077 != 0 {
		t.Fatal("export is not private")
	}
	state, err := readAuthorizedStateFile(output, 1, hardware)
	if err != nil {
		t.Fatal(err)
	}
	state.clear()
	if state, err := readAuthorizedStateFile(output, 2, hardware); err == nil {
		state.clear()
		t.Fatal("another account accepted")
	}
	if err := exportAuthorizedState(dir, output, 1, hardware); err == nil {
		t.Fatal("existing export overwritten")
	}
	wrong := filepath.Join(dir, "wrong-account.json")
	if err := exportAuthorizedState(dir, wrong, 2, hardware); err == nil {
		t.Fatal("another account exported")
	}
	if _, err := os.Lstat(wrong); !os.IsNotExist(err) {
		t.Fatal("failed export left a file")
	}
	after, _ := os.ReadFile(filepath.Join(dir, "account-state.json"))
	if !bytes.Equal(before, after) {
		t.Fatal("export modified saved state")
	}
}

func TestAuthorizationStateIsAccountBoundAndAtomic(t *testing.T) {
	dir := t.TempDir()
	os.Chmod(dir, 0700)
	hardware := []byte{2, 0, 0, 0, 0, 1}
	files := map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}, "lock": {}}
	if err := saveAuthorizedState(dir, 1, hardware, files); err != nil {
		t.Fatal(err)
	}
	saved, err := readAuthorizedState(dir, 1, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer saved.clear()
	if !bytes.Equal(saved.Files["SC Info.sidb"], files["SC Info.sidb"]) {
		t.Fatal("state changed")
	}
	if s, err := readAuthorizedState(dir, 2, hardware); err == nil {
		s.clear()
		t.Fatal("another account state accepted")
	}
	other := bytes.Clone(hardware)
	other[5] ^= 1
	if s, err := readAuthorizedState(dir, 1, other); err == nil {
		s.clear()
		t.Fatal("another identity accepted")
	}
	before, _ := os.ReadFile(filepath.Join(dir, "account-state.json"))
	if err := saveAuthorizedState(dir, 1, hardware, map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}, "../outside": {1}}); err == nil {
		t.Fatal("unsafe state accepted")
	}
	after, _ := os.ReadFile(filepath.Join(dir, "account-state.json"))
	if !bytes.Equal(before, after) {
		t.Fatal("failed save replaced state")
	}
	os.Chmod(filepath.Join(dir, "account-state.json"), 0644)
	if s, err := readAuthorizedState(dir, 1, hardware); err == nil {
		s.clear()
		t.Fatal("public state permissions accepted")
	}
}
func TestAuthorizationBindsAllLicenseDomains(t *testing.T) {
	a := fairplay.PortableEnrollmentInput{SINFSHA256: "a"}
	b := fairplay.PortableEnrollmentInput{SINFSHA256: "b"}
	for _, test := range []struct {
		sinfs   map[string]bool
		domains []fairplay.PortableEnrollmentInput
		want    bool
	}{
		{map[string]bool{"a": true, "b": true}, []fairplay.PortableEnrollmentInput{a, b}, true},
		{map[string]bool{"a": true}, []fairplay.PortableEnrollmentInput{a, a}, true},
		{map[string]bool{"a": true}, []fairplay.PortableEnrollmentInput{a, b}, false},
		{map[string]bool{"a": true, "b": true}, []fairplay.PortableEnrollmentInput{a}, false},
		{map[string]bool{}, nil, false},
	} {
		if capturedSINFMatchesAllDomains(test.sinfs, test.domains) != test.want {
			t.Fatal("license domain set mismatch")
		}
	}
}

func TestDownloadStateRejectsStaleOrForeignProducer(t *testing.T) {
	dir := t.TempDir()
	hardware := []byte{2, 0, 0, 0, 0, 1}
	first := map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}, "SC Info.sidd": {1}}
	next := map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}, "SC Info.sidd": {2}}
	if err := saveDownloadState(dir, 1, hardware, nil, first); err != nil {
		t.Fatal(err)
	}
	original, err := readAuthorizedState(dir, 1, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer original.clear()
	if err := saveDownloadState(dir, 1, hardware, original.Digests, next); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(dir, "account-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, request := range []struct {
		dsid     uint64
		hardware []byte
		expected map[string]string
	}{
		{1, hardware, original.Digests},                 // stale generation
		{1, hardware, nil},                              // account state appeared after a cold generation
		{2, hardware, original.Digests},                 // wrong account
		{1, []byte{2, 0, 0, 0, 0, 2}, original.Digests}, // wrong virtual identity
	} {
		if err := saveDownloadState(dir, request.dsid, request.hardware, request.expected, first); err == nil {
			t.Fatal("stale or foreign download overwrote account state")
		}
		after, err := os.ReadFile(filepath.Join(dir, "account-state.json"))
		if err != nil || !bytes.Equal(before, after) {
			t.Fatal("rejected update changed account state")
		}
	}
	if err := os.Remove(filepath.Join(dir, "account-state.json")); err != nil {
		t.Fatal(err)
	}
	if err := saveDownloadState(dir, 1, hardware, original.Digests, next); err == nil {
		t.Fatal("missing original state accepted")
	}
}
