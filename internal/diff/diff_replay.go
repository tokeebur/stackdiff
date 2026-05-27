package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// ReplayEntry represents a single recorded diff operation in a replay log.
type ReplayEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Operation string            `json:"operation"`
	Address   string            `json:"address"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// ReplayLog holds an ordered sequence of replay entries.
type ReplayLog struct {
	Entries []ReplayEntry `json:"entries"`
}

// NewReplayLog creates an empty ReplayLog.
func NewReplayLog() *ReplayLog {
	return &ReplayLog{}
}

// Record appends a new entry to the replay log.
func (r *ReplayLog) Record(op, address string, meta map[string]string) {
	if r == nil {
		return
	}
	r.Entries = append(r.Entries, ReplayEntry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Address:   address,
		Meta:      meta,
	})
}

// Len returns the number of entries in the replay log.
func (r *ReplayLog) Len() int {
	if r == nil {
		return 0
	}
	return len(r.Entries)
}

// BuildReplayLog constructs a ReplayLog from a Report, recording one entry
// per resource change with the action as the operation name.
func BuildReplayLog(report *Report) *ReplayLog {
	log := NewReplayLog()
	if report == nil {
		return log
	}
	for _, entry := range report.Changes {
		log.Record(string(entry.Action), entry.Address, map[string]string{
			"resource_type": entry.ResourceType,
		})
	}
	return log
}

// WriteReplayLog writes a human-readable replay log to w.
func WriteReplayLog(w io.Writer, log *ReplayLog) error {
	if log == nil {
		_, err := fmt.Fprintln(w, "no replay log available")
		return err
	}
	if len(log.Entries) == 0 {
		_, err := fmt.Fprintln(w, "replay log is empty")
		return err
	}
	sorted := make([]ReplayEntry, len(log.Entries))
	copy(sorted, log.Entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})
	for _, e := range sorted {
		_, err := fmt.Fprintf(w, "[%s] %-10s %s\n",
			e.Timestamp.Format(time.RFC3339), e.Operation, e.Address)
		if err != nil {
			return err
		}
	}
	return nil
}
