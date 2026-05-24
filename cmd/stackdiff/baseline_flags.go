package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/stackdiff/internal/diff"
)

// BaselineConfig holds configuration for baseline save/load operations.
type BaselineConfig struct {
	SavePath string
	LoadPath string
}

// ParseBaselineFlags reads baseline-related flags from the provided FlagSet.
func ParseBaselineFlags(fs *flag.FlagSet) *BaselineConfig {
	cfg := &BaselineConfig{}
	fs.StringVar(&cfg.SavePath, "save-baseline", "", "path to save the current report as a baseline JSON file")
	fs.StringVar(&cfg.LoadPath, "load-baseline", "", "path to a previously saved baseline to compare against")
	return cfg
}

// ApplyBaselineSave writes the report to the configured save path, if set.
func ApplyBaselineSave(cfg *BaselineConfig, r *diff.Report) error {
	if cfg == nil || cfg.SavePath == "" {
		return nil
	}
	if err := diff.SaveBaselineFile(cfg.SavePath, r); err != nil {
		return fmt.Errorf("saving baseline: %w", err)
	}
	fmt.Fprintf(os.Stderr, "baseline saved to %s\n", cfg.SavePath)
	return nil
}

// LoadBaselineReport loads a report from the configured load path, if set.
// Returns nil without error when no load path is configured.
func LoadBaselineReport(cfg *BaselineConfig) (*diff.Report, error) {
	if cfg == nil || cfg.LoadPath == "" {
		return nil, nil
	}
	b, err := diff.LoadBaselineFile(cfg.LoadPath)
	if err != nil {
		return nil, fmt.Errorf("loading baseline: %w", err)
	}
	fmt.Fprintf(os.Stderr, "baseline loaded from %s (created %s)\n", cfg.LoadPath, b.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	return b.Report, nil
}
