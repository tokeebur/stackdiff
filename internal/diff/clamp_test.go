package diff

import (
	"testing"
)

func makeClampReport(entries []ResourceChange) *Report {
	return &Report{Changes: entries}
}

func makeRC2(addr string, attrs []AttrDiff) ResourceChange {
	return ResourceChange{
		Address:      addr,
		ResourceType: "aws_instance",
		Action:       ActionModified,
		ChangedAttrs: attrs,
	}
}

func attrs(keys ...string) []AttrDiff {
	out := make([]AttrDiff, len(keys))
	for i, k := range keys {
		out[i] = AttrDiff{Key: k, OldValue: "a", NewValue: "b"}
	}
	return out
}

func TestClampReport_NilReport(t *testing.T) {
	_, err := ClampReport(nil, ClampConfig{MaxAttrsPerEntry: 2})
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestClampReport_NoLimit(t *testing.T) {
	r := makeClampReport([]ResourceChange{
		makeRC2("aws_instance.a", attrs("x", "y", "z")),
	})
	res, err := ClampReport(r, ClampConfig{MaxAttrsPerEntry: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Report.Changes[0].ChangedAttrs) != 3 {
		t.Errorf("expected 3 attrs, got %d", len(res.Report.Changes[0].ChangedAttrs))
	}
	if len(res.TrimmedKeys) != 0 {
		t.Errorf("expected no trimmed keys")
	}
}

func TestClampReport_BelowLimit(t *testing.T) {
	r := makeClampReport([]ResourceChange{
		makeRC2("aws_instance.a", attrs("x", "y")),
	})
	res, err := ClampReport(r, ClampConfig{MaxAttrsPerEntry: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Report.Changes[0].ChangedAttrs) != 2 {
		t.Errorf("expected 2 attrs, got %d", len(res.Report.Changes[0].ChangedAttrs))
	}
	if len(res.TrimmedKeys) != 0 {
		t.Errorf("expected no trimmed keys")
	}
}

func TestClampReport_ExceedsLimit(t *testing.T) {
	r := makeClampReport([]ResourceChange{
		makeRC2("aws_instance.a", attrs("a", "b", "c", "d", "e")),
	})
	res, err := ClampReport(r, ClampConfig{MaxAttrsPerEntry: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Report.Changes[0].ChangedAttrs) != 3 {
		t.Errorf("expected 3 attrs after clamp, got %d", len(res.Report.Changes[0].ChangedAttrs))
	}
	if res.TrimmedKeys["aws_instance.a"] != 2 {
		t.Errorf("expected 2 trimmed attrs, got %d", res.TrimmedKeys["aws_instance.a"])
	}
}

func TestClampReport_MultipleEntries(t *testing.T) {
	r := makeClampReport([]ResourceChange{
		makeRC2("aws_instance.a", attrs("x", "y", "z")),
		makeRC2("aws_instance.b", attrs("p")),
	})
	res, err := ClampReport(r, ClampConfig{MaxAttrsPerEntry: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Report.Changes[0].ChangedAttrs) != 2 {
		t.Errorf("entry a: expected 2 attrs, got %d", len(res.Report.Changes[0].ChangedAttrs))
	}
	if len(res.Report.Changes[1].ChangedAttrs) != 1 {
		t.Errorf("entry b: expected 1 attr, got %d", len(res.Report.Changes[1].ChangedAttrs))
	}
	if _, ok := res.TrimmedKeys["aws_instance.b"]; ok {
		t.Errorf("entry b should not appear in trimmed keys")
	}
}
