package main

import (
	"flag"
	"os"

	"github.com/yourorg/stackdiff/internal/diff"
)

// ContextConfig holds CLI-parsed context metadata fields.
type ContextConfig struct {
	RunID       string
	Environment string
	TriggeredBy string
	Labels      []string // raw "key=value" strings
}

// ParseContextFlags reads context-related flags from a FlagSet.
func ParseContextFlags(fs *flag.FlagSet, args []string) (*ContextConfig, error) {
	cfg := &ContextConfig{}
	fs.StringVar(&cfg.RunID, "context-run-id", "", "Unique identifier for this diff run")
	fs.StringVar(&cfg.Environment, "context-env", "", "Environment name (e.g. staging, prod)")
	fs.StringVar(&cfg.TriggeredBy, "context-triggered-by", "", "Who or what triggered this run")
	var rawLabels string
	fs.StringVar(&rawLabels, "context-labels", "", "Comma-separated key=value labels")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	cfg.Labels = splitCSV(rawLabels)
	return cfg, nil
}

// BuildContextMetadata converts a ContextConfig into a diff.ContextMetadata.
// Returns nil if cfg is nil.
func BuildContextMetadata(cfg *ContextConfig) *diff.ContextMetadata {
	if cfg == nil {
		return nil
	}
	meta := &diff.ContextMetadata{
		RunID:       cfg.RunID,
		Environment: cfg.Environment,
		TriggeredBy: cfg.TriggeredBy,
	}
	for _, kv := range cfg.Labels {
		parts := splitKV(kv)
		if len(parts) == 2 {
			if meta.Labels == nil {
				meta.Labels = make(map[string]string)
			}
			meta.Labels[parts[0]] = parts[1]
		}
	}
	return meta
}

// splitKV splits a "key=value" string into ["key", "value"].
func splitKV(s string) []string {
	for i, c := range s {
		if c == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

// hostname is a helper used in run context enrichment.
func hostname() string {
	h, _ := os.Hostname()
	return h
}
