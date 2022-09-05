package cmd

import (
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/config"
	"github.com/spf13/cobra"
)

var homePath string

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Tool for generating and running vega network",
	Long:  "Configuration based tool for bootstraping and managing vega network. Primary usages are local development of Vega, testing but also deploy new production network.",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	defaultHomePath, err := config.DefaultNetworkHome()
	if err != nil {
		log.Fatalf("Failed to get default network home: %s", err.Error())
	}

	rootCmd.PersistentFlags().StringVar(&homePath,
		"home-path",
		defaultHomePath,
		"Specify the location of network home directory",
	)

	rootCmd.AddCommand(networkCmd)
	rootCmd.AddCommand(nomadCmd)
	rootCmd.AddCommand(nodesCmd)
	rootCmd.AddCommand(stateCmd)
	rootCmd.AddCommand(ethereumCmd)
	rootCmd.AddCommand(installBinariesCmd)
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(versionCmd)
}
