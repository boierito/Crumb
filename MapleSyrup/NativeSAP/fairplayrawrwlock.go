package rawsource

import (
	"encoding/binary"
	"fmt"
)

// A bounded implementation for one serialized guest CPU. Unsupported lock
// states are rejected instead of being reported as successful.
func (r *rawSourceRuntime) handleSourceRWLock(s *state, pc uint64) (uint64, error) {
	if r.rwlocks == nil {
		r.rwlocks = make(map[uint64]int)
	}
	ptr := s.x[0]
	reject := func() (uint64, error) { return 0, fmt.Errorf("unsupported bootstrap external %#x", pc) }
	if ptr == 0 || ptr%8 != 0 || ptr+200 < ptr {
		return reject()
	}
	data, err := s.outputBytes(ptr, 200)
	if err != nil {
		return 0, err
	}
	defer clear(data)
	if binary.LittleEndian.Uint32(data[8:12]) != 0 {
		return reject()
	}
	signature := binary.LittleEndian.Uint64(data)
	// Darwin rejects an uninitialized zero-signature lock with EINVAL. This
	// is an error return to the guest, not a successful lock operation.
	if signature == 0 && pc != 0xF10000000005C000 {
		return 22, nil
	}
	offset := 32
	if ptr%16 != 0 {
		offset = 40
	}
	initialize := func(attr uint64) (uint64, error) {
		if attr != 0 {
			a, e := s.load(attr, 8)
			if e != nil {
				return 0, e
			}
			shared, e := s.load(attr+8, 4)
			if e != nil {
				return 0, e
			}
			if a != 0x52574c41 || shared != 2 {
				return reject()
			}
		}
		if signature != 0 && signature != 0x2da8b3b4 {
			return reject()
		}
		for _, v := range data[8:] {
			if v != 0 {
				return reject()
			}
		}
		binary.LittleEndian.PutUint64(data[0:8], 0x52574c4b)
		binary.LittleEndian.PutUint32(data[12:16], 0x80000000)
		binary.LittleEndian.PutUint32(data[16:20], 0x20)
		binary.LittleEndian.PutUint32(data[20:24], 0x52574c4b)
		for i := 48; i < 64; i++ {
			data[i] = 0xff
		}
		if offset == 40 {
			for i := 32; i < 40; i++ {
				data[i] = 0xff
			}
		}
		binary.LittleEndian.PutUint32(data[offset:offset+4], 0xc0)
		binary.LittleEndian.PutUint32(data[offset+4:offset+8], 1)
		binary.LittleEndian.PutUint32(data[offset+8:offset+12], 0)
		binary.LittleEndian.PutUint32(data[offset+12:offset+16], 0)
		binary.LittleEndian.PutUint64(data[64:72], ^ptr)
		binary.LittleEndian.PutUint32(data[196:200], 0x52574c4b)
		if err := s.store(ptr, data); err != nil {
			return 0, err
		}
		signature = 0x52574c4b
		r.rwlocks[ptr] = 0
		return 0, nil
	}
	switch pc {
	case 0xF10000000005C000:
		return initialize(s.x[1])
	case 0xF100000000040000, 0xF100000000065000, 0xF10000000000E000:
		if signature == 0x2da8b3b4 && pc != 0xF10000000000E000 {
			if _, err := initialize(0); err != nil {
				return 0, err
			}
		}
		held, known := r.rwlocks[ptr]
		if !known || signature != 0x52574c4b {
			return reject()
		}
		l := binary.LittleEndian.Uint32(data[offset : offset+4])
		seq := binary.LittleEndian.Uint32(data[offset+4 : offset+8])
		u := binary.LittleEndian.Uint32(data[offset+8 : offset+12])
		outstanding := held
		if held == -1 {
			outstanding = 1
		}
		if outstanding < 0 || seq&255 != 1 || u&255 != 0 || l&0xffffff00 != u+uint32(outstanding)*256 || seq != (l&0xffffff00)|1 {
			return reject()
		}
		switch pc {
		case 0xF100000000040000:
			if held < 0 || (l&255 != 0x80 && l&255 != 0xc0) || l&0xffffff00 == 0xffffff00 {
				return reject()
			}
			l = ((l & 0xffffff00) + 256) | 0x80
			seq += 256
			held++
		case 0xF100000000065000:
			if held != 0 || l&255 != 0xc0 || l&0xffffff00 == 0xffffff00 {
				return reject()
			}
			// One emulated guest CPU has one stable, virtual thread identity.
			// It is synchronization metadata, never an account credential.
			if err := s.storeValue(ptr+24, 1, 8); err != nil {
				return 0, err
			}
			l = ((l & 0xffffff00) + 256) | 0x83
			seq += 256
			held = -1
		case 0xF10000000000E000:
			if held == 0 || u > 0xfffffeff {
				return reject()
			}
			if held == -1 {
				if l&255 != 0x83 {
					return reject()
				}
				owner, err := s.load(ptr+24, 8)
				if err != nil {
					return 0, err
				}
				if owner != 1 {
					return reject()
				}
				if err := s.storeValue(ptr+24, 0, 8); err != nil {
					return 0, err
				}
				held = 0
				l = (l & 0xffffff00) | 0xc0
			} else {
				if held < 0 || l&255 != 0x80 {
					return reject()
				}
				held--
				if held == 0 {
					l |= 0x40
				}
			}
			u += 256
		}
		binary.LittleEndian.PutUint32(data[offset:offset+4], l)
		binary.LittleEndian.PutUint32(data[offset+4:offset+8], seq)
		binary.LittleEndian.PutUint32(data[offset+8:offset+12], u)
		if err := s.store(ptr+uint64(offset), data[offset:offset+12]); err != nil {
			return 0, err
		}
		r.rwlocks[ptr] = held
		return 0, nil
	case 0xF10000000004C000:
		held, known := r.rwlocks[ptr]
		if signature == 0x2da8b3b4 && !known {
			if err := s.storeValue(ptr, 0, 8); err != nil {
				return 0, err
			}
			return 0, nil
		}
		if !known || signature != 0x52574c4b {
			return reject()
		}
		if held != 0 {
			return 16, nil
		}
		if err := s.storeValue(ptr, 0, 8); err != nil {
			return 0, err
		}
		delete(r.rwlocks, ptr)
		return 0, nil
	default:
		return reject()
	}
}
