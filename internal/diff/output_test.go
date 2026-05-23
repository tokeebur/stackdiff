package diff

import (
	"strings"
	"testing"
)

func emptyReport() *Report {
	return &Report{}
}

func driftReport() *Report {
	return &Report{
		Added: []ResourceChange{
			{Address: "aws_s3_bucket.new", ResourceType: "aws_s3_bucket"},
		},
		Removed: []ResourceChange{
			{Address: "aws_instance.old", ResourceType: "aws_instance"},
		},
		Modified: []ResourceChange{
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Changes: []AttributeChange{
					{Attribute: "instance_type", OldValue: "t2.micro", NewValue: "t3.small"},
				},
			},
		},
	}
}

func TestWriteReport_TextNoDrift(t *testing.T) {
	var sb strings.Builder
	if err := WriteReport(&sb, emptyReport(), FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", sb.String())
	}
}

func TestWriteReport_TextDrift(t *testing.T) {
	var sb strings.Builder
	if err := WriteReport(&sb, driftReport(), FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "[+] aws_s3_bucket.new") {
		t.Errorf("missing added resource line")
	}
	if !strings.Contains(out, "[-] aws_instance.old") {
		t.Errorf("missing removed resource line")
	}
	if !strings.Contains(out, "[~] aws_instance.web") {
		t.Errorf("missing modified resource line")
	}
	if !strings.Contains(out, "instance_type") {
		t.Errorf("missing attribute change")
	}
}

func TestWriteReport_MarkdownNoDrift(t *testing.T) {
	var sb strings.Builder
	if err := WriteReport(&sb, emptyReport(), FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No drift detected") {
		t.Errorf("expected no-drift message")
	}
}

func TestWriteReport_MarkdownDrift(t *testing.T) {
	var sb strings.Builder
	if err := WriteReport(&sb, driftReport(), FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "## Drift Summary") {
		t.Errorf("missing markdown header")
	}
	if !strings.Contains(out, "### Added") {
		t.Errorf("missing Added section")
	}
	if !strings.Contains(out, "### Removed") {
		t.Errorf("missing Removed section")
	}
	if !strings.Contains(out, "### Modified") {
		t.Errorf("missing Modified section")
	}
}

func TestWriteReport_UnsupportedFormat(t *testing.T) {
	var sb strings.Builder
	err := WriteReport(&sb, emptyReport(), OutputFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
