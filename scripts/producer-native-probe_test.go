package machine

import (
	"bytes"
	"context"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"testing"
	"time"
)

// Exercises actual hash-verified Apple code through the host TCI interpreter.
// DSID 1 and the locally administered MAC are public synthetic inputs. There
// is no Apple login, entitlement request, private keybag or enrolled account.
func TestProducerNativeUnenrolledAndRejectedState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	bundle, err := assets.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	image, err := assets.LoadStoreAgent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	hardware := []byte{2, 0, 0, 0, 0, 1}
	agent, global, err := openStoreAgentGlobalPrepared(ctx, bundle, image, hardware, func(a *StoreAgent) error { return a.guest.services.enableGuestState() })
	if err != nil {
		t.Fatal(err)
	}
	session := &RetainedKBSyncSession{agent: agent, global: global, hardware: bytes.Clone(hardware)}
	defer session.Close()
	needs, err := session.NeedsAccountDiversity(1)
	if err != nil {
		t.Fatal(err)
	}
	if !needs {
		t.Fatal("empty synthetic account unexpectedly provisioned")
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := session.NeedsAccountDiversity(1); err == nil {
		t.Fatal("closed native context used")
	}
	// A header is not an enrollment. The original initializer deletes this
	// truncated record; restoration must reject the reset-to-empty result.
	input := map[string][]byte{"SC Info.sidb": {0, 7, 0, 11}}
	restored, err := OpenRetainedLicenseState(ctx, bundle, hardware, input)
	if restored != nil {
		restored.Close()
		t.Fatal("truncated state restored")
	}
	if err == nil || err.Error() != "producer_state_not_preserved" {
		t.Fatalf("expected native discard rejection, got %v", err)
	}
	if !bytes.Equal(input["SC Info.sidb"], []byte{0, 7, 0, 11}) {
		t.Fatal("capture changed")
	}
}

// Validate the newly identified authorization primitives without sending an
// authorization request or using a real account. Native rejection is distinct
// from an interpreter failure and can never be treated as enrollment.
func TestProducerNativeAuthorizationPrimitives(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	bundle, err := assets.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	producer, downloadSync, err := OpenRetainedKBSync(ctx, bundle, []byte{2, 0, 0, 0, 0, 1}, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	defer clear(downloadSync)
	sync, err := producer.AuthorizationSync(1)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(sync)
	if len(sync) == 0 {
		t.Fatal("native authorization sync empty")
	}
	t.Logf("native authorization sync bytes=%d; contents withheld; no authorization request sent", len(sync))
	status, files, _, err := producer.ImportDiversityBag([]byte("PUBLIC invalid diversity response"))
	defer func() {
		for _, b := range files {
			clear(b)
		}
	}()
	if err != nil {
		t.Fatalf("interpreter failed before native rejection: %v", err)
	}
	if status == 0 {
		t.Fatal("invalid diversity response accepted")
	}
	t.Logf("native invalid diversity response rejected with status=%d", status)
	if err := producer.Close(); err != nil {
		t.Fatal(err)
	}
	if output, err := producer.AuthorizationSync(1); err == nil || output != nil {
		t.Fatal("closed producer generated authorization material")
	}
}
