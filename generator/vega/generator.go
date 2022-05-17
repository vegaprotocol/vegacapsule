package vega

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/ethereum"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.VegaNodePrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (vg ConfigGenerator) Initiate(index int, mode, tendermintHome, nodeWalletPass, vegaWalletPass, ethereumWalletPass string) (*types.VegaNode, error) {
	nodeDir := vg.nodeDir(index)

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	nodeWalletPassFilePath := path.Join(nodeDir, "node-vega-wallet-pass.txt")

	if err := ioutil.WriteFile(nodeWalletPassFilePath, []byte(nodeWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write node wallet passphrase to file: %w", err)
	}

	initOut, err := vg.initiateNode(nodeDir, nodeWalletPassFilePath, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate vega node: %w", err)
	}

	confFilePath := vg.configFilePath(nodeDir)
	origConFilePath := vg.originalConfigFilePath(nodeDir)

	if err := utils.CopyFile(confFilePath, origConFilePath); err != nil {
		return nil, fmt.Errorf("failed to copy initiated config from %q to %q: %w", confFilePath, origConFilePath, err)
	}

	initNode := &types.VegaNode{
		Name:                   fmt.Sprintf("vega-%s-%d", mode, index),
		Mode:                   mode,
		HomeDir:                nodeDir,
		NodeWalletPassFilePath: nodeWalletPassFilePath,
		BinaryPath:             *vg.conf.VegaBinary,
	}

	if mode != types.NodeModeValidator {
		log.Printf("vega config initialized for node %q with id %d, paths: %#v", mode, index, initOut.ConfigFilePath)
		return initNode, nil
	}

	nodeWalletInfo, err := vg.initiateValidatorWallets(nodeDir, tendermintHome, vegaWalletPass, ethereumWalletPass, nodeWalletPassFilePath)
	if err != nil {
		return nil, err
	}
	initNode.NodeWalletInfo = nodeWalletInfo

	log.Printf("vega config initialized for node %q with id %d, paths: %#v", mode, index, initOut.ConfigFilePath)

	return initNode, nil
}

func (vg ConfigGenerator) initiateValidatorWallets(nodeDir, tendermintHome, vegaWalletPass, ethereumWalletPass, nodeWalletPassFilePath string) (*types.NodeWalletInfo, error) {
	walletPassFilePath := path.Join(nodeDir, "vega-wallet-pass.txt")
	ethereumPassFilePath := path.Join(nodeDir, "ethereum-vega-wallet-pass.txt")

	if err := ioutil.WriteFile(walletPassFilePath, []byte(vegaWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write wallet passphrase to file: %w", err)
	}

	if err := ioutil.WriteFile(ethereumPassFilePath, []byte(ethereumWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write ethereum wallet passphrase to file: %w", err)
	}

	vegaOut, err := vg.createWallet(nodeDir, "created-wallet", walletPassFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create vega wallet: %w", err)
	}
	log.Printf("node wallet create out: %#v", vegaOut)

	vegaImportOut, err := vg.importVegaNodeWallet(nodeDir, nodeWalletPassFilePath, walletPassFilePath, vegaOut.Wallet.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to import vega wallet: %w", err)
	}
	log.Printf("node wallet import out: %#v", vegaImportOut)

	ethOut, err := vg.generateNodeWallet(nodeDir, nodeWalletPassFilePath, ethereumPassFilePath, types.NodeWalletChainTypeEthereum)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %q wallet: %w", types.NodeWalletChainTypeEthereum, err)
	}
	log.Printf("ethereum wallet out: %#v", ethOut)

	ethKey, err := ethereum.DescribeKeyPair(ethOut.WalletFilePath, ethereumWalletPass)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain thereum address for the wallet '%s': %w", ethOut.WalletFilePath, err)
	}

	tmtOut, err := vg.importTendermintNodeWallet(nodeDir, nodeWalletPassFilePath, tendermintHome)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tenderming wallet: %w", err)
	}

	log.Printf("tendermint wallet out: %#v", tmtOut)

	return &types.NodeWalletInfo{
		VegaWalletRecoveryPhrase: vegaOut.Wallet.RecoveryPhrase,
		VegaWalletPublicKey:      vegaOut.Key.Public,
		EthereumAddress:          ethKey.Address,
		EthereumPrivateKey:       ethKey.PrivateKey,
	}, nil
}

func (vg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *vg.conf.NodeDirPrefix, i)
	return filepath.Join(vg.homeDir, nodeDirName)
}

func (vg ConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node", "config.toml")
}

func (vg ConfigGenerator) originalConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node", "config-original.toml")
}
