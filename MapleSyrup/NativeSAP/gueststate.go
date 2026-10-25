package machine

import (
	"encoding/binary"
	"io"
	"path"
	"strings"
)

type guestStateFile struct {
	access     uint64
	appendMode bool
	position   int64
	name       string
}
type guestStateFS struct {
	files map[string][]byte
	open  map[uint64]*guestStateFile
	next  uint64
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
	b := make([]byte, 144)
	binary.LittleEndian.PutUint16(b[4:], 0100600)
	binary.LittleEndian.PutUint64(b[8:], 1)
	binary.LittleEndian.PutUint64(b[96:], uint64(len(s.rfs().files[f.name])))
	binary.LittleEndian.PutUint32(b[112:], 4096)
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
		{[]string{"_fstat", "_fstat$INODE64"}, s.fsFstat}, {[]string{"_mkdir"}, s.fsMkdir}, {[]string{"_Gestalt"}, s.guestGestalt},
	} {
		for _, name := range entry.names {
			if address, ok := s.symbols[name]; ok {
				s.entries[address] = shimEntry{name: name, handler: entry.handler}
			} else if _, err := s.addFunction(name, entry.handler); err != nil {
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
func (s *shims) guestStateSnapshot() map[string][]byte {
	out := map[string][]byte{}
	for name, data := range s.rfs().files {
		out[path.Base(name)] = append([]byte(nil), data...)
	}
	return out
}
