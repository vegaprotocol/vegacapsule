package cmd

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/ethereum"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var smartcontractsCmd = &cobra.Command{
	Use:   "smartcontracts",
	Short: "Support interactions with smartcontracts",
}

// Flags
var (
	ethereumAddress string
	ethereumChainID int
)

const (
	defaultEthereumAddress = "ws://127.0.0.1:8545"
	defaultEthereumChainID = 1440
)

func init() {
	smartcontractsCmd.PersistentFlags().StringVar(&ethereumAddress,
		"eth-address",
		defaultEthereumAddress,
		"Specify the ethereum network address",
	)
	smartcontractsCmd.PersistentFlags().IntVar(&ethereumChainID,
		"eth-chain-id",
		defaultEthereumChainID,
		"Specify the ethereum chain ID",
	)

	smartcontractsCmd.AddCommand(smartContractsMultisigCmd)
	smartContractsMultisigCmd.AddCommand(smartContractsMultisigSetupCmd)
}

var smartContractsMultisigCmd = &cobra.Command{
	Use:   "multisig",
	Short: "Manages multisig smartcontract",
}

var smartContractsMultisigSetupCmd = &cobra.Command{
	Use:   "init",
	Short: "Setups the multisig smart contract",
	Long:  `Adds all validators to the multisig smart contract`,
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state get-smartcontracts-addresses")
		}

		if !netState.Running() {
			return networkNotRunningErr("smartcontracts multisig init")
		}

		smartcontracts, err := netState.Config.SmartContractsInfo()
		if err != nil {
			return fmt.Errorf("error getting smart contract informations: %w", err)
		}

		ctx := context.Background()
		client, err := ethereum.NewEthereumClient(ctx, *netState.Config.VegaBinary, ethereumChainID, ethereumAddress, *smartcontracts)
		if err != nil {
			return fmt.Errorf("failed to create ethereum client: %w", err)
		}

		validatorsKeyPairs := getValidatorsEthKeyPairs(netState.GeneratedServices.ListValidators())
		return client.InitMultisig(ctx, *smartcontracts, validatorsKeyPairs)
	},
}

func getValidatorsEthKeyPairs(nodes []types.VegaNodeOutput) []ethereum.KeyPair {
	result := make([]ethereum.KeyPair, len(nodes))
	for idx, node := range nodes {
		result[idx] = ethereum.KeyPair{
			Address:    node.VegaNode.NodeWalletInfo.EthereumAddress,
			PrivateKey: node.VegaNode.NodeWalletInfo.EthereumPrivateKey,
		}
	}

	return result
}
