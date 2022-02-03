package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
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
	Name            string       `hcl:"name,label"`
	GenesisTemplate string       `hcl:"genesis_template"`
	ChainID         string       `hcl:"chain_id"`
	NetworkID       string       `hcl:"network_id"`
	Nodes           []NodeConfig `hcl:"node_set,block"`
}

type NodeConfig struct {
	Name      string         `hcl:"name,label"`
	Mode      string         `hcl:"mode"`
	Count     int            `hcl:"count"`
	Templates TemplateConfig `hcl:"config_templates,block"`
}

type TemplateConfig struct {
	Vega       string `hcl:"vega"`
	Tendermint string `hcl:"tendermint"`
	DataNode   string `hcl:"tendermint"`
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
