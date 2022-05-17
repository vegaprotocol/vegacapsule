package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
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

func applyTemplate(tmplType string, template string, netState *state.NetworkState) error {
	switch tmplType {
	case tendermintTemplateType, vegaTemplateType, dataNodeTemplateType:
		return templateNodeSets(tmplType, template, netState)
	default:
		return fmt.Errorf("template type %q not implemented yet", tmplType)
	}
}

func templateNodeSets(tmplType string, template string, netState *state.NetworkState) error {
	if nodeSetGroupName == "" && nodeSetName == "" {
		return fmt.Errorf("either of 'nodeset-name' or 'nodeset-group-name' flags must be defined to template node set")
	}

	var nodeSets []types.NodeSet
	if nodeSetName != "" {
		ns, err := netState.GeneratedServices.GetNodeSet(nodeSetName)
		if err != nil {
			return err
		}

		nodeSets = append(nodeSets, *ns)
	} else {
		nodeSets = netState.GeneratedServices.GetNodeSetsByGroupName(nodeSetGroupName)
		if len(nodeSets) == 0 {
			return fmt.Errorf("node set group with name %q not found", nodeSetGroupName)
		}
	}

	switch tmplType {
	case tendermintTemplateType:
		template, err := tendermint.NewConfigTemplate(template)
		if err != nil {
			return err
		}

		gen, err := tendermint.NewConfigGenerator(netState.Config, netState.GeneratedServices.NodeSets.ToSlice())
		if err != nil {
			return err
		}

		return templateNodeSetConfig(gen, template, nodeSets)
	case vegaTemplateType, dataNodeTemplateType:
		return fmt.Errorf("template type %q not implemented yet", tmplType)
	}

	return fmt.Errorf("template type %q does not exists", tmplType)
}

func templateNodeSetConfig(gen templator, template *template.Template, nodeSets []types.NodeSet) error {
	var buff *bytes.Buffer
	var err error

	for _, ns := range nodeSets {
		// @TODO Solve this tomorrow
		if withMerge {
			buff, err = gen.TemplateAndMergeConfig(ns, template)
		} else {
			buff, err = gen.TemplateConfig(ns, template)
		}
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("tendermint-%s.conf", ns.Name)

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
