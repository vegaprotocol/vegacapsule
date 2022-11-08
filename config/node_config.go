package config

import (
	"encoding/json"

	"code.vegaprotocol.io/vegacapsule/types"
)

/*
description: |

	Represents, and allows the user to configure, a set of Vega (with Tendermint) and Data Node nodes.
	One node set definition can be used by applied to multiple node sets (see `count` field) and it uses
	templating to distinguish between different nodes and names/ports and other collisions.

example:

	type: hcl
	name: Node set with 2 validator nodes
	value: |
			node_set "validators" {
				count = 2
				mode = "validator"

				node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
				vega_wallet_pass = "w4ll3t-p4ssphr4e3"
				ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

				config_templates {
					vega_file = "./path/vega_validator.tmpl"
					tendermint_file = "./path/tendermint_validator.tmpl"
				}
			}
*/
type NodeConfig struct {
	/*
		description: |
			Name of the node set.
			Nomad instances that are part of these nodes are prefixed with this name.
		example:
			type: hcl
			value: |
				node_set "validators-1" {
					...
				}
	*/
	Name string `hcl:"name,label" cty:"name"`

	/*
		description: |
			Determines what mode the node set should run in.
		values:
			- validator
			- full
	*/
	Mode string `hcl:"mode" cty:"mode"`

	/*
		description: |
			Defines how many node sets with this exact configuration should be created.
	*/
	Count int `hcl:"count" cty:"count"`

	/*
		description: Defines the password for the automatically generated node wallet associated with the created node.
		required_if: mode=validator
	*/
	NodeWalletPass string `hcl:"node_wallet_pass,optional" template:"" cty:"node_wallet_pass"`

	/*
		description: Defines password for automatically generated Ethereum wallet in node wallet.
		required_if: mode=validator
	*/
	EthereumWalletPass string `hcl:"ethereum_wallet_pass,optional" template:"" cty:"ethereum_wallet_pass"`

	/*
		description: Defines password for automatically generated Vega wallet in node wallet.
		required_if: mode=validator
	*/
	VegaWalletPass string `hcl:"vega_wallet_pass,optional" template:"" cty:"vega_wallet_pass"`

	/*
		description: Whether or not Data Node should be deployed on node set.
	*/
	UseDataNode bool `hcl:"use_data_node,optional" cty:"use_data_node"`

	/*
		description: |
					Path to [Visor](https://github.com/vegaprotocol/vega/tree/develop/visor) binary.
					If defined, Visor is automatically used to deploy Vega and Data nodes.
					The relative or absolute path can be used, if only the binary name is defined it automatically looks for it in $PATH.
	*/
	VisorBinary string `hcl:"visor_binary,optional"`

	/*
		description: Templates that allows configurations of Vega, Data, Tendermint and other services.
		example:
				type: hcl
				value: |
					config_templates {
						vega_file = "./path/vega.tmpl"
						tendermint_file = "./path/tendermint.tmpl"
						data_node_file = "./path/data_node.tmpl"
					}
	*/
	ConfigTemplates ConfigTemplates `hcl:"config_templates,block"`

	/*
		description: |
					Alows to define a Vega binary to be used in specific node set only.
					Relative or absolute path can be used or if only binary name is defined it automatically looks up in $PATH.
					This can help with testing of different version compatibilities or protocol upgrade.
		note: Using versions that are not compatible might break the network - therefore this should be use in advance cases only.
	*/
	VegaBinary *string `hcl:"vega_binary_path,optional"`

	PreGenerate *PreGenerate `hcl:"pre_generate,block"`

	PreStartProbe *types.ProbesConfig `hcl:"pre_start_probe,block" template:""`

	ClefWallet *ClefConfig `hcl:"clef_wallet,block" template:""`

	NomadJobTemplate     *string `hcl:"nomad_job_template,optional"`
	NomadJobTemplateFile *string `hcl:"nomad_job_template_file,optional"`
}

type PreGenerate struct {
	Nomad []NomadConfig `hcl:"nomad_job,block"`
}

type ClefConfig struct {
	AccountAddresses []string `hcl:"ethereum_account_addresses" template:""`
	ClefRPCAddr      string   `hcl:"clef_rpc_address" template:""`
}

type ConfigTemplates struct {
	Vega             *string `hcl:"vega,optional"`
	VegaFile         *string `hcl:"vega_file,optional"`
	Tendermint       *string `hcl:"tendermint,optional"`
	TendermintFile   *string `hcl:"tendermint_file,optional"`
	DataNode         *string `hcl:"data_node,optional"`
	DataNodeFile     *string `hcl:"data_node_file,optional"`
	VisorRunConf     *string `hcl:"visor_run_conf,optional"`
	VisorRunConfFile *string `hcl:"visor_run_conf_file,optional"`
	VisorConf        *string `hcl:"visor_conf,optional"`
	VisorConfFile    *string `hcl:"visor_conf_file,optional"`
}

func (nc NodeConfig) Clone() (*NodeConfig, error) {
	origJSON, err := json.Marshal(nc)
	if err != nil {
		return nil, err
	}

	clone := NodeConfig{}
	if err = json.Unmarshal(origJSON, &clone); err != nil {
		return nil, err
	}

	return &clone, nil
}
