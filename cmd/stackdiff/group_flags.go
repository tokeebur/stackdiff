package main

import (
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/your-org/stackdiff/internal/diff"
)

// GroupConfig holds configuration for grouped output.
type GroupConfig struct {
	Enabled bool
}

// ParseGroupFlags reads the -group flag from the provided FlagSet.
func ParseGroupFlags(fs *flag.FlagSet) *GroupConfig {
	cfg := &GroupConfig{}
	fs.BoolVar(&cfg.Enabled, "group", false, "group drift output by resource type")
	return cfg
}

// WriteGroupedReport writes a grouped, human-readable summary to w.
// If cfg is nil or grouping is disabled the function is a no-op and returns nil.
func WriteGroupedReport(w io.Writer, r *diff.Report, cfg *GroupConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	g, err := diff.GroupByType(r)
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}

	if g.TotalEntries() == 0 {
		fmt.Fprintln(w, "No drift detected.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	for _, rtype := range g.SortedTypes() {
		entries := g.Groups[rtype]
		fmt.Fprintf(tw, "\n[%s] (%d)\n", rtype, len(entries))
		for _, e := range entries {
			fmt.Fprintf(tw, "  %-10s\t%s\n", e.Action, e.Address)
		}
	}

	return nil
}
