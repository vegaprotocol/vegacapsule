package types

type VegaNode struct {
	Mode                   string
	HomeDir                string
	NodeWalletPassFilePath string
	NodeWalletInfo         *NodeWalletInfo `json:",omitempty"`
}

type TendermintNode struct {
	NodeID          string
	HomeDir         string
	GenesisFilePath string
}

type DataNode struct {
	HomeDir    string
	BinaryPath string
}

type NodeSet struct {
	GroupName  string
	Name       string
	Mode       string
	Vega       VegaNode
	Tendermint TendermintNode
	DataNode   *DataNode
}

type Wallet struct {
	Name                  string
	HomeDir               string
	Network               string
	ServiceConfigFilePath string
	PublicKeyFilePath     string
	PrivateKeyFilePath    string
}

type Faucet struct {
	Name               string
	HomeDir            string
	PublicKey          string
	ConfigFilePath     string
	WalletFilePath     string
	WalletPassFilePath string
}

type GeneratedServices struct {
	Wallet   *Wallet
	Faucet   *Faucet
	NodeSets []NodeSet
}

type NetworkJobs struct {
	NodesSetsJobIDs []string
	ExtraJobIDs     []string
	FaucetJobID     string
	WalletJobID     string
}

type NodeWalletInfo struct {
	EthereumAddress          string
	EthereumPrivateKey       string
	VegaWalletPublicKey      string
	VegaWalletRecoveryPhrase string
}

const (
	NodeModeValidator           = "validator"
	NodeModeFull                = "full"
	NodeWalletChainTypeVega     = "vega"
	NodeWalletChainTypeEthereum = "ethereum"
)
