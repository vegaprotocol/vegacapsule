package cmd

import (
	"encoding/json"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var nodesLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all node sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("ls")
		}

		nodeSets := networkState.GeneratedServices.NodeSets

		nodeSetsJson, err := json.MarshalIndent(nodeSets, "", "\t")
		if err != nil {
			return fmt.Errorf("failed to marshal validators info: %w", err)
		}

		fmt.Println(string(nodeSetsJson))
		return nil
	},
}
