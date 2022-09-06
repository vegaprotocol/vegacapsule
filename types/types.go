package types

import (
	"fmt"
	"sort"

	"github.com/hashicorp/nomad/api"
)

const (
	NodeModeValidator = "validator"
	NodeModeFull      = "full"

	NodeWalletChainTypeTendermint = "tendermint"
	NodeWalletChainTypeVega       = "vega"
	NodeWalletChainTypeEthereum   = "ethereum"
)

type VegaNodeOutput struct {
	NomadJobName string
	VegaNode
}

type VegaNode struct {
	Name                   string
	Mode                   string
	HomeDir                string
	ConfigFilePath         string
	NodeWalletPassFilePath string

	NodeWalletInfo *NodeWalletInfo `json:",omitempty"`
	BinaryPath     string
}

type TendermintNode struct {
	Name            string
	NodeID          string
	HomeDir         string
	GenesisFilePath string
	BinaryPath      string
}

type DataNode struct {
	Name           string
	HomeDir        string
	BinaryPath     string
	ConfigFilePath string
}

type Visor struct {
	Name       string
	HomeDir    string
	BinaryPath string
}

type RawJobWithNomadJob struct {
	RawJob   string
	NomadJob *api.Job
}

type NomadJob struct {
	ID          string
	NomadJobRaw string
	Ports       []int64
}

type NodeSet struct {
	GroupName string
	Name      string
	Mode      string
	// Index is a node set counter over all created node sets.
	Index int
	// RelativeIndex is a counter relative to current node set group. Related to GroupName.
	RelativeIndex int
	// GroupIndex is a index of the group where this node set belongs to. Related to GroupName.
	GroupIndex      int
	Vega            VegaNode
	Tendermint      TendermintNode
	DataNode        *DataNode
	Visor           *Visor
	NomadJobRaw     *string `json:",omitempty"`
	PreGenerateJobs []NomadJob
}

// PreGenerateJobsIDs returns pre gen jobs ids per specific node set
func (ns NodeSet) PreGenerateJobsIDs() []string {
	preGenJobsIDs := make([]string, 0, len(ns.PreGenerateJobs))
	for _, preGenJob := range ns.PreGenerateJobs {
		preGenJobsIDs = append(preGenJobsIDs, preGenJob.ID)
	}

	return preGenJobsIDs
}

// PreGenerateRawJobs returns pre gen jobs per specific node set
func (ns NodeSet) PreGenerateRawJobs() []string {
	preGenJobs := make([]string, 0, len(ns.PreGenerateJobs))
	for _, preGenJob := range ns.PreGenerateJobs {
		preGenJobs = append(preGenJobs, preGenJob.NomadJobRaw)
	}

	return preGenJobs
}

func (ns NodeSet) IsValidator() bool {
	return ns.Mode == NodeModeValidator
}

type Wallet struct {
	Name               string
	HomeDir            string
	Network            string
	ConfigFilePath     string
	PublicKeyFilePath  string
	PrivateKeyFilePath string
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
	Wallet          *Wallet
	Faucet          *Faucet
	NodeSets        NodeSetMap
	PreGenerateJobs []NomadJob
}

func NewGeneratedServices(w *Wallet, f *Faucet, ns []NodeSet) *GeneratedServices {
	nsm := NodeSetMap{}
	preGenJobs := []NomadJob{}

	for _, ns := range ns {
		nsm[ns.Name] = ns
		preGenJobs = append(preGenJobs, ns.PreGenerateJobs...)
	}

	return &GeneratedServices{
		Wallet:          w,
		Faucet:          f,
		NodeSets:        nsm,
		PreGenerateJobs: preGenJobs,
	}
}

func (gs GeneratedServices) GetNodeSet(name string) (*NodeSet, error) {
	ns, ok := gs.NodeSets[name]
	if !ok {
		return nil, fmt.Errorf("node set with name %q not found", name)
	}
	return &ns, nil
}

// PreGenerateJobsIDs returns pre gen jobs ids across all node sets
func (gs GeneratedServices) PreGenerateJobsIDs() []string {
	preGenJobsIDs := make([]string, 0, len(gs.PreGenerateJobs))
	for _, preGenJob := range gs.PreGenerateJobs {
		preGenJobsIDs = append(preGenJobsIDs, preGenJob.ID)
	}

	return preGenJobsIDs
}

type NodeSetFilter func(ns NodeSet) bool

func NodeSetFilterByNames(names []string) NodeSetFilter {
	return func(ns NodeSet) bool {
		for _, expectedName := range names {
			if ns.Name == expectedName {
				return true
			}
		}
		return false
	}
}

func NodeSetFilterByGroupNames(names []string) NodeSetFilter {
	return func(ns NodeSet) bool {
		for _, expectedName := range names {
			if ns.GroupName == expectedName {
				return true
			}
		}
		return false
	}
}

func NodeSetFilterByGroupName(groupName string) NodeSetFilter {
	return func(ns NodeSet) bool {
		return ns.GroupName == groupName
	}
}

func FilterNodeSets(nodeSets []NodeSet, filters ...NodeSetFilter) []NodeSet {
	var out []NodeSet

	for _, ns := range nodeSets {
		func() {
			for _, filterFunc := range filters {
				if filterFunc == nil {
					return
				}

				if !filterFunc(ns) {
					return
				}
			}

			out = append(out, ns)
		}()
	}

	return out
}

func (gs GeneratedServices) GetNodeSetsByGroupName(groupName string) []NodeSet {
	return FilterNodeSets(gs.NodeSets.ToSlice(), NodeSetFilterByGroupName(groupName))
}

func (gs GeneratedServices) GetValidators() []NodeSet {
	var out []NodeSet

	for _, ns := range gs.NodeSets {
		if ns.IsValidator() {
			out = append(out, ns)
		}
	}

	return out
}

func (gs GeneratedServices) GetNonValidators() []NodeSet {
	var out []NodeSet

	for _, ns := range gs.NodeSets {
		if !ns.IsValidator() {
			out = append(out, ns)
		}
	}

	return out
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

func (nj NetworkJobs) Exists(jobID string) bool {
	if _, ok := nj.NodesSetsJobIDs[jobID]; ok {
		return true
	}
	if _, ok := nj.ExtraJobIDs[jobID]; ok {
		return true
	}
	if nj.FaucetJobID == jobID {
		return true
	}
	if nj.WalletJobID == jobID {
		return true
	}

	return false
}

func (nj NetworkJobs) AddExtraJobIDs(ids []string) {
	if nj.ExtraJobIDs == nil {
		nj.ExtraJobIDs = JobIDMap{}
	}

	for _, id := range ids {
		nj.ExtraJobIDs[id] = true
	}
}

func (nj NetworkJobs) ToSlice() []string {
	out := append(nj.NodesSetsJobIDs.ToSlice(), nj.ExtraJobIDs.ToSlice()...)

	if nj.FaucetJobID != "" {
		out = append(out, nj.FaucetJobID)
	}

	if nj.WalletJobID != "" {
		out = append(out, nj.WalletJobID)
	}

	return out
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

type SmartContractsInfo struct {
	MultisigControl struct {
		EthereumAddress string `json:"Ethereum"`
	} `json:"MultisigControl"`
	EthereumOwner struct {
		Public  string `json:"pub"`
		Private string `json:"priv"`
	} `json:"addr0"`
}
