package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/you/stackdiff/internal/state"
)

const validStateJSON = `{
  "version": 4,
  "terraform_version": "1.5.0",
  "serial": 10,
  "lineage": "abc-123",
  "resources": [
    {
      "mode": "managed",
      "type": "aws_instance",
      "name": "web",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "schema_version": 1,
          "attributes": {
            "id": "i-0abc123",
            "instance_type": "t3.micro"
          },
          "sensitive_attributes": []
        }
      ]
    }
  ]
}`

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "terraform.tfstate")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp state: %v", err)
	}
	return path
}

func TestParseStateFile_Valid(t *testing.T) {
	path := writeTempState(t, validStateJSON)

	s, err := state.ParseStateFile(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if s.Version != 4 {
		t.Errorf("expected version 4, got %d", s.Version)
	}
	if len(s.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(s.Resources))
	}
	if s.Resources[0].Type != "aws_instance" {
		t.Errorf("expected type aws_instance, got %s", s.Resources[0].Type)
	}
}

func TestParseStateFile_UnsupportedVersion(t *testing.T) {
	invalid := `{"version": 3, "terraform_version": "0.12.0", "resources": []}`
	path := writeTempState(t, invalid)

	_, err := state.ParseStateFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported version, got nil")
	}
}

func TestParseStateFile_MissingFile(t *testing.T) {
	_, err := state.ParseStateFile("/nonexistent/path/terraform.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestResourceMap(t *testing.T) {
	path := writeTempState(t, validStateJSON)
	s, err := state.ParseStateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rm := s.ResourceMap()
	if _, ok := rm["aws_instance.web"]; !ok {
		t.Errorf("expected key 'aws_instance.web' in resource map")
	}
}
