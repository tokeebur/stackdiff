package main

import (
	"flag"

	"github.com/yourorg/stackdiff/internal/diff"
)

// NormalizeConfig holds parsed normalization flags.
type NormalizeConfig struct {
	TrimWhitespace bool
	LowercaseKeys  bool
	StripNullAttrs bool
}

// ParseNormalizeFlags reads normalization-related flags from the provided FlagSet.
func ParseNormalizeFlags(fs *flag.FlagSet) *NormalizeConfig {
	cfg := &NormalizeConfig{}
	fs.BoolVar(&cfg.TrimWhitespace, "normalize-trim", true, "trim whitespace from attribute values before comparison")
	fs.BoolVar(&cfg.LowercaseKeys, "normalize-lowercase-keys", false, "lowercase all attribute keys before comparison")
	fs.BoolVar(&cfg.StripNullAttrs, "normalize-strip-null", true, "strip null or empty attributes before comparison")
	return cfg
}

// ApplyNormalize applies normalization to the report if any option is enabled.
// Returns the original report unchanged if cfg is nil.
func ApplyNormalize(r *diff.Report, cfg *NormalizeConfig) *diff.Report {
	if cfg == nil || r == nil {
		return r
	}
	if !cfg.TrimWhitespace && !cfg.LowercaseKeys && !cfg.StripNullAttrs {
		return r
	}
	opts := diff.NormalizeOptions{
		TrimWhitespace: cfg.TrimWhitespace,
		LowercaseKeys:  cfg.LowercaseKeys,
		StripNullAttrs: cfg.StripNullAttrs,
	}
	return diff.NormalizeReport(r, opts)
}
