// Copied into internal/sap/machine: retain the exact KBSync producer for research.
package machine

import (
	"context"
	"fmt"
	"github.com/majd/ipatool/v2/internal/sap/assets"
)

type RetainedKBSyncSession struct {
	agent  *StoreAgent
	global uint32
}

func (r *RetainedKBSyncSession) Close() error {
	if r == nil || r.agent == nil {
		return nil
	}
	if r.agent.guest != nil {
		if fs := r.agent.guest.services.guestStateFS; fs != nil {
			for _, b := range fs.files {
				clear(b)
			}
		}
	}
	return r.agent.Close()
}

// ImportKeybag uses the exact wrapper called by the three native keybag handlers.
func (r *RetainedKBSyncSession) ImportKeybag(keybag []byte) (int32, map[string][]byte, map[string]any, error) {
	if r == nil || r.agent == nil || r.agent.guest == nil || len(keybag) == 0 || len(keybag) > 1<<20 {
		return 0, nil, nil, fmt.Errorf("invalid consumer input")
	}
	g := r.agent.guest
	g.beginCall()
	defer g.clearScratch()
	pointer, err := g.scratch(keybag, uint64(len(keybag)))
	if err != nil {
		return 0, nil, nil, err
	}
	status, err := g.invoke(storeAgentBase+0xd5ce0, uint64(r.global), pointer, uint64(len(keybag)))
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
	return &RetainedKBSyncSession{agent: agent, global: global}, blob, nil
}
