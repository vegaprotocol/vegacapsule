package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	baseOneNode string
	startNode   bool
	count       int
)

var nodesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed to load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes add")
		}

		if count < 1 {
			return fmt.Errorf("count has to be > 0")
		}

		var eg errgroup.Group
		var m sync.Mutex
		newNodeSets := make([]*types.NodeSet, 0, count)

		for i := 0; i < count; i++ {
			i := i + 1
			eg.Go(func() error {
				newNodeSet, err := nodesAddNode(*networkState, i, baseOneNode)
				if err != nil {
					return fmt.Errorf("failed to add new node: %w", err)
				}

				m.Lock()
				newNodeSets = append(newNodeSets, newNodeSet)
				m.Unlock()

				if startNode {
					nomadJobID, err := nodesStartNode(cmd.Context(), newNodeSet, networkState.Config, newNodeSet.Vega.BinaryPath)
					if err != nil {
						return fmt.Errorf("failed start node: %w", err)
					}

					m.Lock()
					networkState.RunningJobs.NodesSetsJobIDs[nomadJobID] = true
					m.Unlock()
				}

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		for _, ns := range newNodeSets {
			networkState.GeneratedServices.NodeSets[ns.Name] = *ns
		}

		if err := networkState.Persist(); err != nil {
			return fmt.Errorf("failed to persist network: %w", err)
		}

		for _, ns := range newNodeSets {
			newNodeJson, err := json.MarshalIndent(ns, "", "\t")
			if err != nil {
				return fmt.Errorf("failed to marshal validators info: %w", err)
			}

			fmt.Println(string(newNodeJson))
		}

		return nil
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

	nodesAddCmd.PersistentFlags().IntVar(&count,
		"count",
		1,
		"Defines how many node sets should be added",
	)
}

func nodesAddNode(state state.NetworkState, index int, baseOneNode string) (*types.NodeSet, error) {
	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	gen, err := generator.New(state.Config, *state.GeneratedServices, nomadRunner, state.VegaChainID)
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

	groupNodeSets := state.GeneratedServices.GetNodeSetsByGroupName(nodeSet.GroupName)

	newNodeSet, err := gen.AddNodeSet(
		len(state.GeneratedServices.NodeSets)-1+index,
		len(groupNodeSets),
		nodeSet.GroupIndex,
		*nodeConfig,
		*nodeSet,
		state.GeneratedServices.Faucet,
	)
	if err != nil {
		return nil, err
	}

	return newNodeSet, nil
}
