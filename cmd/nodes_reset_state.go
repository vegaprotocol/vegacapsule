package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var remoteResetState bool

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

		var cmdResult io.Reader

		if !remoteResetState {
			cmdResult, err = commands.ResetNodeSetsData(
				*netState.Config.VegaBinary,
				netState.GeneratedServices.NodeSets.ToSlice(),
			)
		} else {
			var nomadClient *nomad.Client
			nomadClient, err = nomad.NewClient(nil)
			if err != nil {
				return fmt.Errorf("failed to create nomad client: %w", err)
			}

			commandRunner := nomad.NewCommandRunner(nomadClient)
			cmdResult, err = commandRunner.NetworkUnsafeResetAll(context.Background(), netState.GeneratedServices.NodeSets.ToSlice())
		}

		if err != nil {
			return fmt.Errorf("failed to reset node sets: %w", err)
		}

		if _, err := io.Copy(os.Stdout, cmdResult); err != nil {
			return fmt.Errorf("failed to write command output: %w", err)
		}

		return nil
	},
}

func init() {
	nodesUnsafeResetAllCmd.PersistentFlags().BoolVar(&remoteResetState,
		"remote",
		false,
		"Determines, whether the command should be executed locally or remotelly",
	)
}
