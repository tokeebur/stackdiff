package diff

import (
	"strings"
	"testing"
)

func makeSummaryReport(actions ...string) *Report {
	r := &Report{}
	for i, a := range actions {
		addr := strings.ToLower(a) + fmt.Sprintf("_resource_%d", i)
		r.Changes = append(r.Changes, ResourceChange{
			Address:      addr,
			ResourceType: "aws_instance",
			Action:       a,
			Tags:         []string{"env:prod"},
		})
	}
	return r
}

func TestBuildDiffSummary_NoDrift(t *testing.T) {
	r := &Report{}
	opts := DefaultDiffSummaryOptions()
	out := BuildDiffSummary(r, Stats{}, nil, opts)
	if !strings.Contains(out, "No drift") {
		t.Errorf("expected no drift message, got: %s", out)
	}
}

func TestBuildDiffSummary_NilReport(t *testing.T) {
	out := BuildDiffSummary(nil, Stats{}, nil, DefaultDiffSummaryOptions())
	if !strings.Contains(out, "no report") {
		t.Errorf("expected nil report message, got: %s", out)
	}
}

func TestBuildDiffSummary_WithStats(t *testing.T) {
	r := makeSummaryReport("added", "removed", "modified")
	stats := Stats{Added: 1, Removed: 1, Modified: 1, Total: 3}
	opts := DefaultDiffSummaryOptions()
	out := BuildDiffSummary(r, stats, nil, opts)
	if !strings.Contains(out, "Added:    1") {
		t.Errorf("expected Added stat, got: %s", out)
	}
	if !strings.Contains(out, "Total:    3") {
		t.Errorf("expected Total stat, got: %s", out)
	}
}

func TestBuildDiffSummary_WithScore(t *testing.T) {
	r := makeSummaryReport("added")
	stats := Stats{Added: 1, Total: 1}
	score := &ScoreResult{Score: 42.5, Label: "medium"}
	opts := DefaultDiffSummaryOptions()
	out := BuildDiffSummary(r, stats, score, opts)
	if !strings.Contains(out, "42.5") {
		t.Errorf("expected score in output, got: %s", out)
	}
	if !strings.Contains(out, "medium") {
		t.Errorf("expected label in output, got: %s", out)
	}
}

func TestBuildDiffSummary_Compact(t *testing.T) {
	r := makeSummaryReport("added", "removed")
	stats := Stats{Added: 1, Removed: 1, Total: 2}
	opts := DefaultDiffSummaryOptions()
	opts.Compact = true
	out := BuildDiffSummary(r, stats, nil, opts)
	if strings.Contains(out, "[added]") {
		t.Errorf("compact mode should not list entries, got: %s", out)
	}
}

func TestBuildDiffSummary_ShowTags(t *testing.T) {
	r := makeSummaryReport("modified")
	stats := Stats{Modified: 1, Total: 1}
	opts := DefaultDiffSummaryOptions()
	opts.ShowTags = true
	out := BuildDiffSummary(r, stats, nil, opts)
	if !strings.Contains(out, "env:prod") {
		t.Errorf("expected tag in output, got: %s", out)
	}
}
