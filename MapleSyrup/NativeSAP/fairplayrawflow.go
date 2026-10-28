package rawsource

// Region boundaries describe the same public native landmarks as the pinned
// recovery profile. Registers and memory are carried across every boundary;
// no stage is entered by resetting registers or applying a captured trace.
type nativeRegion struct {
	name       string
	start, end uint64
	visits     int
	maxBlocks  int
}

const nativeEnrollmentStart = uint64(0x100cd1c18)
const nativeEnrollmentEnd = uint64(0x100234cdc)
const nativeFrameLandmark = uint64(0x10008a644)
const maxNativeRecoveryBlocks = 8_000_000

// The pinned caller snapshot has a null saved LR. A return through its exact
// epilogue is a native result, not a missing instruction at address zero.
// It never establishes SINF validity or authorizes plaintext export.
type nativeReturnError struct{ status int32 }

func (*nativeReturnError) Error() string { return "native entry returned without a validated license" }

func nativeSnapshotReturn(s *state, runtime *rawSourceRuntime, pc uint64) error {
	if pc == 0 && runtime.outputScratchEnabled && runtime.outputScratchReady == 15 &&
		runtime.currentInstruction == 0x1000e05e0 && s.sp == 0x700000100000 && s.x[30] == 0 {
		return &nativeReturnError{status: int32(s.x[0])}
	}
	return nil
}

func observeNativePC(runtime *rawSourceRuntime, pc uint64) {
	if runtime == nil || runtime.observe == nil {
		return
	}
	switch pc {
	case nativeEnrollmentStart, 0x10057ae74, 0x1005b7438, 0x1006b10a0, nativeFrameLandmark:
		runtime.observe(pc)
	}
}

// Follow actual native destinations, including calls, returns and indirect
// branches. An AOT trace fragment's stop is not a completed stage. Only natural
// arrival at this region's own boundary can complete it.
func runNativeRegion(s *state, runtime *rawSourceRuntime, pc uint64, region nativeRegion) (uint64, error) {
	if s == nil || runtime == nil || pc != region.start || region.visits < 1 ||
		region.maxBlocks < 1 || region.maxBlocks > maxNativeRecoveryBlocks || runtime.stopPC != 0 {
		return pc, &genericExecutionError{"native_region_invalid"}
	}
	hits := 0
	for count := 0; count < region.maxBlocks; count++ {
		if err := nativeSnapshotReturn(s, runtime, pc); err != nil {
			return pc, err
		}
		if pc == region.end {
			hits++
			if hits == region.visits {
				if runtime.boundary == nil {
					return pc, nil
				}
				accepted, err := runtime.boundary()
				if err != nil {
					return pc, err
				}
				if accepted {
					return pc, nil
				}
				hits = 0
			}
		}
		if pc >= 0xf000000000000000 || pc == 0x720000000000 {
			next, err := runtime.handle(s, pc)
			if err != nil {
				return pc, err
			}
			pc = next
			continue
		}
		next, _, covered, err := resumeSourceBlock(s, runtime, pc, region.end)
		if err != nil {
			return pc, err
		}
		if !covered {
			return pc, &genericExecutionError{"native_region_uncovered"}
		}
		pc = next
	}
	if pc == region.end && hits+1 == region.visits {
		if runtime.boundary == nil {
			return pc, nil
		}
		accepted, err := runtime.boundary()
		if err != nil {
			return pc, err
		}
		if accepted {
			return pc, nil
		}
	}
	return pc, &genericExecutionError{"native_region_block_limit"}
}
