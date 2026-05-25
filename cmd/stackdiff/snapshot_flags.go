package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/you/stackdiff/internal/diff"
)

// SnapshotConfig holds CLI options for snapshot operations.
type SnapshotConfig struct {
	SavePath  string
	LoadPath  string
	Label     string
}

// ParseSnapshotFlags registers and parses snapshot-related CLI flags.
func ParseSnapshotFlags(fs *flag.FlagSet) *SnapshotConfig {
	cfg := &SnapshotConfig{}
	fs.StringVar(&cfg.SavePath, "snapshot-save", "", "path to save snapshot store (JSON)")
	fs.StringVar(&cfg.LoadPath, "snapshot-load", "", "path to load existing snapshot store for append")
	fs.StringVar(&cfg.Label, "snapshot-label", "", "label to attach to the new snapshot")
	return cfg
}

// ApplySnapshotSave appends the current report as a snapshot and writes it to disk.
func ApplySnapshotSave(cfg *SnapshotConfig, r *diff.Report) error {
	if cfg == nil || cfg.SavePath == "" {
		return nil
	}

	var store *diff.SnapshotStore

	if cfg.LoadPath != "" {
		var err error
		store, err = diff.LoadSnapshotFile(cfg.LoadPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load snapshot store %q: %v\n", cfg.LoadPath, err)
			store = &diff.SnapshotStore{}
		}
	} else {
		store = &diff.SnapshotStore{}
	}

	snap, err := diff.NewSnapshot(cfg.Label, r)
	if err != nil {
		return fmt.Errorf("snapshot: %w", err)
	}
	store.AddSnapshot(*snap)

	if err := diff.SaveSnapshotFile(store, cfg.SavePath); err != nil {
		return fmt.Errorf("snapshot: failed to save: %w", err)
	}
	return nil
}
