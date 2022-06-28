package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"code.vegaprotocol.io/vegacapsule/importer"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var (
	networkImportDataFilePath string
)

var netImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import pre-generated keys into the network",
	Long: `The command takes a set of keys for vega nodes and imports them to the previously generated network.

Required values for each node to import keys to the network:

- Node Index - counted from 0
- Tendermint validator private key
- Tendermint node private key
- Ethereum private key
- Vega recovery phrase

Example content of the file with data to import:

  [
    {
      "node_index": "1",
      "ethereum_private_key": "someValue ...",
      "tendermint_validator_private_key": "someValue ...",
      "tendermint_node_private_key": "someValue ...",
      "vega_recovery_phrase": "some value ..."
    },
	{
      "node_index": "1",
      "ethereum_private_key": "someValue ...",
      "tendermint_validator_private_key": "someValue ...",
      "tendermint_node_private_key": "someValue ...",
      "vega_recovery_phrase": "some value ..."
    },
    {
      "node_index": "1",
      "ethereum_private_key": "someValue ...",
      "tendermint_validator_private_key": "someValue ...",
      "tendermint_node_private_key": "someValue ...",
      "vega_recovery_phrase": "some value ..."
    },
	...
  ]
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("start")
		}

		if netState.Running() {
			return fmt.Errorf("the network must be stopped before network data is imported")
		}

		updatedNetState, err := netImport(context.Background(), *netState, networkImportDataFilePath)
		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return updatedNetState.Persist()
	},
}

func init() {
	netImportCmd.PersistentFlags().StringVar(&networkImportDataFilePath,
		"network-data-path",
		"",
		"Path to the file, that includes entire network information",
	)
	netImportCmd.MarkFlagRequired("network-data-path")
}

func netImport(ctx context.Context, state state.NetworkState, dataFilePath string) (*state.NetworkState, error) {
	log.Println("importing network")

	networkData, err := loadNetworkImportData(dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed load network data from given file: %w", err)
	}

	newState, err := importer.ImportDataIntoNetworkValidators(state, *networkData)
	if err != nil {
		return nil, fmt.Errorf("failed to import given data into the network: %w", err)
	}

	return newState, nil
}

func loadNetworkImportData(filePath string) (*importer.NetworkImportdata, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read network data file: %w", err)
	}

	network := &importer.NetworkImportdata{}
	if err := json.Unmarshal(data, network); err != nil {
		return nil, fmt.Errorf("failed to unmarshal network data file: %w", err)
	}

	return network, nil
}
