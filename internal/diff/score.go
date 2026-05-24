package diff

import "fmt"

// DriftScore represents a weighted risk score for a drift report.
type DriftScore struct {
	Total    float64
	Added    float64
	Removed  float64
	Modified float64
	Label    string
}

// ScoreWeights controls how each action type contributes to the total score.
type ScoreWeights struct {
	Added    float64
	Removed  float64
	Modified float64
}

// DefaultScoreWeights returns sensible default weights.
// Removals are weighted highest as they are most destructive.
func DefaultScoreWeights() ScoreWeights {
	return ScoreWeights{
		Added:    1.0,
		Removed:  2.0,
		Modified: 1.5,
	}
}

// ScoreReport computes a DriftScore from a Report using the provided weights.
// Returns an error if the report is nil.
func ScoreReport(r *Report, w ScoreWeights) (*DriftScore, error) {
	if r == nil {
		return nil, fmt.Errorf("score: report must not be nil")
	}

	var added, removed, modified float64

	for _, e := range r.Entries {
		switch e.Action {
		case ActionAdded:
			added += w.Added
		case ActionRemoved:
			removed += w.Removed
		case ActionModified:
			modified += w.Modified * float64(1+len(e.ChangedAttrs))
		}
	}

	total := added + removed + modified

	return &DriftScore{
		Total:    total,
		Added:    added,
		Removed:  removed,
		Modified: modified,
		Label:    scoreLabel(total),
	}, nil
}

func scoreLabel(total float64) string {
	switch {
	case total == 0:
		return "none"
	case total < 5:
		return "low"
	case total < 15:
		return "medium"
	default:
		return "high"
	}
}
