package packageipa

import (
	"archive/zip"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

const SystemPayload = "CrumbSystemPayload"

type SystemJob struct {
	Schema           string `json:"schema"`
	JobID            string `json:"job_id"`
	BundleID         string `json:"bundle_id"`
	Version          string `json:"version"`
	Files            int    `json:"files"`
	RequiresSigning  bool   `json:"requires_signing"`
	LicenseValidated bool   `json:"license_validated"`
}

func systemSafe(name string) bool {
	return safeName(name) && !strings.ContainsRune(name, 0) && len(name) < 4096
}

func systemIndex(z *zip.Reader) (map[string]*zip.File, string, error) {
	files := map[string]*zip.File{}
	seen := map[string]bool{}
	root := ""
	var total uint64
	if len(z.File) > 200000 {
		return nil, "", errors.New("archive entry limit")
	}
	for _, f := range z.File {
		if !systemSafe(f.Name) {
			return nil, "", errors.New("unsafe ZIP path")
		}
		key := strings.ToLower(strings.TrimSuffix(f.Name, "/"))
		if seen[key] {
			return nil, "", errors.New("colliding ZIP paths")
		}
		seen[key] = true
		if f.Mode()&os.ModeSymlink != 0 || !f.Mode().IsRegular() && !f.FileInfo().IsDir() {
			return nil, "", errors.New("unsupported archive file type")
		}
		if f.FileInfo().IsDir() {
			continue
		}
		if f.UncompressedSize64 > 8<<30 || total > 32<<30-f.UncompressedSize64 {
			return nil, "", errors.New("archive size limit")
		}
		total += f.UncompressedSize64
		files[f.Name] = f
		parts := strings.Split(f.Name, "/")
		if len(parts) == 3 && parts[0] == "Payload" && strings.HasSuffix(parts[1], ".app") && parts[2] == "Info.plist" {
			if root != "" {
				return nil, "", errors.New("multiple IPA apps")
			}
			root = "Payload/" + parts[1] + "/"
		}
	}
	if root == "" {
		return nil, "", errors.New("IPA app is missing")
	}
	return files, root, nil
}

func systemZIP(destination string, work func(*zip.Writer) error) (err error) {
	f, e := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	defer func() {
		_ = f.Close()
		if err != nil {
			_ = os.Remove(destination)
		}
	}()
	z := zip.NewWriter(f)
	if err = work(z); err != nil {
		_ = z.Close()
		return err
	}
	if err = z.Close(); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return f.Close()
}

func systemPut(z *zip.Writer, name string, mode fs.FileMode, reader io.Reader) error {
	if !systemSafe(name) {
		return errors.New("invalid output path")
	}
	// Fresh headers exclude original ZIP comments, extras, creator IDs and times.
	h := &zip.FileHeader{Name: name, Method: zip.Deflate}
	h.SetMode(mode.Perm())
	w, e := z.CreateHeader(h)
	if e != nil {
		return e
	}
	_, e = io.Copy(w, reader)
	return e
}

// PrepareSystemHelper snapshots this app's installed, readable bundle and adds
// one purchased IPA. It uses no developer account, signing key or server.
// iOS must accept a NEW signature and installation before the payload is usable.
func PrepareSystemHelper(bundleDirectory, original, destination string) (job SystemJob, err error) {
	if filepath.Ext(bundleDirectory) != ".app" {
		return job, errors.New("invalid helper bundle")
	}
	info, err := os.Lstat(bundleDirectory)
	if err != nil || !info.IsDir() {
		return job, errors.New("helper bundle unavailable")
	}
	abs, err := filepath.Abs(bundleDirectory)
	if err != nil {
		return job, err
	}
	dest, err := filepath.Abs(destination)
	if err != nil {
		return job, err
	}
	if rel, e := filepath.Rel(abs, dest); e != nil || rel == "." || !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return job, errors.New("output must be outside installed bundle")
	}
	source, err := zip.OpenReader(original)
	if err != nil {
		return job, err
	}
	defer source.Close()
	files, root, err := systemIndex(&source.Reader)
	if err != nil {
		return job, err
	}
	data, err := readSmall(files[root+"Info.plist"])
	if err != nil {
		return job, err
	}
	var app Info
	if _, err = plist.Unmarshal(data, &app); err != nil || app.BundleID == "" || app.Version == "" || !systemSafe(app.Executable) || strings.Contains(app.Executable, "/") {
		return job, errors.New("invalid app identity")
	}
	if files[root+app.Executable] == nil {
		return job, errors.New("app executable missing")
	}
	m := files[root+"SC_Info/Manifest.plist"]
	if m == nil {
		return job, errors.New("license manifest missing")
	}
	data, err = readSmall(m)
	if err != nil {
		return job, err
	}
	var manifest struct {
		SinfPaths            []string
		SinfReplicationPaths []string
	}
	if _, err = plist.Unmarshal(data, &manifest); err != nil || len(manifest.SinfPaths) == 0 || len(manifest.SinfReplicationPaths) > 4096 {
		return job, errors.New("invalid license manifest")
	}
	if !systemSafe(manifest.SinfPaths[0]) || files[root+manifest.SinfPaths[0]] == nil {
		return job, errors.New("original SINF missing")
	}
	for _, p := range manifest.SinfReplicationPaths {
		if !systemSafe(p) || strings.HasSuffix(p, "/") {
			return job, errors.New("invalid SINF replication path")
		}
	}
	var nonce [16]byte
	if _, err = rand.Read(nonce[:]); err != nil {
		return job, err
	}
	job = SystemJob{Schema: "crumb-system-fairplay-job-v1", JobID: hex.EncodeToString(nonce[:]), BundleID: app.BundleID, Version: app.Version, RequiresSigning: true}
	helperRoot := "Payload/" + filepath.Base(bundleDirectory) + "/"
	target := helperRoot + SystemPayload + "/Target/"
	err = systemZIP(destination, func(z *zip.Writer) error {
		if e := filepath.WalkDir(bundleDirectory, func(p string, entry fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			rel, e := filepath.Rel(bundleDirectory, p)
			if e != nil {
				return e
			}
			if rel == "." {
				return nil
			}
			if strings.EqualFold(strings.Split(filepath.ToSlash(rel), "/")[0], SystemPayload) {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return errors.New("symlink in helper bundle")
			}
			if entry.IsDir() {
				return nil
			}
			stat, e := entry.Info()
			if e != nil {
				return e
			}
			if !stat.Mode().IsRegular() {
				return errors.New("non-regular helper resource")
			}
			f, e := os.Open(p)
			if e != nil {
				return e
			}
			defer f.Close()
			return systemPut(z, helperRoot+filepath.ToSlash(rel), stat.Mode(), f)
		}); e != nil {
			return e
		}
		for _, f := range source.File {
			if !strings.HasPrefix(f.Name, root) || f.FileInfo().IsDir() {
				continue
			}
			r, e := f.Open()
			if e != nil {
				return e
			}
			e = systemPut(z, target+strings.TrimPrefix(f.Name, root), f.Mode(), r)
			_ = r.Close()
			if e != nil {
				return e
			}
			job.Files++
		}
		for _, relative := range manifest.SinfReplicationPaths {
			if files[root+relative] != nil {
				continue
			}
			r, e := files[root+manifest.SinfPaths[0]].Open()
			if e != nil {
				return e
			}
			e = systemPut(z, target+relative, 0600, r)
			_ = r.Close()
			if e != nil {
				return e
			}
			// Count a replicated path only once; duplicate manifest entries fail.
			files[root+relative] = files[root+manifest.SinfPaths[0]]
		}
		b, e := json.Marshal(job)
		if e != nil {
			return e
		}
		return systemPut(z, helperRoot+SystemPayload+"/job.json", 0600, strings.NewReader(string(b)))
	})
	return job, err
}

func systemPrivatePath(relative string) bool {
	parts := strings.Split(strings.ToLower(filepath.ToSlash(relative)), "/")
	for i, p := range parts {
		if p == "sc_info" || p == "_codesignature" || p == "_masreceipt" || p == "__macosx" || p == "embedded.mobileprovision" || p == "embedded.provisionprofile" || p == "itunesmetadata.plist" || p == ".ds_store" {
			return true
		}
		if p == "receipt" && i > 0 && parts[i-1] == "storekit" {
			return true
		}
	}
	return false
}

// Keep this helper separate from ZIP path validation: legitimate resource names
// are retained verbatim, and no account-derived strings are rewritten in code.
func systemTargetName(relative string) string {
	return path.Join("Payload/Target.app", filepath.ToSlash(relative))
}
