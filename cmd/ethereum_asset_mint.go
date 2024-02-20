package cmd

import (
	"context"
	"fmt"
	"math/big"

	vgethereum "code.vegaprotocol.io/vegacapsule/libs/ethereum"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

type ethereumAssetMintArgs struct {
	amount          int64
	ownerPrivateKey string
	toAddress       string
	assetAddress    string
	networkAddress  string
}

var ethereumAssetMintFlags = struct {
	toAddress   string
	assetSymbol string
	amount      int64
	bridge      string
}{}

func init() {
	ethereumAssetMintCmd.Flags().StringVar(&ethereumAssetMintFlags.assetSymbol, "asset-symbol", "", "symbol of the asset to be minted")
	ethereumAssetMintCmd.Flags().StringVar(&ethereumAssetMintFlags.toAddress, "to-addr", "", "address of where the token will be minted to")
	ethereumAssetMintCmd.Flags().Int64Var(&ethereumAssetMintFlags.amount, "amount", 0, "amount to be minted")
	ethereumAssetMintCmd.Flags().StringVar(&ethereumAssetMintFlags.bridge, "bridge", "primary", "bridge linked to the deposit")
	ethereumAssetMintCmd.MarkFlagRequired("asset-symbol")
	ethereumAssetMintCmd.MarkFlagRequired("to-address")
	ethereumAssetMintCmd.MarkFlagRequired("amount")
}

var ethereumAssetMintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint allows an asset to be minted by a Base Faucet Token contract.",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("ethereum asset mint")
		}

		if !netState.Running() {
			return networkNotRunningErr("ethereum asset mint")
		}

		conf := netState.Config

		asset := conf.GetSmartContractToken(ethereumAssetMintFlags.assetSymbol)
		if asset == nil {
			return fmt.Errorf("failed to get non existing asset: %q", ethereumAssetMintFlags.assetSymbol)
		}

		var (
			networkAddress string
			smartContracts *types.SmartContractsInfo
		)
		switch ethereumAssetMintFlags.bridge {
		case "primary":
			networkAddress = conf.Network.Ethereum.Endpoint
			smartContracts, err = conf.PrimarySmartContractsInfo()
			if err != nil {
				return fmt.Errorf("failed getting primary smart contract informations: %w", err)
			}
		case "secondary":
			networkAddress = conf.Network.SecondaryEthereum.Endpoint
			smartContracts, err = conf.SecondarySmartContractsInfo()
			if err != nil {
				return fmt.Errorf("failed getting secondary smart contract informations: %w", err)
			}
		}

		mintArgs := ethereumAssetMintArgs{
			amount:          ethereumAssetMintFlags.amount,
			ownerPrivateKey: smartContracts.EthereumOwner.Private,
			toAddress:       ethereumAssetMintFlags.toAddress,
			assetAddress:    asset.EthereumAddress,
			networkAddress:  networkAddress,
		}

		return ethereumAssetMint(cmd.Context(), mintArgs)
	},
}

func ethereumAssetMint(ctx context.Context, args ethereumAssetMintArgs) error {
	client, err := vgethereum.NewClient(ctx, args.networkAddress)
	if err != nil {
		return fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	syncTimeout := defeaultSyncTimeout()

	tokenSession, err := client.NewBaseTokenSession(
		ctx,
		args.ownerPrivateKey,
		common.HexToAddress(args.assetAddress),
		&syncTimeout,
	)
	if err != nil {
		return fmt.Errorf("failed to create base token session for %s: %w", args.assetAddress, err)
	}

	tx, err := tokenSession.MintSync(
		common.HexToAddress(args.toAddress),
		big.NewInt(args.amount),
	)
	if err != nil {
		return fmt.Errorf("failed to mint token: %w", err)
	}

	return printEthereumTx(tx)
}
