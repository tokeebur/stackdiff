package diff

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func makeBaselineReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Action:       ActionModified,
				Before:       map[string]string{"ami": "ami-old"},
				After:        map[string]string{"ami": "ami-new"},
			},
		},
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	r := makeBaselineReport()
	var buf bytes.Buffer
	if err := SaveBaseline(&buf, r); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}
	b, err := LoadBaseline(&buf)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if len(b.Report.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(b.Report.Changes))
	}
	if b.Report.Changes[0].Address != "aws_instance.web" {
		t.Errorf("unexpected address: %s", b.Report.Changes[0].Address)
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSaveBaseline_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := SaveBaseline(&buf, nil); err == nil {
		t.Error("expected error for nil report")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	buf := bytes.NewBufferString(`{invalid}`)
	if _, err := LoadBaseline(buf); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadBaseline_MissingReport(t *testing.T) {
	buf := bytes.NewBufferString(`{"created_at":"2024-01-01T00:00:00Z"}`)
	if _, err := LoadBaseline(buf); err == nil {
		t.Error("expected error when report field is missing")
	}
}

func TestSaveAndLoadBaselineFile_RoundTrip(t *testing.T) {
	r := makeBaselineReport()
	tmp := filepath.Join(t.TempDir(), "baseline.json")
	if err := SaveBaselineFile(tmp, r); err != nil {
		t.Fatalf("SaveBaselineFile: %v", err)
	}
	b, err := LoadBaselineFile(tmp)
	if err != nil {
		t.Fatalf("LoadBaselineFile: %v", err)
	}
	if len(b.Report.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(b.Report.Changes))
	}
}

func TestLoadBaselineFile_MissingFile(t *testing.T) {
	if _, err := LoadBaselineFile("/nonexistent/path/baseline.json"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveBaselineFile_BadPath(t *testing.T) {
	r := makeBaselineReport()
	if err := SaveBaselineFile("/nonexistent/dir/baseline.json", r); err == nil {
		t.Error("expected error for bad path")
	}
}

func TestLoadBaselineFile_NotBaseline(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte(`"just a string"`), 0644)
	if _, err := LoadBaselineFile(tmp); err == nil {
		t.Error("expected error for non-baseline JSON")
	}
}
