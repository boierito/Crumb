package rawsource_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
	"testing"
)

func TestRecoverRejectsMissingInputs(t *testing.T) {
	if !rawsource.VerifyResource() {
		t.Fatal("invalid embedded raw resource")
	}
	for _, input := range []struct {
		executable string
		hardware   []byte
		files      map[string][]byte
	}{
		{"/Applications/App.app/App", make([]byte, 6), nil},
		{"/Applications/App.app/App", make([]byte, 5), map[string][]byte{"/App": {1}}},
		{"", make([]byte, 6), map[string][]byte{"/App": {1}}},
	} {
		result, err := rawsource.Recover(input.executable, input.hardware, input.files)
		if err == nil || len(result.Pack) != 0 || len(result.Frame) != 0 || len(result.Source) != 0 {
			t.Fatal("missing state promoted to enrollment")
		}
	}
}

func TestRecoverFailureDoesNotMutateEmbeddedResource(t *testing.T) {
	result, err := rawsource.Recover("/Applications/App.app/App", make([]byte, 6), map[string][]byte{"/App": {1, 2, 3}})
	if err == nil || len(result.Pack) != 0 {
		t.Fatal("invalid native state produced model")
	}
	if !rawsource.VerifyResource() {
		t.Fatal("private run mutated packaged resource")
	}
	diagnostic, ok := rawsource.RecoveryDetails(fmt.Errorf("PRIVATE-WRAPPER: %w", err))
	if !ok || diagnostic.Phase == "" || diagnostic.Failure == "" {
		t.Fatal("runtime failure lacks a classified phase")
	}
}

func TestRecoveryDiagnosticDoesNotRetainPrivateInput(t *testing.T) {
	const canary = "PRIVATE-HARDWARE-AND-PATH-CANARY"
	_, err := rawsource.Recover("/Applications/"+canary+".app/"+canary, make([]byte, 6), map[string][]byte{
		"/" + canary:                         []byte(canary),
		"/Users/Shared/SC Info/SC Info.sidb": []byte(canary),
	})
	diagnostic, ok := rawsource.RecoveryDetails(err)
	if !ok {
		t.Fatal("missing safe diagnostic")
	}
	data, marshalErr := json.Marshal(diagnostic)
	if marshalErr != nil || bytes.Contains(data, []byte(canary)) || bytes.Contains([]byte(err.Error()), []byte(canary)) {
		t.Fatal("private data in runtime diagnostic")
	}
	for key, count := range diagnostic.Lookups {
		switch key {
		case "sinf", "supf", "state", "executable", "other":
		default:
			t.Fatal("unbounded lookup label")
		}
		if count.Attempts < count.Missing || count.Attempts > 1000000 {
			t.Fatal("unbounded lookup count")
		}
	}
	for key := range diagnostic.Lookups {
		delete(diagnostic.Lookups, key)
	}
	copy, ok := rawsource.RecoveryDetails(err)
	if !ok {
		t.Fatal("lost diagnostic after cleanup")
	}
	data, _ = json.Marshal(copy)
	if bytes.Contains(data, []byte(canary)) {
		t.Fatal("private data after cleanup")
	}
	if !rawsource.VerifyResource() {
		t.Fatal("runtime changed embedded profile")
	}
}
