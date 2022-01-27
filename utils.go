package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func executeBinary(binaryPath string, args []string, v interface{}) error {
	command := exec.Command(binaryPath, args...)

	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return fmt.Errorf("%s: %s", stErr.String(), err.Error())
	}

	fmt.Println("executeBinary: ", stdOut.String())

	if err := json.Unmarshal(stdOut.Bytes(), v); err != nil {
		// TODO Maybe failback to text parsing instead??
		return err
	}

	return nil
}
