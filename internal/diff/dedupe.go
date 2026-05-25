package diff

// DedupeReport removes duplicate entries from a Report, keeping the last
// occurrence of any entry with the same Address. Entries are considered
// duplicates when their Address fields are identical.
func DedupeReport(r *Report) *Report {
	if r == nil {
		return nil
	}

	seen := make(map[string]int) // address -> last index in result
	result := make([]ResourceChange, 0, len(r.Changes))

	for _, rc := range r.Changes {
		if idx, exists := seen[rc.Address]; exists {
			// Overwrite the previous occurrence
			result[idx] = rc
		} else {
			seen[rc.Address] = len(result)
			result = append(result, rc)
		}
	}

	return &Report{
		Changes: result,
	}
}

// DedupeCount returns the number of duplicate entries that would be removed
// from the report without modifying it.
func DedupeCount(r *Report) int {
	if r == nil {
		return 0
	}

	seen := make(map[string]struct{})
	duplicates := 0

	for _, rc := range r.Changes {
		if _, exists := seen[rc.Address]; exists {
			duplicates++
		} else {
			seen[rc.Address] = struct{}{}
		}
	}

	return duplicates
}
