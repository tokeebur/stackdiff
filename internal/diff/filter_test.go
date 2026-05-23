package diff

import (
	"testing"
)

func makeRC(address, resourceType string) ResourceChange {
	return ResourceChange{
		Address:      address,
		ResourceType: resourceType,
	}
}

func baseReport() Report {
	return Report{
		Added: []ResourceChange{
			makeRC("aws_instance.web", "aws_instance"),
			makeRC("module.vpc.aws_subnet.public", "aws_subnet"),
		},
		Removed: []ResourceChange{
			makeRC("aws_s3_bucket.logs", "aws_s3_bucket"),
		},
		Modified: []ResourceChange{
			makeRC("aws_instance.worker", "aws_instance"),
			makeRC("module.vpc.aws_vpc.main", "aws_vpc"),
		},
	}
}

func TestFilterReport_NoFilter(t *testing.T) {
	r := baseReport()
	out := FilterReport(r, FilterOptions{})
	if len(out.Added) != 2 || len(out.Removed) != 1 || len(out.Modified) != 2 {
		t.Errorf("expected all resources, got added=%d removed=%d modified=%d",
			len(out.Added), len(out.Removed), len(out.Modified))
	}
}

func TestFilterReport_ByResourceType(t *testing.T) {
	r := baseReport()
	out := FilterReport(r, FilterOptions{ResourceType: "aws_instance"})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added aws_instance, got %d", len(out.Added))
	}
	if len(out.Removed) != 0 {
		t.Errorf("expected 0 removed aws_instance, got %d", len(out.Removed))
	}
	if len(out.Modified) != 1 {
		t.Errorf("expected 1 modified aws_instance, got %d", len(out.Modified))
	}
}

func TestFilterReport_ByAddressPrefix(t *testing.T) {
	r := baseReport()
	out := FilterReport(r, FilterOptions{AddressPrefix: "module.vpc"})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added in module.vpc, got %d", len(out.Added))
	}
	if len(out.Removed) != 0 {
		t.Errorf("expected 0 removed in module.vpc, got %d", len(out.Removed))
	}
	if len(out.Modified) != 1 {
		t.Errorf("expected 1 modified in module.vpc, got %d", len(out.Modified))
	}
}

func TestFilterReport_CombinedFilters(t *testing.T) {
	r := baseReport()
	out := FilterReport(r, FilterOptions{
		ResourceType:  "aws_subnet",
		AddressPrefix: "module.vpc",
	})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 result, got %d", len(out.Added))
	}
	if out.Added[0].Address != "module.vpc.aws_subnet.public" {
		t.Errorf("unexpected address: %s", out.Added[0].Address)
	}
}

func TestFilterReport_NoMatches(t *testing.T) {
	r := baseReport()
	out := FilterReport(r, FilterOptions{ResourceType: "aws_lambda_function"})
	if len(out.Added)+len(out.Removed)+len(out.Modified) != 0 {
		t.Error("expected no results for non-matching type")
	}
}
