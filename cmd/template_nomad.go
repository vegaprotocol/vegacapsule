package cmd

import (
	"fmt"
	"io/ioutil"

	"code.vegaprotocol.io/vegacapsule/generator/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var templateNomadCmd = &cobra.Command{
	Use:   "nomad",
	Short: "Template Nomad job configuration for specific node set",
	RunE: func(cmd *cobra.Command, args []string) error {
		template, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template %q: %w", templatePath, err)
		}

		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed to load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("template nomad")
		}

		return templateNomad(string(template), networkState)
	},
}

func init() {
	templateNomadCmd.PersistentFlags().StringVar(&nodeSetGroupName,
		"nodeset-group-name",
		"",
		"Allows to apply template to all node sets in a specific group",
	)

	templateNomadCmd.PersistentFlags().StringVar(&nodeSetName,
		"nodeset-name",
		"",
		"Allows to apply template to a specific node set",
	)
}

func templateNomad(templateRaw string, netState *state.NetworkState) error {
	nodeSets, err := getNodeSetsByNames(netState)
	if err != nil {
		return err
	}

	for _, ns := range nodeSets {
		buff, err := nomad.GenerateTemplate(templateRaw, ns)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("nomad-%s.hcl", ns.Name)
		if err := outputTemplate(buff, fileName); err != nil {
			return err
		}
	}

	return nil
}
