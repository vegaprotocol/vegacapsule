package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"code.vegaprotocol.io/vegacapsule/libs/ethereum/generated"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type BaseTokenSession struct {
	generated.BaseTokenSession
	syncTimeout time.Duration
	address     common.Address
	privateKey  *ecdsa.PrivateKey
	client      *Client
}

func (ts *BaseTokenSession) Address() common.Address {
	return ts.address
}

func (ts *BaseTokenSession) ApproveSync(spender common.Address, value *big.Int) (*types.Transaction, error) {
	sink := make(chan *generated.BaseTokenApproval)

	sub, err := ts.Contract.WatchApproval(&bind.WatchOpts{}, sink, []common.Address{ts.CallOpts.From}, []common.Address{spender})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for approval: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.Approve(spender, value)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, ts.syncTimeout)
}

func (ts *BaseTokenSession) TransferSync(recipient common.Address, value *big.Int) (*types.Transaction, error) {
	sink := make(chan *generated.BaseTokenTransfer)

	sub, err := ts.Contract.WatchTransfer(&bind.WatchOpts{}, sink, []common.Address{ts.CallOpts.From}, []common.Address{recipient})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for transfer: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.Transfer(recipient, value)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, ts.syncTimeout)
}

func (ts *BaseTokenSession) TransferFromSync(sender common.Address, recipient common.Address, value *big.Int) (*types.Transaction, error) {
	sink := make(chan *generated.BaseTokenTransfer)

	sub, err := ts.Contract.WatchTransfer(&bind.WatchOpts{}, sink, []common.Address{sender}, []common.Address{recipient})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for transfer: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.TransferFrom(sender, recipient, value)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, ts.syncTimeout)
}

func (ts *BaseTokenSession) MintSync(to common.Address, amount *big.Int) (*types.Transaction, error) {
	sink := make(chan *generated.BaseTokenTransfer)

	sub, err := ts.Contract.WatchTransfer(&bind.WatchOpts{}, sink, []common.Address{common.BigToAddress(common.Big0)}, []common.Address{to})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for transfer: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.Mint(to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to mint %s: %w", to, err)
	}

	tx, err = wait(sink, sub, tx, ts.syncTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for mint: %w", err)
	}

	return tx, nil
}

// MintRawSync is an experimental way of minting new tokens. It attempts to execute an on-chain transaction that
// runs a Yul script which loops over the "faucet" method until either of the following happens:
//  1. the target balance is reached,
//  2. the caller runs out of gas.
//
// Use this method only if the "mint" contract method is not available, or a non-owner account is the caller.
// The execution time is very slow so the (context) timeout for it should be higher than for other methods.
func (ts *BaseTokenSession) MintRawSync(ctx context.Context, to common.Address, amount *big.Int) (*big.Int, error) {
	tx, err := ts.MintRaw(to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to mint raw: %w", err)
	}

	minted, err := ts.GetLastTransferValueSync(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for last value: %w", err)
	}

	return minted, nil
}

// MintRaw is an experimental way of minting new tokens. It attempts to execute an on-chain transaction that
// runs a Yul script which loops over the "faucet" method until either of the following happens:
//  1. the target balance is reached,
//  2. the caller runs out of gas.
//
// Use this method only if the "mint" contract method is not available, or a non-owner account is the caller.
func (ts *BaseTokenSession) MintRaw(to common.Address, amount *big.Int) (*types.Transaction, error) {
	if amount == nil || amount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("amount must be not be nil and greater than 0")
	}

	balance, err := ts.BalanceOf(to)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	targetBalance := new(big.Int).Add(balance, amount)

	tx, err := ts.mintRaw(to, targetBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to mint: %w", err)
	}

	return tx, nil
}

func (ts *BaseTokenSession) mintRaw(to common.Address, targetBalance *big.Int) (*types.Transaction, error) {
	signedTx, err := ts.createMintSignedTx(to, targetBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed mint tx: %w", err)
	}

	if err = ts.client.SendTransaction(ts.CallOpts.Context, signedTx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

func (ts *BaseTokenSession) GetLastTransferValueSync(ctx context.Context, signedTx *types.Transaction) (*big.Int, error) {
	receipt, err := bind.WaitMined(ctx, ts.client, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	if len(receipt.Logs) == 0 {
		return nil, fmt.Errorf("no logs in transaction")
	}

	// last transfer event
	log := receipt.Logs[len(receipt.Logs)-1]

	transfer, err := ts.Contract.ParseTransfer(*log)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transfer: %w", err)
	}

	return transfer.Value, nil
}

func (ts *BaseTokenSession) createMintSignedTx(to common.Address, targetBalance *big.Int) (*types.Transaction, error) {
	nonce, err := ts.client.PendingNonceAt(ts.CallOpts.Context, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := ts.client.SuggestGasPrice(ts.CallOpts.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// use maximum gas limit of 30m
	gasLimit := new(big.Int).Mul(big.NewInt(30), new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))

	data, err := ts.prepareMintBytecode(targetBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare data: %w", err)
	}

	tx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit.Uint64(), gasPrice, data)

	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, ts.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

func (ts *BaseTokenSession) prepareMintBytecode(targetBalance *big.Int) ([]byte, error) {
	const script = `0x6000603460be8239805160601c9060145160183384609f565b908181116043575b82806044818088602f3082609f565b5063a9059cbb60e01b600052336004525af1005b63de5f72fd60e01b600052600492918280858180895af15060633086609f565b63de5f72fd60e01b60005291030491815b8381105a61ea6010161560925760019083808481808a5af150016074565b5090915060449050816020565b6024600081926044946370a0823160e01b83526004525afa506024519056`

	hexCode, err := hexutil.Decode(script)
	if err != nil {
		return nil, fmt.Errorf("failed to decode script: %w", err)
	}

	var data []byte
	data = append(data, hexCode...)
	data = append(data, ts.address.Bytes()...)
	data = append(data, common.LeftPadBytes(targetBalance.Bytes(), 32)...)
	return data, nil
}
