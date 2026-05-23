package diff

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseExportFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected ExportFormat
	}{
		{"json", ExportJSON},
		{"JSON", ExportJSON},
		{"markdown", ExportMarkdown},
		{"Markdown", ExportMarkdown},
		{"text", ExportText},
		{"TEXT", ExportText},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseExportFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestParseExportFormat_Invalid(t *testing.T) {
	_, err := ParseExportFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExportReport_NilReport(t *testing.T) {
	err := ExportReport(nil, ExportText, "-", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestExportReport_JSONToWriter(t *testing.T) {
	r := &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: ActionAdded},
		},
	}
	var buf bytes.Buffer
	if err := ExportReport(r, ExportJSON, "-", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(decoded.Changes))
	}
}

func TestExportReport_TextToWriter(t *testing.T) {
	r := &Report{Changes: []ResourceChange{}}
	var buf bytes.Buffer
	if err := ExportReport(r, ExportText, "-", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in output, got: %s", buf.String())
	}
}

func TestExportReport_ToFile(t *testing.T) {
	r := &Report{Changes: []ResourceChange{}}
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "out", "report.json")
	if err := ExportReport(r, ExportJSON, path, os.Stdout); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %s", path)
	}
}
