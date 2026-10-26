package main

import (
    "bytes"
    "encoding/binary"
    "encoding/json"
    "errors"
    "testing"

    "github.com/majd/ipatool/v2/internal/fairplay"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
)

func bridgeFixture() (legacyEnrollmentInput, fairplay.PortableEnrollmentInput) {
    sidb := bytes.Repeat([]byte{0x38}, 882)
    binary.BigEndian.PutUint32(sidb, 0x0007000b)
    return legacyEnrollmentInput{
        hardware: []byte{1,2,3,4,5,6},
        files: map[string][]byte{"SC Info.sidb": sidb, "SC Info.sidb.lock": nil},
    }, fairplay.PortableEnrollmentInput{
        ExecutablePath: "/Fixture.app/Fixture",
        Files: map[string][]byte{"/Fixture": {1,2,3}, "/Fixture.sinf": {1,2}},
    }
}

func TestBridgeAuditDistinguishesLegacyFromModernAndRejectsSINF(t *testing.T) {
    input, domain := bridgeFixture()
    before := bytes.Clone(input.files["SC Info.sidb"])
    hardware := bytes.Clone(input.hardware)
    report := analyzeBridgeProvenance(input, domain, syntheticCausalEngine)
    if !report.SIDBObserved || report.SIDDObserved || report.KnownMacOSGeometry ||
        report.SIDBContract != "macos_v7000b_unclassified_size" ||
        report.SIDBBytes != 882 || report.AccountIdentityBytes != 6 ||
        !report.RuntimeCompleted || !report.SourceFrameVerified || report.SINFModelVerified ||
        report.ModelReady || report.Status != "sinf_rejected" ||
        report.ModelExported || report.PlaintextAuthenticated || report.PrivateStateExported {
        t.Fatalf("legacy state incorrectly promoted: %+v", report)
    }
    if !bytes.Equal(before, input.files["SC Info.sidb"]) || !bytes.Equal(hardware, input.hardware) {
        t.Fatal("bridge audit mutated captured private data")
    }
    encoded, _ := json.Marshal(report)
    for _, secret := range []string{"SC Info.sidb.lock", "hardware_id", "PRIVATE-ACCOUNT", "model_bytes"} {
        if bytes.Contains(encoded, []byte(secret)) { t.Fatal("private details exported") }
    }
}

func TestBridgeAuditFailsClosedOnBadStateAndUnknownRuntime(t *testing.T) {
    input, domain := bridgeFixture()
    input.hardware = input.hardware[:5]
    ran := false
    engine := func(exe string, hw []byte, data map[string][]byte) (rawsource.Result,error) {
        ran = true
        return rawsource.Result{},nil
    }
    invalid := analyzeBridgeProvenance(input,domain,engine)
    if ran || invalid.Status != "inputs_incomplete" || invalid.ModelReady {
        t.Fatal("missing machine identity was accepted")
    }
    input.hardware = []byte{1,2,3,4,5,6}
    rejected := analyzeBridgeProvenance(input,domain,func(string,[]byte,map[string][]byte)(rawsource.Result,error){
        return rawsource.Result{},errors.New("PRIVATE-ACCOUNT-BLOB: should never be in logs")
    })
    if rejected.Status != "runtime_rejected" || rejected.ModelReady || rejected.UnmodeledImport != "" {
        t.Fatalf("rejection handling: %+v", rejected)
    }
    out,_ := json.Marshal(rejected)
    if bytes.Contains(out, []byte("PRIVATE-ACCOUNT")) { t.Fatal("runtime private error leaked") }
}

func TestStaticExternalImportLabelIsStrict(t *testing.T) {
    for _, item := range []struct{ in, want string }{
        {"external-slot-00a", "external-slot-00a"},
        {"external-slot-fff", "external-slot-fff"},
        {"unclassified_external", "unclassified_external"},
        {"external-slot-01000", ""},
        {"external-slot-zz1", ""},
        {"PRIVATE-ACCOUNT-KEY", ""},
    } {
        if got:=auditedImportSlot(item.in); got!=item.want {
            t.Fatalf("%q: %q != %q",item.in,got,item.want)
        }
    }
    report := bridgeAudit{Schema:"fixture"}
    bridgeRuntimeLookups(&report,rawsource.RecoveryDiagnostic{
        StateLookupsByPhase:map[string]rawsource.LookupCount{
            "bootstrap":{Attempts:3,Missing:1},
            "PRIVATE-ACCOUNT":{Attempts:1},
        },
        StateFileTypes:map[string]rawsource.LookupCount{"sidb":{Attempts:4},"SECRET":{Attempts:10}},
    })
    if len(report.StateLookupsByPhase)!=1 || len(report.StateFileLookups)!=1 {
        t.Fatal("untrusted diagnostic category leaked")
    }
}
