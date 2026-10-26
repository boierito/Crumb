package rawsource

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

type nextBlockFixture struct {
	InputRegisters [31]uint64 `json:"input_registers"`
	InputFlags     [4]bool    `json:"input_flags"`
	Registers      [31]uint64 `json:"registers"`
	Flags          [4]bool    `json:"flags"`
	Next           uint64     `json:"next"`
	OutputSHA256   string     `json:"output_sha256"`
}

func nextBlockState(t *testing.T, fixture nextBlockFixture) *state {
	t.Helper()
	pages := map[uint64][]byte{}
	for _, base := range []uint64{0x300000, 0x400000, 0x500000} {
		for offset := uint64(0); offset < 65536; offset += pageSize {
			page := make([]byte, pageSize)
			for j := range page {
				page[j] = byte(base/65536 + offset + uint64(j)*7)
			}
			pages[base+offset] = page
		}
	}
	pages[0x600000] = make([]byte, pageSize)
	binary.LittleEndian.PutUint64(pages[0x600000], 0x200000)
	binary.LittleEndian.PutUint64(pages[0x600000][110*8:], 0x200004)
	m := &model{pages: pages, registers: fixture.InputRegisters,
		n: fixture.InputFlags[0], z: fixture.InputFlags[1], c: fixture.InputFlags[2], v: fixture.InputFlags[3]}
	s := newState(m)
	s.reset()
	return s
}

// Expected full registers, flags and output hashes come from independent
// Unicorn ARM64 execution with synthetic memory. BRAA uses the existing
// runtime convention of branching to the table target without native PAC.
func TestBootstrapSupplement100a4be38IndependentSynthetic(t *testing.T) {
	data, err := os.ReadFile("testdata/bootstrap-next.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []nextBlockFixture
	if err = json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures) != 10 {
		t.Fatal("missing public reference vectors")
	}
	for i, fixture := range fixtures {
		s := nextBlockState(t, fixture)
		runtime := &rawSourceRuntime{}
		next, stopped, covered, err := resumeSourceBlock(s, runtime, 0x100a4be38, 0x1005b7438)
		if err != nil || stopped || !covered || next != fixture.Next {
			t.Fatalf("fixture %d: dispatch mismatch: %v", i, err)
		}
		flags := [4]bool{s.n, s.z, s.c, s.v}
		if s.x != fixture.Registers || flags != fixture.Flags {
			t.Fatalf("fixture %d: registers/flags differ from ARM64", i)
		}
		output, err := s.outputBytes(0x500000, 65536)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(output)
		if hex.EncodeToString(digest[:]) != fixture.OutputSHA256 {
			t.Fatalf("fixture %d: memory differs from ARM64", i)
		}
		if runtime.lastKnownBlock != 0x100a4be38 || runtime.additionalBlocks != 1 {
			t.Fatal("missing dispatcher bookkeeping")
		}
	}
}

func TestBootstrapSupplement100a4be38StopsAndRejectsMissingMemory(t *testing.T) {
	s := nextBlockState(t, nextBlockFixture{})
	before := s.x
	next, stopped, err := rawSourceBootstrapBlock100a4be38(s, &rawSourceRuntime{stopPC: 0x100a4be38})
	if err != nil || !stopped || next != 0x100a4be38 || s.x != before {
		t.Fatal("entry stop modified state")
	}
	if _, _, err = rawSourceBootstrapBlock100a4be38(s, nil); err == nil {
		t.Fatal("unmapped table accepted")
	}
}

func TestBootstrapSupplement100a4be38SignedDispatchIndex(t *testing.T) {
	var fixture nextBlockFixture
	fixture.InputRegisters[0] = 37
	fixture.InputRegisters[1] = 0x10383d
	fixture.InputRegisters[3] = 73
	fixture.InputRegisters[4] = 165
	fixture.InputRegisters[5] = 110
	fixture.InputRegisters[9] = 0xffffffff
	fixture.InputRegisters[11] = 13
	fixture.InputRegisters[13] = 0x300000
	fixture.InputRegisters[14] = 111
	fixture.InputRegisters[15] = 0x14307
	fixture.InputRegisters[16] = 176
	fixture.InputRegisters[22] = 0x500000
	fixture.InputRegisters[23] = 0x400000
	fixture.InputRegisters[30] = 0x600800
	s := nextBlockState(t, fixture)
	// cset=1: (110 xor 0xffffffff), sign extended, is -111.
	binary.LittleEndian.PutUint64(s.model.pages[0x600000][0x800-111*8:], 0x200008)
	next, stopped, err := rawSourceBootstrapBlock100a4be38(s, nil)
	if err != nil || stopped || next != 0x200008 {
		t.Fatalf("signed dispatch: next=%x err=%v", next, err)
	}
}
