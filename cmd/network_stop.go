package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var (
	stopNodesOnly bool
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

		updatedState, err := netStop(context.Background(), *netState)
		if err != nil {
			return fmt.Errorf("failed to stop network: %w", err)
		}

		return updatedState.Persist()
	},
}

func init() {
	netStopCmd.PersistentFlags().BoolVar(&stopNodesOnly,
		"nodes-only",
		false,
		"Only stops running nodes sets in the network.",
	)
}

func netStop(ctx context.Context, state state.NetworkState) (*state.NetworkState, error) {
	log.Println("Stopping network...")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize nomad client: %w", err)
	}

	var logsDir, vegaCapsuleBinary string
	if state.Config != nil {
		logsDir = state.Config.LogsDir()

		if state.Config.VegaCapsuleBinary != nil {
			vegaCapsuleBinary = *state.Config.VegaCapsuleBinary
		}
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, vegaCapsuleBinary, logsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	stoppedJobs, err := nomadRunner.StopNetwork(ctx, state.RunningJobs, stopNodesOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to stop network: %w", err)
	}

	if state.RunningJobs != nil {
		state.RunningJobs.RemoveRunningJobsIDs(stoppedJobs)
	}

	log.Println("Network successfully stopped.")
	return &state, nil
}
