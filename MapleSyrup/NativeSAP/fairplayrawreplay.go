package rawsource

type replayDiagnostic struct {
	Phase                 string              `json:"phase"`
	Failure               string              `json:"failure"`
	FaultInstruction      string              `json:"fault_instruction,omitempty"`
	FaultRegion           string              `json:"fault_region,omitempty"`
	CandidateChecks       int                 `json:"candidate_sinf_checks,omitempty"`
	CandidateSINFVerified bool                `json:"candidate_sinf_verified"`
	RuntimeCompleted      bool                `json:"runtime_completed"`
	ByteLoop              *byteLoopDiagnostic `json:"native_byte_loop,omitempty"`
}

// A second engine always starts from fresh model/files/CPU state. It never
// retries account authorization and never resumes a partially failed block.
// A returned model is still only a candidate; original SINF/IPA gates apply.
func recoverWithReplay(attempt func(bool) (Result, error)) (Result, error) {
	result, err := attempt(false)
	if err == nil {
		return result, nil
	}
	clear(result.Source)
	clear(result.Frame)
	clear(result.Pack)
	initial, ok := RecoveryDetails(err)
	if !ok || (initial.Failure != "raw_memory_unavailable" && initial.Failure != "raw_unknown_control_flow" && initial.Failure != "generic_execution_rejected" && initial.Failure != "native_candidate_limit" && initial.Failure != "native_region_block_limit" && initial.Failure != "generic_instruction_limit") {
		return Result{}, err
	}
	replayed, replayErr := attempt(true)
	report := &replayDiagnostic{}
	if replayErr == nil {
		report.CandidateChecks, report.CandidateSINFVerified = replayed.Runtime.CandidateChecks, replayed.Runtime.CandidateSINFVerified
		report.RuntimeCompleted = true
		replayed.Runtime.Replay = report
		return replayed, nil
	}
	clear(replayed.Source)
	clear(replayed.Frame)
	clear(replayed.Pack)
	if diagnostic, ok := RecoveryDetails(replayErr); ok {
		report.Phase, report.Failure = diagnostic.Phase, diagnostic.Failure
		report.CandidateChecks, report.CandidateSINFVerified = diagnostic.CandidateChecks, diagnostic.CandidateSINFVerified
		report.FaultInstruction, report.FaultRegion = diagnostic.FaultInstruction, diagnostic.FaultRegion
		report.ByteLoop = diagnostic.ByteLoop
	} else {
		report.Failure = "unclassified_replay_rejection"
	}
	initial.Replay = report
	return Result{}, &recoveryError{initial}
}
