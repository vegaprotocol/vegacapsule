package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/nxadm/tail"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var nomadLogColOutDir string

var nomadLogCmd = &cobra.Command{
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

		logsPath := path.Join(path.Dir(cwd), "alloc", "logs")
		match := fmt.Sprintf("%s/*.[stderr|stdout]*[0-9]", logsPath)

		log.Printf("Looking for log files with %q", match)

		logsPaths, err := filepath.Glob(match)
		if err != nil {
			return fmt.Errorf("failed to look for files: %w", err)
		}

		eg, ctx := errgroup.WithContext(cmd.Context())

		for _, logFilePath := range logsPaths {
			logFilePath := logFilePath

			eg.Go(func() error {
				log.Printf("Setting up file listener for %q", logFilePath)

				t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
				if err != nil {
					return fmt.Errorf("failed to setup file listener %q: %q", logFilePath, err)
				}
				defer t.Cleanup()

				logFileName := path.Base(logFilePath)
				destLogFile := path.Join(nomadLogColOutDir, logFileName)

				f, err := os.Create(destLogFile)
				if err != nil {
					return fmt.Errorf("failed to create log file %q: %q", destLogFile, err)
				}
				defer f.Close()

				for {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case l := <-t.Lines:
						if _, err := fmt.Fprintln(f, l.Text); err != nil {
							return fmt.Errorf("failed to write to log file %q: %q", f.Name(), err)
						}
					}
				}
			})
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigs

		log.Printf("Recived signal: %s", sig)

		return nil
	},
}

func init() {
	nomadLogCmd.PersistentFlags().StringVar(&nomadLogColOutDir,
		"out-dir",
		"",
		"Output directory for logs.",
	)

	nomadLogCmd.MarkPersistentFlagRequired("out-dir")
}
