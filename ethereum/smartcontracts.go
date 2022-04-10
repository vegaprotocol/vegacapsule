package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"code.vegaprotocol.io/vega/bridges/multisig"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient struct {
	client  *ethclient.Client
	chainID int64

	multisig *multisig.MultiSigControl

	vegaBinary string
}

func NewEthereumClient(ctx context.Context, vegaBinary string, chainID int, ethAddress string, smartcontracts types.SmartContractsInfo) (*EthereumClient, error) {
	client, err := ethclient.DialContext(ctx, ethAddress)

	if err != nil {
		return nil, fmt.Errorf("failed creating ehtereum client: %w", err)
	}

	if smartcontracts.MultisigControl.EthereumAddress == "" {
		return nil, fmt.Errorf("failed to create ethereum client: the multisig smart contract address is not set, please uptate it in the network configuration")
	}

	multisigControl, err := multisig.NewMultiSigControl(common.HexToAddress(smartcontracts.MultisigControl.EthereumAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed creating multisig control: %w", err)
	}

	return &EthereumClient{
		client:     client,
		multisig:   multisigControl,
		chainID:    int64(chainID),
		vegaBinary: vegaBinary,
	}, nil
}

func (ec EthereumClient) InitMultisig(ctx context.Context, smartcontracts types.SmartContractsInfo, validators KeyPairList) error {
	if smartcontracts.EthereumOwner.Private == "" || smartcontracts.EthereumOwner.Public == "" {
		return fmt.Errorf("failed to init multisig smart contract: missing private or public key of the smart contract owner in the network configuration")
	}

	ownerKeyPair := KeyPair{
		PrivateKey: smartcontracts.EthereumOwner.Private,
		Address:    smartcontracts.EthereumOwner.Public,
	}

	session, err := ec.createMultiSigControlSession(ctx, ownerKeyPair)
	if err != nil {
		return fmt.Errorf("failed to create multisig smart contract session: %w", err)
	}

	validSigner, err := session.IsValidSigner(common.HexToAddress(smartcontracts.EthereumOwner.Public))
	if err != nil {
		return fmt.Errorf("failed to check signer: %w", err)
	}
	if !validSigner {
		return fmt.Errorf("failed to verify signer: %s is not valid signer of messages", smartcontracts.EthereumOwner.Public)
	}

	if err := ec.multisigSetThreshold(ctx, session, 1, KeyPairList{ownerKeyPair}); err != nil {
		return fmt.Errorf("failed to set multisig threshold to 1: %w", err)
	}

	if err := ec.multisigAddSigners(ctx, session, validators, KeyPairList{ownerKeyPair}); err != nil {
		return fmt.Errorf("failed to add signers: %w", err)
	}

	if err := ec.multisigRemoveSigners(ctx, session, smartcontracts.EthereumOwner.Public, KeyPairList{ownerKeyPair}); err != nil {
		return fmt.Errorf("failed to remove contract owner from muiltisig signer: %w", err)
	}

	if err := ec.multisigSetThreshold(ctx, session, 500, validators); err != nil {
		return fmt.Errorf("failed to set multisig threshold to 500: %w", err)
	}

	return nil
}

func (ec EthereumClient) createMultiSigControlSession(ctx context.Context, ownerKeyPair KeyPair) (*multisig.MultiSigControlSession, error) {
	privateKey, err := crypto.HexToECDSA(ownerKeyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert multisig owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ec.chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create multisig authentication: %w", err)
	}

	session := &multisig.MultiSigControlSession{
		Contract: ec.multisig,
		CallOpts: bind.CallOpts{
			From:    common.HexToAddress(ownerKeyPair.Address),
			Context: ctx,
		},
		TransactOpts: *auth,
	}

	return session, nil
}

func (ec EthereumClient) multisigSetThreshold(ctx context.Context, session *multisig.MultiSigControlSession, newThreshold int, signers KeyPairList) error {
	currentThreshold, err := session.GetCurrentThreshold()
	if err != nil {
		return fmt.Errorf("failed to get current multisig threshold: %w", err)
	}
	log.Printf("Current multisig threshold: %d\n", currentThreshold)

	nonce, err := ec.getSessionNonce(session)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	signature, err := setThresholdSignature(ec.vegaBinary, newThreshold, nonce.Uint64(), session.CallOpts.From.Hex(), signers.PrivateKeys())
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

func (ec EthereumClient) multisigAddSigners(ctx context.Context, session *multisig.MultiSigControlSession, validators KeyPairList, signers KeyPairList) error {
	signersCount, err := session.GetValidSignerCount()
	if err != nil {
		return fmt.Errorf("failed to get number of signers for multisig: %w", err)
	}
	log.Printf("Number of signers for multisig: %d\n", signersCount)

	for _, validator := range validators {
		validSigner, err := session.IsValidSigner(common.HexToAddress(validator.Address))
		if err != nil {
			return fmt.Errorf("failed to check signer: %w", err)
		}

		if validSigner {
			log.Printf("%s is already valid signer. No need to add it again", validator.Address)
			continue
		}

		nonce, err := ec.getSessionNonce(session)
		if err != nil {
			return fmt.Errorf("failed to get nonce: %w", err)
		}
		signature, err := addSignerSignature(ec.vegaBinary, validator.Address, nonce.Uint64(), session.CallOpts.From.Hex(), signers.PrivateKeys())
		if err != nil {
			return fmt.Errorf("failed generate the add_signer signature for %s signer: %w", validator.Address, err)
		}
		log.Printf("Computed signature for add_signer for %s: %s\n", validator.Address, signature)

		tx, err := session.AddSigner(common.HexToAddress(validator.Address), nonce, common.FromHex(signature))
		if err != nil {
			return fmt.Errorf("failed to add %s as a multisig signer: %w", validator.Address, err)
		}
		if _, err := bind.WaitMined(ctx, ec.client, tx); err != nil {
			return fmt.Errorf("failed waiting for transaction to be mined: %w", err)
		}

		log.Printf("Added %s as a multisig validator\n", validator.Address)
	}

	signersCount, err = session.GetValidSignerCount()
	if err != nil {
		return fmt.Errorf("failed to get number of signers for multisig: %w", err)
	}
	log.Printf("Updated number of signers for multisig: %d\n", signersCount)

	return nil
}

func (ec EthereumClient) multisigRemoveSigners(ctx context.Context, session *multisig.MultiSigControlSession, oldSigner string, signers KeyPairList) error {
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
	signature, err := removeSignerSignature(ec.vegaBinary, oldSigner, nonce.Uint64(), session.CallOpts.From.Hex(), signers.PrivateKeys())
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

func (ec EthereumClient) getSessionNonce(session *multisig.MultiSigControlSession) (*big.Int, error) {
	nonce, err := ec.client.PendingNonceAt(context.Background(), session.CallOpts.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get multisig owner nonce: %w", err)
	}
	return big.NewInt(int64(nonce)), nil
}
