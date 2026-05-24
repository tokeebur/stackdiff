package diff

import (
	"testing"
)

func makeGroupReport() *Report {
	return &Report{
		Entries: []ReportEntry{
			{Address: "aws_s3_bucket.logs", ResourceType: "aws_s3_bucket", Action: ActionAdded},
			{Address: "aws_s3_bucket.assets", ResourceType: "aws_s3_bucket", Action: ActionRemoved},
			{Address: "aws_iam_role.worker", ResourceType: "aws_iam_role", Action: ActionModified},
			{Address: "aws_iam_role.admin", ResourceType: "aws_iam_role", Action: ActionAdded},
			{Address: "aws_vpc.main", ResourceType: "aws_vpc", Action: ActionNoOp},
		},
	}
}

func TestGroupByType_NilReport(t *testing.T) {
	_, err := GroupByType(nil)
	if err == nil {
		t.Fatal("expected error for nil report, got nil")
	}
}

func TestGroupByType_GroupCount(t *testing.T) {
	g, err := GroupByType(makeGroupReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(g.Groups))
	}
}

func TestGroupByType_EntriesPerGroup(t *testing.T) {
	g, err := GroupByType(makeGroupReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Groups["aws_s3_bucket"]) != 2 {
		t.Errorf("expected 2 s3 entries, got %d", len(g.Groups["aws_s3_bucket"]))
	}
	if len(g.Groups["aws_iam_role"]) != 2 {
		t.Errorf("expected 2 iam_role entries, got %d", len(g.Groups["aws_iam_role"]))
	}
	if len(g.Groups["aws_vpc"]) != 1 {
		t.Errorf("expected 1 vpc entry, got %d", len(g.Groups["aws_vpc"]))
	}
}

func TestGroupByType_SortedWithinGroup(t *testing.T) {
	g, err := GroupByType(makeGroupReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s3 := g.Groups["aws_s3_bucket"]
	if s3[0].Address != "aws_s3_bucket.assets" {
		t.Errorf("expected aws_s3_bucket.assets first, got %s", s3[0].Address)
	}
	if s3[1].Address != "aws_s3_bucket.logs" {
		t.Errorf("expected aws_s3_bucket.logs second, got %s", s3[1].Address)
	}
}

func TestGroupByType_SortedTypes(t *testing.T) {
	g, err := GroupByType(makeGroupReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	types := g.SortedTypes()
	expected := []string{"aws_iam_role", "aws_s3_bucket", "aws_vpc"}
	for i, want := range expected {
		if types[i] != want {
			t.Errorf("index %d: expected %s, got %s", i, want, types[i])
		}
	}
}

func TestGroupByType_TotalEntries(t *testing.T) {
	g, err := GroupByType(makeGroupReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.TotalEntries() != 5 {
		t.Errorf("expected 5 total entries, got %d", g.TotalEntries())
	}
}

func TestGroupByType_EmptyReport(t *testing.T) {
	g, err := GroupByType(&Report{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Groups) != 0 {
		t.Errorf("expected 0 groups for empty report, got %d", len(g.Groups))
	}
	if g.TotalEntries() != 0 {
		t.Errorf("expected 0 total entries, got %d", g.TotalEntries())
	}
}
