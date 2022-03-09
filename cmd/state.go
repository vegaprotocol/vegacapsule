package cmd

import (
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Manages vegacapsule state",
}

func init() {
	stateCmd.AddCommand(stateGetAddresses)
}
