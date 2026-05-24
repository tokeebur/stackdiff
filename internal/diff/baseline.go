package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Baseline represents a saved snapshot of a Report for future comparison.
type Baseline struct {
	CreatedAt time.Time     `json:"created_at"`
	Report    *Report       `json:"report"`
}

// SaveBaseline writes the given report as a JSON baseline to w.
func SaveBaseline(w io.Writer, r *Report) error {
	if r == nil {
		return fmt.Errorf("baseline: report must not be nil")
	}
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Report:    r,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// LoadBaseline reads a JSON baseline from r.
func LoadBaseline(r io.Reader) (*Baseline, error) {
	var b Baseline
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: failed to decode: %w", err)
	}
	if b.Report == nil {
		return nil, fmt.Errorf("baseline: missing report field")
	}
	return &b, nil
}

// SaveBaselineFile writes a baseline to the given file path.
func SaveBaselineFile(path string, r *Report) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: cannot create file %q: %w", path, err)
	}
	defer f.Close()
	return SaveBaseline(f, r)
}

// LoadBaselineFile reads a baseline from the given file path.
func LoadBaselineFile(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: cannot open file %q: %w", path, err)
	}
	defer f.Close()
	return LoadBaseline(f)
}
