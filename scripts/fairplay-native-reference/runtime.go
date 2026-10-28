package rawsource

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type nativeReferenceEvent struct {
	Type         string
	Stop         string `json:",omitempty"`
	BinarySHA256 string `json:",omitempty"`
	CPU          nativeReferenceCPU
	After        *nativeReferenceCPU `json:",omitempty"`
	Pages        map[uint64][]byte   `json:",omitempty"`
	NativeWrites map[uint64]string   `json:",omitempty"`
	Removed      []uint64            `json:",omitempty"`
	Failed       bool                `json:",omitempty"`
}

var nativeReferenceFile *os.File
var nativeReferenceBuffer *bufio.Writer
var nativeReferenceEncoder *json.Encoder
var nativeReferenceAttempted, nativeReferenceExternal, nativeReferenceClosed bool
var nativeReferenceDirty = map[uint64]bool{}
var nativeReferenceBeforeHashes map[uint64][32]byte
var nativeReferencePending nativeReferenceEvent
var nativeReferenceCalls int

func nativeReferencePages(s *state) map[uint64][]byte {
	result := make(map[uint64][]byte, len(s.model.pages))
	for a, b := range s.model.pages {
		if w := s.writable[a]; w != nil {
			b = w
		}
		result[a] = b
	}
	return result
}
func nativeReferenceHashes(s *state) map[uint64][32]byte {
	result := make(map[uint64][32]byte, len(s.model.pages))
	for a, b := range nativeReferencePages(s) {
		result[a] = sha256.Sum256(b)
	}
	return result
}
func nativeReferenceWriteHashes(s *state) map[uint64]string {
	result := map[uint64]string{}
	pages := nativeReferencePages(s)
	for a := range nativeReferenceDirty {
		if b := pages[a]; b != nil {
			h := sha256.Sum256(b)
			result[a] = hex.EncodeToString(h[:])
		}
	}
	return result
}
func nativeReferenceEmit(e nativeReferenceEvent) {
	if nativeReferenceEncoder == nil || nativeReferenceClosed {
		return
	}
	if nativeReferenceEncoder.Encode(e) != nil || nativeReferenceBuffer.Flush() != nil {
		nativeReferenceClosed = true
		nativeReferenceFile.Close()
		fmt.Fprintln(os.Stderr, "native reference trace write failed")
	}
}
func nativeReferenceStart(s *state, pc uint64) {
	if pc != nativeEnrollmentStart {
		return
	}
	path := os.Getenv("CRUMB_NATIVE_REFERENCE_TRACE")
	if path == "" {
		return
	}
	if nativeReferenceAttempted {
		return
	}
	nativeReferenceAttempted = true
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		fmt.Fprintln(os.Stderr, "native reference trace unavailable")
		return
	}
	nativeReferenceFile = f
	nativeReferenceBuffer = bufio.NewWriterSize(f, 1<<20)
	nativeReferenceEncoder = json.NewEncoder(nativeReferenceBuffer)
	nativeReferenceEmit(nativeReferenceEvent{Type: "initial", BinarySHA256: "1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784", CPU: nativeReferenceSnapshot(s, pc), Pages: nativeReferencePages(s)})
}
func nativeReferenceStore(s *state, address uint64, size int) {
	if nativeReferenceEncoder == nil || nativeReferenceClosed || nativeReferenceExternal || s != lastRawSourceState || size <= 0 {
		return
	}
	end := address + uint64(size)
	if end < address {
		return
	}
	for a := address &^ uint64(4095); a < end; a += 4096 {
		nativeReferenceDirty[a] = true
	}
}
func nativeReferenceBefore(s *state, pc uint64) {
	if nativeReferenceEncoder == nil || nativeReferenceClosed {
		return
	}
	nativeReferencePending = nativeReferenceEvent{Type: "external", CPU: nativeReferenceSnapshot(s, pc), NativeWrites: nativeReferenceWriteHashes(s)}
	nativeReferenceBeforeHashes = nativeReferenceHashes(s)
	nativeReferenceExternal = true
}
func nativeReferenceAfter(s *state, next uint64, err error) {
	if nativeReferenceEncoder == nil || nativeReferenceClosed {
		return
	}
	e := nativeReferencePending
	after := nativeReferenceSnapshot(s, next)
	e.After = &after
	e.Pages = map[uint64][]byte{}
	e.Failed = err != nil
	pages := nativeReferencePages(s)
	for a, b := range pages {
		before, exists := nativeReferenceBeforeHashes[a]
		if !exists || sha256.Sum256(b) != before {
			e.Pages[a] = b
		}
	}
	for a := range nativeReferenceBeforeHashes {
		if pages[a] == nil {
			e.Removed = append(e.Removed, a)
		}
	}
	nativeReferenceEmit(e)
	clear(nativeReferenceDirty)
	nativeReferenceBeforeHashes = nil
	nativeReferencePending = nativeReferenceEvent{}
	nativeReferenceExternal = false
	nativeReferenceCalls++
}
func nativeReferenceFinish(s *state, pc uint64, err error) {
	if nativeReferenceEncoder == nil || nativeReferenceClosed {
		return
	}
	stop := "incomplete"
	if err == nil {
		stop = "boundary_reached"
	} else if _, ok := err.(*nativeReturnError); ok {
		stop = "native_return"
	}
	nativeReferenceEmit(nativeReferenceEvent{Type: "final", Stop: stop, CPU: nativeReferenceSnapshot(s, pc), NativeWrites: nativeReferenceWriteHashes(s)})
	nativeReferenceClosed = true
	nativeReferenceBuffer.Flush()
	nativeReferenceFile.Sync()
	nativeReferenceFile.Close()
	fmt.Fprintln(os.Stderr, "native reference trace completed", "external_calls", nativeReferenceCalls)
}

// CPU and memory records contain private licensing material. This file is
// injected only into an explicitly requested host reference build, never iOS.
type nativeReferenceCPU struct {
	PC, SP    uint64
	Registers [31]uint64
	Flags     [4]bool
}

func nativeReferenceSnapshot(s *state, pc uint64) nativeReferenceCPU {
	return nativeReferenceCPU{PC: pc, SP: s.sp, Registers: s.x, Flags: [4]bool{s.n, s.z, s.c, s.v}}
}
