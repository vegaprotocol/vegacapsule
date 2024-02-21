package ethereum_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	vgethereum "code.vegaprotocol.io/vegacapsule/libs/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	erc20BridgeAddress   = common.HexToAddress("0x9708FF7510D4A7B9541e1699d15b53Ecb1AFDc54")
	stakingBridgeAddress = common.HexToAddress("0x9135f5afd6F055e731bca2348429482eE614CFfA")
	tUSDCTokenAddress    = common.HexToAddress("0x1b8a1B6CBE5c93609b46D1829Cc7f3Cb8eeE23a0")
	vegaTokenAddress     = common.HexToAddress("0x67175Da1D5e966e40D11c4B2519392B2058373de")
	contractOwnerAddress = common.HexToAddress("0xEe7D375bcB50C26d52E1A4a472D8822A2A22d94F")

	vegaPubKey              = "321b470cddc840854fb33b34674683cae3201c7b3eceffa927f663549f17484f"
	contractOwnerPrivateKey = "a37f4c2a678aefb5037bf415a826df1540b330b7e471aa54184877ba901b9ef0"
)

type token interface {
	Mint(to common.Address, amount *big.Int) (*types.Transaction, error)
	MintSync(to common.Address, amount *big.Int) (*types.Transaction, error)
	BalanceOf(account common.Address) (*big.Int, error)
	ApproveSync(spender common.Address, value *big.Int) (*types.Transaction, error)
	Address() common.Address
}

func mintTokenAndShowBalances(token token, address common.Address, amount *big.Int) error {
	fmt.Println("---- Minting new token")

	balance, err := token.BalanceOf(address)
	if err != nil {
		return fmt.Errorf("failed to get balance for %s: %s", address.String(), err)
	}
	fmt.Printf("Initial balance of %s is %s \n", address, balance)

	fmt.Printf("Minting token %s amount %s for %s \n", token.Address(), amount, address)
	if _, err := token.MintSync(address, amount); err != nil {
		return fmt.Errorf("failed to call Mint contract: %s", err)
	}

	balance, err = token.BalanceOf(address)
	if err != nil {
		return fmt.Errorf("failed to get balance for %s: %s", address.String(), err)
	}

	fmt.Printf("Balance of %s after mint is %s \n", address, balance)

	fmt.Println("---- Token minted")

	return nil
}

func approveAndDepositToken(token token, bridge *vgethereum.ERC20BridgeSession, amount *big.Int, vegaPubKey string) error {
	fmt.Println("---- Deposit token")

	fmt.Printf("Approving token %s amount %s for %s \n", token.Address(), amount, bridge.Address())
	if _, err := token.ApproveSync(bridge.Address(), amount); err != nil {
		return fmt.Errorf("failed to approve token: %w", err)
	}

	fmt.Printf("Depositing asset %s amout %s Vega pub key %s \n", token.Address(), amount, bridge.Address())

	vegaPubKeyByte32, err := vgethereum.HexStringToByte32Array(vegaPubKey)
	if err != nil {
		return err
	}

	if _, err := bridge.DepositAssetSync(token.Address(), amount, vegaPubKeyByte32); err != nil {
		return fmt.Errorf("failed to deposit asset: %w", err)
	}

	fmt.Println("---- Token deposited")

	return nil
}

func approveAndStakeToken(token token, bridge *vgethereum.StakingBridgeSession, amount *big.Int, vegaPubKey string) error {
	fmt.Println("---- Stake token")

	fmt.Printf("Approving token %s amount %s for %s \n", token.Address(), amount, bridge.Address())
	if _, err := token.ApproveSync(bridge.Address(), amount); err != nil {
		return fmt.Errorf("failed to approve token: %w", err)
	}

	vegaPubKeyByte32, err := vgethereum.HexStringToByte32Array(vegaPubKey)
	if err != nil {
		return err
	}

	fmt.Printf("Staking asset %s amout %s Vega pub key %s \n", token.Address(), amount, bridge.Address())
	if _, err := bridge.Stake(amount, vegaPubKeyByte32); err != nil {
		return fmt.Errorf("failed to stake asset: %w", err)
	}

	fmt.Println("---- Token staked")

	return nil
}

func TestClient(t *testing.T) {
	ethereumAddress := os.Getenv("ETHEREUM_CLIENT_ADDRESS")
	if ethereumAddress == "" {
		t.Skip("skipping Ethereum client test, set environment variable ETHEREUM_CLIENT_ADDRESS to run it")
	}

	amount := big.NewInt(1000000000000000000)

	ctx := context.Background()

	client, err := vgethereum.NewClient(ctx, ethereumAddress)
	if err != nil {
		t.Fatalf("Failed to create Ethereum client: %s", err)
	}

	stakingBridge, err := client.NewStakingBridgeSession(ctx, contractOwnerPrivateKey, stakingBridgeAddress, nil)
	if err != nil {
		t.Fatalf("Failed to create staking bridge: %s", err)
	}

	erc20bridge, err := client.NewERC20BridgeSession(ctx, contractOwnerPrivateKey, erc20BridgeAddress, nil)
	if err != nil {
		t.Fatalf("Failed to create staking bridge: %s", err)
	}

	tUSDCToken, err := client.NewBaseTokenSession(ctx, contractOwnerPrivateKey, tUSDCTokenAddress, nil)
	if err != nil {
		t.Fatalf("Failed to create tUSDC token: %s", err)
	}

	vegaToken, err := client.NewBaseTokenSession(ctx, contractOwnerPrivateKey, vegaTokenAddress, nil)
	if err != nil {
		t.Fatalf("Failed to create vega token: %s", err)
	}

	if err := mintTokenAndShowBalances(tUSDCToken, contractOwnerAddress, amount); err != nil {
		t.Fatalf("Failed to mint and show balances for tUSDCToken: %s", err)
	}

	if err := mintTokenAndShowBalances(vegaToken, contractOwnerAddress, amount); err != nil {
		t.Fatalf("Failed to mint and show balances for vegaToken: %s", err)
	}

	if err := approveAndDepositToken(tUSDCToken, erc20bridge, amount, vegaPubKey); err != nil {
		t.Fatalf("Failed to approve and deposit token on erc20 bridge: %s", err)
	}

	if err := approveAndStakeToken(vegaToken, stakingBridge, amount, vegaPubKey); err != nil {
		t.Fatalf("Failed to approve and stake token on staking bridge: %s", err)
	}

	fmt.Println("Done")
}
