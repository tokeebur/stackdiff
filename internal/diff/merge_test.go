package diff

import (
	"testing"
)

func makeMergeRC(address, rtype, action string) ResourceChange {
	return ResourceChange{
		Address:      address,
		ResourceType: rtype,
		Action:       action,
	}
}

func TestMergeReports_BothNil(t *testing.T) {
	result := MergeReports(nil, nil)
	if result != nil {
		t.Fatalf("expected nil, got %+v", result)
	}
}

func TestMergeReports_BaseNil(t *testing.T) {
	overlay := &Report{Changes: []ResourceChange{makeMergeRC("a.b", "aws_s3_bucket", "added")}}
	result := MergeReports(nil, overlay)
	if result != overlay {
		t.Fatal("expected overlay to be returned unchanged")
	}
}

func TestMergeReports_OverlayNil(t *testing.T) {
	base := &Report{Changes: []ResourceChange{makeMergeRC("a.b", "aws_s3_bucket", "added")}}
	result := MergeReports(base, nil)
	if result != base {
		t.Fatal("expected base to be returned unchanged")
	}
}

func TestMergeReports_NoOverlap(t *testing.T) {
	base := &Report{Changes: []ResourceChange{
		makeMergeRC("res.one", "aws_instance", "added"),
	}}
	overlay := &Report{Changes: []ResourceChange{
		makeMergeRC("res.two", "aws_s3_bucket", "removed"),
	}}
	result := MergeReports(base, overlay)
	if len(result.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result.Changes))
	}
}

func TestMergeReports_OverlayWinsOnConflict(t *testing.T) {
	base := &Report{Changes: []ResourceChange{
		makeMergeRC("res.one", "aws_instance", "added"),
	}}
	overlay := &Report{Changes: []ResourceChange{
		makeMergeRC("res.one", "aws_instance", "modified"),
	}}
	result := MergeReports(base, overlay)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Action != "modified" {
		t.Errorf("expected overlay action 'modified', got %q", result.Changes[0].Action)
	}
}

func TestMergeReports_CombinesDistinctAndOverrides(t *testing.T) {
	base := &Report{Changes: []ResourceChange{
		makeMergeRC("res.one", "aws_instance", "added"),
		makeMergeRC("res.two", "aws_vpc", "removed"),
	}}
	overlay := &Report{Changes: []ResourceChange{
		makeMergeRC("res.one", "aws_instance", "modified"),
		makeMergeRC("res.three", "aws_s3_bucket", "added"),
	}}
	result := MergeReports(base, overlay)
	if len(result.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(result.Changes))
	}
	addressSet := make(map[string]string)
	for _, rc := range result.Changes {
		addressSet[rc.Address] = rc.Action
	}
	if addressSet["res.one"] != "modified" {
		t.Errorf("res.one should be 'modified', got %q", addressSet["res.one"])
	}
	if _, ok := addressSet["res.two"]; !ok {
		t.Error("res.two should be present from base")
	}
	if _, ok := addressSet["res.three"]; !ok {
		t.Error("res.three should be present from overlay")
	}
}
