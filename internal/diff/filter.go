package diff

import "strings"

// FilterOptions controls which resources appear in a Report.
type FilterOptions struct {
	// ResourceType, when non-empty, keeps only resources whose type matches.
	// e.g. "aws_instance"
	ResourceType string

	// AddressPrefix, when non-empty, keeps only resources whose address starts
	// with the given prefix. e.g. "module.vpc"
	AddressPrefix string

	// OnlyDrift, when true, keeps only resources that have changes (added,
	// removed, or modified). No-change resources are excluded.
	OnlyDrift bool
}

// FilterReport returns a new Report containing only the resource changes that
// satisfy all of the provided FilterOptions.
func FilterReport(r Report, opts FilterOptions) Report {
	filtered := Report{
		Added:    make([]ResourceChange, 0),
		Removed:  make([]ResourceChange, 0),
		Modified: make([]ResourceChange, 0),
	}

	for _, rc := range r.Added {
		if matchesFilter(rc, opts) {
			filtered.Added = append(filtered.Added, rc)
		}
	}
	for _, rc := range r.Removed {
		if matchesFilter(rc, opts) {
			filtered.Removed = append(filtered.Removed, rc)
		}
	}
	for _, rc := range r.Modified {
		if matchesFilter(rc, opts) {
			filtered.Modified = append(filtered.Modified, rc)
		}
	}

	return filtered
}

// IsEmpty reports whether opts has no filtering criteria set, meaning every
// resource would pass through unchanged.
func (opts FilterOptions) IsEmpty() bool {
	return opts.ResourceType == "" && opts.AddressPrefix == "" && !opts.OnlyDrift
}

func matchesFilter(rc ResourceChange, opts FilterOptions) bool {
	if opts.ResourceType != "" {
		if rc.ResourceType != opts.ResourceType {
			return false
		}
	}
	if opts.AddressPrefix != "" {
		if !strings.HasPrefix(rc.Address, opts.AddressPrefix) {
			return false
		}
	}
	return true
}
