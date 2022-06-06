package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var (
	nodeSetGroupName    string
	nodeSetName         string
	nodeSetTemplateType string
	nodeSetMode         string

	nodeSetTemplateTypes = []templateKindType{vegaNodeSetTemplateType, tendermintNodeSetTemplateType, dataNodeNodeSetTemplateType}
)

var templateNodeSetsCmd = &cobra.Command{
	Use:   "node-sets",
	Short: "Run config templating for Vega, Tendermit, DataNode node sets",
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
			return networkNotBootstrappedErr("template node-sets")
		}

		return templateNodeSets(templateKindType(nodeSetTemplateType), string(template), networkState)
	},
}

func init() {
	templateNodeSetsCmd.PersistentFlags().StringVar(&nodeSetTemplateType,
		"type",
		"",
		fmt.Sprintf("Template type, one of: %v", nodeSetTemplateTypes),
	)

	templateNodeSetsCmd.PersistentFlags().BoolVar(&withMerge,
		"with-merge",
		false,
		"Defines whether the templated config should be merged with the originally initiated one",
	)

	templateNodeSetsCmd.PersistentFlags().StringVar(&nodeSetGroupName,
		"nodeset-group-name",
		"",
		"Allows to apply template to all node sets in a specific group",
	)

	templateNodeSetsCmd.PersistentFlags().StringVar(&nodeSetName,
		"nodeset-name",
		"",
		"Allows to apply template to a specific node set",
	)

	templateNodeSetsCmd.PersistentFlags().StringVar(&nodeSetMode,
		"nodeset-mode",
		"",
		"Allows to apply template to a specific node set types",
	)

	templateNodeSetsCmd.MarkPersistentFlagRequired("type") // nolint:errcheck
}

type templateFunc func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error)

func templateNodeSets(tmplType templateKindType, templateRaw string, netState *state.NetworkState) error {
	nodeSets, err := filterNodesSets(netState, nodeSetName, nodeSetGroupName, nodeSetMode)
	if err != nil {
		return err
	}

	switch tmplType {
	case tendermintNodeSetTemplateType:
		tmpl, err := tendermint.NewConfigTemplate(templateRaw)
		if err != nil {
			return err
		}

		gen, err := tendermint.NewConfigGenerator(netState.Config, netState.GeneratedServices.NodeSets.ToSlice())
		if err != nil {
			return err
		}

		return templateNodeSetConfig(gen.TemplateConfig, gen.TemplateAndMergeConfig, tmplType, tmpl, nodeSets)
	case vegaNodeSetTemplateType:
		tmpl, err := vega.NewConfigTemplate(templateRaw)
		if err != nil {
			return err
		}

		gen, err := vega.NewConfigGenerator(netState.Config)
		if err != nil {
			return err
		}

		templateF := func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error) {
			return gen.TemplateConfig(ns, netState.GeneratedServices.Faucet, tmpl)
		}

		templateAndMergeF := func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error) {
			return gen.TemplateAndMergeConfig(ns, netState.GeneratedServices.Faucet, tmpl)
		}

		return templateNodeSetConfig(templateF, templateAndMergeF, tmplType, tmpl, nodeSets)
	case dataNodeNodeSetTemplateType:
		tmpl, err := datanode.NewConfigTemplate(templateRaw)
		if err != nil {
			return err
		}

		gen, err := datanode.NewConfigGenerator(netState.Config)
		if err != nil {
			return err
		}

		return templateNodeSetConfig(gen.TemplateConfig, gen.TemplateAndMergeConfig, tmplType, tmpl, nodeSets)
	}

	return fmt.Errorf("template type %q does not exists", tmplType)
}

func templateNodeSetConfig(
	templateF, templateAndMergeF templateFunc,
	tmplType templateKindType,
	template *template.Template,
	nodeSets []types.NodeSet,
) error {
	var buff *bytes.Buffer
	var err error

	for _, ns := range nodeSets {
		if withMerge {
			buff, err = templateAndMergeF(ns, template)
		} else {
			buff, err = templateF(ns, template)
		}
		if err != nil {
			return err
		}

		if templateUpdateNetwork {
			if err := updateTemplateForNode(tmplType, ns.Tendermint.HomeDir, buff); err != nil {
				return fmt.Errorf("failed to update template for node %d: %w", ns.Index, err)
			}
		} else {
			fileName := fmt.Sprintf("%s-%s.conf", tmplType, ns.Name)
			if err := outputTemplate(buff, templateOutDir, fileName, true); err != nil {
				return fmt.Errorf("failed to print generated template for node %d: %w", ns.Index, err)
			}
		}
	}

	return nil
}

func filterNodesSets(netState *state.NetworkState, nodeSetName, nodeSetGroupName, nodeSetMode string) ([]types.NodeSet, error) {
	if nodeSetGroupName == "" && nodeSetName == "" && nodeSetMode == "" {
		return nil, fmt.Errorf("either of 'nodeset-name', 'nodeset-group-name' or 'nodeset-mode' flags must be defined to template node set")
	}

	if nodeSetName != "" {
		ns, err := netState.GeneratedServices.GetNodeSet(nodeSetName)
		if err != nil {
			return nil, err
		}

		return []types.NodeSet{*ns}, nil
	}

	filters := []types.NodeSetFilter{}
	if nodeSetGroupName != "" {
		filters = append(filters, types.NodeSetFilterByGroupName(nodeSetGroupName))
	}

	if nodeSetMode != "" {
		filters = append(filters, types.NodeSetFilterByMode(nodeSetMode))
	}

	nodeSets := types.FilterNodeSets(netState.GeneratedServices.NodeSets.ToSlice(), filters...)
	if len(nodeSets) == 0 {
		return nil, fmt.Errorf("node set group with given criteria [nodeset-name: '%s', nodeset-group-name: '%s', nodeset-mode: '%s'] not found", nodeSetName, nodeSetGroupName, nodeSetMode)
	}

	return nodeSets, nil
}
