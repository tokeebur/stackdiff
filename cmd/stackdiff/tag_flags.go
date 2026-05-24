package main

import (
	"flag"
	"strings"

	"github.com/your-org/stackdiff/internal/diff"
)

// TagConfig holds parsed tagging configuration from CLI flags.
type TagConfig struct {
	Rules []diff.TagRule
}

// ParseTagFlags reads tag-related flags from the provided FlagSet.
// Tags are specified as: --tag "aws_instance::compute" or "module.network::networking"
// Format: [resourceType|addressPrefix]::<tag>
func ParseTagFlags(fs *flag.FlagSet, args []string) (*TagConfig, error) {
	var rawTags string
	fs.StringVar(&rawTags, "tag", "", "Comma-separated tag rules: type=<type>:<tag> or prefix=<prefix>:<tag>")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &TagConfig{}
	for _, entry := range splitCSV(rawTags) {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		rule, ok := parseTagRule(entry)
		if ok {
			cfg.Rules = append(cfg.Rules, rule)
		}
	}
	return cfg, nil
}

// parseTagRule parses a single tag rule string.
// Expected formats:
//   - type=<resourceType>:<tag>
//   - prefix=<addressPrefix>:<tag>
func parseTagRule(s string) (diff.TagRule, bool) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return diff.TagRule{}, false
	}
	kv := strings.SplitN(parts[0], "=", 2)
	tag := strings.TrimSpace(parts[1])
	if len(kv) != 2 || tag == "" {
		return diff.TagRule{}, false
	}
	switch strings.TrimSpace(kv[0]) {
	case "type":
		return diff.TagRule{ResourceType: strings.TrimSpace(kv[1]), Tag: tag}, true
	case "prefix":
		return diff.TagRule{AddressPrefix: strings.TrimSpace(kv[1]), Tag: tag}, true
	}
	return diff.TagRule{}, false
}
