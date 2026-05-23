package diff

import "github.com/user/stackdiff/internal/state"

// CompareToReport runs a full comparison between two parsed states and returns
// a structured Report instead of the raw DiffResult slice.
func CompareToReport(a, b *state.State) Report {
	results := Compare(a, b)

	report := Report{}
	aMap := a.ResourceMap()
	bMap := b.ResourceMap()

	for _, res := range results {
		switch res.ChangeType {
		case "added":
			report.Added = append(report.Added, res.Address)
		case "removed":
			report.Removed = append(report.Removed, res.Address)
		case "modified":
			oldAttrs := flattenAttrs(aMap[res.Address].AttributeValues)
			newAttrs := flattenAttrs(bMap[res.Address].AttributeValues)
			changes := make(map[string]AttributeChange)
			for _, k := range mergeMaps(oldAttrs, newAttrs) {
				ov := oldAttrs[k]
				nv := newAttrs[k]
				if ov != nv {
					changes[k] = AttributeChange{Old: ov, New: nv}
				}
			}
			report.Modified = append(report.Modified, ModifiedResource{
				Address: res.Address,
				Changes: changes,
			})
		case "unchanged":
			report.Unchanged++
		}
	}
	return report
}

// flattenAttrs converts interface{} attribute values to string representations.
func flattenAttrs(attrs map[string]interface{}) map[string]string {
	out := make(map[string]string, len(attrs))
	for k, v := range attrs {
		if v == nil {
			out[k] = ""
		} else {
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out
}
