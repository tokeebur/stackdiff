package diff

import (
	"fmt"
	"strings"
)

// Report holds a structured drift report between two state files.
type Report struct {
	Added    []string
	Removed  []string
	Modified []ModifiedResource
	Unchanged int
}

// ModifiedResource captures which attributes changed for a resource.
type ModifiedResource struct {
	Address string
	Changes map[string]AttributeChange
}

// AttributeChange holds the before and after values of a changed attribute.
type AttributeChange struct {
	Old string
	New string
}

// HasDrift returns true if any resources were added, removed, or modified.
func (r *Report) HasDrift() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

// Summary returns a concise one-line description of the drift.
func (r *Report) Summary() string {
	parts := []string{}
	if n := len(r.Added); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(r.Removed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(r.Modified); n > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", n))
	}
	if len(parts) == 0 {
		return "no drift detected"
	}
	return strings.Join(parts, ", ")
}
