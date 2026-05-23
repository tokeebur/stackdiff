package main

import (
	"fmt"
	"os"

	"github.com/stackdiff/internal/diff"
	"github.com/stackdiff/internal/state"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: stackdiff <state-a> <state-b>")
	}

	stateA, err := state.ParseStateFile(args[0])
	if err != nil {
		return fmt.Errorf("parsing state A (%s): %w", args[0], err)
	}

	stateB, err := state.ParseStateFile(args[1])
	if err != nil {
		return fmt.Errorf("parsing state B (%s): %w", args[1], err)
	}

	result := diff.Compare(stateA, stateB)

	if result.IsEmpty() {
		fmt.Println("No drift detected. States are identical.")
		return nil
	}

	diff.FormatSummary(os.Stdout, result)
	return nil
}
