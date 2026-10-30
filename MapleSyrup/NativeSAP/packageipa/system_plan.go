package packageipa

import (
	"encoding/binary"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type SystemNativeRange struct {
	SliceOffset  uint64 `json:"slice_offset"`
	SliceSize    uint64 `json:"slice_size"`
	CryptOffset  uint64 `json:"crypt_offset"`
	CryptSize    uint64 `json:"crypt_size"`
	CryptidField uint64 `json:"cryptid_field"`
	CPUType      uint32 `json:"cpu_type"`
	CPUSubtype   uint32 `json:"cpu_subtype"`
}
type SystemNativeFile struct {
	RelativePath string              `json:"relative_path"`
	Ranges       []SystemNativeRange `json:"ranges"`
}
type SystemNativePlan struct {
	Schema          string             `json:"schema"`
	Files           []SystemNativeFile `json:"files"`
	EncryptedSlices int                `json:"encrypted_slices"`
}

func PlanSystemTarget(source string) (plan SystemNativePlan, err error) {
	plan.Schema = "crumb-system-fairplay-native-plan-v1"
	stat, e := os.Lstat(source)
	if e != nil || !stat.IsDir() {
		return plan, errors.New("installed target unavailable")
	}
	err = filepath.WalkDir(source, func(p string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		relative, e := filepath.Rel(source, p)
		if e != nil {
			return e
		}
		if relative == "." {
			return nil
		}
		if !systemSafe(filepath.ToSlash(relative)) {
			return errors.New("invalid target resource path")
		}
		if systemPrivatePath(relative) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return errors.New("symlink in installed target")
		}
		if entry.IsDir() {
			return nil
		}
		stat, e := entry.Info()
		if e != nil {
			return e
		}
		if !stat.Mode().IsRegular() {
			return errors.New("unsupported target resource")
		}
		f, e := os.Open(p)
		if e != nil {
			return e
		}
		defer f.Close()
		image, e := systemMachO(f, stat.Size())
		if e != nil {
			return e
		}
		if image == nil {
			return nil
		}
		file := SystemNativeFile{RelativePath: filepath.ToSlash(relative)}
		for _, s := range image.slices {
			for _, c := range s.crypt {
				if c.id == 0 {
					continue
				}
				cpu, subtype := binary.LittleEndian.Uint32(s.header[4:]), binary.LittleEndian.Uint32(s.header[8:])
				if cpu != 0x100000c || s.headerSize != 32 {
					return errors.New("encrypted architecture is unsupported by this iPhone route")
				}
				file.Ranges = append(file.Ranges, SystemNativeRange{uint64(s.off), uint64(s.size), uint64(c.off - s.off), uint64(c.size), uint64(c.field), cpu, subtype})
				plan.EncryptedSlices++
			}
		}
		if len(file.Ranges) > 0 {
			plan.Files = append(plan.Files, file)
		}
		return nil
	})
	if err == nil && plan.EncryptedSlices == 0 {
		err = errors.New("no encrypted target executables")
	}
	return plan, err
}
