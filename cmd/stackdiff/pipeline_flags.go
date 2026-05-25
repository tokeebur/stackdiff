package main

import (
	"flag"

	"github.com/yourorg/stackdiff/internal/diff"
)

// PipelineConfig holds flags that control which pipeline steps are enabled.
type PipelineConfig struct {
	Normalize bool
	Dedupe    bool
	Redact    bool
	Truncate  int
}

// ParsePipelineFlags registers and parses pipeline-related flags from fs.
func ParsePipelineFlags(fs *flag.FlagSet) *PipelineConfig {
	cfg := &PipelineConfig{}
	fs.BoolVar(&cfg.Normalize, "pipeline-normalize", false, "normalize attribute values before output")
	fs.BoolVar(&cfg.Dedupe, "pipeline-dedupe", false, "remove duplicate resource change entries")
	fs.BoolVar(&cfg.Redact, "pipeline-redact", false, "redact sensitive attribute values")
	fs.IntVar(&cfg.Truncate, "pipeline-truncate", 0, "truncate report to N entries (0 = unlimited)")
	return cfg
}

// BuildPipeline constructs a diff.Pipeline from the given PipelineConfig.
// Steps are added in a deterministic order: normalize → dedupe → redact → truncate.
func BuildPipeline(cfg *PipelineConfig) *diff.Pipeline {
	p := diff.NewPipeline()
	if cfg == nil {
		return p
	}
	if cfg.Normalize {
		p.Add(func(r *diff.Report) *diff.Report {
			return diff.NormalizeReport(r, diff.NormalizeOptions{
				TrimWhitespace: true,
				StripNullAttrs: true,
				LowercaseKeys:  false,
			})
		})
	}
	if cfg.Dedupe {
		p.Add(func(r *diff.Report) *diff.Report {
			return diff.DedupeReport(r)
		})
	}
	if cfg.Redact {
		p.Add(func(r *diff.Report) *diff.Report {
			return diff.RedactReport(r, nil)
		})
	}
	if cfg.Truncate > 0 {
		limit := cfg.Truncate
		p.Add(func(r *diff.Report) *diff.Report {
			return diff.TruncateReport(r, limit)
		})
	}
	return p
}
