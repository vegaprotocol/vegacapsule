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
    Hosts = [{{range $i, $v := .NonValidators}}{{if eq $v.GroupName "sentry-0" "sentry-1" "sentry-2"}}{{if ne $i 0}},{{end}}"127.0.0.1:30{{$v.Index}}2"{{end}}{{end}}]
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
    docker_service "postgres-1" {
      image = "vegaprotocol/timescaledb:2.8.0-pg14"
      cmd = "postgres"
      args = []
      env = {
        POSTGRES_USER="vega"
        POSTGRES_PASSWORD="vega"
        POSTGRES_DBS="vega0,vega1,vega2,vega3,vega4,vega5,vega6,vega7,vega8,vega9"
      }
      static_port {
        value = 5232
        to = 5432
      }
      auth_soft_fail = true
    }
  }
  
  genesis_template_file = "./genesis.tmpl"
  smart_contracts_addresses_file = "./public_smart_contracts_addresses.json"

  ## We want 3 validator nodes with each having 2 sentry nodes
  node_set "validator-0" {
    count = 1
    mode = "validator"
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

    config_templates {
      vega_file = "./node_set_templates/sentry/vega_validators.tmpl"
      tendermint_file = "./node_set_templates/sentry/tendermint_validators.tmpl"
    }
  }
  node_set "validator-1" {
    count = 1
    mode = "validator"
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

    config_templates {
      vega_file = "./node_set_templates/sentry/vega_validators.tmpl"
      tendermint_file = "./node_set_templates/sentry/tendermint_validators.tmpl"
    }
  }
  node_set "validator-2" {
    count = 1
    mode = "validator"
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

    config_templates {
      vega_file = "./node_set_templates/sentry/vega_validators.tmpl"
      tendermint_file = "./node_set_templates/sentry/tendermint_validators.tmpl"
    }
  }

  ## Create 2 sentry nodes (each with a data-node) for each validator
  node_set "sentry-0" {
    count = 2
    mode = "full"
	  data_node_binary = "data-node"

    config_templates {
      tendermint_file = "./node_set_templates/sentry/tendermint_sentry-0.tmpl"
      vega_file = "./node_set_templates/sentry/vega_sentry.tmpl"
      data_node_file = "./node_set_templates/default/data_node_full_external_postgres.tmpl"
    }
  }
  node_set "sentry-1" {
    count = 2
    mode = "full"
	  data_node_binary = "data-node"

    config_templates {
      tendermint_file = "./node_set_templates/sentry/tendermint_sentry-1.tmpl"
      vega_file = "./node_set_templates/sentry/vega_sentry.tmpl"
      data_node_file = "./node_set_templates/default/data_node_full_external_postgres.tmpl"
    }
  }
  node_set "sentry-2" {
    count = 2
    mode = "full"
	  data_node_binary = "data-node"

    config_templates {
      tendermint_file = "./node_set_templates/sentry/tendermint_sentry-2.tmpl"
      vega_file = "./node_set_templates/sentry/vega_sentry.tmpl"
      data_node_file = "./node_set_templates/default/data_node_full_external_postgres.tmpl"
    }
  }
}
