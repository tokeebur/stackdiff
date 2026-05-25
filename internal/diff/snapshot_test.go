package diff

import (
	"bytes"
	"testing"
)

func makeSnapshotReport() *Report {
	return &Report{
		Entries: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: "added"},
		},
	}
}

func TestNewSnapshot_NilReport(t *testing.T) {
	_, err := NewSnapshot("test", nil)
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestNewSnapshot_ValidReport(t *testing.T) {
	snap, err := NewSnapshot("v1", makeSnapshotReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Label != "v1" {
		t.Errorf("expected label v1, got %q", snap.Label)
	}
	if snap.Stats == nil {
		t.Error("expected stats to be populated")
	}
	if snap.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestSnapshotStore_Latest_Empty(t *testing.T) {
	store := &SnapshotStore{}
	if store.Latest() != nil {
		t.Error("expected nil for empty store")
	}
}

func TestSnapshotStore_Latest_NonEmpty(t *testing.T) {
	store := &SnapshotStore{}
	snap1, _ := NewSnapshot("a", makeSnapshotReport())
	snap2, _ := NewSnapshot("b", makeSnapshotReport())
	store.AddSnapshot(*snap1)
	store.AddSnapshot(*snap2)
	if store.Latest().Label != "b" {
		t.Errorf("expected latest label b, got %q", store.Latest().Label)
	}
}

func TestSnapshotStore_FindByLabel(t *testing.T) {
	store := &SnapshotStore{}
	snap, _ := NewSnapshot("release-1", makeSnapshotReport())
	store.AddSnapshot(*snap)
	found := store.FindByLabel("release-1")
	if found == nil {
		t.Fatal("expected to find snapshot by label")
	}
	if found.Label != "release-1" {
		t.Errorf("unexpected label: %q", found.Label)
	}
}

func TestSnapshotStore_FindByLabel_Missing(t *testing.T) {
	store := &SnapshotStore{}
	if store.FindByLabel("nope") != nil {
		t.Error("expected nil for missing label")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	store := &SnapshotStore{}
	snap, _ := NewSnapshot("rt", makeSnapshotReport())
	store.AddSnapshot(*snap)

	var buf bytes.Buffer
	if err := SaveSnapshot(store, &buf); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := LoadSnapshot(&buf)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(loaded.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(loaded.Snapshots))
	}
	if loaded.Snapshots[0].Label != "rt" {
		t.Errorf("unexpected label: %q", loaded.Snapshots[0].Label)
	}
}
