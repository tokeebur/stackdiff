package diff

import (
	"testing"
)

func makeStatsReport() Report {
	return Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.a", ResourceType: "aws_instance", Action: ActionAdd},
			{Address: "aws_instance.b", ResourceType: "aws_instance", Action: ActionRemove},
			{Address: "aws_s3_bucket.x", ResourceType: "aws_s3_bucket", Action: ActionModify},
			{Address: "aws_s3_bucket.y", ResourceType: "aws_s3_bucket", Action: ActionModify},
			{Address: "aws_vpc.main", ResourceType: "aws_vpc", Action: ActionAdd},
		},
	}
}

func TestComputeStats_Counts(t *testing.T) {
	r := makeStatsReport()
	s := ComputeStats(r)

	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Modified != 2 {
		t.Errorf("expected Modified=2, got %d", s.Modified)
	}
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
}

func TestComputeStats_Empty(t *testing.T) {
	s := ComputeStats(Report{})
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestStatsByType_Groups(t *testing.T) {
	r := makeStatsReport()
	byType := StatsByType(r)

	if byType["aws_instance"].Added != 1 {
		t.Errorf("expected aws_instance Added=1, got %d", byType["aws_instance"].Added)
	}
	if byType["aws_instance"].Removed != 1 {
		t.Errorf("expected aws_instance Removed=1, got %d", byType["aws_instance"].Removed)
	}
	if byType["aws_s3_bucket"].Modified != 2 {
		t.Errorf("expected aws_s3_bucket Modified=2, got %d", byType["aws_s3_bucket"].Modified)
	}
	if byType["aws_vpc"].Added != 1 {
		t.Errorf("expected aws_vpc Added=1, got %d", byType["aws_vpc"].Added)
	}
	if byType["aws_vpc"].Total != 1 {
		t.Errorf("expected aws_vpc Total=1, got %d", byType["aws_vpc"].Total)
	}
}

func TestStatsByType_Empty(t *testing.T) {
	byType := StatsByType(Report{})
	if len(byType) != 0 {
		t.Errorf("expected empty map, got %d entries", len(byType))
	}
}
