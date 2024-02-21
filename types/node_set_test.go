package types_test

import (
	"testing"

	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/stretchr/testify/assert"
)

func TestFilterNodeSets(t *testing.T) {
	input := []types.NodeSet{
		{Name: "ns1", GroupName: "ns-g-1"},
		{Name: "ns2", GroupName: "ns-g-1"},
		{Name: "ns3", GroupName: "ns-g-1"},
		{Name: "ns4", GroupName: "ns-g-1"},

		{Name: "ns5", GroupName: "ns-g-2"},
		{Name: "ns6", GroupName: "ns-g-2"},

		{Name: "ns7", GroupName: "ns-g-3"},
	}

	testCases := []struct {
		name            string
		filters         []types.NodeSetFilter
		expectedResults []string
	}{
		{
			name:            "filter by names - multiple name sets found",
			filters:         []types.NodeSetFilter{types.NodeSetFilterByNames([]string{"ns1", "ns2", "ns5", "ns10"})},
			expectedResults: []string{"ns1", "ns2", "ns5"},
		},
		{
			name:            "filter by names - single name passed",
			filters:         []types.NodeSetFilter{types.NodeSetFilterByNames([]string{"ns1"})},
			expectedResults: []string{"ns1"},
		},
		{
			name:            "filter by names - no matches",
			filters:         []types.NodeSetFilter{types.NodeSetFilterByNames([]string{"ns10"})},
			expectedResults: []string{},
		},
		{
			name:            "filter by groups names - multiple name sets found",
			filters:         []types.NodeSetFilter{types.NodeSetFilterByGroupNames([]string{"ns-g-1", "ns-g-3", "ns-g-5"})},
			expectedResults: []string{"ns1", "ns2", "ns3", "ns4", "ns7"},
		},
		{
			name:            "filter by groups names - single name passed",
			filters:         []types.NodeSetFilter{types.NodeSetFilterByGroupNames([]string{"ns-g-1"})},
			expectedResults: []string{"ns1", "ns2", "ns3", "ns4"},
		},
		{
			name:            "filter by groups names - no matches",
			filters:         []types.NodeSetFilter{types.NodeSetFilterByGroupNames([]string{"ns-g-4"})},
			expectedResults: []string{},
		},
		{
			name: "multiple-filters - multiple matches",
			filters: []types.NodeSetFilter{
				types.NodeSetFilterByGroupNames([]string{"ns-g-1"}),
				types.NodeSetFilterByNames([]string{"ns1", "ns2", "ns5", "ns7"}),
			},
			expectedResults: []string{"ns1", "ns2"},
		},
		{
			name: "multiple-filters - multiple single match",
			filters: []types.NodeSetFilter{
				types.NodeSetFilterByGroupNames([]string{"ns-g-3"}),
				types.NodeSetFilterByNames([]string{"ns1", "ns2", "ns5", "ns7"}),
			},
			expectedResults: []string{"ns7"},
		},
		{
			name: "multiple-filters - no match",
			filters: []types.NodeSetFilter{
				types.NodeSetFilterByGroupNames([]string{"ns-g-3"}),
				types.NodeSetFilterByNames([]string{"ns1", "ns2", "ns5"}),
			},
			expectedResults: []string{},
		},
		{
			name:            "no filters",
			filters:         []types.NodeSetFilter{},
			expectedResults: []string{"ns1", "ns2", "ns3", "ns4", "ns5", "ns6", "ns7"},
		},

		{
			name:            "nil filters",
			filters:         nil,
			expectedResults: []string{"ns1", "ns2", "ns3", "ns4", "ns5", "ns6", "ns7"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := types.FilterNodeSets(input, tc.filters...)

			assert.Len(t, res, len(tc.expectedResults))

			resultNames := make([]string, len(res))
			for idx, ns := range res {
				resultNames[idx] = ns.Name
			}

			assert.ElementsMatch(t, tc.expectedResults, resultNames)
		})

	}
}
