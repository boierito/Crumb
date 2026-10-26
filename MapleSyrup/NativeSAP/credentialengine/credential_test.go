package fairplay

import (
    "bytes"
    "testing"
)

func TestInvalidCredentialNeverProducesIPA(t *testing.T) {
    var dst bytes.Buffer
    for _, bad := range [][]byte{nil, make([]byte, 15), make([]byte, 17)} {
        result, err := Decrypt("/missing-source.ipa", &dst, bad, false)
        if err == nil || result.Binaries != 0 || result.Pages != 0 || dst.Len() != 0 {
            t.Fatal("invalid credential progressed to IPA output")
        }
    }
}
