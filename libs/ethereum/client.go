package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"code.vegaprotocol.io/vegacapsule/libs/ethereum/generated"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

var defaultSyncDuration = time.Second * 5

type Client struct {
	*ethclient.Client
	chainID *big.Int
}

func NewClient(ctx context.Context, ethereumAddress string) (*Client, error) {
	addr, err := url.Parse(ethereumAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ethereum address: %w", err)
	}

	if addr.Scheme != "ws" && addr.Scheme != "wss" {
		return nil, fmt.Errorf("address scheme needs to be 'ws' or 'wss': %q", addr.Scheme)
	}

	client, err := ethclient.DialContext(ctx, addr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to dial Ethereum client: %s", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return &Client{
		chainID: chainID,
		Client:  client,
	}, nil
}

func (ec *Client) NewERC20BridgeSession(
	ctx context.Context,
	contractOwnerPrivateKey string,
	bridgeAddress common.Address,
	syncTimeout *time.Duration,
) (*ERC20BridgeSession, error) {
	privateKey, err := crypto.HexToECDSA(contractOwnerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert erc20 bridge contract owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ec.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create erc20 bridge contract authentication: %w", err)
	}

	bridge, err := generated.NewERC20Bridge(bridgeAddress, ec.Client)
	if err != nil {
		return nil, fmt.Errorf("failed creating erc20 bridge contract for address %q: %w", bridgeAddress, err)
	}

	if syncTimeout == nil {
		syncTimeout = &defaultSyncDuration
	}

	return &ERC20BridgeSession{
		ERC20BridgeSession: generated.ERC20BridgeSession{
			Contract: bridge,
			CallOpts: bind.CallOpts{
				From:    auth.From,
				Context: ctx,
			},
			TransactOpts: *auth,
		},
		syncTimeout: *syncTimeout,
		address:     bridgeAddress,
	}, nil
}

func (ec *Client) NewStakingBridgeSession(
	ctx context.Context,
	contractOwnerPrivateKey string,
	bridgeAddress common.Address,
	syncTimeout *time.Duration,
) (*StakingBridgeSession, error) {
	privateKey, err := crypto.HexToECDSA(contractOwnerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert staking bridge contract owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ec.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create staking bridge contract authentication: %w", err)
	}

	bridge, err := generated.NewStakingBridge(bridgeAddress, ec.Client)
	if err != nil {
		return nil, fmt.Errorf("failed creating staking bridge contract for address %q: %w", bridgeAddress, err)
	}

	if syncTimeout == nil {
		syncTimeout = &defaultSyncDuration
	}

	return &StakingBridgeSession{
		StakingBridgeSession: generated.StakingBridgeSession{
			Contract: bridge,
			CallOpts: bind.CallOpts{
				From:    auth.From,
				Context: ctx,
			},
			TransactOpts: *auth,
		},
		syncTimeout: *syncTimeout,
		address:     bridgeAddress,
	}, nil
}

func (ec *Client) NewBaseTokenSession(
	ctx context.Context,
	contractOwnerPrivateKey string,
	tokenAddress common.Address,
	syncTimeout *time.Duration,
) (*BaseTokenSession, error) {
	privateKey, err := crypto.HexToECDSA(contractOwnerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base token contract owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ec.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create base token contract authentication: %w", err)
	}

	token, err := generated.NewBaseToken(tokenAddress, ec.Client)
	if err != nil {
		return nil, fmt.Errorf("failed creating base token contract for address %q: %w", tokenAddress, err)
	}

	if syncTimeout == nil {
		syncTimeout = &defaultSyncDuration
	}

	return &BaseTokenSession{
		BaseTokenSession: generated.BaseTokenSession{
			Contract: token,
			CallOpts: bind.CallOpts{
				From:    auth.From,
				Context: ctx,
			},
			TransactOpts: *auth,
		},
		syncTimeout: *syncTimeout,
		address:     tokenAddress,
		privateKey:  privateKey,
		client:      ec,
	}, nil
}

func wait[T any](sink chan T, sub event.Subscription, tx *types.Transaction, timeout time.Duration) (*types.Transaction, error) {
	select {
	case <-sink:
		return tx, nil
	case err := <-sub.Err():
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("transaction time has timed out")
	}
}

func HexStringToByte32Array(str string) ([32]byte, error) {
	value := [32]byte{}

	decoded, err := hex.DecodeString(str)
	if err != nil {
		return value, err
	}

	copy(value[:], decoded)

	return value, nil
}
