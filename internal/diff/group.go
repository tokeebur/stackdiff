package diff

import "sort"

// GroupedReport holds report entries organised by resource type.
type GroupedReport struct {
	Groups map[string][]ReportEntry
}

// GroupByType partitions a Report's entries into a GroupedReport keyed by
// resource type. Returns an error if the report is nil.
func GroupByType(r *Report) (*GroupedReport, error) {
	if r == nil {
		return nil, ErrNilReport
	}

	groups := make(map[string][]ReportEntry)
	for _, entry := range r.Entries {
		groups[entry.ResourceType] = append(groups[entry.ResourceType], entry)
	}

	// Sort entries within each group by address for deterministic output.
	for k := range groups {
		sort.Slice(groups[k], func(i, j int) bool {
			return groups[k][i].Address < groups[k][j].Address
		})
	}

	return &GroupedReport{Groups: groups}, nil
}

// SortedTypes returns the resource type keys in alphabetical order.
func (g *GroupedReport) SortedTypes() []string {
	keys := make([]string, 0, len(g.Groups))
	for k := range g.Groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TotalEntries returns the total number of entries across all groups.
func (g *GroupedReport) TotalEntries() int {
	total := 0
	for _, entries := range g.Groups {
		total += len(entries)
	}
	return total
}
