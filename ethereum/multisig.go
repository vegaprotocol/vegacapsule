package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"

	multisig "code.vegaprotocol.io/vega/core/contracts/multisig_control"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumMultisigClient struct {
	client  *ethclient.Client
	chainID int64

	multisig *multisig.MultisigControl

	vegaBinary string
	vegaHome   string
}

func WaitForNetwork(ctx context.Context, chainID int, ethAddress string) error {
	done := make(chan bool)
	go func(chainID int, ethAddress string) {
		for {
			if client, err := ethclient.DialContext(ctx, ethAddress); err == nil {
				client.Close()
				done <- true
			}
		}
	}(chainID, ethAddress)

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type EthereumMultisigClientParameters struct {
	ChainID            int
	EthereumAddress    string
	SmartcontractsInfo types.SmartContractsInfo

	VegaBinary string
	VegaHome   string
}

func NewEthereumMultisigClient(ctx context.Context, params EthereumMultisigClientParameters) (*EthereumMultisigClient, error) {
	client, err := ethclient.DialContext(ctx, params.EthereumAddress)

	if err != nil {
		return nil, fmt.Errorf("failed creating ehtereum client: %w", err)
	}

	if params.SmartcontractsInfo.MultisigControl.EthereumAddress == "" {
		return nil, fmt.Errorf("failed to create ethereum client: the multisig smart contract address is not set, please uptate it in the network configuration")
	}

	multisigControl, err := multisig.NewMultisigControl(common.HexToAddress(params.SmartcontractsInfo.MultisigControl.EthereumAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed creating multisig control: %w", err)
	}

	return &EthereumMultisigClient{
		client:     client,
		multisig:   multisigControl,
		chainID:    int64(params.ChainID),
		vegaBinary: params.VegaBinary,
		vegaHome:   params.VegaHome,
	}, nil
}

func (ec EthereumMultisigClient) InitMultisig(ctx context.Context, smartcontracts types.SmartContractsInfo, validators SignersList) error {
	if smartcontracts.EthereumOwner.Private == "" || smartcontracts.EthereumOwner.Public == "" {
		return fmt.Errorf("failed to init multisig smart contract: missing private or public key of the smart contract owner in the network configuration")
	}

	if len(validators) == 0 {
		return fmt.Errorf("failed to init multisig smart contract: can not run multisig contract with no validators")
	}

	contractOwner := Signer{
		HomeAddress: validators[0].HomeAddress,
		KeyPair: KeyPair{
			PrivateKey: smartcontracts.EthereumOwner.Private,
			Address:    smartcontracts.EthereumOwner.Public,
		},
	}

	session, err := ec.createMultiSigControlSession(ctx, contractOwner)
	if err != nil {
		return fmt.Errorf("failed to create multisig smart contract session: %w", err)
	}

	validSigner, err := session.IsValidSigner(common.HexToAddress(contractOwner.KeyPair.Address))
	if err != nil {
		return fmt.Errorf("failed to check signer: %w", err)
	}
	if !validSigner {
		return fmt.Errorf("failed to verify signer: %s is not valid signer of messages", contractOwner.KeyPair.Address)
	}

	if err := ec.multisigSetThreshold(ctx, session, 1, SignersList{contractOwner}); err != nil {
		return fmt.Errorf("failed to set multisig threshold to 1: %w", err)
	}

	if err := ec.multisigAddSigners(ctx, session, validators, SignersList{contractOwner}); err != nil {
		return fmt.Errorf("failed to add signers: %w", err)
	}

	if err := ec.multisigRemoveSigner(ctx, session, contractOwner.KeyPair.Address, SignersList{contractOwner}); err != nil {
		return fmt.Errorf("failed to remove contract owner from muiltisig signer: %w", err)
	}

	if err := ec.multisigSetThreshold(ctx, session, 667, validators); err != nil {
		return fmt.Errorf("failed to set multisig threshold to 667: %w", err)
	}

	return nil
}

func (ec EthereumMultisigClient) createMultiSigControlSession(ctx context.Context, ownerKeyPair Signer) (*multisig.MultisigControlSession, error) {
	privateKey, err := crypto.HexToECDSA(ownerKeyPair.KeyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert multisig owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ec.chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create multisig authentication: %w", err)
	}

	session := &multisig.MultisigControlSession{
		Contract: ec.multisig,
		CallOpts: bind.CallOpts{
			From:    common.HexToAddress(ownerKeyPair.KeyPair.Address),
			Context: ctx,
		},
		TransactOpts: *auth,
	}

	return session, nil
}

func (ec EthereumMultisigClient) multisigSetThreshold(ctx context.Context, session *multisig.MultisigControlSession, newThreshold int, signers SignersList) error {
	currentThreshold, err := session.GetCurrentThreshold()
	if err != nil {
		return fmt.Errorf("failed to get current multisig threshold: %w", err)
	}
	log.Printf("Current multisig threshold: %d\n", currentThreshold)

	nonce, err := ec.getSessionNonce(session)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	signature, err := setThresholdSignature(
		ec.vegaBinary,
		newThreshold,
		nonce.Uint64(),
		session.CallOpts.From.Hex(),
		signers,
	)
	if err != nil {
		return fmt.Errorf("failed computing signature: %w", err)
	}

	log.Printf("Computed signature for set_threshold to %d: %s\n", newThreshold, signature)
	tx, err := session.SetThreshold(uint16(newThreshold), nonce, common.FromHex(signature))
	if err != nil {
		return fmt.Errorf("failed setting threshold: %w", err)
	}
	if _, err := bind.WaitMined(ctx, ec.client, tx); err != nil {
		return fmt.Errorf("failed waiting for transaction to be mined: %w", err)
	}

	currentThreshold, err = session.GetCurrentThreshold()
	if err != nil {
		return fmt.Errorf("failed to get current multisig threshold: %w", err)
	}
	log.Printf("Updated multisig threshold: %d\n", currentThreshold)

	return nil
}

func (ec EthereumMultisigClient) multisigAddSigners(ctx context.Context, session *multisig.MultisigControlSession, validators SignersList, signers SignersList) error {
	signersCount, err := session.GetValidSignerCount()
	if err != nil {
		return fmt.Errorf("failed to get number of signers for multisig: %w", err)
	}
	log.Printf("Number of signers for multisig: %d\n", signersCount)

	for _, validator := range validators {
		validSigner, err := session.IsValidSigner(common.HexToAddress(validator.KeyPair.Address))
		if err != nil {
			return fmt.Errorf("failed to check signer: %w", err)
		}

		if validSigner {
			log.Printf("%s is already valid signer. No need to add it again", validator.KeyPair.Address)
			continue
		}

		nonce, err := ec.getSessionNonce(session)
		if err != nil {
			return fmt.Errorf("failed to get nonce: %w", err)
		}
		signature, err := addSignerSignature(
			ec.vegaBinary,
			validator.KeyPair.Address,
			nonce.Uint64(),
			session.CallOpts.From.Hex(),
			signers)
		if err != nil {
			return fmt.Errorf("failed generate the add_signer signature for %s signer: %w", validator.KeyPair.Address, err)
		}
		log.Printf("Computed signature for add_signer for %s: %s\n", validator.KeyPair.Address, signature)

		tx, err := session.AddSigner(common.HexToAddress(validator.KeyPair.Address), nonce, common.FromHex(signature))
		if err != nil {
			return fmt.Errorf("failed to add %s as a multisig signer: %w", validator.KeyPair.Address, err)
		}
		if _, err := bind.WaitMined(ctx, ec.client, tx); err != nil {
			return fmt.Errorf("failed waiting for transaction to be mined: %w", err)
		}

		log.Printf("Added %s as a multisig validator\n", validator.KeyPair.Address)
	}

	signersCount, err = session.GetValidSignerCount()
	if err != nil {
		return fmt.Errorf("failed to get number of signers for multisig: %w", err)
	}
	log.Printf("Updated number of signers for multisig: %d\n", signersCount)

	return nil
}

func (ec EthereumMultisigClient) multisigRemoveSigner(ctx context.Context, session *multisig.MultisigControlSession, oldSigner string, signers SignersList) error {
	validSigner, err := session.IsValidSigner(common.HexToAddress(oldSigner))
	if err != nil {
		return fmt.Errorf("failed to check signer: %w", err)
	}

	if !validSigner {
		return fmt.Errorf("failed to remove signer: %s is not valid signer", oldSigner)
	}

	nonce, err := ec.getSessionNonce(session)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	signature, err := removeSignerSignature(
		ec.vegaBinary,
		oldSigner,
		nonce.Uint64(),
		session.CallOpts.From.Hex(),
		signers)
	if err != nil {
		return fmt.Errorf("failed generate signature: %w", err)
	}
	log.Printf("Computed signature for remove_signer for %s: %s\n", oldSigner, signature)

	tx, err := session.RemoveSigner(common.HexToAddress(oldSigner), nonce, common.FromHex(signature))
	if err != nil {
		return fmt.Errorf("failed to remove signer from multisig control: %w", err)
	}

	if _, err := bind.WaitMined(ctx, ec.client, tx); err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}
	log.Printf("Removed the %s signer from multisig control", oldSigner)

	signersCount, err := session.GetValidSignerCount()
	if err != nil {
		return fmt.Errorf("failed to get number of signers for multisig: %w", err)
	}
	log.Printf("Updated number of signers for multisig: %d\n", signersCount)
	return nil
}

func (ec EthereumMultisigClient) getSessionNonce(session *multisig.MultisigControlSession) (*big.Int, error) {
	nonce, err := ec.client.PendingNonceAt(context.Background(), session.CallOpts.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get multisig owner nonce: %w", err)
	}
	return big.NewInt(int64(nonce)), nil
}
