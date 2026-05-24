package diff

import (
	"sort"
	"time"
)

// TimelineEntry records a snapshot of drift at a point in time.
type TimelineEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Modified  int       `json:"modified"`
	Total     int       `json:"total"`
	Label     string    `json:"label,omitempty"`
}

// Timeline is an ordered sequence of drift snapshots.
type Timeline struct {
	Entries []TimelineEntry `json:"entries"`
}

// AddEntry appends a new entry derived from the given stats.
func (t *Timeline) AddEntry(stats DriftStats, label string) {
	entry := TimelineEntry{
		Timestamp: time.Now().UTC(),
		Added:     stats.Added,
		Removed:   stats.Removed,
		Modified:  stats.Modified,
		Total:     stats.Total,
		Label:     label,
	}
	t.Entries = append(t.Entries, entry)
}

// Len returns the number of entries in the timeline.
func (t *Timeline) Len() int {
	if t == nil {
		return 0
	}
	return len(t.Entries)
}

// SortByTime orders entries from oldest to newest.
func (t *Timeline) SortByTime() {
	if t == nil {
		return
	}
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Timestamp.Before(t.Entries[j].Timestamp)
	})
}

// Latest returns the most recent entry, or nil if empty.
func (t *Timeline) Latest() *TimelineEntry {
	if t == nil || len(t.Entries) == 0 {
		return nil
	}
	return &t.Entries[len(t.Entries)-1]
}
