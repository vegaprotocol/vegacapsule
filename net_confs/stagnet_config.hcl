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
      static_port {
        value = 8545
        to = 8545
      }
    }
  }

  genesis_template_file = "./genesis.tmpl"
  smart_contracts_addresses_file = "./smart_contracts_addresses.json"

  node_set "validators" {
    count = 2
    mode = "validator"
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"
    nomad_job_template_file = "./jobs/node_set.tmpl"

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
log-level = "info"

proxy-app = "tcp://127.0.0.1:266{{.NodeNumber}}8"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}7"
  unsafe = true
  cors-allowed-origins = ["*"]
  cors-allowed-methods = ["HEAD", "GET", "POST", ]
  cors-allowed-headers = ["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time", ]

[p2p]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}6"
  addr-book-strict = false
  max-packet-msg-payload-size = 4096
  pex = false
  allow-duplicate-ip = true

  persistent-peers = "{{- range $i, $peer := .NodePeers -}}
	  {{- if ne $i 0 }},{{end -}}
	  {{- $peer.ID}}@127.0.0.1:266{{$peer.Index}}6
  {{- end -}}"


[mempool]
  size = 10000
  cache-size = 20000

[consensus]
  skip-timeout-commit = false
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
log-level = "info"

proxy-app = "tcp://127.0.0.1:266{{.NodeNumber}}8"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}7"
  unsafe = true
  cors-allowed-origins = ["*"]
  cors-allowed-methods = ["HEAD", "GET", "POST", ]
  cors-allowed-headers = ["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time", ]

[p2p]
  laddr = "tcp://0.0.0.0:266{{.NodeNumber}}6"
  addr-book-strict = false
  max-packet-msg-payload-size = 4096
  pex = false
  allow-duplicate-ip = true
  persistent-peers = "{{- range $i, $peer := .NodePeers -}}
	  {{- if ne $i 0 }},{{end -}}
	  {{- $peer.ID}}@127.0.0.1:266{{$peer.Index}}6
  {{- end -}}"

[mempool]
  size = 10000
  cache-size = 20000

[consensus]
  skip-timeout-commit = false
EOT
    }
  }
}
