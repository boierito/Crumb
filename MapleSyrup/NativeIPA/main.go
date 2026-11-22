// Local C ABI for streaming IPA inspection and preparation. No SAP emulation.
// Derived from ipatool package preparation (MIT); see THIRD_PARTY_NOTICES.md.
package main

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"
	"crumb/nativeipa/packageipa"
)

//export CrumbPrepareIPA
func CrumbPrepareIPA(source *C.char, destination *C.char, input *C.uchar, length C.size_t, output **C.uchar, outputLength *C.size_t) C.int {
	if source == nil || destination == nil || input == nil || length == 0 || length > 16<<20 || output == nil || outputLength == nil {
		return 1
	}
	var parameters packageipa.Input
	if json.Unmarshal(C.GoBytes(unsafe.Pointer(input), C.int(length)), &parameters) != nil {
		return 20
	}
	data, err := packageipa.Prepare(C.GoString(source), C.GoString(destination), parameters)
	if err != nil {
		return C.int(packageipa.FailureCode(err))
	}
	*output = (*C.uchar)(C.CBytes(data))
	*outputLength = C.size_t(len(data))
	return 0
}

//export CrumbInspectIPA
func CrumbInspectIPA(url *C.char, bundle *C.char, output **C.uchar, outputLength *C.size_t) C.int {
	if url == nil || bundle == nil || output == nil || outputLength == nil {
		return 1
	}
	data, err := packageipa.InspectRemote(C.GoString(url), C.GoString(bundle))
	if err != nil {
		return 22
	}
	*output = (*C.uchar)(C.CBytes(data))
	*outputLength = C.size_t(len(data))
	return 0
}

//export CrumbIPAFree
func CrumbIPAFree(pointer *C.uchar, length C.size_t) {
	if pointer != nil {
		C.memset(unsafe.Pointer(pointer), 0, length)
		C.free(unsafe.Pointer(pointer))
	}
}
func main() {}
