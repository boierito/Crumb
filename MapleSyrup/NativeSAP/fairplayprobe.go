package main

// A bounded, non-destructive comparison of genuine on-device StoreAgent state.
// Every variant uses the SAME unmodified private bytes from one capture.
// Only virtual file routing changes. No models, plaintext or secret material
// are written to disk or returned to Swift.
import (
    "fmt"
    "time"

    "github.com/majd/ipatool/v2/internal/fairplay"
    "github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
)

type localProbe struct {
    Mode string `json:"mode"`
    Phase string `json:"phase"`
    Code string `json:"code"`
    LastBlock string `json:"last_translated_block,omitempty"`
    UntranslatedTarget string `json:"untranslated_target,omitempty"`
    TargetClass string `json:"target_class,omitempty"`
    SharedBlocks int `json:"shared_block_executions"`
    FileLookups map[string]rawsource.LookupCount `json:"file_lookups,omitempty"`
    SourceFrameVerified bool `json:"source_frame_verified"`
    DurationMs int64 `json:"milliseconds"`
}

var localProbeModes = []string{"dual", "fpsd-only", "shared-only", "sidb-only"}

// Routing, not mutation: retains opaque bytes and preserved device identity.
// Every variant is based on exactly the same per-domain IPA file inventory.
func probeRuntimeMappings(domainFiles map[string][]byte, stateFiles map[string][]byte, mode string) (map[string][]byte, error) {
    result := make(map[string][]byte, len(domainFiles)+len(stateFiles)*2)
    for name, data := range domainFiles {
        result[name] = data
    }
    switch mode {
    case "dual", "fpsd-only", "shared-only", "sidb-only":
    default:
        return nil, fmt.Errorf("unsupported probe mode")
    }
    for name, data := range stateFiles {
        if mode == "sidb-only" && name != "SC Info.sidb" {
            continue
        }
        if mode == "dual" || mode == "fpsd-only" || mode == "sidb-only" {
            result["/var/db/fpsd/SC Info/"+name] = data
        }
        if mode == "dual" || mode == "shared-only" || mode == "sidb-only" {
            result["/Users/Shared/SC Info/"+name] = data
        }
    }
    return result, nil
}

func probeCaptureEnrollment(source, ipa string) (any, error) {
    input, err := readLegacyEnrollmentInput(source)
    if err != nil {
        return nil, err
    }
    defer input.clear()
    domains, err := fairplay.PortableEnrollmentInputs(ipa)
    if err != nil {
        return nil, enrollmentInventoryFailure(err)
    }
    defer func() {
        for _, domain := range domains {
            for _, data := range domain.Files { clear(data) }
        }
    }()
    if len(domains) == 0 {
        return nil, recoveryFailure("ipa-inventory", "ipa_sinf_missing", nil)
    }
    for _, domain := range domains {
        if !input.sinfs[domain.SINFSHA256] {
            return nil, recoveryFailure("transaction-binding", "capture_ipa_mismatch", nil)
        }
    }

    // Bound the experiment on a mobile device. A large framework-rich IPA
    // cannot cause unbounded repeated recovery attempts.
    count := len(domains)
    if count > 2 { count = 2 }
    deadline := time.Now().Add(45 * time.Second)
    reports := make([]map[string]any, 0, count)
    attempted := 0
    truncated := false
    for index := 0; index < count; index++ {
        domain := domains[index]
        attempts := make([]localProbe, 0, len(localProbeModes))
        for _, mode := range localProbeModes {
            if time.Now().After(deadline) {
                truncated = true
                break
            }
            mapping, e := probeRuntimeMappings(domain.Files, input.files, mode)
            if e != nil { return nil, recoveryFailure("probe-configuration", "probe_mode_invalid", e) }
            t0 := time.Now()
            value, recoveryErr := rawsource.Recover(domain.ExecutablePath, input.hardware, mapping)
            result := localProbe{Mode: mode, Phase: "portable-runtime", Code: "unverified", DurationMs: time.Since(t0).Milliseconds()}
            if recoveryErr != nil {
                if details, found := rawsource.RecoveryDetails(recoveryErr); found {
                    result.Phase = details.Phase
                    result.Code = details.Failure
                    result.LastBlock = details.LastBlock
                    result.UntranslatedTarget = details.UntranslatedTarget
                    result.TargetClass = details.TargetClass
                    result.SharedBlocks = details.SharedBlockExecutions
                    result.FileLookups = details.Lookups
                } else {
                    result.Code = "raw_operation_rejected"
                }
            } else if verifyRecoveredEnrollment(value) == nil {
                result.Phase = "model-validation"
                result.Code = "source_frame_verified"
                result.SourceFrameVerified = true
            } else {
                result.Phase = "model-validation"
                result.Code = "source_frame_mismatch"
            }
            clear(value.Source)
            clear(value.Frame)
            clear(value.Pack)
            attempts = append(attempts, result)
            attempted++
        }
        reports = append(reports, map[string]any{"domain_index": index, "probes": attempts})
        if truncated { break }
    }
    return map[string]any{
        "schema": "crumb-fairplay-local-probes-v1",
        "synthetic": false,
        "private_state_exported": false,
        "model_exported": false,
        "domain_count": len(domains),
        "domains_sampled": len(reports),
        "attempts": attempted,
        "truncated_at_time_budget": truncated,
        "mode_definitions": []string{"dual state aliases", "fpsd state only", "shared state only", "SIDB alone at both aliases"},
        "results": reports,
        "plaintext_authenticated": false,
    }, nil
}
