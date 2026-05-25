package diff

// FilterChain applies a sequence of named filter steps to a Report,
// recording which step (if any) removed each entry.
type FilterChain struct {
	steps []filterStep
}

type filterStep struct {
	name string
	fn   func(*Report) *Report
}

// FilterChainResult holds the final report and a log of entries removed per step.
type FilterChainResult struct {
	Report       *Report
	RemovedByStep map[string]int
}

// NewFilterChain creates an empty FilterChain.
func NewFilterChain() *FilterChain {
	return &FilterChain{}
}

// Add appends a named filter step to the chain.
func (fc *FilterChain) Add(name string, fn func(*Report) *Report) *FilterChain {
	fc.steps = append(fc.steps, filterStep{name: name, fn: fn})
	return fc
}

// Apply runs all steps in order and returns a FilterChainResult.
// If report is nil, a zero-value result is returned.
func (fc *FilterChain) Apply(report *Report) FilterChainResult {
	result := FilterChainResult{
		RemovedByStep: make(map[string]int),
	}
	if report == nil {
		return result
	}

	current := report
	for _, step := range fc.steps {
		before := len(current.Changes)
		current = step.fn(current)
		if current == nil {
			current = &Report{}
		}
		after := len(current.Changes)
		removed := before - after
		if removed < 0 {
			removed = 0
		}
		result.RemovedByStep[step.name] = removed
	}

	result.Report = current
	return result
}

// StepNames returns the names of all registered steps in order.
func (fc *FilterChain) StepNames() []string {
	names := make([]string, len(fc.steps))
	for i, s := range fc.steps {
		names[i] = s.name
	}
	return names
}
