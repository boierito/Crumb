package rawsource

import (
	"encoding/binary"
	"testing"
)

func TestSourceOutputScratchRequiresPinnedContractAndStores(t *testing.T) {
	page := make([]byte, pageSize)
	for index, slot := range sourceOutputSlots {
		binary.LittleEndian.PutUint64(page[8+index*8:], slot.address)
	}
	value := &model{binaryHash: rawSourceBinaryHash, pages: map[uint64][]byte{0x700000100000: page}}
	if !sourceOutputContract(value) {
		t.Fatal("known caller output contract rejected")
	}
	page[8] ^= 1
	if sourceOutputContract(value) {
		t.Fatal("changed output pointer accepted")
	}
	page[8] ^= 1
	old := activeRawSourceRuntime
	defer func() { activeRawSourceRuntime = old }()
	r := &rawSourceRuntime{outputScratchEnabled: true}
	activeRawSourceRuntime = r
	s := newState(value)
	s.reset()
	for _, slot := range sourceOutputSlots {
		if _, err := s.load(slot.address, 1); err == nil {
			t.Fatal("uninitialized output was readable")
		}
		r.currentInstruction = slot.pc + 4
		if err := s.store(slot.address, make([]byte, slot.size)); err == nil {
			t.Fatal("incorrect writer accepted")
		}
		r.currentInstruction = slot.pc
		data := make([]byte, slot.size)
		for i := range data {
			data[i] = byte(i + 1)
		}
		if err := s.store(slot.address, data); err != nil {
			t.Fatal(err)
		}
		if _, err := s.load(slot.address, slot.size); err != nil {
			t.Fatal(err)
		}
		if _, err := s.load(slot.address+uint64(slot.size), 1); err == nil {
			t.Fatal("read outside output slot accepted")
		}
		if err := s.store(slot.address-1, make([]byte, slot.size+1)); err == nil {
			t.Fatal("crossing output slot accepted")
		}
	}
}
