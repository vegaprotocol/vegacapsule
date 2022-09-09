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
	doNotStopOnFailure bool
)

var netStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts existing network",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("start")
		}

		updatedNetState, err := netStart(context.Background(), *netState)
		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return updatedNetState.Persist()
	},
}

func init() {
	netStartCmd.PersistentFlags().BoolVar(&doNotStopOnFailure,
		"do-not-stop-on-failure",
		false,
		"Do not stop partially running network when failed to start",
	)
}

func netStart(ctx context.Context, state state.NetworkState) (*state.NetworkState, error) {
	log.Println("starting network")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	res, err := nomadRunner.StartNetwork(ctx, state.Config, state.GeneratedServices, !doNotStopOnFailure)
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad network: %s", err)
	}
	state.RunningJobs = res

	log.Println("starting network success")

	if err := printNetworkAddresses(ctx, nomadRunner, state.GeneratedServices); err != nil {
		log.Printf("failed to print network addresses - please try to run 'network print-ports' instaed: %s", err)
	}

	return &state, nil
}
