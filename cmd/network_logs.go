package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/spf13/cobra"
)

var (
	followLogs       bool
	logsOrigin       string
	logsOffset       int64
	logsOnlyNodeSets bool
	jobID            string
)

var netLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs from running jobs in network. By default prints logs across all jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("net logs")
		}

		if !netState.Running() {
			return networkNotRunningErr("net logs")
		}

		jobIDs, err := filterJobIDsForLogs(*netState.RunningJobs, logsOnlyNodeSets, jobID)
		if err != nil {
			return err
		}

		logs, err := netLogs(context.Background(), *netState, followLogs, logsOrigin, logsOffset, jobIDs)
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
	netLogsCmd.PersistentFlags().BoolVar(&followLogs,
		"follow",
		false,
		"Allows to configure whether or not the logs should be followed",
	)
	netLogsCmd.PersistentFlags().StringVar(&logsOrigin,
		"origin",
		"start",
		"Origin can be either 'start' or 'end' and it defines from where the offset starts",
	)
	netLogsCmd.PersistentFlags().Int64Var(&logsOffset,
		"offset",
		0,
		"The offset to start streaming data at",
	)
	netLogsCmd.PersistentFlags().StringVar(&jobID,
		"job-id",
		"",
		"ID of the job we want to collect logs from. Leaving empty means all",
	)
	netLogsCmd.PersistentFlags().BoolVar(&logsOnlyNodeSets,
		"nodes-only",
		false,
		"Marks that only logs from all nodes sets should be aggregated. Caution: job-id flag will override this flag",
	)
}

func filterJobIDsForLogs(jobs types.NetworkJobs, nodeSetsOnly bool, jobID string) ([]string, error) {
	if jobID != "" {
		if !jobs.Exists(jobID) {
			return nil, fmt.Errorf("job %q not found", jobID)
		}

		return []string{jobID}, nil
	}

	if nodeSetsOnly {
		return jobs.NodesSetsJobIDs.ToSlice(), nil
	}

	return jobs.ToSlice(), nil
}

func netLogs(
	ctx context.Context,
	state state.NetworkState,
	follow bool,
	origin string,
	offset int64,
	jobIDs []string,
) (io.ReadCloser, error) {
	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	logs, err := nomadClient.LogJobs(ctx, follow, origin, offset, jobIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nomad jobs: %w", err)
	}

	return logs, nil
}
