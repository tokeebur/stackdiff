package diff

import "sort"

// DriftStats holds aggregate counts of resource changes.
type DriftStats struct {
	Added    int
	Removed  int
	Modified int
}

// TypeStats holds per-resource-type drift counts.
type TypeStats struct {
	ResourceType string
	Added        int
	Removed      int
	Modified     int
}

// ComputeStats calculates overall drift counts from a Report.
func ComputeStats(r *Report) DriftStats {
	if r == nil {
		return DriftStats{}
	}
	var s DriftStats
	for _, rc := range r.Changes {
		switch rc.Action {
		case ActionAdded:
			s.Added++
		case ActionRemoved:
			s.Removed++
		case ActionModified:
			s.Modified++
		}
	}
	return s
}

// StatsByType groups drift counts by resource type.
func StatsByType(r *Report) []TypeStats {
	if r == nil {
		return nil
	}
	m := map[string]*TypeStats{}
	for _, rc := range r.Changes {
		ts, ok := m[rc.ResourceType]
		if !ok {
			ts = &TypeStats{ResourceType: rc.ResourceType}
			m[rc.ResourceType] = ts
		}
		switch rc.Action {
		case ActionAdded:
			ts.Added++
		case ActionRemoved:
			ts.Removed++
		case ActionModified:
			ts.Modified++
		}
	}
	result := make([]TypeStats, 0, len(m))
	for _, ts := range m {
		result = append(result, *ts)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ResourceType < result[j].ResourceType
	})
	return result
}
