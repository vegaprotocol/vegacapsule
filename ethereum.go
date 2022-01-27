package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

// TODO this files will probably not be needed

const vegaNodeEthereumMnemonic = "sentence find kit hood will omit awake prize leave bid nation crawl"

func generateEthereumWallet(walletPath string, privateKeyHex string, passphrase string) error {
	ks := keystore.NewKeyStore(walletPath, keystore.StandardScryptN, keystore.StandardScryptP)

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return err
	}

	_, err = ks.ImportECDSA(privateKey, passphrase)
	if err != nil {
		return err
	}

	return nil
}

type keyPair struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

type ethereumKeys struct {
	Master keyPair
	Keys   []keyPair
}

// TODO: change map to struct?
// Only generate master key and derive rest with it in special configs?
func generateEthereumKeys(count int) (*ethereumKeys, error) {
	seed := bip39.NewSeed(vegaNodeEthereumMnemonic, "")

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	mPrivKey, err := masterKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	mPubKey, err := masterKey.ECPubKey()
	if err != nil {
		return nil, err
	}

	ethKeys := &ethereumKeys{
		Master: keyPair{
			Private: crypto.PubkeyToAddress(*mPubKey.ToECDSA()).Hex(),
			Public:  hexutil.Encode(mPrivKey.Serialize())[2:],
		},
		Keys: make([]keyPair, 0, count),
	}

	// this could be method (like: AddNode) instead and can be generating key on demand... instead of upfront
	for i := 0; i < count; i++ {
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		account, err := wallet.Derive(path, false)
		if err != nil {
			return nil, err
		}

		privateKey, err := wallet.PrivateKeyHex(account)
		if err != nil {
			return nil, err
		}

		ethKeys.Keys = append(ethKeys.Keys, keyPair{
			Private: privateKey,
			Public:  account.Address.Hex(),
		})
	}

	return ethKeys, nil
}

type walletMnemonic struct {
	Seed     string `json:"seed"`
	Mnemonic string `json:"mnemonic"`
}

func generateWalletMnemonics(count int, network string) ([]walletMnemonic, error) {
	mnemonics := make([]walletMnemonic, 0, count)

	for i := 0; i < count; i++ {
		entropy := fmt.Sprintf("network %s, node %03d", network, i)
		h := sha256.New()
		h.Write([]byte(entropy))
		entropyHash := h.Sum(nil)

		mnemonic, err := bip39.NewMnemonic(entropyHash)
		if err != nil {
			return nil, err
		}

		mnemonics = append(mnemonics, walletMnemonic{
			Seed:     hex.EncodeToString(entropyHash),
			Mnemonic: mnemonic,
		})
	}

	return mnemonics, nil
}
