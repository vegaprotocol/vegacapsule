package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const defaultCheckpointWaitTimeout = 3 * time.Minute

var (
	checkpointsAmount       int
	printLastCheckpointPath bool
	printListOfCheckpoints  bool
	checkpointWaitTimeout   time.Duration
)

var nodesWaitForCheckpoint = &cobra.Command{
	Use:   "wait-for-checkpoint",
	Short: "The command blocks until the given number of checkpoints are produced.",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

func init() {
	nodesStopCmd.PersistentFlags().IntVar(&checkpointsAmount,
		"checkpoints",
		3,
		"Number for checkpoints to wait for",
	)

	nodesStopCmd.PersistentFlags().BoolVar(&printLastCheckpointPath,
		"print-last-checkpoint-path",
		false,
		"Print only last checkpoint path",
	)

	nodesStopCmd.PersistentFlags().BoolVar(&printListOfCheckpoints,
		"print-list-of-checkpoints",
		false,
		"Print list of checkpoints as a output",
	)

	nodesStopCmd.PersistentFlags().DurationVar(&checkpointWaitTimeout,
		"checkpoints",
		defaultCheckpointWaitTimeout,
		"Timeout for the wait-for-checkpoint command",
	)
}

func findCheckpointsNewerThan(directory string, startDate time.Time) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint directory: %w", err)
	}

	var names []string
	for _, fi := range files {
		fName := fi.Name()
		fileInfo, err := fi.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get info for the %s file: %w", fName, err)
		}

		if fileInfo.Mode().IsRegular() && len(fName) > 3 && fName[len(fName)-3:] == ".cp" && fileInfo.ModTime().After(startDate) {
			names = append(names, fi.Name())
		}
	}

	return names, nil
}
