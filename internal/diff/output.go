package diff

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat controls how the drift report is rendered.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatMarkdown OutputFormat = "markdown"
)

// WriteReport writes the drift report to w in the given format.
func WriteReport(w io.Writer, r *Report, format OutputFormat) error {
	switch format {
	case FormatText:
		return writeText(w, r)
	case FormatMarkdown:
		return writeMarkdown(w, r)
	default:
		return fmt.Errorf("unsupported output format: %q", format)
	}
}

func writeText(w io.Writer, r *Report) error {
	if !r.HasDrift() {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}
	for _, rc := range r.Added {
		fmt.Fprintf(w, "[+] %s (%s)\n", rc.Address, rc.ResourceType)
	}
	for _, rc := range r.Removed {
		fmt.Fprintf(w, "[-] %s (%s)\n", rc.Address, rc.ResourceType)
	}
	for _, rc := range r.Modified {
		fmt.Fprintf(w, "[~] %s (%s)\n", rc.Address, rc.ResourceType)
		for _, ch := range rc.Changes {
			fmt.Fprintf(w, "      %s: %q -> %q\n", ch.Attribute, ch.OldValue, ch.NewValue)
		}
	}
	return nil
}

func writeMarkdown(w io.Writer, r *Report) error {
	if !r.HasDrift() {
		_, err := fmt.Fprintln(w, "_No drift detected._")
		return err
	}
	fmt.Fprintln(w, "## Drift Summary")
	fmt.Fprintln(w, "")
	if len(r.Added) > 0 {
		fmt.Fprintln(w, "### Added")
		for _, rc := range r.Added {
			fmt.Fprintf(w, "- `%s` (%s)\n", rc.Address, rc.ResourceType)
		}
		fmt.Fprintln(w, "")
	}
	if len(r.Removed) > 0 {
		fmt.Fprintln(w, "### Removed")
		for _, rc := range r.Removed {
			fmt.Fprintf(w, "- `%s` (%s)\n", rc.Address, rc.ResourceType)
		}
		fmt.Fprintln(w, "")
	}
	if len(r.Modified) > 0 {
		fmt.Fprintln(w, "### Modified")
		for _, rc := range r.Modified {
			fmt.Fprintf(w, "- `%s` (%s)\n", rc.Address, rc.ResourceType)
			for _, ch := range rc.Changes {
				fmt.Fprintf(w, "  - `%s`: `%s` → `%s`\n",
					ch.Attribute,
					strings.ReplaceAll(ch.OldValue, "`", "'"),
					strings.ReplaceAll(ch.NewValue, "`", "'"),
				)
			}
		}
	}
	return nil
}
