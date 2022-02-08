package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
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
	Network              NetworkConfig `hcl:"network,block"`
}

type NetworkConfig struct {
	Name             string       `hcl:"name,label"`
	GenesisTemplate  string       `hcl:"genesis_template"`
	ChainID          string       `hcl:"chain_id"`
	NetworkID        string       `hcl:"network_id"`
	EthereumEndpoint string       `hcl:"ethereum_endpoint"`
	Nodes            []NodeConfig `hcl:"node_set,block"`
}

type NodeConfig struct {
	Name               string `hcl:"name,label"`
	Mode               string `hcl:"mode"`
	Count              int    `hcl:"count"`
	NodeWalletPass     string `hcl:"node_wallet_pass"`
	EthereumWalletPass string `hcl:"ethereum_wallet_pass"`
	VegaWalletPass     string `hcl:"vega_wallet_pass"`

	Templates TemplateConfig `hcl:"config_templates,block"`
}

type TemplateConfig struct {
	Vega       string `hcl:"vega"`
	Tendermint string `hcl:"tendermint"`
	DataNode   string `hcl:"tendermint"`
}

func (c *Config) Persist() error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(*c, f.Body())
	return ioutil.WriteFile(filepath.Join(c.OutputDir, "config.hcl"), f.Bytes(), 0644)
}

func ParseConfig(conf []byte) (*Config, error) {
	config := &Config{}
	if err := hclsimple.Decode("config.hcl", conf, nil, config); err != nil {
		return nil, fmt.Errorf("failed to load decode configuration: %w", err)
	}

	return config, nil
}

func ParseConfigFile(filePath string) (*Config, error) {
	config := &Config{}
	if err := hclsimple.DecodeFile(filePath, nil, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return config, nil
}
