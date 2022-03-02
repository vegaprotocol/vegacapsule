package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func capsuleHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".vegacapsule"), nil
}

func StartAgent(configPath string) error {
	switch runtime.GOOS {
	case "darwin", "windows", "linux":
	default:
		return fmt.Errorf("unsupported platform %q. Supported platform are: darwin, windows, linux", runtime.GOOS)
	}

	nomadBinary, err := nomadBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get nomad binary: %w", err)
	}

	if _, err := exec.LookPath(nomadBinary); err != nil {
		if err := installNomadBinary(nomadBinary); err != nil {
			return fmt.Errorf("failed to install nomad binary: %w", err)
		}
	}

	if configPath == "" {
		generatedConfigPath, err := generateConfig()
		if err != nil {
			return fmt.Errorf("failed to generate nomad config: %w", err)
		}

		configPath = generatedConfigPath
	}

	args := []string{"agent", "-dev", "-bind", "0.0.0.0", "-config", configPath}

	log.Printf("Starting nomad with with: %s %v", nomadBinary, args)

	command := exec.Command(nomadBinary, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to execute binary %s with error: %s", nomadBinary, err.Error())
	}

	return nil
}
