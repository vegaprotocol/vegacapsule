package types

import (
	"fmt"
	"sort"
)

type VegaNodeOutput struct {
	NomadJobName string
	VegaNode
}

type VegaNode struct {
	Name                   string
	Mode                   string
	HomeDir                string
	NodeWalletPassFilePath string
	NodeWalletInfo         *NodeWalletInfo `json:",omitempty"`
}

type TendermintNode struct {
	Name            string
	NodeID          string
	HomeDir         string
	GenesisFilePath string
}

type DataNode struct {
	Name       string
	HomeDir    string
	BinaryPath string
}

type NodeSet struct {
	GroupName  string
	Name       string
	Mode       string
	Index      int
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

type NodeSetMap map[string]NodeSet

func (nm NodeSetMap) ToSlice() []NodeSet {
	slice := make([]NodeSet, 0, len(nm))
	for _, ns := range nm {
		slice = append(slice, ns)
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Index < slice[j].Index
	})

	return slice
}

type GeneratedServices struct {
	Wallet   *Wallet
	Faucet   *Faucet
	NodeSets NodeSetMap
}

func NewGeneratedServices(w *Wallet, f *Faucet, ns []NodeSet) *GeneratedServices {
	nsm := NodeSetMap{}

	for _, ns := range ns {
		nsm[ns.Name] = ns
	}

	return &GeneratedServices{
		Wallet:   w,
		Faucet:   f,
		NodeSets: nsm,
	}
}

func (gs GeneratedServices) GetNodeSet(name string) (*NodeSet, error) {
	ns, ok := gs.NodeSets[name]
	if !ok {
		return nil, fmt.Errorf("node set with name %q not found", name)
	}
	return &ns, nil
}

func (gs GeneratedServices) ListValidators() []VegaNodeOutput {
	var validators []VegaNodeOutput

	for _, nodeSet := range gs.NodeSets {
		if nodeSet.Mode != NodeModeValidator {
			continue
		}
		validators = append(validators, VegaNodeOutput{
			VegaNode:     nodeSet.Vega,
			NomadJobName: nodeSet.Name,
		})
	}

	return validators
}

type JobIDMap map[string]bool

func (jm JobIDMap) ToSlice() []string {
	slice := make([]string, 0, len(jm))
	for id := range jm {
		slice = append(slice, id)
	}
	return slice
}

type NetworkJobs struct {
	NodesSetsJobIDs JobIDMap
	ExtraJobIDs     JobIDMap
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
