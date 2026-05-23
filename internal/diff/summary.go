package diff

import "fmt"

// Summary holds aggregated counts from a Report.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// String returns a one-line human-readable summary.
func (s Summary) String() string {
	return fmt.Sprintf(
		"Changes: %d total (+%d added, -%d removed, ~%d modified)",
		s.Total, s.Added, s.Removed, s.Modified,
	)
}

// HasDrift returns true when any changes are present.
func (s Summary) HasDrift() bool {
	return s.Total > 0
}

// Summarise computes a Summary from a Report.
func Summarise(r Report) Summary {
	var s Summary
	for _, rc := range r.ResourceChanges {
		switch rc.Action {
		case ActionAdded:
			s.Added++
		case ActionRemoved:
			s.Removed++
		case ActionModified:
			s.Modified++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified
	return s
}
