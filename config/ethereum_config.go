package config

/*
description: |

	Allows to define specific Ethereum network to be used.
	It can either some of the [Public networks](https://ethereum.org/en/developers/docs/networks/#public-networks) or
	local instance of Ganache.

note: |

	The chosen network needs to have deployed [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge) on it.

example:

	type: hcl
	name: Setup for local Ganache
	value: |
			ethereum {
				chain_id   = "1440"
				network_id = "1441"
				endpoint   = "http://127.0.0.1:8545/"
			}
*/
type EthereumConfig struct {
	ChainID   string `hcl:"chain_id"`
	NetworkID string `hcl:"network_id"`
	Endpoint  string `hcl:"endpoint"`
}
