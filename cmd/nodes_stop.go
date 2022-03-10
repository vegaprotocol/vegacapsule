package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

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

		updatedNetworkState, err := nodesStopNode(context.Background(), *networkState, nodeName)
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
	nodesStopCmd.MarkFlagRequired("name")
}

func nodesStopNode(ctx context.Context, state state.NetworkState, name string) (*state.NetworkState, error) {
	log.Printf("stopping %s node set", name)

	_, err := state.GeneratedServices.GetNodeSet(name)
	if err != nil {
		return nil, err
	}

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	if err := nomadRunner.StopJobs(ctx, []string{name}); err != nil {
		return nil, fmt.Errorf("failed to stop nomad job %q: %w", name, err)
	}

	delete(state.RunningJobs.NodesSetsJobIDs, name)

	log.Printf("stopping %s node set success", name)
	return &state, nil
}
