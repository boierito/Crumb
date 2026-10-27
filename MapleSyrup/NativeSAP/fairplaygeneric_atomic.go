package rawsource

import "fmt"

// CAS{A,L,AL}{B,H} executes in the single, serialized guest CPU. Acquire/release
// do not add an observable action in this private sequential memory model.
// CASP and all other LSE encodings deliberately remain unsupported.
func runGenericCAS(s *state, word uint32) (bool, error) {
	if word&0x3fa07c00 != 0x08a07c00 {
		return false, nil
	}
	size := 1 << (word >> 30)
	compare, desired, base := byte(word>>16&31), byte(word&31), byte(word>>5&31)
	address := s.sp
	if base != 31 {
		address = s.x[base]
	}
	if address&uint64(size-1) != 0 {
		return true, fmt.Errorf("unaligned guest compare-and-swap")
	}
	// Capture every input before writing Rs: the operands may alias.
	expected := genericReadRegister(s, compare, 64)
	replacement := genericReadRegister(s, desired, 64)
	mask := ^uint64(0) >> (64 - uint(size)*8)
	old, err := s.load(address, size)
	if err != nil {
		return true, err
	}
	if old == expected&mask {
		if err := s.storeValue(address, replacement, size); err != nil {
			return true, err
		}
	}
	if compare != 31 {
		s.x[compare] = old
	}
	return true, nil
}

// LDP/STP of two 64-bit SIMD registers, with signed offset, pre-index or
// post-index addressing. Other pair encodings remain unsupported.
func runGenericSIMDPair64(s *state, word uint32) (bool, error) {
	if word&0xfe000000 != 0x6c000000 {
		return false, nil
	}
	mode := (word >> 23) & 3
	if mode == 0 { // non-temporal LDNP/STNP has a distinct contract
		return false, nil
	}
	first, second, base := int(word&31), int(word>>10&31), int(word>>5&31)
	load := word&(1<<22) != 0
	if load && first == second {
		return true, fmt.Errorf("unpredictable SIMD pair destination")
	}
	offset := int64((word >> 15) & 127)
	if offset&64 != 0 {
		offset -= 128
	}
	offset *= 8
	address := s.sp
	if base != 31 {
		address = s.x[base]
	}
	updated := address + uint64(offset)
	if mode != 1 {
		address = updated
	}
	if load {
		data, err := s.outputBytes(address, 16)
		if err != nil {
			return true, err
		}
		s.clearVector(first)
		s.clearVector(second)
		copy(s.vectors[first][:8], data[:8])
		copy(s.vectors[second][:8], data[8:])
	} else {
		var data [16]byte
		copy(data[:8], s.vectors[first][:8])
		copy(data[8:], s.vectors[second][:8])
		if err := s.store(address, data[:]); err != nil {
			return true, err
		}
	}
	if mode != 2 {
		if base == 31 {
			s.sp = updated
		} else {
			s.x[base] = updated
		}
	}
	return true, nil
}

// Extended semantics may read instruction words only from the immutable text
// belonging to the profile's verified source binary, never guest-writable code.
func (p *genericInstructionProfile) attachNativeText(source [32]byte, base uint64, text []byte) error {
	if p == nil || p.sourceBinary != source || base%4 != 0 || len(text)%4 != 0 || base+uint64(len(text)) < base {
		return fmt.Errorf("generic native text binding invalid")
	}
	p.nativeBase, p.nativeText = base, text
	return nil
}
