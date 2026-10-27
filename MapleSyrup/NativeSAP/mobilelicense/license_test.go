package mobilelicense

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

type vector struct {
	SINF, Credential, Work, Record, Key, Ciphertext, Plaintext_SHA256 string
	Page                                                              uint32
}

func vectors(t *testing.T) []vector {
	t.Helper()
	data, err := os.ReadFile("testdata/reference.json")
	if err != nil {
		t.Fatal(err)
	}
	var v []vector
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	if len(v) != 12 {
		t.Fatal("reference coverage changed")
	}
	return v
}
func unhex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func TestIndependentMobileReference(t *testing.T) {
	for _, v := range vectors(t) {
		sinf, cred := unhex(t, v.SINF), unhex(t, v.Credential)
		beforeSINF, beforeCred := bytes.Clone(sinf), bytes.Clone(cred)
		work, err := deriveWork16(cred, sinf)
		if err != nil || !bytes.Equal(work, unhex(t, v.Work)) {
			t.Fatal("Python work mismatch", err)
		}
		block, _ := aes.NewCipher(xorBytes(work, workMask))
		record := unhex(t, v.Record)
		key, err := derivePageKey(block, record[:16], v.Page)
		if err != nil || !bytes.Equal(key[:], unhex(t, v.Key)) {
			t.Fatal("Python page key mismatch", err)
		}
		l, err := Open(sinf, cred)
		if err != nil {
			t.Fatal(err)
		}
		ciphertext := unhex(t, v.Ciphertext)
		before := bytes.Clone(ciphertext)
		out, err := l.DecryptPage(ciphertext, record, uint64(v.Page)*4096)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(out)
		if hex.EncodeToString(sum[:]) != v.Plaintext_SHA256 {
			t.Fatal("Python plaintext mismatch")
		}
		if !l.Matches(sinf) || l.Matches(append(bytes.Clone(sinf), 0)) {
			t.Fatal("SINF binding")
		}
		if !bytes.Equal(before, ciphertext) || !bytes.Equal(sinf, beforeSINF) || !bytes.Equal(cred, beforeCred) {
			t.Fatal("inputs modified")
		}
		l.Close()
		l.Close()
		if l.Matches(sinf) || l.work != [16]byte{} {
			t.Fatal("closed material retained")
		}
		if out, err := l.DecryptPage(ciphertext, record, 0); err == nil || out != nil {
			t.Fatal("closed license used")
		}
	}
}
func TestRejectedInputsNeverYieldMaterial(t *testing.T) {
	v := vectors(t)[0]
	sinf, cred := unhex(t, v.SINF), unhex(t, v.Credential)
	for _, n := range []int{0, 1, 7, 8, 15, len(sinf) - 1} {
		if l, e := Open(sinf[:n], cred); e == nil || l != nil {
			t.Fatal("truncated SINF accepted")
		}
	}
	for _, n := range []int{0, 15, 17} {
		if l, e := Open(sinf, make([]byte, n)); e == nil || l != nil {
			t.Fatal("credential size accepted")
		}
	}
	for i := range cred {
		bad := bytes.Clone(cred)
		bad[i] ^= 1
		if l, e := Open(sinf, bad); e == nil || l != nil {
			t.Fatal("wrong credential accepted")
		}
	}
	bad := bytes.Clone(sinf)
	bad[len(bad)-1] = 1
	if l, e := Open(bad, cred); e == nil || l != nil {
		t.Fatal("bad trailer accepted")
	}
	l, e := Open(sinf, cred)
	if e != nil {
		t.Fatal(e)
	}
	defer l.Close()
	for _, args := range []struct {
		page, record int
		offset       uint64
	}{{0, 32, 0}, {4095, 32, 0}, {4096, 31, 0}, {4096, 33, 0}, {4096, 32, 1}, {4096, 32, 1 << 44}} {
		out, e := l.DecryptPage(make([]byte, args.page), make([]byte, args.record), args.offset)
		if e == nil || out != nil {
			t.Fatal("invalid page accepted")
		}
	}
}
func FuzzOpen(f *testing.F) {
	f.Add([]byte{}, []byte{})
	data, err := os.ReadFile("testdata/reference.json")
	if err != nil {
		f.Fatal(err)
	}
	var seeds []vector
	if err := json.Unmarshal(data, &seeds); err != nil {
		f.Fatal(err)
	}
	for _, v := range seeds {
		sinf, _ := hex.DecodeString(v.SINF)
		credential, _ := hex.DecodeString(v.Credential)
		f.Add(sinf, credential)
	}
	f.Fuzz(func(t *testing.T, sinf, credential []byte) {
		l, _ := Open(sinf, credential)
		if l != nil {
			l.Close()
		}
	})
}

// CBC has no authentication tag: do not promote changed candidate bytes to
// verified output merely because page decryption returned successfully.
func TestPageProcessingDoesNotClaimIntegrity(t *testing.T) {
	v := vectors(t)[0]
	l, e := Open(unhex(t, v.SINF), unhex(t, v.Credential))
	if e != nil {
		t.Fatal(e)
	}
	defer l.Close()
	ciphertext, record := unhex(t, v.Ciphertext), unhex(t, v.Record)
	ciphertext[0] ^= 1
	out, e := l.DecryptPage(ciphertext, record, uint64(v.Page)*4096)
	if e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(out)
	if hex.EncodeToString(sum[:]) == v.Plaintext_SHA256 {
		t.Fatal("corruption did not change plaintext")
	}
}
