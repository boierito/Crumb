package rawsource

import "testing"

func TestNativeRegionFollowsRealSourceBranchInsteadOfRecordedTrace(t *testing.T) {
	// Actual 26A428 instructions: add x28,sp,#6,lsl#12; mov x17,#0xc0c7;
	// braa x24,x17. The old outer kernel unconditionally continued at
	// 0x1005f3bdc. Use a different real, translated destination and require
	// execution of its mov x17,#0xdaec; br x9 before the region can finish.
	for _, genericOnly := range []bool{false, true} {
		s := genericSyntheticState()
		s.x[24], s.x[9] = 0x1005b0400, 0x1006b10a0
		runtime := &rawSourceRuntime{genericOnly: genericOnly}
		region := nativeRegion{"source", 0x1005b7a84, 0x1006b10a0, 1, 8}
		pc, err := runNativeRegion(s, runtime, region.start, region)
		if err != nil || pc != region.end || s.x[17] != 0xdaec ||
			s.x[28] != s.sp+(6<<12) || runtime.stopPC != 0 {
			t.Fatalf("native branch destination was not followed: generic=%v err=%v", genericOnly, err)
		}
	}
}

func TestNativeRegionPreservesNaturalBoundaryVisitsAndState(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{
		genericNormal(compactAddImm64, 0, 0, 0, 1), genericBranch(compactBlockRegister, 5, 0),
		genericNormal(compactAddImm64, 1, 1, 0, 1), genericBranch(compactBlockRegister, 6, 0),
	})
	s := genericSyntheticState()
	s.x[0], s.x[1], s.x[5], s.x[6] = 0, 0, base+8, base
	runtime := &rawSourceRuntime{}
	region := nativeRegion{"selection", base, base + 8, 2, 4}
	pc, err := runNativeRegion(s, runtime, base, region)
	if err != nil || pc != region.end || s.x[0] != 2 || s.x[1] != 1 || runtime.stopPC != 0 {
		t.Fatal("stage reset state, skipped a visit or rejected its final-budget arrival", err)
	}
}

func TestNativeRegionCallsAndReturnsPreserveLinkRegister(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{
		genericBranch(compactBlockCallRegister, 30, 0),
		genericBranch(compactBlockRegister, 8, 0),
		genericNormal(compactAddImm64, 7, 7, 0, 1),
		genericBranch(compactBlockReturn, 30, 0),
	})
	s := genericSyntheticState()
	s.x[7], s.x[8], s.x[30] = 0, base+16, base+8
	runtime := &rawSourceRuntime{}
	region := nativeRegion{"source", base, base + 16, 1, 4}
	pc, err := runNativeRegion(s, runtime, base, region)
	if err != nil || pc != region.end || s.x[7] != 1 || s.x[30] != base+4 {
		t.Fatal("call/return corrupted architectural state", err)
	}
}

func TestNativeRegionDoesNotTreatTraceFragmentStopAsStageCompletion(t *testing.T) {
	previous, blocks := genericImage, knownSourceBlocks
	defer func() { genericImage, knownSourceBlocks = previous, blocks }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{
		genericNormal(compactAddImm64, 0, 0, 0, 1),
		genericNormal(compactAddImm64, 1, 1, 0, 1), genericBranch(compactBlockRegister, 5, 0),
	})
	knownSourceBlocks = map[uint64]sourceBlock{base: {run: func(s *state, _ *rawSourceRuntime) (uint64, bool, error) {
		s.x[0]++
		return base + 4, true, nil
	}}}
	s := genericSyntheticState()
	s.x[0], s.x[1], s.x[5] = 0, 0, base+12
	region := nativeRegion{"source", base, base + 12, 1, 3}
	pc, err := runNativeRegion(s, &rawSourceRuntime{}, base, region)
	if err != nil || pc != region.end || s.x[0] != 1 || s.x[1] != 1 {
		t.Fatal("AOT fragment stop was mistaken for native region exit", err)
	}
}

func TestNativeRegionContinuesRealSINFFragmentWithoutRestartingFrame(t *testing.T) {
	const start, fragmentEnd, end = uint64(0x100234cb8), uint64(0x100234cdc), uint64(0x100234ce0)
	m := &model{sp: 0x700000, pages: map[uint64][]byte{0x700000: make([]byte, pageSize), 0x600000: make([]byte, pageSize)}}
	m.pages[0x600000][0] = 0x73
	bias := uint64(0x5ca7911f)
	m.registers[28] = uint64(0x600000) - bias
	if _, exists := knownSourceBlocks[start]; exists {
		s := newState(m)
		s.reset()
		next, stopped, covered, err := resumeSourceBlock(s, &rawSourceRuntime{}, start, nativeFrameLandmark)
		if err != nil || !covered || !stopped || next != fragmentEnd {
			t.Fatal("real cross-phase fragment no longer demonstrates its local stop", err)
		}
	}
	for _, genericOnly := range []bool{false, true} {
		s := newState(m)
		s.reset()
		region := nativeRegion{"selection", start, end, 1, 4}
		next, err := runNativeRegion(s, &rawSourceRuntime{genericOnly: genericOnly}, start, region)
		if err != nil || next != end || s.x[8] != 0x73 {
			t.Fatal("native continuation after the real SINF fragment was skipped", err)
		}
	}
}

func TestNativeRegionRejectsCandidateAndContinuesAtItsActualPC(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{
		genericNormal(compactAddImm64, 0, 0, 0, 1), genericBranch(compactBlockRegister, 5, 0),
		genericNormal(compactAddImm64, 1, 1, 0, 1), genericBranch(compactBlockRegister, 6, 0),
	})
	s := genericSyntheticState()
	s.x[0], s.x[1], s.x[5], s.x[6] = 0, 0, base+8, base
	checks := 0
	runtime := &rawSourceRuntime{boundary: func() (bool, error) { checks++; return checks == 2, nil }}
	region := nativeRegion{"enrollment", base, base + 8, 1, 8}
	pc, err := runNativeRegion(s, runtime, base, region)
	if err != nil || pc != region.end || checks != 2 || s.x[0] != 2 || s.x[1] != 1 {
		t.Fatal("candidate rejection restarted/skipped native execution", err)
	}
}

func TestNativeRegionRejectsUnknownTargetAndBoundsInfiniteLoop(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{genericBranch(compactBlockRegister, 5, 0)})
	for _, target := range []uint64{0x600000000000, base} {
		s := genericSyntheticState()
		s.x[5] = target
		runtime := &rawSourceRuntime{}
		region := nativeRegion{"source", base, base + 4, 1, 3}
		if _, err := runNativeRegion(s, runtime, base, region); err == nil || runtime.stopPC != 0 {
			t.Fatal("unknown target or infinite loop accepted")
		}
		if s.model.pages[0x600000000000] != nil {
			t.Fatal("unknown target created fabricated memory")
		}
	}
}
