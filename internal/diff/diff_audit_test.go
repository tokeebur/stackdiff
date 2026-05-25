package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewAuditLog_Empty(t *testing.T) {
	a := NewAuditLog()
	if a == nil {
		t.Fatal("expected non-nil AuditLog")
	}
	if a.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", a.Len())
	}
}

func TestAuditLog_Record(t *testing.T) {
	a := NewAuditLog()
	a.Record("filter", "by_type", 10, 7)
	if a.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", a.Len())
	}
	e := a.Entries[0]
	if e.Operation != "filter" {
		t.Errorf("expected operation 'filter', got %q", e.Operation)
	}
	if e.Detail != "by_type" {
		t.Errorf("expected detail 'by_type', got %q", e.Detail)
	}
	if e.InputLen != 10 || e.OutputLen != 7 {
		t.Errorf("unexpected lens: %d -> %d", e.InputLen, e.OutputLen)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAuditLog_Record_NilSafe(t *testing.T) {
	var a *AuditLog
	a.Record("op", "", 1, 1) // should not panic
	if a.Len() != 0 {
		t.Error("nil log should report 0 len")
	}
}

func TestAuditLog_MultipleEntries(t *testing.T) {
	a := NewAuditLog()
	a.Record("sort", "", 5, 5)
	a.Record("truncate", "limit=3", 5, 3)
	a.Record("redact", "", 3, 3)
	if a.Len() != 3 {
		t.Errorf("expected 3 entries, got %d", a.Len())
	}
}

func TestWriteAuditLog_Nil(t *testing.T) {
	var buf bytes.Buffer
	var a *AuditLog
	if err := WriteAuditLog(a, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no audit log") {
		t.Errorf("expected 'no audit log' message, got: %q", buf.String())
	}
}

func TestWriteAuditLog_Empty(t *testing.T) {
	var buf bytes.Buffer
	a := NewAuditLog()
	if err := WriteAuditLog(a, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no operations") {
		t.Errorf("expected 'no operations' message, got: %q", buf.String())
	}
}

func TestWriteAuditLog_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	a := NewAuditLog()
	a.Record("filter", "addr", 8, 4)
	a.Record("sort", "", 4, 4)
	if err := WriteAuditLog(a, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "filter") {
		t.Errorf("expected 'filter' in output, got: %q", out)
	}
	if !strings.Contains(out, "sort") {
		t.Errorf("expected 'sort' in output, got: %q", out)
	}
	if !strings.Contains(out, "8 -> 4") {
		t.Errorf("expected '8 -> 4' in output, got: %q", out)
	}
}

func TestAuditOperations_NilLog(t *testing.T) {
	var a *AuditLog
	if ops := AuditOperations(a); ops != nil {
		t.Errorf("expected nil, got %v", ops)
	}
}

func TestAuditOperations_Deduplicated(t *testing.T) {
	a := NewAuditLog()
	a.Record("filter", "a", 5, 4)
	a.Record("sort", "", 4, 4)
	a.Record("filter", "b", 4, 3)
	ops := AuditOperations(a)
	if len(ops) != 2 {
		t.Errorf("expected 2 unique ops, got %d: %v", len(ops), ops)
	}
	if ops[0] != "filter" || ops[1] != "sort" {
		t.Errorf("expected sorted [filter sort], got %v", ops)
	}
}
