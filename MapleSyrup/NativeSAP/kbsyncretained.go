// Copied into internal/sap/machine: retain the exact KBSync producer for research.
package machine

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"sync"
)

type RetainedKBSyncSession struct {
	mu       sync.Mutex
	hardware []byte
	agent    *StoreAgent
	global   uint32
}

func (r *RetainedKBSyncSession) Close() error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	clear(r.hardware)
	r.hardware = nil
	if r.agent == nil {
		return nil
	}
	if r.agent.guest != nil {
		if fs := r.agent.guest.services.guestStateFS; fs != nil {
			for _, b := range fs.files {
				clear(b)
			}
		}
	}
	err := r.agent.Close()
	r.agent = nil
	return err
}

// ImportKeybag uses the exact wrapper called by the three native keybag handlers.
func (r *RetainedKBSyncSession) ImportKeybag(keybag []byte) (int32, map[string][]byte, map[string]any, error) {
	return r.importLicenseBag(storeAgentBase+0xd5ce0, keybag)
}

// ImportDiversityBag uses the separate routine invoked for the authorization
// response's diversitybag field. A purchase/download keybag is not a substitute.
// Only a real server reply bound to this same retained session may be passed.
func (r *RetainedKBSyncSession) ImportDiversityBag(bag []byte) (int32, map[string][]byte, map[string]any, error) {
	return r.importLicenseBag(storeAgentBase+0xcdea0, bag)
}

func (r *RetainedKBSyncSession) importLicenseBag(entry uint64, keybag []byte) (int32, map[string][]byte, map[string]any, error) {
	if r == nil {
		return 0, nil, nil, fmt.Errorf("invalid consumer input")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.agent == nil || r.agent.guest == nil || len(keybag) == 0 || len(keybag) > 1<<20 {
		return 0, nil, nil, fmt.Errorf("invalid consumer input")
	}
	g := r.agent.guest
	g.beginCall()
	defer g.clearScratch()
	pointer, err := g.scratch(keybag, uint64(len(keybag)))
	if err != nil {
		return 0, nil, nil, err
	}
	status, err := g.invoke(entry, uint64(r.global), pointer, uint64(len(keybag)))
	return int32(status), g.services.guestStateSnapshot(), g.services.guestStateTrace(), err
}
func OpenRetainedKBSync(ctx context.Context, bundle assets.Bundle, hardware []byte, dsid uint64) (*RetainedKBSyncSession, []byte, error) {
	image, err := assets.LoadStoreAgent(ctx)
	if err != nil {
		return nil, nil, err
	}
	agent, global, err := openStoreAgentGlobal(ctx, bundle, image, hardware)
	if err != nil {
		return nil, nil, err
	}
	if err = agent.guest.services.enableGuestState(); err != nil {
		agent.Close()
		return nil, nil, err
	}
	blob, err := agent.guest.generateKBSync(storeAgentKBSyncEntry, global, dsid)
	if err != nil {
		agent.Close()
		return nil, nil, err
	}
	return &RetainedKBSyncSession{agent: agent, global: global, hardware: append([]byte(nil), hardware...)}, blob, nil
}

// OpenRetainedLicenseState restores a previously captured legacy producer
// filesystem before the original native global initializer reads any state.
// It does not generate KBSync, contact Apple, re-import a keybag, convert SIDB,
// or assert that the restored account can decode a mobile SINF.
func OpenRetainedLicenseState(ctx context.Context, bundle assets.Bundle, hardware []byte, files map[string][]byte) (*RetainedKBSyncSession, error) {
	state, err := cloneGuestState(files)
	if err != nil {
		return nil, err
	}
	expectedSIDB := sha256.Sum256(state[storeAgentSCInfoPath+"/SC Info.sidb"])
	owned := true
	defer func() {
		if owned {
			for _, b := range state {
				clear(b)
			}
		}
	}()
	image, err := assets.LoadStoreAgent(ctx)
	if err != nil {
		return nil, err
	}
	agent, global, err := openStoreAgentGlobalPrepared(ctx, bundle, image, hardware, func(agent *StoreAgent) error {
		return agent.guest.services.installGuestState(state)
	})
	if err != nil {
		return nil, err
	}
	if sha256.Sum256(state[storeAgentSCInfoPath+"/SC Info.sidb"]) != expectedSIDB {
		_ = agent.Close()
		return nil, fmt.Errorf("producer_state_not_preserved")
	}
	owned = false
	return &RetainedKBSyncSession{agent: agent, global: global, hardware: append([]byte(nil), hardware...)}, nil
}

// NeedsAccountDiversity invokes the exact routine called by
// -[FairPlayHelper needDiversityForDSID:] in the pinned StoreAgent. Its native
// boolean is a provisioning query, never SINF key confirmation. The
// hardware identity is retained from the global initializer.
func (r *RetainedKBSyncSession) NeedsAccountDiversity(dsid uint64) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("producer_session_invalid")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.agent == nil || r.agent.closed || r.agent.guest == nil || dsid == 0 {
		return false, fmt.Errorf("producer_session_invalid")
	}
	block, err := hardwareBlock(r.hardware)
	if err != nil {
		return false, err
	}
	defer clear(block)
	g := r.agent.guest
	g.beginCall()
	defer g.clearScratch()
	pointer, err := g.scratch(block, uint64(len(block)))
	if err != nil {
		return false, err
	}
	status, err := g.invoke(storeAgentBase+0x11d490, uint64(r.global), pointer, dsid)
	if err != nil {
		return false, err
	}
	return uint32(status) != 0, nil
}

// AuthorizationSync generates mode-2 KBSync as used by the original
// authorizeMachine request, not mode-1 download sync. It does not contact Apple
// or claim enrollment. The returned private buffer belongs to the caller.
func (r *RetainedKBSyncSession) AuthorizationSync(dsid uint64) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("producer_session_invalid")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.agent == nil || r.agent.closed || r.agent.guest == nil || dsid == 0 {
		return nil, fmt.Errorf("producer_session_invalid")
	}
	g := r.agent.guest
	g.beginCall()
	defer g.clearScratch()
	pointerField, err := g.scratch(nil, 8)
	if err != nil {
		return nil, err
	}
	lengthField, err := g.scratch(nil, 8)
	if err != nil {
		return nil, err
	}
	// Verified at text+0x59882..0x5988f in the pinned StoreAgent. The length
	// output is uint32; zeroing its upper half preserves consumeOutput bounds.
	status, err := g.invoke(storeAgentKBSyncEntry, uint64(r.global), dsid, 0, 2, pointerField, lengthField)
	if err != nil {
		return nil, err
	}
	output, outputErr := g.consumeOutput(pointerField, lengthField)
	if int32(status) != 0 || outputErr != nil || len(output) == 0 {
		clear(output)
		return nil, fmt.Errorf("producer_authorization_sync_rejected")
	}
	return output, nil
}

// StateSnapshot returns a private deep copy. It is not a credential or a model.
func (r *RetainedKBSyncSession) StateSnapshot() (map[string][]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("producer_session_invalid")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.agent == nil || r.agent.closed || r.agent.guest == nil {
		return nil, fmt.Errorf("producer_session_invalid")
	}
	return r.agent.guest.services.guestStateSnapshot(), nil
}
