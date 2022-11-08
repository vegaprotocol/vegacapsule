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
			All folders generated are placed in the folder with this name.
			All Nomad jobs are prefix with the name.
		example:
			type: hcl
			value: |
				network "name" {
					...
				}
	*/
	Name string `hcl:"name,label"`

	/*
		description: |
			Go template of genesis file that will be used to bootrap the Vega network.
			[Example of templated mainnet genesis file](https://github.com/vegaprotocol/networks/blob/master/mainnet1/genesis.json)
		optional_if: genesis_template_file
		note: |
				It is recommended that you use `genesis_template_file` param instead.
				If both `genesis_template` and `genesis_template_file` are defined, then `genesis_template`
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
	GenesisTemplate *string `hcl:"genesis_template"`

	/*
		description: |
			Same as `genesis_template` but it allows to link the genesis file template as an external file.
		examples:
			- type: hcl
			  value: |
						genesis_template_file = "/your_path/genesis.tmpl"

	*/
	GenesisTemplateFile *string `hcl:"genesis_template_file"`

	/*
		description: |
			Allows to define Ethereum network configuration.
			This is necessery because Vega needs to be connected to [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge)
			or it cannot function otherwise.
		examples:
			- type: hcl
			  value: |
						ethereum {
							...
						}
	*/
	Ethereum EthereumConfig `hcl:"ethereum,block"`

	/*
		description: |
			Smart contract addresses are addresses of [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge) contracts in JSON format.

			These addresses should correspond to the choosen network by [Ethereum network](#EthereumConfig) and
			can be used in various different types of templates in Capsule.
			[Example of smart contract address from mainnet](https://github.com/vegaprotocol/networks/blob/master/mainnet1/smart-contracts.json).
		note: |
				It is recomended to use `smart_contracts_addresses_file` param instead.
				In case both `smart_contracts_addresses` and `smart_contracts_addresses_file` are defined the `genesis_template`
				overrides `smart_contracts_addresses_file`.
		optional_if: smart_contracts_addresses_file
		examples:
			- type: hcl
			  value: |
						smart_contracts_addresses = <<EOH
							{
								"erc20_bridge": "...",
								"staking_bridge": "...",
								...
							}
						EOH
	*/
	SmartContractsAddresses *string `hcl:"smart_contracts_addresses,optional"`

	/*
		description: |
			Same as `smart_contracts_addresses` but it allows to link the smart contracts as an external file.
		examples:
			- type: hcl
			  value: |
						smart_contracts_addresses_file = "/your_path/smart-contratcs.json"
	*/
	SmartContractsAddressesFile *string `hcl:"smart_contracts_addresses_file,optional"`

	/*
		description: |
			Allows to define multiple nodes set and their specific configuration.
			A node set is a representation of Vega and Data Node nodes.
			This is building unit of the Vega network.
		examples:
			- type: hcl
			  name: Validators node set
			  value: |
						node_set "validator-nodes" {
							...
						}
			- type: hcl
			  name: Full nodes node set
			  value: |
						node_set "full-nodes" {
							...
						}
	*/
	Nodes []NodeConfig `hcl:"node_set,block" cty:"node_set"`

	/*
		description: |
			Allows to deploy and configure [Vega Wallet](https://docs.vega.xyz/mainnet/tools/vega-wallet) instance.
			Wallet will not be deployed if this block is not defined.
		examples:
			- type: hcl
			  value: |
						wallet "wallet-name" {
							...
						}
	*/
	Wallet *WalletConfig `hcl:"wallet,block"`

	/*
		description: |
			Allows to deploy and configure [Vega Core Faucet](https://github.com/vegaprotocol/vega/tree/develop/core/faucet#faucet) instance.
			Faucet will not be deployed if this block is not defined.
		examples:
			- type: hcl
			  value: |
						faucet "faucet-name" {
							...
						}
	*/
	Faucet *FaucetConfig `hcl:"faucet,block"`

	/*
		description: |
			Allows to define jobs that should run before the node sets starts.
			It can be used for node sets dependencies like databases or mock Ethereum chain etc..
		examples:
			- type: hcl
			  value: |
						pre_start {
							docker_service "ganache-1" {
								...
							}
							docker_service "postgres-1" {
								...
							}
						}
	*/
	PreStart *PStartConfig `hcl:"pre_start,block"`

	/*
		description: |
			Allows to define jobs that should run after the node sets started.
			It can be used for services that depends not already running network like block explorer or console.
		examples:
			- type: hcl
			  value: |
						post_start {
							docker_service "bloc-explorer-1" {
								...
							}
							docker_service "vega-console-1" {
								...
							}
						}
	*/
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
