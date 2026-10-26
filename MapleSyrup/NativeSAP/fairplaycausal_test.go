package main

import (
    "bytes"
    "encoding/binary"
    "encoding/json"
    "errors"
    "testing"

    "github.com/majd/ipatool/v2/internal/fairplay"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/state"
)

func TestUnfairPlayCanonicalSIDBClassificationIsConservative(t *testing.T) {
    for _, item := range []struct{ n int; magic uint32; want string }{
        {630, 0x0007000b, "macos_v7000b_known_192"},
        {642, 0x0007000b, "macos_v7000b_known_192"},
        {993, 0x0007000b, "macos_v7000b_known_224"},
        {1003, 0x0007000b, "macos_v7000b_known_224"},
        {863, 0x0007000b, "macos_v7000b_unclassified_size"},
        {869, 0x0006000a, "other_header"},
        {0, 0, "unclassified"},
    } {
        data := make([]byte, item.n)
        if len(data) >= 4 { binary.BigEndian.PutUint32(data, item.magic) }
        if got := classifyCanonicalSIDB(data); got != item.want {
            t.Fatalf("n=%d, magic=%x: got %s want %s", item.n, item.magic, got, item.want)
        }
    }
}

func syntheticCausalEngine(_ string, hardware []byte, files map[string][]byte) (rawsource.Result, error) {
    sidb := files["/var/db/fpsd/SC Info/SC Info.sidb"]
    if len(sidb) == 0 || !bytes.Equal(sidb, files["/Users/Shared/SC Info/SC Info.sidb"]) {
        return rawsource.Result{}, errors.New("missing side-by-side synthetic source")
    }
    source := make([]byte, 224)
    for i, v := range sidb { source[i%len(source)] ^= v }
    for i, v := range hardware { source[i] ^= v }
    transformed, err := state.Evaluate(source, nil)
    if err != nil { return rawsource.Result{}, err }
    defer clear(transformed)
    selected, err := frame.Evaluate(transformed[:16], nil)
    if err != nil { return rawsource.Result{}, err }
    pack := bytes.Clone(source)
    return rawsource.Result{Source: source, Frame: selected, Pack: pack}, nil
}

func TestCausalStudyMutatesPrivateCopiesAndNeverReportsMaterial(t *testing.T) {
    const secret = "SECRET-CAPTURE-DATA-DO-NOT-EXPOSE"
    original := bytes.Repeat([]byte{0x5e}, 863)
    // The original is allowed to contain arbitrary real-format opaque bytes.
    copyOriginal := bytes.Clone(original)
    hardware := []byte{1, 2, 3, 4, 5, 6}
    copyHardware := bytes.Clone(hardware)
    input := legacyEnrollmentInput{
        hardware: hardware,
        files: map[string][]byte{"SC Info.sidb": original, "SC Info.sidb.lock": nil},
    }
    domain := fairplay.PortableEnrollmentInput{
        ExecutablePath: "/Payload/Fixture.app/Fixture",
        Files: map[string][]byte{"/Fixture": []byte(secret), "/Fixture.sinf": {1, 2, 3}},
    }
    report, err := studySIDBCausality(input, domain, syntheticCausalEngine)
    if err != nil { t.Fatal(err) }
    if !bytes.Equal(original, copyOriginal) || !bytes.Equal(hardware, copyHardware) {
        t.Fatal("causal mutation altered the original enrollment")
    }
    if len(report.Trials) != 5 || report.SIDBContract != "other_header" ||
        !report.SourceFrameVerified || report.SINFModelVerified ||
        report.PrivateStateExported || report.ModelExported || report.PlaintextAuthenticated ||
        !report.OriginalIPARetained {
        t.Fatalf("unexpected research summary: %+v", report)
    }
    for _, trial := range report.Trials {
        if trial.Status != "completed" || !trial.SourceChanged || !trial.FrameChanged ||
            !trial.ModelChanged || trial.SINFModelVerified {
            t.Fatalf("synthetic input mutation did not alter all downstream material: %+v", trial)
        }
    }
    encoded, err := json.Marshal(report)
    if err != nil { t.Fatal(err) }
    for _, private := range []string{secret, "SECRET", "hardware_id", "model_bytes", "frame_bytes", "account_id"} {
        if bytes.Contains(encoded, []byte(private)) { t.Fatal("private material leaked to the research summary") }
    }
}

func TestCausalStudyLabelsRuntimeRejectionWithoutPrivateCause(t *testing.T) {
    input := legacyEnrollmentInput{hardware: []byte{1, 2, 3, 4, 5, 6},
        files: map[string][]byte{"SC Info.sidb": bytes.Repeat([]byte{0x12}, 863)}}
    domain := fairplay.PortableEnrollmentInput{ExecutablePath: "/Fixture",
        Files: map[string][]byte{"Fixture.sinf": {0}}}
    engine := func(path string, mac []byte, files map[string][]byte) (rawsource.Result, error) {
        if files["/var/db/fpsd/SC Info/SC Info.sidb"][0] != 0x12 {
            return rawsource.Result{}, errors.New("PRIVATE-ACCOUNT-KEY-DO-NOT-LOG")
        }
        return syntheticCausalEngine(path, mac, files)
    }
    report, err := studySIDBCausality(input, domain, engine)
    if err != nil { t.Fatal(err) }
    if report.Trials[0].Status != "runtime_rejected" || report.Trials[0].SINFModelVerified {
        t.Fatalf("runtime failure unexpectedly accepted: %+v", report.Trials[0])
    }
    encoded, _ := json.Marshal(report)
    if bytes.Contains(encoded, []byte("PRIVATE-ACCOUNT")) {
        t.Fatal("raw private exception leaked")
    }
}
