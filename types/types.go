package types

import (
	"fmt"
	"sort"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/hashicorp/nomad/api"
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
	BinaryPath             string
}

type TendermintNode struct {
	Name            string
	NodeID          string
	HomeDir         string
	GenesisFilePath string
	BinaryPath      string
}

type DataNode struct {
	Name       string
	HomeDir    string
	BinaryPath string
}

type NetworkPathsMapping struct {
	TendermintHome string
	VegaHome       string
	DataNodeHome   *string `json:",omitempty"`

	VegaBinary     string
	DataNodeBinary *string `json:",omitempty"`
}

type CommandRunner struct {
	Name        string
	NomadJobRaw string

	PathsMapping NetworkPathsMapping
}

type RawJobWithNomadJob struct {
	RawJob   string
	NomadJob *api.Job
}

type NomadJob struct {
	ID          string
	NomadJobRaw string
}

type NodeSet struct {
	GroupName           string
	Name                string
	Mode                string
	Index               int
	Vega                VegaNode
	Tendermint          TendermintNode
	DataNode            *DataNode
	NomadJobRaw         *string        `json:",omitempty"`
	RemoteCommandRunner *CommandRunner `json:",omitempty"`
	PreGenerateJobs     []NomadJob
}

// PreGenerateJobsIDs returns pre gen jobs ids per specific node set
func (ns NodeSet) PreGenerateJobsIDs() []string {
	preGenJobsIDs := make([]string, 0, len(ns.PreGenerateJobs))
	for _, preGenJob := range ns.PreGenerateJobs {
		preGenJobsIDs = append(preGenJobsIDs, preGenJob.ID)
	}

	return preGenJobsIDs
}

func (ns NodeSet) JobsIDs() []string {
	result := []string{ns.Name}
	result = append(result, ns.PreGenerateJobsIDs()...)

	if ns.RemoteCommandRunner != nil {
		result = append(result, ns.RemoteCommandRunner.Name)
	}

	return result
}

// PreGenerateRawJobs returns pre gen jobs per specific node set
func (ns NodeSet) PreGenerateRawJobs() []string {
	preGenJobs := make([]string, 0, len(ns.PreGenerateJobs))
	for _, preGenJob := range ns.PreGenerateJobs {
		preGenJobs = append(preGenJobs, preGenJob.NomadJobRaw)
	}

	return preGenJobs
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
	Wallet          *Wallet
	Faucet          *Faucet
	NodeSets        NodeSetMap
	PreGenerateJobs []NomadJob
	CommandRunners  []*CommandRunner
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
		if ns.Mode == NodeModeValidator {
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

type NetworkJobState struct {
	Name    string
	Kind    JobKind
	Running bool
}

type JobStateMap map[string]NetworkJobState

func (jm JobStateMap) ToSliceNames() []string {
	slice := make([]string, 0, len(jm))
	for id := range jm {
		slice = append(slice, id)
	}
	return slice
}

func (jm JobStateMap) ToSlice() []NetworkJobState {
	slice := []NetworkJobState{}
	for _, job := range jm {
		slice = append(slice, job)
	}

	return slice
}

type NetworkJobsFilter func(nj NetworkJobs) NetworkJobs

type NetworkJobs struct {
	CommandRunnersJobs JobStateMap
	NodesSetsJobs      JobStateMap
	ExtraJobs          JobStateMap
	FaucetJob          NetworkJobState
	WalletJob          NetworkJobState
}

func NewNetworkJobs() NetworkJobs {
	return NetworkJobs{
		CommandRunnersJobs: JobStateMap{},
		NodesSetsJobs:      JobStateMap{},
		ExtraJobs:          JobStateMap{},
	}
}

func (nj NetworkJobs) Filter(filters []NetworkJobsFilter) NetworkJobs {
	result := nj

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

func (nj NetworkJobs) ToSlice() []NetworkJobState {
	resources := append(nj.NodesSetsJobs.ToSlice(), nj.CommandRunnersJobs.ToSlice()...)
	resources = append(resources, nj.ExtraJobs.ToSlice()...)

	if nj.FaucetJob.Name != "" {
		resources = append(resources, nj.FaucetJob)
	}

	if nj.WalletJob.Name != "" {
		resources = append(resources, nj.WalletJob)
	}

	return resources
}

func (nj *NetworkJobs) Append(job NetworkJobState) {
	switch job.Kind {
	case JobNodeSet:
		if nj.NodesSetsJobs == nil {
			nj.NodesSetsJobs = JobStateMap{}
		}
		nj.NodesSetsJobs[job.Name] = job
	case JobCommandRunner:
		if nj.CommandRunnersJobs == nil {
			nj.CommandRunnersJobs = JobStateMap{}
		}
		nj.CommandRunnersJobs[job.Name] = job
	case JobWallet:
		nj.WalletJob = job
	case JobFaucet:
		nj.FaucetJob = job
	default:
		if nj.ExtraJobs == nil {
			nj.ExtraJobs = JobStateMap{}
		}
		nj.ExtraJobs[job.Name] = job
	}
}

func (nj *NetworkJobs) RemoveJobs(jobs []NetworkJobState) {
	for _, job := range jobs {
		nj.RemoveJob(job)
	}
}

func (nj *NetworkJobs) RemoveJob(job NetworkJobState) {
	switch job.Kind {
	case JobNodeSet:
		delete(nj.NodesSetsJobs, job.Name)
	case JobCommandRunner:
		delete(nj.CommandRunnersJobs, job.Name)
	case JobWallet:
		nj.WalletJob = NetworkJobState{}
	case JobFaucet:
		nj.FaucetJob = NetworkJobState{}
	default:
		delete(nj.ExtraJobs, job.Name)
	}
}

func (nj NetworkJobs) Exists(jobID string) bool {
	if _, ok := nj.NodesSetsJobs[jobID]; ok {
		return true
	}
	if _, ok := nj.ExtraJobs[jobID]; ok {
		return true
	}
	if _, ok := nj.CommandRunnersJobs[jobID]; ok {
		return true
	}
	if nj.FaucetJob.Name == jobID {
		return true
	}
	if nj.WalletJob.Name == jobID {
		return true
	}

	return false
}

func (nj *NetworkJobs) AddExtraJobs(ids []string, kind JobKind) {
	for _, id := range ids {
		nj.Append(NetworkJobState{
			Name:    id,
			Running: true,
			Kind:    kind,
		})
	}
}

func (nj NetworkJobs) ToSliceNames() []string {
	out := []string{}

	for _, job := range nj.ToSlice() {
		out = append(out, job.Name)
	}

	return out
}

type NodeWalletInfo struct {
	EthereumAddress          string
	EthereumPrivateKey       string
	EthereumClefRPCAddress   string
	VegaWalletPublicKey      string
	VegaWalletRecoveryPhrase string
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

const (
	NodeModeValidator           = "validator"
	NodeModeFull                = "full"
	NodeWalletChainTypeVega     = "vega"
	NodeWalletChainTypeEthereum = "ethereum"
)

type JobKind string

const (
	JobNodeSet       JobKind = "node-set"
	JobCommandRunner JobKind = "command-runner"
	JobPreStart      JobKind = "pre-start"
	JobPostStart     JobKind = "post-start"
	JobFaucet        JobKind = "faucet"
	JobWallet        JobKind = "wallet"
	JobPreGenerate   JobKind = "pre-generate"
)

func FilterNetworkJobsByJobKindIn(kinds []JobKind) NetworkJobsFilter {
	return func(nj NetworkJobs) NetworkJobs {
		result := NetworkJobs{
			ExtraJobs: JobStateMap{},
		}

		for _, kind := range kinds {
			switch kind {
			case JobNodeSet:
				result.NodesSetsJobs = nj.NodesSetsJobs
			case JobCommandRunner:
				result.CommandRunnersJobs = nj.CommandRunnersJobs
			case JobWallet:
				result.WalletJob = nj.WalletJob
			case JobFaucet:
				result.FaucetJob = nj.FaucetJob
			default:
				for jobName, job := range nj.ExtraJobs {
					if job.Kind != kind {
						continue
					}
					result.ExtraJobs[jobName] = job
				}
			}
		}

		return result
	}
}

func FilterNetworkJobsByJobKindNotIn(notWantedKinds []JobKind) NetworkJobsFilter {
	allJobsKinds := []JobKind{JobNodeSet, JobCommandRunner, JobPreStart, JobPostStart, JobFaucet, JobWallet}

	wantedKinds := []JobKind{}
	for _, kind := range allJobsKinds {
		if utils.IndexInSlice(notWantedKinds, kind) == -1 {
			wantedKinds = append(wantedKinds, kind)
		}
	}

	return FilterNetworkJobsByJobKindIn(wantedKinds)
}

func FilterNetworkJobsByNames(wantedNames []string) NetworkJobsFilter {
	return func(nj NetworkJobs) NetworkJobs {
		result := NewNetworkJobs()

		for _, job := range nj.ToSlice() {
			if utils.IndexInSlice(wantedNames, job.Name) == -1 {
				continue
			}

			result.Append(job)
		}

		return result
	}
}
