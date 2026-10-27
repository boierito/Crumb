package rawsource

import "encoding/binary"

// These are caller-owned scalar output slots in the pinned public account
// entry snapshot, not license data. The captured model omits their pages.
// Only the exact native return stores may create them, and reads are limited
// to initialized bytes. Unknown snapshots and all other addresses fail closed.
var sourceOutputSlots = [...]struct {
	address, pc uint64
	size        int
}{
	{0x700000240000, 0x1000e0580, 8},
	{0x700000250000, 0x1000e0588, 4},
	{0x700000260000, 0x1000e0594, 8},
	{0x700000270000, 0x1000e059c, 4},
}

func sourceOutputContract(value *model) bool {
	if value == nil || value.binaryHash != rawSourceBinaryHash {
		return false
	}
	page := value.pages[0x700000100000]
	if len(page) != int(pageSize) {
		return false
	}
	for index, slot := range sourceOutputSlots {
		if binary.LittleEndian.Uint64(page[8+index*8:]) != slot.address || value.pages[slot.address] != nil {
			return false
		}
	}
	return true
}

func sourceOutputRead(s *state, address uint64, size int) error {
	r := activeRawSourceRuntime
	if r == nil || !r.outputScratchEnabled {
		return nil
	}
	if size > 0 && address+uint64(size) < address {
		return &sourceMemoryError{address: address}
	}
	if size <= 0 || address >= 0x700000271000 || address+uint64(size) <= 0x700000240000 {
		return nil
	}
	for index, slot := range sourceOutputSlots {
		if size > 0 && address+uint64(size) >= address && address < slot.address+pageSize && address+uint64(size) > slot.address {
			if size < 0 || address < slot.address || address-slot.address+uint64(size) > uint64(slot.size) || r.outputScratchReady&(1<<index) == 0 {
				return &sourceMemoryError{address: address}
			}
			return nil
		}
	}
	return nil
}

func sourceOutputStore(s *state, address uint64, data []byte) error {
	r := activeRawSourceRuntime
	if r == nil || !r.outputScratchEnabled {
		return nil
	}
	if address+uint64(len(data)) < address {
		return &sourceMemoryError{address: address, write: true}
	}
	if len(data) == 0 || address >= 0x700000271000 || address+uint64(len(data)) <= 0x700000240000 {
		return nil
	}
	for index, slot := range sourceOutputSlots {
		if len(data) > 0 && address+uint64(len(data)) >= address && address < slot.address+pageSize && address+uint64(len(data)) > slot.address {
			if address != slot.address || len(data) != slot.size || r.currentInstruction != slot.pc {
				return &sourceMemoryError{address: address, write: true}
			}
			if r.outputScratchReady&(1<<index) == 0 {
				if s.model.pages[slot.address] != nil {
					return &sourceMemoryError{address: address, write: true}
				}
				if err := s.mapZero(slot.address, pageSize); err != nil {
					return err
				}
				r.outputScratchReady |= 1 << index
			}
			return nil
		}
	}
	return nil
}
