package vega

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/ethereum"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

const defaultIsolatedWalletName = "created-wallet"

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

func (vg ConfigGenerator) Initiate(
	index int,
	optVegaBinary *string,
	mode, tendermintHome, nodeWalletPass,
	vegaWalletPass, ethereumWalletPass string,
	clefConf *config.ClefConfig,
) (*types.VegaNode, error) {
	nodeDir := vg.nodeDir(index)

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	nodeWalletPassFilePath := path.Join(nodeDir, "node-vega-wallet-pass.txt")

	if err := os.WriteFile(nodeWalletPassFilePath, []byte(nodeWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write node wallet passphrase to file: %w", err)
	}

	vegaBinary := *vg.conf.VegaBinary
	if optVegaBinary != nil {
		vegaBinary = *optVegaBinary
	}

	initOut, err := vg.initiateNode(vegaBinary, nodeDir, nodeWalletPassFilePath, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate vega node: %w", err)
	}

	confFilePath := ConfigFilePath(nodeDir)
	origConFilePath := originalConfigFilePath(nodeDir)

	if err := utils.CopyFile(confFilePath, origConFilePath); err != nil {
		return nil, fmt.Errorf("failed to copy initiated config from %q to %q: %w", confFilePath, origConFilePath, err)
	}

	initNode := &types.VegaNode{
		GeneratedService: types.GeneratedService{
			Name:           fmt.Sprintf("vega-%s-%d", mode, index),
			HomeDir:        nodeDir,
			ConfigFilePath: confFilePath,
		},
		Mode:                   mode,
		NodeWalletPassFilePath: nodeWalletPassFilePath,
		BinaryPath:             vegaBinary,
	}

	if mode != types.NodeModeValidator {
		log.Printf("vega config initialized for node %q with id %d, paths: %#v", mode, index, initOut.ConfigFilePath)
		return initNode, nil
	}

	nodeWalletInfo, err := vg.initiateValidatorWallets(
		index,
		vegaBinary,
		nodeDir,
		tendermintHome,
		vegaWalletPass,
		ethereumWalletPass,
		nodeWalletPassFilePath,
		clefConf,
	)
	if err != nil {
		return nil, err
	}
	initNode.NodeWalletInfo = nodeWalletInfo

	log.Printf("vega config initialized for node %q with id %d, paths: %#v", mode, index, initOut.ConfigFilePath)

	return initNode, nil
}

func (vg ConfigGenerator) initiateValidatorWallets(
	index int,
	vegaBinary, nodeDir, tendermintHome, vegaWalletPass,
	ethereumWalletPass, nodeWalletPassFilePath string,
	clefConf *config.ClefConfig,
) (*types.NodeWalletInfo, error) {
	walletPassFilePath := path.Join(nodeDir, "vega-wallet-pass.txt")
	ethereumPassFilePath := path.Join(nodeDir, "ethereum-vega-wallet-pass.txt")

	if err := os.WriteFile(walletPassFilePath, []byte(vegaWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write wallet passphrase to file: %w", err)
	}

	if err := os.WriteFile(ethereumPassFilePath, []byte(ethereumWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write ethereum wallet passphrase to file: %w", err)
	}

	vegaOut, err := vg.createWallet(vegaBinary, nodeDir, defaultIsolatedWalletName, walletPassFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create vega wallet: %w", err)
	}
	log.Printf("node wallet create out: %#v", vegaOut)

	vegaImportOut, err := vg.importVegaNodeWallet(vegaBinary, nodeDir, nodeWalletPassFilePath, walletPassFilePath, vegaOut.Wallet.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to import vega wallet: %w", err)
	}
	log.Printf("node wallet import out: %#v", vegaImportOut)

	vegaWalletInfoOut, err := vg.walletInfo(vegaBinary, nodeDir, defaultIsolatedWalletName, walletPassFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get info about vega wallet: %w", err)
	}
	log.Printf("node wallet info out: %#v", vegaWalletInfoOut)

	nwi := &types.NodeWalletInfo{
		VegaWalletID:             vegaWalletInfoOut.ID,
		VegaWalletName:           defaultIsolatedWalletName,
		VegaWalletPassFilePath:   walletPassFilePath,
		VegaWalletRecoveryPhrase: vegaOut.Wallet.RecoveryPhrase,
		VegaWalletPublicKey:      vegaOut.Key.Public,
	}

	if clefConf != nil {
		if err := waitForClef(clefConf.ClefRPCAddr, `{"id": 1, "jsonrpc": "2.0", "method": "account_list"}`, time.Second*30); err != nil {
			return nil, fmt.Errorf("failed to wait for Clef instance: %w", err)
		}

		clefAccountAddr := clefConf.AccountAddresses[index%len(clefConf.AccountAddresses)]

		ethOut, err := vg.importEthereumClefNodeWallet(vegaBinary, nodeDir, nodeWalletPassFilePath, clefAccountAddr, clefConf.ClefRPCAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to import %q wallet: %w", types.NodeWalletChainTypeEthereum, err)
		}
		log.Printf("ethereum wallet out: %#v", ethOut)

		nwi.EthereumAddress = clefAccountAddr
		nwi.EthereumClefRPCAddress = clefConf.ClefRPCAddr
	} else {
		ethOut, err := vg.generateNodeWallet(
			vegaBinary,
			nodeDir,
			nodeWalletPassFilePath,
			ethereumPassFilePath,
			types.NodeWalletChainTypeEthereum,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate %q wallet: %w", types.NodeWalletChainTypeEthereum, err)
		}
		log.Printf("ethereum wallet out: %#v", ethOut)

		ethKey, err := ethereum.DescribeKeyPair(ethOut.WalletFilePath, ethereumWalletPass)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain thereum address for the wallet '%s': %w", ethOut.WalletFilePath, err)
		}

		nwi.EthereumAddress = ethKey.Address
		nwi.EthereumPrivateKey = ethKey.PrivateKey
		nwi.EthereumPassFilePath = ethereumPassFilePath
	}

	tmtOut, err := vg.importTendermintNodeWallet(vegaBinary, nodeDir, nodeWalletPassFilePath, tendermintHome)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tenderming wallet: %w", err)
	}

	log.Printf("tendermint wallet out: %#v", tmtOut)

	return nwi, nil
}

func (vg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *vg.conf.NodeDirPrefix, i)
	return filepath.Join(vg.homeDir, nodeDirName)
}

func ConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node", "config.toml")
}

func originalConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node", "config-original.toml")
}

func IsolatedWalletPath(nodeDir, waleltName string) string {
	return filepath.Join(nodeDir, "data", "wallets", waleltName)
}

func EthereumWalletPath(nodeDir string) string {
	return filepath.Join(nodeDir, "wallets", "ethereum")
}
