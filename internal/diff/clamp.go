package diff

import "fmt"

// ClampConfig controls how many attribute changes are shown per resource entry.
type ClampConfig struct {
	// MaxAttrsPerEntry limits the number of attribute changes shown per resource.
	// A value of 0 means no limit.
	MaxAttrsPerEntry int
}

// ClampResult holds metadata about what was trimmed during clamping.
type ClampResult struct {
	Report      *Report
	TrimmedKeys map[string]int // address -> number of trimmed attrs
}

// ClampReport limits the number of attribute changes shown per resource entry.
// Entries with more attributes than cfg.MaxAttrsPerEntry will have their
// Changed slice trimmed and the overflow count recorded in TrimmedKeys.
func ClampReport(r *Report, cfg ClampConfig) (*ClampResult, error) {
	if r == nil {
		return nil, fmt.Errorf("clamp: report must not be nil")
	}
	if cfg.MaxAttrsPerEntry <= 0 {
		return &ClampResult{Report: r, TrimmedKeys: map[string]int{}}, nil
	}

	trimmed := make(map[string]int)
	outEntries := make([]ResourceChange, 0, len(r.Changes))

	for _, entry := range r.Changes {
		if len(entry.ChangedAttrs) > cfg.MaxAttrsPerEntry {
			trimmed[entry.Address] = len(entry.ChangedAttrs) - cfg.MaxAttrsPerEntry
			clamped := entry
			clamped.ChangedAttrs = make([]AttrDiff, cfg.MaxAttrsPerEntry)
			copy(clamped.ChangedAttrs, entry.ChangedAttrs[:cfg.MaxAttrsPerEntry])
			outEntries = append(outEntries, clamped)
		} else {
			outEntries = append(outEntries, entry)
		}
	}

	return &ClampResult{
		Report:      &Report{Changes: outEntries},
		TrimmedKeys: trimmed,
	}, nil
}
