package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExportFormat represents a supported export format.
type ExportFormat string

const (
	ExportJSON     ExportFormat = "json"
	ExportMarkdown ExportFormat = "markdown"
	ExportText     ExportFormat = "text"
)

// ParseExportFormat parses and validates an export format string.
func ParseExportFormat(s string) (ExportFormat, error) {
	switch ExportFormat(strings.ToLower(strings.TrimSpace(s))) {
	case ExportJSON:
		return ExportJSON, nil
	case ExportMarkdown:
		return ExportMarkdown, nil
	case ExportText:
		return ExportText, nil
	default:
		return "", fmt.Errorf("unsupported export format %q: must be one of json, markdown, text", s)
	}
}

// ExportReport writes the report to the given path using the specified format.
// If path is "-" or empty, output is written to w (stdout fallback).
func ExportReport(r *Report, format ExportFormat, path string, w io.Writer) error {
	if r == nil {
		return fmt.Errorf("export: report must not be nil")
	}

	out := w
	if path != "" && path != "-" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("export: create directories: %w", err)
		}
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("export: create file: %w", err)
		}
		defer f.Close()
		out = f
	}

	switch format {
	case ExportJSON:
		return exportJSON(r, out)
	case ExportMarkdown:
		return WriteReport(r, out, "markdown")
	case ExportText:
		return WriteReport(r, out, "text")
	default:
		return fmt.Errorf("export: unknown format %q", format)
	}
}

func exportJSON(r *Report, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("export: encode JSON: %w", err)
	}
	return nil
}
