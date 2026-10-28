package rawsource

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestGenericLoopObservationAggregatesWithoutCounterValues(t *testing.T) {
	var observation genericLoopObservation
	for _, value := range []uint64{3, 2, 1, 4097, 4096, 0, 1 << 21} {
		observation.observe(value)
	}
	if observation.visits != 7 || observation.sequences != 4 || observation.maxBucket != "over_1m" {
		t.Fatalf("incorrect loop aggregation: %#v", observation)
	}
}

func TestCASRejectsNonCASAndFaultsWithoutPartialWrites(t *testing.T) {
	m := &model{pages: map[uint64][]byte{0x600000: bytes.Repeat([]byte{0xa5}, 4096)}}
	s := newState(m)
	s.reset()
	s.x[19], s.x[11], s.x[12] = 0x600001, 7, 19
	before := s.x
	if ok, err := runGenericCAS(s, 0x88ebfe6c); !ok || err == nil {
		t.Fatal("unaligned CAS accepted")
	}
	if s.x != before || len(s.writable) != 0 {
		t.Fatal("alignment fault modified state")
	}
	s.x[19] = 0x800000
	before = s.x
	if ok, err := runGenericCAS(s, 0x88ebfe6c); !ok || err == nil {
		t.Fatal("unmapped CAS accepted")
	}
	if s.x != before || len(s.writable) != 0 {
		t.Fatal("memory fault modified state")
	}
	for _, word := range []uint32{0, 0xd503201f, 0x48207c82, 0xb8200020, 0xc85f7c20, 0xc8007c20} {
		if ok, err := runGenericCAS(s, word); ok || err != nil {
			t.Fatalf("non-CAS encoding %08x recognized", word)
		}
	}
}

func TestCASNativeDeviceInstructionAndWholeImageCoverage(t *testing.T) {
	if genericImage == nil || len(genericImage.nativeText) == 0 {
		t.Fatal("verified native text not attached")
	}
	const pc = uint64(0x1009885a0)
	if word := binary.LittleEndian.Uint32(genericImage.nativeText[pc-genericImage.nativeBase:]); word != 0x88ebfe6c {
		t.Fatal("native device instruction changed")
	}
	m := &model{pages: map[uint64][]byte{0x600000: make([]byte, 4096)}}
	binary.LittleEndian.PutUint32(m.pages[0x600000], 7)
	s := newState(m)
	s.reset()
	s.x[19], s.x[11], s.x[12] = 0x600000, 7, 19
	runtime := &rawSourceRuntime{stopPC: pc + 4, genericOnly: true}
	next, stopped, err := runGenericInstructionBlock(s, runtime, pc)
	value, _ := s.load(0x600000, 4)
	if err != nil || !stopped || next != pc+4 || s.x[11] != 7 || value != 19 || runtime.genericInstructions != 1 {
		t.Fatal("device CASAL does not execute through generic dispatcher")
	}
	count := 0
	for address := genericImage.base; address < genericImage.base+uint64(genericImage.count)*4; address += 4 {
		if !genericImage.contains(address) || address < genericImage.nativeBase || address-genericImage.nativeBase+4 > uint64(len(genericImage.nativeText)) {
			continue
		}
		word := binary.LittleEndian.Uint32(genericImage.nativeText[address-genericImage.nativeBase:])
		if word&0x3fa07c00 == 0x08a07c00 {
			count++
			if genericImage.record(address)[0] != 0 {
				t.Fatal("CAS coverage census changed")
			}
		}
	}
	if count != 61 {
		t.Fatalf("native CAS census: %d, want 61", count)
	}
}

func TestNativeTextRejectsWrongSource(t *testing.T) {
	p := &genericInstructionProfile{sourceBinary: [32]byte{1}}
	if p.attachNativeText([32]byte{2}, 0x1000, make([]byte, 4)) == nil || len(p.nativeText) != 0 {
		t.Fatal("mismatched instruction source accepted")
	}
}
