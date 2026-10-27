package rawsource

import (
	"encoding/binary"
	"testing"
)

func TestSourceOutputContractMatchesPinnedImage(t *testing.T) {
	value, err := loadModelBytes(rawSourceModel)
	if err != nil || !sourceOutputContract(value) {
		t.Fatalf("pinned caller output contract unavailable: %v", err)
	}
	if genericImage == nil || len(genericImage.nativeText) == 0 {
		t.Fatal("pinned native text unavailable")
	}
	for index, slot := range sourceOutputSlots {
		offset := slot.pc - genericImage.nativeBase
		if offset+4 > uint64(len(genericImage.nativeText)) {
			t.Fatal("output store outside pinned text")
		}
		word := binary.LittleEndian.Uint32(genericImage.nativeText[offset:])
		expected := uint32(0xf9000169) // str x9, [x11]
		if index%2 != 0 {
			expected = 0xb9000149 // str w9, [x10]
		}
		if word != expected {
			t.Fatalf("output store %d does not match pinned instruction", index)
		}
	}
}
