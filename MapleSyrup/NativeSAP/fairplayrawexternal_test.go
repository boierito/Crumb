package rawsource

import "testing"

// Public synthetic import-table entries from pinned 26A428. Never export
// arbitrary guest registers or memory addresses as diagnostics.
func TestExternalImportSlotClassification(t *testing.T) {
    cases := []struct{ pc uint64; want string }{
        {0xf100000000006000, "external-slot-006"}, // free
        {0xf10000000004e000, "external-slot-04e"}, // IOKit CF property
        {0xf1000000000e4000, "external-slot-0e4"}, // chkstk
        {0xf10000000002500c, "unclassified_external"}, // callback, not aligned slot
        {0x720000000000, "unclassified_external"},
        {0x100000000, "unclassified_external"},
        {0xffffffffffffffff, "unclassified_external"},
        {0xf100000001000000, "unclassified_external"},
    }
    for _, tt := range cases {
        if got:=sourceExternalImportSlot(tt.pc);got!=tt.want {
            t.Fatalf("static slot %#x classified %q, wanted %q",tt.pc,got,tt.want)
        }
    }
}

func TestExternalAndStatePhaseDiagnosticsUseOnlySanitizedFields(t *testing.T) {
    r := &rawSourceRuntime{
        unmodeledExternalSlot:"external-slot-04e",
        stateLookupsByPhase:map[string]LookupCount{"selection":{Attempts:3,Missing:1}},
    }
    diag:=runtimeDiagnostic("selection","raw_unsupported_external",r)
    if diag.UnmodeledExternalSlot!="external-slot-04e" ||
        diag.StateLookupsByPhase["selection"].Attempts!=3 {
        t.Fatalf("safe runtime fields not preserved: %+v",diag)
    }
    r.stateLookupsByPhase["selection"]=LookupCount{Attempts:999}
    if diag.StateLookupsByPhase["selection"].Attempts!=3 {
        t.Fatal("diagnostic contains mutable runtime alias")
    }
}
