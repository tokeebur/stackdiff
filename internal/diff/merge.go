package diff

// MergeReports combines two Reports into one, with entries from `overlay`
// taking precedence over entries from `base` when the same address appears
// in both. Entries unique to either report are included as-is.
func MergeReports(base, overlay *Report) *Report {
	if base == nil && overlay == nil {
		return nil
	}
	if base == nil {
		return overlay
	}
	if overlay == nil {
		return base
	}

	// Index base entries by address.
	baseIndex := make(map[string]ResourceChange, len(base.Changes))
	for _, rc := range base.Changes {
		baseIndex[rc.Address] = rc
	}

	// Index overlay entries by address.
	overlayIndex := make(map[string]ResourceChange, len(overlay.Changes))
	for _, rc := range overlay.Changes {
		overlayIndex[rc.Address] = rc
	}

	// Build merged slice: overlay wins on conflict.
	seen := make(map[string]struct{})
	var merged []ResourceChange

	for _, rc := range overlay.Changes {
		merged = append(merged, rc)
		seen[rc.Address] = struct{}{}
	}

	for _, rc := range base.Changes {
		if _, exists := seen[rc.Address]; !exists {
			merged = append(merged, rc)
		}
	}

	_ = baseIndex
	_ = overlayIndex

	return &Report{Changes: merged}
}
