package rawsource

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"sync"
)

var imageTextOnce sync.Once
var imageTextBytes []byte
var imageTextError error

// This resource is extracted from the exact hash-pinned public executable,
// not a device snapshot. Only __TEXT with no chained fixups is eligible.
func immutableImageText() ([]byte, error) {
	imageTextOnce.Do(func() {
		if sha256.Sum256(imageTextCompressed) != imageTextCompressedSHA {
			imageTextError = fmt.Errorf("immutable image resource mismatch")
			return
		}
		imageTextBytes, imageTextError = decodeImmutableText(imageTextCompressed, imageTextRawSHA, imageTextSize)
	})
	return imageTextBytes, imageTextError
}

func decodeImmutableText(encoded []byte, digest [32]byte, size uint64) ([]byte, error) {
	if size == 0 || size > 32<<20 || size%pageSize != 0 {
		return nil, fmt.Errorf("immutable image geometry rejected")
	}
	reader, err := gzip.NewReader(bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("immutable image compression rejected")
	}
	data, readErr := io.ReadAll(io.LimitReader(reader, int64(size)+1))
	closeErr := reader.Close()
	if readErr != nil || closeErr != nil || uint64(len(data)) != size || sha256.Sum256(data) != digest {
		clear(data)
		return nil, fmt.Errorf("immutable image content rejected")
	}
	return data, nil
}

// Preserve every captured page, including the six __TEXT pages that differ
// from the on-disk image. Missing writable data, heap and stack remain absent.
// Each inserted page is private to this recovery and can be cleared normally.
func installImmutableText(value *model, image []byte, address uint64) error {
	if value == nil || value.binaryHash != rawSourceBinaryHash || value.pages == nil ||
		address != 0x100000000 || len(image) == 0 || len(image) > 32<<20 || uint64(len(image))%pageSize != 0 {
		return fmt.Errorf("immutable image mapping rejected")
	}
	for offset := uint64(0); offset < uint64(len(image)); offset += pageSize {
		if value.pages[address+offset] == nil {
			value.pages[address+offset] = bytes.Clone(image[offset : offset+pageSize])
		}
	}
	return nil
}
