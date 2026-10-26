package rawsource

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func bootstrapExitFixture(t *testing.T, index int32) *state {
	t.Helper()
	const stack = uint64(0x71010000)
	const table = uint64(0x600800)
	pages := map[uint64][]byte{stack: bytes.Repeat([]byte{0xa5}, int(pageSize)), 0x600000: make([]byte, pageSize)}
	address := uint64(int64(table) + int64(index)*8)
	if address < 0x600000 || address+8 > 0x601000 {
		t.Fatal("test index outside mapped table")
	}
	binary.LittleEndian.PutUint64(pages[0x600000][address-0x600000:], 0x200008)
	m := &model{pages: pages, sp: stack, n: true, z: false, c: true, v: true}
	for i := range m.registers {
		m.registers[i] = uint64(i)*0x102030405060708 + 0x8877665544332211
	}
	// The W-register addition must wrap before the signed table-index extension.
	bias := uint32(index) - uint32(0x2aa3)
	m.registers[12] = 0xfeedbeef00000000 | uint64(bias)
	m.registers[24] = 0xaabbccdd12345678
	m.registers[30] = table
	s := newState(m)
	s.reset()
	return s
}

func TestBootstrapSupplement1008c292cFullStateAndSignedDispatch(t *testing.T) {
	for _, index := range []int32{-256, -64, -1, 0, 1, 63, 255} {
		s := bootstrapExitFixture(t, index)
		before := s.x
		expected := before
		expected[9] = 0x200008
		expected[17] = 0x1a88
		expected[21] = before[25]
		expected[25] = before[27]
		memory, err := s.outputBytes(s.sp, int(pageSize))
		if err != nil {
			t.Fatal(err)
		}
		binary.LittleEndian.PutUint32(memory[0x11c:], uint32(before[24]))
		runtime := &rawSourceRuntime{}
		next, stopped, covered, err := resumeSourceBlock(s, runtime, 0x1008c292c, 0x1005b7438)
		if err != nil || stopped || !covered || next != 0x200008 {
			t.Fatalf("index %d: dispatch failed: %v", index, err)
		}
		if s.x != expected || !s.n || s.z || !s.c || !s.v || s.sp != 0x71010000 {
			t.Fatalf("index %d: register/flag/SP mismatch", index)
		}
		observed, err := s.outputBytes(s.sp, int(pageSize))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(observed, memory) {
			t.Fatalf("index %d: incorrect stack write or extra mutation", index)
		}
		if runtime.lastKnownBlock != 0x1008c292c || runtime.additionalBlocks != 1 {
			t.Fatal("missing dispatcher bookkeeping")
		}
	}
}

func TestBootstrapSupplement1008c292cInteriorStop(t *testing.T) {
	s := bootstrapExitFixture(t, 0)
	before := s.x
	next, stopped, err := rawSourceBootstrapBlock1008c292c(s, &rawSourceRuntime{stopPC: 0x1008c2930})
	if err != nil || !stopped || next != 0x1008c2930 || s.x != before {
		t.Fatal("interior stop executed beyond stack store")
	}
	word, err := s.load(s.sp+0x11c, 4)
	if err != nil || word != uint64(uint32(before[24])) {
		t.Fatal("stack write missing at interior stop")
	}
}

func TestBootstrapSupplement1008c292cRejectsUnmappedMemory(t *testing.T) {
	for _, missing := range []uint64{0x71010000, 0x600000} {
		s := bootstrapExitFixture(t, 0)
		delete(s.model.pages, missing)
		before := s.x
		if _, _, err := rawSourceBootstrapBlock1008c292c(s, nil); err == nil {
			t.Fatal("unmapped stack or table accepted")
		}
		if s.x[21] != before[21] || s.x[25] != before[25] {
			t.Fatal("continued register restoration after rejected load")
		}
	}
}
