package rawsource

import "testing"

// This actual SINF-entry block is absent from the pinned bootstrap switch.
// It assigns x17 and branches through x9. No account bytes are involved.
const crossPhaseBlock = uint64(0x1005b0400)

func TestSharedBlockMatchesTranslatedOperations(t *testing.T) {
	s := &state{}
	s.x[9] = 0x1005b7438
	runtime := &rawSourceRuntime{}
	next, stopped, covered, err := resumeSourceBlock(s, runtime, crossPhaseBlock, s.x[9])
	if err != nil || !covered || stopped || next != s.x[9] || s.x[17] != 0xdaec {
		t.Fatalf("wrong translated effect: next=%x stopped=%v covered=%v err=%v", next, stopped, covered, err)
	}
	if runtime.stopPC != 0 || runtime.additionalBlocks != 1 || runtime.lastKnownBlock != crossPhaseBlock {
		t.Fatal("lost boundary configuration or dispatch evidence")
	}
	if knownSourceBlocks[crossPhaseBlock].label != "text+0x5b0400" {
		t.Fatal("checkpoint label must be a static source offset")
	}
}

func TestSharedBlockCannotSkipInteriorPhaseBoundary(t *testing.T) {
	s := &state{}
	s.x[9] = 0x1234
	runtime := &rawSourceRuntime{}
	next, stopped, covered, err := resumeSourceBlock(s, runtime, crossPhaseBlock, crossPhaseBlock+4)
	if err != nil || !covered || stopped || next != crossPhaseBlock+4 || s.x[17] != 0xdaec {
		t.Fatal("did not return the interior boundary to the phase loop")
	}
	if runtime.stopPC != 0 {
		t.Fatal("temporary boundary leaked into next phase")
	}
	// A second visit at the boundary must execute its instruction rather than
	// repeatedly stopping there; the selection phase has a two-visit contract.
	// This entry itself is not necessarily a block, so test the same rule at
	// an existing block entry instead of inventing an executable address.
	next, stopped, covered, err = resumeSourceBlock(s, runtime, crossPhaseBlock, crossPhaseBlock)
	if err != nil || !covered || stopped || next != 0x1234 {
		t.Fatal("boundary entry did not execute normally")
	}
}

func TestSharedBlockPreservesExplicitInstructionStop(t *testing.T) {
	s := &state{}
	s.x[9] = 0x1234
	runtime := &rawSourceRuntime{stopPC: crossPhaseBlock + 4}
	next, stopped, covered, err := resumeSourceBlock(s, runtime, crossPhaseBlock, 0x1005b7438)
	if err != nil || !covered || !stopped || next != crossPhaseBlock+4 || runtime.stopPC != crossPhaseBlock+4 {
		t.Fatal("overrode explicit instruction stop")
	}
}

func TestUntranslatedBlockIsNeverAccepted(t *testing.T) {
	s := &state{}
	s.x[9] = 0x1234
	runtime := &rawSourceRuntime{}
	next, stopped, covered, err := resumeSourceBlock(s, runtime, 0xdeadbeef, 0x1005b7438)
	if err == nil || covered || stopped || next != 0xdeadbeef || s.x[9] != 0x1234 || runtime.additionalBlocks != 0 {
		t.Fatal("untranslated target changed state or became success")
	}
}

func TestUnknownCodeTargetIsSanitized(t *testing.T) {
	r := &rawSourceRuntime{}
	r.recordUnknownPC(0x100b30000)
	if r.untranslatedTarget != "text+0xb30000" || r.targetClass != "in_pinned_code_window" {
		t.Fatal("failed to report pinned executable target")
	}
	for _, pointer := range []uint64{0, 0xdeadbeef, 0x200000000, 0x100000001} {
		r.recordUnknownPC(pointer)
		if r.untranslatedTarget != "" || r.targetClass != "outside_pinned_code_window" {
			t.Fatal("exposed an out-of-profile pointer")
		}
	}
}
