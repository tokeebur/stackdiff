package diff

import (
	"testing"
)

func makeReport(added, removed []string, modified []ModifiedResource, unchanged int) Report {
	return Report{
		Added:     added,
		Removed:   removed,
		Modified:  modified,
		Unchanged: unchanged,
	}
}

func TestReport_HasDrift_False(t *testing.T) {
	r := makeReport(nil, nil, nil, 3)
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestReport_HasDrift_Added(t *testing.T) {
	r := makeReport([]string{"aws_instance.foo"}, nil, nil, 0)
	if !r.HasDrift() {
		t.Error("expected drift due to added resource")
	}
}

func TestReport_HasDrift_Removed(t *testing.T) {
	r := makeReport(nil, []string{"aws_instance.bar"}, nil, 0)
	if !r.HasDrift() {
		t.Error("expected drift due to removed resource")
	}
}

func TestReport_HasDrift_Modified(t *testing.T) {
	mod := ModifiedResource{
		Address: "aws_instance.baz",
		Changes: map[string]AttributeChange{
			"ami": {Old: "ami-old", New: "ami-new"},
		},
	}
	r := makeReport(nil, nil, []ModifiedResource{mod}, 1)
	if !r.HasDrift() {
		t.Error("expected drift due to modified resource")
	}
}

func TestReport_Summary_NoDrift(t *testing.T) {
	r := makeReport(nil, nil, nil, 5)
	if got := r.Summary(); got != "no drift detected" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestReport_Summary_Mixed(t *testing.T) {
	mod := ModifiedResource{Address: "aws_instance.x", Changes: map[string]AttributeChange{}}
	r := makeReport(
		[]string{"aws_s3_bucket.a"},
		[]string{"aws_s3_bucket.b", "aws_s3_bucket.c"},
		[]ModifiedResource{mod},
		2,
	)
	got := r.Summary()
	expected := "1 added, 2 removed, 1 modified"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
