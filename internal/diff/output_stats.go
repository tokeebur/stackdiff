package diff

import (
	"fmt"
	"io"
	"sort"
)

// WriteStats writes a drift statistics summary to w in the given format.
// Supported formats: "text", "markdown".
func WriteStats(w io.Writer, r Report, format string) error {
	s := ComputeStats(r)
	byType := StatsByType(r)

	switch format {
	case "markdown":
		return writeStatsMarkdown(w, s, byType)
	default:
		return writeStatsText(w, s, byType)
	}
}

func writeStatsText(w io.Writer, s Stats, byType map[string]Stats) error {
	fmt.Fprintf(w, "Drift Statistics\n")
	fmt.Fprintf(w, "  Total:    %d\n", s.Total)
	fmt.Fprintf(w, "  Added:    %d\n", s.Added)
	fmt.Fprintf(w, "  Removed:  %d\n", s.Removed)
	fmt.Fprintf(w, "  Modified: %d\n", s.Modified)
	if len(byType) > 0 {
		fmt.Fprintf(w, "\nBy Resource Type:\n")
		for _, rt := range sortedTypeKeys(byType) {
			ts := byType[rt]
			fmt.Fprintf(w, "  %-40s add=%-3d remove=%-3d modify=%d\n",
				rt, ts.Added, ts.Removed, ts.Modified)
		}
	}
	return nil
}

func writeStatsMarkdown(w io.Writer, s Stats, byType map[string]Stats) error {
	fmt.Fprintf(w, "## Drift Statistics\n\n")
	fmt.Fprintf(w, "| Metric | Count |\n|--------|-------|\n")
	fmt.Fprintf(w, "| Total | %d |\n", s.Total)
	fmt.Fprintf(w, "| Added | %d |\n", s.Added)
	fmt.Fprintf(w, "| Removed | %d |\n", s.Removed)
	fmt.Fprintf(w, "| Modified | %d |\n", s.Modified)
	if len(byType) > 0 {
		fmt.Fprintf(w, "\n### By Resource Type\n\n")
		fmt.Fprintf(w, "| Type | Added | Removed | Modified |\n|------|-------|---------|----------|\n")
		for _, rt := range sortedTypeKeys(byType) {
			ts := byType[rt]
			fmt.Fprintf(w, "| %s | %d | %d | %d |\n", rt, ts.Added, ts.Removed, ts.Modified)
		}
	}
	return nil
}

func sortedTypeKeys(m map[string]Stats) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
