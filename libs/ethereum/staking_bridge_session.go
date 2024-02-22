package ethereum

import (
	"time"

	"code.vegaprotocol.io/vegacapsule/libs/ethereum/generated"

	"github.com/ethereum/go-ethereum/common"
)

type StakingBridgeSession struct {
	generated.StakingBridgeSession
	syncTimeout time.Duration
	address     common.Address
}

func (ss StakingBridgeSession) Address() common.Address {
	return ss.address
}
