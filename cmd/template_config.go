package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/spf13/cobra"
)

const (
	tendermintTemplateType = "tendermint"

	// @TODO
	genesisTemplateType  = "genesis"
	vegaTemplateType     = "vega"
	dataNodeTemplateType = "datanode"
	nomadTemplateType    = "nomad"
	walletTemplateType   = "wallet"
	faucetTemplateType   = "wallet"
)

var (
	templateTypes = []string{genesisTemplateType, vegaTemplateType, tendermintTemplateType, dataNodeTemplateType, nomadTemplateType}

	templatePath     string
	templateType     string
	nodeSetGroupName string
	nodeSetName      string

	withMerge      bool
	templateOutDir string
)

var templateConfigCmd = &cobra.Command{
	Use:   "template-config",
	Short: "Run config templating and prints templated config to stdout. Useful for debugging.",
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
			return networkNotBootstrappedErr("nodes add")
		}

		return applyTemplate(templateType, string(template), networkState)
	},
}

func init() {
	// general required flags
	templateConfigCmd.PersistentFlags().StringVar(&templatePath,
		"path",
		"",
		"Path to the config that should be templated",
	)
	templateConfigCmd.PersistentFlags().StringVar(&templateType,
		"type",
		"",
		fmt.Sprintf("Template type, one of: %v", templateTypes),
	)

	// general optional flags
	templateConfigCmd.PersistentFlags().BoolVar(&withMerge,
		"with-merge",
		false,
		"Defines whether the templated config should be merged with the originally initiated one",
	)
	templateConfigCmd.PersistentFlags().StringVar(&templateOutDir,
		"out-dir",
		"",
		"Directory where the templated configs will be saved. If empty all will be printed to stdout",
	)

	// node sets optional flags
	templateConfigCmd.PersistentFlags().StringVar(&nodeSetGroupName,
		"nodeset-group-name",
		"",
		"Allows to apply template to all node sets in a specific group",
	)
	templateConfigCmd.PersistentFlags().StringVar(&nodeSetName,
		"nodeset-name",
		"",
		"Allows to apply template to a specific node set",
	)

	templateConfigCmd.MarkPersistentFlagRequired("path")
	templateConfigCmd.MarkPersistentFlagRequired("type")
}

type templator interface {
	TemplateConfig(types.NodeSet, *template.Template) (*bytes.Buffer, error)
	TemplateAndMergeConfig(types.NodeSet, *template.Template) (*bytes.Buffer, error)
}

type templateFunc func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error)

type templatorWrapper struct {
	templateConfig         templateFunc
	templateAndMergeConfig templateFunc
}

func (tw templatorWrapper) TemplateConfig(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error) {
	return tw.templateConfig(ns, tmpl)
}

func (tw templatorWrapper) TemplateAndMergeConfig(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error) {
	return tw.templateAndMergeConfig(ns, tmpl)
}

func applyTemplate(tmplType string, template string, netState *state.NetworkState) error {
	switch tmplType {
	case tendermintTemplateType, vegaTemplateType, dataNodeTemplateType:
		return templateNodeSets(tmplType, template, netState)
	default:
		return fmt.Errorf("template type %q not implemented yet", tmplType)
	}
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

func templateNodeSets(tmplType string, templateRaw string, netState *state.NetworkState) error {
	nodeSets, err := getNodeSetsByNames(netState)
	if err != nil {
		return err
	}

	switch tmplType {
	case tendermintTemplateType:
		tmpl, err := tendermint.NewConfigTemplate(templateRaw)
		if err != nil {
			return err
		}

		gen, err := tendermint.NewConfigGenerator(netState.Config, netState.GeneratedServices.NodeSets.ToSlice())
		if err != nil {
			return err
		}

		return templateNodeSetConfig(gen, tmplType, tmpl, nodeSets)
	case vegaTemplateType:
		tmpl, err := vega.NewConfigTemplate(templateRaw)
		if err != nil {
			return err
		}

		gen, err := vega.NewConfigGenerator(netState.Config)
		if err != nil {
			return err
		}

		genWrapper := templatorWrapper{
			templateConfig: func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error) {
				return gen.TemplateConfig(ns, netState.GeneratedServices.Faucet, tmpl)
			},
			templateAndMergeConfig: func(ns types.NodeSet, tmpl *template.Template) (*bytes.Buffer, error) {
				return gen.TemplateAndMergeConfig(ns, netState.GeneratedServices.Faucet, tmpl)
			},
		}

		return templateNodeSetConfig(genWrapper, tmplType, tmpl, nodeSets)
	case dataNodeTemplateType:
		tmpl, err := datanode.NewConfigTemplate(templateRaw)
		if err != nil {
			return err
		}

		gen, err := datanode.NewConfigGenerator(netState.Config)
		if err != nil {
			return err
		}

		return templateNodeSetConfig(gen, tmplType, tmpl, nodeSets)
	}

	return fmt.Errorf("template type %q does not exists", tmplType)
}

func templateNodeSetConfig(gen templator, tmplType string, template *template.Template, nodeSets []types.NodeSet) error {
	var buff *bytes.Buffer
	var err error

	for _, ns := range nodeSets {
		if withMerge {
			buff, err = gen.TemplateAndMergeConfig(ns, template)
		} else {
			buff, err = gen.TemplateConfig(ns, template)
		}
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s-%s.conf", tmplType, ns.Name)

		if len(templateOutDir) != 0 {
			filePath := path.Join(templateOutDir, fileName)
			f, err := utils.CreateFile(filePath)
			if err != nil {
				return err
			}

			if _, err := f.Write(buff.Bytes()); err != nil {
				return err
			}

			log.Printf("Saving file to %q", filePath)
			continue
		}

		// print to stdout
		fmt.Printf("--- %s ----\n\n", fileName)
		fmt.Println(buff)
		fmt.Printf("\n\n")
	}

	return nil
}
