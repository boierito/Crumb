package machine

import (
	"bytes"
	"errors"
	"testing"
)

func TestDownloadLicenseImportsBothBagsAndReturnsFinalSnapshot(t *testing.T) {
	key, diversity := []byte("key"), []byte("diversity")
	intermediate := []byte("intermediate")
	calls := []string{}
	result, err := consumeLicenseBags(key, diversity,
		func(b []byte) (int32, map[string][]byte, map[string]any, error) {
			if !bytes.Equal(b, key) {
				t.Fatal("keybag substituted")
			}
			calls = append(calls, "key")
			return 0, map[string][]byte{"state": intermediate}, nil, nil
		},
		func(b []byte) (int32, map[string][]byte, map[string]any, error) {
			if !bytes.Equal(b, diversity) {
				t.Fatal("diversity substituted")
			}
			calls = append(calls, "diversity")
			return 0, map[string][]byte{"state": []byte("final")}, nil, nil
		})
	if err != nil || result.Status != 0 || result.KeybagStatus == nil || *result.KeybagStatus != 0 || result.DiversityStatus == nil || *result.DiversityStatus != 0 || !result.DiversityPresent || string(result.Files["state"]) != "final" || len(calls) != 2 || calls[0] != "key" || calls[1] != "diversity" {
		t.Fatal("incomplete/incorrect license import")
	}
	if !bytes.Equal(intermediate, make([]byte, len(intermediate))) {
		t.Fatal("intermediate snapshot retained")
	}
	if string(key) != "key" || string(diversity) != "diversity" {
		t.Fatal("input modified")
	}
}

func TestDownloadLicenseStopsOnFirstRejection(t *testing.T) {
	for _, native := range []bool{false, true} {
		expected := errors.New("emulator failed")
		first := func(b []byte) (int32, map[string][]byte, map[string]any, error) {
			if native {
				return -7, nil, nil, nil
			}
			return 0, nil, nil, expected
		}
		result, err := consumeLicenseBags([]byte{1}, []byte{2}, first, func([]byte) (int32, map[string][]byte, map[string]any, error) {
			t.Fatal("continued rejected import")
			return 0, nil, nil, nil
		})
		if native {
			if result.Status != -7 || result.KeybagStatus == nil || *result.KeybagStatus != -7 || err != nil {
				t.Fatal("native rejection lost")
			}
		} else if !errors.Is(err, expected) || result.KeybagStatus != nil {
			t.Fatal("incomplete execution misclassified")
		}
		if result.DiversityStatus != nil {
			t.Fatal("unexecuted diversity status manufactured")
		}
	}
}

func TestDownloadLicenseSecondFailureCannotBecomeSuccess(t *testing.T) {
	for _, native := range []bool{false, true} {
		expected := errors.New("emulator failed")
		result, err := consumeLicenseBags([]byte{1}, []byte{2}, func([]byte) (int32, map[string][]byte, map[string]any, error) {
			return 0, map[string][]byte{"state": []byte("old")}, nil, nil
		}, func([]byte) (int32, map[string][]byte, map[string]any, error) {
			if native {
				return -8, map[string][]byte{"state": []byte("rejected")}, nil, nil
			}
			return 0, nil, nil, expected
		})
		if result.KeybagStatus == nil || *result.KeybagStatus != 0 {
			t.Fatal("second status overwrote keybag status")
		}
		if native {
			if result.Status != -8 || result.DiversityStatus == nil || *result.DiversityStatus != -8 || err != nil || string(result.Files["state"]) != "rejected" {
				t.Fatal("second native rejection lost")
			}
		} else if !errors.Is(err, expected) || result.DiversityStatus != nil || result.Files != nil {
			t.Fatal("partial import reported as complete")
		}
	}
}

func TestDownloadLicenseWithoutDiversityKeepsOriginalPath(t *testing.T) {
	result, err := consumeLicenseBags([]byte{1}, nil, func([]byte) (int32, map[string][]byte, map[string]any, error) {
		return 0, map[string][]byte{"state": []byte("key-only")}, nil, nil
	}, func([]byte) (int32, map[string][]byte, map[string]any, error) {
		t.Fatal("invented diversity import")
		return 0, nil, nil, nil
	})
	if err != nil || result.Status != 0 || result.DiversityPresent || result.DiversityStatus != nil || string(result.Files["state"]) != "key-only" {
		t.Fatal("legacy path changed")
	}
}
