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
	Example: `
# Generate the nomad configuration for multiple node sets #1
vegacapsule template nomad --path .../node_set.tmpl --nodeset-group-name validators,full

# Generate the nomad configuration for multiple node sets #2
vegacapsule template nomad --path .../node_set.tmpl --nodeset-group-name validators --nodeset-group-name full`,
}

func init() {
	templateNomadCmd.PersistentFlags().StringSliceVar(&nodeSetsGroupsNames,
		"nodeset-group-name",
		[]string{},
		"Allows to apply template to all node sets in a specific group",
	)

	templateNomadCmd.PersistentFlags().StringSliceVar(&nodeSetsNames,
		"nodeset-name",
		[]string{},
		"Allows to apply template to a specific node set",
	)
}

func templateNomad(templateRaw string, netState *state.NetworkState) error {
	nodeSets, err := filterNodesSets(netState, nodeSetsNames, nodeSetsGroupsNames)
	if err != nil {
		return err
	}

	for _, ns := range nodeSets {
		buff, err := nomad.GenerateNodeSetTemplate(templateRaw, ns)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("nomad-%s.hcl", ns.Name)
		if err := outputTemplate(buff, templateOutDir, fileName, true); err != nil {
			return err
		}
	}

	return nil
}
