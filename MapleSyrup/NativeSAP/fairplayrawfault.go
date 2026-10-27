package rawsource

import "errors"

// The address stays inside this package. Error() and JSON diagnostics never
// expose it. Tests can inspect typed faults without parsing private text.
type sourceMemoryError struct {
	address uint64
	write   bool
}

func (*sourceMemoryError) Error() string { return "unmapped guest memory" }

func memoryFaultRegion(err error) string {
	var fault *sourceMemoryError
	if !errors.As(err, &fault) {
		return ""
	}
	switch {
	case fault.address < pageSize:
		return "null_page"
	case fault.address >= 0x100000000 && fault.address < 0x100ff0000:
		return "immutable_image"
	case fault.address >= 0x100ff0000 && fault.address < 0x101088000:
		return "mutable_image"
	case fault.address >= 0x700000000000 && fault.address < 0x720000000000:
		return "virtual_state_window"
	default:
		return "outside_model_windows"
	}
}
