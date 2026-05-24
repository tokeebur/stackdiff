package diff

import (
	"testing"
)

func baseAnnotateReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: ActionAdded},
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: ActionModified},
			{Address: "aws_iam_role.exec", ResourceType: "aws_iam_role", Action: ActionRemoved},
		},
	}
}

func TestAnnotateReport_NilReport(t *testing.T) {
	result := AnnotateReport(nil, []AnnotationRule{{ResourceType: "aws_instance", Label: "compute"}})
	if result != nil {
		t.Fatal("expected nil for nil report")
	}
}

func TestAnnotateReport_NoRules(t *testing.T) {
	r := baseAnnotateReport()
	result := AnnotateReport(r, nil)
	if len(result.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(result.Changes))
	}
	for _, rc := range result.Changes {
		if rc.Annotation != nil {
			t.Errorf("expected no annotation on %s", rc.Address)
		}
	}
}

func TestAnnotateReport_ByResourceType(t *testing.T) {
	r := baseAnnotateReport()
	rules := []AnnotationRule{
		{ResourceType: "aws_instance", Label: "compute", Detail: "EC2 instance"},
	}
	result := AnnotateReport(r, rules)

	for _, rc := range result.Changes {
		if rc.ResourceType == "aws_instance" {
			if rc.Annotation == nil {
				t.Fatal("expected annotation on aws_instance")
			}
			if rc.Annotation.Label != "compute" {
				t.Errorf("expected label 'compute', got %q", rc.Annotation.Label)
			}
		} else if rc.Annotation != nil {
			t.Errorf("unexpected annotation on %s", rc.Address)
		}
	}
}

func TestAnnotateReport_ByAction(t *testing.T) {
	r := baseAnnotateReport()
	rules := []AnnotationRule{
		{Action: "removed", Label: "destructive", Detail: "resource will be destroyed"},
	}
	result := AnnotateReport(r, rules)

	for _, rc := range result.Changes {
		if rc.Action == ActionRemoved {
			if rc.Annotation == nil || rc.Annotation.Label != "destructive" {
				t.Errorf("expected 'destructive' annotation on removed resource %s", rc.Address)
			}
		} else if rc.Annotation != nil {
			t.Errorf("unexpected annotation on %s", rc.Address)
		}
	}
}

func TestAnnotateReport_FirstRuleWins(t *testing.T) {
	r := &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: ActionAdded},
		},
	}
	rules := []AnnotationRule{
		{ResourceType: "aws_instance", Label: "first"},
		{ResourceType: "aws_instance", Label: "second"},
	}
	result := AnnotateReport(r, rules)
	if result.Changes[0].Annotation == nil {
		t.Fatal("expected annotation")
	}
	if result.Changes[0].Annotation.Label != "first" {
		t.Errorf("expected first rule to win, got %q", result.Changes[0].Annotation.Label)
	}
}
