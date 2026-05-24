package main

import (
	"flag"
	"os"

	"github.com/yourorg/stackdiff/internal/diff"
)

// TimelineConfig holds CLI options related to drift timeline tracking.
type TimelineConfig struct {
	SavePath string
	LoadPath string
	Label    string
}

// ParseTimelineFlags reads timeline-related flags from the given FlagSet.
func ParseTimelineFlags(fs *flag.FlagSet) *TimelineConfig {
	cfg := &TimelineConfig{}
	fs.StringVar(&cfg.SavePath, "timeline-save", "", "path to append drift snapshot to timeline file")
	fs.StringVar(&cfg.LoadPath, "timeline-load", "", "path to an existing timeline file to load and display")
	fs.StringVar(&cfg.Label, "timeline-label", "", "label for the current timeline snapshot")
	return cfg
}

// ApplyTimelineSave appends the current stats to the timeline file if configured.
func ApplyTimelineSave(cfg *TimelineConfig, stats diff.DriftStats) error {
	if cfg == nil || cfg.SavePath == "" {
		return nil
	}

	var tl *diff.Timeline
	if _, err := os.Stat(cfg.SavePath); err == nil {
		tl, err = diff.LoadTimelineFile(cfg.SavePath)
		if err != nil {
			return err
		}
	} else {
		tl = &diff.Timeline{}
	}

	tl.AddEntry(stats, cfg.Label)
	return diff.SaveTimelineFile(tl, cfg.SavePath)
}
