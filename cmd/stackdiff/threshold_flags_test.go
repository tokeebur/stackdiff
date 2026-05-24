package main

import (
	"flag"
	"testing"
)

func TestParseThresholdFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	f := ParseThresholdFlags(fs)
	_ = fs.Parse([]string{})

	if f.MaxAdded != 0 || f.MaxRemoved != 0 || f.MaxModified != 0 || f.MaxTotal != 0 {
		t.Errorf("expected all defaults to be 0, got %+v", f)
	}
}

func TestParseThresholdFlags_AllSet(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	f := ParseThresholdFlags(fs)
	_ = fs.Parse([]string{
		"-max-added=3",
		"-max-removed=2",
		"-max-modified=5",
		"-max-total=10",
	})

	if f.MaxAdded != 3 {
		t.Errorf("expected MaxAdded=3, got %d", f.MaxAdded)
	}
	if f.MaxRemoved != 2 {
		t.Errorf("expected MaxRemoved=2, got %d", f.MaxRemoved)
	}
	if f.MaxModified != 5 {
		t.Errorf("expected MaxModified=5, got %d", f.MaxModified)
	}
	if f.MaxTotal != 10 {
		t.Errorf("expected MaxTotal=10, got %d", f.MaxTotal)
	}
}

func TestBuildThresholdConfig_NilInput(t *testing.T) {
	cfg := BuildThresholdConfig(nil)
	if cfg.MaxAdded != 0 || cfg.MaxTotal != 0 {
		t.Errorf("expected zero config from nil flags, got %+v", cfg)
	}
}

func TestBuildThresholdConfig_Values(t *testing.T) {
	f := &ThresholdFlags{MaxAdded: 4, MaxRemoved: 1, MaxModified: 2, MaxTotal: 7}
	cfg := BuildThresholdConfig(f)
	if cfg.MaxAdded != 4 || cfg.MaxRemoved != 1 || cfg.MaxModified != 2 || cfg.MaxTotal != 7 {
		t.Errorf("unexpected config values: %+v", cfg)
	}
}
