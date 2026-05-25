package diff

import (
	"encoding/json"
	"fmt"
	"io"
)

// ExportFormat for context-aware export.
type ContextExportFormat string

const (
	ContextExportJSON ContextExportFormat = "json"
	ContextExportText ContextExportFormat = "text"
)

// ExportDiffContext serialises a DiffContext to the given writer.
// Supported formats: "json", "text".
func ExportDiffContext(dc *DiffContext, format ContextExportFormat, w io.Writer) error {
	if dc == nil {
		return fmt.Errorf("diff context is nil")
	}
	switch format {
	case ContextExportJSON:
		return exportContextJSON(dc, w)
	case ContextExportText:
		return exportContextText(dc, w)
	default:
		return fmt.Errorf("unsupported context export format: %s", format)
	}
}

func exportContextJSON(dc *DiffContext, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(dc)
}

func exportContextText(dc *DiffContext, w io.Writer) error {
	if dc.Metadata != nil {
		fmt.Fprintf(w, "Run ID:      %s\n", dc.Metadata.RunID)
		fmt.Fprintf(w, "Environment: %s\n", dc.Metadata.Environment)
		fmt.Fprintf(w, "Triggered By:%s\n", dc.Metadata.TriggeredBy)
		fmt.Fprintf(w, "Timestamp:   %s\n", dc.Metadata.Timestamp.Format("2006-01-02T15:04:05Z"))
		for k, v := range dc.Metadata.Labels {
			fmt.Fprintf(w, "Label[%s]: %s\n", k, v)
		}
	}
	if dc.Report != nil {
		fmt.Fprintf(w, "Entries: %d\n", len(dc.Report.Entries))
	}
	return nil
}
