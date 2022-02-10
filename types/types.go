package types

type VegaNode struct {
	HomeDir                string
	NodeWalletPassFilePath string
}

type TendermintNode struct {
	HomeDir         string
	GenesisFilePath string
}

type DataNode struct {
	HomeDir    string
	BinaryPath string
}

type NodeSet struct {
	Mode       string
	Vega       VegaNode
	Tendermint TendermintNode
	DataNode   *DataNode
}

const (
	NodeModeValidator           = "validator"
	NodeModeFull                = "full"
	NodeWalletChainTypeVega     = "vega"
	NodeWalletChainTypeEthereum = "ethereum"
)
