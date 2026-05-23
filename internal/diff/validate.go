package diff

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError holds one or more validation failures for a Report.
type ValidationError struct {
	Messages []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("report validation failed: %s", strings.Join(e.Messages, "; "))
}

// ValidateReport checks a Report for internal consistency and returns a
// ValidationError if any issues are found.
func ValidateReport(r *Report) error {
	if r == nil {
		return errors.New("report is nil")
	}

	var msgs []string

	for i, rc := range r.Changes {
		if rc.Address == "" {
			msgs = append(msgs, fmt.Sprintf("change[%d]: address is empty", i))
		}
		if rc.ResourceType == "" {
			msgs = append(msgs, fmt.Sprintf("change[%d] (%s): resource_type is empty", i, rc.Address))
		}
		switch rc.Action {
		case ActionAdd, ActionRemove, ActionModify:
			// valid
		default:
			msgs = append(msgs, fmt.Sprintf("change[%d] (%s): unknown action %q", i, rc.Address, rc.Action))
		}
		if rc.Action == ActionModify && len(rc.Attributes) == 0 {
			msgs = append(msgs, fmt.Sprintf("change[%d] (%s): modified resource has no attribute diffs", i, rc.Address))
		}
	}

	if len(msgs) > 0 {
		return &ValidationError{Messages: msgs}
	}
	return nil
}
