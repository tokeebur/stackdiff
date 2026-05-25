package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeExportContext() *DiffContext {
	r := &Report{
		Entries: []ResourceChange{
			{Address: "aws_s3_bucket.logs", ResourceType: "aws_s3_bucket", Action: "added"},
		},
	}
	meta := &ContextMetadata{
		RunID:       "run-42",
		Environment: "prod",
		TriggeredBy: "human",
		Timestamp:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Labels:      map[string]string{"region": "us-east-1"},
	}
	return NewDiffContext(r, meta)
}

func TestExportDiffContext_NilContext(t *testing.T) {
	var buf bytes.Buffer
	err := ExportDiffContext(nil, ContextExportJSON, &buf)
	if err == nil {
		t.Error("expected error for nil context")
	}
}

func TestExportDiffContext_InvalidFormat(t *testing.T) {
	dc := makeExportContext()
	var buf bytes.Buffer
	err := ExportDiffContext(dc, "xml", &buf)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportDiffContext_JSON(t *testing.T) {
	dc := makeExportContext()
	var buf bytes.Buffer
	err := ExportDiffContext(dc, ContextExportJSON, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["metadata"]; !ok {
		t.Error("expected metadata key in JSON output")
	}
}

func TestExportDiffContext_Text(t *testing.T) {
	dc := makeExportContext()
	var buf bytes.Buffer
	err := ExportDiffContext(dc, ContextExportText, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "run-42") {
		t.Error("expected RunID in text output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected environment in text output")
	}
	if !strings.Contains(out, "Entries: 1") {
		t.Error("expected entry count in text output")
	}
}
