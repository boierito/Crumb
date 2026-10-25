package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/state"
)

func TestLegacyCandidatePreservesWholeStateAndIdentity(t *testing.T) {
	hardware := []byte{1, 2, 3, 4, 5, 6}
	sidb := bytes.Repeat([]byte{17}, 883)
	context, _ := json.Marshal(map[string]any{"schema": "crumb-fairplay-context-v2", "producer": "retained-storeagent", "hardware_id": hardware})
	consumer, _ := json.Marshal(map[string]any{"consumer": "storeagent-10.9-import-keybag", "emulator_completed": true, "import_succeeded": true, "native_status": 0, "state_files": map[string][]byte{"SC Info.sidb": sidb, "SC Info.sidb.lock": {}}})
	input, err := decodeLegacyEnrollment(context, consumer)
	if err != nil {
		t.Fatal(err)
	}
	defer input.clear()
	if !bytes.Equal(input.hardware, hardware) || !bytes.Equal(input.files["SC Info.sidb"], sidb) {
		t.Fatal("identity or opaque SIDB changed")
	}
	if _, err := decodeLegacyEnrollment([]byte(`{"schema":"crumb-fairplay-context-v2","producer":"retained-storeagent","hardware_id":"AQID"}`), consumer); err == nil {
		t.Fatal("accepted wrong hardware length")
	}
	if _, err := readLegacyEnrollmentInput(captureV3Fixture(t, map[string]any{"native_status": 0, "import_succeeded": true})); err == nil {
		t.Fatal("accepted capture with no retained hardware identity")
	}
}

func TestRecoveredFrameGateChecksLaterSlotsAndRejectsMismatch(t *testing.T) {
	source := make([]byte, 224)
	for i := range source {
		source[i] = byte(i)
	}
	transformed, err := state.Evaluate(source, nil)
	if err != nil {
		t.Fatal(err)
	}
	selected, err := frame.Evaluate(transformed[8*16:9*16], nil)
	if err != nil {
		t.Fatal(err)
	}
	value := rawsource.Result{Source: source, Frame: selected}
	if err := verifyRecoveredEnrollment(value); err != nil {
		t.Fatal(err)
	}
	value.Frame = []byte{1, 2, 3}
	if err := verifyRecoveredEnrollment(value); err == nil {
		t.Fatal("accepted nonmatching frame")
	}
}
