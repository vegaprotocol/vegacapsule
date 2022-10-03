package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/spf13/cobra"
)

var (
	nodeName   string
	vegaBinary string
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

		updatedNetworkState, err := nodesStartNode(context.Background(), *networkState, nodeName, vegaBinary)
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
	nodesStartCmd.PersistentFlags().StringVar(&vegaBinary,
		"vega-binary",
		"",
		"Path of Vega binary to be used to start the node",
	)
	nodesStartCmd.MarkFlagRequired("name")
}

func nodesStartNode(ctx context.Context, state state.NetworkState, name, vegaBinary string) (*state.NetworkState, error) {
	log.Printf("starting %s node set", name)

	nodeSet, err := state.GeneratedServices.GetNodeSet(name)
	if err != nil {
		return nil, err
	}

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	if _, err := nomadRunner.RunRawNomadJobs(ctx, nodeSet.PreGenerateRawJobs()); err != nil {
		return nil, fmt.Errorf("failed to start node set %q pre generate jobs: %w", nodeSet.Name, err)
	}

	if vegaBinary != "" {
		vegaBinPath, err := utils.BinaryAbsPath(vegaBinary)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for %q: %w", vegaBinary, err)
		}

		nodeSet.Vega.BinaryPath = vegaBinPath

		if nodeSet.DataNode != nil {
			nodeSet.DataNode.BinaryPath = vegaBinPath
		}
	}

	res, err := nomadRunner.RunNodeSets(ctx, []types.NodeSet{*nodeSet})
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad node set %q : %w", nodeSet.Name, err)
	}

	state.RunningJobs.NodesSetsJobIDs[*res[0].ID] = true
	state.GeneratedServices.NodeSets[nodeSet.Name] = *nodeSet

	log.Printf("starting %s node set success", name)
	return &state, nil
}
