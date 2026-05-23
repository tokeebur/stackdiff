package main

import (
	"flag"
	"strings"

	"github.com/your-org/stackdiff/internal/diff"
)

// Flags holds all parsed CLI flag values.
type Flags struct {
	Format        string
	SortBy        string
	FilterType    string
	FilterAddress string
	IgnoreTypes   string
	IgnoreAttrs   string
	IgnorePrefix  string
	ShowStats     bool
	Validate      bool
}

// ParseFlags registers and parses all stackdiff CLI flags from the provided
// FlagSet, returning the populated Flags struct.
func ParseFlags(fs *flag.FlagSet, args []string) (Flags, error) {
	var f Flags
	fs.StringVar(&f.Format, "format", "text", "Output format: text or markdown")
	fs.StringVar(&f.SortBy, "sort", "address", "Sort results by: address, type, or action")
	fs.StringVar(&f.FilterType, "filter-type", "", "Only show resources matching this type")
	fs.StringVar(&f.FilterAddress, "filter-address", "", "Only show resources whose address starts with this prefix")
	fs.StringVar(&f.IgnoreTypes, "ignore-types", "", "Comma-separated resource types to ignore")
	fs.StringVar(&f.IgnoreAttrs, "ignore-attrs", "", "Comma-separated attribute keys to ignore")
	fs.StringVar(&f.IgnorePrefix, "ignore-prefix", "", "Comma-separated address prefixes to ignore")
	fs.BoolVar(&f.ShowStats, "stats", false, "Print a statistics summary after the report")
	fs.BoolVar(&f.Validate, "validate", false, "Validate the report before output and exit with error if invalid")
	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	return f, nil
}

// BuildIgnoreRules converts parsed flag strings into a slice of IgnoreRule.
func BuildIgnoreRules(f Flags) []diff.IgnoreRule {
	var rules []diff.IgnoreRule
	for _, t := range splitCSV(f.IgnoreTypes) {
		rules = append(rules, diff.IgnoreRule{ResourceType: t})
	}
	for _, p := range splitCSV(f.IgnorePrefix) {
		rules = append(rules, diff.IgnoreRule{AddressPrefix: p})
	}
	for _, a := range splitCSV(f.IgnoreAttrs) {
		rules = append(rules, diff.IgnoreRule{AttributeKey: a})
	}
	return rules
}

// splitCSV splits a comma-separated string into trimmed, non-empty tokens.
func splitCSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
