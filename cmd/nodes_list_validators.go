package cmd

import (
	"encoding/json"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/state"

	"github.com/spf13/cobra"
)

var nodesLsValidatorsCmd = &cobra.Command{
	Use:   "ls-validators",
	Short: "Lists validators from node sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed list validators: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes ls-validators")
		}

		validators := networkState.GeneratedServices.ListValidators()

		validatorsJson, err := json.MarshalIndent(validators, "", "\t")
		if err != nil {
			return fmt.Errorf("failed to marshal validators info: %w", err)
		}

		fmt.Println(string(validatorsJson))
		return nil
	},
}
