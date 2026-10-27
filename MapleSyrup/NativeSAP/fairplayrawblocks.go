// Copied into the isolated rawsource package by prepare-fairplay-raw.py.
package rawsource

import "fmt"

// A PC is only exportable as a label when it is an aligned offset inside the
// publicly pinned executable's code window. Never log arbitrary raw guest
// register values, stack addresses, pointers or private state.
func (r *rawSourceRuntime) recordUnknownPC(pc uint64) {
	if r == nil {
		return
	}
	const base uint64 = 0x100000000
	const window uint64 = 16 << 20
	if pc >= base && pc-base < window && pc%4 == 0 {
		r.untranslatedTarget = fmt.Sprintf("text+0x%x", pc-base)
		r.targetClass = "in_pinned_code_window"
	} else {
		r.untranslatedTarget = ""
		r.targetClass = "outside_pinned_code_window"
	}
}

type sourceBlockFunc func(*state, *rawSourceRuntime) (uint64, bool, error)
type sourceBlock struct {
	run   sourceBlockFunc
	label string
}

func callSourceBlock(s *state, runtime *rawSourceRuntime, pc uint64, block sourceBlockFunc) (uint64, bool, error) {
	if runtime != nil {
		// pc comes from a generated switch case, never a private register.
		runtime.lastKnownBlock = pc
	}
	return block(s, runtime)
}

func sourceBlockLabel(pc uint64) string {
	if block, ok := knownSourceBlocks[pc]; ok {
		return block.label
	}
	if genericImage.contains(pc) && pc >= 0x100000000 && pc-0x100000000 < 16<<20 {
		return fmt.Sprintf("text+0x%x", pc-0x100000000)
	}
	return ""
}

// Supply only an already translated entry. Unknown targets remain failures.
// Returning to the outer loop preserves its block budget and boundary hit count.
func resumeSourceBlock(s *state, runtime *rawSourceRuntime, pc, phaseEnd uint64) (uint64, bool, bool, error) {
	block, covered := knownSourceBlocks[pc]
	if runtime != nil && runtime.genericOnly {
		covered = false
	}
	if !covered && runtime != nil {
		if genericImage == nil {
			return pc, false, false, &genericExecutionError{"generic_profile_missing"}
		}
		if !genericImage.contains(pc) {
			runtime.untranslatedTarget = ""
			runtime.targetClass = "outside_generic_profile"
			return pc, false, false, &genericExecutionError{"generic_pc_outside_profile"}
		}
		block = sourceBlock{run: func(s *state, r *rawSourceRuntime) (uint64, bool, error) { return runGenericInstructionBlock(s, r, pc) }, label: sourceBlockLabel(pc)}
		covered = true
	}
	if !covered || runtime == nil {
		if runtime != nil {
			runtime.recordUnknownPC(pc)
		}
		return pc, false, false, nil
	}
	previousStop := runtime.stopPC
	// A reused block can span a boundary that split another phase's trace.
	// Catch it at the existing per-instruction checks. When entering at the
	// boundary itself, execute normally; selection counts two natural visits.
	injected := previousStop == 0 && pc != phaseEnd
	if injected {
		runtime.stopPC = phaseEnd
	}
	if runtime.additionalBlocks < 1000000 {
		runtime.additionalBlocks++
	}
	next, stopped, err := callSourceBlock(s, runtime, pc, block.run)
	runtime.stopPC = previousStop
	if injected && stopped && next == phaseEnd {
		// Let the original phase loop count this visit and decide completion.
		stopped = false
	}
	return next, stopped, true, err
}
