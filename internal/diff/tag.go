package diff

import "sort"

// TagRule defines a rule for tagging report entries based on resource type or address prefix.
type TagRule struct {
	ResourceType  string
	AddressPrefix string
	Tag           string
}

// ApplyTags annotates each ResourceChange in the report with user-defined tags
// based on matching TagRules. Tags are appended to the entry's Tags slice.
func ApplyTags(report *Report, rules []TagRule) *Report {
	if report == nil {
		return nil
	}
	if len(rules) == 0 {
		return report
	}

	for i := range report.Changes {
		entry := &report.Changes[i]
		for _, rule := range rules {
			if matchesTagRule(entry, rule) {
				if !containsTag(entry.Tags, rule.Tag) {
					entry.Tags = append(entry.Tags, rule.Tag)
				}
			}
		}
		sort.Strings(entry.Tags)
	}
	return report
}

func matchesTagRule(entry *ResourceChange, rule TagRule) bool {
	if rule.ResourceType != "" && entry.ResourceType != rule.ResourceType {
		return false
	}
	if rule.AddressPrefix != "" && !hasPrefix(entry.Address, rule.AddressPrefix) {
		return false
	}
	return true
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
