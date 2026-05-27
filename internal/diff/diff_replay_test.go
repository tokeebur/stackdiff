package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeReplayReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: ActionAdded},
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: ActionRemoved},
			{Address: "aws_sg.default", ResourceType: "aws_sg", Action: ActionModified},
		},
	}
}

func TestNewReplayLog_Empty(t *testing.T) {
	log := NewReplayLog()
	if log == nil {
		t.Fatal("expected non-nil ReplayLog")
	}
	if log.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", log.Len())
	}
}

func TestReplayLog_Record(t *testing.T) {
	log := NewReplayLog()
	log.Record("added", "aws_instance.web", nil)
	if log.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", log.Len())
	}
	if log.Entries[0].Address != "aws_instance.web" {
		t.Errorf("unexpected address: %s", log.Entries[0].Address)
	}
	if log.Entries[0].Operation != "added" {
		t.Errorf("unexpected operation: %s", log.Entries[0].Operation)
	}
}

func TestReplayLog_Record_NilSafe(t *testing.T) {
	var log *ReplayLog
	log.Record("added", "addr", nil) // should not panic
}

func TestBuildReplayLog_NilReport(t *testing.T) {
	log := BuildReplayLog(nil)
	if log == nil {
		t.Fatal("expected non-nil log")
	}
	if log.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", log.Len())
	}
}

func TestBuildReplayLog_EntryCount(t *testing.T) {
	report := makeReplayReport()
	log := BuildReplayLog(report)
	if log.Len() != len(report.Changes) {
		t.Errorf("expected %d entries, got %d", len(report.Changes), log.Len())
	}
}

func TestBuildReplayLog_Operations(t *testing.T) {
	report := makeReplayReport()
	log := BuildReplayLog(report)
	ops := map[string]bool{}
	for _, e := range log.Entries {
		ops[e.Operation] = true
	}
	for _, want := range []string{"added", "removed", "modified"} {
		if !ops[want] {
			t.Errorf("missing operation %q in replay log", want)
		}
	}
}

func TestWriteReplayLog_NilLog(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteReplayLog(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no replay log") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWriteReplayLog_WithEntries(t *testing.T) {
	report := makeReplayReport()
	log := BuildReplayLog(report)
	var buf bytes.Buffer
	if err := WriteReplayLog(&buf, log); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, addr := range []string{"aws_instance.web", "aws_s3_bucket.data", "aws_sg.default"} {
		if !strings.Contains(out, addr) {
			t.Errorf("expected address %q in output", addr)
		}
	}
}
