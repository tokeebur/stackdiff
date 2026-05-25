package diff

import (
	"testing"
	"time"
)

func makeContextReport() *Report {
	return &Report{
		Entries: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: "modified"},
		},
	}
}

func TestNewDiffContext_NilMetadata(t *testing.T) {
	r := makeContextReport()
	dc := NewDiffContext(r, nil)
	if dc == nil {
		t.Fatal("expected non-nil DiffContext")
	}
	if dc.Metadata == nil {
		t.Fatal("expected default metadata to be set")
	}
	if dc.Metadata.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestNewDiffContext_WithMetadata(t *testing.T) {
	r := makeContextReport()
	meta := &ContextMetadata{
		RunID:       "run-123",
		Environment: "staging",
		TriggeredBy: "ci",
		Timestamp:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	dc := NewDiffContext(r, meta)
	if dc.Metadata.RunID != "run-123" {
		t.Errorf("expected RunID run-123, got %s", dc.Metadata.RunID)
	}
	if dc.Metadata.Environment != "staging" {
		t.Errorf("expected environment staging, got %s", dc.Metadata.Environment)
	}
}

func TestNewDiffContext_ZeroTimestampDefaulted(t *testing.T) {
	r := makeContextReport()
	meta := &ContextMetadata{RunID: "abc"}
	dc := NewDiffContext(r, meta)
	if dc.Metadata.Timestamp.IsZero() {
		t.Error("expected zero timestamp to be replaced")
	}
}

func TestDiffContext_WithLabel(t *testing.T) {
	r := makeContextReport()
	dc := NewDiffContext(r, nil)
	dc = dc.WithLabel("team", "platform")
	if dc.Metadata.Labels["team"] != "platform" {
		t.Errorf("expected label team=platform, got %v", dc.Metadata.Labels)
	}
}

func TestDiffContext_WithLabel_NilMetadata(t *testing.T) {
	dc := &DiffContext{}
	dc = dc.WithLabel("env", "prod")
	if dc.Metadata == nil {
		t.Fatal("expected metadata to be initialised")
	}
	if dc.Metadata.Labels["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", dc.Metadata.Labels)
	}
}

func TestDiffContext_WithMultipleLabels(t *testing.T) {
	r := makeContextReport()
	dc := NewDiffContext(r, nil)
	dc.WithLabel("a", "1").WithLabel("b", "2")
	if len(dc.Metadata.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(dc.Metadata.Labels))
	}
}
