package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/yourorg/stackdiff/internal/diff"
)

// SummaryConfig holds parsed flags for the diff summary feature.
type SummaryConfig struct {
	ShowScore bool
	ShowStats bool
	ShowTags  bool
	Compact   bool
}

// ParseSummaryFlags reads summary-related flags from the provided FlagSet.
func ParseSummaryFlags(fs *flag.FlagSet) *SummaryConfig {
	cfg := &SummaryConfig{}
	fs.BoolVar(&cfg.ShowScore, "summary-score", true, "include drift score in summary")
	fs.BoolVar(&cfg.ShowStats, "summary-stats", true, "include counts in summary")
	fs.BoolVar(&cfg.ShowTags, "summary-tags", false, "include resource tags in summary")
	fs.BoolVar(&cfg.Compact, "summary-compact", false, "omit per-resource lines from summary")
	return cfg
}

// WriteDiffSummary computes and writes the diff summary to w.
func WriteDiffSummary(
	w io.Writer,
	r *diff.Report,
	cfg *SummaryConfig,
) error {
	if r == nil {
		return fmt.Errorf("report is nil")
	}

	var scoreCfg *SummaryConfig
	if cfg == nil {
		scoreCfg = &SummaryConfig{
			ShowScore: true,
			ShowStats: true,
		}
	} else {
		scoreCfg = cfg
	}

	stats := diff.ComputeStats(r)

	weights := diff.DefaultScoreWeights()
	scoreResult := diff.ScoreReport(r, weights)

	opts := diff.DiffSummaryOptions{
		ShowScore: scoreCfg.ShowScore,
		ShowStats: scoreCfg.ShowStats,
		ShowTags:  scoreCfg.ShowTags,
		Compact:   scoreCfg.Compact,
	}

	summary := diff.BuildDiffSummary(r, stats, scoreResult, opts)
	_, err := fmt.Fprint(w, summary)
	return err
}
