package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var (
	printLastBlockOnly bool
)

var netLastBlockCmd = &cobra.Command{
	Use:   "last-block",
	Short: "Print last network block even if network is not running",
	Long: `Last block is obtained using the data-node last-block command.
Last block can be obtained as long as postgresql for data-nodes are running`,
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("network addresses")
		}

		lastBlock, err := obtainLastNetworkBlock(networkState.GeneratedServices)
		if err != nil {
			return fmt.Errorf("failed to obtain last network block: %w", err)
		}

		return ptintLastBlock(printLastBlockOnly, lastBlock)
	},
}

func init() {
	netLastBlockCmd.PersistentFlags().BoolVar(&printLastBlockOnly,
		"print-last-block-only",
		false,
		"Print only last block, no other output",
	)
}

func obtainLastNetworkBlock(genServices *types.GeneratedServices) (int, error) {
	dataNodeHomes := []string{}
	dataNodeBinary := ""

	for _, ns := range genServices.NodeSets {
		if ns.Mode != types.NodeModeFull || ns.DataNode == nil {
			continue
		}

		dataNodeHomes = append(dataNodeHomes, ns.DataNode.HomeDir)

		if dataNodeBinary == "" && ns.DataNode.BinaryPath != "" {
			dataNodeBinary = ns.DataNode.BinaryPath
		}
	}

	if len(dataNodeHomes) < 1 {
		return 0, fmt.Errorf("no data node is running")
	}

	if dataNodeBinary == "" {
		return 0, fmt.Errorf("failed to obtain data node binary")
	}

	latestBlock := 0
	for _, homePath := range dataNodeHomes {
		lastBlock, err := commands.LastBlock(dataNodeBinary, homePath)
		if err != nil {
			return 0, fmt.Errorf("failed to get last block from data node: %w", err)
		}

		if latestBlock < lastBlock {
			latestBlock = lastBlock
		}
	}

	return latestBlock, nil
}

func ptintLastBlock(printLastBlockOnly bool, lastBlock int) error {
	if printLastBlockOnly {
		fmt.Printf("%d", lastBlock)
		return nil
	}

	fmt.Printf("Last network block is %d", lastBlock)

	return nil
}
