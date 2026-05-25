package diff

import (
	"testing"
)

func baseMaskReport() *Report {
	return &Report{
		Entries: []ResourceChange{
			{
				Address:      "aws_db_instance.main",
				ResourceType: "aws_db_instance",
				Action:       ActionModified,
				AttributeDiff: map[string]AttrDiff{
					"password":    {Old: "hunter2", New: "secret99"},
					"username":    {Old: "admin", New: "admin"},
					"db_name":     {Old: "prod", New: "prod2"},
					"api_key":     {Old: "key-abc", New: "key-xyz"},
				},
			},
		},
	}
}

func TestMaskReport_NilReport(t *testing.T) {
	result := MaskReport(nil, []MaskRule{{KeyPattern: "password", Replacement: "***"}})
	if result != nil {
		t.Fatal("expected nil")
	}
}

func TestMaskReport_NoRules(t *testing.T) {
	r := baseMaskReport()
	origVal := r.Entries[0].AttributeDiff["password"].Old
	MaskReport(r, nil)
	if r.Entries[0].AttributeDiff["password"].Old != origVal {
		t.Errorf("expected value unchanged, got %q", r.Entries[0].AttributeDiff["password"].Old)
	}
}

func TestMaskReport_MatchesPassword(t *testing.T) {
	r := baseMaskReport()
	MaskReport(r, []MaskRule{{KeyPattern: "password", Replacement: "[MASKED]"}})
	got := r.Entries[0].AttributeDiff["password"]
	if got.Old != "[MASKED]" || got.New != "[MASKED]" {
		t.Errorf("expected masked values, got old=%q new=%q", got.Old, got.New)
	}
}

func TestMaskReport_UnmatchedKeyPreserved(t *testing.T) {
	r := baseMaskReport()
	MaskReport(r, []MaskRule{{KeyPattern: "password", Replacement: "[MASKED]"}})
	got := r.Entries[0].AttributeDiff["db_name"]
	if got.Old != "prod" || got.New != "prod2" {
		t.Errorf("expected unmasked db_name, got old=%q new=%q", got.Old, got.New)
	}
}

func TestMaskReport_MultipleRules(t *testing.T) {
	r := baseMaskReport()
	rules := []MaskRule{
		{KeyPattern: "password", Replacement: "[MASKED]"},
		{KeyPattern: "api_key", Replacement: "[REDACTED]"},
	}
	MaskReport(r, rules)
	attrs := r.Entries[0].AttributeDiff
	if attrs["password"].Old != "[MASKED]" {
		t.Errorf("password not masked")
	}
	if attrs["api_key"].Old != "[REDACTED]" {
		t.Errorf("api_key not redacted")
	}
	if attrs["username"].Old != "admin" {
		t.Errorf("username should be unchanged")
	}
}

func TestMaskReport_InvalidPatternSkipped(t *testing.T) {
	r := baseMaskReport()
	// Invalid regex should not panic; valid rule should still apply.
	rules := []MaskRule{
		{KeyPattern: "[invalid", Replacement: "[BAD]"},
		{KeyPattern: "password", Replacement: "[MASKED]"},
	}
	MaskReport(r, rules)
	if r.Entries[0].AttributeDiff["password"].Old != "[MASKED]" {
		t.Errorf("expected password masked despite invalid preceding rule")
	}
}
