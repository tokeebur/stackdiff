package diff

import (
	"testing"
)

func makeChainReport(addresses ...string) *Report {
	r := &Report{}
	for _, addr := range addresses {
		r.Changes = append(r.Changes, ResourceChange{
			Address:      addr,
			ResourceType: "aws_instance",
			Action:       ActionAdded,
		})
	}
	return r
}

func TestFilterChain_NilReport(t *testing.T) {
	fc := NewFilterChain()
	result := fc.Apply(nil)
	if result.Report != nil {
		t.Errorf("expected nil report, got %v", result.Report)
	}
	if len(result.RemovedByStep) != 0 {
		t.Errorf("expected empty removed map")
	}
}

func TestFilterChain_NoSteps(t *testing.T) {
	report := makeChainReport("a", "b", "c")
	fc := NewFilterChain()
	result := fc.Apply(report)
	if len(result.Report.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(result.Report.Changes))
	}
}

func TestFilterChain_SingleStep_RemovesEntries(t *testing.T) {
	report := makeChainReport("keep.one", "drop.two", "keep.three")
	fc := NewFilterChain()
	fc.Add("drop-prefix", func(r *Report) *Report {
		out := &Report{}
		for _, c := range r.Changes {
			if len(c.Address) >= 4 && c.Address[:4] == "keep" {
				out.Changes = append(out.Changes, c)
			}
		}
		return out
	})
	result := fc.Apply(report)
	if len(result.Report.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(result.Report.Changes))
	}
	if result.RemovedByStep["drop-prefix"] != 1 {
		t.Errorf("expected 1 removed by drop-prefix, got %d", result.RemovedByStep["drop-prefix"])
	}
}

func TestFilterChain_MultipleSteps_AccumulatesRemovals(t *testing.T) {
	report := makeChainReport("a", "b", "c", "d")
	fc := NewFilterChain()
	fc.Add("step1", func(r *Report) *Report {
		return &Report{Changes: r.Changes[:3]} // remove last
	})
	fc.Add("step2", func(r *Report) *Report {
		return &Report{Changes: r.Changes[:1]} // remove 2 more
	})
	result := fc.Apply(report)
	if len(result.Report.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(result.Report.Changes))
	}
	if result.RemovedByStep["step1"] != 1 {
		t.Errorf("expected step1 to remove 1, got %d", result.RemovedByStep["step1"])
	}
	if result.RemovedByStep["step2"] != 2 {
		t.Errorf("expected step2 to remove 2, got %d", result.RemovedByStep["step2"])
	}
}

func TestFilterChain_StepNames(t *testing.T) {
	fc := NewFilterChain()
	fc.Add("alpha", func(r *Report) *Report { return r })
	fc.Add("beta", func(r *Report) *Report { return r })
	names := fc.StepNames()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected step names: %v", names)
	}
}
