package cmd

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var (
	baseOneNode string
	startNode   bool
)

var nodesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed list validators: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes add")
		}

		newNodeSet, err := nodesAddNode(*networkState, baseOneNode)
		if err != nil {
			return fmt.Errorf("failed to add new node: %w", err)
		}

		networkState.GeneratedServices.NodeSets[newNodeSet.Name] = *newNodeSet

		if startNode {
			networkState, err = nodesStartNode(context.Background(), *networkState, newNodeSet.Name)
			if err != nil {
				return fmt.Errorf("failed start node: %w", err)
			}

		}

		return networkState.Persist()
	},
}

func init() {
	nodesAddCmd.PersistentFlags().BoolVar(&startNode,
		"start",
		true,
		"Allows to configure whether or not the new node set should automatically start",
	)
	nodesAddCmd.PersistentFlags().StringVar(&baseOneNode,
		"base-on",
		"",
		"Name of the node set that the new node set should be based on",
	)
	nodesAddCmd.MarkFlagRequired("base-on")
}

func nodesAddNode(state state.NetworkState, baseOneNode string) (*types.NodeSet, error) {
	gen, err := generator.New(state.Config, *state.GeneratedServices)
	if err != nil {
		return nil, err
	}

	nodeSet, err := state.GeneratedServices.GetNodeSet(baseOneNode)
	if err != nil {
		return nil, err
	}

	nodeConfig, err := state.Config.Network.GetNodeConfig(nodeSet.GroupName)
	if err != nil {
		return nil, err
	}

	newNodeSet, err := gen.AddNodeSet(len(state.GeneratedServices.NodeSets), *nodeConfig, *nodeSet, state.GeneratedServices.Faucet)
	if err != nil {
		return nil, err
	}

	return newNodeSet, nil
}
