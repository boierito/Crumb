package rawsource

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// These exercise scheduling and ownership, not FairPlay decryption. The
// actual production callback must still resolve the owner's original SINF.
func TestNativeCandidateValidationIndependentOfCheckpointAndRepeatVisits(t *testing.T) {
	v := nativeCandidateValidation{}
	defer v.clear()
	calls := 0
	validate := func(pack []byte) bool { calls++; return bytes.Equal(pack, []byte{2}) }
	for i := 0; i < 4096; i++ {
		if ok, err := v.observe([]byte{1}, validate); ok || err != nil {
			t.Fatal("repeated rejected candidate accepted or spent visit budget")
		}
	}
	if calls != 1 || v.checks != 1 {
		t.Fatal("duplicate model revalidated")
	}
	pack := []byte{2}
	ok, err := v.observe(pack, validate)
	if !ok || err != nil || v.checks != 2 {
		t.Fatal("original SINF callback not reached independently of checkpoint")
	}
	pack[0] = 3
	if !bytes.Equal(v.accepted, []byte{2}) {
		t.Fatal("accepted model aliases guest scratch")
	}
	retained := v.accepted
	v.clear()
	if retained[0] != 0 || len(v.seen) != 0 {
		t.Fatal("private candidate state retained after cleanup")
	}
}

func TestNativeCandidateValidationBoundsDistinctCryptoWork(t *testing.T) {
	v := nativeCandidateValidation{}
	defer v.clear()
	for i := 0; i <= maxNativeCandidateChecks; i++ {
		pack := make([]byte, 4)
		binary.LittleEndian.PutUint32(pack, uint32(i))
		ok, err := v.observe(pack, func([]byte) bool { return false })
		if ok || (err != nil) != (i == maxNativeCandidateChecks) {
			t.Fatal("candidate budget did not fail closed")
		}
	}
	if v.checks != maxNativeCandidateChecks {
		t.Fatal("invalid candidate count")
	}
}
