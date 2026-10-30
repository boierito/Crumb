package packageipa

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"errors"
	"howett.net/plist"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Every executable remains plaintext after signing; no licensing input may
// reappear. The signed copy intentionally includes the user's public profile.
func validatePlaintextTree(root string, signed bool) error {
	images := 0
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return errors.New("symlink in signing tree")
		}
		rel, e := filepath.Rel(root, p)
		if e != nil {
			return e
		}
		lower := strings.ToLower(filepath.ToSlash(rel))
		for _, part := range strings.Split(lower, "/") {
			if part == "sc_info" || part == "itunesmetadata.plist" || part == "_masreceipt" || part == "crumbsystempayload" {
				return errors.New("license data in decrypted signing tree")
			}
		}
		if !signed && systemPrivatePath(rel) {
			return errors.New("signature in unsigned input")
		}
		if d.IsDir() {
			return nil
		}
		st, e := d.Info()
		if e != nil {
			return e
		}
		if !st.Mode().IsRegular() {
			return errors.New("nonregular signing input")
		}
		f, e := os.Open(p)
		if e != nil {
			return e
		}
		defer f.Close()
		image, e := systemMachO(f, st.Size())
		if e != nil {
			return e
		}
		if image != nil {
			images++
			if systemEncrypted(image) != 0 {
				return errors.New("encrypted executable in install package")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if images == 0 {
		return errors.New("no executable in install package")
	}
	return nil
}

func ArchiveSystemExport(source, baseline, destination string) (map[string]any, error) {
	z, err := zip.OpenReader(baseline)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	files, root, err := systemIndex(&z.Reader)
	if err != nil {
		return nil, err
	}
	appRoot := filepath.Join(source, filepath.FromSlash(root))
	if err = validatePlaintextTree(appRoot, true); err != nil {
		return nil, err
	}
	// All original non-Mach-O resources, including Info.plist, stay byte exact.
	for name, f := range files {
		if !strings.HasPrefix(name, root) {
			continue
		}
		p := filepath.Join(source, filepath.FromSlash(name))
		b, e := os.Open(p)
		if e != nil {
			return nil, e
		}
		stat, e := b.Stat()
		if e != nil {
			b.Close()
			return nil, e
		}
		image, e := systemMachO(b, stat.Size())
		if e != nil {
			b.Close()
			return nil, e
		}
		if image != nil {
			b.Close()
			continue
		}
		a, e := f.Open()
		if e != nil {
			b.Close()
			return nil, e
		}
		ha, hb := sha256.New(), sha256.New()
		_, ea := io.Copy(ha, a)
		_, eb := io.Copy(hb, b)
		a.Close()
		b.Close()
		if ea != nil || eb != nil || !bytes.Equal(ha.Sum(nil), hb.Sum(nil)) {
			return nil, errors.New("signer changed original resource")
		}
	}
	data, err := os.ReadFile(filepath.Join(appRoot, "Info.plist"))
	if err != nil {
		return nil, err
	}
	var info Info
	if _, err = plist.Unmarshal(data, &info); err != nil || info.BundleID == "" {
		return nil, errors.New("invalid signed identity")
	}
	err = systemZIP(destination, func(w *zip.Writer) error {
		return filepath.WalkDir(appRoot, func(p string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if d.IsDir() {
				return nil
			}
			s, e := d.Info()
			if e != nil {
				return e
			}
			rel, e := filepath.Rel(source, p)
			if e != nil {
				return e
			}
			name := filepath.ToSlash(rel)
			if files[name] == nil && !strings.Contains(name, "/_CodeSignature/") && !strings.HasSuffix(name, "/embedded.mobileprovision") {
				return errors.New("unexpected added signing resource")
			}
			f, e := os.Open(p)
			if e != nil {
				return e
			}
			defer f.Close()
			return systemPut(w, name, s.Mode(), f)
		})
	})
	return map[string]any{"bundle_id": info.BundleID, "version": info.Version, "build": info.Build, "installation_verified": false, "contains_signing_identity": true}, err
}
