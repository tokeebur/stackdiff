package diff

import "strings"

// RedactConfig holds patterns for attribute keys whose values should be redacted.
type RedactConfig struct {
	// KeyPatterns is a list of substrings; any attribute key containing one of
	// these substrings (case-insensitive) will have its value replaced.
	KeyPatterns []string
	// Replacement is the string used in place of the real value.
	// Defaults to "[redacted]" if empty.
	Replacement string
}

// RedactReport replaces sensitive attribute values in all report entries
// according to the provided RedactConfig. The original report is mutated in
// place and returned for convenience.
func RedactReport(r *Report, cfg RedactConfig) *Report {
	if r == nil {
		return nil
	}
	if cfg.Replacement == "" {
		cfg.Replacement = "[redacted]"
	}

	for i := range r.Entries {
		e := &r.Entries[i]
		e.AttributeDiff = redactAttrMap(e.AttributeDiff, cfg)
	}
	return r
}

// redactAttrMap returns a new map with sensitive values replaced.
func redactAttrMap(attrs map[string]AttrDiff, cfg RedactConfig) map[string]AttrDiff {
	if attrs == nil {
		return nil
	}
	out := make(map[string]AttrDiff, len(attrs))
	for k, v := range attrs {
		if isSensitiveKey(k, cfg.KeyPatterns) {
			v.OldValue = cfg.Replacement
			v.NewValue = cfg.Replacement
		}
		out[k] = v
	}
	return out
}

// isSensitiveKey returns true when the key contains any of the given patterns
// (case-insensitive comparison).
func isSensitiveKey(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
