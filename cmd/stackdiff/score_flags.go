package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/your-org/stackdiff/internal/diff"
)

// ScoreConfig holds configuration parsed from score-related CLI flags.
type ScoreConfig struct {
	Enabled        bool
	WeightAdded    float64
	WeightRemoved  float64
	WeightModified float64
}

// ParseScoreFlags registers and parses score-related flags from fs.
func ParseScoreFlags(fs *flag.FlagSet) *ScoreConfig {
	c := &ScoreConfig{}
	fs.BoolVar(&c.Enabled, "score", false, "compute and display a drift risk score")
	fs.Float64Var(&c.WeightAdded, "score-weight-added", 1.0, "weight applied per added resource")
	fs.Float64Var(&c.WeightRemoved, "score-weight-removed", 2.0, "weight applied per removed resource")
	fs.Float64Var(&c.WeightModified, "score-weight-modified", 1.5, "weight applied per modified resource (multiplied by changed attribute count)")
	return c
}

// WriteScore computes and writes the drift score to w.
// Returns nil if scoring is disabled or the report is nil.
func WriteScore(w io.Writer, r *diff.Report, c *ScoreConfig) error {
	if c == nil || !c.Enabled {
		return nil
	}
	if r == nil {
		return nil
	}

	weights := diff.ScoreWeights{
		Added:    c.WeightAdded,
		Removed:  c.WeightRemoved,
		Modified: c.WeightModified,
	}

	s, err := diff.ScoreReport(r, weights)
	if err != nil {
		return fmt.Errorf("score: %w", err)
	}

	fmt.Fprintf(w, "\nDrift Risk Score\n")
	fmt.Fprintf(w, "  Total:    %.2f (%s)\n", s.Total, s.Label)
	fmt.Fprintf(w, "  Added:    %.2f\n", s.Added)
	fmt.Fprintf(w, "  Removed:  %.2f\n", s.Removed)
	fmt.Fprintf(w, "  Modified: %.2f\n", s.Modified)

	return nil
}
