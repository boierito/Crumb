package rawsource

import (
	"bytes"
	"testing"

	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	statetransform "github.com/majd/ipatool/v2/wafflebridge/recovery/state"
)

// This uses the actual embedded execution model and the real independent
// state/frame evaluators, with synthetic bytes only. No enrollment is used.
func TestEnrollmentCheckpointPreservesSelectionPairAfterSINFMutation(t *testing.T) {
	recoveryLock.Lock()
	defer recoveryLock.Unlock()

	value, err := loadModelBytes(rawSourceModel)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		for _, page := range value.pages {
			clear(page)
		}
	}()
	// These are public output locations of the hash-pinned 26A428 test fixture.
	// Production recovery derives them through currentRawSource* getters.
	const sourceAddress = uint64(0x7000000fb304)
	const frameAddress = uint64(0x7000000fb724)
	worker := newState(value)
	worker.reset()
	lastRawSourceState = worker
	defer func() {
		lastRawSourceState = nil
		for _, page := range worker.writable {
			clear(page)
		}
	}()
	selectionSource := make([]byte, 224)
	for i := range selectionSource {
		selectionSource[i] = byte(i)
	}
	transformed, err := statetransform.Evaluate(selectionSource, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(transformed)
	selectionFrame, err := frame.Evaluate(transformed[8*16:9*16], nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(selectionFrame)

	if err := worker.store(sourceAddress, selectionSource); err != nil {
		t.Fatal(err)
	}
	if err := worker.store(frameAddress, selectionFrame); err != nil {
		t.Fatal(err)
	}
	checkpoint, err := checkpointEnrollmentOutputs()
	if err != nil {
		t.Fatal(err)
	}
	defer checkpoint.clear()
	if !bytes.Equal(checkpoint.source, selectionSource) || !bytes.Equal(checkpoint.frame, selectionFrame) {
		t.Fatal("selection checkpoint returned different bytes")
	}
	if !verifyEnrollmentCheckpoint(checkpoint) {
		t.Fatal("native-flow observer rejected an independently constructed valid checkpoint")
	}
	invalid := enrollmentOutputCheckpoint{source: bytes.Clone(checkpoint.source), frame: bytes.Clone(checkpoint.frame)}
	invalid.frame[0] ^= 1
	if verifyEnrollmentCheckpoint(invalid) {
		t.Fatal("native-flow observer accepted an invalid checkpoint")
	}
	invalid.clear()

	// Model an SINF phase that overwrites both guest output locations. The
	// checkpoint must retain the earlier, independently valid source/frame.
	if err := worker.store(sourceAddress, bytes.Repeat([]byte{0xA5}, 224)); err != nil {
		t.Fatal(err)
	}
	if err := worker.store(frameAddress, bytes.Repeat([]byte{0x5A}, 16)); err != nil {
		t.Fatal(err)
	}
	runtime := &rawSourceRuntime{}
	observePostSINFOutputs(runtime, checkpoint)
	if !runtime.checkpointCompared || !runtime.sourceChangedAfterSINF || !runtime.frameChangedAfterSINF {
		t.Fatal("post-SINF mutation was not detected without exposing bytes")
	}
	if !bytes.Equal(checkpoint.source, selectionSource) || !bytes.Equal(checkpoint.frame, selectionFrame) {
		t.Fatal("SINF phase mutated the stored selection checkpoint")
	}

	postTransform, err := statetransform.Evaluate(checkpoint.source, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(postTransform)
	matches := 0
	for offset := 0; offset+16 <= len(postTransform); offset += 16 {
		candidate, err := frame.Evaluate(postTransform[offset:offset+16], nil)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(candidate, checkpoint.frame) {
			matches++
		}
		clear(candidate)
	}
	if matches == 0 {
		t.Fatal("the real evaluators rejected the preserved checkpoint")
	}
}

func TestPostSINFComparisonUnavailableIsNotReportedAsStable(t *testing.T) {
	recoveryLock.Lock()
	defer recoveryLock.Unlock()
	previous := lastRawSourceState
	lastRawSourceState = nil
	defer func() { lastRawSourceState = previous }()
	runtime := &rawSourceRuntime{}
	observePostSINFOutputs(runtime, enrollmentOutputCheckpoint{})
	if runtime.checkpointCompared || runtime.sourceChangedAfterSINF || runtime.frameChangedAfterSINF {
		t.Fatal("an unavailable memory comparison was reported as valid")
	}
}

// Real file routing is private. Assert diagnostic counters retain only
// public file-type labels and keep the aggregate state counts unchanged.
func TestStateFileLookupsClassifyRequiredVsOptionalSources(t *testing.T) {
	runtime := &rawSourceRuntime{files: map[string][]byte{}}
	for _, item := range []struct {
		requested string
		found     bool
	}{
		{"/var/db/fpsd/SC Info/SC Info.sidb", true},
		{"/Users/Shared/SC Info/SC Info.sidb", true},
		{"/var/db/fpsd/SC Info/SC Info.sidd", false},
		{"/var/db/fpsd/SC Info/SC Info.sidf", false},
		{"/var/db/fpsd/SC Info/PRIVATE-CANARY.state", false},
	} {
		runtime.recordLookup(item.requested, item.found)
	}
	diag := runtimeDiagnostic("completed", "", runtime)
	if diag.Lookups["state"] != (LookupCount{Attempts: 5, Missing: 3}) ||
		diag.StateFileTypes["sidb"] != (LookupCount{Attempts: 2, Missing: 0}) ||
		diag.StateFileTypes["sidd"] != (LookupCount{Attempts: 1, Missing: 1}) ||
		diag.StateFileTypes["sidf"] != (LookupCount{Attempts: 1, Missing: 1}) ||
		diag.StateFileTypes["other"] != (LookupCount{Attempts: 1, Missing: 1}) ||
		len(diag.StateFileTypes) != 4 {
		t.Fatalf("lost state lookup provenance: %+v", diag)
	}
	// This snapshot must not accidentally retain underlying paths.
	original := &recoveryError{diagnostic: diag}
	copy, ok := RecoveryDetails(original)
	if !ok {
		t.Fatal("safe recovery snapshot unavailable")
	}
	delete(copy.StateFileTypes, "sidb")
	if _, ok := original.diagnostic.StateFileTypes["sidb"]; !ok {
		t.Fatal("mutated original diagnostic through its exported copy")
	}
}

func TestCheckpointSelectsOnlyAnIndependentlyValidSourceWidth(t *testing.T) {
	for _, width := range []int{192, 224} {
		source := make([]byte, 224)
		for i := range source {
			source[i] = byte(i*17 + 3)
		}
		transformed, err := statetransform.Evaluate(source[:width], nil)
		if err != nil {
			t.Fatal(err)
		}
		selection, err := frame.Evaluate(transformed[8*16:9*16], nil)
		clear(transformed)
		if err != nil {
			t.Fatal(err)
		}
		checkpoint := enrollmentOutputCheckpoint{source: source, frame: selection}
		if !normalizeEnrollmentCheckpoint(&checkpoint) || len(checkpoint.source) != width || !verifyEnrollmentCheckpoint(checkpoint) {
			t.Fatalf("valid %d-byte source rejected or changed", width)
		}
		checkpoint.clear()
	}
	invalid := enrollmentOutputCheckpoint{source: make([]byte, 224), frame: make([]byte, 15)}
	if normalizeEnrollmentCheckpoint(&invalid) {
		t.Fatal("invalid frame accepted")
	}
	invalid.clear()
}
