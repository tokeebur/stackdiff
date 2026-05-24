package diff

import (
	"fmt"
	"strings"
)

// DiffSummaryOptions controls how the diff summary is rendered.
type DiffSummaryOptions struct {
	ShowScore    bool
	ShowStats    bool
	ShowTags     bool
	Compact      bool
}

// DefaultDiffSummaryOptions returns sensible defaults.
func DefaultDiffSummaryOptions() DiffSummaryOptions {
	return DiffSummaryOptions{
		ShowScore: true,
		ShowStats: true,
		ShowTags:  false,
		Compact:   false,
	}
}

// BuildDiffSummary produces a human-readable multi-line summary string
// from a Report, stats, and optional score.
func BuildDiffSummary(r *Report, stats Stats, score *ScoreResult, opts DiffSummaryOptions) string {
	if r == nil {
		return "no report available"
	}

	var sb strings.Builder

	if !r.HasDrift() {
		sb.WriteString("✓ No drift detected between state files.\n")
		return sb.String()
	}

	sb.WriteString("✗ Drift detected:\n")

	if opts.ShowStats {
		sb.WriteString(fmt.Sprintf("  Added:    %d\n", stats.Added))
		sb.WriteString(fmt.Sprintf("  Removed:  %d\n", stats.Removed))
		sb.WriteString(fmt.Sprintf("  Modified: %d\n", stats.Modified))
		sb.WriteString(fmt.Sprintf("  Total:    %d\n", stats.Total))
	}

	if opts.ShowScore && score != nil {
		sb.WriteString(fmt.Sprintf("  Score:    %.1f (%s)\n", score.Score, score.Label))
	}

	if !opts.Compact {
		for _, entry := range r.Changes {
			line := fmt.Sprintf("  [%s] %s", entry.Action, entry.Address)
			if opts.ShowTags && len(entry.Tags) > 0 {
				line += fmt.Sprintf(" tags=%s", strings.Join(entry.Tags, ","))
			}
			sb.WriteString(line + "\n")
		}
	}

	return sb.String()
}
