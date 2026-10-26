package recovery

import "fmt"

type compactBlockTerminator uint8

const (
	compactBlockDirect compactBlockTerminator = iota
	compactBlockConditional
	compactBlockRegister
	compactBlockCallRegister
	compactBlockCallDirect
	compactBlockReturn
	compactBlockCompareZero
	compactBlockTestBit
	compactBlockStopDirect
)

type compactCondition uint8

const (
	compactConditionEQ compactCondition = iota
	compactConditionNE
	compactConditionHS
	compactConditionLO
	compactConditionMI
	compactConditionPL
	compactConditionVS
	compactConditionVC
	compactConditionHI
	compactConditionLS
	compactConditionGE
	compactConditionLT
	compactConditionGT
	compactConditionLE
)

type compactBlock struct {
	start        uint64
	instructions []compactInstruction
	terminator   compactBlockTerminator
	condition    compactCondition
	register     uint8
	width        uint8
	bit          uint8
	target       uint64
	alternate    uint64
}

func compactConditionMatches(s *state, condition compactCondition) (bool, error) {
	if condition > compactConditionLE {
		return false, fmt.Errorf("unknown compact block condition %d", condition)
	}
	name := [...]string{"eq", "ne", "hs", "lo", "mi", "pl", "vs", "vc", "hi", "ls", "ge", "lt", "gt", "le"}[condition]
	return s.condition(name), nil
}

func runCompactBlock(s *state, runtime *rawSourceRuntime, block compactBlock) (next uint64, stopped bool, err error) {
	terminalPC := block.start + uint64(len(block.instructions))*4
	if runtime != nil && runtime.stopPC >= block.start && runtime.stopPC <= terminalPC && (runtime.stopPC-block.start)%4 == 0 {
		instructionCount := int((runtime.stopPC - block.start) / 4)
		if instructionCount != 0 {
			prefix := compactKernel{start: block.start, instructions: block.instructions[:instructionCount]}
			if err := runCompactKernel(s, prefix); err != nil {
				return 0, false, err
			}
		}
		return runtime.stopPC, true, nil
	}
	if err := runCompactKernel(s, compactKernel{start: block.start, instructions: block.instructions}); err != nil {
		return 0, false, err
	}
	switch block.terminator {
	case compactBlockDirect:
		return block.target, false, nil
	case compactBlockConditional:
		matches, err := compactConditionMatches(s, block.condition)
		if err != nil {
			return 0, false, err
		}
		if matches {
			return block.target, false, nil
		}
		return block.alternate, false, nil
	case compactBlockRegister:
		return s.x[block.register], false, nil
	case compactBlockCallRegister:
		s.x[30] = terminalPC + 4
		return s.x[block.register], false, nil
	case compactBlockCallDirect:
		s.x[30] = terminalPC + 4
		return block.target, false, nil
	case compactBlockReturn:
		return s.x[30], false, nil
	case compactBlockCompareZero:
		value := s.x[block.register]
		if block.width == 32 {
			value = uint64(uint32(value))
		}
		if block.condition != compactConditionEQ && block.condition != compactConditionNE {
			return 0, false, fmt.Errorf("invalid compare-zero condition %d", block.condition)
		}
		zero := value == 0
		if block.condition == compactConditionEQ && zero || block.condition == compactConditionNE && !zero {
			return block.target, false, nil
		}
		return block.alternate, false, nil
	case compactBlockTestBit:
		if block.condition != compactConditionEQ && block.condition != compactConditionNE {
			return 0, false, fmt.Errorf("invalid test-bit condition %d", block.condition)
		}
		set := s.x[block.register]&(uint64(1)<<block.bit) != 0
		if block.condition == compactConditionNE && set || block.condition == compactConditionEQ && !set {
			return block.target, false, nil
		}
		return block.alternate, false, nil
	case compactBlockStopDirect:
		return block.target, true, nil
	default:
		return 0, false, fmt.Errorf("unknown compact block terminator %d", block.terminator)
	}
}

func compact26A428SINFEntryBlock(pc uint64) (compactBlock, bool) {
	for _, blocks := range []map[uint64]compactBlock{
		compact26A428SINFEntryBlocks00,
		compact26A428SINFEntryBlocks01,
		compact26A428SINFEntryBlocks02,
		compact26A428SINFEntryBlocks03,
		compact26A428SINFEntryBlocks04,
		compact26A428SINFEntryBlocks05,
	} {
		if block, ok := blocks[pc]; ok {
			return block, true
		}
	}
	return compactBlock{}, false
}
