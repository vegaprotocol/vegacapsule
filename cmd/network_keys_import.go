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

const networkKeysImportDescription = `The command takes a set of keys for network nodes and imports them to the previously generated network.

Required values for each node to import keys to the network:

- Node Index - counted from 0
- Tendermint validator private key
- Tendermint node private key
- Ethereum private key
- Vega recovery phrase

Example content of the file with data to import:
  [
    {
      "node_index": 1,
      "ethereum_private_key": "someValue ...",
      "tendermint_validator_private_key": "someValue ...",
      "tendermint_node_private_key": "someValue ...",
      "vega_recovery_phrase": "some value ..."
    },
	{
      "node_index": 2,
      "ethereum_private_key": "someValue ...",
      "tendermint_validator_private_key": "someValue ...",
      "tendermint_node_private_key": "someValue ...",
      "vega_recovery_phrase": "some value ..."
    },
    {
      "node_index": 3,
      "ethereum_private_key": "someValue ...",
      "tendermint_validator_private_key": "someValue ...",
      "tendermint_node_private_key": "someValue ...",
      "vega_recovery_phrase": "some value ..."
    },
	...
  ]`

var (
	keysFilePath string
)

var netKeysImportCmd = &cobra.Command{
	Use:          "import",
	Short:        "Import pre-generated keys into the network",
	Long:         networkKeysImportDescription,
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
			return fmt.Errorf("the network must be stopped before the keys are imported")
		}

		updatedNetState, err := netKeysImport(context.Background(), *netState, keysFilePath)
		if err != nil {
			return fmt.Errorf("failed to import network keys: %w", err)
		}

		return updatedNetState.Persist()
	},
}

func init() {
	netKeysImportCmd.PersistentFlags().StringVar(&keysFilePath,
		"keys-file-path",
		"",
		"Path to the file, that includes network keys",
	)
	netKeysImportCmd.MarkFlagRequired("keys-file-path")
}

func netKeysImport(ctx context.Context, state state.NetworkState, dataFilePath string) (*state.NetworkState, error) {
	log.Println("importing network")

	keysData, err := loadKeysData(dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed load keys from given file: %w", err)
	}

	newState, err := importer.ImportKeysIntoValidatorsWallets(state, *keysData)
	if err != nil {
		return nil, fmt.Errorf("failed to import given keys into the validators wallets: %w", err)
	}

	return newState, nil
}

func loadKeysData(filePath string) (*importer.NodesKeysData, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys from the \"%s\" file: %w", filePath, err)
	}

	network := &importer.NodesKeysData{}
	if err := json.Unmarshal(data, network); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keys: %w", err)
	}

	return network, nil
}
