package diff

import (
	"fmt"

	"github.com/stackdiff/internal/state"
)

// ChangeType represents the type of drift detected.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// ResourceDiff represents a single resource-level change.
type ResourceDiff struct {
	Address    string
	ChangeType ChangeType
	OldAttrs   map[string]interface{}
	NewAttrs   map[string]interface{}
}

// Result holds the full diff between two state files.
type Result struct {
	Changes []ResourceDiff
}

// HasChanges returns true if any drift was detected.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a count of added, removed, and modified resources.
func (r *Result) Summary() (added, removed, modified int) {
	for _, c := range r.Changes {
		switch c.ChangeType {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return
}

// Compare computes the diff between a base and target Terraform state.
func Compare(base, target *state.State) (*Result, error) {
	if base == nil {
		return nil, fmt.Errorf("base state must not be nil")
	}
	if target == nil {
		return nil, fmt.Errorf("target state must not be nil")
	}

	baseMap := base.ResourceMap()
	targetMap := target.ResourceMap()

	var changes []ResourceDiff

	// Detect removed and modified resources.
	for addr, baseRes := range baseMap {
		if targetRes, ok := targetMap[addr]; !ok {
			changes = append(changes, ResourceDiff{
				Address:    addr,
				ChangeType: Removed,
				OldAttrs:   baseRes.AttributeValues,
			})
		} else if attrsChanged(baseRes.AttributeValues, targetRes.AttributeValues) {
			changes = append(changes, ResourceDiff{
				Address:    addr,
				ChangeType: Modified,
				OldAttrs:   baseRes.AttributeValues,
				NewAttrs:   targetRes.AttributeValues,
			})
		}
	}

	// Detect added resources.
	for addr, targetRes := range targetMap {
		if _, ok := baseMap[addr]; !ok {
			changes = append(changes, ResourceDiff{
				Address:    addr,
				ChangeType: Added,
				NewAttrs:   targetRes.AttributeValues,
			})
		}
	}

	return &Result{Changes: changes}, nil
}

// attrsChanged performs a shallow comparison of two attribute maps.
func attrsChanged(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return true
	}
	for k, av := range a {
		bv, ok := b[k]
		if !ok || fmt.Sprintf("%v", av) != fmt.Sprintf("%v", bv) {
			return true
		}
	}
	return false
}
