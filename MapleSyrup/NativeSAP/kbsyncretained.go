// Copied into internal/sap/machine: retain the exact KBSync producer for research.
package machine

import (
	"context"
	"github.com/majd/ipatool/v2/internal/sap/assets"
)

func OpenRetainedKBSync(ctx context.Context, bundle assets.Bundle, hardware []byte, dsid uint64) (*StoreAgent, []byte, error) {
	agent, global, err := openStoreAgentGlobal(ctx, bundle, hardware)
	if err != nil {
		return nil, nil, err
	}
	blob, err := agent.guest.generateKBSync(storeAgentKBSyncEntry, global, dsid)
	if err != nil {
		agent.Close()
		return nil, nil, err
	}
	return agent, blob, nil
}
