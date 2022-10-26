package config

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/types"
)

/*
description: |

	Network configuration allows to customise Vega network into different shapes based on personal needs.
	It allows to configure and deploy different Vega nodes setups (validator, full) and their dependencies (like Ethereum or Postgres).
	It can run custom Docker images before and after the network nodes has started and much more.

example:

	type: hcl
	value: |
		network "testnet" {
			ethereum {
				...
			}

			pre_start {
				...
			}

			genesis_template_file = "..."
			smart_contracts_addresses_file = "..."

			node_set "validator-nodes" {
				...
			}

			node_set "full-nodes" {
				...
			}
		}
*/
type NetworkConfig struct {
	/*
		description: |
			Name of the network.
			All folders generated are placed in folder with this name.
			All Nomad jobs are prefix with the name.
		example:
			type: hcl
			value: |
				network "network_name" {
					...
				}
	*/
	Name string `hcl:"name,label"`

	/*
		description: |
			Template of genesis file that will be used to bootrap the Vega network.
		note: |
				It is recomended to use `genesis_template_file` param instead.
				In case both `genesis_template` and `genesis_template_file` are defined the `genesis_template`
				overrides `genesis_template_file`.
		examples:
			- type: hcl
			  value: |
						genesis_template = <<EOH
							{
								"app_state": {
									...
								}
								..
							}
						EOH

	*/
	GenesisTemplate     *string `hcl:"genesis_template"`
	GenesisTemplateFile *string `hcl:"genesis_template_file"`

	SmartContractsAddresses     *string `hcl:"smart_contracts_addresses,optional"`
	SmartContractsAddressesFile *string `hcl:"smart_contracts_addresses_file,optional"`

	Ethereum EthereumConfig `hcl:"ethereum,block"`

	Nodes []NodeConfig `hcl:"node_set,block" cty:"node_set"`

	Wallet *WalletConfig `hcl:"wallet,block"`
	Faucet *FaucetConfig `hcl:"faucet,block"`

	PreStart  *PStartConfig `hcl:"pre_start,block"`
	PostStart *PStartConfig `hcl:"post_start,block"`

	TokenAddresses map[string]types.SmartContractsToken
}

type PStartConfig struct {
	Docker []DockerConfig `hcl:"docker_service,block"`
}

func (nc NetworkConfig) GetNodeConfig(name string) (*NodeConfig, error) {
	for _, nodeConf := range nc.Nodes {
		if nodeConf.Name == name {
			return &nodeConf, nil
		}
	}

	return nil, fmt.Errorf("node config with name %q not found", name)
}
