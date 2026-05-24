package diff

import (
	"sort"
)

// SortOrder defines how report changes are sorted.
type SortOrder string

const (
	SortByAddress SortOrder = "address"
	SortByType    SortOrder = "type"
	SortByAction  SortOrder = "action"
)

// ValidSortOrders returns all recognised SortOrder values.
func ValidSortOrders() []SortOrder {
	return []SortOrder{SortByAddress, SortByType, SortByAction}
}

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
		sort.Slice(changes, func(i, j int) bool {
			oi := actionSortWeight(changes[i].Action)
			oj := actionSortWeight(changes[j].Action)
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

// actionSortWeight returns the sort priority for a given ChangeAction.
// Lower values sort first: added before removed before modified.
func actionSortWeight(a ChangeAction) int {
	switch a {
	case ActionAdded:
		return 0
	case ActionRemoved:
		return 1
	case ActionModified:
		return 2
	default:
		return 3
	}
}
