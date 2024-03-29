package cmd

import "github.com/spf13/cobra"

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manages nodes sets",
	Long:  "The command allows some operations on nodes sets like adding/removing/listing validators etc.",
}

func init() {
	nodesCmd.AddCommand(nodesLsValidatorsCmd)
	nodesCmd.AddCommand(nodesLsCmd)
	nodesCmd.AddCommand(nodesAddCmd)
	nodesCmd.AddCommand(nodesStartCmd)
	nodesCmd.AddCommand(nodesStopCmd)
	nodesCmd.AddCommand(nodesRemoveCmd)
	nodesCmd.AddCommand(nodesUnsafeResetAllCmd)
	nodesCmd.AddCommand(nodesRestoreCheckpointCmd)
	nodesCmd.AddCommand(nodesProtocolUpgradeCmd)
}
