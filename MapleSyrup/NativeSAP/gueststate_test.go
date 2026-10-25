package machine

import (
	"context"
	"encoding/binary"
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

func TestGuestStateDarwinPathAndCatalogContract(t *testing.T) {
	engine, e := unicorn.New(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	defer engine.Close()
	s, e := newShims(engine, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.enableGuestState(); e != nil {
		t.Fatal(e)
	}
	const base = uint64(0x71000000)
	if e = engine.MemMap(base, 4096); e != nil {
		t.Fatal(e)
	}
	call := func(fn func() error, args ...uint64) uint64 {
		t.Helper()
		regs := []int{unicorn.RegRDI, unicorn.RegRSI, unicorn.RegRDX, unicorn.RegRCX, unicorn.RegR8, unicorn.RegR9}
		for i, v := range args {
			if e := engine.RegWrite(regs[i], v); e != nil {
				t.Fatal(e)
			}
		}
		if e := fn(); e != nil {
			t.Fatal(e)
		}
		v, e := engine.RegRead(unicorn.RegRAX)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	name := storeAgentSCInfoPath + "/public-fixture.sidb"
	if e = s.fsOpen(name, 0x202); e != nil {
		t.Fatal(e)
	}
	fd, _ := engine.RegRead(unicorn.RegRAX)
	if call(s.fsFcntl, fd, 50, base+256) != 0 {
		t.Fatal("F_GETPATH failed")
	}
	pathData, _ := engine.MemRead(base+256, uint64(len(name)+1))
	if string(pathData) != name+"\x00" {
		t.Fatal("F_GETPATH returned incorrect path")
	}
	if call(s.fsFcntl, fd, 2, 1) != 0 || call(s.fsFcntl, fd, 1, 0) != 1 {
		t.Fatal("descriptor flags not retained")
	}
	engine.MemWrite(base, append([]byte(name), 0))
	if call(s.fsLstat, base, base+512) != 0 {
		t.Fatal("lstat rejected an existing virtual file")
	}
	stat, _ := engine.MemRead(base+512, 144)
	if binary.LittleEndian.Uint16(stat[4:]) != 0100600 || binary.LittleEndian.Uint16(stat[6:]) != 1 {
		t.Fatal("Darwin stat layout")
	}
	if call(s.fsPathMakeRef, base, base+768, base+1000) != 0 {
		t.Fatal("FSPathMakeRef failed")
	}
	if call(s.fsRefMakePath, base+768, base+1024, 512) != 0 {
		t.Fatal("FSRefMakePath failed")
	}
	roundtrip, _ := engine.MemRead(base+1024, uint64(len(name)+1))
	if string(roundtrip) != name+"\x00" {
		t.Fatal("FSRef round trip mismatch")
	}
	if call(s.fsGetCatalogInfo, base+768, 0x410, base+1600, 0, 0, 0) != 0 {
		t.Fatal("catalog read failed")
	}
	catalog, _ := engine.MemRead(base+1600, 148)
	if binary.LittleEndian.Uint16(catalog[66:]) != 0100600 || uint64(binary.LittleEndian.Uint32(catalog[8:])) != binary.LittleEndian.Uint64(stat[8:]) {
		t.Fatal("LP64 packed catalog metadata disagrees with stat")
	}
	binary.LittleEndian.PutUint16(catalog[66:], 0100640)
	engine.MemWrite(base+1600, catalog)
	if call(s.fsSetCatalogInfo, base+768, 0x400, base+1600) != 0 {
		t.Fatal("catalog permission update failed")
	}
	if call(s.fsLstat, base, base+512) != 0 {
		t.Fatal("lstat failed")
	}
	stat, _ = engine.MemRead(base+512, 144)
	if binary.LittleEndian.Uint16(stat[4:]) != 0100640 {
		t.Fatal("permission change not reflected by lstat")
	}
	if call(s.fsGetCatalogInfo, base+768, 0x400000, base+1600, 0, 0, 0) == 0 {
		t.Fatal("unsupported security object was silently accepted")
	}
	// An unsupported *write* mask must fail closed while exposing only its
	// safe catalog flags/reason to the private virtual-I/O diagnostic.
	setCatalog := s.traceGuestOperation("_FSSetCatalogInfo", s.fsSetCatalogInfo)
	if call(setCatalog, base+768, 0x410, base+1600) != uint64(0xffffffce) {
		t.Fatal("unsupported catalog write flags were silently accepted")
	}
	trace := s.guestStateTrace()
	events := trace["recent"].([]guestStateEvent)
	last := events[len(events)-1]
	if last.Operation != "_FSSetCatalogInfo" || last.Outcome != "negative_status" ||
		last.Detail != "unsupported_catalog_mask" || last.Mask != 0x410 {
		t.Fatalf("unexpected sanitized catalog trace: %+v", last)
	}
	if call(setCatalog, base+768, 0x400, 0) != uint64(0xffffffce) {
		t.Fatal("missing catalog buffer accepted")
	}
	events = s.guestStateTrace()["recent"].([]guestStateEvent)
	last = events[len(events)-1]
	if last.Detail != "missing_catalog_info" || last.Mask != 0x400 {
		t.Fatalf("missing catalog buffer reason was not retained: %+v", last)
	}
	if call(s.fsFcntl, 999, 50, base+256) != ^uint64(0) {
		t.Fatal("unknown descriptor accepted")
	}
}


// The private trace may record error classes, but never guest paths, pointers,
// file content, account identity, or a complete unbounded syscall history.
func TestGuestStateTraceIsBoundedAndPathFree(t *testing.T) {
	engine, err := unicorn.New(context.Background())
	if err != nil { t.Fatal(err) }
	defer engine.Close()
	s, err := newShims(engine, nil, nil)
	if err != nil { t.Fatal(err) }
	const privatePathCanary = "/Users/Shared/private-account-material"
	if err = s.fsOpen(privatePathCanary, 0x202); err != nil { t.Fatal(err) }
	traceRead := s.traceGuestOperation("_read", func() error { return s.setResult(8192) })
	for i := 0; i < guestStateTraceLimit+5; i++ {
		if err = traceRead(); err != nil { t.Fatal(err) }
	}
	traceMissing := s.traceGuestOperation("_lstat", func() error { return s.fsError(2) })
	if err = traceMissing(); err != nil { t.Fatal(err) }
	trace := s.guestStateTrace()
	if trace["total_calls"] != uint64(guestStateTraceLimit+6) ||
		trace["truncated"] != true { t.Fatal("trace call count or ring limit") }
	events, ok := trace["recent"].([]guestStateEvent)
	if !ok || len(events) != guestStateTraceLimit { t.Fatal("trace not bounded") }
	last := events[len(events)-1]
	if last.Operation != "_lstat" || last.Outcome != "errno" || last.Errno != 2 {
		t.Fatalf("incorrect sanitized trace result: %+v", last)
	}
	for _, event := range events {
		if event.Operation != "_read" && event.Operation != "_lstat" {
			t.Fatal("trace contains unexpected name")
		}
	}
}
