package main

import (
	"debug/macho"
	"fmt"
)

type executableSection struct {
	Name    string `json:"name"`
	Segment string `json:"segment"`
	Address uint64 `json:"address"`
	Bytes   uint64 `json:"bytes"`
}

func executableImage(path string) ([]byte, []bool, uint64, []executableSection, error) {
	file, err := macho.Open(path)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	defer file.Close()
	if file.Cpu != macho.CpuArm64 || uint32(file.SubCpu)&0x00ffffff != 2 {
		return nil, nil, 0, nil, fmt.Errorf("exact thin arm64e image required")
	}
	var sections []executableSection
	var selected []*macho.Section
	var base, end uint64
	for _, section := range file.Sections {
		if section.Flags&0x80000400 == 0 || section.Size == 0 {
			continue
		}
		if section.Addr%4 != 0 || section.Size%4 != 0 || section.Addr+section.Size < section.Addr {
			return nil, nil, 0, nil, fmt.Errorf("invalid executable section")
		}
		if len(selected) == 0 || section.Addr < base {
			base = section.Addr
		}
		if section.Addr+section.Size > end {
			end = section.Addr + section.Size
		}
		selected = append(selected, section)
		sections = append(sections, executableSection{section.Name, section.Seg, section.Addr, section.Size})
	}
	if len(selected) == 0 || end-base > 32<<20 {
		return nil, nil, 0, nil, fmt.Errorf("invalid executable image span")
	}
	image := make([]byte, end-base)
	valid := make([]bool, (end-base)/4)
	for _, section := range selected {
		data, err := section.Data()
		if err != nil || uint64(len(data)) != section.Size {
			return nil, nil, 0, nil, fmt.Errorf("executable section bytes unavailable")
		}
		offset := section.Addr - base
		copy(image[offset:], data)
		for slot := offset / 4; slot < (offset+section.Size)/4; slot++ {
			if valid[slot] {
				return nil, nil, 0, nil, fmt.Errorf("overlapping executable sections")
			}
			valid[slot] = true
		}
	}
	return image, valid, base, sections, nil
}
