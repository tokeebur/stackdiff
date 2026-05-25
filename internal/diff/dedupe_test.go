package diff

import (
	"testing"
)

func makeDedupeReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: "modified"},
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: "added"},
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: "removed"},
			{Address: "aws_lambda.fn", ResourceType: "aws_lambda", Action: "added"},
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: "modified"},
		},
	}
}

func TestDedupeReport_NilReport(t *testing.T) {
	result := DedupeReport(nil)
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestDedupeReport_NoDuplicates(t *testing.T) {
	r := &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.a", ResourceType: "aws_instance", Action: "added"},
			{Address: "aws_instance.b", ResourceType: "aws_instance", Action: "removed"},
		},
	}
	result := DedupeReport(r)
	if len(result.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result.Changes))
	}
}

func TestDedupeReport_RemovesDuplicates(t *testing.T) {
	r := makeDedupeReport()
	result := DedupeReport(r)
	if len(result.Changes) != 3 {
		t.Fatalf("expected 3 unique changes, got %d", len(result.Changes))
	}
}

func TestDedupeReport_KeepsLastOccurrence(t *testing.T) {
	r := makeDedupeReport()
	result := DedupeReport(r)

	for _, rc := range result.Changes {
		if rc.Address == "aws_instance.web" && rc.Action != "removed" {
			t.Errorf("expected last occurrence action 'removed', got '%s'", rc.Action)
		}
		if rc.Address == "aws_s3_bucket.data" && rc.Action != "modified" {
			t.Errorf("expected last occurrence action 'modified', got '%s'", rc.Action)
		}
	}
}

func TestDedupeReport_EmptyChanges(t *testing.T) {
	r := &Report{Changes: []ResourceChange{}}
	result := DedupeReport(r)
	if len(result.Changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(result.Changes))
	}
}

func TestDedupeCount_NilReport(t *testing.T) {
	if n := DedupeCount(nil); n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}

func TestDedupeCount_CountsDuplicates(t *testing.T) {
	r := makeDedupeReport()
	if n := DedupeCount(r); n != 2 {
		t.Fatalf("expected 2 duplicates, got %d", n)
	}
}
