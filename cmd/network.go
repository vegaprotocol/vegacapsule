package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func networkNotBootstrappedErr(cmd string) error {
	return fmt.Errorf("failed to %s network: network not bootstrapped. Use the 'bootstrap' subcommand or provide different network home with the `--home-path` flag", cmd)
}

func networkNotRunningErr(cmd string) error {
	return fmt.Errorf("failed to %s network: network is not running. Use the 'start' subcommand or provide different network home with the `--home-path` flag", cmd)
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manages network",
	Long:  "The command allows common network commands like start/stop/generate/destroy etc..",
	Example: `# Generate the network config files
	vegacapsule network generate -home-path=/var/tmp/veganetwork/testnetwork -config-path=config.hcl
	
	# Starts the network
	vegacapsule network start -home-path=/var/tmp/veganetwork/testnetwork
	
	# Stop the network
	vegacapsule network stop -home-path=/var/tmp/veganetwork/testnetwork
	
	# Resume the network with previous configuration
	vegacapsule network start -home-path=/var/tmp/veganetwork/testnetwork
	
	# Destroy the network
	vegacapsule network destroy -home-path=/var/tmp/veganetwork/testnetwork`,
}

func init() {
	networkCmd.AddCommand(netStartCmd)
	networkCmd.AddCommand(netStopCmd)
	networkCmd.AddCommand(netDestroyCmd)
	networkCmd.AddCommand(netBootstrapCmd)
	networkCmd.AddCommand(netGenerateCmd)
	networkCmd.AddCommand(netLogsCmd)
}
