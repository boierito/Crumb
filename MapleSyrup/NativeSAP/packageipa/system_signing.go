package packageipa

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

type SystemSigningFolder struct {
	BundlePath string `json:"bundle_path"`
	BundleID   string `json:"bundle_id"`
	Version    string `json:"version"`
	Build      string `json:"build"`
}

func UnpackSystemSigning(source, destination string) (result SystemSigningFolder, err error) {
	return unpackSigning(source, destination, true)
}

func UnpackSystemExport(source, destination string) (SystemSigningFolder, error) {
	return unpackSigning(source, destination, false)
}

func unpackSigning(source, destination string, prepared bool) (result SystemSigningFolder, err error) {
	z, err := zip.OpenReader(source)
	if err != nil {
		return result, err
	}
	defer z.Close()
	files, root, err := systemIndex(&z.Reader)
	if err != nil {
		return result, err
	}
	if prepared {
		manifest := files[root+SystemPayload+"/job.json"]
		if manifest == nil {
			return result, errors.New("prepared job is missing")
		}
		data, e := readSmall(manifest)
		if e != nil {
			return result, e
		}
		var job SystemJob
		if json.Unmarshal(data, &job) != nil || job.Schema != "crumb-system-fairplay-job-v1" {
			return result, errors.New("invalid prepared job")
		}
	} else {
		for name := range files {
			if strings.Contains(name, "/"+SystemPayload+"/") {
				return result, errors.New("private helper is not a decrypted export")
			}
		}
	}
	data, e := readSmall(files[root+"Info.plist"])
	if e != nil {
		return result, e
	}
	var app Info
	if _, e = plist.Unmarshal(data, &app); e != nil || app.BundleID == "" {
		return result, errors.New("missing signing identity")
	}
	result.BundleID = app.BundleID
	result.Version = app.Version
	result.Build = app.Build
	if err = os.Mkdir(destination, 0700); err != nil {
		return result, err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(destination)
		}
	}()
	for _, f := range z.File {
		if !strings.HasPrefix(f.Name, root) || f.FileInfo().IsDir() {
			continue
		}
		p := filepath.Join(destination, filepath.FromSlash(f.Name))
		if err = os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			return result, err
		}
		out, e := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, f.Mode().Perm()|0600)
		if e != nil {
			return result, e
		}
		in, e := f.Open()
		if e != nil {
			_ = out.Close()
			return result, e
		}
		_, e = io.Copy(out, in)
		_ = in.Close()
		closeErr := out.Close()
		if e != nil {
			return result, e
		}
		if closeErr != nil {
			return result, closeErr
		}
	}
	if !prepared {
		e = validatePlaintextTree(filepath.Join(destination, filepath.FromSlash(root)), false)
		if e != nil {
			return result, e
		}
	}
	result.BundlePath = strings.TrimSuffix(root, "/")
	return result, nil
}

// ArchiveSystemSigning verifies that the signing library did not modify, remove
// or add ANY purchased-payload resource. New signatures apply only to Crumb.
func ArchiveSystemSigning(source, baseline, destination string) (result map[string]any, err error) {
	z, err := zip.OpenReader(baseline)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	files, root, err := systemIndex(&z.Reader)
	if err != nil {
		return nil, err
	}
	payload := root + SystemPayload + "/"
	data, err := readSmall(files[root+"Info.plist"])
	if err != nil {
		return nil, err
	}
	var original, current Info
	if _, err = plist.Unmarshal(data, &original); err != nil {
		return nil, err
	}
	infoPath := filepath.Join(source, filepath.FromSlash(root), "Info.plist")
	data, err = os.ReadFile(infoPath)
	if err != nil {
		return nil, err
	}
	if _, err = plist.Unmarshal(data, &current); err != nil || current.BundleID != original.BundleID {
		return nil, errors.New("helper bundle identity changed")
	}
	checked := 0
	for name, f := range files {
		if !strings.HasPrefix(name, payload) {
			continue
		}
		p := filepath.Join(source, filepath.FromSlash(name))
		stat, e := os.Lstat(p)
		if e != nil || !stat.Mode().IsRegular() || uint64(stat.Size()) != f.UncompressedSize64 {
			return nil, errors.New("signed payload inventory changed")
		}
		a, e := f.Open()
		if e != nil {
			return nil, e
		}
		b, e := os.Open(p)
		if e != nil {
			_ = a.Close()
			return nil, e
		}
		ha, hb := sha256.New(), sha256.New()
		_, ea := io.Copy(ha, a)
		_, eb := io.Copy(hb, b)
		_ = a.Close()
		_ = b.Close()
		if ea != nil || eb != nil || string(ha.Sum(nil)) != string(hb.Sum(nil)) {
			return nil, errors.New("signer modified purchased payload")
		}
		checked++
	}
	if checked == 0 {
		return nil, errors.New("signed payload is missing")
	}
	err = systemZIP(destination, func(writer *zip.Writer) error {
		return filepath.WalkDir(filepath.Join(source, filepath.FromSlash(root)), func(p string, entry fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return errors.New("symlink in signed package")
			}
			if entry.IsDir() {
				return nil
			}
			stat, e := entry.Info()
			if e != nil {
				return e
			}
			if !stat.Mode().IsRegular() {
				return errors.New("non-regular signed resource")
			}
			relative, e := filepath.Rel(source, p)
			if e != nil {
				return e
			}
			name := filepath.ToSlash(relative)
			if strings.HasPrefix(name, payload) && files[name] == nil {
				return errors.New("signer added purchased payload data")
			}
			f, e := os.Open(p)
			if e != nil {
				return e
			}
			defer f.Close()
			return systemPut(writer, name, stat.Mode(), f)
		})
	})
	return map[string]any{"payload_files_verified": checked, "installation_verified": false, "contains_private_license": true}, err
}
