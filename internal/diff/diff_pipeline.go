package diff

// PipelineStep is a function that transforms a Report.
type PipelineStep func(*Report) *Report

// Pipeline is an ordered sequence of transformation steps applied to a Report.
type Pipeline struct {
	steps []PipelineStep
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends a PipelineStep to the pipeline.
func (p *Pipeline) Add(step PipelineStep) *Pipeline {
	p.steps = append(p.steps, step)
	return p
}

// Run executes all steps in order, passing the report through each.
// If the report becomes nil at any step, execution stops and nil is returned.
func (p *Pipeline) Run(r *Report) *Report {
	if r == nil {
		return nil
	}
	for _, step := range p.steps {
		r = step(r)
		if r == nil {
			return nil
		}
	}
	return r
}

// Len returns the number of steps registered in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.steps)
}
