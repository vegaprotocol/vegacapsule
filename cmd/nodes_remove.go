package cmd

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var nodesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove existing node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed to load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes remove")
		}

		updatedNetworkState, err := nodesStopNode(context.Background(), *networkState, nodeName)
		if err != nil {
			return fmt.Errorf("failed stop node: %w", err)
		}

		updatedNetworkState, err = nodesRemoveNode(*updatedNetworkState, nodeName)
		if err != nil {
			return fmt.Errorf("failed remove node: %w", err)
		}

		return updatedNetworkState.Persist()
	},
}

func init() {
	nodesRemoveCmd.PersistentFlags().StringVar(&nodeName,
		"name",
		"",
		"Name of the node tha should be removed",
	)
	nodesRemoveCmd.MarkFlagRequired("name")
}

func nodesRemoveNode(state state.NetworkState, name string) (*state.NetworkState, error) {
	gen, err := generator.New(state.Config, *state.GeneratedServices)
	if err != nil {
		return nil, err
	}

	nodeSet, err := state.GeneratedServices.GetNodeSet(name)
	if err != nil {
		return nil, err
	}

	if err := gen.RemoveNodeSet(*nodeSet); err != nil {
		return nil, err
	}

	delete(state.GeneratedServices.NodeSets, nodeSet.Name)

	return &state, nil
}
