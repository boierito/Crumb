// JIT-free alternative for a user-supplied, account-matched 16-byte
// credential. Algorithm adapted from EthanArbuckle/fairplay-decrypt (MIT).
// This is distinct from the macOS FPSENR3 portable enrollment model.
package fairplay

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "os"
)

var credentialWorkMask = [16]byte{0xc7,0x62,0x61,0x10,0xce,0xcd,0x15,0x96,0x3d,0x85,0x0f,0x66,0x05,0xf1,0x1c,0x4b}
var credentialConstant = [16]byte{0xb4,0xf2,0x7c,0xb8,0x4e,0xaa,0xff,0x3c,0x13,0x09,0xbe,0x40,0x89,0xf0,0x29,0x72}
const credentialContextBase uint64 = 0x3e4a45b664eba141

// The raw credential never appears in an exported report or error. It must
// already belong to the account that owns the IPA; Crumb does not obtain it.
func loadAccountCredential(path string) ([16]byte, error) {
    var key [16]byte
    info, err := os.Lstat(path)
    if err != nil || !info.Mode().IsRegular() || info.Size() != 16 ||
        info.Mode().Perm()&0077 != 0 {
        return key, errors.New("credential_file_invalid")
    }
    file, err := os.Open(path)
    if err != nil { return key, errors.New("credential_file_unreadable") }
    defer file.Close()
    if _, err = io.ReadFull(file, key[:]); err != nil {
        clear(key[:])
        return [16]byte{}, errors.New("credential_file_unreadable")
    }
    return key, nil
}

// credentialWorkFromSINF validates the encrypted 'priv' structure before
// returning a derivation key. Wrong credentials never become success models.
func credentialWorkFromSINF(credential []byte, sinf []byte) (work [16]byte, err error) {
    if len(credential) != 16 { return work, errors.New("credential_size_invalid") }
    fields, err := parseSINFFields(sinf)
    if err != nil { return work, errors.New("credential_sinf_invalid") }
    nameAtom, err := uniqueAtom(fields, "name")
    if err != nil { return work, errors.New("credential_sinf_name_missing") }
    ivAtom, err := uniqueAtom(fields, "iviv")
    if err != nil || ivAtom.end-ivAtom.start != 16 {
        return work, errors.New("credential_sinf_iv_invalid")
    }
    private, err := uniqueAtom(fields, "priv")
    if err != nil { return work, errors.New("credential_sinf_priv_missing") }
    payload := sinf[private.start:private.end]
    if len(payload) != 440 && len(payload) != 448 {
        return work, errors.New("credential_sinf_priv_unsupported")
    }
    if len(payload) == 440 {
        if !bytes.Equal(payload[432:], make([]byte, 8)) {
            return work, errors.New("credential_sinf_trailer_invalid")
        }
        payload = payload[:432]
    }
    name := sinf[nameAtom.start:nameAtom.end]
    if end := bytes.IndexByte(name, 0); end != -1 { name = name[:end] }
    digest := md5.New()
    digest.Write(name)
    digest.Write(sinf[ivAtom.start:ivAtom.end])
    iv := digest.Sum(nil)
    var aesKey [16]byte
    for i := range aesKey { aesKey[i] = credential[i] ^ byte('0') }
    block, err := aes.NewCipher(aesKey[:])
    clear(aesKey[:])
    if err != nil { return work, errors.New("credential_cipher_invalid") }
    decoded := make([]byte, len(payload))
    defer clear(decoded)
    cipher.NewCBCDecrypter(block, iv).CryptBlocks(decoded, payload)
    if len(decoded) < 48 || string(decoded[0:4]) != "itun" ||
        string(decoded[20:24]) != "key " || string(decoded[44:48]) != "iviv" {
        return work, errors.New("credential_sinf_mismatch")
    }
    for i := range work { work[i] = decoded[24+i] ^ credentialWorkMask[i] }
    return work, nil
}

type credentialPageKeyDeriver struct{ block cipher.Block }

func newCredentialPageKeyDeriver(work [16]byte) *credentialPageKeyDeriver {
    for i := range work { work[i] ^= credentialWorkMask[i] }
    block, _ := aes.NewCipher(work[:])
    clear(work[:])
    return &credentialPageKeyDeriver{block: block}
}

var credentialSelectors = [4][8]uint8{
    {0,1,2,3,4,5,6,7},
    {2,1,0,3,4,7,6,5},
    {2,0,4,1,6,3,7,5},
    {1,2,4,0,7,3,5,6},
}
var credentialRoundConstants = [4]uint32{0x6ee6b74f,0x78b72abf,0xd227c8ea,0x0c0081bc}
var credentialFinalSelectors = [4][4]uint8{{0,1,2,3},{0,2,4,6},{0,3,5,7},{4,5,6,7}}
var credentialFinalConstants = [4]uint64{
    0x4f3dac1fd3361949,0x4f3dac1f631a63a2,0x4f3dac1f87cd2974,0x4f3dac1f7ea63ab1,
}

func credentialFinish(block [16]byte) (out [16]byte) {
    var initial, row [8]uint16
    for i := 0; i < 8; i++ {
        initial[i] = binary.BigEndian.Uint16(block[i*2:]) ^ 0x18fb
        row[i] = initial[i]
    }
    for round := 0; round < 4; round++ {
        var next [8]uint16
        for pair := 0; pair < 4; pair++ {
            a := uint32(row[credentialSelectors[round][2*pair]] ^ 0x18fb)
            b := uint32(row[credentialSelectors[round][2*pair+1]] ^ 0x18fb)
            sum := a + b
            square := uint32(uint64(sum)*uint64(sum)) ^ credentialRoundConstants[pair]
            low := (uint32(0xcb61) - ((sum>>16)*(sum>>16))) ^ square
            next[2*pair] = uint16(low ^ 0x18fb)
            next[2*pair+1] = uint16((square>>16)^0x1e49)
        }
        row = next
    }
    for j := 0; j < 4; j++ {
        indices := credentialFinalSelectors[j]
        left := uint32(row[indices[0]]^0x18fb) | uint32(row[indices[1]]^0x18fb)<<16
        right := uint32(row[indices[2]]^0x18fb) | uint32(row[indices[3]]^0x18fb)<<16
        total := uint64(uint32(left+right))
        x9 := total*total ^ 0x4f3dac1f8847818e
        x19 := credentialFinalConstants[j] ^ x9
        x9 = x19 ^ (x9>>32)
        low := uint16(x9 ^ 0xb4e4)
        high := uint16((x9>>16) ^ 0x4f3d ^ 0x18fb)
        binary.BigEndian.PutUint16(out[j*4:], initial[2*j]^low)
        binary.BigEndian.PutUint16(out[j*4+2:], initial[2*j+1]^high)
    }
    return out
}

func (d *credentialPageKeyDeriver) pageKey(wrapped []byte, cryptOffset uint64) (result [16]byte) {
    if len(wrapped) != 16 || cryptOffset % 4096 != 0 {
        return result
    }
    absolutePage := cryptOffset/4096
    var sources [3][16]byte
    for i := 0; i < 8; i++ {
        sources[0][i] = wrapped[i+8] ^ credentialConstant[i]
        sources[2][i+8] = wrapped[i] ^ credentialConstant[i+8]
    }
    binary.BigEndian.PutUint64(sources[0][8:], (credentialContextBase^(absolutePage<<12))^0x2d43fbf6ed1b8833)
    for i := range sources[1] { sources[1][i] = wrapped[i] ^ credentialConstant[i] }
    binary.BigEndian.PutUint64(sources[2][:8], binary.BigEndian.Uint64(credentialConstant[:8])^(absolutePage<<12))
    var encrypted [3][16]byte
    for i := range sources {
        for j := 0; j < 16; j++ { sources[i][j] ^= credentialConstant[j] }
        d.block.Encrypt(encrypted[i][:], sources[i][:])
    }
    state := credentialFinish(encrypted[0])
    for i := range state { state[i] ^= encrypted[1][i] }
    state = credentialFinish(state)
    for i := range state { state[i] ^= encrypted[2][i] }
    return credentialFinish(state)
}

// A non-secret capability signal. This does not imply that the user's iPhone
// already has a credential or that a particular personalized SINF is supported.
func credentialCapability() string { return fmt.Sprintf("offline-credential-v1") }
