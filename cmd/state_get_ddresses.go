package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var stateGetAddresses = &cobra.Command{
	Use:   "get-smartcontracts-addresses",
	Short: "Print smart contract addresses and keys passes to vegacapsule as a config parameter",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state get-smartcontracts-addresses")
		}

		fmt.Println(netState.Config.Network.SmartContractsAddresses)

		return nil
	},
}
