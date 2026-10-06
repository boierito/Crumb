package sap

import (
	"bytes"
	"context"
	"howett.net/plist"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Manual network test: SAP setup only, never credentials or authentication.
func TestTCIDynamicSAPSmoke(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, "GET", "https://init.itunes.apple.com/bag.xml?guid=020000000001", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("User-Agent", "Configurator/2.17 (Macintosh; OS X 15.2; 24C5089c) AppleWebKit/0620.1.16.11.6")
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatalf("Bag HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > 1<<20 {
		t.Fatal("oversized Bag")
	}
	if start := bytes.Index(data, []byte("<plist")); start >= 0 {
		if end := bytes.Index(data[start:], []byte("</plist>")); end >= 0 {
			data = data[start : start+end+8]
		}
	}
	var root struct {
		URLBag map[string]interface{} `plist:"urlBag"`
	}
	if _, err = plist.Unmarshal(data, &root); err != nil {
		t.Fatal(err)
	}
	versionText, _ := root.URLBag["sign-sap-version"].(string)
	version, err := strconv.ParseUint(versionText, 10, 32)
	if err != nil {
		t.Fatal("invalid SAP version")
	}
	setup, _ := root.URLBag["sign-sap-setup"].(string)
	certificate, _ := root.URLBag["sign-sap-setup-cert"].(string)
	signer, err := NewSigner(ctx, Config{SetupURL: setup, CertificateURL: certificate, Version: uint32(version), HardwareID: []byte{2, 0, 0, 0, 0, 1}})
	if err != nil {
		t.Fatal(err)
	}
	defer signer.Close()
	first, err := signer.Sign([]byte("WaffleStore no-JIT SAP probe"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := signer.Sign([]byte("WaffleStore second no-JIT SAP probe"))
	if err != nil {
		t.Fatal(err)
	}
	if len(first) == 0 || len(second) == 0 || bytes.Equal(first, second) {
		t.Fatal("empty or reused signature")
	}
	t.Logf("SAP signatures generated: %d and %d bytes; contents withheld; Apple login not attempted", len(first), len(second))
}
