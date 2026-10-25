#!/usr/bin/env python3
"""Apply the shared iOS/host compatibility fixes to the pinned FairPlay package."""
import pathlib, sys
p = pathlib.Path(sys.argv[1]); s = p.read_text()
s = s.replace('runtime.GOOS != "darwin" && runtime.GOOS != "linux"', 'runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "ios"')
s = s.replace('runtime.GOOS == "darwin" || runtime.GOOS == "windows"', 'runtime.GOOS == "darwin" || runtime.GOOS == "windows" || runtime.GOOS == "ios"')
p.write_text(s)
p = p.parent / 'decrypt.go'; s = p.read_text()
s = s.replace('record[:16], image.sliceOffset+uint64(image.layout.cryptOffset),', 'record[:16], image.sliceOffset+uint64(image.layout.cryptOffset)+uint64(index*pageBytes),')
s = s.replace('first := output[start : start+pageBytes]', 'auditStart := start\n    if isARM64InstructionSet(image.layout.cpuType) {\n        for auditStart < end {\n            _, words, _ := auditARM64FirstPage(output[auditStart:auditStart+pageBytes], uint64(auditStart), image.layout.instructionRanges)\n            if words > 0 { break }; auditStart += pageBytes\n        }\n        if auditStart >= end { return nil, fmt.Errorf("no encrypted instruction page") }\n    }\n    first := output[auditStart : auditStart+pageBytes]')
s = s.replace('auditARM64FirstPage(first, uint64(start),', 'auditARM64FirstPage(first, uint64(auditStart),')
p.write_text(s)
p = p.parent / 'decrypt_test.go'; s = p.read_text().replace('image.sliceOffset+uint64(layout.cryptOffset),', 'image.sliceOffset+uint64(layout.cryptOffset)+uint64(page*pageBytes),')
p.write_text(s)

# Keep all per-executable license selection in the reviewed companion package.
p = p.parent / 'decrypt.go'; s = p.read_text()
start = s.index('\tif len(inventory.sinfMembers) != 1 {')
end = s.index('\timageMap := map[string][]encryptedImage{}', start)
s = s[:start] + '\tfiles := map[string]*zip.File{}\n\tfor _, file := range source.File { files[file.Name] = file }\n\tmaterials, err := resolveImageMaterials(inventory, files, options.ModelPath)\n\tif err != nil { return result, err }\n\tvar material KeyMaterial\n\tif len(materials) == 1 {\n\t\tfor _, entry := range materials { material = entry.material }\n\t} else {\n\t\tdomains := map[string]any{}\n\t\tfor member, entry := range materials { domains[member] = map[string]any{"sinf_member": entry.sinfMember, "evidence": entry.material.Evidence} }\n\t\tmaterial.Evidence = map[string]any{"provider_kind": "per-executable-enrollment", "model_name": "per-executable", "domains": domains}\n\t}\n' + s[end:]
needle = '\t\ttable, ok := inventory.supfTables[images[0].supfMember]'
assert s.count(needle) == 1
s = s.replace(needle, '\t\tentry := materials[file.Name]\n\t\tmask := appMaskFromLogicalKey(entry.material.LogicalKey)\n\t\tderiver := newPageKeyDeriver(mask)\n' + needle)
s = s.replace('report, err := decryptImage(blob[image.sliceOffset:end], image, records, mask, deriver)', 'report, err := decryptImage(blob[image.sliceOffset:end], image, records, mask, deriver)')
needle = '\t\t\timageReports = append(imageReports, report)'
assert s.count(needle) == 1
s = s.replace(needle, '\t\t\treport.SINFMember = entry.sinfMember\n' + needle)
s = s.replace('result = DecryptResult{options.OutputPath, outputHash, options.ReportPath, len(imageReports), pages, digest(material.LogicalKey[:])}', 'logicalDigest := ""\n\tif len(materials) == 1 { logicalDigest = digest(material.LogicalKey[:]) }\n\tresult = DecryptResult{options.OutputPath, outputHash, options.ReportPath, len(imageReports), pages, logicalDigest}')
p.write_text(s)
p = p.parent / 'report.go'; s = p.read_text()
s = s.replace('type imageReport struct {', 'type imageReport struct {\n SINFMember string `json:"sinf_member"`')
p.write_text(s)
