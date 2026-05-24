package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
)

func TestParseTimelineFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := ParseTimelineFlags(fs)
	_ = fs.Parse([]string{})

	if cfg.SavePath != "" {
		t.Errorf("expected empty SavePath, got %s", cfg.SavePath)
	}
	if cfg.LoadPath != "" {
		t.Errorf("expected empty LoadPath, got %s", cfg.LoadPath)
	}
	if cfg.Label != "" {
		t.Errorf("expected empty Label, got %s", cfg.Label)
	}
}

func TestParseTimelineFlags_AllSet(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := ParseTimelineFlags(fs)
	_ = fs.Parse([]string{"-timeline-save", "/tmp/tl.json", "-timeline-load", "/tmp/tl.json", "-timeline-label", "ci-run"})

	if cfg.SavePath != "/tmp/tl.json" {
		t.Errorf("unexpected SavePath: %s", cfg.SavePath)
	}
	if cfg.Label != "ci-run" {
		t.Errorf("unexpected Label: %s", cfg.Label)
	}
}

func TestApplyTimelineSave_NilConfig(t *testing.T) {
	if err := ApplyTimelineSave(nil, diff.DriftStats{}); err != nil {
		t.Fatalf("expected no error for nil config, got %v", err)
	}
}

func TestApplyTimelineSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "timeline.json")

	cfg := &TimelineConfig{SavePath: path, Label: "test-run"}
	stats := diff.DriftStats{Added: 1, Removed: 0, Modified: 2, Total: 3}

	if err := ApplyTimelineSave(cfg, stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestApplyTimelineSave_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "timeline.json")

	cfg := &TimelineConfig{SavePath: path, Label: "run-1"}
	if err := ApplyTimelineSave(cfg, diff.DriftStats{Total: 1}); err != nil {
		t.Fatalf("first save: %v", err)
	}

	cfg.Label = "run-2"
	if err := ApplyTimelineSave(cfg, diff.DriftStats{Total: 2}); err != nil {
		t.Fatalf("second save: %v", err)
	}

	tl, err := diff.LoadTimelineFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if tl.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", tl.Len())
	}
}
