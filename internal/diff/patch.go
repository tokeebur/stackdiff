package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// PatchEntry represents a single change that could be applied or reviewed.
type PatchEntry struct {
	Address      string            `json:"address"`
	ResourceType string            `json:"resource_type"`
	Action       string            `json:"action"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

// Patch is a collection of PatchEntry items derived from a Report.
type Patch struct {
	Entries []PatchEntry `json:"entries"`
}

// BuildPatch converts a Report into a Patch for serialisation or further processing.
func BuildPatch(r *Report) (*Patch, error) {
	if r == nil {
		return nil, fmt.Errorf("report must not be nil")
	}

	p := &Patch{}

	for _, rc := range r.Added {
		p.Entries = append(p.Entries, PatchEntry{
			Address:      rc.Address,
			ResourceType: rc.ResourceType,
			Action:       "add",
			Attributes:   copyAttrs(rc.Attributes),
		})
	}

	for _, rc := range r.Removed {
		p.Entries = append(p.Entries, PatchEntry{
			Address:      rc.Address,
			ResourceType: rc.ResourceType,
			Action:       "remove",
			Attributes:   copyAttrs(rc.Attributes),
		})
	}

	for _, rc := range r.Modified {
		p.Entries = append(p.Entries, PatchEntry{
			Address:      rc.Address,
			ResourceType: rc.ResourceType,
			Action:       "modify",
			Attributes:   copyAttrs(rc.Attributes),
		})
	}

	sort.Slice(p.Entries, func(i, j int) bool {
		return p.Entries[i].Address < p.Entries[j].Address
	})

	return p, nil
}

// WritePatchJSON serialises the Patch as JSON to the provided writer.
func WritePatchJSON(w io.Writer, p *Patch) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func copyAttrs(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
