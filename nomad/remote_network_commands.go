package nomad

import (
	"context"
	"io"

	"code.vegaprotocol.io/vegacapsule/types"
)

func (runner *CommandExecutor) Execute(ctx context.Context, binary string, args []string, nodeSets []types.NodeSet) (io.Reader, error) {
	command := []commandCallback{
		func(pathsMapping types.NetworkPathsMapping) []string {
			return append([]string{binary}, args...)
		},
	}

	return runner.executeCallbacks(ctx, command, nodeSets)
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

	return runner.executeCallbacks(ctx, []commandCallback{vegaResetCommand, tendermintResetCommand}, nodeSets)
}
