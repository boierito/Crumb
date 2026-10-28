// REVIEW CANDIDATE: not registered for execution or release.
// Exact daemon SHA256: 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784
// Translator revision: f7242cce58077a4af8db738ed295ee23ad8f067d
package rawsource

func rawSourceBootstrapBlock1001005ac(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
	if runtime != nil && runtime.stopPC == 0x1001005ac {
		return 0x1001005ac, true, nil
	}
	// 0x1001005ac: subs x9, x9, #1
	s.x[9] = uint64(s.subtractWithFlags(s.x[9], uint64(0x1), 64))
	if runtime != nil && runtime.stopPC == 0x1001005b0 {
		return 0x1001005b0, true, nil
	}
	// 0x1001005b0: ldrb w10, [x28, x9]
	value1, err := s.load((s.x[28] + s.x[9]), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(uint32(value1))
	if runtime != nil && runtime.stopPC == 0x1001005b4 {
		return 0x1001005b4, true, nil
	}
	// 0x1001005b4: adrp x17, #0x101075000
	s.x[17] = uint64(0x101075000)
	if runtime != nil && runtime.stopPC == 0x1001005b8 {
		return 0x1001005b8, true, nil
	}
	// 0x1001005b8: ldr w17, [x17, #0x588]
	value3, err := s.load(addOffset(s.x[17], 1416), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value3))
	if runtime != nil && runtime.stopPC == 0x1001005bc {
		return 0x1001005bc, true, nil
	}
	// 0x1001005bc: ldr w1, [x23, #0xb40]
	value4, err := s.load(addOffset(s.x[23], 2880), 4)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value4))
	if runtime != nil && runtime.stopPC == 0x1001005c0 {
		return 0x1001005c0, true, nil
	}
	// 0x1001005c0: sub w1, w1, w17
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1001005c4 {
		return 0x1001005c4, true, nil
	}
	// 0x1001005c4: eor w1, w1, w12
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1001005c8 {
		return 0x1001005c8, true, nil
	}
	// 0x1001005c8: mul w1, w1, w13
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) * uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1001005cc {
		return 0x1001005cc, true, nil
	}
	// 0x1001005cc: add x7, x30, w1, uxtb
	s.x[7] = uint64((s.x[30] + uint64(uint8(uint64(uint32(s.x[1]))))))
	if runtime != nil && runtime.stopPC == 0x1001005d0 {
		return 0x1001005d0, true, nil
	}
	// 0x1001005d0: ldrb w7, [x7]
	value9, err := s.load(s.x[7], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value9))
	if runtime != nil && runtime.stopPC == 0x1001005d4 {
		return 0x1001005d4, true, nil
	}
	// 0x1001005d4: eor w7, w7, w14
	s.x[7] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[14])))))
	if runtime != nil && runtime.stopPC == 0x1001005d8 {
		return 0x1001005d8, true, nil
	}
	// 0x1001005d8: ldrb w7, [x27, w7, uxtw]
	value11, err := s.load((s.x[27] + uint64(uint32(uint64(uint32(s.x[7]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value11))
	if runtime != nil && runtime.stopPC == 0x1001005dc {
		return 0x1001005dc, true, nil
	}
	// 0x1001005dc: eor w1, w1, w7
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x1001005e0 {
		return 0x1001005e0, true, nil
	}
	// 0x1001005e0: and w1, w1, #0xff
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x1001005e4 {
		return 0x1001005e4, true, nil
	}
	// 0x1001005e4: ldur x7, [x29, #-0xa0]
	value14, err := s.load(addOffset(s.x[29], -160), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(value14)
	if runtime != nil && runtime.stopPC == 0x1001005e8 {
		return 0x1001005e8, true, nil
	}
	// 0x1001005e8: ldr w7, [x7]
	value15, err := s.load(s.x[7], 4)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value15))
	if runtime != nil && runtime.stopPC == 0x1001005ec {
		return 0x1001005ec, true, nil
	}
	// 0x1001005ec: add w19, w7, w17
	s.x[19] = uint64(uint32((uint64(uint32(s.x[7])) + uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x1001005f0 {
		return 0x1001005f0, true, nil
	}
	// 0x1001005f0: eor w19, w19, w12
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x1001005f4 {
		return 0x1001005f4, true, nil
	}
	// 0x1001005f4: mul w19, w19, w13
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) * uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x1001005f8 {
		return 0x1001005f8, true, nil
	}
	// 0x1001005f8: add x20, x8, w19, uxtb
	s.x[20] = uint64((s.x[8] + uint64(uint8(uint64(uint32(s.x[19]))))))
	if runtime != nil && runtime.stopPC == 0x1001005fc {
		return 0x1001005fc, true, nil
	}
	// 0x1001005fc: add w1, w1, #0x3e2
	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) + uint64(0x3e2))))
	if runtime != nil && runtime.stopPC == 0x100100600 {
		return 0x100100600, true, nil
	}
	// 0x100100600: ldrb w20, [x20]
	value21, err := s.load(s.x[20], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[20] = uint64(uint32(value21))
	if runtime != nil && runtime.stopPC == 0x100100604 {
		return 0x100100604, true, nil
	}
	// 0x100100604: eor w20, w20, w16
	s.x[20] = uint64(uint32((uint64(uint32(s.x[20])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100100608 {
		return 0x100100608, true, nil
	}
	// 0x100100608: ldrb w20, [x25, w20, uxtw]
	value23, err := s.load((s.x[25] + uint64(uint32(uint64(uint32(s.x[20]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[20] = uint64(uint32(value23))
	if runtime != nil && runtime.stopPC == 0x10010060c {
		return 0x10010060c, true, nil
	}
	// 0x10010060c: adrp x28, #0x10105b000
	s.x[28] = uint64(0x10105b000)
	if runtime != nil && runtime.stopPC == 0x100100610 {
		return 0x100100610, true, nil
	}
	// 0x100100610: add x28, x28, #0x6f0
	s.x[28] = uint64((s.x[28] + uint64(0x6f0)))
	if runtime != nil && runtime.stopPC == 0x100100614 {
		return 0x100100614, true, nil
	}
	// 0x100100614: ldr x1, [x28, w1, uxtw #3]
	value26, err := s.load((s.x[28] + shiftLeft(uint64(uint32(uint64(uint32(s.x[1])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(value26)
	if runtime != nil && runtime.stopPC == 0x100100618 {
		return 0x100100618, true, nil
	}
	// 0x100100618: eor w19, w19, w20
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[20])))))
	if runtime != nil && runtime.stopPC == 0x10010061c {
		return 0x10010061c, true, nil
	}
	// 0x10010061c: and w19, w19, #0xff
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) & uint64(0xff))))
	if runtime != nil && runtime.stopPC == 0x100100620 {
		return 0x100100620, true, nil
	}
	// 0x100100620: add w19, w19, #0xcd8
	s.x[19] = uint64(uint32((uint64(uint32(s.x[19])) + uint64(0xcd8))))
	if runtime != nil && runtime.stopPC == 0x100100624 {
		return 0x100100624, true, nil
	}
	// 0x100100624: ldr x19, [x28, w19, uxtw #3]
	value30, err := s.load((s.x[28] + shiftLeft(uint64(uint32(uint64(uint32(s.x[19])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[19] = uint64(value30)
	if runtime != nil && runtime.stopPC == 0x100100628 {
		return 0x100100628, true, nil
	}
	// 0x100100628: and x20, x9, #0xf
	s.x[20] = uint64((s.x[9] & uint64(0xf)))
	if runtime != nil && runtime.stopPC == 0x10010062c {
		return 0x10010062c, true, nil
	}
	// 0x10010062c: sub w17, w7, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[7])) - uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100100630 {
		return 0x100100630, true, nil
	}
	// 0x100100630: eor w17, w17, w12
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))
	if runtime != nil && runtime.stopPC == 0x100100634 {
		return 0x100100634, true, nil
	}
	// 0x100100634: mul w17, w17, w13
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[13])))))
	if runtime != nil && runtime.stopPC == 0x100100638 {
		return 0x100100638, true, nil
	}
	// 0x100100638: add x7, x8, w17, uxtb
	s.x[7] = uint64((s.x[8] + uint64(uint8(uint64(uint32(s.x[17]))))))
	if runtime != nil && runtime.stopPC == 0x10010063c {
		return 0x10010063c, true, nil
	}
	// 0x10010063c: add x21, x20, x0
	s.x[21] = uint64((s.x[20] + s.x[0]))
	if runtime != nil && runtime.stopPC == 0x100100640 {
		return 0x100100640, true, nil
	}
	// 0x100100640: ldrb w7, [x7]
	value37, err := s.load(s.x[7], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value37))
	if runtime != nil && runtime.stopPC == 0x100100644 {
		return 0x100100644, true, nil
	}
	// 0x100100644: eor w7, w7, w16
	s.x[7] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[16])))))
	if runtime != nil && runtime.stopPC == 0x100100648 {
		return 0x100100648, true, nil
	}
	// 0x100100648: ldrb w7, [x25, w7, uxtw]
	value39, err := s.load((s.x[25] + uint64(uint32(uint64(uint32(s.x[7]))))), 1)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value39))
	if runtime != nil && runtime.stopPC == 0x10010064c {
		return 0x10010064c, true, nil
	}
	// 0x10010064c: add x22, x20, x15
	s.x[22] = uint64((s.x[20] + s.x[15]))
	if runtime != nil && runtime.stopPC == 0x100100650 {
		return 0x100100650, true, nil
	}
	// 0x100100650: add x1, x1, x22
	s.x[1] = uint64((s.x[1] + s.x[22]))
	if runtime != nil && runtime.stopPC == 0x100100654 {
		return 0x100100654, true, nil
	}
	// 0x100100654: add x19, x19, x21
	s.x[19] = uint64((s.x[19] + s.x[21]))
	if runtime != nil && runtime.stopPC == 0x100100658 {
		return 0x100100658, true, nil
	}
	// 0x100100658: eor w17, w17, w7
	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[7])))))
	if runtime != nil && runtime.stopPC == 0x10010065c {
		return 0x10010065c, true, nil
	}
	// 0x10010065c: add w17, w2, w17, uxtb
	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) + uint64(uint8(uint64(uint32(s.x[17])))))))
	if runtime != nil && runtime.stopPC == 0x100100660 {
		return 0x100100660, true, nil
	}
	// 0x100100660: ldr x17, [x28, w17, uxtw #3]
	value45, err := s.load((s.x[28] + shiftLeft(uint64(uint32(uint64(uint32(s.x[17])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(value45)
	if runtime != nil && runtime.stopPC == 0x100100664 {
		return 0x100100664, true, nil
	}
	// 0x100100664: ldur x28, [x29, #-0xb0]
	value46, err := s.load(addOffset(s.x[29], -176), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[28] = uint64(value46)
	if runtime != nil && runtime.stopPC == 0x100100668 {
		return 0x100100668, true, nil
	}
	// 0x100100668: ldrb w1, [x1]
	value47, err := s.load(s.x[1], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[1] = uint64(uint32(value47))
	if runtime != nil && runtime.stopPC == 0x10010066c {
		return 0x10010066c, true, nil
	}
	// 0x10010066c: ldrb w7, [x19]
	value48, err := s.load(s.x[19], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[7] = uint64(uint32(value48))
	if runtime != nil && runtime.stopPC == 0x100100670 {
		return 0x100100670, true, nil
	}
	// 0x100100670: add x19, x20, x3
	s.x[19] = uint64((s.x[20] + s.x[3]))
	if runtime != nil && runtime.stopPC == 0x100100674 {
		return 0x100100674, true, nil
	}
	// 0x100100674: add x17, x17, x19
	s.x[17] = uint64((s.x[17] + s.x[19]))
	if runtime != nil && runtime.stopPC == 0x100100678 {
		return 0x100100678, true, nil
	}
	// 0x100100678: ldrb w17, [x17]
	value51, err := s.load(s.x[17], 1)
	if err != nil {
		return 0, false, err
	}
	s.x[17] = uint64(uint32(value51))
	if runtime != nil && runtime.stopPC == 0x10010067c {
		return 0x10010067c, true, nil
	}
	// 0x10010067c: mul w19, w20, w4
	s.x[19] = uint64(uint32((uint64(uint32(s.x[20])) * uint64(uint32(s.x[4])))))
	if runtime != nil && runtime.stopPC == 0x100100680 {
		return 0x100100680, true, nil
	}
	// 0x100100680: eor w10, w10, w19
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[19])))))
	if runtime != nil && runtime.stopPC == 0x100100684 {
		return 0x100100684, true, nil
	}
	// 0x100100684: eor w10, w10, w1
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[1])))))
	if runtime != nil && runtime.stopPC == 0x100100688 {
		return 0x100100688, true, nil
	}
	// 0x100100688: eor w17, w7, w17
	s.x[17] = uint64(uint32((uint64(uint32(s.x[7])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x10010068c {
		return 0x10010068c, true, nil
	}
	// 0x10010068c: eor w10, w10, w17
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[17])))))
	if runtime != nil && runtime.stopPC == 0x100100690 {
		return 0x100100690, true, nil
	}
	// 0x100100690: eor w10, w10, w5
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[5])))))
	if runtime != nil && runtime.stopPC == 0x100100694 {
		return 0x100100694, true, nil
	}
	// 0x100100694: strb w10, [x28, x9]
	if err := s.storeValue((s.x[28] + s.x[9]), uint64(uint32(s.x[10])), 1); err != nil {
		return 0, false, err
	}
	if runtime != nil && runtime.stopPC == 0x100100698 {
		return 0x100100698, true, nil
	}
	// 0x100100698: cset w10, eq
	s.x[10] = uint64(uint32(boolValue(s.condition("eq"))))
	if runtime != nil && runtime.stopPC == 0x10010069c {
		return 0x10010069c, true, nil
	}
	// 0x10010069c: mul w10, w10, w6
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) * uint64(uint32(s.x[6])))))
	if runtime != nil && runtime.stopPC == 0x1001006a0 {
		return 0x1001006a0, true, nil
	}
	// 0x1001006a0: eor w10, w10, w11
	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))
	if runtime != nil && runtime.stopPC == 0x1001006a4 {
		return 0x1001006a4, true, nil
	}
	// 0x1001006a4: ldr x10, [x24, w10, sxtw #3]
	value62, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, false, err
	}
	s.x[10] = uint64(value62)
	if runtime != nil && runtime.stopPC == 0x1001006a8 {
		return 0x1001006a8, true, nil
	}
	// 0x1001006a8: mov x17, #0xfaf4
	s.x[17] = uint64(uint64(0xfaf4))
	if runtime != nil && runtime.stopPC == 0x1001006ac {
		return 0x1001006ac, true, nil
	}
	// 0x1001006ac: braa x10, x17
	next = s.x[10]
	return next, false, nil
}
