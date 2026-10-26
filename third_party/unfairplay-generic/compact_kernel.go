package recovery

import (
	"fmt"
	"math/bits"
)

type compactOpcode uint16

const (
	compactAnd32 compactOpcode = iota
	compactAdd32
	compactLoadSXTW32
	compactXor64
	compactAddLSR64
	compactUMADDL
	compactLSRV64
	compactAddUXTW64
	compactAndLSL64
	compactSub64
	compactAdd64
	compactXor32
	compactStoreSXTW32
	compactAddImm32
	compactCmp32
	compactCCmp32LE
	compactCSet32LT
	compactMul32
	compactLoadSXTW64
	compactMov64
	compactMovImm64
	compactStoreSP32
	compactMovImm32
	compactMovK32
	compactUDiv32
	compactMSub32
	compactLoadSP32
	compactUMull32
	compactLoadUXTW8
	compactLSRImm64
	compactMov32
	compactOr32
	compactSub32
	compactAndLSL32
	compactAddSPShiftImm64
	compactLoadBaseImm32
	compactExtract32
	compactNot32
	compactLSRImm32
	compactLoadSP64
	compactMovImm32High
	compactAddImmShift32
	compactLSLImm32
	compactXorLSL32
	compactStoreBase32
	compactLoadUXTW32
	compactMovNegOne32
	compactMovK64
	compactAddLSL32
	compactXorLSR32
	compactLSRReg32
	compactStoreUXTW8
	compactUMulH64
	compactMAdd32
	compactMAdd64
	compactAdds64
	compactADC32
	compactCMN64
	compactCmpImm32
	compactCSet32HI
	compactStoreBaseImm32
	compactMul64
	compactSubImmShift32
	compactXorImm32
	compactSubLSL32
	compactMovImm32Wide
	compactAnd64
	compactLSLImm64
	compactMovImm64Shift
	compactXorLSL64
	compactCSet32LO
	compactOrLSL32
	compactAndImm32
	compactOrNotLSL32
	compactRORImm32
	compactXorROR32
	compactOrROR32
	compactCSELZeroEq32
	compactCSet32EQ
	compactLoadBase64
	compactLoadPostIndex8
	compactLSLReg32
	compactStoreBase64
	compactCMNImm32
	compactOr64
	compactMovImm64Repeat16
	compactAddLSL64
	compactLoadBaseImm64
	compactAddImm64
	compactLoadBase32
	compactStoreSP64
	compactMovNegImm64
	compactCmp64
	compactSubImm64
	compactStoreXIndex8
	compactCmpImm64
	compactOrImm64
	compactXorImm64
	compactStoreUXTW32
	compactNoop
	compactMovImm48
	compactAddSPImm
	compactSubSPImm
	compactLoadSignedOffset32
	compactLoadSignedOffset64
	compactStoreSignedOffset32
	compactStoreSignedOffset64
	compactLoadSignedOffset8
	compactLoadSignedByte32
	compactStoreSignedOffset8
	compactBic32
	compactBic64
	compactSXT8To32
	compactSXT32To64
	compactNeg32
	compactNeg64
	compactTst32
	compactTstImm32
	compactCSet32
	compactCSEL32
	compactCSEL64
	compactCInc32
	compactCInc64
	compactCNeg32
	compactLoadHalf32
	compactLoadSignedOffsetHalf32
	compactStoreHalf32
	compactStoreSignedOffsetHalf32
	compactPairLoad32
	compactPairLoad64
	compactPairStore32
	compactPairStore64
	compactLoadBaseImm8
	compactLoadSPImm8
	compactLoadXIndex8
	compactStoreBaseImm8
	compactStoreSPImm8
	compactLoadXIndex32
	compactLoadXIndex64
	compactStoreXIndex32
	compactStoreXIndex64
	compactLoadUXTW64
	compactStoreBaseImm64
	compactAddExt32
	compactAddExt64
	compactMovSP
	compactSubSPFromX
	compactOrImm32
	compactOrLSL64
	compactAndLSR32
	compactCmpExtByte32
	compactCmpLSR32
	compactCCmpReg32
	compactCCmpReg64
	compactCCmpImm32
	compactCCmpImm64
	compactAdds32
	compactSubs32
	compactSubs64
	compactSubsImm32
	compactSubsImm64
	compactSMull32
	compactUBFIZ32
	compactUBFIZ64
	compactUBFX32
	compactUBFX64
	compactBFI32
	compactBFXIL32
	compactOrNotShift32
	compactOrNotShift64
	compactNot64
	compactLSLReg64
	compactLSRReg64
	compactAndImm64
	compactAndImmShift64
	compactXorNotImm16Shift64
	compactAddImmShift64
	compactSubImmShift64
	compactMovFromSP
	compactAddLSR32
	compactSubLSL64
	compactSubExt32
	compactLoadXIndexShift32
	compactStoreXIndexShift32
	compactLoadPreIndex32
	compactStorePreIndex64
	compactLoadSignedWord64
	compactAnds32
	compactBics32
	compactCCmnImm32
	compactCLZ32
	compactCMNImm64
	compactMovNegImm48
	compactMovNegImm16Shift64
	compactMovNegRepeat16PlusOne
	compactAndNotImm16Shift64
	compactCASAL32
	compactStorePostIndex32
	compactStorePreIndex32
	compactAddASR32
	compactAddUXTH64
	compactXorLSR64
	compactXorImm16Shift64
	compactLoadSignedWordBase64
	compactLoadHalfXIndex32
	compactLoadSPHalf32
	compactMNeg64
	compactMSub64
	compactSBFIZ64
	compactStoreUXTW64
	compactStorePostIndex8
	compactVectorDup32
	compactLoadVector
	compactASRImm32
	compactAndsLSR32
	compactBicsLSR32
	compactCCmnImm64
	compactCSInc32
	compactTstImm64
	compactEon32
	compactVectorImmediate
	compactVectorBinary
	compactVectorMAC
	compactVectorUnsignedShiftLong
	compactVectorMove64
	compactStoreVector
	compactVectorUnsignedMultiplyLong
	compactVectorEor3
	compactPairStoreVector
	compactVectorTable1
	compactVectorTable4
	compactVectorExtractNarrow
	compactVectorShiftRightNarrow
	compactVectorUnzipEven
	compactVectorNegateFloat64
	compactVectorUnsignedShiftRight
	compactVectorToW
	compactVectorUnsignedShift
	compactVectorMove128
	compactWToVectorLane
	compactAndsImm32
	compactASRReg32
	compactASRImm64
	compactBFI64
	compactBicLSR32
	compactCMNLSL64
	compactCmpUXTH32
	compactEon64
	compactXorASR32
	compactLoadPostIndex32
	compactLoadPostIndex64
	compactLoadHalfUXTW32
	compactMovMaxIntMinusOne
	compactOrLSR32
	compactOrHigh32Imm64
	compactRev32
	compactSBFIZ32
	compactSBFX32
	compactStoreSPHalf32
	compactStorePostIndex64
	compactMovMaxInt
	compactLoadPreIndex64
	compactXorSelectionMask64
	compactSXT16To32
	compactStoreHalfPostIndex32
	compactSubUXTH32
	compactUDiv64
	compactLoadHalfXIndexShift32
	compactSubLSR32
	compactLogicalMask64
	compactTestRegister64
	compactMNeg32
	compactRORImm64
	compactXorROR64
	compactLoadSXTW8
	compactSubLSR64
	compactPairLoadVector
	compactOpcodeCount
)

// compactInstruction stores only the operands needed by the checked recovery
// profile. It is intentionally not a general ARM64 instruction format.
type compactInstruction struct {
	op          compactOpcode
	destination uint8
	left, right uint8
	accumulator uint8
	shift       uint8
	immediate   uint16
}

type compactKernel struct {
	start        uint64
	instructions []compactInstruction
	segments     []compactSegment
	hooks        []compactHook
}

type compactSegment struct {
	index   uint32
	address uint64
}

type compactHookAction uint8

const (
	compactHookDigest compactHookAction = iota
	compactHookOptional
	compactHookEven
	compactHookOdd
	compactHookTail
	compactHookWriter
)

type compactHook struct {
	after  uint64
	action compactHookAction
}

func runCompactHook(s *state, action compactHookAction) error {
	profile, err := executionProfileForModel(s.model)
	if err != nil {
		return err
	}
	switch action {
	case compactHookDigest:
		_, err := runRawSourceKernelLoop(s, "digest-natural")
		return err
	case compactHookOptional:
		if s.x[9] == profile.optionalStart {
			return profile.runOptional(s, fault{})
		}
		return nil
	case compactHookEven:
		_, err := runRawSourceKernelLoop(s, "even")
		return err
	case compactHookOdd:
		_, err := runRawSourceKernelLoop(s, "odd")
		return err
	case compactHookTail:
		_, err := runRawSourceKernelLoop(s, "tail")
		return err
	case compactHookWriter:
		_, err := runRawSourceKernelLoop(s, "writer")
		return err
	default:
		return fmt.Errorf("unknown compact kernel hook %d", action)
	}
}

func runCompactKernel(s *state, kernel compactKernel) error {
	pc := kernel.start
	segment := 0
	hook := 0
	for index, instruction := range kernel.instructions {
		if segment < len(kernel.segments) && kernel.segments[segment].index == uint32(index) {
			pc = kernel.segments[segment].address
			segment++
		}
		d, left, right := int(instruction.destination), int(instruction.left), int(instruction.right)
		switch instruction.op {
		case compactAnd32:
			s.x[d] = uint64(uint32(s.x[left]) & uint32(s.x[right]))
		case compactAdd32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(s.x[right]))
		case compactLoadSXTW32:
			value, err := s.load(s.x[left]+shiftLeft(signExtend32(uint32(s.x[right])), int(instruction.shift), 64), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactXor64:
			s.x[d] = s.x[left] ^ s.x[right]
		case compactAddLSR64:
			s.x[d] = s.x[left] + shiftRight(s.x[right], int(instruction.shift), 64)
		case compactUMADDL:
			s.x[d] = uint64(uint32(s.x[left]))*uint64(uint32(s.x[right])) + s.x[instruction.accumulator]
		case compactLSRV64:
			s.x[d] = shiftRight(s.x[left], int(s.x[right]), 64)
		case compactAddUXTW64:
			s.x[d] = s.x[left] + uint64(uint32(s.x[right]))
		case compactAndLSL64:
			s.x[d] = s.x[left] & shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactLSLImm64:
			s.x[d] = shiftLeft(s.x[left], int(instruction.shift), 64)
		case compactMovImm64Shift:
			s.x[d] = uint64(instruction.immediate) << instruction.shift
		case compactXorLSL64:
			s.x[d] = s.x[left] ^ shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactCSet32LO:
			s.x[d] = uint64(uint32(boolValue(s.condition("lo"))))
		case compactOrLSL32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactAndImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			s.x[d] = uint64(uint32(s.x[left]) & immediate)
		case compactOrNotLSL32:
			s.x[d] = uint64(uint32(s.x[left]) | ^uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactRORImm32:
			s.x[d] = uint64(uint32(rotateRight(uint64(uint32(s.x[left])), int(instruction.shift), 32)))
		case compactRORImm64:
			s.x[d] = rotateRight(s.x[left], int(instruction.shift), 64)
		case compactXorROR64:
			s.x[d] = s.x[left] ^ rotateRight(s.x[right], int(instruction.shift), 64)
		case compactLoadSXTW8:
			address := s.x[left] + uint64(int64(int32(s.x[right])))
			value, err := s.load(address, 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(byte(value))
		case compactSubLSR64:
			s.x[d] = s.x[left] - (s.x[right] >> instruction.shift)
		case compactPairLoadVector:
			base := s.x[left]
			if instruction.left == 31 {
				base = s.sp
			}
			address := addOffset(base, int64(int16(instruction.immediate)))
			if err := s.loadVector(d, address, 16); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if err := s.loadVector(int(instruction.right), address+16, 16); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactXorROR32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(rotateRight(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactOrROR32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(rotateRight(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactCSELZeroEq32:
			if s.condition("eq") {
				s.x[d] = 0
			} else {
				s.x[d] = uint64(uint32(s.x[right]))
			}
		case compactCSet32EQ:
			s.x[d] = uint64(uint32(boolValue(s.condition("eq"))))
		case compactLoadBase64:
			value, err := s.load(s.x[left], 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactLoadPostIndex8:
			value, err := s.load(s.x[left], 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
			s.x[left] += uint64(instruction.immediate)
		case compactLSLReg32:
			s.x[d] = uint64(uint32(shiftLeft(uint64(uint32(s.x[left])), int(s.x[right]), 32)))
		case compactStoreBase64:
			value := uint64(0)
			if instruction.destination < 31 {
				value = s.x[d]
			}
			if err := s.storeValue(s.x[left], value, 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactCMNImm32:
			_ = s.addWithCarry(uint64(uint32(s.x[left])), uint64(uint32(instruction.immediate)), false, 32)
		case compactOr64:
			s.x[d] = s.x[left] | s.x[right]
		case compactMovImm64Repeat16:
			value := uint64(instruction.immediate)
			s.x[d] = value | value<<16 | value<<32 | value<<48
		case compactAddLSL64:
			s.x[d] = s.x[left] + shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactLoadBaseImm64:
			value, err := s.load(addOffset(s.x[left], int64(instruction.immediate)), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactAddImm64:
			s.x[d] = s.x[left] + uint64(instruction.immediate)
		case compactLoadBase32:
			value, err := s.load(s.x[left], 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactStoreSP64:
			value := uint64(0)
			if instruction.left < 31 {
				value = s.x[left]
			}
			if err := s.storeValue(addOffset(s.sp, int64(instruction.immediate)), value, 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactMovNegImm64:
			s.x[d] = uint64(0) - uint64(instruction.immediate)
		case compactCmp64:
			s.compare64(s.x[left], s.x[right])
		case compactSubImm64:
			s.x[d] = s.x[left] - uint64(instruction.immediate)
		case compactStoreXIndex8:
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(byte(s.x[d]))
			}
			if err := s.storeValue(s.x[left]+s.x[right], value, 1); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactCmpImm64:
			s.compare64(s.x[left], uint64(instruction.immediate))
		case compactOrImm64:
			s.x[d] = s.x[left] | uint64(instruction.immediate)
		case compactXorImm64:
			s.x[d] = s.x[left] ^ uint64(instruction.immediate)
		case compactStoreUXTW32:
			address := s.x[left] + shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 64)
			if err := s.storeValue(address, uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactNoop:
		case compactMovImm48:
			s.x[d] = uint64(instruction.immediate) |
				uint64(instruction.left)<<16 |
				uint64(instruction.right)<<24 |
				uint64(instruction.accumulator)<<32 |
				uint64(instruction.shift)<<40
		case compactAddSPImm:
			s.sp += shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactSubSPImm:
			s.sp -= shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactLoadSignedOffset32:
			value, err := s.load(addOffset(s.x[left], int64(int16(instruction.immediate))), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactLoadSignedOffset64:
			value, err := s.load(addOffset(s.x[left], int64(int16(instruction.immediate))), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactStoreSignedOffset32:
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(uint32(s.x[d]))
			}
			if err := s.storeValue(addOffset(s.x[left], int64(int16(instruction.immediate))), value, 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactStoreSignedOffset64:
			value := uint64(0)
			if instruction.destination < 31 {
				value = s.x[d]
			}
			if err := s.storeValue(addOffset(s.x[left], int64(int16(instruction.immediate))), value, 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactLoadSignedOffset8:
			value, err := s.load(addOffset(s.x[left], int64(int16(instruction.immediate))), 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(byte(value))
		case compactLoadSignedByte32:
			value, err := s.load(addOffset(s.x[left], int64(int16(instruction.immediate))), 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(int32(int8(value))))
		case compactStoreSignedOffset8:
			if err := s.storeValue(addOffset(s.x[left], int64(int16(instruction.immediate))), uint64(byte(s.x[d])), 1); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactBic32:
			s.x[d] = uint64(uint32(s.x[left]) &^ uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactBic64:
			s.x[d] = s.x[left] &^ shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactSXT8To32:
			s.x[d] = uint64(uint32(int32(int8(s.x[left]))))
		case compactSXT32To64:
			s.x[d] = uint64(int64(int32(s.x[left])))
		case compactNeg32:
			s.x[d] = uint64(uint32(0) - uint32(s.x[left]))
		case compactNeg64:
			s.x[d] = uint64(0) - s.x[left]
		case compactTst32:
			s.test32(uint32(s.x[left]) & uint32(s.x[right]))
		case compactTstImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			s.test32(uint32(s.x[left]) & immediate)
		case compactCSet32:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(boolValue(matches)))
		case compactCSEL32, compactCSEL64:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			value := uint64(0)
			if instruction.right < 31 {
				value = s.x[right]
			}
			if matches {
				value = 0
				if instruction.left < 31 {
					value = s.x[left]
				}
			}
			if instruction.op == compactCSEL32 {
				value = uint64(uint32(value))
			}
			s.x[d] = value
		case compactCInc32, compactCInc64:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			value := s.x[left]
			if matches {
				value++
			}
			if instruction.op == compactCInc32 {
				value = uint64(uint32(value))
			}
			s.x[d] = value
		case compactCNeg32:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			value := uint32(s.x[left])
			if matches {
				value = 0 - value
			}
			s.x[d] = uint64(value)
		case compactLoadHalf32, compactLoadSignedOffsetHalf32:
			offset := int64(0)
			if instruction.op == compactLoadSignedOffsetHalf32 {
				offset = int64(int16(instruction.immediate))
			}
			value, err := s.load(addOffset(s.x[left], offset), 2)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactStoreHalf32, compactStoreSignedOffsetHalf32:
			offset := int64(0)
			if instruction.op == compactStoreSignedOffsetHalf32 {
				offset = int64(int16(instruction.immediate))
			}
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(uint16(s.x[d]))
			}
			if err := s.storeValue(addOffset(s.x[left], offset), value, 2); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactPairLoad32, compactPairLoad64, compactPairStore32, compactPairStore64:
			base := uint64(0)
			if instruction.left == 31 {
				base = s.sp
			} else {
				base = s.x[left]
			}
			offset := int64(int16(instruction.immediate))
			address := base
			if instruction.accumulator == 0 {
				address = addOffset(base, offset)
			}
			if instruction.accumulator == 1 {
				base = addOffset(base, offset)
				address = base
				if instruction.left == 31 {
					s.sp = base
				} else {
					s.x[left] = base
				}
			}
			width := 8
			if instruction.op == compactPairLoad32 || instruction.op == compactPairStore32 {
				width = 4
			}
			if instruction.op == compactPairLoad32 || instruction.op == compactPairLoad64 {
				first, err := s.load(address, width)
				if err != nil {
					return fmt.Errorf("compact kernel %#x: %w", pc, err)
				}
				second, err := s.load(address+uint64(width), width)
				if err != nil {
					return fmt.Errorf("compact kernel %#x: %w", pc, err)
				}
				if instruction.destination < 31 {
					s.x[d] = first
				}
				if instruction.right < 31 {
					s.x[right] = second
				}
			} else {
				first, second := uint64(0), uint64(0)
				if instruction.destination < 31 {
					first = s.x[d]
				}
				if instruction.right < 31 {
					second = s.x[right]
				}
				if width == 4 {
					first, second = uint64(uint32(first)), uint64(uint32(second))
				}
				if err := s.storeValue(address, first, width); err != nil {
					return fmt.Errorf("compact kernel %#x: %w", pc, err)
				}
				if err := s.storeValue(address+uint64(width), second, width); err != nil {
					return fmt.Errorf("compact kernel %#x: %w", pc, err)
				}
			}
			if instruction.accumulator == 2 {
				base = addOffset(base, offset)
				if instruction.left == 31 {
					s.sp = base
				} else {
					s.x[left] = base
				}
			}
		case compactLoadBaseImm8, compactLoadSPImm8:
			base := s.x[left]
			if instruction.op == compactLoadSPImm8 {
				base = s.sp
			}
			value, err := s.load(addOffset(base, int64(instruction.immediate)), 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(byte(value))
		case compactLoadXIndex8, compactLoadXIndex32, compactLoadXIndex64:
			width := 1
			if instruction.op == compactLoadXIndex32 {
				width = 4
			} else if instruction.op == compactLoadXIndex64 {
				width = 8
			}
			value, err := s.load(s.x[left]+s.x[right], width)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if width == 4 {
				value = uint64(uint32(value))
			}
			s.x[d] = value
		case compactStoreBaseImm8, compactStoreSPImm8:
			base := s.x[left]
			if instruction.op == compactStoreSPImm8 {
				base = s.sp
			}
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(byte(s.x[d]))
			}
			if err := s.storeValue(addOffset(base, int64(instruction.immediate)), value, 1); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactStoreXIndex32, compactStoreXIndex64:
			width := 8
			if instruction.op == compactStoreXIndex32 {
				width = 4
			}
			value := uint64(0)
			if instruction.destination < 31 {
				value = s.x[d]
			}
			if width == 4 {
				value = uint64(uint32(value))
			}
			if err := s.storeValue(s.x[left]+s.x[right], value, width); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactLoadUXTW64:
			value, err := s.load(s.x[left]+shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 64), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactStoreBaseImm64:
			value := uint64(0)
			if instruction.destination < 31 {
				value = s.x[d]
			}
			if err := s.storeValue(addOffset(s.x[left], int64(instruction.immediate)), value, 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactAddExt32, compactAddExt64:
			extended := uint64(byte(s.x[right]))
			if instruction.accumulator == 1 {
				extended = uint64(uint32(s.x[right]))
			} else if instruction.accumulator == 2 {
				extended = uint64(int64(int32(s.x[right])))
			} else if instruction.accumulator == 3 {
				extended = uint64(uint16(s.x[right]))
			}
			value := s.x[left] + shiftLeft(extended, int(instruction.shift), 64)
			if instruction.op == compactAddExt32 {
				value = uint64(uint32(value))
			}
			s.x[d] = value
		case compactMovSP:
			s.sp = s.x[left]
		case compactSubSPFromX:
			s.sp = s.x[left] - shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactOrImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			s.x[d] = uint64(uint32(s.x[left]) | immediate)
		case compactOrLSL64:
			s.x[d] = s.x[left] | shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactAndLSR32:
			s.x[d] = uint64(uint32(s.x[left]) & uint32(shiftRight(s.x[right], int(instruction.shift), 32)))
		case compactCmpExtByte32:
			rightValue := uint32(byte(s.x[right]))
			if instruction.shift != 0 {
				rightValue = uint32(int32(int8(s.x[right])))
			}
			s.compare32(uint32(s.x[left]), rightValue)
		case compactCmpLSR32:
			s.compare32(uint32(s.x[left]), uint32(shiftRight(s.x[right], int(instruction.shift), 32)))
		case compactCCmpReg32, compactCCmpReg64, compactCCmpImm32, compactCCmpImm64:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if !matches {
				s.setNZCV(int(instruction.accumulator))
				break
			}
			rightValue := s.x[right]
			if instruction.op == compactCCmpImm32 || instruction.op == compactCCmpImm64 {
				rightValue = uint64(instruction.immediate)
			}
			if instruction.op == compactCCmpReg32 || instruction.op == compactCCmpImm32 {
				s.compare32(uint32(s.x[left]), uint32(rightValue))
			} else {
				s.compare64(s.x[left], rightValue)
			}
		case compactAdds32:
			s.x[d] = s.addWithCarry(uint64(uint32(s.x[left])), uint64(uint32(s.x[right])), false, 32)
		case compactSubs32:
			s.x[d] = s.subtractWithFlags(uint64(uint32(s.x[left])), uint64(uint32(s.x[right])), 32)
		case compactSubs64:
			s.x[d] = s.subtractWithFlags(s.x[left], s.x[right], 64)
		case compactSubsImm32:
			s.x[d] = s.subtractWithFlags(uint64(uint32(s.x[left])), uint64(instruction.immediate), 32)
		case compactSubsImm64:
			s.x[d] = s.subtractWithFlags(s.x[left], uint64(instruction.immediate), 64)
		case compactSMull32:
			s.x[d] = uint64(int64(int32(s.x[left])) * int64(int32(s.x[right])))
		case compactUBFIZ32, compactUBFIZ64:
			width := 64
			if instruction.op == compactUBFIZ32 {
				width = 32
			}
			s.x[d] = shiftLeft(s.x[left]&bitMask(int(instruction.immediate), width), int(instruction.shift), width)
		case compactUBFX32, compactUBFX64:
			width := 64
			if instruction.op == compactUBFX32 {
				width = 32
			}
			s.x[d] = shiftRight(s.x[left], int(instruction.shift), width) & bitMask(int(instruction.immediate), width)
		case compactBFI32:
			s.x[d] = insertBits(uint64(uint32(s.x[d])), uint64(uint32(s.x[left])), int(instruction.shift), int(instruction.immediate), 32)
		case compactBFXIL32:
			s.x[d] = insertBits(uint64(uint32(s.x[d])), shiftRight(s.x[left], int(instruction.shift), 32), 0, int(instruction.immediate), 32)
		case compactOrNotShift32, compactOrNotShift64:
			width := 64
			if instruction.op == compactOrNotShift32 {
				width = 32
			}
			shifted := shiftLeft(s.x[right], int(instruction.shift), width)
			if instruction.accumulator == 1 {
				shifted = shiftRight(s.x[right], int(instruction.shift), width)
			}
			value := s.x[left] | ^shifted
			if width == 32 {
				value = uint64(uint32(value))
			}
			s.x[d] = value
		case compactNot64:
			s.x[d] = ^s.x[left]
		case compactLSLReg64:
			s.x[d] = shiftLeft(s.x[left], int(s.x[right]), 64)
		case compactLSRReg64:
			s.x[d] = shiftRight(s.x[left], int(s.x[right]), 64)
		case compactAndImm64:
			s.x[d] = s.x[left] & uint64(instruction.immediate)
		case compactAndImmShift64:
			s.x[d] = s.x[left] & (uint64(instruction.immediate) << instruction.shift)
		case compactXorNotImm16Shift64:
			s.x[d] = s.x[left] ^ ^(uint64(instruction.immediate) << instruction.shift)
		case compactAddImmShift64:
			s.x[d] = s.x[left] + shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactSubImmShift64:
			s.x[d] = s.x[left] - shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactMovFromSP:
			s.x[d] = s.sp
		case compactAddLSR32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(shiftRight(s.x[right], int(instruction.shift), 32)))
		case compactSubLSL64:
			s.x[d] = s.x[left] - shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactSubExt32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(byte(s.x[right])))
		case compactLoadXIndexShift32:
			value, err := s.load(s.x[left]+shiftLeft(s.x[right], int(instruction.shift), 64), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactStoreXIndexShift32:
			if err := s.storeValue(s.x[left]+shiftLeft(s.x[right], int(instruction.shift), 64), uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactLoadPreIndex32:
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
			value, err := s.load(s.x[left], 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactStorePreIndex64:
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
			value := uint64(0)
			if instruction.destination < 31 {
				value = s.x[d]
			}
			if err := s.storeValue(s.x[left], value, 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactLoadSignedWord64:
			value, err := s.load(addOffset(s.x[left], int64(int16(instruction.immediate))), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(int64(int32(value)))
		case compactAnds32:
			value := uint32(s.x[left]) & uint32(s.x[right])
			s.x[d] = uint64(value)
			s.test32(value)
		case compactBics32:
			value := uint32(s.x[left]) &^ uint32(s.x[right])
			if instruction.destination < 31 {
				s.x[d] = uint64(value)
			}
			s.test32(value)
		case compactCCmnImm32:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if matches {
				_ = s.addWithCarry(uint64(uint32(s.x[left])), uint64(instruction.immediate), false, 32)
			} else {
				s.setNZCV(int(instruction.accumulator))
			}
		case compactCLZ32:
			s.x[d] = uint64(uint32(bits.LeadingZeros32(uint32(s.x[left]))))
		case compactCMNImm64:
			_ = s.addWithCarry(s.x[left], uint64(instruction.immediate), false, 64)
		case compactMovNegImm48:
			magnitude := uint64(instruction.immediate) |
				uint64(instruction.left)<<16 |
				uint64(instruction.right)<<24 |
				uint64(instruction.accumulator)<<32 |
				uint64(instruction.shift)<<40
			s.x[d] = uint64(0) - magnitude
		case compactMovNegImm16Shift64:
			s.x[d] = uint64(0) - (uint64(instruction.immediate) << instruction.shift)
		case compactMovNegRepeat16PlusOne:
			value := uint64(instruction.immediate)
			magnitude := value | value<<16 | value<<32 | value<<48
			s.x[d] = uint64(0) - (magnitude + 1)
		case compactAndNotImm16Shift64:
			s.x[d] = s.x[left] & ^(uint64(instruction.immediate) << instruction.shift)
		case compactCASAL32:
			value, err := s.load(s.x[instruction.accumulator], 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if uint32(value) == uint32(s.x[left]) {
				if err := s.storeValue(s.x[instruction.accumulator], uint64(uint32(s.x[right])), 4); err != nil {
					return fmt.Errorf("compact kernel %#x: %w", pc, err)
				}
			}
			s.x[left] = uint64(uint32(value))
		case compactStorePostIndex32:
			if err := s.storeValue(s.x[left], uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
		case compactStorePostIndex64:
			if err := s.storeValue(s.x[left], s.x[d], 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
		case compactMovMaxInt:
			s.x[d] = uint64(0x7fffffffffffffff)
		case compactLoadPreIndex64:
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
			value, err := s.load(s.x[left], 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactXorSelectionMask64:
			s.x[d] = s.x[left] ^ uint64(0x3ffffffffffffff0)
		case compactSXT16To32:
			s.x[d] = uint64(uint32(int32(int16(s.x[left]))))
		case compactStoreHalfPostIndex32:
			if err := s.storeValue(s.x[left], uint64(uint16(s.x[d])), 2); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
		case compactSubUXTH32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(uint16(s.x[right])))
		case compactUDiv64:
			s.x[d] = unsignedDivide(s.x[left], s.x[right], 64)
		case compactLoadHalfXIndexShift32:
			value, err := s.load(s.x[left]+shiftLeft(s.x[right], int(instruction.shift), 64), 2)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint16(value))
		case compactStorePreIndex32:
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
			if err := s.storeValue(s.x[left], uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactAddASR32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(shiftRightArithmetic(s.x[right], int(instruction.shift), 32)))
		case compactAddUXTH64:
			s.x[d] = s.x[left] + uint64(uint16(s.x[right]))
		case compactXorLSR64:
			s.x[d] = s.x[left] ^ shiftRight(s.x[right], int(instruction.shift), 64)
		case compactXorImm16Shift64:
			s.x[d] = s.x[left] ^ uint64(instruction.immediate)<<instruction.shift
		case compactLoadSignedWordBase64:
			base := uint64(0)
			if instruction.left == 31 {
				base = s.sp
			} else {
				base = s.x[left]
			}
			value, err := s.load(addOffset(base, int64(int16(instruction.immediate))), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(int64(int32(value)))
		case compactLoadHalfXIndex32:
			value, err := s.load(s.x[left]+s.x[right], 2)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint16(value))
		case compactLoadSPHalf32:
			value, err := s.load(addOffset(s.sp, int64(instruction.immediate)), 2)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint16(value))
		case compactMNeg32:
			s.x[d] = uint64(uint32(0) - uint32(s.x[left])*uint32(s.x[right]))
		case compactMNeg64:
			s.x[d] = uint64(0) - s.x[left]*s.x[right]
		case compactMSub64:
			s.x[d] = s.x[instruction.accumulator] - s.x[left]*s.x[right]
		case compactSBFIZ64:
			s.x[d] = shiftLeft(signExtendBits(s.x[left], int(instruction.immediate), 64), int(instruction.shift), 64)
		case compactStoreUXTW64:
			address := s.x[left] + shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 64)
			value := uint64(0)
			if instruction.destination < 31 {
				value = s.x[d]
			}
			if err := s.storeValue(address, value, 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactStorePostIndex8:
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(byte(s.x[d]))
			}
			if err := s.storeValue(s.x[left], value, 1); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
		case compactVectorDup32:
			s.vectorDup(d, int(instruction.immediate), 32, uint64(uint32(s.x[left])))
		case compactLoadVector:
			base := s.x[left]
			if instruction.left == 31 {
				base = s.sp
			}
			if err := s.loadVector(d, addOffset(base, int64(instruction.immediate)), int(instruction.accumulator)); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactASRImm32:
			s.x[d] = uint64(uint32(shiftRightArithmetic(s.x[left], int(instruction.shift), 32)))
		case compactAndsLSR32:
			value := uint32(s.x[left]) & uint32(shiftRight(s.x[right], int(instruction.shift), 32))
			s.x[d] = uint64(value)
			s.test32(value)
		case compactBicsLSR32:
			value := uint32(s.x[left]) &^ uint32(shiftRight(s.x[right], int(instruction.shift), 32))
			if instruction.destination < 31 {
				s.x[d] = uint64(value)
			}
			s.test32(value)
		case compactCCmnImm64:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if matches {
				_ = s.addWithCarry(s.x[left], uint64(instruction.immediate), false, 64)
			} else {
				s.setNZCV(int(instruction.accumulator))
			}
		case compactCSInc32:
			matches, err := compactConditionMatches(s, compactCondition(instruction.shift))
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if matches {
				s.x[d] = uint64(uint32(s.x[left]))
			} else {
				value := uint32(0)
				if instruction.right < 31 {
					value = uint32(s.x[right])
				}
				s.x[d] = uint64(value + 1)
			}
		case compactTstImm64:
			s.test64(s.x[left] & uint64(instruction.immediate))
		case compactEon32:
			s.x[d] = uint64(uint32(s.x[left]) ^ ^uint32(s.x[right]))
		case compactVectorImmediate:
			s.clearVector(d)
			value := uint64(instruction.immediate) << instruction.shift
			for lane := 0; lane < int(instruction.right); lane++ {
				s.setVectorLane(d, lane, int(instruction.left), value)
			}
		case compactVectorBinary:
			operation := [...]string{"add", "sub", "mul", "eor"}[instruction.accumulator]
			s.vectorBinary(operation, d, left, right, int(instruction.immediate), int(instruction.shift))
		case compactVectorMAC:
			s.vectorMultiplyAccumulate(d, left, right, int(instruction.immediate), int(instruction.shift), instruction.accumulator != 0)
		case compactVectorUnsignedShiftLong:
			lanes := int(byte(instruction.immediate))
			inputBits := int(byte(instruction.immediate >> 8))
			s.vectorUnsignedShiftLong(d, left, lanes, inputBits, int(instruction.shift), instruction.accumulator != 0)
		case compactVectorMove64:
			value := s.vectorLane(left, 0, 64)
			s.clearVector(d)
			s.setVectorLane(d, 0, 64, value)
		case compactStoreVector:
			base := uint64(0)
			if instruction.left == 31 {
				base = s.sp
			} else {
				base = s.x[left]
			}
			if err := s.storeVector(d, addOffset(base, int64(int16(instruction.immediate))), int(instruction.accumulator)); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactVectorUnsignedMultiplyLong:
			s.vectorUnsignedMultiplyLong(d, left, right, int(instruction.immediate), int(instruction.shift), instruction.accumulator != 0)
		case compactVectorEor3:
			s.vectorEor3(d, left, right, int(instruction.accumulator))
		case compactPairStoreVector:
			base := s.x[left]
			address := addOffset(base, int64(int16(instruction.immediate)))
			if err := s.storeVector(d, address, 16); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if err := s.storeVector(int(instruction.right), address+16, 16); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactVectorTable1:
			s.vectorTable(d, []int{left}, right)
		case compactVectorTable4:
			s.vectorTable(d, []int{left, left + 1, left + 2, left + 3}, right)
		case compactVectorExtractNarrow:
			s.vectorExtractNarrow(d, left, int(instruction.immediate), int(instruction.shift))
		case compactVectorShiftRightNarrow:
			s.vectorShiftRightNarrow(d, left, int(instruction.immediate), int(instruction.accumulator), int(instruction.shift))
		case compactVectorUnzipEven:
			s.vectorUnzipEven(d, left, right, int(instruction.immediate), int(instruction.shift))
		case compactVectorNegateFloat64:
			value := s.vectorLane(left, 0, 64) ^ (uint64(1) << 63)
			s.clearVector(d)
			s.setVectorLane(d, 0, 64, value)
		case compactVectorUnsignedShiftRight:
			s.vectorUnsignedShiftRight(d, left, int(instruction.immediate), int(instruction.accumulator), int(instruction.shift))
		case compactVectorToW:
			s.x[d] = uint64(uint32(s.vectorLane(left, int(instruction.immediate), 32)))
		case compactVectorUnsignedShift:
			s.vectorUnsignedShift(d, left, right, int(instruction.immediate), int(instruction.shift))
		case compactVectorMove128:
			s.vectors[d] = s.vectors[left]
		case compactWToVectorLane:
			s.setVectorLane(d, int(instruction.immediate), 32, uint64(uint32(s.x[left])))
		case compactAndsImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			value := uint32(s.x[left]) & immediate
			s.x[d] = uint64(value)
			s.test32(value)
		case compactASRReg32:
			s.x[d] = uint64(uint32(shiftRightArithmetic(s.x[left], int(s.x[right]), 32)))
		case compactASRImm64:
			s.x[d] = shiftRightArithmetic(s.x[left], int(instruction.shift), 64)
		case compactBFI64:
			s.x[d] = insertBits(s.x[d], s.x[left], int(instruction.shift), int(instruction.immediate), 64)
		case compactBicLSR32:
			s.x[d] = uint64(uint32(s.x[left]) &^ uint32(shiftRight(s.x[right], int(instruction.shift), 32)))
		case compactCMNLSL64:
			_ = s.addWithCarry(s.x[left], shiftLeft(s.x[right], int(instruction.shift), 64), false, 64)
		case compactCmpUXTH32:
			s.compare32(uint32(s.x[left]), uint32(uint16(s.x[right])))
		case compactEon64:
			s.x[d] = s.x[left] ^ ^s.x[right]
		case compactXorASR32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(shiftRightArithmetic(s.x[right], int(instruction.shift), 32)))
		case compactLoadPostIndex32, compactLoadPostIndex64:
			width := 8
			if instruction.op == compactLoadPostIndex32 {
				width = 4
			}
			value, err := s.load(s.x[left], width)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			if width == 4 {
				value = uint64(uint32(value))
			}
			s.x[d] = value
			s.x[left] = addOffset(s.x[left], int64(int16(instruction.immediate)))
		case compactLoadHalfUXTW32:
			value, err := s.load(s.x[left]+shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 64), 2)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint16(value))
		case compactMovMaxIntMinusOne:
			s.x[d] = uint64(0x7ffffffffffffffe)
		case compactOrLSR32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(shiftRight(s.x[right], int(instruction.shift), 32)))
		case compactOrHigh32Imm64:
			s.x[d] = s.x[left] | uint64(0xffffffff00000000)
		case compactRev32:
			s.x[d] = uint64(bits.ReverseBytes32(uint32(s.x[left])))
		case compactSBFIZ32:
			s.x[d] = shiftLeft(signExtendBits(s.x[left], int(instruction.immediate), 32), int(instruction.shift), 32)
		case compactSBFX32:
			s.x[d] = signExtendBits(shiftRight(s.x[left], int(instruction.shift), 32), int(instruction.immediate), 32)
		case compactStoreSPHalf32:
			if err := s.storeValue(addOffset(s.sp, int64(instruction.immediate)), uint64(uint16(s.x[d])), 2); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactSub64:
			s.x[d] = s.x[left] - s.x[right]
		case compactAdd64:
			s.x[d] = s.x[left] + s.x[right]
		case compactXor32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(s.x[right]))
		case compactStoreSXTW32:
			address := s.x[left] + shiftLeft(signExtend32(uint32(s.x[right])), int(instruction.shift), 64)
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(uint32(s.x[d]))
			}
			if err := s.storeValue(address, value, 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactAddImm32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(instruction.immediate))
		case compactCmp32:
			s.compare32(uint32(s.x[left]), uint32(s.x[right]))
		case compactCCmp32LE:
			if s.condition("le") {
				s.compare32(uint32(s.x[left]), uint32(s.x[right]))
			} else {
				s.setNZCV(int(instruction.immediate))
			}
		case compactCSet32LT:
			s.x[d] = uint64(uint32(boolValue(s.condition("lt"))))
		case compactMul32:
			s.x[d] = uint64(uint32(s.x[left]) * uint32(s.x[right]))
		case compactLoadSXTW64:
			value, err := s.load(s.x[left]+shiftLeft(signExtend32(uint32(s.x[right])), int(instruction.shift), 64), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactMov64:
			s.x[d] = s.x[left]
		case compactMovImm64:
			s.x[d] = uint64(instruction.immediate)
		case compactStoreSP32:
			value := uint64(0)
			if instruction.left < 31 {
				value = uint64(uint32(s.x[left]))
			}
			if err := s.storeValue(addOffset(s.sp, int64(instruction.immediate)), value, 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactMovImm32:
			s.x[d] = uint64(uint32(instruction.immediate))
		case compactMovK32:
			s.x[d] = uint64(uint32(insertBits(uint64(uint32(s.x[d])), uint64(instruction.immediate), int(instruction.shift), 16, 32)))
		case compactUDiv32:
			s.x[d] = uint64(uint32(unsignedDivide(uint64(uint32(s.x[left])), uint64(uint32(s.x[right])), 32)))
		case compactMSub32:
			s.x[d] = uint64(uint32(s.x[instruction.accumulator]) - uint32(s.x[left])*uint32(s.x[right]))
		case compactLoadSP32:
			value, err := s.load(addOffset(s.sp, int64(instruction.immediate)), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactUMull32:
			s.x[d] = uint64(uint32(s.x[left])) * uint64(uint32(s.x[right]))
		case compactLoadUXTW8:
			value, err := s.load(s.x[left]+uint64(uint32(s.x[right])), 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactLSRImm64:
			s.x[d] = shiftRight(s.x[left], int(instruction.shift), 64)
		case compactMov32:
			s.x[d] = uint64(uint32(s.x[left]))
		case compactOr32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(s.x[right]))
		case compactSub32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(s.x[right]))
		case compactAndLSL32:
			s.x[d] = uint64(uint32(s.x[left]) & uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactAddSPShiftImm64:
			s.x[d] = s.sp + shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactLoadBaseImm32:
			value, err := s.load(addOffset(s.x[left], int64(instruction.immediate)), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactExtract32:
			s.x[d] = uint64(uint32(extractPair(uint64(uint32(s.x[left])), uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactNot32:
			s.x[d] = uint64(^uint32(s.x[left]))
		case compactLSRImm32:
			s.x[d] = uint64(uint32(shiftRight(uint64(uint32(s.x[left])), int(instruction.shift), 32)))
		case compactLoadSP64:
			value, err := s.load(addOffset(s.sp, int64(instruction.immediate)), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactMovImm32High:
			s.x[d] = uint64(uint32(instruction.immediate) << 16)
		case compactAddImmShift32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(instruction.immediate)<<instruction.shift)
		case compactLSLImm32:
			s.x[d] = uint64(uint32(shiftLeft(uint64(uint32(s.x[left])), int(instruction.shift), 32)))
		case compactXorLSL32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactStoreBase32:
			if err := s.storeValue(s.x[left], uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactLoadUXTW32:
			value, err := s.load(s.x[left]+shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 64), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactMovNegOne32:
			s.x[d] = uint64(^uint32(0))
		case compactMovK64:
			s.x[d] = insertBits(s.x[d], uint64(instruction.immediate), int(instruction.shift), 16, 64)
		case compactAddLSL32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactXorLSR32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(shiftRight(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactLSRReg32:
			s.x[d] = uint64(uint32(shiftRight(uint64(uint32(s.x[left])), int(s.x[right]), 32)))
		case compactStoreUXTW8:
			address := s.x[left] + uint64(uint32(s.x[right]))
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(byte(s.x[d]))
			}
			if err := s.storeValue(address, value, 1); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactUMulH64:
			s.x[d] = multiplyHigh64(s.x[left], s.x[right])
		case compactMAdd32:
			s.x[d] = uint64(uint32(s.x[left])*uint32(s.x[right]) + uint32(s.x[instruction.accumulator]))
		case compactMAdd64:
			s.x[d] = s.x[left]*s.x[right] + s.x[instruction.accumulator]
		case compactAdds64:
			s.x[d] = s.addWithCarry(s.x[left], s.x[right], false, 64)
		case compactADC32:
			leftValue, rightValue := uint32(0), uint32(0)
			if instruction.left < 31 {
				leftValue = uint32(s.x[left])
			}
			if instruction.right < 31 {
				rightValue = uint32(s.x[right])
			}
			s.x[d] = uint64(leftValue + rightValue + uint32(boolValue(s.c)))
		case compactCMN64:
			_ = s.addWithCarry(s.x[left], s.x[right], false, 64)
		case compactCmpImm32:
			s.compare32(uint32(s.x[left]), uint32(instruction.immediate))
		case compactCSet32HI:
			s.x[d] = uint64(uint32(boolValue(s.condition("hi"))))
		case compactStoreBaseImm32:
			value := uint64(0)
			if instruction.destination < 31 {
				value = uint64(uint32(s.x[d]))
			}
			if err := s.storeValue(addOffset(s.x[left], int64(instruction.immediate)), value, 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactMul64:
			s.x[d] = s.x[left] * s.x[right]
		case compactSubImmShift32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(instruction.immediate)<<instruction.shift)
		case compactXorImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			s.x[d] = uint64(uint32(s.x[left]) ^ immediate)
		case compactTestRegister64:
			s.test64(s.x[left] & s.x[right])
		case compactLogicalMask64:
			width := instruction.shift
			if width < 2 || width > 64 || width&(width-1) != 0 || instruction.immediate == 0 || instruction.immediate >= uint16(width) || instruction.right >= width || instruction.accumulator > 2 {
				return fmt.Errorf("invalid logical-mask encoding")
			}
			mask := (uint64(1) << instruction.immediate) - 1
			for step := int(width); step < 64; step *= 2 {
				mask |= mask << step
			}
			mask = bits.RotateLeft64(mask, -int(instruction.right))
			switch instruction.accumulator {
			case 0:
				s.x[d] = s.x[left] & mask
			case 1:
				s.x[d] = s.x[left] | mask
			case 2:
				s.x[d] = s.x[left] ^ mask
			}
		case compactSubLSR32:
			s.x[d] = uint64(uint32(s.x[left]) - (uint32(s.x[right]) >> instruction.shift))
		case compactSubLSL32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactMovImm32Wide:
			s.x[d] = uint64(uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24)
		case compactAnd64:
			s.x[d] = s.x[left] & s.x[right]
		default:
			return fmt.Errorf("compact kernel %#x has unknown opcode %d", pc, instruction.op)
		}
		if hook < len(kernel.hooks) && kernel.hooks[hook].after == pc {
			if err := runCompactHook(s, kernel.hooks[hook].action); err != nil {
				return fmt.Errorf("compact kernel hook after %#x: %w", pc, err)
			}
			hook++
		}
		pc += 4
	}
	return nil
}

var compact26A428DigestPrologue = compactKernel{
	start: 0x1007A8214,
	instructions: []compactInstruction{
		{op: compactStoreSP32, left: 0, immediate: 3464},
		{op: compactStoreSP32, left: 16, immediate: 3520},
		{op: compactMovImm32, destination: 14, immediate: 0x92B8},
		{op: compactMovK32, destination: 14, shift: 16, immediate: 0xB039},
		{op: compactAdd32, destination: 14, left: 5, right: 14},
		{op: compactUDiv32, destination: 16, left: 14, right: 11},
		{op: compactMSub32, destination: 14, left: 16, right: 11, accumulator: 14},
		{op: compactLoadSP32, destination: 15, immediate: 3424},
		{op: compactAdd32, destination: 17, left: 15, right: 8},
		{op: compactUMull32, destination: 0, left: 17, right: 9},
		{op: compactLoadUXTW8, destination: 16, left: 30, right: 14},
		{op: compactLSRImm64, destination: 0, left: 0, shift: 34},
		{op: compactMSub32, destination: 17, left: 0, right: 12, accumulator: 17},
	},
}

var compact26A428BootstrapKernel = compactKernel{
	start: 0x10057AE74,
	instructions: []compactInstruction{
		{op: compactAnd32, destination: 11, left: 14, right: 10},
		{op: compactAdd32, destination: 17, left: 20, right: 13},
		{op: compactLoadSXTW32, destination: 2, left: 4, right: 17, shift: 2},
		{op: compactXor64, destination: 2, left: 2, right: 24},
		{op: compactAddLSR64, destination: 2, left: 2, right: 28, shift: 32},
		{op: compactLoadSXTW32, destination: 3, left: 5, right: 17, shift: 2},
		{op: compactXor64, destination: 3, left: 3, right: 24},
		{op: compactUMADDL, destination: 28, left: 3, right: 12, accumulator: 2},
		{op: compactXor64, destination: 2, left: 11, right: 23},
		{op: compactLSRV64, destination: 2, left: 30, right: 2},
		{op: compactXor64, destination: 2, left: 2, right: 6},
		{op: compactLoadSXTW32, destination: 17, left: 22, right: 17, shift: 2},
		{op: compactXor64, destination: 17, left: 17, right: 24},
		{op: compactUMADDL, destination: 17, left: 17, right: 19, accumulator: 2},
		{op: compactAddUXTW64, destination: 17, left: 17, right: 28},
		{op: compactAndLSL64, destination: 2, left: 7, right: 17, shift: 1},
		{op: compactSub64, destination: 17, left: 17, right: 2},
		{op: compactAdd64, destination: 30, left: 17, right: 26},
		{op: compactXor32, destination: 17, left: 30, right: 27},
		{op: compactAdd32, destination: 2, left: 20, right: 21},
		{op: compactStoreSXTW32, destination: 17, left: 4, right: 2, shift: 2},
		{op: compactAddImm32, destination: 2, left: 20, immediate: 1},
		{op: compactAdd32, destination: 17, left: 20, right: 16},
		{op: compactCmp32, left: 17, right: 15},
		{op: compactCCmp32LE, left: 17, right: 1, immediate: 8},
		{op: compactCSet32LT, destination: 17},
		{op: compactMul32, destination: 17, left: 17, right: 0},
		{op: compactXor32, destination: 17, left: 17, right: 9},
		{op: compactLoadSXTW64, destination: 3, left: 25, right: 17, shift: 3},
		{op: compactMov64, destination: 20, left: 2},
		{op: compactMovImm64, destination: 17, immediate: 0xCD95},
	},
}
