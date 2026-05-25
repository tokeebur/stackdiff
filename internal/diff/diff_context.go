package diff

import "time"

// ContextMetadata holds optional metadata attached to a diff run.
type ContextMetadata struct {
	RunID       string            `json:"run_id,omitempty"`
	Environment string            `json:"environment,omitempty"`
	TriggeredBy string            `json:"triggered_by,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// DiffContext wraps a Report with optional metadata for richer output.
type DiffContext struct {
	Report   *Report         `json:"report"`
	Metadata *ContextMetadata `json:"metadata,omitempty"`
}

// NewDiffContext creates a DiffContext from a report and optional metadata.
// If metadata is nil a default one with the current timestamp is used.
func NewDiffContext(r *Report, meta *ContextMetadata) *DiffContext {
	if meta == nil {
		meta = &ContextMetadata{Timestamp: time.Now().UTC()}
	}
	if meta.Timestamp.IsZero() {
		meta.Timestamp = time.Now().UTC()
	}
	return &DiffContext{Report: r, Metadata: meta}
}

// WithLabel adds a label to the context metadata, initialising the map if needed.
func (dc *DiffContext) WithLabel(key, value string) *DiffContext {
	if dc.Metadata == nil {
		dc.Metadata = &ContextMetadata{Timestamp: time.Now().UTC()}
	}
	if dc.Metadata.Labels == nil {
		dc.Metadata.Labels = make(map[string]string)
	}
	dc.Metadata.Labels[key] = value
	return dc
}
