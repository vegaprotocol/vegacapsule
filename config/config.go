package config

import (
	"io/ioutil"
	"path/filepath"

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
	Name             string        `hcl:"name,label"`
	GenesisTemplate  string        `hcl:"genesis_template"`
	ChainID          string        `hcl:"chain_id"`
	NetworkID        string        `hcl:"network_id"`
	EthereumEndpoint string        `hcl:"ethereum_endpoint"`
	Wallet           *WalletConfig `hcl:"wallet,block"`
	Faucet           *FaucetConfig `hcl:"faucet,block"`

	PreStart *PrestartConfig `hcl:"pre_start,block"`

	Nodes []NodeConfig `hcl:"node_set,block"`
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

func (c *Config) Persist() error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(*c, f.Body())
	return ioutil.WriteFile(filepath.Join(c.OutputDir, "config.hcl"), f.Bytes(), 0644)
}
