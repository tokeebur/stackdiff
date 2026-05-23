package diff

import (
	"strings"
)

// IgnoreRule defines a rule for ignoring specific resources or attributes.
type IgnoreRule struct {
	// AddressPrefix ignores any resource whose address starts with this prefix.
	AddressPrefix string
	// ResourceType ignores any resource of this type.
	ResourceType string
	// AttributeKey ignores a specific attribute key across all modified resources.
	AttributeKey string
}

// ApplyIgnoreRules filters out changes from a Report based on the provided ignore rules.
// Resources that match an ignore rule are removed entirely; attribute-level rules
// strip matching keys from Modified entries and drop the entry if no attrs remain.
func ApplyIgnoreRules(r Report, rules []IgnoreRule) Report {
	if len(rules) == 0 {
		return r
	}

	out := Report{}

	for _, rc := range r.Added {
		if !resourceMatchesAny(rc, rules) {
			out.Added = append(out.Added, rc)
		}
	}

	for _, rc := range r.Removed {
		if !resourceMatchesAny(rc, rules) {
			out.Removed = append(out.Removed, rc)
		}
	}

	for _, rc := range r.Modified {
		if resourceMatchesAny(rc, rules) {
			continue
		}
		rc = stripIgnoredAttrs(rc, rules)
		if len(rc.Changes) > 0 {
			out.Modified = append(out.Modified, rc)
		}
	}

	return out
}

func resourceMatchesAny(rc ResourceChange, rules []IgnoreRule) bool {
	for _, rule := range rules {
		if rule.AddressPrefix != "" && strings.HasPrefix(rc.Address, rule.AddressPrefix) {
			return true
		}
		if rule.ResourceType != "" && rc.ResourceType == rule.ResourceType {
			return true
		}
	}
	return false
}

func stripIgnoredAttrs(rc ResourceChange, rules []IgnoreRule) ResourceChange {
	filtered := make(map[string]AttributeChange)
	for k, v := range rc.Changes {
		if !attrMatchesAny(k, rules) {
			filtered[k] = v
		}
	}
	rc.Changes = filtered
	return rc
}

func attrMatchesAny(key string, rules []IgnoreRule) bool {
	for _, rule := range rules {
		if rule.AttributeKey != "" && rule.AttributeKey == key {
			return true
		}
	}
	return false
}
