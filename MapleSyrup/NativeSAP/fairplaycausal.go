package main

// Research-only differential experiment. This never decrypts an IPA, exports
// enrollment material, logs private identities, or promotes a candidate.
// Every perturbation acts on a fresh private copy of the captured state.
import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"

    "github.com/majd/ipatool/v2/internal/fairplay"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
)

// These exact version/size pairs come from UnfairPlay's canonicalSIDB()
// (DrWhax/unfairplay, source commit 12bf939d38f28279de14cb5d318257f2092f2daf).
// A match is only a *format hint*, never evidence of compatible enrollment.
func classifyCanonicalSIDB(blob []byte) string {
    if len(blob) < 4 { return "unclassified" }
    if binary.BigEndian.Uint32(blob[:4]) != 0x0007000b {
        return "other_header"
    }
    switch len(blob) {
    case 630, 642:
        return "macos_v7000b_known_192"
    case 993, 1003:
        return "macos_v7000b_known_224"
    default:
        return "macos_v7000b_unclassified_size"
    }
}

type causalTrial struct {
    Intervention string `json:"intervention"`
    Status string `json:"status"`
    FailurePhase string `json:"failure_phase,omitempty"`
    FailureCode string `json:"failure_code,omitempty"`
    SourceChanged bool `json:"source_changed"`
    FrameChanged bool `json:"frame_changed"`
    ModelChanged bool `json:"model_changed"`
    SINFModelVerified bool `json:"sinf_model_verified"`
}

type causalReport struct {
    Schema string `json:"schema"`
    SIDBContract string `json:"sidb_contract"`
    SIDBBytes int `json:"sidb_bytes"`
    ReplayVerified bool `json:"replay_verified"`
    SourceFrameVerified bool `json:"source_frame_verified"`
    SINFModelVerified bool `json:"sinf_model_verified"`
    Trials []causalTrial `json:"trials"`
    Experimental bool `json:"experimental"`
    OriginalIPARetained bool `json:"original_ipa_retained"`
    PrivateStateExported bool `json:"private_state_exported"`
    ModelExported bool `json:"model_exported"`
    PlaintextAuthenticated bool `json:"plaintext_authenticated"`
}

type causalRecovery func(string, []byte, map[string][]byte) (rawsource.Result, error)

func causalMappings(domain fairplay.PortableEnrollmentInput, files map[string][]byte, sidb []byte) map[string][]byte {
    mapping := make(map[string][]byte, len(domain.Files)+2*len(files))
    for k, v := range domain.Files { mapping[k] = v }
    for name, data := range files {
        if name == "SC Info.sidb" { data = sidb }
        mapping["/var/db/fpsd/SC Info/"+name] = data
        mapping["/Users/Shared/SC Info/"+name] = data
    }
    return mapping
}

func wipeRecovery(value *rawsource.Result) {
    clear(value.Source)
    clear(value.Frame)
    clear(value.Pack)
    *value = rawsource.Result{}
}

func causalPositions(length int) []struct {label string; offset int} {
    proposals := []struct{ label string; offset int }{
        {"sidb_first", 0}, {"sidb_header_end", 4},
        {"sidb_middle", length/2}, {"sidb_last", length-1},
    }
    seen := map[int]bool{}
    result := make([]struct {label string; offset int}, 0, len(proposals))
    for _, p := range proposals {
        if p.offset < 0 || p.offset >= length || seen[p.offset] { continue }
        seen[p.offset] = true
        result = append(result, p)
    }
    return result
}

// Engine injection allows publicly reproducible *causal* unit tests. The
// production caller always passes the real rawsource.Recover and matches
// SINF/IPA identity against the captured private transaction first.
func studySIDBCausality(input legacyEnrollmentInput, domain fairplay.PortableEnrollmentInput, engine causalRecovery) (causalReport, error) {
    sidb := input.files["SC Info.sidb"]
    report := causalReport{
        Schema: "crumb-sidb-causality-v1", SIDBContract: classifyCanonicalSIDB(sidb),
        SIDBBytes: len(sidb), Experimental: true, OriginalIPARetained: true,
    }
    if len(sidb) < 16 || len(sidb) > 4<<20 || len(input.hardware) != 6 || domain.ExecutablePath == "" ||
        len(domain.Files) == 0 || engine == nil {
        return report, errors.New("invalid private research inputs")
    }
    baseline, err := engine(domain.ExecutablePath, input.hardware, causalMappings(domain, input.files, sidb))
    if err != nil {
        wipeRecovery(&baseline)
        return report, fmt.Errorf("original enrollment research execution rejected: %w", err)
    }
    defer wipeRecovery(&baseline)
    report.SourceFrameVerified = verifyRecoveredEnrollment(baseline) == nil
    report.SINFModelVerified, _ = probeModelSINF(domain.Files, baseline.Pack)
    type stimulus struct{ name string; offset int; identity bool }
    stimuli := make([]stimulus, 0, 5)
    for _, item := range causalPositions(len(sidb)) {
        stimuli = append(stimuli, stimulus{name: item.label, offset: item.offset})
    }
    stimuli = append(stimuli, stimulus{name: "hardware_first", offset: 0, identity: true})
    // A deterministic, unmodified replay distinguishes state corruption or
    // identity binding from a failure to reset the interpreter between runs.
    // If this control fails, mutation conclusions are not interpretable.
    trials := append([]stimulus{{name: "baseline_replay", offset: -1}}, stimuli...)
    for _, test := range trials {
        privateSIDB := bytes.Clone(sidb)
        privateHardware := bytes.Clone(input.hardware)
        if test.identity {
            privateHardware[test.offset] ^= 1
        } else if test.offset >= 0 {
            privateSIDB[test.offset] ^= 1
        }
        changed, runErr := engine(domain.ExecutablePath, privateHardware, causalMappings(domain, input.files, privateSIDB))
        entry := causalTrial{Intervention: test.name, Status: "runtime_rejected"}
        if runErr != nil {
            if details, ok := rawsource.RecoveryDetails(runErr); ok {
                // The runtime contract exposes only typed, fixed error
                // classifications, never a private filename, key or address.
                entry.FailurePhase = causalFailurePhase(details.Phase)
                entry.FailureCode = causalFailureCode(details.Failure)
            }
        } else {
            entry.Status = "completed"
            entry.SourceChanged = !bytes.Equal(baseline.Source, changed.Source)
            entry.FrameChanged = !bytes.Equal(baseline.Frame, changed.Frame)
            entry.ModelChanged = !bytes.Equal(baseline.Pack, changed.Pack)
            entry.SINFModelVerified, _ = probeModelSINF(domain.Files, changed.Pack)
        }
        report.Trials = append(report.Trials, entry)
        // If identical inputs cannot be replayed, any later mutation
        // experiment is causally uninterpretable. Never present five
        // synthetic-looking rejections as evidence of input sensitivity.
        if test.name == "baseline_replay" {
            report.ReplayVerified = entry.Status == "completed" &&
                !entry.SourceChanged && !entry.FrameChanged &&
                !entry.ModelChanged &&
                entry.SINFModelVerified == report.SINFModelVerified
        }
        wipeRecovery(&changed)
        clear(privateSIDB)
        clear(privateHardware)
        if test.name == "baseline_replay" && !report.ReplayVerified {
            return report, nil
        }
    }
    return report, nil
}

// Preserve only known, fixed tags from the enclave's public diagnostic ABI.
// Never pass through Error(), guest paths, or arbitrary strings from sources.
func causalFailurePhase(value string) string {
    switch value {
    case "inputs", "resource", "configure", "account-parser",
        "bootstrap", "source", "selection", "selection-output",
        "sinf", "pack":
        return value
    default:
        return "other"
    }
}

func causalFailureCode(value string) string {
    switch value {
    case "raw_invalid_inputs", "raw_resource_invalid", "raw_binary_mismatch",
        "raw_setup_failed", "raw_input_size_limit", "raw_memory_unavailable",
        "raw_unsupported_external", "raw_unknown_control_flow", "raw_step_limit",
        "raw_unsupported_operation", "raw_operation_rejected",
        "raw_selection_output_unavailable", "raw_pack_unavailable",
        "generic_execution_rejected", "generic_instruction_limit",
        "generic_pc_outside_profile", "generic_profile_missing",
        "generic_runtime_missing", "generic_unsupported_instruction":
        return value
    default:
        return "other"
    }
}

// No IPA download or consumer invocation. The caller must already have a
// matching, verified capture and encrypted IPA saved in the app sandbox.
func researchCaptureCausality(source, ipa string) (any, error) {
    input, err := readLegacyEnrollmentInput(source)
    if err != nil { return nil, recoveryFailure("causal-input", "causal_capture_invalid", err) }
    defer input.clear()
    domains, err := fairplay.PortableEnrollmentInputs(ipa)
    if err != nil { return nil, enrollmentInventoryFailure(err) }
    defer func() {
        for _, domain := range domains {
            for _, data := range domain.Files { clear(data) }
        }
    }()
    if len(domains) != 1 || !input.sinfs[domains[0].SINFSHA256] {
        return nil, recoveryFailure("causal-binding", "causal_transaction_mismatch", nil)
    }
    report, err := studySIDBCausality(input, domains[0], rawsource.Recover)
    if err != nil { return nil, recoveryFailure("causal-execution", "causal_runtime_failed", err) }
    return report, nil
}
