package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Config struct {
	OutputDir            string        `hcl:"output_dir"`
	VegaBinary           string        `hcl:"vega_binary_path"`
	Prefix               string        `hcl:"prefix"`
	NodeDirPrefix        string        `hcl:"node_dir_prefix"`
	TendermintNodePrefix string        `hcl:"tendermint_node_prefix"`
	VegaNodePrefix       string        `hcl:"vega_node_prefix"`
	DataNodePrefix       string        `hcl:"data_node_prefix"`
	WalletPrefix         string        `hcl:"wallet_prefix"`
	FaucetPrefix         string        `hcl:"faucet_prefix"`
	Network              NetworkConfig `hcl:"network,block"`
}

type NetworkConfig struct {
	Name            string         `hcl:"name,label"`
	GenesisTemplate string         `hcl:"genesis_template"`
	Ethereum        EthereumConfig `hcl:"ethereum,block"`
	Wallet          *WalletConfig  `hcl:"wallet,block"`
	Faucet          *FaucetConfig  `hcl:"faucet,block"`

	PreStart *PrestartConfig `hcl:"pre_start,block"`

	Nodes []NodeConfig `hcl:"node_set,block"`
}

type EthereumConfig struct {
	ChainID   string `hcl:"chain_id"`
	NetworkID string `hcl:"network_id"`
	Endpoint  string `hcl:"endpoint"`
}

type PrestartConfig struct {
	Docker []DockerConfig `hcl:"docker_service,block"`
}

type DockerConfig struct {
	Name       string   `hcl:"name,label"`
	Image      string   `hcl:"image"`
	Command    string   `hcl:"cmd"`
	Args       []string `hcl:"args"`
	StaticPort int      `hcl:"static_port,optional"`
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
	if !filepath.IsAbs(c.OutputDir) {
		absPath, err := filepath.Abs(c.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for outputDir: %w", err)
		}
		c.OutputDir = absPath
	}

	// Vega binary
	vegaBinPath, err := utils.BinaryAbsPath(c.VegaBinary)
	if err != nil {
		return err
	}
	c.VegaBinary = vegaBinPath

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

	return nil
}

func (c *Config) Persist() error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(*c, f.Body())
	return ioutil.WriteFile(filepath.Join(c.OutputDir, "config.hcl"), f.Bytes(), 0644)
}
