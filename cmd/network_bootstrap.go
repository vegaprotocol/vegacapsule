package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var netBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap generates and starts new network",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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

		saveState := func(updatedNetState *state.NetworkState) {
			if updatedNetState == nil {
				return
			}

			log.Printf("saving network state to the file")
			saveErr := updatedNetState.Persist()
			if saveErr != nil {
				log.Printf("failed to save network state: %s", err)
			}

			// do not shadow the original error as it is more important
			if err == nil {
				err = saveErr
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		updatedNetState, err = netStart(ctx, *updatedNetState)
		// We want state saved even if the network is started with error
		defer saveState(updatedNetState)

		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return err
	},
}

func init() {
	netBootstrapCmd.PersistentFlags().Uint64Var(&timeout,
		"timeout",
		defaultTimeout,
		"Timeout in seconds",
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
