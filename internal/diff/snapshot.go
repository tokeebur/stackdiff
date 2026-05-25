package diff

import (
	"fmt"
	"time"
)

// Snapshot captures a point-in-time view of a Report with metadata.
type Snapshot struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label,omitempty"`
	Report    *Report   `json:"report"`
	Stats     *Stats    `json:"stats"`
}

// SnapshotStore holds an ordered list of snapshots.
type SnapshotStore struct {
	Snapshots []Snapshot `json:"snapshots"`
}

// NewSnapshot creates a Snapshot from a Report, computing stats automatically.
func NewSnapshot(label string, r *Report) (*Snapshot, error) {
	if r == nil {
		return nil, fmt.Errorf("snapshot: report must not be nil")
	}
	stats := ComputeStats(r)
	return &Snapshot{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Report:    r,
		Stats:     &stats,
	}, nil
}

// AddSnapshot appends a snapshot to the store.
func (s *SnapshotStore) AddSnapshot(snap Snapshot) {
	s.Snapshots = append(s.Snapshots, snap)
}

// Latest returns the most recently added snapshot, or nil if empty.
func (s *SnapshotStore) Latest() *Snapshot {
	if len(s.Snapshots) == 0 {
		return nil
	}
	return &s.Snapshots[len(s.Snapshots)-1]
}

// FindByLabel returns the first snapshot matching the given label.
func (s *SnapshotStore) FindByLabel(label string) *Snapshot {
	for i := range s.Snapshots {
		if s.Snapshots[i].Label == label {
			return &s.Snapshots[i]
		}
	}
	return nil
}
