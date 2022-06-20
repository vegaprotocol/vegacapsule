package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"code.vegaprotocol.io/vegacapsule/generator/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
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

		updatedNetworkState, err := templateNomad(string(template), networkState, templateOutDir, templateUpdateNetwork)
		if err != nil {
			return fmt.Errorf("failed to template nomad jobs for nodes ets: %w", err)
		}

		if templateUpdateNetwork {
			log.Printf("Updating nomad template in the network state")
			if err := updatedNetworkState.Persist(); err != nil {
				return fmt.Errorf("failed to save network state: %w", err)
			}
		}

		return nil
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

func templateNomad(templateRaw string, netState *state.NetworkState, templateOutDir string, updateNetwork bool) (*state.NetworkState, error) {
	nodeSets, err := filterNodesSets(netState, nodeSetsNames, nodeSetsGroupsNames)
	if err != nil {
		return nil, err
	}

	newNetworkState := *netState
	for _, ns := range nodeSets {
		buff, err := nomad.GenerateNodeSetTemplate(templateRaw, ns)
		if err != nil {
			return nil, err
		}

		newNetworkState = updateNomadTemplateInTheNetworkState(newNetworkState, ns, buff)
		// Don not print template when updating network
		if updateNetwork {
			continue
		}
		fileName := fmt.Sprintf("nomad-%s.hcl", ns.Name)
		if err := outputTemplate(buff, templateOutDir, fileName, true); err != nil {
			return nil, err
		}
	}

	return &newNetworkState, nil
}

func updateNomadTemplateInTheNetworkState(netState state.NetworkState, modifiedNodeSet types.NodeSet, templateRaw *bytes.Buffer) state.NetworkState {
	for idx, ns := range netState.GeneratedServices.NodeSets {
		if ns.Name != modifiedNodeSet.Name || ns.Index != modifiedNodeSet.Index {
			continue
		}
		ns.NomadJobRaw = utils.StrPoint(templateRaw.String())

		netState.GeneratedServices.NodeSets[idx] = ns
	}

	return netState
}
