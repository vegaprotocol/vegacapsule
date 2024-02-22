package cmd

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"

	"github.com/spf13/cobra"
)

var datanodeBackupDir string

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

		updatedNetworkState, err := nodesStopNode(context.Background(), *networkState, nodeName, true)
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
	nodesRemoveCmd.PersistentFlags().StringVar(&datanodeBackupDir, "datanode-backup-dir", "", "Directory where data node home directories should be backed up before removal")
}

func nodesRemoveNode(state state.NetworkState, name string) (*state.NetworkState, error) {
	gen, err := generator.New(state.Config, *state.GeneratedServices, nomad.NewVoidJobRunner(), state.VegaChainID)
	if err != nil {
		return nil, err
	}

	nodeSet, err := state.GeneratedServices.GetNodeSet(name)
	if err != nil {
		return nil, err
	}

	if len(datanodeBackupDir) > 0 {
		if err := gen.BackupDataNode(*nodeSet, datanodeBackupDir); err != nil {
			return nil, fmt.Errorf("failed to backup data node: %w", err)
		}
	}

	if err := gen.RemoveNodeSet(*nodeSet); err != nil {
		return nil, err
	}

	delete(state.GeneratedServices.NodeSets, nodeSet.Name)

	return &state, nil
}
