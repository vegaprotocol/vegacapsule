package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vega/config"
	"code.vegaprotocol.io/vega/config/encoding"
	"code.vegaprotocol.io/vega/nodewallets"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

const vegaNodeEthereumMnemonic = "sentence find kit hood will omit awake prize leave bid nation crawl"

type keyPair struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

// TODO: change map to struct?
// Only generate master key and derive rest with it in special configs?
func generateEthereumKeys(count int) (map[string]keyPair, error) {
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

	keyPairs := map[string]keyPair{
		"master": {
			Private: crypto.PubkeyToAddress(*mPubKey.ToECDSA()).Hex(),
			Public:  hexutil.Encode(mPrivKey.Serialize())[2:],
		},
	}

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

		keyPairs[strconv.Itoa(i)] = keyPair{
			Private: privateKey,
			Public:  account.Address.Hex(),
		}
	}

	return keyPairs, nil
}

// TODO change map to struct
func generateWalletMnemonics(count int, network string) (map[string]map[string]string, error) {
	mnemonics := make(map[string]map[string]string, count)

	for i := 0; i < count; i++ {
		entropy := fmt.Sprintf("network %s, node %03d", network, i)
		h := sha256.New()
		h.Write([]byte(entropy))
		entropyHash := h.Sum(nil)

		mnemonic, err := bip39.NewMnemonic(entropyHash)
		if err != nil {
			return nil, err
		}

		mnemonics[strconv.Itoa(i)] = map[string]string{
			"seed":     hex.EncodeToString(entropyHash),
			"mnemonic": mnemonic,
		}
	}

	return mnemonics, nil
}

type VegaConfigPaths struct {
	ConfigFilePath           string
	NodeWalletConfigFilePath string
}

// copied from Vega core
func initVegaConfig(modeS, dir, pass string) (*VegaConfigPaths, error) {
	mode, err := encoding.NodeModeFromString(modeS)
	if err != nil {
		return nil, err
	}

	vegaPaths := paths.New(dir)

	// a nodewallet will be required only for a validator node
	var nwRegistry *nodewallets.RegistryLoader
	if mode == encoding.NodeModeValidator {
		nwRegistry, err = nodewallets.NewRegistryLoader(vegaPaths, pass)
		if err != nil {
			return nil, err
		}
	}

	cfgLoader, err := config.InitialiseLoader(vegaPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise configuration loader: %w", err)
	}

	configExists, err := cfgLoader.ConfigExists()
	if err != nil {
		return nil, fmt.Errorf("couldn't verify configuration presence: %w", err)
	}

	if configExists {
		cfgLoader.Remove()
	}

	cfg := config.NewDefaultConfig()
	cfg.NodeMode = mode

	if err := cfgLoader.Save(&cfg); err != nil {
		return nil, fmt.Errorf("couldn't save configuration file: %w", err)
	}

	return &VegaConfigPaths{
		ConfigFilePath:           cfgLoader.ConfigFilePath(),
		NodeWalletConfigFilePath: nwRegistry.RegistryFilePath(),
	}, nil
}

func initateVegaNode(vegaDir string, id int) error {
	nodeDir := path.Join(vegaDir, fmt.Sprintf("node%d", id))

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return err
	}

	paths, err := initVegaConfig("validator", nodeDir, "n0d3w4ll3t-p4ssphr4e3")
	if err != nil {
		return err
	}

	fmt.Printf("vega confi initialised for node id %q, paths: %#v", id, paths)

	// nodeType :="validator"
	// nodeMode := "validator"

	return nil
}

func generateVegaConfig(vegaDir string) error {
	return initateVegaNode(vegaDir, 0)

	// generateEthereumKeys(3)
	// out, err := generateWalletMnemonics(3, "DV")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := json.Marshal(out)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	return nil
}
