package rawsource

import "testing"

func TestAOTMemoryFailureCannotBeAttributedToHistoricalMalloc(t *testing.T) {
	value, err := loadModelBytes(rawSourceModel)
	if err != nil {
		t.Fatal(err)
	}
	s := newState(value)
	s.reset()
	s.x[25] = 0x600000000000 // absent private pointer, deliberately not mapped
	runtime := &rawSourceRuntime{lastExternalSlot: "external-slot-007"}
	_, _, err = rawSourceSINFEntryBlock1000f0a34(s, runtime)
	if err == nil {
		t.Fatal("unmapped input accepted")
	}
	diagnostic, ok := RecoveryDetails(runtimeFailure("selection", err, runtime))
	if !ok || diagnostic.Failure != "raw_memory_unavailable" ||
		diagnostic.FaultInstruction != "text+0xf0a3c" || diagnostic.FaultContext != "aot" ||
		diagnostic.FaultExternalSlot != "" || diagnostic.LastExternalSlot != "external-slot-007" {
		t.Fatalf("fault provenance lost or attributed to a previous operation: %+v", diagnostic)
	}
}
