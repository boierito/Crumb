package main

import (
    "bytes"
    "encoding/json"
    "testing"
)

func TestProbeModesPreservePrivateStateAndChangeOnlyRouting(t *testing.T) {
    sidb := bytes.Repeat([]byte{0x97}, 880)
    lock := []byte{0x33, 0x44}
    domain := map[string][]byte{"/Fixture": {1, 2, 3}, "/Fixture.sinf": {4}}
    state := map[string][]byte{"SC Info.sidb": sidb, "SC Info.sidb.lock": lock}
    for _, test := range []struct {
        mode string
        fpsd bool
        shared bool
        lock bool
    }{
        {"dual", true, true, true},
        {"fpsd-only", true, false, true},
        {"shared-only", false, true, true},
        {"sidb-only", true, true, false},
    } {
        mapping, err := probeRuntimeMappings(domain, state, test.mode)
        if err != nil { t.Fatal(err) }
        if !bytes.Equal(mapping["/Fixture"], domain["/Fixture"]) { t.Fatal("IPA executable altered") }
        for _, alias := range []struct{ prefix string; enabled bool }{
            {"/var/db/fpsd/SC Info/", test.fpsd},
            {"/Users/Shared/SC Info/", test.shared},
        } {
            source, exists := mapping[alias.prefix+"SC Info.sidb"]
            if exists != alias.enabled { t.Fatalf("%s SIDB routing mismatch: %s", test.mode, alias.prefix) }
            if exists && !bytes.Equal(source, sidb) { t.Fatal("SIDB bytes altered") }
            _, presentLock := mapping[alias.prefix+"SC Info.sidb.lock"]
            if presentLock != (alias.enabled && test.lock) { t.Fatal("lock variant mismatch") }
        }
    }
    if _, err := probeRuntimeMappings(domain, state, "unknown"); err == nil {
        t.Fatal("unexpected probe mode accepted")
    }
}

func TestProbeReportHasNoPrivateMaterialFields(t *testing.T) {
    const secret = "PRIVATE-ACCOUNT-KEYBAG-SIDB-HARDWARE"
    report := localProbe{Mode: "dual", Phase: "bootstrap", Code: "raw_unknown_control_flow",
        LastBlock: "text+0xb2f9d8", UntranslatedTarget: "text+0xb30000",
        TargetClass: "in_pinned_code_window"}
    data, err := json.Marshal(report)
    if err != nil || bytes.Contains(data, []byte(secret)) {
        t.Fatal("probe report leaked private input")
    }
    for _, key := range []string{"source", "frame", "pack", "hardware", "sidb", "account", "selector"} {
        if bytes.Contains(data, []byte("\""+key+"\"")) { t.Fatalf("unsafe field %s", key) }
    }
}
