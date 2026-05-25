package diff

import (
	"testing"
)

func baseNormalizeReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Action:       "modified",
				Before:       map[string]string{"name": "  web  ", "count": "null", "zone": "us-east-1"},
				After:        map[string]string{"name": "web", "count": "", "zone": "us-east-1"},
				Diff:         map[string][2]string{"name": {"  web  ", "web"}},
			},
		},
	}
}

func TestNormalizeReport_NilReport(t *testing.T) {
	result := NormalizeReport(nil, DefaultNormalizeOptions)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestNormalizeReport_TrimWhitespace(t *testing.T) {
	r := baseNormalizeReport()
	opts := NormalizeOptions{TrimWhitespace: true, StripNullAttrs: false}
	result := NormalizeReport(r, opts)
	if result.Changes[0].Before["name"] != "web" {
		t.Errorf("expected trimmed name, got %q", result.Changes[0].Before["name"])
	}
}

func TestNormalizeReport_StripNullAttrs(t *testing.T) {
	r := baseNormalizeReport()
	opts := NormalizeOptions{TrimWhitespace: false, StripNullAttrs: true}
	result := NormalizeReport(r, opts)
	if _, ok := result.Changes[0].Before["count"]; ok {
		t.Error("expected 'count' null attr to be stripped from Before")
	}
	if _, ok := result.Changes[0].After["count"]; ok {
		t.Error("expected 'count' empty attr to be stripped from After")
	}
}

func TestNormalizeReport_LowercaseKeys(t *testing.T) {
	r := &Report{
		Changes: []ResourceChange{
			{
				Address:      "aws_s3_bucket.logs",
				ResourceType: "aws_s3_bucket",
				Action:       "modified",
				Before:       map[string]string{"BucketName": "logs"},
				After:        map[string]string{"BucketName": "logs-v2"},
				Diff:         map[string][2]string{"BucketName": {"logs", "logs-v2"}},
			},
		},
	}
	opts := NormalizeOptions{LowercaseKeys: true}
	result := NormalizeReport(r, opts)
	if _, ok := result.Changes[0].Before["bucketname"]; !ok {
		t.Error("expected key to be lowercased in Before")
	}
	if _, ok := result.Changes[0].Diff["bucketname"]; !ok {
		t.Error("expected key to be lowercased in Diff")
	}
}

func TestNormalizeReport_PreservesUnchangedEntries(t *testing.T) {
	r := baseNormalizeReport()
	result := NormalizeReport(r, DefaultNormalizeOptions)
	if len(result.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Address != "aws_instance.web" {
		t.Errorf("unexpected address: %s", result.Changes[0].Address)
	}
}

func TestNormalizeReport_EmptyChanges(t *testing.T) {
	r := &Report{Changes: []ResourceChange{}}
	result := NormalizeReport(r, DefaultNormalizeOptions)
	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
}
