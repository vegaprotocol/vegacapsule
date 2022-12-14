package types

import (
	"sort"
)

type NodeSet struct {
	// description: Name that represents a group of the same node sets.
	GroupName string `cty:"group_name"`
	// description: Name of a specific node set in a node sets group.
	Name string `cty:"name"`
	// description: Mode of the node set. Can be `validator` or `full` (full means a non validating node).
	Mode string `cty:"mode"`
	/*
		description: |
			Index is a position and order in which the node set has been generated respective to all other created node sets.
			It goes from 0-N where N is the total number of node sets.
	*/
	Index int `cty:"index"`
	// description: RelativeIndex is a counter relative to current node set group. Related to GroupName.
	RelativeIndex int
	// description: GroupIndex is an index of the group that this node set belongs to. Related to GroupName.
	GroupIndex int

	// description: Information about generated Vega node.
	Vega VegaNode `cty:"vega"`
	// description: Information about generated Tendermint node.
	Tendermint TendermintNode `cty:"tendermint"`
	// description: Information about generated Data node.
	DataNode *DataNode `cty:"data_node"`
	// description: Information about generated Visor instance.
	Visor *Visor
	// description: Jobs that have been started before the nodes were generated.
	PreGenerateJobs []NomadJob
	// description: Pre start probes.
	PreStartProbe *ProbesConfig `hcl:"pre_start_probe,optional"  template:""`
	// description: Stores custom Nomad job definition of this node set.
	NomadJobRaw *string `json:",omitempty"`
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
