package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"code.vegaprotocol.io/vegacapsule/generator/genesis"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var templateGenesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Template genesis file for network",
	RunE: func(cmd *cobra.Command, args []string) error {
		template, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template %q: %w", templatePath, err)
		}

		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed to load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("template genesis")
		}

		return templateGenesis(string(template), networkState)
	},
	Example: `
# Print the genesis generated from given template to stdout
vegacapsule template genesis --path net_confs/genesis.tmpl

# Save the genesis generated from given template to the './tpl-out' folder
vegacapsule template genesis --path net_confs/genesis.tmpl --out-dir ./tpl-out

# Update the genesis on the previously generated network
go run main.go template genesis --path net_confs/genesis.tmpl --update-network`,
}

func init() {
	templateGenesisCmd.PersistentFlags().BoolVar(&withMerge,
		"with-merge",
		false,
		"Defines whether the templated config should be merged with the originally initiated one",
	)
}

func templateGenesis(templateRaw string, netState *state.NetworkState) error {
	gen, err := genesis.NewGenerator(netState.Config, templateRaw)
	if err != nil {
		return err
	}

	var buff *bytes.Buffer

	if withMerge {
		buff, err = gen.ExecuteTemplate()
	} else {
		var tendermintGen *tendermint.ConfigGenerator
		tendermintGen, err = tendermint.NewConfigGenerator(netState.Config, netState.GeneratedServices.NodeSets.ToSlice())
		if err != nil {
			return err
		}

		buff, err = gen.Generate(netState.GeneratedServices.GetValidators(), tendermintGen.GenesisValidators(), nil)
	}
	if err != nil {
		return err
	}

	if !templateUpdateNetwork {
		return outputTemplate(buff, path.Join(templateOutDir, "genesis.json"), true)
	}

	for _, ns := range netState.GeneratedServices.NodeSets {
		if err := updateTemplateForNode(genesisTemplateType, ns, buff); err != nil {
			return fmt.Errorf("failed to update template for node %d: %w", ns.Index, err)
		}
	}

	return nil
}
