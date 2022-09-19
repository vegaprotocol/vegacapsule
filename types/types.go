package types

import "github.com/hashicorp/nomad/api"

type GeneratedService struct {
	Name           string
	HomeDir        string
	ConfigFilePath string
}

type Wallet struct {
	GeneratedService
	Network            string
	PublicKeyFilePath  string
	PrivateKeyFilePath string
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

	VegaWalletPublicKey      string
	VegaWalletRecoveryPhrase string
	VegaWalletName           string
	VegaWalletPassFilePath   string
}

type VegaNode struct {
	GeneratedService

	Mode                   string
	NodeWalletPassFilePath string

	NodeWalletInfo *NodeWalletInfo `json:",omitempty"`
	BinaryPath     string
}

type TendermintNode struct {
	GeneratedService
	NodeID          string
	GenesisFilePath string
	BinaryPath      string
}

type DataNode struct {
	GeneratedService
	BinaryPath string
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

type VegaNodeOutput struct {
	NomadJobName string
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
}
