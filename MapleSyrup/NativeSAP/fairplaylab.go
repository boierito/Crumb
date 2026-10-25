package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"time"
)

//export WaffleFairPlayLab
func WaffleFairPlayLab(input *C.char, output **C.uchar, length *C.size_t) C.int {
	if input == nil || output == nil || length == nil {
		return 1
	}
	var r labRequest
	if json.Unmarshal([]byte(C.GoString(input)), &r) != nil {
		return 1
	}
	start := time.Now()
	v, e := labRun(r)
	// Do not expose native error strings: they may contain private file paths/material.
	result := map[string]any{"operation": r.Operation, "ok": e == nil, "milliseconds": time.Since(start).Milliseconds()}
	if e != nil {
		stage, code := labErrorDetails(r.Operation, e)
		result["stage"], result["code"] = stage, code
	} else {
		result["result"] = v
	}
	data, e := json.Marshal(result)
	if e != nil {
		return 2
	}
	*output = (*C.uchar)(C.CBytes(data))
	*length = C.size_t(len(data))
	return 0
}
