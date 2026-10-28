package machine

import (
	"archive/zip"
	"context"
	"encoding/json"
	"github.com/majd/ipatool/v2/internal/sap/assets"
	"io"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

// Optional private replay. CI without user inputs skips it; no credentials or
// private capture are embedded in this source or sent to a remote service.
func TestRetainedPrivateCaptureProvisioning(t *testing.T) {
	if os.Getenv("CRUMB_PRIVATE_CAPTURE") == "" {
		t.Skip("private capture not supplied")
	}
	if os.Getenv("CRUMB_PRODUCER_REPORT") == "" {
		t.Fatal("private report destination required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	z, e := zip.OpenReader(os.Getenv("CRUMB_PRIVATE_CAPTURE"))
	if e != nil {
		t.Fatal("capture unavailable")
	}
	defer z.Close()
	var c struct {
		Hardware []byte `json:"hardware_id"`
		DSID     string `json:"dsid"`
	}
	var s struct {
		Files map[string][]byte `json:"state_files"`
	}
	for _, f := range z.File {
		n := path.Base(f.Name)
		if n != "context.json" && n != "consumer-state.json" {
			continue
		}
		r, e := f.Open()
		if e != nil {
			t.Fatal("capture read")
		}
		b, e := io.ReadAll(io.LimitReader(r, 8<<20))
		r.Close()
		if e != nil {
			t.Fatal("capture read")
		}
		if n == "context.json" {
			e = json.Unmarshal(b, &c)
		} else {
			e = json.Unmarshal(b, &s)
		}
		clear(b)
		if e != nil {
			t.Fatal("capture decode")
		}
	}
	defer clear(c.Hardware)
	defer func(files map[string][]byte) {
		for _, b := range files {
			clear(b)
		}
	}(s.Files)
	dsid, e := strconv.ParseUint(c.DSID, 10, 64)
	if e != nil || dsid == 0 || len(c.Hardware) != 6 {
		t.Fatal("capture identity")
	}
	authorized := os.Getenv("CRUMB_PRIVATE_AUTHORIZED_STATE")
	if authorized != "" {
		stored, err := readAuthorizedStateFile(authorized, dsid, c.Hardware)
		if err != nil {
			t.Fatal("authorized state validation failed")
		}
		defer stored.clear()
		s.Files = stored.Files
	}
	bundle, e := assets.Load(ctx)
	if e != nil {
		t.Fatal("assets unavailable")
	}
	producer, e := OpenRetainedLicenseState(ctx, bundle, c.Hardware, s.Files)
	if e != nil {
		t.Fatal("native restore failed")
	}
	defer producer.Close()
	needs, e := producer.NeedsAccountDiversity(dsid)
	if e != nil {
		t.Fatal("native account query failed")
	}
	result := map[string]any{"schema": "crumb-real-producer-host-v1", "actual_capture_used": true, "authorized_state_supplied": authorized != "", "tci_backend": true, "state_restored": true, "native_account_query_completed": true, "needs_account_diversity": needs, "apple_authorization_executed": false, "sinf_validated": false}
	b, e := json.Marshal(result)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(os.Getenv("CRUMB_PRODUCER_REPORT"), b, 0600); e != nil {
		t.Fatal("report failed")
	}
}
