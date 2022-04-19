package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var (
	checkpointFile string
)

var nodesRestoreCheckpointCmd = &cobra.Command{
	Use:   "restore-checkpoint",
	Short: "Restore all Vega nodes state from checkpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("nodes restore-checkpoint")
		}

		if !netState.Running() {
			return networkNotRunningErr("nodes restore-checkpoint")
		}

		if checkpointFile == "" {
			return fmt.Errorf("parameter checkpoint-file can not be empty")
		}

		for _, ns := range netState.GeneratedServices.NodeSets {
			r, err := commands.VegaRestoreCheckpoint(
				*netState.Config.VegaBinary,
				ns.Tendermint.HomeDir,
				checkpointFile,
				ns.Vega.NodeWalletPassFilePath,
			)
			if err != nil {
				return fmt.Errorf("failed to restore node %q from checkpoint: %w", ns.Name, err)
			}

			fmt.Printf("applied transaction for node set %q: %s", ns.Name, r)
		}

		return nil
	},
}

func init() {
	nodesRestoreCheckpointCmd.PersistentFlags().StringVar(&checkpointFile,
		"checkpoint-file",
		"",
		"Path to the checkpoint file",
	)
	nodesRestoreCheckpointCmd.MarkFlagRequired("checkpoint-file")
}
