package machine

import (
	"encoding/binary"
	"io"
	"path"
	"strings"
	"unicode/utf16"

	"github.com/majd/ipatool/v2/internal/sap/unicorn"
)

type guestStateFile struct {
	access          uint64
	appendMode      bool
	descriptorFlags uint64
	position        int64
	name            string
}
type guestStateFS struct {
	permissions map[string]guestPermissions
	finderInfo  map[string][16]byte
	refs        map[[80]byte]string
	files       map[string][]byte
	open        map[uint64]*guestStateFile
	next        uint64
	trace       []guestStateEvent
	traceCalls  uint64
	traceErrno  uint32
	traceHint   string
	traceMask   uint64
}

type guestStateEvent struct {
	Operation string `json:"operation"`
	Outcome   string `json:"outcome"`
	Errno     uint32 `json:"errno,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Mask      uint64 `json:"mask,omitempty"`
}

func (s *shims) rfs() *guestStateFS {
	if s.guestStateFS == nil {
		s.guestStateFS = &guestStateFS{files: map[string][]byte{}, open: map[uint64]*guestStateFile{}, next: 16}
	}
	return s.guestStateFS
}
func guestStatePath(p string) bool {
	return path.Clean(p) == p && path.Dir(p) == storeAgentSCInfoPath && !strings.Contains(path.Base(p), "\\")
}
func (s *shims) fsError(code uint32) error {
	s.rfs().traceErrno = code
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], code)
	if e := s.engine.MemWrite(s.errno, b[:]); e != nil {
		return e
	}
	return s.returnMinusOne()
}
func (s *shims) fsOpen(p string, flags uint64) error {
	if !guestStatePath(p) {
		return s.fsError(13)
	}
	if flags&3 == 3 {
		return s.fsError(22)
	}
	f := s.rfs()
	b, exists := f.files[p]
	if !exists && flags&0x200 == 0 {
		return s.fsError(2)
	}
	if exists && flags&0xa00 == 0xa00 {
		return s.fsError(17)
	}
	if (!exists && len(f.files) >= 32) || len(f.open) >= 32 {
		return s.fsError(24)
	}
	if flags&0x400 != 0 {
		if flags&3 == 0 {
			return s.fsError(9)
		}
		clear(b)
		b = nil
	}
	f.files[p] = b
	f.next++
	f.open[f.next] = &guestStateFile{name: p, access: flags & 3, appendMode: flags&8 != 0}
	return s.setResult(f.next)
}
func (s *shims) fsClose() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	if fd == coreFPFile {
		return s.setResult(0)
	}
	f := s.rfs()
	if _, ok := f.open[fd]; !ok {
		return s.fsError(9)
	}
	delete(f.open, fd)
	return s.setResult(0)
}
func (s *shims) fsRead(fd, ptr, n uint64) error {
	f := s.rfs().open[fd]
	if f == nil {
		return s.fsError(9)
	}
	if f.access == 1 {
		return s.fsError(9)
	}
	if n > 4<<20 {
		return s.fsError(22)
	}
	data := s.rfs().files[f.name]
	if f.position >= int64(len(data)) {
		return s.setResult(0)
	}
	count := min(n, uint64(int64(len(data))-f.position))
	if e := s.engine.MemWrite(ptr, data[f.position:f.position+int64(count)]); e != nil {
		return e
	}
	f.position += int64(count)
	return s.setResult(count)
}
func (s *shims) fsWrite() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	ptr, e := s.argument(1)
	if e != nil {
		return e
	}
	n, e := s.argument(2)
	if e != nil {
		return e
	}
	f := s.rfs().open[fd]
	if f == nil {
		return s.fsError(9)
	}
	if f.access == 0 {
		return s.fsError(9)
	}
	if f.appendMode {
		f.position = int64(len(s.rfs().files[f.name]))
	}
	if n > 1<<20 || f.position < 0 || f.position+int64(n) > 4<<20 {
		return s.fsError(27)
	}
	endSize := f.position + int64(n)
	total := int64(0)
	for name, data := range s.rfs().files {
		if name != f.name {
			total += int64(len(data))
		}
	}
	if total+max(int64(len(s.rfs().files[f.name])), endSize) > 4<<20 {
		return s.fsError(28)
	}
	b, e := s.engine.MemRead(ptr, n)
	if e != nil {
		return e
	}
	defer clear(b)
	data := s.rfs().files[f.name]
	end := f.position + int64(n)
	if end > int64(len(data)) {
		data = append(data, make([]byte, int(end)-len(data))...)
	}
	copy(data[f.position:], b)
	f.position = end
	s.rfs().files[f.name] = data
	return s.setResult(n)
}
func (s *shims) fsLock() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	if fd != coreFPFile && s.rfs().open[fd] == nil {
		return s.fsError(9)
	}
	return s.setResult(0)
}
func (s *shims) fsSeek() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	offset, e := s.argument(1)
	if e != nil {
		return e
	}
	whence, e := s.argument(2)
	if e != nil {
		return e
	}
	f := s.rfs().open[fd]
	if f == nil {
		return s.fsError(9)
	}
	next := int64(offset)
	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		next += f.position
	case io.SeekEnd:
		next += int64(len(s.rfs().files[f.name]))
	default:
		return s.fsError(22)
	}
	if next < 0 || next > 4<<20 {
		return s.fsError(22)
	}
	f.position = next
	return s.setResult(uint64(next))
}
func (s *shims) fsFstat() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	out, e := s.argument(1)
	if e != nil {
		return e
	}
	f := s.rfs().open[fd]
	if f == nil {
		return s.fsError(9)
	}
	perm := s.nodePermissions(f.name)
	b := guestStatBytes(perm.mode, uint64(len(s.rfs().files[f.name])))
	binary.LittleEndian.PutUint64(b[8:], uint64(guestNodeID(f.name)))
	binary.LittleEndian.PutUint32(b[16:], perm.uid)
	binary.LittleEndian.PutUint32(b[20:], perm.gid)
	if e = s.engine.MemWrite(out, b); e != nil {
		return e
	}
	return s.setResult(0)
}
func (s *shims) fsMkdir() error {
	ptr, e := s.argument(0)
	if e != nil {
		return e
	}
	p, e := s.readCString(ptr)
	if e != nil {
		return e
	}
	if p != storeAgentSCInfoPath {
		return s.fsError(13)
	}
	return s.fsError(17)
}

// Enabled only for the retained laboratory producer, after ordinary initialization.
func (s *shims) enableGuestState() error {
	for _, entry := range []struct {
		names   []string
		handler shimHandler
	}{
		{[]string{"_open", "_open$UNIX2003"}, s.guestOpen}, {[]string{"_read", "_read$UNIX2003"}, s.guestRead},
		{[]string{"_write", "_write$UNIX2003"}, s.fsWrite}, {[]string{"_lockf", "_fsync"}, s.fsLock},
		{[]string{"_close", "_close$UNIX2003"}, s.fsClose}, {[]string{"_lseek", "_lseek$UNIX2003"}, s.fsSeek},
		{[]string{"_FSSetCatalogInfo"}, s.fsSetCatalogInfo},
		{[]string{"_FSGetCatalogInfo"}, s.fsGetCatalogInfo},
		{[]string{"_FSPathMakeRef"}, s.fsPathMakeRef},
		{[]string{"_FSRefMakePath"}, s.fsRefMakePath},
		{[]string{"_fcntl", "_fcntl$UNIX2003"}, s.fsFcntl},
		{[]string{"_lstat", "_lstat$INODE64", "_stat", "_stat$INODE64"}, s.fsLstat},
		{[]string{"_fstat", "_fstat$INODE64"}, s.fsFstat}, {[]string{"_mkdir"}, s.fsMkdir}, {[]string{"_Gestalt"}, s.guestGestalt},
	} {
		// An operation label is fixed at compile time; never record user paths
		// or argument/register values from the emulated keybag consumer.
		wrapped := s.traceGuestOperation(entry.names[0], entry.handler)
		for _, name := range entry.names {
			if address, ok := s.symbols[name]; ok {
				s.entries[address] = shimEntry{name: name, handler: wrapped}
			} else if _, err := s.addFunction(name, wrapped); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *shims) guestOpen() error {
	ptr, e := s.argument(0)
	if e != nil {
		return e
	}
	p, e := s.readCString(ptr)
	if e != nil {
		return e
	}
	if !guestStatePath(p) {
		return s.open()
	}
	flags, e := s.argument(1)
	if e != nil {
		return e
	}
	return s.fsOpen(p, flags)
}
func (s *shims) guestRead() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	if fd == coreFPFile {
		return s.read()
	}
	ptr, e := s.argument(1)
	if e != nil {
		return e
	}
	n, e := s.argument(2)
	if e != nil {
		return e
	}
	return s.fsRead(fd, ptr, n)
}
func (s *shims) guestGestalt() error {
	selector, e := s.argument(0)
	if e != nil {
		return e
	}
	out, e := s.argument(1)
	if e != nil {
		return e
	}
	var value uint32
	switch uint32(selector) {
	case 0x73797376:
		value = 0x1090
	case 0x73797331:
		value = 10
	case 0x73797332:
		value = 9
	case 0x73797333:
		value = 0
	default:
		return s.setResult(0xffffea51)
	}
	if out == 0 {
		return s.setResult(0xffffffce)
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], value)
	if e = s.engine.MemWrite(out, b[:]); e != nil {
		return e
	}
	return s.setResult(0)
}

const guestStateTraceLimit = 64

// Retain only operation names, coarse return classes, and errno values.
// A raw return value might contain a pointer or descriptor, so it is never logged.
func (s *shims) traceGuestOperation(name string, handler shimHandler) shimHandler {
	return func() error {
		fs := s.rfs()
		fs.traceErrno = 0
		fs.traceHint = ""
		fs.traceMask = 0
		err := handler()
		entry := guestStateEvent{Operation: name}
		if name == "_FSSetCatalogInfo" {
			entry.Detail = fs.traceHint
			entry.Mask = fs.traceMask
		}
		switch {
		case err != nil:
			entry.Outcome = "emulator_failure"
		case fs.traceErrno != 0:
			entry.Outcome = "errno"
			entry.Errno = fs.traceErrno
		default:
			value, regErr := s.engine.RegRead(unicorn.RegRAX)
			switch {
			case regErr != nil:
				entry.Outcome = "result_unavailable"
			case value == 0:
				entry.Outcome = "zero"
			case int32(uint32(value)) < 0:
				entry.Outcome = "negative_status"
			default:
				entry.Outcome = "nonzero"
			}
		}
		fs.traceCalls++
		if len(fs.trace) == guestStateTraceLimit {
			copy(fs.trace, fs.trace[1:])
			fs.trace[len(fs.trace)-1] = entry
		} else {
			fs.trace = append(fs.trace, entry)
		}
		return err
	}
}

func (s *shims) guestStateTrace() map[string]any {
	fs := s.rfs()
	events := append([]guestStateEvent(nil), fs.trace...)
	return map[string]any{
		"total_calls": fs.traceCalls,
		"truncated": fs.traceCalls > guestStateTraceLimit,
		"recent": events,
	}
}

func (s *shims) guestStateSnapshot() map[string][]byte {
	out := map[string][]byte{}
	for name, data := range s.rfs().files {
		out[path.Base(name)] = append([]byte(nil), data...)
	}
	return out
}

// Darwin stat64 layout, shared by fstat and lstat. These are virtual nodes only.
func guestStatBytes(mode uint16, size uint64) []byte {
	b := make([]byte, 144)
	binary.LittleEndian.PutUint32(b[0:], 1)
	binary.LittleEndian.PutUint16(b[4:], mode)
	binary.LittleEndian.PutUint16(b[6:], 1)
	binary.LittleEndian.PutUint64(b[8:], 1)
	binary.LittleEndian.PutUint64(b[96:], size)
	binary.LittleEndian.PutUint64(b[104:], (size+511)/512)
	binary.LittleEndian.PutUint32(b[112:], 4096)
	return b
}
func (s *shims) fsLstat() error {
	ptr, e := s.argument(0)
	if e != nil {
		return e
	}
	out, e := s.argument(1)
	if e != nil {
		return e
	}
	p, e := s.readCString(ptr)
	if e != nil {
		return e
	}
	var mode uint16
	var size uint64
	switch p {
	case "/", "/Users":
		mode = 0040755
	case "/Users/Shared":
		mode = 0041777
	case storeAgentSCInfoPath:
		mode = 0040700
	default:
		if !guestStatePath(p) {
			return s.fsError(2)
		}
		data, ok := s.rfs().files[p]
		if !ok {
			return s.fsError(2)
		}
		mode = 0100600
		size = uint64(len(data))
	}
	if out == 0 {
		return s.fsError(14)
	}
	perm := s.nodePermissions(p)
	mode = perm.mode
	b := guestStatBytes(mode, size)
	binary.LittleEndian.PutUint64(b[8:], uint64(guestNodeID(p)))
	binary.LittleEndian.PutUint32(b[16:], perm.uid)
	binary.LittleEndian.PutUint32(b[20:], perm.gid)
	if e = s.engine.MemWrite(out, b); e != nil {
		return e
	}
	return s.setResult(0)
}
func (s *shims) fsFcntl() error {
	fd, e := s.argument(0)
	if e != nil {
		return e
	}
	command, e := s.argument(1)
	if e != nil {
		return e
	}
	file := s.rfs().open[fd]
	if file == nil && fd != coreFPFile {
		return s.fsError(9)
	}
	switch command {
	case 1:
		if file == nil {
			return s.setResult(0)
		}
		return s.setResult(file.descriptorFlags)
	case 2:
		flags, e := s.argument(2)
		if e != nil {
			return e
		}
		if file != nil {
			file.descriptorFlags = flags & 1
		}
		return s.setResult(0)
	case 48, 51: // descriptor flags/advisory cache/sync; one isolated guest, no host descriptor.
		return s.setResult(0)
	case 3: // F_GETFL
		if file == nil {
			return s.setResult(0)
		}
		flags := file.access
		if file.appendMode {
			flags |= 8
		}
		return s.setResult(flags)
	case 4: // F_SETFL changes append behavior, never the access mode.
		flags, e := s.argument(2)
		if e != nil {
			return e
		}
		if file != nil {
			file.appendMode = flags&8 != 0
		}
		return s.setResult(0)
	case 8, 9: // F_SETLK/F_SETLKW: no other guest shares this filesystem.
		return s.setResult(0)
	case 50: // F_GETPATH, bounded to Darwin's MAXPATHLEN.
		if file == nil {
			return s.fsError(22)
		}
		out, e := s.argument(2)
		if e != nil {
			return e
		}
		if out == 0 {
			return s.fsError(14)
		}
		if len(file.name) >= 1024 {
			return s.fsError(63)
		}
		if e = s.engine.MemWrite(out, append([]byte(file.name), 0)); e != nil {
			return e
		}
		return s.setResult(0)
	default:
		return s.fsError(22)
	}
}

// FSRef is opaque to callers. This token is valid only within this isolated VFS.
func (s *shims) fsPathMakeRef() error {
	ptr, e := s.argument(0)
	if e != nil {
		return e
	}
	out, e := s.argument(1)
	if e != nil {
		return e
	}
	isDir, e := s.argument(2)
	if e != nil {
		return e
	}
	p, e := s.readCString(ptr)
	if e != nil {
		return e
	}
	directory := p == "/" || p == "/Users" || p == "/Users/Shared" || p == storeAgentSCInfoPath
	if !directory {
		if !guestStatePath(p) {
			return s.setResult(0xffffffd5)
		}
		if _, ok := s.rfs().files[p]; !ok {
			return s.setResult(0xffffffd5)
		}
	} // fnfErr
	if out == 0 {
		return s.setResult(0xffffffce)
	}
	fs := s.rfs()
	if fs.refs == nil {
		fs.refs = map[[80]byte]string{}
	}
	var token [80]byte
	copy(token[:], "CrumbVFS")
	found := false
	for candidate, name := range fs.refs {
		if name == p {
			token = candidate
			found = true
			break
		}
	}
	if !found {
		if len(fs.refs) >= 64 {
			return s.setResult(0xffffffce)
		}
		binary.LittleEndian.PutUint64(token[8:], uint64(len(fs.refs)+1))
		fs.refs[token] = p
	}
	if e = s.engine.MemWrite(out, token[:]); e != nil {
		return e
	}
	if isDir != 0 {
		var value byte
		if directory {
			value = 1
		}
		if e = s.engine.MemWrite(isDir, []byte{value}); e != nil {
			return e
		}
	}
	return s.setResult(0)
}
func (s *shims) fsRefMakePath() error {
	ref, e := s.argument(0)
	if e != nil {
		return e
	}
	out, e := s.argument(1)
	if e != nil {
		return e
	}
	maximum, e := s.argument(2)
	if e != nil {
		return e
	}
	bytes, e := s.engine.MemRead(ref, 80)
	if e != nil {
		return e
	}
	var token [80]byte
	copy(token[:], bytes)
	clear(bytes)
	name, ok := s.rfs().refs[token]
	if !ok || out == 0 || uint64(len(name)+1) > maximum {
		return s.setResult(0xffffffce)
	}
	if e = s.engine.MemWrite(out, append([]byte(name), 0)); e != nil {
		return e
	}
	return s.setResult(0)
}

func guestNodeID(p string) uint32 {
	switch p {
	case "/":
		return 2
	case "/Users":
		return 3
	case "/Users/Shared":
		return 4
	case storeAgentSCInfoPath:
		return 5
	}
	id := uint32(2166136261)
	for _, c := range []byte(p) {
		id ^= uint32(c)
		id *= 16777619
	}
	return id | 0x100
}
func (s *shims) fsGetCatalogInfo() error {
	ref, e := s.argument(0)
	if e != nil {
		return e
	}
	mask, e := s.argument(1)
	if e != nil {
		return e
	}
	catalog, e := s.argument(2)
	if e != nil {
		return e
	}
	outName, e := s.argument(3)
	if e != nil {
		return e
	}
	spec, e := s.argument(4)
	if e != nil {
		return e
	}
	parent, e := s.argument(5)
	if e != nil {
		return e
	}
	if spec != 0 || mask & ^uint64(0xc41e) != 0 {
		return s.setResult(0xffffffce)
	}
	bytes, e := s.engine.MemRead(ref, 80)
	if e != nil {
		return e
	}
	var token [80]byte
	copy(token[:], bytes)
	clear(bytes)
	p, ok := s.rfs().refs[token]
	if !ok {
		return s.setResult(0xffffffce)
	}
	directory := p == "/" || p == "/Users" || p == "/Users/Shared" || p == storeAgentSCInfoPath
	if catalog != 0 {
		b := make([]byte, 148) // Files.h: #pragma pack(2), LP64 FSPermissionInfo includes an 8-byte pointer.
		if mask&2 != 0 && directory {
			binary.LittleEndian.PutUint16(b, 0x10)
		}
		if mask&4 != 0 {
			binary.LittleEndian.PutUint16(b[2:], 1)
		}
		if mask&8 != 0 {
			binary.LittleEndian.PutUint32(b[4:], guestNodeID(path.Dir(p)))
		}
		if mask&0x10 != 0 {
			binary.LittleEndian.PutUint32(b[8:], guestNodeID(p))
		}
		if mask&0x400 != 0 {
			perm := s.nodePermissions(p)
			binary.LittleEndian.PutUint32(b[56:], perm.uid)
			binary.LittleEndian.PutUint32(b[60:], perm.gid)
			binary.LittleEndian.PutUint16(b[66:], perm.mode)
		}
		if mask&0x800 != 0 && s.rfs().finderInfo != nil {
			finder := s.rfs().finderInfo[p]
			copy(b[72:88], finder[:]) // Files.h FSCatalogInfo.finderInfo[16].
		}
		if mask&0x4000 != 0 && !directory {
			size := uint64(len(s.rfs().files[p]))
			binary.LittleEndian.PutUint64(b[108:], size)
			binary.LittleEndian.PutUint64(b[116:], (size+4095)&^4095)
		}
		if e = s.engine.MemWrite(catalog, b); e != nil {
			return e
		}
	} else if mask != 0 {
		return s.setResult(0xffffffce)
	}
	if outName != 0 {
		units := utf16.Encode([]rune(path.Base(p)))
		if len(units) > 255 {
			return s.setResult(0xffffffce)
		}
		b := make([]byte, 512)
		binary.LittleEndian.PutUint16(b, uint16(len(units)))
		for i, v := range units {
			binary.LittleEndian.PutUint16(b[2+2*i:], v)
		}
		if e = s.engine.MemWrite(outName, b); e != nil {
			return e
		}
	}
	if parent != 0 {
		var parentToken [80]byte
		found := false
		parentPath := path.Dir(p)
		for candidate, name := range s.rfs().refs {
			if name == parentPath {
				parentToken = candidate
				found = true
				break
			}
		}
		if !found {
			if len(s.rfs().refs) >= 64 {
				return s.setResult(0xffffffce)
			}
			copy(parentToken[:], "CrumbVFS")
			binary.LittleEndian.PutUint64(parentToken[8:], uint64(len(s.rfs().refs)+1))
			s.rfs().refs[parentToken] = parentPath
		}
		if e = s.engine.MemWrite(parent, parentToken[:]); e != nil {
			return e
		}
	}
	return s.setResult(0)
}

type guestPermissions struct {
	uid, gid uint32
	mode     uint16
}

func (s *shims) nodePermissions(p string) guestPermissions {
	if v, ok := s.rfs().permissions[p]; ok {
		return v
	}
	switch p {
	case "/", "/Users":
		return guestPermissions{mode: 0040755}
	case "/Users/Shared":
		return guestPermissions{mode: 0041777}
	case storeAgentSCInfoPath:
		return guestPermissions{mode: 0040700}
	}
	return guestPermissions{mode: 0100600}
}
func (s *shims) fsSetCatalogInfo() error {
	ref, e := s.argument(0)
	if e != nil {
		return e
	}
	mask, e := s.argument(1)
	if e != nil {
		return e
	}
	info, e := s.argument(2)
	if e != nil {
		return e
	}
	// Metadata flags are public ABI constants, not private material. Record only
	// small expected masks; never publish register values that might be pointers.
	if mask <= 0xffff {
		s.rfs().traceMask = mask
	}
	// Carbon Files.h: 0x400 = permissions, 0x800 = 16-byte Finder info.
	// Only these observed catalog write fields are emulated, in any combination.
	if mask == 0 || mask & ^uint64(0xc00) != 0 {
		s.rfs().traceHint = "unsupported_catalog_mask"
		return s.setResult(0xffffffce)
	}
	if info == 0 {
		s.rfs().traceHint = "missing_catalog_info"
		return s.setResult(0xffffffce)
	}
	b, e := s.engine.MemRead(ref, 80)
	if e != nil {
		return e
	}
	var token [80]byte
	copy(token[:], b)
	clear(b)
	p, ok := s.rfs().refs[token]
	if !ok {
		s.rfs().traceHint = "unknown_catalog_reference"
		return s.setResult(0xffffffce)
	}
	b, e = s.engine.MemRead(info, 148)
	if e != nil {
		return e
	}
	defer clear(b)
	// Honor whichInfo: updating Finder metadata alone must not overwrite
	// permissions, and vice versa. Both are retained in the private VFS.
	if mask&0x400 != 0 {
		value := s.nodePermissions(p)
		value.uid = binary.LittleEndian.Uint32(b[56:])
		value.gid = binary.LittleEndian.Uint32(b[60:])
		value.mode = value.mode&0170000 | binary.LittleEndian.Uint16(b[66:])&07777
		if s.rfs().permissions == nil {
			s.rfs().permissions = map[string]guestPermissions{}
		}
		s.rfs().permissions[p] = value
	}
	if mask&0x800 != 0 {
		var finder [16]byte
		copy(finder[:], b[72:88])
		if s.rfs().finderInfo == nil {
			s.rfs().finderInfo = map[string][16]byte{}
		}
		s.rfs().finderInfo[p] = finder
	}
	return s.setResult(0)
}
