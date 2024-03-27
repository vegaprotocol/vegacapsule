package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/state"

	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Manages vegacapsule state",
}

func init() {
	stateCmd.AddCommand(stateGetSmartContractsAddressesCmd)
	stateCmd.AddCommand(stateGetSecondarySmartContractsAddressesCmd)
}

var stateGetSmartContractsAddressesCmd = &cobra.Command{
	Use:   "get-smartcontracts-addresses",
	Short: "Print primary smartcontracts addresses and keys passed to vegacapsule as a config parameter",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state get-smartcontracts-addresses")
		}

		fmt.Println(*netState.Config.Network.SmartContractsAddresses)

		return nil
	},
}

var stateGetSecondarySmartContractsAddressesCmd = &cobra.Command{
	Use:   "get-secondary-smartcontracts-addresses",
	Short: "Print secondary smartcontracts addresses and keys passed to vegacapsule as a config parameter",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state get-secondary-smartcontracts-addresses")
		}

		fmt.Println(*netState.Config.Network.SecondarySmartContractsAddresses)

		return nil
	},
}
