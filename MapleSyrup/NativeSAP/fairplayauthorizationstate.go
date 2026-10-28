package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type authorizedAccountState struct {
	Schema   string            `json:"schema"`
	Producer string            `json:"producer"`
	DSID     uint64            `json:"dsid"`
	Hardware []byte            `json:"hardware"`
	Files    map[string][]byte `json:"files"`
	Digests  map[string]string `json:"digests"`
}

func (s *authorizedAccountState) clear() {
	clear(s.Hardware)
	for _, b := range s.Files {
		clear(b)
	}
}
func validateAuthorizedState(s authorizedAccountState, dsid uint64, hardware []byte) error {
	if s.Schema != "crumb-local-authorization-state-v1" || s.Producer != "storeagent-10.9" || s.DSID != dsid || dsid == 0 || len(hardware) != 6 || !bytes.Equal(s.Hardware, hardware) || len(s.Files) == 0 || len(s.Files) > 32 || len(s.Files["SC Info.sidb"]) < 4 || len(s.Digests) != len(s.Files) {
		return fmt.Errorf("authorization_state_invalid")
	}
	total := 0
	for name, b := range s.Files {
		total += len(b)
		if name == "" || name == "." || name == ".." || filepath.Base(name) != name || strings.ContainsAny(name, "\\\x00") || total > 4<<20 {
			return fmt.Errorf("authorization_state_invalid")
		}
		sum := sha256.Sum256(b)
		if s.Digests[name] != hex.EncodeToString(sum[:]) {
			return fmt.Errorf("authorization_state_invalid")
		}
	}
	return nil
}
func readAuthorizedState(directory string, dsid uint64, hardware []byte) (result authorizedAccountState, err error) {
	return readAuthorizedStateFile(filepath.Join(directory, "account-state.json"), dsid, hardware)
}

func readAuthorizedStateFile(path string, dsid uint64, hardware []byte) (result authorizedAccountState, err error) {
	info, err := os.Lstat(path)
	if err != nil {
		return result, err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 || info.Size() > 8<<20 {
		return result, fmt.Errorf("authorization_state_invalid")
	}
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer file.Close()
	b, err := io.ReadAll(io.LimitReader(file, (8<<20)+1))
	if err != nil || len(b) > 8<<20 {
		return result, fmt.Errorf("authorization_state_invalid")
	}
	defer clear(b)
	if json.Unmarshal(b, &result) != nil {
		result.clear()
		return authorizedAccountState{}, fmt.Errorf("authorization_state_invalid")
	}
	if err = validateAuthorizedState(result, dsid, hardware); err != nil {
		result.clear()
		return authorizedAccountState{}, err
	}
	return result, nil
}

// Export only the validated licensing state. Passwords, cookies and account
// authentication tokens are not part of this schema. The copy is still private.
func exportAuthorizedState(directory, destination string, dsid uint64, hardware []byte) error {
	value, err := readAuthorizedState(directory, dsid, hardware)
	if err != nil {
		return fmt.Errorf("authorization_cached_state_invalid")
	}
	defer value.clear()
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("authorization_state_export_failed")
	}
	defer clear(data)
	file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("authorization_state_export_failed")
	}
	complete := false
	defer func() {
		file.Close()
		if !complete {
			os.Remove(destination)
		}
	}()
	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("authorization_state_export_failed")
	}
	if err = file.Sync(); err != nil {
		return fmt.Errorf("authorization_state_export_failed")
	}
	if err = file.Close(); err != nil {
		return fmt.Errorf("authorization_state_export_failed")
	}
	complete = true
	return nil
}
func saveAuthorizedState(directory string, dsid uint64, hardware []byte, files map[string][]byte) error {
	info, err := os.Lstat(directory)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm()&0077 != 0 {
		return fmt.Errorf("authorization_directory_invalid")
	}
	value := authorizedAccountState{Schema: "crumb-local-authorization-state-v1", Producer: "storeagent-10.9", DSID: dsid, Hardware: hardware, Files: files, Digests: map[string]string{}}
	for name, b := range files {
		sum := sha256.Sum256(b)
		value.Digests[name] = hex.EncodeToString(sum[:])
	}
	if err = validateAuthorizedState(value, dsid, hardware); err != nil {
		return err
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	defer clear(b)
	temp, err := os.CreateTemp(directory, ".account-state-")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()
	if _, err = temp.Write(b); err != nil {
		return err
	}
	if err = temp.Sync(); err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), filepath.Join(directory, "account-state.json"))
}
