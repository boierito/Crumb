package packageipa

import (
	"archive/zip"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PrepareSystemLocal extracts the purchased target into a NEW private data
// directory. Unlike PrepareSystemHelper this creates no Crumb bundle, signature
// or install job. Whether fairplayd can read this location is a DEVICE question;
// extraction success does not assert authorization or successful decryption.
func PrepareSystemLocal(original, destination string) (job SystemJob, err error) {
	z, err := zip.OpenReader(original)
	if err != nil {
		return job, err
	}
	defer z.Close()
	files, root, err := systemIndex(&z.Reader)
	if err != nil {
		return job, err
	}
	app, copies, err := systemLicensePlan(files, root)
	if err != nil {
		return job, err
	}
	var nonce [16]byte
	if _, err = rand.Read(nonce[:]); err != nil {
		return job, err
	}
	if err = os.Mkdir(destination, 0700); err != nil {
		return job, err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(destination)
		}
	}()
	copyFile := func(name string, source *zip.File) error {
		path := filepath.Join(destination, filepath.FromSlash(name))
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			return e
		}
		out, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, source.Mode().Perm()|0600)
		if e != nil {
			return e
		}
		in, e := source.Open()
		if e != nil {
			_ = out.Close()
			return e
		}
		_, e = io.Copy(out, in)
		_ = in.Close()
		closeErr := out.Close()
		if e != nil {
			return e
		}
		return closeErr
	}
	for _, f := range z.File {
		if !strings.HasPrefix(f.Name, root) || f.FileInfo().IsDir() {
			continue
		}
		if err = copyFile(strings.TrimPrefix(f.Name, root), f); err != nil {
			return job, err
		}
		job.Files++
	}
	for _, c := range copies {
		if err = copyFile(c.Relative, files[c.Source]); err != nil {
			return job, err
		}
	}
	job.Schema = "crumb-system-fairplay-local-v1"
	job.JobID = hex.EncodeToString(nonce[:])
	job.BundleID = app.BundleID
	job.Version = app.Version
	job.Build = app.Build
	return job, nil
}
