// Ahead-of-time continuation observed after the 23042 bootstrap byte loop.
// Exact binary provenance and independent ARM64 synthetic tests are required.
// Exact daemon SHA256: 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784
// Translator revision: f7242cce58077a4af8db738ed295ee23ad8f067d
package rawsource

func rawSourceBootstrapBlock1008c292c(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x1008c292c {
		return 0x1008c292c, true, nil
	}
	// 0x1008c292c: str w24, [sp, #0x11c]
	if err := s.storeValue(addOffset(s.sp, 284), uint64(uint32(s.x[24])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1008c2930 {
		return 0x1008c2930, true, nil
	}
	// 0x1008c2930: mov w9, #0x2aa3
	s.x[9] = uint64(uint32(uint64(0x2aa3)))
	if runtime != nil && runtime.stopPC == 0x1008c2934 {
		return 0x1008c2934, true, nil
	}
	// 0x1008c2934: add w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1008c2938 {
		return 0x1008c2938, true, nil
	}
	// 0x1008c2938: ldr x9, [x30, w9, sxtw #3]
	value3, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value3)
	if runtime != nil && runtime.stopPC == 0x1008c293c {
		return 0x1008c293c, true, nil
	}
	// 0x1008c293c: mov x21, x25
	s.x[21] = uint64(s.x[25])
	if runtime != nil && runtime.stopPC == 0x1008c2940 {
		return 0x1008c2940, true, nil
	}
	// 0x1008c2940: mov x25, x27
	s.x[25] = uint64(s.x[27])
	if runtime != nil && runtime.stopPC == 0x1008c2944 {
		return 0x1008c2944, true, nil
	}
	// 0x1008c2944: mov x17, #0x1a88
	s.x[17] = uint64(uint64(0x1a88))
	if runtime != nil && runtime.stopPC == 0x1008c2948 {
		return 0x1008c2948, true, nil
	}
	// 0x1008c2948: braa x9, x17
	next = s.x[9]
	return next, false, nil
}
