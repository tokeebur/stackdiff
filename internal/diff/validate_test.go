package diff

import (
	"testing"
)

func validReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{
				Address:      "aws_s3_bucket.example",
				ResourceType: "aws_s3_bucket",
				Action:       ActionAdd,
			},
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Action:       ActionModify,
				Attributes: []AttributeDiff{
					{Key: "ami", OldValue: "old", NewValue: "new"},
				},
			},
		},
	}
}

func TestValidateReport_Valid(t *testing.T) {
	if err := ValidateReport(validReport()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateReport_Nil(t *testing.T) {
	err := ValidateReport(nil)
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestValidateReport_EmptyAddress(t *testing.T) {
	r := validReport()
	r.Changes[0].Address = ""
	err := ValidateReport(r)
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestValidateReport_EmptyResourceType(t *testing.T) {
	r := validReport()
	r.Changes[0].ResourceType = ""
	err := ValidateReport(r)
	if err == nil {
		t.Fatal("expected error for empty resource_type")
	}
}

func TestValidateReport_UnknownAction(t *testing.T) {
	r := validReport()
	r.Changes[0].Action = "replace"
	err := ValidateReport(r)
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestValidateReport_ModifyNoAttrs(t *testing.T) {
	r := validReport()
	r.Changes[1].Attributes = nil
	err := ValidateReport(r)
	if err == nil {
		t.Fatal("expected error for modify with no attribute diffs")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Messages) == 0 {
		t.Fatal("expected at least one validation message")
	}
}

func TestValidateReport_EmptyChanges(t *testing.T) {
	r := &Report{Changes: []ResourceChange{}}
	if err := ValidateReport(r); err != nil {
		t.Fatalf("expected no error for report with empty changes, got: %v", err)
	}
}
