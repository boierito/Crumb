// Copied into the pinned internal/sap/machine package.
package machine

import (
	"fmt"
	"path"
	"strings"
)

// cloneGuestState validates the complete saved filesystem before allocating
// private copies. No account record is parsed, modified or synthesized.
func cloneGuestState(files map[string][]byte) (map[string][]byte, error) {
	if len(files) == 0 || len(files) > 32 || len(files["SC Info.sidb"]) < 4 {
		return nil, fmt.Errorf("producer_state_invalid")
	}
	total := 0
	for name, data := range files {
		if name == "." || path.Base(name) != name || strings.ContainsAny(name, "\\\x00") || !guestStatePath(storeAgentSCInfoPath+"/"+name) {
			return nil, fmt.Errorf("producer_state_invalid")
		}
		total += len(data)
		if total > 4<<20 {
			return nil, fmt.Errorf("producer_state_invalid")
		}
	}
	result := make(map[string][]byte, len(files))
	for name, data := range files {
		result[storeAgentSCInfoPath+"/"+name] = append([]byte(nil), data...)
	}
	return result, nil
}

// installGuestState is one-shot and must precede native global initialization.
// The caller owns and wipes files unless installation succeeds.
func (s *shims) installGuestState(files map[string][]byte) error {
	if s.guestStateFS != nil {
		return fmt.Errorf("producer_state_already_initialized")
	}
	if err := s.enableGuestState(); err != nil {
		return err
	}
	s.rfs().files = files
	return nil
}
