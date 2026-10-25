package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
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
		if diagnostic, ok := rawsource.RecoveryDetails(e); ok {
			result["runtime"] = diagnostic
		}
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
