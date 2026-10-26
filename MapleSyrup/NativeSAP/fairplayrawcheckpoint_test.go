package rawsource

import (
	"bytes"
	"testing"

	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/state"
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
	profile, err := executionProfileForModel(value)
	if err != nil {
		t.Fatal(err)
	}
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
	transformed, err := state.Evaluate(selectionSource, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(transformed)
	selectionFrame, err := frame.Evaluate(transformed[8*16:9*16], nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(selectionFrame)

	if err := worker.store(profile.sourceAddress, selectionSource); err != nil {
		t.Fatal(err)
	}
	if err := worker.store(profile.frameAddress, selectionFrame); err != nil {
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

	// Model an SINF phase that overwrites both guest output locations. The
	// checkpoint must retain the earlier, independently valid source/frame.
	if err := worker.store(profile.sourceAddress, bytes.Repeat([]byte{0xA5}, 224)); err != nil {
		t.Fatal(err)
	}
	if err := worker.store(profile.frameAddress, bytes.Repeat([]byte{0x5A}, 16)); err != nil {
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

	postTransform, err := state.Evaluate(checkpoint.source, nil)
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
