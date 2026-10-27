package rawsource

type replayDiagnostic struct {
	Phase            string `json:"phase"`
	Failure          string `json:"failure"`
	FaultInstruction string `json:"fault_instruction,omitempty"`
	FaultRegion      string `json:"fault_region,omitempty"`
	RuntimeCompleted bool   `json:"runtime_completed"`
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
	if !ok || (initial.Failure != "raw_memory_unavailable" && initial.Failure != "raw_unknown_control_flow" && initial.Failure != "generic_execution_rejected" && initial.Failure != "native_candidate_limit") {
		return Result{}, err
	}
	replayed, replayErr := attempt(true)
	report := &replayDiagnostic{}
	if replayErr == nil {
		report.RuntimeCompleted = true
		replayed.Runtime.Replay = report
		return replayed, nil
	}
	clear(replayed.Source)
	clear(replayed.Frame)
	clear(replayed.Pack)
	if diagnostic, ok := RecoveryDetails(replayErr); ok {
		report.Phase, report.Failure = diagnostic.Phase, diagnostic.Failure
		report.FaultInstruction, report.FaultRegion = diagnostic.FaultInstruction, diagnostic.FaultRegion
	} else {
		report.Failure = "unclassified_replay_rejection"
	}
	initial.Replay = report
	return Result{}, &recoveryError{initial}
}
