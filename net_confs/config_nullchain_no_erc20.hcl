vega_binary_path = "vega"

network "testnet" {
  ethereum {
    chain_id   = "1440"
    network_id = "1441"
    endpoint   = "ws://127.0.0.1:8545/"
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
	  docker_service "postgres-1" {
      image = "vegaprotocol/timescaledb:2.8.0-pg14"
      cmd = "postgres"
      args = []
      env = {
        POSTGRES_USER="vega"
        POSTGRES_PASSWORD="vega"
        POSTGRES_DBS="vega,vega0,vega1,vega2,vega3,vega4,vega5,vega6,vega7,vega8,vega9,vega10,vega11,vega12,vega13,vega14,vega15,vega16,vega17,vega18,vega19,vega20,vega21,vega22,vega23,vega24,vega25"
      }
      
      static_port {
        value = 5432
        to = 5432
      }
      resources {
        cpu = 600
        memory = 900
      }
      
      volume_mounts = ["${network_home_path}:${network_home_path}"]

      auth_soft_fail = true
    }
  }

  genesis_template_file = "./node_set_templates/nullchain/genesis_no_erc20.json.tmpl"

  node_set "validators" {
    count = 1
    mode = "validator"
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"
	  use_data_node = true

    config_templates {
		vega_file = "./node_set_templates/nullchain/vega_validator.toml"
		tendermint_file = "./node_set_templates/nullchain/tendermint_validator.toml"
    }
  }

  smart_contracts_addresses_file = "./public_smart_contracts_addresses.json"
}
