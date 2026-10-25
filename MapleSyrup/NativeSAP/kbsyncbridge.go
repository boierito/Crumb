package main

/*
#include <stdint.h>
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"context"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"github.com/majd/ipatool/v2/internal/sap/machine"
	"time"
	"unsafe"
)

type retainedKBSync struct {
	agent        *machine.StoreAgent
	blob, keybag []byte
}

var retainedSessions = map[uint64]*retainedKBSync{}

//export WaffleKBSyncOpen
func WaffleKBSyncOpen(cache *C.char, hardware *C.uchar, length C.size_t, dsid C.uint64_t, handle *C.uint64_t, output **C.uchar, size *C.size_t) C.int {
	mutex.Lock()
	defer mutex.Unlock()
	if cache == nil || hardware == nil || length != 6 || dsid == 0 || handle == nil || output == nil || size == nil {
		return 1
	}
	assets.CacheRoot = C.GoString(cache)
	if assets.CacheRoot == "" {
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	bundle, err := assets.Load(ctx)
	if err != nil {
		return 2
	}
	agent, blob, err := machine.OpenRetainedKBSync(ctx, bundle, C.GoBytes(unsafe.Pointer(hardware), 6), uint64(dsid))
	if err != nil {
		return 8
	}
	nextHandle++
	if nextHandle == 0 {
		agent.Close()
		return 7
	}
	retainedSessions[nextHandle] = &retainedKBSync{agent: agent, blob: blob}
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
	if !bytes.Equal(g.blob, C.GoBytes(unsafe.Pointer(blob), C.int(size))) {
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
	g.agent.Close()
}
