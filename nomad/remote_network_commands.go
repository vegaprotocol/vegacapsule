package nomad

import (
	"context"
	"io"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

func (runner *CommandExecutor) Execute(ctx context.Context, binary string, args []string, nodeSets []types.NodeSet) (io.Reader, error) {
	command := []command{
		func(pathsMapping types.NetworkPathsMapping) []string {
			return append([]string{binary}, args...)
		},
	}

	return runner.executeCommands(ctx, command, nodeSets)
}

func (runner *CommandExecutor) NetworkUnsafeResetAll(ctx context.Context, nodeSets []types.NodeSet) (io.Reader, error) {
	vegaResetCommand := func(mapping types.NetworkPathsMapping) []string {
		return []string{
			mapping.VegaBinary,
			"unsafe_reset_all",
			"--home", mapping.VegaHome,
		}
	}

	tendermintResetCommand := func(mapping types.NetworkPathsMapping) []string {
		return []string{
			mapping.VegaBinary,
			"tm",
			"unsafe_reset_all",
			"--home", mapping.TendermintHome,
		}
	}

	dataNodeResetCommand := func(mapping types.NetworkPathsMapping) []string {
		// No data node is running on the node
		if utils.EmptyStrPoint(mapping.DataNodeBinary) || utils.EmptyStrPoint(mapping.DataNodeHome) {
			return nil
		}

		return []string{
			*mapping.DataNodeBinary,
			"unsafe_reset_all",
			"--home", *mapping.DataNodeHome,
		}
	}

	return runner.executeCommands(ctx, []command{vegaResetCommand, tendermintResetCommand, dataNodeResetCommand}, nodeSets)
}
