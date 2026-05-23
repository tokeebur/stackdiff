package main

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempStateFile(t *testing.T, data map[string]interface{}) string {
	t.Helper()
	f, err := os.CreateTemp("", "state-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatalf("encode state: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRun_MissingArgs(t *testing.T) {
	if err := run([]string{}); err == nil {
		t.Fatal("expected error for missing args")
	}
	if err := run([]string{"only-one"}); err == nil {
		t.Fatal("expected error for single arg")
	}
}

func TestRun_InvalidFile(t *testing.T) {
	if err := run([]string{"nonexistent-a.json", "nonexistent-b.json"}); err == nil {
		t.Fatal("expected error for missing files")
	}
}

func TestRun_IdenticalStates(t *testing.T) {
	state := map[string]interface{}{
		"version": 4,
		"resources": []interface{}{},
	}
	a := writeTempStateFile(t, state)
	b := writeTempStateFile(t, state)
	if err := run([]string{a, b}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_DriftDetected(t *testing.T) {
	stateA := map[string]interface{}{
		"version": 4,
		"resources": []interface{}{
			map[string]interface{}{
				"type": "aws_instance", "name": "web", "provider": "aws",
				"instances": []interface{}{
					map[string]interface{}{"attributes": map[string]interface{}{"id": "i-123"}},
				},
			},
		},
	}
	stateB := map[string]interface{}{
		"version":   4,
		"resources": []interface{}{},
	}
	a := writeTempStateFile(t, stateA)
	b := writeTempStateFile(t, stateB)
	if err := run([]string{a, b}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
