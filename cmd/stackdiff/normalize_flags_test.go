package main

import (
	"flag"
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
)

func TestParseNormalizeFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := ParseNormalizeFlags(fs)
	_ = fs.Parse([]string{})

	if !cfg.TrimWhitespace {
		t.Error("expected TrimWhitespace to default to true")
	}
	if cfg.LowercaseKeys {
		t.Error("expected LowercaseKeys to default to false")
	}
	if !cfg.StripNullAttrs {
		t.Error("expected StripNullAttrs to default to true")
	}
}

func TestParseNormalizeFlags_AllSet(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := ParseNormalizeFlags(fs)
	_ = fs.Parse([]string{
		"-normalize-trim=false",
		"-normalize-lowercase-keys=true",
		"-normalize-strip-null=false",
	})

	if cfg.TrimWhitespace {
		t.Error("expected TrimWhitespace to be false")
	}
	if !cfg.LowercaseKeys {
		t.Error("expected LowercaseKeys to be true")
	}
	if cfg.StripNullAttrs {
		t.Error("expected StripNullAttrs to be false")
	}
}

func TestApplyNormalize_NilConfig(t *testing.T) {
	r := &diff.Report{
		Changes: []diff.ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: "modified"},
		},
	}
	result := ApplyNormalize(r, nil)
	if result != r {
		t.Error("expected original report to be returned when cfg is nil")
	}
}

func TestApplyNormalize_NilReport(t *testing.T) {
	cfg := &NormalizeConfig{TrimWhitespace: true}
	result := ApplyNormalize(nil, cfg)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestApplyNormalize_NoOptionsEnabled(t *testing.T) {
	r := &diff.Report{
		Changes: []diff.ResourceChange{
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: "added"},
		},
	}
	cfg := &NormalizeConfig{TrimWhitespace: false, LowercaseKeys: false, StripNullAttrs: false}
	result := ApplyNormalize(r, cfg)
	if result != r {
		t.Error("expected original report returned when no options are enabled")
	}
}
