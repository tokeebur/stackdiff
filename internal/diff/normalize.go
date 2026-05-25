package diff

import (
	"strings"
)

// NormalizeOptions controls how attribute values are normalized before comparison.
type NormalizeOptions struct {
	TrimWhitespace bool
	LowercaseKeys  bool
	StripNullAttrs bool
}

// DefaultNormalizeOptions returns sensible defaults for normalization.
var DefaultNormalizeOptions = NormalizeOptions{
	TrimWhitespace: true,
	LowercaseKeys:  false,
	StripNullAttrs: true,
}

// NormalizeReport applies normalization rules to all attribute maps in the report.
// This reduces noise from whitespace or null-value drift.
func NormalizeReport(r *Report, opts NormalizeOptions) *Report {
	if r == nil {
		return nil
	}
	normalized := make([]ResourceChange, 0, len(r.Changes))
	for _, rc := range r.Changes {
		rc.Before = normalizeAttrs(rc.Before, opts)
		rc.After = normalizeAttrs(rc.After, opts)
		rc.Diff = normalizeAttrDiff(rc.Diff, opts)
		normalized = append(normalized, rc)
	}
	return &Report{Changes: normalized}
}

func normalizeAttrs(attrs map[string]string, opts NormalizeOptions) map[string]string {
	if attrs == nil {
		return nil
	}
	out := make(map[string]string, len(attrs))
	for k, v := range attrs {
		key := k
		if opts.LowercaseKeys {
			key = strings.ToLower(k)
		}
		val := v
		if opts.TrimWhitespace {
			val = strings.TrimSpace(v)
		}
		if opts.StripNullAttrs && (val == "null" || val == "") {
			continue
		}
		out[key] = val
	}
	return out
}

func normalizeAttrDiff(diff map[string][2]string, opts NormalizeOptions) map[string][2]string {
	if diff == nil {
		return nil
	}
	out := make(map[string][2]string, len(diff))
	for k, pair := range diff {
		key := k
		if opts.LowercaseKeys {
			key = strings.ToLower(k)
		}
		before := pair[0]
		after := pair[1]
		if opts.TrimWhitespace {
			before = strings.TrimSpace(before)
			after = strings.TrimSpace(after)
		}
		out[key] = [2]string{before, after}
	}
	return out
}
