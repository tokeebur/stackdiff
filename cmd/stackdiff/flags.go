package main

import (
	"flag"
	"strings"

	"github.com/yourorg/stackdiff/internal/diff"
)

// CLIFlags holds all parsed command-line options.
type CLIFlags struct {
	OutputFormat  string
	FilterType    string
	FilterPrefix  string
	SortOrder     string
	IgnoreTypes   []string
	IgnoreAttrs   []string
	IgnorePrefixes []string
}

// ParseFlags parses os.Args using the provided FlagSet and returns CLIFlags.
func ParseFlags(fs *flag.FlagSet, args []string) (CLIFlags, error) {
	var f CLIFlags
	var ignoreTypes, ignoreAttrs, ignorePrefixes string

	fs.StringVar(&f.OutputFormat, "format", "text", "Output format: text or markdown")
	fs.StringVar(&f.FilterType, "filter-type", "", "Only show resources of this type")
	fs.StringVar(&f.FilterPrefix, "filter-prefix", "", "Only show resources whose address starts with this prefix")
	fs.StringVar(&f.SortOrder, "sort", "address", "Sort order: address, type, or action")
	fs.StringVar(&ignoreTypes, "ignore-types", "", "Comma-separated resource types to ignore")
	fs.StringVar(&ignoreAttrs, "ignore-attrs", "", "Comma-separated attribute keys to ignore")
	fs.StringVar(&ignorePrefixes, "ignore-prefixes", "", "Comma-separated address prefixes to ignore")

	if err := fs.Parse(args); err != nil {
		return CLIFlags{}, err
	}

	f.IgnoreTypes = splitCSV(ignoreTypes)
	f.IgnoreAttrs = splitCSV(ignoreAttrs)
	f.IgnorePrefixes = splitCSV(ignorePrefixes)
	return f, nil
}

// BuildIgnoreRules converts CLIFlags into a slice of IgnoreRule.
func BuildIgnoreRules(f CLIFlags) []diff.IgnoreRule {
	var rules []diff.IgnoreRule
	for _, t := range f.IgnoreTypes {
		rules = append(rules, diff.IgnoreRule{ResourceType: t})
	}
	for _, a := range f.IgnoreAttrs {
		rules = append(rules, diff.IgnoreRule{AttributeKey: a})
	}
	for _, p := range f.IgnorePrefixes {
		rules = append(rules, diff.IgnoreRule{AddressPrefix: p})
	}
	return rules
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
