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

var (
	stopNodesOnly      bool
	stopWithCmdRunners bool
)

var netStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop existing network",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("stop")
		}

		if err := netStop(context.Background(), netState); err != nil {
			return fmt.Errorf("failed to stop network: %w", err)
		}

		return nil
	},
}

func init() {
	netStopCmd.PersistentFlags().BoolVar(&stopNodesOnly,
		"nodes-only",
		false,
		"Stops all nodes running in the network.",
	)
	netStopCmd.PersistentFlags().BoolVar(&stopWithCmdRunners,
		"with-command-runners",
		false,
		"If this flag is passed command-runner jobs are also stopped",
	)
}

func netStop(ctx context.Context, state *state.NetworkState) error {
	log.Println("stopping network")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return fmt.Errorf("failed to create nomad client: %w", err)
	}

	jobFilters := []types.NetworkJobsFilter{}

	if stopNodesOnly {
		jobFilters = append(jobFilters, types.FilterNetworkJobsByJobKindIn([]types.JobKind{types.JobNodeSet}))
	}

	if !stopWithCmdRunners {
		jobFilters = append(jobFilters, types.FilterNetworkJobsByJobKindNotIn([]types.JobKind{types.JobCommandRunner}))
	}

	var jobs []types.NetworkJobState
	if state.RunningJobs != nil {
		jobs = state.RunningJobs.Filter(jobFilters).ToSlice()
	}

	if len(jobs) == 0 && stopNodesOnly {
		log.Println("All nodes are already stopped")
		return nil
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)
	if err := nomadRunner.StopNetwork(ctx, jobs); err != nil {
		return fmt.Errorf("failed to stop nomad network: %w", err)
	}

	log.Println("stopping network success")
	return nil
}
