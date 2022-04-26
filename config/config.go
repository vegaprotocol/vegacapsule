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
	Name                string         `hcl:"name,label"`
	GenesisTemplate     *string        `hcl:"genesis_template"`
	GenesisTemplateFile *string        `hcl:"genesis_template_file"`
	Ethereum            EthereumConfig `hcl:"ethereum,block"`
	Wallet              *WalletConfig  `hcl:"wallet,block"`
	Faucet              *FaucetConfig  `hcl:"faucet,block"`

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
	StaticPort *StaticPort       `hcl:"static_port,block"`
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

	ConfigTemplates   ConfigTemplates    `hcl:"config_templates,block"`
	NomadJobTemplates *NomadJobTemplates `hcl:"nomad_job_templates,block"`
}

type ConfigTemplates struct {
	Vega       string `hcl:"vega"`
	Tendermint string `hcl:"tendermint"`
	DataNode   string `hcl:"data_node,optional"`
}

type NomadJobTemplates struct {
	VegaFile       *string `hcl:"vega_file"`
	TendermintFile *string `hcl:"tendermint_file"`
	DataNodeFile   *string `hcl:"data_node_file"`
	Vega           *string `hcl:"vega"`
	Tendermint     *string `hcl:"tendermint"`
	DataNode       *string `hcl:"data_node"`
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
		return fmt.Errorf("failed to set absolute paths: %w", err)
	}

	if err := c.validateAndSetGenesis(); err != nil {
		return fmt.Errorf("failed to validate genesis: %w", err)
	}

	if err := c.validateAndSetNodeConfigs(); err != nil {
		return fmt.Errorf("failed to node configs: %w", err)
	}

	if err := c.validateSmartContractsAddresses(); err != nil {
		return fmt.Errorf("invalid configuration for smart contrtacts addresses: %w", err)
	}

	for _, nc := range c.Network.Nodes {
		if nc.NomadJobTemplates == nil {
			continue
		}

		if nc.NomadJobTemplates.Vega != nil {
			fmt.Printf("vega: %s \n", *nc.NomadJobTemplates.Vega)
		}

		if nc.NomadJobTemplates.DataNode != nil {
			fmt.Printf("data-node: %s \n", *nc.NomadJobTemplates.DataNode)
		}

		if nc.NomadJobTemplates.Tendermint != nil {
			fmt.Printf("tendermint: %s \n", *nc.NomadJobTemplates.Tendermint)
		}
	}

	return nil
}

func (c *Config) validateAndSetNodeConfigs() error {
	mErr := utils.NewMultiError()

	for i, nc := range c.Network.Nodes {
		if nc.NomadJobTemplates != nil {
			njt, err := c.validateNomadJobTemplates(nc.NomadJobTemplates)
			if err != nil {
				mErr.Add(fmt.Errorf("failed to validate nomad job templates for %q: %w", nc.Name, err))
				continue
			}

			c.Network.Nodes[i].NomadJobTemplates = njt
		}
	}

	if mErr.HasAny() {
		return mErr
	}

	return nil
}

func (c Config) validateNomadJobTemplates(njt *NomadJobTemplates) (*NomadJobTemplates, error) {
	if njt.DataNode == nil && njt.DataNodeFile != nil {
		dataNodeFile, err := utils.AbsPath(*njt.DataNodeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute file path %q: %w", *njt.DataNodeFile, err)
		}

		dataNode, err := ioutil.ReadFile(dataNodeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", dataNodeFile, err)
		}

		str := string(dataNode)
		njt.DataNode = &str
		njt.DataNodeFile = nil
	}

	if njt.Tendermint == nil && njt.TendermintFile != nil {
		tendermintFile, err := utils.AbsPath(*njt.TendermintFile)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute file path %q: %w", tendermintFile, err)
		}

		tendermint, err := ioutil.ReadFile(tendermintFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", tendermintFile, err)
		}

		str := string(tendermint)
		njt.Tendermint = &str
		njt.TendermintFile = nil
	}

	if njt.Vega == nil && njt.VegaFile != nil {
		vegaFile, err := utils.AbsPath(*njt.VegaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute file path %q: %w", vegaFile, err)
		}

		vega, err := ioutil.ReadFile(vegaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", vegaFile, err)
		}

		str := string(vega)
		njt.Vega = &str
		njt.VegaFile = nil
	}

	return njt, nil
}

func (c *Config) validateAndSetGenesis() error {
	if c.Network.GenesisTemplate != nil {
		return nil
	}

	if c.Network.GenesisTemplateFile != nil {
		genTemplateFile, err := utils.AbsPath(*c.Network.GenesisTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to get absolute file path %q: %w", genTemplateFile, err)
		}

		genTemplate, err := ioutil.ReadFile(genTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", genTemplateFile, err)
		}

		genTemplateStr := string(genTemplate)
		// set file content as template a set file path to nil
		c.Network.GenesisTemplate = &genTemplateStr
		c.Network.GenesisTemplateFile = nil

		return nil
	}

	return fmt.Errorf("missing genesis file template: either genesis_template or genesis_template_file must be defined")
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
