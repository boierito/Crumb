package rawsource_test

import (
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
}
