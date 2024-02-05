package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/installer"
	"code.vegaprotocol.io/vegacapsule/utils"
)

const (
	latestReleaseTag = "v0.74.0-preview.10"
)

var (
	installPath       string
	installReleaseTag string
)

var installBinariesCmd = &cobra.Command{
	Use:   "install-bins",
	Short: "Automatically download and install supported versions of vega, vegawallet and data-node binaries.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if installPath != "" {
			info, err := os.Lstat(installPath)
			if err != nil {
				return fmt.Errorf("failed to get info about install-path %q: %w", installPath, err)
			}

			if !info.IsDir() {
				return fmt.Errorf("install-path should be a should be a directory")
			}
		}

		conf, err := config.DefaultConfig()
		if err != nil {
			return err
		}

		inst := installer.New(conf.BinariesDir(), installPath)

		installedBinsPaths, err := inst.Install(cmd.Context(), getReleaseTag(true))
		if err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}

		if installPath != "" {
			installedBins := make([]string, 0, len(installedBinsPaths))
			for binName := range installedBinsPaths {
				installedBins = append(installedBins, binName)
			}

			if err := utils.BinariesAccessible(installedBins...); err != nil {
				return fmt.Errorf("failed to lookup installed binaries, please check %q is in $PATH: %w", installPath, err)
			}
		}

		return nil
	},
}

func init() {
	installBinariesCmd.PersistentFlags().StringVar(&installPath,
		"install-path",
		"",
		"Install path for the binaries.",
	)
	installBinariesCmd.PersistentFlags().StringVar(&installReleaseTag,
		"install-release-tag",
		latestReleaseTag,
		"Automatically installs specific release tag version of vega, data-node and wallet binaries.",
	)
}

func getReleaseTag(installBinaries bool) string {
	var releaseTag string
	if installReleaseTag != "" {
		releaseTag = installReleaseTag
	} else if installBinaries {
		releaseTag = latestReleaseTag
	}
	return releaseTag
}
