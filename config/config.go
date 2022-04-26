package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Config struct {
	OutputDir            *string       `hcl:"-"`
	VegaBinary           *string       `hcl:"vega_binary_path"`
	Prefix               *string       `hcl:"prefix"`
	NodeDirPrefix        *string       `hcl:"node_dir_prefix"`
	TendermintNodePrefix *string       `hcl:"tendermint_node_prefix"`
	VegaNodePrefix       *string       `hcl:"vega_node_prefix"`
	DataNodePrefix       *string       `hcl:"data_node_prefix"`
	WalletPrefix         *string       `hcl:"wallet_prefix"`
	FaucetPrefix         *string       `hcl:"faucet_prefix"`
	Network              NetworkConfig `hcl:"network,block"`
}

type NetworkConfig struct {
	Name            string         `hcl:"name,label"`
	GenesisTemplate string         `hcl:"genesis_template"`
	Ethereum        EthereumConfig `hcl:"ethereum,block"`
	Wallet          *WalletConfig  `hcl:"wallet,block"`
	Faucet          *FaucetConfig  `hcl:"faucet,block"`

	PreStart *PrestartConfig `hcl:"pre_start,block"`

	Nodes                   []NodeConfig `hcl:"node_set,block"`
	SmartContractsAddresses string       `hcl:"smart_contracts_addresses"`
}

func (nc NetworkConfig) GetNodeConfig(name string) (*NodeConfig, error) {
	for _, nodeConf := range nc.Nodes {
		if nodeConf.Name == name {
			return &nodeConf, nil
		}
	}

	return nil, fmt.Errorf("node config with name %q not found", name)
}

type EthereumConfig struct {
	ChainID   string `hcl:"chain_id"`
	NetworkID string `hcl:"network_id"`
	Endpoint  string `hcl:"endpoint"`
}

type PrestartConfig struct {
	Docker []DockerConfig `hcl:"docker_service,block"`
}

type StaticPort struct {
	To    int `hcl:"to,optional"`
	Value int `hcl:"value"`
}
type DockerConfig struct {
	Name       string            `hcl:"name,label"`
	Image      string            `hcl:"image"`
	Command    string            `hcl:"cmd"`
	Args       []string          `hcl:"args"`
	Env        map[string]string `hcl:"env,optional"`
	StaticPort StaticPort        `hcl:"static_port,block"`
	ToPort     int               `hcl:"to_port,optional"`
}

type WalletConfig struct {
	Name     string `hcl:"name,label"`
	Binary   string `hcl:"binary"`
	Template string `hcl:"template,optional"`
}

type FaucetConfig struct {
	Name     string `hcl:"name,label"`
	Pass     string `hcl:"wallet_pass"`
	Template string `hcl:"template,optional"`
}

type NodeConfig struct {
	Name               string `hcl:"name,label"`
	Mode               string `hcl:"mode"`
	Count              int    `hcl:"count"`
	NodeWalletPass     string `hcl:"node_wallet_pass,optional"`
	EthereumWalletPass string `hcl:"ethereum_wallet_pass,optional"`
	VegaWalletPass     string `hcl:"vega_wallet_pass,optional"`
	DataNodeBinary     string `hcl:"data_node_binary,optional"`

	Templates TemplateConfig `hcl:"config_templates,block"`
}

type TemplateConfig struct {
	Vega       string `hcl:"vega"`
	Tendermint string `hcl:"tendermint"`
	DataNode   string `hcl:"data_node,optional"`
}

func (c *Config) setAbsolutePaths() error {
	// Output directory
	if !filepath.IsAbs(*c.OutputDir) {
		absPath, err := filepath.Abs(*c.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for outputDir: %w", err)
		}
		*c.OutputDir = absPath
	}

	// Vega binary
	vegaBinPath, err := utils.BinaryAbsPath(*c.VegaBinary)
	if err != nil {
		return err
	}
	*c.VegaBinary = vegaBinPath

	// Wallet binary
	if c.Network.Wallet != nil {
		walletBinPath, err := utils.BinaryAbsPath(c.Network.Wallet.Binary)
		if err != nil {
			return err
		}
		c.Network.Wallet.Binary = walletBinPath
	}

	// Data nodes binaries
	for idx, nc := range c.Network.Nodes {
		if nc.DataNodeBinary == "" {
			continue
		}

		dataNodeBinPath, err := utils.BinaryAbsPath(nc.DataNodeBinary)
		if err != nil {
			return err
		}
		c.Network.Nodes[idx].DataNodeBinary = dataNodeBinPath
	}

	return nil
}

func (c *Config) Validate() error {
	if err := c.setAbsolutePaths(); err != nil {
		return err
	}

	if err := c.validateSmartContractsAddresses(); err != nil {
		return fmt.Errorf("invalid configuration for smart contrtacts addresses: %w", err)
	}

	return nil
}

func (c Config) validateSmartContractsAddresses() error {
	_, err := c.SmartContractsInfo()

	if err != nil {
		return fmt.Errorf("failed to check smart contract addreses: %w", err)
	}

	return nil
}

func (c Config) SmartContractsInfo() (*types.SmartContractsInfo, error) {
	smartcontracts := &types.SmartContractsInfo{}

	if err := json.Unmarshal([]byte(c.Network.SmartContractsAddresses), &smartcontracts); err != nil {
		return nil, fmt.Errorf("failed to get smart contracts info: config.network.smart_contracts_addresses format is wrong: %w", err)
	}

	return smartcontracts, nil
}

func (c *Config) Persist() error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(*c, f.Body())
	return ioutil.WriteFile(filepath.Join(*c.OutputDir, "config.hcl"), f.Bytes(), 0644)
}

func DefaultConfig() (*Config, error) {
	outputDir, err := DefaultNetworkHome()
	if err != nil {
		return nil, err
	}

	return &Config{
		OutputDir:            &outputDir,
		Prefix:               utils.StrPoint("st-local"),
		NodeDirPrefix:        utils.StrPoint("node"),
		TendermintNodePrefix: utils.StrPoint("tendermint"),
		VegaNodePrefix:       utils.StrPoint("vega"),
		DataNodePrefix:       utils.StrPoint("data"),
		WalletPrefix:         utils.StrPoint("wallet"),
		FaucetPrefix:         utils.StrPoint("faucet"),
	}, nil
}

func DefaultNetworkHome() (string, error) {
	capsuleHome, err := utils.CapsuleHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(capsuleHome, "testnet"), nil
}
