package cmd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	vgethereum "code.vegaprotocol.io/shared/libs/ethereum"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

type ethereumDepositOptions struct {
	netAddress      string
	ownerPrivateKey string
	bridgeAddress   string
	asset           string
	assetAddress    string
	vegaPubKey      string
	amount          int64
}

var (
	ethereumDepositOpts = ethereumDepositOptions{}
)

func init() {
	// ethereumDepositCmd.Flags().StringVar(&ethereumDepositOpts.ownerPrivateKey, "owner-private-key", "", "private key of the bridge contract owner")
	// ethereumDepositCmd.MarkFlagRequired("owner-private-key")
	// ethereumDepositCmd.Flags().StringVar(&ethereumDepositOpts.bridgeAddress, "bridge-addr", "", "smart contract address of the bridge")
	// ethereumDepositCmd.MarkFlagRequired("bridge-addr")

	ethereumDepositCmd.Flags().StringVar(&ethereumDepositOpts.asset, "asset-addr", "", "address of the asset to be deposited")
	ethereumDepositCmd.Flags().StringVar(&ethereumDepositOpts.vegaPubKey, "pub-key", "", "Vega public key to where the asset will be deposited")
	ethereumDepositCmd.Flags().Int64Var(&ethereumDepositOpts.amount, "amount", 0, "amount to be deposited")
	ethereumDepositCmd.MarkFlagRequired("asset-addr")
	ethereumDepositCmd.MarkFlagRequired("pub-key")
	ethereumDepositCmd.MarkFlagRequired("amount")

	ethereumCmd.AddCommand(ethereumDepositCmd)
}

var ethereumDepositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit allows to deposit an asset to given Vega public key.",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("ethereum multisig init")
		}

		if !netState.Running() {
			return networkNotRunningErr("ethereum multisig init")
		}

		conf := netState.Config

		smartContracts, err := conf.SmartContractsInfo()
		if err != nil {
			return fmt.Errorf("failed getting smart contract informations: %w", err)
		}

		ethereumDepositOpts.netAddress = conf.Network.Ethereum.Endpoint
		ethereumDepositOpts.ownerPrivateKey = smartContracts.EthereumOwner.Private

		return ethereumDeposit(cmd.Context(), ethereumDepositOpts)
	},
}

func ethereumDeposit(ctx context.Context, ethereumOpts ethereumDepositOptions) error {
	defeaultSyncTimeout := time.Second * defaultEthreumWaitTimeout

	client, err := vgethereum.NewClient(ctx, ethereumOpts.netAddress)
	if err != nil {
		return fmt.Errorf("falied to create Ethereum client: %w", err)
	}

	bridgeAddr := common.HexToAddress(ethereumDepositOpts.bridgeAddress)

	bridgeSession, err := client.NewERC20BridgeSession(
		ctx,
		ethereumDepositOpts.ownerPrivateKey,
		bridgeAddr,
		&defeaultSyncTimeout,
	)
	if err != nil {
		return fmt.Errorf("failed to create erc20 bridge session for %s: %w", ethereumDepositOpts.bridgeAddress, err)
	}

	tokenSession, err := client.NewBaseTokenSession(
		ctx,
		ethereumDepositOpts.ownerPrivateKey,
		common.HexToAddress(ethereumDepositOpts.assetAddress),
		&defeaultSyncTimeout,
	)
	if err != nil {
		return fmt.Errorf("failed to create token session for %s: %w", ethereumDepositOpts.assetAddress, err)
	}

	amount := big.NewInt(ethereumDepositOpts.amount)

	if _, err := tokenSession.ApproveSync(bridgeAddr, amount); err != nil {
		return fmt.Errorf("failed to approve asset amount to bridge: %w", err)
	}

	vegaPubKeyArr, err := vgethereum.HexStringToByte32Array(ethereumDepositOpts.vegaPubKey)
	if err != nil {
		return fmt.Errorf("failed to convert Vega pub key string to byte array: %w", err)
	}

	tx, err := bridgeSession.DepositAssetSync(
		common.HexToAddress(ethereumDepositOpts.assetAddress),
		amount,
		vegaPubKeyArr,
	)
	if err != nil {
		return fmt.Errorf("failed to deposit asset: %w", err)
	}

	return printEthereumTx(tx)
}
