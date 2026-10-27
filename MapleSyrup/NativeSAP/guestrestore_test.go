package machine

import (
	"bytes"
	"context"
	"github.com/majd/ipatool/v2/internal/sap/unicorn"
	"testing"
)

func TestGuestStateRestoreOwnershipAndBounds(t *testing.T) {
	input := map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}, "SC Info.sidb.lock": {}}
	saved, err := cloneGuestState(input)
	if err != nil {
		t.Fatal(err)
	}
	copyState := saved[storeAgentSCInfoPath+"/SC Info.sidb"]
	copyState[0] ^= 1
	if input["SC Info.sidb"][0] != 0 {
		t.Fatal("restored state aliases capture")
	}
	clear(copyState)
	if !bytes.Equal(input["SC Info.sidb"], []byte{0, 7, 0, 11}) {
		t.Fatal("capture changed")
	}
	for _, files := range []map[string][]byte{nil, {}, {"SC Info.sidb": {}}, {"SC Info.sidb": {0, 7, 0, 11}, "../escaped": {}}, {"SC Info.sidb": {0, 7, 0, 11}, "nested/file": {}}, {"SC Info.sidb": {0, 7, 0, 11}, "bad\\name": {}}, {"SC Info.sidb": {0, 7, 0, 11}, "bad\x00name": {}}, {"SC Info.sidb": make([]byte, (4<<20)+1)}} {
		if result, err := cloneGuestState(files); err == nil || result != nil {
			t.Fatal("invalid filesystem restored")
		}
	}
}

func TestGuestStateRestoreRejectsReinitialization(t *testing.T) {
	s := &shims{guestStateFS: &guestStateFS{files: map[string][]byte{"existing": {1}}}}
	if err := s.installGuestState(map[string][]byte{}); err == nil {
		t.Fatal("existing state overwritten")
	}
	if s.guestStateFS.files["existing"][0] != 1 {
		t.Fatal("existing state changed")
	}
}

func TestGuestStateClosedProducerRejectsNativeCalls(t *testing.T) {
	session := &RetainedKBSyncSession{hardware: []byte{2, 0, 0, 0, 0, 1}}
	if _, err := session.NeedsAccountDiversity(1); err == nil {
		t.Fatal("closed producer invoked")
	}
	if _, _, _, err := session.ImportKeybag([]byte{1}); err == nil {
		t.Fatal("closed producer imported")
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestGuestStateUnlinkIsolation(t *testing.T) {
	engine, err := unicorn.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	s, err := newShims(engine, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.enableGuestState(); err != nil {
		t.Fatal(err)
	}
	const buffer = uint64(0x70000000)
	if err = engine.MemMap(buffer, 4096); err != nil {
		t.Fatal(err)
	}
	original := []byte("public saved-state fixture")
	s.rfs().files[storeAgentSCInfoPath+"/SC Info.sidb"] = append([]byte(nil), original...)
	unlink := func(name string) uint64 {
		t.Helper()
		if err := engine.MemWrite(buffer, append([]byte(name), 0)); err != nil {
			t.Fatal(err)
		}
		if err := engine.RegWrite(unicorn.RegRDI, buffer); err != nil {
			t.Fatal(err)
		}
		if err := s.fsUnlink(); err != nil {
			t.Fatal(err)
		}
		result, err := engine.RegRead(unicorn.RegRAX)
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	if unlink("/tmp/outside") != ^uint64(0) {
		t.Fatal("host path accepted")
	}
	path := storeAgentSCInfoPath + "/SC Info.sidb"
	s.rfs().open[16] = &guestStateFile{name: path}
	if unlink(path) != ^uint64(0) {
		t.Fatal("unsupported open inode deletion accepted")
	}
	delete(s.rfs().open, 16)
	if unlink(path) != 0 {
		t.Fatal("closed virtual file deletion failed")
	}
	if _, ok := s.rfs().files[path]; ok {
		t.Fatal("virtual file remains")
	}
	if string(original) != "public saved-state fixture" {
		t.Fatal("original modified")
	}
	if unlink(path) != ^uint64(0) {
		t.Fatal("missing file deletion accepted")
	}
}
