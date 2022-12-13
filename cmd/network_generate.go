package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/spf13/cobra"
)

var (
	forceGenerate  bool
	configFilePath string
)

var netGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate new network from configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.ParseConfigFile(configFilePath, homePath, types.DefaultGeneratedServices())
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}

		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}
		netState.Config = conf

		updatedNetState, err := netGenerate(*netState, forceGenerate)
		if err != nil {
			return fmt.Errorf("failed to generate network: %w", err)
		}

		return updatedNetState.Persist()
	},
}

func init() {
	netGenerateCmd.PersistentFlags().BoolVar(&forceGenerate,
		"force",
		false,
		"Force creating even if folders exists",
	)
	netGenerateCmd.PersistentFlags().StringVar(&configFilePath,
		"config-path",
		"",
		"Path to the config file to generate network from",
	)
	netGenerateCmd.MarkFlagRequired("config-path")
}

func netGenerate(state state.NetworkState, force bool) (*state.NetworkState, error) {
	if force {
		if err := os.RemoveAll(*state.Config.OutputDir); err != nil {
			return nil, fmt.Errorf("failed to remove output folder with --force flag: %w", err)
		}
	} else if state.GeneratedServices != nil {
		return nil, fmt.Errorf("failed to generate network: network is already generated")
	}

	if netDirEmpty, _ := utils.DirEmpty(*state.Config.OutputDir, filepath.Base(state.Config.BinariesDir())); !netDirEmpty {
		return nil, fmt.Errorf("output directory %q already exists and it's not empty", *state.Config.OutputDir)
	}

	log.Println("generating network")

	state.VegaChainID = state.Config.Network.Name + "-001"

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	gen, err := generator.New(state.Config, types.GeneratedServices{}, nomadRunner, state.VegaChainID)
	if err != nil {
		return nil, err
	}

	generatedSvcs, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	if err := state.Config.Persist(); err != nil {
		return nil, fmt.Errorf("failed to persist config in output directory %s: %w", *state.Config.OutputDir, err)
	}

	log.Println("generating network success")

	state.GeneratedServices = generatedSvcs

	state.RunningJobs = &types.NetworkJobs{}
	state.RunningJobs.AddExtraJobIDs(generatedSvcs.PreGenerateJobsIDs())
	return &state, nil
}
