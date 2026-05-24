package diff

import (
	"testing"
)

func makeThresholdStats(added, removed, modified int) DriftStats {
	return DriftStats{
		Added:    added,
		Removed:  removed,
		Modified: modified,
	}
}

func TestEvaluateThresholds_NoViolations(t *testing.T) {
	stats := makeThresholdStats(1, 1, 1)
	cfg := ThresholdConfig{MaxAdded: 5, MaxRemoved: 5, MaxModified: 5, MaxTotal: 15}
	result := EvaluateThresholds(stats, cfg)
	if result.Exceeded {
		t.Errorf("expected no violations, got: %v", result.Violations)
	}
}

func TestEvaluateThresholds_AddedExceeds(t *testing.T) {
	stats := makeThresholdStats(10, 0, 0)
	cfg := ThresholdConfig{MaxAdded: 5}
	result := EvaluateThresholds(stats, cfg)
	if !result.Exceeded {
		t.Fatal("expected threshold exceeded")
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
}

func TestEvaluateThresholds_RemovedExceeds(t *testing.T) {
	stats := makeThresholdStats(0, 8, 0)
	cfg := ThresholdConfig{MaxRemoved: 3}
	result := EvaluateThresholds(stats, cfg)
	if !result.Exceeded {
		t.Fatal("expected threshold exceeded")
	}
}

func TestEvaluateThresholds_TotalExceeds(t *testing.T) {
	stats := makeThresholdStats(3, 3, 3)
	cfg := ThresholdConfig{MaxTotal: 8}
	result := EvaluateThresholds(stats, cfg)
	if !result.Exceeded {
		t.Fatal("expected total threshold exceeded")
	}
}

func TestEvaluateThresholds_ZeroLimitsIgnored(t *testing.T) {
	stats := makeThresholdStats(100, 100, 100)
	cfg := ThresholdConfig{} // all zeros = no limits
	result := EvaluateThresholds(stats, cfg)
	if result.Exceeded {
		t.Errorf("expected no violations when limits are zero, got: %v", result.Violations)
	}
}

func TestEvaluateThresholds_MultipleViolations(t *testing.T) {
	stats := makeThresholdStats(10, 10, 10)
	cfg := ThresholdConfig{MaxAdded: 1, MaxRemoved: 1, MaxModified: 1, MaxTotal: 2}
	result := EvaluateThresholds(stats, cfg)
	if len(result.Violations) != 4 {
		t.Errorf("expected 4 violations, got %d: %v", len(result.Violations), result.Violations)
	}
}
