package sap

import (
	"context"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"github.com/majd/ipatool/v2/internal/sap/machine"
	"testing"
	"time"
)

func TestTCIKBSyncSmoke(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	bundle, e := assets.Load(ctx)
	if e != nil {
		t.Fatal(e)
	}
	data, e := machine.GenerateKBSync(ctx, bundle, []byte{2, 0, 0, 0, 0, 1}, 1)
	if e != nil {
		t.Fatal(e)
	}
	if len(data) == 0 {
		t.Fatal("empty kbsync")
	}
	t.Logf("kbsync generated: %d bytes; DSID is synthetic; contents withheld", len(data))
}
