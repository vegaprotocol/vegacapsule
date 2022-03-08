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

var nodeName string

var nodesStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start running node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed list validators: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes start")
		}

		updatedNetworkState, err := nodesStartNode(context.Background(), *networkState, nodeName)
		if err != nil {
			return fmt.Errorf("failed start node: %w", err)
		}

		return updatedNetworkState.Persist()
	},
}

func init() {
	nodesStartCmd.PersistentFlags().StringVar(&nodeName,
		"name",
		"",
		"Name of the node tha should be started",
	)
	nodesStartCmd.MarkFlagRequired("name")
}

func nodesStartNode(ctx context.Context, state state.NetworkState, name string) (*state.NetworkState, error) {
	log.Printf("starting %s node set", name)

	nodeSet, err := state.GeneratedServices.GetNodeSet(name)
	if err != nil {
		return nil, err
	}

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	res, err := nomadRunner.RunNodeSets(ctx, *state.Config.VegaBinary, []types.NodeSet{*nodeSet})
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad network: %s", err)
	}

	state.RunningJobs.NodesSetsJobIDs[*res[0].ID] = true

	log.Printf("starting %s node set success", name)
	return &state, nil
}
