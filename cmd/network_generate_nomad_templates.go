package cmd

import (
	"fmt"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var (
	networkNomadTemplatesOutputDir string
)

var netGenerateNomadTemplatesCmd = &cobra.Command{
	Use:   "generate-nomad-templates",
	Short: "Generate nomad job templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("generate-nomad-templates")
		}

		jobs, err := nomad.GenerateNomadNetworkJobs(netState.Config, netState.GeneratedServices)
		if err != nil {
			return fmt.Errorf("failed to generate nomad jobs for network: %w", err)
		}

		if networkNomadTemplatesOutputDir == "" {
			networkNomadTemplatesOutputDir = filepath.Join(*netState.Config.OutputDir, "nomad-jobs")
		}

		if err := nomad.PersistJobsTemplates(networkNomadTemplatesOutputDir, jobs); err != nil {
			return fmt.Errorf("failed to persist nomad network templates: %w", err)
		}

		fmt.Printf("Templates saved successfully to %s\n", networkNomadTemplatesOutputDir)

		return nil
	},
}

func init() {
	netGenerateNomadTemplatesCmd.PersistentFlags().StringVar(&networkNomadTemplatesOutputDir,
		"output-dir",
		"",
		"Output folder for job teemplates. If empty, the <network-home>/nomad-jobs is used",
	)
}
