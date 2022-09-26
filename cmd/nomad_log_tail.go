package cmd

import (
	"path"

	"code.vegaprotocol.io/vegacapsule/logscollector"
	"github.com/spf13/cobra"
)

var nomadLogsCollectorTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tails logs of a specific Nomad job.",
	RunE: func(cmd *cobra.Command, args []string) error {
		logsDir := path.Join(nomadLogColOutDir, jobID)

		return logscollector.TailLastLogs(logsDir)
	},
}

func init() {
	nomadLogsCollectorTailCmd.PersistentFlags().StringVar(&nomadLogColOutDir,
		"out-dir",
		"",
		"Output directory for logs.",
	)
	nomadLogsCollectorTailCmd.PersistentFlags().StringVar(&jobID,
		"job-id",
		"",
		"ID of the job",
	)

	nomadLogsCollectorTailCmd.MarkPersistentFlagRequired("out-dir") // nolint:errcheck
	nomadLogsCollectorTailCmd.MarkPersistentFlagRequired("job-id")  // nolint:errcheck
}
