package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatSummary writes a human-readable drift summary to w.
func FormatSummary(w io.Writer, result *Result) {
	if !result.HasChanges() {
		fmt.Fprintln(w, "No drift detected. States are identical.")
		return
	}

	// Sort for deterministic output.
	sorted := make([]ResourceDiff, len(result.Changes))
	copy(sorted, result.Changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Address < sorted[j].Address
	})

	fmt.Fprintf(w, "Drift detected: %d change(s)\n\n", len(sorted))

	for _, c := range sorted {
		switch c.ChangeType {
		case Added:
			fmt.Fprintf(w, "  [+] %s (added)\n", c.Address)
			printAttrs(w, c.NewAttrs, "      ")
		case Removed:
			fmt.Fprintf(w, "  [-] %s (removed)\n", c.Address)
			printAttrs(w, c.OldAttrs, "      ")
		case Modified:
			fmt.Fprintf(w, "  [~] %s (modified)\n", c.Address)
			printModifiedAttrs(w, c.OldAttrs, c.NewAttrs, "      ")
		}
		fmt.Fprintln(w)
	}
}

func printAttrs(w io.Writer, attrs map[string]interface{}, indent string) {
	keys := sortedKeys(attrs)
	for _, k := range keys {
		fmt.Fprintf(w, "%s%s = %v\n", indent, k, attrs[k])
	}
}

func printModifiedAttrs(w io.Writer, old, new map[string]interface{}, indent string) {
	keys := sortedKeys(mergeMaps(old, new))
	for _, k := range keys {
		ov, okO := old[k]
		nv, okN := new[k]
		switch {
		case okO && okN && fmt.Sprintf("%v", ov) != fmt.Sprintf("%v", nv):
			fmt.Fprintf(w, "%s%s: %v -> %v\n", indent, k, ov, nv)
		case !okO:
			fmt.Fprintf(w, "%s%s: (new) %v\n", indent, k, nv)
		case !okN:
			fmt.Fprintf(w, "%s%s: %v (removed)\n", indent, k, ov)
		}
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

// suppress unused import warning
var _ = strings.TrimSpace
