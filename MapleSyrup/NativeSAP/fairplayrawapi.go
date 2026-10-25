// Copied into the isolated rawsource package; source inputs are modern native
// enrollment state, not the legacy StoreAgent capture.
package rawsource

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

var recoveryLock sync.Mutex

type Result struct {
	Source, Frame, Pack []byte
}

func VerifyResource() bool {
	sum := sha256.Sum256(rawSourceModel)
	return hex.EncodeToString(sum[:]) == "RAW_RESOURCE_SHA256"
}

// Recover constructs the source field, selection frame and selector-bound pack
// in process. The caller must validate the pack against the target IPA before
// treating it as usable. No material is logged or returned in diagnostics.
func Recover(executable string, hardware []byte, files map[string][]byte) (Result, error) {
	var result Result
	if len(hardware) != 6 || executable == "" || len(files) == 0 || len(files) > 256 {
		return result, fmt.Errorf("invalid modern enrollment inputs")
	}
	recoveryLock.Lock()
	defer recoveryLock.Unlock()
	if !VerifyResource() {
		return result, fmt.Errorf("raw source resource integrity")
	}
	value, err := loadModelBytes(rawSourceModel)
	if err != nil {
		return result, err
	}
	defer func() {
		if activeRawSourceRuntime != nil {
			for _, data := range activeRawSourceRuntime.files {
				clear(data)
			}
			clear(activeRawSourceRuntime.mac)
		}
		if lastRawSourceState != nil {
			for _, data := range lastRawSourceState.writable {
				clear(data)
			}
			clear(lastRawSourceState.x[:])
		}
		for _, data := range value.pages {
			clear(data)
		}
		activeRawSourceRuntime, lastRawSourceState = nil, nil
	}()
	if err := configureRawSourceRuntime(value, "enrollment-sinf", nil, "", executable, 0); err != nil {
		return result, err
	}
	activeRawSourceRuntime.mac = append([]byte(nil), hardware...)
	total := 0
	for name, data := range files {
		total += len(data)
		if name == "" || len(data) > 64<<20 || total > 128<<20 {
			return result, fmt.Errorf("modern enrollment input limit")
		}
		activeRawSourceRuntime.files[name] = append([]byte(nil), data...)
	}
	if _, _, err := runRawSource(value, "enrollment-sinf", false, 0); err != nil {
		// Upstream errors can include private paths, addresses and state.
		return result, fmt.Errorf("modern enrollment derivation failed")
	}
	result.Pack, err = currentRawSourceEnrollmentPack()
	if err != nil {
		return Result{}, fmt.Errorf("modern enrollment pack unavailable")
	}
	result.Source, err = currentRawSourceBytes(224)
	if err != nil {
		clear(result.Pack)
		return Result{}, fmt.Errorf("modern source unavailable")
	}
	result.Frame, err = currentRawSourceFrameBytes()
	if err != nil {
		clear(result.Pack)
		clear(result.Source)
		return Result{}, fmt.Errorf("modern frame unavailable")
	}
	return result, nil
}
