vega_binary_path = "vega"

network "testnet" {
  ethereum {
    chain_id   = "1440"
    network_id = "1441"
    endpoint   = "http://127.0.0.1:8545/"
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
		"snapshot.interval.length": "5",
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
    count = 4
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

[NodeWallet]
	[NodeWallet.ETH]
		Address = "{{.ETHEndpoint}}"

[Processor]
	[Processor.Ratelimit]
		Requests = 10000
		PerNBlocks = 1

[Snapshot]
	Level = "DEBUG"
	KeepRecent = 10
	RetryLimit = 5
	Storage = "GOLevelDB"
	DBPath = ""
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
  cors_allowed_origins = ["*"]
  cors-allowed-methods = ["HEAD", "GET", "POST", ]
  cors-allowed-headers = ["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time", ]

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

    smart_contracts_addresses = <<EOH
{
	"addr0": {
		"priv": "adef89153e4bd6b43876045efdd6818cec359340683edaec5e8588e635e8428b",
		"pub": "0xb89A165EA8b619c14312dB316BaAa80D2a98B493"
	},
	"MultisigControl": {
		"Ethereum": "0xa956B5c58B4Ac8Dd1D44Ade3e8972A16e9C917E4"
	},
	"ERC20_Asset_Pool": {
		"Ethereum": "0x3EA59801698c6820328597F26d29fC3EaAa17AcA"
	},
	"erc20_bridge_1": {
		"Ethereum": "0x0858D9BD11A4F6Bae8b979402550CA6c6ddB8332"
	},
	"erc20_bridge_2": {
		"Ethereum": "0x846087f262859fe6604e2e9f787a9F3f39296Ff8"
	},
	"tBTC": {
		"Ethereum": "0xc6a6000d740707edc35f75f42447320B60450c04",
		"Vega": "0x5cfa87844724df6069b94e4c8a6f03af21907d7bc251593d08e4251043ee9f7c"
	},
	"tDAI": {
		"Ethereum": "0xE25F12E386Cd7F84c41B5210504d9743A35Badda",
		"Vega": "0x6d9d35f657589e40ddfb448b7ad4a7463b66efb307527fedd2aa7df1bbd5ea61"
	},
	"tEURO": {
		"Ethereum": "0x7c23d674fED4500103A0b7e05b4A0da17291FCE9",
		"Vega": "0x8b52d4a3a4b0ffe733cddbc2b67be273816cfeb6ca4c8b339bac03ffba08e4e4"
	},
	"tUSDC": {
		"Ethereum": "0xD76Bd796e117D54044E616ae42A3577256B601D1",
		"Vega": "0x993ed98f4f770d91a796faab1738551193ba45c62341d20597df70fea6704ede"
	},
	"VEGA": {
		"Ethereum": "0xBC944ba38753A6fCAdd634Be98379330dbaB3Eb8",
		"Vega": "0xb4f2726571fbe8e33b442dc92ed2d7f0d810e21835b7371a7915a365f07ccd9b"
	},
	"VEGAv1": {
		"Ethereum": "0xB69a81EE133d8c4dC4AeCB30af93bC8698118ccE",
		"Vega": "0xc1607f28ec1d0a0b36842c8327101b18de2c5f172585870912f5959145a9176c"
	},
	"erc20_vesting": {
		"Ethereum": "0xB9f84835F00C0E4f494C51C945863109cF80754A"
	},
	"staking_bridge": {
		"Ethereum": "0xE4c9fB5955bAa8a7965D000afdCEFF25cfe0E8a3"
	}
}
EOH
}
