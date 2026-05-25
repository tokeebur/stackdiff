package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/your-org/stackdiff/internal/diff"
)

// AuditConfig holds configuration for audit log output.
type AuditConfig struct {
	Enabled    bool
	OutputPath string
}

// ParseAuditFlags registers and parses audit-related CLI flags from fs.
func ParseAuditFlags(fs *flag.FlagSet) *AuditConfig {
	c := &AuditConfig{}
	fs.BoolVar(&c.Enabled, "audit", false, "emit an audit trail of pipeline operations")
	fs.StringVar(&c.OutputPath, "audit-out", "", "path to write audit log (default: stdout)")
	return c
}

// WriteAuditConfig writes the audit log to the configured destination.
// It is a no-op when cfg is nil or audit is not enabled.
func WriteAuditConfig(cfg *AuditConfig, a *diff.AuditLog) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	var w io.Writer = os.Stdout
	if cfg.OutputPath != "" {
		f, err := os.Create(cfg.OutputPath)
		if err != nil {
			return fmt.Errorf("audit: could not open output file %q: %w", cfg.OutputPath, err)
		}
		defer f.Close()
		w = f
	}

	return diff.WriteAuditLog(a, w)
}

// RecordStep is a convenience wrapper that records a pipeline step in the
// audit log when auditing is enabled.
func RecordStep(cfg *AuditConfig, a *diff.AuditLog, op, detail string, before, after int) {
	if cfg == nil || !cfg.Enabled || a == nil {
		return
	}
	a.Record(op, detail, before, after)
}
