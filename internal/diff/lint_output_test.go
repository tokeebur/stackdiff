package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteLintResult_TextNoViolations(t *testing.T) {
	result := &LintResult{}
	var buf bytes.Buffer
	if err := WriteLintResult(&buf, result, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("expected 'no violations', got: %s", buf.String())
	}
}

func TestWriteLintResult_TextWithViolations(t *testing.T) {
	result := &LintResult{
		Violations: []LintViolation{
			{Address: "aws_instance.web", Rule: "some rule", Message: "something is wrong"},
		},
	}
	var buf bytes.Buffer
	if err := WriteLintResult(&buf, result, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "1 violation") {
		t.Errorf("expected violation count, got: %s", out)
	}
	if !strings.Contains(out, "some rule") {
		t.Errorf("expected rule name in output, got: %s", out)
	}
}

func TestWriteLintResult_MarkdownNoViolations(t *testing.T) {
	result := &LintResult{}
	var buf bytes.Buffer
	if err := WriteLintResult(&buf, result, "markdown"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("expected 'no violations', got: %s", buf.String())
	}
}

func TestWriteLintResult_MarkdownWithViolations(t *testing.T) {
	result := &LintResult{
		Violations: []LintViolation{
			{Address: "aws_s3_bucket.logs", Rule: "resource type must not be empty", Message: "aws_s3_bucket.logs: resource type is empty"},
		},
	}
	var buf bytes.Buffer
	if err := WriteLintResult(&buf, result, "markdown"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "## Lint Violations") {
		t.Errorf("expected markdown heading, got: %s", out)
	}
	if !strings.Contains(out, "aws_s3_bucket.logs") {
		t.Errorf("expected address in table, got: %s", out)
	}
}

func TestWriteLintResult_NilResult(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteLintResult(&buf, nil, "text"); err != nil {
		t.Fatalf("unexpected error for nil result: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil result, got: %s", buf.String())
	}
}
