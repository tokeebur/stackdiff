package diff

import (
	"sort"
)

// SortOrder defines how report changes are sorted.
type SortOrder string

const (
	SortByAddress  SortOrder = "address"
	SortByType     SortOrder = "type"
	SortByAction   SortOrder = "action"
)

// SortReport returns a new Report with ResourceChanges sorted by the given order.
// If order is unrecognised, the original order is preserved.
func SortReport(r Report, order SortOrder) Report {
	changes := make([]ResourceChange, len(r.ResourceChanges))
	copy(changes, r.ResourceChanges)

	switch order {
	case SortByAddress:
		sort.Slice(changes, func(i, j int) bool {
			return changes[i].Address < changes[j].Address
		})
	case SortByType:
		sort.Slice(changes, func(i, j int) bool {
			if changes[i].ResourceType == changes[j].ResourceType {
				return changes[i].Address < changes[j].Address
			}
			return changes[i].ResourceType < changes[j].ResourceType
		})
	case SortByAction:
		actionOrder := map[ChangeAction]int{
			ActionAdded:    0,
			ActionRemoved:  1,
			ActionModified: 2,
		}
		sort.Slice(changes, func(i, j int) bool {
			oi := actionOrder[changes[i].Action]
			oj := actionOrder[changes[j].Action]
			if oi == oj {
				return changes[i].Address < changes[j].Address
			}
			return oi < oj
		})
	}

	return Report{
		ResourceChanges: changes,
	}
}
