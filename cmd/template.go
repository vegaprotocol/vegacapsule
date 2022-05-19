package cmd

import (
	"bytes"
	"fmt"
	"log"
	"path"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/spf13/cobra"
)

var (
	templatePath   string
	withMerge      bool
	templateOutDir string
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Allows to template genesis and various types of configurations",
	Long: `Allows templating of genesis and nodes sets configurations like Vega, Tendermint, Nomad.
	It's very usefull for config templates debugging or continuous update on running Capsule network.`,
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

	templateCmd.MarkPersistentFlagRequired("path")
}

func outputTemplate(buff *bytes.Buffer, fileName string) error {
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
		return nil
	}

	// print to stdout
	fmt.Printf("--- %s ----\n\n", fileName)
	fmt.Println(buff)
	fmt.Printf("\n\n")

	return nil
}
