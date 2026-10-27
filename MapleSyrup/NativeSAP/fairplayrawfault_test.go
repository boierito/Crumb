package rawsource

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestMemoryFaultTaxonomyDoesNotExposeAddresses(t *testing.T) {
	for _, trial := range []struct {
		address uint64
		region  string
	}{
		{0xb40, "null_page"}, {0x10008a680, "immutable_image"},
		{0x101075b40, "mutable_image"}, {0x710100001234, "virtual_state_window"},
		{0xabcdef0123456789, "outside_model_windows"},
	} {
		err := fmt.Errorf("wrapper: %w", &sourceMemoryError{address: trial.address})
		if memoryFaultRegion(err) != trial.region {
			t.Fatal("wrong fixed region classification")
		}
		diagnostic, ok := RecoveryDetails(runtimeFailure("sinf", err, nil))
		encoded, _ := json.Marshal(diagnostic)
		if !ok || diagnostic.Failure != "raw_memory_unavailable" || diagnostic.FaultRegion != trial.region ||
			strings.Contains(string(encoded), fmt.Sprintf("%x", trial.address)) || strings.Contains(err.Error(), fmt.Sprintf("%x", trial.address)) {
			t.Fatal("typed memory fault leaked private address or lost classification")
		}
	}
}
