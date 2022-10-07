package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	searchFromBeginning     bool
)

var nodesWaitForCheckpoint = &cobra.Command{
	Use:   "wait-for-checkpoint",
	Short: "The command blocks execution until the given number of checkpoints is produced.",
	Long: `The commands can wait for a checkpoint from the moment of calling it. There is also the
"--search-from-beginning" flag that tells the command to wait until the network produces a 
given number of checkpoints from the beginning.


By default, the command prints a list of paths to the checkpoints that have been found. However, 
you can use the "--print-last-checkpoint-path-only" flag to print the last found checkpoint path only.`,
	Example: `
# Wait for 6 checkpoint up to 8 min
vegacapsule nodes wait-for-checkpoint --checkpoints 6 --timeout 8m

# Wait for network to produce very first checkpoint and print its path to the stdout
vegacapsule nodes wait-for-checkpoint --checkpoints 1 --timeout 8m --print-last-checkpoint-path-only --search-from-beginning`,
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

		startTime := time.Now()
		if searchFromBeginning {
			startTime = time.Time{}
		}

		files, err := waitForCheckpointsInTheNetwork(ctx, netState.GeneratedServices.NodeSets.ToSlice(), checkpointsAmount, startTime)
		if err != nil {
			return fmt.Errorf("failed to wait for checkpoint: %w", err)
		}

		return printOutput(files, printLastCheckpointPath)
	},
}

func init() {
	nodesWaitForCheckpoint.PersistentFlags().IntVar(&checkpointsAmount,
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

	nodesWaitForCheckpoint.PersistentFlags().BoolVar(&searchFromBeginning,
		"search-from-beginning",
		false,
		"It's searching for the N checkpoints from the beginning of the networks",
	)
}

type checkpointsResult struct {
	checkpointsFilesPaths []string
	err                   error
}

func waitForCheckpointsInTheNetwork(ctx context.Context, nodeSets []types.NodeSet, amount int, startTime time.Time) ([]string, error) {
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
			time.Sleep(4 * time.Second)
		}
	}(startTime, nodeSets)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for checkpoint timed out")
	case result := <-result:
		return result.checkpointsFilesPaths, result.err
	}
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
	sort.Slice(files, func(i, j int) bool {
		iInfo, iErr := files[i].Info()
		jInfo, jErr := files[j].Info()

		if jErr != nil || iErr != nil {
			return false
		}

		return iInfo.ModTime().After(jInfo.ModTime())
	})

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
