package main

import (
	"flag"
	"testing"

	"github.com/your-org/stackdiff/internal/diff"
)

func TestParseExportFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseExportFlags(fs, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != diff.ExportText {
		t.Errorf("expected default format 'text', got %q", cfg.Format)
	}
	if cfg.Output != "-" {
		t.Errorf("expected default output '-', got %q", cfg.Output)
	}
}

func TestParseExportFlags_JSON(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseExportFlags(fs, []string{"--format", "json", "--output", "/tmp/out.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != diff.ExportJSON {
		t.Errorf("expected format 'json', got %q", cfg.Format)
	}
	if cfg.Output != "/tmp/out.json" {
		t.Errorf("expected output '/tmp/out.json', got %q", cfg.Output)
	}
}

func TestParseExportFlags_Markdown(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseExportFlags(fs, []string{"--format", "markdown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != diff.ExportMarkdown {
		t.Errorf("expected format 'markdown', got %q", cfg.Format)
	}
}

func TestParseExportFlags_InvalidFormat(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_, err := ParseExportFlags(fs, []string{"--format", "csv"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestApplyExportConfig_NilConfig(t *testing.T) {
	err := ApplyExportConfig(nil, &diff.Report{})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}
