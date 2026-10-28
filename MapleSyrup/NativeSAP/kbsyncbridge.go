package main

/*
#include <stdint.h>
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"github.com/majd/ipatool/v2/internal/sap/machine"
	"os"
	"time"
	"unsafe"
)

type retainedKBSync struct {
	agent            *machine.RetainedKBSyncSession
	blob, keybag     []byte
	accountDirectory string
	accountRestored  bool
	hardwareID       []byte
	accountID        uint64
	expectedDigests  map[string]string
}

var retainedSessions = map[uint64]*retainedKBSync{}

//export WaffleKBSyncOpen
func WaffleKBSyncOpen(cache *C.char, hardware *C.uchar, length C.size_t, dsid C.uint64_t, handle *C.uint64_t, output **C.uchar, size *C.size_t) C.int {
	return openKBSyncWithState(cache, nil, hardware, length, dsid, handle, output, size)
}

//export WaffleKBSyncOpenWithState
func WaffleKBSyncOpenWithState(cache, directory *C.char, hardware *C.uchar, length C.size_t, dsid C.uint64_t, handle *C.uint64_t, output **C.uchar, size *C.size_t) C.int {
	if directory == nil || C.GoString(directory) == "" {
		return 1
	}
	return openKBSyncWithState(cache, directory, hardware, length, dsid, handle, output, size)
}

func openKBSyncWithState(cache, directory *C.char, hardware *C.uchar, length C.size_t, dsid C.uint64_t, handle *C.uint64_t, output **C.uchar, size *C.size_t) C.int {
	mutex.Lock()
	defer mutex.Unlock()
	if cache == nil || hardware == nil || length != 6 || dsid == 0 || handle == nil || output == nil || size == nil {
		return 1
	}
	assets.CacheRoot = C.GoString(cache)
	if assets.CacheRoot == "" {
		return 1
	}
	hardwareID := C.GoBytes(unsafe.Pointer(hardware), 6)
	defer clear(hardwareID)
	accountDirectory := ""
	var stored authorizedAccountState
	restored := false
	if directory != nil {
		accountDirectory = C.GoString(directory)
		var err error
		stored, err = readAuthorizedState(accountDirectory, uint64(dsid), hardwareID)
		if err == nil {
			restored = true
			defer stored.clear()
		} else if !os.IsNotExist(err) {
			return 12
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	bundle, err := assets.Load(ctx)
	if err != nil {
		return 2
	}
	var agent *machine.RetainedKBSyncSession
	var blob []byte
	if restored {
		agent, blob, err = machine.OpenRetainedAuthorizedKBSync(ctx, bundle, hardwareID, uint64(dsid), stored.Files)
	} else {
		agent, blob, err = machine.OpenRetainedKBSync(ctx, bundle, hardwareID, uint64(dsid))
	}
	if err != nil {
		return 8
	}
	nextHandle++
	if nextHandle == 0 {
		clear(blob)
		agent.Close()
		return 7
	}
	var expected map[string]string
	if restored {
		expected = make(map[string]string, len(stored.Digests))
		for name, digest := range stored.Digests {
			expected[name] = digest
		}
	}
	retainedSessions[nextHandle] = &retainedKBSync{agent: agent, blob: blob, accountDirectory: accountDirectory, accountRestored: restored, hardwareID: append([]byte(nil), hardwareID...), accountID: uint64(dsid), expectedDigests: expected}
	*handle = C.uint64_t(nextHandle)
	*output = (*C.uchar)(C.CBytes(blob))
	*size = C.size_t(len(blob))
	return 0
}

//export WaffleKBSyncBind
func WaffleKBSyncBind(handle C.uint64_t, blob *C.uchar, size C.size_t, keybag *C.uchar, keybagSize C.size_t) C.int {
	mutex.Lock()
	defer mutex.Unlock()
	g := retainedSessions[uint64(handle)]
	if g == nil {
		return 7
	}
	if blob == nil || size == 0 || size > 1<<20 || keybag == nil || keybagSize == 0 || keybagSize > 16<<20 {
		return 1
	}
	candidate := C.GoBytes(unsafe.Pointer(blob), C.int(size))
	matched := bytes.Equal(g.blob, candidate)
	clear(candidate)
	if !matched {
		return 10
	}
	clear(g.keybag)
	g.keybag = C.GoBytes(unsafe.Pointer(keybag), C.int(keybagSize))
	// Binding is not consumption: no unverified guest entrypoint is invoked.
	return 0
}

//export WaffleKBSyncClose
func WaffleKBSyncClose(handle C.uint64_t) {
	mutex.Lock()
	defer mutex.Unlock()
	g := retainedSessions[uint64(handle)]
	if g == nil {
		return
	}
	delete(retainedSessions, uint64(handle))
	clear(g.blob)
	clear(g.keybag)
	clear(g.hardwareID)
	clear(g.expectedDigests)
	g.accountID = 0
	g.agent.Close()
}

//export WaffleKBSyncConsume
func WaffleKBSyncConsume(handle C.uint64_t, output **C.uchar, length *C.size_t) C.int {
	mutex.Lock()
	defer mutex.Unlock()
	g := retainedSessions[uint64(handle)]
	if g == nil {
		return 7
	}
	if output == nil || length == nil || len(g.keybag) == 0 {
		return 1
	}
	status, files, ioTrace, err := g.agent.ImportKeybag(g.keybag)
	// Keybags are single-use. Discard the retained copy immediately after the
	// emulated call, not only when the whole download session is closed.
	clear(g.keybag)
	g.keybag = nil
	saved, saveFailed := false, false
	if err == nil && status == 0 && g.accountDirectory != "" {
		saveErr := saveDownloadState(g.accountDirectory, g.accountID, g.hardwareID, g.expectedDigests, files)
		saved, saveFailed = saveErr == nil, saveErr != nil
	}
	result := map[string]any{"consumer": "storeagent-10.9-import-keybag", "native_status": status, "emulator_completed": err == nil, "import_succeeded": err == nil && status == 0, "portable_model_ready": false, "account_state_restored": g.accountRestored, "account_state_saved": saved, "account_state_save_failed": saveFailed, "state_files": files, "virtual_io": ioTrace}
	defer func() {
		for _, b := range files {
			clear(b)
		}
	}()
	if err != nil {
		result["error_category"] = "emulator_failure"
	}
	data, e := json.Marshal(result)
	if e != nil {
		return 2
	}
	defer clear(data)
	*output = (*C.uchar)(C.CBytes(data))
	*length = C.size_t(len(data))
	return 0
}
