package rawsource

import (
    "encoding/binary"
    "testing"
)

// This test contains no Apple account, license, keybag, SIDB or recovered key.
func supplementedBootstrapFixture(t *testing.T) (*state, uint64) {
    t.Helper()
    sp := uint64(0x71010000)
    base := uint64(0x70000000)
    table := base + 0x64e4*8
    target := uint64(0x1005b7438)
    pages := map[uint64][]byte{
        sp: make([]byte, pageSize),
        table &^ (pageSize-1): make([]byte, pageSize),
    }
    binary.LittleEndian.PutUint64(pages[table &^ (pageSize-1)][table&(pageSize-1):], target)
    binary.LittleEndian.PutUint64(pages[sp][0x1b8:], 0x1020304050607080)
    model := &model{pages:pages, sp:sp}
    model.registers[30] = base
    model.registers[12] = 0
    model.registers[4] = 0x12345678
    model.registers[27] = 0x88776655
    model.registers[25] = 0x11223344
    model.registers[20] = 0x66778899
    model.registers[21] = 0xaabbccdd
    model.registers[19] = 0x55667788
    s := newState(model)
    s.reset()
    return s, target
}

func TestPinnedBootstrapSupplementCompleteSynthetic(t *testing.T) {
    const entry uint64 = 0x1008fff2c
    if knownSourceBlocks[entry].run == nil {
        t.Fatal("public hash-pinned supplement missing from runtime registry")
    }
    s, target := supplementedBootstrapFixture(t)
    runtime := &rawSourceRuntime{}
    next, stopped, covered, err := resumeSourceBlock(s, runtime, entry, 0x1005b7438)
    if err != nil || stopped || !covered || next != target {
        t.Fatalf("synthetic supplemental branch: next=%x stopped=%v covered=%v err=%v", next, stopped, covered, err)
    }
    for _, check := range []struct{ address uint64; expected uint32 }{
        {0x18c, 0x12345678}, {0x190, 0x88776655}, {0x194, 0x11223344}, {0x118, 0x66778899},
    } {
        observed, err := s.load(s.sp+check.address, 4)
        if err != nil || uint32(observed) != check.expected {
            t.Fatalf("stack store %#x: %#x %v", check.address, observed, err)
        }
    }
    if s.x[27] != 0xaabbccdd || s.x[25] != 0x55667788 ||
        s.x[2] != 0 || s.x[9] != 0x64d2 || s.x[10] != 0xaf478f47 ||
        s.x[11] != 13 || s.x[14] != 111 || s.x[15] != 0x14307 ||
        s.x[16] != 176 || s.x[0] != 37 || s.x[1] != 0x10383d ||
        s.x[3] != 73 || s.x[4] != 165 || s.x[5] != 110 ||
        s.x[13] != 0x1020304050607080 || s.x[17] != 0x62bb ||
        runtime.lastKnownBlock != entry || runtime.additionalBlocks != 1 {
        t.Fatal("supplement differs from pinned 26A434 public instruction semantics")
    }
}

func TestPinnedBootstrapSupplementInteriorStopSynthetic(t *testing.T) {
    const entry uint64 = 0x1008fff2c
    s, _ := supplementedBootstrapFixture(t)
    runtime := &rawSourceRuntime{stopPC: 0x1008fff50}
    next, stopped, covered, err := resumeSourceBlock(s, runtime, entry, 0x1005b7438)
    if err != nil || !stopped || !covered || next != 0x1008fff50 || s.x[6] != 0 {
        t.Fatalf("interior stop incorrectly ran indirect load: next=%x stopped=%v covered=%v err=%v", next, stopped, covered, err)
    }
    if runtime.stopPC != 0x1008fff50 {
        t.Fatal("runtime's original stop boundary changed")
    }
}
