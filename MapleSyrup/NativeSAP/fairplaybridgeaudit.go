package main

// A read-only, generation-aware producer/consumer bridge audit. Unlike blind
// SIDB mutations, this asks which state was actually provided by StoreAgent,
// which files the modern fairplayd interpreter touched, and whether the
// original target SINF validates. No model, logical key or private blob exits.
import (
    "strings"

    "github.com/majd/ipatool/v2/internal/fairplay"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
)

type bridgeAudit struct {
    Schema string `json:"schema"`
    Producer string `json:"producer"`
    Consumer string `json:"consumer"`
    SIDBContract string `json:"sidb_contract"`
    SIDBBytes int `json:"sidb_bytes"`
    KnownMacOSGeometry bool `json:"known_macos_geometry"`
    SIDDObserved bool `json:"sidd_observed"`
    SIDBObserved bool `json:"sidb_observed"`
    AccountIdentityBytes int `json:"identity_bytes"`
    NativeFiles int `json:"native_file_count"`
    RuntimeCompleted bool `json:"runtime_completed"`
    FailurePhase string `json:"failure_phase,omitempty"`
    FailureCode string `json:"failure_code,omitempty"`
    UnmodeledImport string `json:"unmodeled_import,omitempty"`
    StateLookupsByPhase map[string]rawsource.LookupCount `json:"state_lookups_by_phase,omitempty"`
    StateFileLookups map[string]rawsource.LookupCount `json:"state_file_lookups,omitempty"`
    SourceFrameVerified bool `json:"source_frame_verified"`
    SINFModelVerified bool `json:"sinf_model_verified"`
    ModelReady bool `json:"model_ready"`
    Status string `json:"status"`
    ModelExported bool `json:"model_exported"`
    PrivateStateExported bool `json:"private_state_exported"`
    PlaintextAuthenticated bool `json:"plaintext_authenticated"`
}

// Only the architecture's public virtual-import slot may be exposed. Random
// program counters, arbitrary model errors and guest pointers stay private.
func auditedImportSlot(value string) string {
    if value == "unclassified_external" { return value }
    const prefix = "external-slot-"
    if !strings.HasPrefix(value, prefix) || len(value) != len(prefix)+3 { return "" }
    for _, b := range value[len(prefix):] {
        if !strings.ContainsRune("0123456789abcdef", b) { return "" }
    }
    return value
}

func bridgeRuntimeLookups(report *bridgeAudit, diag rawsource.RecoveryDiagnostic) {
    report.StateLookupsByPhase = make(map[string]rawsource.LookupCount)
    for phase, data := range diag.StateLookupsByPhase {
        switch phase {
        case "account-parser", "bootstrap", "source", "selection", "sinf":
            report.StateLookupsByPhase[phase] = data
        }
    }
    report.StateFileLookups = make(map[string]rawsource.LookupCount)
    for kind, data := range diag.StateFileTypes {
        switch kind {
        case "sidb", "sidd", "sidf", "other":
            report.StateFileLookups[kind] = data
        }
    }
}

// engine is injected so the bridge contract can be tested without account
// material, iOS system calls, JIT, or any vendor cryptographic fixture.
func analyzeBridgeProvenance(input legacyEnrollmentInput, domain fairplay.PortableEnrollmentInput, engine causalRecovery) bridgeAudit {
    sidb := input.files["SC Info.sidb"]
    report := bridgeAudit{
        Schema: "crumb-storeagent-fairplayd-bridge-v1",
        Producer: "storeagent-10.9-import-keybag",
        Consumer: "fairplayd-26a428-recovery",
        SIDBContract: classifyCanonicalSIDB(sidb),
        SIDBBytes: len(sidb),
        SIDBObserved: len(sidb) > 0,
        SIDDObserved: len(input.files["SC Info.sidd"]) > 0,
        AccountIdentityBytes: len(input.hardware),
        NativeFiles: len(input.files),
        Status: "inputs_incomplete",
    }
    switch report.SIDBContract {
    case "macos_v7000b_known_192", "macos_v7000b_known_224":
        report.KnownMacOSGeometry = true
    }
    if len(input.hardware) != 6 || len(sidb) < 4 || domain.ExecutablePath == "" ||
        len(domain.Files) == 0 || engine == nil {
        return report
    }
    value, err := engine(domain.ExecutablePath, input.hardware, causalMappings(domain, input.files, sidb))
    defer wipeRecovery(&value)
    if err != nil {
        report.Status = "runtime_rejected"
        if diag, ok := rawsource.RecoveryDetails(err); ok {
            report.FailurePhase = causalFailurePhase(diag.Phase)
            report.FailureCode = causalFailureCode(diag.Failure)
            report.UnmodeledImport = auditedImportSlot(diag.UnmodeledExternalSlot)
            bridgeRuntimeLookups(&report, diag)
        }
        return report
    }
    report.RuntimeCompleted = true
    bridgeRuntimeLookups(&report, value.Runtime)
    report.SourceFrameVerified = verifyRecoveredEnrollment(value) == nil
    report.SINFModelVerified, _ = probeModelSINF(domain.Files, value.Pack)
    report.ModelReady = report.SourceFrameVerified && report.SINFModelVerified
    report.Status = "sinf_rejected"
    if report.ModelReady { report.Status = "cryptographic_model_validated" }
    return report
}

func researchCaptureBridge(source, ipa string) (any, error) {
    input, err := readLegacyEnrollmentInput(source)
    if err != nil { return nil, recoveryFailure("bridge-input", "bridge_capture_invalid", err) }
    defer input.clear()
    domains, err := fairplay.PortableEnrollmentInputs(ipa)
    if err != nil { return nil, enrollmentInventoryFailure(err) }
    defer func() {
        for _, domain := range domains {
            for _, data := range domain.Files { clear(data) }
        }
    }()
    if len(domains) != 1 || !input.sinfs[domains[0].SINFSHA256] {
        return nil, recoveryFailure("bridge-binding", "bridge_transaction_mismatch", nil)
    }
    report := analyzeBridgeProvenance(input, domains[0], rawsource.Recover)
    return report, nil
}
