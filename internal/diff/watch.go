package diff

import (
	"fmt"
	"io"
	"time"
)

// WatchConfig controls polling behaviour for watch mode.
type WatchConfig struct {
	IntervalSeconds int
	MaxIterations   int // 0 means unlimited
	Quiet           bool
}

// WatchResult holds the outcome of a single poll iteration.
type WatchResult struct {
	Iteration int
	Timestamp time.Time
	Report    *Report
	Stats     DiffStats
	HasDrift  bool
}

// WatchFunc is the callback invoked after each poll. Return true to stop.
type WatchFunc func(result WatchResult) bool

// Watch repeatedly compares two state files at a fixed interval, calling fn
// after each comparison. It stops when fn returns true, MaxIterations is
// reached, or the context is cancelled via the done channel.
func Watch(oldPath, newPath string, cfg WatchConfig, out io.Writer, fn WatchFunc) error {
	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 30
	}

	iteration := 0
	for {
		iteration++

		report, err := pollOnce(oldPath, newPath)
		if err != nil {
			return fmt.Errorf("watch iteration %d: %w", iteration, err)
		}

		stats := ComputeStats(report)
		result := WatchResult{
			Iteration: iteration,
			Timestamp: time.Now().UTC(),
			Report:    report,
			Stats:     stats,
			HasDrift:  report.HasDrift(),
		}

		if !cfg.Quiet {
			fmt.Fprintf(out, "[%s] iteration %d — added:%d removed:%d modified:%d\n",
				result.Timestamp.Format(time.RFC3339),
				iteration,
				stats.Added,
				stats.Removed,
				stats.Modified,
			)
		}

		if stop := fn(result); stop {
			return nil
		}

		if cfg.MaxIterations > 0 && iteration >= cfg.MaxIterations {
			return nil
		}

		time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
	}
}

// pollOnce parses both state files and returns a diff report.
func pollOnce(oldPath, newPath string) (*Report, error) {
	oldState, err := parseStateFilePath(oldPath)
	if err != nil {
		return nil, fmt.Errorf("parsing old state: %w", err)
	}
	newState, err := parseStateFilePath(newPath)
	if err != nil {
		return nil, fmt.Errorf("parsing new state: %w", err)
	}
	return CompareToReport(oldState, newState), nil
}
