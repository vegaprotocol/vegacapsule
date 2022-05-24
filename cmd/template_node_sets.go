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

const (
	vegaNodeSetTemplateType       = "vega"
	tendermintNodeSetTemplateType = "tendermint"
	dataNodeNodeSetTemplateType   = "datanode"
)

var (
	nodeSetGroupName    string
	nodeSetName         string
	nodeSetTemplateType string

	nodeSetTemplateTypes = []string{vegaNodeSetTemplateType, tendermintNodeSetTemplateType, dataNodeNodeSetTemplateType}
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

		return templateNodeSets(nodeSetTemplateType, string(template), networkState)
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

	templateNodeSetsCmd.MarkPersistentFlagRequired("type") // nolint:errcheck
}

type templateFunc func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error)

func templateNodeSets(tmplType string, templateRaw string, netState *state.NetworkState) error {
	nodeSets, err := getNodeSetsByNames(netState)
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
	tmplType string,
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

		fileName := fmt.Sprintf("%s-%s.conf", tmplType, ns.Name)
		if err := outputTemplate(buff, fileName); err != nil {
			return err
		}
	}

	return nil
}

func getNodeSetsByNames(netState *state.NetworkState) ([]types.NodeSet, error) {
	if nodeSetGroupName == "" && nodeSetName == "" {
		return nil, fmt.Errorf("either of 'nodeset-name' or 'nodeset-group-name' flags must be defined to template node set")
	}

	if nodeSetName != "" {
		ns, err := netState.GeneratedServices.GetNodeSet(nodeSetName)
		if err != nil {
			return nil, err
		}

		return []types.NodeSet{*ns}, nil
	}

	nodeSets := netState.GeneratedServices.GetNodeSetsByGroupName(nodeSetGroupName)
	if len(nodeSets) == 0 {
		return nil, fmt.Errorf("node set group with name %q not found", nodeSetGroupName)
	}

	return nodeSets, nil
}
