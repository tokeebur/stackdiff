package diff

import (
	"testing"
)

func baseRenameReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.old_web", ResourceType: "aws_instance", Action: ActionModify},
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: ActionAdd},
			{Address: "aws_instance.db", ResourceType: "aws_instance", Action: ActionRemove},
		},
	}
}

func TestApplyRenames_NilReport(t *testing.T) {
	result := ApplyRenames(nil, RenameConfig{})
	if result != nil {
		t.Fatal("expected nil for nil report")
	}
}

func TestApplyRenames_NoRules(t *testing.T) {
	r := baseRenameReport()
	result := ApplyRenames(r, RenameConfig{})
	if result != r {
		t.Fatal("expected same report when no rules provided")
	}
}

func TestApplyRenames_SingleMatch(t *testing.T) {
	r := baseRenameReport()
	cfg := RenameConfig{
		Rules: []RenameRule{
			{From: "aws_instance.old_web", To: "aws_instance.web"},
		},
	}
	result := ApplyRenames(r, cfg)
	if result.Changes[0].Address != "aws_instance.web" {
		t.Errorf("expected renamed address, got %s", result.Changes[0].Address)
	}
	// Others unchanged
	if result.Changes[1].Address != "aws_s3_bucket.data" {
		t.Errorf("unexpected change to unmatched entry: %s", result.Changes[1].Address)
	}
}

func TestApplyRenames_MultipleRules(t *testing.T) {
	r := baseRenameReport()
	cfg := RenameConfig{
		Rules: []RenameRule{
			{From: "aws_instance.old_web", To: "aws_instance.web"},
			{From: "aws_instance.db", To: "aws_instance.database"},
		},
	}
	result := ApplyRenames(r, cfg)
	if result.Changes[0].Address != "aws_instance.web" {
		t.Errorf("expected aws_instance.web, got %s", result.Changes[0].Address)
	}
	if result.Changes[2].Address != "aws_instance.database" {
		t.Errorf("expected aws_instance.database, got %s", result.Changes[2].Address)
	}
}

func TestApplyRenames_NoMatch(t *testing.T) {
	r := baseRenameReport()
	cfg := RenameConfig{
		Rules: []RenameRule{
			{From: "aws_instance.nonexistent", To: "aws_instance.new"},
		},
	}
	result := ApplyRenames(r, cfg)
	for i, ch := range result.Changes {
		if ch.Address != r.Changes[i].Address {
			t.Errorf("entry %d unexpectedly renamed to %s", i, ch.Address)
		}
	}
}

func TestApplyRenames_OriginalUnmodified(t *testing.T) {
	r := baseRenameReport()
	orig := r.Changes[0].Address
	cfg := RenameConfig{
		Rules: []RenameRule{
			{From: "aws_instance.old_web", To: "aws_instance.web"},
		},
	}
	ApplyRenames(r, cfg)
	if r.Changes[0].Address != orig {
		t.Error("original report was mutated")
	}
}
