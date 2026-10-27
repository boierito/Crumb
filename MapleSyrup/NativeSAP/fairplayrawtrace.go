package rawsource

// Replaces only the generated per-instruction stop predicate. The instruction,
// register effects, memory operations and stop decision remain unchanged.
// Keep the trace writes out of thousands of generated call sites: inlining
// them multiplies compiler/race-instrumentation work across the whole image.
//
//go:noinline
func sourceInstructionStop(runtime *rawSourceRuntime, pc uint64) bool {
	if runtime == nil {
		return false
	}
	runtime.currentInstruction, runtime.executionContext = pc, "aot"
	return runtime.stopPC == pc
}
