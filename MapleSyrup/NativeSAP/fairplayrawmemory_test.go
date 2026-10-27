package rawsource

import (
	"bytes"
	"testing"
)

// Exercise a sequence of actual external ABI calls, not only address arithmetic:
// large malloc -> stat/open/mmap/read of SIDB/SIDD/SINF-like private files ->
// inspect every mapped payload and the original live allocation after all calls.
// Public data only. This proves memory isolation, not FairPlay decryption.
func TestSourceMemoryFileTransactionPreservesAllLiveAllocations(t *testing.T) {
	s := newState(&model{pages: map[uint64][]byte{}})
	r := &rawSourceRuntime{nextAllocation: 0x710100000000, nextFD: 100,
		files: map[string][]byte{}, openFiles: map[uint64]*rawSourceOpenFile{}}
	call := func(pc uint64, args ...uint64) uint64 {
		t.Helper()
		for i, v := range args {
			s.x[i] = v
		}
		s.x[30] = 0x100001000
		if _, err := r.handle(s, pc); err != nil {
			t.Fatal(err)
		}
		return s.x[0]
	}
	large := call(0xF100000000007000, 3<<20)
	sentinel := bytes.Repeat([]byte{0xa5}, 3<<20)
	if err := s.store(large, sentinel); err != nil {
		t.Fatal(err)
	}
	type mapping struct {
		address uint64
		data    []byte
	}
	var mappings []mapping
	for _, name := range []string{"/state/SC Info.sidb", "/state/SC Info.sidd", "/app/App.sinf"} {
		data := bytes.Repeat([]byte(name), (2<<20)/len(name)+1)
		r.files[name] = data
		path := call(0xF100000000007000, uint64(len(name)+1))
		if err := s.store(path, append([]byte(name), 0)); err != nil {
			t.Fatal(err)
		}
		fd := call(0xF10000000001F000, path, 0)
		stat := call(0xF100000000007000, 144)
		if code := call(0xF10000000004D000, fd, stat); code != 0 {
			t.Fatal("stat failed")
		}
		length, err := s.load(stat+96, 8)
		if err != nil || length != uint64(len(data)) {
			t.Fatal("stat length")
		}
		// Model uses 4 KiB pages, while Darwin mmap copies the accessible
		// 16 KiB extent. A partial mapping of a larger file must reserve it.
		requested := uint64(1<<20) + 1
		accessible := (requested + 0x3fff) &^ uint64(0x3fff)
		mapped := call(0xF100000000023000, 0, requested, 1, 2, fd, 0)
		mappings = append(mappings, mapping{mapped, data[:accessible]})
		buffer := call(0xF100000000007000, length)
		if count := call(0xF100000000008000, fd, buffer, length); count != length {
			t.Fatal("read length")
		}
		mappings = append(mappings, mapping{buffer, data})
		call(0xF10000000001A000, fd)
	}
	mappings = append(mappings, mapping{large, sentinel})
	for _, m := range mappings {
		got, err := s.outputBytes(m.address, len(m.data))
		if err != nil || !bytes.Equal(got, m.data) {
			t.Fatal("live allocation corrupted by later allocation/file operation")
		}
	}
}
func TestSourceMemoryCollisionAndOverflowLeaveStateUnchanged(t *testing.T) {
	for _, address := range []uint64{0x710100000000, ^uint64(0) &^ (pageSize - 1)} {
		s := newState(&model{pages: map[uint64][]byte{address: make([]byte, pageSize)}})
		r := &rawSourceRuntime{nextAllocation: address}
		before := len(s.model.pages)
		if _, err := r.allocate(s, 2<<20); err == nil {
			t.Fatal("bad allocation accepted")
		}
		if r.nextAllocation != address || len(s.model.pages) != before {
			t.Fatal("failed allocation mutated state")
		}
	}
}

func TestSourceMemoryInvalidMappingCannotWrapIntoValidAllocation(t *testing.T) {
	for _, size := range []uint64{0, 0x10000001, ^uint64(0)} {
		s := newState(&model{pages: map[uint64][]byte{}})
		r := &rawSourceRuntime{nextAllocation: 0x710100000000, openFiles: map[uint64]*rawSourceOpenFile{1: {data: []byte{1}}}}
		s.x[1], s.x[4] = size, 1
		if _, err := r.handle(s, 0xF100000000023000); err == nil {
			t.Fatal("invalid mapping accepted")
		}
		if r.nextAllocation != 0x710100000000 || len(s.model.pages) != 0 {
			t.Fatal("invalid mapping changed memory")
		}
	}
}
