package diff

import (
	"testing"
)

func makeScoreReport(added, removed, modified int) *Report {
	entries := []ReportEntry{}
	for i := 0; i < added; i++ {
		entries = append(entries, ReportEntry{
			Address:      fmt.Sprintf("aws_instance.added_%d", i),
			ResourceType: "aws_instance",
			Action:       ActionAdded,
		})
	}
	for i := 0; i < removed; i++ {
		entries = append(entries, ReportEntry{
			Address:      fmt.Sprintf("aws_instance.removed_%d", i),
			ResourceType: "aws_instance",
			Action:       ActionRemoved,
		})
	}
	for i := 0; i < modified; i++ {
		entries = append(entries, ReportEntry{
			Address:      fmt.Sprintf("aws_instance.modified_%d", i),
			ResourceType: "aws_instance",
			Action:       ActionModified,
			ChangedAttrs: map[string][2]string{"ami": {"old", "new"}},
		})
	}
	return &Report{Entries: entries}
}

func TestScoreReport_NilReport(t *testing.T) {
	_, err := ScoreReport(nil, DefaultScoreWeights())
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestScoreReport_NoDrift(t *testing.T) {
	r := makeScoreReport(0, 0, 0)
	s, err := ScoreReport(r, DefaultScoreWeights())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Total != 0 {
		t.Errorf("expected total 0, got %f", s.Total)
	}
	if s.Label != "none" {
		t.Errorf("expected label 'none', got %q", s.Label)
	}
}

func TestScoreReport_AddedOnly(t *testing.T) {
	r := makeScoreReport(3, 0, 0)
	s, err := ScoreReport(r, DefaultScoreWeights())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Added != 3.0 {
		t.Errorf("expected added score 3.0, got %f", s.Added)
	}
	if s.Label != "low" {
		t.Errorf("expected label 'low', got %q", s.Label)
	}
}

func TestScoreReport_RemovedWeighted(t *testing.T) {
	r := makeScoreReport(0, 4, 0)
	s, err := ScoreReport(r, DefaultScoreWeights())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 4 removals * 2.0 weight = 8.0
	if s.Removed != 8.0 {
		t.Errorf("expected removed score 8.0, got %f", s.Removed)
	}
	if s.Label != "medium" {
		t.Errorf("expected label 'medium', got %q", s.Label)
	}
}

func TestScoreReport_HighLabel(t *testing.T) {
	r := makeScoreReport(2, 5, 3)
	s, err := ScoreReport(r, DefaultScoreWeights())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Label != "high" {
		t.Errorf("expected label 'high', got %q", s.Label)
	}
}

func TestScoreReport_CustomWeights(t *testing.T) {
	r := makeScoreReport(1, 0, 0)
	w := ScoreWeights{Added: 5.0, Removed: 1.0, Modified: 1.0}
	s, err := ScoreReport(r, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Added != 5.0 {
		t.Errorf("expected added score 5.0, got %f", s.Added)
	}
}
