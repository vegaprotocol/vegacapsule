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

  pre_start {}

  genesis_template_file = "./genesis.tmpl"

  node_set "validators" {
    count = 1
    mode = "validator"
    node_wallet_pass = "n0d3w4ll3t-p4ssphr4e3"
    vega_wallet_pass = "w4ll3t-p4ssphr4e3"
    ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"
	use_data_node = true
    nomad_job_template_file = "./jobs/node_set_nullchain.tmpl"

    config_templates {

// ============================
// ===== VegaNode Config ======
// ============================

      vega = <<-EOT
[Admin]
  [Admin.Server]
    SocketPath = "/tmp/vega-{{.NodeNumber}}.sock"
    Enabled = true

[API]
	Port = 30{{.NodeNumber}}2
	[API.REST]
			Port = 30{{.NodeNumber}}3

[Blockchain]
    ChainProvider = "nullchain"
    [Blockchain.Tendermint]
        RPCAddr = "tcp://127.0.0.1:266{{.NodeNumber}}7"
    [Blockchain.Null]
        Level = "Debug"
        BlockDuration = "1s"
        TransactionsPerBlock = 1
        IP = "0.0.0.0"
        Port = 31{{.NodeNumber}}1
        GenesisFile = "{{.NodeSet.Tendermint.GenesisFilePath}}"

[EvtForward]
	Level = "Info"
	RetryRate = "1s"
	{{if .FaucetPublicKey}}
	BlockchainQueueAllowlist = ["{{ .FaucetPublicKey }}"]
	{{end}}

[Ethereum]
  RPCEndpoint = "{{.ETHEndpoint}}"

[Broker]
  [Broker.Socket]
    Port = 30{{.NodeNumber}}5
    Enabled = true

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
moniker = "{{.TendermintNodePrefix}}-{{.NodeNumber}}"

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


  smart_contracts_addresses = <<-EOT
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
EOT
}
