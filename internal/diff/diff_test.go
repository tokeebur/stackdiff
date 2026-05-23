package diff_test

import (
	"testing"

	"github.com/stackdiff/internal/diff"
	"github.com/stackdiff/internal/state"
)

func makeState(resources []state.Resource) *state.State {
	return &state.State{
		Version:          4,
		TerraformVersion: "1.5.0",
		Resources:        resources,
	}
}

func makeResource(addr, resType string, attrs map[string]interface{}) state.Resource {
	return state.Resource{
		Module: "",
		Type:   resType,
		Name:   addr,
		Instances: []state.Instance{
			{AttributeValues: attrs},
		},
	}
}

func TestCompare_NoChanges(t *testing.T) {
	attrs := map[string]interface{}{"id": "abc", "size": "t2.micro"}
	res := makeResource("web", "aws_instance", attrs)
	base := makeState([]state.Resource{res})
	target := makeState([]state.Resource{res})

	result, err := diff.Compare(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestCompare_AddedResource(t *testing.T) {
	attrs := map[string]interface{}{"id": "xyz"}
	base := makeState(nil)
	target := makeState([]state.Resource{makeResource("db", "aws_db_instance", attrs)})

	result, err := diff.Compare(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changes) != 1 || result.Changes[0].ChangeType != diff.Added {
		t.Errorf("expected 1 added change, got %+v", result.Changes)
	}
}

func TestCompare_RemovedResource(t *testing.T) {
	attrs := map[string]interface{}{"id": "old"}
	base := makeState([]state.Resource{makeResource("web", "aws_instance", attrs)})
	target := makeState(nil)

	result, err := diff.Compare(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changes) != 1 || result.Changes[0].ChangeType != diff.Removed {
		t.Errorf("expected 1 removed change, got %+v", result.Changes)
	}
}

func TestCompare_ModifiedResource(t *testing.T) {
	base := makeState([]state.Resource{makeResource("web", "aws_instance", map[string]interface{}{"size": "t2.micro"})})
	target := makeState([]state.Resource{makeResource("web", "aws_instance", map[string]interface{}{"size": "t3.large"})})

	result, err := diff.Compare(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changes) != 1 || result.Changes[0].ChangeType != diff.Modified {
		t.Errorf("expected 1 modified change, got %+v", result.Changes)
	}
}

func TestCompare_NilInputs(t *testing.T) {
	s := makeState(nil)
	if _, err := diff.Compare(nil, s); err == nil {
		t.Error("expected error for nil base")
	}
	if _, err := diff.Compare(s, nil); err == nil {
		t.Error("expected error for nil target")
	}
}
