package main

import (
	"flag"

	"github.com/your-org/stackdiff/internal/diff"
)

// ThresholdFlags holds raw CLI values for threshold configuration.
type ThresholdFlags struct {
	MaxAdded    int
	MaxRemoved  int
	MaxModified int
	MaxTotal    int
}

// ParseThresholdFlags registers and parses threshold-related CLI flags
// from the provided FlagSet.
func ParseThresholdFlags(fs *flag.FlagSet) *ThresholdFlags {
	f := &ThresholdFlags{}
	fs.IntVar(&f.MaxAdded, "max-added", 0, "max allowed added resources (0 = unlimited)")
	fs.IntVar(&f.MaxRemoved, "max-removed", 0, "max allowed removed resources (0 = unlimited)")
	fs.IntVar(&f.MaxModified, "max-modified", 0, "max allowed modified resources (0 = unlimited)")
	fs.IntVar(&f.MaxTotal, "max-total", 0, "max allowed total drifted resources (0 = unlimited)")
	return f
}

// BuildThresholdConfig converts parsed flags into a diff.ThresholdConfig.
func BuildThresholdConfig(f *ThresholdFlags) diff.ThresholdConfig {
	if f == nil {
		return diff.ThresholdConfig{}
	}
	return diff.ThresholdConfig{
		MaxAdded:    f.MaxAdded,
		MaxRemoved:  f.MaxRemoved,
		MaxModified: f.MaxModified,
		MaxTotal:    f.MaxTotal,
	}
}
