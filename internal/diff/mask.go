package diff

import (
	"regexp"
	"strings"
)

// MaskRule defines a pattern-based masking rule for attribute values.
type MaskRule struct {
	// KeyPattern is a regex matched against attribute keys.
	KeyPattern string
	// Replacement is the string used to replace matched values.
	Replacement string
}

// MaskReport replaces attribute values whose keys match any MaskRule pattern
// with the rule's Replacement string. Unlike RedactReport (which uses a fixed
// redaction marker), MaskReport allows callers to supply arbitrary replacement
// strings — useful for stubbing values in test output or CI summaries.
//
// The original report is mutated in place and also returned for convenience.
// A nil report is returned unchanged.
func MaskReport(report *Report, rules []MaskRule) *Report {
	if report == nil {
		return nil
	}
	if len(rules) == 0 {
		return report
	}

	compiled := compileRules(rules)

	for i := range report.Entries {
		e := &report.Entries[i]
		e.AttributeDiff = maskAttrMap(e.AttributeDiff, compiled)
	}
	return report
}

type compiledRule struct {
	re          *regexp.Regexp
	replacement string
}

func compileRules(rules []MaskRule) []compiledRule {
	out := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile(r.KeyPattern)
		if err != nil {
			// Skip invalid patterns rather than panicking.
			continue
		}
		out = append(out, compiledRule{re: re, replacement: r.Replacement})
	}
	return out
}

func maskAttrMap(attrs map[string]AttrDiff, rules []compiledRule) map[string]AttrDiff {
	if len(attrs) == 0 {
		return attrs
	}
	result := make(map[string]AttrDiff, len(attrs))
	for k, v := range attrs {
		masked := false
		for _, r := range rules {
			if r.re.MatchString(strings.ToLower(k)) {
				result[k] = AttrDiff{Old: r.replacement, New: r.replacement}
				masked = true
				break
			}
		}
		if !masked {
			result[k] = v
		}
	}
	return result
}
