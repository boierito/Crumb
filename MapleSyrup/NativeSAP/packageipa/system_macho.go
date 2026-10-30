package packageipa

// The system path never derives a key or changes cryptid to manufacture success.
// These parsers verify the native output and remove embedded signing identities.
import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type systemRange struct{ off, size int64 }
type systemCrypt struct {
	systemRange
	field int64
	id    uint32
}
type systemSlice struct {
	off, size  int64
	header     []byte
	headerSize int
	crypt      []systemCrypt
	signatures []systemRange
}
type systemImage struct{ slices []systemSlice }

func systemRead(r io.ReaderAt, off, n, size int64) ([]byte, error) {
	if off < 0 || n < 0 || n > 1<<20 || off > size-n {
		return nil, errors.New("invalid Mach-O extent")
	}
	b := make([]byte, n)
	_, err := r.ReadAt(b, off)
	return b, err
}

func systemMachO(r io.ReaderAt, size int64) (*systemImage, error) {
	if size < 4 {
		return nil, nil
	}
	b, e := systemRead(r, 0, 4, size)
	if e != nil {
		return nil, e
	}
	little, big := binary.LittleEndian, binary.BigEndian
	magic := little.Uint32(b)
	var spans []systemRange
	switch magic {
	case 0xfeedface, 0xfeedfacf:
		spans = []systemRange{{0, size}}
	case 0xbebafeca, 0xbfbafeca:
		h, e := systemRead(r, 0, 8, size)
		if e != nil {
			return nil, e
		}
		n := big.Uint32(h[4:])
		stride := int64(20)
		if magic == 0xbfbafeca {
			stride = 32
		}
		if n == 0 || n > 32 {
			return nil, errors.New("unsupported FAT slice count")
		}
		table, e := systemRead(r, 8, int64(n)*stride, size)
		if e != nil {
			return nil, e
		}
		for i := int64(0); i < int64(n); i++ {
			q := table[i*stride:]
			var off, length uint64
			if stride == 20 {
				off = uint64(big.Uint32(q[8:]))
				length = uint64(big.Uint32(q[12:]))
			} else {
				off = big.Uint64(q[8:])
				length = big.Uint64(q[16:])
			}
			if off < uint64(8+int64(n)*stride) || length < 28 || off > uint64(size) || length > uint64(size)-off {
				return nil, errors.New("invalid FAT extent")
			}
			for _, old := range spans {
				if int64(off) < old.off+old.size && old.off < int64(off+length) {
					return nil, errors.New("overlapping FAT slices")
				}
			}
			spans = append(spans, systemRange{int64(off), int64(length)})
		}
	case 0xcefaedfe, 0xcffaedfe, 0xcafebabe, 0xcafebabf:
		return nil, errors.New("unsupported Mach-O byte order")
	default:
		return nil, nil
	}
	image := &systemImage{}
	for _, span := range spans {
		h, e := systemRead(r, span.off, 28, size)
		if e != nil {
			return nil, e
		}
		m := little.Uint32(h)
		hs := 28
		if m == 0xfeedfacf {
			hs = 32
		} else if m != 0xfeedface {
			return nil, errors.New("unsupported FAT member")
		}
		n, commands := little.Uint32(h[16:]), little.Uint32(h[20:])
		if n > 4096 || int64(hs)+int64(commands) > span.size {
			return nil, errors.New("invalid load command count")
		}
		header, e := systemRead(r, span.off, int64(hs)+int64(commands), size)
		if e != nil {
			return nil, e
		}
		s := systemSlice{off: span.off, size: span.size, header: header, headerSize: hs}
		pos := hs
		for j := uint32(0); j < n; j++ {
			if pos+8 > len(header) {
				return nil, errors.New("truncated load command")
			}
			cmd, cs := little.Uint32(header[pos:]), int(little.Uint32(header[pos+4:]))
			if cs < 8 || cs%4 != 0 || cs > len(header)-pos {
				return nil, errors.New("invalid load command size")
			}
			switch cmd {
			case 0x21, 0x2c:
				want := 20
				if cmd == 0x2c {
					want = 24
				}
				if cs != want || len(s.crypt) != 0 {
					return nil, errors.New("invalid encryption command")
				}
				o, l, id := int64(little.Uint32(header[pos+8:])), int64(little.Uint32(header[pos+12:])), little.Uint32(header[pos+16:])
				if o > s.size || l > s.size-o || id != 0 && (l == 0 || o < int64(len(header))) {
					return nil, errors.New("invalid encryption range")
				}
				s.crypt = append(s.crypt, systemCrypt{systemRange{s.off + o, l}, s.off + int64(pos+16), id})
			case 0x1d:
				if cs != 16 || len(s.signatures) != 0 {
					return nil, errors.New("invalid code signature command")
				}
				o, l := int64(little.Uint32(header[pos+8:])), int64(little.Uint32(header[pos+12:]))
				if o > s.size || l > s.size-o || l != 0 && o < int64(len(header)) {
					return nil, errors.New("invalid code signature extent")
				}
				s.signatures = append(s.signatures, systemRange{s.off + o, l})
			case 0x19:
				if cs < 72 {
					return nil, errors.New("invalid segment command")
				}
				o, l := little.Uint64(header[pos+40:]), little.Uint64(header[pos+48:])
				if o > uint64(s.size) || l > uint64(s.size)-o {
					return nil, errors.New("segment exceeds slice")
				}
			}
			pos += cs
		}
		if pos != len(header) {
			return nil, errors.New("unconsumed load commands")
		}
		for _, sig := range s.signatures {
			for _, c := range s.crypt {
				if c.id != 0 && sig.size != 0 && sig.off < c.off+c.size && c.off < sig.off+sig.size {
					return nil, errors.New("signature overlaps ciphertext")
				}
			}
		}
		image.slices = append(image.slices, s)
	}
	return image, nil
}

// A patch reader preserves file sizes and segment offsets. Removing a command
// compacts only the load-command area and zeros the vacated command bytes.
type systemPatch struct {
	systemRange
	data []byte
}
type systemPatched struct {
	io.ReaderAt
	patches []systemPatch
}

func (r systemPatched) ReadAt(p []byte, off int64) (int, error) {
	n, e := r.ReaderAt.ReadAt(p, off)
	for _, patch := range r.patches {
		lo, hi := max(off, patch.off), min(off+int64(n), patch.off+patch.size)
		if lo >= hi {
			continue
		}
		target := p[lo-off : hi-off]
		if patch.data == nil {
			clear(target)
		} else {
			copy(target, patch.data[lo-patch.off:hi-patch.off])
		}
	}
	return n, e
}
func systemRemoveSignatures(r io.ReaderAt, image *systemImage) (io.ReaderAt, int) {
	patched := systemPatched{ReaderAt: r}
	count := 0
	for _, s := range image.slices {
		if len(s.signatures) == 0 {
			continue
		}
		h := make([]byte, len(s.header))
		copy(h, s.header[:s.headerSize])
		read, write, kept := s.headerSize, s.headerSize, uint32(0)
		for read < len(s.header) {
			cmd, n := binary.LittleEndian.Uint32(s.header[read:]), int(binary.LittleEndian.Uint32(s.header[read+4:]))
			if cmd != 0x1d {
				copy(h[write:], s.header[read:read+n])
				write += n
				kept++
			} else {
				count++
			}
			read += n
		}
		binary.LittleEndian.PutUint32(h[16:], kept)
		binary.LittleEndian.PutUint32(h[20:], uint32(write-s.headerSize))
		patched.patches = append(patched.patches, systemPatch{systemRange{s.off, int64(len(h))}, h})
		for _, sig := range s.signatures {
			patched.patches = append(patched.patches, systemPatch{sig, nil})
		}
	}
	return patched, count
}

func systemEncrypted(image *systemImage) int {
	count := 0
	for _, s := range image.slices {
		for _, c := range s.crypt {
			if c.id != 0 {
				count++
			}
		}
	}
	return count
}

func systemVerifyNative(original, candidate io.ReaderAt, size int64, image *systemImage) (int, error) {
	other, e := systemMachO(candidate, size)
	if e != nil || other == nil {
		return 0, errors.New("invalid native Mach-O output")
	}
	if systemEncrypted(other) != 0 {
		return 0, errors.New("native output is still encrypted")
	}
	// A native candidate must preserve every byte outside original crypt ranges
	// and their cryptid fields. The system reader writes no signature changes.
	var allowed []systemRange
	pages := 0
	for _, s := range image.slices {
		for _, c := range s.crypt {
			if c.id != 0 {
				allowed = append(allowed, c.systemRange, systemRange{c.field, 4})
				for off := c.off; off < c.off+c.size; off += 16384 {
					n := min(int64(16384), c.off+c.size-off)
					a := make([]byte, n)
					b := make([]byte, n)
					if _, e := original.ReadAt(a, off); e != nil {
						return 0, e
					}
					if _, e := candidate.ReadAt(b, off); e != nil {
						return 0, e
					}
					if bytes.Equal(a, b) {
						return 0, errors.New("native page did not change")
					}
					pages++
				}
			}
		}
	}
	a, b := make([]byte, 64<<10), make([]byte, 64<<10)
	for off := int64(0); off < size; off += int64(len(a)) {
		n := int(min(int64(len(a)), size-off))
		if _, e := original.ReadAt(a[:n], off); e != nil {
			return 0, e
		}
		if _, e := candidate.ReadAt(b[:n], off); e != nil {
			return 0, e
		}
		for _, span := range allowed {
			lo, hi := max(off, span.off), min(off+int64(n), span.off+span.size)
			if lo < hi {
				clear(a[lo-off : hi-off])
				clear(b[lo-off : hi-off])
			}
		}
		if !bytes.Equal(a[:n], b[:n]) {
			return 0, errors.New("native output changed unprotected bytes")
		}
	}
	return pages, nil
}
