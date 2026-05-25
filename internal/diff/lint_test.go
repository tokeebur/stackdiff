package diff

import "testing"

func baseLintReport() *Report {
	return &Report{
		Entries: []ReportEntry{
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Action:       ActionAdded,
			},
			{
				Address:      "aws_s3_bucket.data",
				ResourceType: "aws_s3_bucket",
				Action:       ActionRemoved,
			},
		},
	}
}

func TestLintReport_NilReport(t *testing.T) {
	result := LintReport(nil, DefaultLintRules())
	if result.HasViolations() {
		t.Errorf("expected no violations for nil report, got %d", len(result.Violations))
	}
}

func TestLintReport_NoViolations(t *testing.T) {
	report := baseLintReport()
	result := LintReport(report, DefaultLintRules())
	if result.HasViolations() {
		t.Errorf("expected no violations, got: %+v", result.Violations)
	}
}

func TestLintReport_EmptyAddress(t *testing.T) {
	report := &Report{
		Entries: []ReportEntry{
			{Address: "", ResourceType: "aws_instance", Action: ActionAdded},
		},
	}
	result := LintReport(report, DefaultLintRules())
	if !result.HasViolations() {
		t.Fatal("expected a violation for empty address")
	}
	found := false
	for _, v := range result.Violations {
		if v.Rule == "address must not be empty" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'address must not be empty' violation, got: %+v", result.Violations)
	}
}

func TestLintReport_EmptyResourceType(t *testing.T) {
	report := &Report{
		Entries: []ReportEntry{
			{Address: "aws_instance.web", ResourceType: "", Action: ActionAdded},
		},
	}
	result := LintReport(report, DefaultLintRules())
	if !result.HasViolations() {
		t.Fatal("expected a violation for empty resource type")
	}
}

func TestLintReport_ModifiedWithNoChangedAttrs(t *testing.T) {
	report := &Report{
		Entries: []ReportEntry{
			{
				Address:           "aws_instance.web",
				ResourceType:      "aws_instance",
				Action:            ActionModified,
				ChangedAttributes: map[string]AttributeChange{},
			},
		},
	}
	result := LintReport(report, DefaultLintRules())
	if !result.HasViolations() {
		t.Fatal("expected violation for modified entry with no changed attributes")
	}
}

func TestLintReport_CustomRule(t *testing.T) {
	customRule := LintRule{
		Description: "address must not start with underscore",
		Check: func(e ReportEntry) string {
			if len(e.Address) > 0 && e.Address[0] == '_' {
				return e.Address + ": address starts with underscore"
			}
			return ""
		},
	}
	report := &Report{
		Entries: []ReportEntry{
			{Address: "_bad_resource", ResourceType: "aws_instance", Action: ActionAdded},
		},
	}
	result := LintReport(report, []LintRule{customRule})
	if !result.HasViolations() {
		t.Fatal("expected violation from custom rule")
	}
	if result.Violations[0].Rule != "address must not start with underscore" {
		t.Errorf("unexpected rule name: %s", result.Violations[0].Rule)
	}
}
