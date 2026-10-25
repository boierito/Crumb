package main

import (
	"errors"
	"strings"
)

// Only fixed stage/code literals are surfaced. The original error can contain
// account paths or state and must never be serialized by the C bridge.
type labStageError struct {
	stage, code string
	cause       error
}

func (e *labStageError) Error() string { return e.code }
func (e *labStageError) Unwrap() error { return e.cause }
func recoveryFailure(stage, code string, cause error) error {
	var classified *labStageError
	if errors.As(cause, &classified) {
		return cause
	}
	return &labStageError{stage: stage, code: code, cause: cause}
}

func labErrorDetails(operation string, err error) (string, string) {
	var classified *labStageError
	if errors.As(err, &classified) {
		return operation + "/" + classified.stage, classified.code
	}
	stage, code := operation, "validation_failed"
	message := strings.ToLower(err.Error())
	if operation == "recover-enrollment" || operation == "recover-capture" {
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
				code = item.code
				break
			}
		}
	}
	for _, category := range []string{"model", "sinf", "supf", "codedirectory", "arm64", "zip", "permission", "space", "profile"} {
		if strings.Contains(message, category) {
			stage = operation + "/" + category
			break
		}
	}
	return stage, code
}
