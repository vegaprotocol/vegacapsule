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
	nodeName  string
	nodeNames *[]string
)

var nodesStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start running node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes start")
		}

		updatedNetworkState, err := nodesStartNode(context.Background(), *networkState, *nodeNames)
		if err != nil {
			return fmt.Errorf("failed start node: %w", err)
		}

		return updatedNetworkState.Persist()
	},
}

func init() {
	nodeNames = nodesStartCmd.PersistentFlags().StringSlice(
		"name",
		[]string{},
		"Name of the node tha should be started",
	)
	nodesStartCmd.MarkFlagRequired("name")
}

func nodesStartNode(ctx context.Context, state state.NetworkState, names []string) (*state.NetworkState, error) {
	nodeSets := []types.NodeSet{}

	for _, n := range names {
		log.Printf("starting %s node set", n)
		nodeSet, err := state.GeneratedServices.GetNodeSet(n)
		if err != nil {
			return nil, err
		}
		nodeSets = append(nodeSets, *nodeSet)
	}

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	for _, nodeSet := range nodeSets {
		if _, err := nomadRunner.RunRawNomadJobs(ctx, nodeSet.PreGenerateRawJobs()); err != nil {
			return nil, fmt.Errorf("failed to start node set %q pre generate jobs: %w", nodeSet.Name, err)
		}
	}

	res, err := nomadRunner.RunNodeSets(ctx, nodeSets)
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad node set: %w", err)
	}

	state.RunningJobs.NodesSetsJobIDs[*res[0].ID] = true

	log.Printf("starting nodes set success")
	return &state, nil
}
