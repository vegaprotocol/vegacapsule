package cmd

import (
	"path"

	"code.vegaprotocol.io/vegacapsule/logscollector"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tool for logs extracting",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state get-smartcontracts-addresses")
		}

		if logsOffset == 0 {
			logsOffset = -1
		}

		return logscollector.Tail(
			path.Join(netState.Config.LogsDir(), jobID),
			logsOffset,
			followLogs,
			false,
		)
	},
}

func init() {
	logsCmd.Flags().BoolVar(&followLogs,
		"follow",
		false,
		"Allows to configure whether or not the logs should be followed",
	)
	logsCmd.Flags().Int64Var(&logsOffset,
		"offset",
		0,
		"The offset to start streaming data at",
	)
	logsCmd.Flags().StringVar(&jobID,
		"job-id",
		"",
		"ID of the job we want to collect logs from. Leaving empty means all",
	)
}
