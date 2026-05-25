package diff

import (
	"testing"
)

func makeTruncateReport(n int) *Report {
	r := &Report{}
	for i := 0; i < n; i++ {
		r.Changes = append(r.Changes, ResourceChange{
			Address:      fmt.Sprintf("aws_instance.res_%d", i),
			ResourceType: "aws_instance",
			Action:       ActionModified,
		})
	}
	return r
}

func TestTruncateReport_NilReport(t *testing.T) {
	_, err := TruncateReport(nil, TruncateOptions{MaxEntries: 5})
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestTruncateReport_NoLimit(t *testing.T) {
	r := makeTruncateReport(10)
	res, err := TruncateReport(r, TruncateOptions{MaxEntries: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Truncated {
		t.Error("expected Truncated=false when MaxEntries=0")
	}
	if res.Kept != 10 {
		t.Errorf("expected Kept=10, got %d", res.Kept)
	}
}

func TestTruncateReport_BelowLimit(t *testing.T) {
	r := makeTruncateReport(3)
	res, err := TruncateReport(r, TruncateOptions{MaxEntries: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Truncated {
		t.Error("expected Truncated=false when entries < MaxEntries")
	}
	if res.Total != 3 || res.Kept != 3 {
		t.Errorf("unexpected counts: total=%d kept=%d", res.Total, res.Kept)
	}
}

func TestTruncateReport_ExceedsLimit(t *testing.T) {
	r := makeTruncateReport(20)
	res, err := TruncateReport(r, TruncateOptions{MaxEntries: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Truncated {
		t.Error("expected Truncated=true")
	}
	if res.Kept != 5 {
		t.Errorf("expected Kept=5, got %d", res.Kept)
	}
	if res.Total != 20 {
		t.Errorf("expected Total=20, got %d", res.Total)
	}
	if len(res.Report.Changes) != 5 {
		t.Errorf("expected 5 changes in result, got %d", len(res.Report.Changes))
	}
}

func TestTruncateReport_MessageContent(t *testing.T) {
	r := makeTruncateReport(8)
	res, _ := TruncateReport(r, TruncateOptions{MaxEntries: 3})
	if res.Message == "" {
		t.Error("expected non-empty message")
	}
}
