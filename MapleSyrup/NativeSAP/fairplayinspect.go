// Copied into the pinned internal/fairplay package by the native build.
package fairplay

// InspectPortableIPA reports geometry only, never selectors, licenses or paths.
func InspectPortableIPA(path string) (map[string]any, error) {
	inventory, archive, err := openIPAInventory(path)
	if err != nil {
		return nil, err
	}
	defer archive.Close()
	images := []map[string]any{}
	pages := uint64(0)
	for _, image := range inventory.images {
		count := uint64(image.layout.cryptSize) / pageBytes
		pages += count
		images = append(images, map[string]any{"cpu_type": image.layout.cpuType, "cpu_subtype": image.layout.cpuSubtype, "encrypted_pages": count, "slice_bytes": image.sliceSize})
	}
	return map[string]any{"encrypted_images": len(images), "encrypted_pages": pages, "sinf_count": len(inventory.sinfMembers), "supf_count": len(inventory.supfTables), "architectures": images, "enrollment_generation": "requires validated enrollment input", "retained_capture_recovery": "experimental"}, nil
}
