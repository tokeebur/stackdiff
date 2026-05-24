package diff

// Annotation holds a human-readable label and optional detail attached to a
// resource change entry in a Report.
type Annotation struct {
	Label  string
	Detail string
}

// AnnotateReport attaches annotations to report entries based on a set of
// AnnotationRules. Each rule maps a resource type prefix to an Annotation.
// Existing annotations on an entry are replaced when a rule matches.
func AnnotateReport(r *Report, rules []AnnotationRule) *Report {
	if r == nil || len(rules) == 0 {
		return r
	}

	annotated := make([]ResourceChange, 0, len(r.Changes))
	for _, rc := range r.Changes {
		for _, rule := range rules {
			if rule.matches(rc) {
				rc.Annotation = &Annotation{
					Label:  rule.Label,
					Detail: rule.Detail,
				}
				break
			}
		}
		annotated = append(annotated, rc)
	}

	return &Report{
		Changes: annotated,
	}
}

// AnnotationRule describes when an annotation should be applied.
type AnnotationRule struct {
	// ResourceType is matched as a prefix against ResourceChange.ResourceType.
	ResourceType string
	// Action restricts the rule to a specific action ("added", "removed",
	// "modified"). An empty string matches any action.
	Action string
	Label  string
	Detail string
}

func (ar AnnotationRule) matches(rc ResourceChange) bool {
	if ar.ResourceType != "" && rc.ResourceType != ar.ResourceType {
		return false
	}
	if ar.Action != "" && string(rc.Action) != ar.Action {
		return false
	}
	return true
}
