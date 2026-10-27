package rawsource

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"testing"
)

// Exercise the entire exact image and the entire public snapshot together.
// These are build-global bytes, not an enrollment or a decrypted IPA fixture.
func TestImmutableImagePreservesSnapshotAndCompletesEveryTextPage(t *testing.T) {
	image, err := immutableImageText()
	if err != nil {
		t.Fatal(err)
	}
	value, err := loadModelBytes(rawSourceModel)
	if err != nil {
		t.Fatal(err)
	}
	original := make(map[uint64][32]byte, len(value.pages))
	for address, page := range value.pages {
		original[address] = sha256.Sum256(page)
	}
	if len(original) != 472 || uint64(len(image)) != imageTextSize {
		t.Fatal("pinned snapshot/image geometry changed")
	}
	if err := installImmutableText(value, image, imageTextAddress); err != nil {
		t.Fatal(err)
	}
	added := 0
	for offset := uint64(0); offset < uint64(len(image)); offset += pageSize {
		address := imageTextAddress + offset
		page := value.pages[address]
		if len(page) != int(pageSize) {
			t.Fatal("image still has a missing page")
		}
		if _, captured := original[address]; !captured {
			added++
			if !bytes.Equal(page, image[offset:offset+pageSize]) {
				t.Fatal("restored page does not match exact image bytes")
			}
		}
	}
	if added != 3853 || len(value.pages) != len(original)+added {
		t.Fatal("image restoration escaped its immutable segment")
	}
	for address, digest := range original {
		if sha256.Sum256(value.pages[address]) != digest {
			t.Fatal("restoration overwrote a captured page")
		}
	}
	// Same-process later recoveries must retain the pristine resource even
	// after the first recovery destroys all its private model pages.
	for _, page := range value.pages {
		clear(page)
	}
	if sha256.Sum256(image) != imageTextRawSHA {
		t.Fatal("model cleanup destroyed the shared image resource")
	}
	second, err := loadModelBytes(rawSourceModel)
	if err != nil || installImmutableText(second, image, imageTextAddress) != nil {
		t.Fatal("second recovery could not restore the same exact image")
	}
	if len(second.pages) != len(original)+added {
		t.Fatal("second recovery changed image coverage")
	}
	worker := newState(second)
	worker.reset()
	if _, err := worker.load(0x600000000000, 8); err == nil {
		t.Fatal("an unknown heap pointer became readable")
	}
}

func TestImmutableImageRejectsAlteredResourceAndMapping(t *testing.T) {
	image := bytes.Repeat([]byte{0x3d}, int(pageSize))
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, _ = writer.Write(image)
	_ = writer.Close()
	encoded := buffer.Bytes()
	digest := sha256.Sum256(image)
	if _, err := decodeImmutableText(encoded, digest, pageSize); err != nil {
		t.Fatal(err)
	}
	wrong := digest
	wrong[0] ^= 1
	for _, trial := range []struct {
		encoded []byte
		digest  [32]byte
		size    uint64
	}{
		{encoded, wrong, pageSize}, {encoded[:len(encoded)-1], digest, pageSize},
		{encoded, digest, 0}, {encoded, digest, pageSize - 1},
		{encoded, digest, 33 << 20}, {encoded, digest, 2 * pageSize},
	} {
		if out, err := decodeImmutableText(trial.encoded, trial.digest, trial.size); err == nil || out != nil {
			t.Fatal("invalid resource accepted")
		}
	}
	value := &model{pages: map[uint64][]byte{}}
	if installImmutableText(value, image, imageTextAddress) == nil || len(value.pages) != 0 {
		t.Fatal("different binary accepted or changed")
	}
	value.binaryHash = rawSourceBinaryHash
	if installImmutableText(value, image, imageTextAddress+pageSize) == nil || len(value.pages) != 0 {
		t.Fatal("different image base accepted or changed")
	}
}
