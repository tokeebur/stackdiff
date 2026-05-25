package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// AuditEntry records a single operation applied to a report during processing.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Detail    string    `json:"detail,omitempty"`
	InputLen  int       `json:"input_len"`
	OutputLen int       `json:"output_len"`
}

// AuditLog holds an ordered list of audit entries for a diff run.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends an audit entry describing an operation and its effect on report size.
func (a *AuditLog) Record(operation, detail string, inputLen, outputLen int) {
	if a == nil {
		return
	}
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		Detail:    detail,
		InputLen:  inputLen,
		OutputLen: outputLen,
	})
}

// Len returns the number of recorded entries.
func (a *AuditLog) Len() int {
	if a == nil {
		return 0
	}
	return len(a.Entries)
}

// WriteAuditLog writes a human-readable audit trail to w.
func WriteAuditLog(a *AuditLog, w io.Writer) error {
	if a == nil {
		_, err := fmt.Fprintln(w, "no audit log available")
		return err
	}
	if len(a.Entries) == 0 {
		_, err := fmt.Fprintln(w, "audit log: no operations recorded")
		return err
	}
	_, err := fmt.Fprintf(w, "audit log (%d operations):\n", len(a.Entries))
	if err != nil {
		return err
	}
	for i, e := range a.Entries {
		detail := ""
		if e.Detail != "" {
			detail = fmt.Sprintf(" [%s]", e.Detail)
		}
		_, err = fmt.Fprintf(w, "  %d. %s%s: %d -> %d entries (%s)\n",
			i+1, e.Operation, detail, e.InputLen, e.OutputLen,
			e.Timestamp.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}
	return nil
}

// AuditOperations returns a deduplicated sorted list of operation names in the log.
func AuditOperations(a *AuditLog) []string {
	if a == nil {
		return nil
	}
	seen := make(map[string]struct{})
	for _, e := range a.Entries {
		seen[e.Operation] = struct{}{}
	}
	ops := make([]string, 0, len(seen))
	for op := range seen {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	return ops
}
