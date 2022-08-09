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
}

func netStop(ctx context.Context, state *state.NetworkState) error {
	log.Println("stopping network")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return fmt.Errorf("failed to create nomad client: %w", err)
	}

	var logsDir string
	if state.Config != nil {
		logsDir = state.Config.LogsDir()
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, logsDir)
	if err != nil {
		return fmt.Errorf("failed to create job runner: %w", err)
	}

	if err := nomadRunner.StopNetwork(ctx, state.RunningJobs, stopNodesOnly); err != nil {
		return fmt.Errorf("failed to stop nomad network: %w", err)
	}

	log.Println("stopping network success")
	return nil
}
