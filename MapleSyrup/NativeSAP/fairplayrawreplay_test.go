package rawsource

import "testing"

func TestRealRecoveryRejectionClearsAttemptStateWithoutChangingInputs(t *testing.T) {
	hardware := []byte{1, 2, 3, 4, 5, 6}
	input := []byte("PRIVATE-INVALID-EXECUTABLE")
	_, err := Recover("/Applications/App.app/App", hardware, map[string][]byte{"/App": input})
	if err == nil || activeRawSourceRuntime != nil || lastRawSourceState != nil {
		t.Fatal("failed recovery retained its active model or runtime")
	}
	if string(input) != "PRIVATE-INVALID-EXECUTABLE" || hardware[0] != 1 || hardware[5] != 6 || !VerifyResource() {
		t.Fatal("attempt cleanup destroyed caller inputs or shared resources")
	}
}

func TestReplayIsFreshBoundedAndReturnsOnlyCandidateMaterial(t *testing.T) {
	calls := 0
	result, err := recoverWithReplay(func(genericOnly bool) (Result, error) {
		calls++
		if !genericOnly {
			return Result{}, recoveryFailed("sinf", "raw_memory_unavailable", nil)
		}
		return Result{Pack: []byte("candidate only")}, nil
	})
	if err != nil || calls != 2 || string(result.Pack) != "candidate only" ||
		result.Runtime.Replay == nil || !result.Runtime.Replay.RuntimeCompleted {
		t.Fatal("clean alternate execution was not returned as a candidate")
	}
}

func TestReplayRejectionClearsPartialOutputsAndKeepsIndependentErrors(t *testing.T) {
	private := []byte("PRIVATE-REPLAY-PARTIAL-OUTPUT")
	calls := 0
	result, err := recoverWithReplay(func(genericOnly bool) (Result, error) {
		calls++
		if !genericOnly {
			return Result{}, recoveryFailed("sinf", "raw_memory_unavailable", nil)
		}
		return Result{Source: private}, recoveryFailed("source", "generic_instruction_limit", nil)
	})
	if err == nil || calls != 2 || len(result.Source)+len(result.Frame)+len(result.Pack) != 0 {
		t.Fatal("failed replay promoted partial material or retried repeatedly")
	}
	for _, b := range private {
		if b != 0 {
			t.Fatal("partial output retained after failed replay")
		}
	}
	diagnostic, ok := RecoveryDetails(err)
	if !ok || diagnostic.Phase != "sinf" || diagnostic.Replay == nil ||
		diagnostic.Replay.Phase != "source" || diagnostic.Replay.Failure != "generic_instruction_limit" || diagnostic.Replay.RuntimeCompleted {
		t.Fatal("replay rejection was treated as equivalent execution")
	}
	diagnostic.Replay.Phase = "modified"
	copy, _ := RecoveryDetails(err)
	if copy.Replay.Phase != "source" {
		t.Fatal("public diagnostic mutated the retained replay")
	}
}

func TestReplayDoesNotRetryInvalidInputsResourcesOrUnsupportedInstructions(t *testing.T) {
	for _, code := range []string{"raw_invalid_inputs", "raw_resource_invalid", "generic_unsupported_instruction", "raw_unsupported_external"} {
		calls := 0
		_, err := recoverWithReplay(func(bool) (Result, error) {
			calls++
			return Result{}, recoveryFailed("test", code, nil)
		})
		if err == nil || calls != 1 {
			t.Fatal("deterministic unsupported operation was retried")
		}
	}
}

func TestReplayIncludesNativeAndInstructionBudgetRejections(t *testing.T) {
	for _, code := range []string{"native_candidate_limit", "native_region_block_limit", "generic_instruction_limit"} {
		calls := 0
		result, err := recoverWithReplay(func(genericOnly bool) (Result, error) {
			calls++
			if !genericOnly {
				return Result{}, recoveryFailed("sinf", code, nil)
			}
			return Result{}, &recoveryError{RecoveryDiagnostic{Phase: "sinf", Failure: "generic_instruction_limit", CandidateChecks: 3, CandidateSINFVerified: true}}
		})
		diagnostic, ok := RecoveryDetails(err)
		if !ok || calls != 2 || diagnostic.Replay == nil || diagnostic.Replay.RuntimeCompleted || diagnostic.Replay.CandidateChecks != 3 || !diagnostic.Replay.CandidateSINFVerified || len(result.Pack) != 0 {
			t.Fatal("budget rejection did not perform one independent, fail-closed replay")
		}
	}
}
