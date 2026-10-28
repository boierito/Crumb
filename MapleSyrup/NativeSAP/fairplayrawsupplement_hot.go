// Ahead-of-time translations of public pinned recovery blocks.
// Source binary SHA-256: 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784
package rawsource

func rawSourceBootstrapBlock1009a1f04(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x1009a1f04 {
		return 0x1009a1f04, true, nil
	}
	// 0x1009a1f04: mov x17, x25
	s.x[17] = uint64(s.x[25])
	if runtime != nil && runtime.stopPC == 0x1009a1f08 {
		return 0x1009a1f08, true, nil
	}
	// 0x1009a1f08: mov x25, x14
	s.x[25] = uint64(s.x[14])
	if runtime != nil && runtime.stopPC == 0x1009a1f0c {
		return 0x1009a1f0c, true, nil
	}
	// 0x1009a1f0c: mov x14, x23
	s.x[14] = uint64(s.x[23])
	if runtime != nil && runtime.stopPC == 0x1009a1f10 {
		return 0x1009a1f10, true, nil
	}
	// 0x1009a1f10: mov x3, x21
	s.x[3] = uint64(s.x[21])
	if runtime != nil && runtime.stopPC == 0x1009a1f14 {
		return 0x1009a1f14, true, nil
	}
	// 0x1009a1f14: mov w9, #0x5464
	s.x[9] = uint64(uint32(uint64(0x5464)))
	if runtime != nil && runtime.stopPC == 0x1009a1f18 {
		return 0x1009a1f18, true, nil
	}
	// 0x1009a1f18: movk w9, #0xba98, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xba98, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f1c {
		return 0x1009a1f1c, true, nil
	}
	// 0x1009a1f1c: add w9, w7, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[7])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f20 {
		return 0x1009a1f20, true, nil
	}
	// 0x1009a1f20: mov w10, #0x9ffd
	s.x[10] = uint64(uint32(uint64(0x9ffd)))
	if runtime != nil && runtime.stopPC == 0x1009a1f24 {
		return 0x1009a1f24, true, nil
	}
	// 0x1009a1f24: movk w10, #0x4567, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x4567, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f28 {
		return 0x1009a1f28, true, nil
	}
	// 0x1009a1f28: and w22, w9, w10
	s.x[22] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f2c {
		return 0x1009a1f2c, true, nil
	}
	// 0x1009a1f2c: mov w9, #0xb23c
	s.x[9] = uint64(uint32(uint64(0xb23c)))
	if runtime != nil && runtime.stopPC == 0x1009a1f30 {
		return 0x1009a1f30, true, nil
	}
	// 0x1009a1f30: movk w9, #0x4805, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x4805, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f34 {
		return 0x1009a1f34, true, nil
	}
	// 0x1009a1f34: add w9, w24, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[24])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f38 {
		return 0x1009a1f38, true, nil
	}
	// 0x1009a1f38: mov w10, #0x971
	s.x[10] = uint64(uint32(uint64(0x971)))
	if runtime != nil && runtime.stopPC == 0x1009a1f3c {
		return 0x1009a1f3c, true, nil
	}
	// 0x1009a1f3c: movk w10, #0x80ff, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x80ff, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f40 {
		return 0x1009a1f40, true, nil
	}
	// 0x1009a1f40: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f44 {
		return 0x1009a1f44, true, nil
	}
	// 0x1009a1f44: mov w15, #0x86fa
	s.x[15] = uint64(uint32(uint64(0x86fa)))
	if runtime != nil && runtime.stopPC == 0x1009a1f48 {
		return 0x1009a1f48, true, nil
	}
	// 0x1009a1f48: movk w15, #0x1382, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x1382, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f4c {
		return 0x1009a1f4c, true, nil
	}
	// 0x1009a1f4c: eor w21, w22, w15
	s.x[21] = uint64(uint32((uint64(uint32(s.x[22])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f50 {
		return 0x1009a1f50, true, nil
	}
	// 0x1009a1f50: mov w15, #0xef60
	s.x[15] = uint64(uint32(uint64(0xef60)))
	if runtime != nil && runtime.stopPC == 0x1009a1f54 {
		return 0x1009a1f54, true, nil
	}
	// 0x1009a1f54: movk w15, #0x30a4, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x30a4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f58 {
		return 0x1009a1f58, true, nil
	}
	// 0x1009a1f58: eor w23, w9, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f5c {
		return 0x1009a1f5c, true, nil
	}
	// 0x1009a1f5c: mov w15, #0xe2da
	s.x[15] = uint64(uint32(uint64(0xe2da)))
	if runtime != nil && runtime.stopPC == 0x1009a1f60 {
		return 0x1009a1f60, true, nil
	}
	// 0x1009a1f60: movk w15, #0x336, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x336, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f64 {
		return 0x1009a1f64, true, nil
	}
	// 0x1009a1f64: eor w24, w9, w15
	s.x[24] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f68 {
		return 0x1009a1f68, true, nil
	}
	// 0x1009a1f68: mov w15, #0xe632
	s.x[15] = uint64(uint32(uint64(0xe632)))
	if runtime != nil && runtime.stopPC == 0x1009a1f6c {
		return 0x1009a1f6c, true, nil
	}
	// 0x1009a1f6c: movk w15, #0x904b, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x904b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f70 {
		return 0x1009a1f70, true, nil
	}
	// 0x1009a1f70: add w24, w24, w15
	s.x[24] = uint64(uint32((uint64(uint32(s.x[24])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f74 {
		return 0x1009a1f74, true, nil
	}
	// 0x1009a1f74: mov w15, #0xffdf
	s.x[15] = uint64(uint32(uint64(0xffdf)))
	if runtime != nil && runtime.stopPC == 0x1009a1f78 {
		return 0x1009a1f78, true, nil
	}
	// 0x1009a1f78: movk w15, #0xdfef, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xdfef, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f7c {
		return 0x1009a1f7c, true, nil
	}
	// 0x1009a1f7c: eor w26, w9, w15
	s.x[26] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f80 {
		return 0x1009a1f80, true, nil
	}
	// 0x1009a1f80: mov w15, #0xeb8c
	s.x[15] = uint64(uint32(uint64(0xeb8c)))
	if runtime != nil && runtime.stopPC == 0x1009a1f84 {
		return 0x1009a1f84, true, nil
	}
	// 0x1009a1f84: movk w15, #0xa3d9, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa3d9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1f88 {
		return 0x1009a1f88, true, nil
	}
	// 0x1009a1f88: add w23, w23, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[23])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f8c {
		return 0x1009a1f8c, true, nil
	}
	// 0x1009a1f8c: eor w9, w9, w23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[23])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f90 {
		return 0x1009a1f90, true, nil
	}
	// 0x1009a1f90: add w10, w10, w21
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f94 {
		return 0x1009a1f94, true, nil
	}
	// 0x1009a1f94: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a1f98 {
		return 0x1009a1f98, true, nil
	}
	// 0x1009a1f98: mov w10, #0xfb35
	s.x[10] = uint64(uint32(uint64(0xfb35)))
	if runtime != nil && runtime.stopPC == 0x1009a1f9c {
		return 0x1009a1f9c, true, nil
	}
	// 0x1009a1f9c: movk w10, #0x4c92, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x4c92, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1fa0 {
		return 0x1009a1fa0, true, nil
	}
	// 0x1009a1fa0: add w10, w26, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fa4 {
		return 0x1009a1fa4, true, nil
	}
	// 0x1009a1fa4: eor w10, w24, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fa8 {
		return 0x1009a1fa8, true, nil
	}
	// 0x1009a1fa8: eor w24, w9, w10
	s.x[24] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fac {
		return 0x1009a1fac, true, nil
	}
	// 0x1009a1fac: mov w9, #0xc275
	s.x[9] = uint64(uint32(uint64(0xc275)))
	if runtime != nil && runtime.stopPC == 0x1009a1fb0 {
		return 0x1009a1fb0, true, nil
	}
	// 0x1009a1fb0: movk w9, #0x95e2, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x95e2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1fb4 {
		return 0x1009a1fb4, true, nil
	}
	// 0x1009a1fb4: add w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fb8 {
		return 0x1009a1fb8, true, nil
	}
	// 0x1009a1fb8: mov w10, #0x454d
	s.x[10] = uint64(uint32(uint64(0x454d)))
	if runtime != nil && runtime.stopPC == 0x1009a1fbc {
		return 0x1009a1fbc, true, nil
	}
	// 0x1009a1fbc: movk w10, #0x3e3d, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x3e3d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1fc0 {
		return 0x1009a1fc0, true, nil
	}
	// 0x1009a1fc0: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fc4 {
		return 0x1009a1fc4, true, nil
	}
	// 0x1009a1fc4: mov w11, #0xb782
	s.x[11] = uint64(uint32(uint64(0xb782)))
	if runtime != nil && runtime.stopPC == 0x1009a1fc8 {
		return 0x1009a1fc8, true, nil
	}
	// 0x1009a1fc8: movk w11, #0x7004, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x7004, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1fcc {
		return 0x1009a1fcc, true, nil
	}
	// 0x1009a1fcc: and w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fd0 {
		return 0x1009a1fd0, true, nil
	}
	// 0x1009a1fd0: mov w15, #0x2749
	s.x[15] = uint64(uint32(uint64(0x2749)))
	if runtime != nil && runtime.stopPC == 0x1009a1fd4 {
		return 0x1009a1fd4, true, nil
	}
	// 0x1009a1fd4: movk w15, #0x7e34, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x7e34, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1fd8 {
		return 0x1009a1fd8, true, nil
	}
	// 0x1009a1fd8: eor w21, w9, w15
	s.x[21] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fdc {
		return 0x1009a1fdc, true, nil
	}
	// 0x1009a1fdc: and w21, w21, w11, lsl #1
	s.x[21] = uint64(uint32((uint64(uint32(s.x[21])) & shiftLeft(uint64(uint32(s.x[11])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a1fe0 {
		return 0x1009a1fe0, true, nil
	}
	// 0x1009a1fe0: eor w11, w21, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a1fe4 {
		return 0x1009a1fe4, true, nil
	}
	// 0x1009a1fe4: lsl w21, w10, #1
	s.x[21] = uint64(uint32(shiftLeft(uint64(uint32(s.x[10])), int(uint64(0x1)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1fe8 {
		return 0x1009a1fe8, true, nil
	}
	// 0x1009a1fe8: mov w15, #0xe59e
	s.x[15] = uint64(uint32(uint64(0xe59e)))
	if runtime != nil && runtime.stopPC == 0x1009a1fec {
		return 0x1009a1fec, true, nil
	}
	// 0x1009a1fec: movk w15, #0x9c73, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x9c73, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a1ff0 {
		return 0x1009a1ff0, true, nil
	}
	// 0x1009a1ff0: eor w23, w21, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a1ff4 {
		return 0x1009a1ff4, true, nil
	}
	// 0x1009a1ff4: and w10, w23, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[23])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a1ff8 {
		return 0x1009a1ff8, true, nil
	}
	// 0x1009a1ff8: mov w15, #0xf2ce
	s.x[15] = uint64(uint32(uint64(0xf2ce)))
	if runtime != nil && runtime.stopPC == 0x1009a1ffc {
		return 0x1009a1ffc, true, nil
	}
	// 0x1009a1ffc: movk w15, #0x4e39, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4e39, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2000 {
		return 0x1009a2000, true, nil
	}
	// 0x1009a2000: and w21, w21, w15
	s.x[21] = uint64(uint32((uint64(uint32(s.x[21])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2004 {
		return 0x1009a2004, true, nil
	}
	// 0x1009a2004: eor w10, w10, w21
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x1009a2008 {
		return 0x1009a2008, true, nil
	}
	// 0x1009a2008: mov w15, #0x1241
	s.x[15] = uint64(uint32(uint64(0x1241)))
	if runtime != nil && runtime.stopPC == 0x1009a200c {
		return 0x1009a200c, true, nil
	}
	// 0x1009a200c: movk w15, #0x4208, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4208, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2010 {
		return 0x1009a2010, true, nil
	}
	// 0x1009a2010: eor w21, w10, w15
	s.x[21] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2014 {
		return 0x1009a2014, true, nil
	}
	// 0x1009a2014: mov w15, #0xe008
	s.x[15] = uint64(uint32(uint64(0xe008)))
	if runtime != nil && runtime.stopPC == 0x1009a2018 {
		return 0x1009a2018, true, nil
	}
	// 0x1009a2018: movk w15, #0x31, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x31, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a201c {
		return 0x1009a201c, true, nil
	}
	// 0x1009a201c: eor w10, w10, w15
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2020 {
		return 0x1009a2020, true, nil
	}
	// 0x1009a2020: and w10, w10, w11, lsl #2
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & shiftLeft(uint64(uint32(s.x[11])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a2024 {
		return 0x1009a2024, true, nil
	}
	// 0x1009a2024: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a2028 {
		return 0x1009a2028, true, nil
	}
	// 0x1009a2028: lsl w11, w21, #2
	s.x[11] = uint64(uint32(shiftLeft(uint64(uint32(s.x[21])), int(uint64(0x2)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a202c {
		return 0x1009a202c, true, nil
	}
	// 0x1009a202c: mov w15, #0xcb3c
	s.x[15] = uint64(uint32(uint64(0xcb3c)))
	if runtime != nil && runtime.stopPC == 0x1009a2030 {
		return 0x1009a2030, true, nil
	}
	// 0x1009a2030: movk w15, #0x38e7, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x38e7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2034 {
		return 0x1009a2034, true, nil
	}
	// 0x1009a2034: eor w23, w11, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2038 {
		return 0x1009a2038, true, nil
	}
	// 0x1009a2038: and w21, w23, w21
	s.x[21] = uint64(uint32((uint64(uint32(s.x[23])) & uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x1009a203c {
		return 0x1009a203c, true, nil
	}
	// 0x1009a203c: mov w15, #0xf2cc
	s.x[15] = uint64(uint32(uint64(0xf2cc)))
	if runtime != nil && runtime.stopPC == 0x1009a2040 {
		return 0x1009a2040, true, nil
	}
	// 0x1009a2040: movk w15, #0x4e39, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4e39, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2044 {
		return 0x1009a2044, true, nil
	}
	// 0x1009a2044: and w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2048 {
		return 0x1009a2048, true, nil
	}
	// 0x1009a2048: eor w11, w21, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a204c {
		return 0x1009a204c, true, nil
	}
	// 0x1009a204c: mov w15, #0x30c3
	s.x[15] = uint64(uint32(uint64(0x30c3)))
	if runtime != nil && runtime.stopPC == 0x1009a2050 {
		return 0x1009a2050, true, nil
	}
	// 0x1009a2050: movk w15, #0x4618, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4618, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2054 {
		return 0x1009a2054, true, nil
	}
	// 0x1009a2054: eor w21, w11, w15
	s.x[21] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2058 {
		return 0x1009a2058, true, nil
	}
	// 0x1009a2058: mov w15, #0xc200
	s.x[15] = uint64(uint32(uint64(0xc200)))
	if runtime != nil && runtime.stopPC == 0x1009a205c {
		return 0x1009a205c, true, nil
	}
	// 0x1009a205c: movk w15, #0x821, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x821, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2060 {
		return 0x1009a2060, true, nil
	}
	// 0x1009a2060: eor w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2064 {
		return 0x1009a2064, true, nil
	}
	// 0x1009a2064: and w11, w11, w10, lsl #4
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & shiftLeft(uint64(uint32(s.x[10])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a2068 {
		return 0x1009a2068, true, nil
	}
	// 0x1009a2068: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a206c {
		return 0x1009a206c, true, nil
	}
	// 0x1009a206c: lsl w11, w21, #4
	s.x[11] = uint64(uint32(shiftLeft(uint64(uint32(s.x[21])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2070 {
		return 0x1009a2070, true, nil
	}
	// 0x1009a2070: mov w15, #0x2cf0
	s.x[15] = uint64(uint32(uint64(0x2cf0)))
	if runtime != nil && runtime.stopPC == 0x1009a2074 {
		return 0x1009a2074, true, nil
	}
	// 0x1009a2074: movk w15, #0xe39f, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xe39f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2078 {
		return 0x1009a2078, true, nil
	}
	// 0x1009a2078: eor w23, w11, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a207c {
		return 0x1009a207c, true, nil
	}
	// 0x1009a207c: and w21, w23, w21
	s.x[21] = uint64(uint32((uint64(uint32(s.x[23])) & uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x1009a2080 {
		return 0x1009a2080, true, nil
	}
	// 0x1009a2080: mov w15, #0xf2c0
	s.x[15] = uint64(uint32(uint64(0xf2c0)))
	if runtime != nil && runtime.stopPC == 0x1009a2084 {
		return 0x1009a2084, true, nil
	}
	// 0x1009a2084: movk w15, #0x4e39, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4e39, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2088 {
		return 0x1009a2088, true, nil
	}
	// 0x1009a2088: and w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a208c {
		return 0x1009a208c, true, nil
	}
	// 0x1009a208c: eor w11, w21, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a2090 {
		return 0x1009a2090, true, nil
	}
	// 0x1009a2090: mov w15, #0xd20f
	s.x[15] = uint64(uint32(uint64(0xd20f)))
	if runtime != nil && runtime.stopPC == 0x1009a2094 {
		return 0x1009a2094, true, nil
	}
	// 0x1009a2094: movk w15, #0xc20, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xc20, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2098 {
		return 0x1009a2098, true, nil
	}
	// 0x1009a2098: eor w21, w11, w15
	s.x[21] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a209c {
		return 0x1009a209c, true, nil
	}
	// 0x1009a209c: mov w15, #0x2000
	s.x[15] = uint64(uint32(uint64(0x2000)))
	if runtime != nil && runtime.stopPC == 0x1009a20a0 {
		return 0x1009a20a0, true, nil
	}
	// 0x1009a20a0: movk w15, #0x4219, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4219, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a20a4 {
		return 0x1009a20a4, true, nil
	}
	// 0x1009a20a4: eor w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a20a8 {
		return 0x1009a20a8, true, nil
	}
	// 0x1009a20a8: and w11, w11, w10, lsl #8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & shiftLeft(uint64(uint32(s.x[10])), 8, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a20ac {
		return 0x1009a20ac, true, nil
	}
	// 0x1009a20ac: mov w15, #0xf2cf
	s.x[15] = uint64(uint32(uint64(0xf2cf)))
	if runtime != nil && runtime.stopPC == 0x1009a20b0 {
		return 0x1009a20b0, true, nil
	}
	// 0x1009a20b0: movk w15, #0x4e39, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4e39, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a20b4 {
		return 0x1009a20b4, true, nil
	}
	// 0x1009a20b4: eor w10, w10, w15
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a20b8 {
		return 0x1009a20b8, true, nil
	}
	// 0x1009a20b8: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a20bc {
		return 0x1009a20bc, true, nil
	}
	// 0x1009a20bc: lsl w11, w21, #8
	s.x[11] = uint64(uint32(shiftLeft(uint64(uint32(s.x[21])), int(uint64(0x8)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a20c0 {
		return 0x1009a20c0, true, nil
	}
	// 0x1009a20c0: mov w15, #0x39f20000
	s.x[15] = uint64(uint32(uint64(0x39f20000)))
	if runtime != nil && runtime.stopPC == 0x1009a20c4 {
		return 0x1009a20c4, true, nil
	}
	// 0x1009a20c4: eor w23, w11, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a20c8 {
		return 0x1009a20c8, true, nil
	}
	// 0x1009a20c8: and w21, w23, w21
	s.x[21] = uint64(uint32((uint64(uint32(s.x[23])) & uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x1009a20cc {
		return 0x1009a20cc, true, nil
	}
	// 0x1009a20cc: and w11, w11, w1
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1009a20d0 {
		return 0x1009a20d0, true, nil
	}
	// 0x1009a20d0: mov w15, #0x46090000
	s.x[15] = uint64(uint32(uint64(0x46090000)))
	if runtime != nil && runtime.stopPC == 0x1009a20d4 {
		return 0x1009a20d4, true, nil
	}
	// 0x1009a20d4: eor w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a20d8 {
		return 0x1009a20d8, true, nil
	}
	// 0x1009a20d8: eor w11, w11, w21
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x1009a20dc {
		return 0x1009a20dc, true, nil
	}
	// 0x1009a20dc: lsl w21, w10, #0x10
	s.x[21] = uint64(uint32(shiftLeft(uint64(uint32(s.x[10])), int(uint64(0x10)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a20e0 {
		return 0x1009a20e0, true, nil
	}
	// 0x1009a20e0: mov w15, #0x72cf0000
	s.x[15] = uint64(uint32(uint64(0x72cf0000)))
	if runtime != nil && runtime.stopPC == 0x1009a20e4 {
		return 0x1009a20e4, true, nil
	}
	// 0x1009a20e4: eor w23, w21, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a20e8 {
		return 0x1009a20e8, true, nil
	}
	// 0x1009a20e8: and w11, w23, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[23])) & uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a20ec {
		return 0x1009a20ec, true, nil
	}
	// 0x1009a20ec: and w21, w21, w1
	s.x[21] = uint64(uint32((uint64(uint32(s.x[21])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1009a20f0 {
		return 0x1009a20f0, true, nil
	}
	// 0x1009a20f0: eor w10, w21, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1009a20f4 {
		return 0x1009a20f4, true, nil
	}
	// 0x1009a20f4: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a20f8 {
		return 0x1009a20f8, true, nil
	}
	// 0x1009a20f8: eor w10, w9, w10, lsl #1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ shiftLeft(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a20fc {
		return 0x1009a20fc, true, nil
	}
	// 0x1009a20fc: mov w9, #0xd77c
	s.x[9] = uint64(uint32(uint64(0xd77c)))
	if runtime != nil && runtime.stopPC == 0x1009a2100 {
		return 0x1009a2100, true, nil
	}
	// 0x1009a2100: movk w9, #0xbe46, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xbe46, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2104 {
		return 0x1009a2104, true, nil
	}
	// 0x1009a2104: eor w21, w10, w9
	s.x[21] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a2108 {
		return 0x1009a2108, true, nil
	}
	// 0x1009a2108: mov w9, #0x7004
	s.x[9] = uint64(uint32(uint64(0x7004)))
	if runtime != nil && runtime.stopPC == 0x1009a210c {
		return 0x1009a210c, true, nil
	}
	// 0x1009a210c: movk w9, #0x7245, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7245, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2110 {
		return 0x1009a2110, true, nil
	}
	// 0x1009a2110: eor w9, w3, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a2114 {
		return 0x1009a2114, true, nil
	}
	// 0x1009a2114: mov w11, #0x8560
	s.x[11] = uint64(uint32(uint64(0x8560)))
	if runtime != nil && runtime.stopPC == 0x1009a2118 {
		return 0x1009a2118, true, nil
	}
	// 0x1009a2118: movk w11, #0xd623, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd623, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a211c {
		return 0x1009a211c, true, nil
	}
	// 0x1009a211c: eor w11, w12, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a2120 {
		return 0x1009a2120, true, nil
	}
	// 0x1009a2120: and w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a2124 {
		return 0x1009a2124, true, nil
	}
	// 0x1009a2124: lsr w11, w10, #0xd
	s.x[11] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(0xd)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2128 {
		return 0x1009a2128, true, nil
	}
	// 0x1009a2128: lsl w26, w21, #0x13
	s.x[26] = uint64(uint32(shiftLeft(uint64(uint32(s.x[21])), int(uint64(0x13)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a212c {
		return 0x1009a212c, true, nil
	}
	// 0x1009a212c: mov w15, #0x61c1
	s.x[15] = uint64(uint32(uint64(0x61c1)))
	if runtime != nil && runtime.stopPC == 0x1009a2130 {
		return 0x1009a2130, true, nil
	}
	// 0x1009a2130: movk w15, #0xf7bc, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xf7bc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2134 {
		return 0x1009a2134, true, nil
	}
	// 0x1009a2134: eor w23, w26, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2138 {
		return 0x1009a2138, true, nil
	}
	// 0x1009a2138: mov w15, #-0x5111
	s.x[15] = uint64(uint32(uint64(0xffffaeef)))
	if runtime != nil && runtime.stopPC == 0x1009a213c {
		return 0x1009a213c, true, nil
	}
	// 0x1009a213c: eor w0, w11, w15
	s.x[0] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2140 {
		return 0x1009a2140, true, nil
	}
	// 0x1009a2140: and w0, w23, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[23])) & uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a2144 {
		return 0x1009a2144, true, nil
	}
	// 0x1009a2144: mov w15, #0x9e3e
	s.x[15] = uint64(uint32(uint64(0x9e3e)))
	if runtime != nil && runtime.stopPC == 0x1009a2148 {
		return 0x1009a2148, true, nil
	}
	// 0x1009a2148: movk w15, #3, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a214c {
		return 0x1009a214c, true, nil
	}
	// 0x1009a214c: and w23, w11, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2150 {
		return 0x1009a2150, true, nil
	}
	// 0x1009a2150: orr w0, w0, w23
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) | uint64(uint32(s.x[23])))))
	if runtime != nil && runtime.stopPC == 0x1009a2154 {
		return 0x1009a2154, true, nil
	}
	// 0x1009a2154: mov w15, #0xc2f0
	s.x[15] = uint64(uint32(uint64(0xc2f0)))
	if runtime != nil && runtime.stopPC == 0x1009a2158 {
		return 0x1009a2158, true, nil
	}
	// 0x1009a2158: movk w15, #0x8ac7, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x8ac7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a215c {
		return 0x1009a215c, true, nil
	}
	// 0x1009a215c: eor w23, w26, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2160 {
		return 0x1009a2160, true, nil
	}
	// 0x1009a2160: mov w15, #0x6f4f
	s.x[15] = uint64(uint32(uint64(0x6f4f)))
	if runtime != nil && runtime.stopPC == 0x1009a2164 {
		return 0x1009a2164, true, nil
	}
	// 0x1009a2164: movk w15, #0x3885, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x3885, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2168 {
		return 0x1009a2168, true, nil
	}
	// 0x1009a2168: add w23, w23, w15
	s.x[23] = uint64(uint32((uint64(uint32(s.x[23])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a216c {
		return 0x1009a216c, true, nil
	}
	// 0x1009a216c: mov w15, #0xdf05
	s.x[15] = uint64(uint32(uint64(0xdf05)))
	if runtime != nil && runtime.stopPC == 0x1009a2170 {
		return 0x1009a2170, true, nil
	}
	// 0x1009a2170: movk w15, #0xce1e, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xce1e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2174 {
		return 0x1009a2174, true, nil
	}
	// 0x1009a2174: eor w15, w26, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2178 {
		return 0x1009a2178, true, nil
	}
	// 0x1009a2178: mov w20, #0x72bc
	s.x[20] = uint64(uint32(uint64(0x72bc)))
	if runtime != nil && runtime.stopPC == 0x1009a217c {
		return 0x1009a217c, true, nil
	}
	// 0x1009a217c: movk w20, #0x7c5c, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0x7c5c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2180 {
		return 0x1009a2180, true, nil
	}
	// 0x1009a2180: add w15, w15, w20
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a2184 {
		return 0x1009a2184, true, nil
	}
	// 0x1009a2184: eor w15, w23, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2188 {
		return 0x1009a2188, true, nil
	}
	// 0x1009a2188: mov w20, #0x1df5
	s.x[20] = uint64(uint32(uint64(0x1df5)))
	if runtime != nil && runtime.stopPC == 0x1009a218c {
		return 0x1009a218c, true, nil
	}
	// 0x1009a218c: movk w20, #0x6fd9, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0x6fd9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2190 {
		return 0x1009a2190, true, nil
	}
	// 0x1009a2190: eor w23, w26, w20
	s.x[23] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a2194 {
		return 0x1009a2194, true, nil
	}
	// 0x1009a2194: mov w20, #0xb04c
	s.x[20] = uint64(uint32(uint64(0xb04c)))
	if runtime != nil && runtime.stopPC == 0x1009a2198 {
		return 0x1009a2198, true, nil
	}
	// 0x1009a2198: movk w20, #0xdd9b, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0xdd9b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a219c {
		return 0x1009a219c, true, nil
	}
	// 0x1009a219c: add w23, w23, w20
	s.x[23] = uint64(uint32((uint64(uint32(s.x[23])) + uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a21a0 {
		return 0x1009a21a0, true, nil
	}
	// 0x1009a21a0: eor w15, w15, w23
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[23])))))
	if runtime != nil && runtime.stopPC == 0x1009a21a4 {
		return 0x1009a21a4, true, nil
	}
	// 0x1009a21a4: mov w20, #0xf6ee
	s.x[20] = uint64(uint32(uint64(0xf6ee)))
	if runtime != nil && runtime.stopPC == 0x1009a21a8 {
		return 0x1009a21a8, true, nil
	}
	// 0x1009a21a8: movk w20, #0x4138, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0x4138, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21ac {
		return 0x1009a21ac, true, nil
	}
	// 0x1009a21ac: eor w23, w0, w20
	s.x[23] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a21b0 {
		return 0x1009a21b0, true, nil
	}
	// 0x1009a21b0: mov w20, #0x23f6
	s.x[20] = uint64(uint32(uint64(0x23f6)))
	if runtime != nil && runtime.stopPC == 0x1009a21b4 {
		return 0x1009a21b4, true, nil
	}
	// 0x1009a21b4: movk w20, #0x37a6, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0x37a6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21b8 {
		return 0x1009a21b8, true, nil
	}
	// 0x1009a21b8: add w23, w23, w20
	s.x[23] = uint64(uint32((uint64(uint32(s.x[23])) + uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a21bc {
		return 0x1009a21bc, true, nil
	}
	// 0x1009a21bc: mov w20, #0x6774
	s.x[20] = uint64(uint32(uint64(0x6774)))
	if runtime != nil && runtime.stopPC == 0x1009a21c0 {
		return 0x1009a21c0, true, nil
	}
	// 0x1009a21c0: movk w20, #0xe30, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0xe30, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21c4 {
		return 0x1009a21c4, true, nil
	}
	// 0x1009a21c4: eor w20, w0, w20
	s.x[20] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a21c8 {
		return 0x1009a21c8, true, nil
	}
	// 0x1009a21c8: mov w16, #0xb270
	s.x[16] = uint64(uint32(uint64(0xb270)))
	if runtime != nil && runtime.stopPC == 0x1009a21cc {
		return 0x1009a21cc, true, nil
	}
	// 0x1009a21cc: movk w16, #0x78ae, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x78ae, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21d0 {
		return 0x1009a21d0, true, nil
	}
	// 0x1009a21d0: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a21d4 {
		return 0x1009a21d4, true, nil
	}
	// 0x1009a21d4: eor w20, w23, w20
	s.x[20] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a21d8 {
		return 0x1009a21d8, true, nil
	}
	// 0x1009a21d8: mov w23, #0xe77c
	s.x[23] = uint64(uint32(uint64(0xe77c)))
	if runtime != nil && runtime.stopPC == 0x1009a21dc {
		return 0x1009a21dc, true, nil
	}
	// 0x1009a21dc: movk w23, #0xbc92, lsl #16
	s.x[23] = uint64(uint32(insertBits(uint64(uint32(s.x[23])), 0xbc92, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21e0 {
		return 0x1009a21e0, true, nil
	}
	// 0x1009a21e0: eor w23, w24, w23
	s.x[23] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[23])))))
	if runtime != nil && runtime.stopPC == 0x1009a21e4 {
		return 0x1009a21e4, true, nil
	}
	// 0x1009a21e4: mov w16, #0xf271
	s.x[16] = uint64(uint32(uint64(0xf271)))
	if runtime != nil && runtime.stopPC == 0x1009a21e8 {
		return 0x1009a21e8, true, nil
	}
	// 0x1009a21e8: movk w16, #0x93b7, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x93b7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21ec {
		return 0x1009a21ec, true, nil
	}
	// 0x1009a21ec: eor w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a21f0 {
		return 0x1009a21f0, true, nil
	}
	// 0x1009a21f0: mov w16, #0x276b
	s.x[16] = uint64(uint32(uint64(0x276b)))
	if runtime != nil && runtime.stopPC == 0x1009a21f4 {
		return 0x1009a21f4, true, nil
	}
	// 0x1009a21f4: movk w16, #0xe529, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe529, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a21f8 {
		return 0x1009a21f8, true, nil
	}
	// 0x1009a21f8: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a21fc {
		return 0x1009a21fc, true, nil
	}
	// 0x1009a21fc: eor w0, w20, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a2200 {
		return 0x1009a2200, true, nil
	}
	// 0x1009a2200: sub w15, w15, w0
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) - uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a2204 {
		return 0x1009a2204, true, nil
	}
	// 0x1009a2204: mov w16, #0x7c1a
	s.x[16] = uint64(uint32(uint64(0x7c1a)))
	if runtime != nil && runtime.stopPC == 0x1009a2208 {
		return 0x1009a2208, true, nil
	}
	// 0x1009a2208: movk w16, #0x87a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x87a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a220c {
		return 0x1009a220c, true, nil
	}
	// 0x1009a220c: eor w0, w11, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2210 {
		return 0x1009a2210, true, nil
	}
	// 0x1009a2210: mov w16, #0x1313
	s.x[16] = uint64(uint32(uint64(0x1313)))
	if runtime != nil && runtime.stopPC == 0x1009a2214 {
		return 0x1009a2214, true, nil
	}
	// 0x1009a2214: movk w16, #0xe11d, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe11d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2218 {
		return 0x1009a2218, true, nil
	}
	// 0x1009a2218: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a221c {
		return 0x1009a221c, true, nil
	}
	// 0x1009a221c: mov w16, #0xc9ef
	s.x[16] = uint64(uint32(uint64(0xc9ef)))
	if runtime != nil && runtime.stopPC == 0x1009a2220 {
		return 0x1009a2220, true, nil
	}
	// 0x1009a2220: movk w16, #0xbbc9, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xbbc9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2224 {
		return 0x1009a2224, true, nil
	}
	// 0x1009a2224: eor w20, w11, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2228 {
		return 0x1009a2228, true, nil
	}
	// 0x1009a2228: mov w16, #0xa6e8
	s.x[16] = uint64(uint32(uint64(0xa6e8)))
	if runtime != nil && runtime.stopPC == 0x1009a222c {
		return 0x1009a222c, true, nil
	}
	// 0x1009a222c: movk w16, #0x52ae, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x52ae, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2230 {
		return 0x1009a2230, true, nil
	}
	// 0x1009a2230: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2234 {
		return 0x1009a2234, true, nil
	}
	// 0x1009a2234: eor w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a2238 {
		return 0x1009a2238, true, nil
	}
	// 0x1009a2238: mov w16, #0xf6df
	s.x[16] = uint64(uint32(uint64(0xf6df)))
	if runtime != nil && runtime.stopPC == 0x1009a223c {
		return 0x1009a223c, true, nil
	}
	// 0x1009a223c: movk w16, #0xb3b0, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xb3b0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2240 {
		return 0x1009a2240, true, nil
	}
	// 0x1009a2240: eor w20, w11, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2244 {
		return 0x1009a2244, true, nil
	}
	// 0x1009a2244: mov w16, #0x99d8
	s.x[16] = uint64(uint32(uint64(0x99d8)))
	if runtime != nil && runtime.stopPC == 0x1009a2248 {
		return 0x1009a2248, true, nil
	}
	// 0x1009a2248: movk w16, #0x5ad7, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x5ad7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a224c {
		return 0x1009a224c, true, nil
	}
	// 0x1009a224c: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2250 {
		return 0x1009a2250, true, nil
	}
	// 0x1009a2250: eor w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a2254 {
		return 0x1009a2254, true, nil
	}
	// 0x1009a2254: lsr w20, w23, #2
	s.x[20] = uint64(uint32(shiftRight(uint64(uint32(s.x[23])), int(uint64(0x2)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2258 {
		return 0x1009a2258, true, nil
	}
	// 0x1009a2258: mov w16, #0x932
	s.x[16] = uint64(uint32(uint64(0x932)))
	if runtime != nil && runtime.stopPC == 0x1009a225c {
		return 0x1009a225c, true, nil
	}
	// 0x1009a225c: movk w16, #0x10df, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x10df, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2260 {
		return 0x1009a2260, true, nil
	}
	// 0x1009a2260: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2264 {
		return 0x1009a2264, true, nil
	}
	// 0x1009a2264: mov w16, #0x2c23
	s.x[16] = uint64(uint32(uint64(0x2c23)))
	if runtime != nil && runtime.stopPC == 0x1009a2268 {
		return 0x1009a2268, true, nil
	}
	// 0x1009a2268: movk w16, #0xe964, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe964, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a226c {
		return 0x1009a226c, true, nil
	}
	// 0x1009a226c: sub w0, w16, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[16])) - uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a2270 {
		return 0x1009a2270, true, nil
	}
	// 0x1009a2270: eor w11, w0, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a2274 {
		return 0x1009a2274, true, nil
	}
	// 0x1009a2274: eor w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2278 {
		return 0x1009a2278, true, nil
	}
	// 0x1009a2278: mov w15, #0xacab
	s.x[15] = uint64(uint32(uint64(0xacab)))
	if runtime != nil && runtime.stopPC == 0x1009a227c {
		return 0x1009a227c, true, nil
	}
	// 0x1009a227c: movk w15, #0x4989, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4989, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2280 {
		return 0x1009a2280, true, nil
	}
	// 0x1009a2280: eor w15, w11, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a2284 {
		return 0x1009a2284, true, nil
	}
	// 0x1009a2284: mov w16, #0x6e6f
	s.x[16] = uint64(uint32(uint64(0x6e6f)))
	if runtime != nil && runtime.stopPC == 0x1009a2288 {
		return 0x1009a2288, true, nil
	}
	// 0x1009a2288: movk w16, #0x6d6d, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x6d6d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a228c {
		return 0x1009a228c, true, nil
	}
	// 0x1009a228c: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2290 {
		return 0x1009a2290, true, nil
	}
	// 0x1009a2290: mov w16, #0x80e1
	s.x[16] = uint64(uint32(uint64(0x80e1)))
	if runtime != nil && runtime.stopPC == 0x1009a2294 {
		return 0x1009a2294, true, nil
	}
	// 0x1009a2294: movk w16, #0xe200, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe200, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2298 {
		return 0x1009a2298, true, nil
	}
	// 0x1009a2298: eor w0, w11, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a229c {
		return 0x1009a229c, true, nil
	}
	// 0x1009a229c: mov w16, #0x4225
	s.x[16] = uint64(uint32(uint64(0x4225)))
	if runtime != nil && runtime.stopPC == 0x1009a22a0 {
		return 0x1009a22a0, true, nil
	}
	// 0x1009a22a0: movk w16, #0xc6e4, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xc6e4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22a4 {
		return 0x1009a22a4, true, nil
	}
	// 0x1009a22a4: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a22a8 {
		return 0x1009a22a8, true, nil
	}
	// 0x1009a22a8: eor w15, w15, w0
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a22ac {
		return 0x1009a22ac, true, nil
	}
	// 0x1009a22ac: mov w16, #0x6f60
	s.x[16] = uint64(uint32(uint64(0x6f60)))
	if runtime != nil && runtime.stopPC == 0x1009a22b0 {
		return 0x1009a22b0, true, nil
	}
	// 0x1009a22b0: movk w16, #0xab8a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xab8a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22b4 {
		return 0x1009a22b4, true, nil
	}
	// 0x1009a22b4: eor w11, w11, w16
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a22b8 {
		return 0x1009a22b8, true, nil
	}
	// 0x1009a22b8: mov w16, #0xada6
	s.x[16] = uint64(uint32(uint64(0xada6)))
	if runtime != nil && runtime.stopPC == 0x1009a22bc {
		return 0x1009a22bc, true, nil
	}
	// 0x1009a22bc: movk w16, #0x8f6e, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x8f6e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22c0 {
		return 0x1009a22c0, true, nil
	}
	// 0x1009a22c0: add w11, w11, w16
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a22c4 {
		return 0x1009a22c4, true, nil
	}
	// 0x1009a22c4: eor w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a22c8 {
		return 0x1009a22c8, true, nil
	}
	// 0x1009a22c8: mov w15, #0xc761
	s.x[15] = uint64(uint32(uint64(0xc761)))
	if runtime != nil && runtime.stopPC == 0x1009a22cc {
		return 0x1009a22cc, true, nil
	}
	// 0x1009a22cc: movk w15, #0x820d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x820d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22d0 {
		return 0x1009a22d0, true, nil
	}
	// 0x1009a22d0: eor w15, w26, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a22d4 {
		return 0x1009a22d4, true, nil
	}
	// 0x1009a22d4: mov w16, #0x6ae0
	s.x[16] = uint64(uint32(uint64(0x6ae0)))
	if runtime != nil && runtime.stopPC == 0x1009a22d8 {
		return 0x1009a22d8, true, nil
	}
	// 0x1009a22d8: movk w16, #0x304f, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x304f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22dc {
		return 0x1009a22dc, true, nil
	}
	// 0x1009a22dc: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a22e0 {
		return 0x1009a22e0, true, nil
	}
	// 0x1009a22e0: mov w16, #0x784e
	s.x[16] = uint64(uint32(uint64(0x784e)))
	if runtime != nil && runtime.stopPC == 0x1009a22e4 {
		return 0x1009a22e4, true, nil
	}
	// 0x1009a22e4: movk w16, #0xf837, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xf837, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22e8 {
		return 0x1009a22e8, true, nil
	}
	// 0x1009a22e8: eor w0, w26, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a22ec {
		return 0x1009a22ec, true, nil
	}
	// 0x1009a22ec: mov w16, #0xd5f1
	s.x[16] = uint64(uint32(uint64(0xd5f1)))
	if runtime != nil && runtime.stopPC == 0x1009a22f0 {
		return 0x1009a22f0, true, nil
	}
	// 0x1009a22f0: movk w16, #0x4a75, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x4a75, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a22f4 {
		return 0x1009a22f4, true, nil
	}
	// 0x1009a22f4: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a22f8 {
		return 0x1009a22f8, true, nil
	}
	// 0x1009a22f8: eor w15, w15, w0
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a22fc {
		return 0x1009a22fc, true, nil
	}
	// 0x1009a22fc: lsr w0, w24, #4
	s.x[0] = uint64(uint32(shiftRight(uint64(uint32(s.x[24])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2300 {
		return 0x1009a2300, true, nil
	}
	// 0x1009a2300: and w0, w0, #0xfffff80
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(0xfffff80))))
	if runtime != nil && runtime.stopPC == 0x1009a2304 {
		return 0x1009a2304, true, nil
	}
	// 0x1009a2304: mov w16, #0xbf2f
	s.x[16] = uint64(uint32(uint64(0xbf2f)))
	if runtime != nil && runtime.stopPC == 0x1009a2308 {
		return 0x1009a2308, true, nil
	}
	// 0x1009a2308: movk w16, #0x513a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x513a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a230c {
		return 0x1009a230c, true, nil
	}
	// 0x1009a230c: eor w26, w26, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[26])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2310 {
		return 0x1009a2310, true, nil
	}
	// 0x1009a2310: mov w16, #0x1292
	s.x[16] = uint64(uint32(uint64(0x1292)))
	if runtime != nil && runtime.stopPC == 0x1009a2314 {
		return 0x1009a2314, true, nil
	}
	// 0x1009a2314: movk w16, #0xe378, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe378, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2318 {
		return 0x1009a2318, true, nil
	}
	// 0x1009a2318: add w26, w26, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a231c {
		return 0x1009a231c, true, nil
	}
	// 0x1009a231c: eor w15, w15, w26
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a2320 {
		return 0x1009a2320, true, nil
	}
	// 0x1009a2320: mov w16, #0xd051
	s.x[16] = uint64(uint32(uint64(0xd051)))
	if runtime != nil && runtime.stopPC == 0x1009a2324 {
		return 0x1009a2324, true, nil
	}
	// 0x1009a2324: movk w16, #0x41d5, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x41d5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2328 {
		return 0x1009a2328, true, nil
	}
	// 0x1009a2328: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a232c {
		return 0x1009a232c, true, nil
	}
	// 0x1009a232c: add w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a2330 {
		return 0x1009a2330, true, nil
	}
	// 0x1009a2330: ubfx w15, w24, #4, #7
	s.x[15] = uint64(uint32((uint64(uint32(s.x[24])) >> 4) & bitMask(7, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2334 {
		return 0x1009a2334, true, nil
	}
	// 0x1009a2334: mov w16, #0xdb76
	s.x[16] = uint64(uint32(uint64(0xdb76)))
	if runtime != nil && runtime.stopPC == 0x1009a2338 {
		return 0x1009a2338, true, nil
	}
	// 0x1009a2338: movk w16, #0xb4ea, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xb4ea, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a233c {
		return 0x1009a233c, true, nil
	}
	// 0x1009a233c: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2340 {
		return 0x1009a2340, true, nil
	}
	// 0x1009a2340: mov w16, #0xfbff
	s.x[16] = uint64(uint32(uint64(0xfbff)))
	if runtime != nil && runtime.stopPC == 0x1009a2344 {
		return 0x1009a2344, true, nil
	}
	// 0x1009a2344: movk w16, #0x76fa, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x76fa, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2348 {
		return 0x1009a2348, true, nil
	}
	// 0x1009a2348: eor w26, w16, w23, lsl #28
	s.x[26] = uint64(uint32((uint64(uint32(s.x[16])) ^ shiftLeft(uint64(uint32(s.x[23])), 28, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a234c {
		return 0x1009a234c, true, nil
	}
	// 0x1009a234c: and w15, w15, w26
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a2350 {
		return 0x1009a2350, true, nil
	}
	// 0x1009a2350: mov w16, #0x1a
	s.x[16] = uint64(uint32(uint64(0x1a)))
	if runtime != nil && runtime.stopPC == 0x1009a2354 {
		return 0x1009a2354, true, nil
	}
	// 0x1009a2354: movk w16, #0x9000, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x9000, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2358 {
		return 0x1009a2358, true, nil
	}
	// 0x1009a2358: and w26, w15, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a235c {
		return 0x1009a235c, true, nil
	}
	// 0x1009a235c: bfi w15, w20, #0x1e, #1
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), uint64(uint32(s.x[20])), 30, 1, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2360 {
		return 0x1009a2360, true, nil
	}
	// 0x1009a2360: mov w16, #0xf0e5
	s.x[16] = uint64(uint32(uint64(0xf0e5)))
	if runtime != nil && runtime.stopPC == 0x1009a2364 {
		return 0x1009a2364, true, nil
	}
	// 0x1009a2364: movk w16, #0x6513, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x6513, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2368 {
		return 0x1009a2368, true, nil
	}
	// 0x1009a2368: eor w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a236c {
		return 0x1009a236c, true, nil
	}
	// 0x1009a236c: mov w16, #0x24a0
	s.x[16] = uint64(uint32(uint64(0x24a0)))
	if runtime != nil && runtime.stopPC == 0x1009a2370 {
		return 0x1009a2370, true, nil
	}
	// 0x1009a2370: movk w16, #0x4b15, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x4b15, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2374 {
		return 0x1009a2374, true, nil
	}
	// 0x1009a2374: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2378 {
		return 0x1009a2378, true, nil
	}
	// 0x1009a2378: and w15, w15, w0
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a237c {
		return 0x1009a237c, true, nil
	}
	// 0x1009a237c: orr w15, w15, w26
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) | uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a2380 {
		return 0x1009a2380, true, nil
	}
	// 0x1009a2380: mov w16, #0x44a4
	s.x[16] = uint64(uint32(uint64(0x44a4)))
	if runtime != nil && runtime.stopPC == 0x1009a2384 {
		return 0x1009a2384, true, nil
	}
	// 0x1009a2384: movk w16, #0x388c, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x388c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2388 {
		return 0x1009a2388, true, nil
	}
	// 0x1009a2388: eor w0, w12, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a238c {
		return 0x1009a238c, true, nil
	}
	// 0x1009a238c: mov w16, #0x7a9f
	s.x[16] = uint64(uint32(uint64(0x7a9f)))
	if runtime != nil && runtime.stopPC == 0x1009a2390 {
		return 0x1009a2390, true, nil
	}
	// 0x1009a2390: movk w16, #0x29dc, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x29dc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2394 {
		return 0x1009a2394, true, nil
	}
	// 0x1009a2394: eor w20, w3, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2398 {
		return 0x1009a2398, true, nil
	}
	// 0x1009a2398: and w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a239c {
		return 0x1009a239c, true, nil
	}
	// 0x1009a239c: mov w16, #0x3e3b
	s.x[16] = uint64(uint32(uint64(0x3e3b)))
	if runtime != nil && runtime.stopPC == 0x1009a23a0 {
		return 0x1009a23a0, true, nil
	}
	// 0x1009a23a0: movk w16, #0x1150, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x1150, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23a4 {
		return 0x1009a23a4, true, nil
	}
	// 0x1009a23a4: and w20, w3, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a23a8 {
		return 0x1009a23a8, true, nil
	}
	// 0x1009a23a8: eor w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a23ac {
		return 0x1009a23ac, true, nil
	}
	// 0x1009a23ac: mov w16, #0x32b2
	s.x[16] = uint64(uint32(uint64(0x32b2)))
	if runtime != nil && runtime.stopPC == 0x1009a23b0 {
		return 0x1009a23b0, true, nil
	}
	// 0x1009a23b0: movk w16, #0xef88, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xef88, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23b4 {
		return 0x1009a23b4, true, nil
	}
	// 0x1009a23b4: eor w20, w10, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a23b8 {
		return 0x1009a23b8, true, nil
	}
	// 0x1009a23b8: mov w16, #0xc5e4
	s.x[16] = uint64(uint32(uint64(0xc5e4)))
	if runtime != nil && runtime.stopPC == 0x1009a23bc {
		return 0x1009a23bc, true, nil
	}
	// 0x1009a23bc: movk w16, #0xfeaf, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xfeaf, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23c0 {
		return 0x1009a23c0, true, nil
	}
	// 0x1009a23c0: eor w26, w0, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a23c4 {
		return 0x1009a23c4, true, nil
	}
	// 0x1009a23c4: and w20, w20, w26
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) & uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a23c8 {
		return 0x1009a23c8, true, nil
	}
	// 0x1009a23c8: mov w16, #0x9680
	s.x[16] = uint64(uint32(uint64(0x9680)))
	if runtime != nil && runtime.stopPC == 0x1009a23cc {
		return 0x1009a23cc, true, nil
	}
	// 0x1009a23cc: movk w16, #0x2075, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x2075, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23d0 {
		return 0x1009a23d0, true, nil
	}
	// 0x1009a23d0: eor w26, w17, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a23d4 {
		return 0x1009a23d4, true, nil
	}
	// 0x1009a23d4: mov w16, #0x60ae
	s.x[16] = uint64(uint32(uint64(0x60ae)))
	if runtime != nil && runtime.stopPC == 0x1009a23d8 {
		return 0x1009a23d8, true, nil
	}
	// 0x1009a23d8: movk w16, #0x87ed, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x87ed, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23dc {
		return 0x1009a23dc, true, nil
	}
	// 0x1009a23dc: and w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a23e0 {
		return 0x1009a23e0, true, nil
	}
	// 0x1009a23e0: eor w0, w20, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a23e4 {
		return 0x1009a23e4, true, nil
	}
	// 0x1009a23e4: mov w20, #0xf564
	s.x[20] = uint64(uint32(uint64(0xf564)))
	if runtime != nil && runtime.stopPC == 0x1009a23e8 {
		return 0x1009a23e8, true, nil
	}
	// 0x1009a23e8: movk w20, #0xa466, lsl #16
	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0xa466, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23ec {
		return 0x1009a23ec, true, nil
	}
	// 0x1009a23ec: and w20, w12, w20
	s.x[20] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a23f0 {
		return 0x1009a23f0, true, nil
	}
	// 0x1009a23f0: mov w16, #0xda8
	s.x[16] = uint64(uint32(uint64(0xda8)))
	if runtime != nil && runtime.stopPC == 0x1009a23f4 {
		return 0x1009a23f4, true, nil
	}
	// 0x1009a23f4: movk w16, #0x8338, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x8338, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a23f8 {
		return 0x1009a23f8, true, nil
	}
	// 0x1009a23f8: eor w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a23fc {
		return 0x1009a23fc, true, nil
	}
	// 0x1009a23fc: eor w9, w20, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a2400 {
		return 0x1009a2400, true, nil
	}
	// 0x1009a2400: mov w16, #0xbf5b
	s.x[16] = uint64(uint32(uint64(0xbf5b)))
	if runtime != nil && runtime.stopPC == 0x1009a2404 {
		return 0x1009a2404, true, nil
	}
	// 0x1009a2404: movk w16, #0x7952, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x7952, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2408 {
		return 0x1009a2408, true, nil
	}
	// 0x1009a2408: eor w20, w0, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a240c {
		return 0x1009a240c, true, nil
	}
	// 0x1009a240c: and w9, w20, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[20])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a2410 {
		return 0x1009a2410, true, nil
	}
	// 0x1009a2410: mov w16, #0x55cb
	s.x[16] = uint64(uint32(uint64(0x55cb)))
	if runtime != nil && runtime.stopPC == 0x1009a2414 {
		return 0x1009a2414, true, nil
	}
	// 0x1009a2414: movk w16, #0xfcb2, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xfcb2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2418 {
		return 0x1009a2418, true, nil
	}
	// 0x1009a2418: eor w20, w17, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a241c {
		return 0x1009a241c, true, nil
	}
	// 0x1009a241c: mov w16, #0xdf23
	s.x[16] = uint64(uint32(uint64(0xdf23)))
	if runtime != nil && runtime.stopPC == 0x1009a2420 {
		return 0x1009a2420, true, nil
	}
	// 0x1009a2420: movk w16, #0xcd7, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xcd7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2424 {
		return 0x1009a2424, true, nil
	}
	// 0x1009a2424: eor w17, w17, w16
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2428 {
		return 0x1009a2428, true, nil
	}
	// 0x1009a2428: mov w16, #0x7737
	s.x[16] = uint64(uint32(uint64(0x7737)))
	if runtime != nil && runtime.stopPC == 0x1009a242c {
		return 0x1009a242c, true, nil
	}
	// 0x1009a242c: movk w16, #0xf8e5, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xf8e5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2430 {
		return 0x1009a2430, true, nil
	}
	// 0x1009a2430: and w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2434 {
		return 0x1009a2434, true, nil
	}
	// 0x1009a2434: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a2438 {
		return 0x1009a2438, true, nil
	}
	// 0x1009a2438: extr w0, w23, w24, #0x19
	s.x[0] = uint64(uint32(extractPair(uint64(uint32(s.x[23])), uint64(uint32(s.x[24])), 25, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a243c {
		return 0x1009a243c, true, nil
	}
	// 0x1009a243c: mov w16, #0xe3f5
	s.x[16] = uint64(uint32(uint64(0xe3f5)))
	if runtime != nil && runtime.stopPC == 0x1009a2440 {
		return 0x1009a2440, true, nil
	}
	// 0x1009a2440: movk w16, #0x7294, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x7294, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2444 {
		return 0x1009a2444, true, nil
	}
	// 0x1009a2444: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2448 {
		return 0x1009a2448, true, nil
	}
	// 0x1009a2448: eor w0, w0, w23, ror #6
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ rotateRight(uint64(uint32(s.x[23])), 6, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a244c {
		return 0x1009a244c, true, nil
	}
	// 0x1009a244c: eor w15, w0, w15, ror #7
	s.x[15] = uint64(uint32((uint64(uint32(s.x[0])) ^ rotateRight(uint64(uint32(s.x[15])), 7, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a2450 {
		return 0x1009a2450, true, nil
	}
	// 0x1009a2450: extr w0, w10, w21, #2
	s.x[0] = uint64(uint32(extractPair(uint64(uint32(s.x[10])), uint64(uint32(s.x[21])), 2, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2454 {
		return 0x1009a2454, true, nil
	}
	// 0x1009a2454: eor w0, w0, w21, ror #22
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ rotateRight(uint64(uint32(s.x[21])), 22, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a2458 {
		return 0x1009a2458, true, nil
	}
	// 0x1009a2458: eor w11, w0, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a245c {
		return 0x1009a245c, true, nil
	}
	// 0x1009a245c: mov w16, #0x6980
	s.x[16] = uint64(uint32(uint64(0x6980)))
	if runtime != nil && runtime.stopPC == 0x1009a2460 {
		return 0x1009a2460, true, nil
	}
	// 0x1009a2460: movk w16, #0xdf8a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xdf8a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2464 {
		return 0x1009a2464, true, nil
	}
	// 0x1009a2464: add w0, w26, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2468 {
		return 0x1009a2468, true, nil
	}
	// 0x1009a2468: mov w16, #0xaa35
	s.x[16] = uint64(uint32(uint64(0xaa35)))
	if runtime != nil && runtime.stopPC == 0x1009a246c {
		return 0x1009a246c, true, nil
	}
	// 0x1009a246c: movk w16, #0x34d, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x34d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2470 {
		return 0x1009a2470, true, nil
	}
	// 0x1009a2470: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2474 {
		return 0x1009a2474, true, nil
	}
	// 0x1009a2474: eor w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a2478 {
		return 0x1009a2478, true, nil
	}
	// 0x1009a2478: mov w16, #0x20dd
	s.x[16] = uint64(uint32(uint64(0x20dd)))
	if runtime != nil && runtime.stopPC == 0x1009a247c {
		return 0x1009a247c, true, nil
	}
	// 0x1009a247c: movk w16, #0xf328, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xf328, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2480 {
		return 0x1009a2480, true, nil
	}
	// 0x1009a2480: add w17, w17, w16
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2484 {
		return 0x1009a2484, true, nil
	}
	// 0x1009a2484: eor w17, w0, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1009a2488 {
		return 0x1009a2488, true, nil
	}
	// 0x1009a2488: mov w16, #0x5c23
	s.x[16] = uint64(uint32(uint64(0x5c23)))
	if runtime != nil && runtime.stopPC == 0x1009a248c {
		return 0x1009a248c, true, nil
	}
	// 0x1009a248c: movk w16, #0xe8ab, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe8ab, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2490 {
		return 0x1009a2490, true, nil
	}
	// 0x1009a2490: eor w0, w15, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2494 {
		return 0x1009a2494, true, nil
	}
	// 0x1009a2494: mov w16, #0x38e4
	s.x[16] = uint64(uint32(uint64(0x38e4)))
	if runtime != nil && runtime.stopPC == 0x1009a2498 {
		return 0x1009a2498, true, nil
	}
	// 0x1009a2498: movk w16, #0xc10, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xc10, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a249c {
		return 0x1009a249c, true, nil
	}
	// 0x1009a249c: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24a0 {
		return 0x1009a24a0, true, nil
	}
	// 0x1009a24a0: mov w16, #0xc34c
	s.x[16] = uint64(uint32(uint64(0xc34c)))
	if runtime != nil && runtime.stopPC == 0x1009a24a4 {
		return 0x1009a24a4, true, nil
	}
	// 0x1009a24a4: movk w16, #0x88c2, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x88c2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a24a8 {
		return 0x1009a24a8, true, nil
	}
	// 0x1009a24a8: eor w20, w15, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24ac {
		return 0x1009a24ac, true, nil
	}
	// 0x1009a24ac: mov w16, #0xa78d
	s.x[16] = uint64(uint32(uint64(0xa78d)))
	if runtime != nil && runtime.stopPC == 0x1009a24b0 {
		return 0x1009a24b0, true, nil
	}
	// 0x1009a24b0: movk w16, #0x6c79, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x6c79, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a24b4 {
		return 0x1009a24b4, true, nil
	}
	// 0x1009a24b4: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24b8 {
		return 0x1009a24b8, true, nil
	}
	// 0x1009a24b8: eor w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a24bc {
		return 0x1009a24bc, true, nil
	}
	// 0x1009a24bc: mov w16, #0x7551
	s.x[16] = uint64(uint32(uint64(0x7551)))
	if runtime != nil && runtime.stopPC == 0x1009a24c0 {
		return 0x1009a24c0, true, nil
	}
	// 0x1009a24c0: movk w16, #0xee25, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xee25, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a24c4 {
		return 0x1009a24c4, true, nil
	}
	// 0x1009a24c4: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24c8 {
		return 0x1009a24c8, true, nil
	}
	// 0x1009a24c8: mov w16, #0x1192
	s.x[16] = uint64(uint32(uint64(0x1192)))
	if runtime != nil && runtime.stopPC == 0x1009a24cc {
		return 0x1009a24cc, true, nil
	}
	// 0x1009a24cc: movk w16, #0xa9e, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a24d0 {
		return 0x1009a24d0, true, nil
	}
	// 0x1009a24d0: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24d4 {
		return 0x1009a24d4, true, nil
	}
	// 0x1009a24d4: eor w15, w0, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a24d8 {
		return 0x1009a24d8, true, nil
	}
	// 0x1009a24d8: mov w16, #0xfb14
	s.x[16] = uint64(uint32(uint64(0xfb14)))
	if runtime != nil && runtime.stopPC == 0x1009a24dc {
		return 0x1009a24dc, true, nil
	}
	// 0x1009a24dc: movk w16, #0x6c82, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x6c82, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a24e0 {
		return 0x1009a24e0, true, nil
	}
	// 0x1009a24e0: eor w0, w24, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24e4 {
		return 0x1009a24e4, true, nil
	}
	// 0x1009a24e4: mov w16, #0x4eb
	s.x[16] = uint64(uint32(uint64(0x4eb)))
	if runtime != nil && runtime.stopPC == 0x1009a24e8 {
		return 0x1009a24e8, true, nil
	}
	// 0x1009a24e8: movk w16, #0x937d, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x937d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a24ec {
		return 0x1009a24ec, true, nil
	}
	// 0x1009a24ec: eor w20, w24, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a24f0 {
		return 0x1009a24f0, true, nil
	}
	// 0x1009a24f0: and w20, w25, w20
	s.x[20] = uint64(uint32((uint64(uint32(s.x[25])) & uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a24f4 {
		return 0x1009a24f4, true, nil
	}
	// 0x1009a24f4: and w0, w14, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a24f8 {
		return 0x1009a24f8, true, nil
	}
	// 0x1009a24f8: orr w0, w0, w20
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) | uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a24fc {
		return 0x1009a24fc, true, nil
	}
	// 0x1009a24fc: mov w16, #0xf40f
	s.x[16] = uint64(uint32(uint64(0xf40f)))
	if runtime != nil && runtime.stopPC == 0x1009a2500 {
		return 0x1009a2500, true, nil
	}
	// 0x1009a2500: movk w16, #0xb844, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xb844, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2504 {
		return 0x1009a2504, true, nil
	}
	// 0x1009a2504: eor w20, w0, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2508 {
		return 0x1009a2508, true, nil
	}
	// 0x1009a2508: mov w16, #0x9ba5
	s.x[16] = uint64(uint32(uint64(0x9ba5)))
	if runtime != nil && runtime.stopPC == 0x1009a250c {
		return 0x1009a250c, true, nil
	}
	// 0x1009a250c: movk w16, #0x3f77, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x3f77, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2510 {
		return 0x1009a2510, true, nil
	}
	// 0x1009a2510: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2514 {
		return 0x1009a2514, true, nil
	}
	// 0x1009a2514: mov w16, #0x89b
	s.x[16] = uint64(uint32(uint64(0x89b)))
	if runtime != nil && runtime.stopPC == 0x1009a2518 {
		return 0x1009a2518, true, nil
	}
	// 0x1009a2518: movk w16, #0x6570, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x6570, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a251c {
		return 0x1009a251c, true, nil
	}
	// 0x1009a251c: eor w26, w0, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2520 {
		return 0x1009a2520, true, nil
	}
	// 0x1009a2520: mov w16, #0x6731
	s.x[16] = uint64(uint32(uint64(0x6731)))
	if runtime != nil && runtime.stopPC == 0x1009a2524 {
		return 0x1009a2524, true, nil
	}
	// 0x1009a2524: movk w16, #0xe243, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe243, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2528 {
		return 0x1009a2528, true, nil
	}
	// 0x1009a2528: add w26, w26, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a252c {
		return 0x1009a252c, true, nil
	}
	// 0x1009a252c: eor w20, w20, w26
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a2530 {
		return 0x1009a2530, true, nil
	}
	// 0x1009a2530: mov w16, #0xe0fc
	s.x[16] = uint64(uint32(uint64(0xe0fc)))
	if runtime != nil && runtime.stopPC == 0x1009a2534 {
		return 0x1009a2534, true, nil
	}
	// 0x1009a2534: movk w16, #0xd24, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xd24, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2538 {
		return 0x1009a2538, true, nil
	}
	// 0x1009a2538: eor w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a253c {
		return 0x1009a253c, true, nil
	}
	// 0x1009a253c: mov w16, #0x8f58
	s.x[16] = uint64(uint32(uint64(0x8f58)))
	if runtime != nil && runtime.stopPC == 0x1009a2540 {
		return 0x1009a2540, true, nil
	}
	// 0x1009a2540: movk w16, #0x8a17, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x8a17, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2544 {
		return 0x1009a2544, true, nil
	}
	// 0x1009a2544: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2548 {
		return 0x1009a2548, true, nil
	}
	// 0x1009a2548: eor w0, w20, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a254c {
		return 0x1009a254c, true, nil
	}
	// 0x1009a254c: mov w16, #0xaacf
	s.x[16] = uint64(uint32(uint64(0xaacf)))
	if runtime != nil && runtime.stopPC == 0x1009a2550 {
		return 0x1009a2550, true, nil
	}
	// 0x1009a2550: movk w16, #0x5a24, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x5a24, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2554 {
		return 0x1009a2554, true, nil
	}
	// 0x1009a2554: eor w20, w11, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2558 {
		return 0x1009a2558, true, nil
	}
	// 0x1009a2558: mov w16, #0x4d89
	s.x[16] = uint64(uint32(uint64(0x4d89)))
	if runtime != nil && runtime.stopPC == 0x1009a255c {
		return 0x1009a255c, true, nil
	}
	// 0x1009a255c: movk w16, #0x19d6, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x19d6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2560 {
		return 0x1009a2560, true, nil
	}
	// 0x1009a2560: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2564 {
		return 0x1009a2564, true, nil
	}
	// 0x1009a2564: mov w16, #0xd0aa
	s.x[16] = uint64(uint32(uint64(0xd0aa)))
	if runtime != nil && runtime.stopPC == 0x1009a2568 {
		return 0x1009a2568, true, nil
	}
	// 0x1009a2568: movk w16, #0x456f, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x456f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a256c {
		return 0x1009a256c, true, nil
	}
	// 0x1009a256c: eor w26, w11, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2570 {
		return 0x1009a2570, true, nil
	}
	// 0x1009a2570: mov w16, #0x37ee
	s.x[16] = uint64(uint32(uint64(0x37ee)))
	if runtime != nil && runtime.stopPC == 0x1009a2574 {
		return 0x1009a2574, true, nil
	}
	// 0x1009a2574: movk w16, #0x69d, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x69d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2578 {
		return 0x1009a2578, true, nil
	}
	// 0x1009a2578: add w26, w26, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a257c {
		return 0x1009a257c, true, nil
	}
	// 0x1009a257c: eor w20, w20, w26
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a2580 {
		return 0x1009a2580, true, nil
	}
	// 0x1009a2580: mov w16, #0x1865
	s.x[16] = uint64(uint32(uint64(0x1865)))
	if runtime != nil && runtime.stopPC == 0x1009a2584 {
		return 0x1009a2584, true, nil
	}
	// 0x1009a2584: movk w16, #0xa4d6, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa4d6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2588 {
		return 0x1009a2588, true, nil
	}
	// 0x1009a2588: eor w11, w11, w16
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a258c {
		return 0x1009a258c, true, nil
	}
	// 0x1009a258c: mov w16, #0xff23
	s.x[16] = uint64(uint32(uint64(0xff23)))
	if runtime != nil && runtime.stopPC == 0x1009a2590 {
		return 0x1009a2590, true, nil
	}
	// 0x1009a2590: movk w16, #0xe724, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe724, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2594 {
		return 0x1009a2594, true, nil
	}
	// 0x1009a2594: add w11, w11, w16
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a2598 {
		return 0x1009a2598, true, nil
	}
	// 0x1009a2598: eor w11, w20, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1009a259c {
		return 0x1009a259c, true, nil
	}
	// 0x1009a259c: mov w16, #0x2a34
	s.x[16] = uint64(uint32(uint64(0x2a34)))
	if runtime != nil && runtime.stopPC == 0x1009a25a0 {
		return 0x1009a25a0, true, nil
	}
	// 0x1009a25a0: movk w16, #0xe746, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe746, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a25a4 {
		return 0x1009a25a4, true, nil
	}
	// 0x1009a25a4: eor w20, w9, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a25a8 {
		return 0x1009a25a8, true, nil
	}
	// 0x1009a25a8: mov w16, #0x5762
	s.x[16] = uint64(uint32(uint64(0x5762)))
	if runtime != nil && runtime.stopPC == 0x1009a25ac {
		return 0x1009a25ac, true, nil
	}
	// 0x1009a25ac: movk w16, #0xb5c0, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xb5c0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a25b0 {
		return 0x1009a25b0, true, nil
	}
	// 0x1009a25b0: add w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a25b4 {
		return 0x1009a25b4, true, nil
	}
	// 0x1009a25b4: mov w16, #0x14d7
	s.x[16] = uint64(uint32(uint64(0x14d7)))
	if runtime != nil && runtime.stopPC == 0x1009a25b8 {
		return 0x1009a25b8, true, nil
	}
	// 0x1009a25b8: movk w16, #0x8b0a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x8b0a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a25bc {
		return 0x1009a25bc, true, nil
	}
	// 0x1009a25bc: eor w26, w9, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a25c0 {
		return 0x1009a25c0, true, nil
	}
	// 0x1009a25c0: mov w16, #0x6983
	s.x[16] = uint64(uint32(uint64(0x6983)))
	if runtime != nil && runtime.stopPC == 0x1009a25c4 {
		return 0x1009a25c4, true, nil
	}
	// 0x1009a25c4: movk w16, #0xd98c, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xd98c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a25c8 {
		return 0x1009a25c8, true, nil
	}
	// 0x1009a25c8: add w26, w26, w16
	s.x[26] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a25cc {
		return 0x1009a25cc, true, nil
	}
	// 0x1009a25cc: eor w20, w20, w26
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[26])))))
	if runtime != nil && runtime.stopPC == 0x1009a25d0 {
		return 0x1009a25d0, true, nil
	}
	// 0x1009a25d0: mov w16, #0xf60f
	s.x[16] = uint64(uint32(uint64(0xf60f)))
	if runtime != nil && runtime.stopPC == 0x1009a25d4 {
		return 0x1009a25d4, true, nil
	}
	// 0x1009a25d4: movk w16, #0xebf3, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xebf3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1009a25d8 {
		return 0x1009a25d8, true, nil
	}
	// 0x1009a25d8: eor w9, w9, w16
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1009a25dc {
		return 0x1009a25dc, true, nil
	}
	// 0x1009a25dc: add w9, w9, w2
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[2])))))
	if runtime != nil && runtime.stopPC == 0x1009a25e0 {
		return 0x1009a25e0, true, nil
	}
	// 0x1009a25e0: eor w9, w20, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a25e4 {
		return 0x1009a25e4, true, nil
	}
	// 0x1009a25e4: umull x20, w28, w30
	s.x[20] = uint64(uint64(uint32(uint64(uint32(s.x[28])))) * uint64(uint32(uint64(uint32(s.x[30])))))
	if runtime != nil && runtime.stopPC == 0x1009a25e8 {
		return 0x1009a25e8, true, nil
	}
	// 0x1009a25e8: lsr x20, x20, #0x25
	s.x[20] = uint64(shiftRight(s.x[20], int(uint64(0x25)), 64))
	if runtime != nil && runtime.stopPC == 0x1009a25ec {
		return 0x1009a25ec, true, nil
	}
	// 0x1009a25ec: msub w20, w20, w5, w28
	s.x[20] = uint64(uint32((uint64(uint32(s.x[28])) - uint64(uint32(s.x[20]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x1009a25f0 {
		return 0x1009a25f0, true, nil
	}
	// 0x1009a25f0: add w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1009a25f4 {
		return 0x1009a25f4, true, nil
	}
	// 0x1009a25f4: umull x11, w6, w30
	s.x[11] = uint64(uint64(uint32(uint64(uint32(s.x[6])))) * uint64(uint32(uint64(uint32(s.x[30])))))
	if runtime != nil && runtime.stopPC == 0x1009a25f8 {
		return 0x1009a25f8, true, nil
	}
	// 0x1009a25f8: lsr x11, x11, #0x25
	s.x[11] = uint64(shiftRight(s.x[11], int(uint64(0x25)), 64))
	if runtime != nil && runtime.stopPC == 0x1009a25fc {
		return 0x1009a25fc, true, nil
	}
	// 0x1009a25fc: msub w11, w11, w5, w6
	s.x[11] = uint64(uint32((uint64(uint32(s.x[6])) - uint64(uint32(s.x[11]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x1009a2600 {
		return 0x1009a2600, true, nil
	}
	// 0x1009a2600: ldr w11, [x8, w11, uxtw #2]
	value447, err := s.load((s.x[8] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value447))
	if runtime != nil && runtime.stopPC == 0x1009a2604 {
		return 0x1009a2604, true, nil
	}
	// 0x1009a2604: ldr w20, [x8, w20, uxtw #2]
	value448, err := s.load((s.x[8] + shiftLeft(uint64(uint32(uint64(uint32(s.x[20])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[20] = uint64(uint32(value448))
	if runtime != nil && runtime.stopPC == 0x1009a2608 {
		return 0x1009a2608, true, nil
	}
	// 0x1009a2608: add w11, w11, w20
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x1009a260c {
		return 0x1009a260c, true, nil
	}
	// 0x1009a260c: add w17, w17, w0
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1009a2610 {
		return 0x1009a2610, true, nil
	}
	// 0x1009a2610: add w11, w11, w17
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1009a2614 {
		return 0x1009a2614, true, nil
	}
	// 0x1009a2614: add w11, w11, w4
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x1009a2618 {
		return 0x1009a2618, true, nil
	}
	// 0x1009a2618: add w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1009a261c {
		return 0x1009a261c, true, nil
	}
	// 0x1009a261c: cmp w27, #0
	s.compare32(uint32(uint64(uint32(s.x[27]))), uint32(uint64(0x0)))
	if runtime != nil && runtime.stopPC == 0x1009a2620 {
		return 0x1009a2620, true, nil
	}
	// 0x1009a2620: cset w15, eq
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	if runtime != nil && runtime.stopPC == 0x1009a2624 {
		return 0x1009a2624, true, nil
	}
	// 0x1009a2624: lsl w17, w15, #9
	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[15])), int(uint64(0x9)), 32)))
	if runtime != nil && runtime.stopPC == 0x1009a2628 {
		return 0x1009a2628, true, nil
	}
	// 0x1009a2628: sub w15, w17, w15, lsl #5
	s.x[15] = uint64(uint32((uint64(uint32(s.x[17])) - shiftLeft(uint64(uint32(s.x[15])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x1009a262c {
		return 0x1009a262c, true, nil
	}
	// 0x1009a262c: eor w15, w15, w7
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x1009a2630 {
		return 0x1009a2630, true, nil
	}
	// 0x1009a2630: ldr x15, [x13, w15, sxtw #3]
	value459, err := s.load((s.x[13] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[15])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(value459)
	if runtime != nil && runtime.stopPC == 0x1009a2634 {
		return 0x1009a2634, true, nil
	}
	// 0x1009a2634: mov x17, #0xffbe
	s.x[17] = uint64(uint64(0xffbe))
	if runtime != nil && runtime.stopPC == 0x1009a2638 {
		return 0x1009a2638, true, nil
	}
	// 0x1009a2638: braa x15, x17
	next = s.x[15]
	return next, false, nil
}

func rawSourceBootstrapBlock100a1278c(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100a1278c {
		return 0x100a1278c, true, nil
	}
	// 0x100a1278c: mov w10, #-0x2899
	s.x[10] = uint64(uint32(uint64(0xffffd767)))
	if runtime != nil && runtime.stopPC == 0x100a12790 {
		return 0x100a12790, true, nil
	}
	// 0x100a12790: add w10, w7, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[7])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100a12794 {
		return 0x100a12794, true, nil
	}
	// 0x100a12794: mov w15, #0x2899
	s.x[15] = uint64(uint32(uint64(0x2899)))
	if runtime != nil && runtime.stopPC == 0x100a12798 {
		return 0x100a12798, true, nil
	}
	// 0x100a12798: add w7, w10, w15
	s.x[7] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a1279c {
		return 0x100a1279c, true, nil
	}
	// 0x100a1279c: mov w15, #0x8a06
	s.x[15] = uint64(uint32(uint64(0x8a06)))
	if runtime != nil && runtime.stopPC == 0x100a127a0 {
		return 0x100a127a0, true, nil
	}
	// 0x100a127a0: movk w15, #0x52c0, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x52c0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127a4 {
		return 0x100a127a4, true, nil
	}
	// 0x100a127a4: add w15, w10, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a127a8 {
		return 0x100a127a8, true, nil
	}
	// 0x100a127a8: eor w15, w15, w19
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[19])))))
	if runtime != nil && runtime.stopPC == 0x100a127ac {
		return 0x100a127ac, true, nil
	}
	// 0x100a127ac: mov w16, #0xc3d0
	s.x[16] = uint64(uint32(uint64(0xc3d0)))
	if runtime != nil && runtime.stopPC == 0x100a127b0 {
		return 0x100a127b0, true, nil
	}
	// 0x100a127b0: movk w16, #0x54e5, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x54e5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127b4 {
		return 0x100a127b4, true, nil
	}
	// 0x100a127b4: eor w17, w19, w16
	s.x[17] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a127b8 {
		return 0x100a127b8, true, nil
	}
	// 0x100a127b8: mov w16, #0x3c30
	s.x[16] = uint64(uint32(uint64(0x3c30)))
	if runtime != nil && runtime.stopPC == 0x100a127bc {
		return 0x100a127bc, true, nil
	}
	// 0x100a127bc: movk w16, #0xab1a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xab1a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127c0 {
		return 0x100a127c0, true, nil
	}
	// 0x100a127c0: add w17, w17, w16
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a127c4 {
		return 0x100a127c4, true, nil
	}
	// 0x100a127c4: mov w16, #0xbedc
	s.x[16] = uint64(uint32(uint64(0xbedc)))
	if runtime != nil && runtime.stopPC == 0x100a127c8 {
		return 0x100a127c8, true, nil
	}
	// 0x100a127c8: movk w16, #0xd006, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xd006, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127cc {
		return 0x100a127cc, true, nil
	}
	// 0x100a127cc: eor w0, w19, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a127d0 {
		return 0x100a127d0, true, nil
	}
	// 0x100a127d0: mov w16, #0x4124
	s.x[16] = uint64(uint32(uint64(0x4124)))
	if runtime != nil && runtime.stopPC == 0x100a127d4 {
		return 0x100a127d4, true, nil
	}
	// 0x100a127d4: movk w16, #0x2ff9, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x2ff9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127d8 {
		return 0x100a127d8, true, nil
	}
	// 0x100a127d8: add w0, w0, w16
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a127dc {
		return 0x100a127dc, true, nil
	}
	// 0x100a127dc: eor w17, w0, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100a127e0 {
		return 0x100a127e0, true, nil
	}
	// 0x100a127e0: mov w16, #0x794
	s.x[16] = uint64(uint32(uint64(0x794)))
	if runtime != nil && runtime.stopPC == 0x100a127e4 {
		return 0x100a127e4, true, nil
	}
	// 0x100a127e4: movk w16, #0xad3f, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xad3f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127e8 {
		return 0x100a127e8, true, nil
	}
	// 0x100a127e8: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a127ec {
		return 0x100a127ec, true, nil
	}
	// 0x100a127ec: eor w15, w17, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a127f0 {
		return 0x100a127f0, true, nil
	}
	// 0x100a127f0: mov w16, #0x3260
	s.x[16] = uint64(uint32(uint64(0x3260)))
	if runtime != nil && runtime.stopPC == 0x100a127f4 {
		return 0x100a127f4, true, nil
	}
	// 0x100a127f4: movk w16, #0xa5da, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa5da, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a127f8 {
		return 0x100a127f8, true, nil
	}
	// 0x100a127f8: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a127fc {
		return 0x100a127fc, true, nil
	}
	// 0x100a127fc: add w24, w11, w15
	s.x[24] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a12800 {
		return 0x100a12800, true, nil
	}
	// 0x100a12800: ldr x10, [x13, w10, sxtw #3]
	value29, err := s.load((s.x[13] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value29)
	if runtime != nil && runtime.stopPC == 0x100a12804 {
		return 0x100a12804, true, nil
	}
	// 0x100a12804: add w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100a12808 {
		return 0x100a12808, true, nil
	}
	// 0x100a12808: mov w11, #0xaf2f
	s.x[11] = uint64(uint32(uint64(0xaf2f)))
	if runtime != nil && runtime.stopPC == 0x100a1280c {
		return 0x100a1280c, true, nil
	}
	// 0x100a1280c: movk w11, #0xd7a8, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd7a8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100a12810 {
		return 0x100a12810, true, nil
	}
	// 0x100a12810: add w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100a12814 {
		return 0x100a12814, true, nil
	}
	// 0x100a12814: sub w27, w27, #1
	s.x[27] = uint64(uint32((uint64(uint32(s.x[27])) - uint64(0x1))))
	if runtime != nil && runtime.stopPC == 0x100a12818 {
		return 0x100a12818, true, nil
	}
	// 0x100a12818: add w28, w28, #0x86
	s.x[28] = uint64(uint32((uint64(uint32(s.x[28])) + uint64(0x86))))
	if runtime != nil && runtime.stopPC == 0x100a1281c {
		return 0x100a1281c, true, nil
	}
	// 0x100a1281c: add w6, w6, #0x86
	s.x[6] = uint64(uint32((uint64(uint32(s.x[6])) + uint64(0x86))))
	if runtime != nil && runtime.stopPC == 0x100a12820 {
		return 0x100a12820, true, nil
	}
	// 0x100a12820: mov x19, x12
	s.x[19] = uint64(s.x[12])
	if runtime != nil && runtime.stopPC == 0x100a12824 {
		return 0x100a12824, true, nil
	}
	// 0x100a12824: mov x12, x3
	s.x[12] = uint64(s.x[3])
	if runtime != nil && runtime.stopPC == 0x100a12828 {
		return 0x100a12828, true, nil
	}
	// 0x100a12828: mov x17, #0x9e6d
	s.x[17] = uint64(uint64(0x9e6d))
	if runtime != nil && runtime.stopPC == 0x100a1282c {
		return 0x100a1282c, true, nil
	}
	// 0x100a1282c: braa x10, x17
	next = s.x[10]
	return next, false, nil
}

func rawSourceBootstrapBlock100a42a30(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100a42a30 {
		return 0x100a42a30, true, nil
	}
	// 0x100a42a30: mul w12, w5, w10
	s.x[12] = uint64(uint32((uint64(uint32(s.x[5])) * uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100a42a34 {
		return 0x100a42a34, true, nil
	}
	// 0x100a42a34: and w17, w12, #0xfffc
	s.x[17] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xfffc))))
	if runtime != nil && runtime.stopPC == 0x100a42a38 {
		return 0x100a42a38, true, nil
	}
	// 0x100a42a38: mul w17, w17, w14
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100a42a3c {
		return 0x100a42a3c, true, nil
	}
	// 0x100a42a3c: lsr w17, w17, #0x16
	s.x[17] = uint64(uint32(shiftRight(uint64(uint32(s.x[17])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a42a40 {
		return 0x100a42a40, true, nil
	}
	// 0x100a42a40: msub w17, w17, w15, w12
	s.x[17] = uint64(uint32((uint64(uint32(s.x[12])) - uint64(uint32(s.x[17]))*uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a42a44 {
		return 0x100a42a44, true, nil
	}
	// 0x100a42a44: add x17, x22, w17, uxth
	s.x[17] = uint64((s.x[22] + uint64(uint16(uint64(uint32(s.x[17]))))))
	if runtime != nil && runtime.stopPC == 0x100a42a48 {
		return 0x100a42a48, true, nil
	}
	// 0x100a42a48: ldrb w17, [x17]
	value6, err := s.load(s.x[17], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value6))
	if runtime != nil && runtime.stopPC == 0x100a42a4c {
		return 0x100a42a4c, true, nil
	}
	// 0x100a42a4c: add w6, w12, #0x25
	s.x[6] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0x25))))
	if runtime != nil && runtime.stopPC == 0x100a42a50 {
		return 0x100a42a50, true, nil
	}
	// 0x100a42a50: and w7, w6, #0xffff
	s.x[7] = uint64(uint32((uint64(uint32(s.x[6])) & uint64(0xffff))))
	if runtime != nil && runtime.stopPC == 0x100a42a54 {
		return 0x100a42a54, true, nil
	}
	// 0x100a42a54: mul w7, w7, w14
	s.x[7] = uint64(uint32((uint64(uint32(s.x[7])) * uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100a42a58 {
		return 0x100a42a58, true, nil
	}
	// 0x100a42a58: lsr w7, w7, #0x16
	s.x[7] = uint64(uint32(shiftRight(uint64(uint32(s.x[7])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a42a5c {
		return 0x100a42a5c, true, nil
	}
	// 0x100a42a5c: msub w6, w7, w15, w6
	s.x[6] = uint64(uint32((uint64(uint32(s.x[6])) - uint64(uint32(s.x[7]))*uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a42a60 {
		return 0x100a42a60, true, nil
	}
	// 0x100a42a60: add w7, w12, #0x4a
	s.x[7] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0x4a))))
	if runtime != nil && runtime.stopPC == 0x100a42a64 {
		return 0x100a42a64, true, nil
	}
	// 0x100a42a64: and w19, w7, #0xfffe
	s.x[19] = uint64(uint32((uint64(uint32(s.x[7])) & uint64(0xfffe))))
	if runtime != nil && runtime.stopPC == 0x100a42a68 {
		return 0x100a42a68, true, nil
	}
	// 0x100a42a68: mul w19, w19, w14
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) * uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100a42a6c {
		return 0x100a42a6c, true, nil
	}
	// 0x100a42a6c: lsr w19, w19, #0x16
	s.x[19] = uint64(uint32(shiftRight(uint64(uint32(s.x[19])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a42a70 {
		return 0x100a42a70, true, nil
	}
	// 0x100a42a70: msub w7, w19, w15, w7
	s.x[7] = uint64(uint32((uint64(uint32(s.x[7])) - uint64(uint32(s.x[19]))*uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a42a74 {
		return 0x100a42a74, true, nil
	}
	// 0x100a42a74: madd w19, w5, w0, w1
	s.x[19] = uint64(uint32((uint64(uint32(s.x[5]))*uint64(uint32(s.x[0])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100a42a78 {
		return 0x100a42a78, true, nil
	}
	// 0x100a42a78: add w5, w5, w16
	s.x[5] = uint64(uint32((uint64(uint32(s.x[5])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100a42a7c {
		return 0x100a42a7c, true, nil
	}
	// 0x100a42a7c: add w5, w5, w11
	s.x[5] = uint64(uint32((uint64(uint32(s.x[5])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100a42a80 {
		return 0x100a42a80, true, nil
	}
	// 0x100a42a80: add x6, x22, w6, uxth
	s.x[6] = uint64((s.x[22] + uint64(uint16(uint64(uint32(s.x[6]))))))
	if runtime != nil && runtime.stopPC == 0x100a42a84 {
		return 0x100a42a84, true, nil
	}
	// 0x100a42a84: ldrb w6, [x6]
	value21, err := s.load(s.x[6], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(uint32(value21))
	if runtime != nil && runtime.stopPC == 0x100a42a88 {
		return 0x100a42a88, true, nil
	}
	// 0x100a42a88: add w12, w12, #0x6f
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0x6f))))
	if runtime != nil && runtime.stopPC == 0x100a42a8c {
		return 0x100a42a8c, true, nil
	}
	// 0x100a42a8c: and w20, w12, #0xffff
	s.x[20] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xffff))))
	if runtime != nil && runtime.stopPC == 0x100a42a90 {
		return 0x100a42a90, true, nil
	}
	// 0x100a42a90: mul w20, w20, w14
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) * uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100a42a94 {
		return 0x100a42a94, true, nil
	}
	// 0x100a42a94: lsr w20, w20, #0x16
	s.x[20] = uint64(uint32(shiftRight(uint64(uint32(s.x[20])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a42a98 {
		return 0x100a42a98, true, nil
	}
	// 0x100a42a98: msub w12, w20, w15, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) - uint64(uint32(s.x[20]))*uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100a42a9c {
		return 0x100a42a9c, true, nil
	}
	// 0x100a42a9c: lsl w17, w17, #0x10
	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[17])), int(uint64(0x10)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a42aa0 {
		return 0x100a42aa0, true, nil
	}
	// 0x100a42aa0: orr w17, w17, w6, lsl #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) | shiftLeft(uint64(uint32(s.x[6])), 8, 32))))
	if runtime != nil && runtime.stopPC == 0x100a42aa4 {
		return 0x100a42aa4, true, nil
	}
	// 0x100a42aa4: add x6, x22, w7, uxth
	s.x[6] = uint64((s.x[22] + uint64(uint16(uint64(uint32(s.x[7]))))))
	if runtime != nil && runtime.stopPC == 0x100a42aa8 {
		return 0x100a42aa8, true, nil
	}
	// 0x100a42aa8: ldrb w6, [x6]
	value30, err := s.load(s.x[6], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(uint32(value30))
	if runtime != nil && runtime.stopPC == 0x100a42aac {
		return 0x100a42aac, true, nil
	}
	// 0x100a42aac: orr w17, w17, w6
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) | uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100a42ab0 {
		return 0x100a42ab0, true, nil
	}
	// 0x100a42ab0: add x12, x22, w12, uxth
	s.x[12] = uint64((s.x[22] + uint64(uint16(uint64(uint32(s.x[12]))))))
	if runtime != nil && runtime.stopPC == 0x100a42ab4 {
		return 0x100a42ab4, true, nil
	}
	// 0x100a42ab4: ldrb w12, [x12]
	value33, err := s.load(s.x[12], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value33))
	if runtime != nil && runtime.stopPC == 0x100a42ab8 {
		return 0x100a42ab8, true, nil
	}
	// 0x100a42ab8: orr w12, w12, w17, lsl #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) | shiftLeft(uint64(uint32(s.x[17])), 8, 32))))
	if runtime != nil && runtime.stopPC == 0x100a42abc {
		return 0x100a42abc, true, nil
	}
	// 0x100a42abc: and w17, w19, #0xfffe
	s.x[17] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(0xfffe))))
	if runtime != nil && runtime.stopPC == 0x100a42ac0 {
		return 0x100a42ac0, true, nil
	}
	// 0x100a42ac0: mul w17, w17, w2
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[2])))))
	if runtime != nil && runtime.stopPC == 0x100a42ac4 {
		return 0x100a42ac4, true, nil
	}
	// 0x100a42ac4: lsr w17, w17, #0x16
	s.x[17] = uint64(uint32(shiftRight(uint64(uint32(s.x[17])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100a42ac8 {
		return 0x100a42ac8, true, nil
	}
	// 0x100a42ac8: msub w17, w17, w3, w19
	s.x[17] = uint64(uint32((uint64(uint32(s.x[19])) - uint64(uint32(s.x[17]))*uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100a42acc {
		return 0x100a42acc, true, nil
	}
	// 0x100a42acc: ubfiz x17, x17, #2, #0x10
	s.x[17] = uint64(((s.x[17]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100a42ad0 {
		return 0x100a42ad0, true, nil
	}
	// 0x100a42ad0: str w12, [x8, x17]
	if err := s.storeValue((s.x[8] + s.x[17]), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100a42ad4 {
		return 0x100a42ad4, true, nil
	}
	// 0x100a42ad4: cmp w5, #0x10
	s.compare32(uint32(uint64(uint32(s.x[5]))), uint32(uint64(0x10)))
	if runtime != nil && runtime.stopPC == 0x100a42ad8 {
		return 0x100a42ad8, true, nil
	}
	// 0x100a42ad8: cset w12, eq
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	if runtime != nil && runtime.stopPC == 0x100a42adc {
		return 0x100a42adc, true, nil
	}
	// 0x100a42adc: mul w12, w12, w4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x100a42ae0 {
		return 0x100a42ae0, true, nil
	}
	// 0x100a42ae0: eor w12, w12, w9
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100a42ae4 {
		return 0x100a42ae4, true, nil
	}
	// 0x100a42ae4: ldr x6, [x30, w12, sxtw #3]
	value45, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(value45)
	if runtime != nil && runtime.stopPC == 0x100a42ae8 {
		return 0x100a42ae8, true, nil
	}
	// 0x100a42ae8: mov w12, #0x10
	s.x[12] = uint64(uint32(uint64(0x10)))
	if runtime != nil && runtime.stopPC == 0x100a42aec {
		return 0x100a42aec, true, nil
	}
	// 0x100a42aec: mov x17, #0x3e0f
	s.x[17] = uint64(uint64(0x3e0f))
	if runtime != nil && runtime.stopPC == 0x100a42af0 {
		return 0x100a42af0, true, nil
	}
	// 0x100a42af0: braa x6, x17
	next = s.x[6]
	return next, false, nil
}

func rawSourceBootstrapBlock100aa281c(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100aa281c {
		return 0x100aa281c, true, nil
	}
	// 0x100aa281c: mov w15, #0x8bea
	s.x[15] = uint64(uint32(uint64(0x8bea)))
	if runtime != nil && runtime.stopPC == 0x100aa2820 {
		return 0x100aa2820, true, nil
	}
	// 0x100aa2820: eor w15, w22, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[22])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2824 {
		return 0x100aa2824, true, nil
	}
	// 0x100aa2824: mov w16, #0x3d
	s.x[16] = uint64(uint32(uint64(0x3d)))
	if runtime != nil && runtime.stopPC == 0x100aa2828 {
		return 0x100aa2828, true, nil
	}
	// 0x100aa2828: mul w5, w15, w16
	s.x[5] = uint64(uint32((uint64(uint32(s.x[15])) * uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa282c {
		return 0x100aa282c, true, nil
	}
	// 0x100aa282c: ldr w20, [sp, #0x118]
	value4, err := s.load(addOffset(s.sp, 280), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[20] = uint64(uint32(value4))
	if runtime != nil && runtime.stopPC == 0x100aa2830 {
		return 0x100aa2830, true, nil
	}
	// 0x100aa2830: sub w20, w20, #0x40
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) - uint64(0x40))))
	if runtime != nil && runtime.stopPC == 0x100aa2834 {
		return 0x100aa2834, true, nil
	}
	// 0x100aa2834: ldr w15, [sp, #0x18c]
	value6, err := s.load(addOffset(s.sp, 396), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value6))
	if runtime != nil && runtime.stopPC == 0x100aa2838 {
		return 0x100aa2838, true, nil
	}
	// 0x100aa2838: add w15, w15, w11
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa283c {
		return 0x100aa283c, true, nil
	}
	// 0x100aa283c: add w15, w15, w5
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100aa2840 {
		return 0x100aa2840, true, nil
	}
	// 0x100aa2840: mov w16, #0x807
	s.x[16] = uint64(uint32(uint64(0x807)))
	if runtime != nil && runtime.stopPC == 0x100aa2844 {
		return 0x100aa2844, true, nil
	}
	// 0x100aa2844: movk w16, #0xdd90, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xdd90, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2848 {
		return 0x100aa2848, true, nil
	}
	// 0x100aa2848: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa284c {
		return 0x100aa284c, true, nil
	}
	// 0x100aa284c: add w4, w15, w9
	s.x[4] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2850 {
		return 0x100aa2850, true, nil
	}
	// 0x100aa2850: mov w9, #0x8341
	s.x[9] = uint64(uint32(uint64(0x8341)))
	if runtime != nil && runtime.stopPC == 0x100aa2854 {
		return 0x100aa2854, true, nil
	}
	// 0x100aa2854: movk w9, #0x2b25, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x2b25, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2858 {
		return 0x100aa2858, true, nil
	}
	// 0x100aa2858: eor w9, w19, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa285c {
		return 0x100aa285c, true, nil
	}
	// 0x100aa285c: mov w15, #0x7cbf
	s.x[15] = uint64(uint32(uint64(0x7cbf)))
	if runtime != nil && runtime.stopPC == 0x100aa2860 {
		return 0x100aa2860, true, nil
	}
	// 0x100aa2860: movk w15, #0xd4da, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xd4da, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2864 {
		return 0x100aa2864, true, nil
	}
	// 0x100aa2864: add w9, w9, w15
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2868 {
		return 0x100aa2868, true, nil
	}
	// 0x100aa2868: mov w15, #0x161d
	s.x[15] = uint64(uint32(uint64(0x161d)))
	if runtime != nil && runtime.stopPC == 0x100aa286c {
		return 0x100aa286c, true, nil
	}
	// 0x100aa286c: movk w15, #0xbe23, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xbe23, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2870 {
		return 0x100aa2870, true, nil
	}
	// 0x100aa2870: eor w15, w19, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2874 {
		return 0x100aa2874, true, nil
	}
	// 0x100aa2874: mov w16, #0xe9e3
	s.x[16] = uint64(uint32(uint64(0xe9e3)))
	if runtime != nil && runtime.stopPC == 0x100aa2878 {
		return 0x100aa2878, true, nil
	}
	// 0x100aa2878: movk w16, #0x41dc, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x41dc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa287c {
		return 0x100aa287c, true, nil
	}
	// 0x100aa287c: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2880 {
		return 0x100aa2880, true, nil
	}
	// 0x100aa2880: eor w9, w9, w15
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2884 {
		return 0x100aa2884, true, nil
	}
	// 0x100aa2884: mov w15, #0x103c
	s.x[15] = uint64(uint32(uint64(0x103c)))
	if runtime != nil && runtime.stopPC == 0x100aa2888 {
		return 0x100aa2888, true, nil
	}
	// 0x100aa2888: movk w15, #0x4325, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4325, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa288c {
		return 0x100aa288c, true, nil
	}
	// 0x100aa288c: eor w15, w19, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2890 {
		return 0x100aa2890, true, nil
	}
	// 0x100aa2890: mov w16, #0xefc4
	s.x[16] = uint64(uint32(uint64(0xefc4)))
	if runtime != nil && runtime.stopPC == 0x100aa2894 {
		return 0x100aa2894, true, nil
	}
	// 0x100aa2894: movk w16, #0xbcda, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xbcda, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2898 {
		return 0x100aa2898, true, nil
	}
	// 0x100aa2898: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa289c {
		return 0x100aa289c, true, nil
	}
	// 0x100aa289c: eor w9, w9, w15
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa28a0 {
		return 0x100aa28a0, true, nil
	}
	// 0x100aa28a0: mov w15, #0xd5c8
	s.x[15] = uint64(uint32(uint64(0xd5c8)))
	if runtime != nil && runtime.stopPC == 0x100aa28a4 {
		return 0x100aa28a4, true, nil
	}
	// 0x100aa28a4: movk w15, #0x3074, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x3074, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28a8 {
		return 0x100aa28a8, true, nil
	}
	// 0x100aa28a8: eor w15, w12, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa28ac {
		return 0x100aa28ac, true, nil
	}
	// 0x100aa28ac: mov w16, #0x2a38
	s.x[16] = uint64(uint32(uint64(0x2a38)))
	if runtime != nil && runtime.stopPC == 0x100aa28b0 {
		return 0x100aa28b0, true, nil
	}
	// 0x100aa28b0: movk w16, #0xcf8b, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xcf8b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28b4 {
		return 0x100aa28b4, true, nil
	}
	// 0x100aa28b4: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa28b8 {
		return 0x100aa28b8, true, nil
	}
	// 0x100aa28b8: mov w16, #0x55ce
	s.x[16] = uint64(uint32(uint64(0x55ce)))
	if runtime != nil && runtime.stopPC == 0x100aa28bc {
		return 0x100aa28bc, true, nil
	}
	// 0x100aa28bc: movk w16, #0x28c2, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x28c2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28c0 {
		return 0x100aa28c0, true, nil
	}
	// 0x100aa28c0: eor w16, w12, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa28c4 {
		return 0x100aa28c4, true, nil
	}
	// 0x100aa28c4: mov w17, #0xaa32
	s.x[17] = uint64(uint32(uint64(0xaa32)))
	if runtime != nil && runtime.stopPC == 0x100aa28c8 {
		return 0x100aa28c8, true, nil
	}
	// 0x100aa28c8: movk w17, #0xd73d, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xd73d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28cc {
		return 0x100aa28cc, true, nil
	}
	// 0x100aa28cc: add w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa28d0 {
		return 0x100aa28d0, true, nil
	}
	// 0x100aa28d0: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa28d4 {
		return 0x100aa28d4, true, nil
	}
	// 0x100aa28d4: mov w16, #0x566
	s.x[16] = uint64(uint32(uint64(0x566)))
	if runtime != nil && runtime.stopPC == 0x100aa28d8 {
		return 0x100aa28d8, true, nil
	}
	// 0x100aa28d8: movk w16, #0xce95, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xce95, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28dc {
		return 0x100aa28dc, true, nil
	}
	// 0x100aa28dc: eor w12, w12, w16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa28e0 {
		return 0x100aa28e0, true, nil
	}
	// 0x100aa28e0: mov w16, #0xfa9a
	s.x[16] = uint64(uint32(uint64(0xfa9a)))
	if runtime != nil && runtime.stopPC == 0x100aa28e4 {
		return 0x100aa28e4, true, nil
	}
	// 0x100aa28e4: movk w16, #0x316a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x316a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28e8 {
		return 0x100aa28e8, true, nil
	}
	// 0x100aa28e8: add w12, w12, w16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa28ec {
		return 0x100aa28ec, true, nil
	}
	// 0x100aa28ec: eor w12, w15, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa28f0 {
		return 0x100aa28f0, true, nil
	}
	// 0x100aa28f0: mov w15, #0xc7b1
	s.x[15] = uint64(uint32(uint64(0xc7b1)))
	if runtime != nil && runtime.stopPC == 0x100aa28f4 {
		return 0x100aa28f4, true, nil
	}
	// 0x100aa28f4: movk w15, #0xe011, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xe011, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa28f8 {
		return 0x100aa28f8, true, nil
	}
	// 0x100aa28f8: ldr w0, [sp, #0x188]
	value55, err := s.load(addOffset(s.sp, 392), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[0] = uint64(uint32(value55))
	if runtime != nil && runtime.stopPC == 0x100aa28fc {
		return 0x100aa28fc, true, nil
	}
	// 0x100aa28fc: eor w15, w0, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2900 {
		return 0x100aa2900, true, nil
	}
	// 0x100aa2900: mov w16, #0x384f
	s.x[16] = uint64(uint32(uint64(0x384f)))
	if runtime != nil && runtime.stopPC == 0x100aa2904 {
		return 0x100aa2904, true, nil
	}
	// 0x100aa2904: movk w16, #0x1fee, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x1fee, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2908 {
		return 0x100aa2908, true, nil
	}
	// 0x100aa2908: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa290c {
		return 0x100aa290c, true, nil
	}
	// 0x100aa290c: mov w16, #0xea9c
	s.x[16] = uint64(uint32(uint64(0xea9c)))
	if runtime != nil && runtime.stopPC == 0x100aa2910 {
		return 0x100aa2910, true, nil
	}
	// 0x100aa2910: movk w16, #0x45de, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x45de, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2914 {
		return 0x100aa2914, true, nil
	}
	// 0x100aa2914: eor w16, w0, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2918 {
		return 0x100aa2918, true, nil
	}
	// 0x100aa2918: mov w17, #0x1564
	s.x[17] = uint64(uint32(uint64(0x1564)))
	if runtime != nil && runtime.stopPC == 0x100aa291c {
		return 0x100aa291c, true, nil
	}
	// 0x100aa291c: movk w17, #0xba21, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xba21, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2920 {
		return 0x100aa2920, true, nil
	}
	// 0x100aa2920: add w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2924 {
		return 0x100aa2924, true, nil
	}
	// 0x100aa2924: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2928 {
		return 0x100aa2928, true, nil
	}
	// 0x100aa2928: mov w16, #0xa84d
	s.x[16] = uint64(uint32(uint64(0xa84d)))
	if runtime != nil && runtime.stopPC == 0x100aa292c {
		return 0x100aa292c, true, nil
	}
	// 0x100aa292c: movk w16, #0x73ec, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x73ec, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2930 {
		return 0x100aa2930, true, nil
	}
	// 0x100aa2930: eor w16, w0, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2934 {
		return 0x100aa2934, true, nil
	}
	// 0x100aa2934: mov w17, #0x57b3
	s.x[17] = uint64(uint32(uint64(0x57b3)))
	if runtime != nil && runtime.stopPC == 0x100aa2938 {
		return 0x100aa2938, true, nil
	}
	// 0x100aa2938: movk w17, #0x8c13, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x8c13, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa293c {
		return 0x100aa293c, true, nil
	}
	// 0x100aa293c: add w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2940 {
		return 0x100aa2940, true, nil
	}
	// 0x100aa2940: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2944 {
		return 0x100aa2944, true, nil
	}
	// 0x100aa2944: mov w16, #0xf2dd
	s.x[16] = uint64(uint32(uint64(0xf2dd)))
	if runtime != nil && runtime.stopPC == 0x100aa2948 {
		return 0x100aa2948, true, nil
	}
	// 0x100aa2948: movk w16, #0xdbba, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xdbba, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa294c {
		return 0x100aa294c, true, nil
	}
	// 0x100aa294c: add w12, w15, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2950 {
		return 0x100aa2950, true, nil
	}
	// 0x100aa2950: add w12, w12, w16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2954 {
		return 0x100aa2954, true, nil
	}
	// 0x100aa2954: mov w15, #0x1fe3
	s.x[15] = uint64(uint32(uint64(0x1fe3)))
	if runtime != nil && runtime.stopPC == 0x100aa2958 {
		return 0x100aa2958, true, nil
	}
	// 0x100aa2958: movk w15, #0x119, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x119, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa295c {
		return 0x100aa295c, true, nil
	}
	// 0x100aa295c: eor w15, w12, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2960 {
		return 0x100aa2960, true, nil
	}
	// 0x100aa2960: mov w16, #0x802
	s.x[16] = uint64(uint32(uint64(0x802)))
	if runtime != nil && runtime.stopPC == 0x100aa2964 {
		return 0x100aa2964, true, nil
	}
	// 0x100aa2964: movk w16, #0xd195, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xd195, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2968 {
		return 0x100aa2968, true, nil
	}
	// 0x100aa2968: mov w17, #0x7f0
	s.x[17] = uint64(uint32(uint64(0x7f0)))
	if runtime != nil && runtime.stopPC == 0x100aa296c {
		return 0x100aa296c, true, nil
	}
	// 0x100aa296c: movk w17, #0xe2e3, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xe2e3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2970 {
		return 0x100aa2970, true, nil
	}
	// 0x100aa2970: eor w17, w12, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2974 {
		return 0x100aa2974, true, nil
	}
	// 0x100aa2974: mov w0, #0x1013
	s.x[0] = uint64(uint32(uint64(0x1013)))
	if runtime != nil && runtime.stopPC == 0x100aa2978 {
		return 0x100aa2978, true, nil
	}
	// 0x100aa2978: movk w0, #0x326f, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x326f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa297c {
		return 0x100aa297c, true, nil
	}
	// 0x100aa297c: add w17, w17, w0
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100aa2980 {
		return 0x100aa2980, true, nil
	}
	// 0x100aa2980: mov w0, #0x7f5
	s.x[0] = uint64(uint32(uint64(0x7f5)))
	if runtime != nil && runtime.stopPC == 0x100aa2984 {
		return 0x100aa2984, true, nil
	}
	// 0x100aa2984: movk w0, #0xf334, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xf334, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2988 {
		return 0x100aa2988, true, nil
	}
	// 0x100aa2988: eor w0, w12, w0
	s.x[0] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100aa298c {
		return 0x100aa298c, true, nil
	}
	// 0x100aa298c: mov w1, #0x1018
	s.x[1] = uint64(uint32(uint64(0x1018)))
	if runtime != nil && runtime.stopPC == 0x100aa2990 {
		return 0x100aa2990, true, nil
	}
	// 0x100aa2990: movk w1, #0x23b8, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x23b8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2994 {
		return 0x100aa2994, true, nil
	}
	// 0x100aa2994: mov w2, #0xf7fb
	s.x[2] = uint64(uint32(uint64(0xf7fb)))
	if runtime != nil && runtime.stopPC == 0x100aa2998 {
		return 0x100aa2998, true, nil
	}
	// 0x100aa2998: movk w2, #0x3fbd, lsl #16
	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0x3fbd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa299c {
		return 0x100aa299c, true, nil
	}
	// 0x100aa299c: eor w2, w12, w2
	s.x[2] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[2])))))
	if runtime != nil && runtime.stopPC == 0x100aa29a0 {
		return 0x100aa29a0, true, nil
	}
	// 0x100aa29a0: add w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa29a4 {
		return 0x100aa29a4, true, nil
	}
	// 0x100aa29a4: eor w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa29a8 {
		return 0x100aa29a8, true, nil
	}
	// 0x100aa29a8: mov w15, #0xe01a
	s.x[15] = uint64(uint32(uint64(0xe01a)))
	if runtime != nil && runtime.stopPC == 0x100aa29ac {
		return 0x100aa29ac, true, nil
	}
	// 0x100aa29ac: movk w15, #0xef31, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xef31, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa29b0 {
		return 0x100aa29b0, true, nil
	}
	// 0x100aa29b0: add w16, w0, w1
	s.x[16] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100aa29b4 {
		return 0x100aa29b4, true, nil
	}
	// 0x100aa29b4: eor w16, w17, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa29b8 {
		return 0x100aa29b8, true, nil
	}
	// 0x100aa29b8: eor w12, w12, w16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa29bc {
		return 0x100aa29bc, true, nil
	}
	// 0x100aa29bc: mov w16, #0x6d7d
	s.x[16] = uint64(uint32(uint64(0x6d7d)))
	if runtime != nil && runtime.stopPC == 0x100aa29c0 {
		return 0x100aa29c0, true, nil
	}
	// 0x100aa29c0: movk w16, #0xf950, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xf950, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa29c4 {
		return 0x100aa29c4, true, nil
	}
	// 0x100aa29c4: add w15, w2, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[2])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa29c8 {
		return 0x100aa29c8, true, nil
	}
	// 0x100aa29c8: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa29cc {
		return 0x100aa29cc, true, nil
	}
	// 0x100aa29cc: eor w19, w12, w15
	s.x[19] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa29d0 {
		return 0x100aa29d0, true, nil
	}
	// 0x100aa29d0: mov w12, #0xe988
	s.x[12] = uint64(uint32(uint64(0xe988)))
	if runtime != nil && runtime.stopPC == 0x100aa29d4 {
		return 0x100aa29d4, true, nil
	}
	// 0x100aa29d4: movk w12, #0x815c, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x815c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa29d8 {
		return 0x100aa29d8, true, nil
	}
	// 0x100aa29d8: ldr w15, [sp, #0x198]
	value111, err := s.load(addOffset(s.sp, 408), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value111))
	if runtime != nil && runtime.stopPC == 0x100aa29dc {
		return 0x100aa29dc, true, nil
	}
	// 0x100aa29dc: add w9, w15, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa29e0 {
		return 0x100aa29e0, true, nil
	}
	// 0x100aa29e0: add w9, w9, w12
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa29e4 {
		return 0x100aa29e4, true, nil
	}
	// 0x100aa29e4: add w26, w9, w11
	s.x[26] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa29e8 {
		return 0x100aa29e8, true, nil
	}
	// 0x100aa29e8: mov w9, #0xf345
	s.x[9] = uint64(uint32(uint64(0xf345)))
	if runtime != nil && runtime.stopPC == 0x100aa29ec {
		return 0x100aa29ec, true, nil
	}
	// 0x100aa29ec: movk w9, #0x6c94, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x6c94, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa29f0 {
		return 0x100aa29f0, true, nil
	}
	// 0x100aa29f0: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa29f4 {
		return 0x100aa29f4, true, nil
	}
	// 0x100aa29f4: mov w11, #0xdbc7
	s.x[11] = uint64(uint32(uint64(0xdbc7)))
	if runtime != nil && runtime.stopPC == 0x100aa29f8 {
		return 0x100aa29f8, true, nil
	}
	// 0x100aa29f8: movk w11, #0x2d2d, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x2d2d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa29fc {
		return 0x100aa29fc, true, nil
	}
	// 0x100aa29fc: add w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a00 {
		return 0x100aa2a00, true, nil
	}
	// 0x100aa2a00: mov w11, #0x912f
	s.x[11] = uint64(uint32(uint64(0x912f)))
	if runtime != nil && runtime.stopPC == 0x100aa2a04 {
		return 0x100aa2a04, true, nil
	}
	// 0x100aa2a04: movk w11, #0xe37f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xe37f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a08 {
		return 0x100aa2a08, true, nil
	}
	// 0x100aa2a08: eor w11, w10, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a0c {
		return 0x100aa2a0c, true, nil
	}
	// 0x100aa2a0c: mov w12, #0xb9ad
	s.x[12] = uint64(uint32(uint64(0xb9ad)))
	if runtime != nil && runtime.stopPC == 0x100aa2a10 {
		return 0x100aa2a10, true, nil
	}
	// 0x100aa2a10: movk w12, #0xa2c6, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xa2c6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a14 {
		return 0x100aa2a14, true, nil
	}
	// 0x100aa2a14: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a18 {
		return 0x100aa2a18, true, nil
	}
	// 0x100aa2a18: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a1c {
		return 0x100aa2a1c, true, nil
	}
	// 0x100aa2a1c: mov w11, #0x3076
	s.x[11] = uint64(uint32(uint64(0x3076)))
	if runtime != nil && runtime.stopPC == 0x100aa2a20 {
		return 0x100aa2a20, true, nil
	}
	// 0x100aa2a20: movk w11, #0xe78e, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xe78e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a24 {
		return 0x100aa2a24, true, nil
	}
	// 0x100aa2a24: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a28 {
		return 0x100aa2a28, true, nil
	}
	// 0x100aa2a28: mov w11, #0x18f6
	s.x[11] = uint64(uint32(uint64(0x18f6)))
	if runtime != nil && runtime.stopPC == 0x100aa2a2c {
		return 0x100aa2a2c, true, nil
	}
	// 0x100aa2a2c: movk w11, #0xa637, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa637, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a30 {
		return 0x100aa2a30, true, nil
	}
	// 0x100aa2a30: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a34 {
		return 0x100aa2a34, true, nil
	}
	// 0x100aa2a34: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a38 {
		return 0x100aa2a38, true, nil
	}
	// 0x100aa2a38: mov w10, #0x1618
	s.x[10] = uint64(uint32(uint64(0x1618)))
	if runtime != nil && runtime.stopPC == 0x100aa2a3c {
		return 0x100aa2a3c, true, nil
	}
	// 0x100aa2a3c: movk w10, #0xd092, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xd092, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a40 {
		return 0x100aa2a40, true, nil
	}
	// 0x100aa2a40: ldr w15, [sp, #0x1a8]
	value137, err := s.load(addOffset(s.sp, 424), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value137))
	if runtime != nil && runtime.stopPC == 0x100aa2a44 {
		return 0x100aa2a44, true, nil
	}
	// 0x100aa2a44: eor w10, w15, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a48 {
		return 0x100aa2a48, true, nil
	}
	// 0x100aa2a48: mov w11, #0xe9e8
	s.x[11] = uint64(uint32(uint64(0xe9e8)))
	if runtime != nil && runtime.stopPC == 0x100aa2a4c {
		return 0x100aa2a4c, true, nil
	}
	// 0x100aa2a4c: movk w11, #0x2f6d, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x2f6d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a50 {
		return 0x100aa2a50, true, nil
	}
	// 0x100aa2a50: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a54 {
		return 0x100aa2a54, true, nil
	}
	// 0x100aa2a54: mov w11, #0x41ec
	s.x[11] = uint64(uint32(uint64(0x41ec)))
	if runtime != nil && runtime.stopPC == 0x100aa2a58 {
		return 0x100aa2a58, true, nil
	}
	// 0x100aa2a58: movk w11, #0x86ad, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x86ad, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a5c {
		return 0x100aa2a5c, true, nil
	}
	// 0x100aa2a5c: eor w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a60 {
		return 0x100aa2a60, true, nil
	}
	// 0x100aa2a60: mov w12, #0xbe14
	s.x[12] = uint64(uint32(uint64(0xbe14)))
	if runtime != nil && runtime.stopPC == 0x100aa2a64 {
		return 0x100aa2a64, true, nil
	}
	// 0x100aa2a64: movk w12, #0x7952, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x7952, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a68 {
		return 0x100aa2a68, true, nil
	}
	// 0x100aa2a68: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a6c {
		return 0x100aa2a6c, true, nil
	}
	// 0x100aa2a6c: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a70 {
		return 0x100aa2a70, true, nil
	}
	// 0x100aa2a70: mov w11, #0xd294
	s.x[11] = uint64(uint32(uint64(0xd294)))
	if runtime != nil && runtime.stopPC == 0x100aa2a74 {
		return 0x100aa2a74, true, nil
	}
	// 0x100aa2a74: movk w11, #0x801c, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x801c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a78 {
		return 0x100aa2a78, true, nil
	}
	// 0x100aa2a78: eor w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a7c {
		return 0x100aa2a7c, true, nil
	}
	// 0x100aa2a7c: mov w12, #0x2d6c
	s.x[12] = uint64(uint32(uint64(0x2d6c)))
	if runtime != nil && runtime.stopPC == 0x100aa2a80 {
		return 0x100aa2a80, true, nil
	}
	// 0x100aa2a80: movk w12, #0x7fe3, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x7fe3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a84 {
		return 0x100aa2a84, true, nil
	}
	// 0x100aa2a84: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a88 {
		return 0x100aa2a88, true, nil
	}
	// 0x100aa2a88: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a8c {
		return 0x100aa2a8c, true, nil
	}
	// 0x100aa2a8c: mov w11, #0x5c49
	s.x[11] = uint64(uint32(uint64(0x5c49)))
	if runtime != nil && runtime.stopPC == 0x100aa2a90 {
		return 0x100aa2a90, true, nil
	}
	// 0x100aa2a90: movk w11, #0x1f0b, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x1f0b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2a94 {
		return 0x100aa2a94, true, nil
	}
	// 0x100aa2a94: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a98 {
		return 0x100aa2a98, true, nil
	}
	// 0x100aa2a98: add w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2a9c {
		return 0x100aa2a9c, true, nil
	}
	// 0x100aa2a9c: mov w10, #0x97aa
	s.x[10] = uint64(uint32(uint64(0x97aa)))
	if runtime != nil && runtime.stopPC == 0x100aa2aa0 {
		return 0x100aa2aa0, true, nil
	}
	// 0x100aa2aa0: movk w10, #0x85d5, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x85d5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2aa4 {
		return 0x100aa2aa4, true, nil
	}
	// 0x100aa2aa4: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2aa8 {
		return 0x100aa2aa8, true, nil
	}
	// 0x100aa2aa8: mov w11, #0x39dd
	s.x[11] = uint64(uint32(uint64(0x39dd)))
	if runtime != nil && runtime.stopPC == 0x100aa2aac {
		return 0x100aa2aac, true, nil
	}
	// 0x100aa2aac: movk w11, #0x8ee, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x8ee, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ab0 {
		return 0x100aa2ab0, true, nil
	}
	// 0x100aa2ab0: mov w12, #0x404c
	s.x[12] = uint64(uint32(uint64(0x404c)))
	if runtime != nil && runtime.stopPC == 0x100aa2ab4 {
		return 0x100aa2ab4, true, nil
	}
	// 0x100aa2ab4: movk w12, #0xa0f6, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xa0f6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ab8 {
		return 0x100aa2ab8, true, nil
	}
	// 0x100aa2ab8: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2abc {
		return 0x100aa2abc, true, nil
	}
	// 0x100aa2abc: mov w15, #0xee3b
	s.x[15] = uint64(uint32(uint64(0xee3b)))
	if runtime != nil && runtime.stopPC == 0x100aa2ac0 {
		return 0x100aa2ac0, true, nil
	}
	// 0x100aa2ac0: movk w15, #0x2dcd, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x2dcd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ac4 {
		return 0x100aa2ac4, true, nil
	}
	// 0x100aa2ac4: add w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ac8 {
		return 0x100aa2ac8, true, nil
	}
	// 0x100aa2ac8: mov w15, #0xf994
	s.x[15] = uint64(uint32(uint64(0xf994)))
	if runtime != nil && runtime.stopPC == 0x100aa2acc {
		return 0x100aa2acc, true, nil
	}
	// 0x100aa2acc: movk w15, #0x481a, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x481a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ad0 {
		return 0x100aa2ad0, true, nil
	}
	// 0x100aa2ad0: eor w15, w9, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ad4 {
		return 0x100aa2ad4, true, nil
	}
	// 0x100aa2ad4: mov w16, #0x57e3
	s.x[16] = uint64(uint32(uint64(0x57e3)))
	if runtime != nil && runtime.stopPC == 0x100aa2ad8 {
		return 0x100aa2ad8, true, nil
	}
	// 0x100aa2ad8: movk w16, #0xc521, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xc521, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2adc {
		return 0x100aa2adc, true, nil
	}
	// 0x100aa2adc: mov w17, #0x7ffb
	s.x[17] = uint64(uint32(uint64(0x7ffb)))
	if runtime != nil && runtime.stopPC == 0x100aa2ae0 {
		return 0x100aa2ae0, true, nil
	}
	// 0x100aa2ae0: movk w17, #0x1ffd, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x1ffd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ae4 {
		return 0x100aa2ae4, true, nil
	}
	// 0x100aa2ae4: eor w17, w9, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ae8 {
		return 0x100aa2ae8, true, nil
	}
	// 0x100aa2ae8: mov w0, #0xd18e
	s.x[0] = uint64(uint32(uint64(0xd18e)))
	if runtime != nil && runtime.stopPC == 0x100aa2aec {
		return 0x100aa2aec, true, nil
	}
	// 0x100aa2aec: movk w0, #0x92c6, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x92c6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2af0 {
		return 0x100aa2af0, true, nil
	}
	// 0x100aa2af0: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2af4 {
		return 0x100aa2af4, true, nil
	}
	// 0x100aa2af4: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2af8 {
		return 0x100aa2af8, true, nil
	}
	// 0x100aa2af8: add w10, w15, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2afc {
		return 0x100aa2afc, true, nil
	}
	// 0x100aa2afc: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b00 {
		return 0x100aa2b00, true, nil
	}
	// 0x100aa2b00: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b04 {
		return 0x100aa2b04, true, nil
	}
	// 0x100aa2b04: mov w10, #0xd4e9
	s.x[10] = uint64(uint32(uint64(0xd4e9)))
	if runtime != nil && runtime.stopPC == 0x100aa2b08 {
		return 0x100aa2b08, true, nil
	}
	// 0x100aa2b08: movk w10, #0xa4e7, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa4e7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b0c {
		return 0x100aa2b0c, true, nil
	}
	// 0x100aa2b0c: add w11, w17, w0
	s.x[11] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b10 {
		return 0x100aa2b10, true, nil
	}
	// 0x100aa2b10: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b14 {
		return 0x100aa2b14, true, nil
	}
	// 0x100aa2b14: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b18 {
		return 0x100aa2b18, true, nil
	}
	// 0x100aa2b18: str w9, [sp, #0x1a8]
	if err := s.storeValue(addOffset(s.sp, 424), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100aa2b1c {
		return 0x100aa2b1c, true, nil
	}
	// 0x100aa2b1c: mov w9, #0x9f5b
	s.x[9] = uint64(uint32(uint64(0x9f5b)))
	if runtime != nil && runtime.stopPC == 0x100aa2b20 {
		return 0x100aa2b20, true, nil
	}
	// 0x100aa2b20: movk w9, #0xa53e, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xa53e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b24 {
		return 0x100aa2b24, true, nil
	}
	// 0x100aa2b24: eor w9, w25, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b28 {
		return 0x100aa2b28, true, nil
	}
	// 0x100aa2b28: mov w10, #0x60a5
	s.x[10] = uint64(uint32(uint64(0x60a5)))
	if runtime != nil && runtime.stopPC == 0x100aa2b2c {
		return 0x100aa2b2c, true, nil
	}
	// 0x100aa2b2c: movk w10, #0x5ac1, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x5ac1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b30 {
		return 0x100aa2b30, true, nil
	}
	// 0x100aa2b30: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b34 {
		return 0x100aa2b34, true, nil
	}
	// 0x100aa2b34: mov w10, #0xd41b
	s.x[10] = uint64(uint32(uint64(0xd41b)))
	if runtime != nil && runtime.stopPC == 0x100aa2b38 {
		return 0x100aa2b38, true, nil
	}
	// 0x100aa2b38: movk w10, #0xc718, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xc718, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b3c {
		return 0x100aa2b3c, true, nil
	}
	// 0x100aa2b3c: eor w10, w25, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b40 {
		return 0x100aa2b40, true, nil
	}
	// 0x100aa2b40: mov w11, #0x2be5
	s.x[11] = uint64(uint32(uint64(0x2be5)))
	if runtime != nil && runtime.stopPC == 0x100aa2b44 {
		return 0x100aa2b44, true, nil
	}
	// 0x100aa2b44: movk w11, #0x38e7, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x38e7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b48 {
		return 0x100aa2b48, true, nil
	}
	// 0x100aa2b48: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b4c {
		return 0x100aa2b4c, true, nil
	}
	// 0x100aa2b4c: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b50 {
		return 0x100aa2b50, true, nil
	}
	// 0x100aa2b50: mov w10, #0x5728
	s.x[10] = uint64(uint32(uint64(0x5728)))
	if runtime != nil && runtime.stopPC == 0x100aa2b54 {
		return 0x100aa2b54, true, nil
	}
	// 0x100aa2b54: movk w10, #0xb236, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xb236, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b58 {
		return 0x100aa2b58, true, nil
	}
	// 0x100aa2b58: eor w10, w25, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b5c {
		return 0x100aa2b5c, true, nil
	}
	// 0x100aa2b5c: mov w11, #0xa8d8
	s.x[11] = uint64(uint32(uint64(0xa8d8)))
	if runtime != nil && runtime.stopPC == 0x100aa2b60 {
		return 0x100aa2b60, true, nil
	}
	// 0x100aa2b60: movk w11, #0x4dc9, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x4dc9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b64 {
		return 0x100aa2b64, true, nil
	}
	// 0x100aa2b64: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b68 {
		return 0x100aa2b68, true, nil
	}
	// 0x100aa2b68: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b6c {
		return 0x100aa2b6c, true, nil
	}
	// 0x100aa2b6c: mov w10, #0x1beb
	s.x[10] = uint64(uint32(uint64(0x1beb)))
	if runtime != nil && runtime.stopPC == 0x100aa2b70 {
		return 0x100aa2b70, true, nil
	}
	// 0x100aa2b70: movk w10, #0xce44, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xce44, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b74 {
		return 0x100aa2b74, true, nil
	}
	// 0x100aa2b74: ldr w15, [sp, #0x184]
	value214, err := s.load(addOffset(s.sp, 388), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value214))
	if runtime != nil && runtime.stopPC == 0x100aa2b78 {
		return 0x100aa2b78, true, nil
	}
	// 0x100aa2b78: eor w10, w15, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b7c {
		return 0x100aa2b7c, true, nil
	}
	// 0x100aa2b7c: mov w11, #0xe415
	s.x[11] = uint64(uint32(uint64(0xe415)))
	if runtime != nil && runtime.stopPC == 0x100aa2b80 {
		return 0x100aa2b80, true, nil
	}
	// 0x100aa2b80: movk w11, #0x31bb, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x31bb, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b84 {
		return 0x100aa2b84, true, nil
	}
	// 0x100aa2b84: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b88 {
		return 0x100aa2b88, true, nil
	}
	// 0x100aa2b88: mov w11, #0x7798
	s.x[11] = uint64(uint32(uint64(0x7798)))
	if runtime != nil && runtime.stopPC == 0x100aa2b8c {
		return 0x100aa2b8c, true, nil
	}
	// 0x100aa2b8c: movk w11, #0x8dd9, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x8dd9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b90 {
		return 0x100aa2b90, true, nil
	}
	// 0x100aa2b90: eor w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2b94 {
		return 0x100aa2b94, true, nil
	}
	// 0x100aa2b94: mov w12, #0x8868
	s.x[12] = uint64(uint32(uint64(0x8868)))
	if runtime != nil && runtime.stopPC == 0x100aa2b98 {
		return 0x100aa2b98, true, nil
	}
	// 0x100aa2b98: movk w12, #0x7226, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x7226, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2b9c {
		return 0x100aa2b9c, true, nil
	}
	// 0x100aa2b9c: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ba0 {
		return 0x100aa2ba0, true, nil
	}
	// 0x100aa2ba0: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ba4 {
		return 0x100aa2ba4, true, nil
	}
	// 0x100aa2ba4: mov w11, #0x701b
	s.x[11] = uint64(uint32(uint64(0x701b)))
	if runtime != nil && runtime.stopPC == 0x100aa2ba8 {
		return 0x100aa2ba8, true, nil
	}
	// 0x100aa2ba8: movk w11, #0x938d, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x938d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2bac {
		return 0x100aa2bac, true, nil
	}
	// 0x100aa2bac: eor w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bb0 {
		return 0x100aa2bb0, true, nil
	}
	// 0x100aa2bb0: mov w12, #0x8fe5
	s.x[12] = uint64(uint32(uint64(0x8fe5)))
	if runtime != nil && runtime.stopPC == 0x100aa2bb4 {
		return 0x100aa2bb4, true, nil
	}
	// 0x100aa2bb4: movk w12, #0x6c72, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x6c72, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2bb8 {
		return 0x100aa2bb8, true, nil
	}
	// 0x100aa2bb8: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bbc {
		return 0x100aa2bbc, true, nil
	}
	// 0x100aa2bbc: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bc0 {
		return 0x100aa2bc0, true, nil
	}
	// 0x100aa2bc0: mov w11, #0x5622
	s.x[11] = uint64(uint32(uint64(0x5622)))
	if runtime != nil && runtime.stopPC == 0x100aa2bc4 {
		return 0x100aa2bc4, true, nil
	}
	// 0x100aa2bc4: movk w11, #0x39c, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x39c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2bc8 {
		return 0x100aa2bc8, true, nil
	}
	// 0x100aa2bc8: add w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bcc {
		return 0x100aa2bcc, true, nil
	}
	// 0x100aa2bcc: add w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bd0 {
		return 0x100aa2bd0, true, nil
	}
	// 0x100aa2bd0: mov w10, #0xd073
	s.x[10] = uint64(uint32(uint64(0xd073)))
	if runtime != nil && runtime.stopPC == 0x100aa2bd4 {
		return 0x100aa2bd4, true, nil
	}
	// 0x100aa2bd4: movk w10, #0xa559, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa559, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2bd8 {
		return 0x100aa2bd8, true, nil
	}
	// 0x100aa2bd8: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bdc {
		return 0x100aa2bdc, true, nil
	}
	// 0x100aa2bdc: mov w11, #0x32df
	s.x[11] = uint64(uint32(uint64(0x32df)))
	if runtime != nil && runtime.stopPC == 0x100aa2be0 {
		return 0x100aa2be0, true, nil
	}
	// 0x100aa2be0: movk w11, #0x39da, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x39da, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2be4 {
		return 0x100aa2be4, true, nil
	}
	// 0x100aa2be4: mov w12, #0x6497
	s.x[12] = uint64(uint32(uint64(0x6497)))
	if runtime != nil && runtime.stopPC == 0x100aa2be8 {
		return 0x100aa2be8, true, nil
	}
	// 0x100aa2be8: movk w12, #0x8962, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x8962, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2bec {
		return 0x100aa2bec, true, nil
	}
	// 0x100aa2bec: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bf0 {
		return 0x100aa2bf0, true, nil
	}
	// 0x100aa2bf0: mov w15, #0x863b
	s.x[15] = uint64(uint32(uint64(0x863b)))
	if runtime != nil && runtime.stopPC == 0x100aa2bf4 {
		return 0x100aa2bf4, true, nil
	}
	// 0x100aa2bf4: movk w15, #0x15e1, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x15e1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2bf8 {
		return 0x100aa2bf8, true, nil
	}
	// 0x100aa2bf8: add w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2bfc {
		return 0x100aa2bfc, true, nil
	}
	// 0x100aa2bfc: mov w15, #0x16c9
	s.x[15] = uint64(uint32(uint64(0x16c9)))
	if runtime != nil && runtime.stopPC == 0x100aa2c00 {
		return 0x100aa2c00, true, nil
	}
	// 0x100aa2c00: movk w15, #0xb8b8, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xb8b8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c04 {
		return 0x100aa2c04, true, nil
	}
	// 0x100aa2c04: eor w15, w9, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c08 {
		return 0x100aa2c08, true, nil
	}
	// 0x100aa2c08: mov w16, #0xf465
	s.x[16] = uint64(uint32(uint64(0xf465)))
	if runtime != nil && runtime.stopPC == 0x100aa2c0c {
		return 0x100aa2c0c, true, nil
	}
	// 0x100aa2c0c: movk w16, #0x243b, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x243b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c10 {
		return 0x100aa2c10, true, nil
	}
	// 0x100aa2c10: mov w17, #0xbf7f
	s.x[17] = uint64(uint32(uint64(0xbf7f)))
	if runtime != nil && runtime.stopPC == 0x100aa2c14 {
		return 0x100aa2c14, true, nil
	}
	// 0x100aa2c14: movk w17, #0xf7ff, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xf7ff, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c18 {
		return 0x100aa2c18, true, nil
	}
	// 0x100aa2c18: eor w17, w9, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c1c {
		return 0x100aa2c1c, true, nil
	}
	// 0x100aa2c1c: mov w0, #0x5dd3
	s.x[0] = uint64(uint32(uint64(0x5dd3)))
	if runtime != nil && runtime.stopPC == 0x100aa2c20 {
		return 0x100aa2c20, true, nil
	}
	// 0x100aa2c20: movk w0, #0x6b7c, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x6b7c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c24 {
		return 0x100aa2c24, true, nil
	}
	// 0x100aa2c24: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c28 {
		return 0x100aa2c28, true, nil
	}
	// 0x100aa2c28: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c2c {
		return 0x100aa2c2c, true, nil
	}
	// 0x100aa2c2c: add w10, w15, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c30 {
		return 0x100aa2c30, true, nil
	}
	// 0x100aa2c30: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c34 {
		return 0x100aa2c34, true, nil
	}
	// 0x100aa2c34: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c38 {
		return 0x100aa2c38, true, nil
	}
	// 0x100aa2c38: mov w10, #0x13a
	s.x[10] = uint64(uint32(uint64(0x13a)))
	if runtime != nil && runtime.stopPC == 0x100aa2c3c {
		return 0x100aa2c3c, true, nil
	}
	// 0x100aa2c3c: movk w10, #0xb36c, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xb36c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c40 {
		return 0x100aa2c40, true, nil
	}
	// 0x100aa2c40: add w11, w17, w0
	s.x[11] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c44 {
		return 0x100aa2c44, true, nil
	}
	// 0x100aa2c44: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c48 {
		return 0x100aa2c48, true, nil
	}
	// 0x100aa2c48: eor w21, w9, w10
	s.x[21] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c4c {
		return 0x100aa2c4c, true, nil
	}
	// 0x100aa2c4c: mov w9, #0x97d2
	s.x[9] = uint64(uint32(uint64(0x97d2)))
	if runtime != nil && runtime.stopPC == 0x100aa2c50 {
		return 0x100aa2c50, true, nil
	}
	// 0x100aa2c50: movk w9, #0x84a1, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x84a1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c54 {
		return 0x100aa2c54, true, nil
	}
	// 0x100aa2c54: ldr w12, [sp, #0x194]
	value270, err := s.load(addOffset(s.sp, 404), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value270))
	if runtime != nil && runtime.stopPC == 0x100aa2c58 {
		return 0x100aa2c58, true, nil
	}
	// 0x100aa2c58: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c5c {
		return 0x100aa2c5c, true, nil
	}
	// 0x100aa2c5c: mov w10, #0x682e
	s.x[10] = uint64(uint32(uint64(0x682e)))
	if runtime != nil && runtime.stopPC == 0x100aa2c60 {
		return 0x100aa2c60, true, nil
	}
	// 0x100aa2c60: movk w10, #0x7b5e, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x7b5e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c64 {
		return 0x100aa2c64, true, nil
	}
	// 0x100aa2c64: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c68 {
		return 0x100aa2c68, true, nil
	}
	// 0x100aa2c68: mov w10, #0x4639
	s.x[10] = uint64(uint32(uint64(0x4639)))
	if runtime != nil && runtime.stopPC == 0x100aa2c6c {
		return 0x100aa2c6c, true, nil
	}
	// 0x100aa2c6c: movk w10, #0x9832, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x9832, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c70 {
		return 0x100aa2c70, true, nil
	}
	// 0x100aa2c70: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c74 {
		return 0x100aa2c74, true, nil
	}
	// 0x100aa2c74: mov w11, #0xb9c7
	s.x[11] = uint64(uint32(uint64(0xb9c7)))
	if runtime != nil && runtime.stopPC == 0x100aa2c78 {
		return 0x100aa2c78, true, nil
	}
	// 0x100aa2c78: movk w11, #0x67cd, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x67cd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c7c {
		return 0x100aa2c7c, true, nil
	}
	// 0x100aa2c7c: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c80 {
		return 0x100aa2c80, true, nil
	}
	// 0x100aa2c80: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c84 {
		return 0x100aa2c84, true, nil
	}
	// 0x100aa2c84: mov w10, #0xcd83
	s.x[10] = uint64(uint32(uint64(0xcd83)))
	if runtime != nil && runtime.stopPC == 0x100aa2c88 {
		return 0x100aa2c88, true, nil
	}
	// 0x100aa2c88: movk w10, #0xcc83, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xcc83, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c8c {
		return 0x100aa2c8c, true, nil
	}
	// 0x100aa2c8c: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c90 {
		return 0x100aa2c90, true, nil
	}
	// 0x100aa2c90: mov w11, #0x327d
	s.x[11] = uint64(uint32(uint64(0x327d)))
	if runtime != nil && runtime.stopPC == 0x100aa2c94 {
		return 0x100aa2c94, true, nil
	}
	// 0x100aa2c94: movk w11, #0x337c, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x337c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2c98 {
		return 0x100aa2c98, true, nil
	}
	// 0x100aa2c98: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2c9c {
		return 0x100aa2c9c, true, nil
	}
	// 0x100aa2c9c: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ca0 {
		return 0x100aa2ca0, true, nil
	}
	// 0x100aa2ca0: mov w10, #0x2314
	s.x[10] = uint64(uint32(uint64(0x2314)))
	if runtime != nil && runtime.stopPC == 0x100aa2ca4 {
		return 0x100aa2ca4, true, nil
	}
	// 0x100aa2ca4: movk w10, #0x8018, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x8018, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ca8 {
		return 0x100aa2ca8, true, nil
	}
	// 0x100aa2ca8: eor w10, w14, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cac {
		return 0x100aa2cac, true, nil
	}
	// 0x100aa2cac: mov w11, #0xdcec
	s.x[11] = uint64(uint32(uint64(0xdcec)))
	if runtime != nil && runtime.stopPC == 0x100aa2cb0 {
		return 0x100aa2cb0, true, nil
	}
	// 0x100aa2cb0: movk w11, #0x7fe7, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x7fe7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2cb4 {
		return 0x100aa2cb4, true, nil
	}
	// 0x100aa2cb4: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cb8 {
		return 0x100aa2cb8, true, nil
	}
	// 0x100aa2cb8: mov w11, #0xbb3f
	s.x[11] = uint64(uint32(uint64(0xbb3f)))
	if runtime != nil && runtime.stopPC == 0x100aa2cbc {
		return 0x100aa2cbc, true, nil
	}
	// 0x100aa2cbc: movk w11, #0xf290, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xf290, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2cc0 {
		return 0x100aa2cc0, true, nil
	}
	// 0x100aa2cc0: eor w11, w14, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cc4 {
		return 0x100aa2cc4, true, nil
	}
	// 0x100aa2cc4: mov w12, #0x44c1
	s.x[12] = uint64(uint32(uint64(0x44c1)))
	if runtime != nil && runtime.stopPC == 0x100aa2cc8 {
		return 0x100aa2cc8, true, nil
	}
	// 0x100aa2cc8: movk w12, #0xd6f, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xd6f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ccc {
		return 0x100aa2ccc, true, nil
	}
	// 0x100aa2ccc: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cd0 {
		return 0x100aa2cd0, true, nil
	}
	// 0x100aa2cd0: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cd4 {
		return 0x100aa2cd4, true, nil
	}
	// 0x100aa2cd4: mov w11, #0x8443
	s.x[11] = uint64(uint32(uint64(0x8443)))
	if runtime != nil && runtime.stopPC == 0x100aa2cd8 {
		return 0x100aa2cd8, true, nil
	}
	// 0x100aa2cd8: movk w11, #0xa298, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa298, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2cdc {
		return 0x100aa2cdc, true, nil
	}
	// 0x100aa2cdc: eor w11, w14, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ce0 {
		return 0x100aa2ce0, true, nil
	}
	// 0x100aa2ce0: mov w12, #0x7bbd
	s.x[12] = uint64(uint32(uint64(0x7bbd)))
	if runtime != nil && runtime.stopPC == 0x100aa2ce4 {
		return 0x100aa2ce4, true, nil
	}
	// 0x100aa2ce4: movk w12, #0x5d67, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x5d67, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ce8 {
		return 0x100aa2ce8, true, nil
	}
	// 0x100aa2ce8: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cec {
		return 0x100aa2cec, true, nil
	}
	// 0x100aa2cec: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cf0 {
		return 0x100aa2cf0, true, nil
	}
	// 0x100aa2cf0: mov w11, #0x3118
	s.x[11] = uint64(uint32(uint64(0x3118)))
	if runtime != nil && runtime.stopPC == 0x100aa2cf4 {
		return 0x100aa2cf4, true, nil
	}
	// 0x100aa2cf4: movk w11, #0xb480, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb480, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2cf8 {
		return 0x100aa2cf8, true, nil
	}
	// 0x100aa2cf8: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2cfc {
		return 0x100aa2cfc, true, nil
	}
	// 0x100aa2cfc: add w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d00 {
		return 0x100aa2d00, true, nil
	}
	// 0x100aa2d00: mov w10, #0xb824
	s.x[10] = uint64(uint32(uint64(0xb824)))
	if runtime != nil && runtime.stopPC == 0x100aa2d04 {
		return 0x100aa2d04, true, nil
	}
	// 0x100aa2d04: movk w10, #0xca29, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xca29, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d08 {
		return 0x100aa2d08, true, nil
	}
	// 0x100aa2d08: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d0c {
		return 0x100aa2d0c, true, nil
	}
	// 0x100aa2d0c: mov w11, #0xbf94
	s.x[11] = uint64(uint32(uint64(0xbf94)))
	if runtime != nil && runtime.stopPC == 0x100aa2d10 {
		return 0x100aa2d10, true, nil
	}
	// 0x100aa2d10: movk w11, #0x2189, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x2189, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d14 {
		return 0x100aa2d14, true, nil
	}
	// 0x100aa2d14: mov w12, #0x2b2f
	s.x[12] = uint64(uint32(uint64(0x2b2f)))
	if runtime != nil && runtime.stopPC == 0x100aa2d18 {
		return 0x100aa2d18, true, nil
	}
	// 0x100aa2d18: movk w12, #0x5e49, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x5e49, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d1c {
		return 0x100aa2d1c, true, nil
	}
	// 0x100aa2d1c: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d20 {
		return 0x100aa2d20, true, nil
	}
	// 0x100aa2d20: mov w14, #0x2c99
	s.x[14] = uint64(uint32(uint64(0x2c99)))
	if runtime != nil && runtime.stopPC == 0x100aa2d24 {
		return 0x100aa2d24, true, nil
	}
	// 0x100aa2d24: movk w14, #0xb5e9, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xb5e9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d28 {
		return 0x100aa2d28, true, nil
	}
	// 0x100aa2d28: add w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d2c {
		return 0x100aa2d2c, true, nil
	}
	// 0x100aa2d2c: mov w14, #0x94bd
	s.x[14] = uint64(uint32(uint64(0x94bd)))
	if runtime != nil && runtime.stopPC == 0x100aa2d30 {
		return 0x100aa2d30, true, nil
	}
	// 0x100aa2d30: movk w14, #0x3fce, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x3fce, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d34 {
		return 0x100aa2d34, true, nil
	}
	// 0x100aa2d34: eor w14, w9, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d38 {
		return 0x100aa2d38, true, nil
	}
	// 0x100aa2d38: mov w15, #0x930b
	s.x[15] = uint64(uint32(uint64(0x930b)))
	if runtime != nil && runtime.stopPC == 0x100aa2d3c {
		return 0x100aa2d3c, true, nil
	}
	// 0x100aa2d3c: movk w15, #0xd46e, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xd46e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d40 {
		return 0x100aa2d40, true, nil
	}
	// 0x100aa2d40: mov w16, #0xfffe
	s.x[16] = uint64(uint32(uint64(0xfffe)))
	if runtime != nil && runtime.stopPC == 0x100aa2d44 {
		return 0x100aa2d44, true, nil
	}
	// 0x100aa2d44: movk w16, #0xbff1, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xbff1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d48 {
		return 0x100aa2d48, true, nil
	}
	// 0x100aa2d48: eor w16, w9, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d4c {
		return 0x100aa2d4c, true, nil
	}
	// 0x100aa2d4c: mov w17, #0xf84a
	s.x[17] = uint64(uint32(uint64(0xf84a)))
	if runtime != nil && runtime.stopPC == 0x100aa2d50 {
		return 0x100aa2d50, true, nil
	}
	// 0x100aa2d50: movk w17, #0x5451, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x5451, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d54 {
		return 0x100aa2d54, true, nil
	}
	// 0x100aa2d54: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d58 {
		return 0x100aa2d58, true, nil
	}
	// 0x100aa2d58: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d5c {
		return 0x100aa2d5c, true, nil
	}
	// 0x100aa2d5c: add w10, w14, w15
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d60 {
		return 0x100aa2d60, true, nil
	}
	// 0x100aa2d60: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d64 {
		return 0x100aa2d64, true, nil
	}
	// 0x100aa2d64: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d68 {
		return 0x100aa2d68, true, nil
	}
	// 0x100aa2d68: mov w10, #0xe420
	s.x[10] = uint64(uint32(uint64(0xe420)))
	if runtime != nil && runtime.stopPC == 0x100aa2d6c {
		return 0x100aa2d6c, true, nil
	}
	// 0x100aa2d6c: movk w10, #0xc44f, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xc44f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d70 {
		return 0x100aa2d70, true, nil
	}
	// 0x100aa2d70: add w11, w16, w17
	s.x[11] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d74 {
		return 0x100aa2d74, true, nil
	}
	// 0x100aa2d74: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d78 {
		return 0x100aa2d78, true, nil
	}
	// 0x100aa2d78: eor w25, w9, w10
	s.x[25] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d7c {
		return 0x100aa2d7c, true, nil
	}
	// 0x100aa2d7c: mov w9, #0x657
	s.x[9] = uint64(uint32(uint64(0x657)))
	if runtime != nil && runtime.stopPC == 0x100aa2d80 {
		return 0x100aa2d80, true, nil
	}
	// 0x100aa2d80: movk w9, #0x91a7, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x91a7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d84 {
		return 0x100aa2d84, true, nil
	}
	// 0x100aa2d84: ldr w12, [sp, #0x190]
	value346, err := s.load(addOffset(s.sp, 400), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value346))
	if runtime != nil && runtime.stopPC == 0x100aa2d88 {
		return 0x100aa2d88, true, nil
	}
	// 0x100aa2d88: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d8c {
		return 0x100aa2d8c, true, nil
	}
	// 0x100aa2d8c: mov w10, #0xf9a9
	s.x[10] = uint64(uint32(uint64(0xf9a9)))
	if runtime != nil && runtime.stopPC == 0x100aa2d90 {
		return 0x100aa2d90, true, nil
	}
	// 0x100aa2d90: movk w10, #0x6e58, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x6e58, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2d94 {
		return 0x100aa2d94, true, nil
	}
	// 0x100aa2d94: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2d98 {
		return 0x100aa2d98, true, nil
	}
	// 0x100aa2d98: mov w10, #0x2025
	s.x[10] = uint64(uint32(uint64(0x2025)))
	if runtime != nil && runtime.stopPC == 0x100aa2d9c {
		return 0x100aa2d9c, true, nil
	}
	// 0x100aa2d9c: movk w10, #0x292d, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x292d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2da0 {
		return 0x100aa2da0, true, nil
	}
	// 0x100aa2da0: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2da4 {
		return 0x100aa2da4, true, nil
	}
	// 0x100aa2da4: mov w11, #0xdfdb
	s.x[11] = uint64(uint32(uint64(0xdfdb)))
	if runtime != nil && runtime.stopPC == 0x100aa2da8 {
		return 0x100aa2da8, true, nil
	}
	// 0x100aa2da8: movk w11, #0xd6d2, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd6d2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2dac {
		return 0x100aa2dac, true, nil
	}
	// 0x100aa2dac: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2db0 {
		return 0x100aa2db0, true, nil
	}
	// 0x100aa2db0: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2db4 {
		return 0x100aa2db4, true, nil
	}
	// 0x100aa2db4: mov w10, #0x3a1a
	s.x[10] = uint64(uint32(uint64(0x3a1a)))
	if runtime != nil && runtime.stopPC == 0x100aa2db8 {
		return 0x100aa2db8, true, nil
	}
	// 0x100aa2db8: movk w10, #0x689a, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x689a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2dbc {
		return 0x100aa2dbc, true, nil
	}
	// 0x100aa2dbc: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2dc0 {
		return 0x100aa2dc0, true, nil
	}
	// 0x100aa2dc0: mov w11, #0xc5e6
	s.x[11] = uint64(uint32(uint64(0xc5e6)))
	if runtime != nil && runtime.stopPC == 0x100aa2dc4 {
		return 0x100aa2dc4, true, nil
	}
	// 0x100aa2dc4: movk w11, #0x9765, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9765, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2dc8 {
		return 0x100aa2dc8, true, nil
	}
	// 0x100aa2dc8: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2dcc {
		return 0x100aa2dcc, true, nil
	}
	// 0x100aa2dcc: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2dd0 {
		return 0x100aa2dd0, true, nil
	}
	// 0x100aa2dd0: mov w10, #0x4a19
	s.x[10] = uint64(uint32(uint64(0x4a19)))
	if runtime != nil && runtime.stopPC == 0x100aa2dd4 {
		return 0x100aa2dd4, true, nil
	}
	// 0x100aa2dd4: movk w10, #0xe94, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xe94, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2dd8 {
		return 0x100aa2dd8, true, nil
	}
	// 0x100aa2dd8: eor w10, w24, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ddc {
		return 0x100aa2ddc, true, nil
	}
	// 0x100aa2ddc: mov w11, #0x529b
	s.x[11] = uint64(uint32(uint64(0x529b)))
	if runtime != nil && runtime.stopPC == 0x100aa2de0 {
		return 0x100aa2de0, true, nil
	}
	// 0x100aa2de0: movk w11, #0x4df9, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x4df9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2de4 {
		return 0x100aa2de4, true, nil
	}
	// 0x100aa2de4: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2de8 {
		return 0x100aa2de8, true, nil
	}
	// 0x100aa2de8: mov w11, #0x4df6
	s.x[11] = uint64(uint32(uint64(0x4df6)))
	if runtime != nil && runtime.stopPC == 0x100aa2dec {
		return 0x100aa2dec, true, nil
	}
	// 0x100aa2dec: movk w11, #0xfca5, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xfca5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2df0 {
		return 0x100aa2df0, true, nil
	}
	// 0x100aa2df0: eor w11, w24, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2df4 {
		return 0x100aa2df4, true, nil
	}
	// 0x100aa2df4: mov w12, #0x5576
	s.x[12] = uint64(uint32(uint64(0x5576)))
	if runtime != nil && runtime.stopPC == 0x100aa2df8 {
		return 0x100aa2df8, true, nil
	}
	// 0x100aa2df8: movk w12, #0xbfc8, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xbfc8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2dfc {
		return 0x100aa2dfc, true, nil
	}
	// 0x100aa2dfc: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e00 {
		return 0x100aa2e00, true, nil
	}
	// 0x100aa2e00: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e04 {
		return 0x100aa2e04, true, nil
	}
	// 0x100aa2e04: mov w11, #0xfcfb
	s.x[11] = uint64(uint32(uint64(0xfcfb)))
	if runtime != nil && runtime.stopPC == 0x100aa2e08 {
		return 0x100aa2e08, true, nil
	}
	// 0x100aa2e08: movk w11, #0x9eb3, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9eb3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e0c {
		return 0x100aa2e0c, true, nil
	}
	// 0x100aa2e0c: eor w11, w24, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e10 {
		return 0x100aa2e10, true, nil
	}
	// 0x100aa2e10: mov w12, #0xe479
	s.x[12] = uint64(uint32(uint64(0xe479)))
	if runtime != nil && runtime.stopPC == 0x100aa2e14 {
		return 0x100aa2e14, true, nil
	}
	// 0x100aa2e14: movk w12, #0xddde, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xddde, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e18 {
		return 0x100aa2e18, true, nil
	}
	// 0x100aa2e18: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e1c {
		return 0x100aa2e1c, true, nil
	}
	// 0x100aa2e1c: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e20 {
		return 0x100aa2e20, true, nil
	}
	// 0x100aa2e20: mov w11, #0xe744
	s.x[11] = uint64(uint32(uint64(0xe744)))
	if runtime != nil && runtime.stopPC == 0x100aa2e24 {
		return 0x100aa2e24, true, nil
	}
	// 0x100aa2e24: movk w11, #0xd74f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd74f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e28 {
		return 0x100aa2e28, true, nil
	}
	// 0x100aa2e28: add w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e2c {
		return 0x100aa2e2c, true, nil
	}
	// 0x100aa2e2c: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e30 {
		return 0x100aa2e30, true, nil
	}
	// 0x100aa2e30: mov w10, #0x6791
	s.x[10] = uint64(uint32(uint64(0x6791)))
	if runtime != nil && runtime.stopPC == 0x100aa2e34 {
		return 0x100aa2e34, true, nil
	}
	// 0x100aa2e34: movk w10, #0xa3aa, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa3aa, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e38 {
		return 0x100aa2e38, true, nil
	}
	// 0x100aa2e38: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e3c {
		return 0x100aa2e3c, true, nil
	}
	// 0x100aa2e3c: mov w11, #0x361b
	s.x[11] = uint64(uint32(uint64(0x361b)))
	if runtime != nil && runtime.stopPC == 0x100aa2e40 {
		return 0x100aa2e40, true, nil
	}
	// 0x100aa2e40: movk w11, #0x6b7a, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x6b7a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e44 {
		return 0x100aa2e44, true, nil
	}
	// 0x100aa2e44: mov w12, #0x13e
	s.x[12] = uint64(uint32(uint64(0x13e)))
	if runtime != nil && runtime.stopPC == 0x100aa2e48 {
		return 0x100aa2e48, true, nil
	}
	// 0x100aa2e48: movk w12, #0xfcdd, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xfcdd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e4c {
		return 0x100aa2e4c, true, nil
	}
	// 0x100aa2e4c: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e50 {
		return 0x100aa2e50, true, nil
	}
	// 0x100aa2e50: mov w14, #0x50b6
	s.x[14] = uint64(uint32(uint64(0x50b6)))
	if runtime != nil && runtime.stopPC == 0x100aa2e54 {
		return 0x100aa2e54, true, nil
	}
	// 0x100aa2e54: movk w14, #0x340d, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x340d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e58 {
		return 0x100aa2e58, true, nil
	}
	// 0x100aa2e58: add w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e5c {
		return 0x100aa2e5c, true, nil
	}
	// 0x100aa2e5c: mov w14, #0x352c
	s.x[14] = uint64(uint32(uint64(0x352c)))
	if runtime != nil && runtime.stopPC == 0x100aa2e60 {
		return 0x100aa2e60, true, nil
	}
	// 0x100aa2e60: movk w14, #0x97b7, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x97b7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e64 {
		return 0x100aa2e64, true, nil
	}
	// 0x100aa2e64: eor w14, w9, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e68 {
		return 0x100aa2e68, true, nil
	}
	// 0x100aa2e68: mov w15, #0x64a8
	s.x[15] = uint64(uint32(uint64(0x64a8)))
	if runtime != nil && runtime.stopPC == 0x100aa2e6c {
		return 0x100aa2e6c, true, nil
	}
	// 0x100aa2e6c: movk w15, #0x5f67, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x5f67, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e70 {
		return 0x100aa2e70, true, nil
	}
	// 0x100aa2e70: mov w16, #0xfdf7
	s.x[16] = uint64(uint32(uint64(0xfdf7)))
	if runtime != nil && runtime.stopPC == 0x100aa2e74 {
		return 0x100aa2e74, true, nil
	}
	// 0x100aa2e74: movk w16, #0xffef, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xffef, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e78 {
		return 0x100aa2e78, true, nil
	}
	// 0x100aa2e78: eor w16, w9, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e7c {
		return 0x100aa2e7c, true, nil
	}
	// 0x100aa2e7c: mov w17, #0xac7d
	s.x[17] = uint64(uint32(uint64(0xac7d)))
	if runtime != nil && runtime.stopPC == 0x100aa2e80 {
		return 0x100aa2e80, true, nil
	}
	// 0x100aa2e80: movk w17, #0x373f, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x373f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2e84 {
		return 0x100aa2e84, true, nil
	}
	// 0x100aa2e84: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e88 {
		return 0x100aa2e88, true, nil
	}
	// 0x100aa2e88: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e8c {
		return 0x100aa2e8c, true, nil
	}
	// 0x100aa2e8c: add w10, w14, w15
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e90 {
		return 0x100aa2e90, true, nil
	}
	// 0x100aa2e90: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e94 {
		return 0x100aa2e94, true, nil
	}
	// 0x100aa2e94: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2e98 {
		return 0x100aa2e98, true, nil
	}
	// 0x100aa2e98: mov w10, #0xb21c
	s.x[10] = uint64(uint32(uint64(0xb21c)))
	if runtime != nil && runtime.stopPC == 0x100aa2e9c {
		return 0x100aa2e9c, true, nil
	}
	// 0x100aa2e9c: movk w10, #0xe73f, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xe73f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ea0 {
		return 0x100aa2ea0, true, nil
	}
	// 0x100aa2ea0: add w11, w16, w17
	s.x[11] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ea4 {
		return 0x100aa2ea4, true, nil
	}
	// 0x100aa2ea4: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ea8 {
		return 0x100aa2ea8, true, nil
	}
	// 0x100aa2ea8: eor w27, w9, w10
	s.x[27] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2eac {
		return 0x100aa2eac, true, nil
	}
	// 0x100aa2eac: mov w9, #0x7bcc
	s.x[9] = uint64(uint32(uint64(0x7bcc)))
	if runtime != nil && runtime.stopPC == 0x100aa2eb0 {
		return 0x100aa2eb0, true, nil
	}
	// 0x100aa2eb0: movk w9, #0xba32, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xba32, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2eb4 {
		return 0x100aa2eb4, true, nil
	}
	// 0x100aa2eb4: eor w9, w3, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2eb8 {
		return 0x100aa2eb8, true, nil
	}
	// 0x100aa2eb8: mov w10, #0x8434
	s.x[10] = uint64(uint32(uint64(0x8434)))
	if runtime != nil && runtime.stopPC == 0x100aa2ebc {
		return 0x100aa2ebc, true, nil
	}
	// 0x100aa2ebc: movk w10, #0x45cd, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x45cd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ec0 {
		return 0x100aa2ec0, true, nil
	}
	// 0x100aa2ec0: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ec4 {
		return 0x100aa2ec4, true, nil
	}
	// 0x100aa2ec4: mov w10, #0x3f82
	s.x[10] = uint64(uint32(uint64(0x3f82)))
	if runtime != nil && runtime.stopPC == 0x100aa2ec8 {
		return 0x100aa2ec8, true, nil
	}
	// 0x100aa2ec8: movk w10, #0xc107, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xc107, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ecc {
		return 0x100aa2ecc, true, nil
	}
	// 0x100aa2ecc: eor w10, w3, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ed0 {
		return 0x100aa2ed0, true, nil
	}
	// 0x100aa2ed0: mov w11, #0xc07e
	s.x[11] = uint64(uint32(uint64(0xc07e)))
	if runtime != nil && runtime.stopPC == 0x100aa2ed4 {
		return 0x100aa2ed4, true, nil
	}
	// 0x100aa2ed4: movk w11, #0x3ef8, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x3ef8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ed8 {
		return 0x100aa2ed8, true, nil
	}
	// 0x100aa2ed8: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2edc {
		return 0x100aa2edc, true, nil
	}
	// 0x100aa2edc: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ee0 {
		return 0x100aa2ee0, true, nil
	}
	// 0x100aa2ee0: mov w10, #0xc12e
	s.x[10] = uint64(uint32(uint64(0xc12e)))
	if runtime != nil && runtime.stopPC == 0x100aa2ee4 {
		return 0x100aa2ee4, true, nil
	}
	// 0x100aa2ee4: movk w10, #0xad16, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xad16, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ee8 {
		return 0x100aa2ee8, true, nil
	}
	// 0x100aa2ee8: eor w10, w3, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2eec {
		return 0x100aa2eec, true, nil
	}
	// 0x100aa2eec: mov w11, #0x3ed2
	s.x[11] = uint64(uint32(uint64(0x3ed2)))
	if runtime != nil && runtime.stopPC == 0x100aa2ef0 {
		return 0x100aa2ef0, true, nil
	}
	// 0x100aa2ef0: movk w11, #0x52e9, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x52e9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2ef4 {
		return 0x100aa2ef4, true, nil
	}
	// 0x100aa2ef4: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100aa2ef8 {
		return 0x100aa2ef8, true, nil
	}
	// 0x100aa2ef8: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2efc {
		return 0x100aa2efc, true, nil
	}
	// 0x100aa2efc: mov w10, #0x8560
	s.x[10] = uint64(uint32(uint64(0x8560)))
	if runtime != nil && runtime.stopPC == 0x100aa2f00 {
		return 0x100aa2f00, true, nil
	}
	// 0x100aa2f00: movk w10, #0xd623, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xd623, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100aa2f04 {
		return 0x100aa2f04, true, nil
	}
	// 0x100aa2f04: ldr w11, [sp, #0x178]
	value442, err := s.load(addOffset(s.sp, 376), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value442))
	if runtime != nil && runtime.stopPC == 0x100aa2f08 {
		return 0x100aa2f08, true, nil
	}
	// 0x100aa2f08: add w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100aa2f0c {
		return 0x100aa2f0c, true, nil
	}
	// 0x100aa2f0c: add w28, w10, w9
	s.x[28] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100aa2f10 {
		return 0x100aa2f10, true, nil
	}
	// 0x100aa2f10: ldr x9, [x13, w22, uxtw #3]
	value445, err := s.load((s.x[13] + shiftLeft(uint64(uint32(uint64(uint32(s.x[22])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value445)
	if runtime != nil && runtime.stopPC == 0x100aa2f14 {
		return 0x100aa2f14, true, nil
	}
	// 0x100aa2f14: mov x30, x13
	s.x[30] = uint64(s.x[13])
	if runtime != nil && runtime.stopPC == 0x100aa2f18 {
		return 0x100aa2f18, true, nil
	}
	// 0x100aa2f18: ldr x22, [sp, #0x1c0]
	value447, err := s.load(addOffset(s.sp, 448), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[22] = uint64(value447)
	if runtime != nil && runtime.stopPC == 0x100aa2f1c {
		return 0x100aa2f1c, true, nil
	}
	// 0x100aa2f1c: ldr x23, [sp, #0x1b0]
	value448, err := s.load(addOffset(s.sp, 432), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[23] = uint64(value448)
	if runtime != nil && runtime.stopPC == 0x100aa2f20 {
		return 0x100aa2f20, true, nil
	}
	// 0x100aa2f20: ldr w24, [sp, #0x11c]
	value449, err := s.load(addOffset(s.sp, 284), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[24] = uint64(uint32(value449))
	if runtime != nil && runtime.stopPC == 0x100aa2f24 {
		return 0x100aa2f24, true, nil
	}
	// 0x100aa2f24: mov x17, #0x261f
	s.x[17] = uint64(uint64(0x261f))
	if runtime != nil && runtime.stopPC == 0x100aa2f28 {
		return 0x100aa2f28, true, nil
	}
	// 0x100aa2f28: braa x9, x17
	next = s.x[9]
	return next, false, nil
}

func rawSourceBootstrapBlock100b0b1d0(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100b0b1d0 {
		return 0x100b0b1d0, true, nil
	}
	// 0x100b0b1d0: mov w9, #-0x35f5
	s.x[9] = uint64(uint32(uint64(0xffffca0b)))
	if runtime != nil && runtime.stopPC == 0x100b0b1d4 {
		return 0x100b0b1d4, true, nil
	}
	// 0x100b0b1d4: add w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b1d8 {
		return 0x100b0b1d8, true, nil
	}
	// 0x100b0b1d8: orr w10, w9, #0x20
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) | uint64(0x20))))
	if runtime != nil && runtime.stopPC == 0x100b0b1dc {
		return 0x100b0b1dc, true, nil
	}
	// 0x100b0b1dc: mov w9, #0x86
	s.x[9] = uint64(uint32(uint64(0x86)))
	if runtime != nil && runtime.stopPC == 0x100b0b1e0 {
		return 0x100b0b1e0, true, nil
	}
	// 0x100b0b1e0: mul w14, w12, w9
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b1e4 {
		return 0x100b0b1e4, true, nil
	}
	// 0x100b0b1e4: mov w9, #0x1b38
	s.x[9] = uint64(uint32(uint64(0x1b38)))
	if runtime != nil && runtime.stopPC == 0x100b0b1e8 {
		return 0x100b0b1e8, true, nil
	}
	// 0x100b0b1e8: add w15, w14, w9
	s.x[15] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b1ec {
		return 0x100b0b1ec, true, nil
	}
	// 0x100b0b1ec: mov w9, #0x7433
	s.x[9] = uint64(uint32(uint64(0x7433)))
	if runtime != nil && runtime.stopPC == 0x100b0b1f0 {
		return 0x100b0b1f0, true, nil
	}
	// 0x100b0b1f0: and w16, w15, #0xfffe
	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xfffe))))
	if runtime != nil && runtime.stopPC == 0x100b0b1f4 {
		return 0x100b0b1f4, true, nil
	}
	// 0x100b0b1f4: mul w16, w16, w9
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b1f8 {
		return 0x100b0b1f8, true, nil
	}
	// 0x100b0b1f8: lsr w16, w16, #0x16
	s.x[16] = uint64(uint32(shiftRight(uint64(uint32(s.x[16])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b1fc {
		return 0x100b0b1fc, true, nil
	}
	// 0x100b0b1fc: mov w5, #0x8d
	s.x[5] = uint64(uint32(uint64(0x8d)))
	if runtime != nil && runtime.stopPC == 0x100b0b200 {
		return 0x100b0b200, true, nil
	}
	// 0x100b0b200: mov w17, #0x2206
	s.x[17] = uint64(uint32(uint64(0x2206)))
	if runtime != nil && runtime.stopPC == 0x100b0b204 {
		return 0x100b0b204, true, nil
	}
	// 0x100b0b204: add w0, w14, w17
	s.x[0] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100b0b208 {
		return 0x100b0b208, true, nil
	}
	// 0x100b0b208: msub w14, w16, w5, w15
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) - uint64(uint32(s.x[16]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100b0b20c {
		return 0x100b0b20c, true, nil
	}
	// 0x100b0b20c: and w15, w0, #0xfffe
	s.x[15] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(0xfffe))))
	if runtime != nil && runtime.stopPC == 0x100b0b210 {
		return 0x100b0b210, true, nil
	}
	// 0x100b0b210: mul w15, w15, w9
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b214 {
		return 0x100b0b214, true, nil
	}
	// 0x100b0b214: lsr w16, w15, #0x16
	s.x[16] = uint64(uint32(shiftRight(uint64(uint32(s.x[15])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b218 {
		return 0x100b0b218, true, nil
	}
	// 0x100b0b218: msub w16, w16, w5, w0
	s.x[16] = uint64(uint32((uint64(uint32(s.x[0])) - uint64(uint32(s.x[16]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100b0b21c {
		return 0x100b0b21c, true, nil
	}
	// 0x100b0b21c: mov w1, #-0x6de0
	s.x[1] = uint64(uint32(uint64(0xffff9220)))
	if runtime != nil && runtime.stopPC == 0x100b0b220 {
		return 0x100b0b220, true, nil
	}
	// 0x100b0b220: ubfiz x2, x14, #2, #0x10
	s.x[2] = uint64(((s.x[14]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100b0b224 {
		return 0x100b0b224, true, nil
	}
	// 0x100b0b224: add w14, w10, w1
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100b0b228 {
		return 0x100b0b228, true, nil
	}
	// 0x100b0b228: mul w14, w14, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) * uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100b0b22c {
		return 0x100b0b22c, true, nil
	}
	// 0x100b0b22c: add w17, w14, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100b0b230 {
		return 0x100b0b230, true, nil
	}
	// 0x100b0b230: and w1, w17, #0xffff
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(0xffff))))
	if runtime != nil && runtime.stopPC == 0x100b0b234 {
		return 0x100b0b234, true, nil
	}
	// 0x100b0b234: ldr w2, [x8, x2]
	value25, err := s.load((s.x[8] + s.x[2]), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[2] = uint64(uint32(value25))
	if runtime != nil && runtime.stopPC == 0x100b0b238 {
		return 0x100b0b238, true, nil
	}
	// 0x100b0b238: mul w1, w1, w9
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b23c {
		return 0x100b0b23c, true, nil
	}
	// 0x100b0b23c: lsr w1, w1, #0x16
	s.x[1] = uint64(uint32(shiftRight(uint64(uint32(s.x[1])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b240 {
		return 0x100b0b240, true, nil
	}
	// 0x100b0b240: msub w17, w1, w5, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) - uint64(uint32(s.x[1]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100b0b244 {
		return 0x100b0b244, true, nil
	}
	// 0x100b0b244: ubfiz x17, x17, #2, #0x10
	s.x[17] = uint64(((s.x[17]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100b0b248 {
		return 0x100b0b248, true, nil
	}
	// 0x100b0b248: ldr w17, [x8, x17]
	value30, err := s.load((s.x[8] + s.x[17]), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value30))
	if runtime != nil && runtime.stopPC == 0x100b0b24c {
		return 0x100b0b24c, true, nil
	}
	// 0x100b0b24c: lsr w15, w15, #0x17
	s.x[15] = uint64(uint32(shiftRight(uint64(uint32(s.x[15])), int(uint64(0x17)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b250 {
		return 0x100b0b250, true, nil
	}
	// 0x100b0b250: mov w1, #0x11a
	s.x[1] = uint64(uint32(uint64(0x11a)))
	if runtime != nil && runtime.stopPC == 0x100b0b254 {
		return 0x100b0b254, true, nil
	}
	// 0x100b0b254: msub w15, w15, w1, w0
	s.x[15] = uint64(uint32((uint64(uint32(s.x[0])) - uint64(uint32(s.x[15]))*uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100b0b258 {
		return 0x100b0b258, true, nil
	}
	// 0x100b0b258: and w15, w15, #0xfffe
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xfffe))))
	if runtime != nil && runtime.stopPC == 0x100b0b25c {
		return 0x100b0b25c, true, nil
	}
	// 0x100b0b25c: subs w0, w15, #0x8d
	s.x[0] = uint64(uint32(s.subtractWithFlags(uint64(uint32(s.x[15])), uint64(0x8d), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b260 {
		return 0x100b0b260, true, nil
	}
	// 0x100b0b260: lsr w1, w2, #3
	s.x[1] = uint64(uint32(shiftRight(uint64(uint32(s.x[2])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b264 {
		return 0x100b0b264, true, nil
	}
	// 0x100b0b264: ubfiz x16, x16, #2, #0x10
	s.x[16] = uint64(((s.x[16]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100b0b268 {
		return 0x100b0b268, true, nil
	}
	// 0x100b0b268: ldr w16, [x8, x16]
	value38, err := s.load((s.x[8] + s.x[16]), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value38))
	if runtime != nil && runtime.stopPC == 0x100b0b26c {
		return 0x100b0b26c, true, nil
	}
	// 0x100b0b26c: csel w15, w15, w0, lo
	s.x[15] = uint64(uint32(selectValue(s.condition("lo"), uint64(uint32(s.x[15])), uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100b0b270 {
		return 0x100b0b270, true, nil
	}
	// 0x100b0b270: ror w0, w2, #0x12
	s.x[0] = uint64(uint32(rotateRight(uint64(uint32(s.x[2])), int(uint64(0x12)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b274 {
		return 0x100b0b274, true, nil
	}
	// 0x100b0b274: ldr w15, [x8, w15, uxtw #2]
	value41, err := s.load((s.x[8] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value41))
	if runtime != nil && runtime.stopPC == 0x100b0b278 {
		return 0x100b0b278, true, nil
	}
	// 0x100b0b278: mov w3, #0x4954
	s.x[3] = uint64(uint32(uint64(0x4954)))
	if runtime != nil && runtime.stopPC == 0x100b0b27c {
		return 0x100b0b27c, true, nil
	}
	// 0x100b0b27c: movk w3, #0x457, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0x457, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b280 {
		return 0x100b0b280, true, nil
	}
	// 0x100b0b280: and w3, w3, w2, lsr #2
	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) & shiftRight(uint64(uint32(s.x[2])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100b0b284 {
		return 0x100b0b284, true, nil
	}
	// 0x100b0b284: mov w4, #0xa4aa
	s.x[4] = uint64(uint32(uint64(0xa4aa)))
	if runtime != nil && runtime.stopPC == 0x100b0b288 {
		return 0x100b0b288, true, nil
	}
	// 0x100b0b288: movk w4, #0x222b, lsl #16
	s.x[4] = uint64(uint32(insertBits(uint64(uint32(s.x[4])), 0x222b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b28c {
		return 0x100b0b28c, true, nil
	}
	// 0x100b0b28c: sub w1, w1, w3
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) - uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100b0b290 {
		return 0x100b0b290, true, nil
	}
	// 0x100b0b290: add w1, w1, w4
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) + uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x100b0b294 {
		return 0x100b0b294, true, nil
	}
	// 0x100b0b294: lsr w3, w16, #0xa
	s.x[3] = uint64(uint32(shiftRight(uint64(uint32(s.x[16])), int(uint64(0xa)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b298 {
		return 0x100b0b298, true, nil
	}
	// 0x100b0b298: mov w6, #0x1f68
	s.x[6] = uint64(uint32(uint64(0x1f68)))
	if runtime != nil && runtime.stopPC == 0x100b0b29c {
		return 0x100b0b29c, true, nil
	}
	// 0x100b0b29c: add w6, w14, w6
	s.x[6] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2a0 {
		return 0x100b0b2a0, true, nil
	}
	// 0x100b0b2a0: and w7, w6, #0xffff
	s.x[7] = uint64(uint32((uint64(uint32(s.x[6])) & uint64(0xffff))))
	if runtime != nil && runtime.stopPC == 0x100b0b2a4 {
		return 0x100b0b2a4, true, nil
	}
	// 0x100b0b2a4: eor w0, w0, w2, ror #7
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ rotateRight(uint64(uint32(s.x[2])), 7, 32))))
	if runtime != nil && runtime.stopPC == 0x100b0b2a8 {
		return 0x100b0b2a8, true, nil
	}
	// 0x100b0b2a8: mul w2, w7, w9
	s.x[2] = uint64(uint32((uint64(uint32(s.x[7])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2ac {
		return 0x100b0b2ac, true, nil
	}
	// 0x100b0b2ac: lsr w2, w2, #0x16
	s.x[2] = uint64(uint32(shiftRight(uint64(uint32(s.x[2])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b2b0 {
		return 0x100b0b2b0, true, nil
	}
	// 0x100b0b2b0: extr w15, w15, w16, #0x13
	s.x[15] = uint64(uint32(extractPair(uint64(uint32(s.x[15])), uint64(uint32(s.x[16])), 19, 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b2b4 {
		return 0x100b0b2b4, true, nil
	}
	// 0x100b0b2b4: extr w17, w17, w16, #0x11
	s.x[17] = uint64(uint32(extractPair(uint64(uint32(s.x[17])), uint64(uint32(s.x[16])), 17, 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b2b8 {
		return 0x100b0b2b8, true, nil
	}
	// 0x100b0b2b8: msub w2, w2, w5, w6
	s.x[2] = uint64(uint32((uint64(uint32(s.x[6])) - uint64(uint32(s.x[2]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2bc {
		return 0x100b0b2bc, true, nil
	}
	// 0x100b0b2bc: mov w6, #0x2312
	s.x[6] = uint64(uint32(uint64(0x2312)))
	if runtime != nil && runtime.stopPC == 0x100b0b2c0 {
		return 0x100b0b2c0, true, nil
	}
	// 0x100b0b2c0: add w6, w14, w6
	s.x[6] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2c4 {
		return 0x100b0b2c4, true, nil
	}
	// 0x100b0b2c4: and w7, w6, #0xffff
	s.x[7] = uint64(uint32((uint64(uint32(s.x[6])) & uint64(0xffff))))
	if runtime != nil && runtime.stopPC == 0x100b0b2c8 {
		return 0x100b0b2c8, true, nil
	}
	// 0x100b0b2c8: mul w7, w7, w9
	s.x[7] = uint64(uint32((uint64(uint32(s.x[7])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2cc {
		return 0x100b0b2cc, true, nil
	}
	// 0x100b0b2cc: eor w0, w0, w1
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2d0 {
		return 0x100b0b2d0, true, nil
	}
	// 0x100b0b2d0: lsr w1, w7, #0x16
	s.x[1] = uint64(uint32(shiftRight(uint64(uint32(s.x[7])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b2d4 {
		return 0x100b0b2d4, true, nil
	}
	// 0x100b0b2d4: msub w1, w1, w5, w6
	s.x[1] = uint64(uint32((uint64(uint32(s.x[6])) - uint64(uint32(s.x[1]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2d8 {
		return 0x100b0b2d8, true, nil
	}
	// 0x100b0b2d8: mov w6, #0x7d72
	s.x[6] = uint64(uint32(uint64(0x7d72)))
	if runtime != nil && runtime.stopPC == 0x100b0b2dc {
		return 0x100b0b2dc, true, nil
	}
	// 0x100b0b2dc: movk w6, #0x54, lsl #16
	s.x[6] = uint64(uint32(insertBits(uint64(uint32(s.x[6])), 0x54, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b2e0 {
		return 0x100b0b2e0, true, nil
	}
	// 0x100b0b2e0: and w16, w6, w16, lsr #9
	s.x[16] = uint64(uint32((uint64(uint32(s.x[6])) & shiftRight(uint64(uint32(s.x[16])), 9, 32))))
	if runtime != nil && runtime.stopPC == 0x100b0b2e4 {
		return 0x100b0b2e4, true, nil
	}
	// 0x100b0b2e4: mov w6, #0x3eb9
	s.x[6] = uint64(uint32(uint64(0x3eb9)))
	if runtime != nil && runtime.stopPC == 0x100b0b2e8 {
		return 0x100b0b2e8, true, nil
	}
	// 0x100b0b2e8: movk w6, #0x42aa, lsl #16
	s.x[6] = uint64(uint32(insertBits(uint64(uint32(s.x[6])), 0x42aa, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b2ec {
		return 0x100b0b2ec, true, nil
	}
	// 0x100b0b2ec: eor w0, w0, w4
	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2f0 {
		return 0x100b0b2f0, true, nil
	}
	// 0x100b0b2f0: sub w16, w3, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[3])) - uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2f4 {
		return 0x100b0b2f4, true, nil
	}
	// 0x100b0b2f4: add w16, w16, w6
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2f8 {
		return 0x100b0b2f8, true, nil
	}
	// 0x100b0b2f8: eor w15, w17, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100b0b2fc {
		return 0x100b0b2fc, true, nil
	}
	// 0x100b0b2fc: eor w16, w16, w6
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100b0b300 {
		return 0x100b0b300, true, nil
	}
	// 0x100b0b300: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100b0b304 {
		return 0x100b0b304, true, nil
	}
	// 0x100b0b304: add w15, w15, w0
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100b0b308 {
		return 0x100b0b308, true, nil
	}
	// 0x100b0b308: ubfiz x16, x1, #2, #0x10
	s.x[16] = uint64(((s.x[1]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100b0b30c {
		return 0x100b0b30c, true, nil
	}
	// 0x100b0b30c: str w15, [x8, x16]
	if err := s.storeValue((s.x[8] + s.x[16]), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100b0b310 {
		return 0x100b0b310, true, nil
	}
	// 0x100b0b310: ubfiz x17, x2, #2, #0x10
	s.x[17] = uint64(((s.x[2]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100b0b314 {
		return 0x100b0b314, true, nil
	}
	// 0x100b0b314: ldr w17, [x8, x17]
	value81, err := s.load((s.x[8] + s.x[17]), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value81))
	if runtime != nil && runtime.stopPC == 0x100b0b318 {
		return 0x100b0b318, true, nil
	}
	// 0x100b0b318: mov w0, #0x1ab2
	s.x[0] = uint64(uint32(uint64(0x1ab2)))
	if runtime != nil && runtime.stopPC == 0x100b0b31c {
		return 0x100b0b31c, true, nil
	}
	// 0x100b0b31c: add w14, w14, w0
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100b0b320 {
		return 0x100b0b320, true, nil
	}
	// 0x100b0b320: and w0, w14, #0xffff
	s.x[0] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(0xffff))))
	if runtime != nil && runtime.stopPC == 0x100b0b324 {
		return 0x100b0b324, true, nil
	}
	// 0x100b0b324: mul w9, w0, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[0])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b328 {
		return 0x100b0b328, true, nil
	}
	// 0x100b0b328: lsr w9, w9, #0x16
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x16)), 32)))
	if runtime != nil && runtime.stopPC == 0x100b0b32c {
		return 0x100b0b32c, true, nil
	}
	// 0x100b0b32c: msub w9, w9, w5, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) - uint64(uint32(s.x[9]))*uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100b0b330 {
		return 0x100b0b330, true, nil
	}
	// 0x100b0b330: add w14, w15, w17
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100b0b334 {
		return 0x100b0b334, true, nil
	}
	// 0x100b0b334: add w12, w12, #1
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0x1))))
	if runtime != nil && runtime.stopPC == 0x100b0b338 {
		return 0x100b0b338, true, nil
	}
	// 0x100b0b338: ubfiz x9, x9, #2, #0x10
	s.x[9] = uint64(((s.x[9]) & bitMask(16, 64)) << 2)
	if runtime != nil && runtime.stopPC == 0x100b0b33c {
		return 0x100b0b33c, true, nil
	}
	// 0x100b0b33c: str w14, [x8, x16]
	if err := s.storeValue((s.x[8] + s.x[16]), uint64(uint32(s.x[14])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100b0b340 {
		return 0x100b0b340, true, nil
	}
	// 0x100b0b340: ldr w9, [x8, x9]
	value92, err := s.load((s.x[8] + s.x[9]), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value92))
	if runtime != nil && runtime.stopPC == 0x100b0b344 {
		return 0x100b0b344, true, nil
	}
	// 0x100b0b344: add w9, w14, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100b0b348 {
		return 0x100b0b348, true, nil
	}
	// 0x100b0b348: cmp w12, #0x40
	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(0x40)))
	if runtime != nil && runtime.stopPC == 0x100b0b34c {
		return 0x100b0b34c, true, nil
	}
	// 0x100b0b34c: str w9, [x8, x16]
	if err := s.storeValue((s.x[8] + s.x[16]), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100b0b350 {
		return 0x100b0b350, true, nil
	}
	// 0x100b0b350: cset w9, eq
	s.x[9] = uint64(uint32(boolValue(s.condition("eq"))))
	if runtime != nil && runtime.stopPC == 0x100b0b354 {
		return 0x100b0b354, true, nil
	}
	// 0x100b0b354: mov w14, #0x3bd
	s.x[14] = uint64(uint32(uint64(0x3bd)))
	if runtime != nil && runtime.stopPC == 0x100b0b358 {
		return 0x100b0b358, true, nil
	}
	// 0x100b0b358: mul w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100b0b35c {
		return 0x100b0b35c, true, nil
	}
	// 0x100b0b35c: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100b0b360 {
		return 0x100b0b360, true, nil
	}
	// 0x100b0b360: ldr x9, [x30, w9, sxtw #3]
	value100, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value100)
	if runtime != nil && runtime.stopPC == 0x100b0b364 {
		return 0x100b0b364, true, nil
	}
	// 0x100b0b364: mov x17, #0x3e0f
	s.x[17] = uint64(uint64(0x3e0f))
	if runtime != nil && runtime.stopPC == 0x100b0b368 {
		return 0x100b0b368, true, nil
	}
	// 0x100b0b368: braa x9, x17
	next = s.x[9]
	return next, false, nil
}
