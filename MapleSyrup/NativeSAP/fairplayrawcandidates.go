package rawsource

import (
	"bytes"
	"crypto/sha256"
)

// This bounds cryptographic work on distinct, naturally produced native models;
// it does not bound repeated visits to an interior instruction landmark. All
// private fingerprints are discarded with the execution and never reported.
const maxNativeCandidateChecks = 256

type nativeCandidateValidation struct {
	seen     map[[32]byte]struct{}
	accepted []byte
	checks   int
}

func (v *nativeCandidateValidation) observe(pack []byte, validate func([]byte) bool) (bool, error) {
	if len(v.accepted) != 0 {
		return true, nil
	}
	if len(pack) == 0 || validate == nil {
		return false, nil
	}
	digest := sha256.Sum256(pack)
	if _, found := v.seen[digest]; found {
		return false, nil
	}
	if v.checks >= maxNativeCandidateChecks {
		return false, &genericExecutionError{"native_candidate_limit"}
	}
	if v.seen == nil {
		v.seen = make(map[[32]byte]struct{})
	}
	v.seen[digest] = struct{}{}
	v.checks++
	if !validate(pack) {
		return false, nil
	}
	v.accepted = bytes.Clone(pack)
	return true, nil
}

func (v *nativeCandidateValidation) clear() {
	clear(v.accepted)
	clear(v.seen)
	v.accepted, v.seen = nil, nil
}
