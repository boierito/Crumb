package rawsource

// Replaces only the generated per-instruction stop predicate. The instruction,
// register effects, memory operations and stop decision remain unchanged.
func sourceInstructionStop(runtime *rawSourceRuntime, pc uint64) bool {
	if runtime == nil {
		return false
	}
	runtime.currentInstruction, runtime.executionContext = pc, "aot"
	return runtime.stopPC == pc
}
