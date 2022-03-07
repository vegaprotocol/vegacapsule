output_dir             = "./testnet"
vega_binary_path       = "vega"
prefix                 = "st-local"
node_dir_prefix        = "node"
tendermint_node_prefix = "tendermint"
vega_node_prefix       = "vega"
data_node_prefix       = "data"
wallet_prefix          = "wallet"
faucet_prefix          = "faucet"

network "testnet" {
	ethereum {
    chain_id   = "1440"
    network_id = "1441"
    endpoint   = "http://127.0.0.1:8545/"
  }
  
  faucet "faucet-1" {
	  wallet_pass = "f4uc3tw4ll3t-v3g4-p4ssphr4e3"

	  template = <<-EOT
[Node]
  Port = 3002
  IP = "127.0.0.1"
EOT
  }

  wallet "wallet-1" {
    binary = "vegawallet"
    
    template = <<-EOT
Name = "DV"
Level = "info"
TokenExpiry = "168h0m0s"
Port = 1789
Host = "0.0.0.0"

[API]
  [API.GRPC]
    Hosts = [{{range $i, $v := .Validators}}{{if ne $i 0}},{{end}}"127.0.0.1:30{{$i}}2"{{end}}]
    Retries = 5
EOT
  }

  pre_start {
    docker_service "ganache-1" {
      image = "ghcr.io/vegaprotocol/devops-infra/ganache:latest"
      cmd = "ganache-cli"
      args = [
        "--blockTime", "1",
        "--chainId", "1440",
        "--networkId", "1441",
        "-h", "0.0.0.0",
        "-p", "8545",
        "-m", "cherry manage trip absorb logic half number test shed logic purpose rifle",
        "--db", "/app/ganache-db",
      ]
      static_port = 8545
    }
  }

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
		"blockchains.ethereumConfig": "{\"network_id\": \"{{ .NetworkID }}\", \"chain_id\": \"{{ .ChainID }}\", \"collateral_bridge_contract\": { \"address\": \"{{.GetEthContractAddr "erc20_bridge_1"}}\" }, \"confirmations\": 3, \"staking_bridge_contract\": { \"address\": \"{{.GetEthContractAddr "staking_bridge"}}\", \"deployment_block_height\": 0}, \"token_vesting_contract\": { \"address\": \"{{.GetEthContractAddr "erc20_vesting"}}\", \"deployment_block_height\": 0 }, \"multisig_control_contract\": { \"address\": \"{{.GetEthContractAddr "MultisigControl"}}\", \"deployment_block_height\": 0 }}",
		"governance.proposal.asset.minClose": "5s",
		"governance.proposal.asset.minEnact": "5s",
		"governance.proposal.asset.requiredParticipation": "0.00000015",
		"governance.proposal.market.minClose": "5s",
		"governance.proposal.market.minEnact": "5s",
		"governance.proposal.market.requiredParticipation": "0.00000015",
		"governance.proposal.updateMarket.minClose": "5s",
		"governance.proposal.updateMarket.minEnact": "5s",
		"governance.proposal.updateMarket.requiredParticipation": "0.00000015",
		"governance.proposal.updateNetParam.minClose": "5s",
		"governance.proposal.updateNetParam.minEnact": "5s",
		"governance.proposal.updateNetParam.requiredParticipation": "0.00000015",
		"market.auction.minimumDuration": "1s",
		"market.fee.factors.infrastructureFee": "0.001",
		"market.fee.factors.makerFee": "0.004",
		"market.monitor.price.updateFrequency": "4s",
		"market.liquidity.stakeToCcySiskas": "0.2",
		"market.liquidity.targetstake.triggering.ratio": "0",
		"network.checkpoint.timeElapsedBetweenCheckpoints": "10s",
		"reward.staking.delegation.competitionLevel": "1.1",
		"reward.staking.delegation.delegatorShare": "0.883",
		"reward.staking.delegation.maxPayoutPerParticipant": "0",
		"reward.staking.delegation.minimumValidatorStake": "0",
		"reward.staking.delegation.payoutDelay": "0s",
		"reward.staking.delegation.payoutFraction": ".1",
		"spam.protection.delegation.min.tokens": "1000000000000000000",
		"spam.protection.max.delegations": "3",
		"spam.protection.max.proposals": "3",
		"spam.protection.max.votes": "3",
		"spam.protection.proposal.min.tokens": "1000000000000000000",
		"spam.protection.voting.min.tokens": "1000000000000000000",
		"validators.delegation.minAmount": "1000000000000000000",
		"validators.epoch.length": "5s",
		"validators.vote.required": "0.67",
		"reward.staking.delegation.minValidators": "3",
		"reward.staking.delegation.optimalStakeMultiplier": "3.0",
		"reward.asset": "{{.GetVegaContractID "VEGA"}}",
		"governance.proposal.asset.maxClose": "8760h0m0s",
		"governance.proposal.asset.maxEnact": "8760h0m0s",
		"governance.proposal.asset.minProposerBalance": "2000000000000000000",
		"governance.proposal.asset.minVoterBalance": "2000000000000000000",
		"governance.proposal.asset.requiredMajority": "0.66",
		"governance.proposal.market.maxClose": "8760h0m0s",
		"governance.proposal.market.maxEnact": "8760h0m0s",
		"governance.proposal.market.minProposerBalance": "2000000000000000000",
		"governance.proposal.market.minVoterBalance": "2000000000000000000",
		"governance.proposal.market.requiredMajority": "0.66",
		"governance.proposal.updateMarket.maxClose": "8760h0m0s",
		"governance.proposal.updateMarket.maxEnact": "8760h0m0s",
		"governance.proposal.updateMarket.minProposerBalance": "2000000000000000000",
		"governance.proposal.updateMarket.minVoterBalance": "2000000000000000000",
		"governance.proposal.updateMarket.requiredMajority": "0.66",
		"governance.proposal.updateNetParam.maxClose": "8760h0m0s",
		"governance.proposal.updateNetParam.maxEnact": "8760h0m0s",
		"governance.proposal.updateNetParam.minProposerBalance": "2000000000000000000",
		"governance.proposal.updateNetParam.minVoterBalance": "2000000000000000000",
		"governance.proposal.updateNetParam.requiredMajority": "0.66",
		"market.auction.maximumDuration": "168h0m0s",
		"market.liquidity.bondPenaltyParameter": "0.1",
		"market.liquidity.maximumLiquidityFeeFactorLevel": "0.3",
		"market.liquidityProvision.shapes.maxSize": "10",
		"market.margin.scalingFactors": "{\"search_level\": 1.1, \"initial_margin\": 1.2, \"collateral_release\": 1.4}",
		"market.stake.target.scalingFactor": "5",
		"market.stake.target.timeWindow": "10s",
		"market.value.windowLength": "60s",
		"network.checkpoint.networkEndOfLifeDate": "2021-12-15T17:00:00Z",
		"governance.proposal.freeform.minClose": "5s",
		"governance.proposal.freeform.requiredParticipation": "0.0000000000000000000000000015"
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
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

    config_templates {

// ============================
// ===== VegaNode Config ======
// ============================

      vega = <<-EOT
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
	{{if .FaucetPublicKey}}
	BlockchainQueueAllowlist = ["{{ .FaucetPublicKey }}"]
	{{end}}

[NodeWallet]
	[NodeWallet.ETH]
		Address = "{{.ETHEndpoint}}"

[Processor]
	[Processor.Ratelimit]
		Requests = 10000
		PerNBlocks = 1
EOT

// ============================
// ==== Tendermint Config =====
// ============================

	  tendermint = <<-EOT
log_level = "info"

proxy_app = "tcp://127.0.0.1:266{{.NodeNumber}}8"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}7"
  unsafe = true

[p2p]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}6"
  addr_book_strict = false
  max_packet_msg_payload_size = 4096
  pex = false
  allow_duplicate_ip = true

  persistent_peers = "{{- range $i, $peer := .NodePeers -}}
	  {{- if ne $i 0 }},{{end -}}
	  {{- $peer.ID}}@127.0.0.1:266{{$peer.Index}}6
  {{- end -}}"


[mempool]
  size = 10000
  cache_size = 20000

[consensus]
  skip_timeout_commit = false
EOT
    }
  }

  node_set "full" {
    count = 1
    mode = "full"
	  data_node_binary = "data-node"

    config_templates {

// ============================
// ===== VegaNode Config ======
// ============================

      vega = <<-EOT
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

[Broker]
  [Broker.Socket]
    Port = 30{{.NodeNumber}}5
    Enabled = true
EOT

// ============================
// ===== DataNode Config ======
// ============================

      data_node = <<-EOT

GatewayEnabled = true
[SqlStore]
  Port = 5{{.NodeNumber}}32

[API]
  Level = "Info"
  Port = 30{{.NodeNumber}}7
  CoreNodeGRPCPort = 30{{.NodeNumber}}2

[Pprof]
  Level = "Info"
  Enabled = true
  Port = 6{{.NodeNumber}}60
  ProfilesDir = "{{.NodeHomeDir}}"

[Gateway]
  Level = "Info"
  [Gateway.Node]
    Port = 30{{.NodeNumber}}7
  [Gateway.GraphQL]
    Port = 30{{.NodeNumber}}8
  [Gateway.REST]
    Port = 30{{.NodeNumber}}9
	
[Metrics]
  Level = "Info"
  Timeout = "5s"
  Port = 21{{.NodeNumber}}2
  Enabled = false
[Broker]
  Level = "Info"
  UseEventFile = false
  [Broker.SocketConfig]
    Port = 30{{.NodeNumber}}5

EOT

// ============================
// ==== Tendermint Config =====
// ============================

	  tendermint = <<-EOT
log_level = "info"

proxy_app = "tcp://127.0.0.1:266{{.NodeNumber}}8"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}7"
  unsafe = true

[p2p]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}6"
  addr_book_strict = false
  max_packet_msg_payload_size = 4096
  pex = false
  allow_duplicate_ip = true
  persistent_peers = "{{- range $i, $peer := .NodePeers -}}
	  {{- if ne $i 0 }},{{end -}}
	  {{- $peer.ID}}@127.0.0.1:266{{$peer.Index}}6
  {{- end -}}"

[mempool]
  size = 10000
  cache_size = 20000

[consensus]
  skip_timeout_commit = false
EOT
    }
  }
}
