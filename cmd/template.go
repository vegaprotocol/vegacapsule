package cmd

import (
	"bytes"
	"fmt"
	"log"
	"path"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/spf13/cobra"

	datanodegen "code.vegaprotocol.io/vegacapsule/generator/datanode"
	genesisgen "code.vegaprotocol.io/vegacapsule/generator/genesis"
	tmgen "code.vegaprotocol.io/vegacapsule/generator/tendermint"
	vegagen "code.vegaprotocol.io/vegacapsule/generator/vega"
)

type templateKindType string

const (
	vegaNodeSetTemplateType       templateKindType = "vega"
	tendermintNodeSetTemplateType templateKindType = "tendermint"
	dataNodeNodeSetTemplateType   templateKindType = "data-node"

	genesisTemplateType templateKindType = "genesis"
)

var (
	templatePath          string
	withMerge             bool
	templateOutDir        string
	templateUpdateNetwork bool
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Allows to template genesis and various types of configurations",
	Long: `The function allows templating for genesis and node-sets configurations
like Vega, Tendermint, and Nomad. It's useful for config templates debugging or
continuous updates on running Capsule network.

By default, the command prints generated templates to the stdout. If the "out-dir"
flag is specified, the command saves generated templates to the given directory.

If you have a network generated and running, you can update templates for all nodes
by specifying the "update-network" flag.

If you set the "update-network" flag, the command does not print a template to stdout
or save it to the "out-dir" folder - files in the "network-home" are modified only.`,
}

func init() {
	templateCmd.AddCommand(templateNodeSetsCmd)
	templateCmd.AddCommand(templateGenesisCmd)
	templateCmd.AddCommand(templateNomadCmd)

	templateCmd.PersistentFlags().StringVar(&templatePath,
		"path",
		"",
		"Path to the config that should be templated",
	)

	templateCmd.PersistentFlags().StringVar(&templateOutDir,
		"out-dir",
		"",
		"Directory where the templated configs will be saved. If empty all will be printed to stdout",
	)

	templateCmd.PersistentFlags().BoolVar(&templateUpdateNetwork,
		"update-network",
		false,
		"Flag defines if previously generated configuration for network should be updated",
	)

	templateCmd.MarkPersistentFlagRequired("path") // nolint:errcheck
}

// updateTemplateForNode writes given template to the given node
func updateTemplateForNode(kind templateKindType, nodeHomePath string, buff *bytes.Buffer) error {
	configTypeFilePathMap := map[templateKindType]string{
		vegaNodeSetTemplateType:       vegagen.ConfigFilePath(""),
		tendermintNodeSetTemplateType: tmgen.ConfigFilePath(""),
		dataNodeNodeSetTemplateType:   datanodegen.ConfigFilePath(""),
		genesisTemplateType:           genesisgen.ConfigFilePath(""),
	}

	configFilePath, templateSupported := configTypeFilePathMap[kind]
	if !templateSupported {
		return fmt.Errorf("failed to update the %v template for node %s: template type not supported", kind, nodeHomePath)
	}

	return outputTemplate(buff, nodeHomePath, configFilePath, false)
}

func outputTemplate(buff *bytes.Buffer, templateOutDir, fileName string, writeToStdOut bool) error {
	if len(templateOutDir) != 0 {
		filePath := path.Join(templateOutDir, fileName)
		f, err := utils.CreateFile(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.Write(buff.Bytes()); err != nil {
			return err
		}

		log.Printf("Saving file to %q", filePath)
		return nil
	}

	if !writeToStdOut {
		return nil
	}

	// print to stdout
	fmt.Printf("--- %s ----\n\n", fileName)
	fmt.Println(buff)
	fmt.Printf("\n\n")

	return nil
}
