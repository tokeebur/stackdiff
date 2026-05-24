package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// SaveTimeline serialises the timeline to the given writer as JSON.
func SaveTimeline(t *Timeline, w io.Writer) error {
	if t == nil {
		return fmt.Errorf("timeline is nil")
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(t)
}

// LoadTimeline deserialises a timeline from the given reader.
func LoadTimeline(r io.Reader) (*Timeline, error) {
	var t Timeline
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return nil, fmt.Errorf("decode timeline: %w", err)
	}
	return &t, nil
}

// SaveTimelineFile writes the timeline to the named file.
func SaveTimelineFile(t *Timeline, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create timeline file: %w", err)
	}
	defer f.Close()
	return SaveTimeline(t, f)
}

// LoadTimelineFile reads a timeline from the named file.
func LoadTimelineFile(path string) (*Timeline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open timeline file: %w", err)
	}
	defer f.Close()
	return LoadTimeline(f)
}
