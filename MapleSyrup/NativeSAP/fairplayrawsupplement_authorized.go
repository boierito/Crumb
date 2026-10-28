// Ahead-of-time translations of public pinned account-state processing code.
// Source binary SHA-256: 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784
package rawsource

func rawSourceBootstrapBlock100234cdc(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100234cdc {
		return 0x100234cdc, true, nil
	}
	// 0x100234cdc: ldrb w8, [x28, x8]
	value0, err := s.load((s.x[28] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value0))
	if runtime != nil && runtime.stopPC == 0x100234ce0 {
		return 0x100234ce0, true, nil
	}
	// 0x100234ce0: mov w9, #0x292c
	s.x[9] = uint64(uint32(uint64(0x292c)))
	if runtime != nil && runtime.stopPC == 0x100234ce4 {
		return 0x100234ce4, true, nil
	}
	// 0x100234ce4: movk w9, #0x29ad, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x29ad, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234ce8 {
		return 0x100234ce8, true, nil
	}
	// 0x100234ce8: add w11, w10, w9
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234cec {
		return 0x100234cec, true, nil
	}
	// 0x100234cec: mov w9, #0x1a
	s.x[9] = uint64(uint32(uint64(0x1a)))
	if runtime != nil && runtime.stopPC == 0x100234cf0 {
		return 0x100234cf0, true, nil
	}
	// 0x100234cf0: str w11, [sp, #0x404]
	if err := s.storeValue(addOffset(s.sp, 1028), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234cf4 {
		return 0x100234cf4, true, nil
	}
	// 0x100234cf4: eor w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234cf8 {
		return 0x100234cf8, true, nil
	}
	// 0x100234cf8: mov w10, #0xc5
	s.x[10] = uint64(uint32(uint64(0xc5)))
	if runtime != nil && runtime.stopPC == 0x100234cfc {
		return 0x100234cfc, true, nil
	}
	// 0x100234cfc: eor x10, x8, x10
	s.x[10] = uint64((s.x[8] ^ s.x[10]))
	if runtime != nil && runtime.stopPC == 0x100234d00 {
		return 0x100234d00, true, nil
	}
	// 0x100234d00: ldr x21, [sp, #0x398]
	value9, err := s.load(addOffset(s.sp, 920), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[21] = uint64(value9)
	if runtime != nil && runtime.stopPC == 0x100234d04 {
		return 0x100234d04, true, nil
	}
	// 0x100234d04: ldrb w10, [x21, x10]
	value10, err := s.load((s.x[21] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value10))
	if runtime != nil && runtime.stopPC == 0x100234d08 {
		return 0x100234d08, true, nil
	}
	// 0x100234d08: eor w9, w9, w8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100234d0c {
		return 0x100234d0c, true, nil
	}
	// 0x100234d0c: mov w11, #0x48
	s.x[11] = uint64(uint32(uint64(0x48)))
	if runtime != nil && runtime.stopPC == 0x100234d10 {
		return 0x100234d10, true, nil
	}
	// 0x100234d10: eor w8, w8, w11
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100234d14 {
		return 0x100234d14, true, nil
	}
	// 0x100234d14: and w11, w12, #0xfffffffd
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xfffffffd))))
	if runtime != nil && runtime.stopPC == 0x100234d18 {
		return 0x100234d18, true, nil
	}
	// 0x100234d18: mov w7, #0x911e
	s.x[7] = uint64(uint32(uint64(0x911e)))
	if runtime != nil && runtime.stopPC == 0x100234d1c {
		return 0x100234d1c, true, nil
	}
	// 0x100234d1c: movk w7, #0x5ca7, lsl #16
	s.x[7] = uint64(uint32(insertBits(uint64(uint32(s.x[7])), 0x5ca7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234d20 {
		return 0x100234d20, true, nil
	}
	// 0x100234d20: eor x11, x11, x7
	s.x[11] = uint64((s.x[11] ^ s.x[7]))
	if runtime != nil && runtime.stopPC == 0x100234d24 {
		return 0x100234d24, true, nil
	}
	// 0x100234d24: str x11, [sp, #0x478]
	if err := s.storeValue(addOffset(s.sp, 1144), s.x[11], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234d28 {
		return 0x100234d28, true, nil
	}
	// 0x100234d28: ldrb w11, [x28, x11]
	value19, err := s.load((s.x[28] + s.x[11]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value19))
	if runtime != nil && runtime.stopPC == 0x100234d2c {
		return 0x100234d2c, true, nil
	}
	// 0x100234d2c: mov w13, #0x32
	s.x[13] = uint64(uint32(uint64(0x32)))
	if runtime != nil && runtime.stopPC == 0x100234d30 {
		return 0x100234d30, true, nil
	}
	// 0x100234d30: eor x11, x11, x13
	s.x[11] = uint64((s.x[11] ^ s.x[13]))
	if runtime != nil && runtime.stopPC == 0x100234d34 {
		return 0x100234d34, true, nil
	}
	// 0x100234d34: ldr x20, [sp, #0x3a0]
	value22, err := s.load(addOffset(s.sp, 928), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[20] = uint64(value22)
	if runtime != nil && runtime.stopPC == 0x100234d38 {
		return 0x100234d38, true, nil
	}
	// 0x100234d38: ldr w11, [x20, x11, lsl #2]
	value23, err := s.load((s.x[20] + shiftLeft(s.x[11], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value23))
	if runtime != nil && runtime.stopPC == 0x100234d3c {
		return 0x100234d3c, true, nil
	}
	// 0x100234d3c: ror w13, w11, #4
	s.x[13] = uint64(uint32(rotateRight(uint64(uint32(s.x[11])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x100234d40 {
		return 0x100234d40, true, nil
	}
	// 0x100234d40: add w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234d44 {
		return 0x100234d44, true, nil
	}
	// 0x100234d44: add w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234d48 {
		return 0x100234d48, true, nil
	}
	// 0x100234d48: mov w9, #-0xe
	s.x[9] = uint64(uint32(uint64(0xfffffff2)))
	if runtime != nil && runtime.stopPC == 0x100234d4c {
		return 0x100234d4c, true, nil
	}
	// 0x100234d4c: and w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234d50 {
		return 0x100234d50, true, nil
	}
	// 0x100234d50: mov w19, #0x911d
	s.x[19] = uint64(uint32(uint64(0x911d)))
	if runtime != nil && runtime.stopPC == 0x100234d54 {
		return 0x100234d54, true, nil
	}
	// 0x100234d54: movk w19, #0x5ca7, lsl #16
	s.x[19] = uint64(uint32(insertBits(uint64(uint32(s.x[19])), 0x5ca7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234d58 {
		return 0x100234d58, true, nil
	}
	// 0x100234d58: eor x10, x9, x19
	s.x[10] = uint64((s.x[9] ^ s.x[19]))
	if runtime != nil && runtime.stopPC == 0x100234d5c {
		return 0x100234d5c, true, nil
	}
	// 0x100234d5c: str x10, [sp, #0x4a0]
	if err := s.storeValue(addOffset(s.sp, 1184), s.x[10], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234d60 {
		return 0x100234d60, true, nil
	}
	// 0x100234d60: mov w9, #-0xb
	s.x[9] = uint64(uint32(uint64(0xfffffff5)))
	if runtime != nil && runtime.stopPC == 0x100234d64 {
		return 0x100234d64, true, nil
	}
	// 0x100234d64: and w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234d68 {
		return 0x100234d68, true, nil
	}
	// 0x100234d68: ldrb w10, [x28, x10]
	value35, err := s.load((s.x[28] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value35))
	if runtime != nil && runtime.stopPC == 0x100234d6c {
		return 0x100234d6c, true, nil
	}
	// 0x100234d6c: eor x11, x9, x7
	s.x[11] = uint64((s.x[9] ^ s.x[7]))
	if runtime != nil && runtime.stopPC == 0x100234d70 {
		return 0x100234d70, true, nil
	}
	// 0x100234d70: str x11, [sp, #0x488]
	if err := s.storeValue(addOffset(s.sp, 1160), s.x[11], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234d74 {
		return 0x100234d74, true, nil
	}
	// 0x100234d74: mov w9, #0x5f
	s.x[9] = uint64(uint32(uint64(0x5f)))
	if runtime != nil && runtime.stopPC == 0x100234d78 {
		return 0x100234d78, true, nil
	}
	// 0x100234d78: eor x9, x10, x9
	s.x[9] = uint64((s.x[10] ^ s.x[9]))
	if runtime != nil && runtime.stopPC == 0x100234d7c {
		return 0x100234d7c, true, nil
	}
	// 0x100234d7c: ldrb w10, [x28, x11]
	value40, err := s.load((s.x[28] + s.x[11]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value40))
	if runtime != nil && runtime.stopPC == 0x100234d80 {
		return 0x100234d80, true, nil
	}
	// 0x100234d80: ldr x22, [sp, #0x390]
	value41, err := s.load(addOffset(s.sp, 912), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[22] = uint64(value41)
	if runtime != nil && runtime.stopPC == 0x100234d84 {
		return 0x100234d84, true, nil
	}
	// 0x100234d84: ldr w15, [x22, x9, lsl #2]
	value42, err := s.load((s.x[22] + shiftLeft(s.x[9], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value42))
	if runtime != nil && runtime.stopPC == 0x100234d88 {
		return 0x100234d88, true, nil
	}
	// 0x100234d88: sub w8, w8, #9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(0x9))))
	if runtime != nil && runtime.stopPC == 0x100234d8c {
		return 0x100234d8c, true, nil
	}
	// 0x100234d8c: mov w9, #0x3a
	s.x[9] = uint64(uint32(uint64(0x3a)))
	if runtime != nil && runtime.stopPC == 0x100234d90 {
		return 0x100234d90, true, nil
	}
	// 0x100234d90: eor x9, x10, x9
	s.x[9] = uint64((s.x[10] ^ s.x[9]))
	if runtime != nil && runtime.stopPC == 0x100234d94 {
		return 0x100234d94, true, nil
	}
	// 0x100234d94: ldr w9, [x20, x9, lsl #2]
	value46, err := s.load((s.x[20] + shiftLeft(s.x[9], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value46))
	if runtime != nil && runtime.stopPC == 0x100234d98 {
		return 0x100234d98, true, nil
	}
	// 0x100234d98: and w17, w8, #0xff
	s.x[17] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100234d9c {
		return 0x100234d9c, true, nil
	}
	// 0x100234d9c: mov w8, #-0xc
	s.x[8] = uint64(uint32(uint64(0xfffffff4)))
	if runtime != nil && runtime.stopPC == 0x100234da0 {
		return 0x100234da0, true, nil
	}
	// 0x100234da0: and w8, w12, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100234da4 {
		return 0x100234da4, true, nil
	}
	// 0x100234da4: eor x10, x8, x27
	s.x[10] = uint64((s.x[8] ^ s.x[27]))
	if runtime != nil && runtime.stopPC == 0x100234da8 {
		return 0x100234da8, true, nil
	}
	// 0x100234da8: str x10, [sp, #0x4a8]
	if err := s.storeValue(addOffset(s.sp, 1192), s.x[10], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234dac {
		return 0x100234dac, true, nil
	}
	// 0x100234dac: mov w8, #-6
	s.x[8] = uint64(uint32(uint64(0xfffffffa)))
	if runtime != nil && runtime.stopPC == 0x100234db0 {
		return 0x100234db0, true, nil
	}
	// 0x100234db0: and w8, w12, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100234db4 {
		return 0x100234db4, true, nil
	}
	// 0x100234db4: eor x30, x8, x19
	s.x[30] = uint64((s.x[8] ^ s.x[19]))
	if runtime != nil && runtime.stopPC == 0x100234db8 {
		return 0x100234db8, true, nil
	}
	// 0x100234db8: str x30, [sp, #0x438]
	if err := s.storeValue(addOffset(s.sp, 1080), s.x[30], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234dbc {
		return 0x100234dbc, true, nil
	}
	// 0x100234dbc: ror w11, w9, #4
	s.x[11] = uint64(uint32(rotateRight(uint64(uint32(s.x[9])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x100234dc0 {
		return 0x100234dc0, true, nil
	}
	// 0x100234dc0: mov w26, #0xf730
	s.x[26] = uint64(uint32(uint64(0xf730)))
	if runtime != nil && runtime.stopPC == 0x100234dc4 {
		return 0x100234dc4, true, nil
	}
	// 0x100234dc4: movk w26, #0xbb2c, lsl #16
	s.x[26] = uint64(uint32(insertBits(uint64(uint32(s.x[26])), 0xbb2c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234dc8 {
		return 0x100234dc8, true, nil
	}
	// 0x100234dc8: eor w8, w26, w11, lsl #4
	s.x[8] = uint64(uint32((uint64(uint32(s.x[26])) ^ shiftLeft(uint64(uint32(s.x[11])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100234dcc {
		return 0x100234dcc, true, nil
	}
	// 0x100234dcc: ldrb w9, [x28, x10]
	value60, err := s.load((s.x[28] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value60))
	if runtime != nil && runtime.stopPC == 0x100234dd0 {
		return 0x100234dd0, true, nil
	}
	// 0x100234dd0: mov w10, #0x43
	s.x[10] = uint64(uint32(uint64(0x43)))
	if runtime != nil && runtime.stopPC == 0x100234dd4 {
		return 0x100234dd4, true, nil
	}
	// 0x100234dd4: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234dd8 {
		return 0x100234dd8, true, nil
	}
	// 0x100234dd8: mov x4, x12
	s.x[4] = uint64(s.x[12])
	if runtime != nil && runtime.stopPC == 0x100234ddc {
		return 0x100234ddc, true, nil
	}
	// 0x100234ddc: mov w12, #0x26
	s.x[12] = uint64(uint32(uint64(0x26)))
	if runtime != nil && runtime.stopPC == 0x100234de0 {
		return 0x100234de0, true, nil
	}
	// 0x100234de0: eor x12, x9, x12
	s.x[12] = uint64((s.x[9] ^ s.x[12]))
	if runtime != nil && runtime.stopPC == 0x100234de4 {
		return 0x100234de4, true, nil
	}
	// 0x100234de4: ldrb w12, [x21, x12]
	value66, err := s.load((s.x[21] + s.x[12]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value66))
	if runtime != nil && runtime.stopPC == 0x100234de8 {
		return 0x100234de8, true, nil
	}
	// 0x100234de8: mov w14, #0xab
	s.x[14] = uint64(uint32(uint64(0xab)))
	if runtime != nil && runtime.stopPC == 0x100234dec {
		return 0x100234dec, true, nil
	}
	// 0x100234dec: eor w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100234df0 {
		return 0x100234df0, true, nil
	}
	// 0x100234df0: add w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234df4 {
		return 0x100234df4, true, nil
	}
	// 0x100234df4: add w9, w9, w12
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234df8 {
		return 0x100234df8, true, nil
	}
	// 0x100234df8: sub w9, w9, #9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(0x9))))
	if runtime != nil && runtime.stopPC == 0x100234dfc {
		return 0x100234dfc, true, nil
	}
	// 0x100234dfc: and w10, w4, #0xfffffff7
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffff7))))
	if runtime != nil && runtime.stopPC == 0x100234e00 {
		return 0x100234e00, true, nil
	}
	// 0x100234e00: mov w1, #0x911c
	s.x[1] = uint64(uint32(uint64(0x911c)))
	if runtime != nil && runtime.stopPC == 0x100234e04 {
		return 0x100234e04, true, nil
	}
	// 0x100234e04: movk w1, #0x5ca7, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x5ca7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e08 {
		return 0x100234e08, true, nil
	}
	// 0x100234e08: eor x10, x10, x1
	s.x[10] = uint64((s.x[10] ^ s.x[1]))
	if runtime != nil && runtime.stopPC == 0x100234e0c {
		return 0x100234e0c, true, nil
	}
	// 0x100234e0c: str x10, [sp, #0x480]
	if err := s.storeValue(addOffset(s.sp, 1152), s.x[10], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234e10 {
		return 0x100234e10, true, nil
	}
	// 0x100234e10: ldrb w10, [x28, x10]
	value77, err := s.load((s.x[28] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value77))
	if runtime != nil && runtime.stopPC == 0x100234e14 {
		return 0x100234e14, true, nil
	}
	// 0x100234e14: mov w12, #0x8a
	s.x[12] = uint64(uint32(uint64(0x8a)))
	if runtime != nil && runtime.stopPC == 0x100234e18 {
		return 0x100234e18, true, nil
	}
	// 0x100234e18: eor x10, x10, x12
	s.x[10] = uint64((s.x[10] ^ s.x[12]))
	if runtime != nil && runtime.stopPC == 0x100234e1c {
		return 0x100234e1c, true, nil
	}
	// 0x100234e1c: ldr x6, [sp, #0x3a8]
	value80, err := s.load(addOffset(s.sp, 936), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(value80)
	if runtime != nil && runtime.stopPC == 0x100234e20 {
		return 0x100234e20, true, nil
	}
	// 0x100234e20: ldr w10, [x6, x10, lsl #2]
	value81, err := s.load((s.x[6] + shiftLeft(s.x[10], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value81))
	if runtime != nil && runtime.stopPC == 0x100234e24 {
		return 0x100234e24, true, nil
	}
	// 0x100234e24: mov w12, #0x8683
	s.x[12] = uint64(uint32(uint64(0x8683)))
	if runtime != nil && runtime.stopPC == 0x100234e28 {
		return 0x100234e28, true, nil
	}
	// 0x100234e28: movk w12, #0xda0d, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xda0d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e2c {
		return 0x100234e2c, true, nil
	}
	// 0x100234e2c: eor w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234e30 {
		return 0x100234e30, true, nil
	}
	// 0x100234e30: mov w14, #0x797d
	s.x[14] = uint64(uint32(uint64(0x797d)))
	if runtime != nil && runtime.stopPC == 0x100234e34 {
		return 0x100234e34, true, nil
	}
	// 0x100234e34: movk w14, #0x25f2, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x25f2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e38 {
		return 0x100234e38, true, nil
	}
	// 0x100234e38: add w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100234e3c {
		return 0x100234e3c, true, nil
	}
	// 0x100234e3c: mov w14, #0xfcbe
	s.x[14] = uint64(uint32(uint64(0xfcbe)))
	if runtime != nil && runtime.stopPC == 0x100234e40 {
		return 0x100234e40, true, nil
	}
	// 0x100234e40: movk w14, #0xaf94, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xaf94, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e44 {
		return 0x100234e44, true, nil
	}
	// 0x100234e44: eor w14, w10, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100234e48 {
		return 0x100234e48, true, nil
	}
	// 0x100234e48: mov w16, #0x342
	s.x[16] = uint64(uint32(uint64(0x342)))
	if runtime != nil && runtime.stopPC == 0x100234e4c {
		return 0x100234e4c, true, nil
	}
	// 0x100234e4c: movk w16, #0x506b, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x506b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e50 {
		return 0x100234e50, true, nil
	}
	// 0x100234e50: add w14, w14, w16
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100234e54 {
		return 0x100234e54, true, nil
	}
	// 0x100234e54: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100234e58 {
		return 0x100234e58, true, nil
	}
	// 0x100234e58: mov w14, #0x705a
	s.x[14] = uint64(uint32(uint64(0x705a)))
	if runtime != nil && runtime.stopPC == 0x100234e5c {
		return 0x100234e5c, true, nil
	}
	// 0x100234e5c: movk w14, #0x8f40, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x8f40, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e60 {
		return 0x100234e60, true, nil
	}
	// 0x100234e60: eor w14, w10, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100234e64 {
		return 0x100234e64, true, nil
	}
	// 0x100234e64: mov w16, #0x8fa6
	s.x[16] = uint64(uint32(uint64(0x8fa6)))
	if runtime != nil && runtime.stopPC == 0x100234e68 {
		return 0x100234e68, true, nil
	}
	// 0x100234e68: movk w16, #0x70bf, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x70bf, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e6c {
		return 0x100234e6c, true, nil
	}
	// 0x100234e6c: add w14, w14, w16
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100234e70 {
		return 0x100234e70, true, nil
	}
	// 0x100234e70: eor w14, w12, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100234e74 {
		return 0x100234e74, true, nil
	}
	// 0x100234e74: and w16, w9, #0xff
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100234e78 {
		return 0x100234e78, true, nil
	}
	// 0x100234e78: and w9, w4, #0xfffffff3
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffff3))))
	if runtime != nil && runtime.stopPC == 0x100234e7c {
		return 0x100234e7c, true, nil
	}
	// 0x100234e7c: eor x9, x9, x1
	s.x[9] = uint64((s.x[9] ^ s.x[1]))
	if runtime != nil && runtime.stopPC == 0x100234e80 {
		return 0x100234e80, true, nil
	}
	// 0x100234e80: str x9, [sp, #0x490]
	if err := s.storeValue(addOffset(s.sp, 1168), s.x[9], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234e84 {
		return 0x100234e84, true, nil
	}
	// 0x100234e84: ldrb w9, [x28, x9]
	value106, err := s.load((s.x[28] + s.x[9]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value106))
	if runtime != nil && runtime.stopPC == 0x100234e88 {
		return 0x100234e88, true, nil
	}
	// 0x100234e88: mov w23, #0x299c
	s.x[23] = uint64(uint32(uint64(0x299c)))
	if runtime != nil && runtime.stopPC == 0x100234e8c {
		return 0x100234e8c, true, nil
	}
	// 0x100234e8c: movk w23, #0xeb64, lsl #16
	s.x[23] = uint64(uint32(insertBits(uint64(uint32(s.x[23])), 0xeb64, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234e90 {
		return 0x100234e90, true, nil
	}
	// 0x100234e90: eor w0, w23, w10, lsl #2
	s.x[0] = uint64(uint32((uint64(uint32(s.x[23])) ^ shiftLeft(uint64(uint32(s.x[10])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100234e94 {
		return 0x100234e94, true, nil
	}
	// 0x100234e94: lsl w10, w9, #1
	s.x[10] = uint64(uint32(shiftLeft(uint64(uint32(s.x[9])), int(uint64(0x1)), 32)))
	if runtime != nil && runtime.stopPC == 0x100234e98 {
		return 0x100234e98, true, nil
	}
	// 0x100234e98: and w10, w10, #0xf8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0xf8))))
	if runtime != nil && runtime.stopPC == 0x100234e9c {
		return 0x100234e9c, true, nil
	}
	// 0x100234e9c: sub w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234ea0 {
		return 0x100234ea0, true, nil
	}
	// 0x100234ea0: add w9, w9, #0x7c
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0x7c))))
	if runtime != nil && runtime.stopPC == 0x100234ea4 {
		return 0x100234ea4, true, nil
	}
	// 0x100234ea4: mov w10, #0xd9
	s.x[10] = uint64(uint32(uint64(0xd9)))
	if runtime != nil && runtime.stopPC == 0x100234ea8 {
		return 0x100234ea8, true, nil
	}
	// 0x100234ea8: eor x9, x9, x10
	s.x[9] = uint64((s.x[9] ^ s.x[10]))
	if runtime != nil && runtime.stopPC == 0x100234eac {
		return 0x100234eac, true, nil
	}
	// 0x100234eac: ldr w9, [x6, x9, lsl #2]
	value116, err := s.load((s.x[6] + shiftLeft(s.x[9], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value116))
	if runtime != nil && runtime.stopPC == 0x100234eb0 {
		return 0x100234eb0, true, nil
	}
	// 0x100234eb0: mov w10, #0x6c1f
	s.x[10] = uint64(uint32(uint64(0x6c1f)))
	if runtime != nil && runtime.stopPC == 0x100234eb4 {
		return 0x100234eb4, true, nil
	}
	// 0x100234eb4: movk w10, #0x6f6f, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x6f6f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234eb8 {
		return 0x100234eb8, true, nil
	}
	// 0x100234eb8: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234ebc {
		return 0x100234ebc, true, nil
	}
	// 0x100234ebc: mov w12, #0x93e1
	s.x[12] = uint64(uint32(uint64(0x93e1)))
	if runtime != nil && runtime.stopPC == 0x100234ec0 {
		return 0x100234ec0, true, nil
	}
	// 0x100234ec0: movk w12, #0x9090, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x9090, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234ec4 {
		return 0x100234ec4, true, nil
	}
	// 0x100234ec4: add w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234ec8 {
		return 0x100234ec8, true, nil
	}
	// 0x100234ec8: mov w12, #0x9126
	s.x[12] = uint64(uint32(uint64(0x9126)))
	if runtime != nil && runtime.stopPC == 0x100234ecc {
		return 0x100234ecc, true, nil
	}
	// 0x100234ecc: movk w12, #0xc83a, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xc83a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234ed0 {
		return 0x100234ed0, true, nil
	}
	// 0x100234ed0: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234ed4 {
		return 0x100234ed4, true, nil
	}
	// 0x100234ed4: mov w2, #0x6eda
	s.x[2] = uint64(uint32(uint64(0x6eda)))
	if runtime != nil && runtime.stopPC == 0x100234ed8 {
		return 0x100234ed8, true, nil
	}
	// 0x100234ed8: movk w2, #0x37c5, lsl #16
	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0x37c5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234edc {
		return 0x100234edc, true, nil
	}
	// 0x100234edc: add w12, w12, w2
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[2])))))
	if runtime != nil && runtime.stopPC == 0x100234ee0 {
		return 0x100234ee0, true, nil
	}
	// 0x100234ee0: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234ee4 {
		return 0x100234ee4, true, nil
	}
	// 0x100234ee4: mov w12, #0xf75e
	s.x[12] = uint64(uint32(uint64(0xf75e)))
	if runtime != nil && runtime.stopPC == 0x100234ee8 {
		return 0x100234ee8, true, nil
	}
	// 0x100234ee8: movk w12, #0x5d8c, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x5d8c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234eec {
		return 0x100234eec, true, nil
	}
	// 0x100234eec: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234ef0 {
		return 0x100234ef0, true, nil
	}
	// 0x100234ef0: mov w2, #0x8a2
	s.x[2] = uint64(uint32(uint64(0x8a2)))
	if runtime != nil && runtime.stopPC == 0x100234ef4 {
		return 0x100234ef4, true, nil
	}
	// 0x100234ef4: movk w2, #0xa273, lsl #16
	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0xa273, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234ef8 {
		return 0x100234ef8, true, nil
	}
	// 0x100234ef8: add w12, w12, w2
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[2])))))
	if runtime != nil && runtime.stopPC == 0x100234efc {
		return 0x100234efc, true, nil
	}
	// 0x100234efc: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234f00 {
		return 0x100234f00, true, nil
	}
	// 0x100234f00: eor w9, w23, w9, lsl #2
	s.x[9] = uint64(uint32((uint64(uint32(s.x[23])) ^ shiftLeft(uint64(uint32(s.x[9])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100234f04 {
		return 0x100234f04, true, nil
	}
	// 0x100234f04: and w12, w4, #0xfffffffb
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffffb))))
	if runtime != nil && runtime.stopPC == 0x100234f08 {
		return 0x100234f08, true, nil
	}
	// 0x100234f08: eor x12, x12, x1
	s.x[12] = uint64((s.x[12] ^ s.x[1]))
	if runtime != nil && runtime.stopPC == 0x100234f0c {
		return 0x100234f0c, true, nil
	}
	// 0x100234f0c: str x12, [sp, #0x448]
	if err := s.storeValue(addOffset(s.sp, 1096), s.x[12], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234f10 {
		return 0x100234f10, true, nil
	}
	// 0x100234f10: ldrb w12, [x28, x12]
	value141, err := s.load((s.x[28] + s.x[12]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value141))
	if runtime != nil && runtime.stopPC == 0x100234f14 {
		return 0x100234f14, true, nil
	}
	// 0x100234f14: mov w2, #0x42
	s.x[2] = uint64(uint32(uint64(0x42)))
	if runtime != nil && runtime.stopPC == 0x100234f18 {
		return 0x100234f18, true, nil
	}
	// 0x100234f18: eor x12, x12, x2
	s.x[12] = uint64((s.x[12] ^ s.x[2]))
	if runtime != nil && runtime.stopPC == 0x100234f1c {
		return 0x100234f1c, true, nil
	}
	// 0x100234f1c: ldr w2, [x6, x12, lsl #2]
	value144, err := s.load((s.x[6] + shiftLeft(s.x[12], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[2] = uint64(uint32(value144))
	if runtime != nil && runtime.stopPC == 0x100234f20 {
		return 0x100234f20, true, nil
	}
	// 0x100234f20: and w12, w4, #0xfffffff1
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffff1))))
	if runtime != nil && runtime.stopPC == 0x100234f24 {
		return 0x100234f24, true, nil
	}
	// 0x100234f24: eor x24, x12, x7
	s.x[24] = uint64((s.x[12] ^ s.x[7]))
	if runtime != nil && runtime.stopPC == 0x100234f28 {
		return 0x100234f28, true, nil
	}
	// 0x100234f28: str x24, [sp, #0x470]
	if err := s.storeValue(addOffset(s.sp, 1136), s.x[24], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234f2c {
		return 0x100234f2c, true, nil
	}
	// 0x100234f2c: mov w12, #0xea67
	s.x[12] = uint64(uint32(uint64(0xea67)))
	if runtime != nil && runtime.stopPC == 0x100234f30 {
		return 0x100234f30, true, nil
	}
	// 0x100234f30: movk w12, #0x5488, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x5488, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234f34 {
		return 0x100234f34, true, nil
	}
	// 0x100234f34: eor w12, w2, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100234f38 {
		return 0x100234f38, true, nil
	}
	// 0x100234f38: mov w3, #0x1599
	s.x[3] = uint64(uint32(uint64(0x1599)))
	if runtime != nil && runtime.stopPC == 0x100234f3c {
		return 0x100234f3c, true, nil
	}
	// 0x100234f3c: movk w3, #0xab77, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0xab77, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234f40 {
		return 0x100234f40, true, nil
	}
	// 0x100234f40: add w12, w12, w3
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100234f44 {
		return 0x100234f44, true, nil
	}
	// 0x100234f44: mov w3, #0xd9b4
	s.x[3] = uint64(uint32(uint64(0xd9b4)))
	if runtime != nil && runtime.stopPC == 0x100234f48 {
		return 0x100234f48, true, nil
	}
	// 0x100234f48: movk w3, #0xca60, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0xca60, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234f4c {
		return 0x100234f4c, true, nil
	}
	// 0x100234f4c: eor w3, w2, w3
	s.x[3] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100234f50 {
		return 0x100234f50, true, nil
	}
	// 0x100234f50: mov w5, #0x264c
	s.x[5] = uint64(uint32(uint64(0x264c)))
	if runtime != nil && runtime.stopPC == 0x100234f54 {
		return 0x100234f54, true, nil
	}
	// 0x100234f54: movk w5, #0x359f, lsl #16
	s.x[5] = uint64(uint32(insertBits(uint64(uint32(s.x[5])), 0x359f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234f58 {
		return 0x100234f58, true, nil
	}
	// 0x100234f58: add w3, w3, w5
	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100234f5c {
		return 0x100234f5c, true, nil
	}
	// 0x100234f5c: eor w12, w12, w3
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100234f60 {
		return 0x100234f60, true, nil
	}
	// 0x100234f60: mov w3, #0x39b4
	s.x[3] = uint64(uint32(uint64(0x39b4)))
	if runtime != nil && runtime.stopPC == 0x100234f64 {
		return 0x100234f64, true, nil
	}
	// 0x100234f64: movk w3, #0x6431, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0x6431, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234f68 {
		return 0x100234f68, true, nil
	}
	// 0x100234f68: eor w3, w2, w3
	s.x[3] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100234f6c {
		return 0x100234f6c, true, nil
	}
	// 0x100234f6c: mov w5, #0xc64c
	s.x[5] = uint64(uint32(uint64(0xc64c)))
	if runtime != nil && runtime.stopPC == 0x100234f70 {
		return 0x100234f70, true, nil
	}
	// 0x100234f70: movk w5, #0x9bce, lsl #16
	s.x[5] = uint64(uint32(insertBits(uint64(uint32(s.x[5])), 0x9bce, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234f74 {
		return 0x100234f74, true, nil
	}
	// 0x100234f74: add w3, w3, w5
	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100234f78 {
		return 0x100234f78, true, nil
	}
	// 0x100234f78: eor w12, w12, w3
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100234f7c {
		return 0x100234f7c, true, nil
	}
	// 0x100234f7c: ldrb w3, [x28, x24]
	value168, err := s.load((s.x[28] + s.x[24]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value168))
	if runtime != nil && runtime.stopPC == 0x100234f80 {
		return 0x100234f80, true, nil
	}
	// 0x100234f80: mov w5, #0x52
	s.x[5] = uint64(uint32(uint64(0x52)))
	if runtime != nil && runtime.stopPC == 0x100234f84 {
		return 0x100234f84, true, nil
	}
	// 0x100234f84: eor x3, x3, x5
	s.x[3] = uint64((s.x[3] ^ s.x[5]))
	if runtime != nil && runtime.stopPC == 0x100234f88 {
		return 0x100234f88, true, nil
	}
	// 0x100234f88: ldr w3, [x20, x3, lsl #2]
	value171, err := s.load((s.x[20] + shiftLeft(s.x[3], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value171))
	if runtime != nil && runtime.stopPC == 0x100234f8c {
		return 0x100234f8c, true, nil
	}
	// 0x100234f8c: ror w3, w3, #4
	s.x[3] = uint64(uint32(rotateRight(uint64(uint32(s.x[3])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x100234f90 {
		return 0x100234f90, true, nil
	}
	// 0x100234f90: eor w15, w15, w17
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100234f94 {
		return 0x100234f94, true, nil
	}
	// 0x100234f94: eor w17, w26, w3, lsl #4
	s.x[17] = uint64(uint32((uint64(uint32(s.x[26])) ^ shiftLeft(uint64(uint32(s.x[3])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100234f98 {
		return 0x100234f98, true, nil
	}
	// 0x100234f98: mov w5, #0xfec0
	s.x[5] = uint64(uint32(uint64(0xfec0)))
	if runtime != nil && runtime.stopPC == 0x100234f9c {
		return 0x100234f9c, true, nil
	}
	// 0x100234f9c: movk w5, #0x315e, lsl #16
	s.x[5] = uint64(uint32(insertBits(uint64(uint32(s.x[5])), 0x315e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234fa0 {
		return 0x100234fa0, true, nil
	}
	// 0x100234fa0: add w9, w9, w5
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100234fa4 {
		return 0x100234fa4, true, nil
	}
	// 0x100234fa4: eor w9, w15, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234fa8 {
		return 0x100234fa8, true, nil
	}
	// 0x100234fa8: mov w15, w4
	s.x[15] = uint64(uint32(uint64(uint32(s.x[4]))))
	if runtime != nil && runtime.stopPC == 0x100234fac {
		return 0x100234fac, true, nil
	}
	// 0x100234fac: eor x15, x15, x1
	s.x[15] = uint64((s.x[15] ^ s.x[1]))
	if runtime != nil && runtime.stopPC == 0x100234fb0 {
		return 0x100234fb0, true, nil
	}
	// 0x100234fb0: str x15, [sp, #0x440]
	if err := s.storeValue(addOffset(s.sp, 1088), s.x[15], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100234fb4 {
		return 0x100234fb4, true, nil
	}
	// 0x100234fb4: eor w9, w9, w3
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100234fb8 {
		return 0x100234fb8, true, nil
	}
	// 0x100234fb8: ldrb w15, [x28, x15]
	value183, err := s.load((s.x[28] + s.x[15]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value183))
	if runtime != nil && runtime.stopPC == 0x100234fbc {
		return 0x100234fbc, true, nil
	}
	// 0x100234fbc: mov w24, #0x85df
	s.x[24] = uint64(uint32(uint64(0x85df)))
	if runtime != nil && runtime.stopPC == 0x100234fc0 {
		return 0x100234fc0, true, nil
	}
	// 0x100234fc0: movk w24, #0xdac1, lsl #16
	s.x[24] = uint64(uint32(insertBits(uint64(uint32(s.x[24])), 0xdac1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234fc4 {
		return 0x100234fc4, true, nil
	}
	// 0x100234fc4: add w10, w10, w24
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[24])))))
	if runtime != nil && runtime.stopPC == 0x100234fc8 {
		return 0x100234fc8, true, nil
	}
	// 0x100234fc8: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234fcc {
		return 0x100234fcc, true, nil
	}
	// 0x100234fcc: mov w1, #0xf8d0
	s.x[1] = uint64(uint32(uint64(0xf8d0)))
	if runtime != nil && runtime.stopPC == 0x100234fd0 {
		return 0x100234fd0, true, nil
	}
	// 0x100234fd0: movk w1, #0x87da, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x87da, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234fd4 {
		return 0x100234fd4, true, nil
	}
	// 0x100234fd4: add w10, w17, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100234fd8 {
		return 0x100234fd8, true, nil
	}
	// 0x100234fd8: eor w25, w9, w10
	s.x[25] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234fdc {
		return 0x100234fdc, true, nil
	}
	// 0x100234fdc: eor x9, x15, #0xc
	s.x[9] = uint64((s.x[15] ^ uint64(0xc)))
	if runtime != nil && runtime.stopPC == 0x100234fe0 {
		return 0x100234fe0, true, nil
	}
	// 0x100234fe0: ldr w15, [x6, x9, lsl #2]
	value193, err := s.load((s.x[6] + shiftLeft(s.x[9], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value193))
	if runtime != nil && runtime.stopPC == 0x100234fe4 {
		return 0x100234fe4, true, nil
	}
	// 0x100234fe4: mov w9, #0xc4e0
	s.x[9] = uint64(uint32(uint64(0xc4e0)))
	if runtime != nil && runtime.stopPC == 0x100234fe8 {
		return 0x100234fe8, true, nil
	}
	// 0x100234fe8: movk w9, #0xb22b, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xb22b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234fec {
		return 0x100234fec, true, nil
	}
	// 0x100234fec: eor w9, w15, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100234ff0 {
		return 0x100234ff0, true, nil
	}
	// 0x100234ff0: mov w10, #0x3b20
	s.x[10] = uint64(uint32(uint64(0x3b20)))
	if runtime != nil && runtime.stopPC == 0x100234ff4 {
		return 0x100234ff4, true, nil
	}
	// 0x100234ff4: movk w10, #0x4dd4, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x4dd4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100234ff8 {
		return 0x100234ff8, true, nil
	}
	// 0x100234ff8: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100234ffc {
		return 0x100234ffc, true, nil
	}
	// 0x100234ffc: mov w10, #0x99d3
	s.x[10] = uint64(uint32(uint64(0x99d3)))
	if runtime != nil && runtime.stopPC == 0x100235000 {
		return 0x100235000, true, nil
	}
	// 0x100235000: movk w10, #0x658c, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x658c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235004 {
		return 0x100235004, true, nil
	}
	// 0x100235004: eor w10, w15, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235008 {
		return 0x100235008, true, nil
	}
	// 0x100235008: mov w17, #0x662d
	s.x[17] = uint64(uint32(uint64(0x662d)))
	if runtime != nil && runtime.stopPC == 0x10023500c {
		return 0x10023500c, true, nil
	}
	// 0x10023500c: movk w17, #0x9a73, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x9a73, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235010 {
		return 0x100235010, true, nil
	}
	// 0x100235010: add w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235014 {
		return 0x100235014, true, nil
	}
	// 0x100235014: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235018 {
		return 0x100235018, true, nil
	}
	// 0x100235018: mov w10, #0x5754
	s.x[10] = uint64(uint32(uint64(0x5754)))
	if runtime != nil && runtime.stopPC == 0x10023501c {
		return 0x10023501c, true, nil
	}
	// 0x10023501c: movk w10, #0x2d7e, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x2d7e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235020 {
		return 0x100235020, true, nil
	}
	// 0x100235020: eor w10, w15, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235024 {
		return 0x100235024, true, nil
	}
	// 0x100235024: mov w17, #0xa8ac
	s.x[17] = uint64(uint32(uint64(0xa8ac)))
	if runtime != nil && runtime.stopPC == 0x100235028 {
		return 0x100235028, true, nil
	}
	// 0x100235028: movk w17, #0xd281, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xd281, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023502c {
		return 0x10023502c, true, nil
	}
	// 0x10023502c: add w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235030 {
		return 0x100235030, true, nil
	}
	// 0x100235030: eor w3, w9, w10
	s.x[3] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235034 {
		return 0x100235034, true, nil
	}
	// 0x100235034: mov w9, #-0xa
	s.x[9] = uint64(uint32(uint64(0xfffffff6)))
	if runtime != nil && runtime.stopPC == 0x100235038 {
		return 0x100235038, true, nil
	}
	// 0x100235038: and w9, w4, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023503c {
		return 0x10023503c, true, nil
	}
	// 0x10023503c: eor x10, x9, x19
	s.x[10] = uint64((s.x[9] ^ s.x[19]))
	if runtime != nil && runtime.stopPC == 0x100235040 {
		return 0x100235040, true, nil
	}
	// 0x100235040: str x10, [sp, #0x468]
	if err := s.storeValue(addOffset(s.sp, 1128), s.x[10], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235044 {
		return 0x100235044, true, nil
	}
	// 0x100235044: eor w9, w11, w16
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235048 {
		return 0x100235048, true, nil
	}
	// 0x100235048: ldrb w10, [x28, x10]
	value219, err := s.load((s.x[28] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value219))
	if runtime != nil && runtime.stopPC == 0x10023504c {
		return 0x10023504c, true, nil
	}
	// 0x10023504c: add w11, w0, w5
	s.x[11] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100235050 {
		return 0x100235050, true, nil
	}
	// 0x100235050: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235054 {
		return 0x100235054, true, nil
	}
	// 0x100235054: and w11, w4, #0xfffffff8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffff8))))
	if runtime != nil && runtime.stopPC == 0x100235058 {
		return 0x100235058, true, nil
	}
	// 0x100235058: eor x16, x11, x27
	s.x[16] = uint64((s.x[11] ^ s.x[27]))
	if runtime != nil && runtime.stopPC == 0x10023505c {
		return 0x10023505c, true, nil
	}
	// 0x10023505c: str x16, [sp, #0x428]
	if err := s.storeValue(addOffset(s.sp, 1064), s.x[16], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235060 {
		return 0x100235060, true, nil
	}
	// 0x100235060: mov w11, #0x3b
	s.x[11] = uint64(uint32(uint64(0x3b)))
	if runtime != nil && runtime.stopPC == 0x100235064 {
		return 0x100235064, true, nil
	}
	// 0x100235064: eor x10, x10, x11
	s.x[10] = uint64((s.x[10] ^ s.x[11]))
	if runtime != nil && runtime.stopPC == 0x100235068 {
		return 0x100235068, true, nil
	}
	// 0x100235068: ldr w10, [x22, x10, lsl #2]
	value227, err := s.load((s.x[22] + shiftLeft(s.x[10], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value227))
	if runtime != nil && runtime.stopPC == 0x10023506c {
		return 0x10023506c, true, nil
	}
	// 0x10023506c: add w8, w8, w1
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235070 {
		return 0x100235070, true, nil
	}
	// 0x100235070: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235074 {
		return 0x100235074, true, nil
	}
	// 0x100235074: ldrb w9, [x28, x16]
	value230, err := s.load((s.x[28] + s.x[16]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value230))
	if runtime != nil && runtime.stopPC == 0x100235078 {
		return 0x100235078, true, nil
	}
	// 0x100235078: add w11, w14, w24
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[24])))))
	if runtime != nil && runtime.stopPC == 0x10023507c {
		return 0x10023507c, true, nil
	}
	// 0x10023507c: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235080 {
		return 0x100235080, true, nil
	}
	// 0x100235080: mov w11, #0x79
	s.x[11] = uint64(uint32(uint64(0x79)))
	if runtime != nil && runtime.stopPC == 0x100235084 {
		return 0x100235084, true, nil
	}
	// 0x100235084: eor w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235088 {
		return 0x100235088, true, nil
	}
	// 0x100235088: eor w0, w8, w10
	s.x[0] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023508c {
		return 0x10023508c, true, nil
	}
	// 0x10023508c: eor x8, x9, #0x1c
	s.x[8] = uint64((s.x[9] ^ uint64(0x1c)))
	if runtime != nil && runtime.stopPC == 0x100235090 {
		return 0x100235090, true, nil
	}
	// 0x100235090: ldrb w8, [x21, x8]
	value237, err := s.load((s.x[21] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value237))
	if runtime != nil && runtime.stopPC == 0x100235094 {
		return 0x100235094, true, nil
	}
	// 0x100235094: mov w10, #0x91
	s.x[10] = uint64(uint32(uint64(0x91)))
	if runtime != nil && runtime.stopPC == 0x100235098 {
		return 0x100235098, true, nil
	}
	// 0x100235098: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023509c {
		return 0x10023509c, true, nil
	}
	// 0x10023509c: add w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002350a0 {
		return 0x1002350a0, true, nil
	}
	// 0x1002350a0: ldrb w10, [x28, x30]
	value241, err := s.load((s.x[28] + s.x[30]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value241))
	if runtime != nil && runtime.stopPC == 0x1002350a4 {
		return 0x1002350a4, true, nil
	}
	// 0x1002350a4: mov w11, #0x6c
	s.x[11] = uint64(uint32(uint64(0x6c)))
	if runtime != nil && runtime.stopPC == 0x1002350a8 {
		return 0x1002350a8, true, nil
	}
	// 0x1002350a8: eor x10, x10, x11
	s.x[10] = uint64((s.x[10] ^ s.x[11]))
	if runtime != nil && runtime.stopPC == 0x1002350ac {
		return 0x1002350ac, true, nil
	}
	// 0x1002350ac: ldr w10, [x22, x10, lsl #2]
	value244, err := s.load((s.x[22] + shiftLeft(s.x[10], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value244))
	if runtime != nil && runtime.stopPC == 0x1002350b0 {
		return 0x1002350b0, true, nil
	}
	// 0x1002350b0: and w11, w4, #0xfffffff9
	s.x[11] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffff9))))
	if runtime != nil && runtime.stopPC == 0x1002350b4 {
		return 0x1002350b4, true, nil
	}
	// 0x1002350b4: eor x11, x11, x7
	s.x[11] = uint64((s.x[11] ^ s.x[7]))
	if runtime != nil && runtime.stopPC == 0x1002350b8 {
		return 0x1002350b8, true, nil
	}
	// 0x1002350b8: str x11, [sp, #0x458]
	if err := s.storeValue(addOffset(s.sp, 1112), s.x[11], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002350bc {
		return 0x1002350bc, true, nil
	}
	// 0x1002350bc: ldrb w11, [x28, x11]
	value248, err := s.load((s.x[28] + s.x[11]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value248))
	if runtime != nil && runtime.stopPC == 0x1002350c0 {
		return 0x1002350c0, true, nil
	}
	// 0x1002350c0: mov w14, #0x6e
	s.x[14] = uint64(uint32(uint64(0x6e)))
	if runtime != nil && runtime.stopPC == 0x1002350c4 {
		return 0x1002350c4, true, nil
	}
	// 0x1002350c4: eor x11, x11, x14
	s.x[11] = uint64((s.x[11] ^ s.x[14]))
	if runtime != nil && runtime.stopPC == 0x1002350c8 {
		return 0x1002350c8, true, nil
	}
	// 0x1002350c8: ldr w11, [x20, x11, lsl #2]
	value251, err := s.load((s.x[20] + shiftLeft(s.x[11], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value251))
	if runtime != nil && runtime.stopPC == 0x1002350cc {
		return 0x1002350cc, true, nil
	}
	// 0x1002350cc: add w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002350d0 {
		return 0x1002350d0, true, nil
	}
	// 0x1002350d0: ror w9, w11, #4
	s.x[9] = uint64(uint32(rotateRight(uint64(uint32(s.x[11])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002350d4 {
		return 0x1002350d4, true, nil
	}
	// 0x1002350d4: eor w11, w23, w2, lsl #2
	s.x[11] = uint64(uint32((uint64(uint32(s.x[23])) ^ shiftLeft(uint64(uint32(s.x[2])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x1002350d8 {
		return 0x1002350d8, true, nil
	}
	// 0x1002350d8: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002350dc {
		return 0x1002350dc, true, nil
	}
	// 0x1002350dc: add w11, w11, w5
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x1002350e0 {
		return 0x1002350e0, true, nil
	}
	// 0x1002350e0: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002350e4 {
		return 0x1002350e4, true, nil
	}
	// 0x1002350e4: eor w9, w26, w9, lsl #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[26])) ^ shiftLeft(uint64(uint32(s.x[9])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002350e8 {
		return 0x1002350e8, true, nil
	}
	// 0x1002350e8: add w9, w9, w1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002350ec {
		return 0x1002350ec, true, nil
	}
	// 0x1002350ec: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002350f0 {
		return 0x1002350f0, true, nil
	}
	// 0x1002350f0: str w4, [sp, #0x464]
	if err := s.storeValue(addOffset(s.sp, 1124), uint64(uint32(s.x[4])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002350f4 {
		return 0x1002350f4, true, nil
	}
	// 0x1002350f4: and w10, w4, #0xfffffffc
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffffc))))
	if runtime != nil && runtime.stopPC == 0x1002350f8 {
		return 0x1002350f8, true, nil
	}
	// 0x1002350f8: eor x11, x10, x27
	s.x[11] = uint64((s.x[10] ^ s.x[27]))
	if runtime != nil && runtime.stopPC == 0x1002350fc {
		return 0x1002350fc, true, nil
	}
	// 0x1002350fc: str x11, [sp, #0x430]
	if err := s.storeValue(addOffset(s.sp, 1072), s.x[11], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235100 {
		return 0x100235100, true, nil
	}
	// 0x100235100: sub w8, w8, #9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(0x9))))
	if runtime != nil && runtime.stopPC == 0x100235104 {
		return 0x100235104, true, nil
	}
	// 0x100235104: and w8, w8, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235108 {
		return 0x100235108, true, nil
	}
	// 0x100235108: add w10, w12, w24
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[24])))))
	if runtime != nil && runtime.stopPC == 0x10023510c {
		return 0x10023510c, true, nil
	}
	// 0x10023510c: eor w8, w10, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235110 {
		return 0x100235110, true, nil
	}
	// 0x100235110: ldrb w10, [x28, x11]
	value269, err := s.load((s.x[28] + s.x[11]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value269))
	if runtime != nil && runtime.stopPC == 0x100235114 {
		return 0x100235114, true, nil
	}
	// 0x100235114: eor w30, w9, w8
	s.x[30] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235118 {
		return 0x100235118, true, nil
	}
	// 0x100235118: mov w7, #0xb7
	s.x[7] = uint64(uint32(uint64(0xb7)))
	if runtime != nil && runtime.stopPC == 0x10023511c {
		return 0x10023511c, true, nil
	}
	// 0x10023511c: eor x8, x10, x7
	s.x[8] = uint64((s.x[10] ^ s.x[7]))
	if runtime != nil && runtime.stopPC == 0x100235120 {
		return 0x100235120, true, nil
	}
	// 0x100235120: ldrb w8, [x21, x8]
	value273, err := s.load((s.x[21] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value273))
	if runtime != nil && runtime.stopPC == 0x100235124 {
		return 0x100235124, true, nil
	}
	// 0x100235124: mov w9, #0xd2
	s.x[9] = uint64(uint32(uint64(0xd2)))
	if runtime != nil && runtime.stopPC == 0x100235128 {
		return 0x100235128, true, nil
	}
	// 0x100235128: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023512c {
		return 0x10023512c, true, nil
	}
	// 0x10023512c: mov w11, #0x3a
	s.x[11] = uint64(uint32(uint64(0x3a)))
	if runtime != nil && runtime.stopPC == 0x100235130 {
		return 0x100235130, true, nil
	}
	// 0x100235130: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235134 {
		return 0x100235134, true, nil
	}
	// 0x100235134: add w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235138 {
		return 0x100235138, true, nil
	}
	// 0x100235138: add w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x10023513c {
		return 0x10023513c, true, nil
	}
	// 0x10023513c: and w9, w4, #0xfffffffe
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & uint64(0xfffffffe))))
	if runtime != nil && runtime.stopPC == 0x100235140 {
		return 0x100235140, true, nil
	}
	// 0x100235140: eor x9, x9, x19
	s.x[9] = uint64((s.x[9] ^ s.x[19]))
	if runtime != nil && runtime.stopPC == 0x100235144 {
		return 0x100235144, true, nil
	}
	// 0x100235144: str x9, [sp, #0x450]
	if err := s.storeValue(addOffset(s.sp, 1104), s.x[9], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235148 {
		return 0x100235148, true, nil
	}
	// 0x100235148: ldrb w9, [x28, x9]
	value283, err := s.load((s.x[28] + s.x[9]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value283))
	if runtime != nil && runtime.stopPC == 0x10023514c {
		return 0x10023514c, true, nil
	}
	// 0x10023514c: mov w10, #0x62
	s.x[10] = uint64(uint32(uint64(0x62)))
	if runtime != nil && runtime.stopPC == 0x100235150 {
		return 0x100235150, true, nil
	}
	// 0x100235150: eor x9, x9, x10
	s.x[9] = uint64((s.x[9] ^ s.x[10]))
	if runtime != nil && runtime.stopPC == 0x100235154 {
		return 0x100235154, true, nil
	}
	// 0x100235154: ldr w9, [x22, x9, lsl #2]
	value286, err := s.load((s.x[22] + shiftLeft(s.x[9], 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value286))
	if runtime != nil && runtime.stopPC == 0x100235158 {
		return 0x100235158, true, nil
	}
	// 0x100235158: eor w10, w23, w15, lsl #2
	s.x[10] = uint64(uint32((uint64(uint32(s.x[23])) ^ shiftLeft(uint64(uint32(s.x[15])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x10023515c {
		return 0x10023515c, true, nil
	}
	// 0x10023515c: eor w9, w13, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235160 {
		return 0x100235160, true, nil
	}
	// 0x100235160: add w10, w10, w5
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100235164 {
		return 0x100235164, true, nil
	}
	// 0x100235164: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235168 {
		return 0x100235168, true, nil
	}
	// 0x100235168: eor w10, w26, w13, lsl #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[26])) ^ shiftLeft(uint64(uint32(s.x[13])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x10023516c {
		return 0x10023516c, true, nil
	}
	// 0x10023516c: sub w8, w8, #9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(0x9))))
	if runtime != nil && runtime.stopPC == 0x100235170 {
		return 0x100235170, true, nil
	}
	// 0x100235170: and w8, w8, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235174 {
		return 0x100235174, true, nil
	}
	// 0x100235174: add w11, w3, w24
	s.x[11] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[24])))))
	if runtime != nil && runtime.stopPC == 0x100235178 {
		return 0x100235178, true, nil
	}
	// 0x100235178: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023517c {
		return 0x10023517c, true, nil
	}
	// 0x10023517c: add w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235180 {
		return 0x100235180, true, nil
	}
	// 0x100235180: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235184 {
		return 0x100235184, true, nil
	}
	// 0x100235184: eor w3, w9, w8
	s.x[3] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235188 {
		return 0x100235188, true, nil
	}
	// 0x100235188: adrp x8, #0x101075000
	s.x[8] = uint64(0x101075000)
	if runtime != nil && runtime.stopPC == 0x10023518c {
		return 0x10023518c, true, nil
	}
	// 0x10023518c: ldr w21, [x8, #0x588]
	value300, err := s.load(addOffset(s.x[8], 1416), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[21] = uint64(uint32(value300))
	if runtime != nil && runtime.stopPC == 0x100235190 {
		return 0x100235190, true, nil
	}
	// 0x100235190: adrp x8, #0x101075000
	s.x[8] = uint64(0x101075000)
	if runtime != nil && runtime.stopPC == 0x100235194 {
		return 0x100235194, true, nil
	}
	// 0x100235194: ldr w10, [x8, #0xb40]
	value302, err := s.load(addOffset(s.x[8], 2880), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value302))
	if runtime != nil && runtime.stopPC == 0x100235198 {
		return 0x100235198, true, nil
	}
	// 0x100235198: str w10, [sp, #0x424]
	if err := s.storeValue(addOffset(s.sp, 1060), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x10023519c {
		return 0x10023519c, true, nil
	}
	// 0x10023519c: mov w8, #0xd1a6
	s.x[8] = uint64(uint32(uint64(0xd1a6)))
	if runtime != nil && runtime.stopPC == 0x1002351a0 {
		return 0x1002351a0, true, nil
	}
	// 0x1002351a0: movk w8, #0x8eae, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x8eae, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002351a4 {
		return 0x1002351a4, true, nil
	}
	// 0x1002351a4: eor w9, w21, w8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002351a8 {
		return 0x1002351a8, true, nil
	}
	// 0x1002351a8: eor w8, w9, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002351ac {
		return 0x1002351ac, true, nil
	}
	// 0x1002351ac: mov w10, #0xd4cf
	s.x[10] = uint64(uint32(uint64(0xd4cf)))
	if runtime != nil && runtime.stopPC == 0x1002351b0 {
		return 0x1002351b0, true, nil
	}
	// 0x1002351b0: movk w10, #0xa5c, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa5c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002351b4 {
		return 0x1002351b4, true, nil
	}
	// 0x1002351b4: mul w20, w8, w10
	s.x[20] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002351b8 {
		return 0x1002351b8, true, nil
	}
	// 0x1002351b8: mov w1, #0xd4cf
	s.x[1] = uint64(uint32(uint64(0xd4cf)))
	if runtime != nil && runtime.stopPC == 0x1002351bc {
		return 0x1002351bc, true, nil
	}
	// 0x1002351bc: movk w1, #0xa5c, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa5c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002351c0 {
		return 0x1002351c0, true, nil
	}
	// 0x1002351c0: and x2, x20, #0xff
	s.x[2] = uint64((s.x[20] & uint64(0xff)))
	if runtime != nil && runtime.stopPC == 0x1002351c4 {
		return 0x1002351c4, true, nil
	}
	// 0x1002351c4: ldur x8, [x29, #-0xb8]
	value314, err := s.load(addOffset(s.x[29], -184), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(value314)
	if runtime != nil && runtime.stopPC == 0x1002351c8 {
		return 0x1002351c8, true, nil
	}
	// 0x1002351c8: add x8, x8, x2
	s.x[8] = uint64((s.x[8] + s.x[2]))
	if runtime != nil && runtime.stopPC == 0x1002351cc {
		return 0x1002351cc, true, nil
	}
	// 0x1002351cc: ldrb w8, [x8]
	value316, err := s.load(s.x[8], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value316))
	if runtime != nil && runtime.stopPC == 0x1002351d0 {
		return 0x1002351d0, true, nil
	}
	// 0x1002351d0: mov w10, #0x5d
	s.x[10] = uint64(uint32(uint64(0x5d)))
	if runtime != nil && runtime.stopPC == 0x1002351d4 {
		return 0x1002351d4, true, nil
	}
	// 0x1002351d4: eor w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002351d8 {
		return 0x1002351d8, true, nil
	}
	// 0x1002351d8: ldur x10, [x29, #-0xa8]
	value319, err := s.load(addOffset(s.x[29], -168), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value319)
	if runtime != nil && runtime.stopPC == 0x1002351dc {
		return 0x1002351dc, true, nil
	}
	// 0x1002351dc: ldrb w8, [x10, w8, uxtw]
	value320, err := s.load((s.x[10] + uint64(uint32(uint64(uint32(s.x[8]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value320))
	if runtime != nil && runtime.stopPC == 0x1002351e0 {
		return 0x1002351e0, true, nil
	}
	// 0x1002351e0: eor w8, w20, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002351e4 {
		return 0x1002351e4, true, nil
	}
	// 0x1002351e4: mov w10, #0x152c
	s.x[10] = uint64(uint32(uint64(0x152c)))
	if runtime != nil && runtime.stopPC == 0x1002351e8 {
		return 0x1002351e8, true, nil
	}
	// 0x1002351e8: add w8, w10, w8, uxtb
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint8(uint64(uint32(s.x[8])))))))
	if runtime != nil && runtime.stopPC == 0x1002351ec {
		return 0x1002351ec, true, nil
	}
	// 0x1002351ec: adrp x17, #0x10105b000
	s.x[17] = uint64(0x10105b000)
	if runtime != nil && runtime.stopPC == 0x1002351f0 {
		return 0x1002351f0, true, nil
	}
	// 0x1002351f0: add x17, x17, #0x6f0
	s.x[17] = uint64((s.x[17] + uint64(0x6f0)))
	if runtime != nil && runtime.stopPC == 0x1002351f4 {
		return 0x1002351f4, true, nil
	}
	// 0x1002351f4: ldr x8, [x17, w8, uxtw #3]
	value326, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(value326)
	if runtime != nil && runtime.stopPC == 0x1002351f8 {
		return 0x1002351f8, true, nil
	}
	// 0x1002351f8: mov x10, #-0x433e
	s.x[10] = uint64(uint64(0xffffffffffffbcc2))
	if runtime != nil && runtime.stopPC == 0x1002351fc {
		return 0x1002351fc, true, nil
	}
	// 0x1002351fc: movk x10, #0xa28e, lsl #16
	s.x[10] = uint64(insertBits(s.x[10], 0xa28e, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x100235200 {
		return 0x100235200, true, nil
	}
	// 0x100235200: add x22, x8, x10
	s.x[22] = uint64((s.x[8] + s.x[10]))
	if runtime != nil && runtime.stopPC == 0x100235204 {
		return 0x100235204, true, nil
	}
	// 0x100235204: ldr w8, [sp, #0x2c8]
	value330, err := s.load(addOffset(s.sp, 712), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value330))
	if runtime != nil && runtime.stopPC == 0x100235208 {
		return 0x100235208, true, nil
	}
	// 0x100235208: eor w16, w8, w3
	s.x[16] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x10023520c {
		return 0x10023520c, true, nil
	}
	// 0x10023520c: lsr w8, w16, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[16])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235210 {
		return 0x100235210, true, nil
	}
	// 0x100235210: ldr w8, [x22, w8, uxtw #2]
	value333, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value333))
	if runtime != nil && runtime.stopPC == 0x100235214 {
		return 0x100235214, true, nil
	}
	// 0x100235214: ldur x10, [x29, #-0xa0]
	value334, err := s.load(addOffset(s.x[29], -160), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value334)
	if runtime != nil && runtime.stopPC == 0x100235218 {
		return 0x100235218, true, nil
	}
	// 0x100235218: ldr w10, [x10]
	value335, err := s.load(s.x[10], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value335))
	if runtime != nil && runtime.stopPC == 0x10023521c {
		return 0x10023521c, true, nil
	}
	// 0x10023521c: eor w12, w9, w10
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235220 {
		return 0x100235220, true, nil
	}
	// 0x100235220: ldr w9, [sp, #0x36c]
	value337, err := s.load(addOffset(s.sp, 876), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value337))
	if runtime != nil && runtime.stopPC == 0x100235224 {
		return 0x100235224, true, nil
	}
	// 0x100235224: str w0, [sp, #0x3fc]
	if err := s.storeValue(addOffset(s.sp, 1020), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235228 {
		return 0x100235228, true, nil
	}
	// 0x100235228: eor w19, w9, w0
	s.x[19] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x10023522c {
		return 0x10023522c, true, nil
	}
	// 0x10023522c: ldr w9, [sp, #0x370]
	value340, err := s.load(addOffset(s.sp, 880), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value340))
	if runtime != nil && runtime.stopPC == 0x100235230 {
		return 0x100235230, true, nil
	}
	// 0x100235230: eor w15, w9, w25
	s.x[15] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[25])))))
	if runtime != nil && runtime.stopPC == 0x100235234 {
		return 0x100235234, true, nil
	}
	// 0x100235234: ubfx w13, w15, #0x10, #8
	s.x[13] = uint64(uint32((uint64(uint32(s.x[15])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235238 {
		return 0x100235238, true, nil
	}
	// 0x100235238: mov w27, #0x7a55
	s.x[27] = uint64(uint32(uint64(0x7a55)))
	if runtime != nil && runtime.stopPC == 0x10023523c {
		return 0x10023523c, true, nil
	}
	// 0x10023523c: movk w27, #0xbc9e, lsl #16
	s.x[27] = uint64(uint32(insertBits(uint64(uint32(s.x[27])), 0xbc9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235240 {
		return 0x100235240, true, nil
	}
	// 0x100235240: eor w9, w13, w27
	s.x[9] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235244 {
		return 0x100235244, true, nil
	}
	// 0x100235244: ubfx w11, w19, #8, #8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[19])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235248 {
		return 0x100235248, true, nil
	}
	// 0x100235248: ldr w10, [sp, #0x360]
	value347, err := s.load(addOffset(s.sp, 864), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value347))
	if runtime != nil && runtime.stopPC == 0x10023524c {
		return 0x10023524c, true, nil
	}
	// 0x10023524c: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235250 {
		return 0x100235250, true, nil
	}
	// 0x100235250: sub w9, w13, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[13])) - uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235254 {
		return 0x100235254, true, nil
	}
	// 0x100235254: mov w0, #0x95
	s.x[0] = uint64(uint32(uint64(0x95)))
	if runtime != nil && runtime.stopPC == 0x100235258 {
		return 0x100235258, true, nil
	}
	// 0x100235258: eor w14, w13, w0
	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x10023525c {
		return 0x10023525c, true, nil
	}
	// 0x10023525c: mul w23, w12, w1
	s.x[23] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235260 {
		return 0x100235260, true, nil
	}
	// 0x100235260: and x12, x23, #0xff
	s.x[12] = uint64((s.x[23] & uint64(0xff)))
	if runtime != nil && runtime.stopPC == 0x100235264 {
		return 0x100235264, true, nil
	}
	// 0x100235264: str x12, [sp, #0x410]
	if err := s.storeValue(addOffset(s.sp, 1040), s.x[12], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235268 {
		return 0x100235268, true, nil
	}
	// 0x100235268: ldp x4, x13, [x29, #-0xd0]
	value355a, err := s.load(addOffset(s.x[29], -208), 8)
	if err != nil {
		return 0, false, err
	}
	value355b, err := s.load((addOffset(s.x[29], -208))+8, 8)
	if err != nil {
		return 0, false, err
	}
	s.x[4] = uint64(value355a)
	s.x[13] = uint64(value355b)
	if runtime != nil && runtime.stopPC == 0x10023526c {
		return 0x10023526c, true, nil
	}
	// 0x10023526c: add x12, x13, x12
	s.x[12] = uint64((s.x[13] + s.x[12]))
	if runtime != nil && runtime.stopPC == 0x100235270 {
		return 0x100235270, true, nil
	}
	// 0x100235270: ldrb w12, [x12]
	value357, err := s.load(s.x[12], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value357))
	if runtime != nil && runtime.stopPC == 0x100235274 {
		return 0x100235274, true, nil
	}
	// 0x100235274: mov w6, #0x17
	s.x[6] = uint64(uint32(uint64(0x17)))
	if runtime != nil && runtime.stopPC == 0x100235278 {
		return 0x100235278, true, nil
	}
	// 0x100235278: eor w12, w12, w6
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x10023527c {
		return 0x10023527c, true, nil
	}
	// 0x10023527c: ldrb w12, [x4, w12, uxtw]
	value360, err := s.load((s.x[4] + uint64(uint32(uint64(uint32(s.x[12]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value360))
	if runtime != nil && runtime.stopPC == 0x100235280 {
		return 0x100235280, true, nil
	}
	// 0x100235280: eor w12, w23, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235284 {
		return 0x100235284, true, nil
	}
	// 0x100235284: and w12, w12, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235288 {
		return 0x100235288, true, nil
	}
	// 0x100235288: add w12, w12, #0xa89
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0xa89))))
	if runtime != nil && runtime.stopPC == 0x10023528c {
		return 0x10023528c, true, nil
	}
	// 0x10023528c: ldr x12, [x17, w12, uxtw #3]
	value364, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(value364)
	if runtime != nil && runtime.stopPC == 0x100235290 {
		return 0x100235290, true, nil
	}
	// 0x100235290: mov x1, #-0x9f66
	s.x[1] = uint64(uint64(0xffffffffffff609a))
	if runtime != nil && runtime.stopPC == 0x100235294 {
		return 0x100235294, true, nil
	}
	// 0x100235294: movk x1, #0xaf03, lsl #16
	s.x[1] = uint64(insertBits(s.x[1], 0xaf03, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x100235298 {
		return 0x100235298, true, nil
	}
	// 0x100235298: add x24, x12, x1
	s.x[24] = uint64((s.x[12] + s.x[1]))
	if runtime != nil && runtime.stopPC == 0x10023529c {
		return 0x10023529c, true, nil
	}
	// 0x10023529c: ldr w12, [x24, w14, uxtw #2]
	value368, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value368))
	if runtime != nil && runtime.stopPC == 0x1002352a0 {
		return 0x1002352a0, true, nil
	}
	// 0x1002352a0: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002352a4 {
		return 0x1002352a4, true, nil
	}
	// 0x1002352a4: str x2, [sp, #0x418]
	if err := s.storeValue(addOffset(s.sp, 1048), s.x[2], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002352a8 {
		return 0x1002352a8, true, nil
	}
	// 0x1002352a8: add x10, x13, x2
	s.x[10] = uint64((s.x[13] + s.x[2]))
	if runtime != nil && runtime.stopPC == 0x1002352ac {
		return 0x1002352ac, true, nil
	}
	// 0x1002352ac: ldrb w10, [x10]
	value372, err := s.load(s.x[10], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value372))
	if runtime != nil && runtime.stopPC == 0x1002352b0 {
		return 0x1002352b0, true, nil
	}
	// 0x1002352b0: mov w13, #0xb1
	s.x[13] = uint64(uint32(uint64(0xb1)))
	if runtime != nil && runtime.stopPC == 0x1002352b4 {
		return 0x1002352b4, true, nil
	}
	// 0x1002352b4: eor w10, w10, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002352b8 {
		return 0x1002352b8, true, nil
	}
	// 0x1002352b8: ldrb w10, [x4, w10, uxtw]
	value375, err := s.load((s.x[4] + uint64(uint32(uint64(uint32(s.x[10]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value375))
	if runtime != nil && runtime.stopPC == 0x1002352bc {
		return 0x1002352bc, true, nil
	}
	// 0x1002352bc: eor w8, w9, w8, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002352c0 {
		return 0x1002352c0, true, nil
	}
	// 0x1002352c0: eor w9, w20, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002352c4 {
		return 0x1002352c4, true, nil
	}
	// 0x1002352c4: and w9, w9, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002352c8 {
		return 0x1002352c8, true, nil
	}
	// 0x1002352c8: add w9, w9, #0xd3c
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0xd3c))))
	if runtime != nil && runtime.stopPC == 0x1002352cc {
		return 0x1002352cc, true, nil
	}
	// 0x1002352cc: ldr x9, [x17, w9, uxtw #3]
	value380, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value380)
	if runtime != nil && runtime.stopPC == 0x1002352d0 {
		return 0x1002352d0, true, nil
	}
	// 0x1002352d0: mov w13, #0x33
	s.x[13] = uint64(uint32(uint64(0x33)))
	if runtime != nil && runtime.stopPC == 0x1002352d4 {
		return 0x1002352d4, true, nil
	}
	// 0x1002352d4: eor w10, w11, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002352d8 {
		return 0x1002352d8, true, nil
	}
	// 0x1002352d8: mov x11, #-0x535b
	s.x[11] = uint64(uint64(0xffffffffffffaca5))
	if runtime != nil && runtime.stopPC == 0x1002352dc {
		return 0x1002352dc, true, nil
	}
	// 0x1002352dc: movk x11, #0x8280, lsl #16
	s.x[11] = uint64(insertBits(s.x[11], 0x8280, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x1002352e0 {
		return 0x1002352e0, true, nil
	}
	// 0x1002352e0: add x26, x9, x11
	s.x[26] = uint64((s.x[9] + s.x[11]))
	if runtime != nil && runtime.stopPC == 0x1002352e4 {
		return 0x1002352e4, true, nil
	}
	// 0x1002352e4: ldr w9, [x26, w10, uxtw #2]
	value386, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value386))
	if runtime != nil && runtime.stopPC == 0x1002352e8 {
		return 0x1002352e8, true, nil
	}
	// 0x1002352e8: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002352ec {
		return 0x1002352ec, true, nil
	}
	// 0x1002352ec: ldr w10, [sp, #0x388]
	value388, err := s.load(addOffset(s.sp, 904), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value388))
	if runtime != nil && runtime.stopPC == 0x1002352f0 {
		return 0x1002352f0, true, nil
	}
	// 0x1002352f0: eor w10, w30, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[30])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002352f4 {
		return 0x1002352f4, true, nil
	}
	// 0x1002352f4: and w10, w10, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002352f8 {
		return 0x1002352f8, true, nil
	}
	// 0x1002352f8: ldur x11, [x29, #-0xe8]
	value391, err := s.load(addOffset(s.x[29], -232), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value391)
	if runtime != nil && runtime.stopPC == 0x1002352fc {
		return 0x1002352fc, true, nil
	}
	// 0x1002352fc: add x11, x11, x2
	s.x[11] = uint64((s.x[11] + s.x[2]))
	if runtime != nil && runtime.stopPC == 0x100235300 {
		return 0x100235300, true, nil
	}
	// 0x100235300: ldrb w11, [x11]
	value393, err := s.load(s.x[11], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value393))
	if runtime != nil && runtime.stopPC == 0x100235304 {
		return 0x100235304, true, nil
	}
	// 0x100235304: mov w12, #0x29
	s.x[12] = uint64(uint32(uint64(0x29)))
	if runtime != nil && runtime.stopPC == 0x100235308 {
		return 0x100235308, true, nil
	}
	// 0x100235308: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023530c {
		return 0x10023530c, true, nil
	}
	// 0x10023530c: ldur x12, [x29, #-0xf0]
	value396, err := s.load(addOffset(s.x[29], -240), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(value396)
	if runtime != nil && runtime.stopPC == 0x100235310 {
		return 0x100235310, true, nil
	}
	// 0x100235310: ldrb w11, [x12, w11, uxtw]
	value397, err := s.load((s.x[12] + uint64(uint32(uint64(uint32(s.x[11]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value397))
	if runtime != nil && runtime.stopPC == 0x100235314 {
		return 0x100235314, true, nil
	}
	// 0x100235314: eor w11, w20, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235318 {
		return 0x100235318, true, nil
	}
	// 0x100235318: and w11, w11, #0xff
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x10023531c {
		return 0x10023531c, true, nil
	}
	// 0x10023531c: add w11, w11, #0xbbe
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(0xbbe))))
	if runtime != nil && runtime.stopPC == 0x100235320 {
		return 0x100235320, true, nil
	}
	// 0x100235320: ldr x11, [x17, w11, uxtw #3]
	value401, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value401)
	if runtime != nil && runtime.stopPC == 0x100235324 {
		return 0x100235324, true, nil
	}
	// 0x100235324: mov x12, #-0x9c76
	s.x[12] = uint64(uint64(0xffffffffffff638a))
	if runtime != nil && runtime.stopPC == 0x100235328 {
		return 0x100235328, true, nil
	}
	// 0x100235328: movk x12, #0xb4d7, lsl #16
	s.x[12] = uint64(insertBits(s.x[12], 0xb4d7, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x10023532c {
		return 0x10023532c, true, nil
	}
	// 0x10023532c: add x2, x11, x12
	s.x[2] = uint64((s.x[11] + s.x[12]))
	if runtime != nil && runtime.stopPC == 0x100235330 {
		return 0x100235330, true, nil
	}
	// 0x100235330: ldr w10, [x2, w10, uxtw #2]
	value405, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value405))
	if runtime != nil && runtime.stopPC == 0x100235334 {
		return 0x100235334, true, nil
	}
	// 0x100235334: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235338 {
		return 0x100235338, true, nil
	}
	// 0x100235338: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023533c {
		return 0x10023533c, true, nil
	}
	// 0x10023533c: mov w4, #0x48a9
	s.x[4] = uint64(uint32(uint64(0x48a9)))
	if runtime != nil && runtime.stopPC == 0x100235340 {
		return 0x100235340, true, nil
	}
	// 0x100235340: movk w4, #0xaf, lsl #16
	s.x[4] = uint64(uint32(insertBits(uint64(uint32(s.x[4])), 0xaf, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235344 {
		return 0x100235344, true, nil
	}
	// 0x100235344: and w9, w4, w10, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[10])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235348 {
		return 0x100235348, true, nil
	}
	// 0x100235348: mov w10, #0x8bd8
	s.x[10] = uint64(uint32(uint64(0x8bd8)))
	if runtime != nil && runtime.stopPC == 0x10023534c {
		return 0x10023534c, true, nil
	}
	// 0x10023534c: movk w10, #0x73a0, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x73a0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235350 {
		return 0x100235350, true, nil
	}
	// 0x100235350: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235354 {
		return 0x100235354, true, nil
	}
	// 0x100235354: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235358 {
		return 0x100235358, true, nil
	}
	// 0x100235358: ldr w9, [sp, #0x2c4]
	value415, err := s.load(addOffset(s.sp, 708), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value415))
	if runtime != nil && runtime.stopPC == 0x10023535c {
		return 0x10023535c, true, nil
	}
	// 0x10023535c: str w30, [sp, #0x3f8]
	if err := s.storeValue(addOffset(s.sp, 1016), uint64(uint32(s.x[30])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235360 {
		return 0x100235360, true, nil
	}
	// 0x100235360: eor w9, w9, w30
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[30])))))
	if runtime != nil && runtime.stopPC == 0x100235364 {
		return 0x100235364, true, nil
	}
	// 0x100235364: ubfx w10, w19, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[19])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235368 {
		return 0x100235368, true, nil
	}
	// 0x100235368: eor w11, w10, w27
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x10023536c {
		return 0x10023536c, true, nil
	}
	// 0x10023536c: ubfx w12, w9, #8, #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235370 {
		return 0x100235370, true, nil
	}
	// 0x100235370: ldr w14, [sp, #0x3c8]
	value421, err := s.load(addOffset(s.sp, 968), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value421))
	if runtime != nil && runtime.stopPC == 0x100235374 {
		return 0x100235374, true, nil
	}
	// 0x100235374: eor w14, w14, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235378 {
		return 0x100235378, true, nil
	}
	// 0x100235378: sub w11, w10, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023537c {
		return 0x10023537c, true, nil
	}
	// 0x10023537c: eor w11, w11, w14
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235380 {
		return 0x100235380, true, nil
	}
	// 0x100235380: eor w10, w10, w0
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235384 {
		return 0x100235384, true, nil
	}
	// 0x100235384: ldr w10, [x24, w10, uxtw #2]
	value426, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value426))
	if runtime != nil && runtime.stopPC == 0x100235388 {
		return 0x100235388, true, nil
	}
	// 0x100235388: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023538c {
		return 0x10023538c, true, nil
	}
	// 0x10023538c: ldr w11, [sp, #0x3d0]
	value428, err := s.load(addOffset(s.sp, 976), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value428))
	if runtime != nil && runtime.stopPC == 0x100235390 {
		return 0x100235390, true, nil
	}
	// 0x100235390: str w25, [sp, #0x400]
	if err := s.storeValue(addOffset(s.sp, 1024), uint64(uint32(s.x[25])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100235394 {
		return 0x100235394, true, nil
	}
	// 0x100235394: eor w11, w25, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235398 {
		return 0x100235398, true, nil
	}
	// 0x100235398: lsr w11, w11, #0x18
	s.x[11] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x10023539c {
		return 0x10023539c, true, nil
	}
	// 0x10023539c: eor w11, w11, #0x3e
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(0x3e))))
	if runtime != nil && runtime.stopPC == 0x1002353a0 {
		return 0x1002353a0, true, nil
	}
	// 0x1002353a0: ldr w11, [x22, w11, uxtw #2]
	value433, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value433))
	if runtime != nil && runtime.stopPC == 0x1002353a4 {
		return 0x1002353a4, true, nil
	}
	// 0x1002353a4: eor w10, w10, w11, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ rotateRight(uint64(uint32(s.x[11])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002353a8 {
		return 0x1002353a8, true, nil
	}
	// 0x1002353a8: eor w11, w12, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002353ac {
		return 0x1002353ac, true, nil
	}
	// 0x1002353ac: ldr w11, [x26, w11, uxtw #2]
	value436, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value436))
	if runtime != nil && runtime.stopPC == 0x1002353b0 {
		return 0x1002353b0, true, nil
	}
	// 0x1002353b0: ldr w12, [sp, #0x38c]
	value437, err := s.load(addOffset(s.sp, 908), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value437))
	if runtime != nil && runtime.stopPC == 0x1002353b4 {
		return 0x1002353b4, true, nil
	}
	// 0x1002353b4: str w3, [sp, #0x3f4]
	if err := s.storeValue(addOffset(s.sp, 1012), uint64(uint32(s.x[3])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002353b8 {
		return 0x1002353b8, true, nil
	}
	// 0x1002353b8: eor w12, w3, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002353bc {
		return 0x1002353bc, true, nil
	}
	// 0x1002353bc: and w12, w12, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002353c0 {
		return 0x1002353c0, true, nil
	}
	// 0x1002353c0: eor w12, w12, w7
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x1002353c4 {
		return 0x1002353c4, true, nil
	}
	// 0x1002353c4: ldr w12, [x2, w12, uxtw #2]
	value442, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value442))
	if runtime != nil && runtime.stopPC == 0x1002353c8 {
		return 0x1002353c8, true, nil
	}
	// 0x1002353c8: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002353cc {
		return 0x1002353cc, true, nil
	}
	// 0x1002353cc: and w12, w4, w12, lsr #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[12])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002353d0 {
		return 0x1002353d0, true, nil
	}
	// 0x1002353d0: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002353d4 {
		return 0x1002353d4, true, nil
	}
	// 0x1002353d4: eor w11, w10, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002353d8 {
		return 0x1002353d8, true, nil
	}
	// 0x1002353d8: lsr w10, w9, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002353dc {
		return 0x1002353dc, true, nil
	}
	// 0x1002353dc: ldr w10, [x22, w10, uxtw #2]
	value448, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value448))
	if runtime != nil && runtime.stopPC == 0x1002353e0 {
		return 0x1002353e0, true, nil
	}
	// 0x1002353e0: ubfx w12, w16, #0x10, #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002353e4 {
		return 0x1002353e4, true, nil
	}
	// 0x1002353e4: eor w14, w12, w27
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x1002353e8 {
		return 0x1002353e8, true, nil
	}
	// 0x1002353e8: sub w14, w12, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) - uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002353ec {
		return 0x1002353ec, true, nil
	}
	// 0x1002353ec: eor w12, w12, w0
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002353f0 {
		return 0x1002353f0, true, nil
	}
	// 0x1002353f0: ldr w12, [x24, w12, uxtw #2]
	value453, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value453))
	if runtime != nil && runtime.stopPC == 0x1002353f4 {
		return 0x1002353f4, true, nil
	}
	// 0x1002353f4: ldr w17, [sp, #0x35c]
	value454, err := s.load(addOffset(s.sp, 860), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value454))
	if runtime != nil && runtime.stopPC == 0x1002353f8 {
		return 0x1002353f8, true, nil
	}
	// 0x1002353f8: eor w14, w17, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002353fc {
		return 0x1002353fc, true, nil
	}
	// 0x1002353fc: ubfx w3, w15, #8, #8
	s.x[3] = uint64(uint32((uint64(uint32(s.x[15])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235400 {
		return 0x100235400, true, nil
	}
	// 0x100235400: eor w14, w14, w3
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100235404 {
		return 0x100235404, true, nil
	}
	// 0x100235404: eor w10, w14, w10, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ rotateRight(uint64(uint32(s.x[10])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235408 {
		return 0x100235408, true, nil
	}
	// 0x100235408: eor w14, w3, w13
	s.x[14] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x10023540c {
		return 0x10023540c, true, nil
	}
	// 0x10023540c: ldr w14, [x26, w14, uxtw #2]
	value460, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value460))
	if runtime != nil && runtime.stopPC == 0x100235410 {
		return 0x100235410, true, nil
	}
	// 0x100235410: eor w14, w12, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235414 {
		return 0x100235414, true, nil
	}
	// 0x100235414: lsr w3, w19, #0x18
	s.x[3] = uint64(uint32(shiftRight(uint64(uint32(s.x[19])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235418 {
		return 0x100235418, true, nil
	}
	// 0x100235418: and w12, w19, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x10023541c {
		return 0x10023541c, true, nil
	}
	// 0x10023541c: ldr w12, [x2, w12, uxtw #2]
	value464, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value464))
	if runtime != nil && runtime.stopPC == 0x100235420 {
		return 0x100235420, true, nil
	}
	// 0x100235420: ldr w7, [x22, w3, uxtw #2]
	value465, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value465))
	if runtime != nil && runtime.stopPC == 0x100235424 {
		return 0x100235424, true, nil
	}
	// 0x100235424: eor w14, w14, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235428 {
		return 0x100235428, true, nil
	}
	// 0x100235428: eor w3, w10, w14
	s.x[3] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10023542c {
		return 0x10023542c, true, nil
	}
	// 0x10023542c: ubfx w9, w9, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235430 {
		return 0x100235430, true, nil
	}
	// 0x100235430: ubfx w10, w16, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235434 {
		return 0x100235434, true, nil
	}
	// 0x100235434: eor w14, w9, w27
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235438 {
		return 0x100235438, true, nil
	}
	// 0x100235438: ldr w16, [sp, #0x350]
	value471, err := s.load(addOffset(s.sp, 848), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value471))
	if runtime != nil && runtime.stopPC == 0x10023543c {
		return 0x10023543c, true, nil
	}
	// 0x10023543c: eor w16, w16, w10
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235440 {
		return 0x100235440, true, nil
	}
	// 0x100235440: sub w14, w9, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235444 {
		return 0x100235444, true, nil
	}
	// 0x100235444: eor w14, w16, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235448 {
		return 0x100235448, true, nil
	}
	// 0x100235448: eor w14, w14, w7, ror #23
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ rotateRight(uint64(uint32(s.x[7])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x10023544c {
		return 0x10023544c, true, nil
	}
	// 0x10023544c: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235450 {
		return 0x100235450, true, nil
	}
	// 0x100235450: ldr w9, [x24, w9, uxtw #2]
	value477, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value477))
	if runtime != nil && runtime.stopPC == 0x100235454 {
		return 0x100235454, true, nil
	}
	// 0x100235454: eor w10, w10, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235458 {
		return 0x100235458, true, nil
	}
	// 0x100235458: ldr w10, [x26, w10, uxtw #2]
	value479, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value479))
	if runtime != nil && runtime.stopPC == 0x10023545c {
		return 0x10023545c, true, nil
	}
	// 0x10023545c: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235460 {
		return 0x100235460, true, nil
	}
	// 0x100235460: and w10, w15, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235464 {
		return 0x100235464, true, nil
	}
	// 0x100235464: ldr w10, [x2, w10, uxtw #2]
	value482, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value482))
	if runtime != nil && runtime.stopPC == 0x100235468 {
		return 0x100235468, true, nil
	}
	// 0x100235468: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023546c {
		return 0x10023546c, true, nil
	}
	// 0x10023546c: eor w16, w14, w9
	s.x[16] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235470 {
		return 0x100235470, true, nil
	}
	// 0x100235470: mov w9, #0x4120
	s.x[9] = uint64(uint32(uint64(0x4120)))
	if runtime != nil && runtime.stopPC == 0x100235474 {
		return 0x100235474, true, nil
	}
	// 0x100235474: movk w9, #0x1c84, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x1c84, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235478 {
		return 0x100235478, true, nil
	}
	// 0x100235478: and w9, w9, w8, lsl #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & shiftLeft(uint64(uint32(s.x[8])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x10023547c {
		return 0x10023547c, true, nil
	}
	// 0x10023547c: sub w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235480 {
		return 0x100235480, true, nil
	}
	// 0x100235480: mov w9, #0x2090
	s.x[9] = uint64(uint32(uint64(0x2090)))
	if runtime != nil && runtime.stopPC == 0x100235484 {
		return 0x100235484, true, nil
	}
	// 0x100235484: movk w9, #0xe42, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xe42, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235488 {
		return 0x100235488, true, nil
	}
	// 0x100235488: add w15, w8, w9
	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023548c {
		return 0x10023548c, true, nil
	}
	// 0x10023548c: mov w5, #0xdd
	s.x[5] = uint64(uint32(uint64(0xdd)))
	if runtime != nil && runtime.stopPC == 0x100235490 {
		return 0x100235490, true, nil
	}
	// 0x100235490: eor w8, w5, w15, lsr #24
	s.x[8] = uint64(uint32((uint64(uint32(s.x[5])) ^ shiftRight(uint64(uint32(s.x[15])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235494 {
		return 0x100235494, true, nil
	}
	// 0x100235494: ldr w8, [x22, w8, uxtw #2]
	value494, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value494))
	if runtime != nil && runtime.stopPC == 0x100235498 {
		return 0x100235498, true, nil
	}
	// 0x100235498: ldr w9, [sp, #0x348]
	value495, err := s.load(addOffset(s.sp, 840), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value495))
	if runtime != nil && runtime.stopPC == 0x10023549c {
		return 0x10023549c, true, nil
	}
	// 0x10023549c: eor w14, w9, w8, ror #23
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002354a0 {
		return 0x1002354a0, true, nil
	}
	// 0x1002354a0: mov w8, #0x38ab
	s.x[8] = uint64(uint32(uint64(0x38ab)))
	if runtime != nil && runtime.stopPC == 0x1002354a4 {
		return 0x1002354a4, true, nil
	}
	// 0x1002354a4: movk w8, #0x6a62, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x6a62, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354a8 {
		return 0x1002354a8, true, nil
	}
	// 0x1002354a8: eor w8, w11, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002354ac {
		return 0x1002354ac, true, nil
	}
	// 0x1002354ac: ubfx w25, w8, #0x10, #8
	s.x[25] = uint64(uint32((uint64(uint32(s.x[8])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354b0 {
		return 0x1002354b0, true, nil
	}
	// 0x1002354b0: eor w9, w25, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002354b4 {
		return 0x1002354b4, true, nil
	}
	// 0x1002354b4: ldr w30, [x24, w9, uxtw #2]
	value502, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[30] = uint64(uint32(value502))
	if runtime != nil && runtime.stopPC == 0x1002354b8 {
		return 0x1002354b8, true, nil
	}
	// 0x1002354b8: mov w9, #0x731b
	s.x[9] = uint64(uint32(uint64(0x731b)))
	if runtime != nil && runtime.stopPC == 0x1002354bc {
		return 0x1002354bc, true, nil
	}
	// 0x1002354bc: movk w9, #0x8fa0, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x8fa0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354c0 {
		return 0x1002354c0, true, nil
	}
	// 0x1002354c0: eor w9, w14, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002354c4 {
		return 0x1002354c4, true, nil
	}
	// 0x1002354c4: mov w17, #0xe2e6
	s.x[17] = uint64(uint32(uint64(0xe2e6)))
	if runtime != nil && runtime.stopPC == 0x1002354c8 {
		return 0x1002354c8, true, nil
	}
	// 0x1002354c8: movk w17, #0xd783, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xd783, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354cc {
		return 0x1002354cc, true, nil
	}
	// 0x1002354cc: eor w7, w30, w17
	s.x[7] = uint64(uint32((uint64(uint32(s.x[30])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002354d0 {
		return 0x1002354d0, true, nil
	}
	// 0x1002354d0: and w7, w9, w7
	s.x[7] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x1002354d4 {
		return 0x1002354d4, true, nil
	}
	// 0x1002354d4: mov w9, #0x3a72
	s.x[9] = uint64(uint32(uint64(0x3a72)))
	if runtime != nil && runtime.stopPC == 0x1002354d8 {
		return 0x1002354d8, true, nil
	}
	// 0x1002354d8: movk w9, #0xfe09, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xfe09, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354dc {
		return 0x1002354dc, true, nil
	}
	// 0x1002354dc: eor w9, w30, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[30])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002354e0 {
		return 0x1002354e0, true, nil
	}
	// 0x1002354e0: mov w17, #0x6471
	s.x[17] = uint64(uint32(uint64(0x6471)))
	if runtime != nil && runtime.stopPC == 0x1002354e4 {
		return 0x1002354e4, true, nil
	}
	// 0x1002354e4: movk w17, #0x403e, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x403e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354e8 {
		return 0x1002354e8, true, nil
	}
	// 0x1002354e8: eor w19, w14, w17
	s.x[19] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002354ec {
		return 0x1002354ec, true, nil
	}
	// 0x1002354ec: and w9, w19, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002354f0 {
		return 0x1002354f0, true, nil
	}
	// 0x1002354f0: and w10, w4, w10, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[10])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002354f4 {
		return 0x1002354f4, true, nil
	}
	// 0x1002354f4: mov w17, #0xcb3e
	s.x[17] = uint64(uint32(uint64(0xcb3e)))
	if runtime != nil && runtime.stopPC == 0x1002354f8 {
		return 0x1002354f8, true, nil
	}
	// 0x1002354f8: movk w17, #0x3875, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x3875, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002354fc {
		return 0x1002354fc, true, nil
	}
	// 0x1002354fc: eor w19, w16, w17
	s.x[19] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235500 {
		return 0x100235500, true, nil
	}
	// 0x100235500: eor w16, w19, w10
	s.x[16] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235504 {
		return 0x100235504, true, nil
	}
	// 0x100235504: and w10, w16, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235508 {
		return 0x100235508, true, nil
	}
	// 0x100235508: ldr w10, [x2, w10, uxtw #2]
	value523, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value523))
	if runtime != nil && runtime.stopPC == 0x10023550c {
		return 0x10023550c, true, nil
	}
	// 0x10023550c: lsr w1, w10, #4
	s.x[1] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235510 {
		return 0x100235510, true, nil
	}
	// 0x100235510: mov w17, #0x8a9
	s.x[17] = uint64(uint32(uint64(0x8a9)))
	if runtime != nil && runtime.stopPC == 0x100235514 {
		return 0x100235514, true, nil
	}
	// 0x100235514: movk w17, #0x26, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x26, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235518 {
		return 0x100235518, true, nil
	}
	// 0x100235518: and w28, w1, w17
	s.x[28] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10023551c {
		return 0x10023551c, true, nil
	}
	// 0x10023551c: mov w17, #0x18f9
	s.x[17] = uint64(uint32(uint64(0x18f9)))
	if runtime != nil && runtime.stopPC == 0x100235520 {
		return 0x100235520, true, nil
	}
	// 0x100235520: movk w17, #0x9936, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x9936, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235524 {
		return 0x100235524, true, nil
	}
	// 0x100235524: and w17, w10, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235528 {
		return 0x100235528, true, nil
	}
	// 0x100235528: eor w17, w28, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[28])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10023552c {
		return 0x10023552c, true, nil
	}
	// 0x10023552c: mov w28, #0x4000
	s.x[28] = uint64(uint32(uint64(0x4000)))
	if runtime != nil && runtime.stopPC == 0x100235530 {
		return 0x100235530, true, nil
	}
	// 0x100235530: movk w28, #0x89, lsl #16
	s.x[28] = uint64(uint32(insertBits(uint64(uint32(s.x[28])), 0x89, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235534 {
		return 0x100235534, true, nil
	}
	// 0x100235534: and w1, w1, w28
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x100235538 {
		return 0x100235538, true, nil
	}
	// 0x100235538: mov w28, #0xe706
	s.x[28] = uint64(uint32(uint64(0xe706)))
	if runtime != nil && runtime.stopPC == 0x10023553c {
		return 0x10023553c, true, nil
	}
	// 0x10023553c: movk w28, #0x66c9, lsl #16
	s.x[28] = uint64(uint32(insertBits(uint64(uint32(s.x[28])), 0x66c9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235540 {
		return 0x100235540, true, nil
	}
	// 0x100235540: and w10, w10, w28
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x100235544 {
		return 0x100235544, true, nil
	}
	// 0x100235544: mov w28, #0x51a3
	s.x[28] = uint64(uint32(uint64(0x51a3)))
	if runtime != nil && runtime.stopPC == 0x100235548 {
		return 0x100235548, true, nil
	}
	// 0x100235548: movk w28, #0x9fd, lsl #16
	s.x[28] = uint64(uint32(insertBits(uint64(uint32(s.x[28])), 0x9fd, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023554c {
		return 0x10023554c, true, nil
	}
	// 0x10023554c: eor w10, w10, w28
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x100235550 {
		return 0x100235550, true, nil
	}
	// 0x100235550: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235554 {
		return 0x100235554, true, nil
	}
	// 0x100235554: eor w1, w17, #0xffffff9f
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(0xffffff9f))))
	if runtime != nil && runtime.stopPC == 0x100235558 {
		return 0x100235558, true, nil
	}
	// 0x100235558: and w10, w1, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023555c {
		return 0x10023555c, true, nil
	}
	// 0x10023555c: mov w1, #0x858
	s.x[1] = uint64(uint32(uint64(0x858)))
	if runtime != nil && runtime.stopPC == 0x100235560 {
		return 0x100235560, true, nil
	}
	// 0x100235560: movk w1, #0x9002, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x9002, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235564 {
		return 0x100235564, true, nil
	}
	// 0x100235564: and w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235568 {
		return 0x100235568, true, nil
	}
	// 0x100235568: orr w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) | uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10023556c {
		return 0x10023556c, true, nil
	}
	// 0x10023556c: ubfx w17, w16, #8, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235570 {
		return 0x100235570, true, nil
	}
	// 0x100235570: eor w1, w17, w13
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235574 {
		return 0x100235574, true, nil
	}
	// 0x100235574: ldr w1, [x26, w1, uxtw #2]
	value550, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value550))
	if runtime != nil && runtime.stopPC == 0x100235578 {
		return 0x100235578, true, nil
	}
	// 0x100235578: mov w28, #0xe895
	s.x[28] = uint64(uint32(uint64(0xe895)))
	if runtime != nil && runtime.stopPC == 0x10023557c {
		return 0x10023557c, true, nil
	}
	// 0x10023557c: movk w28, #0x3061, lsl #16
	s.x[28] = uint64(uint32(insertBits(uint64(uint32(s.x[28])), 0x3061, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235580 {
		return 0x100235580, true, nil
	}
	// 0x100235580: and w28, w30, w28
	s.x[28] = uint64(uint32((uint64(uint32(s.x[30])) & uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x100235584 {
		return 0x100235584, true, nil
	}
	// 0x100235584: adrp x30, #0x10105b000
	s.x[30] = uint64(0x10105b000)
	if runtime != nil && runtime.stopPC == 0x100235588 {
		return 0x100235588, true, nil
	}
	// 0x100235588: add x30, x30, #0x6f0
	s.x[30] = uint64((s.x[30] + uint64(0x6f0)))
	if runtime != nil && runtime.stopPC == 0x10023558c {
		return 0x10023558c, true, nil
	}
	// 0x10023558c: eor w1, w28, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[28])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235590 {
		return 0x100235590, true, nil
	}
	// 0x100235590: and w12, w4, w12, lsr #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[12])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235594 {
		return 0x100235594, true, nil
	}
	// 0x100235594: mov w28, #0x544b
	s.x[28] = uint64(uint32(uint64(0x544b)))
	if runtime != nil && runtime.stopPC == 0x100235598 {
		return 0x100235598, true, nil
	}
	// 0x100235598: movk w28, #0xdcc9, lsl #16
	s.x[28] = uint64(uint32(insertBits(uint64(uint32(s.x[28])), 0xdcc9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023559c {
		return 0x10023559c, true, nil
	}
	// 0x10023559c: eor w3, w3, w28
	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x1002355a0 {
		return 0x1002355a0, true, nil
	}
	// 0x1002355a0: eor w28, w3, w12
	s.x[28] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002355a4 {
		return 0x1002355a4, true, nil
	}
	// 0x1002355a4: mov w12, #0x276b
	s.x[12] = uint64(uint32(uint64(0x276b)))
	if runtime != nil && runtime.stopPC == 0x1002355a8 {
		return 0x1002355a8, true, nil
	}
	// 0x1002355a8: movk w12, #0xd675, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xd675, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002355ac {
		return 0x1002355ac, true, nil
	}
	// 0x1002355ac: and w12, w14, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002355b0 {
		return 0x1002355b0, true, nil
	}
	// 0x1002355b0: and w14, w28, #0xff
	s.x[14] = uint64(uint32((uint64(uint32(s.x[28])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002355b4 {
		return 0x1002355b4, true, nil
	}
	// 0x1002355b4: ldr w14, [x2, w14, uxtw #2]
	value566, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value566))
	if runtime != nil && runtime.stopPC == 0x1002355b8 {
		return 0x1002355b8, true, nil
	}
	// 0x1002355b8: eor w1, w1, w14
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002355bc {
		return 0x1002355bc, true, nil
	}
	// 0x1002355bc: eor w12, w1, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002355c0 {
		return 0x1002355c0, true, nil
	}
	// 0x1002355c0: eor w12, w12, w7
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x1002355c4 {
		return 0x1002355c4, true, nil
	}
	// 0x1002355c4: and w14, w4, w14, lsr #4
	s.x[14] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[14])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002355c8 {
		return 0x1002355c8, true, nil
	}
	// 0x1002355c8: eor w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002355cc {
		return 0x1002355cc, true, nil
	}
	// 0x1002355cc: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002355d0 {
		return 0x1002355d0, true, nil
	}
	// 0x1002355d0: eor w12, w25, w27
	s.x[12] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x1002355d4 {
		return 0x1002355d4, true, nil
	}
	// 0x1002355d4: sub w12, w25, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[25])) - uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002355d8 {
		return 0x1002355d8, true, nil
	}
	// 0x1002355d8: eor w12, w17, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002355dc {
		return 0x1002355dc, true, nil
	}
	// 0x1002355dc: mov w14, #0x5bd4
	s.x[14] = uint64(uint32(uint64(0x5bd4)))
	if runtime != nil && runtime.stopPC == 0x1002355e0 {
		return 0x1002355e0, true, nil
	}
	// 0x1002355e0: movk w14, #0x3459, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x3459, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002355e4 {
		return 0x1002355e4, true, nil
	}
	// 0x1002355e4: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002355e8 {
		return 0x1002355e8, true, nil
	}
	// 0x1002355e8: eor w14, w9, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002355ec {
		return 0x1002355ec, true, nil
	}
	// 0x1002355ec: lsr w9, w3, #0x18
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[3])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002355f0 {
		return 0x1002355f0, true, nil
	}
	// 0x1002355f0: ldr w9, [x22, w9, uxtw #2]
	value581, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value581))
	if runtime != nil && runtime.stopPC == 0x1002355f4 {
		return 0x1002355f4, true, nil
	}
	// 0x1002355f4: ubfx w12, w15, #0x10, #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[15])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002355f8 {
		return 0x1002355f8, true, nil
	}
	// 0x1002355f8: eor w17, w12, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x1002355fc {
		return 0x1002355fc, true, nil
	}
	// 0x1002355fc: sub w17, w12, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[12])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235600 {
		return 0x100235600, true, nil
	}
	// 0x100235600: eor w12, w12, w0
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235604 {
		return 0x100235604, true, nil
	}
	// 0x100235604: ldr w12, [x24, w12, uxtw #2]
	value586, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value586))
	if runtime != nil && runtime.stopPC == 0x100235608 {
		return 0x100235608, true, nil
	}
	// 0x100235608: ldr w1, [sp, #0x3cc]
	value587, err := s.load(addOffset(s.sp, 972), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value587))
	if runtime != nil && runtime.stopPC == 0x10023560c {
		return 0x10023560c, true, nil
	}
	// 0x10023560c: eor w17, w1, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235610 {
		return 0x100235610, true, nil
	}
	// 0x100235610: eor w12, w17, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235614 {
		return 0x100235614, true, nil
	}
	// 0x100235614: eor w9, w12, w9, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ rotateRight(uint64(uint32(s.x[9])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235618 {
		return 0x100235618, true, nil
	}
	// 0x100235618: ubfx w12, w8, #8, #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[8])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x10023561c {
		return 0x10023561c, true, nil
	}
	// 0x10023561c: eor w17, w12, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235620 {
		return 0x100235620, true, nil
	}
	// 0x100235620: ldr w17, [x26, w17, uxtw #2]
	value593, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value593))
	if runtime != nil && runtime.stopPC == 0x100235624 {
		return 0x100235624, true, nil
	}
	// 0x100235624: eor w12, w17, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235628 {
		return 0x100235628, true, nil
	}
	// 0x100235628: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023562c {
		return 0x10023562c, true, nil
	}
	// 0x10023562c: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235630 {
		return 0x100235630, true, nil
	}
	// 0x100235630: ubfx w9, w28, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[28])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235634 {
		return 0x100235634, true, nil
	}
	// 0x100235634: eor w12, w9, w27
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235638 {
		return 0x100235638, true, nil
	}
	// 0x100235638: sub w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023563c {
		return 0x10023563c, true, nil
	}
	// 0x10023563c: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235640 {
		return 0x100235640, true, nil
	}
	// 0x100235640: ldr w9, [x24, w9, uxtw #2]
	value601, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value601))
	if runtime != nil && runtime.stopPC == 0x100235644 {
		return 0x100235644, true, nil
	}
	// 0x100235644: ldr w17, [sp, #0x344]
	value602, err := s.load(addOffset(s.sp, 836), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value602))
	if runtime != nil && runtime.stopPC == 0x100235648 {
		return 0x100235648, true, nil
	}
	// 0x100235648: eor w12, w17, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023564c {
		return 0x10023564c, true, nil
	}
	// 0x10023564c: ubfx w17, w15, #8, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235650 {
		return 0x100235650, true, nil
	}
	// 0x100235650: eor w12, w12, w17
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235654 {
		return 0x100235654, true, nil
	}
	// 0x100235654: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235658 {
		return 0x100235658, true, nil
	}
	// 0x100235658: lsr w12, w19, #0x18
	s.x[12] = uint64(uint32(shiftRight(uint64(uint32(s.x[19])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x10023565c {
		return 0x10023565c, true, nil
	}
	// 0x10023565c: ldr w12, [x22, w12, uxtw #2]
	value608, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value608))
	if runtime != nil && runtime.stopPC == 0x100235660 {
		return 0x100235660, true, nil
	}
	// 0x100235660: eor w12, w9, w12, ror #23
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ rotateRight(uint64(uint32(s.x[12])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235664 {
		return 0x100235664, true, nil
	}
	// 0x100235664: eor w9, w17, w13
	s.x[9] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235668 {
		return 0x100235668, true, nil
	}
	// 0x100235668: ldr w17, [x26, w9, uxtw #2]
	value611, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value611))
	if runtime != nil && runtime.stopPC == 0x10023566c {
		return 0x10023566c, true, nil
	}
	// 0x10023566c: and w9, w11, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235670 {
		return 0x100235670, true, nil
	}
	// 0x100235670: mov w11, #0x93
	s.x[11] = uint64(uint32(uint64(0x93)))
	if runtime != nil && runtime.stopPC == 0x100235674 {
		return 0x100235674, true, nil
	}
	// 0x100235674: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235678 {
		return 0x100235678, true, nil
	}
	// 0x100235678: ldr w9, [x2, w9, uxtw #2]
	value615, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value615))
	if runtime != nil && runtime.stopPC == 0x10023567c {
		return 0x10023567c, true, nil
	}
	// 0x10023567c: eor w11, w12, w17
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235680 {
		return 0x100235680, true, nil
	}
	// 0x100235680: mov w12, #0x12fd
	s.x[12] = uint64(uint32(uint64(0x12fd)))
	if runtime != nil && runtime.stopPC == 0x100235684 {
		return 0x100235684, true, nil
	}
	// 0x100235684: movk w12, #0xd84d, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xd84d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235688 {
		return 0x100235688, true, nil
	}
	// 0x100235688: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023568c {
		return 0x10023568c, true, nil
	}
	// 0x10023568c: eor w12, w11, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235690 {
		return 0x100235690, true, nil
	}
	// 0x100235690: ubfx w11, w16, #0x10, #8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235694 {
		return 0x100235694, true, nil
	}
	// 0x100235694: ubfx w16, w28, #8, #8
	s.x[16] = uint64(uint32((uint64(uint32(s.x[28])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235698 {
		return 0x100235698, true, nil
	}
	// 0x100235698: eor w17, w11, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x10023569c {
		return 0x10023569c, true, nil
	}
	// 0x10023569c: ldr w1, [sp, #0x340]
	value624, err := s.load(addOffset(s.sp, 832), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value624))
	if runtime != nil && runtime.stopPC == 0x1002356a0 {
		return 0x1002356a0, true, nil
	}
	// 0x1002356a0: eor w1, w1, w16
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1002356a4 {
		return 0x1002356a4, true, nil
	}
	// 0x1002356a4: sub w17, w11, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[11])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002356a8 {
		return 0x1002356a8, true, nil
	}
	// 0x1002356a8: eor w17, w1, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002356ac {
		return 0x1002356ac, true, nil
	}
	// 0x1002356ac: lsr w8, w8, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002356b0 {
		return 0x1002356b0, true, nil
	}
	// 0x1002356b0: ldr w8, [x22, w8, uxtw #2]
	value629, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value629))
	if runtime != nil && runtime.stopPC == 0x1002356b4 {
		return 0x1002356b4, true, nil
	}
	// 0x1002356b4: eor w11, w11, w0
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002356b8 {
		return 0x1002356b8, true, nil
	}
	// 0x1002356b8: ldr w11, [x24, w11, uxtw #2]
	value631, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value631))
	if runtime != nil && runtime.stopPC == 0x1002356bc {
		return 0x1002356bc, true, nil
	}
	// 0x1002356bc: eor w16, w16, w13
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002356c0 {
		return 0x1002356c0, true, nil
	}
	// 0x1002356c0: ldr w16, [x26, w16, uxtw #2]
	value633, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value633))
	if runtime != nil && runtime.stopPC == 0x1002356c4 {
		return 0x1002356c4, true, nil
	}
	// 0x1002356c4: eor w8, w17, w8, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002356c8 {
		return 0x1002356c8, true, nil
	}
	// 0x1002356c8: eor w11, w11, w16
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1002356cc {
		return 0x1002356cc, true, nil
	}
	// 0x1002356cc: eor w16, w8, w11
	s.x[16] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002356d0 {
		return 0x1002356d0, true, nil
	}
	// 0x1002356d0: mov w8, #0x2d0d
	s.x[8] = uint64(uint32(uint64(0x2d0d)))
	if runtime != nil && runtime.stopPC == 0x1002356d4 {
		return 0x1002356d4, true, nil
	}
	// 0x1002356d4: movk w8, #0x9fae, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x9fae, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002356d8 {
		return 0x1002356d8, true, nil
	}
	// 0x1002356d8: eor w8, w8, w14, ror #9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ rotateRight(uint64(uint32(s.x[14])), 9, 32))))
	if runtime != nil && runtime.stopPC == 0x1002356dc {
		return 0x1002356dc, true, nil
	}
	// 0x1002356dc: mov w11, #0xc373
	s.x[11] = uint64(uint32(uint64(0xc373)))
	if runtime != nil && runtime.stopPC == 0x1002356e0 {
		return 0x1002356e0, true, nil
	}
	// 0x1002356e0: movk w11, #0x5fa4, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x5fa4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002356e4 {
		return 0x1002356e4, true, nil
	}
	// 0x1002356e4: eor w11, w11, w8, ror #23
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002356e8 {
		return 0x1002356e8, true, nil
	}
	// 0x1002356e8: mov w8, #0x7ad0
	s.x[8] = uint64(uint32(uint64(0x7ad0)))
	if runtime != nil && runtime.stopPC == 0x1002356ec {
		return 0x1002356ec, true, nil
	}
	// 0x1002356ec: movk w8, #0x11fe, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x11fe, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002356f0 {
		return 0x1002356f0, true, nil
	}
	// 0x1002356f0: eor w8, w10, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002356f4 {
		return 0x1002356f4, true, nil
	}
	// 0x1002356f4: lsr w10, w8, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002356f8 {
		return 0x1002356f8, true, nil
	}
	// 0x1002356f8: ldr w10, [x22, w10, uxtw #2]
	value647, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value647))
	if runtime != nil && runtime.stopPC == 0x1002356fc {
		return 0x1002356fc, true, nil
	}
	// 0x1002356fc: ldr w14, [sp, #0x33c]
	value648, err := s.load(addOffset(s.sp, 828), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value648))
	if runtime != nil && runtime.stopPC == 0x100235700 {
		return 0x100235700, true, nil
	}
	// 0x100235700: eor w10, w14, w10, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ rotateRight(uint64(uint32(s.x[10])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235704 {
		return 0x100235704, true, nil
	}
	// 0x100235704: ubfx w14, w11, #0x10, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235708 {
		return 0x100235708, true, nil
	}
	// 0x100235708: eor w17, w14, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x10023570c {
		return 0x10023570c, true, nil
	}
	// 0x10023570c: sub w17, w14, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235710 {
		return 0x100235710, true, nil
	}
	// 0x100235710: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235714 {
		return 0x100235714, true, nil
	}
	// 0x100235714: and w15, w15, #0xff
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235718 {
		return 0x100235718, true, nil
	}
	// 0x100235718: eor w15, w15, #0x38
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(0x38))))
	if runtime != nil && runtime.stopPC == 0x10023571c {
		return 0x10023571c, true, nil
	}
	// 0x10023571c: ldr w15, [x2, w15, uxtw #2]
	value656, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value656))
	if runtime != nil && runtime.stopPC == 0x100235720 {
		return 0x100235720, true, nil
	}
	// 0x100235720: mov w17, #0x8eb5
	s.x[17] = uint64(uint32(uint64(0x8eb5)))
	if runtime != nil && runtime.stopPC == 0x100235724 {
		return 0x100235724, true, nil
	}
	// 0x100235724: movk w17, #0xb886, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xb886, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235728 {
		return 0x100235728, true, nil
	}
	// 0x100235728: eor w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10023572c {
		return 0x10023572c, true, nil
	}
	// 0x10023572c: eor w3, w16, w15
	s.x[3] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100235730 {
		return 0x100235730, true, nil
	}
	// 0x100235730: and w15, w4, w15, lsr #4
	s.x[15] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[15])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235734 {
		return 0x100235734, true, nil
	}
	// 0x100235734: eor w14, w14, w0
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235738 {
		return 0x100235738, true, nil
	}
	// 0x100235738: ldr w14, [x24, w14, uxtw #2]
	value663, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value663))
	if runtime != nil && runtime.stopPC == 0x10023573c {
		return 0x10023573c, true, nil
	}
	// 0x10023573c: eor w16, w3, w15
	s.x[16] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100235740 {
		return 0x100235740, true, nil
	}
	// 0x100235740: mov w15, #0x6d78
	s.x[15] = uint64(uint32(uint64(0x6d78)))
	if runtime != nil && runtime.stopPC == 0x100235744 {
		return 0x100235744, true, nil
	}
	// 0x100235744: movk w15, #0xd85a, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xd85a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235748 {
		return 0x100235748, true, nil
	}
	// 0x100235748: eor w10, w10, w15
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x10023574c {
		return 0x10023574c, true, nil
	}
	// 0x10023574c: ubfx w15, w16, #8, #8
	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235750 {
		return 0x100235750, true, nil
	}
	// 0x100235750: eor w17, w15, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235754 {
		return 0x100235754, true, nil
	}
	// 0x100235754: ldr w17, [x26, w17, uxtw #2]
	value670, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value670))
	if runtime != nil && runtime.stopPC == 0x100235758 {
		return 0x100235758, true, nil
	}
	// 0x100235758: eor w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10023575c {
		return 0x10023575c, true, nil
	}
	// 0x10023575c: mov w14, #0xd4b1
	s.x[14] = uint64(uint32(uint64(0xd4b1)))
	if runtime != nil && runtime.stopPC == 0x100235760 {
		return 0x100235760, true, nil
	}
	// 0x100235760: movk w14, #0x91e3, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x91e3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235764 {
		return 0x100235764, true, nil
	}
	// 0x100235764: eor w14, w15, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235768 {
		return 0x100235768, true, nil
	}
	// 0x100235768: eor w14, w17, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10023576c {
		return 0x10023576c, true, nil
	}
	// 0x10023576c: and w15, w14, w10
	s.x[15] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235770 {
		return 0x100235770, true, nil
	}
	// 0x100235770: add w10, w14, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235774 {
		return 0x100235774, true, nil
	}
	// 0x100235774: sub w15, w10, w15, lsl #1
	s.x[15] = uint64(uint32((uint64(uint32(s.x[10])) - shiftLeft(uint64(uint32(s.x[15])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100235778 {
		return 0x100235778, true, nil
	}
	// 0x100235778: and w9, w4, w9, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[9])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x10023577c {
		return 0x10023577c, true, nil
	}
	// 0x10023577c: eor w14, w12, w9
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235780 {
		return 0x100235780, true, nil
	}
	// 0x100235780: lsr w9, w12, #0x18
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[12])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235784 {
		return 0x100235784, true, nil
	}
	// 0x100235784: ldr w9, [x22, w9, uxtw #2]
	value682, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value682))
	if runtime != nil && runtime.stopPC == 0x100235788 {
		return 0x100235788, true, nil
	}
	// 0x100235788: ubfx w10, w8, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x10023578c {
		return 0x10023578c, true, nil
	}
	// 0x10023578c: eor w12, w10, w27
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235790 {
		return 0x100235790, true, nil
	}
	// 0x100235790: sub w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235794 {
		return 0x100235794, true, nil
	}
	// 0x100235794: ldr w17, [sp, #0x338]
	value686, err := s.load(addOffset(s.sp, 824), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value686))
	if runtime != nil && runtime.stopPC == 0x100235798 {
		return 0x100235798, true, nil
	}
	// 0x100235798: eor w12, w17, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023579c {
		return 0x10023579c, true, nil
	}
	// 0x10023579c: ubfx w17, w11, #8, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[11])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002357a0 {
		return 0x1002357a0, true, nil
	}
	// 0x1002357a0: eor w12, w12, w17
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002357a4 {
		return 0x1002357a4, true, nil
	}
	// 0x1002357a4: eor w9, w12, w9, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ rotateRight(uint64(uint32(s.x[9])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002357a8 {
		return 0x1002357a8, true, nil
	}
	// 0x1002357a8: eor w10, w10, w0
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002357ac {
		return 0x1002357ac, true, nil
	}
	// 0x1002357ac: ldr w10, [x24, w10, uxtw #2]
	value692, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value692))
	if runtime != nil && runtime.stopPC == 0x1002357b0 {
		return 0x1002357b0, true, nil
	}
	// 0x1002357b0: and w12, w16, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002357b4 {
		return 0x1002357b4, true, nil
	}
	// 0x1002357b4: ldr w12, [x2, w12, uxtw #2]
	value694, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value694))
	if runtime != nil && runtime.stopPC == 0x1002357b8 {
		return 0x1002357b8, true, nil
	}
	// 0x1002357b8: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002357bc {
		return 0x1002357bc, true, nil
	}
	// 0x1002357bc: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002357c0 {
		return 0x1002357c0, true, nil
	}
	// 0x1002357c0: eor w10, w17, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002357c4 {
		return 0x1002357c4, true, nil
	}
	// 0x1002357c4: ldr w10, [x26, w10, uxtw #2]
	value698, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value698))
	if runtime != nil && runtime.stopPC == 0x1002357c8 {
		return 0x1002357c8, true, nil
	}
	// 0x1002357c8: eor w25, w9, w10
	s.x[25] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002357cc {
		return 0x1002357cc, true, nil
	}
	// 0x1002357cc: ubfx w9, w16, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002357d0 {
		return 0x1002357d0, true, nil
	}
	// 0x1002357d0: eor w10, w9, w0
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002357d4 {
		return 0x1002357d4, true, nil
	}
	// 0x1002357d4: ldr w10, [x24, w10, uxtw #2]
	value702, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value702))
	if runtime != nil && runtime.stopPC == 0x1002357d8 {
		return 0x1002357d8, true, nil
	}
	// 0x1002357d8: ubfx w16, w14, #8, #8
	s.x[16] = uint64(uint32((uint64(uint32(s.x[14])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002357dc {
		return 0x1002357dc, true, nil
	}
	// 0x1002357dc: eor w17, w16, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002357e0 {
		return 0x1002357e0, true, nil
	}
	// 0x1002357e0: ldr w17, [x26, w17, uxtw #2]
	value705, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value705))
	if runtime != nil && runtime.stopPC == 0x1002357e4 {
		return 0x1002357e4, true, nil
	}
	// 0x1002357e4: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002357e8 {
		return 0x1002357e8, true, nil
	}
	// 0x1002357e8: and w17, w8, #0xff
	s.x[17] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002357ec {
		return 0x1002357ec, true, nil
	}
	// 0x1002357ec: ldr w17, [x2, w17, uxtw #2]
	value708, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value708))
	if runtime != nil && runtime.stopPC == 0x1002357f0 {
		return 0x1002357f0, true, nil
	}
	// 0x1002357f0: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002357f4 {
		return 0x1002357f4, true, nil
	}
	// 0x1002357f4: and w17, w4, w17, lsr #4
	s.x[17] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[17])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002357f8 {
		return 0x1002357f8, true, nil
	}
	// 0x1002357f8: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002357fc {
		return 0x1002357fc, true, nil
	}
	// 0x1002357fc: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235800 {
		return 0x100235800, true, nil
	}
	// 0x100235800: eor w16, w9, w27
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235804 {
		return 0x100235804, true, nil
	}
	// 0x100235804: sub w9, w9, w16
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235808 {
		return 0x100235808, true, nil
	}
	// 0x100235808: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023580c {
		return 0x10023580c, true, nil
	}
	// 0x10023580c: lsr w10, w11, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235810 {
		return 0x100235810, true, nil
	}
	// 0x100235810: ldr w10, [x22, w10, uxtw #2]
	value717, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value717))
	if runtime != nil && runtime.stopPC == 0x100235814 {
		return 0x100235814, true, nil
	}
	// 0x100235814: ldr w16, [sp, #0x334]
	value718, err := s.load(addOffset(s.sp, 820), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value718))
	if runtime != nil && runtime.stopPC == 0x100235818 {
		return 0x100235818, true, nil
	}
	// 0x100235818: eor w10, w16, w10, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) ^ rotateRight(uint64(uint32(s.x[10])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x10023581c {
		return 0x10023581c, true, nil
	}
	// 0x10023581c: mov w16, #0xd678
	s.x[16] = uint64(uint32(uint64(0xd678)))
	if runtime != nil && runtime.stopPC == 0x100235820 {
		return 0x100235820, true, nil
	}
	// 0x100235820: movk w16, #0x3923, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x3923, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235824 {
		return 0x100235824, true, nil
	}
	// 0x100235824: eor w9, w9, w16
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235828 {
		return 0x100235828, true, nil
	}
	// 0x100235828: eor w19, w10, w9
	s.x[19] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023582c {
		return 0x10023582c, true, nil
	}
	// 0x10023582c: and w9, w14, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235830 {
		return 0x100235830, true, nil
	}
	// 0x100235830: ldr w9, [x2, w9, uxtw #2]
	value725, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value725))
	if runtime != nil && runtime.stopPC == 0x100235834 {
		return 0x100235834, true, nil
	}
	// 0x100235834: and w10, w4, w9, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[9])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235838 {
		return 0x100235838, true, nil
	}
	// 0x100235838: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023583c {
		return 0x10023583c, true, nil
	}
	// 0x10023583c: mov w10, #0xc9a0
	s.x[10] = uint64(uint32(uint64(0xc9a0)))
	if runtime != nil && runtime.stopPC == 0x100235840 {
		return 0x100235840, true, nil
	}
	// 0x100235840: movk w10, #0xf323, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xf323, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235844 {
		return 0x100235844, true, nil
	}
	// 0x100235844: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235848 {
		return 0x100235848, true, nil
	}
	// 0x100235848: eor w16, w9, w15
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x10023584c {
		return 0x10023584c, true, nil
	}
	// 0x10023584c: and w9, w4, w12, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[12])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235850 {
		return 0x100235850, true, nil
	}
	// 0x100235850: eor w15, w25, w9
	s.x[15] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235854 {
		return 0x100235854, true, nil
	}
	// 0x100235854: lsr w9, w19, #0x10
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[19])), int(uint64(0x10)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235858 {
		return 0x100235858, true, nil
	}
	// 0x100235858: and w9, w9, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0x8))))
	if runtime != nil && runtime.stopPC == 0x10023585c {
		return 0x10023585c, true, nil
	}
	// 0x10023585c: mov w10, #0xe4e1
	s.x[10] = uint64(uint32(uint64(0xe4e1)))
	if runtime != nil && runtime.stopPC == 0x100235860 {
		return 0x100235860, true, nil
	}
	// 0x100235860: movk w10, #0x6e17, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x6e17, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235864 {
		return 0x100235864, true, nil
	}
	// 0x100235864: eor w10, w15, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235868 {
		return 0x100235868, true, nil
	}
	// 0x100235868: lsr w10, w10, w9
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(uint32(s.x[9]))), 32)))
	if runtime != nil && runtime.stopPC == 0x10023586c {
		return 0x10023586c, true, nil
	}
	// 0x10023586c: mov w12, #0x9d0c
	s.x[12] = uint64(uint32(uint64(0x9d0c)))
	if runtime != nil && runtime.stopPC == 0x100235870 {
		return 0x100235870, true, nil
	}
	// 0x100235870: movk w12, #0xdadc, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xdadc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235874 {
		return 0x100235874, true, nil
	}
	// 0x100235874: and w12, w12, w10, lsl #1
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & shiftLeft(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100235878 {
		return 0x100235878, true, nil
	}
	// 0x100235878: sub w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023587c {
		return 0x10023587c, true, nil
	}
	// 0x10023587c: mov w12, #0x4e86
	s.x[12] = uint64(uint32(uint64(0x4e86)))
	if runtime != nil && runtime.stopPC == 0x100235880 {
		return 0x100235880, true, nil
	}
	// 0x100235880: movk w12, #0xed6e, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xed6e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235884 {
		return 0x100235884, true, nil
	}
	// 0x100235884: add w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235888 {
		return 0x100235888, true, nil
	}
	// 0x100235888: eor w9, w9, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(0x8))))
	if runtime != nil && runtime.stopPC == 0x10023588c {
		return 0x10023588c, true, nil
	}
	// 0x10023588c: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235890 {
		return 0x100235890, true, nil
	}
	// 0x100235890: lsr w12, w10, w9
	s.x[12] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(uint32(s.x[9]))), 32)))
	if runtime != nil && runtime.stopPC == 0x100235894 {
		return 0x100235894, true, nil
	}
	// 0x100235894: ubfx w9, w14, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235898 {
		return 0x100235898, true, nil
	}
	// 0x100235898: mov w10, #0x7a55
	s.x[10] = uint64(uint32(uint64(0x7a55)))
	if runtime != nil && runtime.stopPC == 0x10023589c {
		return 0x10023589c, true, nil
	}
	// 0x10023589c: movk w10, #0x3c9e, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x3c9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002358a0 {
		return 0x1002358a0, true, nil
	}
	// 0x1002358a0: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002358a4 {
		return 0x1002358a4, true, nil
	}
	// 0x1002358a4: sub w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002358a8 {
		return 0x1002358a8, true, nil
	}
	// 0x1002358a8: and w10, w10, #0x7fffffff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0x7fffffff))))
	if runtime != nil && runtime.stopPC == 0x1002358ac {
		return 0x1002358ac, true, nil
	}
	// 0x1002358ac: ubfx w8, w8, #8, #8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002358b0 {
		return 0x1002358b0, true, nil
	}
	// 0x1002358b0: ldr w14, [sp, #0x330]
	value757, err := s.load(addOffset(s.sp, 816), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value757))
	if runtime != nil && runtime.stopPC == 0x1002358b4 {
		return 0x1002358b4, true, nil
	}
	// 0x1002358b4: eor w14, w14, w8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002358b8 {
		return 0x1002358b8, true, nil
	}
	// 0x1002358b8: eor w10, w14, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002358bc {
		return 0x1002358bc, true, nil
	}
	// 0x1002358bc: lsr w14, w3, #0x18
	s.x[14] = uint64(uint32(shiftRight(uint64(uint32(s.x[3])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002358c0 {
		return 0x1002358c0, true, nil
	}
	// 0x1002358c0: ldr w14, [x22, w14, uxtw #2]
	value761, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value761))
	if runtime != nil && runtime.stopPC == 0x1002358c4 {
		return 0x1002358c4, true, nil
	}
	// 0x1002358c4: eor w10, w10, w14, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ rotateRight(uint64(uint32(s.x[14])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002358c8 {
		return 0x1002358c8, true, nil
	}
	// 0x1002358c8: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002358cc {
		return 0x1002358cc, true, nil
	}
	// 0x1002358cc: ldr w9, [x24, w9, uxtw #2]
	value764, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value764))
	if runtime != nil && runtime.stopPC == 0x1002358d0 {
		return 0x1002358d0, true, nil
	}
	// 0x1002358d0: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002358d4 {
		return 0x1002358d4, true, nil
	}
	// 0x1002358d4: eor w8, w8, w13
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002358d8 {
		return 0x1002358d8, true, nil
	}
	// 0x1002358d8: ldr w8, [x26, w8, uxtw #2]
	value767, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value767))
	if runtime != nil && runtime.stopPC == 0x1002358dc {
		return 0x1002358dc, true, nil
	}
	// 0x1002358dc: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002358e0 {
		return 0x1002358e0, true, nil
	}
	// 0x1002358e0: and w9, w11, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002358e4 {
		return 0x1002358e4, true, nil
	}
	// 0x1002358e4: ldr w9, [x2, w9, uxtw #2]
	value770, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value770))
	if runtime != nil && runtime.stopPC == 0x1002358e8 {
		return 0x1002358e8, true, nil
	}
	// 0x1002358e8: mov w10, #0xc2bb
	s.x[10] = uint64(uint32(uint64(0xc2bb)))
	if runtime != nil && runtime.stopPC == 0x1002358ec {
		return 0x1002358ec, true, nil
	}
	// 0x1002358ec: movk w10, #0x5c5d, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x5c5d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002358f0 {
		return 0x1002358f0, true, nil
	}
	// 0x1002358f0: eor w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002358f4 {
		return 0x1002358f4, true, nil
	}
	// 0x1002358f4: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002358f8 {
		return 0x1002358f8, true, nil
	}
	// 0x1002358f8: lsr w10, w8, #0x1b
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x1b)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002358fc {
		return 0x1002358fc, true, nil
	}
	// 0x1002358fc: mov w11, #0x4925
	s.x[11] = uint64(uint32(uint64(0x4925)))
	if runtime != nil && runtime.stopPC == 0x100235900 {
		return 0x100235900, true, nil
	}
	// 0x100235900: movk w11, #0x2492, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x2492, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235904 {
		return 0x100235904, true, nil
	}
	// 0x100235904: umull x11, w10, w11
	s.x[11] = uint64(uint64(uint32(uint64(uint32(s.x[10])))) * uint64(uint32(uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235908 {
		return 0x100235908, true, nil
	}
	// 0x100235908: lsr x11, x11, #0x20
	s.x[11] = uint64(shiftRight(s.x[11], int(uint64(0x20)), 64))
	if runtime != nil && runtime.stopPC == 0x10023590c {
		return 0x10023590c, true, nil
	}
	// 0x10023590c: lsl w14, w11, #6
	s.x[14] = uint64(uint32(shiftLeft(uint64(uint32(s.x[11])), int(uint64(0x6)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235910 {
		return 0x100235910, true, nil
	}
	// 0x100235910: sub w14, w14, w11, lsl #3
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) - shiftLeft(uint64(uint32(s.x[11])), 3, 32))))
	if runtime != nil && runtime.stopPC == 0x100235914 {
		return 0x100235914, true, nil
	}
	// 0x100235914: and w9, w4, w9, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[9])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235918 {
		return 0x100235918, true, nil
	}
	// 0x100235918: eor w11, w8, w9
	s.x[11] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023591c {
		return 0x10023591c, true, nil
	}
	// 0x10023591c: lsr w8, w8, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235920 {
		return 0x100235920, true, nil
	}
	// 0x100235920: mov w9, #0x25
	s.x[9] = uint64(uint32(uint64(0x25)))
	if runtime != nil && runtime.stopPC == 0x100235924 {
		return 0x100235924, true, nil
	}
	// 0x100235924: mul w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235928 {
		return 0x100235928, true, nil
	}
	// 0x100235928: lsr w9, w9, #8
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x8)), 32)))
	if runtime != nil && runtime.stopPC == 0x10023592c {
		return 0x10023592c, true, nil
	}
	// 0x10023592c: mov w10, #0x38
	s.x[10] = uint64(uint32(uint64(0x38)))
	if runtime != nil && runtime.stopPC == 0x100235930 {
		return 0x100235930, true, nil
	}
	// 0x100235930: msub w8, w9, w10, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[9]))*uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235934 {
		return 0x100235934, true, nil
	}
	// 0x100235934: and w9, w14, #0xf8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(0xf8))))
	if runtime != nil && runtime.stopPC == 0x100235938 {
		return 0x100235938, true, nil
	}
	// 0x100235938: add w9, w9, w8, uxtb
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint8(uint64(uint32(s.x[8])))))))
	if runtime != nil && runtime.stopPC == 0x10023593c {
		return 0x10023593c, true, nil
	}
	// 0x10023593c: ubfx w8, w16, #0x10, #8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235940 {
		return 0x100235940, true, nil
	}
	// 0x100235940: eor w10, w8, w27
	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235944 {
		return 0x100235944, true, nil
	}
	// 0x100235944: ubfx w14, w19, #8, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[19])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235948 {
		return 0x100235948, true, nil
	}
	// 0x100235948: ldr w17, [sp, #0x3e4]
	value795, err := s.load(addOffset(s.sp, 996), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value795))
	if runtime != nil && runtime.stopPC == 0x10023594c {
		return 0x10023594c, true, nil
	}
	// 0x10023594c: eor w17, w14, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235950 {
		return 0x100235950, true, nil
	}
	// 0x100235950: sub w10, w8, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235954 {
		return 0x100235954, true, nil
	}
	// 0x100235954: eor w10, w17, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235958 {
		return 0x100235958, true, nil
	}
	// 0x100235958: mov w17, #0xb3
	s.x[17] = uint64(uint32(uint64(0xb3)))
	if runtime != nil && runtime.stopPC == 0x10023595c {
		return 0x10023595c, true, nil
	}
	// 0x10023595c: eor w17, w17, w25, lsr #24
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ shiftRight(uint64(uint32(s.x[25])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235960 {
		return 0x100235960, true, nil
	}
	// 0x100235960: ldr w17, [x22, w17, uxtw #2]
	value801, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value801))
	if runtime != nil && runtime.stopPC == 0x100235964 {
		return 0x100235964, true, nil
	}
	// 0x100235964: eor w10, w10, w17, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ rotateRight(uint64(uint32(s.x[17])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235968 {
		return 0x100235968, true, nil
	}
	// 0x100235968: eor w14, w14, w13
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x10023596c {
		return 0x10023596c, true, nil
	}
	// 0x10023596c: ldr w14, [x26, w14, uxtw #2]
	value804, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value804))
	if runtime != nil && runtime.stopPC == 0x100235970 {
		return 0x100235970, true, nil
	}
	// 0x100235970: eor w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235974 {
		return 0x100235974, true, nil
	}
	// 0x100235974: eor w8, w8, w0
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235978 {
		return 0x100235978, true, nil
	}
	// 0x100235978: ldr w8, [x24, w8, uxtw #2]
	value807, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value807))
	if runtime != nil && runtime.stopPC == 0x10023597c {
		return 0x10023597c, true, nil
	}
	// 0x10023597c: eor w3, w10, w8
	s.x[3] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235980 {
		return 0x100235980, true, nil
	}
	// 0x100235980: ubfx w10, w15, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235984 {
		return 0x100235984, true, nil
	}
	// 0x100235984: eor w8, w10, w6
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100235988 {
		return 0x100235988, true, nil
	}
	// 0x100235988: mov w14, #0x7a42
	s.x[14] = uint64(uint32(uint64(0x7a42)))
	if runtime != nil && runtime.stopPC == 0x10023598c {
		return 0x10023598c, true, nil
	}
	// 0x10023598c: movk w14, #0xbc9e, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xbc9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235990 {
		return 0x100235990, true, nil
	}
	// 0x100235990: eor w14, w10, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235994 {
		return 0x100235994, true, nil
	}
	// 0x100235994: sub w14, w8, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235998 {
		return 0x100235998, true, nil
	}
	// 0x100235998: lsr w8, w19, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[19])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x10023599c {
		return 0x10023599c, true, nil
	}
	// 0x10023599c: ldr w8, [x22, w8, uxtw #2]
	value816, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value816))
	if runtime != nil && runtime.stopPC == 0x1002359a0 {
		return 0x1002359a0, true, nil
	}
	// 0x1002359a0: ubfx w17, w11, #0x10, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[11])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002359a4 {
		return 0x1002359a4, true, nil
	}
	// 0x1002359a4: eor w1, w17, w0
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002359a8 {
		return 0x1002359a8, true, nil
	}
	// 0x1002359a8: ldr w1, [x24, w1, uxtw #2]
	value819, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value819))
	if runtime != nil && runtime.stopPC == 0x1002359ac {
		return 0x1002359ac, true, nil
	}
	// 0x1002359ac: and w7, w16, #0xff
	s.x[7] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002359b0 {
		return 0x1002359b0, true, nil
	}
	// 0x1002359b0: ldr w7, [x2, w7, uxtw #2]
	value821, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[7])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value821))
	if runtime != nil && runtime.stopPC == 0x1002359b4 {
		return 0x1002359b4, true, nil
	}
	// 0x1002359b4: ldr w6, [sp, #0x3dc]
	value822, err := s.load(addOffset(s.sp, 988), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(uint32(value822))
	if runtime != nil && runtime.stopPC == 0x1002359b8 {
		return 0x1002359b8, true, nil
	}
	// 0x1002359b8: eor w8, w6, w8, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[6])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x1002359bc {
		return 0x1002359bc, true, nil
	}
	// 0x1002359bc: eor w8, w8, w1
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002359c0 {
		return 0x1002359c0, true, nil
	}
	// 0x1002359c0: eor w8, w8, w7
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x1002359c4 {
		return 0x1002359c4, true, nil
	}
	// 0x1002359c4: and w1, w4, w7, lsr #4
	s.x[1] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[7])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002359c8 {
		return 0x1002359c8, true, nil
	}
	// 0x1002359c8: eor w8, w8, w1
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002359cc {
		return 0x1002359cc, true, nil
	}
	// 0x1002359cc: and w12, w12, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002359d0 {
		return 0x1002359d0, true, nil
	}
	// 0x1002359d0: eor w8, w8, w12
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002359d4 {
		return 0x1002359d4, true, nil
	}
	// 0x1002359d4: eor w1, w17, w27
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x1002359d8 {
		return 0x1002359d8, true, nil
	}
	// 0x1002359d8: sub w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002359dc {
		return 0x1002359dc, true, nil
	}
	// 0x1002359dc: eor w8, w8, w17
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002359e0 {
		return 0x1002359e0, true, nil
	}
	// 0x1002359e0: eor w12, w12, w13
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002359e4 {
		return 0x1002359e4, true, nil
	}
	// 0x1002359e4: ldr w12, [x26, w12, uxtw #2]
	value834, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value834))
	if runtime != nil && runtime.stopPC == 0x1002359e8 {
		return 0x1002359e8, true, nil
	}
	// 0x1002359e8: mov w17, #0x2068
	s.x[17] = uint64(uint32(uint64(0x2068)))
	if runtime != nil && runtime.stopPC == 0x1002359ec {
		return 0x1002359ec, true, nil
	}
	// 0x1002359ec: movk w17, #0x3da3, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x3da3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002359f0 {
		return 0x1002359f0, true, nil
	}
	// 0x1002359f0: eor w8, w8, w17
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002359f4 {
		return 0x1002359f4, true, nil
	}
	// 0x1002359f4: eor w8, w8, w12
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002359f8 {
		return 0x1002359f8, true, nil
	}
	// 0x1002359f8: eor w9, w9, w5
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x1002359fc {
		return 0x1002359fc, true, nil
	}
	// 0x1002359fc: ldr w9, [x22, w9, uxtw #2]
	value840, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value840))
	if runtime != nil && runtime.stopPC == 0x100235a00 {
		return 0x100235a00, true, nil
	}
	// 0x100235a00: ldr w12, [sp, #0x328]
	value841, err := s.load(addOffset(s.sp, 808), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value841))
	if runtime != nil && runtime.stopPC == 0x100235a04 {
		return 0x100235a04, true, nil
	}
	// 0x100235a04: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235a08 {
		return 0x100235a08, true, nil
	}
	// 0x100235a08: ubfx w14, w16, #8, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235a0c {
		return 0x100235a0c, true, nil
	}
	// 0x100235a0c: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235a10 {
		return 0x100235a10, true, nil
	}
	// 0x100235a10: eor w9, w12, w9, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ rotateRight(uint64(uint32(s.x[9])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235a14 {
		return 0x100235a14, true, nil
	}
	// 0x100235a14: mov w12, #0x82
	s.x[12] = uint64(uint32(uint64(0x82)))
	if runtime != nil && runtime.stopPC == 0x100235a18 {
		return 0x100235a18, true, nil
	}
	// 0x100235a18: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235a1c {
		return 0x100235a1c, true, nil
	}
	// 0x100235a1c: ldr w12, [x24, w10, uxtw #2]
	value848, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value848))
	if runtime != nil && runtime.stopPC == 0x100235a20 {
		return 0x100235a20, true, nil
	}
	// 0x100235a20: eor w10, w14, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235a24 {
		return 0x100235a24, true, nil
	}
	// 0x100235a24: ldr w14, [x26, w10, uxtw #2]
	value850, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value850))
	if runtime != nil && runtime.stopPC == 0x100235a28 {
		return 0x100235a28, true, nil
	}
	// 0x100235a28: and w10, w19, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235a2c {
		return 0x100235a2c, true, nil
	}
	// 0x100235a2c: ldr w10, [x2, w10, uxtw #2]
	value852, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value852))
	if runtime != nil && runtime.stopPC == 0x100235a30 {
		return 0x100235a30, true, nil
	}
	// 0x100235a30: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235a34 {
		return 0x100235a34, true, nil
	}
	// 0x100235a34: eor w12, w12, w10
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235a38 {
		return 0x100235a38, true, nil
	}
	// 0x100235a38: eor w25, w9, w12
	s.x[25] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235a3c {
		return 0x100235a3c, true, nil
	}
	// 0x100235a3c: ubfx w9, w19, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[19])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235a40 {
		return 0x100235a40, true, nil
	}
	// 0x100235a40: eor w12, w9, w27
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235a44 {
		return 0x100235a44, true, nil
	}
	// 0x100235a44: sub w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235a48 {
		return 0x100235a48, true, nil
	}
	// 0x100235a48: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235a4c {
		return 0x100235a4c, true, nil
	}
	// 0x100235a4c: ldr w9, [x24, w9, uxtw #2]
	value860, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value860))
	if runtime != nil && runtime.stopPC == 0x100235a50 {
		return 0x100235a50, true, nil
	}
	// 0x100235a50: ldr w14, [sp, #0x3bc]
	value861, err := s.load(addOffset(s.sp, 956), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value861))
	if runtime != nil && runtime.stopPC == 0x100235a54 {
		return 0x100235a54, true, nil
	}
	// 0x100235a54: eor w12, w14, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235a58 {
		return 0x100235a58, true, nil
	}
	// 0x100235a58: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235a5c {
		return 0x100235a5c, true, nil
	}
	// 0x100235a5c: lsr w12, w16, #0x18
	s.x[12] = uint64(uint32(shiftRight(uint64(uint32(s.x[16])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235a60 {
		return 0x100235a60, true, nil
	}
	// 0x100235a60: ldr w12, [x22, w12, uxtw #2]
	value865, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value865))
	if runtime != nil && runtime.stopPC == 0x100235a64 {
		return 0x100235a64, true, nil
	}
	// 0x100235a64: ubfx w14, w11, #8, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235a68 {
		return 0x100235a68, true, nil
	}
	// 0x100235a68: eor w14, w14, w13
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235a6c {
		return 0x100235a6c, true, nil
	}
	// 0x100235a6c: ldr w14, [x26, w14, uxtw #2]
	value868, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value868))
	if runtime != nil && runtime.stopPC == 0x100235a70 {
		return 0x100235a70, true, nil
	}
	// 0x100235a70: eor w9, w9, w12, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ rotateRight(uint64(uint32(s.x[12])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235a74 {
		return 0x100235a74, true, nil
	}
	// 0x100235a74: and w12, w15, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235a78 {
		return 0x100235a78, true, nil
	}
	// 0x100235a78: mov w6, #0xd9
	s.x[6] = uint64(uint32(uint64(0xd9)))
	if runtime != nil && runtime.stopPC == 0x100235a7c {
		return 0x100235a7c, true, nil
	}
	// 0x100235a7c: eor w12, w12, w6
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100235a80 {
		return 0x100235a80, true, nil
	}
	// 0x100235a80: ldr w12, [x2, w12, uxtw #2]
	value873, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value873))
	if runtime != nil && runtime.stopPC == 0x100235a84 {
		return 0x100235a84, true, nil
	}
	// 0x100235a84: eor w14, w14, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235a88 {
		return 0x100235a88, true, nil
	}
	// 0x100235a88: mov w15, #0x2cfa
	s.x[15] = uint64(uint32(uint64(0x2cfa)))
	if runtime != nil && runtime.stopPC == 0x100235a8c {
		return 0x100235a8c, true, nil
	}
	// 0x100235a8c: movk w15, #0x458d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x458d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235a90 {
		return 0x100235a90, true, nil
	}
	// 0x100235a90: eor w14, w14, w15
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100235a94 {
		return 0x100235a94, true, nil
	}
	// 0x100235a94: eor w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235a98 {
		return 0x100235a98, true, nil
	}
	// 0x100235a98: and w12, w4, w12, lsr #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[12])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235a9c {
		return 0x100235a9c, true, nil
	}
	// 0x100235a9c: lsr w14, w9, #0x18
	s.x[14] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235aa0 {
		return 0x100235aa0, true, nil
	}
	// 0x100235aa0: ldr w15, [x22, w14, uxtw #2]
	value881, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value881))
	if runtime != nil && runtime.stopPC == 0x100235aa4 {
		return 0x100235aa4, true, nil
	}
	// 0x100235aa4: eor w14, w9, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235aa8 {
		return 0x100235aa8, true, nil
	}
	// 0x100235aa8: ubfx w9, w8, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235aac {
		return 0x100235aac, true, nil
	}
	// 0x100235aac: eor w12, w9, w27
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235ab0 {
		return 0x100235ab0, true, nil
	}
	// 0x100235ab0: sub w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235ab4 {
		return 0x100235ab4, true, nil
	}
	// 0x100235ab4: ldr w16, [sp, #0x320]
	value886, err := s.load(addOffset(s.sp, 800), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value886))
	if runtime != nil && runtime.stopPC == 0x100235ab8 {
		return 0x100235ab8, true, nil
	}
	// 0x100235ab8: eor w12, w16, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235abc {
		return 0x100235abc, true, nil
	}
	// 0x100235abc: eor w12, w12, w15, ror #23
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ rotateRight(uint64(uint32(s.x[15])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235ac0 {
		return 0x100235ac0, true, nil
	}
	// 0x100235ac0: and w15, w11, #0xff
	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235ac4 {
		return 0x100235ac4, true, nil
	}
	// 0x100235ac4: ldr w15, [x2, w15, uxtw #2]
	value890, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value890))
	if runtime != nil && runtime.stopPC == 0x100235ac8 {
		return 0x100235ac8, true, nil
	}
	// 0x100235ac8: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235acc {
		return 0x100235acc, true, nil
	}
	// 0x100235acc: ldr w9, [x24, w9, uxtw #2]
	value892, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value892))
	if runtime != nil && runtime.stopPC == 0x100235ad0 {
		return 0x100235ad0, true, nil
	}
	// 0x100235ad0: eor w12, w12, w9
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235ad4 {
		return 0x100235ad4, true, nil
	}
	// 0x100235ad4: eor w16, w3, w15
	s.x[16] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100235ad8 {
		return 0x100235ad8, true, nil
	}
	// 0x100235ad8: and w9, w4, w15, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[15])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235adc {
		return 0x100235adc, true, nil
	}
	// 0x100235adc: eor w9, w16, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235ae0 {
		return 0x100235ae0, true, nil
	}
	// 0x100235ae0: and w10, w4, w10, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[10])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235ae4 {
		return 0x100235ae4, true, nil
	}
	// 0x100235ae4: eor w3, w25, w10
	s.x[3] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235ae8 {
		return 0x100235ae8, true, nil
	}
	// 0x100235ae8: ubfx w10, w3, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235aec {
		return 0x100235aec, true, nil
	}
	// 0x100235aec: eor w15, w10, #6
	s.x[15] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(0x6))))
	if runtime != nil && runtime.stopPC == 0x100235af0 {
		return 0x100235af0, true, nil
	}
	// 0x100235af0: ldr w15, [x26, w15, uxtw #2]
	value901, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value901))
	if runtime != nil && runtime.stopPC == 0x100235af4 {
		return 0x100235af4, true, nil
	}
	// 0x100235af4: and w17, w9, #0xff
	s.x[17] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235af8 {
		return 0x100235af8, true, nil
	}
	// 0x100235af8: mov w28, #0x88
	s.x[28] = uint64(uint32(uint64(0x88)))
	if runtime != nil && runtime.stopPC == 0x100235afc {
		return 0x100235afc, true, nil
	}
	// 0x100235afc: eor w17, w17, w28
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x100235b00 {
		return 0x100235b00, true, nil
	}
	// 0x100235b00: ldr w17, [x2, w17, uxtw #2]
	value905, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value905))
	if runtime != nil && runtime.stopPC == 0x100235b04 {
		return 0x100235b04, true, nil
	}
	// 0x100235b04: eor w12, w12, w17
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235b08 {
		return 0x100235b08, true, nil
	}
	// 0x100235b08: eor w15, w12, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100235b0c {
		return 0x100235b0c, true, nil
	}
	// 0x100235b0c: and w12, w4, w17, lsr #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[17])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235b10 {
		return 0x100235b10, true, nil
	}
	// 0x100235b10: eor w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235b14 {
		return 0x100235b14, true, nil
	}
	// 0x100235b14: mov w10, #0xca
	s.x[10] = uint64(uint32(uint64(0xca)))
	if runtime != nil && runtime.stopPC == 0x100235b18 {
		return 0x100235b18, true, nil
	}
	// 0x100235b18: eor w10, w10, w16, lsr #24
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[16])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235b1c {
		return 0x100235b1c, true, nil
	}
	// 0x100235b1c: ldr w10, [x22, w10, uxtw #2]
	value912, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value912))
	if runtime != nil && runtime.stopPC == 0x100235b20 {
		return 0x100235b20, true, nil
	}
	// 0x100235b20: ubfx w16, w14, #0x10, #8
	s.x[16] = uint64(uint32((uint64(uint32(s.x[14])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235b24 {
		return 0x100235b24, true, nil
	}
	// 0x100235b24: eor w17, w16, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235b28 {
		return 0x100235b28, true, nil
	}
	// 0x100235b28: sub w17, w16, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235b2c {
		return 0x100235b2c, true, nil
	}
	// 0x100235b2c: eor w16, w16, w0
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235b30 {
		return 0x100235b30, true, nil
	}
	// 0x100235b30: ldr w16, [x24, w16, uxtw #2]
	value917, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value917))
	if runtime != nil && runtime.stopPC == 0x100235b34 {
		return 0x100235b34, true, nil
	}
	// 0x100235b34: ldr w1, [sp, #0x318]
	value918, err := s.load(addOffset(s.sp, 792), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value918))
	if runtime != nil && runtime.stopPC == 0x100235b38 {
		return 0x100235b38, true, nil
	}
	// 0x100235b38: eor w17, w1, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235b3c {
		return 0x100235b3c, true, nil
	}
	// 0x100235b3c: ubfx w1, w8, #8, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[8])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235b40 {
		return 0x100235b40, true, nil
	}
	// 0x100235b40: eor w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235b44 {
		return 0x100235b44, true, nil
	}
	// 0x100235b44: eor w10, w17, w10, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[10])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235b48 {
		return 0x100235b48, true, nil
	}
	// 0x100235b48: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235b4c {
		return 0x100235b4c, true, nil
	}
	// 0x100235b4c: eor w16, w1, w13
	s.x[16] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235b50 {
		return 0x100235b50, true, nil
	}
	// 0x100235b50: ldr w16, [x26, w16, uxtw #2]
	value925, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value925))
	if runtime != nil && runtime.stopPC == 0x100235b54 {
		return 0x100235b54, true, nil
	}
	// 0x100235b54: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235b58 {
		return 0x100235b58, true, nil
	}
	// 0x100235b58: lsr w11, w11, #8
	s.x[11] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x8)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235b5c {
		return 0x100235b5c, true, nil
	}
	// 0x100235b5c: eor w11, w14, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235b60 {
		return 0x100235b60, true, nil
	}
	// 0x100235b60: ubfx w14, w14, #8, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235b64 {
		return 0x100235b64, true, nil
	}
	// 0x100235b64: ubfx w16, w9, #0x10, #8
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235b68 {
		return 0x100235b68, true, nil
	}
	// 0x100235b68: mov w17, #0xd6
	s.x[17] = uint64(uint32(uint64(0xd6)))
	if runtime != nil && runtime.stopPC == 0x100235b6c {
		return 0x100235b6c, true, nil
	}
	// 0x100235b6c: eor w17, w16, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235b70 {
		return 0x100235b70, true, nil
	}
	// 0x100235b70: mov w1, #0x7a83
	s.x[1] = uint64(uint32(uint64(0x7a83)))
	if runtime != nil && runtime.stopPC == 0x100235b74 {
		return 0x100235b74, true, nil
	}
	// 0x100235b74: movk w1, #0xbc9e, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xbc9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235b78 {
		return 0x100235b78, true, nil
	}
	// 0x100235b78: eor w1, w16, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235b7c {
		return 0x100235b7c, true, nil
	}
	// 0x100235b7c: ldr w5, [sp, #0x314]
	value936, err := s.load(addOffset(s.sp, 788), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[5] = uint64(uint32(value936))
	if runtime != nil && runtime.stopPC == 0x100235b80 {
		return 0x100235b80, true, nil
	}
	// 0x100235b80: eor w7, w5, w14
	s.x[7] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235b84 {
		return 0x100235b84, true, nil
	}
	// 0x100235b84: sub w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235b88 {
		return 0x100235b88, true, nil
	}
	// 0x100235b88: eor w17, w7, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235b8c {
		return 0x100235b8c, true, nil
	}
	// 0x100235b8c: mov w1, #0xa2
	s.x[1] = uint64(uint32(uint64(0xa2)))
	if runtime != nil && runtime.stopPC == 0x100235b90 {
		return 0x100235b90, true, nil
	}
	// 0x100235b90: eor w1, w1, w25, lsr #24
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ shiftRight(uint64(uint32(s.x[25])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235b94 {
		return 0x100235b94, true, nil
	}
	// 0x100235b94: ldr w1, [x22, w1, uxtw #2]
	value942, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value942))
	if runtime != nil && runtime.stopPC == 0x100235b98 {
		return 0x100235b98, true, nil
	}
	// 0x100235b98: mov w5, #0x43
	s.x[5] = uint64(uint32(uint64(0x43)))
	if runtime != nil && runtime.stopPC == 0x100235b9c {
		return 0x100235b9c, true, nil
	}
	// 0x100235b9c: eor w16, w16, w5
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100235ba0 {
		return 0x100235ba0, true, nil
	}
	// 0x100235ba0: ldr w16, [x24, w16, uxtw #2]
	value945, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value945))
	if runtime != nil && runtime.stopPC == 0x100235ba4 {
		return 0x100235ba4, true, nil
	}
	// 0x100235ba4: eor w14, w14, w13
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235ba8 {
		return 0x100235ba8, true, nil
	}
	// 0x100235ba8: ldr w14, [x26, w14, uxtw #2]
	value947, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value947))
	if runtime != nil && runtime.stopPC == 0x100235bac {
		return 0x100235bac, true, nil
	}
	// 0x100235bac: eor w17, w17, w1, ror #23
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[1])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235bb0 {
		return 0x100235bb0, true, nil
	}
	// 0x100235bb0: eor w14, w16, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235bb4 {
		return 0x100235bb4, true, nil
	}
	// 0x100235bb4: and w16, w8, #0xff
	s.x[16] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235bb8 {
		return 0x100235bb8, true, nil
	}
	// 0x100235bb8: ldr w1, [x2, w16, uxtw #2]
	value951, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value951))
	if runtime != nil && runtime.stopPC == 0x100235bbc {
		return 0x100235bbc, true, nil
	}
	// 0x100235bbc: eor w14, w14, w1
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235bc0 {
		return 0x100235bc0, true, nil
	}
	// 0x100235bc0: eor w16, w17, w14
	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235bc4 {
		return 0x100235bc4, true, nil
	}
	// 0x100235bc4: lsr w8, w8, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235bc8 {
		return 0x100235bc8, true, nil
	}
	// 0x100235bc8: ldr w8, [x22, w8, uxtw #2]
	value955, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value955))
	if runtime != nil && runtime.stopPC == 0x100235bcc {
		return 0x100235bcc, true, nil
	}
	// 0x100235bcc: ubfx w9, w9, #8, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235bd0 {
		return 0x100235bd0, true, nil
	}
	// 0x100235bd0: mov w14, #0xaa
	s.x[14] = uint64(uint32(uint64(0xaa)))
	if runtime != nil && runtime.stopPC == 0x100235bd4 {
		return 0x100235bd4, true, nil
	}
	// 0x100235bd4: eor w14, w9, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235bd8 {
		return 0x100235bd8, true, nil
	}
	// 0x100235bd8: ldr w14, [x26, w14, uxtw #2]
	value959, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value959))
	if runtime != nil && runtime.stopPC == 0x100235bdc {
		return 0x100235bdc, true, nil
	}
	// 0x100235bdc: ldr w17, [sp, #0x3b4]
	value960, err := s.load(addOffset(s.sp, 948), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value960))
	if runtime != nil && runtime.stopPC == 0x100235be0 {
		return 0x100235be0, true, nil
	}
	// 0x100235be0: eor w8, w17, w8, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235be4 {
		return 0x100235be4, true, nil
	}
	// 0x100235be4: eor w8, w8, w14
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235be8 {
		return 0x100235be8, true, nil
	}
	// 0x100235be8: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235bec {
		return 0x100235bec, true, nil
	}
	// 0x100235bec: ubfx w9, w3, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[3])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235bf0 {
		return 0x100235bf0, true, nil
	}
	// 0x100235bf0: mov w5, #0xed
	s.x[5] = uint64(uint32(uint64(0xed)))
	if runtime != nil && runtime.stopPC == 0x100235bf4 {
		return 0x100235bf4, true, nil
	}
	// 0x100235bf4: eor w14, w9, w5
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100235bf8 {
		return 0x100235bf8, true, nil
	}
	// 0x100235bf8: ldr w14, [x24, w14, uxtw #2]
	value967, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value967))
	if runtime != nil && runtime.stopPC == 0x100235bfc {
		return 0x100235bfc, true, nil
	}
	// 0x100235bfc: and w11, w11, #0xff
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235c00 {
		return 0x100235c00, true, nil
	}
	// 0x100235c00: ldr w11, [x2, w11, uxtw #2]
	value969, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value969))
	if runtime != nil && runtime.stopPC == 0x100235c04 {
		return 0x100235c04, true, nil
	}
	// 0x100235c04: eor w8, w8, w14
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235c08 {
		return 0x100235c08, true, nil
	}
	// 0x100235c08: and w14, w4, w11, lsr #4
	s.x[14] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[11])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235c0c {
		return 0x100235c0c, true, nil
	}
	// 0x100235c0c: eor w11, w11, w14
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235c10 {
		return 0x100235c10, true, nil
	}
	// 0x100235c10: eor w14, w9, #0x78
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(0x78))))
	if runtime != nil && runtime.stopPC == 0x100235c14 {
		return 0x100235c14, true, nil
	}
	// 0x100235c14: mov w17, #0x7a2d
	s.x[17] = uint64(uint32(uint64(0x7a2d)))
	if runtime != nil && runtime.stopPC == 0x100235c18 {
		return 0x100235c18, true, nil
	}
	// 0x100235c18: movk w17, #0xbc9e, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xbc9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235c1c {
		return 0x100235c1c, true, nil
	}
	// 0x100235c1c: eor w9, w9, w17
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235c20 {
		return 0x100235c20, true, nil
	}
	// 0x100235c20: sub w9, w14, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) - uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235c24 {
		return 0x100235c24, true, nil
	}
	// 0x100235c24: eor w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235c28 {
		return 0x100235c28, true, nil
	}
	// 0x100235c28: and w11, w3, #0xff
	s.x[11] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235c2c {
		return 0x100235c2c, true, nil
	}
	// 0x100235c2c: eor w11, w11, w5
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100235c30 {
		return 0x100235c30, true, nil
	}
	// 0x100235c30: ldr w11, [x2, w11, uxtw #2]
	value981, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value981))
	if runtime != nil && runtime.stopPC == 0x100235c34 {
		return 0x100235c34, true, nil
	}
	// 0x100235c34: mov w14, #0xe168
	s.x[14] = uint64(uint32(uint64(0xe168)))
	if runtime != nil && runtime.stopPC == 0x100235c38 {
		return 0x100235c38, true, nil
	}
	// 0x100235c38: movk w14, #0x2cec, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x2cec, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235c3c {
		return 0x100235c3c, true, nil
	}
	// 0x100235c3c: eor w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235c40 {
		return 0x100235c40, true, nil
	}
	// 0x100235c40: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235c44 {
		return 0x100235c44, true, nil
	}
	// 0x100235c44: eor w19, w10, w11
	s.x[19] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235c48 {
		return 0x100235c48, true, nil
	}
	// 0x100235c48: and w9, w4, w11, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[11])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235c4c {
		return 0x100235c4c, true, nil
	}
	// 0x100235c4c: mov w10, #0x7275
	s.x[10] = uint64(uint32(uint64(0x7275)))
	if runtime != nil && runtime.stopPC == 0x100235c50 {
		return 0x100235c50, true, nil
	}
	// 0x100235c50: movk w10, #0x499, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x499, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235c54 {
		return 0x100235c54, true, nil
	}
	// 0x100235c54: eor w10, w19, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235c58 {
		return 0x100235c58, true, nil
	}
	// 0x100235c58: eor w3, w10, w9
	s.x[3] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235c5c {
		return 0x100235c5c, true, nil
	}
	// 0x100235c5c: and w9, w4, w1, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[1])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235c60 {
		return 0x100235c60, true, nil
	}
	// 0x100235c60: mov w10, #0xf6f2
	s.x[10] = uint64(uint32(uint64(0xf6f2)))
	if runtime != nil && runtime.stopPC == 0x100235c64 {
		return 0x100235c64, true, nil
	}
	// 0x100235c64: movk w10, #0x8e7c, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x8e7c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235c68 {
		return 0x100235c68, true, nil
	}
	// 0x100235c68: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235c6c {
		return 0x100235c6c, true, nil
	}
	// 0x100235c6c: eor w11, w16, w9
	s.x[11] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235c70 {
		return 0x100235c70, true, nil
	}
	// 0x100235c70: lsr w9, w8, #0x18
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235c74 {
		return 0x100235c74, true, nil
	}
	// 0x100235c74: ldr w9, [x22, w9, uxtw #2]
	value998, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value998))
	if runtime != nil && runtime.stopPC == 0x100235c78 {
		return 0x100235c78, true, nil
	}
	// 0x100235c78: ubfx w10, w11, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235c7c {
		return 0x100235c7c, true, nil
	}
	// 0x100235c7c: eor w14, w10, w0
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235c80 {
		return 0x100235c80, true, nil
	}
	// 0x100235c80: ldr w14, [x24, w14, uxtw #2]
	value1001, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1001))
	if runtime != nil && runtime.stopPC == 0x100235c84 {
		return 0x100235c84, true, nil
	}
	// 0x100235c84: ldr w17, [sp, #0x3c0]
	value1002, err := s.load(addOffset(s.sp, 960), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1002))
	if runtime != nil && runtime.stopPC == 0x100235c88 {
		return 0x100235c88, true, nil
	}
	// 0x100235c88: eor w9, w17, w9, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[9])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235c8c {
		return 0x100235c8c, true, nil
	}
	// 0x100235c8c: eor w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235c90 {
		return 0x100235c90, true, nil
	}
	// 0x100235c90: mov w14, #0xe801
	s.x[14] = uint64(uint32(uint64(0xe801)))
	if runtime != nil && runtime.stopPC == 0x100235c94 {
		return 0x100235c94, true, nil
	}
	// 0x100235c94: movk w14, #0xd50, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xd50, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235c98 {
		return 0x100235c98, true, nil
	}
	// 0x100235c98: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235c9c {
		return 0x100235c9c, true, nil
	}
	// 0x100235c9c: eor w12, w15, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235ca0 {
		return 0x100235ca0, true, nil
	}
	// 0x100235ca0: ubfx w14, w3, #8, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[3])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235ca4 {
		return 0x100235ca4, true, nil
	}
	// 0x100235ca4: eor w17, w14, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235ca8 {
		return 0x100235ca8, true, nil
	}
	// 0x100235ca8: ldr w17, [x26, w17, uxtw #2]
	value1011, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1011))
	if runtime != nil && runtime.stopPC == 0x100235cac {
		return 0x100235cac, true, nil
	}
	// 0x100235cac: and w1, w12, #0xff
	s.x[1] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235cb0 {
		return 0x100235cb0, true, nil
	}
	// 0x100235cb0: ldr w1, [x2, w1, uxtw #2]
	value1013, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value1013))
	if runtime != nil && runtime.stopPC == 0x100235cb4 {
		return 0x100235cb4, true, nil
	}
	// 0x100235cb4: eor w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235cb8 {
		return 0x100235cb8, true, nil
	}
	// 0x100235cb8: eor w9, w9, w17
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235cbc {
		return 0x100235cbc, true, nil
	}
	// 0x100235cbc: and w17, w4, w1, lsr #4
	s.x[17] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[1])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235cc0 {
		return 0x100235cc0, true, nil
	}
	// 0x100235cc0: eor w14, w17, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235cc4 {
		return 0x100235cc4, true, nil
	}
	// 0x100235cc4: eor w17, w10, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235cc8 {
		return 0x100235cc8, true, nil
	}
	// 0x100235cc8: sub w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235ccc {
		return 0x100235ccc, true, nil
	}
	// 0x100235ccc: eor w10, w14, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235cd0 {
		return 0x100235cd0, true, nil
	}
	// 0x100235cd0: mov w14, #0xd0
	s.x[14] = uint64(uint32(uint64(0xd0)))
	if runtime != nil && runtime.stopPC == 0x100235cd4 {
		return 0x100235cd4, true, nil
	}
	// 0x100235cd4: eor w14, w14, w15, lsr #24
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ shiftRight(uint64(uint32(s.x[15])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235cd8 {
		return 0x100235cd8, true, nil
	}
	// 0x100235cd8: ldr w14, [x22, w14, uxtw #2]
	value1023, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1023))
	if runtime != nil && runtime.stopPC == 0x100235cdc {
		return 0x100235cdc, true, nil
	}
	// 0x100235cdc: eor w15, w9, w10
	s.x[15] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235ce0 {
		return 0x100235ce0, true, nil
	}
	// 0x100235ce0: ubfx w9, w8, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235ce4 {
		return 0x100235ce4, true, nil
	}
	// 0x100235ce4: eor w10, w9, w27
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235ce8 {
		return 0x100235ce8, true, nil
	}
	// 0x100235ce8: sub w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235cec {
		return 0x100235cec, true, nil
	}
	// 0x100235cec: ldr w17, [sp, #0x308]
	value1028, err := s.load(addOffset(s.sp, 776), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1028))
	if runtime != nil && runtime.stopPC == 0x100235cf0 {
		return 0x100235cf0, true, nil
	}
	// 0x100235cf0: eor w10, w17, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235cf4 {
		return 0x100235cf4, true, nil
	}
	// 0x100235cf4: eor w10, w10, w14, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ rotateRight(uint64(uint32(s.x[14])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235cf8 {
		return 0x100235cf8, true, nil
	}
	// 0x100235cf8: eor w9, w9, w0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235cfc {
		return 0x100235cfc, true, nil
	}
	// 0x100235cfc: ldr w9, [x24, w9, uxtw #2]
	value1032, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1032))
	if runtime != nil && runtime.stopPC == 0x100235d00 {
		return 0x100235d00, true, nil
	}
	// 0x100235d00: eor w10, w10, w9
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235d04 {
		return 0x100235d04, true, nil
	}
	// 0x100235d04: ubfx w9, w11, #8, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235d08 {
		return 0x100235d08, true, nil
	}
	// 0x100235d08: eor w9, w9, w13
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235d0c {
		return 0x100235d0c, true, nil
	}
	// 0x100235d0c: ldr w14, [x26, w9, uxtw #2]
	value1036, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1036))
	if runtime != nil && runtime.stopPC == 0x100235d10 {
		return 0x100235d10, true, nil
	}
	// 0x100235d10: and w9, w3, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235d14 {
		return 0x100235d14, true, nil
	}
	// 0x100235d14: ldr w9, [x2, w9, uxtw #2]
	value1038, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1038))
	if runtime != nil && runtime.stopPC == 0x100235d18 {
		return 0x100235d18, true, nil
	}
	// 0x100235d18: eor w14, w14, w9
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235d1c {
		return 0x100235d1c, true, nil
	}
	// 0x100235d1c: eor w14, w10, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235d20 {
		return 0x100235d20, true, nil
	}
	// 0x100235d20: ubfx w10, w12, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235d24 {
		return 0x100235d24, true, nil
	}
	// 0x100235d24: eor w17, w10, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235d28 {
		return 0x100235d28, true, nil
	}
	// 0x100235d28: ubfx w1, w8, #8, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[8])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235d2c {
		return 0x100235d2c, true, nil
	}
	// 0x100235d2c: ldr w5, [sp, #0x3d8]
	value1044, err := s.load(addOffset(s.sp, 984), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[5] = uint64(uint32(value1044))
	if runtime != nil && runtime.stopPC == 0x100235d30 {
		return 0x100235d30, true, nil
	}
	// 0x100235d30: eor w7, w5, w1
	s.x[7] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235d34 {
		return 0x100235d34, true, nil
	}
	// 0x100235d34: sub w17, w10, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235d38 {
		return 0x100235d38, true, nil
	}
	// 0x100235d38: eor w17, w7, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235d3c {
		return 0x100235d3c, true, nil
	}
	// 0x100235d3c: eor w10, w10, w0
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235d40 {
		return 0x100235d40, true, nil
	}
	// 0x100235d40: ldr w10, [x24, w10, uxtw #2]
	value1049, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1049))
	if runtime != nil && runtime.stopPC == 0x100235d44 {
		return 0x100235d44, true, nil
	}
	// 0x100235d44: eor w10, w17, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235d48 {
		return 0x100235d48, true, nil
	}
	// 0x100235d48: eor w17, w6, w19, lsr #24
	s.x[17] = uint64(uint32((uint64(uint32(s.x[6])) ^ shiftRight(uint64(uint32(s.x[19])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235d4c {
		return 0x100235d4c, true, nil
	}
	// 0x100235d4c: ldr w17, [x22, w17, uxtw #2]
	value1052, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1052))
	if runtime != nil && runtime.stopPC == 0x100235d50 {
		return 0x100235d50, true, nil
	}
	// 0x100235d50: eor w10, w10, w17, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ rotateRight(uint64(uint32(s.x[17])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235d54 {
		return 0x100235d54, true, nil
	}
	// 0x100235d54: eor w17, w1, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235d58 {
		return 0x100235d58, true, nil
	}
	// 0x100235d58: ldr w17, [x26, w17, uxtw #2]
	value1055, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1055))
	if runtime != nil && runtime.stopPC == 0x100235d5c {
		return 0x100235d5c, true, nil
	}
	// 0x100235d5c: and w1, w11, #0xff
	s.x[1] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235d60 {
		return 0x100235d60, true, nil
	}
	// 0x100235d60: ldr w1, [x2, w1, uxtw #2]
	value1057, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value1057))
	if runtime != nil && runtime.stopPC == 0x100235d64 {
		return 0x100235d64, true, nil
	}
	// 0x100235d64: eor w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235d68 {
		return 0x100235d68, true, nil
	}
	// 0x100235d68: mov w5, #0xddc4
	s.x[5] = uint64(uint32(uint64(0xddc4)))
	if runtime != nil && runtime.stopPC == 0x100235d6c {
		return 0x100235d6c, true, nil
	}
	// 0x100235d6c: movk w5, #0x2c11, lsl #16
	s.x[5] = uint64(uint32(insertBits(uint64(uint32(s.x[5])), 0x2c11, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235d70 {
		return 0x100235d70, true, nil
	}
	// 0x100235d70: eor w17, w17, w5
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100235d74 {
		return 0x100235d74, true, nil
	}
	// 0x100235d74: eor w25, w10, w17
	s.x[25] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235d78 {
		return 0x100235d78, true, nil
	}
	// 0x100235d78: ubfx w10, w3, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235d7c {
		return 0x100235d7c, true, nil
	}
	// 0x100235d7c: ubfx w12, w12, #8, #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235d80 {
		return 0x100235d80, true, nil
	}
	// 0x100235d80: eor w17, w10, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235d84 {
		return 0x100235d84, true, nil
	}
	// 0x100235d84: ldr w3, [sp, #0x304]
	value1066, err := s.load(addOffset(s.sp, 772), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1066))
	if runtime != nil && runtime.stopPC == 0x100235d88 {
		return 0x100235d88, true, nil
	}
	// 0x100235d88: eor w3, w3, w12
	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235d8c {
		return 0x100235d8c, true, nil
	}
	// 0x100235d8c: sub w17, w10, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235d90 {
		return 0x100235d90, true, nil
	}
	// 0x100235d90: eor w17, w3, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235d94 {
		return 0x100235d94, true, nil
	}
	// 0x100235d94: mov w3, #0x53
	s.x[3] = uint64(uint32(uint64(0x53)))
	if runtime != nil && runtime.stopPC == 0x100235d98 {
		return 0x100235d98, true, nil
	}
	// 0x100235d98: eor w16, w3, w16, lsr #24
	s.x[16] = uint64(uint32((uint64(uint32(s.x[3])) ^ shiftRight(uint64(uint32(s.x[16])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x100235d9c {
		return 0x100235d9c, true, nil
	}
	// 0x100235d9c: ldr w16, [x22, w16, uxtw #2]
	value1072, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value1072))
	if runtime != nil && runtime.stopPC == 0x100235da0 {
		return 0x100235da0, true, nil
	}
	// 0x100235da0: eor w10, w10, w0
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235da4 {
		return 0x100235da4, true, nil
	}
	// 0x100235da4: ldr w10, [x24, w10, uxtw #2]
	value1074, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1074))
	if runtime != nil && runtime.stopPC == 0x100235da8 {
		return 0x100235da8, true, nil
	}
	// 0x100235da8: eor w16, w17, w16, ror #23
	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[16])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235dac {
		return 0x100235dac, true, nil
	}
	// 0x100235dac: eor w12, w12, w13
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235db0 {
		return 0x100235db0, true, nil
	}
	// 0x100235db0: ldr w12, [x26, w12, uxtw #2]
	value1077, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value1077))
	if runtime != nil && runtime.stopPC == 0x100235db4 {
		return 0x100235db4, true, nil
	}
	// 0x100235db4: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235db8 {
		return 0x100235db8, true, nil
	}
	// 0x100235db8: eor w10, w16, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235dbc {
		return 0x100235dbc, true, nil
	}
	// 0x100235dbc: and w8, w8, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235dc0 {
		return 0x100235dc0, true, nil
	}
	// 0x100235dc0: ldr w8, [x2, w8, uxtw #2]
	value1081, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1081))
	if runtime != nil && runtime.stopPC == 0x100235dc4 {
		return 0x100235dc4, true, nil
	}
	// 0x100235dc4: mov w12, #0xb409
	s.x[12] = uint64(uint32(uint64(0xb409)))
	if runtime != nil && runtime.stopPC == 0x100235dc8 {
		return 0x100235dc8, true, nil
	}
	// 0x100235dc8: movk w12, #0x1919, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x1919, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235dcc {
		return 0x100235dcc, true, nil
	}
	// 0x100235dcc: eor w12, w8, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235dd0 {
		return 0x100235dd0, true, nil
	}
	// 0x100235dd0: eor w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235dd4 {
		return 0x100235dd4, true, nil
	}
	// 0x100235dd4: and w10, w4, w1, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[1])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235dd8 {
		return 0x100235dd8, true, nil
	}
	// 0x100235dd8: and w8, w4, w8, lsr #4
	s.x[8] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[8])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235ddc {
		return 0x100235ddc, true, nil
	}
	// 0x100235ddc: eor w19, w12, w8
	s.x[19] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235de0 {
		return 0x100235de0, true, nil
	}
	// 0x100235de0: ubfx w8, w19, #0x10, #8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[19])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235de4 {
		return 0x100235de4, true, nil
	}
	// 0x100235de4: eor w16, w8, w0
	s.x[16] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235de8 {
		return 0x100235de8, true, nil
	}
	// 0x100235de8: ldr w17, [x24, w16, uxtw #2]
	value1091, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1091))
	if runtime != nil && runtime.stopPC == 0x100235dec {
		return 0x100235dec, true, nil
	}
	// 0x100235dec: eor w16, w25, w10
	s.x[16] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235df0 {
		return 0x100235df0, true, nil
	}
	// 0x100235df0: ubfx w1, w16, #8, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235df4 {
		return 0x100235df4, true, nil
	}
	// 0x100235df4: eor w10, w1, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235df8 {
		return 0x100235df8, true, nil
	}
	// 0x100235df8: ldr w10, [x26, w10, uxtw #2]
	value1095, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1095))
	if runtime != nil && runtime.stopPC == 0x100235dfc {
		return 0x100235dfc, true, nil
	}
	// 0x100235dfc: eor w10, w17, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235e00 {
		return 0x100235e00, true, nil
	}
	// 0x100235e00: lsr w17, w15, #0x18
	s.x[17] = uint64(uint32(shiftRight(uint64(uint32(s.x[15])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235e04 {
		return 0x100235e04, true, nil
	}
	// 0x100235e04: eor w17, w17, #0x3e
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(0x3e))))
	if runtime != nil && runtime.stopPC == 0x100235e08 {
		return 0x100235e08, true, nil
	}
	// 0x100235e08: ldr w17, [x22, w17, uxtw #2]
	value1099, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1099))
	if runtime != nil && runtime.stopPC == 0x100235e0c {
		return 0x100235e0c, true, nil
	}
	// 0x100235e0c: ldr w3, [sp, #0x300]
	value1100, err := s.load(addOffset(s.sp, 768), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1100))
	if runtime != nil && runtime.stopPC == 0x100235e10 {
		return 0x100235e10, true, nil
	}
	// 0x100235e10: eor w17, w3, w17, ror #23
	s.x[17] = uint64(uint32((uint64(uint32(s.x[3])) ^ rotateRight(uint64(uint32(s.x[17])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235e14 {
		return 0x100235e14, true, nil
	}
	// 0x100235e14: eor w17, w17, w10
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235e18 {
		return 0x100235e18, true, nil
	}
	// 0x100235e18: lsr w11, w11, #8
	s.x[11] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x8)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235e1c {
		return 0x100235e1c, true, nil
	}
	// 0x100235e1c: and w9, w4, w9, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[9])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235e20 {
		return 0x100235e20, true, nil
	}
	// 0x100235e20: mov w10, #0x5ce0
	s.x[10] = uint64(uint32(uint64(0x5ce0)))
	if runtime != nil && runtime.stopPC == 0x100235e24 {
		return 0x100235e24, true, nil
	}
	// 0x100235e24: movk w10, #0xfd27, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xfd27, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235e28 {
		return 0x100235e28, true, nil
	}
	// 0x100235e28: eor w3, w14, w10
	s.x[3] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235e2c {
		return 0x100235e2c, true, nil
	}
	// 0x100235e2c: eor w10, w3, w9
	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235e30 {
		return 0x100235e30, true, nil
	}
	// 0x100235e30: eor w9, w10, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235e34 {
		return 0x100235e34, true, nil
	}
	// 0x100235e34: and w9, w9, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235e38 {
		return 0x100235e38, true, nil
	}
	// 0x100235e38: ldr w9, [x2, w9, uxtw #2]
	value1111, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1111))
	if runtime != nil && runtime.stopPC == 0x100235e3c {
		return 0x100235e3c, true, nil
	}
	// 0x100235e3c: eor w11, w9, w1
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235e40 {
		return 0x100235e40, true, nil
	}
	// 0x100235e40: eor w11, w17, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100235e44 {
		return 0x100235e44, true, nil
	}
	// 0x100235e44: and w9, w4, w9, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[9])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235e48 {
		return 0x100235e48, true, nil
	}
	// 0x100235e48: eor w14, w8, w27
	s.x[14] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235e4c {
		return 0x100235e4c, true, nil
	}
	// 0x100235e4c: sub w8, w8, w14
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235e50 {
		return 0x100235e50, true, nil
	}
	// 0x100235e50: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235e54 {
		return 0x100235e54, true, nil
	}
	// 0x100235e54: eor w11, w11, w8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235e58 {
		return 0x100235e58, true, nil
	}
	// 0x100235e58: ubfx w7, w15, #0x10, #8
	s.x[7] = uint64(uint32((uint64(uint32(s.x[15])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235e5c {
		return 0x100235e5c, true, nil
	}
	// 0x100235e5c: mov w6, #0x90
	s.x[6] = uint64(uint32(uint64(0x90)))
	if runtime != nil && runtime.stopPC == 0x100235e60 {
		return 0x100235e60, true, nil
	}
	// 0x100235e60: eor w8, w7, w6
	s.x[8] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100235e64 {
		return 0x100235e64, true, nil
	}
	// 0x100235e64: mov w9, #0x7ac5
	s.x[9] = uint64(uint32(uint64(0x7ac5)))
	if runtime != nil && runtime.stopPC == 0x100235e68 {
		return 0x100235e68, true, nil
	}
	// 0x100235e68: movk w9, #0xbc9e, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xbc9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235e6c {
		return 0x100235e6c, true, nil
	}
	// 0x100235e6c: eor w9, w7, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235e70 {
		return 0x100235e70, true, nil
	}
	// 0x100235e70: sub w9, w8, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235e74 {
		return 0x100235e74, true, nil
	}
	// 0x100235e74: ubfx w8, w10, #0x10, #8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235e78 {
		return 0x100235e78, true, nil
	}
	// 0x100235e78: eor w14, w8, w27
	s.x[14] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235e7c {
		return 0x100235e7c, true, nil
	}
	// 0x100235e7c: ubfx w17, w15, #8, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235e80 {
		return 0x100235e80, true, nil
	}
	// 0x100235e80: ldr w1, [sp, #0x3e0]
	value1129, err := s.load(addOffset(s.sp, 992), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value1129))
	if runtime != nil && runtime.stopPC == 0x100235e84 {
		return 0x100235e84, true, nil
	}
	// 0x100235e84: eor w1, w17, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235e88 {
		return 0x100235e88, true, nil
	}
	// 0x100235e88: sub w14, w8, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235e8c {
		return 0x100235e8c, true, nil
	}
	// 0x100235e8c: eor w14, w1, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235e90 {
		return 0x100235e90, true, nil
	}
	// 0x100235e90: eor w8, w8, w0
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235e94 {
		return 0x100235e94, true, nil
	}
	// 0x100235e94: ldr w8, [x24, w8, uxtw #2]
	value1134, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1134))
	if runtime != nil && runtime.stopPC == 0x100235e98 {
		return 0x100235e98, true, nil
	}
	// 0x100235e98: eor w8, w14, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100235e9c {
		return 0x100235e9c, true, nil
	}
	// 0x100235e9c: lsr w14, w25, #0x18
	s.x[14] = uint64(uint32(shiftRight(uint64(uint32(s.x[25])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235ea0 {
		return 0x100235ea0, true, nil
	}
	// 0x100235ea0: ldr w14, [x22, w14, uxtw #2]
	value1137, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1137))
	if runtime != nil && runtime.stopPC == 0x100235ea4 {
		return 0x100235ea4, true, nil
	}
	// 0x100235ea4: eor w8, w8, w14, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ rotateRight(uint64(uint32(s.x[14])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235ea8 {
		return 0x100235ea8, true, nil
	}
	// 0x100235ea8: eor w14, w17, w6
	s.x[14] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100235eac {
		return 0x100235eac, true, nil
	}
	// 0x100235eac: ldr w17, [x26, w14, uxtw #2]
	value1140, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1140))
	if runtime != nil && runtime.stopPC == 0x100235eb0 {
		return 0x100235eb0, true, nil
	}
	// 0x100235eb0: and w14, w19, #0xff
	s.x[14] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235eb4 {
		return 0x100235eb4, true, nil
	}
	// 0x100235eb4: ldr w14, [x2, w14, uxtw #2]
	value1142, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1142))
	if runtime != nil && runtime.stopPC == 0x100235eb8 {
		return 0x100235eb8, true, nil
	}
	// 0x100235eb8: eor w17, w17, w14
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100235ebc {
		return 0x100235ebc, true, nil
	}
	// 0x100235ebc: eor w8, w8, w17
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235ec0 {
		return 0x100235ec0, true, nil
	}
	// 0x100235ec0: ubfx w10, w10, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235ec4 {
		return 0x100235ec4, true, nil
	}
	// 0x100235ec4: ubfx w17, w16, #0x10, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235ec8 {
		return 0x100235ec8, true, nil
	}
	// 0x100235ec8: eor w1, w17, w27
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235ecc {
		return 0x100235ecc, true, nil
	}
	// 0x100235ecc: ldr w5, [sp, #0x2fc]
	value1148, err := s.load(addOffset(s.sp, 764), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[5] = uint64(uint32(value1148))
	if runtime != nil && runtime.stopPC == 0x100235ed0 {
		return 0x100235ed0, true, nil
	}
	// 0x100235ed0: eor w25, w5, w10
	s.x[25] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235ed4 {
		return 0x100235ed4, true, nil
	}
	// 0x100235ed4: sub w1, w17, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235ed8 {
		return 0x100235ed8, true, nil
	}
	// 0x100235ed8: eor w1, w25, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[25])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235edc {
		return 0x100235edc, true, nil
	}
	// 0x100235edc: lsr w12, w12, #0x18
	s.x[12] = uint64(uint32(shiftRight(uint64(uint32(s.x[12])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235ee0 {
		return 0x100235ee0, true, nil
	}
	// 0x100235ee0: ldr w12, [x22, w12, uxtw #2]
	value1153, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value1153))
	if runtime != nil && runtime.stopPC == 0x100235ee4 {
		return 0x100235ee4, true, nil
	}
	// 0x100235ee4: eor w12, w1, w12, ror #23
	s.x[12] = uint64(uint32((uint64(uint32(s.x[1])) ^ rotateRight(uint64(uint32(s.x[12])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235ee8 {
		return 0x100235ee8, true, nil
	}
	// 0x100235ee8: eor w17, w17, w0
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235eec {
		return 0x100235eec, true, nil
	}
	// 0x100235eec: ldr w17, [x24, w17, uxtw #2]
	value1156, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1156))
	if runtime != nil && runtime.stopPC == 0x100235ef0 {
		return 0x100235ef0, true, nil
	}
	// 0x100235ef0: eor w12, w12, w17
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235ef4 {
		return 0x100235ef4, true, nil
	}
	// 0x100235ef4: eor w10, w10, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235ef8 {
		return 0x100235ef8, true, nil
	}
	// 0x100235ef8: ldr w10, [x26, w10, uxtw #2]
	value1159, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1159))
	if runtime != nil && runtime.stopPC == 0x100235efc {
		return 0x100235efc, true, nil
	}
	// 0x100235efc: and w15, w15, #0xff
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235f00 {
		return 0x100235f00, true, nil
	}
	// 0x100235f00: mov w17, #0x3b
	s.x[17] = uint64(uint32(uint64(0x3b)))
	if runtime != nil && runtime.stopPC == 0x100235f04 {
		return 0x100235f04, true, nil
	}
	// 0x100235f04: eor w15, w15, w17
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f08 {
		return 0x100235f08, true, nil
	}
	// 0x100235f08: ldr w15, [x2, w15, uxtw #2]
	value1163, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value1163))
	if runtime != nil && runtime.stopPC == 0x100235f0c {
		return 0x100235f0c, true, nil
	}
	// 0x100235f0c: eor w10, w10, w15
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100235f10 {
		return 0x100235f10, true, nil
	}
	// 0x100235f10: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235f14 {
		return 0x100235f14, true, nil
	}
	// 0x100235f14: and w12, w4, w15, lsr #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[15])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235f18 {
		return 0x100235f18, true, nil
	}
	// 0x100235f18: eor w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100235f1c {
		return 0x100235f1c, true, nil
	}
	// 0x100235f1c: mov w10, #5
	s.x[10] = uint64(uint32(uint64(0x5)))
	if runtime != nil && runtime.stopPC == 0x100235f20 {
		return 0x100235f20, true, nil
	}
	// 0x100235f20: eor w10, w7, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235f24 {
		return 0x100235f24, true, nil
	}
	// 0x100235f24: ldr w10, [x24, w10, uxtw #2]
	value1170, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1170))
	if runtime != nil && runtime.stopPC == 0x100235f28 {
		return 0x100235f28, true, nil
	}
	// 0x100235f28: ldr w15, [sp, #0x2f8]
	value1171, err := s.load(addOffset(s.sp, 760), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value1171))
	if runtime != nil && runtime.stopPC == 0x100235f2c {
		return 0x100235f2c, true, nil
	}
	// 0x100235f2c: eor w9, w15, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100235f30 {
		return 0x100235f30, true, nil
	}
	// 0x100235f30: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235f34 {
		return 0x100235f34, true, nil
	}
	// 0x100235f34: lsr w10, w3, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[3])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235f38 {
		return 0x100235f38, true, nil
	}
	// 0x100235f38: ldr w10, [x22, w10, uxtw #2]
	value1175, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1175))
	if runtime != nil && runtime.stopPC == 0x100235f3c {
		return 0x100235f3c, true, nil
	}
	// 0x100235f3c: ubfx w15, w19, #8, #8
	s.x[15] = uint64(uint32((uint64(uint32(s.x[19])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235f40 {
		return 0x100235f40, true, nil
	}
	// 0x100235f40: eor w17, w15, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100235f44 {
		return 0x100235f44, true, nil
	}
	// 0x100235f44: ldr w17, [x26, w17, uxtw #2]
	value1178, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1178))
	if runtime != nil && runtime.stopPC == 0x100235f48 {
		return 0x100235f48, true, nil
	}
	// 0x100235f48: eor w9, w9, w17
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f4c {
		return 0x100235f4c, true, nil
	}
	// 0x100235f4c: ror w10, w10, #0x17
	s.x[10] = uint64(uint32(rotateRight(uint64(uint32(s.x[10])), int(uint64(0x17)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235f50 {
		return 0x100235f50, true, nil
	}
	// 0x100235f50: mov w17, #0xb91e
	s.x[17] = uint64(uint32(uint64(0xb91e)))
	if runtime != nil && runtime.stopPC == 0x100235f54 {
		return 0x100235f54, true, nil
	}
	// 0x100235f54: movk w17, #0x6059, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x6059, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235f58 {
		return 0x100235f58, true, nil
	}
	// 0x100235f58: and w17, w10, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f5c {
		return 0x100235f5c, true, nil
	}
	// 0x100235f5c: eor w9, w9, w17
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f60 {
		return 0x100235f60, true, nil
	}
	// 0x100235f60: mov w17, #0x5ffb
	s.x[17] = uint64(uint32(uint64(0x5ffb)))
	if runtime != nil && runtime.stopPC == 0x100235f64 {
		return 0x100235f64, true, nil
	}
	// 0x100235f64: movk w17, #0x4750, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x4750, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235f68 {
		return 0x100235f68, true, nil
	}
	// 0x100235f68: eor w17, w10, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f6c {
		return 0x100235f6c, true, nil
	}
	// 0x100235f6c: ldr w1, [sp, #0x384]
	value1188, err := s.load(addOffset(s.sp, 900), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value1188))
	if runtime != nil && runtime.stopPC == 0x100235f70 {
		return 0x100235f70, true, nil
	}
	// 0x100235f70: and w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235f74 {
		return 0x100235f74, true, nil
	}
	// 0x100235f74: eor w9, w9, w17
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f78 {
		return 0x100235f78, true, nil
	}
	// 0x100235f78: mov w17, #0x9c9d
	s.x[17] = uint64(uint32(uint64(0x9c9d)))
	if runtime != nil && runtime.stopPC == 0x100235f7c {
		return 0x100235f7c, true, nil
	}
	// 0x100235f7c: movk w17, #0x25db, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x25db, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235f80 {
		return 0x100235f80, true, nil
	}
	// 0x100235f80: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f84 {
		return 0x100235f84, true, nil
	}
	// 0x100235f84: and w16, w16, #0xff
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100235f88 {
		return 0x100235f88, true, nil
	}
	// 0x100235f88: ldr w16, [x2, w16, uxtw #2]
	value1195, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value1195))
	if runtime != nil && runtime.stopPC == 0x100235f8c {
		return 0x100235f8c, true, nil
	}
	// 0x100235f8c: ldr w17, [sp, #0x380]
	value1196, err := s.load(addOffset(s.sp, 896), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1196))
	if runtime != nil && runtime.stopPC == 0x100235f90 {
		return 0x100235f90, true, nil
	}
	// 0x100235f90: and w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100235f94 {
		return 0x100235f94, true, nil
	}
	// 0x100235f94: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235f98 {
		return 0x100235f98, true, nil
	}
	// 0x100235f98: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235f9c {
		return 0x100235f9c, true, nil
	}
	// 0x100235f9c: and w10, w4, w16, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[16])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235fa0 {
		return 0x100235fa0, true, nil
	}
	// 0x100235fa0: eor w17, w15, w10
	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235fa4 {
		return 0x100235fa4, true, nil
	}
	// 0x100235fa4: mov w10, #0xc7c9
	s.x[10] = uint64(uint32(uint64(0xc7c9)))
	if runtime != nil && runtime.stopPC == 0x100235fa8 {
		return 0x100235fa8, true, nil
	}
	// 0x100235fa8: movk w10, #0x70bc, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x70bc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235fac {
		return 0x100235fac, true, nil
	}
	// 0x100235fac: eor w15, w12, w10
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235fb0 {
		return 0x100235fb0, true, nil
	}
	// 0x100235fb0: ubfx w10, w15, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[15])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100235fb4 {
		return 0x100235fb4, true, nil
	}
	// 0x100235fb4: eor w16, w10, w0
	s.x[16] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100235fb8 {
		return 0x100235fb8, true, nil
	}
	// 0x100235fb8: ldr w16, [x24, w16, uxtw #2]
	value1207, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[16])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value1207))
	if runtime != nil && runtime.stopPC == 0x100235fbc {
		return 0x100235fbc, true, nil
	}
	// 0x100235fbc: eor w1, w10, w27
	s.x[1] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100235fc0 {
		return 0x100235fc0, true, nil
	}
	// 0x100235fc0: sub w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235fc4 {
		return 0x100235fc4, true, nil
	}
	// 0x100235fc4: eor w10, w16, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235fc8 {
		return 0x100235fc8, true, nil
	}
	// 0x100235fc8: mov w16, #0x50c7
	s.x[16] = uint64(uint32(uint64(0x50c7)))
	if runtime != nil && runtime.stopPC == 0x100235fcc {
		return 0x100235fcc, true, nil
	}
	// 0x100235fcc: movk w16, #0x2a9d, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x2a9d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235fd0 {
		return 0x100235fd0, true, nil
	}
	// 0x100235fd0: eor w16, w11, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100235fd4 {
		return 0x100235fd4, true, nil
	}
	// 0x100235fd4: lsr w11, w16, #0x18
	s.x[11] = uint64(uint32(shiftRight(uint64(uint32(s.x[16])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100235fd8 {
		return 0x100235fd8, true, nil
	}
	// 0x100235fd8: ldr w11, [x22, w11, uxtw #2]
	value1215, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value1215))
	if runtime != nil && runtime.stopPC == 0x100235fdc {
		return 0x100235fdc, true, nil
	}
	// 0x100235fdc: ldr w1, [sp, #0x2f4]
	value1216, err := s.load(addOffset(s.sp, 756), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value1216))
	if runtime != nil && runtime.stopPC == 0x100235fe0 {
		return 0x100235fe0, true, nil
	}
	// 0x100235fe0: eor w11, w1, w11, ror #23
	s.x[11] = uint64(uint32((uint64(uint32(s.x[1])) ^ rotateRight(uint64(uint32(s.x[11])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100235fe4 {
		return 0x100235fe4, true, nil
	}
	// 0x100235fe4: mov w1, #0x61e7
	s.x[1] = uint64(uint32(uint64(0x61e7)))
	if runtime != nil && runtime.stopPC == 0x100235fe8 {
		return 0x100235fe8, true, nil
	}
	// 0x100235fe8: movk w1, #0x1205, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x1205, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100235fec {
		return 0x100235fec, true, nil
	}
	// 0x100235fec: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100235ff0 {
		return 0x100235ff0, true, nil
	}
	// 0x100235ff0: eor w11, w11, w10
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235ff4 {
		return 0x100235ff4, true, nil
	}
	// 0x100235ff4: and w10, w4, w14, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[14])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100235ff8 {
		return 0x100235ff8, true, nil
	}
	// 0x100235ff8: eor w19, w8, w10
	s.x[19] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100235ffc {
		return 0x100235ffc, true, nil
	}
	// 0x100235ffc: ubfx w10, w19, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[19])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236000 {
		return 0x100236000, true, nil
	}
	// 0x100236000: mov w14, #0x92
	s.x[14] = uint64(uint32(uint64(0x92)))
	if runtime != nil && runtime.stopPC == 0x100236004 {
		return 0x100236004, true, nil
	}
	// 0x100236004: eor w14, w10, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236008 {
		return 0x100236008, true, nil
	}
	// 0x100236008: ldr w14, [x26, w14, uxtw #2]
	value1227, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1227))
	if runtime != nil && runtime.stopPC == 0x10023600c {
		return 0x10023600c, true, nil
	}
	// 0x10023600c: mov w1, #0xd410
	s.x[1] = uint64(uint32(uint64(0xd410)))
	if runtime != nil && runtime.stopPC == 0x100236010 {
		return 0x100236010, true, nil
	}
	// 0x100236010: movk w1, #0x91e3, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x91e3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236014 {
		return 0x100236014, true, nil
	}
	// 0x100236014: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236018 {
		return 0x100236018, true, nil
	}
	// 0x100236018: eor w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10023601c {
		return 0x10023601c, true, nil
	}
	// 0x10023601c: eor w14, w10, w11
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236020 {
		return 0x100236020, true, nil
	}
	// 0x100236020: sub w14, w11, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) - uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236024 {
		return 0x100236024, true, nil
	}
	// 0x100236024: neg w1, w10
	s.x[1] = uint64(uint32(-(uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236028 {
		return 0x100236028, true, nil
	}
	// 0x100236028: eor w14, w1, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10023602c {
		return 0x10023602c, true, nil
	}
	// 0x10023602c: mov w1, #0xc416
	s.x[1] = uint64(uint32(uint64(0xc416)))
	if runtime != nil && runtime.stopPC == 0x100236030 {
		return 0x100236030, true, nil
	}
	// 0x100236030: movk w1, #0x56ee, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x56ee, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236034 {
		return 0x100236034, true, nil
	}
	// 0x100236034: and w1, w1, w10, lsl #1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) & shiftLeft(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236038 {
		return 0x100236038, true, nil
	}
	// 0x100236038: sub w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x10023603c {
		return 0x10023603c, true, nil
	}
	// 0x10023603c: mov w6, #0x620b
	s.x[6] = uint64(uint32(uint64(0x620b)))
	if runtime != nil && runtime.stopPC == 0x100236040 {
		return 0x100236040, true, nil
	}
	// 0x100236040: movk w6, #0xab77, lsl #16
	s.x[6] = uint64(uint32(insertBits(uint64(uint32(s.x[6])), 0xab77, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236044 {
		return 0x100236044, true, nil
	}
	// 0x100236044: add w10, w10, w6
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x100236048 {
		return 0x100236048, true, nil
	}
	// 0x100236048: eor w10, w14, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023604c {
		return 0x10023604c, true, nil
	}
	// 0x10023604c: mov w14, #0xe2b8
	s.x[14] = uint64(uint32(uint64(0xe2b8)))
	if runtime != nil && runtime.stopPC == 0x100236050 {
		return 0x100236050, true, nil
	}
	// 0x100236050: movk w14, #0x8e87, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x8e87, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236054 {
		return 0x100236054, true, nil
	}
	// 0x100236054: eor w9, w9, w14
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236058 {
		return 0x100236058, true, nil
	}
	// 0x100236058: eor w3, w9, w17
	s.x[3] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10023605c {
		return 0x10023605c, true, nil
	}
	// 0x10023605c: ubfx w14, w3, #0x10, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[3])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236060 {
		return 0x100236060, true, nil
	}
	// 0x100236060: eor w17, w14, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236064 {
		return 0x100236064, true, nil
	}
	// 0x100236064: ubfx w1, w16, #8, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236068 {
		return 0x100236068, true, nil
	}
	// 0x100236068: ldr w5, [sp, #0x2f0]
	value1251, err := s.load(addOffset(s.sp, 752), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[5] = uint64(uint32(value1251))
	if runtime != nil && runtime.stopPC == 0x10023606c {
		return 0x10023606c, true, nil
	}
	// 0x10023606c: eor w7, w5, w1
	s.x[7] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236070 {
		return 0x100236070, true, nil
	}
	// 0x100236070: sub w17, w14, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236074 {
		return 0x100236074, true, nil
	}
	// 0x100236074: eor w17, w7, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236078 {
		return 0x100236078, true, nil
	}
	// 0x100236078: eor w14, w14, w0
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x10023607c {
		return 0x10023607c, true, nil
	}
	// 0x10023607c: ldr w14, [x24, w14, uxtw #2]
	value1256, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1256))
	if runtime != nil && runtime.stopPC == 0x100236080 {
		return 0x100236080, true, nil
	}
	// 0x100236080: eor w14, w17, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236084 {
		return 0x100236084, true, nil
	}
	// 0x100236084: mov w17, #0x8c
	s.x[17] = uint64(uint32(uint64(0x8c)))
	if runtime != nil && runtime.stopPC == 0x100236088 {
		return 0x100236088, true, nil
	}
	// 0x100236088: eor w8, w17, w8, lsr #24
	s.x[8] = uint64(uint32((uint64(uint32(s.x[17])) ^ shiftRight(uint64(uint32(s.x[8])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x10023608c {
		return 0x10023608c, true, nil
	}
	// 0x10023608c: ldr w8, [x22, w8, uxtw #2]
	value1260, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1260))
	if runtime != nil && runtime.stopPC == 0x100236090 {
		return 0x100236090, true, nil
	}
	// 0x100236090: eor w17, w1, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236094 {
		return 0x100236094, true, nil
	}
	// 0x100236094: ldr w17, [x26, w17, uxtw #2]
	value1262, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1262))
	if runtime != nil && runtime.stopPC == 0x100236098 {
		return 0x100236098, true, nil
	}
	// 0x100236098: eor w8, w14, w8, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[14])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x10023609c {
		return 0x10023609c, true, nil
	}
	// 0x10023609c: and w12, w12, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002360a0 {
		return 0x1002360a0, true, nil
	}
	// 0x1002360a0: mov w14, #0xf1
	s.x[14] = uint64(uint32(uint64(0xf1)))
	if runtime != nil && runtime.stopPC == 0x1002360a4 {
		return 0x1002360a4, true, nil
	}
	// 0x1002360a4: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002360a8 {
		return 0x1002360a8, true, nil
	}
	// 0x1002360a8: ldr w12, [x2, w12, uxtw #2]
	value1267, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value1267))
	if runtime != nil && runtime.stopPC == 0x1002360ac {
		return 0x1002360ac, true, nil
	}
	// 0x1002360ac: eor w14, w17, w12
	s.x[14] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002360b0 {
		return 0x1002360b0, true, nil
	}
	// 0x1002360b0: eor w8, w8, w14
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002360b4 {
		return 0x1002360b4, true, nil
	}
	// 0x1002360b4: and w14, w3, #0xff
	s.x[14] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002360b8 {
		return 0x1002360b8, true, nil
	}
	// 0x1002360b8: ldr w14, [x2, w14, uxtw #2]
	value1271, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1271))
	if runtime != nil && runtime.stopPC == 0x1002360bc {
		return 0x1002360bc, true, nil
	}
	// 0x1002360bc: and w17, w4, w14, lsr #4
	s.x[17] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[14])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002360c0 {
		return 0x1002360c0, true, nil
	}
	// 0x1002360c0: eor w14, w14, w17
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002360c4 {
		return 0x1002360c4, true, nil
	}
	// 0x1002360c4: eor w10, w10, w6
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x1002360c8 {
		return 0x1002360c8, true, nil
	}
	// 0x1002360c8: mov w17, #0xa7c8
	s.x[17] = uint64(uint32(uint64(0xa7c8)))
	if runtime != nil && runtime.stopPC == 0x1002360cc {
		return 0x1002360cc, true, nil
	}
	// 0x1002360cc: movk w17, #0x19ea, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x19ea, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002360d0 {
		return 0x1002360d0, true, nil
	}
	// 0x1002360d0: eor w14, w14, w17
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002360d4 {
		return 0x1002360d4, true, nil
	}
	// 0x1002360d4: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002360d8 {
		return 0x1002360d8, true, nil
	}
	// 0x1002360d8: eor w11, w14, w10
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002360dc {
		return 0x1002360dc, true, nil
	}
	// 0x1002360dc: ubfx w10, w16, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002360e0 {
		return 0x1002360e0, true, nil
	}
	// 0x1002360e0: eor w14, w10, w0
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002360e4 {
		return 0x1002360e4, true, nil
	}
	// 0x1002360e4: ldr w14, [x24, w14, uxtw #2]
	value1282, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1282))
	if runtime != nil && runtime.stopPC == 0x1002360e8 {
		return 0x1002360e8, true, nil
	}
	// 0x1002360e8: eor w17, w10, w27
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x1002360ec {
		return 0x1002360ec, true, nil
	}
	// 0x1002360ec: sub w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002360f0 {
		return 0x1002360f0, true, nil
	}
	// 0x1002360f0: eor w10, w14, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002360f4 {
		return 0x1002360f4, true, nil
	}
	// 0x1002360f4: lsr w9, w9, #0x18
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002360f8 {
		return 0x1002360f8, true, nil
	}
	// 0x1002360f8: ldr w9, [x22, w9, uxtw #2]
	value1287, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1287))
	if runtime != nil && runtime.stopPC == 0x1002360fc {
		return 0x1002360fc, true, nil
	}
	// 0x1002360fc: ldr w14, [sp, #0x2ec]
	value1288, err := s.load(addOffset(s.sp, 748), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1288))
	if runtime != nil && runtime.stopPC == 0x100236100 {
		return 0x100236100, true, nil
	}
	// 0x100236100: eor w9, w14, w9, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) ^ rotateRight(uint64(uint32(s.x[9])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100236104 {
		return 0x100236104, true, nil
	}
	// 0x100236104: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236108 {
		return 0x100236108, true, nil
	}
	// 0x100236108: ubfx w14, w15, #8, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x10023610c {
		return 0x10023610c, true, nil
	}
	// 0x10023610c: eor w10, w14, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236110 {
		return 0x100236110, true, nil
	}
	// 0x100236110: ldr w10, [x26, w10, uxtw #2]
	value1293, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1293))
	if runtime != nil && runtime.stopPC == 0x100236114 {
		return 0x100236114, true, nil
	}
	// 0x100236114: and w17, w19, #0xff
	s.x[17] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236118 {
		return 0x100236118, true, nil
	}
	// 0x100236118: eor w17, w17, w28
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[28])))))
	if runtime != nil && runtime.stopPC == 0x10023611c {
		return 0x10023611c, true, nil
	}
	// 0x10023611c: ldr w17, [x2, w17, uxtw #2]
	value1296, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1296))
	if runtime != nil && runtime.stopPC == 0x100236120 {
		return 0x100236120, true, nil
	}
	// 0x100236120: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236124 {
		return 0x100236124, true, nil
	}
	// 0x100236124: mov w1, #0x598d
	s.x[1] = uint64(uint32(uint64(0x598d)))
	if runtime != nil && runtime.stopPC == 0x100236128 {
		return 0x100236128, true, nil
	}
	// 0x100236128: movk w1, #0x847a, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x847a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023612c {
		return 0x10023612c, true, nil
	}
	// 0x10023612c: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236130 {
		return 0x100236130, true, nil
	}
	// 0x100236130: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236134 {
		return 0x100236134, true, nil
	}
	// 0x100236134: and w9, w4, w17, lsr #4
	s.x[9] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[17])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100236138 {
		return 0x100236138, true, nil
	}
	// 0x100236138: eor w14, w14, w9
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023613c {
		return 0x10023613c, true, nil
	}
	// 0x10023613c: ubfx w9, w19, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[19])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236140 {
		return 0x100236140, true, nil
	}
	// 0x100236140: ubfx w17, w3, #8, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[3])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236144 {
		return 0x100236144, true, nil
	}
	// 0x100236144: mov w1, #0xaf
	s.x[1] = uint64(uint32(uint64(0xaf)))
	if runtime != nil && runtime.stopPC == 0x100236148 {
		return 0x100236148, true, nil
	}
	// 0x100236148: eor w1, w9, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x10023614c {
		return 0x10023614c, true, nil
	}
	// 0x10023614c: mov w3, #0x7afa
	s.x[3] = uint64(uint32(uint64(0x7afa)))
	if runtime != nil && runtime.stopPC == 0x100236150 {
		return 0x100236150, true, nil
	}
	// 0x100236150: movk w3, #0xbc9e, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0xbc9e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236154 {
		return 0x100236154, true, nil
	}
	// 0x100236154: eor w3, w9, w3
	s.x[3] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100236158 {
		return 0x100236158, true, nil
	}
	// 0x100236158: ldr w5, [sp, #0x2e8]
	value1311, err := s.load(addOffset(s.sp, 744), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[5] = uint64(uint32(value1311))
	if runtime != nil && runtime.stopPC == 0x10023615c {
		return 0x10023615c, true, nil
	}
	// 0x10023615c: eor w7, w5, w17
	s.x[7] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236160 {
		return 0x100236160, true, nil
	}
	// 0x100236160: sub w1, w1, w3
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) - uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100236164 {
		return 0x100236164, true, nil
	}
	// 0x100236164: eor w1, w7, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236168 {
		return 0x100236168, true, nil
	}
	// 0x100236168: lsr w15, w15, #0x18
	s.x[15] = uint64(uint32(shiftRight(uint64(uint32(s.x[15])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x10023616c {
		return 0x10023616c, true, nil
	}
	// 0x10023616c: ldr w15, [x22, w15, uxtw #2]
	value1316, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value1316))
	if runtime != nil && runtime.stopPC == 0x100236170 {
		return 0x100236170, true, nil
	}
	// 0x100236170: mov w3, #0x3a
	s.x[3] = uint64(uint32(uint64(0x3a)))
	if runtime != nil && runtime.stopPC == 0x100236174 {
		return 0x100236174, true, nil
	}
	// 0x100236174: eor w9, w9, w3
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100236178 {
		return 0x100236178, true, nil
	}
	// 0x100236178: ldr w9, [x24, w9, uxtw #2]
	value1319, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1319))
	if runtime != nil && runtime.stopPC == 0x10023617c {
		return 0x10023617c, true, nil
	}
	// 0x10023617c: eor w15, w1, w15, ror #23
	s.x[15] = uint64(uint32((uint64(uint32(s.x[1])) ^ rotateRight(uint64(uint32(s.x[15])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100236180 {
		return 0x100236180, true, nil
	}
	// 0x100236180: eor w9, w15, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236184 {
		return 0x100236184, true, nil
	}
	// 0x100236184: eor w15, w17, w13
	s.x[15] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236188 {
		return 0x100236188, true, nil
	}
	// 0x100236188: ldr w15, [x26, w15, uxtw #2]
	value1323, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value1323))
	if runtime != nil && runtime.stopPC == 0x10023618c {
		return 0x10023618c, true, nil
	}
	// 0x10023618c: eor w9, w9, w15
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100236190 {
		return 0x100236190, true, nil
	}
	// 0x100236190: and w15, w16, #0xff
	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236194 {
		return 0x100236194, true, nil
	}
	// 0x100236194: ldr w15, [x2, w15, uxtw #2]
	value1326, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[15])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value1326))
	if runtime != nil && runtime.stopPC == 0x100236198 {
		return 0x100236198, true, nil
	}
	// 0x100236198: mov w16, #0xc6ac
	s.x[16] = uint64(uint32(uint64(0xc6ac)))
	if runtime != nil && runtime.stopPC == 0x10023619c {
		return 0x10023619c, true, nil
	}
	// 0x10023619c: movk w16, #0x7d18, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x7d18, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002361a0 {
		return 0x1002361a0, true, nil
	}
	// 0x1002361a0: eor w16, w15, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1002361a4 {
		return 0x1002361a4, true, nil
	}
	// 0x1002361a4: eor w9, w9, w16
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1002361a8 {
		return 0x1002361a8, true, nil
	}
	// 0x1002361a8: and w12, w4, w12, lsr #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[12])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002361ac {
		return 0x1002361ac, true, nil
	}
	// 0x1002361ac: and w16, w4, w15, lsr #4
	s.x[16] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[15])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002361b0 {
		return 0x1002361b0, true, nil
	}
	// 0x1002361b0: mov w15, #0x8362
	s.x[15] = uint64(uint32(uint64(0x8362)))
	if runtime != nil && runtime.stopPC == 0x1002361b4 {
		return 0x1002361b4, true, nil
	}
	// 0x1002361b4: movk w15, #0x949c, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x949c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002361b8 {
		return 0x1002361b8, true, nil
	}
	// 0x1002361b8: eor w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1002361bc {
		return 0x1002361bc, true, nil
	}
	// 0x1002361bc: eor w15, w8, w12
	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002361c0 {
		return 0x1002361c0, true, nil
	}
	// 0x1002361c0: eor w12, w9, w16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1002361c4 {
		return 0x1002361c4, true, nil
	}
	// 0x1002361c4: ubfx w16, w12, #0x10, #8
	s.x[16] = uint64(uint32((uint64(uint32(s.x[12])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002361c8 {
		return 0x1002361c8, true, nil
	}
	// 0x1002361c8: eor w17, w16, w0
	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002361cc {
		return 0x1002361cc, true, nil
	}
	// 0x1002361cc: ldr w17, [x24, w17, uxtw #2]
	value1340, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1340))
	if runtime != nil && runtime.stopPC == 0x1002361d0 {
		return 0x1002361d0, true, nil
	}
	// 0x1002361d0: ubfx w1, w15, #8, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[15])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002361d4 {
		return 0x1002361d4, true, nil
	}
	// 0x1002361d4: eor w3, w1, w13
	s.x[3] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002361d8 {
		return 0x1002361d8, true, nil
	}
	// 0x1002361d8: ldr w3, [x26, w3, uxtw #2]
	value1343, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1343))
	if runtime != nil && runtime.stopPC == 0x1002361dc {
		return 0x1002361dc, true, nil
	}
	// 0x1002361dc: eor w17, w17, w3
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x1002361e0 {
		return 0x1002361e0, true, nil
	}
	// 0x1002361e0: eor w14, w10, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002361e4 {
		return 0x1002361e4, true, nil
	}
	// 0x1002361e4: and w3, w14, #0xff
	s.x[3] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002361e8 {
		return 0x1002361e8, true, nil
	}
	// 0x1002361e8: ldr w3, [x2, w3, uxtw #2]
	value1347, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1347))
	if runtime != nil && runtime.stopPC == 0x1002361ec {
		return 0x1002361ec, true, nil
	}
	// 0x1002361ec: eor w17, w17, w3
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x1002361f0 {
		return 0x1002361f0, true, nil
	}
	// 0x1002361f0: eor w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002361f4 {
		return 0x1002361f4, true, nil
	}
	// 0x1002361f4: and w1, w4, w3, lsr #4
	s.x[1] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[3])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002361f8 {
		return 0x1002361f8, true, nil
	}
	// 0x1002361f8: eor w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002361fc {
		return 0x1002361fc, true, nil
	}
	// 0x1002361fc: lsr w1, w11, #0x18
	s.x[1] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236200 {
		return 0x100236200, true, nil
	}
	// 0x100236200: ldr w1, [x22, w1, uxtw #2]
	value1353, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value1353))
	if runtime != nil && runtime.stopPC == 0x100236204 {
		return 0x100236204, true, nil
	}
	// 0x100236204: ldr w3, [sp, #0x2e4]
	value1354, err := s.load(addOffset(s.sp, 740), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1354))
	if runtime != nil && runtime.stopPC == 0x100236208 {
		return 0x100236208, true, nil
	}
	// 0x100236208: eor w1, w3, w1, ror #23
	s.x[1] = uint64(uint32((uint64(uint32(s.x[3])) ^ rotateRight(uint64(uint32(s.x[1])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x10023620c {
		return 0x10023620c, true, nil
	}
	// 0x10023620c: eor w17, w1, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236210 {
		return 0x100236210, true, nil
	}
	// 0x100236210: eor w1, w16, w27
	s.x[1] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236214 {
		return 0x100236214, true, nil
	}
	// 0x100236214: sub w16, w16, w1
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236218 {
		return 0x100236218, true, nil
	}
	// 0x100236218: eor w16, w17, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x10023621c {
		return 0x10023621c, true, nil
	}
	// 0x10023621c: lsr w10, w10, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236220 {
		return 0x100236220, true, nil
	}
	// 0x100236220: ldr w10, [x22, w10, uxtw #2]
	value1361, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1361))
	if runtime != nil && runtime.stopPC == 0x100236224 {
		return 0x100236224, true, nil
	}
	// 0x100236224: ldr w17, [sp, #0x2e0]
	value1362, err := s.load(addOffset(s.sp, 736), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1362))
	if runtime != nil && runtime.stopPC == 0x100236228 {
		return 0x100236228, true, nil
	}
	// 0x100236228: eor w10, w17, w10, ror #23
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) ^ rotateRight(uint64(uint32(s.x[10])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x10023622c {
		return 0x10023622c, true, nil
	}
	// 0x10023622c: ubfx w17, w11, #0x10, #8
	s.x[17] = uint64(uint32((uint64(uint32(s.x[11])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236230 {
		return 0x100236230, true, nil
	}
	// 0x100236230: eor w1, w17, w27
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236234 {
		return 0x100236234, true, nil
	}
	// 0x100236234: sub w1, w17, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) - uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236238 {
		return 0x100236238, true, nil
	}
	// 0x100236238: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x10023623c {
		return 0x10023623c, true, nil
	}
	// 0x10023623c: ubfx w1, w12, #8, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[12])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236240 {
		return 0x100236240, true, nil
	}
	// 0x100236240: eor w3, w1, w13
	s.x[3] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236244 {
		return 0x100236244, true, nil
	}
	// 0x100236244: ldr w3, [x26, w3, uxtw #2]
	value1370, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1370))
	if runtime != nil && runtime.stopPC == 0x100236248 {
		return 0x100236248, true, nil
	}
	// 0x100236248: eor w10, w10, w3
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x10023624c {
		return 0x10023624c, true, nil
	}
	// 0x10023624c: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236250 {
		return 0x100236250, true, nil
	}
	// 0x100236250: eor w17, w17, w0
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x100236254 {
		return 0x100236254, true, nil
	}
	// 0x100236254: ldr w17, [x24, w17, uxtw #2]
	value1374, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1374))
	if runtime != nil && runtime.stopPC == 0x100236258 {
		return 0x100236258, true, nil
	}
	// 0x100236258: mov w1, #0xdeed
	s.x[1] = uint64(uint32(uint64(0xdeed)))
	if runtime != nil && runtime.stopPC == 0x10023625c {
		return 0x10023625c, true, nil
	}
	// 0x10023625c: movk w1, #0xb99f, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xb99f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236260 {
		return 0x100236260, true, nil
	}
	// 0x100236260: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236264 {
		return 0x100236264, true, nil
	}
	// 0x100236264: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236268 {
		return 0x100236268, true, nil
	}
	// 0x100236268: mov w17, #0x49fc
	s.x[17] = uint64(uint32(uint64(0x49fc)))
	if runtime != nil && runtime.stopPC == 0x10023626c {
		return 0x10023626c, true, nil
	}
	// 0x10023626c: movk w17, #0x404, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x404, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236270 {
		return 0x100236270, true, nil
	}
	// 0x100236270: and w17, w10, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236274 {
		return 0x100236274, true, nil
	}
	// 0x100236274: mov w1, #0xa872
	s.x[1] = uint64(uint32(uint64(0xa872)))
	if runtime != nil && runtime.stopPC == 0x100236278 {
		return 0x100236278, true, nil
	}
	// 0x100236278: movk w1, #0x9cf1, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x9cf1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023627c {
		return 0x10023627c, true, nil
	}
	// 0x10023627c: eor w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236280 {
		return 0x100236280, true, nil
	}
	// 0x100236280: mov w1, #0xb603
	s.x[1] = uint64(uint32(uint64(0xb603)))
	if runtime != nil && runtime.stopPC == 0x100236284 {
		return 0x100236284, true, nil
	}
	// 0x100236284: movk w1, #0xfbfb, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xfbfb, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236288 {
		return 0x100236288, true, nil
	}
	// 0x100236288: and w1, w10, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x10023628c {
		return 0x10023628c, true, nil
	}
	// 0x10023628c: mov w3, #0xdbff
	s.x[3] = uint64(uint32(uint64(0xdbff)))
	if runtime != nil && runtime.stopPC == 0x100236290 {
		return 0x100236290, true, nil
	}
	// 0x100236290: movk w3, #0x8c3c, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0x8c3c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236294 {
		return 0x100236294, true, nil
	}
	// 0x100236294: eor w1, w1, w3
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100236298 {
		return 0x100236298, true, nil
	}
	// 0x100236298: and w17, w17, w1
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x10023629c {
		return 0x10023629c, true, nil
	}
	// 0x10023629c: mov w1, #0x1601
	s.x[1] = uint64(uint32(uint64(0x1601)))
	if runtime != nil && runtime.stopPC == 0x1002362a0 {
		return 0x1002362a0, true, nil
	}
	// 0x1002362a0: movk w1, #0x630a, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x630a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002362a4 {
		return 0x1002362a4, true, nil
	}
	// 0x1002362a4: and w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x1002362a8 {
		return 0x1002362a8, true, nil
	}
	// 0x1002362a8: orr w10, w17, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[17])) | uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002362ac {
		return 0x1002362ac, true, nil
	}
	// 0x1002362ac: and w17, w15, #0xff
	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002362b0 {
		return 0x1002362b0, true, nil
	}
	// 0x1002362b0: ldr w17, [x2, w17, uxtw #2]
	value1397, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1397))
	if runtime != nil && runtime.stopPC == 0x1002362b4 {
		return 0x1002362b4, true, nil
	}
	// 0x1002362b4: and w1, w4, w17, lsr #4
	s.x[1] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[17])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1002362b8 {
		return 0x1002362b8, true, nil
	}
	// 0x1002362b8: eor w17, w1, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002362bc {
		return 0x1002362bc, true, nil
	}
	// 0x1002362bc: lsr w9, w9, #0x18
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002362c0 {
		return 0x1002362c0, true, nil
	}
	// 0x1002362c0: mov w1, #0x49
	s.x[1] = uint64(uint32(uint64(0x49)))
	if runtime != nil && runtime.stopPC == 0x1002362c4 {
		return 0x1002362c4, true, nil
	}
	// 0x1002362c4: eor w8, w1, w8, lsr #24
	s.x[8] = uint64(uint32((uint64(uint32(s.x[1])) ^ shiftRight(uint64(uint32(s.x[8])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x1002362c8 {
		return 0x1002362c8, true, nil
	}
	// 0x1002362c8: ldr w8, [x22, w8, uxtw #2]
	value1403, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1403))
	if runtime != nil && runtime.stopPC == 0x1002362cc {
		return 0x1002362cc, true, nil
	}
	// 0x1002362cc: ldr w9, [x22, w9, uxtw #2]
	value1404, err := s.load((s.x[22] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1404))
	if runtime != nil && runtime.stopPC == 0x1002362d0 {
		return 0x1002362d0, true, nil
	}
	// 0x1002362d0: ubfx w1, w14, #0x10, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[14])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002362d4 {
		return 0x1002362d4, true, nil
	}
	// 0x1002362d4: eor w3, w1, w0
	s.x[3] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002362d8 {
		return 0x1002362d8, true, nil
	}
	// 0x1002362d8: ldr w3, [x24, w3, uxtw #2]
	value1407, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[3] = uint64(uint32(value1407))
	if runtime != nil && runtime.stopPC == 0x1002362dc {
		return 0x1002362dc, true, nil
	}
	// 0x1002362dc: ubfx w15, w15, #0x10, #8
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002362e0 {
		return 0x1002362e0, true, nil
	}
	// 0x1002362e0: eor w7, w15, w0
	s.x[7] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002362e4 {
		return 0x1002362e4, true, nil
	}
	// 0x1002362e4: ldr w7, [x24, w7, uxtw #2]
	value1410, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[7])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value1410))
	if runtime != nil && runtime.stopPC == 0x1002362e8 {
		return 0x1002362e8, true, nil
	}
	// 0x1002362e8: ubfx w22, w11, #8, #8
	s.x[22] = uint64(uint32((uint64(uint32(s.x[11])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002362ec {
		return 0x1002362ec, true, nil
	}
	// 0x1002362ec: eor w19, w22, w13
	s.x[19] = uint64(uint32((uint64(uint32(s.x[22])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1002362f0 {
		return 0x1002362f0, true, nil
	}
	// 0x1002362f0: ldr w24, [x26, w19, uxtw #2]
	value1413, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[19])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[24] = uint64(uint32(value1413))
	if runtime != nil && runtime.stopPC == 0x1002362f4 {
		return 0x1002362f4, true, nil
	}
	// 0x1002362f4: eor w19, w17, w10
	s.x[19] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002362f8 {
		return 0x1002362f8, true, nil
	}
	// 0x1002362f8: ubfx w10, w14, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002362fc {
		return 0x1002362fc, true, nil
	}
	// 0x1002362fc: eor w14, w10, w13
	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236300 {
		return 0x100236300, true, nil
	}
	// 0x100236300: ldr w14, [x26, w14, uxtw #2]
	value1417, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1417))
	if runtime != nil && runtime.stopPC == 0x100236304 {
		return 0x100236304, true, nil
	}
	// 0x100236304: and w12, w12, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236308 {
		return 0x100236308, true, nil
	}
	// 0x100236308: ldr w12, [x2, w12, uxtw #2]
	value1419, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value1419))
	if runtime != nil && runtime.stopPC == 0x10023630c {
		return 0x10023630c, true, nil
	}
	// 0x10023630c: and w11, w11, #0xff
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236310 {
		return 0x100236310, true, nil
	}
	// 0x100236310: ldr w17, [x2, w11, uxtw #2]
	value1421, err := s.load((s.x[2] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1421))
	if runtime != nil && runtime.stopPC == 0x100236314 {
		return 0x100236314, true, nil
	}
	// 0x100236314: eor w11, w1, w27
	s.x[11] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236318 {
		return 0x100236318, true, nil
	}
	// 0x100236318: sub w11, w1, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[1])) - uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023631c {
		return 0x10023631c, true, nil
	}
	// 0x10023631c: ldr w13, [sp, #0x2dc]
	value1424, err := s.load(addOffset(s.sp, 732), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1424))
	if runtime != nil && runtime.stopPC == 0x100236320 {
		return 0x100236320, true, nil
	}
	// 0x100236320: eor w11, w13, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236324 {
		return 0x100236324, true, nil
	}
	// 0x100236324: eor w11, w11, w22
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[22])))))
	if runtime != nil && runtime.stopPC == 0x100236328 {
		return 0x100236328, true, nil
	}
	// 0x100236328: eor w8, w11, w8, ror #23
	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) ^ rotateRight(uint64(uint32(s.x[8])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x10023632c {
		return 0x10023632c, true, nil
	}
	// 0x10023632c: eor w8, w8, w3
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100236330 {
		return 0x100236330, true, nil
	}
	// 0x100236330: eor w8, w8, w12
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236334 {
		return 0x100236334, true, nil
	}
	// 0x100236334: eor w8, w8, w24
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[24])))))
	if runtime != nil && runtime.stopPC == 0x100236338 {
		return 0x100236338, true, nil
	}
	// 0x100236338: ldr w11, [sp, #0x3b8]
	value1431, err := s.load(addOffset(s.sp, 952), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value1431))
	if runtime != nil && runtime.stopPC == 0x10023633c {
		return 0x10023633c, true, nil
	}
	// 0x10023633c: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236340 {
		return 0x100236340, true, nil
	}
	// 0x100236340: eor w11, w15, w27
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236344 {
		return 0x100236344, true, nil
	}
	// 0x100236344: sub w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) - uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236348 {
		return 0x100236348, true, nil
	}
	// 0x100236348: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023634c {
		return 0x10023634c, true, nil
	}
	// 0x10023634c: eor w10, w10, w7
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x100236350 {
		return 0x100236350, true, nil
	}
	// 0x100236350: eor w9, w10, w9, ror #23
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ rotateRight(uint64(uint32(s.x[9])), 23, 32))))
	if runtime != nil && runtime.stopPC == 0x100236354 {
		return 0x100236354, true, nil
	}
	// 0x100236354: eor w10, w14, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236358 {
		return 0x100236358, true, nil
	}
	// 0x100236358: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x10023635c {
		return 0x10023635c, true, nil
	}
	// 0x10023635c: ldr x10, [sp, #0x410]
	value1440, err := s.load(addOffset(s.sp, 1040), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1440)
	if runtime != nil && runtime.stopPC == 0x100236360 {
		return 0x100236360, true, nil
	}
	// 0x100236360: ldur x11, [x29, #-0xe8]
	value1441, err := s.load(addOffset(s.x[29], -232), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value1441)
	if runtime != nil && runtime.stopPC == 0x100236364 {
		return 0x100236364, true, nil
	}
	// 0x100236364: add x10, x11, x10
	s.x[10] = uint64((s.x[11] + s.x[10]))
	if runtime != nil && runtime.stopPC == 0x100236368 {
		return 0x100236368, true, nil
	}
	// 0x100236368: ldrb w10, [x10]
	value1443, err := s.load(s.x[10], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1443))
	if runtime != nil && runtime.stopPC == 0x10023636c {
		return 0x10023636c, true, nil
	}
	// 0x10023636c: mov w11, #0x29
	s.x[11] = uint64(uint32(uint64(0x29)))
	if runtime != nil && runtime.stopPC == 0x100236370 {
		return 0x100236370, true, nil
	}
	// 0x100236370: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236374 {
		return 0x100236374, true, nil
	}
	// 0x100236374: ldur x11, [x29, #-0xf0]
	value1446, err := s.load(addOffset(s.x[29], -240), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value1446)
	if runtime != nil && runtime.stopPC == 0x100236378 {
		return 0x100236378, true, nil
	}
	// 0x100236378: ldrb w10, [x11, w10, uxtw]
	value1447, err := s.load((s.x[11] + uint64(uint32(uint64(uint32(s.x[10]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1447))
	if runtime != nil && runtime.stopPC == 0x10023637c {
		return 0x10023637c, true, nil
	}
	// 0x10023637c: eor w13, w23, w10
	s.x[13] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236380 {
		return 0x100236380, true, nil
	}
	// 0x100236380: ldr w10, [sp, #0x424]
	value1449, err := s.load(addOffset(s.sp, 1060), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1449))
	if runtime != nil && runtime.stopPC == 0x100236384 {
		return 0x100236384, true, nil
	}
	// 0x100236384: add w11, w10, w21
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x100236388 {
		return 0x100236388, true, nil
	}
	// 0x100236388: and w10, w4, w12, lsr #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[12])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x10023638c {
		return 0x10023638c, true, nil
	}
	// 0x10023638c: and w14, w4, w17, lsr #4
	s.x[14] = uint64(uint32((uint64(uint32(s.x[4])) & shiftRight(uint64(uint32(s.x[17])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100236390 {
		return 0x100236390, true, nil
	}
	// 0x100236390: mov w12, #0xf28e
	s.x[12] = uint64(uint32(uint64(0xf28e)))
	if runtime != nil && runtime.stopPC == 0x100236394 {
		return 0x100236394, true, nil
	}
	// 0x100236394: movk w12, #0xf722, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xf722, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236398 {
		return 0x100236398, true, nil
	}
	// 0x100236398: eor w12, w8, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10023639c {
		return 0x10023639c, true, nil
	}
	// 0x10023639c: eor w8, w12, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002363a0 {
		return 0x1002363a0, true, nil
	}
	// 0x1002363a0: lsr w10, w12, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[12])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002363a4 {
		return 0x1002363a4, true, nil
	}
	// 0x1002363a4: mov w12, #0xc930
	s.x[12] = uint64(uint32(uint64(0xc930)))
	if runtime != nil && runtime.stopPC == 0x1002363a8 {
		return 0x1002363a8, true, nil
	}
	// 0x1002363a8: movk w12, #0xc66a, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xc66a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002363ac {
		return 0x1002363ac, true, nil
	}
	// 0x1002363ac: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1002363b0 {
		return 0x1002363b0, true, nil
	}
	// 0x1002363b0: and w9, w13, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002363b4 {
		return 0x1002363b4, true, nil
	}
	// 0x1002363b4: add w9, w9, #0x82a
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0x82a))))
	if runtime != nil && runtime.stopPC == 0x1002363b8 {
		return 0x1002363b8, true, nil
	}
	// 0x1002363b8: mov x26, x30
	s.x[26] = uint64(s.x[30])
	if runtime != nil && runtime.stopPC == 0x1002363bc {
		return 0x1002363bc, true, nil
	}
	// 0x1002363bc: ldr x9, [x30, w9, uxtw #3]
	value1464, err := s.load((s.x[30] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1464)
	if runtime != nil && runtime.stopPC == 0x1002363c0 {
		return 0x1002363c0, true, nil
	}
	// 0x1002363c0: mov x13, #-0x690a
	s.x[13] = uint64(uint64(0xffffffffffff96f6))
	if runtime != nil && runtime.stopPC == 0x1002363c4 {
		return 0x1002363c4, true, nil
	}
	// 0x1002363c4: movk x13, #0xccb9, lsl #16
	s.x[13] = uint64(insertBits(s.x[13], 0xccb9, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x1002363c8 {
		return 0x1002363c8, true, nil
	}
	// 0x1002363c8: add x9, x9, x13
	s.x[9] = uint64((s.x[9] + s.x[13]))
	if runtime != nil && runtime.stopPC == 0x1002363cc {
		return 0x1002363cc, true, nil
	}
	// 0x1002363cc: mov w0, #0xa3
	s.x[0] = uint64(uint32(uint64(0xa3)))
	if runtime != nil && runtime.stopPC == 0x1002363d0 {
		return 0x1002363d0, true, nil
	}
	// 0x1002363d0: eor w13, w0, w16, lsr #24
	s.x[13] = uint64(uint32((uint64(uint32(s.x[0])) ^ shiftRight(uint64(uint32(s.x[16])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x1002363d4 {
		return 0x1002363d4, true, nil
	}
	// 0x1002363d4: ldr w30, [x9, w13, uxtw #2]
	value1470, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[13])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[30] = uint64(uint32(value1470))
	if runtime != nil && runtime.stopPC == 0x1002363d8 {
		return 0x1002363d8, true, nil
	}
	// 0x1002363d8: eor w24, w12, w14
	s.x[24] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002363dc {
		return 0x1002363dc, true, nil
	}
	// 0x1002363dc: mov w13, #0x84
	s.x[13] = uint64(uint32(uint64(0x84)))
	if runtime != nil && runtime.stopPC == 0x1002363e0 {
		return 0x1002363e0, true, nil
	}
	// 0x1002363e0: eor w13, w13, w19, lsr #24
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ shiftRight(uint64(uint32(s.x[19])), 24, 32))))
	if runtime != nil && runtime.stopPC == 0x1002363e4 {
		return 0x1002363e4, true, nil
	}
	// 0x1002363e4: ldr w25, [x9, w13, uxtw #2]
	value1474, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[13])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[25] = uint64(uint32(value1474))
	if runtime != nil && runtime.stopPC == 0x1002363e8 {
		return 0x1002363e8, true, nil
	}
	// 0x1002363e8: ldr w10, [x9, w10, uxtw #2]
	value1475, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1475))
	if runtime != nil && runtime.stopPC == 0x1002363ec {
		return 0x1002363ec, true, nil
	}
	// 0x1002363ec: str w10, [sp, #0x424]
	if err := s.storeValue(addOffset(s.sp, 1060), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002363f0 {
		return 0x1002363f0, true, nil
	}
	// 0x1002363f0: lsr w10, w12, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[12])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002363f4 {
		return 0x1002363f4, true, nil
	}
	// 0x1002363f4: ldr w9, [x9, w10, uxtw #2]
	value1478, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1478))
	if runtime != nil && runtime.stopPC == 0x1002363f8 {
		return 0x1002363f8, true, nil
	}
	// 0x1002363f8: str w9, [sp, #0x410]
	if err := s.storeValue(addOffset(s.sp, 1040), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002363fc {
		return 0x1002363fc, true, nil
	}
	// 0x1002363fc: mov w9, #0xdd5d
	s.x[9] = uint64(uint32(uint64(0xdd5d)))
	if runtime != nil && runtime.stopPC == 0x100236400 {
		return 0x100236400, true, nil
	}
	// 0x100236400: movk w9, #0x2a2a, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x2a2a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236404 {
		return 0x100236404, true, nil
	}
	// 0x100236404: eor w10, w19, w9
	s.x[10] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236408 {
		return 0x100236408, true, nil
	}
	// 0x100236408: mov w27, #0xa6
	s.x[27] = uint64(uint32(uint64(0xa6)))
	if runtime != nil && runtime.stopPC == 0x10023640c {
		return 0x10023640c, true, nil
	}
	// 0x10023640c: eor w9, w11, w27
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236410 {
		return 0x100236410, true, nil
	}
	// 0x100236410: mov w11, #0xd4cf
	s.x[11] = uint64(uint32(uint64(0xd4cf)))
	if runtime != nil && runtime.stopPC == 0x100236414 {
		return 0x100236414, true, nil
	}
	// 0x100236414: movk w11, #0xa5c, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa5c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236418 {
		return 0x100236418, true, nil
	}
	// 0x100236418: mul w12, w9, w11
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023641c {
		return 0x10023641c, true, nil
	}
	// 0x10023641c: ldur x9, [x29, #-0xb8]
	value1488, err := s.load(addOffset(s.x[29], -184), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1488)
	if runtime != nil && runtime.stopPC == 0x100236420 {
		return 0x100236420, true, nil
	}
	// 0x100236420: add x9, x9, w12, uxtb
	s.x[9] = uint64((s.x[9] + uint64(uint8(uint64(uint32(s.x[12]))))))
	if runtime != nil && runtime.stopPC == 0x100236424 {
		return 0x100236424, true, nil
	}
	// 0x100236424: ldrb w17, [x9]
	value1490, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1490))
	if runtime != nil && runtime.stopPC == 0x100236428 {
		return 0x100236428, true, nil
	}
	// 0x100236428: mov w9, #0x5d
	s.x[9] = uint64(uint32(uint64(0x5d)))
	if runtime != nil && runtime.stopPC == 0x10023642c {
		return 0x10023642c, true, nil
	}
	// 0x10023642c: eor w9, w17, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236430 {
		return 0x100236430, true, nil
	}
	// 0x100236430: ldur x11, [x29, #-0xa8]
	value1493, err := s.load(addOffset(s.x[29], -168), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value1493)
	if runtime != nil && runtime.stopPC == 0x100236434 {
		return 0x100236434, true, nil
	}
	// 0x100236434: ldrb w9, [x11, w9, uxtw]
	value1494, err := s.load((s.x[11] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1494))
	if runtime != nil && runtime.stopPC == 0x100236438 {
		return 0x100236438, true, nil
	}
	// 0x100236438: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023643c {
		return 0x10023643c, true, nil
	}
	// 0x10023643c: and w9, w9, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236440 {
		return 0x100236440, true, nil
	}
	// 0x100236440: add w9, w9, #0x547
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0x547))))
	if runtime != nil && runtime.stopPC == 0x100236444 {
		return 0x100236444, true, nil
	}
	// 0x100236444: ldr x9, [x26, w9, uxtw #3]
	value1498, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1498)
	if runtime != nil && runtime.stopPC == 0x100236448 {
		return 0x100236448, true, nil
	}
	// 0x100236448: mov x11, #-0x8563
	s.x[11] = uint64(uint64(0xffffffffffff7a9d))
	if runtime != nil && runtime.stopPC == 0x10023644c {
		return 0x10023644c, true, nil
	}
	// 0x10023644c: movk x11, #0xd7f8, lsl #16
	s.x[11] = uint64(insertBits(s.x[11], 0xd7f8, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x100236450 {
		return 0x100236450, true, nil
	}
	// 0x100236450: add x9, x9, x11
	s.x[9] = uint64((s.x[9] + s.x[11]))
	if runtime != nil && runtime.stopPC == 0x100236454 {
		return 0x100236454, true, nil
	}
	// 0x100236454: ubfx w11, w24, #0x10, #8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[24])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236458 {
		return 0x100236458, true, nil
	}
	// 0x100236458: ldr w2, [x9, w11, uxtw #2]
	value1503, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[11])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[2] = uint64(uint32(value1503))
	if runtime != nil && runtime.stopPC == 0x10023645c {
		return 0x10023645c, true, nil
	}
	// 0x10023645c: ubfx w14, w16, #0x10, #8
	s.x[14] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236460 {
		return 0x100236460, true, nil
	}
	// 0x100236460: mov w13, #0xd9
	s.x[13] = uint64(uint32(uint64(0xd9)))
	if runtime != nil && runtime.stopPC == 0x100236464 {
		return 0x100236464, true, nil
	}
	// 0x100236464: eor w14, w14, w13
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236468 {
		return 0x100236468, true, nil
	}
	// 0x100236468: ldr w14, [x9, w14, uxtw #2]
	value1507, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[14])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value1507))
	if runtime != nil && runtime.stopPC == 0x10023646c {
		return 0x10023646c, true, nil
	}
	// 0x10023646c: ubfx w1, w10, #0x10, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[10])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236470 {
		return 0x100236470, true, nil
	}
	// 0x100236470: ldr w5, [x9, w1, uxtw #2]
	value1509, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[5] = uint64(uint32(value1509))
	if runtime != nil && runtime.stopPC == 0x100236474 {
		return 0x100236474, true, nil
	}
	// 0x100236474: ubfx w1, w8, #0x10, #8
	s.x[1] = uint64(uint32((uint64(uint32(s.x[8])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236478 {
		return 0x100236478, true, nil
	}
	// 0x100236478: ldr w23, [x9, w1, uxtw #2]
	value1511, err := s.load((s.x[9] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[23] = uint64(uint32(value1511))
	if runtime != nil && runtime.stopPC == 0x10023647c {
		return 0x10023647c, true, nil
	}
	// 0x10023647c: ldur x7, [x29, #-0x100]
	value1512, err := s.load(addOffset(s.x[29], -256), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(value1512)
	if runtime != nil && runtime.stopPC == 0x100236480 {
		return 0x100236480, true, nil
	}
	// 0x100236480: ldr x13, [sp, #0x418]
	value1513, err := s.load(addOffset(s.sp, 1048), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value1513)
	if runtime != nil && runtime.stopPC == 0x100236484 {
		return 0x100236484, true, nil
	}
	// 0x100236484: add x9, x7, x13
	s.x[9] = uint64((s.x[7] + s.x[13]))
	if runtime != nil && runtime.stopPC == 0x100236488 {
		return 0x100236488, true, nil
	}
	// 0x100236488: ldrb w9, [x9]
	value1515, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1515))
	if runtime != nil && runtime.stopPC == 0x10023648c {
		return 0x10023648c, true, nil
	}
	// 0x10023648c: mov w11, #0xfb
	s.x[11] = uint64(uint32(uint64(0xfb)))
	if runtime != nil && runtime.stopPC == 0x100236490 {
		return 0x100236490, true, nil
	}
	// 0x100236490: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236494 {
		return 0x100236494, true, nil
	}
	// 0x100236494: ldur x28, [x29, #-0xc0]
	value1518, err := s.load(addOffset(s.x[29], -192), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[28] = uint64(value1518)
	if runtime != nil && runtime.stopPC == 0x100236498 {
		return 0x100236498, true, nil
	}
	// 0x100236498: ldrb w9, [x28, w9, uxtw]
	value1519, err := s.load((s.x[28] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1519))
	if runtime != nil && runtime.stopPC == 0x10023649c {
		return 0x10023649c, true, nil
	}
	// 0x10023649c: eor w9, w20, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002364a0 {
		return 0x1002364a0, true, nil
	}
	// 0x1002364a0: and w9, w9, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002364a4 {
		return 0x1002364a4, true, nil
	}
	// 0x1002364a4: add w9, w9, #0x27
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0x27))))
	if runtime != nil && runtime.stopPC == 0x1002364a8 {
		return 0x1002364a8, true, nil
	}
	// 0x1002364a8: ldr x9, [x26, w9, uxtw #3]
	value1523, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1523)
	if runtime != nil && runtime.stopPC == 0x1002364ac {
		return 0x1002364ac, true, nil
	}
	// 0x1002364ac: mov x1, #-0x97bf
	s.x[1] = uint64(uint64(0xffffffffffff6841))
	if runtime != nil && runtime.stopPC == 0x1002364b0 {
		return 0x1002364b0, true, nil
	}
	// 0x1002364b0: movk x1, #0x81ae, lsl #16
	s.x[1] = uint64(insertBits(s.x[1], 0x81ae, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x1002364b4 {
		return 0x1002364b4, true, nil
	}
	// 0x1002364b4: add x1, x9, x1
	s.x[1] = uint64((s.x[9] + s.x[1]))
	if runtime != nil && runtime.stopPC == 0x1002364b8 {
		return 0x1002364b8, true, nil
	}
	// 0x1002364b8: ubfx w9, w8, #8, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002364bc {
		return 0x1002364bc, true, nil
	}
	// 0x1002364bc: ldr w19, [x1, w9, uxtw #2]
	value1528, err := s.load((s.x[1] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[19] = uint64(uint32(value1528))
	if runtime != nil && runtime.stopPC == 0x1002364c0 {
		return 0x1002364c0, true, nil
	}
	// 0x1002364c0: ubfx w9, w24, #8, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[24])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002364c4 {
		return 0x1002364c4, true, nil
	}
	// 0x1002364c4: ldr w9, [x1, w9, uxtw #2]
	value1530, err := s.load((s.x[1] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1530))
	if runtime != nil && runtime.stopPC == 0x1002364c8 {
		return 0x1002364c8, true, nil
	}
	// 0x1002364c8: ubfx w3, w16, #8, #8
	s.x[3] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002364cc {
		return 0x1002364cc, true, nil
	}
	// 0x1002364cc: eor w3, w3, w0
	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) ^ uint64(uint32(s.x[0])))))
	if runtime != nil && runtime.stopPC == 0x1002364d0 {
		return 0x1002364d0, true, nil
	}
	// 0x1002364d0: ldr w6, [x1, w3, uxtw #2]
	value1533, err := s.load((s.x[1] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[6] = uint64(uint32(value1533))
	if runtime != nil && runtime.stopPC == 0x1002364d4 {
		return 0x1002364d4, true, nil
	}
	// 0x1002364d4: ubfx w3, w10, #8, #8
	s.x[3] = uint64(uint32((uint64(uint32(s.x[10])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002364d8 {
		return 0x1002364d8, true, nil
	}
	// 0x1002364d8: ldr w22, [x1, w3, uxtw #2]
	value1535, err := s.load((s.x[1] + shiftLeft(uint64(uint32(uint64(uint32(s.x[3])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[22] = uint64(uint32(value1535))
	if runtime != nil && runtime.stopPC == 0x1002364dc {
		return 0x1002364dc, true, nil
	}
	// 0x1002364dc: eor w17, w17, #4
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(0x4))))
	if runtime != nil && runtime.stopPC == 0x1002364e0 {
		return 0x1002364e0, true, nil
	}
	// 0x1002364e0: ldur x1, [x29, #-0xa8]
	value1537, err := s.load(addOffset(s.x[29], -168), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(value1537)
	if runtime != nil && runtime.stopPC == 0x1002364e4 {
		return 0x1002364e4, true, nil
	}
	// 0x1002364e4: ldrb w17, [x1, w17, uxtw]
	value1538, err := s.load((s.x[1] + uint64(uint32(uint64(uint32(s.x[17]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value1538))
	if runtime != nil && runtime.stopPC == 0x1002364e8 {
		return 0x1002364e8, true, nil
	}
	// 0x1002364e8: eor w12, w12, w17
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002364ec {
		return 0x1002364ec, true, nil
	}
	// 0x1002364ec: and w10, w10, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002364f0 {
		return 0x1002364f0, true, nil
	}
	// 0x1002364f0: and w12, w12, #0xff
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002364f4 {
		return 0x1002364f4, true, nil
	}
	// 0x1002364f4: add w12, w12, #0xaa8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0xaa8))))
	if runtime != nil && runtime.stopPC == 0x1002364f8 {
		return 0x1002364f8, true, nil
	}
	// 0x1002364f8: ldr x12, [x26, w12, uxtw #3]
	value1543, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(value1543)
	if runtime != nil && runtime.stopPC == 0x1002364fc {
		return 0x1002364fc, true, nil
	}
	// 0x1002364fc: mov x17, #-0xefa6
	s.x[17] = uint64(uint64(0xffffffffffff105a))
	if runtime != nil && runtime.stopPC == 0x100236500 {
		return 0x100236500, true, nil
	}
	// 0x100236500: movk x17, #0xcc12, lsl #16
	s.x[17] = uint64(insertBits(s.x[17], 0xcc12, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x100236504 {
		return 0x100236504, true, nil
	}
	// 0x100236504: add x17, x12, x17
	s.x[17] = uint64((s.x[12] + s.x[17]))
	if runtime != nil && runtime.stopPC == 0x100236508 {
		return 0x100236508, true, nil
	}
	// 0x100236508: ldr w21, [x17, w10, uxtw #2]
	value1547, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[21] = uint64(uint32(value1547))
	if runtime != nil && runtime.stopPC == 0x10023650c {
		return 0x10023650c, true, nil
	}
	// 0x10023650c: and w8, w8, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236510 {
		return 0x100236510, true, nil
	}
	// 0x100236510: ldr w8, [x17, w8, uxtw #2]
	value1549, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1549))
	if runtime != nil && runtime.stopPC == 0x100236514 {
		return 0x100236514, true, nil
	}
	// 0x100236514: and w10, w24, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[24])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236518 {
		return 0x100236518, true, nil
	}
	// 0x100236518: ldr w24, [x17, w10, uxtw #2]
	value1551, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[24] = uint64(uint32(value1551))
	if runtime != nil && runtime.stopPC == 0x10023651c {
		return 0x10023651c, true, nil
	}
	// 0x10023651c: and w10, w16, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236520 {
		return 0x100236520, true, nil
	}
	// 0x100236520: mov w16, #0xe3
	s.x[16] = uint64(uint32(uint64(0xe3)))
	if runtime != nil && runtime.stopPC == 0x100236524 {
		return 0x100236524, true, nil
	}
	// 0x100236524: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100236528 {
		return 0x100236528, true, nil
	}
	// 0x100236528: ldr w10, [x17, w10, uxtw #2]
	value1555, err := s.load((s.x[17] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 2, 64)), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1555))
	if runtime != nil && runtime.stopPC == 0x10023652c {
		return 0x10023652c, true, nil
	}
	// 0x10023652c: str w10, [sp, #0x408]
	if err := s.storeValue(addOffset(s.sp, 1032), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236530 {
		return 0x100236530, true, nil
	}
	// 0x100236530: ldr w10, [sp, #0x2cc]
	value1557, err := s.load(addOffset(s.sp, 716), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1557))
	if runtime != nil && runtime.stopPC == 0x100236534 {
		return 0x100236534, true, nil
	}
	// 0x100236534: ldr w16, [sp, #0x4b4]
	value1558, err := s.load(addOffset(s.sp, 1204), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value1558))
	if runtime != nil && runtime.stopPC == 0x100236538 {
		return 0x100236538, true, nil
	}
	// 0x100236538: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x10023653c {
		return 0x10023653c, true, nil
	}
	// 0x10023653c: ldr w16, [sp, #0x2d4]
	value1560, err := s.load(addOffset(s.sp, 724), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value1560))
	if runtime != nil && runtime.stopPC == 0x100236540 {
		return 0x100236540, true, nil
	}
	// 0x100236540: eor w10, w16, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236544 {
		return 0x100236544, true, nil
	}
	// 0x100236544: eor w10, w10, w9
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236548 {
		return 0x100236548, true, nil
	}
	// 0x100236548: mov w17, #0x7a0d
	s.x[17] = uint64(uint32(uint64(0x7a0d)))
	if runtime != nil && runtime.stopPC == 0x10023654c {
		return 0x10023654c, true, nil
	}
	// 0x10023654c: movk w17, #0xab36, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xab36, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236550 {
		return 0x100236550, true, nil
	}
	// 0x100236550: add w16, w25, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[25])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236554 {
		return 0x100236554, true, nil
	}
	// 0x100236554: adrp x25, #0x101075000
	s.x[25] = uint64(0x101075000)
	if runtime != nil && runtime.stopPC == 0x100236558 {
		return 0x100236558, true, nil
	}
	// 0x100236558: eor w10, w10, w16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x10023655c {
		return 0x10023655c, true, nil
	}
	// 0x10023655c: eor w10, w10, w9, lsr #2
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[9])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100236560 {
		return 0x100236560, true, nil
	}
	// 0x100236560: eor w9, w10, w9, lsr #3
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[9])), 3, 32))))
	if runtime != nil && runtime.stopPC == 0x100236564 {
		return 0x100236564, true, nil
	}
	// 0x100236564: mov w1, #0x1edc
	s.x[1] = uint64(uint32(uint64(0x1edc)))
	if runtime != nil && runtime.stopPC == 0x100236568 {
		return 0x100236568, true, nil
	}
	// 0x100236568: movk w1, #0x7a2, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x7a2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023656c {
		return 0x10023656c, true, nil
	}
	// 0x10023656c: and w10, w1, w14, lsl #1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[1])) & shiftLeft(uint64(uint32(s.x[14])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236570 {
		return 0x100236570, true, nil
	}
	// 0x100236570: mov w3, #0xc54
	s.x[3] = uint64(uint32(uint64(0xc54)))
	if runtime != nil && runtime.stopPC == 0x100236574 {
		return 0x100236574, true, nil
	}
	// 0x100236574: movk w3, #0x382, lsl #16
	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0x382, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236578 {
		return 0x100236578, true, nil
	}
	// 0x100236578: eor w10, w10, w3
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x10023657c {
		return 0x10023657c, true, nil
	}
	// 0x10023657c: mov w4, #0xf954
	s.x[4] = uint64(uint32(uint64(0xf954)))
	if runtime != nil && runtime.stopPC == 0x100236580 {
		return 0x100236580, true, nil
	}
	// 0x100236580: movk w4, #0x461a, lsl #16
	s.x[4] = uint64(uint32(insertBits(uint64(uint32(s.x[4])), 0x461a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236584 {
		return 0x100236584, true, nil
	}
	// 0x100236584: eor w14, w14, w4
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x100236588 {
		return 0x100236588, true, nil
	}
	// 0x100236588: add w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10023658c {
		return 0x10023658c, true, nil
	}
	// 0x10023658c: eor w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236590 {
		return 0x100236590, true, nil
	}
	// 0x100236590: mov w14, #0x663
	s.x[14] = uint64(uint32(uint64(0x663)))
	if runtime != nil && runtime.stopPC == 0x100236594 {
		return 0x100236594, true, nil
	}
	// 0x100236594: movk w14, #0x5cb0, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x5cb0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236598 {
		return 0x100236598, true, nil
	}
	// 0x100236598: and w10, w14, w10, lsr #1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) & shiftRight(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x10023659c {
		return 0x10023659c, true, nil
	}
	// 0x10023659c: eor w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002365a0 {
		return 0x1002365a0, true, nil
	}
	// 0x1002365a0: eor w16, w9, w8
	s.x[16] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002365a4 {
		return 0x1002365a4, true, nil
	}
	// 0x1002365a4: and w8, w16, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002365a8 {
		return 0x1002365a8, true, nil
	}
	// 0x1002365a8: eor w9, w8, #0x60
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0x60))))
	if runtime != nil && runtime.stopPC == 0x1002365ac {
		return 0x1002365ac, true, nil
	}
	// 0x1002365ac: mov w10, #0x34
	s.x[10] = uint64(uint32(uint64(0x34)))
	if runtime != nil && runtime.stopPC == 0x1002365b0 {
		return 0x1002365b0, true, nil
	}
	// 0x1002365b0: eor w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002365b4 {
		return 0x1002365b4, true, nil
	}
	// 0x1002365b4: sub w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002365b8 {
		return 0x1002365b8, true, nil
	}
	// 0x1002365b8: eor w8, w8, #0xfffffff8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0xfffffff8))))
	if runtime != nil && runtime.stopPC == 0x1002365bc {
		return 0x1002365bc, true, nil
	}
	// 0x1002365bc: add w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002365c0 {
		return 0x1002365c0, true, nil
	}
	// 0x1002365c0: ldur x9, [x29, #-0xf8]
	value1593, err := s.load(addOffset(s.x[29], -248), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1593)
	if runtime != nil && runtime.stopPC == 0x1002365c4 {
		return 0x1002365c4, true, nil
	}
	// 0x1002365c4: add x9, x9, x13
	s.x[9] = uint64((s.x[9] + s.x[13]))
	if runtime != nil && runtime.stopPC == 0x1002365c8 {
		return 0x1002365c8, true, nil
	}
	// 0x1002365c8: ldrb w9, [x9]
	value1595, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1595))
	if runtime != nil && runtime.stopPC == 0x1002365cc {
		return 0x1002365cc, true, nil
	}
	// 0x1002365cc: mov w10, #0xb5
	s.x[10] = uint64(uint32(uint64(0xb5)))
	if runtime != nil && runtime.stopPC == 0x1002365d0 {
		return 0x1002365d0, true, nil
	}
	// 0x1002365d0: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002365d4 {
		return 0x1002365d4, true, nil
	}
	// 0x1002365d4: ldr x10, [sp, #0x4c8]
	value1598, err := s.load(addOffset(s.sp, 1224), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1598)
	if runtime != nil && runtime.stopPC == 0x1002365d8 {
		return 0x1002365d8, true, nil
	}
	// 0x1002365d8: ldrb w9, [x10, w9, uxtw]
	value1599, err := s.load((s.x[10] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1599))
	if runtime != nil && runtime.stopPC == 0x1002365dc {
		return 0x1002365dc, true, nil
	}
	// 0x1002365dc: eor w9, w20, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002365e0 {
		return 0x1002365e0, true, nil
	}
	// 0x1002365e0: and w9, w9, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002365e4 {
		return 0x1002365e4, true, nil
	}
	// 0x1002365e4: add w9, w9, #0xe94
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0xe94))))
	if runtime != nil && runtime.stopPC == 0x1002365e8 {
		return 0x1002365e8, true, nil
	}
	// 0x1002365e8: ldr x9, [x26, w9, uxtw #3]
	value1603, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1603)
	if runtime != nil && runtime.stopPC == 0x1002365ec {
		return 0x1002365ec, true, nil
	}
	// 0x1002365ec: mov x10, #-0x693f
	s.x[10] = uint64(uint64(0xffffffffffff96c1))
	if runtime != nil && runtime.stopPC == 0x1002365f0 {
		return 0x1002365f0, true, nil
	}
	// 0x1002365f0: movk x10, #0xec5a, lsl #16
	s.x[10] = uint64(insertBits(s.x[10], 0xec5a, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x1002365f4 {
		return 0x1002365f4, true, nil
	}
	// 0x1002365f4: add x0, x9, x10
	s.x[0] = uint64((s.x[9] + s.x[10]))
	if runtime != nil && runtime.stopPC == 0x1002365f8 {
		return 0x1002365f8, true, nil
	}
	// 0x1002365f8: ldrb w8, [x0, w8, uxtw]
	value1607, err := s.load((s.x[0] + uint64(uint32(uint64(uint32(s.x[8]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1607))
	if runtime != nil && runtime.stopPC == 0x1002365fc {
		return 0x1002365fc, true, nil
	}
	// 0x1002365fc: lsr w9, w8, #3
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236600 {
		return 0x100236600, true, nil
	}
	// 0x100236600: orr w8, w9, w8, lsl #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) | shiftLeft(uint64(uint32(s.x[8])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236604 {
		return 0x100236604, true, nil
	}
	// 0x100236604: mov w12, #0x3b
	s.x[12] = uint64(uint32(uint64(0x3b)))
	if runtime != nil && runtime.stopPC == 0x100236608 {
		return 0x100236608, true, nil
	}
	// 0x100236608: and w9, w12, w8, lsr #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) & shiftRight(uint64(uint32(s.x[8])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x10023660c {
		return 0x10023660c, true, nil
	}
	// 0x10023660c: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236610 {
		return 0x100236610, true, nil
	}
	// 0x100236610: eor w8, w8, #0xffffff81
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0xffffff81))))
	if runtime != nil && runtime.stopPC == 0x100236614 {
		return 0x100236614, true, nil
	}
	// 0x100236614: ldur x9, [x29, #-0xb0]
	value1614, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1614)
	if runtime != nil && runtime.stopPC == 0x100236618 {
		return 0x100236618, true, nil
	}
	// 0x100236618: ldr x10, [sp, #0x428]
	value1615, err := s.load(addOffset(s.sp, 1064), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1615)
	if runtime != nil && runtime.stopPC == 0x10023661c {
		return 0x10023661c, true, nil
	}
	// 0x10023661c: strb w8, [x9, x10]
	if err := s.storeValue((s.x[9] + s.x[10]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236620 {
		return 0x100236620, true, nil
	}
	// 0x100236620: ldur x15, [x29, #-0xa0]
	value1617, err := s.load(addOffset(s.x[29], -160), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(value1617)
	if runtime != nil && runtime.stopPC == 0x100236624 {
		return 0x100236624, true, nil
	}
	// 0x100236624: ldr w8, [x15]
	value1618, err := s.load(s.x[15], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1618))
	if runtime != nil && runtime.stopPC == 0x100236628 {
		return 0x100236628, true, nil
	}
	// 0x100236628: eor w8, w8, w27
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x10023662c {
		return 0x10023662c, true, nil
	}
	// 0x10023662c: ldr w9, [x25, #0x588]
	value1620, err := s.load(addOffset(s.x[25], 1416), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1620))
	if runtime != nil && runtime.stopPC == 0x100236630 {
		return 0x100236630, true, nil
	}
	// 0x100236630: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100236634 {
		return 0x100236634, true, nil
	}
	// 0x100236634: mov w13, #0xd4cf
	s.x[13] = uint64(uint32(uint64(0xd4cf)))
	if runtime != nil && runtime.stopPC == 0x100236638 {
		return 0x100236638, true, nil
	}
	// 0x100236638: movk w13, #0xa5c, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xa5c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023663c {
		return 0x10023663c, true, nil
	}
	// 0x10023663c: mul w8, w8, w13
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236640 {
		return 0x100236640, true, nil
	}
	// 0x100236640: ldur x9, [x29, #-0xd8]
	value1625, err := s.load(addOffset(s.x[29], -216), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1625)
	if runtime != nil && runtime.stopPC == 0x100236644 {
		return 0x100236644, true, nil
	}
	// 0x100236644: add x9, x9, w8, uxtb
	s.x[9] = uint64((s.x[9] + uint64(uint8(uint64(uint32(s.x[8]))))))
	if runtime != nil && runtime.stopPC == 0x100236648 {
		return 0x100236648, true, nil
	}
	// 0x100236648: ldrb w9, [x9]
	value1627, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1627))
	if runtime != nil && runtime.stopPC == 0x10023664c {
		return 0x10023664c, true, nil
	}
	// 0x10023664c: mov w11, #0x84
	s.x[11] = uint64(uint32(uint64(0x84)))
	if runtime != nil && runtime.stopPC == 0x100236650 {
		return 0x100236650, true, nil
	}
	// 0x100236650: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236654 {
		return 0x100236654, true, nil
	}
	// 0x100236654: ldr x10, [sp, #0x3e8]
	value1630, err := s.load(addOffset(s.sp, 1000), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1630)
	if runtime != nil && runtime.stopPC == 0x100236658 {
		return 0x100236658, true, nil
	}
	// 0x100236658: ldrb w9, [x10, w9, uxtw]
	value1631, err := s.load((s.x[10] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1631))
	if runtime != nil && runtime.stopPC == 0x10023665c {
		return 0x10023665c, true, nil
	}
	// 0x10023665c: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236660 {
		return 0x100236660, true, nil
	}
	// 0x100236660: mov w9, #0x1292
	s.x[9] = uint64(uint32(uint64(0x1292)))
	if runtime != nil && runtime.stopPC == 0x100236664 {
		return 0x100236664, true, nil
	}
	// 0x100236664: add w8, w9, w8, uxtb
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint8(uint64(uint32(s.x[8])))))))
	if runtime != nil && runtime.stopPC == 0x100236668 {
		return 0x100236668, true, nil
	}
	// 0x100236668: ldr x8, [x26, w8, uxtw #3]
	value1635, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(value1635)
	if runtime != nil && runtime.stopPC == 0x10023666c {
		return 0x10023666c, true, nil
	}
	// 0x10023666c: lsr w9, w16, #0x18
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[16])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236670 {
		return 0x100236670, true, nil
	}
	// 0x100236670: eor x9, x9, #8
	s.x[9] = uint64((s.x[9] ^ uint64(0x8)))
	if runtime != nil && runtime.stopPC == 0x100236674 {
		return 0x100236674, true, nil
	}
	// 0x100236674: mov x10, #-0x26b2
	s.x[10] = uint64(uint64(0xffffffffffffd94e))
	if runtime != nil && runtime.stopPC == 0x100236678 {
		return 0x100236678, true, nil
	}
	// 0x100236678: movk x10, #0xf7ae, lsl #16
	s.x[10] = uint64(insertBits(s.x[10], 0xf7ae, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x10023667c {
		return 0x10023667c, true, nil
	}
	// 0x10023667c: add x20, x8, x10
	s.x[20] = uint64((s.x[8] + s.x[10]))
	if runtime != nil && runtime.stopPC == 0x100236680 {
		return 0x100236680, true, nil
	}
	// 0x100236680: ldrb w8, [x20, x9]
	value1641, err := s.load((s.x[20] + s.x[9]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1641))
	if runtime != nil && runtime.stopPC == 0x100236684 {
		return 0x100236684, true, nil
	}
	// 0x100236684: lsr w9, w8, #3
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236688 {
		return 0x100236688, true, nil
	}
	// 0x100236688: orr w8, w9, w8, lsl #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) | shiftLeft(uint64(uint32(s.x[8])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x10023668c {
		return 0x10023668c, true, nil
	}
	// 0x10023668c: mov w9, #0x68
	s.x[9] = uint64(uint32(uint64(0x68)))
	if runtime != nil && runtime.stopPC == 0x100236690 {
		return 0x100236690, true, nil
	}
	// 0x100236690: bfxil w9, w8, #5, #3
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), uint64(uint32(s.x[8]))>>5, 0, 3, 32)))
	if runtime != nil && runtime.stopPC == 0x100236694 {
		return 0x100236694, true, nil
	}
	// 0x100236694: and w10, w8, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236698 {
		return 0x100236698, true, nil
	}
	// 0x100236698: eor w8, w8, w10, lsr #1
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x10023669c {
		return 0x10023669c, true, nil
	}
	// 0x10023669c: mov w10, #0x43
	s.x[10] = uint64(uint32(uint64(0x43)))
	if runtime != nil && runtime.stopPC == 0x1002366a0 {
		return 0x1002366a0, true, nil
	}
	// 0x1002366a0: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002366a4 {
		return 0x1002366a4, true, nil
	}
	// 0x1002366a4: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002366a8 {
		return 0x1002366a8, true, nil
	}
	// 0x1002366a8: ldur x9, [x29, #-0xb0]
	value1651, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1651)
	if runtime != nil && runtime.stopPC == 0x1002366ac {
		return 0x1002366ac, true, nil
	}
	// 0x1002366ac: ldr x10, [sp, #0x448]
	value1652, err := s.load(addOffset(s.sp, 1096), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1652)
	if runtime != nil && runtime.stopPC == 0x1002366b0 {
		return 0x1002366b0, true, nil
	}
	// 0x1002366b0: strb w8, [x9, x10]
	if err := s.storeValue((s.x[9] + s.x[10]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002366b4 {
		return 0x1002366b4, true, nil
	}
	// 0x1002366b4: ldr w8, [sp, #0x4b8]
	value1654, err := s.load(addOffset(s.sp, 1208), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1654))
	if runtime != nil && runtime.stopPC == 0x1002366b8 {
		return 0x1002366b8, true, nil
	}
	// 0x1002366b8: ldr w9, [sp, #0x2d0]
	value1655, err := s.load(addOffset(s.sp, 720), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1655))
	if runtime != nil && runtime.stopPC == 0x1002366bc {
		return 0x1002366bc, true, nil
	}
	// 0x1002366bc: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002366c0 {
		return 0x1002366c0, true, nil
	}
	// 0x1002366c0: ldr w9, [sp, #0x2d8]
	value1657, err := s.load(addOffset(s.sp, 728), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1657))
	if runtime != nil && runtime.stopPC == 0x1002366c4 {
		return 0x1002366c4, true, nil
	}
	// 0x1002366c4: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1002366c8 {
		return 0x1002366c8, true, nil
	}
	// 0x1002366c8: add w9, w30, w17
	s.x[9] = uint64(uint32((uint64(uint32(s.x[30])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1002366cc {
		return 0x1002366cc, true, nil
	}
	// 0x1002366cc: ldur x30, [x29, #-0xb8]
	value1660, err := s.load(addOffset(s.x[29], -184), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[30] = uint64(value1660)
	if runtime != nil && runtime.stopPC == 0x1002366d0 {
		return 0x1002366d0, true, nil
	}
	// 0x1002366d0: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002366d4 {
		return 0x1002366d4, true, nil
	}
	// 0x1002366d4: and w9, w1, w2, lsl #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[1])) & shiftLeft(uint64(uint32(s.x[2])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1002366d8 {
		return 0x1002366d8, true, nil
	}
	// 0x1002366d8: eor w9, w9, w3
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x1002366dc {
		return 0x1002366dc, true, nil
	}
	// 0x1002366dc: eor w10, w2, w4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x1002366e0 {
		return 0x1002366e0, true, nil
	}
	// 0x1002366e0: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002366e4 {
		return 0x1002366e4, true, nil
	}
	// 0x1002366e4: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002366e8 {
		return 0x1002366e8, true, nil
	}
	// 0x1002366e8: and w9, w14, w9, lsr #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) & shiftRight(uint64(uint32(s.x[9])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1002366ec {
		return 0x1002366ec, true, nil
	}
	// 0x1002366ec: eor w8, w8, w19
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[19])))))
	if runtime != nil && runtime.stopPC == 0x1002366f0 {
		return 0x1002366f0, true, nil
	}
	// 0x1002366f0: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002366f4 {
		return 0x1002366f4, true, nil
	}
	// 0x1002366f4: eor w8, w8, w19, lsr #2
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[19])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x1002366f8 {
		return 0x1002366f8, true, nil
	}
	// 0x1002366f8: eor w8, w8, w19, lsr #3
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[19])), 3, 32))))
	if runtime != nil && runtime.stopPC == 0x1002366fc {
		return 0x1002366fc, true, nil
	}
	// 0x1002366fc: eor w21, w8, w21
	s.x[21] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[21])))))
	if runtime != nil && runtime.stopPC == 0x100236700 {
		return 0x100236700, true, nil
	}
	// 0x100236700: lsr w8, w21, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[21])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236704 {
		return 0x100236704, true, nil
	}
	// 0x100236704: eor x8, x8, #0x1c
	s.x[8] = uint64((s.x[8] ^ uint64(0x1c)))
	if runtime != nil && runtime.stopPC == 0x100236708 {
		return 0x100236708, true, nil
	}
	// 0x100236708: ldrb w8, [x20, x8]
	value1675, err := s.load((s.x[20] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1675))
	if runtime != nil && runtime.stopPC == 0x10023670c {
		return 0x10023670c, true, nil
	}
	// 0x10023670c: lsr w9, w8, #3
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236710 {
		return 0x100236710, true, nil
	}
	// 0x100236710: orr w8, w9, w8, lsl #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) | shiftLeft(uint64(uint32(s.x[8])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236714 {
		return 0x100236714, true, nil
	}
	// 0x100236714: and w9, w8, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236718 {
		return 0x100236718, true, nil
	}
	// 0x100236718: eor w8, w8, w9, lsr #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[9])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x10023671c {
		return 0x10023671c, true, nil
	}
	// 0x10023671c: eor w8, w8, w9, lsr #1
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[9])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236720 {
		return 0x100236720, true, nil
	}
	// 0x100236720: eor w8, w8, #0xaaaaaaaa
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0xaaaaaaaa))))
	if runtime != nil && runtime.stopPC == 0x100236724 {
		return 0x100236724, true, nil
	}
	// 0x100236724: ldp x9, x4, [x29, #-0xb0]
	value1682a, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	value1682b, err := s.load((addOffset(s.x[29], -176))+8, 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1682a)
	s.x[4] = uint64(value1682b)
	if runtime != nil && runtime.stopPC == 0x100236728 {
		return 0x100236728, true, nil
	}
	// 0x100236728: ldr x10, [sp, #0x440]
	value1683, err := s.load(addOffset(s.sp, 1088), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1683)
	if runtime != nil && runtime.stopPC == 0x10023672c {
		return 0x10023672c, true, nil
	}
	// 0x10023672c: strb w8, [x9, x10]
	if err := s.storeValue((s.x[9] + s.x[10]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236730 {
		return 0x100236730, true, nil
	}
	// 0x100236730: and w8, w21, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[21])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236734 {
		return 0x100236734, true, nil
	}
	// 0x100236734: eor x8, x8, x11
	s.x[8] = uint64((s.x[8] ^ s.x[11]))
	if runtime != nil && runtime.stopPC == 0x100236738 {
		return 0x100236738, true, nil
	}
	// 0x100236738: mov w2, #0x84
	s.x[2] = uint64(uint32(uint64(0x84)))
	if runtime != nil && runtime.stopPC == 0x10023673c {
		return 0x10023673c, true, nil
	}
	// 0x10023673c: ldrb w8, [x0, x8]
	value1688, err := s.load((s.x[0] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1688))
	if runtime != nil && runtime.stopPC == 0x100236740 {
		return 0x100236740, true, nil
	}
	// 0x100236740: lsr w9, w8, #3
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236744 {
		return 0x100236744, true, nil
	}
	// 0x100236744: orr w8, w9, w8, lsl #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) | shiftLeft(uint64(uint32(s.x[8])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236748 {
		return 0x100236748, true, nil
	}
	// 0x100236748: and w9, w12, w8, lsr #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) & shiftRight(uint64(uint32(s.x[8])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x10023674c {
		return 0x10023674c, true, nil
	}
	// 0x10023674c: mov w3, #0x3b
	s.x[3] = uint64(uint32(uint64(0x3b)))
	if runtime != nil && runtime.stopPC == 0x100236750 {
		return 0x100236750, true, nil
	}
	// 0x100236750: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236754 {
		return 0x100236754, true, nil
	}
	// 0x100236754: eor w8, w8, #0x70
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0x70))))
	if runtime != nil && runtime.stopPC == 0x100236758 {
		return 0x100236758, true, nil
	}
	// 0x100236758: ldur x9, [x29, #-0xb0]
	value1695, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1695)
	if runtime != nil && runtime.stopPC == 0x10023675c {
		return 0x10023675c, true, nil
	}
	// 0x10023675c: ldr x10, [sp, #0x430]
	value1696, err := s.load(addOffset(s.sp, 1072), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1696)
	if runtime != nil && runtime.stopPC == 0x100236760 {
		return 0x100236760, true, nil
	}
	// 0x100236760: strb w8, [x9, x10]
	if err := s.storeValue((s.x[9] + s.x[10]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236764 {
		return 0x100236764, true, nil
	}
	// 0x100236764: ldr w8, [x15]
	value1698, err := s.load(s.x[15], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1698))
	if runtime != nil && runtime.stopPC == 0x100236768 {
		return 0x100236768, true, nil
	}
	// 0x100236768: eor w8, w8, w27
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x10023676c {
		return 0x10023676c, true, nil
	}
	// 0x10023676c: ldr w9, [x25, #0x588]
	value1700, err := s.load(addOffset(s.x[25], 1416), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1700))
	if runtime != nil && runtime.stopPC == 0x100236770 {
		return 0x100236770, true, nil
	}
	// 0x100236770: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100236774 {
		return 0x100236774, true, nil
	}
	// 0x100236774: mul w8, w8, w13
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236778 {
		return 0x100236778, true, nil
	}
	// 0x100236778: mov w12, #0xd4cf
	s.x[12] = uint64(uint32(uint64(0xd4cf)))
	if runtime != nil && runtime.stopPC == 0x10023677c {
		return 0x10023677c, true, nil
	}
	// 0x10023677c: movk w12, #0xa5c, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xa5c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236780 {
		return 0x100236780, true, nil
	}
	// 0x100236780: mov x19, x30
	s.x[19] = uint64(s.x[30])
	if runtime != nil && runtime.stopPC == 0x100236784 {
		return 0x100236784, true, nil
	}
	// 0x100236784: add x9, x30, w8, uxtb
	s.x[9] = uint64((s.x[30] + uint64(uint8(uint64(uint32(s.x[8]))))))
	if runtime != nil && runtime.stopPC == 0x100236788 {
		return 0x100236788, true, nil
	}
	// 0x100236788: ldrb w9, [x9]
	value1707, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1707))
	if runtime != nil && runtime.stopPC == 0x10023678c {
		return 0x10023678c, true, nil
	}
	// 0x10023678c: mov w10, #0x5d
	s.x[10] = uint64(uint32(uint64(0x5d)))
	if runtime != nil && runtime.stopPC == 0x100236790 {
		return 0x100236790, true, nil
	}
	// 0x100236790: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236794 {
		return 0x100236794, true, nil
	}
	// 0x100236794: ldrb w9, [x4, w9, uxtw]
	value1710, err := s.load((s.x[4] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1710))
	if runtime != nil && runtime.stopPC == 0x100236798 {
		return 0x100236798, true, nil
	}
	// 0x100236798: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023679c {
		return 0x10023679c, true, nil
	}
	// 0x10023679c: ubfx w9, w16, #0x10, #8
	s.x[9] = uint64(uint32((uint64(uint32(s.x[16])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x1002367a0 {
		return 0x1002367a0, true, nil
	}
	// 0x1002367a0: mov w10, #0x37
	s.x[10] = uint64(uint32(uint64(0x37)))
	if runtime != nil && runtime.stopPC == 0x1002367a4 {
		return 0x1002367a4, true, nil
	}
	// 0x1002367a4: eor x9, x9, x10
	s.x[9] = uint64((s.x[9] ^ s.x[10]))
	if runtime != nil && runtime.stopPC == 0x1002367a8 {
		return 0x1002367a8, true, nil
	}
	// 0x1002367a8: and w8, w8, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002367ac {
		return 0x1002367ac, true, nil
	}
	// 0x1002367ac: add w8, w8, #0x3e
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(0x3e))))
	if runtime != nil && runtime.stopPC == 0x1002367b0 {
		return 0x1002367b0, true, nil
	}
	// 0x1002367b0: ldr x8, [x26, w8, uxtw #3]
	value1717, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(value1717)
	if runtime != nil && runtime.stopPC == 0x1002367b4 {
		return 0x1002367b4, true, nil
	}
	// 0x1002367b4: mov x10, #-0xb9aa
	s.x[10] = uint64(uint64(0xffffffffffff4656))
	if runtime != nil && runtime.stopPC == 0x1002367b8 {
		return 0x1002367b8, true, nil
	}
	// 0x1002367b8: movk x10, #0xf734, lsl #16
	s.x[10] = uint64(insertBits(s.x[10], 0xf734, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x1002367bc {
		return 0x1002367bc, true, nil
	}
	// 0x1002367bc: add x8, x8, x10
	s.x[8] = uint64((s.x[8] + s.x[10]))
	if runtime != nil && runtime.stopPC == 0x1002367c0 {
		return 0x1002367c0, true, nil
	}
	// 0x1002367c0: ldrb w9, [x8, x9]
	value1721, err := s.load((s.x[8] + s.x[9]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1721))
	if runtime != nil && runtime.stopPC == 0x1002367c4 {
		return 0x1002367c4, true, nil
	}
	// 0x1002367c4: mov w10, #0x54
	s.x[10] = uint64(uint32(uint64(0x54)))
	if runtime != nil && runtime.stopPC == 0x1002367c8 {
		return 0x1002367c8, true, nil
	}
	// 0x1002367c8: eor w10, w10, w16, lsr #16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[16])), 16, 32))))
	if runtime != nil && runtime.stopPC == 0x1002367cc {
		return 0x1002367cc, true, nil
	}
	// 0x1002367cc: add w10, w10, #0x68
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(0x68))))
	if runtime != nil && runtime.stopPC == 0x1002367d0 {
		return 0x1002367d0, true, nil
	}
	// 0x1002367d0: mov w11, #0x59
	s.x[11] = uint64(uint32(uint64(0x59)))
	if runtime != nil && runtime.stopPC == 0x1002367d4 {
		return 0x1002367d4, true, nil
	}
	// 0x1002367d4: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002367d8 {
		return 0x1002367d8, true, nil
	}
	// 0x1002367d8: eor w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1002367dc {
		return 0x1002367dc, true, nil
	}
	// 0x1002367dc: ldur x10, [x29, #-0xb0]
	value1728, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1728)
	if runtime != nil && runtime.stopPC == 0x1002367e0 {
		return 0x1002367e0, true, nil
	}
	// 0x1002367e0: ldr x11, [sp, #0x438]
	value1729, err := s.load(addOffset(s.sp, 1080), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value1729)
	if runtime != nil && runtime.stopPC == 0x1002367e4 {
		return 0x1002367e4, true, nil
	}
	// 0x1002367e4: strb w9, [x10, x11]
	if err := s.storeValue((s.x[10] + s.x[11]), uint64(uint32(s.x[9])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002367e8 {
		return 0x1002367e8, true, nil
	}
	// 0x1002367e8: lsr w9, w21, #4
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[21])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002367ec {
		return 0x1002367ec, true, nil
	}
	// 0x1002367ec: and w9, w9, #0xf0
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xf0))))
	if runtime != nil && runtime.stopPC == 0x1002367f0 {
		return 0x1002367f0, true, nil
	}
	// 0x1002367f0: ubfx w10, w21, #0xc, #4
	s.x[10] = uint64(uint32((uint64(uint32(s.x[21])) >> 12) & bitMask(4, 32)))
	if runtime != nil && runtime.stopPC == 0x1002367f4 {
		return 0x1002367f4, true, nil
	}
	// 0x1002367f4: mov w11, #0xf6
	s.x[11] = uint64(uint32(uint64(0xf6)))
	if runtime != nil && runtime.stopPC == 0x1002367f8 {
		return 0x1002367f8, true, nil
	}
	// 0x1002367f8: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002367fc {
		return 0x1002367fc, true, nil
	}
	// 0x1002367fc: mov w11, #0x5f
	s.x[11] = uint64(uint32(uint64(0x5f)))
	if runtime != nil && runtime.stopPC == 0x100236800 {
		return 0x100236800, true, nil
	}
	// 0x100236800: eor w9, w9, w11
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236804 {
		return 0x100236804, true, nil
	}
	// 0x100236804: and w9, w10, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236808 {
		return 0x100236808, true, nil
	}
	// 0x100236808: eor w9, w9, #0xfffffffd
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(0xfffffffd))))
	if runtime != nil && runtime.stopPC == 0x10023680c {
		return 0x10023680c, true, nil
	}
	// 0x10023680c: lsl w10, w21, #0x14
	s.x[10] = uint64(uint32(shiftLeft(uint64(uint32(s.x[21])), int(uint64(0x14)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236810 {
		return 0x100236810, true, nil
	}
	// 0x100236810: eor w10, w10, #0x80000000
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(0x80000000))))
	if runtime != nil && runtime.stopPC == 0x100236814 {
		return 0x100236814, true, nil
	}
	// 0x100236814: extr w9, w9, w10, #0x1c
	s.x[9] = uint64(uint32(extractPair(uint64(uint32(s.x[9])), uint64(uint32(s.x[10])), 28, 32)))
	if runtime != nil && runtime.stopPC == 0x100236818 {
		return 0x100236818, true, nil
	}
	// 0x100236818: ldr w10, [x25, #0x588]
	value1743, err := s.load(addOffset(s.x[25], 1416), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1743))
	if runtime != nil && runtime.stopPC == 0x10023681c {
		return 0x10023681c, true, nil
	}
	// 0x10023681c: ldr w11, [x15]
	value1744, err := s.load(s.x[15], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value1744))
	if runtime != nil && runtime.stopPC == 0x100236820 {
		return 0x100236820, true, nil
	}
	// 0x100236820: sub w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) - uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236824 {
		return 0x100236824, true, nil
	}
	// 0x100236824: eor w10, w10, w27
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[27])))))
	if runtime != nil && runtime.stopPC == 0x100236828 {
		return 0x100236828, true, nil
	}
	// 0x100236828: mov x30, x7
	s.x[30] = uint64(s.x[7])
	if runtime != nil && runtime.stopPC == 0x10023682c {
		return 0x10023682c, true, nil
	}
	// 0x10023682c: mul w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) * uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236830 {
		return 0x100236830, true, nil
	}
	// 0x100236830: add x11, x7, w10, uxtb
	s.x[11] = uint64((s.x[7] + uint64(uint8(uint64(uint32(s.x[10]))))))
	if runtime != nil && runtime.stopPC == 0x100236834 {
		return 0x100236834, true, nil
	}
	// 0x100236834: ldrb w11, [x11]
	value1750, err := s.load(s.x[11], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value1750))
	if runtime != nil && runtime.stopPC == 0x100236838 {
		return 0x100236838, true, nil
	}
	// 0x100236838: mov w12, #0xfb
	s.x[12] = uint64(uint32(uint64(0xfb)))
	if runtime != nil && runtime.stopPC == 0x10023683c {
		return 0x10023683c, true, nil
	}
	// 0x10023683c: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236840 {
		return 0x100236840, true, nil
	}
	// 0x100236840: ldrb w11, [x28, w11, uxtw]
	value1753, err := s.load((s.x[28] + uint64(uint32(uint64(uint32(s.x[11]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value1753))
	if runtime != nil && runtime.stopPC == 0x100236844 {
		return 0x100236844, true, nil
	}
	// 0x100236844: ldur x28, [x29, #-0xb0]
	value1754, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[28] = uint64(value1754)
	if runtime != nil && runtime.stopPC == 0x100236848 {
		return 0x100236848, true, nil
	}
	// 0x100236848: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023684c {
		return 0x10023684c, true, nil
	}
	// 0x10023684c: and x9, x9, #0xff
	s.x[9] = uint64((s.x[9] & uint64(0xff)))
	if runtime != nil && runtime.stopPC == 0x100236850 {
		return 0x100236850, true, nil
	}
	// 0x100236850: mov w11, #0x2a
	s.x[11] = uint64(uint32(uint64(0x2a)))
	if runtime != nil && runtime.stopPC == 0x100236854 {
		return 0x100236854, true, nil
	}
	// 0x100236854: eor x9, x9, x11
	s.x[9] = uint64((s.x[9] ^ s.x[11]))
	if runtime != nil && runtime.stopPC == 0x100236858 {
		return 0x100236858, true, nil
	}
	// 0x100236858: and w10, w10, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x10023685c {
		return 0x10023685c, true, nil
	}
	// 0x10023685c: add w10, w10, #0xf1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(0xf1))))
	if runtime != nil && runtime.stopPC == 0x100236860 {
		return 0x100236860, true, nil
	}
	// 0x100236860: ldr x10, [x26, w10, uxtw #3]
	value1761, err := s.load((s.x[26] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1761)
	if runtime != nil && runtime.stopPC == 0x100236864 {
		return 0x100236864, true, nil
	}
	// 0x100236864: mov x11, #-0xf277
	s.x[11] = uint64(uint64(0xffffffffffff0d89))
	if runtime != nil && runtime.stopPC == 0x100236868 {
		return 0x100236868, true, nil
	}
	// 0x100236868: movk x11, #0xa9f8, lsl #16
	s.x[11] = uint64(insertBits(s.x[11], 0xa9f8, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x10023686c {
		return 0x10023686c, true, nil
	}
	// 0x10023686c: add x12, x10, x11
	s.x[12] = uint64((s.x[10] + s.x[11]))
	if runtime != nil && runtime.stopPC == 0x100236870 {
		return 0x100236870, true, nil
	}
	// 0x100236870: ldrb w9, [x12, x9]
	value1765, err := s.load((s.x[12] + s.x[9]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1765))
	if runtime != nil && runtime.stopPC == 0x100236874 {
		return 0x100236874, true, nil
	}
	// 0x100236874: eor w9, w9, #0x3c
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(0x3c))))
	if runtime != nil && runtime.stopPC == 0x100236878 {
		return 0x100236878, true, nil
	}
	// 0x100236878: ldr x10, [sp, #0x478]
	value1767, err := s.load(addOffset(s.sp, 1144), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1767)
	if runtime != nil && runtime.stopPC == 0x10023687c {
		return 0x10023687c, true, nil
	}
	// 0x10023687c: strb w9, [x28, x10]
	if err := s.storeValue((s.x[28] + s.x[10]), uint64(uint32(s.x[9])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236880 {
		return 0x100236880, true, nil
	}
	// 0x100236880: mov w9, #0x9a46
	s.x[9] = uint64(uint32(uint64(0x9a46)))
	if runtime != nil && runtime.stopPC == 0x100236884 {
		return 0x100236884, true, nil
	}
	// 0x100236884: movk w9, #0x7bd6, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7bd6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236888 {
		return 0x100236888, true, nil
	}
	// 0x100236888: eor w9, w5, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023688c {
		return 0x10023688c, true, nil
	}
	// 0x10023688c: mov w10, #0x65ba
	s.x[10] = uint64(uint32(uint64(0x65ba)))
	if runtime != nil && runtime.stopPC == 0x100236890 {
		return 0x100236890, true, nil
	}
	// 0x100236890: movk w10, #0x8429, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x8429, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236894 {
		return 0x100236894, true, nil
	}
	// 0x100236894: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236898 {
		return 0x100236898, true, nil
	}
	// 0x100236898: mov w10, #0x62c3
	s.x[10] = uint64(uint32(uint64(0x62c3)))
	if runtime != nil && runtime.stopPC == 0x10023689c {
		return 0x10023689c, true, nil
	}
	// 0x10023689c: movk w10, #0x3f8a, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x3f8a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368a0 {
		return 0x1002368a0, true, nil
	}
	// 0x1002368a0: eor w10, w5, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002368a4 {
		return 0x1002368a4, true, nil
	}
	// 0x1002368a4: mov w11, #0x9d3d
	s.x[11] = uint64(uint32(uint64(0x9d3d)))
	if runtime != nil && runtime.stopPC == 0x1002368a8 {
		return 0x1002368a8, true, nil
	}
	// 0x1002368a8: movk w11, #0xc075, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xc075, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368ac {
		return 0x1002368ac, true, nil
	}
	// 0x1002368ac: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002368b0 {
		return 0x1002368b0, true, nil
	}
	// 0x1002368b0: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002368b4 {
		return 0x1002368b4, true, nil
	}
	// 0x1002368b4: mov w10, #0xebf
	s.x[10] = uint64(uint32(uint64(0xebf)))
	if runtime != nil && runtime.stopPC == 0x1002368b8 {
		return 0x1002368b8, true, nil
	}
	// 0x1002368b8: movk w10, #0x8197, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x8197, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368bc {
		return 0x1002368bc, true, nil
	}
	// 0x1002368bc: eor w10, w5, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[5])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002368c0 {
		return 0x1002368c0, true, nil
	}
	// 0x1002368c0: mov w11, #0xf141
	s.x[11] = uint64(uint32(uint64(0xf141)))
	if runtime != nil && runtime.stopPC == 0x1002368c4 {
		return 0x1002368c4, true, nil
	}
	// 0x1002368c4: movk w11, #0x7e68, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x7e68, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368c8 {
		return 0x1002368c8, true, nil
	}
	// 0x1002368c8: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002368cc {
		return 0x1002368cc, true, nil
	}
	// 0x1002368cc: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002368d0 {
		return 0x1002368d0, true, nil
	}
	// 0x1002368d0: mov w10, #0x8e19
	s.x[10] = uint64(uint32(uint64(0x8e19)))
	if runtime != nil && runtime.stopPC == 0x1002368d4 {
		return 0x1002368d4, true, nil
	}
	// 0x1002368d4: movk w10, #0x2ee0, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x2ee0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368d8 {
		return 0x1002368d8, true, nil
	}
	// 0x1002368d8: eor w10, w23, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002368dc {
		return 0x1002368dc, true, nil
	}
	// 0x1002368dc: mov w11, #0x71e7
	s.x[11] = uint64(uint32(uint64(0x71e7)))
	if runtime != nil && runtime.stopPC == 0x1002368e0 {
		return 0x1002368e0, true, nil
	}
	// 0x1002368e0: movk w11, #0xd11f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd11f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368e4 {
		return 0x1002368e4, true, nil
	}
	// 0x1002368e4: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002368e8 {
		return 0x1002368e8, true, nil
	}
	// 0x1002368e8: mov w11, #0xa300
	s.x[11] = uint64(uint32(uint64(0xa300)))
	if runtime != nil && runtime.stopPC == 0x1002368ec {
		return 0x1002368ec, true, nil
	}
	// 0x1002368ec: movk w11, #0x745f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x745f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368f0 {
		return 0x1002368f0, true, nil
	}
	// 0x1002368f0: eor w11, w23, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1002368f4 {
		return 0x1002368f4, true, nil
	}
	// 0x1002368f4: mov w13, #0x5d00
	s.x[13] = uint64(uint32(uint64(0x5d00)))
	if runtime != nil && runtime.stopPC == 0x1002368f8 {
		return 0x1002368f8, true, nil
	}
	// 0x1002368f8: movk w13, #0x8ba0, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x8ba0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1002368fc {
		return 0x1002368fc, true, nil
	}
	// 0x1002368fc: add w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236900 {
		return 0x100236900, true, nil
	}
	// 0x100236900: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236904 {
		return 0x100236904, true, nil
	}
	// 0x100236904: mov w11, #0xdb23
	s.x[11] = uint64(uint32(uint64(0xdb23)))
	if runtime != nil && runtime.stopPC == 0x100236908 {
		return 0x100236908, true, nil
	}
	// 0x100236908: movk w11, #0x9f74, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9f74, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10023690c {
		return 0x10023690c, true, nil
	}
	// 0x10023690c: eor w11, w23, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236910 {
		return 0x100236910, true, nil
	}
	// 0x100236910: mov w13, #0x24dd
	s.x[13] = uint64(uint32(uint64(0x24dd)))
	if runtime != nil && runtime.stopPC == 0x100236914 {
		return 0x100236914, true, nil
	}
	// 0x100236914: movk w13, #0x608b, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x608b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236918 {
		return 0x100236918, true, nil
	}
	// 0x100236918: add w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x10023691c {
		return 0x10023691c, true, nil
	}
	// 0x10023691c: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236920 {
		return 0x100236920, true, nil
	}
	// 0x100236920: ldr w11, [sp, #0x3c4]
	value1809, err := s.load(addOffset(s.sp, 964), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value1809))
	if runtime != nil && runtime.stopPC == 0x100236924 {
		return 0x100236924, true, nil
	}
	// 0x100236924: ldr w13, [sp, #0x4bc]
	value1810, err := s.load(addOffset(s.sp, 1212), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1810))
	if runtime != nil && runtime.stopPC == 0x100236928 {
		return 0x100236928, true, nil
	}
	// 0x100236928: eor w11, w13, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x10023692c {
		return 0x10023692c, true, nil
	}
	// 0x10023692c: ldr w13, [sp, #0x424]
	value1812, err := s.load(addOffset(s.sp, 1060), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1812))
	if runtime != nil && runtime.stopPC == 0x100236930 {
		return 0x100236930, true, nil
	}
	// 0x100236930: add w13, w13, w17
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236934 {
		return 0x100236934, true, nil
	}
	// 0x100236934: eor w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236938 {
		return 0x100236938, true, nil
	}
	// 0x100236938: eor w11, w11, w6
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x10023693c {
		return 0x10023693c, true, nil
	}
	// 0x10023693c: eor w11, w11, w6, lsr #2
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ shiftRight(uint64(uint32(s.x[6])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100236940 {
		return 0x100236940, true, nil
	}
	// 0x100236940: eor w11, w11, w6, lsr #3
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ shiftRight(uint64(uint32(s.x[6])), 3, 32))))
	if runtime != nil && runtime.stopPC == 0x100236944 {
		return 0x100236944, true, nil
	}
	// 0x100236944: eor w11, w11, w24
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[24])))))
	if runtime != nil && runtime.stopPC == 0x100236948 {
		return 0x100236948, true, nil
	}
	// 0x100236948: mov w1, #0x5a8
	s.x[1] = uint64(uint32(uint64(0x5a8)))
	if runtime != nil && runtime.stopPC == 0x10023694c {
		return 0x10023694c, true, nil
	}
	// 0x10023694c: movk w1, #0x499d, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x499d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236950 {
		return 0x100236950, true, nil
	}
	// 0x100236950: add w9, w9, w1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236954 {
		return 0x100236954, true, nil
	}
	// 0x100236954: eor w11, w11, w9
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236958 {
		return 0x100236958, true, nil
	}
	// 0x100236958: and w9, w14, w9, lsr #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[14])) & shiftRight(uint64(uint32(s.x[9])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x10023695c {
		return 0x10023695c, true, nil
	}
	// 0x10023695c: eor w11, w11, w9
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236960 {
		return 0x100236960, true, nil
	}
	// 0x100236960: ldr w9, [sp, #0x3d4]
	value1825, err := s.load(addOffset(s.sp, 980), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value1825))
	if runtime != nil && runtime.stopPC == 0x100236964 {
		return 0x100236964, true, nil
	}
	// 0x100236964: ldr w13, [sp, #0x4c0]
	value1826, err := s.load(addOffset(s.sp, 1216), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1826))
	if runtime != nil && runtime.stopPC == 0x100236968 {
		return 0x100236968, true, nil
	}
	// 0x100236968: eor w9, w13, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x10023696c {
		return 0x10023696c, true, nil
	}
	// 0x10023696c: ldr w13, [sp, #0x410]
	value1828, err := s.load(addOffset(s.sp, 1040), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1828))
	if runtime != nil && runtime.stopPC == 0x100236970 {
		return 0x100236970, true, nil
	}
	// 0x100236970: add w13, w13, w17
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100236974 {
		return 0x100236974, true, nil
	}
	// 0x100236974: eor w9, w9, w13
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236978 {
		return 0x100236978, true, nil
	}
	// 0x100236978: eor w9, w9, w22
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[22])))))
	if runtime != nil && runtime.stopPC == 0x10023697c {
		return 0x10023697c, true, nil
	}
	// 0x10023697c: eor w9, w9, w22, lsr #2
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ shiftRight(uint64(uint32(s.x[22])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100236980 {
		return 0x100236980, true, nil
	}
	// 0x100236980: eor w9, w9, w22, lsr #3
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ shiftRight(uint64(uint32(s.x[22])), 3, 32))))
	if runtime != nil && runtime.stopPC == 0x100236984 {
		return 0x100236984, true, nil
	}
	// 0x100236984: add w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100236988 {
		return 0x100236988, true, nil
	}
	// 0x100236988: ldr w13, [sp, #0x408]
	value1835, err := s.load(addOffset(s.sp, 1032), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1835))
	if runtime != nil && runtime.stopPC == 0x10023698c {
		return 0x10023698c, true, nil
	}
	// 0x10023698c: eor w13, w13, w10
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236990 {
		return 0x100236990, true, nil
	}
	// 0x100236990: and w10, w14, w10, lsr #1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[14])) & shiftRight(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236994 {
		return 0x100236994, true, nil
	}
	// 0x100236994: eor w10, w13, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236998 {
		return 0x100236998, true, nil
	}
	// 0x100236998: ubfx w13, w11, #8, #8
	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x10023699c {
		return 0x10023699c, true, nil
	}
	// 0x10023699c: eor x13, x13, #0xc
	s.x[13] = uint64((s.x[13] ^ uint64(0xc)))
	if runtime != nil && runtime.stopPC == 0x1002369a0 {
		return 0x1002369a0, true, nil
	}
	// 0x1002369a0: ldrb w13, [x12, x13]
	value1841, err := s.load((s.x[12] + s.x[13]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1841))
	if runtime != nil && runtime.stopPC == 0x1002369a4 {
		return 0x1002369a4, true, nil
	}
	// 0x1002369a4: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002369a8 {
		return 0x1002369a8, true, nil
	}
	// 0x1002369a8: eor w10, w13, #0x3c
	s.x[10] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(0x3c))))
	if runtime != nil && runtime.stopPC == 0x1002369ac {
		return 0x1002369ac, true, nil
	}
	// 0x1002369ac: ldr x13, [sp, #0x488]
	value1844, err := s.load(addOffset(s.sp, 1160), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value1844)
	if runtime != nil && runtime.stopPC == 0x1002369b0 {
		return 0x1002369b0, true, nil
	}
	// 0x1002369b0: strb w10, [x28, x13]
	if err := s.storeValue((s.x[28] + s.x[13]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002369b4 {
		return 0x1002369b4, true, nil
	}
	// 0x1002369b4: and w10, w9, #0xff
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1002369b8 {
		return 0x1002369b8, true, nil
	}
	// 0x1002369b8: mov w13, #0xcf
	s.x[13] = uint64(uint32(uint64(0xcf)))
	if runtime != nil && runtime.stopPC == 0x1002369bc {
		return 0x1002369bc, true, nil
	}
	// 0x1002369bc: eor x10, x10, x13
	s.x[10] = uint64((s.x[10] ^ s.x[13]))
	if runtime != nil && runtime.stopPC == 0x1002369c0 {
		return 0x1002369c0, true, nil
	}
	// 0x1002369c0: ldrb w10, [x0, x10]
	value1849, err := s.load((s.x[0] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1849))
	if runtime != nil && runtime.stopPC == 0x1002369c4 {
		return 0x1002369c4, true, nil
	}
	// 0x1002369c4: lsr w13, w10, #3
	s.x[13] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002369c8 {
		return 0x1002369c8, true, nil
	}
	// 0x1002369c8: orr w10, w13, w10, lsl #5
	s.x[10] = uint64(uint32((uint64(uint32(s.x[13])) | shiftLeft(uint64(uint32(s.x[10])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x1002369cc {
		return 0x1002369cc, true, nil
	}
	// 0x1002369cc: and w13, w3, w10, lsr #1
	s.x[13] = uint64(uint32((uint64(uint32(s.x[3])) & shiftRight(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1002369d0 {
		return 0x1002369d0, true, nil
	}
	// 0x1002369d0: mov w14, #0x23
	s.x[14] = uint64(uint32(uint64(0x23)))
	if runtime != nil && runtime.stopPC == 0x1002369d4 {
		return 0x1002369d4, true, nil
	}
	// 0x1002369d4: eor w13, w13, w14
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002369d8 {
		return 0x1002369d8, true, nil
	}
	// 0x1002369d8: mov w14, #0x55
	s.x[14] = uint64(uint32(uint64(0x55)))
	if runtime != nil && runtime.stopPC == 0x1002369dc {
		return 0x1002369dc, true, nil
	}
	// 0x1002369dc: eor w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1002369e0 {
		return 0x1002369e0, true, nil
	}
	// 0x1002369e0: and w14, w13, w10
	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002369e4 {
		return 0x1002369e4, true, nil
	}
	// 0x1002369e4: add w10, w13, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1002369e8 {
		return 0x1002369e8, true, nil
	}
	// 0x1002369e8: sub w10, w10, w14, lsl #1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) - shiftLeft(uint64(uint32(s.x[14])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1002369ec {
		return 0x1002369ec, true, nil
	}
	// 0x1002369ec: ldr x13, [sp, #0x498]
	value1860, err := s.load(addOffset(s.sp, 1176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value1860)
	if runtime != nil && runtime.stopPC == 0x1002369f0 {
		return 0x1002369f0, true, nil
	}
	// 0x1002369f0: strb w10, [x28, x13]
	if err := s.storeValue((s.x[28] + s.x[13]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1002369f4 {
		return 0x1002369f4, true, nil
	}
	// 0x1002369f4: lsr w10, w11, #0x18
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x1002369f8 {
		return 0x1002369f8, true, nil
	}
	// 0x1002369f8: mov w13, #0x9d
	s.x[13] = uint64(uint32(uint64(0x9d)))
	if runtime != nil && runtime.stopPC == 0x1002369fc {
		return 0x1002369fc, true, nil
	}
	// 0x1002369fc: eor x10, x10, x13
	s.x[10] = uint64((s.x[10] ^ s.x[13]))
	if runtime != nil && runtime.stopPC == 0x100236a00 {
		return 0x100236a00, true, nil
	}
	// 0x100236a00: ldrb w10, [x20, x10]
	value1865, err := s.load((s.x[20] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1865))
	if runtime != nil && runtime.stopPC == 0x100236a04 {
		return 0x100236a04, true, nil
	}
	// 0x100236a04: lsr w13, w10, #3
	s.x[13] = uint64(uint32(shiftRight(uint64(uint32(s.x[10])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236a08 {
		return 0x100236a08, true, nil
	}
	// 0x100236a08: orr w10, w13, w10, lsl #5
	s.x[10] = uint64(uint32((uint64(uint32(s.x[13])) | shiftLeft(uint64(uint32(s.x[10])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236a0c {
		return 0x100236a0c, true, nil
	}
	// 0x100236a0c: and w13, w10, #0xff
	s.x[13] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236a10 {
		return 0x100236a10, true, nil
	}
	// 0x100236a10: eor w10, w10, w13, lsr #5
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[13])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236a14 {
		return 0x100236a14, true, nil
	}
	// 0x100236a14: eor w10, w10, w13, lsr #1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[13])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236a18 {
		return 0x100236a18, true, nil
	}
	// 0x100236a18: mov w13, #0xca
	s.x[13] = uint64(uint32(uint64(0xca)))
	if runtime != nil && runtime.stopPC == 0x100236a1c {
		return 0x100236a1c, true, nil
	}
	// 0x100236a1c: eor w10, w10, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236a20 {
		return 0x100236a20, true, nil
	}
	// 0x100236a20: ldr x13, [sp, #0x480]
	value1873, err := s.load(addOffset(s.sp, 1152), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value1873)
	if runtime != nil && runtime.stopPC == 0x100236a24 {
		return 0x100236a24, true, nil
	}
	// 0x100236a24: strb w10, [x28, x13]
	if err := s.storeValue((s.x[28] + s.x[13]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236a28 {
		return 0x100236a28, true, nil
	}
	// 0x100236a28: lsr w10, w11, #0x10
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[11])), int(uint64(0x10)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236a2c {
		return 0x100236a2c, true, nil
	}
	// 0x100236a2c: ubfx w13, w11, #0x10, #8
	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236a30 {
		return 0x100236a30, true, nil
	}
	// 0x100236a30: mov w14, #0x64
	s.x[14] = uint64(uint32(uint64(0x64)))
	if runtime != nil && runtime.stopPC == 0x100236a34 {
		return 0x100236a34, true, nil
	}
	// 0x100236a34: eor x13, x13, x14
	s.x[13] = uint64((s.x[13] ^ s.x[14]))
	if runtime != nil && runtime.stopPC == 0x100236a38 {
		return 0x100236a38, true, nil
	}
	// 0x100236a38: ldrb w13, [x8, x13]
	value1879, err := s.load((s.x[8] + s.x[13]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value1879))
	if runtime != nil && runtime.stopPC == 0x100236a3c {
		return 0x100236a3c, true, nil
	}
	// 0x100236a3c: eor w10, w10, #7
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(0x7))))
	if runtime != nil && runtime.stopPC == 0x100236a40 {
		return 0x100236a40, true, nil
	}
	// 0x100236a40: add w10, w10, #0x68
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(0x68))))
	if runtime != nil && runtime.stopPC == 0x100236a44 {
		return 0x100236a44, true, nil
	}
	// 0x100236a44: mov w14, #0xab
	s.x[14] = uint64(uint32(uint64(0xab)))
	if runtime != nil && runtime.stopPC == 0x100236a48 {
		return 0x100236a48, true, nil
	}
	// 0x100236a48: eor w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236a4c {
		return 0x100236a4c, true, nil
	}
	// 0x100236a4c: eor w10, w10, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236a50 {
		return 0x100236a50, true, nil
	}
	// 0x100236a50: ldr x13, [sp, #0x468]
	value1885, err := s.load(addOffset(s.sp, 1128), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value1885)
	if runtime != nil && runtime.stopPC == 0x100236a54 {
		return 0x100236a54, true, nil
	}
	// 0x100236a54: strb w10, [x28, x13]
	if err := s.storeValue((s.x[28] + s.x[13]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236a58 {
		return 0x100236a58, true, nil
	}
	// 0x100236a58: ubfx w10, w16, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[16])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236a5c {
		return 0x100236a5c, true, nil
	}
	// 0x100236a5c: mov w13, #0xa
	s.x[13] = uint64(uint32(uint64(0xa)))
	if runtime != nil && runtime.stopPC == 0x100236a60 {
		return 0x100236a60, true, nil
	}
	// 0x100236a60: eor x10, x10, x13
	s.x[10] = uint64((s.x[10] ^ s.x[13]))
	if runtime != nil && runtime.stopPC == 0x100236a64 {
		return 0x100236a64, true, nil
	}
	// 0x100236a64: ldrb w10, [x12, x10]
	value1890, err := s.load((s.x[12] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1890))
	if runtime != nil && runtime.stopPC == 0x100236a68 {
		return 0x100236a68, true, nil
	}
	// 0x100236a68: eor w10, w10, #0x7c
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(0x7c))))
	if runtime != nil && runtime.stopPC == 0x100236a6c {
		return 0x100236a6c, true, nil
	}
	// 0x100236a6c: ldr x13, [sp, #0x458]
	value1892, err := s.load(addOffset(s.sp, 1112), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value1892)
	if runtime != nil && runtime.stopPC == 0x100236a70 {
		return 0x100236a70, true, nil
	}
	// 0x100236a70: strb w10, [x28, x13]
	if err := s.storeValue((s.x[28] + s.x[13]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236a74 {
		return 0x100236a74, true, nil
	}
	// 0x100236a74: ubfx w10, w9, #8, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) >> 8) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236a78 {
		return 0x100236a78, true, nil
	}
	// 0x100236a78: mov w13, #0x32
	s.x[13] = uint64(uint32(uint64(0x32)))
	if runtime != nil && runtime.stopPC == 0x100236a7c {
		return 0x100236a7c, true, nil
	}
	// 0x100236a7c: eor x10, x10, x13
	s.x[10] = uint64((s.x[10] ^ s.x[13]))
	if runtime != nil && runtime.stopPC == 0x100236a80 {
		return 0x100236a80, true, nil
	}
	// 0x100236a80: ldrb w10, [x12, x10]
	value1897, err := s.load((s.x[12] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1897))
	if runtime != nil && runtime.stopPC == 0x100236a84 {
		return 0x100236a84, true, nil
	}
	// 0x100236a84: eor w10, w10, #0x88888888
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(0x88888888))))
	if runtime != nil && runtime.stopPC == 0x100236a88 {
		return 0x100236a88, true, nil
	}
	// 0x100236a88: ldr x12, [sp, #0x470]
	value1899, err := s.load(addOffset(s.sp, 1136), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(value1899)
	if runtime != nil && runtime.stopPC == 0x100236a8c {
		return 0x100236a8c, true, nil
	}
	// 0x100236a8c: strb w10, [x28, x12]
	if err := s.storeValue((s.x[28] + s.x[12]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236a90 {
		return 0x100236a90, true, nil
	}
	// 0x100236a90: ubfx w10, w9, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236a94 {
		return 0x100236a94, true, nil
	}
	// 0x100236a94: mov w12, #0xc1
	s.x[12] = uint64(uint32(uint64(0xc1)))
	if runtime != nil && runtime.stopPC == 0x100236a98 {
		return 0x100236a98, true, nil
	}
	// 0x100236a98: eor x10, x10, x12
	s.x[10] = uint64((s.x[10] ^ s.x[12]))
	if runtime != nil && runtime.stopPC == 0x100236a9c {
		return 0x100236a9c, true, nil
	}
	// 0x100236a9c: ldrb w10, [x8, x10]
	value1904, err := s.load((s.x[8] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1904))
	if runtime != nil && runtime.stopPC == 0x100236aa0 {
		return 0x100236aa0, true, nil
	}
	// 0x100236aa0: mov w12, #0xa2
	s.x[12] = uint64(uint32(uint64(0xa2)))
	if runtime != nil && runtime.stopPC == 0x100236aa4 {
		return 0x100236aa4, true, nil
	}
	// 0x100236aa4: eor w12, w12, w9, lsr #16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ shiftRight(uint64(uint32(s.x[9])), 16, 32))))
	if runtime != nil && runtime.stopPC == 0x100236aa8 {
		return 0x100236aa8, true, nil
	}
	// 0x100236aa8: add w12, w12, #0x68
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0x68))))
	if runtime != nil && runtime.stopPC == 0x100236aac {
		return 0x100236aac, true, nil
	}
	// 0x100236aac: mov w13, #0xf4
	s.x[13] = uint64(uint32(uint64(0xf4)))
	if runtime != nil && runtime.stopPC == 0x100236ab0 {
		return 0x100236ab0, true, nil
	}
	// 0x100236ab0: eor w12, w12, w13
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236ab4 {
		return 0x100236ab4, true, nil
	}
	// 0x100236ab4: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236ab8 {
		return 0x100236ab8, true, nil
	}
	// 0x100236ab8: ldr x12, [sp, #0x4a0]
	value1911, err := s.load(addOffset(s.sp, 1184), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(value1911)
	if runtime != nil && runtime.stopPC == 0x100236abc {
		return 0x100236abc, true, nil
	}
	// 0x100236abc: strb w10, [x28, x12]
	if err := s.storeValue((s.x[28] + s.x[12]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236ac0 {
		return 0x100236ac0, true, nil
	}
	// 0x100236ac0: ubfx w10, w21, #0x10, #8
	s.x[10] = uint64(uint32((uint64(uint32(s.x[21])) >> 16) & bitMask(8, 32)))
	if runtime != nil && runtime.stopPC == 0x100236ac4 {
		return 0x100236ac4, true, nil
	}
	// 0x100236ac4: mov w12, #0x6a
	s.x[12] = uint64(uint32(uint64(0x6a)))
	if runtime != nil && runtime.stopPC == 0x100236ac8 {
		return 0x100236ac8, true, nil
	}
	// 0x100236ac8: eor x10, x10, x12
	s.x[10] = uint64((s.x[10] ^ s.x[12]))
	if runtime != nil && runtime.stopPC == 0x100236acc {
		return 0x100236acc, true, nil
	}
	// 0x100236acc: ldrb w8, [x8, x10]
	value1916, err := s.load((s.x[8] + s.x[10]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1916))
	if runtime != nil && runtime.stopPC == 0x100236ad0 {
		return 0x100236ad0, true, nil
	}
	// 0x100236ad0: mov w10, #9
	s.x[10] = uint64(uint32(uint64(0x9)))
	if runtime != nil && runtime.stopPC == 0x100236ad4 {
		return 0x100236ad4, true, nil
	}
	// 0x100236ad4: eor w10, w10, w21, lsr #16
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ shiftRight(uint64(uint32(s.x[21])), 16, 32))))
	if runtime != nil && runtime.stopPC == 0x100236ad8 {
		return 0x100236ad8, true, nil
	}
	// 0x100236ad8: add w10, w10, #0x68
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(0x68))))
	if runtime != nil && runtime.stopPC == 0x100236adc {
		return 0x100236adc, true, nil
	}
	// 0x100236adc: eor w8, w10, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100236ae0 {
		return 0x100236ae0, true, nil
	}
	// 0x100236ae0: eor w8, w8, #0xfffffffd
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0xfffffffd))))
	if runtime != nil && runtime.stopPC == 0x100236ae4 {
		return 0x100236ae4, true, nil
	}
	// 0x100236ae4: ldr x10, [sp, #0x450]
	value1922, err := s.load(addOffset(s.sp, 1104), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value1922)
	if runtime != nil && runtime.stopPC == 0x100236ae8 {
		return 0x100236ae8, true, nil
	}
	// 0x100236ae8: strb w8, [x28, x10]
	if err := s.storeValue((s.x[28] + s.x[10]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236aec {
		return 0x100236aec, true, nil
	}
	// 0x100236aec: lsr w8, w9, #0x18
	s.x[8] = uint64(uint32(shiftRight(uint64(uint32(s.x[9])), int(uint64(0x18)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236af0 {
		return 0x100236af0, true, nil
	}
	// 0x100236af0: mov w9, #0xac
	s.x[9] = uint64(uint32(uint64(0xac)))
	if runtime != nil && runtime.stopPC == 0x100236af4 {
		return 0x100236af4, true, nil
	}
	// 0x100236af4: eor x8, x8, x9
	s.x[8] = uint64((s.x[8] ^ s.x[9]))
	if runtime != nil && runtime.stopPC == 0x100236af8 {
		return 0x100236af8, true, nil
	}
	// 0x100236af8: ldrb w8, [x20, x8]
	value1927, err := s.load((s.x[20] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1927))
	if runtime != nil && runtime.stopPC == 0x100236afc {
		return 0x100236afc, true, nil
	}
	// 0x100236afc: lsr w9, w8, #3
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236b00 {
		return 0x100236b00, true, nil
	}
	// 0x100236b00: orr w8, w9, w8, lsl #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) | shiftLeft(uint64(uint32(s.x[8])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236b04 {
		return 0x100236b04, true, nil
	}
	// 0x100236b04: and w9, w8, #0xff
	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236b08 {
		return 0x100236b08, true, nil
	}
	// 0x100236b08: eor w8, w8, w9, lsr #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[9])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236b0c {
		return 0x100236b0c, true, nil
	}
	// 0x100236b0c: eor w8, w8, w9, lsr #1
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftRight(uint64(uint32(s.x[9])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236b10 {
		return 0x100236b10, true, nil
	}
	// 0x100236b10: mov w9, #5
	s.x[9] = uint64(uint32(uint64(0x5)))
	if runtime != nil && runtime.stopPC == 0x100236b14 {
		return 0x100236b14, true, nil
	}
	// 0x100236b14: eor w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236b18 {
		return 0x100236b18, true, nil
	}
	// 0x100236b18: ldr x9, [sp, #0x490]
	value1935, err := s.load(addOffset(s.sp, 1168), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value1935)
	if runtime != nil && runtime.stopPC == 0x100236b1c {
		return 0x100236b1c, true, nil
	}
	// 0x100236b1c: strb w8, [x28, x9]
	if err := s.storeValue((s.x[28] + s.x[9]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236b20 {
		return 0x100236b20, true, nil
	}
	// 0x100236b20: and w8, w11, #0xff
	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100236b24 {
		return 0x100236b24, true, nil
	}
	// 0x100236b24: mov w9, #0xbf
	s.x[9] = uint64(uint32(uint64(0xbf)))
	if runtime != nil && runtime.stopPC == 0x100236b28 {
		return 0x100236b28, true, nil
	}
	// 0x100236b28: eor x8, x8, x9
	s.x[8] = uint64((s.x[8] ^ s.x[9]))
	if runtime != nil && runtime.stopPC == 0x100236b2c {
		return 0x100236b2c, true, nil
	}
	// 0x100236b2c: ldrb w8, [x0, x8]
	value1940, err := s.load((s.x[0] + s.x[8]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(uint32(value1940))
	if runtime != nil && runtime.stopPC == 0x100236b30 {
		return 0x100236b30, true, nil
	}
	// 0x100236b30: lsr w9, w8, #3
	s.x[9] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x3)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236b34 {
		return 0x100236b34, true, nil
	}
	// 0x100236b34: orr w8, w9, w8, lsl #5
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) | shiftLeft(uint64(uint32(s.x[8])), 5, 32))))
	if runtime != nil && runtime.stopPC == 0x100236b38 {
		return 0x100236b38, true, nil
	}
	// 0x100236b38: mov w9, #0x5110
	s.x[9] = uint64(uint32(uint64(0x5110)))
	if runtime != nil && runtime.stopPC == 0x100236b3c {
		return 0x100236b3c, true, nil
	}
	// 0x100236b3c: movk w9, #0x7d57, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7d57, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b40 {
		return 0x100236b40, true, nil
	}
	// 0x100236b40: ldr w12, [sp, #0x464]
	value1945, err := s.load(addOffset(s.sp, 1124), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value1945))
	if runtime != nil && runtime.stopPC == 0x100236b44 {
		return 0x100236b44, true, nil
	}
	// 0x100236b44: eor w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100236b48 {
		return 0x100236b48, true, nil
	}
	// 0x100236b48: mov w10, #0xaef0
	s.x[10] = uint64(uint32(uint64(0xaef0)))
	if runtime != nil && runtime.stopPC == 0x100236b4c {
		return 0x100236b4c, true, nil
	}
	// 0x100236b4c: movk w10, #0x82a8, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x82a8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b50 {
		return 0x100236b50, true, nil
	}
	// 0x100236b50: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236b54 {
		return 0x100236b54, true, nil
	}
	// 0x100236b54: mov w10, #0xb5c3
	s.x[10] = uint64(uint32(uint64(0xb5c3)))
	if runtime != nil && runtime.stopPC == 0x100236b58 {
		return 0x100236b58, true, nil
	}
	// 0x100236b58: movk w10, #0xccee, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xccee, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b5c {
		return 0x100236b5c, true, nil
	}
	// 0x100236b5c: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236b60 {
		return 0x100236b60, true, nil
	}
	// 0x100236b60: mov w11, #0x4a3d
	s.x[11] = uint64(uint32(uint64(0x4a3d)))
	if runtime != nil && runtime.stopPC == 0x100236b64 {
		return 0x100236b64, true, nil
	}
	// 0x100236b64: movk w11, #0x3311, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x3311, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b68 {
		return 0x100236b68, true, nil
	}
	// 0x100236b68: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236b6c {
		return 0x100236b6c, true, nil
	}
	// 0x100236b6c: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236b70 {
		return 0x100236b70, true, nil
	}
	// 0x100236b70: mov w10, #0x75cf
	s.x[10] = uint64(uint32(uint64(0x75cf)))
	if runtime != nil && runtime.stopPC == 0x100236b74 {
		return 0x100236b74, true, nil
	}
	// 0x100236b74: movk w10, #0xed1e, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xed1e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b78 {
		return 0x100236b78, true, nil
	}
	// 0x100236b78: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236b7c {
		return 0x100236b7c, true, nil
	}
	// 0x100236b7c: mov w11, #0x8a31
	s.x[11] = uint64(uint32(uint64(0x8a31)))
	if runtime != nil && runtime.stopPC == 0x100236b80 {
		return 0x100236b80, true, nil
	}
	// 0x100236b80: movk w11, #0x12e1, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x12e1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b84 {
		return 0x100236b84, true, nil
	}
	// 0x100236b84: add w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236b88 {
		return 0x100236b88, true, nil
	}
	// 0x100236b88: eor w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236b8c {
		return 0x100236b8c, true, nil
	}
	// 0x100236b8c: mov w10, #0x145f
	s.x[10] = uint64(uint32(uint64(0x145f)))
	if runtime != nil && runtime.stopPC == 0x100236b90 {
		return 0x100236b90, true, nil
	}
	// 0x100236b90: movk w10, #0xd162, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xd162, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236b94 {
		return 0x100236b94, true, nil
	}
	// 0x100236b94: add w9, w9, w10
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236b98 {
		return 0x100236b98, true, nil
	}
	// 0x100236b98: mov w10, #0xc9d7
	s.x[10] = uint64(uint32(uint64(0xc9d7)))
	if runtime != nil && runtime.stopPC == 0x100236b9c {
		return 0x100236b9c, true, nil
	}
	// 0x100236b9c: movk w10, #0xb691, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xb691, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236ba0 {
		return 0x100236ba0, true, nil
	}
	// 0x100236ba0: eor w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236ba4 {
		return 0x100236ba4, true, nil
	}
	// 0x100236ba4: mov w11, #0x7ccd
	s.x[11] = uint64(uint32(uint64(0x7ccd)))
	if runtime != nil && runtime.stopPC == 0x100236ba8 {
		return 0x100236ba8, true, nil
	}
	// 0x100236ba8: movk w11, #0x8b45, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x8b45, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236bac {
		return 0x100236bac, true, nil
	}
	// 0x100236bac: and w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236bb0 {
		return 0x100236bb0, true, nil
	}
	// 0x100236bb0: mov w12, #0x78cd
	s.x[12] = uint64(uint32(uint64(0x78cd)))
	if runtime != nil && runtime.stopPC == 0x100236bb4 {
		return 0x100236bb4, true, nil
	}
	// 0x100236bb4: movk w12, #0xa211, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xa211, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236bb8 {
		return 0x100236bb8, true, nil
	}
	// 0x100236bb8: eor w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236bbc {
		return 0x100236bbc, true, nil
	}
	// 0x100236bbc: and w12, w12, w11, lsl #1
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & shiftLeft(uint64(uint32(s.x[11])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236bc0 {
		return 0x100236bc0, true, nil
	}
	// 0x100236bc0: eor w11, w12, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236bc4 {
		return 0x100236bc4, true, nil
	}
	// 0x100236bc4: lsl w12, w10, #1
	s.x[12] = uint64(uint32(shiftLeft(uint64(uint32(s.x[10])), int(uint64(0x1)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236bc8 {
		return 0x100236bc8, true, nil
	}
	// 0x100236bc8: mov w13, #0x6a34
	s.x[13] = uint64(uint32(uint64(0x6a34)))
	if runtime != nil && runtime.stopPC == 0x100236bcc {
		return 0x100236bcc, true, nil
	}
	// 0x100236bcc: movk w13, #0x7ba9, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x7ba9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236bd0 {
		return 0x100236bd0, true, nil
	}
	// 0x100236bd0: eor w13, w12, w13
	s.x[13] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236bd4 {
		return 0x100236bd4, true, nil
	}
	// 0x100236bd4: and w10, w13, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236bd8 {
		return 0x100236bd8, true, nil
	}
	// 0x100236bd8: mov w14, #0xb51a
	s.x[14] = uint64(uint32(uint64(0xb51a)))
	if runtime != nil && runtime.stopPC == 0x100236bdc {
		return 0x100236bdc, true, nil
	}
	// 0x100236bdc: movk w14, #0x3dd4, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x3dd4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236be0 {
		return 0x100236be0, true, nil
	}
	// 0x100236be0: and w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236be4 {
		return 0x100236be4, true, nil
	}
	// 0x100236be4: eor w10, w10, w12
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236be8 {
		return 0x100236be8, true, nil
	}
	// 0x100236be8: mov w12, #0x950a
	s.x[12] = uint64(uint32(uint64(0x950a)))
	if runtime != nil && runtime.stopPC == 0x100236bec {
		return 0x100236bec, true, nil
	}
	// 0x100236bec: movk w12, #0x454, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x454, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236bf0 {
		return 0x100236bf0, true, nil
	}
	// 0x100236bf0: eor w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236bf4 {
		return 0x100236bf4, true, nil
	}
	// 0x100236bf4: mov w13, #0x2010
	s.x[13] = uint64(uint32(uint64(0x2010)))
	if runtime != nil && runtime.stopPC == 0x100236bf8 {
		return 0x100236bf8, true, nil
	}
	// 0x100236bf8: movk w13, #0x3900, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x3900, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236bfc {
		return 0x100236bfc, true, nil
	}
	// 0x100236bfc: eor w10, w10, w13
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c00 {
		return 0x100236c00, true, nil
	}
	// 0x100236c00: and w10, w10, w11, lsl #2
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & shiftLeft(uint64(uint32(s.x[11])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x100236c04 {
		return 0x100236c04, true, nil
	}
	// 0x100236c04: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236c08 {
		return 0x100236c08, true, nil
	}
	// 0x100236c08: lsl w11, w12, #2
	s.x[11] = uint64(uint32(shiftLeft(uint64(uint32(s.x[12])), int(uint64(0x2)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236c0c {
		return 0x100236c0c, true, nil
	}
	// 0x100236c0c: mov w13, #0xd468
	s.x[13] = uint64(uint32(uint64(0xd468)))
	if runtime != nil && runtime.stopPC == 0x100236c10 {
		return 0x100236c10, true, nil
	}
	// 0x100236c10: movk w13, #0xf752, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xf752, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c14 {
		return 0x100236c14, true, nil
	}
	// 0x100236c14: eor w13, w11, w13
	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c18 {
		return 0x100236c18, true, nil
	}
	// 0x100236c18: and w12, w13, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236c1c {
		return 0x100236c1c, true, nil
	}
	// 0x100236c1c: mov w13, #0xb518
	s.x[13] = uint64(uint32(uint64(0xb518)))
	if runtime != nil && runtime.stopPC == 0x100236c20 {
		return 0x100236c20, true, nil
	}
	// 0x100236c20: movk w13, #0x3dd4, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x3dd4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c24 {
		return 0x100236c24, true, nil
	}
	// 0x100236c24: and w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c28 {
		return 0x100236c28, true, nil
	}
	// 0x100236c28: eor w11, w12, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236c2c {
		return 0x100236c2c, true, nil
	}
	// 0x100236c2c: mov w12, #0x2112
	s.x[12] = uint64(uint32(uint64(0x2112)))
	if runtime != nil && runtime.stopPC == 0x100236c30 {
		return 0x100236c30, true, nil
	}
	// 0x100236c30: movk w12, #0x884, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x884, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c34 {
		return 0x100236c34, true, nil
	}
	// 0x100236c34: eor w12, w11, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236c38 {
		return 0x100236c38, true, nil
	}
	// 0x100236c38: mov w13, #0x9400
	s.x[13] = uint64(uint32(uint64(0x9400)))
	if runtime != nil && runtime.stopPC == 0x100236c3c {
		return 0x100236c3c, true, nil
	}
	// 0x100236c3c: movk w13, #0x3550, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x3550, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c40 {
		return 0x100236c40, true, nil
	}
	// 0x100236c40: eor w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c44 {
		return 0x100236c44, true, nil
	}
	// 0x100236c44: and w11, w11, w10, lsl #4
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & shiftLeft(uint64(uint32(s.x[10])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x100236c48 {
		return 0x100236c48, true, nil
	}
	// 0x100236c48: eor w10, w11, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236c4c {
		return 0x100236c4c, true, nil
	}
	// 0x100236c4c: lsl w11, w12, #4
	s.x[11] = uint64(uint32(shiftLeft(uint64(uint32(s.x[12])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236c50 {
		return 0x100236c50, true, nil
	}
	// 0x100236c50: mov w13, #0x51a0
	s.x[13] = uint64(uint32(uint64(0x51a0)))
	if runtime != nil && runtime.stopPC == 0x100236c54 {
		return 0x100236c54, true, nil
	}
	// 0x100236c54: movk w13, #0xdd4b, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xdd4b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c58 {
		return 0x100236c58, true, nil
	}
	// 0x100236c58: eor w13, w11, w13
	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c5c {
		return 0x100236c5c, true, nil
	}
	// 0x100236c5c: and w12, w13, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236c60 {
		return 0x100236c60, true, nil
	}
	// 0x100236c60: mov w13, #0xb500
	s.x[13] = uint64(uint32(uint64(0xb500)))
	if runtime != nil && runtime.stopPC == 0x100236c64 {
		return 0x100236c64, true, nil
	}
	// 0x100236c64: movk w13, #0x3dd4, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x3dd4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c68 {
		return 0x100236c68, true, nil
	}
	// 0x100236c68: and w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c6c {
		return 0x100236c6c, true, nil
	}
	// 0x100236c6c: eor w11, w12, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236c70 {
		return 0x100236c70, true, nil
	}
	// 0x100236c70: mov w12, #0xa41a
	s.x[12] = uint64(uint32(uint64(0xa41a)))
	if runtime != nil && runtime.stopPC == 0x100236c74 {
		return 0x100236c74, true, nil
	}
	// 0x100236c74: movk w12, #0x2094, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x2094, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c78 {
		return 0x100236c78, true, nil
	}
	// 0x100236c78: eor w12, w11, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236c7c {
		return 0x100236c7c, true, nil
	}
	// 0x100236c7c: mov w13, #0x1100
	s.x[13] = uint64(uint32(uint64(0x1100)))
	if runtime != nil && runtime.stopPC == 0x100236c80 {
		return 0x100236c80, true, nil
	}
	// 0x100236c80: movk w13, #0x1d40, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x1d40, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236c84 {
		return 0x100236c84, true, nil
	}
	// 0x100236c84: eor w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236c88 {
		return 0x100236c88, true, nil
	}
	// 0x100236c88: and w11, w11, w10, lsl #8
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & shiftLeft(uint64(uint32(s.x[10])), 8, 32))))
	if runtime != nil && runtime.stopPC == 0x100236c8c {
		return 0x100236c8c, true, nil
	}
	// 0x100236c8c: eor w10, w10, w14
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236c90 {
		return 0x100236c90, true, nil
	}
	// 0x100236c90: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236c94 {
		return 0x100236c94, true, nil
	}
	// 0x100236c94: lsl w11, w12, #8
	s.x[11] = uint64(uint32(shiftLeft(uint64(uint32(s.x[12])), int(uint64(0x8)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236c98 {
		return 0x100236c98, true, nil
	}
	// 0x100236c98: mov w13, #0x54b50000
	s.x[13] = uint64(uint32(uint64(0x54b50000)))
	if runtime != nil && runtime.stopPC == 0x100236c9c {
		return 0x100236c9c, true, nil
	}
	// 0x100236c9c: eor w13, w11, w13
	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236ca0 {
		return 0x100236ca0, true, nil
	}
	// 0x100236ca0: and w12, w13, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236ca4 {
		return 0x100236ca4, true, nil
	}
	// 0x100236ca4: mov w14, #0x3dd40000
	s.x[14] = uint64(uint32(uint64(0x3dd40000)))
	if runtime != nil && runtime.stopPC == 0x100236ca8 {
		return 0x100236ca8, true, nil
	}
	// 0x100236ca8: and w11, w11, w14
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236cac {
		return 0x100236cac, true, nil
	}
	// 0x100236cac: mov w13, #0x29400000
	s.x[13] = uint64(uint32(uint64(0x29400000)))
	if runtime != nil && runtime.stopPC == 0x100236cb0 {
		return 0x100236cb0, true, nil
	}
	// 0x100236cb0: eor w11, w11, w13
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236cb4 {
		return 0x100236cb4, true, nil
	}
	// 0x100236cb4: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100236cb8 {
		return 0x100236cb8, true, nil
	}
	// 0x100236cb8: lsl w12, w10, #0x10
	s.x[12] = uint64(uint32(shiftLeft(uint64(uint32(s.x[10])), int(uint64(0x10)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236cbc {
		return 0x100236cbc, true, nil
	}
	// 0x100236cbc: mov w13, #0x351a0000
	s.x[13] = uint64(uint32(uint64(0x351a0000)))
	if runtime != nil && runtime.stopPC == 0x100236cc0 {
		return 0x100236cc0, true, nil
	}
	// 0x100236cc0: eor w13, w12, w13
	s.x[13] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100236cc4 {
		return 0x100236cc4, true, nil
	}
	// 0x100236cc4: and w11, w13, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236cc8 {
		return 0x100236cc8, true, nil
	}
	// 0x100236cc8: and w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236ccc {
		return 0x100236ccc, true, nil
	}
	// 0x100236ccc: eor w10, w12, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236cd0 {
		return 0x100236cd0, true, nil
	}
	// 0x100236cd0: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236cd4 {
		return 0x100236cd4, true, nil
	}
	// 0x100236cd4: ldr w12, [sp, #0x378]
	value2046, err := s.load(addOffset(s.sp, 888), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(uint32(value2046))
	if runtime != nil && runtime.stopPC == 0x100236cd8 {
		return 0x100236cd8, true, nil
	}
	// 0x100236cd8: mov w13, #0xf1e1
	s.x[13] = uint64(uint32(uint64(0xf1e1)))
	if runtime != nil && runtime.stopPC == 0x100236cdc {
		return 0x100236cdc, true, nil
	}
	// 0x100236cdc: movk w13, #0xb41e, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xb41e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236ce0 {
		return 0x100236ce0, true, nil
	}
	// 0x100236ce0: cmp w12, w13
	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[13]))))
	if runtime != nil && runtime.stopPC == 0x100236ce4 {
		return 0x100236ce4, true, nil
	}
	// 0x100236ce4: eor w9, w9, w10, lsl #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ shiftLeft(uint64(uint32(s.x[10])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x100236ce8 {
		return 0x100236ce8, true, nil
	}
	// 0x100236ce8: lsr w10, w8, #1
	s.x[10] = uint64(uint32(shiftRight(uint64(uint32(s.x[8])), int(uint64(0x1)), 32)))
	if runtime != nil && runtime.stopPC == 0x100236cec {
		return 0x100236cec, true, nil
	}
	// 0x100236cec: and w10, w10, #0x3f
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0x3f))))
	if runtime != nil && runtime.stopPC == 0x100236cf0 {
		return 0x100236cf0, true, nil
	}
	// 0x100236cf0: orr w10, w10, w2
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) | uint64(uint32(s.x[2])))))
	if runtime != nil && runtime.stopPC == 0x100236cf4 {
		return 0x100236cf4, true, nil
	}
	// 0x100236cf4: eor w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236cf8 {
		return 0x100236cf8, true, nil
	}
	// 0x100236cf8: cset w10, lo
	s.x[10] = uint64(uint32(boolValue(s.condition("lo"))))
	if runtime != nil && runtime.stopPC == 0x100236cfc {
		return 0x100236cfc, true, nil
	}
	// 0x100236cfc: eor w8, w8, #0xfffffffb
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(0xfffffffb))))
	if runtime != nil && runtime.stopPC == 0x100236d00 {
		return 0x100236d00, true, nil
	}
	// 0x100236d00: ldr x11, [sp, #0x4a8]
	value2057, err := s.load(addOffset(s.sp, 1192), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(value2057)
	if runtime != nil && runtime.stopPC == 0x100236d04 {
		return 0x100236d04, true, nil
	}
	// 0x100236d04: strb w8, [x28, x11]
	if err := s.storeValue((s.x[28] + s.x[11]), uint64(uint32(s.x[8])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100236d08 {
		return 0x100236d08, true, nil
	}
	// 0x100236d08: mov w8, #0x2611
	s.x[8] = uint64(uint32(uint64(0x2611)))
	if runtime != nil && runtime.stopPC == 0x100236d0c {
		return 0x100236d0c, true, nil
	}
	// 0x100236d0c: movk w8, #0x903b, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x903b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d10 {
		return 0x100236d10, true, nil
	}
	// 0x100236d10: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100236d14 {
		return 0x100236d14, true, nil
	}
	// 0x100236d14: mov w11, #0x5e0c
	s.x[11] = uint64(uint32(uint64(0x5e0c)))
	if runtime != nil && runtime.stopPC == 0x100236d18 {
		return 0x100236d18, true, nil
	}
	// 0x100236d18: movk w11, #0xa9af, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa9af, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d1c {
		return 0x100236d1c, true, nil
	}
	// 0x100236d1c: add w8, w8, w11
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d20 {
		return 0x100236d20, true, nil
	}
	// 0x100236d20: mov w11, #0xfb95
	s.x[11] = uint64(uint32(uint64(0xfb95)))
	if runtime != nil && runtime.stopPC == 0x100236d24 {
		return 0x100236d24, true, nil
	}
	// 0x100236d24: movk w11, #0x1a88, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x1a88, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d28 {
		return 0x100236d28, true, nil
	}
	// 0x100236d28: eor w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d2c {
		return 0x100236d2c, true, nil
	}
	// 0x100236d2c: mov w14, #0x8390
	s.x[14] = uint64(uint32(uint64(0x8390)))
	if runtime != nil && runtime.stopPC == 0x100236d30 {
		return 0x100236d30, true, nil
	}
	// 0x100236d30: movk w14, #0x231c, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x231c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d34 {
		return 0x100236d34, true, nil
	}
	// 0x100236d34: add w11, w11, w14
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236d38 {
		return 0x100236d38, true, nil
	}
	// 0x100236d38: eor w8, w8, w11
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d3c {
		return 0x100236d3c, true, nil
	}
	// 0x100236d3c: mov w11, #0xcb7d
	s.x[11] = uint64(uint32(uint64(0xcb7d)))
	if runtime != nil && runtime.stopPC == 0x100236d40 {
		return 0x100236d40, true, nil
	}
	// 0x100236d40: movk w11, #0x107f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x107f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d44 {
		return 0x100236d44, true, nil
	}
	// 0x100236d44: eor w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d48 {
		return 0x100236d48, true, nil
	}
	// 0x100236d48: mov w14, #0xb368
	s.x[14] = uint64(uint32(uint64(0xb368)))
	if runtime != nil && runtime.stopPC == 0x100236d4c {
		return 0x100236d4c, true, nil
	}
	// 0x100236d4c: movk w14, #0x29eb, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x29eb, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d50 {
		return 0x100236d50, true, nil
	}
	// 0x100236d50: add w11, w11, w14
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100236d54 {
		return 0x100236d54, true, nil
	}
	// 0x100236d54: eor w8, w8, w11
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d58 {
		return 0x100236d58, true, nil
	}
	// 0x100236d58: mov w11, #0x82fd
	s.x[11] = uint64(uint32(uint64(0x82fd)))
	if runtime != nil && runtime.stopPC == 0x100236d5c {
		return 0x100236d5c, true, nil
	}
	// 0x100236d5c: movk w11, #0x10c6, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x10c6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100236d60 {
		return 0x100236d60, true, nil
	}
	// 0x100236d60: add w8, w8, w11
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d64 {
		return 0x100236d64, true, nil
	}
	// 0x100236d64: cmp w8, w13
	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[13]))))
	if runtime != nil && runtime.stopPC == 0x100236d68 {
		return 0x100236d68, true, nil
	}
	// 0x100236d68: cset w11, lo
	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))
	if runtime != nil && runtime.stopPC == 0x100236d6c {
		return 0x100236d6c, true, nil
	}
	// 0x100236d6c: eor w11, w10, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d70 {
		return 0x100236d70, true, nil
	}
	// 0x100236d70: cmp w8, w12
	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[12]))))
	if runtime != nil && runtime.stopPC == 0x100236d74 {
		return 0x100236d74, true, nil
	}
	// 0x100236d74: cset w8, lo
	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))
	if runtime != nil && runtime.stopPC == 0x100236d78 {
		return 0x100236d78, true, nil
	}
	// 0x100236d78: cmp w11, #0
	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(0x0)))
	if runtime != nil && runtime.stopPC == 0x100236d7c {
		return 0x100236d7c, true, nil
	}
	// 0x100236d7c: csel w8, w10, w8, ne
	s.x[8] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[10])), uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100236d80 {
		return 0x100236d80, true, nil
	}
	// 0x100236d80: mov w10, #0x7f56
	s.x[10] = uint64(uint32(uint64(0x7f56)))
	if runtime != nil && runtime.stopPC == 0x100236d84 {
		return 0x100236d84, true, nil
	}
	// 0x100236d84: mul w8, w8, w10
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236d88 {
		return 0x100236d88, true, nil
	}
	// 0x100236d88: ldr w11, [sp, #0x4c4]
	value2091, err := s.load(addOffset(s.sp, 1220), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value2091))
	if runtime != nil && runtime.stopPC == 0x100236d8c {
		return 0x100236d8c, true, nil
	}
	// 0x100236d8c: eor w8, w8, w11
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100236d90 {
		return 0x100236d90, true, nil
	}
	// 0x100236d90: mov w10, #0x9315
	s.x[10] = uint64(uint32(uint64(0x9315)))
	if runtime != nil && runtime.stopPC == 0x100236d94 {
		return 0x100236d94, true, nil
	}
	// 0x100236d94: eor w11, w11, w10
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100236d98 {
		return 0x100236d98, true, nil
	}
	// 0x100236d98: ldur x10, [x29, #-0xe0]
	value2095, err := s.load(addOffset(s.x[29], -224), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value2095)
	if runtime != nil && runtime.stopPC == 0x100236d9c {
		return 0x100236d9c, true, nil
	}
	// 0x100236d9c: ldr x8, [x10, w8, sxtw #3]
	value2096, err := s.load((s.x[10] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(value2096)
	if runtime != nil && runtime.stopPC == 0x100236da0 {
		return 0x100236da0, true, nil
	}
	// 0x100236da0: mov x17, #0xe6c7
	s.x[17] = uint64(uint64(0xe6c7))
	if runtime != nil && runtime.stopPC == 0x100236da4 {
		return 0x100236da4, true, nil
	}
	// 0x100236da4: braa x8, x17
	next = s.x[8]
	return next, false, nil
}

func rawSourceBootstrapBlock100128818(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100128818 {
		return 0x100128818, true, nil
	}
	// 0x100128818: mov w8, #0xd844
	s.x[8] = uint64(uint32(uint64(0xd844)))
	if runtime != nil && runtime.stopPC == 0x10012881c {
		return 0x10012881c, true, nil
	}
	// 0x10012881c: movk w8, #0xad60, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xad60, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100128820 {
		return 0x100128820, true, nil
	}
	// 0x100128820: ldr w11, [sp, #0x404]
	value2, err := s.load(addOffset(s.sp, 1028), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[11] = uint64(uint32(value2))
	if runtime != nil && runtime.stopPC == 0x100128824 {
		return 0x100128824, true, nil
	}
	// 0x100128824: add w8, w11, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100128828 {
		return 0x100128828, true, nil
	}
	// 0x100128828: mov w10, #0xfe8f
	s.x[10] = uint64(uint32(uint64(0xfe8f)))
	if runtime != nil && runtime.stopPC == 0x10012882c {
		return 0x10012882c, true, nil
	}
	// 0x10012882c: movk w10, #0x28f1, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x28f1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100128830 {
		return 0x100128830, true, nil
	}
	// 0x100128830: and w10, w8, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100128834 {
		return 0x100128834, true, nil
	}
	// 0x100128834: mov w8, #0x87e5
	s.x[8] = uint64(uint32(uint64(0x87e5)))
	if runtime != nil && runtime.stopPC == 0x100128838 {
		return 0x100128838, true, nil
	}
	// 0x100128838: movk w8, #0xc66b, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xc66b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10012883c {
		return 0x10012883c, true, nil
	}
	// 0x10012883c: eor w12, w9, w8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100128840 {
		return 0x100128840, true, nil
	}
	// 0x100128840: mov w8, #0x81d6
	s.x[8] = uint64(uint32(uint64(0x81d6)))
	if runtime != nil && runtime.stopPC == 0x100128844 {
		return 0x100128844, true, nil
	}
	// 0x100128844: movk w8, #0x4296, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x4296, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100128848 {
		return 0x100128848, true, nil
	}
	// 0x100128848: eor w8, w10, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x10012884c {
		return 0x10012884c, true, nil
	}
	// 0x10012884c: ldr w9, [sp, #0x3f8]
	value13, err := s.load(addOffset(s.sp, 1016), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value13))
	if runtime != nil && runtime.stopPC == 0x100128850 {
		return 0x100128850, true, nil
	}
	// 0x100128850: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100128854 {
		return 0x100128854, true, nil
	}
	// 0x100128854: str w8, [sp, #0x4b4]
	if err := s.storeValue(addOffset(s.sp, 1204), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100128858 {
		return 0x100128858, true, nil
	}
	// 0x100128858: mov w8, #0x9cfd
	s.x[8] = uint64(uint32(uint64(0x9cfd)))
	if runtime != nil && runtime.stopPC == 0x10012885c {
		return 0x10012885c, true, nil
	}
	// 0x10012885c: movk w8, #0xa099, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xa099, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100128860 {
		return 0x100128860, true, nil
	}
	// 0x100128860: ldr w9, [sp, #0x3f4]
	value18, err := s.load(addOffset(s.sp, 1012), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value18))
	if runtime != nil && runtime.stopPC == 0x100128864 {
		return 0x100128864, true, nil
	}
	// 0x100128864: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100128868 {
		return 0x100128868, true, nil
	}
	// 0x100128868: str w8, [sp, #0x4b8]
	if err := s.storeValue(addOffset(s.sp, 1208), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x10012886c {
		return 0x10012886c, true, nil
	}
	// 0x10012886c: mov w8, #0x8d51
	s.x[8] = uint64(uint32(uint64(0x8d51)))
	if runtime != nil && runtime.stopPC == 0x100128870 {
		return 0x100128870, true, nil
	}
	// 0x100128870: movk w8, #0xb9f8, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xb9f8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100128874 {
		return 0x100128874, true, nil
	}
	// 0x100128874: ldr w9, [sp, #0x3fc]
	value23, err := s.load(addOffset(s.sp, 1020), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value23))
	if runtime != nil && runtime.stopPC == 0x100128878 {
		return 0x100128878, true, nil
	}
	// 0x100128878: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x10012887c {
		return 0x10012887c, true, nil
	}
	// 0x10012887c: str w8, [sp, #0x4bc]
	if err := s.storeValue(addOffset(s.sp, 1212), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100128880 {
		return 0x100128880, true, nil
	}
	// 0x100128880: mov w8, #0x8ccb
	s.x[8] = uint64(uint32(uint64(0x8ccb)))
	if runtime != nil && runtime.stopPC == 0x100128884 {
		return 0x100128884, true, nil
	}
	// 0x100128884: movk w8, #0x1e5d, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x1e5d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100128888 {
		return 0x100128888, true, nil
	}
	// 0x100128888: ldr w9, [sp, #0x400]
	value28, err := s.load(addOffset(s.sp, 1024), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(uint32(value28))
	if runtime != nil && runtime.stopPC == 0x10012888c {
		return 0x10012888c, true, nil
	}
	// 0x10012888c: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100128890 {
		return 0x100128890, true, nil
	}
	// 0x100128890: str w8, [sp, #0x4c0]
	if err := s.storeValue(addOffset(s.sp, 1216), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100128894 {
		return 0x100128894, true, nil
	}
	// 0x100128894: mov w8, #0x2b41
	s.x[8] = uint64(uint32(uint64(0x2b41)))
	if runtime != nil && runtime.stopPC == 0x100128898 {
		return 0x100128898, true, nil
	}
	// 0x100128898: movk w8, #0x468b, lsl #16
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x468b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10012889c {
		return 0x10012889c, true, nil
	}
	// 0x10012889c: cmp w11, w8
	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[8]))))
	if runtime != nil && runtime.stopPC == 0x1001288a0 {
		return 0x1001288a0, true, nil
	}
	// 0x1001288a0: cset w8, hi
	s.x[8] = uint64(uint32(boolValue(s.condition("hi"))))
	if runtime != nil && runtime.stopPC == 0x1001288a4 {
		return 0x1001288a4, true, nil
	}
	// 0x1001288a4: mov w9, #0xa7ad
	s.x[9] = uint64(uint32(uint64(0xa7ad)))
	if runtime != nil && runtime.stopPC == 0x1001288a8 {
		return 0x1001288a8, true, nil
	}
	// 0x1001288a8: mul w8, w8, w9
	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1001288ac {
		return 0x1001288ac, true, nil
	}
	// 0x1001288ac: mov w9, #0x872e
	s.x[9] = uint64(uint32(uint64(0x872e)))
	if runtime != nil && runtime.stopPC == 0x1001288b0 {
		return 0x1001288b0, true, nil
	}
	// 0x1001288b0: movk w9, #0x29ad, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x29ad, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1001288b4 {
		return 0x1001288b4, true, nil
	}
	// 0x1001288b4: eor w9, w11, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1001288b8 {
		return 0x1001288b8, true, nil
	}
	// 0x1001288b8: eor w8, w9, w8
	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x1001288bc {
		return 0x1001288bc, true, nil
	}
	// 0x1001288bc: ldur x9, [x29, #-0xe0]
	value41, err := s.load(addOffset(s.x[29], -224), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[9] = uint64(value41)
	if runtime != nil && runtime.stopPC == 0x1001288c0 {
		return 0x1001288c0, true, nil
	}
	// 0x1001288c0: ldr x8, [x9, w8, sxtw #3]
	value42, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[8] = uint64(value42)
	if runtime != nil && runtime.stopPC == 0x1001288c4 {
		return 0x1001288c4, true, nil
	}
	// 0x1001288c4: mov x17, #0x33ea
	s.x[17] = uint64(uint64(0x33ea))
	if runtime != nil && runtime.stopPC == 0x1001288c8 {
		return 0x1001288c8, true, nil
	}
	// 0x1001288c8: braa x8, x17
	next = s.x[8]
	return next, false, nil
}

func rawSourceBootstrapBlock10048866c(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x10048866c {
		return 0x10048866c, true, nil
	}
	// 0x10048866c: mov w10, #0x2fc7
	s.x[10] = uint64(uint32(uint64(0x2fc7)))
	if runtime != nil && runtime.stopPC == 0x100488670 {
		return 0x100488670, true, nil
	}
	// 0x100488670: movk w10, #0x8ef4, lsl #16
	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x8ef4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488674 {
		return 0x100488674, true, nil
	}
	// 0x100488674: mov w11, #0xf6b9
	s.x[11] = uint64(uint32(uint64(0xf6b9)))
	if runtime != nil && runtime.stopPC == 0x100488678 {
		return 0x100488678, true, nil
	}
	// 0x100488678: movk w11, #0x6c5b, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x6c5b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10048867c {
		return 0x10048867c, true, nil
	}
	// 0x10048867c: ldr x15, [sp, #0xd98]
	value4, err := s.load(addOffset(s.sp, 3480), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(value4)
	if runtime != nil && runtime.stopPC == 0x100488680 {
		return 0x100488680, true, nil
	}
	// 0x100488680: eor w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100488684 {
		return 0x100488684, true, nil
	}
	// 0x100488684: mov w12, #0x947
	s.x[12] = uint64(uint32(uint64(0x947)))
	if runtime != nil && runtime.stopPC == 0x100488688 {
		return 0x100488688, true, nil
	}
	// 0x100488688: movk w12, #0x93a4, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x93a4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10048868c {
		return 0x10048868c, true, nil
	}
	// 0x10048868c: mov w13, #0x8d6b
	s.x[13] = uint64(uint32(uint64(0x8d6b)))
	if runtime != nil && runtime.stopPC == 0x100488690 {
		return 0x100488690, true, nil
	}
	// 0x100488690: movk w13, #0x69e7, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x69e7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488694 {
		return 0x100488694, true, nil
	}
	// 0x100488694: eor w13, w15, w13
	s.x[13] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100488698 {
		return 0x100488698, true, nil
	}
	// 0x100488698: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10048869c {
		return 0x10048869c, true, nil
	}
	// 0x10048869c: mov w12, #0x7295
	s.x[12] = uint64(uint32(uint64(0x7295)))
	if runtime != nil && runtime.stopPC == 0x1004886a0 {
		return 0x1004886a0, true, nil
	}
	// 0x1004886a0: movk w12, #0x9618, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x9618, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886a4 {
		return 0x1004886a4, true, nil
	}
	// 0x1004886a4: mov w14, #0x17
	s.x[14] = uint64(uint32(uint64(0x17)))
	if runtime != nil && runtime.stopPC == 0x1004886a8 {
		return 0x1004886a8, true, nil
	}
	// 0x1004886a8: movk w14, #0xcbc9, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xcbc9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886ac {
		return 0x1004886ac, true, nil
	}
	// 0x1004886ac: eor w14, w15, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1004886b0 {
		return 0x1004886b0, true, nil
	}
	// 0x1004886b0: mov w15, #0xffe9
	s.x[15] = uint64(uint32(uint64(0xffe9)))
	if runtime != nil && runtime.stopPC == 0x1004886b4 {
		return 0x1004886b4, true, nil
	}
	// 0x1004886b4: movk w15, #0x3436, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x3436, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886b8 {
		return 0x1004886b8, true, nil
	}
	// 0x1004886b8: add w12, w13, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1004886bc {
		return 0x1004886bc, true, nil
	}
	// 0x1004886bc: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1004886c0 {
		return 0x1004886c0, true, nil
	}
	// 0x1004886c0: add w12, w14, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1004886c4 {
		return 0x1004886c4, true, nil
	}
	// 0x1004886c4: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1004886c8 {
		return 0x1004886c8, true, nil
	}
	// 0x1004886c8: mov w12, #0x7fc0
	s.x[12] = uint64(uint32(uint64(0x7fc0)))
	if runtime != nil && runtime.stopPC == 0x1004886cc {
		return 0x1004886cc, true, nil
	}
	// 0x1004886cc: movk w12, #0x993a, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x993a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886d0 {
		return 0x1004886d0, true, nil
	}
	// 0x1004886d0: add w10, w8, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x1004886d4 {
		return 0x1004886d4, true, nil
	}
	// 0x1004886d4: eor w12, w10, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1004886d8 {
		return 0x1004886d8, true, nil
	}
	// 0x1004886d8: mov w13, #0x88a1
	s.x[13] = uint64(uint32(uint64(0x88a1)))
	if runtime != nil && runtime.stopPC == 0x1004886dc {
		return 0x1004886dc, true, nil
	}
	// 0x1004886dc: movk w13, #0x6a, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x6a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886e0 {
		return 0x1004886e0, true, nil
	}
	// 0x1004886e0: add w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1004886e4 {
		return 0x1004886e4, true, nil
	}
	// 0x1004886e4: eor w12, w11, w13
	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1004886e8 {
		return 0x1004886e8, true, nil
	}
	// 0x1004886e8: mov w13, #0x15ba
	s.x[13] = uint64(uint32(uint64(0x15ba)))
	if runtime != nil && runtime.stopPC == 0x1004886ec {
		return 0x1004886ec, true, nil
	}
	// 0x1004886ec: movk w13, #0xb6cc, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xb6cc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886f0 {
		return 0x1004886f0, true, nil
	}
	// 0x1004886f0: mov w14, #0x8575
	s.x[14] = uint64(uint32(uint64(0x8575)))
	if runtime != nil && runtime.stopPC == 0x1004886f4 {
		return 0x1004886f4, true, nil
	}
	// 0x1004886f4: movk w14, #0x4a88, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x4a88, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004886f8 {
		return 0x1004886f8, true, nil
	}
	// 0x1004886f8: mov w15, #0x186e
	s.x[15] = uint64(uint32(uint64(0x186e)))
	if runtime != nil && runtime.stopPC == 0x1004886fc {
		return 0x1004886fc, true, nil
	}
	// 0x1004886fc: movk w15, #0xfc2e, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xfc2e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488700 {
		return 0x100488700, true, nil
	}
	// 0x100488700: eor w14, w11, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x100488704 {
		return 0x100488704, true, nil
	}
	// 0x100488704: mov w16, #0x10cc
	s.x[16] = uint64(uint32(uint64(0x10cc)))
	if runtime != nil && runtime.stopPC == 0x100488708 {
		return 0x100488708, true, nil
	}
	// 0x100488708: movk w16, #0x7c40, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x7c40, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10048870c {
		return 0x10048870c, true, nil
	}
	// 0x10048870c: eor w16, w11, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100488710 {
		return 0x100488710, true, nil
	}
	// 0x100488710: mov w17, #0x8dd5
	s.x[17] = uint64(uint32(uint64(0x8dd5)))
	if runtime != nil && runtime.stopPC == 0x100488714 {
		return 0x100488714, true, nil
	}
	// 0x100488714: movk w17, #0xcae6, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xcae6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488718 {
		return 0x100488718, true, nil
	}
	// 0x100488718: mov w1, #0x7fff
	s.x[1] = uint64(uint32(uint64(0x7fff)))
	if runtime != nil && runtime.stopPC == 0x10048871c {
		return 0x10048871c, true, nil
	}
	// 0x10048871c: movk w1, #0x7ffb, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x7ffb, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488720 {
		return 0x100488720, true, nil
	}
	// 0x100488720: eor w1, w11, w1
	s.x[1] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100488724 {
		return 0x100488724, true, nil
	}
	// 0x100488724: add w14, w14, w15
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100488728 {
		return 0x100488728, true, nil
	}
	// 0x100488728: mov w15, #0xe2e8
	s.x[15] = uint64(uint32(uint64(0xe2e8)))
	if runtime != nil && runtime.stopPC == 0x10048872c {
		return 0x10048872c, true, nil
	}
	// 0x10048872c: movk w15, #0xc95d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xc95d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488730 {
		return 0x100488730, true, nil
	}
	// 0x100488730: add w12, w12, w13
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100488734 {
		return 0x100488734, true, nil
	}
	// 0x100488734: eor w11, w12, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100488738 {
		return 0x100488738, true, nil
	}
	// 0x100488738: add w12, w16, w17
	s.x[12] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10048873c {
		return 0x10048873c, true, nil
	}
	// 0x10048873c: eor w12, w14, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100488740 {
		return 0x100488740, true, nil
	}
	// 0x100488740: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100488744 {
		return 0x100488744, true, nil
	}
	// 0x100488744: add w12, w1, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[1])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x100488748 {
		return 0x100488748, true, nil
	}
	// 0x100488748: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x10048874c {
		return 0x10048874c, true, nil
	}
	// 0x10048874c: mov w12, #0x1922
	s.x[12] = uint64(uint32(uint64(0x1922)))
	if runtime != nil && runtime.stopPC == 0x100488750 {
		return 0x100488750, true, nil
	}
	// 0x100488750: movk w12, #0x872c, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x872c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488754 {
		return 0x100488754, true, nil
	}
	// 0x100488754: mov w13, #0x62e7
	s.x[13] = uint64(uint32(uint64(0x62e7)))
	if runtime != nil && runtime.stopPC == 0x100488758 {
		return 0x100488758, true, nil
	}
	// 0x100488758: movk w13, #0x4959, lsl #16
	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x4959, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10048875c {
		return 0x10048875c, true, nil
	}
	// 0x10048875c: mov w14, #0x8aea
	s.x[14] = uint64(uint32(uint64(0x8aea)))
	if runtime != nil && runtime.stopPC == 0x100488760 {
		return 0x100488760, true, nil
	}
	// 0x100488760: movk w14, #0xdf39, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xdf39, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488764 {
		return 0x100488764, true, nil
	}
	// 0x100488764: eor w15, w11, w13
	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100488768 {
		return 0x100488768, true, nil
	}
	// 0x100488768: eor w14, w11, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x10048876c {
		return 0x10048876c, true, nil
	}
	// 0x10048876c: mov w16, #0x6c38
	s.x[16] = uint64(uint32(uint64(0x6c38)))
	if runtime != nil && runtime.stopPC == 0x100488770 {
		return 0x100488770, true, nil
	}
	// 0x100488770: movk w16, #0xa7ea, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa7ea, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488774 {
		return 0x100488774, true, nil
	}
	// 0x100488774: mov w17, #0xc506
	s.x[17] = uint64(uint32(uint64(0xc506)))
	if runtime != nil && runtime.stopPC == 0x100488778 {
		return 0x100488778, true, nil
	}
	// 0x100488778: movk w17, #0x575c, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x575c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10048877c {
		return 0x10048877c, true, nil
	}
	// 0x10048877c: eor w17, w11, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100488780 {
		return 0x100488780, true, nil
	}
	// 0x100488780: mov w1, #0x23dc
	s.x[1] = uint64(uint32(uint64(0x23dc)))
	if runtime != nil && runtime.stopPC == 0x100488784 {
		return 0x100488784, true, nil
	}
	// 0x100488784: movk w1, #0x2f8f, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x2f8f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x100488788 {
		return 0x100488788, true, nil
	}
	// 0x100488788: add w14, w14, w16
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x10048878c {
		return 0x10048878c, true, nil
	}
	// 0x10048878c: add w16, w17, w1
	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100488790 {
		return 0x100488790, true, nil
	}
	// 0x100488790: eor w14, w14, w16
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100488794 {
		return 0x100488794, true, nil
	}
	// 0x100488794: mov w16, #0x2d0b
	s.x[16] = uint64(uint32(uint64(0x2d0b)))
	if runtime != nil && runtime.stopPC == 0x100488798 {
		return 0x100488798, true, nil
	}
	// 0x100488798: movk w16, #0xc13c, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xc13c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x10048879c {
		return 0x10048879c, true, nil
	}
	// 0x10048879c: eor w16, w11, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1004887a0 {
		return 0x1004887a0, true, nil
	}
	// 0x1004887a0: mov w17, #0xcbd7
	s.x[17] = uint64(uint32(uint64(0xcbd7)))
	if runtime != nil && runtime.stopPC == 0x1004887a4 {
		return 0x1004887a4, true, nil
	}
	// 0x1004887a4: movk w17, #0xb9ef, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xb9ef, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004887a8 {
		return 0x1004887a8, true, nil
	}
	// 0x1004887a8: add w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1004887ac {
		return 0x1004887ac, true, nil
	}
	// 0x1004887ac: eor w14, w14, w16
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1004887b0 {
		return 0x1004887b0, true, nil
	}
	// 0x1004887b0: mov w16, #0x7bd1
	s.x[16] = uint64(uint32(uint64(0x7bd1)))
	if runtime != nil && runtime.stopPC == 0x1004887b4 {
		return 0x1004887b4, true, nil
	}
	// 0x1004887b4: movk w16, #0xce75, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xce75, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004887b8 {
		return 0x1004887b8, true, nil
	}
	// 0x1004887b8: ldr x17, [sp, #0xce0]
	value83, err := s.load(addOffset(s.sp, 3296), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(value83)
	if runtime != nil && runtime.stopPC == 0x1004887bc {
		return 0x1004887bc, true, nil
	}
	// 0x1004887bc: ldrsw x17, [x17]
	value84, err := s.load(s.x[17], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(signExtend32(uint32(value84)))
	if runtime != nil && runtime.stopPC == 0x1004887c0 {
		return 0x1004887c0, true, nil
	}
	// 0x1004887c0: ldr x1, [sp, #0xcd8]
	value85, err := s.load(addOffset(s.sp, 3288), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(value85)
	if runtime != nil && runtime.stopPC == 0x1004887c4 {
		return 0x1004887c4, true, nil
	}
	// 0x1004887c4: ldr x1, [x1]
	value86, err := s.load(s.x[1], 8)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(value86)
	if runtime != nil && runtime.stopPC == 0x1004887c8 {
		return 0x1004887c8, true, nil
	}
	// 0x1004887c8: add w14, w14, w16
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1004887cc {
		return 0x1004887cc, true, nil
	}
	// 0x1004887cc: add x14, x0, x14
	s.x[14] = uint64((s.x[0] + s.x[14]))
	if runtime != nil && runtime.stopPC == 0x1004887d0 {
		return 0x1004887d0, true, nil
	}
	// 0x1004887d0: mov x16, #-0x6d80
	s.x[16] = uint64(uint64(0xffffffffffff9280))
	if runtime != nil && runtime.stopPC == 0x1004887d4 {
		return 0x1004887d4, true, nil
	}
	// 0x1004887d4: movk x16, #0xc654, lsl #16
	s.x[16] = uint64(insertBits(s.x[16], 0xc654, 16, 16, 64))
	if runtime != nil && runtime.stopPC == 0x1004887d8 {
		return 0x1004887d8, true, nil
	}
	// 0x1004887d8: and x16, x17, x16
	s.x[16] = uint64((s.x[17] & s.x[16]))
	if runtime != nil && runtime.stopPC == 0x1004887dc {
		return 0x1004887dc, true, nil
	}
	// 0x1004887dc: add x16, x1, x16
	s.x[16] = uint64((s.x[1] + s.x[16]))
	if runtime != nil && runtime.stopPC == 0x1004887e0 {
		return 0x1004887e0, true, nil
	}
	// 0x1004887e0: ldr w16, [x16]
	value93, err := s.load(s.x[16], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[16] = uint64(uint32(value93))
	if runtime != nil && runtime.stopPC == 0x1004887e4 {
		return 0x1004887e4, true, nil
	}
	// 0x1004887e4: eor w16, w14, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1004887e8 {
		return 0x1004887e8, true, nil
	}
	// 0x1004887e8: add w15, w15, w19
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[19])))))
	if runtime != nil && runtime.stopPC == 0x1004887ec {
		return 0x1004887ec, true, nil
	}
	// 0x1004887ec: and w16, w16, #0x7fffffff
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0x7fffffff))))
	if runtime != nil && runtime.stopPC == 0x1004887f0 {
		return 0x1004887f0, true, nil
	}
	// 0x1004887f0: mov w17, #0xbb2f
	s.x[17] = uint64(uint32(uint64(0xbb2f)))
	if runtime != nil && runtime.stopPC == 0x1004887f4 {
		return 0x1004887f4, true, nil
	}
	// 0x1004887f4: movk w17, #0xbf5a, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xbf5a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1004887f8 {
		return 0x1004887f8, true, nil
	}
	// 0x1004887f8: mul w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) * uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1004887fc {
		return 0x1004887fc, true, nil
	}
	// 0x1004887fc: eor w16, w16, w16, lsr #16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ shiftRight(uint64(uint32(s.x[16])), 16, 32))))
	if runtime != nil && runtime.stopPC == 0x100488800 {
		return 0x100488800, true, nil
	}
	// 0x100488800: ldrb w15, [x22, w15, uxtw]
	value101, err := s.load((s.x[22] + uint64(uint32(uint64(uint32(s.x[15]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[15] = uint64(uint32(value101))
	if runtime != nil && runtime.stopPC == 0x100488804 {
		return 0x100488804, true, nil
	}
	// 0x100488804: mul w16, w16, w17
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) * uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100488808 {
		return 0x100488808, true, nil
	}
	// 0x100488808: lsr x17, x16, #0x18
	s.x[17] = uint64(shiftRight(s.x[16], int(uint64(0x18)), 64))
	if runtime != nil && runtime.stopPC == 0x10048880c {
		return 0x10048880c, true, nil
	}
	// 0x10048880c: ldr x1, [sp, #0x990]
	value104, err := s.load(addOffset(s.sp, 2448), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(value104)
	if runtime != nil && runtime.stopPC == 0x100488810 {
		return 0x100488810, true, nil
	}
	// 0x100488810: add x1, x1, x17
	s.x[1] = uint64((s.x[1] + s.x[17]))
	if runtime != nil && runtime.stopPC == 0x100488814 {
		return 0x100488814, true, nil
	}
	// 0x100488814: ldrb w1, [x1]
	value106, err := s.load(s.x[1], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value106))
	if runtime != nil && runtime.stopPC == 0x100488818 {
		return 0x100488818, true, nil
	}
	// 0x100488818: ldr x2, [sp, #0x9f8]
	value107, err := s.load(addOffset(s.sp, 2552), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[2] = uint64(value107)
	if runtime != nil && runtime.stopPC == 0x10048881c {
		return 0x10048881c, true, nil
	}
	// 0x10048881c: add x2, x17, x2
	s.x[2] = uint64((s.x[17] + s.x[2]))
	if runtime != nil && runtime.stopPC == 0x100488820 {
		return 0x100488820, true, nil
	}
	// 0x100488820: ldrb w2, [x2, #2]
	value109, err := s.load(addOffset(s.x[2], 2), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[2] = uint64(uint32(value109))
	if runtime != nil && runtime.stopPC == 0x100488824 {
		return 0x100488824, true, nil
	}
	// 0x100488824: mov w3, #0x2f
	s.x[3] = uint64(uint32(uint64(0x2f)))
	if runtime != nil && runtime.stopPC == 0x100488828 {
		return 0x100488828, true, nil
	}
	// 0x100488828: mul w3, w17, w3
	s.x[3] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x10048882c {
		return 0x10048882c, true, nil
	}
	// 0x10048882c: ldr x4, [sp, #0x998]
	value112, err := s.load(addOffset(s.sp, 2456), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[4] = uint64(value112)
	if runtime != nil && runtime.stopPC == 0x100488830 {
		return 0x100488830, true, nil
	}
	// 0x100488830: add x17, x17, x4
	s.x[17] = uint64((s.x[17] + s.x[4]))
	if runtime != nil && runtime.stopPC == 0x100488834 {
		return 0x100488834, true, nil
	}
	// 0x100488834: ldrb w17, [x17, #1]
	value114, err := s.load(addOffset(s.x[17], 1), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value114))
	if runtime != nil && runtime.stopPC == 0x100488838 {
		return 0x100488838, true, nil
	}
	// 0x100488838: eor w15, w15, w1
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x10048883c {
		return 0x10048883c, true, nil
	}
	// 0x10048883c: eor w17, w2, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100488840 {
		return 0x100488840, true, nil
	}
	// 0x100488840: eor w15, w15, w17
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100488844 {
		return 0x100488844, true, nil
	}
	// 0x100488844: eor w16, w16, w3
	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[3])))))
	if runtime != nil && runtime.stopPC == 0x100488848 {
		return 0x100488848, true, nil
	}
	// 0x100488848: eor w15, w15, w16
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x10048884c {
		return 0x10048884c, true, nil
	}
	// 0x10048884c: eor w15, w15, #0x1c
	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(0x1c))))
	if runtime != nil && runtime.stopPC == 0x100488850 {
		return 0x100488850, true, nil
	}
	// 0x100488850: cmp w11, w13
	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[13]))))
	if runtime != nil && runtime.stopPC == 0x100488854 {
		return 0x100488854, true, nil
	}
	// 0x100488854: strb w15, [x14]
	if err := s.storeValue(s.x[14], uint64(uint32(s.x[15])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100488858 {
		return 0x100488858, true, nil
	}
	// 0x100488858: cset w13, eq
	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))))
	if runtime != nil && runtime.stopPC == 0x10048885c {
		return 0x10048885c, true, nil
	}
	// 0x10048885c: mul w13, w13, w9
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) * uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100488860 {
		return 0x100488860, true, nil
	}
	// 0x100488860: eor w13, w13, w8
	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[8])))))
	if runtime != nil && runtime.stopPC == 0x100488864 {
		return 0x100488864, true, nil
	}
	// 0x100488864: ldr x14, [sp, #0xdf0]
	value126, err := s.load(addOffset(s.sp, 3568), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(value126)
	if runtime != nil && runtime.stopPC == 0x100488868 {
		return 0x100488868, true, nil
	}
	// 0x100488868: ldr x13, [x14, w13, sxtw #3]
	value127, err := s.load((s.x[14] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(value127)
	if runtime != nil && runtime.stopPC == 0x10048886c {
		return 0x10048886c, true, nil
	}
	// 0x10048886c: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100488870 {
		return 0x100488870, true, nil
	}
	// 0x100488870: str x11, [sp, #0xd98]
	if err := s.storeValue(addOffset(s.sp, 3480), s.x[11], 8); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100488874 {
		return 0x100488874, true, nil
	}
	// 0x100488874: mov x17, #0xf758
	s.x[17] = uint64(uint64(0xf758))
	if runtime != nil && runtime.stopPC == 0x100488878 {
		return 0x100488878, true, nil
	}
	// 0x100488878: braa x13, x17
	next = s.x[13]
	return next, false, nil
}

func rawSourceBootstrapBlock1005ea024(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x1005ea024 {
		return 0x1005ea024, true, nil
	}
	// 0x1005ea024: mov w6, #0
	s.x[6] = uint64(uint32(uint64(0x0)))
	if runtime != nil && runtime.stopPC == 0x1005ea028 {
		return 0x1005ea028, true, nil
	}
	// 0x1005ea028: mov w9, #0xf899
	s.x[9] = uint64(uint32(uint64(0xf899)))
	if runtime != nil && runtime.stopPC == 0x1005ea02c {
		return 0x1005ea02c, true, nil
	}
	// 0x1005ea02c: eor w7, w10, w9
	s.x[7] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1005ea030 {
		return 0x1005ea030, true, nil
	}
	// 0x1005ea030: mov w9, #0x8dfd
	s.x[9] = uint64(uint32(uint64(0x8dfd)))
	if runtime != nil && runtime.stopPC == 0x1005ea034 {
		return 0x1005ea034, true, nil
	}
	// 0x1005ea034: movk w9, #0xc9d5, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xc9d5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea038 {
		return 0x1005ea038, true, nil
	}
	// 0x1005ea038: mov w11, #0xeb74
	s.x[11] = uint64(uint32(uint64(0xeb74)))
	if runtime != nil && runtime.stopPC == 0x1005ea03c {
		return 0x1005ea03c, true, nil
	}
	// 0x1005ea03c: movk w11, #0xbe86, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xbe86, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea040 {
		return 0x1005ea040, true, nil
	}
	// 0x1005ea040: add w9, w28, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[28])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x1005ea044 {
		return 0x1005ea044, true, nil
	}
	// 0x1005ea044: eor w11, w9, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea048 {
		return 0x1005ea048, true, nil
	}
	// 0x1005ea048: mov w12, #0x1de6
	s.x[12] = uint64(uint32(uint64(0x1de6)))
	if runtime != nil && runtime.stopPC == 0x1005ea04c {
		return 0x1005ea04c, true, nil
	}
	// 0x1005ea04c: movk w12, #0x19ab, lsl #16
	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x19ab, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea050 {
		return 0x1005ea050, true, nil
	}
	// 0x1005ea050: and w12, w9, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea054 {
		return 0x1005ea054, true, nil
	}
	// 0x1005ea054: mov w14, #0xd9f4
	s.x[14] = uint64(uint32(uint64(0xd9f4)))
	if runtime != nil && runtime.stopPC == 0x1005ea058 {
		return 0x1005ea058, true, nil
	}
	// 0x1005ea058: movk w14, #0x9d82, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x9d82, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea05c {
		return 0x1005ea05c, true, nil
	}
	// 0x1005ea05c: eor w14, w9, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea060 {
		return 0x1005ea060, true, nil
	}
	// 0x1005ea060: and w14, w14, w12, lsl #1
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) & shiftLeft(uint64(uint32(s.x[12])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1005ea064 {
		return 0x1005ea064, true, nil
	}
	// 0x1005ea064: eor w12, w14, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea068 {
		return 0x1005ea068, true, nil
	}
	// 0x1005ea068: lsl w14, w11, #1
	s.x[14] = uint64(uint32(shiftLeft(uint64(uint32(s.x[11])), int(uint64(0x1)), 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea06c {
		return 0x1005ea06c, true, nil
	}
	// 0x1005ea06c: mov w15, #0xed24
	s.x[15] = uint64(uint32(uint64(0xed24)))
	if runtime != nil && runtime.stopPC == 0x1005ea070 {
		return 0x1005ea070, true, nil
	}
	// 0x1005ea070: movk w15, #0x4e5b, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4e5b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea074 {
		return 0x1005ea074, true, nil
	}
	// 0x1005ea074: eor w15, w14, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea078 {
		return 0x1005ea078, true, nil
	}
	// 0x1005ea078: and w11, w15, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea07c {
		return 0x1005ea07c, true, nil
	}
	// 0x1005ea07c: mov w15, #0x5f93
	s.x[15] = uint64(uint32(uint64(0x5f93)))
	if runtime != nil && runtime.stopPC == 0x1005ea080 {
		return 0x1005ea080, true, nil
	}
	// 0x1005ea080: movk w15, #0xa72d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa72d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea084 {
		return 0x1005ea084, true, nil
	}
	// 0x1005ea084: add w15, w7, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[7])) + uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea088 {
		return 0x1005ea088, true, nil
	}
	// 0x1005ea088: and w14, w15, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea08c {
		return 0x1005ea08c, true, nil
	}
	// 0x1005ea08c: eor w11, w14, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea090 {
		return 0x1005ea090, true, nil
	}
	// 0x1005ea090: mov w14, #0x1292
	s.x[14] = uint64(uint32(uint64(0x1292)))
	if runtime != nil && runtime.stopPC == 0x1005ea094 {
		return 0x1005ea094, true, nil
	}
	// 0x1005ea094: movk w14, #0xa124, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xa124, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea098 {
		return 0x1005ea098, true, nil
	}
	// 0x1005ea098: eor w14, w11, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea09c {
		return 0x1005ea09c, true, nil
	}
	// 0x1005ea09c: mov w15, #0xe400
	s.x[15] = uint64(uint32(uint64(0xe400)))
	if runtime != nil && runtime.stopPC == 0x1005ea0a0 {
		return 0x1005ea0a0, true, nil
	}
	// 0x1005ea0a0: movk w15, #0x608, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x608, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0a4 {
		return 0x1005ea0a4, true, nil
	}
	// 0x1005ea0a4: eor w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0a8 {
		return 0x1005ea0a8, true, nil
	}
	// 0x1005ea0a8: and w11, w11, w12, lsl #2
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & shiftLeft(uint64(uint32(s.x[12])), 2, 32))))
	if runtime != nil && runtime.stopPC == 0x1005ea0ac {
		return 0x1005ea0ac, true, nil
	}
	// 0x1005ea0ac: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0b0 {
		return 0x1005ea0b0, true, nil
	}
	// 0x1005ea0b0: lsl w12, w14, #2
	s.x[12] = uint64(uint32(shiftLeft(uint64(uint32(s.x[14])), int(uint64(0x2)), 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0b4 {
		return 0x1005ea0b4, true, nil
	}
	// 0x1005ea0b4: mov w15, #0xda48
	s.x[15] = uint64(uint32(uint64(0xda48)))
	if runtime != nil && runtime.stopPC == 0x1005ea0b8 {
		return 0x1005ea0b8, true, nil
	}
	// 0x1005ea0b8: movk w15, #0x9cb7, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x9cb7, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0bc {
		return 0x1005ea0bc, true, nil
	}
	// 0x1005ea0bc: eor w15, w12, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0c0 {
		return 0x1005ea0c0, true, nil
	}
	// 0x1005ea0c0: and w14, w15, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0c4 {
		return 0x1005ea0c4, true, nil
	}
	// 0x1005ea0c4: mov w15, #0xf690
	s.x[15] = uint64(uint32(uint64(0xf690)))
	if runtime != nil && runtime.stopPC == 0x1005ea0c8 {
		return 0x1005ea0c8, true, nil
	}
	// 0x1005ea0c8: movk w15, #0xa72d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa72d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0cc {
		return 0x1005ea0cc, true, nil
	}
	// 0x1005ea0cc: and w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0d0 {
		return 0x1005ea0d0, true, nil
	}
	// 0x1005ea0d0: eor w12, w14, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0d4 {
		return 0x1005ea0d4, true, nil
	}
	// 0x1005ea0d4: mov w14, #0x2492
	s.x[14] = uint64(uint32(uint64(0x2492)))
	if runtime != nil && runtime.stopPC == 0x1005ea0d8 {
		return 0x1005ea0d8, true, nil
	}
	// 0x1005ea0d8: movk w14, #0x2308, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x2308, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0dc {
		return 0x1005ea0dc, true, nil
	}
	// 0x1005ea0dc: eor w14, w12, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0e0 {
		return 0x1005ea0e0, true, nil
	}
	// 0x1005ea0e0: mov w15, #0xd200
	s.x[15] = uint64(uint32(uint64(0xd200)))
	if runtime != nil && runtime.stopPC == 0x1005ea0e4 {
		return 0x1005ea0e4, true, nil
	}
	// 0x1005ea0e4: movk w15, #0x8425, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x8425, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0e8 {
		return 0x1005ea0e8, true, nil
	}
	// 0x1005ea0e8: eor w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0ec {
		return 0x1005ea0ec, true, nil
	}
	// 0x1005ea0ec: and w12, w12, w11, lsl #4
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & shiftLeft(uint64(uint32(s.x[11])), 4, 32))))
	if runtime != nil && runtime.stopPC == 0x1005ea0f0 {
		return 0x1005ea0f0, true, nil
	}
	// 0x1005ea0f0: eor w11, w12, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea0f4 {
		return 0x1005ea0f4, true, nil
	}
	// 0x1005ea0f4: lsl w12, w14, #4
	s.x[12] = uint64(uint32(shiftLeft(uint64(uint32(s.x[14])), int(uint64(0x4)), 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea0f8 {
		return 0x1005ea0f8, true, nil
	}
	// 0x1005ea0f8: mov w15, #0x6920
	s.x[15] = uint64(uint32(uint64(0x6920)))
	if runtime != nil && runtime.stopPC == 0x1005ea0fc {
		return 0x1005ea0fc, true, nil
	}
	// 0x1005ea0fc: movk w15, #0x72df, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x72df, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea100 {
		return 0x1005ea100, true, nil
	}
	// 0x1005ea100: eor w15, w12, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea104 {
		return 0x1005ea104, true, nil
	}
	// 0x1005ea104: and w14, w15, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea108 {
		return 0x1005ea108, true, nil
	}
	// 0x1005ea108: mov w15, #0xf680
	s.x[15] = uint64(uint32(uint64(0xf680)))
	if runtime != nil && runtime.stopPC == 0x1005ea10c {
		return 0x1005ea10c, true, nil
	}
	// 0x1005ea10c: movk w15, #0xa72d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa72d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea110 {
		return 0x1005ea110, true, nil
	}
	// 0x1005ea110: and w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea114 {
		return 0x1005ea114, true, nil
	}
	// 0x1005ea114: eor w12, w14, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea118 {
		return 0x1005ea118, true, nil
	}
	// 0x1005ea118: mov w14, #0x9692
	s.x[14] = uint64(uint32(uint64(0x9692)))
	if runtime != nil && runtime.stopPC == 0x1005ea11c {
		return 0x1005ea11c, true, nil
	}
	// 0x1005ea11c: movk w14, #0x8520, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x8520, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea120 {
		return 0x1005ea120, true, nil
	}
	// 0x1005ea120: eor w14, w12, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea124 {
		return 0x1005ea124, true, nil
	}
	// 0x1005ea124: mov w15, #0x6000
	s.x[15] = uint64(uint32(uint64(0x6000)))
	if runtime != nil && runtime.stopPC == 0x1005ea128 {
		return 0x1005ea128, true, nil
	}
	// 0x1005ea128: movk w15, #0x220d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x220d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea12c {
		return 0x1005ea12c, true, nil
	}
	// 0x1005ea12c: eor w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea130 {
		return 0x1005ea130, true, nil
	}
	// 0x1005ea130: and w12, w12, w11, lsl #8
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & shiftLeft(uint64(uint32(s.x[11])), 8, 32))))
	if runtime != nil && runtime.stopPC == 0x1005ea134 {
		return 0x1005ea134, true, nil
	}
	// 0x1005ea134: mov w15, #0xf692
	s.x[15] = uint64(uint32(uint64(0xf692)))
	if runtime != nil && runtime.stopPC == 0x1005ea138 {
		return 0x1005ea138, true, nil
	}
	// 0x1005ea138: movk w15, #0xa72d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa72d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea13c {
		return 0x1005ea13c, true, nil
	}
	// 0x1005ea13c: eor w11, w11, w15
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea140 {
		return 0x1005ea140, true, nil
	}
	// 0x1005ea140: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea144 {
		return 0x1005ea144, true, nil
	}
	// 0x1005ea144: lsl w12, w14, #8
	s.x[12] = uint64(uint32(shiftLeft(uint64(uint32(s.x[14])), int(uint64(0x8)), 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea148 {
		return 0x1005ea148, true, nil
	}
	// 0x1005ea148: mov w15, #0x2df60000
	s.x[15] = uint64(uint32(uint64(0x2df60000)))
	if runtime != nil && runtime.stopPC == 0x1005ea14c {
		return 0x1005ea14c, true, nil
	}
	// 0x1005ea14c: eor w15, w12, w15
	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea150 {
		return 0x1005ea150, true, nil
	}
	// 0x1005ea150: and w14, w15, w14
	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea154 {
		return 0x1005ea154, true, nil
	}
	// 0x1005ea154: mov w15, #0x272d0000
	s.x[15] = uint64(uint32(uint64(0x272d0000)))
	if runtime != nil && runtime.stopPC == 0x1005ea158 {
		return 0x1005ea158, true, nil
	}
	// 0x1005ea158: and w12, w12, w15
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea15c {
		return 0x1005ea15c, true, nil
	}
	// 0x1005ea15c: mov w16, #0x2090000
	s.x[16] = uint64(uint32(uint64(0x2090000)))
	if runtime != nil && runtime.stopPC == 0x1005ea160 {
		return 0x1005ea160, true, nil
	}
	// 0x1005ea160: eor w12, w12, w16
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1005ea164 {
		return 0x1005ea164, true, nil
	}
	// 0x1005ea164: eor w12, w12, w14
	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea168 {
		return 0x1005ea168, true, nil
	}
	// 0x1005ea168: lsl w14, w11, #0x10
	s.x[14] = uint64(uint32(shiftLeft(uint64(uint32(s.x[11])), int(uint64(0x10)), 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea16c {
		return 0x1005ea16c, true, nil
	}
	// 0x1005ea16c: mov w16, #0x76920000
	s.x[16] = uint64(uint32(uint64(0x76920000)))
	if runtime != nil && runtime.stopPC == 0x1005ea170 {
		return 0x1005ea170, true, nil
	}
	// 0x1005ea170: eor w16, w14, w16
	s.x[16] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x1005ea174 {
		return 0x1005ea174, true, nil
	}
	// 0x1005ea174: and w12, w16, w12
	s.x[12] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea178 {
		return 0x1005ea178, true, nil
	}
	// 0x1005ea178: and w14, w14, w15
	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) & uint64(uint32(s.x[15])))))
	if runtime != nil && runtime.stopPC == 0x1005ea17c {
		return 0x1005ea17c, true, nil
	}
	// 0x1005ea17c: eor w11, w14, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea180 {
		return 0x1005ea180, true, nil
	}
	// 0x1005ea180: eor w11, w11, w12
	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1005ea184 {
		return 0x1005ea184, true, nil
	}
	// 0x1005ea184: eor w9, w9, w11, lsl #1
	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ shiftLeft(uint64(uint32(s.x[11])), 1, 32))))
	if runtime != nil && runtime.stopPC == 0x1005ea188 {
		return 0x1005ea188, true, nil
	}
	// 0x1005ea188: mov w11, #0x75a2
	s.x[11] = uint64(uint32(uint64(0x75a2)))
	if runtime != nil && runtime.stopPC == 0x1005ea18c {
		return 0x1005ea18c, true, nil
	}
	// 0x1005ea18c: movk w11, #0xcdd3, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xcdd3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea190 {
		return 0x1005ea190, true, nil
	}
	// 0x1005ea190: eor w12, w9, w11
	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea194 {
		return 0x1005ea194, true, nil
	}
	// 0x1005ea194: mov w15, #0x9e3a
	s.x[15] = uint64(uint32(uint64(0x9e3a)))
	if runtime != nil && runtime.stopPC == 0x1005ea198 {
		return 0x1005ea198, true, nil
	}
	// 0x1005ea198: movk w15, #0x4cad, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4cad, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea19c {
		return 0x1005ea19c, true, nil
	}
	// 0x1005ea19c: mov w9, #0xf146
	s.x[9] = uint64(uint32(uint64(0xf146)))
	if runtime != nil && runtime.stopPC == 0x1005ea1a0 {
		return 0x1005ea1a0, true, nil
	}
	// 0x1005ea1a0: movk w9, #0x9c62, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x9c62, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1a4 {
		return 0x1005ea1a4, true, nil
	}
	// 0x1005ea1a4: str w9, [x8, #0x18c]
	if err := s.storeValue(addOffset(s.x[8], 396), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea1a8 {
		return 0x1005ea1a8, true, nil
	}
	// 0x1005ea1a8: mov w9, #0x6a71
	s.x[9] = uint64(uint32(uint64(0x6a71)))
	if runtime != nil && runtime.stopPC == 0x1005ea1ac {
		return 0x1005ea1ac, true, nil
	}
	// 0x1005ea1ac: movk w9, #0xbfe4, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xbfe4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1b0 {
		return 0x1005ea1b0, true, nil
	}
	// 0x1005ea1b0: str w9, [x8, #0x1fc]
	if err := s.storeValue(addOffset(s.x[8], 508), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea1b4 {
		return 0x1005ea1b4, true, nil
	}
	// 0x1005ea1b4: mov w9, #0x193a
	s.x[9] = uint64(uint32(uint64(0x193a)))
	if runtime != nil && runtime.stopPC == 0x1005ea1b8 {
		return 0x1005ea1b8, true, nil
	}
	// 0x1005ea1b8: movk w9, #0xe22b, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xe22b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1bc {
		return 0x1005ea1bc, true, nil
	}
	// 0x1005ea1bc: str w9, [x8, #0x154]
	if err := s.storeValue(addOffset(s.x[8], 340), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea1c0 {
		return 0x1005ea1c0, true, nil
	}
	// 0x1005ea1c0: mov w9, #0xf460
	s.x[9] = uint64(uint32(uint64(0xf460)))
	if runtime != nil && runtime.stopPC == 0x1005ea1c4 {
		return 0x1005ea1c4, true, nil
	}
	// 0x1005ea1c4: movk w9, #0x2e54, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x2e54, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1c8 {
		return 0x1005ea1c8, true, nil
	}
	// 0x1005ea1c8: str w9, [x8, #0x11c]
	if err := s.storeValue(addOffset(s.x[8], 284), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea1cc {
		return 0x1005ea1cc, true, nil
	}
	// 0x1005ea1cc: mov w9, #0xc9a3
	s.x[9] = uint64(uint32(uint64(0xc9a3)))
	if runtime != nil && runtime.stopPC == 0x1005ea1d0 {
		return 0x1005ea1d0, true, nil
	}
	// 0x1005ea1d0: movk w9, #0x1ca6, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x1ca6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1d4 {
		return 0x1005ea1d4, true, nil
	}
	// 0x1005ea1d4: str w9, [x8, #0x138]
	if err := s.storeValue(addOffset(s.x[8], 312), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea1d8 {
		return 0x1005ea1d8, true, nil
	}
	// 0x1005ea1d8: mov w16, #0x7549
	s.x[16] = uint64(uint32(uint64(0x7549)))
	if runtime != nil && runtime.stopPC == 0x1005ea1dc {
		return 0x1005ea1dc, true, nil
	}
	// 0x1005ea1dc: movk w16, #0xa5ff, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa5ff, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1e0 {
		return 0x1005ea1e0, true, nil
	}
	// 0x1005ea1e0: mov w9, #0xb333
	s.x[9] = uint64(uint32(uint64(0xb333)))
	if runtime != nil && runtime.stopPC == 0x1005ea1e4 {
		return 0x1005ea1e4, true, nil
	}
	// 0x1005ea1e4: movk w9, #0x7b5a, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7b5a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1e8 {
		return 0x1005ea1e8, true, nil
	}
	// 0x1005ea1e8: str w9, [x8, #0x218]
	if err := s.storeValue(addOffset(s.x[8], 536), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea1ec {
		return 0x1005ea1ec, true, nil
	}
	// 0x1005ea1ec: mov w9, #0xd863
	s.x[9] = uint64(uint32(uint64(0xd863)))
	if runtime != nil && runtime.stopPC == 0x1005ea1f0 {
		return 0x1005ea1f0, true, nil
	}
	// 0x1005ea1f0: movk w9, #0xeebe, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xeebe, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1f4 {
		return 0x1005ea1f4, true, nil
	}
	// 0x1005ea1f4: mov w17, #0xc68
	s.x[17] = uint64(uint32(uint64(0xc68)))
	if runtime != nil && runtime.stopPC == 0x1005ea1f8 {
		return 0x1005ea1f8, true, nil
	}
	// 0x1005ea1f8: movk w17, #0x19e5, lsl #16
	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x19e5, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea1fc {
		return 0x1005ea1fc, true, nil
	}
	// 0x1005ea1fc: mov w11, #0xee69
	s.x[11] = uint64(uint32(uint64(0xee69)))
	if runtime != nil && runtime.stopPC == 0x1005ea200 {
		return 0x1005ea200, true, nil
	}
	// 0x1005ea200: movk w11, #0xc97c, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xc97c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea204 {
		return 0x1005ea204, true, nil
	}
	// 0x1005ea204: str w11, [x8, #0x174]
	if err := s.storeValue(addOffset(s.x[8], 372), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea208 {
		return 0x1005ea208, true, nil
	}
	// 0x1005ea208: mov w11, #0x8093
	s.x[11] = uint64(uint32(uint64(0x8093)))
	if runtime != nil && runtime.stopPC == 0x1005ea20c {
		return 0x1005ea20c, true, nil
	}
	// 0x1005ea20c: movk w11, #0x6414, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x6414, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea210 {
		return 0x1005ea210, true, nil
	}
	// 0x1005ea210: str w11, [x8, #0x1a8]
	if err := s.storeValue(addOffset(s.x[8], 424), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea214 {
		return 0x1005ea214, true, nil
	}
	// 0x1005ea214: mov w14, #0x6016
	s.x[14] = uint64(uint32(uint64(0x6016)))
	if runtime != nil && runtime.stopPC == 0x1005ea218 {
		return 0x1005ea218, true, nil
	}
	// 0x1005ea218: movk w14, #0xcbbf, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xcbbf, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea21c {
		return 0x1005ea21c, true, nil
	}
	// 0x1005ea21c: mov w11, #0x966a
	s.x[11] = uint64(uint32(uint64(0x966a)))
	if runtime != nil && runtime.stopPC == 0x1005ea220 {
		return 0x1005ea220, true, nil
	}
	// 0x1005ea220: movk w11, #0xba26, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xba26, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea224 {
		return 0x1005ea224, true, nil
	}
	// 0x1005ea224: str w11, [x8, #0x190]
	if err := s.storeValue(addOffset(s.x[8], 400), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea228 {
		return 0x1005ea228, true, nil
	}
	// 0x1005ea228: mov w11, #0xcd77
	s.x[11] = uint64(uint32(uint64(0xcd77)))
	if runtime != nil && runtime.stopPC == 0x1005ea22c {
		return 0x1005ea22c, true, nil
	}
	// 0x1005ea22c: movk w11, #0xb53f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb53f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea230 {
		return 0x1005ea230, true, nil
	}
	// 0x1005ea230: str w11, [x8, #0x170]
	if err := s.storeValue(addOffset(s.x[8], 368), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea234 {
		return 0x1005ea234, true, nil
	}
	// 0x1005ea234: mov w0, #0xcc16
	s.x[0] = uint64(uint32(uint64(0xcc16)))
	if runtime != nil && runtime.stopPC == 0x1005ea238 {
		return 0x1005ea238, true, nil
	}
	// 0x1005ea238: movk w0, #0x7ce1, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x7ce1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea23c {
		return 0x1005ea23c, true, nil
	}
	// 0x1005ea23c: mov w11, #0xec65
	s.x[11] = uint64(uint32(uint64(0xec65)))
	if runtime != nil && runtime.stopPC == 0x1005ea240 {
		return 0x1005ea240, true, nil
	}
	// 0x1005ea240: movk w11, #0x5f2f, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x5f2f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea244 {
		return 0x1005ea244, true, nil
	}
	// 0x1005ea244: str w11, [x8, #0x100]
	if err := s.storeValue(addOffset(s.x[8], 256), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea248 {
		return 0x1005ea248, true, nil
	}
	// 0x1005ea248: mov w1, #0xb628
	s.x[1] = uint64(uint32(uint64(0xb628)))
	if runtime != nil && runtime.stopPC == 0x1005ea24c {
		return 0x1005ea24c, true, nil
	}
	// 0x1005ea24c: movk w1, #0xf9e1, lsl #16
	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xf9e1, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea250 {
		return 0x1005ea250, true, nil
	}
	// 0x1005ea250: mov w11, #0x106e
	s.x[11] = uint64(uint32(uint64(0x106e)))
	if runtime != nil && runtime.stopPC == 0x1005ea254 {
		return 0x1005ea254, true, nil
	}
	// 0x1005ea254: movk w11, #0x2e30, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x2e30, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea258 {
		return 0x1005ea258, true, nil
	}
	// 0x1005ea258: mov w2, #0xdc9e
	s.x[2] = uint64(uint32(uint64(0xdc9e)))
	if runtime != nil && runtime.stopPC == 0x1005ea25c {
		return 0x1005ea25c, true, nil
	}
	// 0x1005ea25c: movk w2, #0x574f, lsl #16
	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0x574f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea260 {
		return 0x1005ea260, true, nil
	}
	// 0x1005ea260: stp w16, w2, [x8, #0xac]
	if err := s.storeValue(addOffset(s.x[8], 172), uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 172))+4, uint64(uint32(s.x[2])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea264 {
		return 0x1005ea264, true, nil
	}
	// 0x1005ea264: mov w16, #0x7a95
	s.x[16] = uint64(uint32(uint64(0x7a95)))
	if runtime != nil && runtime.stopPC == 0x1005ea268 {
		return 0x1005ea268, true, nil
	}
	// 0x1005ea268: movk w16, #0xd103, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xd103, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea26c {
		return 0x1005ea26c, true, nil
	}
	// 0x1005ea26c: str w16, [x8, #0x158]
	if err := s.storeValue(addOffset(s.x[8], 344), uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea270 {
		return 0x1005ea270, true, nil
	}
	// 0x1005ea270: mov w16, #0x350f
	s.x[16] = uint64(uint32(uint64(0x350f)))
	if runtime != nil && runtime.stopPC == 0x1005ea274 {
		return 0x1005ea274, true, nil
	}
	// 0x1005ea274: movk w16, #0xb255, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xb255, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea278 {
		return 0x1005ea278, true, nil
	}
	// 0x1005ea278: str w16, [x8, #0x1ac]
	if err := s.storeValue(addOffset(s.x[8], 428), uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea27c {
		return 0x1005ea27c, true, nil
	}
	// 0x1005ea27c: mov w16, #0x30fd
	s.x[16] = uint64(uint32(uint64(0x30fd)))
	if runtime != nil && runtime.stopPC == 0x1005ea280 {
		return 0x1005ea280, true, nil
	}
	// 0x1005ea280: movk w16, #0x437a, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x437a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea284 {
		return 0x1005ea284, true, nil
	}
	// 0x1005ea284: str w16, [x8, #0x1c4]
	if err := s.storeValue(addOffset(s.x[8], 452), uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea288 {
		return 0x1005ea288, true, nil
	}
	// 0x1005ea288: mov w16, #0x187e
	s.x[16] = uint64(uint32(uint64(0x187e)))
	if runtime != nil && runtime.stopPC == 0x1005ea28c {
		return 0x1005ea28c, true, nil
	}
	// 0x1005ea28c: movk w16, #0x66d4, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x66d4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea290 {
		return 0x1005ea290, true, nil
	}
	// 0x1005ea290: str w16, [x8, #0x200]
	if err := s.storeValue(addOffset(s.x[8], 512), uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea294 {
		return 0x1005ea294, true, nil
	}
	// 0x1005ea294: mov w16, #0x7927
	s.x[16] = uint64(uint32(uint64(0x7927)))
	if runtime != nil && runtime.stopPC == 0x1005ea298 {
		return 0x1005ea298, true, nil
	}
	// 0x1005ea298: movk w16, #0x31da, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x31da, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea29c {
		return 0x1005ea29c, true, nil
	}
	// 0x1005ea29c: stp w0, w16, [x8, #0xe4]
	if err := s.storeValue(addOffset(s.x[8], 228), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 228))+4, uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2a0 {
		return 0x1005ea2a0, true, nil
	}
	// 0x1005ea2a0: mov w16, #0x20a0
	s.x[16] = uint64(uint32(uint64(0x20a0)))
	if runtime != nil && runtime.stopPC == 0x1005ea2a4 {
		return 0x1005ea2a4, true, nil
	}
	// 0x1005ea2a4: movk w16, #0x8b02, lsl #16
	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x8b02, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2a8 {
		return 0x1005ea2a8, true, nil
	}
	// 0x1005ea2a8: mov w0, #0xc045
	s.x[0] = uint64(uint32(uint64(0xc045)))
	if runtime != nil && runtime.stopPC == 0x1005ea2ac {
		return 0x1005ea2ac, true, nil
	}
	// 0x1005ea2ac: movk w0, #0xd18f, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xd18f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2b0 {
		return 0x1005ea2b0, true, nil
	}
	// 0x1005ea2b0: str w0, [x8, #0x1e8]
	if err := s.storeValue(addOffset(s.x[8], 488), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2b4 {
		return 0x1005ea2b4, true, nil
	}
	// 0x1005ea2b4: mov w0, #0x4a47
	s.x[0] = uint64(uint32(uint64(0x4a47)))
	if runtime != nil && runtime.stopPC == 0x1005ea2b8 {
		return 0x1005ea2b8, true, nil
	}
	// 0x1005ea2b8: movk w0, #0xf3d9, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xf3d9, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2bc {
		return 0x1005ea2bc, true, nil
	}
	// 0x1005ea2bc: str w0, [x8, #0x1e0]
	if err := s.storeValue(addOffset(s.x[8], 480), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2c0 {
		return 0x1005ea2c0, true, nil
	}
	// 0x1005ea2c0: mov w0, #0xd1f3
	s.x[0] = uint64(uint32(uint64(0xd1f3)))
	if runtime != nil && runtime.stopPC == 0x1005ea2c4 {
		return 0x1005ea2c4, true, nil
	}
	// 0x1005ea2c4: movk w0, #0x10ed, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x10ed, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2c8 {
		return 0x1005ea2c8, true, nil
	}
	// 0x1005ea2c8: str w0, [x8, #0x120]
	if err := s.storeValue(addOffset(s.x[8], 288), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2cc {
		return 0x1005ea2cc, true, nil
	}
	// 0x1005ea2cc: mov w0, #0x56bb
	s.x[0] = uint64(uint32(uint64(0x56bb)))
	if runtime != nil && runtime.stopPC == 0x1005ea2d0 {
		return 0x1005ea2d0, true, nil
	}
	// 0x1005ea2d0: movk w0, #0xdbb6, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xdbb6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2d4 {
		return 0x1005ea2d4, true, nil
	}
	// 0x1005ea2d4: str w0, [x8, #0x1cc]
	if err := s.storeValue(addOffset(s.x[8], 460), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2d8 {
		return 0x1005ea2d8, true, nil
	}
	// 0x1005ea2d8: mov w0, #0x9b11
	s.x[0] = uint64(uint32(uint64(0x9b11)))
	if runtime != nil && runtime.stopPC == 0x1005ea2dc {
		return 0x1005ea2dc, true, nil
	}
	// 0x1005ea2dc: movk w0, #0x380c, lsl #16
	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x380c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2e0 {
		return 0x1005ea2e0, true, nil
	}
	// 0x1005ea2e0: stp w15, w0, [x8]
	if err := s.storeValue(s.x[8], uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((s.x[8])+4, uint64(uint32(s.x[0])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2e4 {
		return 0x1005ea2e4, true, nil
	}
	// 0x1005ea2e4: mov w15, #0xf77c
	s.x[15] = uint64(uint32(uint64(0xf77c)))
	if runtime != nil && runtime.stopPC == 0x1005ea2e8 {
		return 0x1005ea2e8, true, nil
	}
	// 0x1005ea2e8: movk w15, #0x811c, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x811c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2ec {
		return 0x1005ea2ec, true, nil
	}
	// 0x1005ea2ec: str w15, [x8, #0x1e4]
	if err := s.storeValue(addOffset(s.x[8], 484), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2f0 {
		return 0x1005ea2f0, true, nil
	}
	// 0x1005ea2f0: mov w15, #0x9809
	s.x[15] = uint64(uint32(uint64(0x9809)))
	if runtime != nil && runtime.stopPC == 0x1005ea2f4 {
		return 0x1005ea2f4, true, nil
	}
	// 0x1005ea2f4: movk w15, #0x1e4c, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x1e4c, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea2f8 {
		return 0x1005ea2f8, true, nil
	}
	// 0x1005ea2f8: str w15, [x8, #0x104]
	if err := s.storeValue(addOffset(s.x[8], 260), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea2fc {
		return 0x1005ea2fc, true, nil
	}
	// 0x1005ea2fc: mov w15, #0x37d0
	s.x[15] = uint64(uint32(uint64(0x37d0)))
	if runtime != nil && runtime.stopPC == 0x1005ea300 {
		return 0x1005ea300, true, nil
	}
	// 0x1005ea300: movk w15, #0x8be6, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x8be6, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea304 {
		return 0x1005ea304, true, nil
	}
	// 0x1005ea304: stp w17, w15, [x8, #0x3c]
	if err := s.storeValue(addOffset(s.x[8], 60), uint64(uint32(s.x[17])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 60))+4, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea308 {
		return 0x1005ea308, true, nil
	}
	// 0x1005ea308: mov w15, #0xfa12
	s.x[15] = uint64(uint32(uint64(0xfa12)))
	if runtime != nil && runtime.stopPC == 0x1005ea30c {
		return 0x1005ea30c, true, nil
	}
	// 0x1005ea30c: movk w15, #0xcc6e, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xcc6e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea310 {
		return 0x1005ea310, true, nil
	}
	// 0x1005ea310: str w15, [x8, #0x204]
	if err := s.storeValue(addOffset(s.x[8], 516), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea314 {
		return 0x1005ea314, true, nil
	}
	// 0x1005ea314: mov w15, #0xbff4
	s.x[15] = uint64(uint32(uint64(0xbff4)))
	if runtime != nil && runtime.stopPC == 0x1005ea318 {
		return 0x1005ea318, true, nil
	}
	// 0x1005ea318: movk w15, #0xa261, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa261, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea31c {
		return 0x1005ea31c, true, nil
	}
	// 0x1005ea31c: str w15, [x8, #0x1c8]
	if err := s.storeValue(addOffset(s.x[8], 456), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea320 {
		return 0x1005ea320, true, nil
	}
	// 0x1005ea320: mov w15, #0x2fb8
	s.x[15] = uint64(uint32(uint64(0x2fb8)))
	if runtime != nil && runtime.stopPC == 0x1005ea324 {
		return 0x1005ea324, true, nil
	}
	// 0x1005ea324: movk w15, #0x23c8, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x23c8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea328 {
		return 0x1005ea328, true, nil
	}
	// 0x1005ea328: str w15, [x8, #0x15c]
	if err := s.storeValue(addOffset(s.x[8], 348), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea32c {
		return 0x1005ea32c, true, nil
	}
	// 0x1005ea32c: mov w15, #0xa427
	s.x[15] = uint64(uint32(uint64(0xa427)))
	if runtime != nil && runtime.stopPC == 0x1005ea330 {
		return 0x1005ea330, true, nil
	}
	// 0x1005ea330: movk w15, #0xfe31, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xfe31, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea334 {
		return 0x1005ea334, true, nil
	}
	// 0x1005ea334: str w15, [x8, #0x194]
	if err := s.storeValue(addOffset(s.x[8], 404), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea338 {
		return 0x1005ea338, true, nil
	}
	// 0x1005ea338: mov w15, #0x795d
	s.x[15] = uint64(uint32(uint64(0x795d)))
	if runtime != nil && runtime.stopPC == 0x1005ea33c {
		return 0x1005ea33c, true, nil
	}
	// 0x1005ea33c: movk w15, #0x808d, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x808d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea340 {
		return 0x1005ea340, true, nil
	}
	// 0x1005ea340: stp w1, w15, [x8, #0x58]
	if err := s.storeValue(addOffset(s.x[8], 88), uint64(uint32(s.x[1])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 88))+4, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea344 {
		return 0x1005ea344, true, nil
	}
	// 0x1005ea344: mov w15, #0xffe9
	s.x[15] = uint64(uint32(uint64(0xffe9)))
	if runtime != nil && runtime.stopPC == 0x1005ea348 {
		return 0x1005ea348, true, nil
	}
	// 0x1005ea348: movk w15, #0xdfca, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xdfca, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea34c {
		return 0x1005ea34c, true, nil
	}
	// 0x1005ea34c: str w15, [x8, #0x13c]
	if err := s.storeValue(addOffset(s.x[8], 316), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea350 {
		return 0x1005ea350, true, nil
	}
	// 0x1005ea350: mov w15, #0xd211
	s.x[15] = uint64(uint32(uint64(0xd211)))
	if runtime != nil && runtime.stopPC == 0x1005ea354 {
		return 0x1005ea354, true, nil
	}
	// 0x1005ea354: movk w15, #0x82c8, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x82c8, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea358 {
		return 0x1005ea358, true, nil
	}
	// 0x1005ea358: str w15, [x8, #0x60]
	if err := s.storeValue(addOffset(s.x[8], 96), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea35c {
		return 0x1005ea35c, true, nil
	}
	// 0x1005ea35c: mov w15, #0xf34c
	s.x[15] = uint64(uint32(uint64(0xf34c)))
	if runtime != nil && runtime.stopPC == 0x1005ea360 {
		return 0x1005ea360, true, nil
	}
	// 0x1005ea360: movk w15, #0x5497, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x5497, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea364 {
		return 0x1005ea364, true, nil
	}
	// 0x1005ea364: str w15, [x8, #0x21c]
	if err := s.storeValue(addOffset(s.x[8], 540), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea368 {
		return 0x1005ea368, true, nil
	}
	// 0x1005ea368: mov w15, #0xf190
	s.x[15] = uint64(uint32(uint64(0xf190)))
	if runtime != nil && runtime.stopPC == 0x1005ea36c {
		return 0x1005ea36c, true, nil
	}
	// 0x1005ea36c: movk w15, #0x7eb2, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x7eb2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea370 {
		return 0x1005ea370, true, nil
	}
	// 0x1005ea370: str w15, [x8, #0x7c]
	if err := s.storeValue(addOffset(s.x[8], 124), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea374 {
		return 0x1005ea374, true, nil
	}
	// 0x1005ea374: mov w15, #0xdaaa
	s.x[15] = uint64(uint32(uint64(0xdaaa)))
	if runtime != nil && runtime.stopPC == 0x1005ea378 {
		return 0x1005ea378, true, nil
	}
	// 0x1005ea378: movk w15, #0x285a, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x285a, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea37c {
		return 0x1005ea37c, true, nil
	}
	// 0x1005ea37c: str w15, [x8, #0x140]
	if err := s.storeValue(addOffset(s.x[8], 320), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea380 {
		return 0x1005ea380, true, nil
	}
	// 0x1005ea380: mov w15, #0xe5ee
	s.x[15] = uint64(uint32(uint64(0xe5ee)))
	if runtime != nil && runtime.stopPC == 0x1005ea384 {
		return 0x1005ea384, true, nil
	}
	// 0x1005ea384: movk w15, #0x316b, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x316b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea388 {
		return 0x1005ea388, true, nil
	}
	// 0x1005ea388: str w15, [x8, #0x124]
	if err := s.storeValue(addOffset(s.x[8], 292), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea38c {
		return 0x1005ea38c, true, nil
	}
	// 0x1005ea38c: mov w15, #0x7bb5
	s.x[15] = uint64(uint32(uint64(0x7bb5)))
	if runtime != nil && runtime.stopPC == 0x1005ea390 {
		return 0x1005ea390, true, nil
	}
	// 0x1005ea390: movk w15, #0x5d5b, lsl #16
	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x5d5b, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea394 {
		return 0x1005ea394, true, nil
	}
	// 0x1005ea394: stp w14, w15, [x8, #0x90]
	if err := s.storeValue(addOffset(s.x[8], 144), uint64(uint32(s.x[14])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 144))+4, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea398 {
		return 0x1005ea398, true, nil
	}
	// 0x1005ea398: mov w14, #0x7b55
	s.x[14] = uint64(uint32(uint64(0x7b55)))
	if runtime != nil && runtime.stopPC == 0x1005ea39c {
		return 0x1005ea39c, true, nil
	}
	// 0x1005ea39c: movk w14, #0x433f, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x433f, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3a0 {
		return 0x1005ea3a0, true, nil
	}
	// 0x1005ea3a0: str w14, [x8, #0xec]
	if err := s.storeValue(addOffset(s.x[8], 236), uint64(uint32(s.x[14])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3a4 {
		return 0x1005ea3a4, true, nil
	}
	// 0x1005ea3a4: mov w14, #0xe1f6
	s.x[14] = uint64(uint32(uint64(0xe1f6)))
	if runtime != nil && runtime.stopPC == 0x1005ea3a8 {
		return 0x1005ea3a8, true, nil
	}
	// 0x1005ea3a8: movk w14, #0x6f2d, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x6f2d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3ac {
		return 0x1005ea3ac, true, nil
	}
	// 0x1005ea3ac: stp w9, w14, [x8, #0x74]
	if err := s.storeValue(addOffset(s.x[8], 116), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 116))+4, uint64(uint32(s.x[14])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3b0 {
		return 0x1005ea3b0, true, nil
	}
	// 0x1005ea3b0: mov w9, #0xd4ed
	s.x[9] = uint64(uint32(uint64(0xd4ed)))
	if runtime != nil && runtime.stopPC == 0x1005ea3b4 {
		return 0x1005ea3b4, true, nil
	}
	// 0x1005ea3b4: movk w9, #0xb23d, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xb23d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3b8 {
		return 0x1005ea3b8, true, nil
	}
	// 0x1005ea3b8: str w9, [x8, #0x220]
	if err := s.storeValue(addOffset(s.x[8], 544), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3bc {
		return 0x1005ea3bc, true, nil
	}
	// 0x1005ea3bc: mov w9, #0x18ec
	s.x[9] = uint64(uint32(uint64(0x18ec)))
	if runtime != nil && runtime.stopPC == 0x1005ea3c0 {
		return 0x1005ea3c0, true, nil
	}
	// 0x1005ea3c0: movk w9, #0x58fc, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x58fc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3c4 {
		return 0x1005ea3c4, true, nil
	}
	// 0x1005ea3c4: str w9, [x8, #0xd0]
	if err := s.storeValue(addOffset(s.x[8], 208), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3c8 {
		return 0x1005ea3c8, true, nil
	}
	// 0x1005ea3c8: mov w9, #0x8fda
	s.x[9] = uint64(uint32(uint64(0x8fda)))
	if runtime != nil && runtime.stopPC == 0x1005ea3cc {
		return 0x1005ea3cc, true, nil
	}
	// 0x1005ea3cc: movk w9, #0x383e, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x383e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3d0 {
		return 0x1005ea3d0, true, nil
	}
	// 0x1005ea3d0: stp w16, w9, [x8, #0xc8]
	if err := s.storeValue(addOffset(s.x[8], 200), uint64(uint32(s.x[16])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 200))+4, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3d4 {
		return 0x1005ea3d4, true, nil
	}
	// 0x1005ea3d4: mov w9, #0x2b57
	s.x[9] = uint64(uint32(uint64(0x2b57)))
	if runtime != nil && runtime.stopPC == 0x1005ea3d8 {
		return 0x1005ea3d8, true, nil
	}
	// 0x1005ea3d8: movk w9, #0x3ed4, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x3ed4, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3dc {
		return 0x1005ea3dc, true, nil
	}
	// 0x1005ea3dc: str w9, [x8, #0x108]
	if err := s.storeValue(addOffset(s.x[8], 264), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3e0 {
		return 0x1005ea3e0, true, nil
	}
	// 0x1005ea3e0: mov w9, #0x38f1
	s.x[9] = uint64(uint32(uint64(0x38f1)))
	if runtime != nil && runtime.stopPC == 0x1005ea3e4 {
		return 0x1005ea3e4, true, nil
	}
	// 0x1005ea3e4: movk w9, #0x65c0, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x65c0, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3e8 {
		return 0x1005ea3e8, true, nil
	}
	// 0x1005ea3e8: str w9, [x8, #0xb4]
	if err := s.storeValue(addOffset(s.x[8], 180), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3ec {
		return 0x1005ea3ec, true, nil
	}
	// 0x1005ea3ec: mov w9, #0xde95
	s.x[9] = uint64(uint32(uint64(0xde95)))
	if runtime != nil && runtime.stopPC == 0x1005ea3f0 {
		return 0x1005ea3f0, true, nil
	}
	// 0x1005ea3f0: movk w9, #0x7251, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7251, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea3f4 {
		return 0x1005ea3f4, true, nil
	}
	// 0x1005ea3f4: str w9, [x8, #0x98]
	if err := s.storeValue(addOffset(s.x[8], 152), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea3f8 {
		return 0x1005ea3f8, true, nil
	}
	// 0x1005ea3f8: mov w9, #0xf12
	s.x[9] = uint64(uint32(uint64(0xf12)))
	if runtime != nil && runtime.stopPC == 0x1005ea3fc {
		return 0x1005ea3fc, true, nil
	}
	// 0x1005ea3fc: movk w9, #0x1a8e, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x1a8e, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea400 {
		return 0x1005ea400, true, nil
	}
	// 0x1005ea400: str w9, [x8, #0x178]
	if err := s.storeValue(addOffset(s.x[8], 376), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea404 {
		return 0x1005ea404, true, nil
	}
	// 0x1005ea404: mov w9, #0x74c6
	s.x[9] = uint64(uint32(uint64(0x74c6)))
	if runtime != nil && runtime.stopPC == 0x1005ea408 {
		return 0x1005ea408, true, nil
	}
	// 0x1005ea408: movk w9, #0xe0bc, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xe0bc, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea40c {
		return 0x1005ea40c, true, nil
	}
	// 0x1005ea40c: str w9, [x8, #0x1b0]
	if err := s.storeValue(addOffset(s.x[8], 432), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea410 {
		return 0x1005ea410, true, nil
	}
	// 0x1005ea410: mov w9, #0x9b27
	s.x[9] = uint64(uint32(uint64(0x9b27)))
	if runtime != nil && runtime.stopPC == 0x1005ea414 {
		return 0x1005ea414, true, nil
	}
	// 0x1005ea414: movk w9, #0x9c95, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x9c95, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea418 {
		return 0x1005ea418, true, nil
	}
	// 0x1005ea418: stp w11, w9, [x8, #0x20]
	if err := s.storeValue(addOffset(s.x[8], 32), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 32))+4, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea41c {
		return 0x1005ea41c, true, nil
	}
	// 0x1005ea41c: mov w9, #0x70aa
	s.x[9] = uint64(uint32(uint64(0x70aa)))
	if runtime != nil && runtime.stopPC == 0x1005ea420 {
		return 0x1005ea420, true, nil
	}
	// 0x1005ea420: movk w9, #0x96ea, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x96ea, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea424 {
		return 0x1005ea424, true, nil
	}
	// 0x1005ea424: str w9, [x8, #0x28]
	if err := s.storeValue(addOffset(s.x[8], 40), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea428 {
		return 0x1005ea428, true, nil
	}
	// 0x1005ea428: mov w9, #0xe794
	s.x[9] = uint64(uint32(uint64(0xe794)))
	if runtime != nil && runtime.stopPC == 0x1005ea42c {
		return 0x1005ea42c, true, nil
	}
	// 0x1005ea42c: movk w9, #0xd094, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xd094, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea430 {
		return 0x1005ea430, true, nil
	}
	// 0x1005ea430: str w9, [x8, #0x1ec]
	if err := s.storeValue(addOffset(s.x[8], 492), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea434 {
		return 0x1005ea434, true, nil
	}
	// 0x1005ea434: mov w9, #0xdb8d
	s.x[9] = uint64(uint32(uint64(0xdb8d)))
	if runtime != nil && runtime.stopPC == 0x1005ea438 {
		return 0x1005ea438, true, nil
	}
	// 0x1005ea438: movk w9, #0xae73, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xae73, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea43c {
		return 0x1005ea43c, true, nil
	}
	// 0x1005ea43c: str w9, [x8, #0x224]
	if err := s.storeValue(addOffset(s.x[8], 548), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea440 {
		return 0x1005ea440, true, nil
	}
	// 0x1005ea440: mov w9, #0x5743
	s.x[9] = uint64(uint32(uint64(0x5743)))
	if runtime != nil && runtime.stopPC == 0x1005ea444 {
		return 0x1005ea444, true, nil
	}
	// 0x1005ea444: movk w9, #0xace3, lsl #16
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xace3, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea448 {
		return 0x1005ea448, true, nil
	}
	// 0x1005ea448: mov w11, #0x1299
	s.x[11] = uint64(uint32(uint64(0x1299)))
	if runtime != nil && runtime.stopPC == 0x1005ea44c {
		return 0x1005ea44c, true, nil
	}
	// 0x1005ea44c: movk w11, #0xc91d, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xc91d, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea450 {
		return 0x1005ea450, true, nil
	}
	// 0x1005ea450: str w11, [x8, #0x208]
	if err := s.storeValue(addOffset(s.x[8], 520), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea454 {
		return 0x1005ea454, true, nil
	}
	// 0x1005ea454: mov w11, #0xfa28
	s.x[11] = uint64(uint32(uint64(0xfa28)))
	if runtime != nil && runtime.stopPC == 0x1005ea458 {
		return 0x1005ea458, true, nil
	}
	// 0x1005ea458: movk w11, #0x8969, lsl #16
	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x8969, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea45c {
		return 0x1005ea45c, true, nil
	}
	// 0x1005ea45c: ldr w13, [sp, #0x18c]
	value270, err := s.load(addOffset(s.sp, 396), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[13] = uint64(uint32(value270))
	if runtime != nil && runtime.stopPC == 0x1005ea460 {
		return 0x1005ea460, true, nil
	}
	// 0x1005ea460: add w11, w13, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1005ea464 {
		return 0x1005ea464, true, nil
	}
	// 0x1005ea464: mov w14, #0x93cd
	s.x[14] = uint64(uint32(uint64(0x93cd)))
	if runtime != nil && runtime.stopPC == 0x1005ea468 {
		return 0x1005ea468, true, nil
	}
	// 0x1005ea468: movk w14, #0xb5fe, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xb5fe, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea46c {
		return 0x1005ea46c, true, nil
	}
	// 0x1005ea46c: add w24, w26, w14
	s.x[24] = uint64(uint32((uint64(uint32(s.x[26])) + uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1005ea470 {
		return 0x1005ea470, true, nil
	}
	// 0x1005ea470: mov w14, #0xe6b6
	s.x[14] = uint64(uint32(uint64(0xe6b6)))
	if runtime != nil && runtime.stopPC == 0x1005ea474 {
		return 0x1005ea474, true, nil
	}
	// 0x1005ea474: movk w14, #0x8eeb, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x8eeb, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea478 {
		return 0x1005ea478, true, nil
	}
	// 0x1005ea478: str w14, [x8, #0x44]
	if err := s.storeValue(addOffset(s.x[8], 68), uint64(uint32(s.x[14])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea47c {
		return 0x1005ea47c, true, nil
	}
	// 0x1005ea47c: mov w14, #0x6e9c
	s.x[14] = uint64(uint32(uint64(0x6e9c)))
	if runtime != nil && runtime.stopPC == 0x1005ea480 {
		return 0x1005ea480, true, nil
	}
	// 0x1005ea480: movk w14, #0x9ae2, lsl #16
	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x9ae2, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea484 {
		return 0x1005ea484, true, nil
	}
	// 0x1005ea484: stp w9, w14, [x8, #8]
	if err := s.storeValue(addOffset(s.x[8], 8), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, false, err
	}
	if err := s.storeValue((addOffset(s.x[8], 8))+4, uint64(uint32(s.x[14])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x1005ea488 {
		return 0x1005ea488, true, nil
	}
	// 0x1005ea488: ldr x10, [x30, w10, sxtw #3]
	value281, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value281)
	if runtime != nil && runtime.stopPC == 0x1005ea48c {
		return 0x1005ea48c, true, nil
	}
	// 0x1005ea48c: mov w27, #0x3f
	s.x[27] = uint64(uint32(uint64(0x3f)))
	if runtime != nil && runtime.stopPC == 0x1005ea490 {
		return 0x1005ea490, true, nil
	}
	// 0x1005ea490: mov w28, #0x2312
	s.x[28] = uint64(uint32(uint64(0x2312)))
	if runtime != nil && runtime.stopPC == 0x1005ea494 {
		return 0x1005ea494, true, nil
	}
	// 0x1005ea494: mov x13, x30
	s.x[13] = uint64(s.x[30])
	if runtime != nil && runtime.stopPC == 0x1005ea498 {
		return 0x1005ea498, true, nil
	}
	// 0x1005ea498: mov w30, #0x6b1f
	s.x[30] = uint64(uint32(uint64(0x6b1f)))
	if runtime != nil && runtime.stopPC == 0x1005ea49c {
		return 0x1005ea49c, true, nil
	}
	// 0x1005ea49c: movk w30, #0x3a19, lsl #16
	s.x[30] = uint64(uint32(insertBits(uint64(uint32(s.x[30])), 0x3a19, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea4a0 {
		return 0x1005ea4a0, true, nil
	}
	// 0x1005ea4a0: mov w1, #0x4e390000
	s.x[1] = uint64(uint32(uint64(0x4e390000)))
	if runtime != nil && runtime.stopPC == 0x1005ea4a4 {
		return 0x1005ea4a4, true, nil
	}
	// 0x1005ea4a4: mov w2, #0x8b5b
	s.x[2] = uint64(uint32(uint64(0x8b5b)))
	if runtime != nil && runtime.stopPC == 0x1005ea4a8 {
		return 0x1005ea4a8, true, nil
	}
	// 0x1005ea4a8: movk w2, #0xb975, lsl #16
	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0xb975, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea4ac {
		return 0x1005ea4ac, true, nil
	}
	// 0x1005ea4ac: mov w4, #0x46db
	s.x[4] = uint64(uint32(uint64(0x46db)))
	if runtime != nil && runtime.stopPC == 0x1005ea4b0 {
		return 0x1005ea4b0, true, nil
	}
	// 0x1005ea4b0: movk w4, #0x5898, lsl #16
	s.x[4] = uint64(uint32(insertBits(uint64(uint32(s.x[4])), 0x5898, 16, 16, 32)))
	if runtime != nil && runtime.stopPC == 0x1005ea4b4 {
		return 0x1005ea4b4, true, nil
	}
	// 0x1005ea4b4: mov x19, x21
	s.x[19] = uint64(s.x[21])
	if runtime != nil && runtime.stopPC == 0x1005ea4b8 {
		return 0x1005ea4b8, true, nil
	}
	// 0x1005ea4b8: ldr w21, [sp, #0x1a8]
	value293, err := s.load(addOffset(s.sp, 424), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[21] = uint64(uint32(value293))
	if runtime != nil && runtime.stopPC == 0x1005ea4bc {
		return 0x1005ea4bc, true, nil
	}
	// 0x1005ea4bc: ldr w23, [sp, #0x190]
	value294, err := s.load(addOffset(s.sp, 400), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[23] = uint64(uint32(value294))
	if runtime != nil && runtime.stopPC == 0x1005ea4c0 {
		return 0x1005ea4c0, true, nil
	}
	// 0x1005ea4c0: ldr w14, [sp, #0x194]
	value295, err := s.load(addOffset(s.sp, 404), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[14] = uint64(uint32(value295))
	if runtime != nil && runtime.stopPC == 0x1005ea4c4 {
		return 0x1005ea4c4, true, nil
	}
	// 0x1005ea4c4: mov x17, #0x9e6d
	s.x[17] = uint64(uint64(0x9e6d))
	if runtime != nil && runtime.stopPC == 0x1005ea4c8 {
		return 0x1005ea4c8, true, nil
	}
	// 0x1005ea4c8: braa x10, x17
	next = s.x[10]
	return next, false, nil
}

func rawSourceBootstrapBlock100a06948(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x100a06948 {
		return 0x100a06948, true, nil
	}
	// 0x100a06948: mov w5, #0
	s.x[5] = uint64(uint32(uint64(0x0)))
	if runtime != nil && runtime.stopPC == 0x100a0694c {
		return 0x100a0694c, true, nil
	}
	// 0x100a0694c: mov w9, #-0x34f3
	s.x[9] = uint64(uint32(uint64(0xffffcb0d)))
	if runtime != nil && runtime.stopPC == 0x100a06950 {
		return 0x100a06950, true, nil
	}
	// 0x100a06950: mov w10, #0x5f47
	s.x[10] = uint64(uint32(uint64(0x5f47)))
	if runtime != nil && runtime.stopPC == 0x100a06954 {
		return 0x100a06954, true, nil
	}
	// 0x100a06954: add w9, w12, w9
	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[9])))))
	if runtime != nil && runtime.stopPC == 0x100a06958 {
		return 0x100a06958, true, nil
	}
	// 0x100a06958: add w10, w9, w10
	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))
	if runtime != nil && runtime.stopPC == 0x100a0695c {
		return 0x100a0695c, true, nil
	}
	// 0x100a0695c: ldr x12, [x30, w9, sxtw #3]
	value5, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[12] = uint64(value5)
	if runtime != nil && runtime.stopPC == 0x100a06960 {
		return 0x100a06960, true, nil
	}
	// 0x100a06960: mov w11, #0x200a
	s.x[11] = uint64(uint32(uint64(0x200a)))
	if runtime != nil && runtime.stopPC == 0x100a06964 {
		return 0x100a06964, true, nil
	}
	// 0x100a06964: orr w11, w10, w11
	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) | uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x100a06968 {
		return 0x100a06968, true, nil
	}
	// 0x100a06968: mov w10, #0x94
	s.x[10] = uint64(uint32(uint64(0x94)))
	if runtime != nil && runtime.stopPC == 0x100a0696c {
		return 0x100a0696c, true, nil
	}
	// 0x100a0696c: mov w14, #0xe071
	s.x[14] = uint64(uint32(uint64(0xe071)))
	if runtime != nil && runtime.stopPC == 0x100a06970 {
		return 0x100a06970, true, nil
	}
	// 0x100a06970: mov w15, #0x49
	s.x[15] = uint64(uint32(uint64(0x49)))
	if runtime != nil && runtime.stopPC == 0x100a06974 {
		return 0x100a06974, true, nil
	}
	// 0x100a06974: mov w16, #-0xa43a
	s.x[16] = uint64(uint32(uint64(0xffff5bc6)))
	if runtime != nil && runtime.stopPC == 0x100a06978 {
		return 0x100a06978, true, nil
	}
	// 0x100a06978: mov w0, #0x86
	s.x[0] = uint64(uint32(uint64(0x86)))
	if runtime != nil && runtime.stopPC == 0x100a0697c {
		return 0x100a0697c, true, nil
	}
	// 0x100a0697c: mov w1, #0x2312
	s.x[1] = uint64(uint32(uint64(0x2312)))
	if runtime != nil && runtime.stopPC == 0x100a06980 {
		return 0x100a06980, true, nil
	}
	// 0x100a06980: mov w2, #0x7433
	s.x[2] = uint64(uint32(uint64(0x7433)))
	if runtime != nil && runtime.stopPC == 0x100a06984 {
		return 0x100a06984, true, nil
	}
	// 0x100a06984: mov w3, #0x8d
	s.x[3] = uint64(uint32(uint64(0x8d)))
	if runtime != nil && runtime.stopPC == 0x100a06988 {
		return 0x100a06988, true, nil
	}
	// 0x100a06988: mov w4, #0x80d1
	s.x[4] = uint64(uint32(uint64(0x80d1)))
	if runtime != nil && runtime.stopPC == 0x100a0698c {
		return 0x100a0698c, true, nil
	}
	// 0x100a0698c: str w26, [sp, #0x198]
	if err := s.storeValue(addOffset(s.sp, 408), uint64(uint32(s.x[26])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100a06990 {
		return 0x100a06990, true, nil
	}
	// 0x100a06990: str w21, [sp, #0x188]
	if err := s.storeValue(addOffset(s.sp, 392), uint64(uint32(s.x[21])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100a06994 {
		return 0x100a06994, true, nil
	}
	// 0x100a06994: str w25, [sp, #0x184]
	if err := s.storeValue(addOffset(s.sp, 388), uint64(uint32(s.x[25])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100a06998 {
		return 0x100a06998, true, nil
	}
	// 0x100a06998: str w28, [sp, #0x178]
	if err := s.storeValue(addOffset(s.sp, 376), uint64(uint32(s.x[28])), 4); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100a0699c {
		return 0x100a0699c, true, nil
	}
	// 0x100a0699c: mov x17, #0x3e0f
	s.x[17] = uint64(uint64(0x3e0f))
	if runtime != nil && runtime.stopPC == 0x100a069a0 {
		return 0x100a069a0, true, nil
	}
	// 0x100a069a0: braa x12, x17
	next = s.x[12]
	return next, false, nil
}
