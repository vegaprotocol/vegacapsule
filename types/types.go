package types

type VegaNode struct {
	HomeDir                string
	NodeWalletPassFilePath string
}

type TendermintNode struct {
	HomeDir         string
	GenesisFilePath string
}

type NodeSet struct {
	Mode       string
	Vega       VegaNode
	Tendermint TendermintNode
}

const (
	NodeModeValidator           = "validator"
	NodeModeFull                = "full"
	NodeWalletChainTypeVega     = "vega"
	NodeWalletChainTypeEthereum = "ethereum"
)
