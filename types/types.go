package types

import "github.com/hashicorp/nomad/api"

// description: Represents any generated Capsule service.
type GeneratedService struct {
	// description: Name of the service.
	Name string `cty:"name"`
	// description: Path to home directory of the service.
	HomeDir string `cty:"home_dir"`
	// description: Path to service configuration.
	ConfigFilePath string `cty:"config_file_path"`
}

type Wallet struct {
	GeneratedService
	Network             string
	BinaryPath          string
	TokenPassphrasePath string
}

type Faucet struct {
	GeneratedService
	PublicKey          string
	WalletFilePath     string
	WalletPassFilePath string
	BinaryPath         string
}

type Binary struct {
	GeneratedService
	BinaryPath string
	Args       []string
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

	// description: Name of the Vega wallet.
	VegaWalletName string
	// description: ID of Vega wallet.
	VegaWalletID string
	// description: Public key used from the Vega wallet.
	VegaWalletPublicKey string
	// description: Recovery phrase from the Vega wallet.
	VegaWalletRecoveryPhrase string
	// description: File path of the Vega wallet passphrase.
	VegaWalletPassFilePath string
}

// description: Represents generated Vega node.
type VegaNode struct {
	// description: General information about the node.
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

// description: Represents generated Tendermint node.
type TendermintNode struct {
	// description: General information about the node.
	GeneratedService `cty:"service"`

	// description: ID of the Tendermint node.
	NodeID string `cty:"node_id"`
	// description: File path of the genesis file used to bootstrap the network.
	GenesisFilePath string
	// description: Path to binary used to generate and run the node.
	BinaryPath string
	// description: Generated public key of the Tendermint validator.
	ValidatorPublicKey string
}

type DataNode struct {
	// description: General information about the node.
	GeneratedService `cty:"service"`
	// description: Path to binary used to generate and run the node.
	BinaryPath string
	// description: Unique IPFS swarm key for this network
	UniqueSwarmKey string
}

type Visor struct {
	// description: General information about Visor.
	GeneratedService
	// description: Path to binary used to generate and run the node.
	BinaryPath string
}

type RawJobWithNomadJob struct {
	RawJob   string
	NomadJob *api.Job
}

// description: Represents a raw Nomad job.
type NomadJob struct {
	// description: Binary selected ID - name of the job.
	ID string
	// description: Nomad job definition.
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
