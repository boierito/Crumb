package rawsource

import "fmt"

// Preserve the reference allocator's minimum/spacing for small requests, but
// reserve the entire mapped extent. Mapping existing pages must never alias a
// prior live allocation or static model page. No absent guest pointer is mapped
// here: only an explicit supported allocation request creates memory.
func allocateSourceMemory(r *rawSourceRuntime, s *state, size uint64) (uint64, error) {
	if size > 0x10000000 {
		return 0, fmt.Errorf("invalid bootstrap allocation size")
	}
	mapped := (max(size, 1) + pageSize - 1) &^ (pageSize - 1)
	if mapped < 0x100000 {
		mapped = 0x100000
	}
	address := r.nextAllocation
	if address == 0 || address&(pageSize-1) != 0 || address+mapped < address {
		return 0, fmt.Errorf("invalid bootstrap allocation extent")
	}
	// Refuse collision atomically. In particular mapZero's existing-page
	// behavior must not convert a collision into an apparently valid malloc.
	for page := address; page < address+mapped; page += pageSize {
		if s.model.pages[page] != nil {
			return 0, fmt.Errorf("bootstrap allocation collision")
		}
	}
	if err := s.mapZero(address, mapped); err != nil {
		return 0, err
	}
	r.nextAllocation = address + mapped
	return address, nil
}
