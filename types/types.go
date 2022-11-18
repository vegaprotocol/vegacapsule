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
	BinaryPath         string
}

// description: Information about node wallets.
type NodeWalletInfo struct {
	/*
		description: Ethereum account address.
		note: Only available when Key Store wallet is used.
	*/
	EthereumAddress string

	EthereumPrivateKey string
	// description: Path to file where Ethereum wallet key is stored.
	EthereumPassFilePath string
	/*
		description: Address of Clef wallet.
		note: Only available when Clef wallet is used.
	*/
	EthereumClefRPCAddress string

	VegaWalletID             string
	VegaWalletPublicKey      string
	VegaWalletRecoveryPhrase string
	VegaWalletName           string
	VegaWalletPassFilePath   string
}

// description: Represents generated Vega node.
type VegaNode struct {
	// description: Path to binary used to generate and run the node.
	GeneratedService `cty:"service"`

	// description: Mode of the node - `validator` or `full`.
	Mode string `cty:"mode"`

	/*
		description: Path to generated node wallet passphrase file.
		note: Only present if `mode = validator`.
	*/
	NodeWalletPassFilePath string

	/*
		description: Information about generated/imported node wallets.
		note: Only present if `mode = validator`.
	*/
	NodeWalletInfo *NodeWalletInfo `json:",omitempty"`
	// description: Path to binary used to generate and run the node.
	BinaryPath string
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
