package types

import "github.com/hashicorp/nomad/api"

type GeneratedService struct {
	Name           string `cty:"name"`
	HomeDir        string `cty:"home_dir"`
	ConfigFilePath string `cty:"config_file_path"`
}

type Wallet struct {
	GeneratedService
	Network            string
	PublicKeyFilePath  string
	PrivateKeyFilePath string
	BinaryPath         string
}

type Faucet struct {
	GeneratedService
	PublicKey          string
	WalletFilePath     string
	WalletPassFilePath string
}

type NodeWalletInfo struct {
	EthereumPassFilePath   string
	EthereumAddress        string
	EthereumPrivateKey     string
	EthereumClefRPCAddress string

	VegaWalletID             string
	VegaWalletPublicKey      string
	VegaWalletRecoveryPhrase string
	VegaWalletName           string
	VegaWalletPassFilePath   string
}

type VegaNode struct {
	GeneratedService `cty:"service"`

	Mode                   string `cty:"mode"`
	NodeWalletPassFilePath string

	NodeWalletInfo *NodeWalletInfo `json:",omitempty"`
	BinaryPath     string
}

type TendermintNode struct {
	GeneratedService   `cty:"service"`
	NodeID             string `cty:"node_id"`
	GenesisFilePath    string
	BinaryPath         string
	ValidatorPublicKey string
}

type DataNode struct {
	GeneratedService `cty:"service"`
	BinaryPath       string
}

type Visor struct {
	GeneratedService
	BinaryPath string
}

type RawJobWithNomadJob struct {
	RawJob   string
	NomadJob *api.Job
}

type NomadJob struct {
	ID          string
	NomadJobRaw string
}

type TendermintOutput struct {
	NodeID             string
	ValidatorPublicKey string
}

type VegaNodeOutput struct {
	NomadJobName string
	Tendermint   TendermintOutput
	VegaNode
}

type SmartContractsInfo struct {
	MultisigControl struct {
		EthereumAddress string `json:"Ethereum"`
	} `json:"MultisigControl"`
	EthereumOwner struct {
		Public  string `json:"pub"`
		Private string `json:"priv"`
	} `json:"addr0"`
	ERC20Bridge struct {
		EthereumAddress string `json:"Ethereum"`
	} `json:"erc20_bridge_1"`
	StakingBridge struct {
		EthereumAddress string `json:"Ethereum"`
	} `json:"staking_bridge"`
}

type SmartContractsToken struct {
	EthereumAddress string `json:"Ethereum"`
	VegaAddress     string `json:"Vega"`
}

type HTTPProbe struct {
	URL string `hcl:"url" template:""`
}

type TCPProbe struct {
	Address string `hcl:"address" template:""`
}

type PostgresProbe struct {
	Connection string `hcl:"connection" template:""`
	Query      string `hcl:"query" template:""`
}

type ProbesConfig struct {
	HTTP     *HTTPProbe     `hcl:"http,block" template:""`
	TCP      *TCPProbe      `hcl:"tcp,block" template:""`
	Postgres *PostgresProbe `hcl:"postgres,block" template:""`
}
