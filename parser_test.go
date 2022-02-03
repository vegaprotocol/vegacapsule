package main

import (
	"fmt"
	"log"
	"testing"
)

var testConfig = `
output_dir             = "/Users/karelmoravec/vega/vegacomposer/testnet"
prefix                 = "st-local"
node_dir_prefix        = "node"
tendermint_node_prefix = "tendermint-node"
vega_node_prefix       = "vega-node"
data_node_prefix       = "data-node"
vega_binary_path       = "/Users/karelmoravec/go/bin/vega"

network "testnet" {
  chain_id               = "1440"
  network_id             = "1441"
  genesis_template = <<EOH
{
	"app_state": {
	  "assets": {
		"fBTC": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "BTC (fake)",
		  "symbol": "fBTC",
		  "total_supply": "21000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "1000000"
			}
		  }
		},
		"fDAI": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "DAI (fake)",
		  "symbol": "fDAI",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "10000000000"
			}
		  }
		},
		"fEURO": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "EURO (fake)",
		  "symbol": "fEURO",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "10000000000"
			}
		  }
		},
		"fUSDC": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "USDC (fake)",
		  "symbol": "fUSDC",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "1000000000000"
			}
		  }
		},
		"XYZalpha": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "XYZ (α alpha)",
		  "symbol": "XYZalpha",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "100000000000"
			}
		  }
		},
		"XYZbeta": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "XYZ (β beta)",
		  "symbol": "XYZbeta",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "100000000000"
			}
		  }
		},
		"XYZgamma": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "XYZ (γ gamma)",
		  "symbol": "XYZgamma",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "100000000000"
			}
		  }
		},
		"XYZdelta": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "XYZ (δ delta)",
		  "symbol": "XYZdelta",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "100000000000"
			}
		  }
		},
		"XYZepsilon": {
		  "min_lp_stake": "1",
		  "decimals": 5,
		  "name": "XYZ (ε epsilon)",
		  "symbol": "XYZepsilon",
		  "total_supply": "1000000000",
		  "source": {
			"builtin_asset": {
			  "max_faucet_amount_mint": "100000000000"
			}
		  }
		},
		"{{.GetVegaContractID "tBTC"}}": {
			"min_lp_stake": "1",
			"decimals": 5,
			"name": "BTC (local)",
			"symbol": "tBTC",
			"total_supply": "0",
			"source": {
				"erc20": {
					"contract_address": "{{.GetEthContractAddr "tBTC"}}"
				}
			}
		},
		"{{.GetVegaContractID "tDAI"}}": {
			"min_lp_stake": "1",
			"decimals": 5,
			"name": "DAI (local)",
			"symbol": "tDAI",
			"total_supply": "0",
			"source": {
				"erc20": {
					"contract_address": "{{.GetEthContractAddr "tDAI"}}"
				}
			}
		},
		"{{.GetVegaContractID "tEURO"}}": {
			"min_lp_stake": "1",
			"decimals": 5,
			"name": "EURO (local)",
			"symbol": "tEURO",
			"total_supply": "0",
			"source": {
				"erc20": {
					"contract_address": "{{.GetEthContractAddr "tEURO"}}"
				}
			}
		},
		"{{.GetVegaContractID "tUSDC"}}": {
			"min_lp_stake": "1",
			"decimals": 5,
			"name": "USDC (local)",
			"symbol": "tUSDC",
			"total_supply": "0",
			"source": {
				"erc20": {
				"contract_address": "{{.GetEthContractAddr "tUSDC"}}"
				}
			}
		},
		"{{.GetVegaContractID "VEGA"}}": {
			"min_lp_stake": "1",
			"decimals": 18,
			"name": "Vega",
			"symbol": "VEGA",
			"total_supply": "64999723000000000000000000",
			"source": {
				"erc20": {
				"contract_address": "{{.GetEthContractAddr "VEGA"}}"
				}
			}
		}
	  },
	  "network": {
		"ReplayAttackThreshold": 30
	  },
	  "network_parameters": {
		"blockchains.ethereumConfig": "{\"network_id\": \"{{ .NetworkID }}\", \"chain_id\": \"{{ .ChainID }}\", \"collateral_bridge_contract\": { \"address\": \"{{.GetEthContractAddr "erc20_bridge_1"}}\" }, \"confirmations\": 3, \"staking_bridge_contract\": { \"address\": \"{{.GetEthContractAddr "staking_bridge"}}\", \"deployment_block_height\": 0}, \"token_vesting_contract\": { \"address\": \"{{.GetEthContractAddr "erc20_vesting"}}\", \"deployment_block_height\": 0 }}",
		"governance.proposal.asset.minClose": "2s",
		"governance.proposal.asset.minEnact": "2s",
		"governance.proposal.asset.requiredParticipation": "0.00000000000000000000000015",
		"governance.proposal.market.minClose": "2s",
		"governance.proposal.market.minEnact": "2s",
		"governance.proposal.market.requiredParticipation": "0.00000000000000000000000015",
		"governance.proposal.updateMarket.minClose": "2s",
		"governance.proposal.updateMarket.minEnact": "2s",
		"governance.proposal.updateMarket.requiredParticipation": "0.00000000000000000000000015",
		"governance.proposal.updateNetParam.minClose": "2s",
		"governance.proposal.updateNetParam.minEnact": "2s",
		"governance.proposal.updateNetParam.requiredParticipation": "0.00000000000000000000000015",
		"market.auction.minimumDuration": "3s",
		"market.fee.factors.infrastructureFee": "0.001",
		"market.fee.factors.makerFee": "0.004",
		"market.monitor.price.updateFrequency": "1s",
		"market.liquidity.stakeToCcySiskas": "0.3",
		"market.liquidity.targetstake.triggering.ratio": "0.7",
		"network.checkpoint.timeElapsedBetweenCheckpoints": "10s",
		"reward.asset": "{{.GetVegaContractID "VEGA"}}",
		"reward.staking.delegation.competitionLevel": "3.1",
		"reward.staking.delegation.delegatorShare": "0.883",
		"reward.staking.delegation.maxPayoutPerParticipant": "700000000000000000000",
		"reward.staking.delegation.minimumValidatorStake": "3000000000000000000000",
		"reward.staking.delegation.payoutDelay": "5m",
		"reward.staking.delegation.payoutFraction": "0.007",
		"spam.protection.delegation.min.tokens": "1000000000000000000",
		"spam.protection.max.delegations": "390",
		"spam.protection.max.proposals": "100",
		"spam.protection.max.votes": "100",
		"spam.protection.proposal.min.tokens": "1000000000000000000",
		"spam.protection.voting.min.tokens": "1000000000000000000",
		"validators.delegation.minAmount": "100000000000000000",
		"validators.epoch.length": "5s",
		"validators.vote.required": "0.67"
	  },
	  "network_limits": {
		"propose_asset_enabled": true,
		"propose_asset_enabled_from": "2021-09-01T00:00:00Z",
		"propose_market_enabled": true,
		"propose_market_enabled_from": "2021-09-01T00:00:00Z"
	  }
	},
	"consensus_params": {
	  "block": {
		"time_iota_ms": "1"
	  }
	}
}
  EOH

  node_set "validators" {
    count = 2
    mode = "validator"
    config_templates {
      vega = <<EOH
[API]
	Port = 30{{.NodeNumber}}2
	[API.REST]
			Port = 30{{.NodeNumber}}3

[Blockchain]
	[Blockchain.Tendermint]
		ClientAddr = "tcp://127.0.0.1:266{{.NodeNumber}}7"
		ServerAddr = "0.0.0.0"
		ServerPort = 266{{.NodeNumber}}8
	[Blockchain.Null]
		Port = 31{{.NodeNumber}}1

[EvtForward]
	Level = "Info"
	RetryRate = "1s"

[NodeWallet]
	[NodeWallet.ETH]
		Address = "{{.ETHEndpoint}}"

[Processor]
	[Processor.Ratelimit]
		Requests = 10000
		PerNBlocks = 1
{{if eq .Type "full"}}
[Broker]
	[Broker.Socket]
		Address = "{{.Prefix}}-{{.DataNodePrefix}}{{.NodeNumber}}"
		Port = 3005
		Enabled = true
{{end}}
      EOH
      tendermint = <<EOH
log_level = "debug"

proxy_app = "tcp://127.0.0.1:266{{.NodeNumber}}8"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
laddr = "tcp://0.0.0.0:266{{.NodeNumber}}7"
unsafe = true

[p2p]
laddr = "tcp://0.0.0.0:266{{.NodeNumber}}6"
addr_book_strict = true
max_packet_msg_payload_size = 4096
pex = false
allow_duplicate_ip = false
persistent_peers = "{{range $i, $v := .NodeIDs}}{{if ne $i 0}},{{end}}{{$v}}@127.0.0.1:266{{$i}}6{{end}}"

[mempool]
size = 10000
cache_size = 20000

[consensus]
skip_timeout_commit = false
	    EOH
    }
  }
}
`

func TestParseConfig(t *testing.T) {
	c, err := ParseConfig([]byte(testConfig))
	if err != nil {
		t.Error(err)
	}

	var nValidators, nNonValidators int
	for _, n := range c.Network.Nodes {
		switch nodeMode(n.Mode) {
		case NodeModeValidator:
			nValidators += n.Count
		case NodeModeFull:
			nNonValidators += n.Count
		default:
			log.Fatalf("unknown node mode: %s", n.Mode)
		}
	}

	fmt.Printf("generating init testnet: nValidators: %d, nNonValidators: %d\n\n", nValidators, nNonValidators)
	fmt.Printf("genesis: %s\n", c.Network.GenesisTemplate)

	for _, n := range c.Network.Nodes {
		for i := 0; i < n.Count; i++ {
			fmt.Println("")
			fmt.Printf("rewriting Vega config for number %d, %s \n", i, n.Templates.Vega)
			fmt.Printf("rewriting Tendermint config for number %d, %s \n", i, n.Templates.Tendermint)
		}
	}

	fmt.Println("done...")
}
