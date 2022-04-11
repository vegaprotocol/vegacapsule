package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
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
	if filepath.IsAbs(lPath) {
		return lPath, nil
	}

	aPath, err := filepath.Abs(lPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %q: %w", lPath, err)
	}

	return aPath, nil
}

func VegaNodeHomePath(networkHomePath string, nodeIdx int) string {
	return filepath.Join(networkHomePath, "vega", fmt.Sprintf("node%d", nodeIdx))
}

func StrPoint(s string) *string {
	return &s
}

func IntPoint(i int) *int {
	return &i
}
