package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/majd/ipatool/v2/internal/fairplay"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/frame"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/sidb"
	"github.com/majd/ipatool/v2/wafflebridge/recovery/state"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type labRequest struct {
	Operation, Source, Destination, Report, Model string
	Input, Frame                                  []byte
}

func labRun(r labRequest) (any, error) {
	switch r.Operation {
	case "export":
		file, e := os.Create(r.Destination)
		if e != nil {
			return nil, e
		}
		defer file.Close()
		if e = file.Chmod(0600); e != nil {
			return nil, e
		}
		writer := zip.NewWriter(file)
		e = filepath.WalkDir(r.Source, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlink")
			}
			name := entry.Name()
			if name != "manifest.json" && name != "kbsync.bin" && name != "keybag.bin" && !strings.HasPrefix(name, "sinf-") {
				return fmt.Errorf("unexpected capture file")
			}
			relative, e := filepath.Rel(r.Source, path)
			if e != nil {
				return e
			}
			dst, e := writer.Create(filepath.ToSlash(relative))
			if e != nil {
				return e
			}
			src, e := os.Open(path)
			if e != nil {
				return e
			}
			defer src.Close()
			_, e = io.Copy(dst, src)
			return e
		})
		closeErr := writer.Close()
		if e != nil {
			return nil, e
		}
		if closeErr != nil {
			return nil, closeErr
		}
		return map[string]any{"private_capture": true}, nil
	case "selftest":
		if !state.VerifyResource() || !frame.VerifyResource() || !sidb.VerifyResource() {
			return nil, fmt.Errorf("resource integrity")
		}
		a, e := state.Evaluate(make([]byte, 224), nil)
		if e != nil {
			return nil, e
		}
		b, e := frame.Evaluate(make([]byte, 16), nil)
		if e != nil {
			return nil, e
		}
		f := sha256.Sum256([]byte("ipatool public SIDB frame fixture"))
		i := sha256.Sum256([]byte("ipatool public SIDB block fixture"))
		c, e := sidb.Evaluate(i[:16], f[:16])
		if e != nil {
			return nil, e
		}
		hashes := []string{}
		for _, v := range [][]byte{a, b, c} {
			h := sha256.Sum256(v)
			hashes = append(hashes, hex.EncodeToString(h[:]))
		}
		expected := []string{"dbfb58ced0660063f9de17b93eb18997986d166ada93c3ef42c4315fd257f640", "7d3b5648332a3ca5df0971442b796a4df472c09e63e9861756d982c2ff3d5fb4", "fa60630407415ac643265feb92f467340dd4bc63fb3f37cf308a6619436cb898"}
		for n := range hashes {
			if hashes[n] != expected[n] {
				return nil, fmt.Errorf("fixture mismatch")
			}
		}
		return map[string]any{"state": "passed", "frame": "passed", "sidb": "passed", "synthetic": true, "revision": "ef98227"}, nil
	case "inspect":
		return fairplay.InspectPortableIPA(r.Source)
	case "state":
		return state.Evaluate(r.Input, nil)
	case "frame":
		return frame.Evaluate(r.Input, nil)
	case "sidb":
		return sidb.Recover(r.Frame, 2)
	case "decrypt":
		if r.Model == "" {
			return nil, fmt.Errorf("model required")
		}
		v, e := fairplay.DecryptIPA(fairplay.DecryptOptions{InputPath: r.Source, OutputPath: r.Destination, ReportPath: r.Report, ModelPath: r.Model})
		if e != nil {
			return nil, e
		}
		return map[string]any{"images": v.EncryptedImages, "pages": v.DecryptedPages, "completed": true, "plaintext_authenticated": false}, nil
	}
	return nil, fmt.Errorf("unsupported operation")
}
