package diff

import (
	"testing"
)

func baseRedactReport() *Report {
	return &Report{
		Entries: []ReportEntry{
			{
				Address:      "aws_db_instance.main",
				ResourceType: "aws_db_instance",
				Action:       ActionModified,
				AttributeDiff: map[string]AttrDiff{
					"username":   {OldValue: "admin", NewValue: "root"},
					"password":   {OldValue: "secret123", NewValue: "newsecret"},
					"db_name":    {OldValue: "mydb", NewValue: "mydb"},
					"secret_key": {OldValue: "abc", NewValue: "xyz"},
				},
			},
		},
	}
}

func TestRedactReport_NilReport(t *testing.T) {
	result := RedactReport(nil, RedactConfig{KeyPatterns: []string{"password"}})
	if result != nil {
		t.Error("expected nil for nil report")
	}
}

func TestRedactReport_NoPatterns(t *testing.T) {
	r := baseRedactReport()
	RedactReport(r, RedactConfig{})
	// No patterns — values should be unchanged.
	if r.Entries[0].AttributeDiff["password"].OldValue != "secret123" {
		t.Errorf("expected password unchanged, got %q", r.Entries[0].AttributeDiff["password"].OldValue)
	}
}

func TestRedactReport_MatchesPassword(t *testing.T) {
	r := baseRedactReport()
	RedactReport(r, RedactConfig{KeyPatterns: []string{"password"}})
	ad := r.Entries[0].AttributeDiff
	if ad["password"].OldValue != "[redacted]" {
		t.Errorf("expected [redacted], got %q", ad["password"].OldValue)
	}
	if ad["password"].NewValue != "[redacted]" {
		t.Errorf("expected [redacted], got %q", ad["password"].NewValue)
	}
	// Non-sensitive key should be untouched.
	if ad["username"].OldValue != "admin" {
		t.Errorf("expected username untouched, got %q", ad["username"].OldValue)
	}
}

func TestRedactReport_MultiplePatterns(t *testing.T) {
	r := baseRedactReport()
	RedactReport(r, RedactConfig{KeyPatterns: []string{"password", "secret"}})
	ad := r.Entries[0].AttributeDiff
	for _, key := range []string{"password", "secret_key"} {
		if ad[key].OldValue != "[redacted]" {
			t.Errorf("key %q: expected [redacted], got %q", key, ad[key].OldValue)
		}
	}
	if ad["db_name"].OldValue != "mydb" {
		t.Errorf("db_name should be untouched")
	}
}

func TestRedactReport_CustomReplacement(t *testing.T) {
	r := baseRedactReport()
	RedactReport(r, RedactConfig{KeyPatterns: []string{"password"}, Replacement: "***"})
	if got := r.Entries[0].AttributeDiff["password"].OldValue; got != "***" {
		t.Errorf("expected ***, got %q", got)
	}
}

func TestRedactReport_CaseInsensitiveKey(t *testing.T) {
	r := &Report{
		Entries: []ReportEntry{
			{
				Address:      "aws_iam_user.dev",
				ResourceType: "aws_iam_user",
				Action:       ActionModified,
				AttributeDiff: map[string]AttrDiff{
					"ACCESS_KEY": {OldValue: "AKIA123", NewValue: "AKIA456"},
				},
			},
		},
	}
	RedactReport(r, RedactConfig{KeyPatterns: []string{"access_key"}})
	if got := r.Entries[0].AttributeDiff["ACCESS_KEY"].OldValue; got != "[redacted]" {
		t.Errorf("expected case-insensitive match, got %q", got)
	}
}
