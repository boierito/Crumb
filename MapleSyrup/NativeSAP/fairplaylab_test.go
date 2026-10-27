package main

import "testing"

func TestPortableDeviceSelftestContract(t *testing.T) {
	if _, err := labRun(labRequest{Operation: "selftest"}); err != nil {
		t.Fatal(err)
	}
	if _, err := labRun(labRequest{Operation: "decrypt"}); err == nil {
		t.Fatal("missing model accepted")
	}
	if _, err := labRun(labRequest{Operation: "state", Input: make([]byte, 208)}); err == nil {
		t.Fatal("unsupported state size accepted")
	}
	if _, err := labRun(labRequest{Operation: "sidb"}); err == nil {
		t.Fatal("missing private frame accepted")
	}
}

func TestLegacyProducerProbeRequiresNativeImplementation(t *testing.T) {
	if legacyProducerProbe != nil {
		t.Skip("native producer is registered")
	}
	if _, err := labRun(labRequest{Operation: "probe-legacy-producer"}); err == nil {
		t.Fatal("unavailable producer reported success")
	}
}
