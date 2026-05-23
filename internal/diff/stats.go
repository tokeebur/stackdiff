package diff

// Stats holds aggregate counts derived from a Report.
type Stats struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// ComputeStats calculates drift statistics from a Report.
func ComputeStats(r Report) Stats {
	var s Stats
	for _, rc := range r.Changes {
		switch rc.Action {
		case ActionAdd:
			s.Added++
		case ActionRemove:
			s.Removed++
		case ActionModify:
			s.Modified++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified
	return s
}

// StatsByType returns a map of resource type to Stats.
func StatsByType(r Report) map[string]Stats {
	result := make(map[string]Stats)
	for _, rc := range r.Changes {
		s := result[rc.ResourceType]
		switch rc.Action {
		case ActionAdd:
			s.Added++
		case ActionRemove:
			s.Removed++
		case ActionModify:
			s.Modified++
		}
		s.Total = s.Added + s.Removed + s.Modified
		result[rc.ResourceType] = s
	}
	return result
}
