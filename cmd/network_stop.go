package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
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

func netStop(ctx context.Context, state *state.NetworkState) error {
	log.Println("stopping network")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	if err := nomadRunner.StopNetwork(ctx, state.RunningJobs); err != nil {
		return fmt.Errorf("failed to stop nomad network: %w", err)
	}

	log.Println("stopping network success")
	return nil
}
