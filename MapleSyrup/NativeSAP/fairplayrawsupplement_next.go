// Ahead-of-time block from exact system profile; synthetic differential tests
// cover registers, flags, indirect targets and memory. Not proof of enrollment.
// Exact daemon SHA256: 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784
// Translator revision: f7242cce58077a4af8db738ed295ee23ad8f067d
package rawsource

func rawSourceBootstrapBlock100a4be38(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100a4be38 {
		return 0x100a4be38, true, nil
	}
	// 0x100a4be38: add w17, w24, w10
	s.x[17] = uint64(uint32((uint64(uint32(s.x[24])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100a4be3c {
		return 0x100a4be3c, true, nil
	}
	// 0x100a4be3c: ldrb w17, [x13, w17, uxtw]
	value1, err := s.load((s.x[13] + uint64(uint32(uint64(uint32(s.x[17]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1))
	if runtime != nil && runtime.stopPC == 0x100a4be40 {
		return 0x100a4be40, true, nil
	}
	// 0x100a4be40: and w6, w2, #0xff
	s.x[6] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100a4be44 {
		return 0x100a4be44, true, nil
	}
	// 0x100a4be44: sub w2, w24, #1
	s.x[2] = uint64(uint32((uint64(uint32(s.x[24])) - uint64(0x1))))
	if runtime != nil && runtime.stopPC == 0x100a4be48 {
		return 0x100a4be48, true, nil
	}
	// 0x100a4be48: and w2, w2, #0xf
	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(0xf))))
	if runtime != nil && runtime.stopPC == 0x100a4be4c {
		return 0x100a4be4c, true, nil
	}
	// 0x100a4be4c: eor w2, w2, w11
	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100a4be50 {
		return 0x100a4be50, true, nil
	}
	// 0x100a4be50: mul w7, w2, w14
	s.x[7] = uint64(uint32((uint64(uint32(s.x[2])) * uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100a4be54 {
		return 0x100a4be54, true, nil
	}
	// 0x100a4be54: mul w2, w2, w15
	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) * uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a4be58 {
		return 0x100a4be58, true, nil
	}
	// 0x100a4be58: lsr w2, w2, #0x11
	s.x[2] = uint64(uint32(shiftRight(uint64(uint32(s.x[2])), int(uint64(0x11)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a4be5c {
		return 0x100a4be5c, true, nil
	}
	// 0x100a4be5c: msub w2, w2, w16, w7
	s.x[2] = uint64(uint32((uint64(uint32(s.x[7])) - uint64(uint32(s.x[2]))*uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a4be60 {
		return 0x100a4be60, true, nil
	}
	// 0x100a4be60: add x2, x23, w2, uxth
	s.x[2] = uint64((s.x[23] + uint64(uint16(uint64(uint32(s.x[2]))))))
	if runtime != nil && runtime.stopPC == 0x100a4be64 {
		return 0x100a4be64, true, nil
	}
	// 0x100a4be64: ldrb w2, [x2]
	value11, err := s.load(s.x[2], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[2] = uint64(uint32(value11))
	if runtime != nil && runtime.stopPC == 0x100a4be68 {
		return 0x100a4be68, true, nil
	}
	// 0x100a4be68: lsr w7, w2, #7
	s.x[7] = uint64(uint32(shiftRight(uint64(uint32(s.x[2])), int(uint64(0x7)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a4be6c {
		return 0x100a4be6c, true, nil
	}
	// 0x100a4be6c: orr w7, w7, w2, lsl #1
	s.x[7] = uint64(uint32((uint64(uint32(s.x[7])) | shiftLeft(uint64(uint32(s.x[2])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100a4be70 {
		return 0x100a4be70, true, nil
	}
	// 0x100a4be70: lsr w2, w2, #3
	s.x[2] = uint64(uint32(shiftRight(uint64(uint32(s.x[2])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a4be74 {
		return 0x100a4be74, true, nil
	}
	// 0x100a4be74: and w19, w2, #8
	s.x[19] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(0x8))))
	if runtime != nil && runtime.stopPC == 0x100a4be78 {
		return 0x100a4be78, true, nil
	}
	// 0x100a4be78: add w2, w6, #1
	s.x[2] = uint64(uint32((uint64(uint32(s.x[6])) + uint64(0x1))))
	if runtime != nil && runtime.stopPC == 0x100a4be7c {
		return 0x100a4be7c, true, nil
	}
	// 0x100a4be7c: orr w19, w19, #0xffffffc1
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) | uint64(0xffffffc1))))
	if runtime != nil && runtime.stopPC == 0x100a4be80 {
		return 0x100a4be80, true, nil
	}
	// 0x100a4be80: tst w2, #0xc0
	s.test32(uint32(uint64(uint32(s.x[2])) & uint64(0xc0)))
	if runtime != nil && runtime.stopPC == 0x100a4be84 {
		return 0x100a4be84, true, nil
	}
	// 0x100a4be84: cset w20, eq
	s.x[20] = uint64(uint32(boolValue(s.condition("eq"))))
	if runtime != nil && runtime.stopPC == 0x100a4be88 {
		return 0x100a4be88, true, nil
	}
	// 0x100a4be88: mul w21, w6, w0
	s.x[21] = uint64(uint32((uint64(uint32(s.x[6])) * uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100a4be8c {
		return 0x100a4be8c, true, nil
	}
	// 0x100a4be8c: mul w6, w6, w1
	s.x[6] = uint64(uint32((uint64(uint32(s.x[6])) * uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100a4be90 {
		return 0x100a4be90, true, nil
	}
	// 0x100a4be90: lsr w6, w6, #0x15
	s.x[6] = uint64(uint32(shiftRight(uint64(uint32(s.x[6])), int(uint64(0x15)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a4be94 {
		return 0x100a4be94, true, nil
	}
	// 0x100a4be94: msub w6, w6, w3, w21
	s.x[6] = uint64(uint32((uint64(uint32(s.x[21])) - uint64(uint32(s.x[6]))*uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100a4be98 {
		return 0x100a4be98, true, nil
	}
	// 0x100a4be98: eor w17, w17, w7
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x100a4be9c {
		return 0x100a4be9c, true, nil
	}
	// 0x100a4be9c: eor w7, w19, w4
	s.x[7] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x100a4bea0 {
		return 0x100a4bea0, true, nil
	}
	// 0x100a4bea0: eor w17, w17, w7
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x100a4bea4 {
		return 0x100a4bea4, true, nil
	}
	// 0x100a4bea4: add x6, x22, w6, uxth
	s.x[6] = uint64((s.x[22] + uint64(uint16(uint64(uint32(s.x[6]))))))
	if runtime != nil && runtime.stopPC == 0x100a4bea8 {
		return 0x100a4bea8, true, nil
	}
	// 0x100a4bea8: strb w17, [x6]
	if err := s.storeValue(s.x[6], uint64(uint32(s.x[17])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100a4beac {
		return 0x100a4beac, true, nil
	}
	// 0x100a4beac: mul w17, w20, w5
	s.x[17] = uint64(uint32((uint64(uint32(s.x[20])) * uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100a4beb0 {
		return 0x100a4beb0, true, nil
	}
	// 0x100a4beb0: eor w17, w17, w9
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100a4beb4 {
		return 0x100a4beb4, true, nil
	}
	// 0x100a4beb4: ldr x6, [x30, w17, sxtw #3]
	value31, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[17])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(value31)
	if runtime != nil && runtime.stopPC == 0x100a4beb8 {
		return 0x100a4beb8, true, nil
	}
	// 0x100a4beb8: add w24, w24, #1
	s.x[24] = uint64(uint32((uint64(uint32(s.x[24])) + uint64(0x1))))
	if runtime != nil && runtime.stopPC == 0x100a4bebc {
		return 0x100a4bebc, true, nil
	}
	// 0x100a4bebc: mov x17, #0x62bb
	s.x[17] = uint64(uint64(0x62bb))
	if runtime != nil && runtime.stopPC == 0x100a4bec0 {
		return 0x100a4bec0, true, nil
	}
	// 0x100a4bec0: braa x6, x17
	next = s.x[6]
	return next, false, nil
}
