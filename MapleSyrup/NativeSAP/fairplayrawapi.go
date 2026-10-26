// Copied into the isolated rawsource package; source inputs are modern native
// enrollment state, not the legacy StoreAgent capture.
package rawsource

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var recoveryLock sync.Mutex

type Result struct {
	Source, Frame, Pack []byte
	Runtime             RecoveryDiagnostic `json:"-"`
}

// RecoveryDiagnostic contains only fixed labels and aggregate lookup counts.
// No requested path, address, state, selector, or underlying error is retained.
type LookupCount struct {
	Attempts int `json:"attempts"`
	Missing  int `json:"missing"`
}
type RecoveryDiagnostic struct {
	Phase                  string                 `json:"phase"`
	UnmodeledExternalSlot  string                 `json:"unmodeled_external_slot,omitempty"`
	Failure                string                 `json:"failure"`
	Lookups                map[string]LookupCount `json:"file_lookups,omitempty"`
	StateFileTypes         map[string]LookupCount `json:"state_file_lookups,omitempty"`
	LastBlock              string                 `json:"last_translated_block,omitempty"`
	SharedBlockExecutions  int                    `json:"shared_block_executions,omitempty"`
	UntranslatedTarget     string                 `json:"untranslated_target,omitempty"`
	TargetClass            string                 `json:"target_class,omitempty"`
	ExecutionBackend       string                 `json:"execution_backend,omitempty"`
	GenericInstructions    uint64                 `json:"generic_instructions,omitempty"`
	GenericLastInstruction string                 `json:"generic_last_instruction,omitempty"`
	GenericBlocks          int                    `json:"generic_blocks,omitempty"`
	CheckpointCompared      bool                   `json:"checkpoint_compared_after_sinf"`
	SourceChangedAfterSINF  bool                   `json:"source_changed_after_sinf"`
	FrameChangedAfterSINF   bool                   `json:"frame_changed_after_sinf"`
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
	result.StateFileTypes = make(map[string]LookupCount, len(classified.diagnostic.StateFileTypes))
	for key, count := range classified.diagnostic.StateFileTypes {
		result.StateFileTypes[key] = count
	}
	return result, true
}
func recoveryFailed(phase, code string, runtime *rawSourceRuntime) error {
	return &recoveryError{runtimeDiagnostic(phase, code, runtime)}
}
func runtimeDiagnostic(phase, code string, runtime *rawSourceRuntime) RecoveryDiagnostic {
	diagnostic := RecoveryDiagnostic{Phase: phase, Failure: code, Lookups: map[string]LookupCount{}}
	if runtime != nil {
		diagnostic.LastBlock = sourceBlockLabel(runtime.lastKnownBlock)
		diagnostic.SharedBlockExecutions = runtime.additionalBlocks
		diagnostic.ExecutionBackend = "aot"
		if runtime.genericBlocks > 0 {
			diagnostic.ExecutionBackend = "generic-profile"
		}
		diagnostic.GenericInstructions = runtime.genericInstructions
		diagnostic.GenericLastInstruction = sourceBlockLabel(runtime.genericLastPC)
		diagnostic.GenericBlocks = runtime.genericBlocks
		diagnostic.CheckpointCompared = runtime.checkpointCompared
		diagnostic.SourceChangedAfterSINF = runtime.sourceChangedAfterSINF
		diagnostic.FrameChangedAfterSINF = runtime.frameChangedAfterSINF
		if code == "raw_unsupported_external" {
			diagnostic.UnmodeledExternalSlot = runtime.unmodeledExternalSlot
		}
		if code == "raw_unknown_control_flow" || strings.HasPrefix(code, "generic_") {
			diagnostic.UntranslatedTarget = runtime.untranslatedTarget
			diagnostic.TargetClass = runtime.targetClass
		}
		for key, count := range runtime.lookupCounts {
			diagnostic.Lookups[key] = count
		}
		if len(runtime.stateFileLookupCounts) > 0 {
			diagnostic.StateFileTypes = make(map[string]LookupCount, len(runtime.stateFileLookupCounts))
			for key, count := range runtime.stateFileLookupCounts {
				diagnostic.StateFileTypes[key] = count
			}
		}
	}
	return diagnostic
}
func runtimeFailure(phase string, err error, runtime *rawSourceRuntime) error {
	var generic *genericExecutionError
	if errors.As(err, &generic) {
		return recoveryFailed(phase, generic.code, runtime)
	}
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
	if category != "state" {
		return
	}
	// Names are discarded immediately; persist only this finite public
	// extension taxonomy, never account-specific file names or paths.
	kind := "other"
	switch {
	case strings.HasSuffix(lower, ".sidb"):
		kind = "sidb"
	case strings.HasSuffix(lower, ".sidd"):
		kind = "sidd"
	case strings.HasSuffix(lower, ".sidf"):
		kind = "sidf"
	}
	if r.stateFileLookupCounts == nil {
		r.stateFileLookupCounts = map[string]LookupCount{}
	}
	detail := r.stateFileLookupCounts[kind]
	if detail.Attempts < 1000000 {
		detail.Attempts++
	}
	if !found && detail.Missing < 1000000 {
		detail.Missing++
	}
	r.stateFileLookupCounts[kind] = detail
}

// The selection CFG materializes the frame during its second frame call.
// Capture source/frame at this phase boundary, before SINF changes guest state.
// This does not substitute any frame or relax downstream cryptographic checks.
type enrollmentOutputCheckpoint struct {
	source, frame []byte
}

func (c *enrollmentOutputCheckpoint) clear() {
	clear(c.source)
	clear(c.frame)
	c.source, c.frame = nil, nil
}

func checkpointEnrollmentOutputs() (enrollmentOutputCheckpoint, error) {
	var checkpoint enrollmentOutputCheckpoint
	var err error
	checkpoint.source, err = currentRawSourceBytes(224)
	if err != nil {
		return enrollmentOutputCheckpoint{}, err
	}
	checkpoint.frame, err = currentRawSourceFrameBytes()
	if err != nil {
		checkpoint.clear()
		return enrollmentOutputCheckpoint{}, err
	}
	return checkpoint, nil
}

// Comparisons are local only. Unavailable comparison stays explicitly false;
// no secret bytes, secret digests, guest addresses or paths leave the runtime.
func observePostSINFOutputs(runtime *rawSourceRuntime, checkpoint enrollmentOutputCheckpoint) {
	if runtime == nil {
		return
	}
	postSource, sourceErr := currentRawSourceBytes(224)
	postFrame, frameErr := currentRawSourceFrameBytes()
	defer clear(postSource)
	defer clear(postFrame)
	if sourceErr != nil || frameErr != nil {
		return
	}
	runtime.checkpointCompared = true
	runtime.sourceChangedAfterSINF = !bytes.Equal(checkpoint.source, postSource)
	runtime.frameChangedAfterSINF = !bytes.Equal(checkpoint.frame, postFrame)
}

// Preserve the upstream natural path and bounds while identifying its phases.
func runEnrollmentSource(value *model, runtime *rawSourceRuntime) (enrollmentOutputCheckpoint, error) {
	if value.binaryHash != rawSourceBinaryHash {
		return enrollmentOutputCheckpoint{}, recoveryFailed("resource", "raw_binary_mismatch", runtime)
	}
	worker := newState(value)
	worker.reset()
	lastRawSourceState = worker
	if _, _, err := runRawSourceAccountEntryCFG(worker, runtime, 1000000); err != nil {
		return enrollmentOutputCheckpoint{}, runtimeFailure("account-parser", err, runtime)
	}
	if _, _, err := runRawSourceBootstrapCFG(worker, runtime, 1000000); err != nil {
		return enrollmentOutputCheckpoint{}, runtimeFailure("bootstrap", err, runtime)
	}
	if _, err := runRawSourceKernelLoop(worker, "outer"); err != nil {
		return enrollmentOutputCheckpoint{}, runtimeFailure("source", err, runtime)
	}
	if _, _, err := runRawSourceSelectionCFG(worker, runtime, 2000000); err != nil {
		return enrollmentOutputCheckpoint{}, runtimeFailure("selection", err, runtime)
	}
	checkpoint, err := checkpointEnrollmentOutputs()
	if err != nil {
		return enrollmentOutputCheckpoint{}, recoveryFailed("selection-output", "raw_selection_output_unavailable", runtime)
	}
	if _, _, err := runRawSourceSINFEntryCFG(worker, runtime, 2000000); err != nil {
		checkpoint.clear()
		return enrollmentOutputCheckpoint{}, runtimeFailure("sinf", err, runtime)
	}
	observePostSINFOutputs(runtime, checkpoint)
	return checkpoint, nil
}

func VerifyResource() bool {
	sum := sha256.Sum256(rawSourceModel)
	return hex.EncodeToString(sum[:]) == "RAW_RESOURCE_SHA256" && genericImage != nil
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
	checkpoint, err := runEnrollmentSource(value, activeRawSourceRuntime)
	if err != nil {
		return result, err
	}
	result.Pack, err = currentRawSourceEnrollmentPack()
	if err != nil {
		checkpoint.clear()
		return Result{}, recoveryFailed("pack", "raw_pack_unavailable", activeRawSourceRuntime)
	}
	result.Source, result.Frame = checkpoint.source, checkpoint.frame
	result.Runtime = runtimeDiagnostic("completed", "", activeRawSourceRuntime)
	return result, nil
}

// sourceExternalImportSlot classifies the *public, profile-defined* 26A428
// Mach-O external stub namespace, not an arbitrary guest address. The high
// bits and alignment are checked before deriving the bounded slot index.
// A synthetic/invalid guest branch can never become a raw pointer in logs.
func sourceExternalImportSlot(pc uint64) string {
    const base uint64 = 0xf100000000000000
    if pc < base || pc >= base+(1<<24) || (pc-base)%0x1000 != 0 {
        return "unclassified_external"
    }
    return fmt.Sprintf("external-slot-%03x", (pc-base)/0x1000)
}
