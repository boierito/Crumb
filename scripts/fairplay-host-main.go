// Built with the production portable bridge by build-fairplay-host.sh.
package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

func main() {
	ipa := flag.String("ipa", "", "Existing encrypted IPA")
	capture := flag.String("capture", "", "Matching private capture ZIP")
	output := flag.String("output", "", "New directory for validated output only")
	accountState := flag.String("authorized-state", "", "Optional private authorized-state JSON from the same Crumb account")
	flag.Parse()
	if *ipa == "" || *capture == "" || *output == "" || flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}
	report := map[string]any{"schema": "crumb-host-integration-v1", "scope": "production-portable-capture-replay", "apple_authorization_executed": false, "completed": false, "plaintext_authenticated": false}
	report["authorized_state_supplied"] = *accountState != ""
	stages := []map[string]any{}
	failed := false
	for _, request := range []labRequest{
		{Operation: "selftest"},
		{Operation: "inspect", Source: *ipa},
		{Operation: "analyze-capture", Source: *capture},
		{Operation: "verify-capture-ipa", Source: *capture, IPA: *ipa},
		{Operation: "recover-capture", Source: *capture, IPA: *ipa, Destination: *output},
	} {
		start := time.Now()
		var result any
		var err error
		if request.Operation == "recover-capture" && *accountState != "" {
			result, err = recoverCaptureWithAccountState(*capture, *ipa, *output, *accountState)
		} else {
			result, err = labRun(request)
		}
		stage := map[string]any{"operation": request.Operation, "milliseconds": time.Since(start).Milliseconds(), "ok": err == nil}
		if err != nil {
			stage["failure"] = labFailureFields(request.Operation, err)
			failed = true
		} else {
			stage["result"] = result
		}
		stages = append(stages, stage)
		if failed {
			break
		}
		if request.Operation == "recover-capture" {
			if fields, ok := result.(map[string]any); ok && fields["completed"] == true {
				report["completed"] = true
			}
		}
	}
	report["stages"] = stages
	if json.NewEncoder(os.Stdout).Encode(report) != nil {
		os.Exit(2)
	}
	if failed {
		os.Exit(1)
	}
}
