package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func ExecuteBinary(binaryPath string, args []string, v interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute binary %s %v with error: %s: %s", binaryPath, args, stErr.String(), err.Error())
	}

	if v == nil {
		return stdOut.Bytes(), nil
	}

	if err := json.Unmarshal(stdOut.Bytes(), v); err != nil {
		// TODO Maybe failback to text parsing instead??
		return nil, err
	}

	return nil, nil
}

func BinaryAbsPath(p string) (string, error) {
	lPath, err := exec.LookPath(p)
	if err != nil {
		return "", fmt.Errorf("failed to look up path for %q: %w", p, err)
	}

	return AbsPath(lPath)
}

// Returns and error if given binaries are not accesible trough PATH
func BinariesAccessible(binaries ...string) error {
	var eg errgroup.Group

	for _, bin := range binaries {
		bin := bin
		eg.Go(func() error {
			if _, err := BinaryAbsPath(bin); err != nil {
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}

func AbsPath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}

	aPath, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %q: %w", p, err)
	}

	return aPath, nil
}

func AbsPathWithPrefix(prefix, p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}

	aPath, err := filepath.Abs(path.Join(prefix, p))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %q: %w", p, err)
	}

	return aPath, nil
}

func VegaNodeHomePath(networkHomePath string, nodeIdx int) string {
	return filepath.Join(networkHomePath, "vega", fmt.Sprintf("node%d", nodeIdx))
}

func ToPoint[T any](val T) *T {
	return &val
}
