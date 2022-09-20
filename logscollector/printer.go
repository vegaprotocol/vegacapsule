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

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/nxadm/tail"
)

type nomadLogFile struct {
	path      string
	name      string
	taskName  string
	createdAt time.Time
}

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

		if !withLogger && strings.Contains(taskName, "logger") {
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
	for k := range logFilePerTaskName {
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
}
