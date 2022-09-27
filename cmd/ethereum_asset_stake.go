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

var (
	ethereumAssetStakeFlags = struct {
		vegaPubKey  string
		assetSymbol string
		amount      int64
	}{}
)

func init() {
	ethereumAssetStakeCmd.Flags().StringVar(&ethereumAssetStakeFlags.assetSymbol, "asset-symbol", "", "symbol of the asset to be staked")
	ethereumAssetStakeCmd.Flags().StringVar(&ethereumAssetStakeFlags.vegaPubKey, "pub-key", "", "Vega public key to where the asset will be staked")
	ethereumAssetStakeCmd.Flags().Int64Var(&ethereumAssetStakeFlags.amount, "amount", 0, "amount to be staked")
	ethereumAssetStakeCmd.MarkFlagRequired("asset-symbol")
	ethereumAssetStakeCmd.MarkFlagRequired("pub-key")
	ethereumAssetStakeCmd.MarkFlagRequired("amount")
}

var ethereumAssetStakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake allows an asset to be staked to a given Vega public key.",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("ethereum asset stake")
		}

		if !netState.Running() {
			return networkNotRunningErr("ethereum asset stake")
		}

		conf := netState.Config

		smartContracts, err := conf.SmartContractsInfo()
		if err != nil {
			return fmt.Errorf("failed getting smart contract informations: %w", err)
		}

		asset := conf.GetSmartContractToken(ethereumAssetStakeFlags.assetSymbol)
		if asset == nil {
			return fmt.Errorf("failed to get non existing asset: %q", ethereumAssetStakeFlags.assetSymbol)
		}

		netAddr, err := ethereumEndpointAddress(conf)
		if err != nil {
			return fmt.Errorf("failed to parse Ethereum network address: %w", err)
		}

		stakeArgs := ethereumAssetDepositOrStakeArgs{
			amount:          ethereumAssetStakeFlags.amount,
			vegaPubKey:      ethereumAssetStakeFlags.vegaPubKey,
			ownerPrivateKey: smartContracts.EthereumOwner.Private,
			bridgeAddress:   smartContracts.StakingBridge.EthereumAddress,
			assetAddress:    asset.EthereumAddress,
			networkAddress:  netAddr,
		}

		return ethereumAssetStake(cmd.Context(), stakeArgs)
	},
}

func ethereumAssetStake(ctx context.Context, args ethereumAssetDepositOrStakeArgs) error {
	client, err := vgethereum.NewClient(ctx, args.networkAddress)
	if err != nil {
		return fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	syncTimeout := defeaultSyncTimeout()
	bridgeAddr := common.HexToAddress(args.bridgeAddress)

	bridgeSession, err := client.NewStakingBridgeSession(
		ctx,
		args.ownerPrivateKey,
		bridgeAddr,
		&syncTimeout,
	)
	if err != nil {
		return fmt.Errorf("failed to create staking bridge session for %s: %w", args.bridgeAddress, err)
	}

	tokenSession, err := client.NewBaseTokenSession(
		ctx,
		args.ownerPrivateKey,
		common.HexToAddress(args.assetAddress),
		&syncTimeout,
	)
	if err != nil {
		return fmt.Errorf("failed to create token session for %s: %w", args.assetAddress, err)
	}

	amount := big.NewInt(args.amount)

	if _, err := tokenSession.ApproveSync(bridgeAddr, amount); err != nil {
		return fmt.Errorf("failef to approve asset amount to bridge: %w", err)
	}

	vegaPubKeyArr, err := vgethereum.HexStringToByte32Array(args.vegaPubKey)
	if err != nil {
		return fmt.Errorf("failed to convert Vega pub key string to byte array: %w", err)
	}

	tx, err := bridgeSession.Stake(
		big.NewInt(args.amount),
		vegaPubKeyArr,
	)
	if err != nil {
		return fmt.Errorf("failed to stake asset: %w", err)
	}

	return printEthereumTx(tx)
}
