package cmd

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/installer"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/spf13/cobra"
)

var installBinaries bool

var netBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap generates and starts new network",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.ParseConfigFile(configFilePath, homePath, types.DefaultGeneratedServices())
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}

		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		defer func() {
			cmd.SilenceUsage = true
		}()

		conf.OutputDir = &homePath

		releaseTag := getReleaseTag(installBinaries)
		if releaseTag != "" {
			inst := installer.New(conf.BinariesDir(), installPath)

			installedBinsPaths, err := inst.Install(cmd.Context(), releaseTag)
			if err != nil {
				return fmt.Errorf("failed to install dependencies: %w", err)
			}

			conf.SetBinaryPaths(installedBinsPaths)
		}

		netState.Config = conf
		updatedNetState, err := netGenerate(*netState, forceGenerate)
		if err != nil {
			return fmt.Errorf("failed to generate network: %w", err)
		}

		if err := updatedNetState.Persist(); err != nil {
			return fmt.Errorf("failed to persist network state after bootstrap/generate command: %w", err)
		}

		updatedNetState, err = netStart(context.Background(), *updatedNetState)
		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return updatedNetState.Persist()
	},
}

func init() {
	netBootstrapCmd.PersistentFlags().BoolVar(&forceGenerate,
		"force",
		false,
		"Force creating even if folders exists",
	)
	netBootstrapCmd.PersistentFlags().BoolVar(&installBinaries,
		"install",
		false,
		"Automatically installs latest version of vega, data-node and wallet binaries.",
	)
	netBootstrapCmd.PersistentFlags().StringVar(&installReleaseTag,
		"install-release-tag",
		"",
		"Installs specific release tag version of vega, data-node and wallet binaries.",
	)
	netBootstrapCmd.PersistentFlags().StringVar(&configFilePath,
		"config-path",
		"",
		"Path to the config file to generate network from",
	)
	netBootstrapCmd.PersistentFlags().BoolVar(&doNotStopAllJobsOnFailure,
		"do-not-stop-on-failure",
		false,
		"Do not stop partially running network when failed to start",
	)
	netBootstrapCmd.MarkFlagRequired("config-path")
}
