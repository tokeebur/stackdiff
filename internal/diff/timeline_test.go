package diff

import (
	"bytes"
	"testing"
	"time"
)

func makeStats(added, removed, modified int) DriftStats {
	return DriftStats{
		Added:    added,
		Removed:  removed,
		Modified: modified,
		Total:    added + removed + modified,
	}
}

func TestTimeline_AddEntry(t *testing.T) {
	var tl Timeline
	tl.AddEntry(makeStats(1, 2, 3), "run-1")
	if tl.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tl.Len())
	}
	e := tl.Entries[0]
	if e.Added != 1 || e.Removed != 2 || e.Modified != 3 || e.Total != 6 {
		t.Errorf("unexpected entry values: %+v", e)
	}
	if e.Label != "run-1" {
		t.Errorf("expected label run-1, got %s", e.Label)
	}
}

func TestTimeline_Latest_Empty(t *testing.T) {
	var tl Timeline
	if tl.Latest() != nil {
		t.Fatal("expected nil for empty timeline")
	}
}

func TestTimeline_Latest_NonEmpty(t *testing.T) {
	var tl Timeline
	tl.AddEntry(makeStats(1, 0, 0), "first")
	tl.AddEntry(makeStats(0, 1, 0), "last")
	if tl.Latest().Label != "last" {
		t.Errorf("expected last, got %s", tl.Latest().Label)
	}
}

func TestTimeline_SortByTime(t *testing.T) {
	var tl Timeline
	now := time.Now().UTC()
	tl.Entries = []TimelineEntry{
		{Timestamp: now.Add(2 * time.Second), Label: "c"},
		{Timestamp: now, Label: "a"},
		{Timestamp: now.Add(time.Second), Label: "b"},
	}
	tl.SortByTime()
	if tl.Entries[0].Label != "a" || tl.Entries[1].Label != "b" || tl.Entries[2].Label != "c" {
		t.Errorf("unexpected order after sort: %v", tl.Entries)
	}
}

func TestSaveAndLoadTimeline_RoundTrip(t *testing.T) {
	var tl Timeline
	tl.AddEntry(makeStats(2, 1, 3), "snapshot")

	var buf bytes.Buffer
	if err := SaveTimeline(&tl, &buf); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadTimeline(&buf)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", loaded.Len())
	}
	if loaded.Entries[0].Label != "snapshot" {
		t.Errorf("label mismatch: %s", loaded.Entries[0].Label)
	}
}

func TestSaveTimeline_Nil(t *testing.T) {
	var buf bytes.Buffer
	if err := SaveTimeline(nil, &buf); err == nil {
		t.Fatal("expected error for nil timeline")
	}
}
