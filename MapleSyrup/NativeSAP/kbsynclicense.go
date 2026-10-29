package machine

// LicenseBagsResult describes the whole download reply, not just the first
// successful import. Status is the first failing native status, or zero only
// when all supplied bags were accepted. Files always belong to the caller.
type LicenseBagsResult struct {
	Status           int32
	KeybagStatus     *int32
	DiversityStatus  *int32
	DiversityPresent bool
	Files            map[string][]byte
	Trace            map[string]any
}

type licenseBagImport func([]byte) (int32, map[string][]byte, map[string]any, error)

func consumeLicenseBags(keybag, diversity []byte, importKey, importDiversity licenseBagImport) (LicenseBagsResult, error) {
	result := LicenseBagsResult{DiversityPresent: len(diversity) != 0}
	status, files, trace, err := importKey(keybag)
	result.Status, result.Files, result.Trace = status, files, trace
	if err != nil {
		return result, err
	}
	keyStatus := status
	result.KeybagStatus = &keyStatus
	if status != 0 || len(diversity) == 0 {
		return result, nil
	}
	status, files, trace, err = importDiversity(diversity)
	// Neither intermediate snapshot may escape as a successful final state.
	for _, b := range result.Files {
		clear(b)
	}
	result.Status, result.Files, result.Trace = status, files, trace
	if err == nil {
		result.DiversityStatus = &status
	}
	return result, err
}
