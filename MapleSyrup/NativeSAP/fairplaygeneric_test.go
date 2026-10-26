package rawsource

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func genericTestImage(t *testing.T, base uint64, records [][genericRecordBytes]byte) *genericInstructionProfile {
	t.Helper()
	data := make([]byte, 96)
	copy(data, "CRIR1\x00\x00\x00")
	binary.LittleEndian.PutUint32(data[8:], 2)
	binary.LittleEndian.PutUint32(data[12:], genericRecordBytes)
	source := [32]byte{1, 2, 3}
	copy(data[16:], source[:])
	abi, err := hex.DecodeString(genericKernelABI)
	if err != nil {
		t.Fatal(err)
	}
	copy(data[48:], abi)
	binary.LittleEndian.PutUint64(data[80:], base)
	binary.LittleEndian.PutUint32(data[88:], uint32(len(records)))
	for _, record := range records {
		data = append(data, record[:]...)
	}
	sum := sha256.Sum256(data)
	data = append(data, sum[:]...)
	var compressed bytes.Buffer
	stream := gzip.NewWriter(&compressed)
	if _, err = stream.Write(data); err != nil {
		t.Fatal(err)
	}
	if err = stream.Close(); err != nil {
		t.Fatal(err)
	}
	image, err := loadGenericInstructionProfile(compressed.Bytes(), source, genericKernelABI)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = loadGenericInstructionProfile(compressed.Bytes(), [32]byte{99}, genericKernelABI); err == nil {
		t.Fatal("wrong source accepted")
	}
	if _, err = loadGenericInstructionProfile(compressed.Bytes(), source, "wrong-ABI"); err == nil {
		t.Fatal("wrong ABI accepted")
	}
	data[len(data)-1] ^= 1
	var broken bytes.Buffer
	stream = gzip.NewWriter(&broken)
	stream.Write(data)
	stream.Close()
	if _, err = loadGenericInstructionProfile(broken.Bytes(), source, genericKernelABI); err == nil {
		t.Fatal("corrupt profile accepted")
	}
	return image
}
func genericNormal(op compactOpcode, d, l, r byte, immediate uint16) (record [genericRecordBytes]byte) {
	record[0] = 1
	binary.LittleEndian.PutUint16(record[1:3], uint16(op))
	record[3] = d
	record[4] = l
	record[5] = r
	binary.LittleEndian.PutUint16(record[8:10], immediate)
	return
}
func genericBranch(term compactBlockTerminator, reg byte, target uint64) (record [genericRecordBytes]byte) {
	record[0] = 2
	record[1] = byte(term)
	record[5] = reg
	binary.LittleEndian.PutUint64(record[8:], target)
	return
}
func genericSyntheticState() *state {
	m := &model{pages: map[uint64][]byte{0x700000: make([]byte, pageSize)}, sp: 0x700100}
	s := newState(m)
	s.reset()
	return s
}

func TestGenericUnknownAddressesAndIndirectInteriorLandings(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	for _, base := range []uint64{0x12345000, 0x98765000} {
		records := [][genericRecordBytes]byte{
			genericNormal(compactMovImm64, 0, 0, 0, 40), genericNormal(compactAddImm64, 0, 0, 0, 2), genericNormal(compactStoreSP64, 0, 0, 0, 4), genericBranch(compactBlockRegister, 5, 0), {}, {}, genericNormal(compactAddImm64, 0, 0, 0, 1), genericBranch(compactBlockReturn, 0, 0),
		}
		genericImage = genericTestImage(t, base, records)
		s := genericSyntheticState()
		s.x[5] = base + 24
		s.x[30] = 0x200000
		r := &rawSourceRuntime{}
		next, stop, err := runGenericInstructionBlock(s, r, base)
		if err != nil || stop || next != base+24 || s.x[0] != 42 {
			t.Fatalf("first indirect chunk: %x %v", next, err)
		}
		stack, err := s.load(s.sp+4, 8)
		if err != nil || stack != 42 {
			t.Fatal("generic memory store mismatch")
		}
		next, stop, err = runGenericInstructionBlock(s, r, next)
		if err != nil || stop || next != 0x200000 || s.x[0] != 43 || r.genericBlocks != 2 || r.genericInstructions != 6 {
			t.Fatalf("interior landing: %x %v", next, err)
		}
	}
}
func TestGenericCallThroughLRPreservesOriginalTarget(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{genericBranch(compactBlockCallRegister, 30, 0)})
	s := genericSyntheticState()
	s.x[30] = 0x200008
	next, stop, err := runGenericInstructionBlock(s, &rawSourceRuntime{}, base)
	if err != nil || stop || next != 0x200008 || s.x[30] != base+4 {
		t.Fatal("BLR LR did not capture target before updating link register")
	}
}
func TestGenericConditionalAndStops(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	record := genericBranch(compactBlockConditional, 0, 0x200000)
	record[4] = byte(compactConditionEQ)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{record})
	for _, zero := range []bool{false, true} {
		s := genericSyntheticState()
		s.z = zero
		expected := base + 4
		if zero {
			expected = 0x200000
		}
		next, _, err := runGenericInstructionBlock(s, &rawSourceRuntime{}, base)
		if err != nil || next != expected || s.z != zero {
			t.Fatal("conditional branch changed flags or selected wrong target")
		}
	}
	s := genericSyntheticState()
	before := s.x
	next, stop, err := runGenericInstructionBlock(s, &rawSourceRuntime{stopPC: base}, base)
	if err != nil || !stop || next != base || s.x != before {
		t.Fatal("stopPC executed an instruction")
	}
}
func TestGenericRejectsUnsupportedAndOutsideSlots(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{{}})
	for _, pc := range []uint64{base, base + 1, base + 4, base - 4} {
		s := genericSyntheticState()
		before := s.x
		if _, _, err := runGenericInstructionBlock(s, &rawSourceRuntime{}, pc); err == nil {
			t.Fatal("unsupported/outside instruction accepted")
		}
		if s.x != before {
			t.Fatal("rejected instruction modified registers")
		}
	}
}
func TestGenericSequentialBudget(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	records := make([][genericRecordBytes]byte, 5000)
	for i := range records {
		records[i] = genericNormal(compactNoop, 0, 0, 0, 0)
	}
	genericImage = genericTestImage(t, base, records)
	r := &rawSourceRuntime{}
	next, stop, err := runGenericInstructionBlock(genericSyntheticState(), r, base)
	if err != nil || stop || next != base+4096*4 || r.genericInstructions != 4096 {
		t.Fatal("sequential execution ignored chunk budget")
	}
}

func TestGenericCumulativeBudgetAcrossChunks(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{genericNormal(compactNoop, 0, 0, 0, 0), genericNormal(compactNoop, 0, 0, 0, 0)})
	runtime := &rawSourceRuntime{genericInstructions: genericMaxInstructions - 1}
	if _, _, err := runGenericInstructionBlock(genericSyntheticState(), runtime, base); err == nil {
		t.Fatal("cumulative execution budget ignored")
	}
	if runtime.genericInstructions != genericMaxInstructions {
		t.Fatal("executed beyond cumulative instruction budget")
	}
}

func TestGenericNonExecutableGapsAreNotCode(t *testing.T) {
	previous := genericImage
	defer func() { genericImage = previous }()
	const base = uint64(0x12345000)
	gap := [genericRecordBytes]byte{6}
	genericImage = genericTestImage(t, base, [][genericRecordBytes]byte{genericNormal(compactNoop, 0, 0, 0, 0), gap, genericBranch(compactBlockReturn, 0, 0)})
	if genericImage.contains(base+4) || !genericImage.contains(base) || !genericImage.contains(base+8) {
		t.Fatal("section gap treated as executable code")
	}
	if genericImage.executable != 2 || genericImage.supported != 2 {
		t.Fatal("section gaps counted as unsupported instructions")
	}
	if _, _, err := runGenericInstructionBlock(genericSyntheticState(), &rawSourceRuntime{}, base+4); err == nil {
		t.Fatal("non-executable gap executed")
	}
}
