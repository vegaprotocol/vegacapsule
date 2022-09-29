package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"

	vegapaths "code.vegaprotocol.io/vega/paths"
)

const defaultCheckpointWaitTimeout = 3 * time.Minute

var (
	checkpointsAmount       int
	printLastCheckpointPath bool
	checkpointWaitTimeout   time.Duration
)

var nodesWaitForCheckpoint = &cobra.Command{
	Use:   "wait-for-checkpoint",
	Short: "The command blocks execution until the given number of checkpoints is produced.",
	Long: `# todod add
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("nodes wait-for-checkpoint")
		}

		if !netState.Running() {
			return networkNotRunningErr("nodes wait-for-checkpoint")
		}
		ctx, cancel := context.WithTimeout(context.Background(), checkpointWaitTimeout)
		defer cancel()

		files, err := waitForCheckpointsInTheNetwork(ctx, netState.GeneratedServices.NodeSets.ToSlice(), checkpointsAmount)
		if err != nil {
			return fmt.Errorf("failed to wait for checkpoint: %w", err)
		}
		printOutput(files, printLastCheckpointPath)
		return nil
	},
}

type checkpointsResult struct {
	checkpointsFilesPaths []string
	err                   error
}

func waitForCheckpointsInTheNetwork(ctx context.Context, nodeSets []types.NodeSet, amount int) ([]string, error) {
	startTime := time.Now()

	result := make(chan checkpointsResult, 1)

	go func(startTime time.Time, nodeSets []types.NodeSet) {
		for {
			for _, ns := range nodeSets {
				vegaPaths := vegapaths.New(ns.Vega.HomeDir)
				checkpointDirectory := vegaPaths.StatePathFor(vegapaths.CheckpointStateHome)

				checkpointsFiles, err := findCheckpointsNewerThan(checkpointDirectory, startTime)

				if err != nil {
					result <- checkpointsResult{
						err: err,
					}
					return
				}

				if len(checkpointsFiles) >= amount {
					result <- checkpointsResult{
						checkpointsFilesPaths: checkpointsFiles,
					}
					return
				}
			}
		}
	}(startTime, nodeSets)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for checkpoint timed out")
	case result := <-result:
		return result.checkpointsFilesPaths, result.err
	}
}

func init() {
	nodesStopCmd.PersistentFlags().IntVar(&checkpointsAmount,
		"checkpoints",
		3,
		"Number for checkpoints to wait for",
	)

	nodesWaitForCheckpoint.PersistentFlags().BoolVar(&printLastCheckpointPath,
		"print-last-checkpoint-path-only",
		false,
		"Print only last checkpoint path",
	)

	nodesWaitForCheckpoint.PersistentFlags().DurationVar(&checkpointWaitTimeout,
		"timeout",
		defaultCheckpointWaitTimeout,
		"Timeout for the wait-for-checkpoint command",
	)
}

func isCheckpointFile(file os.DirEntry) bool {
	fInfo, err := file.Info()

	return err == nil && fInfo.Mode().IsRegular() && len(file.Name()) > 3 && file.Name()[len(file.Name())-3:] == ".cp"
}

func findCheckpointsNewerThan(directory string, startDate time.Time) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint directory: %w", err)
	}

	var checkpointFiles []string
	for _, fi := range files {
		fName := fi.Name()
		fileInfo, err := fi.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get info for the %s file: %w", fName, err)
		}

		if isCheckpointFile(fi) && fileInfo.ModTime().After(startDate) {
			checkpointFiles = append(checkpointFiles, filepath.Join(directory, fi.Name()))
		}
	}

	return checkpointFiles, nil
}

func printOutput(checkpointFiles []string, lastCheckpointOnly bool) error {
	if len(checkpointFiles) < 1 {
		return fmt.Errorf("no checkpoint found")
	}

	if lastCheckpointOnly {
		fmt.Println(checkpointFiles[len(checkpointFiles)-1])
		return nil
	}
	result := struct {
		Checkpoints []string `json:"checkpoints"`
	}{
		Checkpoints: checkpointFiles,
	}

	resultBytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to parse wait-for-checkpoint result: %w", err)
	}

	fmt.Println(string(resultBytes))
	return nil
}
