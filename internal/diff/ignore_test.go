package diff

import (
	"testing"
)

func baseIgnoreReport() Report {
	return Report{
		Added: []ResourceChange{
			{Address: "aws_s3_bucket.logs", ResourceType: "aws_s3_bucket"},
			{Address: "aws_iam_role.worker", ResourceType: "aws_iam_role"},
		},
		Removed: []ResourceChange{
			{Address: "aws_instance.old", ResourceType: "aws_instance"},
		},
		Modified: []ResourceChange{
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Changes: map[string]AttributeChange{
					"tags":          {Old: "v1", New: "v2"},
					"instance_type": {Old: "t2.micro", New: "t3.micro"},
				},
			},
		},
	}
}

func TestApplyIgnoreRules_NoRules(t *testing.T) {
	r := baseIgnoreReport()
	out := ApplyIgnoreRules(r, nil)
	if len(out.Added) != 2 || len(out.Removed) != 1 || len(out.Modified) != 1 {
		t.Errorf("expected report unchanged, got added=%d removed=%d modified=%d",
			len(out.Added), len(out.Removed), len(out.Modified))
	}
}

func TestApplyIgnoreRules_ByResourceType(t *testing.T) {
	r := baseIgnoreReport()
	rules := []IgnoreRule{{ResourceType: "aws_s3_bucket"}}
	out := ApplyIgnoreRules(r, rules)
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added resource, got %d", len(out.Added))
	}
	if out.Added[0].ResourceType == "aws_s3_bucket" {
		t.Error("aws_s3_bucket should have been ignored")
	}
}

func TestApplyIgnoreRules_ByAddressPrefix(t *testing.T) {
	r := baseIgnoreReport()
	rules := []IgnoreRule{{AddressPrefix: "aws_iam_"}}
	out := ApplyIgnoreRules(r, rules)
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added resource, got %d", len(out.Added))
	}
}

func TestApplyIgnoreRules_ByAttributeKey(t *testing.T) {
	r := baseIgnoreReport()
	rules := []IgnoreRule{{AttributeKey: "tags"}}
	out := ApplyIgnoreRules(r, rules)
	if len(out.Modified) != 1 {
		t.Fatalf("expected 1 modified resource, got %d", len(out.Modified))
	}
	if _, ok := out.Modified[0].Changes["tags"]; ok {
		t.Error("tags attribute should have been stripped")
	}
	if _, ok := out.Modified[0].Changes["instance_type"]; !ok {
		t.Error("instance_type attribute should remain")
	}
}

func TestApplyIgnoreRules_AttributeKeyRemovesEntireEntry(t *testing.T) {
	r := Report{
		Modified: []ResourceChange{
			{
				Address:      "aws_instance.web",
				ResourceType: "aws_instance",
				Changes: map[string]AttributeChange{
					"tags": {Old: "v1", New: "v2"},
				},
			},
		},
	}
	rules := []IgnoreRule{{AttributeKey: "tags"}}
	out := ApplyIgnoreRules(r, rules)
	if len(out.Modified) != 0 {
		t.Errorf("expected modified to be empty after stripping only attr, got %d", len(out.Modified))
	}
}
