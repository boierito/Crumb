package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadStateRejectsStaleOrForeignProducer(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0700); err != nil {
		t.Fatal(err)
	}
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

func TestAuthorizationExportSeparatesMissingInvalidAndDestinationFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0700); err != nil {
		t.Fatal(err)
	}
	hardware := []byte{2, 0, 0, 0, 0, 1}
	destination := filepath.Join(dir, "export.json")
	if err := exportAuthorizedState(dir, destination, 1, hardware); err == nil || err.Error() != "authorization_state_missing" {
		t.Fatalf("missing state: %v", err)
	}
	if err := saveAuthorizedState(dir, 1, hardware, map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(dir, "account-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := exportAuthorizedState(dir, destination, 2, hardware); err == nil || err.Error() != "authorization_cached_state_invalid" {
		t.Fatalf("foreign state: %v", err)
	}
	if err := exportAuthorizedState(dir, filepath.Join(dir, "missing", "export.json"), 1, hardware); err == nil || err.Error() != "authorization_state_export_failed" {
		t.Fatalf("destination failure: %v", err)
	}
	if err := exportAuthorizedState(dir, destination, 1, hardware); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(dir, "account-state.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("export changed internal state")
	}
}
