package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// SaveSnapshot serialises a SnapshotStore to the given writer.
func SaveSnapshot(store *SnapshotStore, w io.Writer) error {
	if store == nil {
		return fmt.Errorf("snapshot: store must not be nil")
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(store)
}

// LoadSnapshot deserialises a SnapshotStore from the given reader.
func LoadSnapshot(r io.Reader) (*SnapshotStore, error) {
	var store SnapshotStore
	if err := json.NewDecoder(r).Decode(&store); err != nil {
		return nil, fmt.Errorf("snapshot: failed to decode store: %w", err)
	}
	return &store, nil
}

// SaveSnapshotFile writes a SnapshotStore to a file path.
func SaveSnapshotFile(store *SnapshotStore, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: cannot create file %q: %w", path, err)
	}
	defer f.Close()
	return SaveSnapshot(store, f)
}

// LoadSnapshotFile reads a SnapshotStore from a file path.
func LoadSnapshotFile(path string) (*SnapshotStore, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: cannot open file %q: %w", path, err)
	}
	defer f.Close()
	return LoadSnapshot(f)
}
