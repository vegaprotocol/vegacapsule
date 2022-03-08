package cmd

import (
	nmrunner "code.vegaprotocol.io/vegacapsule/nomad/runner"
	"github.com/spf13/cobra"
)

var nomadConfigPath string

var nomadCmd = &cobra.Command{
	Use:   "nomad",
	Short: "Starts Nomad instance locally",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nmrunner.StartAgent(nomadConfigPath)
	},
}

func init() {
	nomadCmd.PersistentFlags().StringVar(&nomadConfigPath,
		"nomad-config-path",
		"",
		"Allows to use Nomad configuration",
	)
}
