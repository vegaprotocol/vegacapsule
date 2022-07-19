package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var stopWithPreGenenerate bool

var nodesStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop running node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes stop")
		}

		nodeSet, err := networkState.GeneratedServices.GetNodeSet(nodeName)
		if err != nil {
			return err
		}

		wantedNamesToStop := []string{nodeSet.Name}
		if stopWithPreGenenerate {
			wantedNamesToStop = append(wantedNamesToStop, nodeSet.PreGenerateJobsIDs()...)
		}
		if stopWithCmdRunners {
			wantedNamesToStop = append(wantedNamesToStop, nodeSet.RemoteCommandRunner.Name)
		}

		jobs := networkState.RunningJobs.GetByNames(wantedNamesToStop)

		updatedNetworkState, err := nodesStopNode(context.Background(), *networkState, nodeName, jobs, false)
		if err != nil {
			return fmt.Errorf("failed stop node: %w", err)
		}

		return updatedNetworkState.Persist()
	},
}

func init() {
	nodesStopCmd.PersistentFlags().StringVar(&nodeName,
		"name",
		"",
		"Name of the node tha should be stopped",
	)
	nodesStopCmd.PersistentFlags().BoolVar(&stopWithPreGenenerate,
		"with-pre-generate",
		true,
		"Whether or not the pre-generate jobs should be also stopped",
	)
	nodesStopCmd.PersistentFlags().BoolVar(&stopWithCmdRunners,
		"with-command-runners",
		false,
		"Whether or not the command-runner job should be also stopped",
	)
	nodesStopCmd.MarkFlagRequired("name")
}

func nodesStopNode(ctx context.Context, state state.NetworkState, name string, jobs types.JobStateMap, ignoreIfStopped bool) (*state.NetworkState, error) {
	log.Printf("stopping the %s node set", name)

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	if len(jobs.ToSlice()) == 0 {
		if ignoreIfStopped {
			return &state, nil
		}

		return nil, fmt.Errorf("given node set is not running")
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	if err := nomadRunner.StopJobs(ctx, jobs.ToSliceNames()); err != nil {
		return nil, fmt.Errorf("failed to stop nomad job %q: %w", name, err)
	}

	for _, job := range jobs {
		state.RunningJobs.RemoveJobs([]types.NetworkJobState{job})
	}

	log.Printf("stopping %s node set success", name)
	return &state, nil
}
