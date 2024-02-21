package cmd

import (
	"fmt"
	"io"
	"os"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/state"

	"github.com/spf13/cobra"
)

var nodesUnsafeResetAllCmd = &cobra.Command{
	Use:   "unsafe-reset-all",
	Short: "(unsafe) Reset all nodes (Vega and Tendermint) state (checkpoints, snapshots etc..)",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state unsafe-reset-all")
		}

		r, err := commands.ResetNodeSetsData(
			*netState.Config.VegaBinary,
			netState.GeneratedServices.NodeSets.ToSlice(),
		)
		if err != nil {
			return fmt.Errorf("failed to reset node sets: %w", err)
		}

		if _, err := io.Copy(os.Stdout, r); err != nil {
			return fmt.Errorf("failed to write command output: %w", err)
		}

		return nil
	},
}
