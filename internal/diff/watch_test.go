package diff

import (
	"bytes"
	"testing"
)

func TestWatch_MaxIterations(t *testing.T) {
	oldPath := writeTempWatchState(t, map[string]map[string]string{
		"aws_instance.web": {"ami": "ami-111"},
	})
	newPath := writeTempWatchState(t, map[string]map[string]string{
		"aws_instance.web": {"ami": "ami-222"},
	})

	var buf bytes.Buffer
	cfg := WatchConfig{IntervalSeconds: 0, MaxIterations: 2, Quiet: false}

	callCount := 0
	err := Watch(oldPath, newPath, cfg, &buf, func(r WatchResult) bool {
		callCount++
		if r.Report == nil {
			t.Error("expected non-nil report")
		}
		if !r.HasDrift {
			t.Error("expected drift to be detected")
		}
		return false
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 iterations, got %d", callCount)
	}
}

func TestWatch_StopEarly(t *testing.T) {
	oldPath := writeTempWatchState(t, map[string]map[string]string{
		"aws_s3_bucket.data": {"region": "us-east-1"},
	})
	newPath := writeTempWatchState(t, map[string]map[string]string{
		"aws_s3_bucket.data": {"region": "eu-west-1"},
	})

	var buf bytes.Buffer
	cfg := WatchConfig{IntervalSeconds: 0, MaxIterations: 10, Quiet: true}

	callCount := 0
	err := Watch(oldPath, newPath, cfg, &buf, func(r WatchResult) bool {
		callCount++
		return callCount >= 3 // stop after 3
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 iterations before stop, got %d", callCount)
	}
}

func TestWatch_NoDrift(t *testing.T) {
	attrs := map[string]map[string]string{
		"aws_instance.web": {"ami": "ami-abc"},
	}
	oldPath := writeTempWatchState(t, attrs)
	newPath := writeTempWatchState(t, attrs)

	var buf bytes.Buffer
	cfg := WatchConfig{IntervalSeconds: 0, MaxIterations: 1, Quiet: true}

	err := Watch(oldPath, newPath, cfg, &buf, func(r WatchResult) bool {
		if r.HasDrift {
			t.Error("expected no drift for identical states")
		}
		return false
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatch_InvalidFile(t *testing.T) {
	cfg := WatchConfig{IntervalSeconds: 0, MaxIterations: 1, Quiet: true}
	var buf bytes.Buffer
	err := Watch("/nonexistent/old.tfstate", "/nonexistent/new.tfstate", cfg, &buf, func(r WatchResult) bool {
		return false
	})
	if err == nil {
		t.Error("expected error for missing files")
	}
}
