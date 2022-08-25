package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"code.vegaprotocol.io/vegacapsule/logscollector"
	"github.com/spf13/cobra"
)

var nomadLogColOutDir string

var nomadLogsCollectorCmd = &cobra.Command{
	Use:   "logscollector",
	Short: "Starts a log collection program that should be run as a logging sidecar inside Nomad job.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		if err := os.MkdirAll(nomadLogColOutDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to make output directory %q: %w", nomadLogColOutDir, err)
		}

		collector := logscollector.New(path.Join(path.Dir(cwd), "alloc", "logs"), nomadLogColOutDir)

		ctx, cancel := context.WithCancel(cmd.Context())

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.Printf("Received signal: %s", sig)
			cancel()
		}()

		return collector.Run(ctx)
	},
}

func init() {
	nomadLogsCollectorCmd.PersistentFlags().StringVar(&nomadLogColOutDir,
		"out-dir",
		"",
		"Output directory for logs.",
	)

	nomadLogsCollectorCmd.MarkPersistentFlagRequired("out-dir") // nolint:errcheck
}
