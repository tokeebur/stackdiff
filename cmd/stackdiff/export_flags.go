package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/stackdiff/internal/diff"
)

// ExportConfig holds export-related CLI options.
type ExportConfig struct {
	Format diff.ExportFormat
	Output string // file path or "-" for stdout
}

// ParseExportFlags reads export-related flags from the provided FlagSet.
// It expects --format and --output flags to be registered.
func ParseExportFlags(fs *flag.FlagSet, args []string) (*ExportConfig, error) {
	formatStr := fs.String("format", "text", "Output format: text, markdown, json")
	output := fs.String("output", "-", "Output file path (use '-' for stdout)")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parse flags: %w", err)
	}

	fmt, err := diff.ParseExportFormat(*formatStr)
	if err != nil {
		return nil, err
	}

	return &ExportConfig{
		Format: fmt,
		Output: *output,
	}, nil
}

// ApplyExportConfig writes the report using the given ExportConfig.
func ApplyExportConfig(cfg *ExportConfig, r *diff.Report) error {
	if cfg == nil {
		return fmt.Errorf("export config must not be nil")
	}
	return diff.ExportReport(r, cfg.Format, cfg.Output, os.Stdout)
}
