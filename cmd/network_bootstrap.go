package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var netBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap generates and starts new network",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.ParseConfigFile(configFilePath, homePath)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}

		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		conf.OutputDir = &homePath

		netState.Config = conf

		updatedNetState, err := netGenerate(*netState, forceGenerate)
		if err != nil {
			return fmt.Errorf("failed to generate network: %w", err)
		}

		if err := updatedNetState.Persist(); err != nil {
			return fmt.Errorf("failed to persist network state after bootstrap/generate command: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		updatedNetState, err = netStart(ctx, *updatedNetState)
		// We want state saved even if the network is started with error
		defer saveNetworkState(updatedNetState)

		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return err
	},
}

// saveNetworkState saves state to disk. The function accepts previousError which comes
// from the network operations like start, bootstrap, etc.
func saveNetworkState(updatedNetState *state.NetworkState) {
	if updatedNetState == nil {
		return
	}

	log.Printf("saving network state to the file")
	err := updatedNetState.Persist()
	if err != nil {
		log.Fatalf("failed to save network state: %s", err)
	}
}

func init() {
	netBootstrapCmd.PersistentFlags().DurationVar(&timeout,
		"timeout",
		defaultTimeout,
		"Bootstrap timeout",
	)
	netBootstrapCmd.PersistentFlags().BoolVar(&forceGenerate,
		"force",
		false,
		"Force creating even if folders exists",
	)
	netBootstrapCmd.PersistentFlags().StringVar(&configFilePath,
		"config-path",
		"",
		"Path to the config file to generate network from",
	)
	netBootstrapCmd.MarkFlagRequired("config-path")
}
