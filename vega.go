package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/zannen/toml"

	"github.com/imdario/mergo"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vega/config"
)

const (
	NodeModeValidator           = "validator"
	NodeModeFull                = "full"
	NodeWalletChainTypeVega     = "vega"
	NodeWalletChainTypeEthereum = "ethereum"
)

type InitVegaNode struct {
	HomeDir                string
	NodeWalletPassFilePath string
}

type VegaTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	ETHEndpoint          string
	NodeMode             string
	NodeNumber           int
}

type VegaConfigGenerator struct {
	conf    *Config
	homeDir string
}

func NewVegaConfigGenerator(conf *Config) (*VegaConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(conf.OutputDir, conf.VegaNodePrefix))
	if err != nil {
		return nil, err
	}

	return &VegaConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (vg VegaConfigGenerator) Initiate(index int, mode, tendermintHome, nodeWalletPass, vegaWalletPass, ethereumWalletPass string) (*InitVegaNode, error) {
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

	if mode != NodeModeValidator {
		log.Printf("vega config initialized for node %q with id %d, paths: %#v", mode, index, initOut.ConfigFilePath)

		return nil, nil
	}

	if err := vg.initiateValidatorWallets(nodeDir, tendermintHome, vegaWalletPass, ethereumWalletPass, nodeWalletPassFilePath); err != nil {
		return nil, err
	}

	log.Printf("vega config initialized for node %q with id %d, paths: %#v", mode, index, initOut.ConfigFilePath)

	return &InitVegaNode{
		HomeDir:                nodeDir,
		NodeWalletPassFilePath: nodeWalletPassFilePath,
	}, nil
}

func (vg VegaConfigGenerator) initiateValidatorWallets(nodeDir, tendermintHome, vegaWalletPass, ethereumWalletPass, nodeWalletPassFilePath string) error {
	walletPassFilePath := path.Join(nodeDir, "vega-wallet-pass.txt")
	ethereumPassFilePath := path.Join(nodeDir, "ethereum-vega-wallet-pass.txt")

	if err := ioutil.WriteFile(walletPassFilePath, []byte(vegaWalletPass), 0644); err != nil {
		return fmt.Errorf("failed to write wallet passphrase to file: %w", err)
	}

	if err := ioutil.WriteFile(ethereumPassFilePath, []byte(ethereumWalletPass), 0644); err != nil {
		return fmt.Errorf("failed to write ethereum wallet passphrase to file: %w", err)
	}

	vegaOut, err := vg.generateNodeWallet(nodeDir, nodeWalletPassFilePath, walletPassFilePath, NodeWalletChainTypeVega)
	if err != nil {
		return fmt.Errorf("failed to generate vega wallet: %w", err)
	}

	log.Printf("node wallet out: %#v", vegaOut)

	ethOut, err := vg.generateNodeWallet(nodeDir, nodeWalletPassFilePath, ethereumPassFilePath, NodeWalletChainTypeEthereum)
	if err != nil {
		return fmt.Errorf("failed to generate vega wallet: %w", err)
	}

	log.Printf("ethereum wallet out: %#v", ethOut)

	tmtOut, err := vg.importTendermintNodeWallet(nodeDir, nodeWalletPassFilePath, tendermintHome)
	if err != nil {
		return fmt.Errorf("failed to generate tenderming wallet: %w", err)
	}

	log.Printf("tendermint wallet out: %#v", tmtOut)

	return nil
}

func (vg VegaConfigGenerator) OverrideConfig(index int, mode string, configTemplate *template.Template) error {
	templateCtx := VegaTemplateContext{
		Prefix:               vg.conf.Prefix,
		TendermintNodePrefix: vg.conf.TendermintNodePrefix,
		VegaNodePrefix:       vg.conf.VegaNodePrefix,
		DataNodePrefix:       vg.conf.DataNodePrefix,
		ETHEndpoint:          vg.conf.Network.EthereumEndpoint,
		NodeMode:             mode,
		NodeNumber:           index,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	overrideConfig := config.Config{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	configFilePath := vg.configFilePath(vg.nodeDir(index))

	vegaConfig := config.NewDefaultConfig()
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

func (vg VegaConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", vg.conf.NodeDirPrefix, i)
	return filepath.Join(vg.homeDir, nodeDirName)
}

func (vg VegaConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node", "config.toml")
}

type vegaNode struct {
	NodeMode               string
	NodeHome               string
	WalletPassFilePath     string
	NodeWalletPassFilePath string
	EthereumPassFilePath   string
	VegaWallet             *generateNodeWalletOutput
	EthereumWallet         *generateNodeWalletOutput
	TendermintWallet       *importNodeWalletOutput
}
