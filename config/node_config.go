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
		description: Templates that can be used for configurations of Vega and Data nodes, Tendermint and other services.
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
					Allows user to define a Vega binary to be used in specific node set only.
					A relative or absolute path can be used. If only the binary name is defined, it automatically looks for it in $PATH.
					This can help with testing different version compatibilities or a protocol upgrade.
		note: Using versions that are not compatible could break the network - therefore this should be used in advanced cases only.
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

/*
description: |

	Allow to add configuration template for certain services deployed by Capsule.
	Learn more about how configuration templating work here
*/
type ConfigTemplates struct {
	/*
		description: Go template of Vega config.
		optional_if: vega_file
		note: |
				It is recommended that you use `vega_file` param instead.
				If both `vega` and `vega_file` are defined, then `vega`
				overrides `vega_file`.
		examples:
			- type: hcl
			  value: |
						vega = <<EOH
							...
						EOH

	*/
	Vega *string `hcl:"vega,optional"`

	/*
		description: |
			Same as `vega` but it allows the user to link the Vega config template as an external file.
		examples:
			- type: hcl
			  value: |
						vega_file = "/your_path/vega_config.tmpl"

	*/
	VegaFile *string `hcl:"vega_file,optional"`

	/*
		description: Go template of Tendermint config.
		optional_if: tendermint_file
		note: |
				It is recommended that you use `tendermint_file` param instead.
				If both `tendermint` and `tendermint_file` are defined, then `tendermint`
				overrides `tendermint_file`.
		examples:
			- type: hcl
			  value: |
						tendermint = <<EOH
							...
						EOH

	*/
	Tendermint *string `hcl:"tendermint,optional"`
	/*
		description: |
			Same as `tendermint` but it allows the user to link the Tendermint config template as an external file.
		examples:
			- type: hcl
			  value: |
						tendermint_file = "/your_path/tendermint_config.tmpl"

	*/
	TendermintFile *string `hcl:"tendermint_file,optional"`

	/*
		description: Go template of Data Node config.
		optional_if: data_node_file
		note: |
				It is recommended that you use `data_node_file` param instead.
				If both `data_node` and `data_node_file` are defined, then `data_node`
				overrides `data_node_file`.
		example:
			type: hcl
			value: |
					data_node = <<EOH
						...
					EOH

	*/
	DataNode *string `hcl:"data_node,optional"`

	/*
		description: |
			Same as `data_node` but it allows the user to link the Data Node config template as an external file.
		example:
			type: hcl
			value: |
					data_node_file = "/your_path/data_node_config.tmpl"

	*/
	DataNodeFile *string `hcl:"data_node_file,optional"`

	/*
		description: |
						Go template of Visor genesis run config.
						Current Vega binary is automatically copied to the Visor genesis folder by Capsule
						so it can be used from this template.
		optional_if: visor_run_conf_file
		note: |
				It is recommended that you use `visor_run_conf_file` param instead.
				If both `visor_run_conf` and `visor_run_conf_file` are defined, then `visor_run_conf`
				overrides `visor_run_conf_file`.
		example:
			type: hcl
			value: |
					visor_run_conf = <<EOH
						...
					EOH

	*/
	VisorRunConf *string `hcl:"visor_run_conf,optional"`
	/*
		description: |
			Same as `visor_run_conf` but it allows the user to link the Visor genesis run config template as an external file.
		example:
			type: hcl
			value: |
					visor_run_conf_file = "/your_path/visor_run_config.tmpl"

	*/
	VisorRunConfFile *string `hcl:"visor_run_conf_file,optional"`

	/*
		description: Go template of Visor config.
		optional_if: visor_conf_file
		note: |
				It is recommended that you use `visor_conf_file` param instead.
				If both `visor_conf` and `visor_conf_file` are defined, then `visor_conf`
				overrides `visor_conf_file`.
		example:
			type: hcl
			value: |
					visor_conf = <<EOH
						...
					EOH

	*/
	VisorConf *string `hcl:"visor_conf,optional"`

	/*
		description: |
			Same as `visor_conf` but it allows the user to link the Visor genesis run config template as an external file.
		example:
			type: hcl
			value: |
					visor_conf_file = "/your_path/visor_config.tmpl"

	*/
	VisorConfFile *string `hcl:"visor_conf_file,optional"`
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
