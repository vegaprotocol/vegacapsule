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
		installPath, err := getInstallPath(installPath)
		if err != nil {
			return err
		}

		return nmrunner.StartAgent(nomadConfigPath, installPath)
	},
}

func init() {
	nomadCmd.PersistentFlags().StringVar(&nomadConfigPath,
		"nomad-config-path",
		"",
		"Allows to use Nomad configuration",
	)
	nomadCmd.PersistentFlags().StringVar(&installPath,
		"install-path",
		"",
		"Install path for the Nomad binary. Uses GOBIN environment variable by default.",
	)

	nomadCmd.AddCommand(nomadLogsCollectorCmd)
}
