package diff

import (
	"testing"
)

func makeSortReport() Report {
	return Report{
		ResourceChanges: []ResourceChange{
			{Address: "module.vpc.aws_subnet.private", ResourceType: "aws_subnet", Action: ActionModified},
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: ActionAdded},
			{Address: "aws_s3_bucket.logs", ResourceType: "aws_s3_bucket", Action: ActionRemoved},
			{Address: "aws_instance.db", ResourceType: "aws_instance", Action: ActionModified},
		},
	}
}

func TestSortReport_ByAddress(t *testing.T) {
	r := SortReport(makeSortReport(), SortByAddress)
	addrs := make([]string, len(r.ResourceChanges))
	for i, rc := range r.ResourceChanges {
		addrs[i] = rc.Address
	}
	expected := []string{
		"aws_instance.db",
		"aws_instance.web",
		"aws_s3_bucket.logs",
		"module.vpc.aws_subnet.private",
	}
	for i, a := range expected {
		if addrs[i] != a {
			t.Errorf("index %d: got %q, want %q", i, addrs[i], a)
		}
	}
}

func TestSortReport_ByType(t *testing.T) {
	r := SortReport(makeSortReport(), SortByType)
	if r.ResourceChanges[0].ResourceType != "aws_instance" {
		t.Errorf("expected aws_instance first, got %q", r.ResourceChanges[0].ResourceType)
	}
	if r.ResourceChanges[2].ResourceType != "aws_s3_bucket" {
		t.Errorf("expected aws_s3_bucket at index 2, got %q", r.ResourceChanges[2].ResourceType)
	}
}

func TestSortReport_ByAction(t *testing.T) {
	r := SortReport(makeSortReport(), SortByAction)
	if r.ResourceChanges[0].Action != ActionAdded {
		t.Errorf("expected first action to be Added, got %q", r.ResourceChanges[0].Action)
	}
	if r.ResourceChanges[1].Action != ActionRemoved {
		t.Errorf("expected second action to be Removed, got %q", r.ResourceChanges[1].Action)
	}
}

func TestSortReport_UnknownOrder_PreservesLength(t *testing.T) {
	r := SortReport(makeSortReport(), SortOrder("unknown"))
	if len(r.ResourceChanges) != 4 {
		t.Errorf("expected 4 changes, got %d", len(r.ResourceChanges))
	}
}

func TestSortReport_DoesNotMutateOriginal(t *testing.T) {
	orig := makeSortReport()
	firstAddr := orig.ResourceChanges[0].Address
	SortReport(orig, SortByAddress)
	if orig.ResourceChanges[0].Address != firstAddr {
		t.Error("SortReport mutated the original report")
	}
}
