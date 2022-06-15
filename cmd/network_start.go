package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

const defaultTimeout = 300

var (
	timeout uint64
)

var netStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts existing network",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("start")
		}

		saveState := func(updatedNetState *state.NetworkState) {
			if updatedNetState == nil {
				return
			}

			log.Printf("saving network state to the file")
			saveErr := updatedNetState.Persist()
			if saveErr != nil {
				log.Printf("failed to save network state: %s", err)
			}

			// do not shadow the original error as it is more important
			if err == nil {
				err = saveErr
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		updatedNetState, err := netStart(ctx, *netState)
		// We want state saved even if the network is not started properly
		defer saveState(updatedNetState)
		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return err
	},
}

func init() {
	netStartCmd.PersistentFlags().Uint64Var(&timeout,
		"timeout",
		defaultTimeout,
		"Timeout in seconds",
	)
}

func netStart(ctx context.Context, state state.NetworkState) (*state.NetworkState, error) {
	log.Println("starting network")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	res, err := nomadRunner.StartNetwork(ctx, state.Config, state.GeneratedServices)
	// if network state returned, save it in current state
	if res != nil {
		state.RunningJobs = res
	}

	if err != nil {
		return nil, fmt.Errorf("failed to start nomad network: %s", err)
	}

	log.Println("starting network success")
	return &state, nil
}
