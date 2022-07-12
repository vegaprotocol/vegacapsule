package cmd

import (
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manages network keys",
	Long:  "The command allows to manage network keys",
}

func init() {
	keysCmd.AddCommand(netKeysImportCmd)
}
