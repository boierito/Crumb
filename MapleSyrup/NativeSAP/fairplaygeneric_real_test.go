package rawsource

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// Preserve the signed-dispatch device regression after its AOT promotion.
// Force the generic fallback as well as exercising the configured dispatcher.
func TestGenericRealProfileCurrentDeviceLanding(t *testing.T) {
	const pc = uint64(0x100a06948)
	const sp = uint64(0x71010000)
	const table = uint64(0x600800)
	if genericImage == nil || !genericImage.contains(pc) {
		t.Fatal("real generic profile missing current device continuation")
	}
	_, hasAOT := knownSourceBlocks[pc]
	for _, genericOnly := range []bool{false, true} {
		for _, index := range []int32{-256, -1, 0, 255} {
			stack := bytes.Repeat([]byte{0xa5}, int(pageSize))
			dispatch := make([]byte, pageSize)
			offset := int64(0x800) + int64(index)*8
			binary.LittleEndian.PutUint64(dispatch[offset:], 0x1005b7438)
			m := &model{pages: map[uint64][]byte{sp: stack, 0x600000: dispatch}, sp: sp, n: true, z: false, c: true, v: true}
			for i := range m.registers {
				m.registers[i] = uint64(i)*0x102030405060708 + 0x8877665544332211
			}
			bias := uint32(index) - uint32(0xffffcb0d)
			m.registers[12] = uint64(bias) | 0xfeedbeef00000000
			m.registers[30] = table
			expected := m.registers
			expected[5] = 0
			expected[9] = uint64(uint32(index))
			expected[10] = 148
			expected[11] = uint64((uint32(index) + 24391) | 8202)
			expected[12] = 0x1005b7438
			expected[14] = 57457
			expected[15] = 73
			expected[16] = 0xffff5bc6
			expected[0] = 134
			expected[1] = 8978
			expected[2] = 29747
			expected[3] = 141
			expected[4] = 32977
			expected[17] = 15887
			memory := append([]byte(nil), stack...)
			for _, pair := range [][2]int{{0x198, 26}, {0x188, 21}, {0x184, 25}, {0x178, 28}} {
				binary.LittleEndian.PutUint32(memory[pair[0]:], uint32(m.registers[pair[1]]))
			}
			s := newState(m)
			s.reset()
			runtime := &rawSourceRuntime{genericOnly: genericOnly}
			next, stopped, covered, err := resumeSourceBlock(s, runtime, pc, 0x1005b7438)
			if err != nil || stopped || !covered || next != 0x1005b7438 {
				t.Fatalf("index %d: generic continuation failed: %v", index, err)
			}
			if s.x != expected || s.sp != sp || !s.n || s.z || !s.c || !s.v {
				t.Fatalf("index %d: register/flags mismatch", index)
			}
			output, err := s.outputBytes(sp, int(pageSize))
			if err != nil || !bytes.Equal(output, memory) {
				t.Fatalf("index %d: stack mismatch", index)
			}
			if genericOnly || !hasAOT {
				if runtime.genericBlocks != 1 || runtime.genericInstructions != 23 {
					t.Fatal("generic interpreter was not the executing backend")
				}
			} else if runtime.genericBlocks != 0 || runtime.genericInstructions != 0 {
				t.Fatal("registered AOT dispatcher unexpectedly used generic execution")
			}
		}
	}
}

func TestGenericDispatcherPreservesInteriorPhaseBoundary(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{genericNormal(compactMovImm64, 0, 0, 0, 7), genericNormal(compactAddImm64, 0, 0, 0, 1), genericBranch(compactBlockRegister, 5, 0)})
	s := genericSyntheticState()
	s.x[5] = 0x200000
	runtime := &rawSourceRuntime{}
	next, stopped, covered, err := resumeSourceBlock(s, runtime, base, base+4)
	if err != nil || stopped || !covered || next != base+4 || s.x[0] != 7 || runtime.stopPC != 0 {
		t.Fatal("generic chunk skipped or leaked phase boundary")
	}
	next, stopped, covered, err = resumeSourceBlock(s, runtime, base+4, base+4)
	if err != nil || stopped || !covered || next != 0x200000 || s.x[0] != 8 || runtime.stopPC != 0 {
		t.Fatal("generic boundary entry could not resume")
	}
}
