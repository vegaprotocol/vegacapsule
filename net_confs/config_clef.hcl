vega_binary_path = "vega"

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
      image = "vegaprotocol/ganache:latest"
      cmd = "ganache-cli"
      args = [
        "--blockTime", "1",
        "--chainId", "1440",
        "--networkId", "1441",
        "-h", "0.0.0.0",
        "-p", "8545",
        "-m", "ozone access unlock valid olympic save include omit supply green clown session",
        "--db", "/app/ganache-db",
      ]
      static_port {
        value = 8545
        to = 8545
      }
      auth_soft_fail = true
    }
  }
  
  genesis_template_file = "./genesis.tmpl"
  smart_contracts_addresses_file = "./public_smart_contracts_addresses.json"

  node_set "validators" {
    count = 2
    mode = "validator"
    
    pre_generate {
      nomad_job "clef" {
        job_template_file = "../jobs/clef.tmpl"
      }
    }

    clef_wallet {
      ethereum_account_addresses = ["0x6fddbbaa5bc0d53d91d05157c05d86864993cd16","0x3168b80d6efc229a635aa0401600074bb40257c4"]
      clef_rpc_address = "http://localhost:855{{ .NodeNumber }}"
    }

    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

    config_templates {
      vega_file = "./node_set_templates/default/vega_validators_clef.tmpl"
      tendermint_file = "./node_set_templates/default/tendermint_validators.tmpl"
    }
  }

  node_set "full" {
    count = 1
    mode = "full"
	  data_node_binary = "data-node"

    config_templates {
      vega_file = "./node_set_templates/default/vega_full.tmpl"
      tendermint_file = "./node_set_templates/default/tendermint_full.tmpl"
      data_node_file = "./node_set_templates/default/data_node_full.tmpl"
    }
  }
}
