package diff

import (
	"bytes"
	"encoding/json"
	"testing"
)

func makePatchReport() *Report {
	return &Report{
		Added: []ResourceChange{
			{Address: "aws_s3_bucket.new", ResourceType: "aws_s3_bucket", Attributes: map[string]string{"bucket": "my-bucket"}},
		},
		Removed: []ResourceChange{
			{Address: "aws_instance.old", ResourceType: "aws_instance", Attributes: map[string]string{"ami": "ami-123"}},
		},
		Modified: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Attributes: map[string]string{"instance_type": "t3.medium"}},
		},
	}
}

func TestBuildPatch_NilReport(t *testing.T) {
	_, err := BuildPatch(nil)
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestBuildPatch_EntryCount(t *testing.T) {
	p, err := BuildPatch(makePatchReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(p.Entries))
	}
}

func TestBuildPatch_Actions(t *testing.T) {
	p, err := BuildPatch(makePatchReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	actions := map[string]bool{}
	for _, e := range p.Entries {
		actions[e.Action] = true
	}
	for _, want := range []string{"add", "remove", "modify"} {
		if !actions[want] {
			t.Errorf("expected action %q in patch", want)
		}
	}
}

func TestBuildPatch_SortedByAddress(t *testing.T) {
	p, err := BuildPatch(makePatchReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(p.Entries); i++ {
		if p.Entries[i].Address < p.Entries[i-1].Address {
			t.Errorf("entries not sorted: %s before %s", p.Entries[i-1].Address, p.Entries[i].Address)
		}
	}
}

func TestWritePatchJSON_ValidOutput(t *testing.T) {
	p, _ := BuildPatch(makePatchReport())
	var buf bytes.Buffer
	if err := WritePatchJSON(&buf, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Patch
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded.Entries) != 3 {
		t.Errorf("expected 3 decoded entries, got %d", len(decoded.Entries))
	}
}

func TestBuildPatch_EmptyReport(t *testing.T) {
	p, err := BuildPatch(&Report{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries) != 0 {
		t.Errorf("expected 0 entries for empty report, got %d", len(p.Entries))
	}
}
