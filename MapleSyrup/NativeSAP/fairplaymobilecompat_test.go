package fairplay

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/majd/ipatool/v2/wafflebridge/mobilelicense"
)

// Compare the actual IPA processor's key encoding/page engine with the
// independently generated EthanArbuckle mobile fixtures and the separate
// mobile consumer. This covers 440/448-byte priv records and absolute page
// indices through uint32 max. No real account credential is used.
func TestPortablePagesMatchIndependentMobileLicenseVectors(t *testing.T) {
	var vectors []struct {
		SINF, Credential, Work, Record, Key, Ciphertext string
		Page                                            uint32
		PlaintextSHA256                                 string `json:"plaintext_sha256"`
	}
	data, err := os.ReadFile("testdata/mobile-reference.json")
	if err != nil || json.Unmarshal(data, &vectors) != nil || len(vectors) != 12 {
		t.Fatal("independent mobile fixtures unavailable")
	}
	decode := func(value string) []byte {
		t.Helper()
		out, err := hex.DecodeString(value)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	mask := decode("c7626110cecd15963d850f6605f11c4b")
	for _, vector := range vectors {
		work, record := decode(vector.Work), decode(vector.Record)
		var inner [16]byte
		for i := range inner {
			inner[i] = work[i] ^ mask[i]
		}
		deriver := newPageKeyDeriver(appMaskFromLogicalKey(logicalFromInnerKey(inner)))
		offset := uint64(vector.Page) * 4096
		key := deriver.pageKey(record[:16], offset)
		if !bytes.Equal(key[:], decode(vector.Key)) {
			t.Fatal("portable and mobile page-key derivations differ")
		}
		block, err := aes.NewCipher(key[:])
		if err != nil {
			t.Fatal(err)
		}
		ciphertext := decode(vector.Ciphertext)
		plaintext := make([]byte, len(ciphertext))
		cipher.NewCBCDecrypter(block, record[16:]).CryptBlocks(plaintext, ciphertext)
		digest := sha256.Sum256(plaintext)
		if hex.EncodeToString(digest[:]) != vector.PlaintextSHA256 {
			t.Fatal("portable page plaintext differs from independent reference")
		}
		license, err := mobilelicense.Open(decode(vector.SINF), decode(vector.Credential))
		if err != nil {
			t.Fatal(err)
		}
		mobile, err := license.DecryptPage(ciphertext, record, offset)
		license.Close()
		if err != nil || !bytes.Equal(mobile, plaintext) {
			t.Fatal("separate mobile consumer and IPA engine disagree")
		}
		wrong := deriver.pageKey(record[:16], uint64(vector.Page^1)*4096)
		if wrong == key {
			t.Fatal("page engine ignored the absolute page index")
		}
		clear(inner[:])
		clear(key[:])
		clear(work)
		clear(plaintext)
		clear(mobile)
	}
}
