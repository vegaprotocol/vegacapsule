package vega

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/zannen/toml"

	"github.com/imdario/mergo"

	"code.vegaprotocol.io/shared/paths"
	vgconfig "code.vegaprotocol.io/vega/config"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/ethereum"
	"code.vegaprotocol.io/vegacapsule/types"
)

type ConfigTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	ETHEndpoint          string
	NodeMode             string
	FaucetPublicKey      string
	NodeNumber           int
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

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

	initNode := &types.VegaNode{
		Mode:                   mode,
		HomeDir:                nodeDir,
		NodeWalletPassFilePath: nodeWalletPassFilePath,
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

func (vg ConfigGenerator) OverwriteConfig(index int, mode string, fc *types.Faucet, configTemplate *template.Template) error {
	templateCtx := ConfigTemplateContext{
		Prefix:               *vg.conf.Prefix,
		TendermintNodePrefix: *vg.conf.TendermintNodePrefix,
		VegaNodePrefix:       *vg.conf.VegaNodePrefix,
		DataNodePrefix:       *vg.conf.DataNodePrefix,
		ETHEndpoint:          vg.conf.Network.Ethereum.Endpoint,
		NodeMode:             mode,
		NodeNumber:           index,
	}

	if fc != nil {
		templateCtx.FaucetPublicKey = fc.PublicKey
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	overrideConfig := vgconfig.Config{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	configFilePath := vg.configFilePath(vg.nodeDir(index))

	vegaConfig := vgconfig.NewDefaultConfig()
	if err := paths.ReadStructuredFile(configFilePath, &vegaConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configFilePath, err)
	}

	if err := mergo.Merge(&overrideConfig, vegaConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(configFilePath, overrideConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", configFilePath, err)
	}

	return nil
}

func (vg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *vg.conf.NodeDirPrefix, i)
	return filepath.Join(vg.homeDir, nodeDirName)
}

func (vg ConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node", "config.toml")
}
