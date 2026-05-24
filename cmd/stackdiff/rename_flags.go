package main

import (
	"flag"
	"strings"

	"github.com/your-org/stackdiff/internal/diff"
)

// RenameFlags holds raw CLI values for rename rules.
type RenameFlags struct {
	// RenameCSV is a comma-separated list of "old=new" address pairs.
	RenameCSV string
}

// ParseRenameFlags registers and parses rename-related flags from the given FlagSet.
func ParseRenameFlags(fs *flag.FlagSet) *RenameFlags {
	f := &RenameFlags{}
	fs.StringVar(&f.RenameCSV, "rename", "", "comma-separated list of address renames in old=new format")
	return f
}

// BuildRenameConfig converts parsed flags into a diff.RenameConfig.
// Each entry in RenameCSV must be of the form "old_address=new_address".
// Malformed entries (missing '=') are silently skipped.
func BuildRenameConfig(f *RenameFlags) diff.RenameConfig {
	if f == nil || strings.TrimSpace(f.RenameCSV) == "" {
		return diff.RenameConfig{}
	}

	var rules []diff.RenameRule
	for _, pair := range splitCSV(f.RenameCSV) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		if from == "" || to == "" {
			continue
		}
		rules = append(rules, diff.RenameRule{From: from, To: to})
	}
	return diff.RenameConfig{Rules: rules}
}
