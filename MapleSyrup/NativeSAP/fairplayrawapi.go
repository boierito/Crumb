// Copied into the isolated rawsource package; source inputs are modern native
// enrollment state, not the legacy StoreAgent capture.
package rawsource

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
)

var recoveryLock sync.Mutex

type Result struct {
	Source, Frame, Pack []byte
}

// RecoveryDiagnostic contains only fixed labels and aggregate lookup counts.
// No requested path, address, state, selector, or underlying error is retained.
type LookupCount struct {
	Attempts int `json:"attempts"`
	Missing  int `json:"missing"`
}
type RecoveryDiagnostic struct {
	Phase                 string                 `json:"phase"`
	Failure               string                 `json:"failure"`
	Lookups               map[string]LookupCount `json:"file_lookups,omitempty"`
	LastBlock             string                 `json:"last_translated_block,omitempty"`
	SharedBlockExecutions int                    `json:"shared_block_executions,omitempty"`
	UntranslatedTarget string                 `json:"untranslated_target,omitempty"`
	TargetClass         string                 `json:"target_class,omitempty"`
}
type recoveryError struct{ diagnostic RecoveryDiagnostic }

func (e *recoveryError) Error() string { return e.diagnostic.Failure }
func RecoveryDetails(err error) (RecoveryDiagnostic, bool) {
	var classified *recoveryError
	if !errors.As(err, &classified) {
		return RecoveryDiagnostic{}, false
	}
	result := classified.diagnostic
	result.Lookups = make(map[string]LookupCount, len(classified.diagnostic.Lookups))
	for key, count := range classified.diagnostic.Lookups {
		result.Lookups[key] = count
	}
	return result, true
}
func recoveryFailed(phase, code string, runtime *rawSourceRuntime) error {
	diagnostic := RecoveryDiagnostic{Phase: phase, Failure: code, Lookups: map[string]LookupCount{}}
	if runtime != nil {
		diagnostic.LastBlock = knownSourceBlocks[runtime.lastKnownBlock].label
		diagnostic.SharedBlockExecutions = runtime.additionalBlocks
		if code == "raw_unknown_control_flow" {
			diagnostic.UntranslatedTarget = runtime.untranslatedTarget
			diagnostic.TargetClass = runtime.targetClass
		}
		for key, count := range runtime.lookupCounts {
			diagnostic.Lookups[key] = count
		}
	}
	return &recoveryError{diagnostic}
}
func runtimeFailure(phase string, err error, runtime *rawSourceRuntime) error {
	code := "raw_operation_rejected"
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "outside the model") || strings.Contains(message, "unmapped"):
		code = "raw_memory_unavailable"
	case strings.Contains(message, "unsupported bootstrap external"):
		code = "raw_unsupported_external"
	case strings.Contains(message, "unknown block"):
		code = "raw_unknown_control_flow"
	case strings.Contains(message, "exceeded") || strings.Contains(message, "bounded exit"):
		code = "raw_step_limit"
	case strings.Contains(message, "unsupported"):
		code = "raw_unsupported_operation"
	}
	return recoveryFailed(phase, code, runtime)
}
func (r *rawSourceRuntime) recordLookup(requested string, found bool) {
	category := "other"
	lower := strings.ToLower(requested)
	switch {
	case strings.HasSuffix(lower, ".sinf"):
		category = "sinf"
	case strings.HasSuffix(lower, ".supf"):
		category = "supf"
	case strings.Contains(lower, "/sc info/") || strings.HasSuffix(lower, ".sidb") || strings.HasSuffix(lower, ".sidd"):
		category = "state"
	default:
		// The executable mapping is the only slash-prefixed basename mapping.
		for suffix := range r.files {
			if strings.HasPrefix(suffix, "/") && strings.Count(suffix, "/") == 1 && strings.HasSuffix(requested, suffix) {
				category = "executable"
				break
			}
		}
	}
	if r.lookupCounts == nil {
		r.lookupCounts = map[string]LookupCount{}
	}
	count := r.lookupCounts[category]
	// Saturate rather than retain an unbounded per-request trace.
	if count.Attempts < 1000000 {
		count.Attempts++
	}
	if !found && count.Missing < 1000000 {
		count.Missing++
	}
	r.lookupCounts[category] = count
}

// Preserve the upstream natural path and bounds while identifying its phases.
func runEnrollmentSource(value *model, runtime *rawSourceRuntime) error {
	if value.binaryHash != rawSourceBinaryHash {
		return recoveryFailed("resource", "raw_binary_mismatch", runtime)
	}
	worker := newState(value)
	worker.reset()
	lastRawSourceState = worker
	if _, _, err := runRawSourceAccountEntryCFG(worker, runtime, 1000000); err != nil {
		return runtimeFailure("account-parser", err, runtime)
	}
	if _, _, err := runRawSourceBootstrapCFG(worker, runtime, 1000000); err != nil {
		return runtimeFailure("bootstrap", err, runtime)
	}
	if _, err := runRawSourceKernelLoop(worker, "outer"); err != nil {
		return runtimeFailure("source", err, runtime)
	}
	if _, _, err := runRawSourceSelectionCFG(worker, runtime, 2000000); err != nil {
		return runtimeFailure("selection", err, runtime)
	}
	if _, _, err := runRawSourceSINFEntryCFG(worker, runtime, 2000000); err != nil {
		return runtimeFailure("sinf", err, runtime)
	}
	return nil
}

func VerifyResource() bool {
	sum := sha256.Sum256(rawSourceModel)
	return hex.EncodeToString(sum[:]) == "RAW_RESOURCE_SHA256"
}

// Recover constructs the source field, selection frame and selector-bound pack
// in process. The caller must validate the pack against the target IPA before
// treating it as usable. No material is logged or returned in diagnostics.
func Recover(executable string, hardware []byte, files map[string][]byte) (Result, error) {
	var result Result
	if len(hardware) != 6 || executable == "" || len(files) == 0 || len(files) > 256 {
		return result, recoveryFailed("inputs", "raw_invalid_inputs", nil)
	}
	recoveryLock.Lock()
	defer recoveryLock.Unlock()
	if !VerifyResource() {
		return result, recoveryFailed("resource", "raw_resource_invalid", nil)
	}
	value, err := loadModelBytes(rawSourceModel)
	if err != nil {
		return result, recoveryFailed("resource", "raw_resource_invalid", nil)
	}
	defer func() {
		if activeRawSourceRuntime != nil {
			for _, data := range activeRawSourceRuntime.files {
				clear(data)
			}
			clear(activeRawSourceRuntime.mac)
			for _, item := range activeRawSourceRuntime.fakeCFValues {
				clear(item.data)
			}
		}
		if lastRawSourceState != nil {
			for _, data := range lastRawSourceState.writable {
				clear(data)
			}
			clear(lastRawSourceState.x[:])
			clear(lastRawSourceState.vectors[:])
		}
		for _, data := range value.pages {
			clear(data)
		}
		activeRawSourceRuntime, lastRawSourceState = nil, nil
	}()
	if err := configureRawSourceRuntime(value, "enrollment-sinf", nil, "", executable, 0); err != nil {
		return result, recoveryFailed("configure", "raw_setup_failed", activeRawSourceRuntime)
	}
	activeRawSourceRuntime.mac = append([]byte(nil), hardware...)
	total := 0
	for name, data := range files {
		total += len(data)
		if name == "" || len(data) > 64<<20 || total > 128<<20 {
			return result, recoveryFailed("inputs", "raw_input_size_limit", activeRawSourceRuntime)
		}
		activeRawSourceRuntime.files[name] = append([]byte(nil), data...)
	}
	if err := runEnrollmentSource(value, activeRawSourceRuntime); err != nil {
		return result, err
	}
	result.Pack, err = currentRawSourceEnrollmentPack()
	if err != nil {
		return Result{}, recoveryFailed("pack", "raw_pack_unavailable", activeRawSourceRuntime)
	}
	result.Source, err = currentRawSourceBytes(224)
	if err != nil {
		clear(result.Pack)
		return Result{}, recoveryFailed("source-output", "raw_source_unavailable", activeRawSourceRuntime)
	}
	result.Frame, err = currentRawSourceFrameBytes()
	if err != nil {
		clear(result.Pack)
		clear(result.Source)
		return Result{}, recoveryFailed("frame-output", "raw_frame_unavailable", activeRawSourceRuntime)
	}
	return result, nil
}
