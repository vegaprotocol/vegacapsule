package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func executeBinary(binaryPath string, args []string, v interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("%s: %s", stErr.String(), err.Error())
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

func createSmartContractsDir(OutputDir string, smartContractsDir string) error {
	scd := path.Join(OutputDir, smartContractsDir)

	if err := os.MkdirAll(scd, os.ModePerm); err != nil {
		return err
	}
	return nil
}
