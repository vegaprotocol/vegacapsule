package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/nxadm/tail"
	"github.com/spf13/cobra"
)

type nomadLogFile struct {
	path      string
	name      string
	taskName  string
	createdAt time.Time
}

var nomadLogsCollectorTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tails logs of a specific Nomad job.",
	RunE: func(cmd *cobra.Command, args []string) error {
		logsDir := path.Join(nomadLogColOutDir, jobID)
		fileExists, err := utils.FileExists(logsDir)
		if err != nil {
			return err
		}

		if !fileExists {
			return fmt.Errorf("jobs %q logs directory %q does not exists", jobID, logsDir)
		}

		match := fmt.Sprintf(`%s/*.[stderr|stdout]*.log`, logsDir)

		logsPaths, err := filepath.Glob(match)
		if err != nil {
			return fmt.Errorf("failed to look for files: %w", err)
		}

		re := regexp.MustCompile("(.*[stderr|stdout])-(.*).log")

		logFilePerTaskName := map[string]nomadLogFile{}

		for _, logPath := range logsPaths {
			logFile := filepath.Base(logPath)

			subMatch := re.FindStringSubmatch(logFile)
			if len(subMatch) != 3 {
				continue
			}

			taskName, createdAtStr := subMatch[1], subMatch[2]

			createdAt, err := time.Parse(time.RFC3339, createdAtStr)
			if err != nil {
				log.Printf("failed to parse time from log file %q: %s", logFile, err)
				continue
			}

			existingLogFile, ok := logFilePerTaskName[taskName]
			if ok && existingLogFile.createdAt.After(createdAt) {
				continue
			}

			logFilePerTaskName[taskName] = nomadLogFile{
				taskName:  taskName,
				name:      logFile,
				path:      logPath,
				createdAt: createdAt,
			}
		}

		keys := make([]string, 0, len(logFilePerTaskName))
		for k, _ := range logFilePerTaskName {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			logFile := logFilePerTaskName[key]

			fileInfo, err := os.Stat(logFile.path)
			if err != nil {
				return err
			}

			if fileInfo.Size() == 0 {
				continue
			}

			var offset int64 = 4000
			if offset > fileInfo.Size() {
				offset = fileInfo.Size()
			}

			t, err := tail.TailFile(logFile.path, tail.Config{
				Follow: false,
				Poll:   true,
				Location: &tail.SeekInfo{
					Offset: -offset,
					Whence: io.SeekEnd,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to tail file %q: %q", logFile.path, err)
			}

			fmt.Printf("-------- %s:\n", logFile.name)

			for l := range t.Lines {
				if l == nil {
					continue
				}
				if l.Err != nil {
					t.Cleanup()
					return err
				}

				if _, err := fmt.Fprintln(os.Stdout, l.Text); err != nil {
					return fmt.Errorf("failed to write to log file %q: %q", logFile.path, err)
				}
			}
			fmt.Println()
			t.Cleanup()
		}

		return nil
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
