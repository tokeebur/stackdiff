package diff

import (
	"fmt"
	"io"
	"strings"
)

// WriteLintResult writes a human-readable lint summary to w.
// format should be "text" or "markdown".
func WriteLintResult(w io.Writer, result *LintResult, format string) error {
	if result == nil {
		return nil
	}
	switch strings.ToLower(format) {
	case "markdown":
		return writeLintMarkdown(w, result)
	default:
		return writeLintText(w, result)
	}
}

func writeLintText(w io.Writer, result *LintResult) error {
	if !result.HasViolations() {
		_, err := fmt.Fprintln(w, "lint: no violations found")
		return err
	}
	_, err := fmt.Fprintf(w, "lint: %d violation(s) found\n", len(result.Violations))
	if err != nil {
		return err
	}
	for _, v := range result.Violations {
		_, err = fmt.Fprintf(w, "  [%s] %s\n", v.Rule, v.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLintMarkdown(w io.Writer, result *LintResult) error {
	if !result.HasViolations() {
		_, err := fmt.Fprintln(w, "**Lint:** no violations found")
		return err
	}
	_, err := fmt.Fprintf(w, "## Lint Violations (%d)\n\n", len(result.Violations))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "| Address | Rule | Message |")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "|---------|------|---------|")
	if err != nil {
		return err
	}
	for _, v := range result.Violations {
		_, err = fmt.Fprintf(w, "| %s | %s | %s |\n", v.Address, v.Rule, v.Message)
		if err != nil {
			return err
		}
	}
	return nil
}
