package types

import "sort"

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
	PreStartProbe   string `hcl:"pre_start_probe,optional"  template:""`
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
