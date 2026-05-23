package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// TerraformState represents the top-level structure of a Terraform state file.
type TerraformState struct {
	Version          int       `json:"version"`
	TerraformVersion string    `json:"terraform_version"`
	Serial           int64     `json:"serial"`
	Lineage          string    `json:"lineage"`
	Resources        []Resource `json:"resources"`
}

// Resource represents a single resource block in the state file.
type Resource struct {
	Module    string     `json:"module,omitempty"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

// Instance holds the attributes of a resource instance.
type Instance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
	SensitiveAttributes []interface{}   `json:"sensitive_attributes"`
}

// ResourceKey returns a unique string key for a resource.
func (r Resource) ResourceKey() string {
	if r.Module != "" {
		return fmt.Sprintf("%s.%s.%s", r.Module, r.Type, r.Name)
	}
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

// ParseStateFile reads and parses a Terraform state file from the given path.
func ParseStateFile(path string) (*TerraformState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file %q: %w", path, err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file %q: %w", path, err)
	}

	if state.Version != 4 {
		return nil, fmt.Errorf("unsupported state file version %d (only version 4 is supported)", state.Version)
	}

	return &state, nil
}

// ResourceMap returns a map of resource keys to Resource structs for quick lookup.
func (s *TerraformState) ResourceMap() map[string]Resource {
	rm := make(map[string]Resource, len(s.Resources))
	for _, r := range s.Resources {
		rm[r.ResourceKey()] = r
	}
	return rm
}
