package packageipa

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

type SystemExportReport struct {
	Schema                  string `json:"schema"`
	Files                   int    `json:"files"`
	MachOFiles              int    `json:"macho_files"`
	EncryptedSlices         int    `json:"encrypted_slices_processed"`
	ChangedPages            int    `json:"changed_16k_pages"`
	RemovedFiles            int    `json:"known_metadata_entries_excluded"`
	RemovedSignatures       int    `json:"embedded_signatures_removed"`
	CheckedIdentifiers      int    `json:"known_identifiers_checked"`
	RequiresSigning         bool   `json:"requires_signing"`
	PlaintextAuthenticated  bool   `json:"plaintext_authenticated"`
	NoAccountDataGuaranteed bool   `json:"no_account_data_guaranteed"`
}

func systemPatterns(identifiers []string) ([][]byte, error) {
	if len(identifiers) > 32 {
		return nil, errors.New("too many private identifiers")
	}
	var patterns [][]byte
	for _, s := range identifiers {
		if len(s) < 6 || len(s) > 4096 {
			return nil, errors.New("invalid private identifier length")
		}
		for _, value := range []string{s, strings.ToLower(s)} {
			patterns = append(patterns, []byte(value), []byte(base64.StdEncoding.EncodeToString([]byte(value))))
			units := utf16.Encode([]rune(value))
			le, be := make([]byte, 2*len(units)), make([]byte, 2*len(units))
			for i, u := range units {
				binary.LittleEndian.PutUint16(le[2*i:], u)
				binary.BigEndian.PutUint16(be[2*i:], u)
			}
			patterns = append(patterns, le, be)
		}
	}
	return patterns, nil
}

func systemContains(data []byte, patterns [][]byte) bool {
	for _, p := range patterns {
		if bytes.Contains(data, p) {
			return true
		}
	}
	return false
}

func systemCopyChecked(writer io.Writer, reader io.Reader, patterns [][]byte) error {
	maxPattern := 0
	for _, p := range patterns {
		maxPattern = max(maxPattern, len(p))
	}
	buffer := make([]byte, 64<<10)
	var tail []byte
	for {
		n, e := reader.Read(buffer)
		if n > 0 {
			window := append(tail, buffer[:n]...)
			if systemContains(window, patterns) {
				return errors.New("known private identifier remains; export rejected")
			}
			if _, err := writer.Write(buffer[:n]); err != nil {
				return err
			}
			if maxPattern > 0 {
				tail = bytes.Clone(window[max(0, len(window)-maxPattern+1):])
			}
		}
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}
	}
}

// ExportSystemResult verifies native overrides against the installed ciphertext,
// then emits a fresh unsigned archive. It does NOT certify absence of arbitrary
// app-defined watermarks, obfuscated identifiers or unknown account data.
func ExportSystemResult(source, overrides, destination string, identifiers []string) (report SystemExportReport, err error) {
	report = SystemExportReport{Schema: "crumb-system-fairplay-export-v1", RequiresSigning: true, CheckedIdentifiers: len(identifiers)}
	patterns, err := systemPatterns(identifiers)
	if err != nil {
		return report, err
	}
	for _, p := range []string{source, overrides} {
		info, e := os.Lstat(p)
		if e != nil || !info.IsDir() {
			return report, errors.New("private export directory unavailable")
		}
	}
	for _, root := range []string{source, overrides} {
		rel, e := filepath.Rel(root, destination)
		if e != nil || rel == "." || !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return report, errors.New("export output overlaps input")
		}
	}
	err = systemZIP(destination, func(z *zip.Writer) error {
		return filepath.WalkDir(source, func(p string, entry fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			relative, e := filepath.Rel(source, p)
			if e != nil {
				return e
			}
			if relative == "." {
				return nil
			}
			if !systemSafe(filepath.ToSlash(relative)) {
				return errors.New("invalid installed resource path")
			}
			if systemPrivatePath(relative) {
				report.RemovedFiles++
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return errors.New("symlink in exported app")
			}
			if entry.IsDir() {
				return nil
			}
			info, e := entry.Info()
			if e != nil {
				return e
			}
			if !info.Mode().IsRegular() {
				return errors.New("non-regular exported resource")
			}
			if systemContains([]byte(filepath.ToSlash(relative)), patterns) {
				return errors.New("known private identifier in resource name")
			}
			original, e := os.Open(p)
			if e != nil {
				return e
			}
			defer original.Close()
			var reader io.ReaderAt = original
			image, e := systemMachO(original, info.Size())
			if e != nil {
				return e
			}
			if image != nil {
				report.MachOFiles++
				if encrypted := systemEncrypted(image); encrypted > 0 {
					candidatePath := filepath.Join(overrides, relative)
					stat, e := os.Lstat(candidatePath)
					if e != nil || !stat.Mode().IsRegular() || stat.Size() != info.Size() {
						return errors.New("complete native output missing")
					}
					candidate, e := os.Open(candidatePath)
					if e != nil {
						return e
					}
					defer candidate.Close()
					pages, e := systemVerifyNative(original, candidate, info.Size(), image)
					if e != nil {
						return e
					}
					report.EncryptedSlices += encrypted
					report.ChangedPages += pages
					reader = candidate
					image, e = systemMachO(candidate, info.Size())
					if e != nil {
						return e
					}
				}
				var count int
				reader, count = systemRemoveSignatures(reader, image)
				report.RemovedSignatures += count
				clean, e := systemMachO(reader, info.Size())
				if e != nil || clean == nil || systemEncrypted(clean) != 0 {
					return errors.New("sanitized executable validation failed")
				}
				for _, s := range clean.slices {
					if len(s.signatures) > 0 {
						return errors.New("embedded signature remains")
					}
				}
			}
			// New ZIP headers discard account-related archive metadata and comments.
			h := &zip.FileHeader{Name: systemTargetName(relative), Method: zip.Deflate}
			h.SetMode(info.Mode().Perm())
			writer, e := z.CreateHeader(h)
			if e != nil {
				return e
			}
			if e = systemCopyChecked(writer, io.NewSectionReader(reader, 0, info.Size()), patterns); e != nil {
				return e
			}
			report.Files++
			return nil
		})
	})
	if err != nil {
		return report, err
	}
	// A package without any observed encrypted slice is not this operation's
	// successful result. Remove it instead of producing a misleading export.
	if report.EncryptedSlices == 0 {
		_ = os.Remove(destination)
		return report, errors.New("no encrypted executable processed")
	}
	return report, nil
}
