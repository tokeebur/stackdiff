package diff

import "fmt"

// TruncateOptions controls how a Report is truncated.
type TruncateOptions struct {
	// MaxEntries is the maximum number of ResourceChange entries to keep.
	// Zero or negative means no limit.
	MaxEntries int
}

// TruncateResult holds the output of a truncation operation.
type TruncateResult struct {
	Report    *Report
	Total     int
	Kept      int
	Truncated bool
	Message   string
}

// TruncateReport limits the number of entries in a Report according to opts.
// Entries are kept in their current order (apply SortReport first if ordering
// matters). A nil report returns an error.
func TruncateReport(r *Report, opts TruncateOptions) (*TruncateResult, error) {
	if r == nil {
		return nil, fmt.Errorf("truncate: report must not be nil")
	}

	total := len(r.Changes)

	if opts.MaxEntries <= 0 || total <= opts.MaxEntries {
		return &TruncateResult{
			Report:    r,
			Total:     total,
			Kept:      total,
			Truncated: false,
			Message:   fmt.Sprintf("showing all %d entries", total),
		}, nil
	}

	kept := opts.MaxEntries
	truncated := &Report{
		Changes: make([]ResourceChange, kept),
	}
	copy(truncated.Changes, r.Changes[:kept])

	return &TruncateResult{
		Report:    truncated,
		Total:     total,
		Kept:      kept,
		Truncated: true,
		Message:   fmt.Sprintf("showing %d of %d entries (truncated)", kept, total),
	}, nil
}
