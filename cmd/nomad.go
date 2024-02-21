package cmd

import (
	"fmt"
	"os"

	nmrunner "code.vegaprotocol.io/vegacapsule/nomad/runner"

	"github.com/spf13/cobra"
)

var nomadConfigPath string

func getInstallPath(installPath string) (string, error) {
	if len(installPath) != 0 {
		return installPath, nil
	}

	installPath = os.Getenv("GOBIN")
	if len(installPath) == 0 {
		return "", fmt.Errorf("GOBIN environment variable has not been found - please set install-path flag instead")
	}

	return installPath, nil
}

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
	nomadCmd.Flags().StringVar(&nomadConfigPath,
		"nomad-config-path",
		"",
		"Allows to use Nomad configuration",
	)
	nomadCmd.Flags().StringVar(&installPath,
		"install-path",
		"",
		"Install path for the Nomad binary. Uses GOBIN environment variable by default.",
	)

	nomadCmd.AddCommand(nomadLogsCollectorCmd)
}
