package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"strings"
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
		result["stage"] = r.Operation
		result["code"] = "validation_failed"
		message := strings.ToLower(e.Error())
		if r.Operation == "recover-enrollment" || r.Operation == "recover-capture" {
			for _, item := range []struct{ fragment, code string }{
				{"schema required", "modern_native_capture_required"},
				{"legacy state not compatible", "legacy_state_not_compatible"},
				{"not bound to captured transaction", "capture_ipa_mismatch"},
				{"state and one installed sinf required", "missing_native_enrollment_state"},
				{"state and enrollment frame disagree", "enrollment_frame_mismatch"},
				{"target ipa validation failed", "model_ipa_mismatch"},
				{"modern enrollment derivation failed", "enrollment_derivation_failed"},
				{"capture integrity", "capture_integrity_failed"},
			} {
				if strings.Contains(message, item.fragment) {
					result["code"] = item.code
					break
				}
			}
		}
		for _, category := range []string{"model", "sinf", "supf", "codedirectory", "arm64", "zip", "permission", "space", "profile"} {
			if strings.Contains(message, category) {
				result["stage"] = r.Operation + "/" + category
				break
			}
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
