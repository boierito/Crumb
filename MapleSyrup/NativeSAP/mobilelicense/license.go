// Package mobilelicense implements the portable mobile license consumer only.
// It cannot obtain an account credential. Callers must supply material from a
// separately verified producer, never from SIDB slicing or a guessed key.
package mobilelicense

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
)

// License holds SINF-checked working material, not an authenticated application.
// Its zero value is unusable. Do not copy it or use it concurrently with Close.
// No methods expose credentials, work keys, page keys or diagnostic digests.
type License struct {
	work  [16]byte
	sinf  [32]byte
	ready bool
}

// Open checks the original SINF's supported private-record structure using the
// supplied account credential. This is structural key confirmation, not a MAC
// or signature check. It neither modifies inputs nor retains the credential.
func Open(sinf, credential []byte) (*License, error) {
	if len(sinf) == 0 || len(sinf) > 1<<20 || len(credential) != 16 {
		return nil, errors.New("mobile_license_invalid_input")
	}
	work, err := deriveWork16(credential, sinf)
	if err != nil {
		return nil, errors.New("mobile_license_sinf_rejected")
	}
	defer clear(work)
	license := &License{sinf: sha256.Sum256(sinf), ready: true}
	copy(license.work[:], work)
	return license, nil
}

// Matches prevents accidentally reusing per-app material for another SINF.
func (l *License) Matches(sinf []byte) bool {
	return l != nil && l.ready && l.sinf == sha256.Sum256(sinf)
}

// Close clears retained working material. Go's AES implementation may retain
// expanded keys transiently during an operation; this is not a secure-erasure
// guarantee for the process heap.
func (l *License) Close() {
	if l != nil {
		clear(l.work[:])
		clear(l.sinf[:])
		l.ready = false
	}
}

// DecryptPage returns candidate plaintext for one 4096-byte page and its exact
// 32-byte SUPF record. absoluteOffset includes the FAT slice base, matching the
// mobile reference. The caller must validate integrity before exporting it.
// Errors return no plaintext. Neither the ciphertext nor record is modified.
func (l *License) DecryptPage(ciphertext, record []byte, absoluteOffset uint64) ([]byte, error) {
	if l == nil || !l.ready {
		return nil, errors.New("mobile_license_unavailable")
	}
	if len(ciphertext) != 4096 || len(record) != 32 || absoluteOffset%4096 != 0 || absoluteOffset/4096 > 0xffffffff {
		return nil, errors.New("mobile_page_invalid_input")
	}
	effective := xorBytes(l.work[:], workMask)
	defer clear(effective)
	block, err := aes.NewCipher(effective)
	if err != nil {
		return nil, errors.New("mobile_page_cipher_failed")
	}
	key, err := derivePageKey(block, record[:16], uint32(absoluteOffset/4096))
	if err != nil {
		return nil, errors.New("mobile_page_record_rejected")
	}
	defer clear(key[:])
	pageBlock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, errors.New("mobile_page_cipher_failed")
	}
	plaintext := make([]byte, 4096)
	cipher.NewCBCDecrypter(pageBlock, record[16:]).CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}
