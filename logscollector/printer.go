package logscollector

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/nxadm/tail"
	"golang.org/x/sync/errgroup"
)

type nomadLogFile struct {
	path      string
	name      string
	taskName  string
	createdAt time.Time
}

var logsFileRegex = regexp.MustCompile("(.*[stderr|stdout])-(.*).log")

func TailLastLogs(logsDir string) error {
	return Tail(logsDir, 4000, false, false)
}

func Tail(logsDir string, offset int64, follow, withLogger bool) error {
	fileExists, err := utils.FileExists(logsDir)
	if err != nil {
		return err
	}

	if !fileExists {
		return fmt.Errorf("logs directory %q does not exists", logsDir)
	}

	logFilePerTaskName, err := getLogsFilesPerTaskName(logsDir, withLogger)
	if err != nil {
		return fmt.Errorf("failed to get logs per task name: %w", err)
	}

	var eg errgroup.Group
	for _, key := range logFilePerTaskName.SortedKeys() {
		logFile := logFilePerTaskName[key]

		if follow {
			eg.Go(func() error {
				return printLogFile(logFile, offset, follow)
			})
		} else {
			if err := printLogFile(logFile, offset, follow); err != nil {
				return err
			}
		}
	}

	return eg.Wait()
}

func printLogFile(logFile nomadLogFile, offset int64, follow bool) error {
	fileInfo, err := os.Stat(logFile.path)
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		return nil
	}

	var seekInfo *tail.SeekInfo
	if offset > fileInfo.Size() {
		offset = fileInfo.Size()
	}

	if offset > 0 {
		seekInfo = &tail.SeekInfo{
			Offset: -offset,
			Whence: io.SeekEnd,
		}
	}

	t, err := tail.TailFile(logFile.path, tail.Config{
		Follow:   follow,
		Poll:     true,
		Location: seekInfo,
	})
	if err != nil {
		return fmt.Errorf("failed to tail file %q: %q", logFile.path, err)
	}

	if !follow {
		fmt.Printf("-------- %s:\n", logFile.name)
	}

	for l := range t.Lines {
		if l == nil {
			continue
		}
		if l.Err != nil {
			t.Cleanup()
			return err
		}

		text := l.Text
		if follow {
			text = fmt.Sprintf("%s:    %s", logFile.taskName, text)
		}

		if _, err := fmt.Fprintln(os.Stdout, text); err != nil {
			return fmt.Errorf("failed to write to log file %q: %q", logFile.path, err)
		}
	}
	fmt.Println()
	t.Cleanup()

	return nil
}

type logsFiles map[string]nomadLogFile

func (lf logsFiles) SortedKeys() []string {
	keys := make([]string, 0, len(lf))
	for k := range lf {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getLogsFilesPerTaskName(logsDir string, withLogger bool) (logsFiles, error) {
	match := fmt.Sprintf(`%s/*.[stderr|stdout]*.log`, logsDir)

	logsPaths, err := filepath.Glob(match)
	if err != nil {
		return nil, fmt.Errorf("failed to look for files: %w", err)
	}

	logFilePerTaskName := logsFiles{}

	for _, logPath := range logsPaths {
		logFile := filepath.Base(logPath)

		subMatch := logsFileRegex.FindStringSubmatch(logFile)
		if len(subMatch) != 3 {
			continue
		}

		taskName, createdAtStr := subMatch[1], subMatch[2]

		createdAt, err := time.Parse(timeFormat, createdAtStr)
		if err != nil {
			log.Printf("failed to parse time from log file %q: %s", logFile, err)
			continue
		}

		existingLogFile, ok := logFilePerTaskName[taskName]
		if ok && existingLogFile.createdAt.After(createdAt) {
			continue
		}

		if !withLogger && strings.Contains(taskName, types.NomadLogsCollectorTaskName) {
			continue
		}

		logFilePerTaskName[taskName] = nomadLogFile{
			taskName:  taskName,
			name:      logFile,
			path:      logPath,
			createdAt: createdAt,
		}
	}

	return logFilePerTaskName, nil
}
