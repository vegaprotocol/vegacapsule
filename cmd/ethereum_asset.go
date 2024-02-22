package cmd

import (
	"fmt"
	"time"

	vgtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

var ethereumAssetCmd = &cobra.Command{
	Use:   "asset",
	Short: "Allows to deposit/stake/mint tokens through smartcontract",
}

func init() {
	ethereumAssetCmd.AddCommand(ethereumAssetStakeCmd)
	ethereumAssetCmd.AddCommand(ethereumAssetDepositCmd)
	ethereumAssetCmd.AddCommand(ethereumAssetMintCmd)
}

func printEthereumTx(tx *vgtypes.Transaction) error {
	txJSON, err := tx.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal transaction to JSON: %w", err)
	}

	fmt.Printf("Transaction: %s", txJSON)

	return nil
}

func defeaultSyncTimeout() time.Duration {
	return time.Second * defaultEthreumWaitTimeout
}
