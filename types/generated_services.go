package types

import "fmt"

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

func (gs GeneratedServices) GetByName(name string) []GeneratedService {
	if gs.Wallet != nil && gs.Wallet.Name == name {
		return []GeneratedService{gs.Wallet.GeneratedService}
	}

	if gs.Faucet != nil && gs.Faucet.Name == name {
		return []GeneratedService{gs.Faucet.GeneratedService}
	}

	if ns, ok := gs.NodeSets[name]; ok {
		out := []GeneratedService{ns.Vega.GeneratedService}
		if ns.DataNode != nil {
			out = append(out, ns.DataNode.GeneratedService)
		}
		return out
	}

	return nil
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
