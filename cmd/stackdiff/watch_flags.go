package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/stackdiff/internal/diff"
)

// WatchFlags holds parsed watch-mode options.
type WatchFlags struct {
	Enabled         bool
	IntervalSeconds int
	MaxIterations   int
	Quiet           bool
	ExitOnDrift     bool
}

// ParseWatchFlags reads watch-related flags from the provided FlagSet.
func ParseWatchFlags(fs *flag.FlagSet) *WatchFlags {
	wf := &WatchFlags{}
	fs.BoolVar(&wf.Enabled, "watch", false, "enable watch mode: re-compare states on a fixed interval")
	fs.IntVar(&wf.IntervalSeconds, "watch-interval", 30, "seconds between each comparison in watch mode")
	fs.IntVar(&wf.MaxIterations, "watch-max", 0, "maximum number of watch iterations (0 = unlimited)")
	fs.BoolVar(&wf.Quiet, "watch-quiet", false, "suppress per-iteration status lines in watch mode")
	fs.BoolVar(&wf.ExitOnDrift, "watch-exit-on-drift", false, "exit with code 1 as soon as drift is detected")
	return wf
}

// RunWatch executes watch mode using the supplied paths and flags, writing
// output to out. Returns a non-zero exit code when ExitOnDrift triggers.
func RunWatch(oldPath, newPath string, wf *WatchFlags, out io.Writer) int {
	if wf == nil || !wf.Enabled {
		return 0
	}

	cfg := diff.WatchConfig{
		IntervalSeconds: wf.IntervalSeconds,
		MaxIterations:   wf.MaxIterations,
		Quiet:           wf.Quiet,
	}

	exitCode := 0
	err := diff.Watch(oldPath, newPath, cfg, out, func(r diff.WatchResult) bool {
		if r.HasDrift && wf.ExitOnDrift {
			fmt.Fprintln(out, "drift detected — exiting")
			exitCode = 1
			return true // stop
		}
		return false
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "watch error: %v\n", err)
		return 1
	}
	return exitCode
}
