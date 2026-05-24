package diff

// RenameRule describes a resource address rename mapping.
type RenameRule struct {
	From string
	To   string
}

// RenameConfig holds a set of rename rules to apply to a report.
type RenameConfig struct {
	Rules []RenameRule
}

// ApplyRenames rewrites resource addresses in the report according to the
// provided rename rules. Each rule maps an old address to a new address.
// Only entries whose address exactly matches a From value are renamed.
// The original report is not modified; a new Report is returned.
func ApplyRenames(r *Report, cfg RenameConfig) *Report {
	if r == nil {
		return nil
	}
	if len(cfg.Rules) == 0 {
		return r
	}

	// Build a lookup map for O(1) access.
	ruleMap := make(map[string]string, len(cfg.Rules))
	for _, rule := range cfg.Rules {
		if rule.From != "" && rule.To != "" {
			ruleMap[rule.From] = rule.To
		}
	}

	updated := make([]ResourceChange, 0, len(r.Changes))
	for _, entry := range r.Changes {
		if to, ok := ruleMap[entry.Address]; ok {
			entry.Address = to
		}
		updated = append(updated, entry)
	}

	return &Report{Changes: updated}
}
