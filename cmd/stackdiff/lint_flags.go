package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/stackdiff/internal/diff"
)

// LintConfig holds configuration parsed from lint-related CLI flags.
type LintConfig struct {
	Enabled bool
	Format  string
	FailOn  bool
}

// ParseLintFlags registers and parses lint-related flags from fs.
func ParseLintFlags(fs *flag.FlagSet) *LintConfig {
	cfg := &LintConfig{}
	fs.BoolVar(&cfg.Enabled, "lint", false, "run lint rules against the diff report")
	fs.StringVar(&cfg.Format, "lint-format", "text", "output format for lint results (text|markdown)")
	fs.BoolVar(&cfg.FailOn, "lint-fail", false, "exit with non-zero status if lint violations are found")
	return cfg
}

// ApplyLint runs lint rules against report if enabled, writes results, and
// optionally exits with a non-zero code when violations are found.
func ApplyLint(cfg *LintConfig, report *diff.Report) {
	if cfg == nil || !cfg.Enabled {
		return
	}
	rules := diff.DefaultLintRules()
	result := diff.LintReport(report, rules)
	if err := diff.WriteLintResult(os.Stdout, result, cfg.Format); err != nil {
		fmt.Fprintf(os.Stderr, "lint output error: %v\n", err)
		return
	}
	if cfg.FailOn && result.HasViolations() {
		os.Exit(2)
	}
}
