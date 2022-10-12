package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"code.vegaprotocol.io/vegacapsule/utils"
)

var nomadConfigTemplate = `
plugin "docker" {
	config {
		auth {
			%s
		}

		extra_labels = ["job_name", "job_id"]

		volumes {
			enabled = true
		}
	}
}
client {
	cpu_total_compute = 20000
	memory_total_mb = 16000
}
`

func templateConfig() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf(nomadConfigTemplate, `helper = "osxkeychain"`), nil
	case "windows":
		return fmt.Sprintf(nomadConfigTemplate, `helper = "desktop.exe"`), nil
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		authConfig := fmt.Sprintf(`config = "%s"`, filepath.Join(homeDir, ".docker", "config.json"))
		return fmt.Sprintf(nomadConfigTemplate, authConfig), nil
	default:
		return "", fmt.Errorf("platform not supported")
	}
}

func generateConfig() (string, error) {
	conf, err := templateConfig()
	if err != nil {
		return "", fmt.Errorf("failed to template config: %w", err)
	}

	homeDir, err := capsuleHome()
	if err != nil {
		return "", fmt.Errorf("failed to get vegacapsule home: %w", err)
	}

	confPath := filepath.Join(homeDir, "nomad-config.hcl")
	configFile, err := utils.CreateFile(confPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q: %w", confPath, err)
	}
	defer configFile.Close()

	if _, err := configFile.WriteString(conf); err != nil {
		return "", fmt.Errorf("failed to write to file %q: %w", confPath, err)
	}

	return confPath, nil
}
