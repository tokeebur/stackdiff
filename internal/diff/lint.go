package diff

import "fmt"

// LintRule describes a rule that can be checked against a report entry.
type LintRule struct {
	// Description is a human-readable explanation of the rule.
	Description string
	// Check returns a non-empty message if the entry violates the rule.
	Check func(entry ReportEntry) string
}

// LintViolation records a rule violation for a specific resource.
type LintViolation struct {
	Address     string
	Rule        string
	Message     string
}

// LintResult holds all violations found during a lint pass.
type LintResult struct {
	Violations []LintViolation
}

// HasViolations returns true if any rule violations were found.
func (r *LintResult) HasViolations() bool {
	return r != nil && len(r.Violations) > 0
}

// DefaultLintRules returns a set of sensible built-in lint rules.
func DefaultLintRules() []LintRule {
	return []LintRule{
		{
			Description: "address must not be empty",
			Check: func(e ReportEntry) string {
				if e.Address == "" {
					return "resource address is empty"
				}
				return ""
			},
		},
		{
			Description: "resource type must not be empty",
			Check: func(e ReportEntry) string {
				if e.ResourceType == "" {
					return fmt.Sprintf("%s: resource type is empty", e.Address)
				}
				return ""
			},
		},
		{
			Description: "modified entry must have at least one changed attribute",
			Check: func(e ReportEntry) string {
				if e.Action == ActionModified && len(e.ChangedAttributes) == 0 {
					return fmt.Sprintf("%s: modified resource has no changed attributes", e.Address)
				}
				return ""
			},
		},
	}
}

// LintReport applies the given rules to every entry in the report and
// returns a LintResult containing all violations found.
func LintReport(report *Report, rules []LintRule) *LintResult {
	result := &LintResult{}
	if report == nil {
		return result
	}
	for _, entry := range report.Entries {
		for _, rule := range rules {
			if msg := rule.Check(entry); msg != "" {
				result.Violations = append(result.Violations, LintViolation{
					Address: entry.Address,
					Rule:    rule.Description,
					Message: msg,
				})
			}
		}
	}
	return result
}
