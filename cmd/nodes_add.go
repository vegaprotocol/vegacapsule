package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	baseOnNode     string
	baseOnGroup    string
	startNode      bool
	resultsOutPath string
	count          int
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
				newNodeSet, err := nodesAddNode(*networkState, i, baseOnNode, baseOnGroup)
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

		outputStringJSON, err := json.MarshalIndent(newNodeSets, "", "\t")
		if err != nil {
			return fmt.Errorf("failed to marshal validators info: %w", err)
		}

		fmt.Println(string(outputStringJSON))
		if resultsOutPath != "" {
			if err := os.WriteFile(resultsOutPath, outputStringJSON, 0o666); err != nil {
				return fmt.Errorf("failed to save results about new nodes into the file: %w", err)
			}
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
	nodesAddCmd.PersistentFlags().StringVar(&baseOnNode,
		"base-on",
		"",
		"Name of the node set that the new node set should be based on",
	)
	nodesAddCmd.PersistentFlags().StringVar(&baseOnGroup,
		"base-on-group",
		"",
		"Name of the group that the new node set should be based on",
	)
	nodesAddCmd.PersistentFlags().IntVar(&count,
		"count",
		1,
		"Defines how many node sets should be added",
	)

	nodesAddCmd.PersistentFlags().StringVar(&resultsOutPath,
		"out-path",
		"",
		"If not empty, details about added nodes are saved in the given file",
	)
}

func nodesAddNode(state state.NetworkState, index int, baseOnNode, baseOnGroup string) (*types.NodeSet, error) {
	if baseOnNode != "" && baseOnGroup != "" {
		return nil, fmt.Errorf("provide either value for --base-on or --base-on-group, not both values")
	}

	if baseOnNode == "" && baseOnGroup == "" {
		return nil, fmt.Errorf("value for --base-on-node or --base-on-group must be provided")
	}

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	capsuleBinary := ""
	if state.Config.VegaCapsuleBinary != nil {
		capsuleBinary = *state.Config.VegaCapsuleBinary
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, capsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	gen, err := generator.New(state.Config, *state.GeneratedServices, nomadRunner, state.VegaChainID)
	if err != nil {
		return nil, err
	}

	var (
		nodeSet    *types.NodeSet
		groupName  string
		groupIndex int = -1
	)
	if baseOnNode != "" {
		nodeSet, err = state.GeneratedServices.GetNodeSet(baseOnNode)
		if err != nil {
			return nil, fmt.Errorf("failed to get node set by name: %w", err)
		}

		groupName = nodeSet.GroupName
		groupIndex = nodeSet.GroupIndex
	} else {
		for groupIdx, group := range state.Config.Network.Nodes {
			if group.Name == baseOnGroup {
				groupIndex = groupIdx
				break
			}
		}

		if groupIndex < 0 {
			return nil, fmt.Errorf("the %s nodes group not found", baseOnGroup)
		}

		nodes := state.GeneratedServices.GetNodeSetsByGroupName(baseOnGroup)
		if len(nodes) < 1 {
			// Nodes within given group does not exists, fallback to the first available node
			nodes = state.GeneratedServices.NodeSets.ToSlice()
		}

		nodeSet = &(nodes[0])
		groupName = baseOnGroup
	}

	if err != nil {
		return nil, err
	}

	nodeConfig, err := state.Config.Network.GetNodeConfig(groupName)
	if err != nil {
		return nil, err
	}

	groupNodeSets := state.GeneratedServices.GetNodeSetsByGroupName(groupName)

	newNodeSet, err := gen.AddNodeSet(
		len(state.GeneratedServices.NodeSets)-1+index,
		len(groupNodeSets),
		groupIndex,
		*nodeConfig,
		*nodeSet,
		state.GeneratedServices.Faucet,
	)
	if err != nil {
		return nil, err
	}

	return newNodeSet, nil
}
