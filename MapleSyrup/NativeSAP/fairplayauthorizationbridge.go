package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"github.com/majd/ipatool/v2/internal/sap/machine"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/rawsource"
)

type authorizationRequest struct {
	Operation, Cache, Source, IPA, DSID, GUID, Directory, Destination string
	Handle                                                            uint64
	Keybag, Diversitybag                                              []byte
}
type authorizationSession struct {
	producer                    *machine.RetainedKBSyncSession
	hardware                    []byte
	dsid                        uint64
	domains                     []fairplay.PortableEnrollmentInput
	ipa, directory, destination string
	needs                       bool
	created                     time.Time
}

func (s *authorizationSession) close() {
	s.producer.Close()
	clear(s.hardware)
	for _, domain := range s.domains {
		for _, b := range domain.Files {
			clear(b)
		}
	}
}

var authorizationSessions = map[uint64]*authorizationSession{}

func openAuthorization(r authorizationRequest) (any, error) {
	for handle, s := range authorizationSessions {
		if time.Since(s.created) > 10*time.Minute {
			s.close()
			delete(authorizationSessions, handle)
		}
	}
	if len(authorizationSessions) != 0 {
		return nil, fmt.Errorf("authorization_busy")
	}
	dsid, err := strconv.ParseUint(r.DSID, 10, 64)
	hardware, hexErr := hex.DecodeString(r.GUID)
	if err != nil || dsid == 0 || hexErr != nil || len(hardware) != 6 || r.Cache == "" || r.Directory == "" || r.Destination == "" {
		clear(hardware)
		return nil, fmt.Errorf("authorization_identity_invalid")
	}
	owned := true
	defer func() {
		if owned {
			clear(hardware)
		}
	}()
	domains, err := fairplay.PortableEnrollmentInputs(r.IPA)
	if err != nil || len(domains) == 0 {
		return nil, fmt.Errorf("authorization_ipa_invalid")
	}
	domainOwned := true
	defer func() {
		if domainOwned {
			for _, domain := range domains {
				for _, b := range domain.Files {
					clear(b)
				}
			}
		}
	}()
	var files map[string][]byte
	stored, err := readAuthorizedState(r.Directory, dsid, hardware)
	if err == nil {
		files = stored.Files
		defer stored.clear()
	} else if os.IsNotExist(err) {
		capture, e := readLegacyEnrollmentInput(r.Source)
		if e != nil {
			return nil, fmt.Errorf("authorization_capture_invalid")
		}
		defer capture.clear()
		if capture.dsid != dsid || !bytes.Equal(capture.hardware, hardware) || !capturedSINFMatchesAllDomains(capture.sinfs, domains) {
			return nil, fmt.Errorf("authorization_capture_account_mismatch")
		}
		files = capture.files
	} else {
		return nil, fmt.Errorf("authorization_cached_state_invalid")
	}
	assets.CacheRoot = r.Cache
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	bundle, err := assets.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("authorization_assets_unavailable")
	}
	producer, err := machine.OpenRetainedLicenseState(ctx, bundle, hardware, files)
	if err != nil {
		return nil, fmt.Errorf("authorization_restore_rejected")
	}
	producerOwned := true
	defer func() {
		if producerOwned {
			producer.Close()
		}
	}()
	needs, err := producer.NeedsAccountDiversity(dsid)
	if err != nil {
		return nil, fmt.Errorf("authorization_query_failed")
	}
	var sync []byte
	if needs {
		sync, err = producer.AuthorizationSync(dsid)
		if err != nil {
			return nil, fmt.Errorf("authorization_sync_failed")
		}
	}
	// This is a PRIVATE wire reply, consumed only by the dedicated Swift bridge.
	// It is never passed to the laboratory diagnostic logger.
	nextHandle++
	if nextHandle == 0 {
		clear(sync)
		return nil, fmt.Errorf("authorization_handle_exhausted")
	}
	authorizationSessions[nextHandle] = &authorizationSession{producer: producer, hardware: hardware, dsid: dsid, domains: domains, ipa: r.IPA, directory: r.Directory, destination: r.Destination, needs: needs, created: time.Now()}
	producerOwned = false
	domainOwned = false
	owned = false
	return map[string]any{"handle": nextHandle, "kbsync": sync, "needs_diversity": needs}, nil
}

func finishAuthorization(r authorizationRequest) (any, error) {
	s := authorizationSessions[r.Handle]
	if s == nil {
		return nil, fmt.Errorf("authorization_handle_invalid")
	}
	delete(authorizationSessions, r.Handle)
	defer s.close()
	if time.Since(s.created) > 10*time.Minute {
		return nil, fmt.Errorf("authorization_session_expired")
	}
	if s.needs {
		if (len(r.Keybag) == 0 && len(r.Diversitybag) == 0) || len(r.Keybag) > 1<<20 || len(r.Diversitybag) > 1<<20 {
			return nil, fmt.Errorf("authorization_license_incomplete")
		}
		// The original handler independently skips each absent field. Import
		// only genuine reply bytes; provisioning and SINF gates still follow.
		if len(r.Keybag) != 0 {
			status, snapshot, _, err := s.producer.ImportKeybag(r.Keybag)
			for _, b := range snapshot {
				clear(b)
			}
			if err != nil || status != 0 {
				return nil, fmt.Errorf("authorization_keybag_rejected")
			}
		}
		if len(r.Diversitybag) != 0 {
			status, snapshot, _, err := s.producer.ImportDiversityBag(r.Diversitybag)
			for _, b := range snapshot {
				clear(b)
			}
			if err != nil || status != 0 {
				return nil, fmt.Errorf("authorization_diversity_rejected")
			}
		}
	} else if len(r.Keybag) != 0 || len(r.Diversitybag) != 0 {
		return nil, fmt.Errorf("authorization_unexpected_reply")
	}
	needs, err := s.producer.NeedsAccountDiversity(s.dsid)
	if err != nil {
		return nil, fmt.Errorf("authorization_query_failed")
	}
	files, err := s.producer.StateSnapshot()
	if err != nil {
		return nil, fmt.Errorf("authorization_snapshot_failed")
	}
	defer func() {
		for _, b := range files {
			clear(b)
		}
	}()
	// Persist account state even if model compatibility fails later; no IPA path
	// or per-app key is used as the storage identity. Original capture is intact.
	if err = saveAuthorizedState(s.directory, s.dsid, s.hardware, files); err != nil {
		return nil, fmt.Errorf("authorization_state_save_failed")
	}
	report := map[string]any{"schema": "crumb-account-authorization-v1", "account_state_saved": true, "license_imports_performed": s.needs, "keybag_imported": s.needs && len(r.Keybag) != 0, "diversitybag_imported": s.needs && len(r.Diversitybag) != 0, "needs_account_diversity": needs, "encrypted_domains": len(s.domains), "sinf_validation_attempted": false, "sinf_model_verified": false, "source_frame_verified": false, "completed": false, "plaintext_authenticated": false}
	if needs {
		report["status"] = "diversity_still_required"
		return report, nil
	}
	recovered := make([]rawsource.Result, 0, len(s.domains))
	defer func() {
		for i := range recovered {
			wipeRecovery(&recovered[i])
		}
	}()
	for _, domain := range s.domains {
		value, e := rawsource.RecoverValidated(domain.ExecutablePath, s.hardware, causalMappings(domain, files, files["SC Info.sidb"]), func(pack []byte) bool {
			report["sinf_validation_attempted"] = true
			valid, _ := probeModelSINF(domain.Files, pack)
			return valid
		})
		if e != nil {
			report["status"] = "model_recovery_rejected"
			report["failure"] = labFailureFields("authorization-sinf", e)
			return report, nil
		}
		recovered = append(recovered, value)
		if e = verifyRecoveredEnrollment(value); e != nil {
			report["status"] = "source_frame_rejected"
			return report, nil
		}
		report["sinf_validation_attempted"] = true
		valid, _ := probeModelSINF(domain.Files, value.Pack)
		if !valid {
			report["status"] = "sinf_model_rejected"
			return report, nil
		}
	}
	report["source_frame_verified"] = true
	report["sinf_model_verified"] = true
	result, e := finishRecoveredEnrollments(recovered, s.ipa, s.destination)
	if e != nil {
		report["status"] = "ipa_processing_rejected"
		report["failure"] = labFailureFields("authorization-decrypt", e)
		return report, nil
	}
	report["completed"] = true
	report["status"] = "processed_unverified"
	report["decryption"] = result
	return report, nil
}

//export WaffleFairPlayAuthorization
func WaffleFairPlayAuthorization(input *C.char, output **C.uchar, length *C.size_t) C.int {
	if input == nil || output == nil || length == nil {
		return 1
	}
	raw := []byte(C.GoString(input))
	defer clear(raw)
	if len(raw) > 4<<20 {
		return 1
	}
	var r authorizationRequest
	if json.Unmarshal(raw, &r) != nil {
		return 1
	}
	defer clear(r.Keybag)
	defer clear(r.Diversitybag)
	mutex.Lock()
	defer mutex.Unlock()
	var value any
	var err error
	switch r.Operation {
	case "open":
		value, err = openAuthorization(r)
	case "finish":
		value, err = finishAuthorization(r)
	case "close":
		if s := authorizationSessions[r.Handle]; s != nil {
			s.close()
			delete(authorizationSessions, r.Handle)
		}
		value = map[string]bool{"closed": true}
	default:
		err = fmt.Errorf("authorization_operation_invalid")
	}
	envelope := map[string]any{"ok": err == nil}
	if err != nil {
		envelope["code"] = authorizationErrorCode(err)
	} else {
		envelope["result"] = value
	}
	b, e := json.Marshal(envelope)
	if m, ok := value.(map[string]any); ok {
		if secret, ok := m["kbsync"].([]byte); ok {
			clear(secret)
		}
	}
	if e != nil {
		return 2
	}
	defer clear(b)
	*output = (*C.uchar)(C.CBytes(b))
	*length = C.size_t(len(b))
	return 0
}
func authorizationErrorCode(err error) string {
	switch err.Error() {
	case "authorization_busy", "authorization_identity_invalid", "authorization_ipa_invalid", "authorization_capture_invalid", "authorization_capture_account_mismatch", "authorization_cached_state_invalid", "authorization_assets_unavailable", "authorization_restore_rejected", "authorization_query_failed", "authorization_sync_failed", "authorization_handle_exhausted", "authorization_handle_invalid", "authorization_session_expired", "authorization_license_incomplete", "authorization_keybag_rejected", "authorization_diversity_rejected", "authorization_unexpected_reply", "authorization_snapshot_failed", "authorization_state_save_failed", "authorization_operation_invalid":
		return err.Error()
	}
	return "authorization_failed"
}
