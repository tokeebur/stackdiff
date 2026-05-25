package diff

import (
	"testing"
)

func makePipelineReport(addresses ...string) *Report {
	var entries []ResourceChange
	for _, a := range addresses {
		entries = append(entries, ResourceChange{
			Address:      a,
			ResourceType: "aws_instance",
			Action:       ActionModified,
		})
	}
	return &Report{Changes: entries}
}

func TestPipeline_NilReport(t *testing.T) {
	p := NewPipeline()
	p.Add(func(r *Report) *Report { return r })
	if got := p.Run(nil); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestPipeline_NoSteps(t *testing.T) {
	p := NewPipeline()
	r := makePipelineReport("a", "b")
	got := p.Run(r)
	if got == nil || len(got.Changes) != 2 {
		t.Errorf("expected 2 changes unchanged, got %v", got)
	}
}

func TestPipeline_SingleStep_Transforms(t *testing.T) {
	p := NewPipeline()
	p.Add(func(r *Report) *Report {
		filtered := []ResourceChange{}
		for _, c := range r.Changes {
			if c.Address != "drop" {
				filtered = append(filtered, c)
			}
		}
		r.Changes = filtered
		return r
	})
	r := makePipelineReport("keep", "drop", "also-keep")
	got := p.Run(r)
	if len(got.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(got.Changes))
	}
}

func TestPipeline_MultipleSteps_Ordered(t *testing.T) {
	order := []string{}
	p := NewPipeline()
	p.Add(func(r *Report) *Report { order = append(order, "first"); return r })
	p.Add(func(r *Report) *Report { order = append(order, "second"); return r })
	p.Add(func(r *Report) *Report { order = append(order, "third"); return r })
	p.Run(makePipelineReport("x"))
	if len(order) != 3 || order[0] != "first" || order[1] != "second" || order[2] != "third" {
		t.Errorf("unexpected execution order: %v", order)
	}
}

func TestPipeline_StepReturnsNil_StopsExecution(t *testing.T) {
	called := false
	p := NewPipeline()
	p.Add(func(r *Report) *Report { return nil })
	p.Add(func(r *Report) *Report { called = true; return r })
	got := p.Run(makePipelineReport("a"))
	if got != nil {
		t.Errorf("expected nil report after nil step")
	}
	if called {
		t.Errorf("subsequent step should not have been called")
	}
}

func TestPipeline_Len(t *testing.T) {
	p := NewPipeline()
	if p.Len() != 0 {
		t.Errorf("expected 0, got %d", p.Len())
	}
	p.Add(func(r *Report) *Report { return r })
	p.Add(func(r *Report) *Report { return r })
	if p.Len() != 2 {
		t.Errorf("expected 2, got %d", p.Len())
	}
}
