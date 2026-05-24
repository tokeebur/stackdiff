package diff

import "fmt"

// ThresholdConfig defines limits that trigger warnings or errors
// when the number of drifted resources exceeds expectations.
type ThresholdConfig struct {
	MaxAdded    int
	MaxRemoved  int
	MaxModified int
	MaxTotal    int
}

// ThresholdResult holds the outcome of a threshold evaluation.
type ThresholdResult struct {
	Exceeded   bool
	Violations []string
}

// EvaluateThresholds checks a DriftStats against the provided ThresholdConfig
// and returns a ThresholdResult describing any violations.
func EvaluateThresholds(stats DriftStats, cfg ThresholdConfig) ThresholdResult {
	var violations []string

	if cfg.MaxAdded > 0 && stats.Added > cfg.MaxAdded {
		violations = append(violations,
			fmt.Sprintf("added resources %d exceeds max %d", stats.Added, cfg.MaxAdded))
	}
	if cfg.MaxRemoved > 0 && stats.Removed > cfg.MaxRemoved {
		violations = append(violations,
			fmt.Sprintf("removed resources %d exceeds max %d", stats.Removed, cfg.MaxRemoved))
	}
	if cfg.MaxModified > 0 && stats.Modified > cfg.MaxModified {
		violations = append(violations,
			fmt.Sprintf("modified resources %d exceeds max %d", stats.Modified, cfg.MaxModified))
	}

	total := stats.Added + stats.Removed + stats.Modified
	if cfg.MaxTotal > 0 && total > cfg.MaxTotal {
		violations = append(violations,
			fmt.Sprintf("total drifted resources %d exceeds max %d", total, cfg.MaxTotal))
	}

	return ThresholdResult{
		Exceeded:   len(violations) > 0,
		Violations: violations,
	}
}
