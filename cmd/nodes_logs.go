package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var (
	followLogs bool
	logsOrigin string
	logsOffset int64
	jobID      string
)

var nodesLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print nodes sets logs. By default prints logs accross all node sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("state unsafe-reset-all")
		}

		logs, err := nodesLogs(context.Background(), *netState, followLogs, logsOrigin, logsOffset, jobID)
		if err != nil {
			return err
		}

		if _, err := io.Copy(os.Stdout, logs); err != nil {
			return fmt.Errorf("failed to write command output: %w", err)
		}

		return nil
	},
}

func init() {
	nodesLogsCmd.PersistentFlags().BoolVar(&followLogs,
		"follow",
		false,
		"Allows to configure whether or not the logs should be followed",
	)
	nodesLogsCmd.PersistentFlags().StringVar(&logsOrigin,
		"origin",
		"start",
		"Origin can be either 'start' or 'end' and it defines from where the offset starts",
	)
	nodesLogsCmd.PersistentFlags().Int64Var(&logsOffset,
		"offset",
		0,
		"The offset to start streaming data at",
	)
	nodesLogsCmd.PersistentFlags().StringVar(&jobID,
		"id",
		"",
		"ID of the node set we want to collect logs from. Leaving empty means all",
	)
}

func nodesLogs(
	ctx context.Context,
	state state.NetworkState,
	follow bool,
	origin string,
	offset int64,
	jobID string,
) (io.ReadCloser, error) {
	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	if jobID != "" && !state.RunningJobs.NodesSetsJobIDs[jobID] {
		return nil, fmt.Errorf("job %q not found", jobID)
	}

	jobIDs := []string{jobID}
	if jobIDs[0] == "" {
		jobIDs = state.RunningJobs.NodesSetsJobIDs.ToSlice()
	}

	logs, err := nomadClient.LogJobs(ctx, follow, origin, offset, jobIDs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
