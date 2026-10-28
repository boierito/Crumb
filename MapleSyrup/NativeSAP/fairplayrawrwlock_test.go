package rawsource

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestNativeReferenceLifecycle(t *testing.T) {
	const ptr = uint64(0x600000)
	states := []string{
		"4b4c5752000000000000000000000080200000004b4c57520000000000000000c0000000010000000000000000000000ffffffffffffffffffffffffffffffff",
		"4b4c5752000000000000000000000080200000004b4c5752000000000000000080010000010100000000000000000000ffffffffffffffffffffffffffffffff",
		"4b4c5752000000000000000000000080200000004b4c57520000000000000000c0010000010100000001000000000000ffffffffffffffffffffffffffffffff",
		"00000000000000000000000000000080200000004b4c57520000000000000000c0010000010100000001000000000000ffffffffffffffffffffffffffffffff",
	}
	m := &model{pages: map[uint64][]byte{ptr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	runtime := &rawSourceRuntime{rwlocks: map[uint64]int{}}
	calls := []uint64{0xf10000000005c000, 0xf100000000040000, 0xf10000000000e000, 0xf10000000004c000}
	for i, pc := range calls {
		s.x[0], s.x[1] = ptr, 0
		status, err := runtime.handleSourceRWLock(s, pc)
		if err != nil || status != 0 {
			t.Fatalf("stage %d failed: %v", i, err)
		}
		got, err := s.outputBytes(ptr, 64)
		want, _ := hex.DecodeString(states[i])
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("stage %d differs from native", i)
		}
	}
}
func TestUnmodeledLockInputsFailClosed(t *testing.T) {
	const ptr = uint64(0x600000)
	m := &model{pages: map[uint64][]byte{ptr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	runtime := &rawSourceRuntime{rwlocks: map[uint64]int{}}
	s.x[0], s.x[1] = ptr, 1
	if _, err := runtime.handleSourceRWLock(s, 0xf10000000005c000); err == nil {
		t.Fatal("unsupported attr accepted")
	}
	if code, err := runtime.handleSourceRWLock(s, 0xf100000000065000); err != nil || code != 22 {
		t.Fatal("uninitialized writer must return Darwin EINVAL")
	}
	b, _ := s.outputBytes(ptr, 200)
	if !bytes.Equal(b, make([]byte, 200)) {
		t.Fatal("rejection mutated guest memory")
	}
}
func TestRepeatedReadersMatchNativeAndBusyDestroy(t *testing.T) {
	const ptr = uint64(0x600000)
	m := &model{pages: map[uint64][]byte{ptr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	runtime := &rawSourceRuntime{rwlocks: map[uint64]int{}}
	init := uint64(0xf10000000005c000)
	rd := uint64(0xf100000000040000)
	un := uint64(0xf10000000000e000)
	destroy := uint64(0xf10000000004c000)
	s.x[0], s.x[1] = ptr, 0
	if _, e := runtime.handleSourceRWLock(s, init); e != nil {
		t.Fatal(e)
	}
	calls := []uint64{rd, rd, un, un, rd, un}
	wants := []string{"80010000010100000000000000000000", "80020000010200000000000000000000", "80020000010200000001000000000000", "c0020000010200000002000000000000", "80030000010300000002000000000000", "c0030000010300000003000000000000"}
	for i, call := range calls {
		s.x[0] = ptr
		v, e := runtime.handleSourceRWLock(s, call)
		if e != nil || v != 0 {
			t.Fatalf("stage %d: %v", i, e)
		}
		got, e := s.outputBytes(ptr+32, 16)
		want, _ := hex.DecodeString(wants[i])
		if e != nil || !bytes.Equal(got, want) {
			t.Fatalf("stage %d disagrees with native", i)
		}
		if i == 1 {
			s.x[0] = ptr
			v, e = runtime.handleSourceRWLock(s, destroy)
			if e != nil || v != 16 {
				t.Fatal("busy lock destroyed")
			}
		}
	}
	s.x[0] = ptr
	v, e := runtime.handleSourceRWLock(s, destroy)
	if e != nil || v != 0 {
		t.Fatal("unlocked destroy failed")
	}
}
func TestMisalignedDefaultAttributesMatchNative(t *testing.T) {
	const ptr = uint64(0x600008)
	const attr = uint64(0x601000)
	m := &model{pages: map[uint64][]byte{0x600000: make([]byte, 4096), attr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	runtime := &rawSourceRuntime{rwlocks: map[uint64]int{}}
	if e := s.storeValue(attr, 0x52574c41, 8); e != nil {
		t.Fatal(e)
	}
	if e := s.storeValue(attr+8, 2, 4); e != nil {
		t.Fatal(e)
	}
	calls := []uint64{0xf10000000005c000, 0xf100000000040000, 0xf10000000000e000, 0xf10000000004c000}
	wants := []string{"4b4c5752000000000000000000000080200000004b4c57520000000000000000ffffffffffffffffc0000000010000000000000000000000ffffffffffffffff", "4b4c5752000000000000000000000080200000004b4c57520000000000000000ffffffffffffffff80010000010100000000000000000000ffffffffffffffff", "4b4c5752000000000000000000000080200000004b4c57520000000000000000ffffffffffffffffc0010000010100000001000000000000ffffffffffffffff", "00000000000000000000000000000080200000004b4c57520000000000000000ffffffffffffffffc0010000010100000001000000000000ffffffffffffffff"}
	for i, pc := range calls {
		s.x[0] = ptr
		s.x[1] = attr
		v, e := runtime.handleSourceRWLock(s, pc)
		if e != nil || v != 0 {
			t.Fatalf("stage %d: %v", i, e)
		}
		got, e := s.outputBytes(ptr, 64)
		want, _ := hex.DecodeString(wants[i])
		if e != nil || !bytes.Equal(got, want) {
			t.Fatalf("stage %d mismatched native", i)
		}
	}
}
func TestStaticInitializerLazyRead(t *testing.T) {
	const ptr = uint64(0x600000)
	m := &model{pages: map[uint64][]byte{ptr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	runtime := &rawSourceRuntime{rwlocks: map[uint64]int{}}
	if e := s.storeValue(ptr, 0x2da8b3b4, 8); e != nil {
		t.Fatal(e)
	}
	s.x[0], s.x[1] = ptr, 0xdeadbeef
	if v, e := runtime.handleSourceRWLock(s, 0xf100000000040000); e != nil || v != 0 {
		t.Fatal("static read failed", e)
	}
	if v, e := runtime.handleSourceRWLock(s, 0xf10000000004c000); e != nil || v != 16 {
		t.Fatal("busy static lock was destroyed")
	}
}

func TestNativeRepeatedWriterTransitions(t *testing.T) {
	const ptr = uint64(0x600000)
	m := &model{pages: map[uint64][]byte{ptr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	r := &rawSourceRuntime{}
	s.x[0], s.x[1] = ptr, 0
	if _, e := r.handleSourceRWLock(s, 0xf10000000005c000); e != nil {
		t.Fatal(e)
	}
	calls := []uint64{0xf100000000065000, 0xf10000000000e000, 0xf100000000065000, 0xf10000000000e000, 0xf100000000065000, 0xf10000000000e000}
	want := []string{"83010000010100000000000000000000", "c0010000010100000001000000000000", "83020000010200000001000000000000", "c0020000010200000002000000000000", "83030000010300000002000000000000", "c0030000010300000003000000000000"}
	for i, pc := range calls {
		s.x[0] = ptr
		v, e := r.handleSourceRWLock(s, pc)
		if e != nil || v != 0 {
			t.Fatalf("stage %d: %v", i, e)
		}
		got, _ := s.outputBytes(ptr+32, 16)
		w, _ := hex.DecodeString(want[i])
		if !bytes.Equal(got, w) {
			t.Fatalf("stage %d differs from native", i)
		}
		if i%2 == 0 {
			s.x[0] = ptr
			v, e = r.handleSourceRWLock(s, 0xf10000000004c000)
			if e != nil || v != 16 {
				t.Fatal("writer must make destroy busy")
			}
		}
	}
	s.x[0] = ptr
	if v, e := r.handleSourceRWLock(s, 0xf10000000004c000); e != nil || v != 0 {
		t.Fatal("destroy rejected")
	}
}

func TestZeroSignatureReturnsNativeErrorWithoutMutation(t *testing.T) {
	const ptr = uint64(0x600000)
	m := &model{pages: map[uint64][]byte{ptr: make([]byte, 4096)}}
	s := newState(m)
	s.reset()
	r := &rawSourceRuntime{}
	for _, pc := range []uint64{0xf100000000065000, 0xf100000000040000, 0xf10000000000e000, 0xf10000000004c000} {
		s.x[0] = ptr
		code, err := r.handleSourceRWLock(s, pc)
		if err != nil || code != 22 {
			t.Fatal("expected Darwin EINVAL")
		}
	}
	got, e := s.outputBytes(ptr, 200)
	if e != nil || !bytes.Equal(got, make([]byte, 200)) {
		t.Fatal("invalid lock mutated")
	}
}
