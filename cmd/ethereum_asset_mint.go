package cmd

import (
	"context"
	"fmt"
	"math/big"

	vgethereum "code.vegaprotocol.io/shared/libs/ethereum"
	"code.vegaprotocol.io/vegacapsule/state"
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

var (
	ethereumAssetMintFlags = struct {
		toAddress   string
		assetSymbol string
		amount      int64
	}{}
)

func init() {
	ethereumAssetMintCmd.Flags().StringVar(&ethereumAssetMintFlags.assetSymbol, "asset-symbol", "", "symbol of the asset to be minted")
	ethereumAssetMintCmd.Flags().StringVar(&ethereumAssetMintFlags.toAddress, "to-addr", "", "address of where the token will be minted to")
	ethereumAssetMintCmd.Flags().Int64Var(&ethereumAssetMintFlags.amount, "amount", 0, "amount to be minted")
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

		smartContracts, err := conf.SmartContractsInfo()
		if err != nil {
			return fmt.Errorf("failed getting smart contract informations: %w", err)
		}

		asset := conf.GetSmartContractToken(ethereumAssetMintFlags.assetSymbol)
		if asset == nil {
			return fmt.Errorf("failed to get non existing asset: %q", ethereumAssetMintFlags.assetSymbol)
		}

		netAddr, err := ethereumEndpointAddress(conf)
		if err != nil {
			return fmt.Errorf("failed to parse Ethereum network address: %w", err)
		}

		mintArgs := ethereumAssetMintArgs{
			amount:          ethereumAssetMintFlags.amount,
			ownerPrivateKey: smartContracts.EthereumOwner.Private,
			toAddress:       ethereumAssetMintFlags.toAddress,
			assetAddress:    asset.EthereumAddress,
			networkAddress:  netAddr,
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
