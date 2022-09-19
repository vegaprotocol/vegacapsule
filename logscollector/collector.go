package logscollector

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nxadm/tail"
	"golang.org/x/sync/errgroup"
)

type Collector struct {
	logsDir   string
	outputDir string

	filesToCollect chan string
}

func New(logsDir, outputDir string) *Collector {
	return &Collector{
		filesToCollect: make(chan string, 20),
		logsDir:        logsDir,
		outputDir:      outputDir,
	}
}

func (lc Collector) watchForCreatedFiles(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(lc.logsDir); err != nil {
		return err
	}
	defer close(lc.filesToCollect)

	log.Printf("starting files watcher in %q", lc.logsDir)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				lc.filesToCollect <- event.Name
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return err
		}
	}
}

func (lc Collector) collectLogs(ctx context.Context, logFilePath string) error {
	log.Printf("Setting up file listener for %q", logFilePath)
	defer log.Printf("Shutting down file listener for %q", logFilePath)

	t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
	if err != nil {
		return fmt.Errorf("failed to setup file listener %q: %q", logFilePath, err)
	}
	defer t.Cleanup()

	logFileNameBase := path.Base(logFilePath)
	logFileName := strings.TrimSuffix(logFileNameBase, filepath.Ext(logFileNameBase))
	destLogFile := path.Join(lc.outputDir, fmt.Sprintf("%s-%s.log", logFileName, time.Now().Format(time.RFC3339)))

	f, err := os.Create(destLogFile)
	if err != nil {
		return fmt.Errorf("failed to create log file %q: %q", destLogFile, err)
	}
	defer f.Close()

	stopChan := make(chan struct{}, 1)
	defer close(stopChan)

	go func() {
		<-ctx.Done()
		if err := t.StopAtEOF(); err != nil {
			t.Stop() //nolint
		}

		time.Sleep(time.Second * 5)
		stopChan <- struct{}{}
	}()

	for {
		select {
		case l := <-t.Lines:
			if l == nil {
				return nil
			}
			if l.Err != nil {
				return err
			}

			if _, err := fmt.Fprintln(f, l.Text); err != nil {
				return fmt.Errorf("failed to write to log file %q: %q", f.Name(), err)
			}
		case <-stopChan:
			return nil
		}

	}
}

func (lc Collector) Run(ctx context.Context) error {
	match := fmt.Sprintf("%s/*.[stderr|stdout]*[0-9]", lc.logsDir)

	log.Printf("Looking for log files with %q", match)

	logsPaths, err := filepath.Glob(match)
	if err != nil {
		return fmt.Errorf("failed to look for files: %w", err)
	}

	for _, logPath := range logsPaths {
		lc.filesToCollect <- logPath
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := lc.watchForCreatedFiles(ctx); err != nil {
			return fmt.Errorf("failed to watch for created files: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		for logFilePath := range lc.filesToCollect {
			logFilePath := logFilePath

			eg.Go(func() error {
				if err := lc.collectLogs(ctx, logFilePath); err != nil {
					return fmt.Errorf("failec to collect logs for %q: %w", logFilePath, err)
				}

				return nil
			})
		}

		return nil
	})

	return eg.Wait()
}
