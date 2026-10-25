package machine

import (
	"context"
	"github.com/majd/ipatool/v2/internal/sap/unicorn"
	"testing"
)

func TestGuestStateIsolationAndIO(t *testing.T) {
	engine, err := unicorn.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	s, err := newShims(engine, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if s.guestStateFS != nil {
		t.Fatal("virtual state must be opt in")
	}
	if err = s.enableGuestState(); err != nil {
		t.Fatal(err)
	}
	const buffer = uint64(0x70000000)
	if err = engine.MemMap(buffer, 4096); err != nil {
		t.Fatal(err)
	}
	call := func(fn func() error, args ...uint64) uint64 {
		t.Helper()
		regs := []int{unicorn.RegRDI, unicorn.RegRSI, unicorn.RegRDX}
		for i, v := range args {
			if err := engine.RegWrite(regs[i], v); err != nil {
				t.Fatal(err)
			}
		}
		if err := fn(); err != nil {
			t.Fatal(err)
		}
		v, err := engine.RegRead(unicorn.RegRAX)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	if err = s.fsOpen("/tmp/outside", 0x202); err != nil {
		t.Fatal(err)
	}
	v, _ := engine.RegRead(unicorn.RegRAX)
	if v != ^uint64(0) {
		t.Fatal("outside path accepted")
	}
	if err = s.fsOpen(storeAgentSCInfoPath+"/SC Info.sidb", 0x202); err != nil {
		t.Fatal(err)
	}
	fd, _ := engine.RegRead(unicorn.RegRAX)
	engine.MemWrite(buffer, []byte("public-test-data"))
	if call(s.fsWrite, fd, buffer, 16) != 16 {
		t.Fatal("write failed")
	}
	if call(s.fsSeek, fd, 0, 0) != 0 {
		t.Fatal("seek failed")
	}
	if err = s.fsRead(fd, buffer+128, 16); err != nil {
		t.Fatal(err)
	}
	data, _ := engine.MemRead(buffer+128, 16)
	if string(data) != "public-test-data" {
		t.Fatal("read mismatch")
	}
	if call(s.fsWrite, fd, buffer, 1<<21) != ^uint64(0) {
		t.Fatal("oversized write accepted")
	}
	if call(s.fsLock, fd) != 0 {
		t.Fatal("single-guest lock failed")
	}
	if call(s.fsClose, fd) != 0 {
		t.Fatal("close failed")
	}
	if call(s.fsWrite, fd, buffer, 1) != ^uint64(0) {
		t.Fatal("closed descriptor accepted")
	}
	snapshot := s.guestStateSnapshot()
	snapshot["SC Info.sidb"][0] = 0
	if s.rfs().files[storeAgentSCInfoPath+"/SC Info.sidb"][0] == 0 {
		t.Fatal("snapshot aliases state")
	}
	for _, b := range snapshot {
		clear(b)
	}
}
func TestGuestStatePathBoundary(t *testing.T) {
	for _, p := range []string{"/tmp/x", storeAgentSCInfoPath + "/../outside", storeAgentSCInfoPath + "/nested/file", storeAgentSCInfoPath + "/bad\\path"} {
		if guestStatePath(p) {
			t.Fatal("path escaped")
		}
	}
}
