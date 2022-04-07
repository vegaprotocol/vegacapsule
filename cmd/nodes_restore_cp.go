package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
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

		ns := selectNodeSetForCheckpointRestore(netState)
		if ns == nil {
			return fmt.Errorf("no running node set found")
		}

		r, err := commands.VegaRestoreCheckpoint(
			*netState.Config.VegaBinary,
			ns.Vega.HomeDir,
			checkpointFile,
			ns.Vega.NodeWalletPassFilePath,
		)
		if err != nil {
			return fmt.Errorf("failed to restart from checkpoint node sets: %w", err)
		}

		fmt.Printf("applied transaction for node set %q: %s", ns.Name, r)

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

func selectNodeSetForCheckpointRestore(state *state.NetworkState) *types.NodeSet {
	for _, jobID := range state.RunningJobs.NodesSetsJobIDs.ToSlice() {
		nodeSet, _ := state.GeneratedServices.GetNodeSet(jobID)
		if nodeSet != nil && nodeSet.Mode == types.NodeModeValidator {
			return nodeSet
		}
	}

	return nil
}
