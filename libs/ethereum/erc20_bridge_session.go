package ethereum

import (
	"fmt"
	"math/big"
	"time"

	"code.vegaprotocol.io/vegacapsule/libs/ethereum/generated"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20BridgeSession struct {
	generated.ERC20BridgeSession
	syncTimeout time.Duration
	address     common.Address
}

func (bs ERC20BridgeSession) Address() common.Address {
	return bs.address
}

func (bs ERC20BridgeSession) DepositAssetSync(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	sink := make(chan *generated.ERC20BridgeAssetDeposited)

	sub, err := bs.Contract.WatchAssetDeposited(&bind.WatchOpts{}, sink, []common.Address{}, []common.Address{asset_source})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for deposit: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := bs.DepositAsset(asset_source, amount, vega_public_key)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, bs.syncTimeout)
}
