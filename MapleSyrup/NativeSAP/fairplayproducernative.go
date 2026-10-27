package main

import (
	"context"
	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"github.com/majd/ipatool/v2/internal/sap/machine"
	"time"
)

func init() { legacyProducerProbe = probeLegacyProducer }

// This uses the original legacy producer, never a modern enrollment conversion.
// A successful native provisioning query does not validate the SINF or pages.
func probeLegacyProducer(request labRequest) (any, error) {
	input, err := readLegacyEnrollmentInput(request.Source)
	if err != nil {
		return nil, err
	}
	defer input.clear()
	if input.dsid == 0 {
		return nil, recoveryFailure("identity", "capture_dsid_missing", nil)
	}
	domains, err := fairplay.PortableEnrollmentInputs(request.IPA)
	if err != nil {
		return nil, enrollmentInventoryFailure(err)
	}
	defer func() {
		for _, domain := range domains {
			for _, data := range domain.Files {
				clear(data)
			}
		}
	}()
	if !capturedSINFMatchesDomains(input.sinfs, domains) {
		return nil, recoveryFailure("binding", "producer_transaction_mismatch", nil)
	}
	if request.Cache == "" {
		return nil, recoveryFailure("producer", "producer_cache_missing", nil)
	}
	// Shared with SAP/KBSync callers: protect the explicit sandbox cache root and
	// serialize this native operation against concurrent native bridge calls.
	mutex.Lock()
	defer mutex.Unlock()
	assets.CacheRoot = request.Cache
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	bundle, err := assets.Load(ctx)
	if err != nil {
		return nil, recoveryFailure("assets", "producer_assets_unavailable", err)
	}
	producer, err := machine.OpenRetainedLicenseState(ctx, bundle, input.hardware, input.files)
	if err != nil {
		code := "producer_state_restore_failed"
		if err.Error() == "producer_state_not_preserved" {
			code = "producer_state_not_preserved"
		}
		return nil, recoveryFailure("producer-initialization", code, err)
	}
	defer producer.Close()
	needs, err := producer.NeedsAccountDiversity(input.dsid)
	if err != nil {
		return nil, recoveryFailure("producer-account-query", "producer_account_query_failed", err)
	}
	return map[string]any{
		"schema":                              "crumb-legacy-producer-probe-v1",
		"producer":                            "storeagent-10.9",
		"capture_ipa_bound":                   true,
		"state_restored":                      true,
		"sidb_preserved_after_initialization": true,
		"account_query_completed":             true,
		"needs_account_diversity":             needs,
		"sinf_model_verified":                 false,
		"ready_to_decrypt":                    false,
		"plaintext_authenticated":             false,
	}, nil
}
